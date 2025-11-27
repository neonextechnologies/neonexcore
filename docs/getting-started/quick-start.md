# Quick Start

This guide will help you create your first Neonex Core application in under 5 minutes.

## Step 1: Create a New Project

```bash
neonex new my-app
cd my-app
```

This creates a complete project structure:

```
my-app/
├── cmd/
├── internal/
├── modules/
├── pkg/
├── main.go
├── go.mod
├── .air.toml
├── .env.example
└── README.md
```

## Step 2: Install Dependencies

```bash
go mod download
```

## Step 3: Run the Application

```bash
neonex serve
```

Or with hot reload:

```bash
neonex serve --hot
```

The server will start at `http://localhost:8080`.

## Step 4: Test the API

Open your browser or use curl:

```bash
curl http://localhost:8080
```

You should see:

```json
{
  "framework": "Neonex Core",
  "version": "0.1-alpha",
  "status": "running"
}
```

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
