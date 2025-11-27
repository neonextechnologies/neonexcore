# Database Seeders

Populate your database with initial data using Neonex Core's seeding system.

---

## Overview

**Seeders** provide a way to populate your database with:

- ✅ Initial/default data
- ✅ Test data for development
- ✅ Demo data
- ✅ Lookup tables
- ✅ Configuration data

**Benefits:**
- Consistent data across environments
- Easy database reset for testing
- Reproducible development setup
- Quick demo environment setup

---

## Seeder Interface

### Definition

```go
// pkg/database/seeder.go
type Seeder interface {
    Seed(ctx context.Context, db *gorm.DB) error
}
```

### Implementation

```go
// modules/user/seeder.go
package user

import (
    "context"
    "gorm.io/gorm"
)

type UserSeeder struct{}

func NewUserSeeder() *UserSeeder {
    return &UserSeeder{}
}

func (s *UserSeeder) Seed(ctx context.Context, db *gorm.DB) error {
    users := []User{
        {
            Name:     "Admin User",
            Email:    "admin@example.com",
            Password: "$2a$10$...", // bcrypt hash
            Role:     "admin",
            Active:   true,
        },
        {
            Name:     "Test User",
            Email:    "test@example.com",
            Password: "$2a$10$...",
            Role:     "user",
            Active:   true,
        },
    }
    
    for _, user := range users {
        // Check if exists
        var existing User
        if err := db.Where("email = ?", user.Email).First(&existing).Error; err == nil {
            continue // Skip if exists
        }
        
        // Create user
        if err := db.Create(&user).Error; err != nil {
            return err
        }
    }
    
    return nil
}
```

---

## Creating Seeders

### Basic Seeder

```go
// modules/product/seeder.go
package product

import (
    "context"
    "gorm.io/gorm"
)

type ProductSeeder struct{}

func NewProductSeeder() *ProductSeeder {
    return &ProductSeeder{}
}

func (s *ProductSeeder) Seed(ctx context.Context, db *gorm.DB) error {
    products := []Product{
        {Name: "Laptop", Price: 999.99, Stock: 50},
        {Name: "Mouse", Price: 29.99, Stock: 200},
        {Name: "Keyboard", Price: 79.99, Stock: 150},
    }
    
    return db.Create(&products).Error
}
```

### Idempotent Seeder

Seeder that can run multiple times safely:

```go
func (s *ProductSeeder) Seed(ctx context.Context, db *gorm.DB) error {
    products := []Product{
        {Name: "Laptop", Price: 999.99, Stock: 50},
        {Name: "Mouse", Price: 29.99, Stock: 200},
    }
    
    for _, product := range products {
        // Use FirstOrCreate for idempotency
        result := db.Where("name = ?", product.Name).FirstOrCreate(&product)
        if result.Error != nil {
            return result.Error
        }
    }
    
    return nil
}
```

### Seeder with Dependencies

```go
type OrderSeeder struct {
    userRepo    UserRepository
    productRepo ProductRepository
}

func NewOrderSeeder(userRepo UserRepository, productRepo ProductRepository) *OrderSeeder {
    return &OrderSeeder{
        userRepo:    userRepo,
        productRepo: productRepo,
    }
}

func (s *OrderSeeder) Seed(ctx context.Context, db *gorm.DB) error {
    // Get existing users
    users, err := s.userRepo.FindAll(ctx)
    if err != nil {
        return err
    }
    
    // Get existing products
    products, err := s.productRepo.FindAll(ctx)
    if err != nil {
        return err
    }
    
    // Create orders
    for i, user := range users {
        order := Order{
            UserID: user.ID,
            Items: []OrderItem{
                {ProductID: products[i%len(products)].ID, Quantity: 1},
            },
            Total: products[i%len(products)].Price,
        }
        
        if err := db.Create(&order).Error; err != nil {
            return err
        }
    }
    
    return nil
}
```

---

## Running Seeders

### Automatic (Framework Handles)

Seeders are automatically run on application startup after migrations.

### Manual Execution

```go
// In your application code
func RunSeeders(db *gorm.DB) error {
    ctx := context.Background()
    
    // Create seeder manager
    manager := database.NewSeederManager(db)
    
    // Register seeders
    manager.Register(user.NewUserSeeder())
    manager.Register(product.NewProductSeeder())
    manager.Register(category.NewCategorySeeder())
    
    // Run all seeders
    return manager.Run(ctx)
}
```

### Command Line

```bash
# Create a seed command
go run main.go seed

# Or with specific seeder
go run main.go seed --only users
```

---

## Seeder Patterns

### Factory Pattern

```go
type UserFactory struct {
    faker *faker.Faker
}

func NewUserFactory() *UserFactory {
    return &UserFactory{
        faker: faker.New(),
    }
}

func (f *UserFactory) Create() *User {
    return &User{
        Name:     f.faker.Person().Name(),
        Email:    f.faker.Internet().Email(),
        Password: "$2a$10$...", // hashed "password"
        Role:     "user",
        Active:   true,
    }
}

func (f *UserFactory) CreateMany(count int) []User {
    users := make([]User, count)
    for i := 0; i < count; i++ {
        users[i] = *f.Create()
    }
    return users
}
```

**Usage:**
```go
func (s *UserSeeder) Seed(ctx context.Context, db *gorm.DB) error {
    factory := NewUserFactory()
    users := factory.CreateMany(100)
    return db.Create(&users).Error
}
```

### Bulk Insert

```go
func (s *ProductSeeder) Seed(ctx context.Context, db *gorm.DB) error {
    products := make([]Product, 10000)
    
    for i := 0; i < 10000; i++ {
        products[i] = Product{
            Name:  fmt.Sprintf("Product %d", i+1),
            Price: float64(10 + i%100),
            Stock: 100,
        }
    }
    
    // Batch insert (100 at a time)
    return db.CreateInBatches(products, 100).Error
}
```

### Conditional Seeding

```go
func (s *UserSeeder) Seed(ctx context.Context, db *gorm.DB) error {
    // Check environment
    env := os.Getenv("APP_ENV")
    
    if env == "production" {
        // Only seed essential data in production
        return s.seedAdminUser(db)
    }
    
    // Seed test data in development
    return s.seedTestUsers(db)
}

func (s *UserSeeder) seedAdminUser(db *gorm.DB) error {
    admin := User{
        Email: "admin@example.com",
        Role:  "admin",
    }
    return db.FirstOrCreate(&admin, "email = ?", admin.Email).Error
}

func (s *UserSeeder) seedTestUsers(db *gorm.DB) error {
    // Create 50 test users
    factory := NewUserFactory()
    users := factory.CreateMany(50)
    return db.Create(&users).Error
}
```

---

## Real-World Examples

### E-commerce Seeder

```go
type EcommerceSeeder struct {
    db *gorm.DB
}

func (s *EcommerceSeeder) Seed(ctx context.Context, db *gorm.DB) error {
    s.db = db
    
    // Seed in order
    if err := s.seedCategories(); err != nil {
        return err
    }
    
    if err := s.seedProducts(); err != nil {
        return err
    }
    
    if err := s.seedUsers(); err != nil {
        return err
    }
    
    if err := s.seedOrders(); err != nil {
        return err
    }
    
    return nil
}

func (s *EcommerceSeeder) seedCategories() error {
    categories := []Category{
        {Name: "Electronics", Slug: "electronics"},
        {Name: "Clothing", Slug: "clothing"},
        {Name: "Books", Slug: "books"},
    }
    return s.db.Create(&categories).Error
}

func (s *EcommerceSeeder) seedProducts() error {
    var electronics Category
    s.db.Where("slug = ?", "electronics").First(&electronics)
    
    products := []Product{
        {
            Name:       "Laptop",
            Price:      999.99,
            Stock:      50,
            CategoryID: electronics.ID,
        },
        {
            Name:       "Smartphone",
            Price:      599.99,
            Stock:      100,
            CategoryID: electronics.ID,
        },
    }
    return s.db.Create(&products).Error
}
```

### Lookup Tables

```go
type CountrySeeder struct{}

func (s *CountrySeeder) Seed(ctx context.Context, db *gorm.DB) error {
    countries := []Country{
        {Code: "US", Name: "United States"},
        {Code: "UK", Name: "United Kingdom"},
        {Code: "TH", Name: "Thailand"},
        {Code: "JP", Name: "Japan"},
    }
    
    for _, country := range countries {
        db.Where("code = ?", country.Code).FirstOrCreate(&country)
    }
    
    return nil
}
```

### Multi-tenant Data

```go
type TenantSeeder struct{}

func (s *TenantSeeder) Seed(ctx context.Context, db *gorm.DB) error {
    tenants := []Tenant{
        {Name: "Acme Corp", Domain: "acme.example.com"},
        {Name: "Tech Inc", Domain: "tech.example.com"},
    }
    
    for _, tenant := range tenants {
        // Create tenant
        if err := db.FirstOrCreate(&tenant, "domain = ?", tenant.Domain).Error; err != nil {
            return err
        }
        
        // Create tenant admin
        admin := User{
            TenantID: tenant.ID,
            Email:    fmt.Sprintf("admin@%s", tenant.Domain),
            Role:     "admin",
        }
        if err := db.FirstOrCreate(&admin, "email = ?", admin.Email).Error; err != nil {
            return err
        }
    }
    
    return nil
}
```

---

## Testing with Seeders

### Setup Test Database

```go
func SetupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    require.NoError(t, err)
    
    // Migrate
    db.AutoMigrate(&User{}, &Product{}, &Order{})
    
    // Seed
    ctx := context.Background()
    seeder := NewUserSeeder()
    err = seeder.Seed(ctx, db)
    require.NoError(t, err)
    
    return db
}

func TestUserService(t *testing.T) {
    db := SetupTestDB(t)
    defer os.Remove("test.db")
    
    // Test with seeded data
    var users []User
    db.Find(&users)
    assert.Greater(t, len(users), 0)
}
```

---

## Best Practices

### ✅ DO:

**1. Use FirstOrCreate for Idempotency**
```go
// Good: Can run multiple times
db.Where("email = ?", user.Email).FirstOrCreate(&user)
```

**2. Seed in Dependency Order**
```go
// Good: Categories before Products
seedCategories()
seedProducts() // Uses categories
```

**3. Use Transactions for Consistency**
```go
// Good: All or nothing
return db.Transaction(func(tx *gorm.DB) error {
    if err := seedUsers(tx); err != nil {
        return err
    }
    return seedOrders(tx)
})
```

**4. Check Environment**
```go
// Good: Different data per environment
if env == "production" {
    return seedMinimalData(db)
}
return seedTestData(db)
```

### ❌ DON'T:

**1. Seed Sensitive Data**
```go
// Bad: Real passwords/tokens
Password: "admin123"  // Use hashed passwords

// Good
Password: "$2a$10$..." // Bcrypt hash
```

**2. Create Duplicate Data**
```go
// Bad: No uniqueness check
db.Create(&user)  // Creates duplicate

// Good
db.FirstOrCreate(&user, "email = ?", user.Email)
```

**3. Ignore Errors**
```go
// Bad
db.Create(&users)

// Good
if err := db.Create(&users).Error; err != nil {
    return err
}
```

---

## Advanced Techniques

### CSV/JSON Import

```go
func (s *ProductSeeder) Seed(ctx context.Context, db *gorm.DB) error {
    // Read CSV file
    data, err := os.ReadFile("seeds/products.csv")
    if err != nil {
        return err
    }
    
    reader := csv.NewReader(bytes.NewReader(data))
    records, err := reader.ReadAll()
    if err != nil {
        return err
    }
    
    // Parse and insert
    for _, record := range records[1:] { // Skip header
        product := Product{
            Name:  record[0],
            Price: parseFloat(record[1]),
            Stock: parseInt(record[2]),
        }
        if err := db.Create(&product).Error; err != nil {
            return err
        }
    }
    
    return nil
}
```

### Faker Integration

```go
import "github.com/go-faker/faker/v4"

func (s *UserSeeder) Seed(ctx context.Context, db *gorm.DB) error {
    for i := 0; i < 100; i++ {
        user := User{
            Name:     faker.Name(),
            Email:    faker.Email(),
            Phone:    faker.Phonenumber(),
            Address:  faker.GetRealAddress().Address,
            Password: "$2a$10$...",
        }
        
        if err := db.Create(&user).Error; err != nil {
            return err
        }
    }
    
    return nil
}
```

---

## Troubleshooting

### Foreign Key Violations

**Problem:**
```
Error: foreign key constraint failed
```

**Solution:**
```go
// Seed in correct order
seedUsers()    // First
seedProducts() // Second
seedOrders()   // Last (depends on users and products)
```

### Duplicate Key Errors

**Problem:**
```
Error: duplicate key value
```

**Solution:**
```go
// Use FirstOrCreate or check existence
db.Where("email = ?", user.Email).FirstOrCreate(&user)
```

### Performance Issues

**Problem:** Seeding is slow

**Solution:**
```go
// Use batch inserts
db.CreateInBatches(users, 1000)

// Or disable transactions temporarily
db.Session(&gorm.Session{
    SkipDefaultTransaction: true,
}).CreateInBatches(users, 1000)
```

---

## Next Steps

- [**Migrations**](migrations.md) - Schema management
- [**Transactions**](transactions.md) - Transaction handling
- [**Repositories**](repositories.md) - Data access
- [**Testing**](../development/testing.md) - Testing with seeded data

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
