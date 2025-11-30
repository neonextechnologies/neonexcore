# Performance Optimization

Optimize Neonex Core applications for maximum performance.

---

## Overview

Key performance areas:
- **Database** - Query optimization, indexing, connection pooling
- **HTTP** - Response caching, compression, keep-alive
- **Concurrency** - Goroutines, channels, worker pools
- **Memory** - Reduce allocations, garbage collection
- **CPU** - Algorithm efficiency, profiling

---

## Database Optimization

### Use Indexes

```go
type User struct {
    ID    uint   `gorm:"primarykey"`
    Email string `gorm:"uniqueIndex;size:100"`  // Fast lookups
    Name  string `gorm:"index"`                 // Fast searches
    Age   int    `gorm:"index"`                 // Range queries
}

// Composite index for common queries
type Order struct {
    ID       uint      `gorm:"primarykey"`
    UserID   uint      `gorm:"index:idx_user_status"`
    Status   string    `gorm:"index:idx_user_status;size:20"`
    CreatedAt time.Time `gorm:"index"`
}
```

### Avoid N+1 Queries

```go
// ❌ Bad: N+1 queries
users, _ := db.Find(&users).Error
for _, user := range users {
    db.Find(&user.Orders)  // N queries
}

// ✅ Good: Preload relationships
db.Preload("Orders").Find(&users)

// ✅ Good: Preload with conditions
db.Preload("Orders", "status = ?", "completed").Find(&users)

// ✅ Good: Nested preloading
db.Preload("Orders.Items").Find(&users)
```

### Use Select to Limit Fields

```go
// ❌ Bad: Select all columns
db.Find(&users)

// ✅ Good: Select only needed fields
db.Select("id", "name", "email").Find(&users)

// ✅ Good: Omit large fields
db.Omit("bio", "profile_picture").Find(&users)
```

### Batch Operations

```go
// ❌ Bad: Individual inserts
for _, user := range users {
    db.Create(&user)  // N queries
}

// ✅ Good: Batch insert
db.CreateInBatches(users, 100)  // Batch of 100

// ✅ Good: Transaction for related data
db.Transaction(func(tx *gorm.DB) error {
    tx.Create(&user)
    tx.Create(&profile)
    return nil
})
```

### Connection Pooling

```go
sqlDB, err := db.DB()

// Set max connections
sqlDB.SetMaxOpenConns(25)
sqlDB.SetMaxIdleConns(5)
sqlDB.SetConnMaxLifetime(5 * time.Minute)
sqlDB.SetConnMaxIdleTime(30 * time.Second)

// Monitor pool stats
stats := sqlDB.Stats()
log.Printf("Open connections: %d", stats.OpenConnections)
log.Printf("In use: %d", stats.InUse)
log.Printf("Idle: %d", stats.Idle)
```

### Query Caching

```go
var cache = make(map[string]interface{})
var mu sync.RWMutex

func GetUserByIDCached(id uint) (*User, error) {
    key := fmt.Sprintf("user:%d", id)
    
    // Check cache
    mu.RLock()
    if cached, ok := cache[key]; ok {
        mu.RUnlock()
        return cached.(*User), nil
    }
    mu.RUnlock()
    
    // Query database
    var user User
    if err := db.First(&user, id).Error; err != nil {
        return nil, err
    }
    
    // Store in cache
    mu.Lock()
    cache[key] = &user
    mu.Unlock()
    
    return &user, nil
}
```

---

## HTTP Optimization

### Enable Compression

```go
import "github.com/gofiber/fiber/v2/middleware/compress"

app.Use(compress.New(compress.Config{
    Level: compress.LevelBestSpeed,  // Fast compression
}))
```

### Response Caching

```go
import "github.com/gofiber/fiber/v2/middleware/cache"

app.Use(cache.New(cache.Config{
    Next: func(c *fiber.Ctx) bool {
        return c.Query("refresh") == "true"
    },
    Expiration:   30 * time.Minute,
    CacheControl: true,
}))
```

### Keep-Alive Connections

```go
app := fiber.New(fiber.Config{
    DisableKeepalive: false,  // Enable keep-alive
    ReadTimeout:      10 * time.Second,
    WriteTimeout:     10 * time.Second,
    IdleTimeout:      120 * time.Second,
})
```

### Reduce Response Size

```go
// ✅ Good: Return only needed fields
type UserResponse struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func getUser(c *fiber.Ctx) error {
    user, _ := service.GetUser(id)
    
    response := UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }
    
    return c.JSON(response)
}
```

### Pagination

```go
type PaginationParams struct {
    Page    int `query:"page"`
    PerPage int `query:"per_page"`
}

func GetUsers(c *fiber.Ctx) error {
    params := new(PaginationParams)
    c.QueryParser(params)
    
    if params.Page < 1 {
        params.Page = 1
    }
    if params.PerPage < 1 || params.PerPage > 100 {
        params.PerPage = 10
    }
    
    offset := (params.Page - 1) * params.PerPage
    
    var users []User
    var total int64
    
    db.Model(&User{}).Count(&total)
    db.Limit(params.PerPage).Offset(offset).Find(&users)
    
    return c.JSON(fiber.Map{
        "data": users,
        "meta": fiber.Map{
            "total":    total,
            "page":     params.Page,
            "per_page": params.PerPage,
        },
    })
}
```

---

## Concurrency Optimization

### Worker Pool Pattern

```go
type Job struct {
    ID   int
    Data interface{}
}

func worker(id int, jobs <-chan Job, results chan<- error) {
    for job := range jobs {
        // Process job
        err := processJob(job)
        results <- err
    }
}

func ProcessBatch(items []interface{}) []error {
    numWorkers := runtime.NumCPU()
    jobs := make(chan Job, len(items))
    results := make(chan error, len(items))
    
    // Start workers
    for w := 1; w <= numWorkers; w++ {
        go worker(w, jobs, results)
    }
    
    // Send jobs
    for i, item := range items {
        jobs <- Job{ID: i, Data: item}
    }
    close(jobs)
    
    // Collect results
    var errors []error
    for range items {
        if err := <-results; err != nil {
            errors = append(errors, err)
        }
    }
    
    return errors
}
```

### Concurrent API Calls

```go
func GetUserData(userID uint) (*UserData, error) {
    var wg sync.WaitGroup
    var profile Profile
    var orders []Order
    var err1, err2 error
    
    // Fetch profile
    wg.Add(1)
    go func() {
        defer wg.Done()
        err1 = db.Where("user_id = ?", userID).First(&profile).Error
    }()
    
    // Fetch orders
    wg.Add(1)
    go func() {
        defer wg.Done()
        err2 = db.Where("user_id = ?", userID).Find(&orders).Error
    }()
    
    wg.Wait()
    
    if err1 != nil || err2 != nil {
        return nil, errors.New("failed to fetch user data")
    }
    
    return &UserData{
        Profile: profile,
        Orders:  orders,
    }, nil
}
```

### Rate Limiting with Channels

```go
type RateLimiter struct {
    tokens chan struct{}
}

func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
    rl := &RateLimiter{
        tokens: make(chan struct{}, rate),
    }
    
    // Refill tokens
    go func() {
        ticker := time.NewTicker(interval)
        for range ticker.C {
            select {
            case rl.tokens <- struct{}{}:
            default:
            }
        }
    }()
    
    // Initial tokens
    for i := 0; i < rate; i++ {
        rl.tokens <- struct{}{}
    }
    
    return rl
}

func (rl *RateLimiter) Allow() bool {
    select {
    case <-rl.tokens:
        return true
    default:
        return false
    }
}

// Usage
limiter := NewRateLimiter(100, time.Second)

func handler(c *fiber.Ctx) error {
    if !limiter.Allow() {
        return c.Status(429).JSON(fiber.Map{
            "error": "Rate limit exceeded",
        })
    }
    
    // Handle request
    return c.SendString("OK")
}
```

---

## Memory Optimization

### Reduce Allocations

```go
// ❌ Bad: Multiple allocations
func BuildString() string {
    s := ""
    for i := 0; i < 1000; i++ {
        s += fmt.Sprintf("%d ", i)  // Allocates each time
    }
    return s
}

// ✅ Good: Single allocation
func BuildString() string {
    var sb strings.Builder
    sb.Grow(5000)  // Pre-allocate
    for i := 0; i < 1000; i++ {
        sb.WriteString(fmt.Sprintf("%d ", i))
    }
    return sb.String()
}
```

### Reuse Objects (sync.Pool)

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func ProcessData(data []byte) []byte {
    // Get buffer from pool
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()
    
    // Use buffer
    buf.Write(data)
    // Process...
    
    return buf.Bytes()
}
```

### Avoid Memory Leaks

```go
// ✅ Good: Close resources
func ProcessFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()  // Always close
    
    // Process file
    return nil
}

// ✅ Good: Cancel contexts
func MakeRequest(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()  // Prevent leak
    
    // Make request
    return nil
}
```

---

## Profiling

### CPU Profiling

```go
import (
    "runtime/pprof"
    "os"
)

func main() {
    // Start CPU profiling
    f, _ := os.Create("cpu.prof")
    defer f.Close()
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
    
    // Run application
    app := core.NewApp()
    app.Run()
}
```

**Analyze:**
```bash
go tool pprof cpu.prof
(pprof) top10
(pprof) list FunctionName
```

### Memory Profiling

```go
import "net/http/pprof"

func main() {
    // Enable pprof endpoints
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()
    
    app := core.NewApp()
    app.Run()
}
```

**Analyze:**
```bash
# Heap profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutines
go tool pprof http://localhost:6060/debug/pprof/goroutine

# Allocations
go tool pprof http://localhost:6060/debug/pprof/allocs
```

### Benchmarking

```go
func BenchmarkGetUser(b *testing.B) {
    setup()
    defer teardown()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        service.GetUser(1)
    }
}

func BenchmarkGetUserParallel(b *testing.B) {
    setup()
    defer teardown()
    
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            service.GetUser(1)
        }
    })
}
```

**Run:**
```bash
go test -bench=. -benchmem
go test -bench=GetUser -cpuprofile=cpu.prof
go test -bench=GetUser -memprofile=mem.prof
```

---

## Best Practices

### ✅ DO: Measure Before Optimizing

```bash
# Profile production
go tool pprof http://production-server/debug/pprof/profile

# Find bottlenecks
(pprof) top10
```

### ✅ DO: Use Appropriate Data Structures

```go
// Fast lookups: map
users := make(map[uint]*User)

// Ordered data: slice
items := make([]Item, 0, capacity)

// Queue: channel
queue := make(chan Job, 100)

// Set: map[T]struct{}
seen := make(map[string]struct{})
```

### ❌ DON'T: Premature Optimization

```go
// Bad: Over-optimized before measuring
func complexOptimization() {
    // Lots of complex code
}

// Good: Simple first, optimize if needed
func simpleImplementation() {
    // Clear, simple code
}
```

---

## Performance Checklist

- [ ] Database indexes on frequently queried fields
- [ ] Connection pooling configured
- [ ] Preload relationships (avoid N+1)
- [ ] Batch operations for bulk data
- [ ] HTTP compression enabled
- [ ] Response caching where appropriate
- [ ] Pagination for large datasets
- [ ] Worker pools for concurrent tasks
- [ ] Profiling in production
- [ ] Benchmarks for critical paths

---

## Next Steps

- [**Debugging**](../development/debugging.md) - Profile and debug
- [**Testing**](../development/testing.md) - Benchmark tests
- [**Security**](security.md) - Secure optimizations
- [**Monitoring**](../deployment/monitoring.md) - Track performance

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
