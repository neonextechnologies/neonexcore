# Custom Modules

Learn how to create custom modules in Neonex Core.

---

## Overview

Modules are self-contained features that encapsulate related functionality. Each module can have:
- **Models** - Data structures
- **Repository** - Data access
- **Service** - Business logic
- **Controller** - HTTP handlers
- **Routes** - API endpoints
- **Seeders** - Sample data

---

## Create Module

### Using CLI

```powershell
# Create new module
neonex module create <name>

# Example
neonex module create product
```

### Manual Creation

```
modules/product/
├── model.go
├── repository.go
├── service.go
├── controller.go
├── routes.go
├── seeder.go
├── di.go
└── module.json
```

---

## Module Structure

### module.json

```json
{
  "name": "product",
  "version": "1.0.0",
  "description": "Product management module",
  "dependencies": [],
  "routes": true,
  "migrations": true,
  "seeders": true
}
```

### Model

```go
// modules/product/model.go
package product

import "gorm.io/gorm"

type Product struct {
    gorm.Model
    Name        string  `gorm:"size:255;not null"`
    Description string  `gorm:"type:text"`
    Price       float64 `gorm:"type:decimal(10,2)"`
    Stock       int
    Active      bool `gorm:"default:true"`
}
```

### Repository

```go
// modules/product/repository.go
package product

import (
    "context"
    "neonexcore/pkg/database"
    "gorm.io/gorm"
)

type Repository interface {
    database.Repository[Product]
    FindActive(ctx context.Context) ([]*Product, error)
}

type repository struct {
    database.BaseRepository[Product]
}

func NewRepository(db *gorm.DB) Repository {
    return &repository{
        BaseRepository: database.NewBaseRepository[Product](db),
    }
}

func (r *repository) FindActive(ctx context.Context) ([]*Product, error) {
    var products []*Product
    err := r.DB().WithContext(ctx).
        Where("active = ?", true).
        Find(&products).Error
    return products, err
}
```

### Service

```go
// modules/product/service.go
package product

import (
    "context"
    "neonexcore/pkg/logger"
)

type Service interface {
    Create(ctx context.Context, req *CreateProductRequest) (*Product, error)
    GetAll(ctx context.Context) ([]*Product, error)
}

type service struct {
    repo   Repository
    logger logger.Logger
}

func NewService(repo Repository, logger logger.Logger) Service {
    return &service{
        repo:   repo,
        logger: logger,
    }
}

func (s *service) Create(ctx context.Context, req *CreateProductRequest) (*Product, error) {
    product := &Product{
        Name:        req.Name,
        Description: req.Description,
        Price:       req.Price,
        Stock:       req.Stock,
    }
    
    if err := s.repo.Create(ctx, product); err != nil {
        s.logger.Error("Failed to create product", logger.Fields{
            "error": err.Error(),
        })
        return nil, err
    }
    
    return product, nil
}
```

### Controller

```go
// modules/product/controller.go
package product

import (
    "github.com/gofiber/fiber/v2"
    "neonexcore/pkg/logger"
)

type controller struct {
    service Service
    logger  logger.Logger
}

func NewController(service Service, logger logger.Logger) *controller {
    return &controller{
        service: service,
        logger:  logger,
    }
}

func (ctrl *controller) Create(c *fiber.Ctx) error {
    var req CreateProductRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }
    
    product, err := ctrl.service.Create(c.Context(), &req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.Status(201).JSON(product)
}
```

### Routes

```go
// modules/product/routes.go
package product

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, ctrl *controller) {
    products := router.Group("/products")
    
    products.Post("/", ctrl.Create)
    products.Get("/", ctrl.GetAll)
    products.Get("/:id", ctrl.GetByID)
    products.Put("/:id", ctrl.Update)
    products.Delete("/:id", ctrl.Delete)
}
```

### Dependency Injection

```go
// modules/product/di.go
package product

import (
    "neonexcore/internal/core"
    "neonexcore/pkg/logger"
    "gorm.io/gorm"
)

func RegisterDependencies(c *core.Container) {
    // Repository
    core.RegisterTransient(c, "product.repository", func() Repository {
        db := core.Resolve[*gorm.DB](c, "database")
        return NewRepository(db)
    })
    
    // Service
    core.RegisterTransient(c, "product.service", func() Service {
        repo := core.Resolve[Repository](c, "product.repository")
        logger := core.Resolve[logger.Logger](c, "logger")
        return NewService(repo, logger)
    })
    
    // Controller
    core.RegisterTransient(c, "product.controller", func() *controller {
        service := core.Resolve[Service](c, "product.service")
        logger := core.Resolve[logger.Logger](c, "logger")
        return NewController(service, logger)
    })
}
```

---

## Module Registration

### Register in App

```go
// cmd/neonex/main.go
import "neonexcore/modules/product"

func main() {
    app := core.NewApp()
    
    // Register module
    product.RegisterDependencies(app.Container)
    
    // Register routes
    ctrl := core.Resolve[*product.controller](app.Container, "product.controller")
    product.RegisterRoutes(app.Router, ctrl)
    
    app.Run()
}
```

---

## Module Communication

### Service-to-Service

```go
// Module A calls Module B
type OrderService struct {
    productService product.Service
}

func (s *OrderService) CreateOrder(ctx context.Context, productID uint) error {
    // Call product service
    product, err := s.productService.GetByID(ctx, productID)
    if err != nil {
        return err
    }
    
    // Process order...
    return nil
}
```

### Events

```go
// Emit event
type ProductCreatedEvent struct {
    ProductID uint
    Name      string
}

func (s *service) Create(ctx context.Context, req *CreateProductRequest) (*Product, error) {
    product, err := s.repo.Create(ctx, &Product{...})
    if err != nil {
        return nil, err
    }
    
    // Emit event
    s.eventBus.Publish("product.created", ProductCreatedEvent{
        ProductID: product.ID,
        Name:      product.Name,
    })
    
    return product, nil
}
```

---

## Best Practices

### ✅ DO:
- Keep modules independent
- Use clear naming conventions
- Document public APIs
- Implement proper error handling
- Write comprehensive tests

### ❌ DON'T:
- Create circular dependencies
- Mix concerns across layers
- Expose internal implementation
- Skip validation

---

## Next Steps

- [**Middleware**](middleware.md) - Custom middleware
- [**Error Handling**](error-handling.md) - Error patterns
- [**Performance**](performance.md) - Optimization
- [**Module System**](../core-concepts/module-system.md) - Core concepts

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
