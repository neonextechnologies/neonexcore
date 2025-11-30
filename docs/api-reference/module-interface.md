# Module Interface API

Complete API reference for implementing modules in Neonex Core.

---

## Overview

The Module Interface defines the contract that all Neonex Core modules must implement. It provides lifecycle hooks, dependency registration, and route configuration methods.

**Package:** `neonexcore/internal/core`

**Key Concepts:**
- 🎯 **Convention over Configuration** - Standard interface for all modules
- 🔌 **Auto-Discovery** - Modules automatically found and loaded
- 🔄 **Lifecycle Management** - Predictable initialization flow
- 📦 **Self-Contained** - Each module manages its own dependencies
- 🌐 **Route Registration** - Modules define their own HTTP endpoints

---

## Module Interface

### Definition

```go
type Module interface {
    Name() string
    Init()
    Routes(*fiber.App, *Container)
    RegisterServices(*Container)
}
```

The core interface that every module must implement.

---

## Interface Methods

### Name

```go
Name() string
```

Returns the unique identifier for the module.

**Returns:**
- `string` - Module name (lowercase, no spaces)

**Rules:**
- Must be unique across all modules
- Should be lowercase
- No special characters except hyphen/underscore
- Should match directory name

**Examples:**

```go
func (m *UserModule) Name() string {
    return "user"
}

func (m *ProductModule) Name() string {
    return "product"
}

func (m *OrderManagementModule) Name() string {
    return "order-management"
}
```

**Used For:**
- Module identification in logs
- Dependency resolution
- Route grouping
- Configuration lookup

---

### Init

```go
Init()
```

Called once during application initialization, before services are registered.

**Purpose:**
- Initialize module-level configuration
- Set up module state
- Validate module requirements
- Load module metadata

**Timing:**
- Called after module instantiation
- Called before `RegisterServices()`
- Called before `Routes()`
- Only called once per application lifecycle

**Examples:**

**Simple Initialization:**
```go
func (m *UserModule) Init() {
    fmt.Printf("User module initialized\n")
}
```

**Load Configuration:**
```go
func (m *ProductModule) Init() {
    config, err := LoadConfigFromFile("modules/product/config.json")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    m.config = config
    
    log.Printf("Product module initialized with config: %+v", config)
}
```

**Validate Requirements:**
```go
func (m *PaymentModule) Init() {
    if os.Getenv("STRIPE_API_KEY") == "" {
        log.Fatal("Payment module requires STRIPE_API_KEY")
    }
    
    if os.Getenv("STRIPE_WEBHOOK_SECRET") == "" {
        log.Fatal("Payment module requires STRIPE_WEBHOOK_SECRET")
    }
    
    log.Println("Payment module initialized successfully")
}
```

**Set Up Module State:**
```go
func (m *CacheModule) Init() {
    m.cache = make(map[string]interface{})
    m.mu = &sync.RWMutex{}
    m.ttl = 5 * time.Minute
    
    log.Println("Cache module initialized")
}
```

**Start Background Workers:**
```go
func (m *NotificationModule) Init() {
    m.queue = make(chan Notification, 100)
    
    // Start background worker
    go m.processNotifications()
    
    log.Println("Notification worker started")
}

func (m *NotificationModule) processNotifications() {
    for notif := range m.queue {
        m.send(notif)
    }
}
```

---

### RegisterServices

```go
RegisterServices(*Container)
```

Registers module dependencies in the DI container.

**Parameters:**
- `*Container` - The application's dependency injection container

**Purpose:**
- Register repositories
- Register services
- Register controllers
- Register validators
- Register any module-specific components

**Timing:**
- Called after `Init()`
- Called before `Routes()`
- Called once during startup

**Examples:**

**Basic Registration:**
```go
func (m *UserModule) RegisterServices(c *core.Container) {
    // Repository
    c.Provide(func() Repository {
        db := core.Resolve[*gorm.DB](c)
        return NewRepository(db)
    }, core.Singleton)
    
    // Service
    c.Provide(func() Service {
        repo := core.Resolve[Repository](c)
        return NewService(repo)
    }, core.Singleton)
    
    // Controller
    c.Provide(func() *Controller {
        service := core.Resolve[Service](c)
        return NewController(service)
    }, core.Singleton)
}
```

**With Logger:**
```go
func (m *ProductModule) RegisterServices(c *core.Container) {
    c.Provide(func() Repository {
        db := core.Resolve[*gorm.DB](c)
        return NewRepository(db)
    }, core.Singleton)
    
    c.Provide(func() Service {
        repo := core.Resolve[Repository](c)
        logger := core.Resolve[logger.Logger](c)
        return NewService(repo, logger)
    }, core.Singleton)
    
    c.Provide(func() *Controller {
        service := core.Resolve[Service](c)
        logger := core.Resolve[logger.Logger](c)
        return NewController(service, logger)
    }, core.Singleton)
}
```

**With Multiple Dependencies:**
```go
func (m *OrderModule) RegisterServices(c *core.Container) {
    // Repository
    c.Provide(func() Repository {
        db := core.Resolve[*gorm.DB](c)
        return NewRepository(db)
    }, core.Singleton)
    
    // Service with multiple dependencies
    c.Provide(func() Service {
        repo := core.Resolve[Repository](c)
        logger := core.Resolve[logger.Logger](c)
        cache := core.Resolve[CacheService](c)
        config := core.Resolve[*Config](c)
        return NewService(repo, logger, cache, config)
    }, core.Singleton)
    
    // Controller
    c.Provide(func() *Controller {
        service := core.Resolve[Service](c)
        return NewController(service)
    }, core.Singleton)
}
```

**Delegate to Separate Function:**
```go
func (m *UserModule) RegisterServices(c *core.Container) {
    // Delegate to dedicated DI file
    RegisterDependencies(c)
}

// In di.go
func RegisterDependencies(c *core.Container) {
    registerRepository(c)
    registerService(c)
    registerController(c)
    registerValidators(c)
}
```

---

### Routes

```go
Routes(*fiber.App, *Container)
```

Registers HTTP routes for the module.

**Parameters:**
- `*fiber.App` - Fiber application instance
- `*Container` - DI container to resolve controllers

**Purpose:**
- Define HTTP endpoints
- Map routes to controllers
- Add route-specific middleware
- Group related endpoints

**Timing:**
- Called after `RegisterServices()`
- Called before server starts
- Called once during startup

**Examples:**

**Basic Routes:**
```go
func (m *UserModule) Routes(app *fiber.App, c *core.Container) {
    ctrl := core.Resolve[*Controller](c)
    
    api := app.Group("/api/users")
    api.Get("/", ctrl.GetAll)
    api.Get("/:id", ctrl.GetByID)
    api.Post("/", ctrl.Create)
    api.Put("/:id", ctrl.Update)
    api.Delete("/:id", ctrl.Delete)
}
```

**With Middleware:**
```go
func (m *ProductModule) Routes(app *fiber.App, c *core.Container) {
    ctrl := core.Resolve[*Controller](c)
    
    api := app.Group("/api/products")
    
    // Public routes
    api.Get("/", ctrl.GetAll)
    api.Get("/:id", ctrl.GetByID)
    
    // Protected routes
    admin := api.Group("/", middleware.Auth(), middleware.RequireRole("admin"))
    admin.Post("/", ctrl.Create)
    admin.Put("/:id", ctrl.Update)
    admin.Delete("/:id", ctrl.Delete)
}
```

**Multiple Controllers:**
```go
func (m *OrderModule) Routes(app *fiber.App, c *core.Container) {
    orderCtrl := core.Resolve[*OrderController](c)
    itemCtrl := core.Resolve[*OrderItemController](c)
    
    api := app.Group("/api")
    
    // Order routes
    orders := api.Group("/orders")
    orders.Get("/", orderCtrl.GetAll)
    orders.Get("/:id", orderCtrl.GetByID)
    orders.Post("/", orderCtrl.Create)
    
    // Order items routes
    items := api.Group("/orders/:orderId/items")
    items.Get("/", itemCtrl.GetAll)
    items.Post("/", itemCtrl.Add)
}
```

**Versioned API:**
```go
func (m *UserModule) Routes(app *fiber.App, c *core.Container) {
    v1Ctrl := core.Resolve[*V1Controller](c)
    v2Ctrl := core.Resolve[*V2Controller](c)
    
    // API v1
    v1 := app.Group("/api/v1/users")
    v1.Get("/", v1Ctrl.GetAll)
    v1.Post("/", v1Ctrl.Create)
    
    // API v2
    v2 := app.Group("/api/v2/users")
    v2.Get("/", v2Ctrl.GetAll)
    v2.Post("/", v2Ctrl.Create)
}
```

**Delegate to Separate Function:**
```go
func (m *UserModule) Routes(app *fiber.App, c *core.Container) {
    ctrl := core.Resolve[*Controller](c)
    RegisterRoutes(app, ctrl)
}

// In routes.go
func RegisterRoutes(app *fiber.App, ctrl *Controller) {
    api := app.Group("/api/users")
    api.Get("/", ctrl.GetAll)
    api.Get("/:id", ctrl.GetByID)
    api.Post("/", ctrl.Create)
    api.Put("/:id", ctrl.Update)
    api.Delete("/:id", ctrl.Delete)
}
```

---

## Complete Implementation Examples

### Minimal Module

```go
// modules/health/health.go
package health

import (
    "github.com/gofiber/fiber/v2"
    "neonexcore/internal/core"
)

type HealthModule struct{}

func New() *HealthModule {
    return &HealthModule{}
}

func (m *HealthModule) Name() string {
    return "health"
}

func (m *HealthModule) Init() {
    fmt.Println("Health module initialized")
}

func (m *HealthModule) RegisterServices(c *core.Container) {
    // No services needed
}

func (m *HealthModule) Routes(app *fiber.App, c *core.Container) {
    app.Get("/health", func(ctx *fiber.Ctx) error {
        return ctx.JSON(fiber.Map{
            "status":  "healthy",
            "module": "health",
        })
    })
}
```

### Standard CRUD Module

```go
// modules/product/product.go
package product

import (
    "github.com/gofiber/fiber/v2"
    "neonexcore/internal/core"
)

type ProductModule struct {
    config *Config
}

func New() *ProductModule {
    return &ProductModule{}
}

func (m *ProductModule) Name() string {
    return "product"
}

func (m *ProductModule) Init() {
    // Load configuration
    m.config = LoadConfig()
    fmt.Printf("Product module initialized with config: %+v\n", m.config)
}

func (m *ProductModule) RegisterServices(c *core.Container) {
    // Repository
    c.Provide(func() Repository {
        db := core.Resolve[*gorm.DB](c)
        return NewRepository(db)
    }, core.Singleton)
    
    // Service
    c.Provide(func() Service {
        repo := core.Resolve[Repository](c)
        logger := core.Resolve[logger.Logger](c)
        return NewService(repo, logger)
    }, core.Singleton)
    
    // Controller
    c.Provide(func() *Controller {
        service := core.Resolve[Service](c)
        return NewController(service)
    }, core.Singleton)
}

func (m *ProductModule) Routes(app *fiber.App, c *core.Container) {
    ctrl := core.Resolve[*Controller](c)
    
    api := app.Group("/api/products")
    api.Get("/", ctrl.GetAll)
    api.Get("/:id", ctrl.GetByID)
    api.Post("/", ctrl.Create)
    api.Put("/:id", ctrl.Update)
    api.Delete("/:id", ctrl.Delete)
}
```

### Complex Module with Features

```go
// modules/analytics/analytics.go
package analytics

import (
    "github.com/gofiber/fiber/v2"
    "neonexcore/internal/core"
    "time"
)

type AnalyticsModule struct {
    config  *Config
    enabled bool
    queue   chan Event
}

func New() *AnalyticsModule {
    return &AnalyticsModule{
        queue: make(chan Event, 1000),
    }
}

func (m *AnalyticsModule) Name() string {
    return "analytics"
}

func (m *AnalyticsModule) Init() {
    // Load configuration
    m.config = LoadConfig()
    m.enabled = m.config.Enabled
    
    if !m.enabled {
        fmt.Println("Analytics module disabled")
        return
    }
    
    // Start background processor
    go m.processEvents()
    
    fmt.Println("Analytics module initialized and running")
}

func (m *AnalyticsModule) RegisterServices(c *core.Container) {
    if !m.enabled {
        return
    }
    
    // Repository
    c.Provide(func() Repository {
        db := core.Resolve[*gorm.DB](c)
        return NewRepository(db)
    }, core.Singleton)
    
    // Tracker service
    c.Provide(func() TrackerService {
        repo := core.Resolve[Repository](c)
        return NewTrackerService(repo, m.queue)
    }, core.Singleton)
    
    // Analytics service
    c.Provide(func() AnalyticsService {
        repo := core.Resolve[Repository](c)
        return NewAnalyticsService(repo)
    }, core.Singleton)
    
    // Controller
    c.Provide(func() *Controller {
        service := core.Resolve[AnalyticsService](c)
        return NewController(service)
    }, core.Singleton)
}

func (m *AnalyticsModule) Routes(app *fiber.App, c *core.Container) {
    if !m.enabled {
        return
    }
    
    ctrl := core.Resolve[*Controller](c)
    
    api := app.Group("/api/analytics")
    api.Get("/summary", ctrl.GetSummary)
    api.Get("/events", ctrl.GetEvents)
    api.Get("/users/:id", ctrl.GetUserAnalytics)
}

func (m *AnalyticsModule) processEvents() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    batch := make([]Event, 0, 100)
    
    for {
        select {
        case event := <-m.queue:
            batch = append(batch, event)
            if len(batch) >= 100 {
                m.processBatch(batch)
                batch = batch[:0]
            }
        case <-ticker.C:
            if len(batch) > 0 {
                m.processBatch(batch)
                batch = batch[:0]
            }
        }
    }
}
```

---

## Module Lifecycle

### Initialization Sequence

```
1. Application starts
   ↓
2. Module instances created (New() called)
   ↓
3. Module.Init() called for each module
   ↓
4. Module.RegisterServices() called for each module
   ↓
5. Module.Routes() called for each module
   ↓
6. Server starts
   ↓
7. Modules handle requests (Runtime)
```

### Example with Logging

```go
func (m *UserModule) Init() {
    log.Println("[1/4] User module: Init")
    m.config = LoadConfig()
}

func (m *UserModule) RegisterServices(c *core.Container) {
    log.Println("[2/4] User module: RegisterServices")
    c.Provide(func() Service {
        log.Println("  → Creating UserService")
        return NewService()
    }, core.Singleton)
}

func (m *UserModule) Routes(app *fiber.App, c *core.Container) {
    log.Println("[3/4] User module: Routes")
    ctrl := core.Resolve[*Controller](c)
    app.Get("/users", ctrl.GetAll)
}

// Output:
// [1/4] User module: Init
// [2/4] User module: RegisterServices
// [3/4] User module: Routes
//   → Creating UserService
// [4/4] Server ready at :8080
```

---

## Best Practices

### ✅ DO: Keep Init Lightweight

```go
// Good: Quick initialization
func (m *Module) Init() {
    m.config = LoadConfig()
    fmt.Println("Module initialized")
}

// Bad: Blocking operations
func (m *Module) Init() {
    time.Sleep(30 * time.Second)  // ❌ Blocks startup
    m.setupHeavyResource()         // ❌ Delays application
}
```

### ✅ DO: Register All Dependencies

```go
// Good: Complete registration
func (m *Module) RegisterServices(c *core.Container) {
    c.Provide(func() Repository { ... }, core.Singleton)
    c.Provide(func() Service { ... }, core.Singleton)
    c.Provide(func() *Controller { ... }, core.Singleton)
}
```

### ✅ DO: Group Related Routes

```go
// Good: Organized routes
func (m *Module) Routes(app *fiber.App, c *core.Container) {
    api := app.Group("/api/users")
    
    // Public
    api.Get("/", ctrl.GetAll)
    
    // Protected
    protected := api.Use(middleware.Auth())
    protected.Post("/", ctrl.Create)
}
```

### ✅ DO: Handle Optional Features

```go
// Good: Conditional features
func (m *Module) Init() {
    m.enabled = os.Getenv("FEATURE_ENABLED") == "true"
}

func (m *Module) RegisterServices(c *core.Container) {
    if !m.enabled {
        return  // Skip registration
    }
    // Register services...
}
```

### ❌ DON'T: Panic in Init

```go
// Bad: Panic fails entire app
func (m *Module) Init() {
    if err := m.setup(); err != nil {
        panic(err)  // ❌ Crashes app
    }
}

// Good: Log and continue or fail gracefully
func (m *Module) Init() {
    if err := m.setup(); err != nil {
        log.Printf("Warning: Module setup failed: %v", err)
        m.enabled = false
    }
}
```

### ❌ DON'T: Resolve in RegisterServices

```go
// Bad: Early resolution
func (m *Module) RegisterServices(c *core.Container) {
    db := core.Resolve[*gorm.DB](c)  // ❌ May not be ready
    
    c.Provide(func() Service {
        return NewService(db)
    }, core.Singleton)
}

// Good: Lazy resolution
func (m *Module) RegisterServices(c *core.Container) {
    c.Provide(func() Service {
        db := core.Resolve[*gorm.DB](c)  // ✅ Resolved when needed
        return NewService(db)
    }, core.Singleton)
}
```

---

## Module Registration

### Auto-Discovery

Modules in `modules/` directory are automatically discovered:

```
modules/
├── user/
│   └── user.go      # Must have New() function
├── product/
│   └── product.go   # Must have New() function
└── order/
    └── order.go     # Must have New() function
```

Each must export a `New()` function:

```go
func New() *UserModule {
    return &UserModule{}
}
```

### Manual Registration

Register modules explicitly in `main.go`:

```go
func main() {
    app := core.NewApp()
    
    // Register modules manually
    app.Registry.Register(user.New())
    app.Registry.Register(product.New())
    app.Registry.Register(order.New())
    
    app.Run()
}
```

---

## Testing Modules

```go
func TestUserModule(t *testing.T) {
    // Create module
    module := New()
    
    // Test Name
    assert.Equal(t, "user", module.Name())
    
    // Test Init
    module.Init()
    assert.NotNil(t, module.config)
    
    // Test RegisterServices
    container := core.NewContainer()
    module.RegisterServices(container)
    
    // Verify service registered
    service := core.Resolve[Service](container)
    assert.NotNil(t, service)
    
    // Test Routes
    app := fiber.New()
    module.Routes(app, container)
    
    // Test route exists
    req := httptest.NewRequest("GET", "/api/users", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, 200, resp.StatusCode)
}
```

---

## Related Documentation

- [**Container API**](container.md) - Dependency injection
- [**Core API**](core.md) - Core framework
- [**Module System Guide**](../core-concepts/module-system.md) - Module concepts
- [**Creating Modules**](../advanced/custom-modules.md) - Build custom modules

---

## Source Code

- [GitHub: module.go](https://github.com/neonextechnologies/neonexcore/blob/main/internal/core/module.go)
- [GitHub: Example Modules](https://github.com/neonextechnologies/neonexcore/tree/main/modules)

---

**Need help?** Check our [FAQ](../resources/faq.md) or [get support](../resources/support.md).
