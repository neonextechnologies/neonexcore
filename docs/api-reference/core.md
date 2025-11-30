# Core API

Complete API reference for Neonex Core framework components.

---

## Overview

The Core API provides the foundation of the Neonex Core framework, including application initialization, module management, and core utilities.

**Package:** `neonexcore/internal/core`

**Components:**
- 🚀 **App** - Application lifecycle management
- 📦 **Container** - Dependency injection
- 🧩 **Module** - Module interface and registry
- 🗄️ **Database** - Database connection and management
- 📝 **Logger** - Logging system

---

## App

### Definition

```go
type App struct {
    Registry  *ModuleRegistry
    Container *Container
    Config    *Config
    Logger    logger.Logger
    Migrator  *database.Migrator
}
```

Main application instance that orchestrates framework components.

---

### NewApp

```go
func NewApp() *App
```

Creates and initializes a new application instance.

**Returns:**
- `*App` - Configured application

**Example:**
```go
func main() {
    app := core.NewApp()
    
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

**Initialization Steps:**
1. Loads environment variables from `.env`
2. Creates module registry
3. Initializes DI container
4. Sets up logger
5. Prepares database connection

---

### Run

```go
func (a *App) Run() error
```

Starts the application with full lifecycle initialization.

**Returns:**
- `error` - Error if startup fails

**Example:**
```go
app := core.NewApp()

if err := app.Run(); err != nil {
    log.Fatalf("Application failed to start: %v", err)
}
```

**Execution Flow:**
```
1. Load configuration
2. Initialize logger
3. Connect to database
4. Discover modules
5. Register module services
6. Run migrations
7. Run seeders
8. Register routes
9. Start HTTP server
```

---

### StartHTTP

```go
func (a *App) StartHTTP() error
```

Starts the HTTP server.

**Returns:**
- `error` - Error if server fails to start

**Example:**
```go
app := core.NewApp()

// Custom initialization
app.InitDatabase()
app.RegisterModules()

// Start server
port := os.Getenv("APP_PORT")
if port == "" {
    port = "8080"
}

log.Printf("Starting server on port %s", port)
if err := app.StartHTTP(); err != nil {
    log.Fatal(err)
}
```

**Server Configuration:**
```go
fiber.Config{
    AppName:      "Neonex Core",
    ServerHeader: "Neonex",
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
}
```

---

### InitDatabase

```go
func (a *App) InitDatabase() error
```

Initializes database connection.

**Returns:**
- `error` - Error if connection fails

**Example:**
```go
app := core.NewApp()

if err := app.InitDatabase(); err != nil {
    log.Fatalf("Database connection failed: %v", err)
}

log.Println("Database connected successfully")
```

**Configuration:**
Reads from environment variables:
```bash
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=myapp
DB_USERNAME=postgres
DB_PASSWORD=secret
```

---

### RegisterModels

```go
func (a *App) RegisterModels(models ...interface{})
```

Registers models for auto-migration.

**Parameters:**
- `models` - Model structs to register

**Example:**
```go
import (
    "myapp/modules/user"
    "myapp/modules/product"
)

app := core.NewApp()

// Register models
app.RegisterModels(
    &user.User{},
    &user.Profile{},
    &product.Product{},
    &product.Category{},
)

// Models will be migrated when app starts
```

---

### AutoMigrate

```go
func (a *App) AutoMigrate() error
```

Runs database migrations for registered models.

**Returns:**
- `error` - Error if migration fails

**Example:**
```go
app := core.NewApp()
app.InitDatabase()

// Register models
app.RegisterModels(&User{}, &Product{})

// Run migrations
if err := app.AutoMigrate(); err != nil {
    log.Fatalf("Migration failed: %v", err)
}

log.Println("Migration completed successfully")
```

**Behavior:**
- Creates tables if they don't exist
- Adds new columns
- Creates indexes
- Does NOT delete columns or tables
- Safe to run multiple times

---

### Shutdown

```go
func (a *App) Shutdown() error
```

Gracefully shuts down the application.

**Returns:**
- `error` - Error if shutdown fails

**Example:**
```go
app := core.NewApp()

// Setup signal handling
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt, syscall.SIGTERM)

go func() {
    <-c
    log.Println("Shutting down gracefully...")
    
    if err := app.Shutdown(); err != nil {
        log.Printf("Shutdown error: %v", err)
    }
    
    os.Exit(0)
}()

app.Run()
```

**Shutdown Steps:**
1. Stop accepting new requests
2. Wait for active requests to complete
3. Close database connections
4. Flush log buffers
5. Clean up resources

---

## ModuleRegistry

### Definition

```go
type ModuleRegistry struct {
    modules map[string]Module
    mu      sync.RWMutex
}
```

Manages module registration and lifecycle.

---

### Register

```go
func (r *ModuleRegistry) Register(module Module)
```

Registers a module with the registry.

**Parameters:**
- `module` - Module instance to register

**Example:**
```go
registry := core.NewModuleRegistry()

// Register modules
registry.Register(user.New())
registry.Register(product.New())
registry.Register(order.New())
```

---

### DiscoverModules

```go
func (r *ModuleRegistry) DiscoverModules() error
```

Automatically discovers and registers modules from `modules/` directory.

**Returns:**
- `error` - Error if discovery fails

**Example:**
```go
registry := core.NewModuleRegistry()

// Auto-discover modules
if err := registry.DiscoverModules(); err != nil {
    log.Fatalf("Module discovery failed: %v", err)
}

log.Printf("Discovered %d modules", len(registry.GetModules()))
```

**Convention:**
- Looks in `modules/` directory
- Each subdirectory is a module
- Module must have `New()` function
- Function must return type implementing `Module` interface

```
modules/
├── user/
│   └── user.go     // func New() *UserModule
├── product/
│   └── product.go  // func New() *ProductModule
└── order/
    └── order.go    // func New() *OrderModule
```

---

### GetModules

```go
func (r *ModuleRegistry) GetModules() []Module
```

Returns all registered modules.

**Returns:**
- `[]Module` - Slice of modules

**Example:**
```go
modules := registry.GetModules()

for _, module := range modules {
    fmt.Printf("Module: %s\n", module.Name())
}
```

---

### GetModule

```go
func (r *ModuleRegistry) GetModule(name string) (Module, bool)
```

Retrieves a specific module by name.

**Parameters:**
- `name` - Module name

**Returns:**
- `Module` - Module instance
- `bool` - true if found, false otherwise

**Example:**
```go
userModule, exists := registry.GetModule("user")
if !exists {
    log.Fatal("User module not found")
}

fmt.Printf("Found module: %s\n", userModule.Name())
```

---

## Config

### Definition

```go
type Config struct {
    App      AppConfig
    Database DatabaseConfig
    Logger   LoggerConfig
    Server   ServerConfig
}
```

Application configuration structure.

---

### AppConfig

```go
type AppConfig struct {
    Name        string
    Environment string  // development, staging, production
    Debug       bool
    Port        string
}
```

Application-level configuration.

**Example:**
```go
config.App = AppConfig{
    Name:        "My App",
    Environment: "production",
    Debug:       false,
    Port:        "8080",
}
```

---

### DatabaseConfig

```go
type DatabaseConfig struct {
    Driver   string  // postgres, mysql, sqlite
    Host     string
    Port     string
    Database string
    Username string
    Password string
    SSLMode  string
}
```

Database configuration.

**Example:**
```go
config.Database = DatabaseConfig{
    Driver:   "postgres",
    Host:     "localhost",
    Port:     "5432",
    Database: "myapp",
    Username: "postgres",
    Password: "secret",
    SSLMode:  "disable",
}
```

---

### LoggerConfig

```go
type LoggerConfig struct {
    Level      string  // debug, info, warn, error, fatal
    Format     string  // json, text
    Output     string  // console, file, both
    FilePath   string
    MaxSize    int
    MaxBackups int
    MaxAge     int
}
```

Logger configuration.

**Example:**
```go
config.Logger = LoggerConfig{
    Level:      "info",
    Format:     "json",
    Output:     "both",
    FilePath:   "logs/app.log",
    MaxSize:    100,  // MB
    MaxBackups: 5,
    MaxAge:     30,   // days
}
```

---

### ServerConfig

```go
type ServerConfig struct {
    Port         string
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    IdleTimeout  time.Duration
}
```

HTTP server configuration.

**Example:**
```go
config.Server = ServerConfig{
    Port:         "8080",
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  120 * time.Second,
}
```

---

### LoadConfig

```go
func LoadConfig() (*Config, error)
```

Loads configuration from environment variables and files.

**Returns:**
- `*Config` - Loaded configuration
- `error` - Error if loading fails

**Example:**
```go
config, err := core.LoadConfig()
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}

fmt.Printf("App: %s\n", config.App.Name)
fmt.Printf("Environment: %s\n", config.App.Environment)
fmt.Printf("Port: %s\n", config.Server.Port)
```

**Priority:**
1. Environment variables (highest)
2. `.env` file
3. Default values (lowest)

---

## Utility Functions

### WithTransaction

```go
func WithTransaction(ctx context.Context, db *gorm.DB, fn func(*gorm.DB) error) error
```

Executes a function within a database transaction.

**Parameters:**
- `ctx` - Context for cancellation
- `db` - Database instance
- `fn` - Function to execute in transaction

**Returns:**
- `error` - Error from transaction or function

**Example:**
```go
err := core.WithTransaction(ctx, db, func(tx *gorm.DB) error {
    // Create user
    user := &User{Name: "John"}
    if err := tx.Create(user).Error; err != nil {
        return err  // Rollback
    }
    
    // Create profile
    profile := &Profile{UserID: user.ID}
    if err := tx.Create(profile).Error; err != nil {
        return err  // Rollback
    }
    
    return nil  // Commit
})

if err != nil {
    log.Printf("Transaction failed: %v", err)
}
```

**Behavior:**
- Automatically begins transaction
- Commits if function returns nil
- Rolls back if function returns error
- Rolls back on panic

---

### Resolve

```go
func Resolve[T any](c *Container) T
```

Resolves a service from the DI container (type-safe).

**Type Parameters:**
- `T` - Service type

**Parameters:**
- `c` - Container instance

**Returns:**
- `T` - Service instance

**Example:**
```go
// Resolve database
db := core.Resolve[*gorm.DB](container)

// Resolve service
userService := core.Resolve[user.Service](container)

// Resolve repository
userRepo := core.Resolve[user.Repository](container)

// Resolve logger
log := core.Resolve[logger.Logger](container)
```

**See:** [Container API](container.md) for more details.

---

## Complete Application Example

### Full Setup

```go
// main.go
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "neonexcore/internal/core"
    "neonexcore/modules/user"
    "neonexcore/modules/product"
)

func main() {
    // Create application
    app := core.NewApp()
    
    // Register modules manually (or use auto-discovery)
    app.Registry.Register(user.New())
    app.Registry.Register(product.New())
    
    // Register models for migration
    app.RegisterModels(
        &user.User{},
        &user.Profile{},
        &product.Product{},
        &product.Category{},
    )
    
    // Setup graceful shutdown
    setupShutdown(app)
    
    // Start application
    log.Println("Starting Neonex Core...")
    if err := app.Run(); err != nil {
        log.Fatalf("Application failed: %v", err)
    }
}

func setupShutdown(app *core.App) {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        <-c
        log.Println("\nReceived shutdown signal...")
        
        if err := app.Shutdown(); err != nil {
            log.Printf("Shutdown error: %v", err)
        }
        
        log.Println("Shutdown complete")
        os.Exit(0)
    }()
}
```

### Custom Initialization

```go
func main() {
    app := core.NewApp()
    
    // Custom configuration
    app.Config.App.Port = "3000"
    app.Config.Logger.Level = "debug"
    
    // Initialize components manually
    if err := app.InitDatabase(); err != nil {
        log.Fatal(err)
    }
    
    // Discover modules
    if err := app.Registry.DiscoverModules(); err != nil {
        log.Fatal(err)
    }
    
    // Register services
    for _, module := range app.Registry.GetModules() {
        module.RegisterServices(app.Container)
    }
    
    // Custom setup
    setupCustomMiddleware(app)
    setupCustomRoutes(app)
    
    // Start server
    log.Printf("Server starting on port %s", app.Config.App.Port)
    if err := app.StartHTTP(); err != nil {
        log.Fatal(err)
    }
}
```

---

## Environment Variables

### Required Variables

```bash
# Application
APP_NAME=my-app
APP_ENV=production
APP_PORT=8080

# Database
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=myapp
DB_USERNAME=postgres
DB_PASSWORD=secret

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=both
LOG_FILE=logs/app.log
```

### Optional Variables

```bash
# Debug
DEBUG=false

# Database
DB_SSL_MODE=disable
DB_MAX_CONNECTIONS=100
DB_MAX_IDLE=10

# Server
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s
SERVER_IDLE_TIMEOUT=120s

# Logging
LOG_MAX_SIZE=100
LOG_MAX_BACKUPS=5
LOG_MAX_AGE=30
```

---

## Best Practices

### ✅ DO: Handle Errors

```go
// Good: Check all errors
if err := app.InitDatabase(); err != nil {
    log.Fatalf("Database init failed: %v", err)
}

if err := app.AutoMigrate(); err != nil {
    log.Fatalf("Migration failed: %v", err)
}
```

### ✅ DO: Graceful Shutdown

```go
// Good: Handle signals
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt, syscall.SIGTERM)

go func() {
    <-c
    app.Shutdown()
    os.Exit(0)
}()
```

### ✅ DO: Use Environment Variables

```go
// Good: Configurable
port := os.Getenv("APP_PORT")
if port == "" {
    port = "8080"
}
```

### ❌ DON'T: Ignore Errors

```go
// Bad: Silent failure
app.InitDatabase()  // What if it fails?
app.Run()
```

### ❌ DON'T: Hardcode Configuration

```go
// Bad: Hardcoded
app.Config.Database.Host = "localhost"
app.Config.Database.Password = "secret123"

// Good: From environment
app.Config.Database.Host = os.Getenv("DB_HOST")
app.Config.Database.Password = os.Getenv("DB_PASSWORD")
```

---

## Related Documentation

- [**Container API**](container.md) - Dependency injection
- [**Module Interface**](module-interface.md) - Module system
- [**Getting Started**](../getting-started/quick-start.md) - Quick start guide
- [**Configuration**](../getting-started/configuration.md) - Configuration guide

---

## Source Code

- [GitHub: app.go](https://github.com/neonextechnologies/neonexcore/blob/main/internal/core/app.go)
- [GitHub: registry.go](https://github.com/neonextechnologies/neonexcore/blob/main/internal/core/registry.go)
- [GitHub: config.go](https://github.com/neonextechnologies/neonexcore/blob/main/internal/config/)

---

**Need help?** Check our [FAQ](../resources/faq.md) or [get support](../resources/support.md).
