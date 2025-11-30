# Logger API

Complete API reference for Neonex Core's structured logging system powered by Zap.

---

## Overview

Neonex Core provides a high-performance, structured logging system built on top of **Uber's Zap logger**. It supports multiple output formats, log levels, and destinations.

**Package:** `neonexcore/pkg/logger`

**Features:**
- ⚡ **High Performance** - Zero-allocation logging in production
- 📊 **Structured Logging** - Key-value pairs instead of string formatting
- 🎯 **Multiple Levels** - Debug, Info, Warn, Error, Fatal
- 📁 **Multiple Outputs** - Console, file, or both simultaneously
- 🔄 **Auto-Rotation** - Automatic log file rotation by size/time
- 🎨 **Flexible Formats** - JSON for production, text for development

---

## Logger Interface

### Definition

```go
type Logger interface {
    Debug(msg string, fields ...Fields)
    Info(msg string, fields ...Fields)
    Warn(msg string, fields ...Fields)
    Error(msg string, fields ...Fields)
    Fatal(msg string, fields ...Fields)
    WithFields(fields Fields) Logger
    Sync() error
}
```

Main logger interface for application-wide logging.

---

### Fields Type

```go
type Fields map[string]interface{}
```

Map of key-value pairs for structured logging.

**Example:**
```go
logger.Info("User created", logger.Fields{
    "user_id": 123,
    "email":   "user@example.com",
    "role":    "admin",
})
```

---

## Log Levels

### Debug

```go
Debug(msg string, fields ...Fields)
```

Logs detailed information for debugging purposes.

**When to use:**
- Development debugging
- Detailed trace information
- Variable inspection
- Flow tracking

**Example:**
```go
logger.Debug("Processing request", logger.Fields{
    "method":     "GET",
    "path":       "/api/users",
    "query":      r.URL.Query(),
    "headers":    r.Header,
    "start_time": time.Now(),
})
```

**Output (Development):**
```
2024-11-30T10:30:15.123Z  DEBUG  Processing request  {"method":"GET","path":"/api/users"}
```

---

### Info

```go
Info(msg string, fields ...Fields)
```

Logs general informational messages about application state.

**When to use:**
- Application lifecycle events
- Successful operations
- Important state changes
- Business events

**Examples:**

```go
// Startup
logger.Info("Application started", logger.Fields{
    "port":    8080,
    "env":     "production",
    "version": "1.0.0",
})

// Business events
logger.Info("Order placed", logger.Fields{
    "order_id": order.ID,
    "user_id":  order.UserID,
    "total":    order.Total,
    "items":    len(order.Items),
})

// User actions
logger.Info("User logged in", logger.Fields{
    "user_id": user.ID,
    "email":   user.Email,
    "ip":      clientIP,
})
```

**Output (Production JSON):**
```json
{
  "level": "info",
  "ts": "2024-11-30T10:30:15.123Z",
  "msg": "Order placed",
  "order_id": 12345,
  "user_id": 67,
  "total": 199.99,
  "items": 3
}
```

---

### Warn

```go
Warn(msg string, fields ...Fields)
```

Logs warning messages for potentially harmful situations.

**When to use:**
- Deprecated API usage
- Fallback to default values
- Recoverable errors
- Performance issues
- Unexpected but handled situations

**Examples:**

```go
// Deprecated feature
logger.Warn("Using deprecated API", logger.Fields{
    "endpoint":   "/api/v1/users",
    "deprecated": true,
    "use":        "/api/v2/users",
})

// Configuration fallback
logger.Warn("Config missing, using default", logger.Fields{
    "key":     "max_connections",
    "default": 100,
})

// Performance warning
logger.Warn("Slow query detected", logger.Fields{
    "query":    "SELECT * FROM users",
    "duration": "5.2s",
    "limit":    "1s",
})

// Business logic warning
logger.Warn("Low inventory", logger.Fields{
    "product_id": product.ID,
    "stock":      product.Stock,
    "threshold":  10,
})
```

---

### Error

```go
Error(msg string, fields ...Fields)
```

Logs error messages for serious problems that need attention.

**When to use:**
- Database errors
- External API failures
- Validation errors
- Business rule violations
- Any error that needs investigation

**Examples:**

```go
// Database error
logger.Error("Database query failed", logger.Fields{
    "error":     err.Error(),
    "query":     "INSERT INTO users",
    "operation": "create_user",
})

// External API error
logger.Error("Payment gateway failed", logger.Fields{
    "error":        err.Error(),
    "gateway":      "stripe",
    "amount":       order.Total,
    "order_id":     order.ID,
    "retry_count":  3,
})

// Business logic error
logger.Error("Insufficient inventory", logger.Fields{
    "product_id": productID,
    "requested":  quantity,
    "available":  stock,
    "order_id":   orderID,
})

// Validation error
logger.Error("Invalid input", logger.Fields{
    "error":  "email format invalid",
    "field":  "email",
    "value":  input.Email,
    "user_id": userID,
})
```

---

### Fatal

```go
Fatal(msg string, fields ...Fields)
```

Logs critical errors and **terminates the application** with `os.Exit(1)`.

**⚠️ Warning:** This method exits the application immediately after logging.

**When to use:**
- Startup failures (database connection, config loading)
- Critical resource unavailability
- Unrecoverable errors
- Security breaches

**Examples:**

```go
// Database connection failure
if err := db.Connect(); err != nil {
    logger.Fatal("Failed to connect to database", logger.Fields{
        "error": err.Error(),
        "host":  config.DBHost,
        "port":  config.DBPort,
    })
    // Application exits here
}

// Configuration error
if config.SecretKey == "" {
    logger.Fatal("Secret key not configured", logger.Fields{
        "env": os.Getenv("APP_ENV"),
    })
}

// Critical resource
if err := initCriticalService(); err != nil {
    logger.Fatal("Critical service initialization failed", logger.Fields{
        "error":   err.Error(),
        "service": "payment_processor",
    })
}
```

**Use sparingly!** Prefer `Error` for most cases.

---

### WithFields

```go
WithFields(fields Fields) Logger
```

Creates a contextual logger with pre-filled fields for all subsequent log calls.

**Returns:**
- `Logger` - New logger instance with attached fields

**When to use:**
- Request-scoped logging (request ID, user ID)
- Module-specific logging
- Operation tracking
- Consistent context across multiple log statements

**Examples:**

**Request-Scoped Logger:**
```go
// In HTTP middleware
func RequestLoggerMiddleware(c *fiber.Ctx) error {
    requestLogger := logger.WithFields(logger.Fields{
        "request_id": c.Locals("request_id"),
        "method":     c.Method(),
        "path":       c.Path(),
        "ip":         c.IP(),
    })
    
    c.Locals("logger", requestLogger)
    return c.Next()
}

// In handler
func (ctrl *Controller) CreateUser(c *fiber.Ctx) error {
    log := c.Locals("logger").(logger.Logger)
    
    log.Info("Creating user")  
    // Automatically includes: request_id, method, path, ip
    
    user, err := ctrl.service.Create(c.Context(), req)
    if err != nil {
        log.Error("Failed to create user", logger.Fields{
            "error": err.Error(),
        })
        // Automatically includes: request_id, method, path, ip, error
        return err
    }
    
    log.Info("User created successfully", logger.Fields{
        "user_id": user.ID,
    })
    // Automatically includes: request_id, method, path, ip, user_id
    
    return c.JSON(user)
}
```

**User-Scoped Logger:**
```go
type UserService struct {
    logger logger.Logger
}

func (s *UserService) ProcessUserAction(userID uint) error {
    // Create user-scoped logger
    userLogger := s.logger.WithFields(logger.Fields{
        "user_id": userID,
    })
    
    userLogger.Info("Starting user action")
    // Includes: user_id
    
    if err := s.validateUser(userID); err != nil {
        userLogger.Error("Validation failed", logger.Fields{
            "error": err.Error(),
        })
        // Includes: user_id, error
        return err
    }
    
    userLogger.Info("Action completed successfully")
    // Includes: user_id
    
    return nil
}
```

**Operation Tracking:**
```go
func (s *OrderService) ProcessOrder(orderID uint) error {
    // Operation-scoped logger
    opLogger := s.logger.WithFields(logger.Fields{
        "operation": "process_order",
        "order_id":  orderID,
    })
    
    opLogger.Info("Starting order processing")
    
    opLogger.Debug("Validating order")
    if err := s.validateOrder(orderID); err != nil {
        opLogger.Error("Validation failed", logger.Fields{"error": err.Error()})
        return err
    }
    
    opLogger.Debug("Charging payment")
    if err := s.chargePayment(orderID); err != nil {
        opLogger.Error("Payment failed", logger.Fields{"error": err.Error()})
        return err
    }
    
    opLogger.Debug("Updating inventory")
    if err := s.updateInventory(orderID); err != nil {
        opLogger.Error("Inventory update failed", logger.Fields{"error": err.Error()})
        return err
    }
    
    opLogger.Info("Order processed successfully")
    return nil
}
// All logs automatically include: operation, order_id
```

---

### Sync

```go
Sync() error
```

Flushes any buffered log entries to their destinations.

**Returns:**
- `error` - Error if sync fails, `nil` on success

**When to use:**
- Before application shutdown
- After critical operations
- In defer statements for cleanup

**Examples:**

```go
// Application shutdown
func main() {
    logger := logger.NewLogger()
    defer logger.Sync()  // Ensure logs are flushed
    
    // Application code...
}

// Critical operation
func ProcessCriticalData() error {
    defer logger.Sync()  // Flush logs before returning
    
    logger.Info("Starting critical operation")
    // ... process data ...
    logger.Info("Critical operation completed")
    
    return nil
}

// Error handling with sync
func HandleFatalError(err error) {
    logger.Fatal("Fatal error occurred", logger.Fields{
        "error": err.Error(),
    })
    logger.Sync()  // Ensure fatal log is written
    os.Exit(1)
}
```

---

## Configuration

### Setup Logger

```go
func Setup(cfg Config) error
```

Initializes the global logger with specified configuration.

**Parameters:**
- `cfg` - Logger configuration

**Example:**

```go
import "neonexcore/pkg/logger"

func main() {
    cfg := logger.Config{
        Level:      "info",
        Format:     "json",
        Output:     "both",
        FilePath:   "logs/app.log",
        MaxSize:    100,  // MB
        MaxBackups: 5,
        MaxAge:     30,   // days
    }
    
    if err := logger.Setup(cfg); err != nil {
        log.Fatalf("Failed to setup logger: %v", err)
    }
    
    logger.Info("Logger initialized")
}
```

---

### Config Struct

```go
type Config struct {
    Level      string  // debug, info, warn, error, fatal
    Format     string  // json, text
    Output     string  // console, file, both
    FilePath   string  // Path to log file
    MaxSize    int     // Maximum size in MB before rotation
    MaxBackups int     // Maximum number of old files to keep
    MaxAge     int     // Maximum days to keep old files
}
```

**Example Configurations:**

**Development:**
```go
cfg := logger.Config{
    Level:  "debug",
    Format: "text",
    Output: "console",
}
```

**Production:**
```go
cfg := logger.Config{
    Level:      "info",
    Format:     "json",
    Output:     "both",
    FilePath:   "/var/log/app/neonex.log",
    MaxSize:    100,
    MaxBackups: 10,
    MaxAge:     30,
}
```

**Testing:**
```go
cfg := logger.Config{
    Level:  "warn",
    Format: "text",
    Output: "console",
}
```

---

## Usage Examples

### Basic Logging

```go
import "neonexcore/pkg/logger"

func main() {
    logger.Info("Application starting")
    
    logger.Debug("Configuration loaded", logger.Fields{
        "config_file": "config.yaml",
        "env":         "production",
    })
    
    logger.Warn("Cache miss", logger.Fields{
        "key": "user:123",
    })
    
    logger.Error("Failed to send email", logger.Fields{
        "to":    "user@example.com",
        "error": err.Error(),
    })
}
```

### HTTP Request Logging

```go
func LogHTTPRequest(c *fiber.Ctx) error {
    start := time.Now()
    
    logger.Info("HTTP Request", logger.Fields{
        "method":     c.Method(),
        "path":       c.Path(),
        "ip":         c.IP(),
        "user_agent": c.Get("User-Agent"),
    })
    
    err := c.Next()
    
    logger.Info("HTTP Response", logger.Fields{
        "method":   c.Method(),
        "path":     c.Path(),
        "status":   c.Response().StatusCode(),
        "duration": time.Since(start).String(),
    })
    
    return err
}
```

### Database Query Logging

```go
func (r *repository) FindByID(ctx context.Context, id uint) (*User, error) {
    start := time.Now()
    
    logger.Debug("Executing query", logger.Fields{
        "operation": "FindByID",
        "table":     "users",
        "id":        id,
    })
    
    var user User
    err := r.db.WithContext(ctx).First(&user, id).Error
    
    duration := time.Since(start)
    
    if err != nil {
        logger.Error("Query failed", logger.Fields{
            "operation": "FindByID",
            "error":     err.Error(),
            "duration":  duration.String(),
        })
        return nil, err
    }
    
    logger.Debug("Query completed", logger.Fields{
        "operation": "FindByID",
        "duration":  duration.String(),
        "found":     user.ID > 0,
    })
    
    return &user, nil
}
```

### Service Layer Logging

```go
type UserService struct {
    repo   Repository
    logger logger.Logger
}

func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    s.logger.Info("Creating user", logger.Fields{
        "email": req.Email,
    })
    
    // Validation
    if err := s.validate(req); err != nil {
        s.logger.Warn("Validation failed", logger.Fields{
            "email": req.Email,
            "error": err.Error(),
        })
        return nil, err
    }
    
    // Create user
    user := &User{
        Name:  req.Name,
        Email: req.Email,
    }
    
    if err := s.repo.Create(ctx, user); err != nil {
        s.logger.Error("Failed to create user", logger.Fields{
            "email": req.Email,
            "error": err.Error(),
        })
        return nil, err
    }
    
    s.logger.Info("User created successfully", logger.Fields{
        "user_id": user.ID,
        "email":   user.Email,
    })
    
    return user, nil
}
```

---

## Best Practices

### ✅ DO: Use Structured Fields

```go
// Good: Structured
logger.Info("User login", logger.Fields{
    "user_id": user.ID,
    "email":   user.Email,
    "ip":      clientIP,
})

// Bad: String formatting
logger.Info(fmt.Sprintf("User %d (%s) logged in from %s", 
    user.ID, user.Email, clientIP))
```

### ✅ DO: Log Important Events

```go
// Good: Log business events
logger.Info("Order placed", logger.Fields{
    "order_id": order.ID,
    "total":    order.Total,
})

logger.Info("Payment processed", logger.Fields{
    "payment_id": payment.ID,
    "amount":     payment.Amount,
})
```

### ✅ DO: Use Appropriate Levels

```go
// Good: Correct level usage
logger.Debug("Cache hit")         // Development info
logger.Info("User logged in")     // Normal operation
logger.Warn("Deprecated API")     // Warning
logger.Error("Database timeout")  // Error needs fix
logger.Fatal("Config missing")    // Cannot continue
```

### ✅ DO: Use Contextual Loggers

```go
// Good: Context preserved
requestLogger := logger.WithFields(logger.Fields{
    "request_id": requestID,
})

requestLogger.Info("Processing")
requestLogger.Info("Completed")
// Both include request_id
```

### ❌ DON'T: Log Sensitive Data

```go
// Bad: Sensitive data exposed
logger.Info("User authenticated", logger.Fields{
    "password": user.Password,  // ❌ Never log passwords
    "credit_card": payment.CC,   // ❌ Never log CC numbers
})

// Good: Safe logging
logger.Info("User authenticated", logger.Fields{
    "user_id": user.ID,
    "method":  "password",
})
```

### ❌ DON'T: Log Too Much

```go
// Bad: Excessive logging
for _, item := range items {
    logger.Debug("Processing item", logger.Fields{
        "item": item,  // ❌ Too verbose
    })
}

// Good: Summary logging
logger.Debug("Processing items", logger.Fields{
    "count": len(items),
})
```

---

## Related Documentation

- [**Logging Overview**](../logging/overview.md) - Logging concepts
- [**Logging Configuration**](../logging/configuration.md) - Setup guide
- [**Logging Usage**](../logging/usage.md) - Usage examples
- [**Best Practices**](../logging/best-practices.md) - Logging patterns

---

## External Resources

- [Uber Zap](https://github.com/uber-go/zap) - Underlying logger
- [Structured Logging](https://www.honeycomb.io/blog/structured-logging-and-your-team)
- [12-Factor Logs](https://12factor.net/logs)

---

**Need help?** Check our [FAQ](../resources/faq.md) or [get support](../resources/support.md).
