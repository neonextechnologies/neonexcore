# Database Overview

Learn about Neonex Core's database layer powered by GORM with multi-driver support and advanced ORM features.

---

## Overview

Neonex Core provides a **comprehensive database layer** built on top of GORM, offering:

- ✅ **Multi-Database Support** - SQLite, PostgreSQL, MySQL, Turso
- ✅ **Generic Repository Pattern** - Type-safe CRUD operations
- ✅ **Auto-Migration** - Automatic schema synchronization
- ✅ **Transaction Management** - ACID-compliant transactions
- ✅ **Database Seeding** - Initial data management
- ✅ **Connection Pooling** - Optimized performance
- ✅ **Zero Configuration** - Works with SQLite out of the box

---

## Supported Databases

### SQLite (Default)

**Perfect for:**
- Development
- Prototyping
- Small applications
- Embedded systems

**Configuration:**
```bash
DB_DRIVER=sqlite
DB_NAME=neonex.db
```

**Pros:**
- ✅ Zero configuration
- ✅ No server required
- ✅ File-based
- ✅ Fast for small datasets

**Cons:**
- ❌ Not for high concurrency
- ❌ Limited scalability
- ❌ Single write at a time

### PostgreSQL

**Perfect for:**
- Production applications
- High concurrency
- Complex queries
- Large datasets

**Configuration:**
```bash
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=neonex
DB_USER=postgres
DB_PASSWORD=secret
DB_SSLMODE=disable
```

**Pros:**
- ✅ ACID compliant
- ✅ Advanced features
- ✅ High performance
- ✅ Excellent concurrency

### MySQL

**Perfect for:**
- Web applications
- High read workloads
- Replication setups

**Configuration:**
```bash
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_NAME=neonex
DB_USER=root
DB_PASSWORD=secret
DB_CHARSET=utf8mb4
DB_PARSE_TIME=true
```

**Pros:**
- ✅ Mature and stable
- ✅ Wide adoption
- ✅ Good performance
- ✅ Replication support

### Turso (LibSQL)

**Perfect for:**
- Edge computing
- Serverless applications
- Global distribution
- Low latency

**Configuration:**
```bash
DB_DRIVER=turso
DB_URL=libsql://your-db.turso.io
DB_TOKEN=your-auth-token
```

**Pros:**
- ✅ Edge deployment
- ✅ SQLite compatible
- ✅ Global replication
- ✅ Low latency

---

## Architecture

### Database Layer Structure

```
Application
    ↓
Service Layer (Business Logic)
    ↓
Repository Layer (Data Access)
    ↓
GORM (ORM)
    ↓
Database Driver (SQLite/PostgreSQL/MySQL/Turso)
    ↓
Database
```

### Components

| Component | Purpose | Location |
|-----------|---------|----------|
| **Connection Manager** | Database connections | `internal/config/database.go` |
| **Migrator** | Schema management | `pkg/database/migrator.go` |
| **Repository** | Generic CRUD | `pkg/database/repository.go` |
| **Transaction Manager** | ACID transactions | `pkg/database/transaction.go` |
| **Seeder** | Initial data | `pkg/database/seeder.go` |

---

## Quick Start

### 1. Configure Database

```bash
# .env file
DB_DRIVER=sqlite
DB_NAME=app.db
```

### 2. Define Model

```go
// modules/product/model.go
package product

import (
    "gorm.io/gorm"
    "time"
)

type Product struct {
    ID          uint           `gorm:"primarykey" json:"id"`
    Name        string         `gorm:"size:100;not null" json:"name"`
    Description string         `gorm:"size:500" json:"description"`
    Price       float64        `gorm:"not null" json:"price"`
    Stock       int            `gorm:"default:0" json:"stock"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
```

### 3. Create Repository

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
}

type repository struct {
    *database.BaseRepository[Product]
}

func NewRepository(db *gorm.DB) Repository {
    return &repository{
        BaseRepository: database.NewBaseRepository[Product](db),
    }
}
```

### 4. Use in Service

```go
// modules/product/service.go
package product

func (s *service) CreateProduct(ctx context.Context, product *Product) error {
    return s.repo.Create(ctx, product)
}

func (s *service) GetProduct(ctx context.Context, id uint) (*Product, error) {
    return s.repo.FindByID(ctx, id)
}
```

### 5. Auto-Migration

```go
// Automatically handled by framework
// Models registered in module.go are auto-migrated on startup
```

---

## Features

### Generic Repository

Type-safe CRUD operations with Go generics:

```go
// Create
product := &Product{Name: "Laptop", Price: 999.99}
err := repo.Create(ctx, product)

// Read
product, err := repo.FindByID(ctx, 1)
products, err := repo.FindAll(ctx)

// Update
product.Price = 899.99
err := repo.Update(ctx, product)

// Delete
err := repo.Delete(ctx, 1)

// Query
products, err := repo.FindByCondition(ctx, "price > ?", 500)

// Paginate
products, total, err := repo.Paginate(ctx, 1, 20)
```

### Auto-Migration

Automatic schema synchronization:

```go
// Framework handles this automatically
// Define models in your modules
// They're auto-migrated on application start

type Product struct {
    ID        uint      `gorm:"primarykey"`
    Name      string    `gorm:"size:100;not null"`
    Price     float64   `gorm:"not null"`
    CreatedAt time.Time
}
```

### Transaction Support

ACID-compliant transactions:

```go
err := txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
    // Create order
    if err := orderRepo.Create(ctx, order); err != nil {
        return err
    }
    
    // Update inventory
    if err := productRepo.UpdateStock(ctx, productID, -quantity); err != nil {
        return err
    }
    
    return nil
})
```

### Database Seeding

Populate initial data:

```go
type ProductSeeder struct {
    db *gorm.DB
}

func (s *ProductSeeder) Seed(ctx context.Context, db *gorm.DB) error {
    products := []Product{
        {Name: "Product 1", Price: 99.99},
        {Name: "Product 2", Price: 149.99},
    }
    
    return db.Create(&products).Error
}
```

---

## GORM Features

### Relationships

```go
// One-to-Many
type User struct {
    ID     uint
    Orders []Order
}

type Order struct {
    ID     uint
    UserID uint
    User   User
}

// Many-to-Many
type Product struct {
    ID         uint
    Categories []Category `gorm:"many2many:product_categories;"`
}

type Category struct {
    ID       uint
    Products []Product `gorm:"many2many:product_categories;"`
}
```

### Hooks

```go
// Before/After Create
func (p *Product) BeforeCreate(tx *gorm.DB) error {
    p.Slug = generateSlug(p.Name)
    return nil
}

func (p *Product) AfterCreate(tx *gorm.DB) error {
    // Send notification
    return nil
}
```

### Scopes

```go
// Define scope
func Active(db *gorm.DB) *gorm.DB {
    return db.Where("active = ?", true)
}

// Use scope
var users []User
db.Scopes(Active).Find(&users)
```

### Raw SQL

```go
// Raw query
var result struct {
    Total int
    Avg   float64
}

db.Raw("SELECT COUNT(*) as total, AVG(price) as avg FROM products").
    Scan(&result)

// Exec
db.Exec("UPDATE products SET stock = stock - 1 WHERE id = ?", productID)
```

---

## Performance Optimization

### Connection Pooling

```go
// internal/config/database.go
sqlDB, _ := db.DB()

// Set max idle connections
sqlDB.SetMaxIdleConns(10)

// Set max open connections
sqlDB.SetMaxOpenConns(100)

// Set max lifetime
sqlDB.SetConnMaxLifetime(time.Hour)
```

### Indexing

```go
type Product struct {
    ID       uint   `gorm:"primarykey"`
    Name     string `gorm:"size:100;not null;index"`
    Category string `gorm:"size:50;index"`
    Price    float64 `gorm:"index"`
}
```

### Preloading

```go
// Eager loading
var users []User
db.Preload("Orders").Find(&users)

// Nested preloading
db.Preload("Orders.Items").Find(&users)

// Conditional preloading
db.Preload("Orders", "status = ?", "completed").Find(&users)
```

### Select Fields

```go
// Select specific fields
var products []Product
db.Select("id", "name", "price").Find(&products)

// Omit fields
db.Omit("Password").Find(&users)
```

---

## Best Practices

### ✅ DO:

**1. Use Transactions for Multiple Operations**
```go
// Good: All or nothing
txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
    // Multiple operations
})
```

**2. Use Context**
```go
// Good: Cancellation support
repo.FindAll(ctx)
```

**3. Handle Soft Deletes**
```go
// Good: Use DeletedAt
type Model struct {
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

**4. Use Prepared Statements**
```go
// Good: Parameterized queries
db.Where("name = ?", name).Find(&products)
```

### ❌ DON'T:

**1. Don't Use String Concatenation**
```go
// Bad: SQL injection risk
db.Where("name = '" + name + "'").Find(&products)
```

**2. Don't Ignore Errors**
```go
// Bad
db.Create(&product)

// Good
if err := db.Create(&product).Error; err != nil {
    return err
}
```

**3. Don't Use Global DB**
```go
// Bad: Global state
var DB *gorm.DB

// Good: Dependency injection
func NewRepository(db *gorm.DB) Repository
```

---

## Migration from Other ORMs

### From Raw SQL

```go
// Before (raw SQL)
rows, err := db.Query("SELECT * FROM users WHERE active = ?", true)

// After (GORM)
var users []User
db.Where("active = ?", true).Find(&users)
```

### From Other Go ORMs

Most ORMs follow similar patterns. Key differences:

- **sqlx** → Use GORM's struct scanning
- **sqlc** → GORM generates methods automatically
- **ent** → Similar schema definition, different syntax

---

## Troubleshooting

### Connection Issues

**Problem:** Cannot connect to database

**Solution:**
```bash
# Check config
cat .env | grep DB_

# Test connection
go run main.go  # Check logs
```

### Migration Failures

**Problem:** Migration failed

**Solution:**
```bash
# Check model tags
# Fix any syntax errors in GORM tags
# Run migration manually
```

### Performance Issues

**Problem:** Slow queries

**Solution:**
- Add indexes
- Use pagination
- Enable query logging
- Use explain analyze

---

## Next Steps

- [**Configuration**](configuration.md) - Database configuration details
- [**Migrations**](migrations.md) - Schema management
- [**Seeders**](seeders.md) - Initial data setup
- [**Transactions**](transactions.md) - Transaction handling
- [**Repositories**](repositories.md) - Advanced repository patterns

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
