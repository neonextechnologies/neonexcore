# Debugging

Learn how to debug Neonex Core applications effectively.

---

## Overview

Debugging strategies for Neonex Core:

- **Print Debugging** - Quick insights with log statements
- **Delve Debugger** - Step-through debugging
- **VS Code Debugging** - Visual debugging experience
- **Request Logging** - HTTP request inspection
- **Database Query Logging** - SQL query analysis
- **Performance Profiling** - Find bottlenecks

---

## Quick Debugging

### Print Debugging

```go
// Simple logging
log.Println("User ID:", user.ID)
log.Printf("Processing order: %+v\n", order)

// With Zap logger
logger.Debug("Entering function",
    zap.String("user_id", userID),
    zap.Int("order_count", len(orders)),
)

logger.Info("Request received",
    zap.String("method", c.Method()),
    zap.String("path", c.Path()),
)
```

### Debug Environment Variable

```bash
# Enable debug mode
DEBUG=true neonex serve

# Or in .env
DEBUG=true
LOG_LEVEL=debug
```

---

## Delve Debugger

### Installation

```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Verify installation
dlv version
```

### Start Debugging

```bash
# Debug current package
dlv debug

# Debug with arguments
dlv debug -- --port 3000

# Debug specific file
dlv exec ./main.go

# Attach to running process
dlv attach <PID>
```

### Common Commands

```bash
# Set breakpoint
(dlv) break main.main
(dlv) break modules/user/service.go:25

# Continue execution
(dlv) continue
(dlv) c

# Step into function
(dlv) step
(dlv) s

# Step over function
(dlv) next
(dlv) n

# Print variable
(dlv) print user
(dlv) p user.ID

# List source code
(dlv) list
(dlv) ls

# Show stack trace
(dlv) stack
(dlv) bt

# Show goroutines
(dlv) goroutines

# Exit debugger
(dlv) exit
(dlv) quit
```

### Breakpoint Examples

```go
// Set breakpoint in code
import "runtime/debug"

func CreateUser(user *User) error {
    debug.SetTraceback("all")  // Detailed stack traces
    
    // Your code here
    if user.Name == "" {
        // Breakpoint: dlv break here
        return errors.New("name required")
    }
    
    return nil
}
```

---

## VS Code Debugging

### Setup Launch Configuration

Create `.vscode/launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Server",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/main.go",
            "env": {
                "APP_ENV": "development",
                "LOG_LEVEL": "debug"
            },
            "args": []
        },
        {
            "name": "Debug Test",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "${fileDirname}",
            "env": {
                "DB_DRIVER": "sqlite",
                "DB_NAME": ":memory:"
            }
        },
        {
            "name": "Attach to Process",
            "type": "go",
            "request": "attach",
            "mode": "local",
            "processId": "${command:pickProcess}"
        }
    ]
}
```

### Using VS Code Debugger

**Set Breakpoints:**
1. Click left margin of line number
2. Red dot appears
3. Press F5 to start debugging

**Debug Actions:**
- **F5** - Continue
- **F10** - Step Over
- **F11** - Step Into
- **Shift+F11** - Step Out
- **Shift+F5** - Stop
- **Ctrl+Shift+F5** - Restart

**Debug Panel:**
- **Variables** - Inspect local/global vars
- **Watch** - Monitor expressions
- **Call Stack** - See execution path
- **Breakpoints** - Manage breakpoints

### Example Session

```go
// modules/user/service.go
func (s *service) CreateUser(user *User) error {
    // Set breakpoint here (click line number)
    if user.Name == "" {
        return errors.New("name required")
    }
    
    // Breakpoint here too
    err := s.repo.Create(user)
    if err != nil {
        return err
    }
    
    return nil
}
```

**Debug Flow:**
1. Set breakpoints on lines 3 and 8
2. Press F5 to start
3. Make HTTP request to create user
4. Execution pauses at breakpoint
5. Inspect `user` variable
6. Press F10 to step through
7. Watch variables change

---

## HTTP Request Debugging

### Fiber Middleware Logging

```go
// Enable request logging
app.Use(logger.New(logger.Config{
    Format: "[${time}] ${status} - ${method} ${path} ${latency}\n",
    TimeFormat: "2006-01-02 15:04:05",
    TimeZone: "Local",
}))
```

### Custom Debug Middleware

```go
// middleware/debug.go
func DebugMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Log request
        log.Printf("→ %s %s", c.Method(), c.Path())
        log.Printf("  Headers: %v", c.GetReqHeaders())
        log.Printf("  Body: %s", string(c.Body()))
        
        // Continue
        err := c.Next()
        
        // Log response
        log.Printf("← Status: %d", c.Response().StatusCode())
        log.Printf("  Body: %s", string(c.Response().Body()))
        
        return err
    }
}

// Use in main.go
if os.Getenv("DEBUG") == "true" {
    app.Use(middleware.DebugMiddleware())
}
```

### cURL Verbose Mode

```bash
# Verbose request
curl -v http://localhost:8080/users

# Include headers
curl -i http://localhost:8080/users

# Save response
curl -o response.json http://localhost:8080/users

# Time request
curl -w "Time: %{time_total}s\n" http://localhost:8080/users
```

---

## Database Query Debugging

### Enable GORM Logger

```go
// Enable SQL logging
db, err := gorm.Open(sqlite.Open("neonex.db"), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})

// Custom logger
newLogger := logger.New(
    log.New(os.Stdout, "\r\n", log.LstdFlags),
    logger.Config{
        SlowThreshold:             200 * time.Millisecond,
        LogLevel:                  logger.Info,
        IgnoreRecordNotFoundError: false,
        Colorful:                  true,
    },
)

db, err := gorm.Open(sqlite.Open("neonex.db"), &gorm.Config{
    Logger: newLogger,
})
```

**Output:**
```sql
[2024-11-30 10:30:00] [2.12ms] SELECT * FROM `users` WHERE `id` = 1
[2024-11-30 10:30:01] [1.45ms] INSERT INTO `users` (`name`,`email`) VALUES ('John','john@example.com')
[2024-11-30 10:30:02] [0.89ms] UPDATE `users` SET `name`='Jane' WHERE `id` = 1
```

### Query Debugging

```go
// Debug single query
result := db.Debug().Where("id = ?", 1).First(&user)

// Show query without executing
sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
    return tx.Where("id = ?", 1).First(&user)
})
log.Println("SQL:", sql)
```

---

## Error Debugging

### Stack Traces

```go
import (
    "runtime/debug"
    "log"
)

func handleError(err error) {
    if err != nil {
        log.Printf("Error: %v\n", err)
        log.Printf("Stack trace:\n%s", debug.Stack())
    }
}

// Usage
if err := service.CreateUser(user); err != nil {
    handleError(err)
}
```

### Custom Error Types

```go
// pkg/errors/errors.go
type AppError struct {
    Code    int
    Message string
    Err     error
    Stack   string
}

func (e *AppError) Error() string {
    return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
}

func NewError(code int, message string, err error) *AppError {
    return &AppError{
        Code:    code,
        Message: message,
        Err:     err,
        Stack:   string(debug.Stack()),
    }
}

// Usage
if err := db.Create(&user).Error; err != nil {
    return NewError(500, "Failed to create user", err)
}
```

---

## Performance Debugging

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
    
    // Your application code
    app := core.NewApp()
    app.Run()
}
```

**Analyze:**
```bash
# View profile
go tool pprof cpu.prof

# Top functions
(pprof) top

# List function source
(pprof) list FunctionName

# Web interface
(pprof) web
```

### Memory Profiling

```go
import (
    "runtime/pprof"
    "os"
)

func dumpMemProfile() {
    f, _ := os.Create("mem.prof")
    defer f.Close()
    pprof.WriteHeapProfile(f)
}

// Add HTTP endpoint
import "net/http/pprof"

func main() {
    // Enable pprof endpoints
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    app := core.NewApp()
    app.Run()
}
```

**View in browser:**
```bash
# Open pprof
http://localhost:6060/debug/pprof/

# Heap profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutines
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

### Trace Analysis

```go
import (
    "runtime/trace"
    "os"
)

func main() {
    // Start trace
    f, _ := os.Create("trace.out")
    defer f.Close()
    trace.Start(f)
    defer trace.Stop()
    
    app := core.NewApp()
    app.Run()
}
```

**View trace:**
```bash
go tool trace trace.out
```

---

## Common Issues

### Issue 1: Port Already in Use

**Error:**
```
Error: listen tcp :8080: bind: address already in use
```

**Debug:**
```bash
# Find process
lsof -i :8080           # Mac/Linux
netstat -ano | findstr :8080  # Windows

# Kill process
kill -9 <PID>           # Mac/Linux
taskkill /F /PID <PID>  # Windows
```

### Issue 2: Database Locked

**Error:**
```
Error: database is locked
```

**Debug:**
```go
// Check connections
db, _ := database.DB()
stats := db.Stats()
log.Printf("Open connections: %d", stats.OpenConnections)

// Set connection pool
db.SetMaxOpenConns(1)  // SQLite single connection
db.SetMaxIdleConns(1)
```

### Issue 3: Nil Pointer

**Error:**
```
panic: runtime error: invalid memory address or nil pointer dereference
```

**Debug:**
```go
// Add nil checks
func ProcessUser(user *User) error {
    if user == nil {
        return errors.New("user is nil")
    }
    
    // Safe to use user
    log.Println("User:", user.Name)
    return nil
}
```

### Issue 4: Goroutine Leak

**Debug:**
```go
import "runtime"

// Before
before := runtime.NumGoroutine()

// Your code
go someFunction()

// After
after := runtime.NumGoroutine()
if after > before+1 {
    log.Printf("Goroutine leak detected: %d", after-before)
}
```

---

## Debugging Tips

### 1. Start Simple

```go
// Add log at function entry
func CreateUser(user *User) error {
    log.Printf("CreateUser called with: %+v", user)
    defer log.Println("CreateUser finished")
    
    // Your code
    return nil
}
```

### 2. Binary Search Debugging

```go
log.Println("1. Before database call")
err := db.Create(&user).Error
log.Println("2. After database call")

if err != nil {
    log.Println("3. Error branch:", err)
    return err
}

log.Println("4. Success branch")
return nil
```

### 3. Use Debug Build

```bash
# Build with debugging symbols
go build -gcflags="all=-N -l" -o debug-app main.go

# Run with race detector
go run -race main.go

# Run tests with race detector
go test -race ./...
```

### 4. Environment-Specific Debugging

```go
// Debug only in development
if os.Getenv("APP_ENV") == "development" {
    log.Printf("Debug: user=%+v, order=%+v", user, order)
}
```

---

## Tools

### Recommended Tools

| Tool | Purpose | Installation |
|------|---------|--------------|
| **Delve** | Go debugger | `go install github.com/go-delve/delve/cmd/dlv@latest` |
| **pprof** | Profiling | Built-in |
| **Postman** | API testing | [Download](https://www.postman.com/downloads/) |
| **DB Browser** | SQLite viewer | [Download](https://sqlitebrowser.org/) |
| **httpie** | HTTP client | `pip install httpie` |

### Browser Extensions

- **JSON Viewer** - Format JSON responses
- **REST Client** - Test APIs from browser
- **Redux DevTools** - Debug state (if using frontend)

---

## Next Steps

- [**Testing**](testing.md) - Write tests to catch bugs
- [**Best Practices**](best-practices.md) - Prevent common issues
- [**Logging**](../logging/overview.md) - Effective logging strategies
- [**Performance**](../advanced/performance.md) - Optimize your app

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
