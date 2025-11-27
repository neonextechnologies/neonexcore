# Repository Pattern

Learn how Neonex Core implements the Repository Pattern for clean data access abstraction.

---

## Overview

The **Repository Pattern** provides an abstraction layer between business logic and data access. Neonex Core uses **generic repositories** powered by Go generics for type-safe, reusable data operations.

**Benefits:**
- ✅ **Separation of Concerns** - Business logic independent of data access
- ✅ **Testability** - Easy to mock repositories
- ✅ **Consistency** - Uniform data access across modules
- ✅ **Type Safety** - Generic-based type checking
- ✅ **Reusability** - DRY principle with base repository

---

## Base Repository

### Generic Repository Interface

```go
// pkg/database/repository.go
package database

type Repository[T any] interface {
    Create(ctx context.Context, entity *T) error
    CreateBatch(ctx context.Context, entities []*T) error
    Update(ctx context.Context, entity *T) error
    Delete(ctx context.Context, id interface{}) error
    FindByID(ctx context.Context, id interface{}) (*T, error)
    FindAll(ctx context.Context) ([]*T, error)
    FindByCondition(ctx context.Context, condition interface{}, args ...interface{}) ([]*T, error)
    FindOne(ctx context.Context, condition interface{}, args ...interface{}) (*T, error)
    Count(ctx context.Context, condition interface{}, args ...interface{}) (int64, error)
    Paginate(ctx context.Context, page, pageSize int) ([]*T, int64, error)
}
```

### Implementation

```go
type BaseRepository[T any] struct {
    db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
    return &BaseRepository[T]{db: db}
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
    return r.db.WithContext(ctx).Create(entity).Error
}

func (r *BaseRepository[T]) FindByID(ctx context.Context, id interface{}) (*T, error) {
    var entity T
    err := r.db.WithContext(ctx).First(&entity, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &entity, nil
}

func (r *BaseRepository[T]) FindAll(ctx context.Context) ([]*T, error) {
    var entities []*T
    err := r.db.WithContext(ctx).Find(&entities).Error
    return entities, err
}
```

---

## Module Repository

### Creating a Repository

```go
// modules/user/repository.go
package user

import (
    "context"
    "neonexcore/pkg/database"
    "gorm.io/gorm"
)

// Define interface with custom methods
type Repository interface {
    database.Repository[User]  // Embed base interface
    
    // Custom methods
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindActive(ctx context.Context) ([]*User, error)
    UpdateLastLogin(ctx context.Context, id uint) error
}

// Implementation
type repository struct {
    *database.BaseRepository[User]
}

func NewRepository(db *gorm.DB) Repository {
    return &repository{
        BaseRepository: database.NewBaseRepository[User](db),
    }
}

// Implement custom methods
func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
    return r.FindOne(ctx, "email = ?", email)
}

func (r *repository) FindActive(ctx context.Context) ([]*User, error) {
    return r.FindByCondition(ctx, "active = ?", true)
}

func (r *repository) UpdateLastLogin(ctx context.Context, id uint) error {
    return r.GetDB().WithContext(ctx).
        Model(&User{}).
        Where("id = ?", id).
        Update("last_login_at", time.Now()).Error
}
```

---

## CRUD Operations

### Create

```go
// Single entity
user := &User{
    Name:  "John Doe",
    Email: "john@example.com",
}
err := repo.Create(ctx, user)

// Batch insert
users := []*User{
    {Name: "User 1", Email: "user1@example.com"},
    {Name: "User 2", Email: "user2@example.com"},
}
err := repo.CreateBatch(ctx, users)
```

### Read

```go
// Find by ID
user, err := repo.FindByID(ctx, 1)

// Find all
users, err := repo.FindAll(ctx)

// Find by condition
activeUsers, err := repo.FindByCondition(ctx, "active = ?", true)

// Find one
user, err := repo.FindOne(ctx, "email = ?", "john@example.com")

// Count
count, err := repo.Count(ctx, "role = ?", "admin")

// Paginate
users, total, err := repo.Paginate(ctx, 1, 20) // page 1, 20 per page
```

### Update

```go
user.Name = "Jane Doe"
err := repo.Update(ctx, user)
```

### Delete

```go
// Soft delete (if model has DeletedAt)
err := repo.Delete(ctx, 1)

// Hard delete
err := repo.GetDB().Unscoped().Delete(&User{}, 1).Error
```

---

## Advanced Queries

### Complex Conditions

```go
func (r *repository) FindByRole(ctx context.Context, role string) ([]*User, error) {
    var users []*User
    err := r.GetDB().WithContext(ctx).
        Where("role = ? AND active = ?", role, true).
        Order("created_at DESC").
        Find(&users).Error
    return users, err
}
```

### Joins

```go
func (r *repository) FindWithOrders(ctx context.Context, userID uint) (*User, error) {
    var user User
    err := r.GetDB().WithContext(ctx).
        Preload("Orders").
        First(&user, userID).Error
    return &user, err
}
```

### Aggregations

```go
func (r *repository) GetStats(ctx context.Context) (*Stats, error) {
    var stats Stats
    err := r.GetDB().WithContext(ctx).
        Model(&User{}).
        Select("COUNT(*) as total, COUNT(CASE WHEN active = true THEN 1 END) as active").
        Scan(&stats).Error
    return &stats, err
}
```

### Raw SQL

```go
func (r *repository) ExecuteRaw(ctx context.Context, query string, args ...interface{}) error {
    return r.GetDB().WithContext(ctx).Exec(query, args...).Error
}

func (r *repository) QueryRaw(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
    return r.GetDB().WithContext(ctx).Raw(query, args...).Scan(dest).Error
}
```

---

## Transactions

### Using BaseRepository with Transaction

```go
func (s *service) TransferUser(ctx context.Context, fromID, toID uint) error {
    return database.WithTransaction(ctx, s.db, func(tx *gorm.DB) error {
        // Create repo with transaction
        txRepo := s.repo.(*repository).WithTx(tx)
        
        from, err := txRepo.FindByID(ctx, fromID)
        if err != nil {
            return err
        }
        
        from.Balance -= 100
        if err := txRepo.Update(ctx, from); err != nil {
            return err
        }
        
        to, err := txRepo.FindByID(ctx, toID)
        if err != nil {
            return err
        }
        
        to.Balance += 100
        return txRepo.Update(ctx, to)
    })
}
```

### WithTx Helper

```go
func (r *BaseRepository[T]) WithTx(tx *gorm.DB) *BaseRepository[T] {
    return &BaseRepository[T]{db: tx}
}
```

---

## Best Practices

### ✅ DO:

**1. Use Interfaces**
```go
// Good: Define interface
type UserRepository interface {
    FindByEmail(email string) (*User, error)
}

// Good: Program to interface
type UserService struct {
    repo UserRepository
}
```

**2. Add Context**
```go
// Good: Pass context for cancellation
func (r *repository) FindAll(ctx context.Context) ([]*User, error) {
    return r.BaseRepository.FindAll(ctx)
}
```

**3. Handle Errors**
```go
// Good: Check for specific errors
user, err := repo.FindByID(ctx, id)
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, ErrUserNotFound
    }
    return nil, err
}
```

**4. Use Generics**
```go
// Good: Type-safe repository
type Repository[T any] interface {
    Create(ctx context.Context, entity *T) error
}
```

### ❌ DON'T:

**1. Return GORM Errors Directly**
```go
// Bad: Expose implementation details
return r.db.First(&user, id).Error

// Good: Wrap errors
if err := r.db.First(&user, id).Error; err != nil {
    return fmt.Errorf("failed to find user: %w", err)
}
```

**2. Put Business Logic in Repository**
```go
// Bad: Business logic in repository
func (r *repository) CreateUserWithDiscount(user *User) error {
    user.Discount = 0.1  // Business logic!
    return r.Create(user)
}

// Good: Keep it in service layer
```

**3. Use Global DB**
```go
// Bad: Global variable
var DB *gorm.DB

// Good: Inject dependency
func NewRepository(db *gorm.DB) Repository {
    return &repository{db: db}
}
```

---

## Testing

### Mock Repository

```go
// test/mock_repository.go
type MockUserRepository struct {
    users map[uint]*User
    err   error
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*User, error) {
    if m.err != nil {
        return nil, m.err
    }
    return m.users[id], nil
}

func (m *MockUserRepository) Create(ctx context.Context, user *User) error {
    if m.err != nil {
        return m.err
    }
    user.ID = uint(len(m.users) + 1)
    m.users[user.ID] = user
    return nil
}
```

### Test with Mock

```go
func TestUserService_CreateUser(t *testing.T) {
    // Setup mock
    mockRepo := &MockUserRepository{
        users: make(map[uint]*User),
    }
    
    service := NewUserService(mockRepo)
    
    // Test
    user := &User{Name: "Test", Email: "test@example.com"}
    err := service.CreateUser(context.Background(), user)
    
    assert.NoError(t, err)
    assert.Equal(t, uint(1), user.ID)
}
```

---

## Complete Example

### Full Repository Implementation

```go
// modules/product/repository.go
package product

import (
    "context"
    "fmt"
    "neonexcore/pkg/database"
    "gorm.io/gorm"
)

type Repository interface {
    database.Repository[Product]
    
    FindByCategory(ctx context.Context, category string) ([]*Product, error)
    FindInPriceRange(ctx context.Context, min, max float64) ([]*Product, error)
    FindPopular(ctx context.Context, limit int) ([]*Product, error)
    UpdateStock(ctx context.Context, id uint, quantity int) error
    Search(ctx context.Context, query string) ([]*Product, error)
}

type repository struct {
    *database.BaseRepository[Product]
}

func NewRepository(db *gorm.DB) Repository {
    return &repository{
        BaseRepository: database.NewBaseRepository[Product](db),
    }
}

func (r *repository) FindByCategory(ctx context.Context, category string) ([]*Product, error) {
    return r.FindByCondition(ctx, "category = ?", category)
}

func (r *repository) FindInPriceRange(ctx context.Context, min, max float64) ([]*Product, error) {
    var products []*Product
    err := r.GetDB().WithContext(ctx).
        Where("price BETWEEN ? AND ?", min, max).
        Find(&products).Error
    return products, err
}

func (r *repository) FindPopular(ctx context.Context, limit int) ([]*Product, error) {
    var products []*Product
    err := r.GetDB().WithContext(ctx).
        Order("views DESC").
        Limit(limit).
        Find(&products).Error
    return products, err
}

func (r *repository) UpdateStock(ctx context.Context, id uint, quantity int) error {
    return r.GetDB().WithContext(ctx).
        Model(&Product{}).
        Where("id = ?", id).
        Update("stock", gorm.Expr("stock + ?", quantity)).Error
}

func (r *repository) Search(ctx context.Context, query string) ([]*Product, error) {
    var products []*Product
    searchQuery := fmt.Sprintf("%%%s%%", query)
    err := r.GetDB().WithContext(ctx).
        Where("name LIKE ? OR description LIKE ?", searchQuery, searchQuery).
        Find(&products).Error
    return products, err
}
```

---

## Next Steps

- [**Service Layer**](service-layer.md) - Business logic on top of repositories
- [**Transactions**](../database/transactions.md) - Managing database transactions
- [**Module System**](module-system.md) - Integrate repositories with modules
- [**Testing**](../development/testing.md) - Test repositories and mocks

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
