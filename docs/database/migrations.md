# Database Migrations

Automatic schema synchronization and migration management in Neonex Core.

---

## Overview

Neonex Core provides **automatic database migrations** powered by GORM's AutoMigrate. The framework:

- ✅ Automatically creates tables from models
- ✅ Adds missing columns
- ✅ Creates indexes
- ✅ Updates column types
- ✅ Preserves existing data

**Note:** AutoMigrate does NOT delete columns or tables. Manual intervention needed for destructive changes.

---

## How It Works

### Automatic Registration

Models are automatically registered and migrated on application startup:

```go
// Framework handles this automatically
// 1. Discover modules
// 2. Collect models from each module
// 3. Register with migrator
// 4. Run AutoMigrate
```

### Model Registration

Define models in your module:

```go
// modules/product/model.go
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

Framework discovers and migrates automatically on startup.

---

## GORM Tags

### Common Tags

| Tag | Description | Example |
|-----|-------------|---------|
| `primarykey` | Primary key | `gorm:"primarykey"` |
| `size` | Column size | `gorm:"size:100"` |
| `not null` | NOT NULL constraint | `gorm:"not null"` |
| `unique` | UNIQUE constraint | `gorm:"unique"` |
| `default` | Default value | `gorm:"default:0"` |
| `index` | Create index | `gorm:"index"` |
| `type` | Column type | `gorm:"type:text"` |
| `autoIncrement` | Auto increment | `gorm:"autoIncrement"` |
| `-` | Ignore field | `gorm:"-"` |
| `->` | Read-only | `gorm:"<-:false"` |
| `<-` | Write-only | `gorm:"<-:create"` |

### Examples

```go
type User struct {
    // Primary key
    ID uint `gorm:"primarykey"`
    
    // String with size
    Name string `gorm:"size:100;not null"`
    
    // Unique email
    Email string `gorm:"size:255;unique;not null"`
    
    // Default value
    Role string `gorm:"size:50;default:'user'"`
    
    // Index
    Status string `gorm:"size:20;index"`
    
    // Composite index
    Category string `gorm:"size:50;index:idx_category_price"`
    Price    float64 `gorm:"index:idx_category_price"`
    
    // JSON column
    Metadata string `gorm:"type:json"`
    
    // Ignore field
    TempData string `gorm:"-"`
    
    // Read-only
    CreatedAt time.Time `gorm:"<-:create"`
    
    // Soft delete
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

---

## Migration Process

### Startup Flow

```
1. Application Start
   ↓
2. Load Modules
   ↓
3. Collect Models
   ↓
4. Register with Migrator
   ↓
5. Run AutoMigrate
   ↓
6. Log Results
```

### Console Output

```bash
🔄 Running auto-migration for 5 models...
✅ Migrated: users
✅ Migrated: products
✅ Migrated: orders
✅ Migrated: categories
✅ Migrated: reviews
✅ Database migration completed
```

---

## Model Definitions

### Basic Model

```go
type Product struct {
    ID          uint      `gorm:"primarykey" json:"id"`
    Name        string    `gorm:"size:100;not null" json:"name"`
    Price       float64   `gorm:"not null" json:"price"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### With Soft Delete

```go
import "gorm.io/gorm"

type Product struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    Name      string         `gorm:"size:100;not null" json:"name"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
```

### With Relationships

```go
// One-to-Many
type User struct {
    ID     uint    `gorm:"primarykey"`
    Name   string  `gorm:"size:100"`
    Orders []Order `gorm:"foreignKey:UserID"`
}

type Order struct {
    ID      uint    `gorm:"primarykey"`
    UserID  uint    `gorm:"not null;index"`
    User    User    `gorm:"constraint:OnDelete:CASCADE"`
    Total   float64 `gorm:"not null"`
}

// Many-to-Many
type Product struct {
    ID         uint       `gorm:"primarykey"`
    Name       string     `gorm:"size:100"`
    Categories []Category `gorm:"many2many:product_categories;"`
}

type Category struct {
    ID       uint      `gorm:"primarykey"`
    Name     string    `gorm:"size:50"`
    Products []Product `gorm:"many2many:product_categories;"`
}
```

### Custom Table Name

```go
type Product struct {
    ID   uint   `gorm:"primarykey"`
    Name string `gorm:"size:100"`
}

func (Product) TableName() string {
    return "products_v2"
}
```

---

## Manual Migration Control

### Skip Auto-Migration

```go
// main.go
func main() {
    app := core.NewApp()
    
    // Skip migration
    os.Setenv("SKIP_MIGRATION", "true")
    
    app.Run()
}
```

### Manual Migration

```go
func (a *App) ManualMigrate() error {
    db := config.DB.GetDB()
    
    // Migrate specific models
    return db.AutoMigrate(
        &User{},
        &Product{},
        &Order{},
    )
}
```

### Check Table Exists

```go
func (a *App) CheckTable() {
    db := config.DB.GetDB()
    
    if db.Migrator().HasTable(&User{}) {
        fmt.Println("✅ Users table exists")
    } else {
        fmt.Println("❌ Users table missing")
    }
}
```

---

## Schema Changes

### Adding New Field

**Before:**
```go
type Product struct {
    ID    uint   `gorm:"primarykey"`
    Name  string `gorm:"size:100"`
    Price float64
}
```

**After:**
```go
type Product struct {
    ID          uint   `gorm:"primarykey"`
    Name        string `gorm:"size:100"`
    Price       float64
    Description string `gorm:"size:500"` // New field
}
```

**Result:** AutoMigrate adds `description` column automatically.

### Modifying Field

**Before:**
```go
Name string `gorm:"size:50"`
```

**After:**
```go
Name string `gorm:"size:100"` // Increased size
```

**Result:** AutoMigrate updates column type.

### Adding Index

**Before:**
```go
Email string `gorm:"size:255"`
```

**After:**
```go
Email string `gorm:"size:255;index"` // Add index
```

**Result:** AutoMigrate creates index.

---

## Limitations

### What AutoMigrate Does NOT Do

❌ **Delete Columns**
```go
// Removing a field does NOT delete the column
// Manual SQL required
```

❌ **Delete Tables**
```go
// Removing a model does NOT drop the table
// Manual intervention needed
```

❌ **Rename Columns**
```go
// Renaming requires manual migration
ALTER TABLE products RENAME COLUMN old_name TO new_name;
```

❌ **Complex Constraints**
```go
// CHECK constraints need manual SQL
ALTER TABLE products ADD CONSTRAINT price_positive CHECK (price > 0);
```

### Manual SQL for Destructive Changes

```go
func DropColumn(db *gorm.DB) error {
    return db.Exec("ALTER TABLE products DROP COLUMN old_field").Error
}

func RenameColumn(db *gorm.DB) error {
    return db.Exec("ALTER TABLE products RENAME COLUMN old_name TO new_name").Error
}

func DropTable(db *gorm.DB) error {
    return db.Migrator().DropTable(&OldModel{})
}
```

---

## Migration Strategies

### Development

```go
// Reset database (drop and recreate)
func ResetDatabase(db *gorm.DB) error {
    migrator := database.NewMigrator(db)
    
    // Drop all tables
    if err := migrator.DropTables(); err != nil {
        return err
    }
    
    // Recreate tables
    return migrator.AutoMigrate()
}
```

### Staging/Production

```go
// Safe migration (only additive changes)
func SafeMigrate(db *gorm.DB) error {
    // Backup first
    if err := BackupDatabase(); err != nil {
        return err
    }
    
    // Run migration
    migrator := database.NewMigrator(db)
    return migrator.AutoMigrate()
}
```

---

## Version Control

### Migration History

Track schema versions manually:

```go
type Migration struct {
    ID        uint      `gorm:"primarykey"`
    Version   string    `gorm:"size:50;not null;unique"`
    AppliedAt time.Time `gorm:"not null"`
}

func RecordMigration(db *gorm.DB, version string) error {
    migration := &Migration{
        Version:   version,
        AppliedAt: time.Now(),
    }
    return db.Create(migration).Error
}
```

### Schema Snapshots

```bash
# Export schema
pg_dump -s mydb > schema_v1.sql

# Compare schemas
diff schema_v1.sql schema_v2.sql
```

---

## Testing Migrations

### Test in Isolated Database

```go
func TestMigration(t *testing.T) {
    // Create test database
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    require.NoError(t, err)
    
    // Run migration
    migrator := database.NewMigrator(db)
    migrator.RegisterModels(&User{}, &Product{})
    err = migrator.AutoMigrate()
    require.NoError(t, err)
    
    // Verify tables
    assert.True(t, db.Migrator().HasTable(&User{}))
    assert.True(t, db.Migrator().HasTable(&Product{}))
    
    // Cleanup
    os.Remove("test.db")
}
```

---

## Best Practices

### ✅ DO:

**1. Use Soft Deletes**
```go
// Good: Preserve data
type Model struct {
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

**2. Add Indexes for Foreign Keys**
```go
// Good: Query performance
UserID uint `gorm:"not null;index"`
```

**3. Set NOT NULL for Required Fields**
```go
// Good: Data integrity
Email string `gorm:"size:255;not null"`
```

**4. Use Appropriate Sizes**
```go
// Good: Storage optimization
Name string `gorm:"size:100"` // Not size:255 for everything
```

### ❌ DON'T:

**1. Remove Fields Without Backup**
```go
// Bad: Data loss risk
// type User struct {
//     ImportantField string // Removing this
// }
```

**2. Change Primary Key Type**
```go
// Bad: Breaking change
// ID uint   → ID string
```

**3. Ignore Migration Errors**
```go
// Bad
migrator.AutoMigrate() // No error check

// Good
if err := migrator.AutoMigrate(); err != nil {
    log.Fatal(err)
}
```

---

## Troubleshooting

### Migration Failed

**Problem:**
```
Error: failed to migrate model
```

**Solutions:**
1. Check GORM tags syntax
2. Verify database permissions
3. Check for conflicting constraints
4. Review database logs

### Column Already Exists

**Problem:**
```
Error: column already exists
```

**Solution:**
- AutoMigrate is idempotent
- Safe to run multiple times
- Check for manual schema changes

### Type Mismatch

**Problem:**
```
Error: cannot convert column type
```

**Solution:**
```sql
-- Manual type conversion
ALTER TABLE products ALTER COLUMN price TYPE DECIMAL(10,2);
```

---

## Advanced Topics

### Custom Migrator

```go
type CustomMigrator struct {
    *database.Migrator
}

func (m *CustomMigrator) AutoMigrate() error {
    // Pre-migration tasks
    m.BackupSchema()
    
    // Run migration
    if err := m.Migrator.AutoMigrate(); err != nil {
        return err
    }
    
    // Post-migration tasks
    m.CreateCustomIndexes()
    m.AddConstraints()
    
    return nil
}
```

### Migration Callbacks

```go
func (u *User) AfterMigrate(tx *gorm.DB) error {
    // Create default admin user
    return tx.FirstOrCreate(&User{
        Email: "admin@example.com",
        Role:  "admin",
    }).Error
}
```

---

## Next Steps

- [**Seeders**](seeders.md) - Populate initial data
- [**Repositories**](repositories.md) - Data access patterns
- [**Transactions**](transactions.md) - Transaction management
- [**Configuration**](configuration.md) - Database setup

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
