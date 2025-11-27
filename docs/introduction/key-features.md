# Key Features

Neonex Core provides a comprehensive set of features for modern application development.

## 🎯 Modular Architecture

### Auto-Discovery

Modules are automatically discovered and loaded based on a simple JSON configuration:

```json
{
  "name": "user",
  "version": "1.0.0",
  "enabled": true
}
```

No need to manually register modules in your bootstrap code.

### Self-Contained Modules

Each module is completely self-contained with its own:

* Routes and controllers
* Services and business logic
* Repositories and data access
* Models and migrations
* Seeders and test data

### Example Module Structure

```
modules/user/
├── user.go          # Module entry point
├── model.go         # GORM models
├── repository.go    # Data access layer
├── service.go       # Business logic
├── controller.go    # HTTP handlers
├── routes.go        # Route definitions
├── di.go           # Dependency injection
├── seeder.go       # Database seeders
└── module.json     # Module metadata
```

## ⚡ High Performance

### Built on Fiber

Neonex Core uses [Fiber](https://gofiber.io/), one of the fastest web frameworks:

* **Zero memory allocation** in hot paths
* **Express-like API** familiar to web developers
* **Fast routing** with zero allocations
* **Low memory footprint**

### Performance Benchmarks

```
Requests/sec:  100,000+
Latency:       <1ms (p50)
Memory:        ~10MB base
Throughput:    >2GB/s
```

## 🗄️ Advanced ORM Layer

### Generic Repository Pattern

```go
type Repository[T any] struct {
    DB *gorm.DB
}

// Works with any model
userRepo := NewRepository[User](db)
productRepo := NewRepository[Product](db)
```

### Built-in Operations

* CRUD operations
* Pagination and filtering
* Soft deletes
* Transactions
* Eager loading
* Custom queries

### Multi-Database Support

* SQLite (default, zero-config)
* MySQL / MariaDB
* PostgreSQL
* Turso (Edge database)

## 💉 Dependency Injection

### Type-Safe DI Container

```go
// Register services
container.Singleton("userRepo", func(c *Container) interface{} {
    return NewUserRepository(db)
})

container.Singleton("userService", func(c *Container) interface{} {
    repo := c.Resolve("userRepo").(*UserRepository)
    return NewUserService(repo)
})

// Resolve with type safety
service := container.Resolve("userService").(*UserService)
```

### Lifecycle Management

* **Singleton** - Single instance shared across the application
* **Transient** - New instance on every resolution
* **Scoped** - Instance per request (coming soon)

## 📝 Structured Logging

### Multiple Log Levels

```go
logger.Debug("Detailed debug info")
logger.Info("General information")
logger.Warn("Warning messages")
logger.Error("Error occurred")
logger.Fatal("Critical failure")
```

### Rich Context

```go
logger.Info("User logged in", logger.Fields{
    "user_id": 123,
    "ip": "192.168.1.1",
    "action": "login",
    "timestamp": time.Now(),
})
```

### Multiple Outputs

* Console (colored, human-readable)
* File (with automatic rotation)
* JSON (machine-parseable)
* Custom writers

### HTTP Request Logging

Automatic logging of all HTTP requests with:

* Method, path, status code
* Request duration
* IP address and user agent
* Request ID for tracing
* Query parameters

## 🛠️ CLI Tools

### Project Scaffolding

```bash
neonex new my-project
```

Creates a complete project structure with all necessary files.

### Module Generation

```bash
neonex module create product
```

Generates a full module with:

* Model, repository, service, controller
* Routes and dependency injection
* Seeder and tests
* Module configuration

### Development Server

```bash
neonex serve --hot
```

Starts development server with hot reload enabled.

## 🔥 Hot Reload

### Instant Feedback

Changes are automatically detected and the application rebuilds:

1. Save your file
2. Auto-rebuild (1-3 seconds)
3. Test immediately

### Configurable Watching

* Watch specific directories
* Exclude patterns
* Custom build commands
* Delay configuration

### Integration

Works seamlessly with:

* VSCode
* GoLand / IntelliJ
* Command line
* Make / Task runners

## 🔄 Auto-Migration

### Zero-Config Migrations

```go
// Register models
app.RegisterModels(
    &user.User{},
    &product.Product{},
)

// Run migrations
app.AutoMigrate()
```

### Smart Updates

* Adds new tables automatically
* Adds new columns to existing tables
* Preserves existing data
* No manual migration files needed

## 🌱 Database Seeding

### Simple Seeders

```go
type UserSeeder struct{}

func (s *UserSeeder) Seed(ctx context.Context, db *gorm.DB) error {
    users := []User{
        {Name: "Alice", Email: "alice@example.com"},
        {Name: "Bob", Email: "bob@example.com"},
    }
    return db.Create(&users).Error
}
```

### Seeder Management

```go
seeder := database.NewSeederManager(db)
seeder.Register(&UserSeeder{})
seeder.Register(&ProductSeeder{})
seeder.Run(ctx)
```

## 🏗️ Transaction Manager

### ACID Compliance

```go
err := txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
    // All operations use tx
    if err := userRepo.Create(tx, user); err != nil {
        return err // Auto rollback
    }
    if err := orderRepo.Create(tx, order); err != nil {
        return err // Auto rollback
    }
    return nil // Auto commit
})
```

### Automatic Management

* Auto-commit on success
* Auto-rollback on error
* Nested transaction support
* Context propagation

## 📦 Zero Configuration

### Works Out of the Box

```go
func main() {
    app := core.NewApp()
    app.InitDatabase()  // Uses SQLite by default
    app.StartHTTP()     // Starts on :8080
}
```

### Environment-Based Config

```bash
DB_DRIVER=mysql
DB_HOST=localhost
LOG_LEVEL=debug
HTTP_PORT=3000
```

## Next Steps

* [Explore the Architecture](architecture.md)
* [Install Neonex Core](../getting-started/installation.md)
* [Build your first app](../getting-started/quick-start.md)
