# Project Structure

Learn about the directory organization and file structure of a Neonex Core application.

---

## Overview

Neonex Core follows the **Standard Go Project Layout** with additional conventions for modular architecture:

```
my-app/
├── cmd/                    # Application entry points
│   └── neonex/            # CLI commands
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── core/             # Core framework components
│   └── module/           # Module interface
├── modules/              # Application modules
│   ├── auth/            # Authentication module
│   ├── user/            # User module
│   └── product/         # Product module (example)
├── pkg/                  # Public libraries
│   ├── database/        # Database utilities
│   ├── http/            # HTTP server
│   ├── logger/          # Logging system
│   └── utils/           # Utility functions
├── main.go              # Application entry point
├── go.mod               # Go module definition
├── .env                 # Environment variables
├── .air.toml            # Hot reload configuration
└── Makefile             # Build automation
```

---

## Directory Breakdown

### `cmd/` - Command Line Interface

Contains executable commands for your application.

```
cmd/
└── neonex/
    ├── main.go         # CLI entry point
    ├── new.go          # Project scaffolding
    ├── serve.go        # Development server
    └── module.go       # Module management
```

**Purpose:**
- CLI tools for scaffolding and development
- Separate from main application code
- Each subdirectory is a different executable

### `internal/` - Private Application Code

Code that cannot be imported by other projects.

```
internal/
├── config/
│   └── database.go      # Database configuration
├── core/
│   ├── app.go          # Application initialization
│   ├── container.go    # Dependency injection container
│   ├── modulemap.go    # Module registry
│   └── registry.go     # Service registry
└── module/
    └── module.go       # Module interface definition
```

**Purpose:**
- Core framework logic
- Configuration management
- Module system implementation
- Not accessible outside this project

### `modules/` - Application Modules

Feature-based module organization.

```
modules/
└── user/               # Example: User module
    ├── model.go       # GORM model
    ├── repository.go  # Data access layer
    ├── service.go     # Business logic
    ├── controller.go  # HTTP handlers
    ├── routes.go      # Route definitions
    ├── di.go          # Dependency injection
    ├── seeder.go      # Database seeding
    ├── user.go        # Module entry point
    └── module.json    # Module metadata
```

**Module Structure:**
Each module follows a layered architecture:

1. **Model Layer** - Database entities
2. **Repository Layer** - Data access
3. **Service Layer** - Business logic
4. **Controller Layer** - HTTP handlers
5. **Routes** - API endpoints
6. **DI** - Dependency wiring
7. **Seeder** - Test/initial data

### `pkg/` - Public Libraries

Reusable packages that can be imported by other projects.

```
pkg/
├── database/
│   ├── migrator.go      # Auto-migration
│   ├── repository.go    # Generic repository
│   ├── seeder.go        # Seeding interface
│   └── transaction.go   # Transaction manager
├── http/
│   └── server.go        # Fiber server wrapper
├── logger/
│   ├── logger.go        # Logger interface
│   ├── formatter.go     # Log formatters
│   ├── writer.go        # File writers
│   ├── middleware.go    # HTTP middleware
│   └── config.go        # Logger config
└── utils/
    └── helpers.go       # Utility functions
```

**Purpose:**
- Framework utilities
- Database abstractions
- HTTP server management
- Logging system
- Can be extracted to separate packages

---

## Key Files

### `main.go` - Application Entry Point

```go
package main

import (
    "log"
    "github.com/YOUR_USERNAME/neonexcore/internal/core"
)

func main() {
    app := core.NewApp()
    
    if err := app.Run(); err != nil {
        log.Fatalf("Failed to start application: %v", err)
    }
}
```

The entry point:
- Initializes the application
- Loads configuration
- Discovers modules
- Starts HTTP server

### `go.mod` - Module Dependencies

```go
module github.com/YOUR_USERNAME/neonexcore

go 1.21

require (
    github.com/gofiber/fiber/v2 v2.52.9
    gorm.io/gorm v1.25.12
    gorm.io/driver/sqlite v1.5.6
    github.com/spf13/cobra v1.8.0
)
```

Manages:
- Go version requirement
- External dependencies
- Module path

### `.env` - Environment Configuration

```bash
# Application
APP_NAME=neonex-core
APP_ENV=development
APP_PORT=8080

# Database
DB_DRIVER=sqlite
DB_HOST=localhost
DB_PORT=5432
DB_NAME=neonex.db
DB_USER=
DB_PASSWORD=

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
LOG_FILE=logs/app.log
```

Configuration for:
- Application settings
- Database connection
- Logging behavior
- Never commit `.env` to version control

### `.air.toml` - Hot Reload Configuration

```toml
[build]
  cmd = "go build -o ./tmp/main.exe ."
  bin = "tmp/main.exe"
  include_ext = ["go", "tpl", "tmpl", "html"]
  exclude_dir = ["tmp", "vendor"]
  delay = 1000
```

Controls:
- File watching patterns
- Build commands
- Reload behavior
- Development workflow

### `Makefile` - Build Automation

```makefile
.PHONY: dev build test clean

dev:
	@air

build:
	@go build -o bin/app main.go

test:
	@go test -v ./...

clean:
	@rm -rf tmp/ bin/
```

Common tasks:
- Development mode
- Building binaries
- Running tests
- Cleaning artifacts

---

## Module Organization

### Creating a New Module

When you run `neonex module create product`, it generates:

```
modules/product/
├── model.go          # Product struct with GORM tags
├── repository.go     # ProductRepository interface + implementation
├── service.go        # ProductService with business logic
├── controller.go     # ProductController with HTTP handlers
├── routes.go         # RegisterRoutes(router fiber.Router)
├── di.go            # RegisterDependencies(c *core.Container)
├── seeder.go        # ProductSeeder for test data
├── product.go       # Module entry point (implements module.Module)
└── module.json      # Metadata (name, version, dependencies)
```

### Module Lifecycle

1. **Discovery** - Framework scans `modules/` directory
2. **Registration** - Calls `RegisterDependencies()` for DI
3. **Migration** - Registers models for auto-migration
4. **Seeding** - Registers seeders for data population
5. **Routing** - Calls `RegisterRoutes()` to mount endpoints

---

## Best Practices

### ✅ DO:

- **Keep modules independent** - Avoid tight coupling between modules
- **Use interfaces** - Define contracts in repository/service layers
- **Follow naming conventions** - `UserRepository`, `UserService`, `UserController`
- **Separate concerns** - Repository for data, Service for logic, Controller for HTTP
- **Use dependency injection** - Register all dependencies in `di.go`

### ❌ DON'T:

- **Don't import between modules** - Use shared packages in `pkg/` instead
- **Don't mix layers** - Don't call repository directly from controller
- **Don't hardcode config** - Use environment variables
- **Don't commit secrets** - Add `.env` to `.gitignore`
- **Don't nest modules** - Keep flat structure in `modules/`

---

## Growing Your Project

### Adding Shared Logic

Create packages in `pkg/` for code reused across modules:

```
pkg/
├── auth/           # Authentication utilities
├── validation/     # Input validation
├── middleware/     # Custom HTTP middleware
└── errors/         # Error handling
```

### Internal Packages

Use `internal/` for private shared code:

```
internal/
├── dto/           # Data transfer objects
├── constants/     # Application constants
└── helpers/       # Internal utilities
```

### Testing Structure

Mirror your code structure in tests:

```
modules/user/
├── model_test.go
├── repository_test.go
├── service_test.go
└── controller_test.go
```

---

## Next Steps

- [Configuration Guide](configuration.md) - Learn about environment variables
- [CLI Tools](../cli-tools/overview.md) - Master the `neonex` command
- [Module System](../core-concepts/module-system.md) - Deep dive into modules

---

**Need help?** Check the [FAQ](../resources/faq.md) or [ask the community](../resources/support.md).
