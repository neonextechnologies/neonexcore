# neonex module

Manage application modules: create new modules, list existing ones, and more.

---

## Synopsis

```bash
neonex module <command> [arguments] [flags]
```

Module management commands for scaffolding and organizing your application features.

**Available Commands:**
- `create` - Generate a new module with complete CRUD structure
- `list` - Display all modules in the project

---

## neonex module create

Generate a new feature module with full CRUD implementation.

### Usage

```bash
neonex module create <module-name>
```

### What It Creates

Generates **9 files** in `modules/<module-name>/`:

| File | Purpose | Lines |
|------|---------|-------|
| `model.go` | GORM model with database tags | ~30 |
| `repository.go` | Data access layer with interface | ~80 |
| `service.go` | Business logic layer | ~100 |
| `controller.go` | HTTP request handlers | ~120 |
| `routes.go` | API route definitions | ~20 |
| `di.go` | Dependency injection setup | ~25 |
| `seeder.go` | Database seeding | ~40 |
| `{module}.go` | Module entry point | ~45 |
| `module.json` | Module metadata | ~10 |

**Total:** ~470 lines of production-ready code!

### Example

```bash
# Create a product module
neonex module create product
```

**Output:**
```
Creating module: product
✓ Created model.go
✓ Created repository.go
✓ Created service.go
✓ Created controller.go
✓ Created routes.go
✓ Created di.go
✓ Created seeder.go
✓ Created product.go
✓ Created module.json

Module created successfully!

Next steps:
  1. Edit modules/product/model.go to define your schema
  2. Run: neonex serve --hot
  3. Test: curl http://localhost:8080/product/

Endpoints:
  GET    /product/       List all items
  GET    /product/:id    Get by ID
  POST   /product/       Create new
  PUT    /product/:id    Update
  DELETE /product/:id    Delete
```

---

## Generated Files

### model.go - GORM Model

```go
package product

import (
    "gorm.io/gorm"
    "time"
)

type Product struct {
    ID          uint           `gorm:"primarykey" json:"id"`
    Name        string         `gorm:"size:100;not null" json:"name"`
    Description string         `gorm:"size:500" json:"description"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Product) TableName() string {
    return "products"
}
```

**Features:**
- ✅ GORM tags for validation
- ✅ Soft delete support
- ✅ JSON serialization
- ✅ Custom table name
- ✅ Timestamps

### repository.go - Data Access

```go
package product

import (
    "github.com/YOUR_USERNAME/neonexcore/pkg/database"
    "gorm.io/gorm"
)

type Repository interface {
    Create(product *Product) error
    FindByID(id uint) (*Product, error)
    FindAll() ([]Product, error)
    Update(product *Product) error
    Delete(id uint) error
}

type repository struct {
    *database.Repository[Product]
}

func NewRepository(db *gorm.DB) Repository {
    return &repository{
        Repository: database.NewRepository[Product](db),
    }
}

func (r *repository) Create(product *Product) error {
    return r.Repository.Create(product)
}

func (r *repository) FindByID(id uint) (*Product, error) {
    return r.Repository.FindByID(id)
}

func (r *repository) FindAll() ([]Product, error) {
    return r.Repository.FindAll()
}

func (r *repository) Update(product *Product) error {
    return r.Repository.Update(product)
}

func (r *repository) Delete(id uint) error {
    return r.Repository.Delete(id)
}
```

**Features:**
- ✅ Interface-based design
- ✅ Generic repository pattern
- ✅ Clean separation of concerns
- ✅ Easy to mock for testing

### service.go - Business Logic

```go
package product

import (
    "errors"
    "gorm.io/gorm"
)

type Service interface {
    CreateProduct(product *Product) error
    GetProduct(id uint) (*Product, error)
    GetAllProducts() ([]Product, error)
    UpdateProduct(id uint, product *Product) error
    DeleteProduct(id uint) error
}

type service struct {
    repo Repository
}

func NewService(repo Repository) Service {
    return &service{repo: repo}
}

func (s *service) CreateProduct(product *Product) error {
    if product.Name == "" {
        return errors.New("name is required")
    }
    return s.repo.Create(product)
}

func (s *service) GetProduct(id uint) (*Product, error) {
    return s.repo.FindByID(id)
}

func (s *service) GetAllProducts() ([]Product, error) {
    return s.repo.FindAll()
}

func (s *service) UpdateProduct(id uint, product *Product) error {
    existing, err := s.repo.FindByID(id)
    if err != nil {
        return err
    }
    
    product.ID = existing.ID
    return s.repo.Update(product)
}

func (s *service) DeleteProduct(id uint) error {
    return s.repo.Delete(id)
}
```

**Features:**
- ✅ Input validation
- ✅ Error handling
- ✅ Business rules enforcement
- ✅ Testable logic

### controller.go - HTTP Handlers

```go
package product

import (
    "strconv"
    "github.com/gofiber/fiber/v2"
)

type Controller struct {
    service Service
}

func NewController(service Service) *Controller {
    return &Controller{service: service}
}

func (c *Controller) GetAll(ctx *fiber.Ctx) error {
    products, err := c.service.GetAllProducts()
    if err != nil {
        return ctx.Status(500).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    return ctx.JSON(products)
}

func (c *Controller) GetByID(ctx *fiber.Ctx) error {
    id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
    if err != nil {
        return ctx.Status(400).JSON(fiber.Map{
            "error": "Invalid ID",
        })
    }
    
    product, err := c.service.GetProduct(uint(id))
    if err != nil {
        return ctx.Status(404).JSON(fiber.Map{
            "error": "Product not found",
        })
    }
    
    return ctx.JSON(product)
}

func (c *Controller) Create(ctx *fiber.Ctx) error {
    var product Product
    if err := ctx.BodyParser(&product); err != nil {
        return ctx.Status(400).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    
    if err := c.service.CreateProduct(&product); err != nil {
        return ctx.Status(400).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return ctx.Status(201).JSON(product)
}

func (c *Controller) Update(ctx *fiber.Ctx) error {
    id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
    if err != nil {
        return ctx.Status(400).JSON(fiber.Map{
            "error": "Invalid ID",
        })
    }
    
    var product Product
    if err := ctx.BodyParser(&product); err != nil {
        return ctx.Status(400).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    
    if err := c.service.UpdateProduct(uint(id), &product); err != nil {
        return ctx.Status(400).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return ctx.JSON(product)
}

func (c *Controller) Delete(ctx *fiber.Ctx) error {
    id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
    if err != nil {
        return ctx.Status(400).JSON(fiber.Map{
            "error": "Invalid ID",
        })
    }
    
    if err := c.service.DeleteProduct(uint(id)); err != nil {
        return ctx.Status(400).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    return ctx.SendStatus(204)
}
```

**Features:**
- ✅ RESTful API design
- ✅ Proper HTTP status codes
- ✅ Request validation
- ✅ Error responses
- ✅ JSON serialization

### routes.go - API Routes

```go
package product

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, controller *Controller) {
    group := router.Group("/product")
    
    group.Get("/", controller.GetAll)
    group.Get("/:id", controller.GetByID)
    group.Post("/", controller.Create)
    group.Put("/:id", controller.Update)
    group.Delete("/:id", controller.Delete)
}
```

**Endpoints:**
- `GET /product/` - List all
- `GET /product/:id` - Get one
- `POST /product/` - Create
- `PUT /product/:id` - Update
- `DELETE /product/:id` - Delete

### di.go - Dependency Injection

```go
package product

import (
    "github.com/YOUR_USERNAME/neonexcore/internal/core"
    "gorm.io/gorm"
)

func RegisterDependencies(c *core.Container) error {
    // Register Repository
    c.Singleton(func() Repository {
        db := c.MustResolve(&gorm.DB{}).(*gorm.DB)
        return NewRepository(db)
    })
    
    // Register Service
    c.Singleton(func() Service {
        repo := c.MustResolve((*Repository)(nil)).(Repository)
        return NewService(repo)
    })
    
    // Register Controller
    c.Singleton(func() *Controller {
        service := c.MustResolve((*Service)(nil)).(Service)
        return NewController(service)
    })
    
    return nil
}
```

**Features:**
- ✅ Singleton services
- ✅ Automatic dependency resolution
- ✅ Type-safe container

### seeder.go - Database Seeding

```go
package product

import (
    "gorm.io/gorm"
)

type ProductSeeder struct {
    db *gorm.DB
}

func NewSeeder(db *gorm.DB) *ProductSeeder {
    return &ProductSeeder{db: db}
}

func (s *ProductSeeder) Seed() error {
    products := []Product{
        {Name: "Product 1", Description: "Description 1"},
        {Name: "Product 2", Description: "Description 2"},
        {Name: "Product 3", Description: "Description 3"},
    }
    
    for _, product := range products {
        if err := s.db.Create(&product).Error; err != nil {
            return err
        }
    }
    
    return nil
}
```

### product.go - Module Entry Point

```go
package product

import (
    "github.com/YOUR_USERNAME/neonexcore/internal/core"
    "github.com/YOUR_USERNAME/neonexcore/internal/module"
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
)

type ProductModule struct{}

func New() module.Module {
    return &ProductModule{}
}

func (m *ProductModule) Register(c *core.Container) error {
    return RegisterDependencies(c)
}

func (m *ProductModule) Boot(c *core.Container) error {
    return nil
}

func (m *ProductModule) Routes(router fiber.Router, c *core.Container) error {
    controller := c.MustResolve((*Controller)(nil)).(*Controller)
    RegisterRoutes(router, controller)
    return nil
}

func (m *ProductModule) Models() []interface{} {
    return []interface{}{&Product{}}
}

func (m *ProductModule) Seeders(db *gorm.DB) []interface{} {
    return []interface{}{NewSeeder(db)}
}
```

### module.json - Metadata

```json
{
  "name": "product",
  "version": "1.0.0",
  "description": "Product management module",
  "dependencies": []
}
```

---

## neonex module list

Display all modules in the project.

### Usage

```bash
neonex module list
```

**Output:**
```
Discovered modules:

  user (v1.0.0)
    Path: modules/user
    Description: User management module
    
  product (v1.0.0)
    Path: modules/product
    Description: Product management module
    
  order (v1.0.0)
    Path: modules/order
    Description: Order management module

Total: 3 modules
```

### Details Shown

For each module:
- **Name** - Module identifier
- **Version** - From module.json
- **Path** - File system location
- **Description** - Module purpose
- **Dependencies** - Required modules

---

## Module Architecture

### Layered Structure

```
Controller → Service → Repository → Database
   ↓          ↓          ↓
 HTTP      Business    Data
Layer      Logic      Access
```

**Flow:**
1. **Controller** receives HTTP request
2. **Service** processes business logic
3. **Repository** accesses database
4. Response flows back up

### Dependency Flow

```
Module Entry (product.go)
    ↓
Dependency Injection (di.go)
    ↓
Controller ← Service ← Repository
    ↓
Routes (routes.go)
```

---

## Customization

### Modify Generated Code

After generation, customize to your needs:

```bash
# Generate module
neonex module create blog

# Edit model
nano modules/blog/model.go
```

**Example: Add fields**
```go
type Blog struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    Title     string    `gorm:"size:200;not null" json:"title"`
    Content   string    `gorm:"type:text" json:"content"`
    Author    string    `gorm:"size:100" json:"author"`
    Published bool      `json:"published"`
    Tags      []string  `gorm:"type:json" json:"tags"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Add Custom Methods

**Repository:**
```go
func (r *repository) FindByAuthor(author string) ([]Blog, error) {
    var blogs []Blog
    err := r.DB.Where("author = ?", author).Find(&blogs).Error
    return blogs, err
}
```

**Service:**
```go
func (s *service) GetBlogsByAuthor(author string) ([]Blog, error) {
    return s.repo.FindByAuthor(author)
}
```

**Controller:**
```go
func (c *Controller) GetByAuthor(ctx *fiber.Ctx) error {
    author := ctx.Query("author")
    blogs, err := c.service.GetBlogsByAuthor(author)
    if err != nil {
        return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return ctx.JSON(blogs)
}
```

**Routes:**
```go
group.Get("/by-author", controller.GetByAuthor)
```

---

## Examples

### Example 1: E-commerce Modules

```bash
# Create core modules
neonex module create product
neonex module create category
neonex module create cart
neonex module create order
neonex module create customer

# Start server
neonex serve --hot

# Test
curl http://localhost:8080/product/
curl http://localhost:8080/category/
```

### Example 2: Blog System

```bash
# Generate modules
neonex module create blog
neonex module create comment
neonex module create tag

# Customize blog module
# Add relationships in model.go
```

### Example 3: Multi-tenant SaaS

```bash
# Tenant modules
neonex module create tenant
neonex module create subscription
neonex module create billing
neonex module create analytics
```

---

## Best Practices

### Naming Conventions

✅ **Good:**
```bash
neonex module create user
neonex module create product
neonex module create order_item
```

❌ **Avoid:**
```bash
neonex module create User        # CamelCase
neonex module create product-details  # hyphens
neonex module create PRODUCT     # uppercase
```

### Module Organization

```
modules/
├── core/              # Core business modules
│   ├── user/
│   ├── auth/
│   └── profile/
├── ecommerce/        # Feature-specific
│   ├── product/
│   ├── cart/
│   └── order/
└── admin/            # Admin modules
    ├── dashboard/
    └── settings/
```

### Keep Modules Independent

❌ **Don't:**
```go
// In modules/product/service.go
import "github.com/myapp/modules/user"  // Tight coupling
```

✅ **Do:**
```go
// Use shared interfaces in pkg/
import "github.com/myapp/pkg/auth"
```

---

## Troubleshooting

### Module Not Detected

**Problem:**
```bash
Error: Module 'product' not found
```

**Solution:**
```bash
# Check structure
ls modules/product/

# Verify module.json exists
cat modules/product/module.json

# Check product.go has New() function
grep "func New()" modules/product/product.go
```

### Compilation Errors

**Problem:**
```bash
Error: undefined: Product
```

**Solution:**
```bash
# Re-run module creation
neonex module create product

# Or fix imports
go mod tidy
```

### Routes Not Working

**Problem:**
```bash
404 Not Found on /product/
```

**Solution:**
1. Check `RegisterRoutes()` is called
2. Verify module is discovered
3. Check server logs
4. Restart server

---

## Next Steps

- [**Module System**](../core-concepts/module-system.md) - Deep dive into modules
- [**Dependency Injection**](../core-concepts/dependency-injection.md) - Master DI
- [**Repository Pattern**](../core-concepts/repository-pattern.md) - Data access
- [**Service Layer**](../core-concepts/service-layer.md) - Business logic

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
