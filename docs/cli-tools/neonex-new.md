# neonex new

Create a new Neonex Core project with complete structure and configuration.

---

## Synopsis

```bash
neonex new <project-name> [flags]
```

Scaffolds a complete project with:
- ✅ Directory structure (cmd, internal, modules, pkg)
- ✅ Go module files (go.mod, go.sum)
- ✅ Configuration files (.env, .air.toml, Makefile)
- ✅ Example module (user)
- ✅ Git repository
- ✅ Documentation (README.md)

---

## Usage

### Basic Command

```bash
neonex new my-app
```

**Output:**
```
Creating project: my-app
✓ Created directory structure
✓ Generated go.mod
✓ Created main.go
✓ Generated .env file
✓ Created .air.toml
✓ Generated Makefile
✓ Created user module
✓ Initialized git repository

Project created successfully!

Next steps:
  cd my-app
  go mod download
  neonex serve --hot

Happy coding! 🚀
```

### With Custom Path

```bash
# Create in specific directory
neonex new ~/projects/my-app

# Create in current directory
neonex new .

# Create nested path (will create parent directories)
neonex new backend/api/my-app
```

---

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--no-git` | bool | false | Skip git initialization |
| `--no-example` | bool | false | Don't create example module |
| `--module-path` | string | auto | Custom Go module path |
| `--go-version` | string | 1.21 | Go version requirement |

### Examples with Flags

```bash
# Skip git init
neonex new my-app --no-git

# No example module
neonex new my-app --no-example

# Custom module path
neonex new my-app --module-path github.com/myuser/my-app

# Specific Go version
neonex new my-app --go-version 1.22
```

---

## Generated Structure

### Complete Directory Tree

```
my-app/
├── cmd/
│   └── neonex/              # CLI commands (optional)
├── internal/                # Private application code
│   ├── config/
│   │   └── database.go     # Database configuration
│   ├── core/
│   │   ├── app.go          # Application initialization
│   │   ├── container.go    # DI container
│   │   ├── modulemap.go    # Module registry
│   │   └── registry.go     # Service registry
│   └── module/
│       └── module.go       # Module interface
├── modules/                 # Application modules
│   └── user/               # Example user module
│       ├── controller.go
│       ├── di.go
│       ├── model.go
│       ├── module.json
│       ├── repository.go
│       ├── routes.go
│       ├── seeder.go
│       ├── service.go
│       └── user.go
├── pkg/                    # Public libraries
│   ├── database/
│   │   ├── migrator.go    # Auto-migration
│   │   ├── repository.go  # Generic repository
│   │   ├── seeder.go      # Seeding interface
│   │   └── transaction.go # Transaction manager
│   ├── http/
│   │   └── server.go      # Fiber server wrapper
│   ├── logger/
│   │   ├── config.go
│   │   ├── formatter.go
│   │   ├── logger.go
│   │   ├── middleware.go
│   │   └── writer.go
│   └── utils/
│       └── helpers.go
├── .air.toml               # Hot reload config
├── .env                    # Environment variables
├── .env.example            # Environment template
├── .gitignore              # Git ignore rules
├── go.mod                  # Go module
├── main.go                 # Application entry point
├── Makefile                # Build automation
└── README.md               # Project documentation
```

---

## Generated Files

### main.go

```go
package main

import (
    "log"
    "github.com/YOUR_USERNAME/my-app/internal/core"
)

func main() {
    app := core.NewApp()
    
    if err := app.Run(); err != nil {
        log.Fatalf("Failed to start application: %v", err)
    }
}
```

### go.mod

```go
module github.com/YOUR_USERNAME/my-app

go 1.21

require (
    github.com/gofiber/fiber/v2 v2.52.9
    github.com/joho/godotenv v1.5.1
    gorm.io/driver/sqlite v1.5.6
    gorm.io/gorm v1.25.12
)
```

### .env

```bash
# Application
APP_NAME=my-app
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

### .env.example

Template for team members:

```bash
# Application
APP_NAME=my-app
APP_ENV=development
APP_PORT=8080

# Database (configure based on your setup)
DB_DRIVER=sqlite
DB_HOST=
DB_PORT=
DB_NAME=neonex.db
DB_USER=
DB_PASSWORD=

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
LOG_FILE=logs/app.log
```

### .air.toml

Hot reload configuration:

```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main.exe"
  cmd = "go build -o ./tmp/main.exe ."
  delay = 1000
  exclude_dir = ["tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
```

### Makefile

Build automation:

```makefile
.PHONY: help dev build test clean install air

help:
	@echo "Available commands:"
	@echo "  make dev        - Start development server with hot reload"
	@echo "  make build      - Build the application"
	@echo "  make test       - Run tests"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make install    - Install dependencies"
	@echo "  make air        - Install Air for hot reload"

dev:
	@echo "Starting development server..."
	@air

build:
	@echo "Building application..."
	@go build -o bin/app main.go

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning..."
	@rm -rf tmp/ bin/ *.db

install:
	@echo "Installing dependencies..."
	@go mod download

air:
	@echo "Installing Air..."
	@go install github.com/air-verse/air@latest
```

### .gitignore

```gitignore
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
bin/
tmp/

# Test binary
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# Go workspace file
go.work

# Environment files
.env
.env.local

# Database files
*.db
*.sqlite
*.sqlite3

# Logs
logs/
*.log

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db
```

### README.md

```markdown
# my-app

A Neonex Core application.

## Quick Start

### Prerequisites
- Go 1.21 or higher

### Installation

\`\`\`bash
# Install dependencies
go mod download

# Run with hot reload
neonex serve --hot

# Or use make
make dev
\`\`\`

The server will start at `http://localhost:8080`

### Available Endpoints

- `GET /user/` - List all users
- `GET /user/:id` - Get user by ID
- `POST /user/` - Create new user
- `PUT /user/:id` - Update user
- `DELETE /user/:id` - Delete user

### Development

\`\`\`bash
# Start development server
make dev

# Build for production
make build

# Run tests
make test

# Clean artifacts
make clean
\`\`\`

### Project Structure

\`\`\`
my-app/
├── internal/        # Private application code
├── modules/         # Feature modules
├── pkg/            # Public libraries
└── main.go         # Entry point
\`\`\`

## Documentation

- [Neonex Core Docs](https://github.com/neonextechnologies/neonexcore)
- [Go Documentation](https://golang.org/doc/)

## License

MIT
```

---

## Post-Creation Steps

### 1. Navigate to Project

```bash
cd my-app
```

### 2. Install Dependencies

```bash
# Download all dependencies
go mod download

# Or using make
make install
```

### 3. Configure Environment

```bash
# Copy example config
cp .env.example .env

# Edit configuration
nano .env  # or your preferred editor
```

### 4. Start Development

```bash
# With hot reload
neonex serve --hot

# Or using make
make dev

# Or directly
go run main.go
```

### 5. Test the API

```bash
# List users
curl http://localhost:8080/user/

# Create user
curl -X POST http://localhost:8080/user/ \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com"}'
```

---

## Customization

### Custom Module Path

By default, uses project name as module path. Override with:

```bash
neonex new my-app --module-path github.com/myorg/my-app
```

**Generated go.mod:**
```go
module github.com/myorg/my-app

go 1.21
// ...
```

### Skip Git Initialization

If you already have a git repository:

```bash
neonex new my-app --no-git
```

### Skip Example Module

For a minimal project:

```bash
neonex new my-app --no-example
```

**Result:** No `modules/user/` directory

---

## Examples

### Example 1: Simple API

```bash
# Create project
neonex new todo-api
cd todo-api

# Generate module
neonex module create todo

# Start server
neonex serve --hot

# Test
curl http://localhost:8080/todo/
```

### Example 2: Microservice

```bash
# Create with custom module path
neonex new order-service --module-path github.com/mycompany/order-service
cd order-service

# Generate modules
neonex module create order
neonex module create payment
neonex module create shipping

# Configure for production
nano .env

# Build
go build -o bin/order-service main.go
```

### Example 3: Monorepo

```bash
# Create multiple services
mkdir backend
cd backend

neonex new api-gateway
neonex new user-service
neonex new product-service
neonex new order-service

# Each service is independent
cd api-gateway && neonex serve --hot
```

---

## Troubleshooting

### Permission Denied

**Problem:**
```bash
mkdir: cannot create directory 'my-app': Permission denied
```

**Solution:**
```bash
# Choose different location
cd ~/projects
neonex new my-app

# Or use sudo (not recommended)
sudo neonex new my-app
sudo chown -R $USER:$USER my-app
```

### Module Path Conflict

**Problem:**
```bash
go: module github.com/user/my-app: git ls-remote failed
```

**Solution:**
```bash
# Use custom path
neonex new my-app --module-path example.com/my-app

# Or fix later in go.mod
cd my-app
nano go.mod  # Edit module path
go mod tidy
```

### Directory Already Exists

**Problem:**
```bash
Error: Directory 'my-app' already exists
```

**Solution:**
```bash
# Use different name
neonex new my-app-v2

# Or remove existing
rm -rf my-app
neonex new my-app
```

---

## Advanced Usage

### Create in Current Directory

```bash
# Initialize in existing folder
mkdir my-project
cd my-project
neonex new .
```

### Custom Go Version

```bash
# Use Go 1.22
neonex new my-app --go-version 1.22

# Check go.mod
cat my-app/go.mod
# module my-app
# go 1.22
```

### Batch Project Creation

```bash
# Create multiple projects
for name in api-gateway user-service product-service; do
  neonex new $name
done
```

---

## Next Steps

- [**neonex serve**](neonex-serve.md) - Start development server
- [**neonex module**](neonex-module.md) - Add modules to your project
- [**Project Structure**](../getting-started/project-structure.md) - Understand the layout
- [**Configuration**](../getting-started/configuration.md) - Configure your app

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
