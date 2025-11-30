# Repository API

Complete API reference for Neonex Core's Generic Repository Pattern.

---

## Overview

The Repository API provides a type-safe, generic abstraction layer for data access operations. Built on top of GORM, it offers common CRUD operations with compile-time type checking using Go generics.

**Package:** `neonexcore/pkg/database`

**Features:**
- 🎯 **Type-Safe** - Generic-based, compile-time type checking
- 🔄 **Reusable** - DRY principle with base repository
- 📦 **Extensible** - Easy to add custom methods
- ⚡ **Performance** - Optimized GORM queries
- 🧪 **Testable** - Easy to mock for testing
- 🔍 **Context-Aware** - Supports cancellation and timeouts

---

## Repository Interface

### Definition

```go
type Repository[T any] interface {
    // Create operations
    Create(ctx context.Context, entity *T) error
    CreateBatch(ctx context.Context, entities []*T) error
    
    // Read operations
    FindByID(ctx context.Context, id interface{}) (*T, error)
    FindAll(ctx context.Context) ([]*T, error)
    FindByCondition(ctx context.Context, condition interface{}, args ...interface{}) ([]*T, error)
    FindOne(ctx context.Context, condition interface{}, args ...interface{}) (*T, error)
    
    // Update operations
    Update(ctx context.Context, entity *T) error
    UpdateFields(ctx context.Context, id interface{}, fields map[string]interface{}) error
    
    // Delete operations
    Delete(ctx context.Context, id interface{}) error
    DeleteBatch(ctx context.Context, ids []interface{}) error
    
    // Query operations
    Count(ctx context.Context, condition interface{}, args ...interface{}) (int64, error)
    Exists(ctx context.Context, condition interface{}, args ...interface{}) (bool, error)
    Paginate(ctx context.Context, page, pageSize int) ([]*T, int64, error)
    
    // Transaction support
    WithTx(tx *gorm.DB) Repository[T]
    GetDB() *gorm.DB
}
```

Generic repository interface with type parameter `T` for the entity type.

---

## Base Repository

### Struct Definition

```go
type BaseRepository[T any] struct {
    db *gorm.DB
}
```

Base implementation providing all common repository operations.

### Constructor

```go
func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T]
```

Creates a new base repository instance.

**Parameters:**
- `db` - GORM database instance

**Returns:**
- `*BaseRepository[T]` - New repository

**Example:**
```go
// Generic repository for User entity
userRepo := database.NewBaseRepository[User](db)

// Generic repository for Product entity
productRepo := database.NewBaseRepository[Product](db)
```

---

## Create Operations

### Create

```go
Create(ctx context.Context, entity *T) error
```

Inserts a single entity into the database.

**Parameters:**
- `ctx` - Context for cancellation/timeout
- `entity` - Pointer to entity to create

**Returns:**
- `error` - Error if creation fails

**Example:**
```go
user := &User{
    Name:  "John Doe",
    Email: "john@example.com",
}

err := repo.Create(ctx, user)
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}

// user.ID now contains the auto-generated ID
fmt.Printf("Created user ID: %d\n", user.ID)
```

**Behavior:**
- Auto-generates ID (if using auto-increment)
- Sets CreatedAt timestamp (if model has it)
- Returns populated entity with ID
- Validates constraints

---

### CreateBatch

```go
CreateBatch(ctx context.Context, entities []*T) error
```

Inserts multiple entities in a single database operation.

**Parameters:**
- `ctx` - Context
- `entities` - Slice of entity pointers

**Returns:**
- `error` - Error if batch creation fails

**Example:**
```go
users := []*User{
    {Name: "User 1", Email: "user1@example.com"},
    {Name: "User 2", Email: "user2@example.com"},
    {Name: "User 3", Email: "user3@example.com"},
}

err := repo.CreateBatch(ctx, users)
if err != nil {
    return fmt.Errorf("batch creation failed: %w", err)
}

// All users now have IDs
for _, user := range users {
    fmt.Printf("Created user ID: %d\n", user.ID)
}
```

**Performance:**
- More efficient than multiple Create() calls
- Single database round-trip
- Transaction-safe

---

## Read Operations

### FindByID

```go
FindByID(ctx context.Context, id interface{}) (*T, error)
```

Finds an entity by its primary key.

**Parameters:**
- `ctx` - Context
- `id` - Primary key value (uint, string, etc.)

**Returns:**
- `*T` - Found entity or nil
- `error` - Error if query fails

**Example:**
```go
user, err := repo.FindByID(ctx, 123)
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, fmt.Errorf("user not found")
    }
    return nil, err
}

fmt.Printf("Found: %s\n", user.Name)
```

**Behavior:**
- Returns `gorm.ErrRecordNotFound` if not found
- Includes soft-deleted check (if using DeletedAt)
- Supports composite primary keys

---

### FindAll

```go
FindAll(ctx context.Context) ([]*T, error)
```

Retrieves all entities from the table.

**Parameters:**
- `ctx` - Context

**Returns:**
- `[]*T` - Slice of entities
- `error` - Error if query fails

**Example:**
```go
users, err := repo.FindAll(ctx)
if err != nil {
    return nil, err
}

fmt.Printf("Found %d users\n", len(users))
for _, user := range users {
    fmt.Printf("- %s\n", user.Name)
}
```

**⚠️ Warning:** Use with caution on large tables. Consider pagination.

---

### FindByCondition

```go
FindByCondition(ctx context.Context, condition interface{}, args ...interface{}) ([]*T, error)
```

Finds entities matching specified conditions.

**Parameters:**
- `ctx` - Context
- `condition` - WHERE clause (string or map)
- `args` - Values for placeholders

**Returns:**
- `[]*T` - Matching entities
- `error` - Error if query fails

**Examples:**

**String Condition:**
```go
// Find active users
users, err := repo.FindByCondition(ctx, "active = ?", true)

// Find by email domain
users, err := repo.FindByCondition(ctx, "email LIKE ?", "%@example.com")

// Multiple conditions
users, err := repo.FindByCondition(ctx, 
    "role = ? AND created_at > ?", 
    "admin", 
    time.Now().AddDate(0, -1, 0),
)
```

**Map Condition:**
```go
// Find by exact match
users, err := repo.FindByCondition(ctx, map[string]interface{}{
    "role":   "admin",
    "active": true,
})
```

---

### FindOne

```go
FindOne(ctx context.Context, condition interface{}, args ...interface{}) (*T, error)
```

Finds first entity matching condition.

**Parameters:**
- `ctx` - Context
- `condition` - WHERE clause
- `args` - Values for placeholders

**Returns:**
- `*T` - First matching entity or nil
- `error` - Error if query fails

**Example:**
```go
// Find by email
user, err := repo.FindOne(ctx, "email = ?", "john@example.com")
if err != nil {
    return nil, err
}

// Find by username
user, err := repo.FindOne(ctx, "username = ?", "john_doe")
```

---

### Count

```go
Count(ctx context.Context, condition interface{}, args ...interface{}) (int64, error)
```

Counts entities matching condition.

**Parameters:**
- `ctx` - Context
- `condition` - WHERE clause (optional)
- `args` - Values for placeholders

**Returns:**
- `int64` - Number of matching entities
- `error` - Error if query fails

**Examples:**
```go
// Count all
total, err := repo.Count(ctx, "")

// Count active users
activeCount, err := repo.Count(ctx, "active = ?", true)

// Count by role
adminCount, err := repo.Count(ctx, "role = ?", "admin")

fmt.Printf("Total: %d, Active: %d, Admins: %d\n", 
    total, activeCount, adminCount)
```

---

### Exists

```go
Exists(ctx context.Context, condition interface{}, args ...interface{}) (bool, error)
```

Checks if any entity matches condition.

**Parameters:**
- `ctx` - Context
- `condition` - WHERE clause
- `args` - Values for placeholders

**Returns:**
- `bool` - true if exists, false otherwise
- `error` - Error if query fails

**Examples:**
```go
// Check if email exists
exists, err := repo.Exists(ctx, "email = ?", "test@example.com")
if exists {
    return errors.New("email already taken")
}

// Check if username exists
exists, err := repo.Exists(ctx, "username = ?", "john_doe")
```

---

### Paginate

```go
Paginate(ctx context.Context, page, pageSize int) ([]*T, int64, error)
```

Retrieves paginated results with total count.

**Parameters:**
- `ctx` - Context
- `page` - Page number (1-indexed)
- `pageSize` - Number of items per page

**Returns:**
- `[]*T` - Entities for the page
- `int64` - Total count across all pages
- `error` - Error if query fails

**Example:**
```go
page := 1
pageSize := 20

users, total, err := repo.Paginate(ctx, page, pageSize)
if err != nil {
    return nil, err
}

fmt.Printf("Page %d of %d\n", page, (total+int64(pageSize)-1)/int64(pageSize))
fmt.Printf("Showing %d-%d of %d users\n", 
    (page-1)*pageSize+1, 
    (page-1)*pageSize+len(users), 
    total,
)
```

**Response Format:**
```go
type PaginatedResponse struct {
    Data       []*User `json:"data"`
    Page       int     `json:"page"`
    PageSize   int     `json:"page_size"`
    Total      int64   `json:"total"`
    TotalPages int     `json:"total_pages"`
}
```

---

## Update Operations

### Update

```go
Update(ctx context.Context, entity *T) error
```

Updates all fields of an existing entity.

**Parameters:**
- `ctx` - Context
- `entity` - Entity with updated values (must have ID)

**Returns:**
- `error` - Error if update fails

**Example:**
```go
// Get user
user, _ := repo.FindByID(ctx, 123)

// Modify
user.Name = "Jane Doe"
user.Email = "jane@example.com"

// Update
err := repo.Update(ctx, user)
if err != nil {
    return fmt.Errorf("update failed: %w", err)
}
```

**Behavior:**
- Updates all non-zero fields
- Sets UpdatedAt timestamp
- Requires entity to have ID
- Uses optimistic locking if available

---

### UpdateFields

```go
UpdateFields(ctx context.Context, id interface{}, fields map[string]interface{}) error
```

Updates specific fields without loading the entity.

**Parameters:**
- `ctx` - Context
- `id` - Primary key
- `fields` - Map of field names to values

**Returns:**
- `error` - Error if update fails

**Example:**
```go
// Update specific fields
err := repo.UpdateFields(ctx, 123, map[string]interface{}{
    "name":       "New Name",
    "updated_at": time.Now(),
})

// Increment counter
err := repo.UpdateFields(ctx, 123, map[string]interface{}{
    "login_count": gorm.Expr("login_count + 1"),
})

// Update timestamp
err := repo.UpdateFields(ctx, 123, map[string]interface{}{
    "last_login_at": time.Now(),
})
```

**Performance:**
- More efficient than Update() for partial updates
- No need to load entity first
- Atomic operation

---

## Delete Operations

### Delete

```go
Delete(ctx context.Context, id interface{}) error
```

Deletes an entity by ID (soft delete if model supports it).

**Parameters:**
- `ctx` - Context
- `id` - Primary key

**Returns:**
- `error` - Error if deletion fails

**Examples:**
```go
// Soft delete (if model has DeletedAt field)
err := repo.Delete(ctx, 123)

// Hard delete (permanently remove)
err := repo.GetDB().Unscoped().Delete(&User{}, 123).Error
```

**Behavior:**
- Soft delete if model has `gorm.DeletedAt` field
- Hard delete otherwise
- Sets DeletedAt timestamp for soft deletes
- Does not physically remove row for soft deletes

---

### DeleteBatch

```go
DeleteBatch(ctx context.Context, ids []interface{}) error
```

Deletes multiple entities by their IDs.

**Parameters:**
- `ctx` - Context
- `ids` - Slice of primary keys

**Returns:**
- `error` - Error if deletion fails

**Example:**
```go
ids := []interface{}{1, 2, 3, 4, 5}
err := repo.DeleteBatch(ctx, ids)
if err != nil {
    return fmt.Errorf("batch delete failed: %w", err)
}

fmt.Printf("Deleted %d users\n", len(ids))
```

---

## Transaction Support

### WithTx

```go
WithTx(tx *gorm.DB) Repository[T]
```

Returns a new repository instance using the provided transaction.

**Parameters:**
- `tx` - GORM transaction instance

**Returns:**
- `Repository[T]` - Repository using transaction

**Example:**
```go
err := database.WithTransaction(ctx, db, func(tx *gorm.DB) error {
    // Create transactional repository
    txRepo := repo.WithTx(tx)
    
    // All operations use transaction
    user := &User{Name: "Test"}
    if err := txRepo.Create(ctx, user); err != nil {
        return err  // Rolls back
    }
    
    profile := &Profile{UserID: user.ID}
    if err := profileRepo.WithTx(tx).Create(ctx, profile); err != nil {
        return err  // Rolls back
    }
    
    return nil  // Commits
})
```

---

### GetDB

```go
GetDB() *gorm.DB
```

Returns the underlying GORM database instance.

**Returns:**
- `*gorm.DB` - Database instance

**Example:**
```go
db := repo.GetDB()

// Custom query
var results []CustomResult
err := db.Raw("SELECT ... FROM users WHERE ...").Scan(&results).Error

// Complex join
err := db.
    Table("users").
    Select("users.*, profiles.*").
    Joins("LEFT JOIN profiles ON users.id = profiles.user_id").
    Where("users.active = ?", true).
    Scan(&results).Error
```

---

## Custom Repository Example

### Extending Base Repository

```go
// modules/user/repository.go
package user

import (
    "context"
    "neonexcore/pkg/database"
    "gorm.io/gorm"
)

// Interface with custom methods
type Repository interface {
    database.Repository[User]  // Embed base interface
    
    // Custom methods
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindActive(ctx context.Context) ([]*User, error)
    FindByRole(ctx context.Context, role string) ([]*User, error)
    UpdateLastLogin(ctx context.Context, id uint) error
    Search(ctx context.Context, query string) ([]*User, error)
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

// Custom method implementations
func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
    return r.FindOne(ctx, "email = ?", email)
}

func (r *repository) FindActive(ctx context.Context) ([]*User, error) {
    return r.FindByCondition(ctx, "active = ?", true)
}

func (r *repository) FindByRole(ctx context.Context, role string) ([]*User, error) {
    return r.FindByCondition(ctx, "role = ?", role)
}

func (r *repository) UpdateLastLogin(ctx context.Context, id uint) error {
    return r.UpdateFields(ctx, id, map[string]interface{}{
        "last_login_at": time.Now(),
    })
}

func (r *repository) Search(ctx context.Context, query string) ([]*User, error) {
    searchQuery := fmt.Sprintf("%%%s%%", query)
    return r.FindByCondition(ctx, 
        "name LIKE ? OR email LIKE ?", 
        searchQuery, searchQuery,
    )
}
```

---

## Best Practices

### ✅ DO: Use Context

```go
// Good: Pass context
user, err := repo.FindByID(ctx, id)

// Bad: No context
user, err := repo.FindByID(context.Background(), id)
```

### ✅ DO: Handle Errors

```go
// Good: Proper error handling
user, err := repo.FindByID(ctx, id)
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, ErrUserNotFound
    }
    return nil, fmt.Errorf("database error: %w", err)
}
```

### ✅ DO: Use Transactions for Multiple Operations

```go
// Good: Atomic operations
err := database.WithTransaction(ctx, db, func(tx *gorm.DB) error {
    txRepo := repo.WithTx(tx)
    
    if err := txRepo.Create(ctx, user); err != nil {
        return err
    }
    
    if err := txRepo.Create(ctx, profile); err != nil {
        return err
    }
    
    return nil
})
```

### ❌ DON'T: Ignore Record Not Found

```go
// Bad: Treats not found as error
user, err := repo.FindByID(ctx, id)
if err != nil {
    return err  // Wrong!
}

// Good: Handle not found separately
user, err := repo.FindByID(ctx, id)
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return ErrNotFound
    }
    return err
}
```

---

## Related Documentation

- [**Repository Pattern Guide**](../core-concepts/repository-pattern.md)
- [**Transactions**](../database/transactions.md)
- [**GORM Documentation**](https://gorm.io/docs/)

---

**Need help?** [FAQ](../resources/faq.md) | [Support](../resources/support.md)
