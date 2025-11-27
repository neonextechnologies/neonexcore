# Neonex Logger System

ระบบ Logging ที่สมบูรณ์แบบสำหรับ Neonex Core

## Features

- ✅ **Multiple Log Levels** - DEBUG, INFO, WARN, ERROR, FATAL
- ✅ **Flexible Formatters** - Text และ JSON format
- ✅ **Multiple Outputs** - Console, File, หรือทั้งสองอย่าง
- ✅ **File Rotation** - อัตโนมัติตาม size และ age
- ✅ **Structured Logging** - Fields support
- ✅ **HTTP Middleware** - Request/Response logging
- ✅ **Request ID Tracking** - ติดตาม request แต่ละอัน
- ✅ **Colored Output** - สีสันในคอนโซล
- ✅ **Caller Information** - แสดงไฟล์และบรรทัดที่ log
- ✅ **Context Support** - รองรับ context.Context

## Quick Start

### Basic Usage

```go
import "neonexcore/pkg/logger"

func main() {
    // Simple logging
    logger.Info("Application started")
    logger.Warn("This is a warning")
    logger.Error("Something went wrong")
    
    // With fields
    logger.Info("User logged in", logger.Fields{
        "user_id": 123,
        "ip": "192.168.1.1",
    })
}
```

### Configuration

```go
import "neonexcore/pkg/logger"

func main() {
    // Load from environment
    config := logger.LoadConfig()
    
    // Or create custom config
    config := logger.Config{
        Level:        "debug",
        Format:       "json",
        Output:       "both", // console, file, or both
        FilePath:     "logs/app.log",
        MaxSize:      100, // MB
        MaxBackups:   7,
        MaxAge:       30, // days
        EnableCaller: true,
        EnableColor:  true,
    }
    
    // Setup logger
    if err := logger.Setup(config); err != nil {
        panic(err)
    }
}
```

### Environment Variables

```bash
LOG_LEVEL=debug
LOG_FORMAT=json
LOG_OUTPUT=both
LOG_FILE_PATH=logs/app.log
```

## Log Levels

```go
logger.SetGlobalLevel(logger.DebugLevel)
logger.SetGlobalLevel(logger.InfoLevel)
logger.SetGlobalLevel(logger.WarnLevel)
logger.SetGlobalLevel(logger.ErrorLevel)
```

## Formatters

### Text Formatter (Default)

```go
logger.SetGlobalFormatter(logger.NewTextFormatter())
```

Output:
```
[2024-11-25 10:30:45] INFO  [user.go:23] User logged in | user_id=123, action=login
[2024-11-25 10:30:46] WARN  [auth.go:45] Invalid token | token=xxx, ip=192.168.1.1
[2024-11-25 10:30:47] ERROR [db.go:89] Connection failed | error=timeout
```

### JSON Formatter

```go
logger.SetGlobalFormatter(logger.NewJSONFormatter())
```

Output:
```json
{"time":"2024-11-25T10:30:45.123Z","level":"INFO","message":"User logged in","caller":"user.go:23","user_id":123,"action":"login"}
{"time":"2024-11-25T10:30:46.456Z","level":"WARN","message":"Invalid token","caller":"auth.go:45","token":"xxx","ip":"192.168.1.1"}
```

## File Writer with Rotation

```go
fileWriter, err := logger.NewFileWriter(logger.FileWriterConfig{
    Filename:     "logs/app.log",
    MaxSize:      100,  // 100 MB
    MaxBackups:   7,    // Keep 7 backup files
    MaxAge:       30,   // Keep logs for 30 days
    RotateOnDate: true, // Rotate daily
})

if err != nil {
    panic(err)
}

logger.AddGlobalWriter(fileWriter)
```

## Structured Logging

```go
// Create logger with default fields
userLogger := logger.With(logger.Fields{
    "module": "user",
    "service": "authentication",
})

// All logs from userLogger will include these fields
userLogger.Info("Login attempt", logger.Fields{
    "user_id": 123,
    "ip": "192.168.1.1",
})

userLogger.Warn("Failed login", logger.Fields{
    "user_id": 456,
    "reason": "invalid_password",
})
```

## HTTP Middleware

### Setup in Fiber

```go
import (
    "neonexcore/pkg/logger"
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()
    
    // Add logger middleware
    app.Use(logger.RequestIDMiddleware(myLogger))
    app.Use(logger.HTTPMiddleware(myLogger))
    
    // Your routes...
    app.Listen(":8080")
}
```

### Use Logger in Handlers

```go
func myHandler(c *fiber.Ctx) error {
    // Get logger with request context
    log := logger.GetLogger(c)
    
    log.Info("Processing request", logger.Fields{
        "user_id": getUserID(c),
    })
    
    return c.JSON(fiber.Map{"status": "ok"})
}
```

## Integration with Neonex Core

ใน `main.go`:

```go
package main

import (
    "neonexcore/internal/core"
    "neonexcore/pkg/logger"
)

func main() {
    app := core.NewApp()
    
    // Initialize Logger
    loggerConfig := logger.LoadConfig()
    if err := app.InitLogger(loggerConfig); err != nil {
        panic(err)
    }
    
    // Initialize Database
    if err := app.InitDatabase(); err != nil {
        panic(err)
    }
    
    // ... rest of initialization
    
    app.StartHTTP()
}
```

Logger จะถูก integrate อัตโนมัติและจะ log:
- Database initialization
- Model registration
- Migration progress
- Module loading
- HTTP requests/responses
- Server startup

## Example Log Output

### Console (Text Format with Colors)

```
[2024-11-25 10:30:00] INFO  Logger initialized | level=info, format=text, output=console
[2024-11-25 10:30:00] INFO  Database initialized | driver=sqlite
[2024-11-25 10:30:00] INFO  Models registered for migration | count=1
[2024-11-25 10:30:00] INFO  Running auto-migration...
[2024-11-25 10:30:00] INFO  Auto-migration completed
[2024-11-25 10:30:00] INFO  Neonex Core booting...
[2024-11-25 10:30:00] INFO  Registering modules...
[2024-11-25 10:30:00] INFO  HTTP server starting | port=8080
[2024-11-25 10:30:05] INFO  HTTP Request | method=GET, path=/user/, status=200, duration=45, ip=127.0.0.1
[2024-11-25 10:30:06] WARN  HTTP Request | method=GET, path=/user/999, status=404, duration=12, ip=127.0.0.1
```

### File (JSON Format)

```json
{"time":"2024-11-25T10:30:00.123Z","level":"INFO","message":"Logger initialized","level":"info","format":"text","output":"console"}
{"time":"2024-11-25T10:30:00.234Z","level":"INFO","message":"Database initialized","driver":"sqlite"}
{"time":"2024-11-25T10:30:05.567Z","level":"INFO","message":"HTTP Request","method":"GET","path":"/user/","status":200,"duration":45,"ip":"127.0.0.1","request_id":"20241125103005-abc123"}
```

## Testing

Run the test file:

```bash
go run test_logger.go
```

This will test:
1. Text format output
2. JSON format output
3. Logger with fields
4. Different log levels
5. File output with rotation
6. Structured logging

## File Structure

```
pkg/logger/
├── logger.go      # Main logger interface and implementation
├── formatter.go   # Text and JSON formatters
├── writer.go      # File writer with rotation
├── middleware.go  # HTTP middleware for Fiber
└── config.go      # Configuration management
```

## Best Practices

1. **Use Structured Logging**
   ```go
   // Good
   logger.Info("User created", logger.Fields{"user_id": 123, "email": "user@example.com"})
   
   // Avoid
   logger.Info(fmt.Sprintf("User %d created with email %s", 123, "user@example.com"))
   ```

2. **Create Module-Specific Loggers**
   ```go
   var log = logger.With(logger.Fields{"module": "user"})
   
   func CreateUser() {
       log.Info("Creating user...")
   }
   ```

3. **Use Appropriate Log Levels**
   - DEBUG: รายละเอียดสำหรับ debugging
   - INFO: ข้อมูลทั่วไป, ปกติ
   - WARN: เตือน, ไม่ร้ายแรง
   - ERROR: ข้อผิดพลาด, ควรแก้ไข
   - FATAL: ข้อผิดพลาดร้ายแรง, จะหยุดโปรแกรม

4. **Always Include Context**
   ```go
   logger.Error("Database query failed", logger.Fields{
       "query": "SELECT * FROM users",
       "error": err.Error(),
       "table": "users",
   })
   ```

## Advanced Usage

### Custom Formatter

```go
type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logger.Entry) ([]byte, error) {
    return []byte(fmt.Sprintf("[%s] %s\n", entry.Level, entry.Message)), nil
}

logger.SetGlobalFormatter(&CustomFormatter{})
```

### Multiple Writers

```go
fileWriter1, _ := logger.NewFileWriter(logger.FileWriterConfig{
    Filename: "logs/app.log",
})

fileWriter2, _ := logger.NewFileWriter(logger.FileWriterConfig{
    Filename: "logs/error.log",
})

multiWriter := logger.NewMultiWriter(os.Stdout, fileWriter1, fileWriter2)
logger.AddGlobalWriter(multiWriter)
```

### Context-Aware Logging

```go
func processRequest(ctx context.Context) {
    log := logger.WithContext(ctx)
    log.Info("Processing...")
}
```

---

**Logger พร้อมใช้งานแล้ว!** 🎉
