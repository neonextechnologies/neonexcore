# Logging Overview

Comprehensive logging system for Neonex Core applications.

---

## Introduction

Neonex Core provides a **structured logging system** built on top of popular Go logging libraries, offering:

- **Structured Logging** - JSON/text formatted logs
- **Multiple Levels** - Debug, Info, Warn, Error, Fatal
- **Contextual Fields** - Add metadata to log entries
- **Multiple Outputs** - Console, file, external services
- **Performance** - Minimal overhead with lazy evaluation

---

## Logger Architecture

### Components

```
┌──────────────────────────────────────┐
│         Application Code              │
└──────────────┬───────────────────────┘
               │
               ▼
┌──────────────────────────────────────┐
│       pkg/logger Interface            │
│  - Info(), Error(), Debug()           │
│  - WithFields(), WithContext()        │
└──────────────┬───────────────────────┘
               │
               ▼
┌──────────────────────────────────────┐
│       Logger Implementation           │
│  - zap/zerolog/logrus adapters        │
└──────────────┬───────────────────────┘
               │
               ▼
┌──────────────────────────────────────┐
│          Output Destinations          │
│  - Console (stdout/stderr)            │
│  - Files (rotated logs)               │
│  - External (ELK, Datadog, etc.)      │
└──────────────────────────────────────┘
```

---

## Quick Start

### Basic Logging

```go
package user

import "neonexcore/pkg/logger"

type service struct {
    logger logger.Logger
}

func (s *service) CreateUser(ctx context.Context, req *CreateUserRequest) error {
    // Info log
    s.logger.Info("Creating user", logger.Fields{
        "email": req.Email,
        "name":  req.Name,
    })
    
    user := &User{
        Name:  req.Name,
        Email: req.Email,
    }
    
    if err := s.repo.Create(ctx, user); err != nil {
        // Error log
        s.logger.Error("Failed to create user", logger.Fields{
            "error": err.Error(),
            "email": req.Email,
        })
        return err
    }
    
    // Success log
    s.logger.Info("User created successfully", logger.Fields{
        "user_id": user.ID,
        "email":   user.Email,
    })
    
    return nil
}
```

---

## Log Levels

### Level Hierarchy

| Level | Description | Use Case | Output |
|-------|-------------|----------|--------|
| **Debug** | Detailed diagnostic | Development debugging | Dev only |
| **Info** | General information | Normal operations | All envs |
| **Warn** | Warning messages | Potential issues | All envs |
| **Error** | Error conditions | Operation failures | All envs |
| **Fatal** | Critical errors | App termination | All envs |

### Setting Log Level

```go
// config/.env
LOG_LEVEL=info      # Production
# LOG_LEVEL=debug   # Development
```

### Usage by Level

```go
// Debug - Verbose diagnostic information
s.logger.Debug("Processing request", logger.Fields{
    "user_id":    userID,
    "request_id": requestID,
    "payload":    payload, // Detailed data
})

// Info - Important business events
s.logger.Info("User logged in", logger.Fields{
    "user_id": userID,
    "ip":      clientIP,
})

// Warn - Non-critical issues
s.logger.Warn("Slow query detected", logger.Fields{
    "query":    query,
    "duration": duration,
})

// Error - Operation failures
s.logger.Error("Payment processing failed", logger.Fields{
    "order_id": orderID,
    "error":    err.Error(),
})

// Fatal - Unrecoverable errors (exits app)
if err := db.Connect(); err != nil {
    s.logger.Fatal("Database connection failed", logger.Fields{
        "error": err.Error(),
    })
}
```

---

## Structured Logging

### Fields

```go
// Add context with fields
s.logger.Info("Order processed", logger.Fields{
    "order_id":   12345,
    "user_id":    678,
    "amount":     99.99,
    "currency":   "USD",
    "status":     "completed",
    "processing_time": time.Since(start),
})
```

**Output (JSON):**
```json
{
  "level": "info",
  "timestamp": "2024-01-15T10:30:45Z",
  "message": "Order processed",
  "order_id": 12345,
  "user_id": 678,
  "amount": 99.99,
  "currency": "USD",
  "status": "completed",
  "processing_time": "245ms"
}
```

### Contextual Logger

```go
// Create logger with common fields
userLogger := s.logger.WithFields(logger.Fields{
    "user_id": userID,
    "session": sessionID,
})

// All logs from this logger include user_id and session
userLogger.Info("Profile updated")
userLogger.Info("Password changed")
userLogger.Error("Invalid request", logger.Fields{
    "error": err.Error(),
})
```

---

## Log Formats

### JSON Format (Production)

```go
// config/.env
LOG_FORMAT=json
```

**Output:**
```json
{"level":"info","ts":"2024-01-15T10:30:45.123Z","msg":"User created","user_id":123,"email":"user@example.com"}
```

**Benefits:**
- Machine-readable
- Easy to parse
- Structured data
- Integration with log aggregators

### Text Format (Development)

```go
// config/.env
LOG_FORMAT=text
```

**Output:**
```
2024-01-15T10:30:45.123Z  INFO  User created  user_id=123 email=user@example.com
```

**Benefits:**
- Human-readable
- Better for console viewing
- Easier debugging locally

---

## Logger Interface

### Core Interface

```go
// pkg/logger/logger.go
package logger

type Logger interface {
    // Log methods
    Debug(msg string, fields ...Fields)
    Info(msg string, fields ...Fields)
    Warn(msg string, fields ...Fields)
    Error(msg string, fields ...Fields)
    Fatal(msg string, fields ...Fields)
    
    // Context methods
    WithFields(fields Fields) Logger
    WithContext(ctx context.Context) Logger
}

type Fields map[string]interface{}
```

### Implementation

```go
// Using Zap logger
type zapLogger struct {
    logger *zap.Logger
}

func NewZapLogger(config *Config) Logger {
    zapConfig := zap.NewProductionConfig()
    
    if config.Format == "text" {
        zapConfig = zap.NewDevelopmentConfig()
    }
    
    zapConfig.Level = zap.NewAtomicLevelAt(parseLevel(config.Level))
    
    logger, _ := zapConfig.Build()
    return &zapLogger{logger: logger}
}

func (l *zapLogger) Info(msg string, fields ...Fields) {
    zapFields := convertFields(fields...)
    l.logger.Info(msg, zapFields...)
}
```

---

## Configuration

### Environment Variables

```bash
# Log level
LOG_LEVEL=info          # debug, info, warn, error, fatal

# Log format
LOG_FORMAT=json         # json, text

# Log output
LOG_OUTPUT=stdout       # stdout, stderr, file

# Log file settings (if LOG_OUTPUT=file)
LOG_FILE_PATH=logs/app.log
LOG_FILE_MAX_SIZE=100   # MB
LOG_FILE_MAX_BACKUPS=3
LOG_FILE_MAX_AGE=28     # days
LOG_FILE_COMPRESS=true
```

### Configuration Struct

```go
// internal/config/logger.go
package config

type LoggerConfig struct {
    Level      string `env:"LOG_LEVEL" envDefault:"info"`
    Format     string `env:"LOG_FORMAT" envDefault:"json"`
    Output     string `env:"LOG_OUTPUT" envDefault:"stdout"`
    FilePath   string `env:"LOG_FILE_PATH" envDefault:"logs/app.log"`
    MaxSize    int    `env:"LOG_FILE_MAX_SIZE" envDefault:"100"`
    MaxBackups int    `env:"LOG_FILE_MAX_BACKUPS" envDefault:"3"`
    MaxAge     int    `env:"LOG_FILE_MAX_AGE" envDefault:"28"`
    Compress   bool   `env:"LOG_FILE_COMPRESS" envDefault:"true"`
}
```

---

## Output Destinations

### Console Output

```go
// Default: stdout
LOG_OUTPUT=stdout

// Error stream
LOG_OUTPUT=stderr
```

### File Output

```go
// Single file
LOG_OUTPUT=file
LOG_FILE_PATH=logs/app.log

// Rotated logs
LOG_FILE_MAX_SIZE=100      # Rotate at 100MB
LOG_FILE_MAX_BACKUPS=3     # Keep 3 old files
LOG_FILE_MAX_AGE=28        # Delete files older than 28 days
LOG_FILE_COMPRESS=true     # Compress old files
```

**File structure:**
```
logs/
├── app.log           # Current log
├── app-2024-01-14.log
├── app-2024-01-13.log.gz
└── app-2024-01-12.log.gz
```

### Multiple Outputs

```go
// Write to both console and file
logger := logger.NewMultiLogger(
    logger.NewConsoleLogger(consoleConfig),
    logger.NewFileLogger(fileConfig),
)
```

---

## Best Practices

### ✅ DO:

**1. Use Structured Fields**
```go
// Good: Structured
logger.Info("User login", logger.Fields{
    "user_id": 123,
    "ip":      "192.168.1.1",
})

// Bad: String concatenation
logger.Info(fmt.Sprintf("User %d logged in from %s", 123, "192.168.1.1"))
```

**2. Log Important Events**
```go
// Log business events
logger.Info("Payment processed", logger.Fields{
    "order_id": orderID,
    "amount":   amount,
})
```

**3. Add Context**
```go
// Create contextual logger
reqLogger := logger.WithFields(logger.Fields{
    "request_id": requestID,
    "user_id":    userID,
})
```

### ❌ DON'T:

**1. Log Sensitive Data**
```go
// Bad: Logging passwords
logger.Info("User created", logger.Fields{
    "email":    email,
    "password": password, // ❌ Never log passwords
})

// Good: Omit sensitive data
logger.Info("User created", logger.Fields{
    "email": email,
})
```

**2. Excessive Logging**
```go
// Bad: Too verbose
for i := 0; i < 1000; i++ {
    logger.Debug("Processing item", logger.Fields{"index": i}) // ❌
}

// Good: Log summary
logger.Debug("Processed items", logger.Fields{"count": 1000})
```

**3. Ignore Errors**
```go
// Bad: Silent failure
if err != nil {
    return err // ❌ No log
}

// Good: Log errors
if err != nil {
    logger.Error("Operation failed", logger.Fields{
        "error": err.Error(),
    })
    return err
}
```

---

## Performance Considerations

### Lazy Evaluation

```go
// Fields are only evaluated if log level is enabled
logger.Debug("Complex operation", logger.Fields{
    "data": expensiveComputation(), // Only called if debug enabled
})
```

### Conditional Logging

```go
// Check level before expensive operations
if logger.IsDebugEnabled() {
    data := gatherDebugInfo() // Expensive
    logger.Debug("Debug info", logger.Fields{"data": data})
}
```

### Async Logging

```go
// Buffer logs for async writing (reduces latency)
logger := logger.NewAsyncLogger(config)
defer logger.Sync() // Flush on shutdown
```

---

## Integration Examples

### HTTP Request Logging

```go
func LoggingMiddleware(logger logger.Logger) fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        
        // Add request ID to context
        requestID := uuid.New().String()
        reqLogger := logger.WithFields(logger.Fields{
            "request_id": requestID,
        })
        
        // Log request
        reqLogger.Info("Request started", logger.Fields{
            "method": c.Method(),
            "path":   c.Path(),
        })
        
        // Process request
        err := c.Next()
        
        // Log response
        reqLogger.Info("Request completed", logger.Fields{
            "status":   c.Response().StatusCode(),
            "duration": time.Since(start),
        })
        
        return err
    }
}
```

### Database Query Logging

```go
func (r *repository) Create(ctx context.Context, user *User) error {
    r.logger.Debug("Creating user", logger.Fields{
        "name":  user.Name,
        "email": user.Email,
    })
    
    start := time.Now()
    err := r.db.WithContext(ctx).Create(user).Error
    duration := time.Since(start)
    
    if err != nil {
        r.logger.Error("User creation failed", logger.Fields{
            "error":    err.Error(),
            "duration": duration,
        })
        return err
    }
    
    r.logger.Info("User created", logger.Fields{
        "user_id":  user.ID,
        "duration": duration,
    })
    
    return nil
}
```

---

## Next Steps

- [**Configuration**](configuration.md) - Logger configuration details
- [**Usage**](usage.md) - Advanced logging techniques
- [**Middleware**](middleware.md) - HTTP logging middleware
- [**Best Practices**](best-practices.md) - Logging guidelines

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
