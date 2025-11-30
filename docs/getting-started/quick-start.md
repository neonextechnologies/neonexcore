# Quick Start

Build your first high-performance API with Neonex Core in under 5 minutes! ⚡

---

## Overview

This quick start guide will walk you through:

1. ✨ Creating a new project
2. 🚀 Running the application
3. 🧩 Creating your first module
4. 🔌 Testing the API endpoints
5. 📊 Viewing structured logs

**Time required:** ~5 minutes

**Prerequisites:** [Neonex Core installed](installation.md)

---

## Step 1: Create a New Project

Create a new Neonex Core project using the CLI:

```bash
# Create project
neonex new my-app

# Output:
# ✨ Creating new Neonex Core project...
# 📁 my-app/
# ├── 📂 cmd/neonex/
# ├── 📂 internal/
# ├── 📂 modules/
# ├── 📂 pkg/
# ├── 📄 main.go
# ├── 📄 go.mod
# ├── 📄 .air.toml
# ├── 📄 .env.example
# └── 📄 README.md
# ✅ Project created successfully!

# Navigate to project
cd my-app
```

This creates a complete project with:
- **Modular structure** for scalable architecture
- **Dependency injection** container pre-configured
- **Database** connection ready (GORM)
- **Logging** system (Zap) set up
- **Hot reload** configuration (.air.toml)
- **Environment** configuration (.env.example)

---

## Step 2: Install Dependencies

Download all required Go modules:

```bash
# Install dependencies
go mod download

# Output:
# go: downloading github.com/gofiber/fiber/v2 v2.52.0
# go: downloading gorm.io/gorm v1.25.5
# go: downloading go.uber.org/zap v1.26.0
# ... (other dependencies)
```

> **💡 Tip:** This step is automatic if you use `neonex serve` directly.

---

## Step 3: Run the Application

Start your application with one command:

### Option 1: Standard Mode

```bash
neonex serve
```

**Output:**
```
🚀 Starting Neonex Core v0.1-alpha...
📂 Loading configuration from .env
🗄️  Connecting to database (SQLite: neonex.db)
✅ Database connected successfully
📦 Discovering modules...
✅ User module registered (v1.0.0)
✅ Auth module registered (v1.0.0)
🔧 Running auto-migrations...
✅ Migrations completed
🌱 Running seeders...
✅ Seeders completed
🌐 HTTP server starting...
⚡ Server ready at http://localhost:8080
📝 Logs: logs/neonex.log
```

### Option 2: Hot Reload Mode (Recommended for Development)

Automatically restart on file changes:

```bash
neonex serve --hot

# Or use Air directly
air
```

**Output:**
```
🔥 Hot reload enabled (Air)
🚀 Starting Neonex Core...
⚡ Server ready at http://localhost:8080
👀 Watching for file changes...
```

### Option 3: Custom Port

```bash
# Run on different port
neonex serve --port 3000

# Run with environment
neonex serve --env production
```

---

## Step 4: Test the API

### Test Default Endpoint

Open your browser or use curl:

```bash
# Test with curl
curl http://localhost:8080

# Or with PowerShell
Invoke-RestMethod http://localhost:8080
```

**Response:**
```json
{
  "framework": "Neonex Core",
  "version": "0.1-alpha",
  "status": "running",
  "engine": "Fiber (fasthttp)",
  "timestamp": "2024-01-15T10:30:00Z",
  "uptime": "5s",
  "environment": "development"
}
```

### Test Built-in User API

The default project includes a User module:

```bash
# Get all users
curl http://localhost:8080/api/users

# Response:
# {
#   "data": [],
#   "meta": {
#     "total": 0,
#     "page": 1,
#     "per_page": 10
#   }
# }

# Create a user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "secret123"
  }'

# Response:
# {
#   "id": 1,
#   "name": "John Doe",
#   "email": "john@example.com",
#   "created_at": "2024-01-15T10:30:00Z"
# }

# Get user by ID
curl http://localhost:8080/api/users/1

# Update user
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Doe"}'

# Delete user
curl -X DELETE http://localhost:8080/api/users/1
```

---

## Step 5: Create Your First Module

```bash
neonex module create product
```

This generates a complete CRUD module with:

* Model and database table
* Repository for data access
* Service for business logic
* Controller for HTTP handlers
* Routes automatically registered
* Dependency injection configured

## Step 6: Register the Module

Edit `main.go`:

```go
package main

import (
    "neonexcore/internal/core"
    "my-app/modules/product"  // Add this
)

func main() {
    // Register module
    core.ModuleMap["product"] = func() core.Module { 
        return product.New() 
    }
    
    app := core.NewApp()
    
    // Register model for migration
    app.RegisterModels(&product.Product{})
    
    // ... rest of the code
}
```

## Step 7: Test Your Module

Restart the server and test the endpoints:

### Get All Products

```bash
curl http://localhost:8080/product/
```

### Create a Product

```bash
curl -X POST http://localhost:8080/product/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "description": "Gaming laptop",
    "is_active": true
  }'
```

### Get Product by ID

```bash
curl http://localhost:8080/product/1
```

## What's Next?

Now that you have a working application, explore more features:

* [Project Structure](project-structure.md) - Understand the folder organization
* [Module System](../core-concepts/module-system.md) - Learn about modules
* [Database](../database/overview.md) - Work with the database
* [Logging](../logging/overview.md) - Add structured logging
* [Hot Reload](../development/hot-reload.md) - Speed up development

## Common Tasks

### Add More Modules

```bash
neonex module create category
neonex module create order
neonex module create customer
```

### Change Database

Edit `.env`:

```bash
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=secret
DB_DATABASE=myapp
```

### Add Custom Routes

In your module's `routes.go`:

```go
func RegisterRoutes(app *fiber.App, container *core.Container) {
    ctrl := container.Resolve("product.Controller").(*Controller)
    
    group := app.Group("/product")
    group.Get("/", ctrl.GetAll)
    group.Get("/featured", ctrl.GetFeatured)  // Custom route
    group.Post("/", ctrl.Create)
}
```

### Enable Debug Logging

Edit `.env`:

```bash
LOG_LEVEL=debug
LOG_FORMAT=text
LOG_OUTPUT=both
```

## Tips

* Use `neonex serve --hot` for development
* Run `neonex module list` to see all modules
* Check `logs/` folder for application logs
* Use `.env` for environment-specific configuration

## Getting Help

* [Documentation](../README.md)
* [GitHub Issues](https://github.com/neonextechnologies/neonexcore/issues)
* [Discussions](https://github.com/neonextechnologies/neonexcore/discussions)
