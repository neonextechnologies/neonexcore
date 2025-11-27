# Dependency Injection

Master Neonex Core's type-safe dependency injection container for managing service lifecycles and dependencies.

---

## Overview

Neonex Core includes a built-in **Dependency Injection (DI) Container** that manages object creation and lifetime. It supports two lifecycle scopes: **Singleton** and **Transient**.

**Benefits:**
- ✅ **Loose Coupling** - Depend on interfaces, not implementations
- ✅ **Testability** - Easy to mock dependencies
- ✅ **Lifecycle Management** - Automatic singleton management
- ✅ **Type Safety** - Generic-based resolution
- ✅ **No Reflection** - Fast performance

---

## Container Basics

### Creating a Container

```go
import "neonexcore/internal/core"

container := core.NewContainer()
```

### Provider Types

| Type | Behavior | Use Case |
|------|----------|----------|
| **Singleton** | Created once, reused | Database connections, services, configs |
| **Transient** | Created every time | Requests, temporary objects |

---

## Registration

### Singleton Registration

Created once and reused across the application:

```go
// Register database connection (Singleton)
container.Provide(func() *gorm.DB {
    db, _ := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
    return db
}, core.Singleton)

// Register service (Singleton)
container.Provide(func() UserService {
    db := core.Resolve[*gorm.DB](container)
    repo := NewUserRepository(db)
    return NewUserService(repo)
}, core.Singleton)
```

**When to use:**
- Database connections
- Configuration objects
- Stateless services
- Caches
- Connection pools

### Transient Registration

Created new every time it's requested:

```go
// Register request handler (Transient)
container.Provide(func() *RequestHandler {
    return &RequestHandler{
        ID: generateUniqueID(),
    }
}, core.Transient)
```

**When to use:**
- Request-scoped objects
- Objects with mutable state
- Short-lived operations
- Context-specific instances

---

## Resolution

### Generic Resolution

Type-safe resolution using Go generics:

```go
// Resolve database
db := core.Resolve[*gorm.DB](container)

// Resolve service
userService := core.Resolve[UserService](container)

// Resolve repository
userRepo := core.Resolve[UserRepository](container)
```

### Resolution in Modules

```go
// modules/user/di.go
package user

func RegisterDependencies(c *core.Container) {
    // 1. Register Repository
    c.Provide(func() Repository {
        db := core.Resolve[*gorm.DB](c)
        return NewRepository(db)
    }, core.Singleton)
    
    // 2. Register Service (depends on Repository)
    c.Provide(func() Service {
        repo := core.Resolve[Repository](c)
        return NewService(repo)
    }, core.Singleton)
    
    // 3. Register Controller (depends on Service)
    c.Provide(func() *Controller {
        service := core.Resolve[Service](c)
        return NewController(service)
    }, core.Singleton)
}
```

---

## Common Patterns

### Pattern 1: Interface-Based Registration

```go
// Define interface
type UserRepository interface {
    FindByID(id uint) (*User, error)
    Create(user *User) error
}

// Implement interface
type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

// Register interface, not implementation
container.Provide(func() UserRepository {
    db := core.Resolve[*gorm.DB](container)
    return NewUserRepository(db)
}, core.Singleton)

// Resolve by interface
repo := core.Resolve[UserRepository](container)
```

**Benefits:**
- Easy to swap implementations
- Better for testing (mock interfaces)
- Follows SOLID principles

### Pattern 2: Factory Pattern

```go
// Factory function
type ServiceFactory func() Service

// Register factory
container.Provide(func() ServiceFactory {
    return func() Service {
        repo := core.Resolve[Repository](container)
        return NewService(repo)
    }
}, core.Singleton)

// Use factory
factory := core.Resolve[ServiceFactory](container)
service1 := factory()
service2 := factory()
```

### Pattern 3: Configuration Injection

```go
// Register configuration
container.Provide(func() *Config {
    return LoadConfigFromEnv()
}, core.Singleton)

// Inject into service
container.Provide(func() Service {
    config := core.Resolve[*Config](container)
    repo := core.Resolve[Repository](container)
    return NewService(repo, config)
}, core.Singleton)
```

---

## Real-World Example

### Complete Module DI Setup

```go
// modules/product/di.go
package product

import (
    "neonexcore/internal/core"
    "neonexcore/pkg/logger"
    "gorm.io/gorm"
)

func RegisterDependencies(c *core.Container) error {
    // 1. Register Repository
    c.Provide(func() Repository {
        db := core.Resolve[*gorm.DB](c)
        return NewRepository(db)
    }, core.Singleton)
    
    // 2. Register Service with Logger
    c.Provide(func() Service {
        repo := core.Resolve[Repository](c)
        log := core.Resolve[logger.Logger](c)
        return NewService(repo, log)
    }, core.Singleton)
    
    // 3. Register Controller
    c.Provide(func() *Controller {
        service := core.Resolve[Service](c)
        return NewController(service)
    }, core.Singleton)
    
    // 4. Register Validator (Transient)
    c.Provide(func() *Validator {
        return &Validator{}
    }, core.Transient)
    
    return nil
}
```

### Service with Dependencies

```go
// modules/product/service.go
package product

import "neonexcore/pkg/logger"

type Service interface {
    CreateProduct(product *Product) error
    GetProduct(id uint) (*Product, error)
}

type service struct {
    repo   Repository
    logger logger.Logger
}

func NewService(repo Repository, log logger.Logger) Service {
    return &service{
        repo:   repo,
        logger: log,
    }
}

func (s *service) CreateProduct(product *Product) error {
    s.logger.Info("Creating product", logger.Fields{
        "name": product.Name,
    })
    
    if err := s.repo.Create(product); err != nil {
        s.logger.Error("Failed to create product", logger.Fields{
            "error": err.Error(),
        })
        return err
    }
    
    s.logger.Info("Product created", logger.Fields{
        "id": product.ID,
    })
    
    return nil
}
```

---

## Advanced Usage

### Conditional Registration

```go
func RegisterDependencies(c *core.Container) {
    // Register based on environment
    if os.Getenv("USE_REDIS") == "true" {
        c.Provide(func() CacheService {
            return NewRedisCache()
        }, core.Singleton)
    } else {
        c.Provide(func() CacheService {
            return NewMemoryCache()
        }, core.Singleton)
    }
}
```

### Multiple Implementations

```go
// Register named services
type SQLRepository struct{}
type MongoRepository struct{}

// Use different types
container.Provide(func() *SQLRepository {
    return &SQLRepository{}
}, core.Singleton)

container.Provide(func() *MongoRepository {
    return &MongoRepository{}
}, core.Singleton)

// Resolve specific implementation
sqlRepo := core.Resolve[*SQLRepository](container)
mongoRepo := core.Resolve[*MongoRepository](container)
```

### Lazy Loading

```go
// Service is only created when first resolved
container.Provide(func() ExpensiveService {
    // This runs only once (Singleton)
    // and only when first requested
    return InitializeExpensiveService()
}, core.Singleton)
```

---

## Testing with DI

### Mocking Dependencies

```go
// test/user_service_test.go
package test

import "testing"

// Mock repository
type mockRepository struct {
    users map[uint]*User
}

func (m *mockRepository) FindByID(id uint) (*User, error) {
    return m.users[id], nil
}

// Test with mock
func TestUserService(t *testing.T) {
    // Create container
    container := core.NewContainer()
    
    // Register mock
    container.Provide(func() Repository {
        return &mockRepository{
            users: map[uint]*User{
                1: {ID: 1, Name: "Test"},
            },
        }
    }, core.Singleton)
    
    // Register service
    container.Provide(func() Service {
        repo := core.Resolve[Repository](container)
        return NewService(repo)
    }, core.Singleton)
    
    // Test
    service := core.Resolve[Service](container)
    user, err := service.GetUser(1)
    
    if err != nil {
        t.Fatal(err)
    }
    if user.Name != "Test" {
        t.Errorf("Expected 'Test', got '%s'", user.Name)
    }
}
```

---

## Best Practices

### ✅ DO:

**1. Register Interfaces**
```go
// Good
c.Provide(func() UserRepository {
    return NewUserRepository()
}, core.Singleton)
```

**2. Use Constructor Functions**
```go
// Good
func NewUserService(repo Repository) Service {
    return &userService{repo: repo}
}
```

**3. Resolve in Factories**
```go
// Good
c.Provide(func() Service {
    repo := core.Resolve[Repository](c)
    return NewService(repo)
}, core.Singleton)
```

### ❌ DON'T:

**1. Store Container in Structs**
```go
// Bad
type Service struct {
    container *core.Container  // ❌ Service Locator anti-pattern
}
```

**2. Resolve Outside Factories**
```go
// Bad
repo := core.Resolve[Repository](container)
c.Provide(func() Service {
    return NewService(repo)  // ❌ Resolved too early
}, core.Singleton)
```

**3. Mix Singletons and Transients Incorrectly**
```go
// Bad: Transient depending on Singleton is OK
// Bad: Singleton depending on Transient is NOT OK
```

---

## Troubleshooting

### Nil Pointer Panic

**Problem:**
```
panic: runtime error: invalid memory address or nil pointer dereference
```

**Cause:** Dependency not registered before resolution.

**Solution:**
```go
// Ensure dependency is registered first
c.Provide(func() Repository {
    return NewRepository()
}, core.Singleton)

// Then register dependent service
c.Provide(func() Service {
    repo := core.Resolve[Repository](c)
    return NewService(repo)
}, core.Singleton)
```

### Type Not Found

**Problem:**
```
Error: type not found in container
```

**Solution:** Check exact type match:
```go
// Registration
c.Provide(func() UserRepository { ... }, core.Singleton)

// Resolution - must match exactly
repo := core.Resolve[UserRepository](c)  // ✅
repo := core.Resolve[*UserRepository](c) // ❌ Different type
```

---

## Next Steps

- [**Module System**](module-system.md) - Organize modules with DI
- [**Repository Pattern**](repository-pattern.md) - Data access with DI
- [**Service Layer**](service-layer.md) - Business logic with DI
- [**Testing Guide**](../development/testing.md) - Test with mocks

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
