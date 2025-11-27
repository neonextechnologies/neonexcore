# Advanced Repository Patterns

Explore advanced repository patterns and custom query techniques in Neonex Core.

---

## Overview

**Repositories** provide an abstraction layer between business logic and data access. This guide covers advanced patterns beyond basic CRUD.

**Benefits:**
- **Encapsulation** - Hide query complexity
- **Reusability** - Share common queries
- **Testability** - Mock data layer easily
- **Maintainability** - Centralize data access

---

## Custom Repository Methods

### Basic Extension

```go
// modules/user/repository.go
package user

import (
    "context"
    "neonexcore/pkg/database"
    "gorm.io/gorm"
)

type Repository interface {
    database.Repository[User]
    
    // Custom methods
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindActive(ctx context.Context) ([]*User, error)
    Search(ctx context.Context, query string) ([]*User, error)
}

type repository struct {
    database.BaseRepository[User]
}

func NewRepository(db *gorm.DB) Repository {
    return &repository{
        BaseRepository: database.NewBaseRepository[User](db),
    }
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
    var user User
    err := r.DB().WithContext(ctx).
        Where("email = ?", email).
        First(&user).Error
    
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *repository) FindActive(ctx context.Context) ([]*User, error) {
    var users []*User
    err := r.DB().WithContext(ctx).
        Where("active = ?", true).
        Order("created_at DESC").
        Find(&users).Error
    
    return users, err
}

func (r *repository) Search(ctx context.Context, query string) ([]*User, error) {
    var users []*User
    pattern := "%" + query + "%"
    
    err := r.DB().WithContext(ctx).
        Where("name LIKE ? OR email LIKE ?", pattern, pattern).
        Find(&users).Error
    
    return users, err
}
```

---

## Specification Pattern

### Query Specifications

```go
// pkg/database/specification.go
package database

import "gorm.io/gorm"

type Specification interface {
    Apply(db *gorm.DB) *gorm.DB
}

// AND specification
type AndSpecification struct {
    specs []Specification
}

func And(specs ...Specification) Specification {
    return &AndSpecification{specs: specs}
}

func (s *AndSpecification) Apply(db *gorm.DB) *gorm.DB {
    for _, spec := range s.specs {
        db = spec.Apply(db)
    }
    return db
}

// OR specification
type OrSpecification struct {
    specs []Specification
}

func Or(specs ...Specification) Specification {
    return &OrSpecification{specs: specs}
}

func (s *OrSpecification) Apply(db *gorm.DB) *gorm.DB {
    conditions := make([]interface{}, 0, len(s.specs))
    
    for _, spec := range s.specs {
        subDB := db.Session(&gorm.Session{NewDB: true})
        conditions = append(conditions, spec.Apply(subDB))
    }
    
    return db.Where(db.Or(conditions...))
}
```

### User Specifications

```go
// modules/user/specifications.go
package user

import (
    "gorm.io/gorm"
    "neonexcore/pkg/database"
)

// Active users
type ActiveSpec struct{}

func (s ActiveSpec) Apply(db *gorm.DB) *gorm.DB {
    return db.Where("active = ?", true)
}

// By role
type RoleSpec struct {
    Role string
}

func (s RoleSpec) Apply(db *gorm.DB) *gorm.DB {
    return db.Where("role = ?", s.Role)
}

// Created after date
type CreatedAfterSpec struct {
    Date time.Time
}

func (s CreatedAfterSpec) Apply(db *gorm.DB) *gorm.DB {
    return db.Where("created_at > ?", s.Date)
}

// Email verified
type EmailVerifiedSpec struct{}

func (s EmailVerifiedSpec) Apply(db *gorm.DB) *gorm.DB {
    return db.Where("email_verified_at IS NOT NULL")
}
```

### Using Specifications

```go
// Repository method
func (r *repository) FindBySpec(ctx context.Context, spec database.Specification) ([]*User, error) {
    var users []*User
    db := spec.Apply(r.DB().WithContext(ctx))
    err := db.Find(&users).Error
    return users, err
}

// Service usage
func (s *service) GetActiveAdmins(ctx context.Context) ([]*User, error) {
    spec := database.And(
        ActiveSpec{},
        RoleSpec{Role: "admin"},
    )
    return s.repo.FindBySpec(ctx, spec)
}

func (s *service) GetRecentVerifiedUsers(ctx context.Context) ([]*User, error) {
    spec := database.And(
        EmailVerifiedSpec{},
        CreatedAfterSpec{Date: time.Now().AddDate(0, -1, 0)},
    )
    return s.repo.FindBySpec(ctx, spec)
}
```

---

## Query Builder Pattern

### Fluent Query Builder

```go
// pkg/database/query_builder.go
package database

type QueryBuilder[T any] struct {
    db *gorm.DB
}

func NewQueryBuilder[T any](db *gorm.DB) *QueryBuilder[T] {
    return &QueryBuilder[T]{db: db}
}

func (qb *QueryBuilder[T]) Where(query string, args ...interface{}) *QueryBuilder[T] {
    qb.db = qb.db.Where(query, args...)
    return qb
}

func (qb *QueryBuilder[T]) Order(value string) *QueryBuilder[T] {
    qb.db = qb.db.Order(value)
    return qb
}

func (qb *QueryBuilder[T]) Limit(limit int) *QueryBuilder[T] {
    qb.db = qb.db.Limit(limit)
    return qb
}

func (qb *QueryBuilder[T]) Offset(offset int) *QueryBuilder[T] {
    qb.db = qb.db.Offset(offset)
    return qb
}

func (qb *QueryBuilder[T]) Preload(query string, args ...interface{}) *QueryBuilder[T] {
    qb.db = qb.db.Preload(query, args...)
    return qb
}

func (qb *QueryBuilder[T]) Find() ([]T, error) {
    var results []T
    err := qb.db.Find(&results).Error
    return results, err
}

func (qb *QueryBuilder[T]) First() (*T, error) {
    var result T
    err := qb.db.First(&result).Error
    return &result, err
}

func (qb *QueryBuilder[T]) Count() (int64, error) {
    var count int64
    err := qb.db.Count(&count).Error
    return count, err
}
```

### Usage Example

```go
// Repository method
func (r *repository) Query() *database.QueryBuilder[User] {
    return database.NewQueryBuilder[User](r.DB())
}

// Service usage
func (s *service) FindUsers(ctx context.Context, filters UserFilters) ([]*User, error) {
    query := s.repo.Query()
    
    if filters.Active != nil {
        query = query.Where("active = ?", *filters.Active)
    }
    
    if filters.Role != "" {
        query = query.Where("role = ?", filters.Role)
    }
    
    if filters.Search != "" {
        pattern := "%" + filters.Search + "%"
        query = query.Where("name LIKE ? OR email LIKE ?", pattern, pattern)
    }
    
    if filters.Limit > 0 {
        query = query.Limit(filters.Limit)
    }
    
    if filters.Offset > 0 {
        query = query.Offset(filters.Offset)
    }
    
    users, err := query.Order("created_at DESC").Find()
    if err != nil {
        return nil, err
    }
    
    // Convert []User to []*User
    result := make([]*User, len(users))
    for i := range users {
        result[i] = &users[i]
    }
    
    return result, nil
}
```

---

## Pagination

### Paginated Results

```go
// pkg/database/pagination.go
package database

type PaginationParams struct {
    Page     int
    PageSize int
    Sort     string
    Order    string
}

type PaginatedResult[T any] struct {
    Data       []T   `json:"data"`
    Total      int64 `json:"total"`
    Page       int   `json:"page"`
    PageSize   int   `json:"page_size"`
    TotalPages int   `json:"total_pages"`
}

func Paginate[T any](db *gorm.DB, params PaginationParams) (*PaginatedResult[T], error) {
    var data []T
    var total int64
    
    // Count total
    if err := db.Model(new(T)).Count(&total).Error; err != nil {
        return nil, err
    }
    
    // Calculate pagination
    if params.PageSize <= 0 {
        params.PageSize = 10
    }
    if params.Page <= 0 {
        params.Page = 1
    }
    
    offset := (params.Page - 1) * params.PageSize
    totalPages := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))
    
    // Apply sorting
    order := params.Sort
    if params.Order != "" {
        order += " " + params.Order
    }
    
    // Fetch data
    if err := db.Order(order).Offset(offset).Limit(params.PageSize).Find(&data).Error; err != nil {
        return nil, err
    }
    
    return &PaginatedResult[T]{
        Data:       data,
        Total:      total,
        Page:       params.Page,
        PageSize:   params.PageSize,
        TotalPages: totalPages,
    }, nil
}
```

### Repository Usage

```go
// Repository method
func (r *repository) Paginate(ctx context.Context, params database.PaginationParams) (*database.PaginatedResult[User], error) {
    db := r.DB().WithContext(ctx)
    return database.Paginate[User](db, params)
}

// Service usage
func (s *service) ListUsers(ctx context.Context, page, pageSize int) (*database.PaginatedResult[User], error) {
    params := database.PaginationParams{
        Page:     page,
        PageSize: pageSize,
        Sort:     "created_at",
        Order:    "DESC",
    }
    
    return s.repo.Paginate(ctx, params)
}
```

---

## Complex Queries

### Subqueries

```go
// Find users with more than 5 orders
func (r *repository) FindActiveCustomers(ctx context.Context) ([]*User, error) {
    var users []*User
    
    subQuery := r.DB().Model(&Order{}).
        Select("user_id").
        Group("user_id").
        Having("COUNT(*) > ?", 5)
    
    err := r.DB().WithContext(ctx).
        Where("id IN (?)", subQuery).
        Find(&users).Error
    
    return users, err
}
```

### Joins

```go
// Find users with their order count
func (r *repository) FindUsersWithOrderCount(ctx context.Context) ([]UserWithOrderCount, error) {
    type UserWithOrderCount struct {
        User
        OrderCount int
    }
    
    var results []UserWithOrderCount
    
    err := r.DB().WithContext(ctx).
        Model(&User{}).
        Select("users.*, COUNT(orders.id) as order_count").
        Joins("LEFT JOIN orders ON orders.user_id = users.id").
        Group("users.id").
        Scan(&results).Error
    
    return results, err
}
```

### Aggregations

```go
// Get user statistics
func (r *repository) GetStatistics(ctx context.Context) (*UserStats, error) {
    var stats UserStats
    
    err := r.DB().WithContext(ctx).
        Model(&User{}).
        Select(`
            COUNT(*) as total,
            COUNT(CASE WHEN active = true THEN 1 END) as active,
            COUNT(CASE WHEN email_verified_at IS NOT NULL THEN 1 END) as verified
        `).
        Scan(&stats).Error
    
    return &stats, err
}
```

---

## Raw SQL Queries

### Raw Query

```go
// Complex raw query
func (r *repository) FindTopSpenders(ctx context.Context, limit int) ([]*User, error) {
    var users []*User
    
    sql := `
        SELECT u.* 
        FROM users u
        INNER JOIN (
            SELECT user_id, SUM(amount) as total
            FROM orders
            WHERE status = 'completed'
            GROUP BY user_id
            ORDER BY total DESC
            LIMIT ?
        ) top_users ON u.id = top_users.user_id
    `
    
    err := r.DB().WithContext(ctx).Raw(sql, limit).Scan(&users).Error
    return users, err
}
```

### Named Parameters

```go
// Using named parameters (PostgreSQL)
func (r *repository) SearchByMultipleCriteria(ctx context.Context, criteria SearchCriteria) ([]*User, error) {
    var users []*User
    
    sql := `
        SELECT * FROM users
        WHERE 
            (@name::text IS NULL OR name ILIKE '%' || @name || '%')
            AND (@email::text IS NULL OR email ILIKE '%' || @email || '%')
            AND (@min_age::int IS NULL OR age >= @min_age)
            AND (@max_age::int IS NULL OR age <= @max_age)
    `
    
    err := r.DB().WithContext(ctx).
        Raw(sql, sql.Named("name", criteria.Name),
            sql.Named("email", criteria.Email),
            sql.Named("min_age", criteria.MinAge),
            sql.Named("max_age", criteria.MaxAge)).
        Scan(&users).Error
    
    return users, err
}
```

---

## Batch Operations

### Batch Insert

```go
// Insert many records efficiently
func (r *repository) CreateBatch(ctx context.Context, users []*User) error {
    return r.DB().WithContext(ctx).
        CreateInBatches(users, 1000).Error
}
```

### Batch Update

```go
// Update multiple records
func (r *repository) UpdateBatch(ctx context.Context, ids []uint, updates map[string]interface{}) error {
    return r.DB().WithContext(ctx).
        Model(&User{}).
        Where("id IN ?", ids).
        Updates(updates).Error
}
```

### Upsert (Insert or Update)

```go
// Upsert single record
func (r *repository) Upsert(ctx context.Context, user *User) error {
    return r.DB().WithContext(ctx).
        Clauses(clause.OnConflict{
            Columns:   []clause.Column{{Name: "email"}},
            DoUpdates: clause.AssignmentColumns([]string{"name", "updated_at"}),
        }).
        Create(user).Error
}

// Upsert multiple records
func (r *repository) UpsertBatch(ctx context.Context, users []*User) error {
    return r.DB().WithContext(ctx).
        Clauses(clause.OnConflict{
            Columns:   []clause.Column{{Name: "email"}},
            DoUpdates: clause.AssignmentColumns([]string{"name", "updated_at"}),
        }).
        Create(users).Error
}
```

---

## Soft Deletes

### Custom Soft Delete

```go
// Model with soft delete
type User struct {
    ID        uint           `gorm:"primarykey"`
    Name      string
    Email     string
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Include deleted records
func (r *repository) FindWithDeleted(ctx context.Context, id uint) (*User, error) {
    var user User
    err := r.DB().WithContext(ctx).
        Unscoped(). // Include soft-deleted records
        First(&user, id).Error
    return &user, err
}

// Permanently delete
func (r *repository) ForceDelete(ctx context.Context, id uint) error {
    return r.DB().WithContext(ctx).
        Unscoped().
        Delete(&User{}, id).Error
}

// Restore deleted record
func (r *repository) Restore(ctx context.Context, id uint) error {
    return r.DB().WithContext(ctx).
        Model(&User{}).
        Unscoped().
        Where("id = ?", id).
        Update("deleted_at", nil).Error
}
```

---

## Caching Layer

### Repository with Cache

```go
// Repository with Redis cache
type cachedRepository struct {
    database.BaseRepository[User]
    cache *redis.Client
}

func NewCachedRepository(db *gorm.DB, cache *redis.Client) Repository {
    return &cachedRepository{
        BaseRepository: database.NewBaseRepository[User](db),
        cache:          cache,
    }
}

func (r *cachedRepository) FindByID(ctx context.Context, id uint) (*User, error) {
    // Try cache first
    cacheKey := fmt.Sprintf("user:%d", id)
    
    cached, err := r.cache.Get(ctx, cacheKey).Result()
    if err == nil {
        var user User
        if err := json.Unmarshal([]byte(cached), &user); err == nil {
            return &user, nil
        }
    }
    
    // Cache miss, query database
    user, err := r.BaseRepository.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Cache result
    data, _ := json.Marshal(user)
    r.cache.Set(ctx, cacheKey, data, 15*time.Minute)
    
    return user, nil
}

func (r *cachedRepository) Update(ctx context.Context, user *User) error {
    // Update database
    if err := r.BaseRepository.Update(ctx, user); err != nil {
        return err
    }
    
    // Invalidate cache
    cacheKey := fmt.Sprintf("user:%d", user.ID)
    r.cache.Del(ctx, cacheKey)
    
    return nil
}
```

---

## Testing Repositories

### Mock Repository

```go
// Mock repository for testing
type mockRepository struct {
    users []*User
}

func NewMockRepository() Repository {
    return &mockRepository{
        users: []*User{},
    }
}

func (m *mockRepository) Create(ctx context.Context, user *User) error {
    user.ID = uint(len(m.users) + 1)
    m.users = append(m.users, user)
    return nil
}

func (m *mockRepository) FindByID(ctx context.Context, id uint) (*User, error) {
    for _, u := range m.users {
        if u.ID == id {
            return u, nil
        }
    }
    return nil, gorm.ErrRecordNotFound
}

func (m *mockRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
    for _, u := range m.users {
        if u.Email == email {
            return u, nil
        }
    }
    return nil, gorm.ErrRecordNotFound
}
```

### Repository Tests

```go
func TestRepository_FindByEmail(t *testing.T) {
    db := setupTestDB(t)
    defer teardownTestDB(t, db)
    
    repo := NewRepository(db)
    
    // Create test user
    user := &User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    err := repo.Create(context.Background(), user)
    require.NoError(t, err)
    
    // Test find by email
    found, err := repo.FindByEmail(context.Background(), "john@example.com")
    require.NoError(t, err)
    assert.Equal(t, user.ID, found.ID)
    assert.Equal(t, user.Email, found.Email)
    
    // Test not found
    _, err = repo.FindByEmail(context.Background(), "notfound@example.com")
    assert.Error(t, err)
}
```

---

## Best Practices

### ✅ DO:

**1. Use Context**
```go
func (r *repository) FindByID(ctx context.Context, id uint) (*User, error) {
    var user User
    err := r.DB().WithContext(ctx).First(&user, id).Error
    return &user, err
}
```

**2. Handle Errors**
```go
func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
    var user User
    err := r.DB().WithContext(ctx).Where("email = ?", email).First(&user).Error
    
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrUserNotFound
        }
        return nil, fmt.Errorf("find user by email: %w", err)
    }
    
    return &user, nil
}
```

**3. Use Preloading**
```go
func (r *repository) FindWithRelations(ctx context.Context, id uint) (*User, error) {
    var user User
    err := r.DB().WithContext(ctx).
        Preload("Profile").
        Preload("Orders").
        First(&user, id).Error
    return &user, err
}
```

### ❌ DON'T:

**1. N+1 Queries**
```go
// Bad: N+1 problem
users, _ := repo.FindAll(ctx)
for _, user := range users {
    orders, _ := orderRepo.FindByUserID(ctx, user.ID) // ❌
}

// Good: Use preload
users, _ := repo.FindAllWithOrders(ctx) // Preload("Orders")
```

**2. Business Logic in Repository**
```go
// Bad: Business logic in repository
func (r *repository) CreateUserWithValidation(ctx context.Context, user *User) error {
    if user.Age < 18 { // ❌ Business logic
        return errors.New("too young")
    }
    return r.DB().Create(user).Error
}

// Good: Keep repository simple
func (r *repository) Create(ctx context.Context, user *User) error {
    return r.DB().WithContext(ctx).Create(user).Error
}
```

---

## Next Steps

- [**Transactions**](transactions.md) - Managing database transactions
- [**Service Layer**](../core-concepts/service-layer.md) - Business logic layer
- [**Testing**](../development/testing.md) - Testing strategies
- [**Performance**](../advanced/performance.md) - Query optimization

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
