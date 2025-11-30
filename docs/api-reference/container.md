# Container API

Complete API reference for Neonex Core's type-safe Dependency Injection Container.

---

## Overview

The Container provides **type-safe dependency injection** using Go generics. It manages service lifecycles, resolves dependencies automatically, and ensures thread-safe concurrent access.

**Package:** `neonexcore/internal/core`

**Key Features:**
- 🎯 **Type-Safe** - No reflection, compile-time type checking
- 🔄 **Lifecycle Management** - Singleton and Transient scopes
- 🧵 **Thread-Safe** - Concurrent access without race conditions
- ⚡ **Fast** - Zero-allocation resolution for singletons
- 🧪 **Testable** - Easy mocking for unit tests

---

## Types

### Container

```go
type Container struct {
    services  map[string]ServiceDefinition
    instances map[string]interface{}
    mu        sync.RWMutex
}
```

Main dependency injection container that stores service definitions and manages instances.

**Fields:**
- `services` - Map of service type names to their factory definitions
- `instances` - Cache of singleton instances
- `mu` - Read-write mutex for thread-safe operations

**Example:**
```go
container := core.NewContainer()
fmt.Printf("Container created: %+v\n", container)
```

---

### ServiceDefinition

```go
type ServiceDefinition struct {
    Factory   interface{}
    Lifecycle Lifecycle
}
```

Defines how a service should be created and its lifecycle.

**Fields:**
- `Factory` - Function that creates the service (signature: `func() T`)
- `Lifecycle` - Either `Singleton` or `Transient`

---

### Lifecycle

```go
type Lifecycle int

const (
    Singleton Lifecycle = iota  // Created once, shared across app
    Transient                    // Created every time requested
)
```

Defines service instance lifecycle behavior.

**Singleton:**
- Created only once on first resolution
- Cached and reused for all subsequent requests
- Best for: Database connections, configurations, stateless services

**Transient:**
- New instance created for every resolution
- No caching
- Best for: Request handlers, stateful objects, temporary operations

---

## Core Functions

### NewContainer

```go
func NewContainer() *Container
```

Creates and returns a new Container instance with initialized internal maps.

**Returns:**
- `*Container` - Newly created container ready for use

**Example:**
```go
// Create new container
container := core.NewContainer()

// Container is ready to register and resolve services
fmt.Println("Container initialized")
```

**Thread Safety:** Safe to create multiple containers concurrently.

---

### Provide

```go
func (c *Container) Provide(factory interface{}, lifecycle Lifecycle) error
```

Registers a service factory function with the specified lifecycle in the container.

**Parameters:**
- `factory` - Factory function that creates the service (must be `func() T`)
- `lifecycle` - Service lifecycle (`core.Singleton` or `core.Transient`)

**Returns:**
- `error` - Error if factory signature is invalid, `nil` on success

**Factory Function Signature:**
```go
func() T  // Where T is any type
```

**Examples:**

**Singleton Registration:**
```go
// Database connection (shared instance)
container.Provide(func() *gorm.DB {
    db, _ := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
    return db
}, core.Singleton)

// Configuration (shared)
container.Provide(func() *Config {
    return &Config{
        Host: "localhost",
        Port: 8080,
    }
}, core.Singleton)

// Logger (shared)
container.Provide(func() logger.Logger {
    return logger.NewLogger()
}, core.Singleton)
```

**Transient Registration:**
```go
// Request handler (new instance each time)
container.Provide(func() *RequestHandler {
    return &RequestHandler{
        ID:        generateUniqueID(),
        Timestamp: time.Now(),
    }
}, core.Transient)

// Temporary processor
container.Provide(func() *DataProcessor {
    return &DataProcessor{
        Buffer: make([]byte, 1024),
    }
}, core.Transient)
```

**With Dependencies:**
```go
// Repository depends on database
container.Provide(func() UserRepository {
    db := core.Resolve[*gorm.DB](container)
    return NewUserRepository(db)
}, core.Singleton)

// Service depends on repository and logger
container.Provide(func() UserService {
    repo := core.Resolve[UserRepository](container)
    log := core.Resolve[logger.Logger](container)
    return NewUserService(repo, log)
}, core.Singleton)

// Controller depends on service
container.Provide(func() *UserController {
    service := core.Resolve[UserService](container)
    return NewUserController(service)
}, core.Singleton)
```

**Error Handling:**
```go
err := container.Provide(func() *Service {
    return &Service{}
}, core.Singleton)

if err != nil {
    log.Fatalf("Failed to register service: %v", err)
}
```

---

### Resolve

```go
func Resolve[T any](c *Container) T
```

Resolves and returns a service instance of type `T` from the container.

**Type Parameters:**
- `T` - The service type to resolve

**Parameters:**
- `c` - Container instance

**Returns:**
- `T` - Service instance

**Behavior:**
- **Singleton**: Returns cached instance (creates on first call)
- **Transient**: Creates and returns new instance every time

**Examples:**

**Basic Resolution:**
```go
// Resolve database
db := core.Resolve[*gorm.DB](container)
fmt.Printf("Database: %+v\n", db)

// Resolve logger
log := core.Resolve[logger.Logger](container)
log.Info("Logger resolved")

// Resolve config
config := core.Resolve[*Config](container)
fmt.Printf("Port: %d\n", config.Port)
```

**Interface Resolution:**
```go
// Resolve by interface type
repo := core.Resolve[UserRepository](container)
users, err := repo.FindAll(ctx)

cache := core.Resolve[CacheService](container)
cache.Set("key", "value")
```

**In Module Registration:**
```go
func RegisterDependencies(c *core.Container) {
    // Register repository
    c.Provide(func() Repository {
        db := core.Resolve[*gorm.DB](c)
        return NewRepository(db)
    }, core.Singleton)
    
    // Register service
    c.Provide(func() Service {
        repo := core.Resolve[Repository](c)
        return NewService(repo)
    }, core.Singleton)
}
```

**Multiple Resolutions:**
```go
// Get multiple dependencies
db := core.Resolve[*gorm.DB](container)
logger := core.Resolve[logger.Logger](container)
cache := core.Resolve[CacheService](container)

service := NewComplexService(db, logger, cache)
```

**Panic on Not Found:**
```go
// Will panic if type not registered
defer func() {
    if r := recover(); r != nil {
        log.Printf("Service not found: %v", r)
    }
}()

service := core.Resolve[UnregisteredService](container)
```

---

### Has

```go
func (c *Container) Has(typeName string) bool
```

Checks if a service type is registered in the container.

**Parameters:**
- `typeName` - Fully qualified type name (e.g., `"*gorm.DB"`)

**Returns:**
- `bool` - `true` if registered, `false` otherwise

**Examples:**

```go
// Check if service exists
if container.Has("*gorm.DB") {
    db := core.Resolve[*gorm.DB](container)
    // Use database
}

// Conditional resolution
if container.Has("CacheService") {
    cache := core.Resolve[CacheService](container)
    cache.Clear()
} else {
    log.Println("Cache service not available")
}

// Feature flags
if container.Has("AnalyticsService") {
    analytics := core.Resolve[AnalyticsService](container)
    analytics.Track("event")
}
```

---

### Reset

```go
func (c *Container) Reset()
```

Clears all singleton instances from the cache while keeping service registrations.

**Use Cases:**
- Testing: Reset state between tests
- Dynamic reconfiguration
- Memory cleanup

**Examples:**

```go
// Testing scenario
func TestWithFreshContainer(t *testing.T) {
    container := setupContainer()
    
    // Run first test
    service := core.Resolve[Service](container)
    service.DoWork()
    
    // Reset for second test
    container.Reset()
    
    // New instance will be created
    service2 := core.Resolve[Service](container)
    // service != service2 (new instance)
}

// Memory cleanup
func CleanupSingletons() {
    container.Reset()
    runtime.GC()
}
```

**Note:** Registrations remain intact - only cached instances are cleared.

---

## Advanced Usage

### Conditional Registration

Register different implementations based on environment or configuration:

```go
func RegisterCache(c *core.Container) {
    if os.Getenv("USE_REDIS") == "true" {
        // Redis cache in production
        c.Provide(func() CacheService {
            return NewRedisCache(
                os.Getenv("REDIS_HOST"),
                os.Getenv("REDIS_PORT"),
            )
        }, core.Singleton)
    } else {
        // Memory cache in development
        c.Provide(func() CacheService {
            return NewMemoryCache()
        }, core.Singleton)
    }
}

// Usage
RegisterCache(container)
cache := core.Resolve[CacheService](container) // Gets correct implementation
```

### Factory Pattern

Use factories to create multiple instances with different configurations:

```go
// Register factory function
container.Provide(func() ServiceFactory {
    return func(config string) Service {
        return NewService(config)
    }
}, core.Singleton)

// Use factory
factory := core.Resolve[ServiceFactory](container)
serviceA := factory("configA")
serviceB := factory("configB")
```

### Multiple Implementations

Register different implementations as different types:

```go
// SQL implementation
container.Provide(func() *SQLRepository {
    db := core.Resolve[*gorm.DB](container)
    return NewSQLRepository(db)
}, core.Singleton)

// MongoDB implementation  
container.Provide(func() *MongoRepository {
    client := core.Resolve[*mongo.Client](container)
    return NewMongoRepository(client)
}, core.Singleton)

// Use specific implementation
sqlRepo := core.Resolve[*SQLRepository](container)
mongoRepo := core.Resolve[*MongoRepository](container)
```

### Lazy Initialization

Services are only created when first requested:

```go
// Expensive service only registered
container.Provide(func() *ExpensiveService {
    log.Println("Creating expensive service...")
    time.Sleep(5 * time.Second) // Expensive initialization
    return &ExpensiveService{}
}, core.Singleton)

// Not created yet
log.Println("Service registered")

// Created on first resolve
service := core.Resolve[*ExpensiveService](container)
log.Println("Service created and ready")
```

### Scoped Containers

Create child containers for scoped scenarios:

```go
// Request-scoped container
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    // Create scoped container
    requestContainer := core.NewContainer()
    
    // Register request-specific services
    requestContainer.Provide(func() *http.Request {
        return r
    }, core.Singleton)
    
    requestContainer.Provide(func() RequestContext {
        req := core.Resolve[*http.Request](requestContainer)
        return NewRequestContext(req)
    }, core.Singleton)
    
    // Use scoped services
    ctx := core.Resolve[RequestContext](requestContainer)
    // Process request
}
```

---

## Best Practices

### ✅ DO: Register Interfaces

Register and resolve by interface types for loose coupling:

```go
// Good: Interface-based
c.Provide(func() UserRepository {
    db := core.Resolve[*gorm.DB](c)
    return NewUserRepository(db)
}, core.Singleton)

// Resolve by interface
repo := core.Resolve[UserRepository](container)
```

**Benefits:**
- Easy to swap implementations
- Better for testing (mock interfaces)
- Follows dependency inversion principle

### ✅ DO: Use Appropriate Lifecycle

Choose lifecycle based on service characteristics:

```go
// Singleton - stateless, shared resources
c.Provide(func() *gorm.DB { ... }, core.Singleton)
c.Provide(func() *Config { ... }, core.Singleton)
c.Provide(func() logger.Logger { ... }, core.Singleton)
c.Provide(func() UserService { ... }, core.Singleton)

// Transient - stateful, request-specific
c.Provide(func() *RequestHandler { ... }, core.Transient)
c.Provide(func() *TemporaryBuffer { ... }, core.Transient)
```

### ✅ DO: Resolve in Factory Functions

Always resolve dependencies inside factory functions:

```go
// Good: Lazy resolution
c.Provide(func() UserService {
    repo := core.Resolve[UserRepository](c)  // Resolved when service created
    logger := core.Resolve[logger.Logger](c)
    return NewUserService(repo, logger)
}, core.Singleton)

// Bad: Early resolution
repo := core.Resolve[UserRepository](c)  // Resolved immediately
c.Provide(func() UserService {
    return NewUserService(repo, nil)  // repo might not be ready
}, core.Singleton)
```

### ✅ DO: Use Constructor Functions

Create services using constructor functions:

```go
// Good: Constructor function
func NewUserService(repo UserRepository, log logger.Logger) UserService {
    return &userService{
        repo:   repo,
        logger: log,
    }
}

c.Provide(func() UserService {
    return NewUserService(
        core.Resolve[UserRepository](c),
        core.Resolve[logger.Logger](c),
    )
}, core.Singleton)
```

### ❌ DON'T: Store Container in Services

Avoid Service Locator anti-pattern:

```go
// Bad: Service Locator
type Service struct {
    container *core.Container  // ❌ Don't do this
}

func (s *Service) DoWork() {
    repo := core.Resolve[Repository](s.container)  // ❌ Anti-pattern
}

// Good: Dependency Injection
type service struct {
    repo Repository  // ✅ Injected dependency
}

func (s *service) DoWork() {
    s.repo.FindAll()  // ✅ Use injected dependency
}
```

### ❌ DON'T: Create Circular Dependencies

Avoid services that depend on each other:

```go
// Bad: Circular dependency
c.Provide(func() ServiceA {
    b := core.Resolve[ServiceB](c)
    return NewServiceA(b)
}, core.Singleton)

c.Provide(func() ServiceB {
    a := core.Resolve[ServiceA](c)  // ❌ Circular!
    return NewServiceB(a)
}, core.Singleton)

// Good: Extract shared logic
c.Provide(func() SharedLogic { ... }, core.Singleton)

c.Provide(func() ServiceA {
    shared := core.Resolve[SharedLogic](c)
    return NewServiceA(shared)
}, core.Singleton)

c.Provide(func() ServiceB {
    shared := core.Resolve[SharedLogic](c)
    return NewServiceB(shared)
}, core.Singleton)
```

### ❌ DON'T: Mix Transient and Singleton Incorrectly

Be careful with lifecycle dependencies:

```go
// OK: Transient depends on Singleton
c.Provide(func() Database { ... }, core.Singleton)
c.Provide(func() *Handler {
    db := core.Resolve[Database](c)  // ✅ OK
    return &Handler{db: db}
}, core.Transient)

// Problematic: Singleton depends on Transient
c.Provide(func() *Handler { ... }, core.Transient)
c.Provide(func() Service {
    handler := core.Resolve[*Handler](c)  // ⚠️ Gets first transient instance only
    return &Service{handler: handler}
}, core.Singleton)
```

---

## Thread Safety

The Container is **fully thread-safe** for concurrent registration and resolution:

```go
var wg sync.WaitGroup

// Concurrent registration (safe)
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        container.Provide(func() *Service {
            return &Service{ID: id}
        }, core.Transient)
    }(i)
}

// Concurrent resolution (safe)
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        db := core.Resolve[*gorm.DB](container)
        db.Exec("SELECT 1")
    }()
}

wg.Wait()
```

**Guarantees:**
- No race conditions
- Singleton created only once even with concurrent access
- Safe to use from multiple goroutines

---

## Performance Characteristics

### Singleton Performance

```go
// First resolution: Creates instance
start := time.Now()
db1 := core.Resolve[*gorm.DB](container)
fmt.Printf("First resolve: %v\n", time.Since(start))
// Output: First resolve: ~1ms

// Subsequent resolutions: Returns cached instance
start = time.Now()
db2 := core.Resolve[*gorm.DB](container)
fmt.Printf("Cached resolve: %v\n", time.Since(start))
// Output: Cached resolve: ~0.001ms (1000x faster)

// Same instance
fmt.Println(db1 == db2) // true
```

### Transient Performance

```go
// Every resolution creates new instance
start := time.Now()
handler1 := core.Resolve[*Handler](container)
fmt.Printf("Transient resolve 1: %v\n", time.Since(start))
// Output: ~0.1ms

start = time.Now()
handler2 := core.Resolve[*Handler](container)
fmt.Printf("Transient resolve 2: %v\n", time.Since(start))
// Output: ~0.1ms

// Different instances
fmt.Println(handler1 == handler2) // false
```

### Memory Usage

- **Singleton**: One instance stored in memory
- **Transient**: No caching, instances eligible for GC immediately after use

---

## Testing with Container

### Unit Testing with Mocks

```go
func TestUserService_CreateUser(t *testing.T) {
    // Create test container
    container := core.NewContainer()
    
    // Register mock repository
    container.Provide(func() UserRepository {
        return &MockUserRepository{
            users: map[uint]*User{},
        }
    }, core.Singleton)
    
    // Register mock logger
    container.Provide(func() logger.Logger {
        return &MockLogger{}
    }, core.Singleton)
    
    // Register service under test
    container.Provide(func() UserService {
        repo := core.Resolve[UserRepository](container)
        log := core.Resolve[logger.Logger](container)
        return NewUserService(repo, log)
    }, core.Singleton)
    
    // Test
    service := core.Resolve[UserService](container)
    user := &User{Name: "Test", Email: "test@example.com"}
    err := service.CreateUser(context.Background(), user)
    
    // Assert
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
}
```

### Table-Driven Tests

```go
func TestUserService(t *testing.T) {
    tests := []struct {
        name    string
        setup   func(*core.Container)
        test    func(*testing.T, UserService)
    }{
        {
            name: "creates user successfully",
            setup: func(c *core.Container) {
                c.Provide(func() UserRepository {
                    return &MockUserRepository{}
                }, core.Singleton)
            },
            test: func(t *testing.T, s UserService) {
                err := s.CreateUser(ctx, &User{Name: "Test"})
                assert.NoError(t, err)
            },
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            container := core.NewContainer()
            tt.setup(container)
            
            container.Provide(func() UserService {
                repo := core.Resolve[UserRepository](container)
                return NewUserService(repo)
            }, core.Singleton)
            
            service := core.Resolve[UserService](container)
            tt.test(t, service)
        })
    }
}
```

---

## Common Patterns

### Complete Module Registration

```go
// modules/user/di.go
package user

import "neonexcore/internal/core"

func RegisterDependencies(c *core.Container) {
    // 1. Repository layer
    c.Provide(func() Repository {
        db := core.Resolve[*gorm.DB](c)
        return NewRepository(db)
    }, core.Singleton)
    
    // 2. Service layer
    c.Provide(func() Service {
        repo := core.Resolve[Repository](c)
        logger := core.Resolve[logger.Logger](c)
        cache := core.Resolve[CacheService](c)
        return NewService(repo, logger, cache)
    }, core.Singleton)
    
    // 3. Controller layer
    c.Provide(func() *Controller {
        service := core.Resolve[Service](c)
        return NewController(service)
    }, core.Singleton)
    
    // 4. Validators (Transient)
    c.Provide(func() *Validator {
        return &Validator{}
    }, core.Transient)
}
```

### Configuration-Based Registration

```go
type AppConfig struct {
    UseRedis      bool
    UseElastic    bool
    EnableMetrics bool
}

func RegisterServices(c *core.Container, config *AppConfig) {
    // Cache service based on config
    if config.UseRedis {
        c.Provide(func() CacheService {
            return NewRedisCache()
        }, core.Singleton)
    } else {
        c.Provide(func() CacheService {
            return NewMemoryCache()
        }, core.Singleton)
    }
    
    // Search service based on config
    if config.UseElastic {
        c.Provide(func() SearchService {
            return NewElasticSearch()
        }, core.Singleton)
    } else {
        c.Provide(func() SearchService {
            return NewBasicSearch()
        }, core.Singleton)
    }
    
    // Metrics (optional)
    if config.EnableMetrics {
        c.Provide(func() MetricsCollector {
            return NewPrometheusCollector()
        }, core.Singleton)
    }
}
```

---

## Error Handling

### Registration Errors

```go
err := container.Provide(invalidFactory, core.Singleton)
if err != nil {
    log.Fatalf("Registration failed: %v", err)
}
```

### Resolution Errors (Panic)

```go
defer func() {
    if r := recover(); r != nil {
        log.Printf("Resolution failed: %v", r)
        // Handle panic gracefully
    }
}()

service := core.Resolve[UnregisteredService](container)
```

### Safe Resolution

```go
func SafeResolve[T any](c *core.Container) (T, error) {
    var result T
    
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("resolution failed: %v", r)
        }
    }()
    
    result = core.Resolve[T](c)
    return result, nil
}

// Usage
service, err := SafeResolve[UserService](container)
if err != nil {
    log.Printf("Failed to resolve: %v", err)
}
```

---

## Related Documentation

- [**Core API**](core.md) - Core framework APIs
- [**Module Interface**](module-interface.md) - Module system
- [**Dependency Injection Guide**](../core-concepts/dependency-injection.md) - DI concepts and patterns
- [**Module System**](../core-concepts/module-system.md) - Module architecture
- [**Testing Guide**](../development/testing.md) - Testing with DI

---

## External Resources

- [Martin Fowler - Inversion of Control](https://martinfowler.com/bliki/InversionOfControl.html)
- [Dependency Injection Principles](https://en.wikipedia.org/wiki/Dependency_injection)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)

---

## Source Code

View the implementation:
- [GitHub: container.go](https://github.com/neonextechnologies/neonexcore/blob/main/internal/core/container.go)
- [GitHub: container_test.go](https://github.com/neonextechnologies/neonexcore/blob/main/internal/core/container_test.go)

---

**Need help?** Check our [FAQ](../resources/faq.md) or [get support](../resources/support.md).
