# Application Lifecycle

Understanding how Neonex Core applications initialize, run, and shutdown.

---

## Overview

The application lifecycle defines the sequence of events from startup to shutdown. Understanding this flow helps you:

- ✅ Initialize resources at the right time
- ✅ Handle dependencies properly
- ✅ Implement graceful shutdown
- ✅ Debug initialization issues
- ✅ Optimize startup performance

---

## Lifecycle Phases

### Complete Flow

```
1. Application Start
   ↓
2. Load Environment (.env)
   ↓
3. Initialize Logger
   ↓
4. Connect to Database
   ↓
5. Discover Modules
   ↓
6. Register Dependencies (DI)
   ↓
7. Auto-Migrate Models
   ↓
8. Run Seeders
   ↓
9. Register Routes
   ↓
10. Start HTTP Server
   ↓
11. Handle Requests (Runtime)
   ↓
12. Graceful Shutdown
```

---

## Phase-by-Phase Breakdown

### Phase 1: Application Start

```go
// main.go
package main

import (
    "neonexcore/internal/core"
)

func main() {
    // Create application instance
    app := core.NewApp()
    
    // Start application
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

**What happens:**
- Application instance created
- Module registry initialized
- DI container created

### Phase 2: Load Environment

```go
func NewApp() *App {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
    
    return &App{
        Registry:  NewModuleRegistry(),
        Container: NewContainer(),
    }
}
```

**Environment variables loaded:**
- Database configuration
- Server settings
- Logger configuration
- Module settings

### Phase 3: Initialize Logger

```go
func (a *App) InitLogger(cfg logger.Config) error {
    if err := logger.Setup(cfg); err != nil {
        return err
    }
    
    a.Logger = logger.NewLogger()
    a.Logger.Info("Logger initialized", logger.Fields{
        "level":  cfg.Level,
        "format": cfg.Format,
    })
    
    return nil
}
```

**Logger setup:**
- Create log file (if configured)
- Set log level
- Initialize formatters
- Configure output destinations

### Phase 4: Connect to Database

```go
func (a *App) InitDatabase() error {
    dbConfig := config.LoadDatabaseConfig()
    
    db, err := config.InitDatabase(dbConfig)
    if err != nil {
        return fmt.Errorf("database connection failed: %w", err)
    }
    
    a.Migrator = database.NewMigrator(db)
    a.Logger.Info("Database connected", logger.Fields{
        "driver": dbConfig.Driver,
    })
    
    return nil
}
```

**Database connection:**
- Read database config
- Open connection
- Test connection
- Create migrator instance

### Phase 5: Discover Modules

```go
func (r *ModuleRegistry) DiscoverModules() {
    entries, err := os.ReadDir("./modules")
    if err != nil {
        return
    }
    
    for _, entry := range entries {
        if entry.IsDir() {
            moduleName := entry.Name()
            
            // Load module (by convention)
            // Expects: modules/{name}/{name}.go with New() function
            module := loadModule(moduleName)
            if module != nil {
                r.Register(module)
            }
        }
    }
}
```

**Module discovery:**
- Scan `modules/` directory
- Find module packages
- Load module instances
- Register in registry

### Phase 6: Register Dependencies

```go
func (a *App) RegisterServices() {
    // Register core services
    a.Container.Provide(func() *gorm.DB {
        return config.DB.GetDB()
    }, core.Singleton)
    
    a.Container.Provide(func() logger.Logger {
        return a.Logger
    }, core.Singleton)
    
    // Register module services
    a.Registry.RegisterModuleServices(a.Container)
}
```

**Dependency registration:**
- Core services (DB, Logger)
- Module repositories
- Module services
- Module controllers

### Phase 7: Auto-Migrate Models

```go
func (a *App) AutoMigrate() error {
    // Collect models from all modules
    models := a.Registry.CollectModels()
    
    // Register with migrator
    a.Migrator.RegisterModels(models...)
    
    // Run migration
    a.Logger.Info("Running auto-migration...")
    if err := a.Migrator.AutoMigrate(); err != nil {
        return err
    }
    
    a.Logger.Info("Migration completed", logger.Fields{
        "count": len(models),
    })
    
    return nil
}
```

**Migration:**
- Collect all module models
- Generate migration DDL
- Apply to database
- Log results

### Phase 8: Run Seeders

```go
func (a *App) RunSeeders() error {
    seeders := a.Registry.CollectSeeders()
    
    for _, seeder := range seeders {
        a.Logger.Info("Running seeder", logger.Fields{
            "name": reflect.TypeOf(seeder).String(),
        })
        
        if err := seeder.Seed(); err != nil {
            a.Logger.Error("Seeder failed", logger.Fields{
                "error": err.Error(),
            })
            return err
        }
    }
    
    return nil
}
```

**Seeding:**
- Collect module seeders
- Run each seeder
- Populate initial data
- Log progress

### Phase 9: Register Routes

```go
func (a *App) RegisterRoutes(app *fiber.App) {
    // Default routes
    app.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "framework": "Neonex Core",
            "version":   "0.1-alpha",
        })
    })
    
    // Module routes
    a.Registry.LoadRoutes(app, a.Container)
}
```

**Route registration:**
- Set up default routes
- Load module routes
- Configure middleware
- Mount route groups

### Phase 10: Start HTTP Server

```go
func (a *App) StartHTTP() {
    app := fiber.New(fiber.Config{
        AppName: "Neonex Core",
    })
    
    // Add middleware
    app.Use(logger.RequestIDMiddleware(a.Logger))
    app.Use(logger.HTTPMiddleware(a.Logger))
    
    // Register routes
    a.RegisterRoutes(app)
    
    // Start server
    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8080"
    }
    
    a.Logger.Info("Starting HTTP server", logger.Fields{
        "port": port,
    })
    
    if err := app.Listen(":" + port); err != nil {
        a.Logger.Fatal("Server failed", logger.Fields{
            "error": err.Error(),
        })
    }
}
```

**Server startup:**
- Create Fiber instance
- Add middleware
- Register routes
- Listen on port

### Phase 11: Runtime

```
Request → Middleware → Router → Controller → Service → Repository → DB
                                                                    ↓
Response ← Controller ← Service ← Repository ← DB Result
```

**Request handling:**
- Receive HTTP request
- Process through middleware
- Route to controller
- Execute business logic
- Return response

### Phase 12: Graceful Shutdown

```go
func (a *App) Shutdown() error {
    a.Logger.Info("Shutting down gracefully...")
    
    // Close database connections
    if db := config.DB.GetDB(); db != nil {
        sqlDB, _ := db.DB()
        if err := sqlDB.Close(); err != nil {
            a.Logger.Error("Database close failed", logger.Fields{
                "error": err.Error(),
            })
        }
    }
    
    // Close logger files
    if closer, ok := a.Logger.(io.Closer); ok {
        closer.Close()
    }
    
    a.Logger.Info("Shutdown complete")
    return nil
}
```

**Shutdown process:**
- Stop accepting requests
- Finish pending requests
- Close database connections
- Flush log buffers
- Clean up resources

---

## Module Lifecycle

### Module Initialization Flow

```go
type Module interface {
    Name() string
    Init()                                    // Called during initialization
    RegisterServices(*Container)              // Register dependencies
    Routes(*fiber.App, *Container)           // Register routes
}
```

**Module lifecycle:**

1. **Instantiation** - `module.New()` called
2. **Registration** - Added to registry
3. **Initialization** - `Init()` called
4. **Service Registration** - `RegisterServices()` called
5. **Route Registration** - `Routes()` called
6. **Runtime** - Handle requests

### Example Module Lifecycle

```go
// modules/user/user.go
type UserModule struct {
    config *Config
}

// 1. Instantiation
func New() *UserModule {
    return &UserModule{}
}

// 2. Initialization
func (m *UserModule) Init() {
    m.config = LoadConfig()
    fmt.Println("User module initialized")
}

// 3. Service Registration
func (m *UserModule) RegisterServices(c *core.Container) {
    RegisterDependencies(c)
}

// 4. Route Registration
func (m *UserModule) Routes(app *fiber.App, c *core.Container) {
    controller := core.Resolve[*Controller](c)
    RegisterRoutes(app, controller)
}
```

---

## Hooks and Events

### Custom Lifecycle Hooks

```go
type Module interface {
    Name() string
    Init()
    BeforeStart()  // Called before server starts
    AfterStart()   // Called after server starts
    OnShutdown()   // Called during shutdown
    RegisterServices(*Container)
    Routes(*fiber.App, *Container)
}
```

### Implementation

```go
func (m *UserModule) BeforeStart() {
    // Warm up caches
    m.cache.Preload()
    
    // Check external services
    if err := m.checkDependencies(); err != nil {
        log.Fatal(err)
    }
}

func (m *UserModule) AfterStart() {
    // Start background jobs
    go m.startCleanupJob()
    
    // Send startup notification
    m.notifyStartup()
}

func (m *UserModule) OnShutdown() {
    // Stop background jobs
    m.stopJobs()
    
    // Save state
    m.saveState()
}
```

---

## Error Handling During Startup

### Fail Fast Strategy

```go
func (a *App) Run() error {
    // Critical: Database must connect
    if err := a.InitDatabase(); err != nil {
        return fmt.Errorf("fatal: %w", err)
    }
    
    // Critical: Migrations must succeed
    if err := a.AutoMigrate(); err != nil {
        return fmt.Errorf("fatal: %w", err)
    }
    
    // Non-critical: Seeders can fail
    if err := a.RunSeeders(); err != nil {
        a.Logger.Warn("Seeding failed, continuing...", logger.Fields{
            "error": err.Error(),
        })
    }
    
    // Start server
    a.StartHTTP()
    return nil
}
```

### Retry Strategy

```go
func (a *App) InitDatabaseWithRetry() error {
    maxRetries := 5
    retryDelay := 2 * time.Second
    
    for i := 0; i < maxRetries; i++ {
        err := a.InitDatabase()
        if err == nil {
            return nil
        }
        
        a.Logger.Warn("Database connection failed, retrying...", logger.Fields{
            "attempt": i + 1,
            "max":     maxRetries,
        })
        
        time.Sleep(retryDelay)
    }
    
    return errors.New("database connection failed after retries")
}
```

---

## Monitoring Lifecycle

### Startup Metrics

```go
func (a *App) Run() error {
    startTime := time.Now()
    
    a.Logger.Info("Application starting...")
    
    if err := a.InitDatabase(); err != nil {
        return err
    }
    
    moduleCount := len(a.Registry.GetModules())
    
    if err := a.StartHTTP(); err != nil {
        return err
    }
    
    duration := time.Since(startTime)
    
    a.Logger.Info("Application started", logger.Fields{
        "startup_time": duration.String(),
        "modules":      moduleCount,
    })
    
    return nil
}
```

### Health Checks

```go
app.Get("/health", func(c *fiber.Ctx) error {
    // Check database
    db := config.DB.GetDB()
    sqlDB, _ := db.DB()
    if err := sqlDB.Ping(); err != nil {
        return c.Status(503).JSON(fiber.Map{
            "status":   "unhealthy",
            "database": "down",
        })
    }
    
    return c.JSON(fiber.Map{
        "status":   "healthy",
        "uptime":   time.Since(startTime).String(),
        "database": "up",
    })
})
```

---

## Best Practices

### ✅ DO:

**1. Initialize in Correct Order**
```go
// Good: Logger before database
app.InitLogger()
app.InitDatabase()  // Can log errors
```

**2. Handle Errors Gracefully**
```go
// Good: Check errors at each step
if err := app.InitDatabase(); err != nil {
    log.Fatalf("Database initialization failed: %v", err)
}
```

**3. Log Lifecycle Events**
```go
// Good: Log important milestones
app.Logger.Info("Starting module initialization...")
```

### ❌ DON'T:

**1. Ignore Startup Errors**
```go
// Bad: Silent failure
app.InitDatabase()  // What if it fails?
```

**2. Block Startup**
```go
// Bad: Long-running operation during startup
func (m *Module) Init() {
    time.Sleep(30 * time.Second)  // ❌ Blocks startup
}
```

---

## Next Steps

- [**Module System**](module-system.md) - Module initialization
- [**Dependency Injection**](dependency-injection.md) - Service registration
- [**Database**](../database/overview.md) - Database initialization
- [**Logging**](../logging/overview.md) - Logger setup

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
