# Architecture

Learn about the architectural design and principles behind Neonex Core.

---

## Overview

Neonex Core implements a **Modular Monolith** architecture, combining the simplicity of monolithic applications with the modularity of microservices.

```
┌───────────────────────────────────────────────────────────┐
│                    Application Layer                       │
│                  (Your Business Logic)                     │
├───────────────────────────────────────────────────────────┤
│              Module Layer (Auto-Discovery)                 │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐    │
│  │ User    │  │ Product │  │ Order   │  │ Payment │    │
│  │ Module  │  │ Module  │  │ Module  │  │ Module  │    │
│  └─────────┘  └─────────┘  └─────────┘  └─────────┘    │
├───────────────────────────────────────────────────────────┤
│                     Core Framework                         │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐         │
│  │ Container  │  │ Registry   │  │ Lifecycle  │         │
│  │ (DI)       │  │ (Modules)  │  │ Management │         │
│  └────────────┘  └────────────┘  └────────────┘         │
├───────────────────────────────────────────────────────────┤
│                  Infrastructure Layer                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌─────────┐ │
│  │ HTTP     │  │ Database │  │ Logger   │  │ Cache   │ │
│  │ (Fiber)  │  │ (GORM)   │  │ (Zap)    │  │ (Redis) │ │
│  └──────────┘  └──────────┘  └──────────┘  └─────────┘ │
└───────────────────────────────────────────────────────────┘
```

---

## Design Principles

### 1. Modular First

**Principle:** Every feature is a self-contained module.

**Benefits:**
- Clear boundaries between features
- Independent development and testing
- Easy to add, remove, or replace modules
- Reduced cognitive load

**Implementation:**
```go
// Each module is independent
modules/
├── user/          # User management
├── product/       # Product catalog
├── order/         # Order processing
└── payment/       # Payment handling
```

### 2. Dependency Injection

**Principle:** Dependencies are injected, not created.

**Benefits:**
- Loose coupling between components
- Easy to test with mocks
- Clear dependency graph
- Flexible configuration

**Implementation:**
```go
// Register dependencies
container.Singleton("userRepo", func(c *Container) interface{} {
    return NewUserRepository(db)
})

// Inject dependencies
type UserService struct {
    repo UserRepository  // Injected, not created
}
```

### 3. Repository Pattern

**Principle:** Separate data access from business logic.

**Benefits:**
- Testable business logic
- Swappable data sources
- Consistent data access patterns
- Query optimization in one place

**Implementation:**
```go
// Repository interface
type UserRepository interface {
    FindByID(ctx context.Context, id uint) (*User, error)
    Create(ctx context.Context, user *User) error
}

// Service uses interface
type UserService struct {
    repo UserRepository  // Interface, not concrete type
}
```

### 4. Clean Architecture

**Principle:** Business logic is independent of frameworks.

**Layers:**
```
Presentation → Application → Domain → Infrastructure
   (HTTP)        (Service)   (Models)    (Database)
```

**Benefits:**
- Framework-independent business logic
- Easy to test without infrastructure
- Can swap Fiber for another HTTP framework
- Can swap GORM for another ORM

---

## Component Architecture

### Module System

```go
type Module interface {
    Name() string           // Module identifier
    Version() string        // Semantic version
    Dependencies() []string // Required modules
    Register(*Container)    // Register services
    Boot(*App)             // Initialize module
}
```

**Lifecycle:**
1. **Discovery** - Find all modules in `modules/` directory
2. **Resolution** - Resolve dependencies in correct order
3. **Registration** - Register services in DI container
4. **Boot** - Initialize module (routes, migrations, etc.)

### Dependency Injection Container

```go
type Container struct {
    services  map[string]ServiceDefinition
    instances map[string]interface{}
    mu        sync.RWMutex
}
```

**Capabilities:**
- Singleton services (shared instance)
- Transient services (new instance each time)
- Factory functions with dependency resolution
- Thread-safe concurrent access

### HTTP Layer (Fiber)

```
Request → Middleware → Router → Controller → Service → Repository → Database
                                      ↓
                                  Response
```

**Features:**
- Fast routing with zero allocations
- Middleware chain
- Context passing
- Error handling
- Request validation

### Database Layer (GORM)

```go
// Generic repository
type Repository[T any] struct {
    db *gorm.DB
}

// Type-safe operations
userRepo := NewRepository[User](db)
user, err := userRepo.FindByID(ctx, 1)
```

**Features:**
- Generic repository pattern
- Transaction management
- Soft deletes
- Eager loading
- Custom queries
- Multi-database support

### Logging System

```
Event → Logger → Encoder → Writer → Output
                              ↓
                     [Console | File | JSON]
```

**Features:**
- Structured logging (key-value pairs)
- Multiple log levels
- Multiple outputs simultaneously
- Automatic file rotation
- Contextual loggers

---

## Request Lifecycle

### HTTP Request Flow

```
1. HTTP Request arrives
   ↓
2. Global Middleware
   - CORS
   - Request ID
   - Logging
   - Rate Limiting
   ↓
3. Router matches route
   ↓
4. Route Middleware
   - Authentication
   - Authorization
   - Validation
   ↓
5. Controller receives request
   ↓
6. Service processes business logic
   ↓
7. Repository accesses database
   ↓
8. Response flows back through middleware
   ↓
9. Client receives response
```

### Example with Code

```go
// 1. Request arrives
POST /api/users

// 2. Global middleware
func RequestID(c *fiber.Ctx) error {
    c.Locals("request_id", uuid.New())
    return c.Next()
}

// 3 & 4. Route with middleware
router.Post("/users", 
    middleware.Auth(),        // Route middleware
    controller.CreateUser,    // Controller
)

// 5. Controller
func (ctrl *UserController) CreateUser(c *fiber.Ctx) error {
    var req CreateUserRequest
    c.BodyParser(&req)
    
    // 6. Service
    user, err := ctrl.service.Create(c.Context(), &req)
    if err != nil {
        return err
    }
    
    return c.JSON(user)
}

// 6. Service (business logic)
func (s *UserService) Create(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // Validate
    if err := s.validator.Validate(req); err != nil {
        return nil, err
    }
    
    // Hash password
    hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
    
    user := &User{
        Name:     req.Name,
        Email:    req.Email,
        Password: string(hash),
    }
    
    // 7. Repository (database access)
    if err := s.repo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    return user, nil
}
```

---

## Module Communication

### Direct Service Calls

```go
// Order module calls Product module
type OrderService struct {
    productService *ProductService  // Injected
}

func (s *OrderService) CreateOrder(items []OrderItem) error {
    for _, item := range items {
        // Direct call to another module
        product, err := s.productService.GetByID(item.ProductID)
        if err != nil {
            return err
        }
        // ... process order
    }
}
```

### Event-Based Communication

```go
// User module emits event
eventBus.Publish("user.created", UserCreatedEvent{
    UserID: user.ID,
    Email:  user.Email,
})

// Order module listens
eventBus.Subscribe("user.created", func(event UserCreatedEvent) {
    // Create welcome order, send email, etc.
})
```

---

## Scalability Patterns

### Horizontal Scaling

```
Load Balancer
    ↓
┌─────────┬─────────┬─────────┐
│ App 1   │ App 2   │ App 3   │  (Multiple instances)
└────┬────┴────┬────┴────┬────┘
     └─────────┼─────────┘
               ↓
          Database
```

**Requirements:**
- Stateless application (session in Redis/DB)
- Shared database or read replicas
- Cache invalidation strategy

### Vertical Scaling

```
Single Powerful Server
    ↓
Neonex Core
(Handles thousands of req/sec)
```

**When to use:**
- Simpler to manage
- Lower operational complexity
- Cost-effective for medium traffic
- Go's concurrency handles load well

### Migration to Microservices

When you outgrow the monolith:

```
Monolith                    Microservices
┌───────────┐              ┌────────┐
│  User     │              │ User   │
│  Product  │    →→→       │ Service│
│  Order    │              └────────┘
│  Payment  │              ┌────────┐
└───────────┘              │Product │
                           │Service │
                           └────────┘
                           ┌────────┐
                           │ Order  │
                           │Service │
                           └────────┘
```

**Migration Path:**
1. Extract module as separate service
2. Add API layer between services
3. Deploy independently
4. Gradually decompose monolith

---

## Design Patterns Used

### Creational Patterns

**1. Factory Pattern**
```go
func NewUserService(repo UserRepository, logger Logger) *UserService {
    return &UserService{
        repo:   repo,
        logger: logger,
    }
}
```

**2. Singleton Pattern**
```go
// DI container ensures single instance
container.Singleton("database", func() *gorm.DB {
    return setupDB()
})
```

### Structural Patterns

**1. Repository Pattern**
```go
type Repository[T any] interface {
    FindByID(ctx context.Context, id uint) (*T, error)
    Create(ctx context.Context, entity *T) error
}
```

**2. Adapter Pattern**
```go
// Fiber adapter for HTTP
type FiberAdapter struct {
    app *fiber.App
}

func (a *FiberAdapter) Start(port string) error {
    return a.app.Listen(":" + port)
}
```

### Behavioral Patterns

**1. Chain of Responsibility**
```go
// Middleware chain
app.Use(RequestID())
app.Use(Logger())
app.Use(Auth())
```

**2. Observer Pattern**
```go
// Event system
eventBus.Subscribe("user.created", handler1)
eventBus.Subscribe("user.created", handler2)
```

**3. Strategy Pattern**
```go
// Different cache strategies
type CacheStrategy interface {
    Get(key string) (interface{}, error)
    Set(key string, value interface{}) error
}
```

---

## Performance Considerations

### Zero-Allocation Routing

Fiber uses fasthttp which avoids allocations in hot paths:

```go
// No allocations
c.Params("id")      // Returns string view, no copy
c.Query("page")     // Returns string view, no copy
```

### Connection Pooling

```go
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### Caching Strategy

```
Request → Check Cache → Found? Return : Query DB → Cache Result → Return
```

### Goroutine Usage

```go
// Concurrent processing
var wg sync.WaitGroup
for _, item := range items {
    wg.Add(1)
    go func(item Item) {
        defer wg.Done()
        processItem(item)
    }(item)
}
wg.Wait()
```

---

## Security Architecture

### Authentication Flow

```
1. User submits credentials
   ↓
2. Validate credentials
   ↓
3. Generate JWT token
   ↓
4. Return token to client
   ↓
5. Client includes token in requests
   ↓
6. Middleware validates token
   ↓
7. Extract user info to context
```

### Authorization

```go
// Role-based access control
func RequirePermission(permission string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        user := c.Locals("user").(*User)
        if !user.HasPermission(permission) {
            return c.Status(403).JSON(fiber.Map{
                "error": "Forbidden",
            })
        }
        return c.Next()
    }
}
```

### Input Validation

```
Request → Parse Body → Validate → Sanitize → Process
```

---

## Testing Strategy

### Unit Tests

```go
func TestUserService_Create(t *testing.T) {
    // Mock dependencies
    mockRepo := &MockUserRepository{}
    mockLogger := &MockLogger{}
    
    service := NewUserService(mockRepo, mockLogger)
    
    // Test
    user, err := service.Create(ctx, req)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

### Integration Tests

```go
func TestUserAPI_Create(t *testing.T) {
    // Setup test database
    db := setupTestDB()
    
    // Create app
    app := setupTestApp(db)
    
    // Make request
    req := httptest.NewRequest("POST", "/api/users", body)
    resp, _ := app.Test(req)
    
    // Assert
    assert.Equal(t, 201, resp.StatusCode)
}
```

---

## Configuration Management

### Environment-Based

```
.env.development    → Development settings
.env.staging        → Staging settings
.env.production     → Production settings
```

### Hierarchical Config

```go
type Config struct {
    App      AppConfig
    Database DatabaseConfig
    Logger   LoggerConfig
    Cache    CacheConfig
}
```

### Type-Safe Access

```go
config := LoadConfig()
port := config.App.Port           // Type-safe
dbHost := config.Database.Host    // Type-safe
```

---

## Monitoring & Observability

### Logging

```go
logger.Info("User created", logger.Fields{
    "user_id": user.ID,
    "email":   user.Email,
    "ip":      clientIP,
})
```

### Metrics (Future)

```go
metrics.Counter("user.created").Inc()
metrics.Histogram("request.duration").Observe(duration)
```

### Tracing (Future)

```go
span := tracer.StartSpan("CreateUser")
defer span.Finish()
```

### Health Checks

```go
GET /health
{
    "status": "ok",
    "database": "connected",
    "cache": "connected"
}
```

---

## Next Steps

- [**Installation**](../getting-started/installation.md) - Install Neonex Core
- [**Quick Start**](../getting-started/quick-start.md) - Build your first app
- [**Module System**](../core-concepts/module-system.md) - Learn about modules
- [**Dependency Injection**](../core-concepts/dependency-injection.md) - Master DI

---

**Questions?** Check the [FAQ](../resources/faq.md) or [get support](../resources/support.md).
