# Neonex Core Framework

<div align="center">

![Neonex Core](https://img.shields.io/badge/Neonex-Core-blue?style=for-the-badge)
![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

**A Modern, Modular Go Framework with Built-in ORM**

[Features](#-features) • [Quick Start](#-quick-start) • [Documentation](#-documentation) • [Architecture](#-architecture)

</div>

---

## 🚀 Features

- **🎯 Modular Architecture** - Auto-discovery module system with dependency injection
- **⚡ High Performance** - Built on Fiber (fasthttp) for blazing fast HTTP handling
- **🗄️ Advanced ORM Layer** - Generic repository pattern with GORM
- **💉 Dependency Injection** - Type-safe DI container with Singleton/Transient scopes
- **🔄 Auto-Migration** - Database schema management out of the box
- **🌱 Database Seeding** - Initial data management system
- **🔌 Multi-Database Support** - SQLite, MySQL, PostgreSQL, Turso
- **🏗️ Transaction Manager** - ACID-compliant transaction handling
- **📦 Zero Configuration** - Works with SQLite out of the box

## 📦 Quick Start

### Prerequisites
- Go 1.21 or higher

### Installation

```bash
# Clone the repository
git clone https://github.com/YOUR_USERNAME/neonexcore.git
cd neonexcore

# Download dependencies
go mod download

# Run the application
go run main.go
```

The server will start at `http://localhost:8080`

### Test the API

```bash
# Get all users
curl http://localhost:8080/user/

# Get user by ID
curl http://localhost:8080/user/1

# Search users
curl http://localhost:8080/user/search?q=alice

# Create user
curl -X POST http://localhost:8080/user/ \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com","password":"secret","age":25}'
```

## 🏗️ Architecture

```
neonexcore/
├── cmd/                    # CLI commands (future)
├── internal/
│   ├── config/            # Configuration management
│   └── core/              # Framework core
│       ├── app.go         # Application orchestrator
│       ├── container.go   # DI container
│       ├── modulemap.go   # Module registry
│       └── registry.go    # Module discovery
├── modules/               # Application modules
│   └── user/             # Example user module
│       ├── controller.go # HTTP handlers
│       ├── di.go         # Dependency injection
│       ├── model.go      # GORM model
│       ├── repository.go # Data access layer
│       ├── routes.go     # Route definitions
│       ├── seeder.go     # Seed data
│       ├── service.go    # Business logic
│       └── module.json   # Module metadata
├── pkg/
│   ├── database/         # ORM utilities
│   │   ├── repository.go # Generic repository
│   │   ├── transaction.go# Transaction manager
│   │   ├── migrator.go   # Auto-migration
│   │   └── seeder.go     # Seeder system
│   └── logger/           # Logging utilities
└── main.go               # Application entry point
```

## 📚 Documentation

### Creating a New Module

1. **Create module directory:**
```bash
mkdir -p modules/mymodule
```

2. **Create `module.json`:**
```json
{
  "name": "mymodule",
  "enabled": true
}
```

3. **Implement Module interface:**
```go
package mymodule

type MyModule struct{}

func New() *MyModule { return &MyModule{} }

func (m *MyModule) Name() string { return "mymodule" }
func (m *MyModule) Init() { /* initialization */ }
func (m *MyModule) Routes(app *fiber.App, c *core.Container) { /* routes */ }
func (m *MyModule) RegisterServices(c *core.Container) { /* DI setup */ }
```

4. **Register in `main.go`:**
```go
core.ModuleMap["mymodule"] = func() core.Module { return mymodule.New() }
```

### Using the Repository Pattern

```go
// Define model
type Product struct {
    ID    uint   `gorm:"primarykey"`
    Name  string `gorm:"size:255"`
    Price float64
}

// Create repository
repo := database.NewBaseRepository[Product](db)

// Use CRUD operations
products, _ := repo.FindAll(ctx)
product, _ := repo.FindByID(ctx, 1)
repo.Create(ctx, &newProduct)
repo.Update(ctx, &product)
```

### Transaction Example

```go
txManager := database.NewTxManager(db)

err := txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
    // All operations here are in a transaction
    repo := repository.WithTx(tx)
    repo.Create(ctx, &user1)
    repo.Create(ctx, &user2)
    return nil // Auto commit on success
})
```

## 🔧 Configuration

Environment variables (optional):
```bash
DB_DRIVER=sqlite           # sqlite, mysql, postgres, turso
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=
DB_DATABASE=neonex.db
```

## 🎯 Roadmap

- [x] Core framework structure
- [x] Module system with auto-discovery
- [x] Generic DI container
- [x] ORM layer with repository pattern
- [x] Transaction management
- [x] Auto-migration
- [ ] CLI tools (`neonex new`, `neonex serve`)
- [ ] Structured logging system
- [ ] Hot reload support
- [ ] Authentication/Authorization
- [ ] Middleware system
- [ ] Validation framework
- [ ] API documentation generator
- [ ] Plugin ecosystem

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🌟 Acknowledgments

- Built with [Fiber](https://github.com/gofiber/fiber) - Fast HTTP framework
- Powered by [GORM](https://gorm.io) - ORM library for Go
- Inspired by Laravel, NestJS, and Spring Boot

---

<div align="center">

**Made with ❤️ for the Go community**

[Report Bug](https://github.com/YOUR_USERNAME/neonexcore/issues) • [Request Feature](https://github.com/YOUR_USERNAME/neonexcore/issues)

</div>
