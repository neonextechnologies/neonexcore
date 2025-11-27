# Performance Optimization

Techniques for optimizing Neonex Core applications.

---

## Database Optimization

### Use Indexes

```go
type User struct {
    ID    uint   `gorm:"primarykey"`
    Email string `gorm:"index:idx_email,unique"`
    Name  string `gorm:"index"`
}
```

### Preload Relationships

```go
db.Preload("Orders").Preload("Profile").Find(&users)
```

### Batch Operations

```go
db.CreateInBatches(users, 1000)
```

## HTTP Performance

### Compression

```go
import "github.com/gofiber/fiber/v2/middleware/compress"

app.Use(compress.New())
```

### Caching

```go
import "github.com/gofiber/fiber/v2/middleware/cache"

app.Use(cache.New(cache.Config{
    Expiration: 5 * time.Minute,
}))
```

### Connection Pooling

```go
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

## Code Optimization

### Use Goroutines

```go
go processInBackground(data)
```

### Context Timeouts

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

---

## Next Steps

- [**Custom Modules**](custom-modules.md)
- [**Middleware**](middleware.md)
- [**Security**](security.md)