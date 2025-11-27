# HTTP Logging Middleware

Comprehensive guide to logging HTTP requests and responses in Neonex Core.

---

## Overview

HTTP logging middleware automatically logs all incoming requests and outgoing responses, providing visibility into your application's HTTP layer.

**Features:**
- Request/response logging
- Performance metrics
- Error tracking
- Custom field injection
- Configurable verbosity

---

## Basic Middleware

### Simple Request Logging

```go
// pkg/http/middleware/logging.go
package middleware

import (
    "time"
    "github.com/gofiber/fiber/v2"
    "neonexcore/pkg/logger"
)

func Logger(log logger.Logger) fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        
        // Log request
        log.Info("Request started", logger.Fields{
            "method": c.Method(),
            "path":   c.Path(),
            "ip":     c.IP(),
        })
        
        // Process request
        err := c.Next()
        
        // Calculate duration
        duration := time.Since(start)
        
        // Log response
        log.Info("Request completed", logger.Fields{
            "method":   c.Method(),
            "path":     c.Path(),
            "status":   c.Response().StatusCode(),
            "duration": duration,
        })
        
        return err
    }
}
```

### Registration

```go
// pkg/http/server.go
func (s *Server) SetupMiddleware() {
    // Register logging middleware
    s.app.Use(middleware.Logger(s.logger))
}
```

---

## Advanced Middleware

### Detailed Request/Response Logging

```go
// pkg/http/middleware/logging.go
package middleware

import (
    "time"
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "neonexcore/pkg/logger"
)

type LoggingConfig struct {
    Logger          logger.Logger
    SkipPaths       []string
    LogRequestBody  bool
    LogResponseBody bool
    LogHeaders      bool
}

func LoggerWithConfig(config LoggingConfig) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Skip configured paths
        if shouldSkip(c.Path(), config.SkipPaths) {
            return c.Next()
        }
        
        // Generate request ID
        requestID := uuid.New().String()
        c.Locals("request_id", requestID)
        
        // Create request logger
        reqLogger := config.Logger.WithFields(logger.Fields{
            "request_id": requestID,
        })
        
        // Build request fields
        fields := logger.Fields{
            "method":     c.Method(),
            "path":       c.Path(),
            "ip":         c.IP(),
            "user_agent": c.Get("User-Agent"),
        }
        
        // Add headers if configured
        if config.LogHeaders {
            fields["headers"] = headersToMap(c)
        }
        
        // Add request body if configured
        if config.LogRequestBody && len(c.Body()) > 0 {
            fields["request_body"] = string(c.Body())
        }
        
        // Log request
        reqLogger.Info("HTTP Request", fields)
        
        // Process request
        start := time.Now()
        err := c.Next()
        duration := time.Since(start)
        
        // Build response fields
        respFields := logger.Fields{
            "status":   c.Response().StatusCode(),
            "duration": duration,
            "size":     len(c.Response().Body()),
        }
        
        // Add response body if configured
        if config.LogResponseBody {
            respFields["response_body"] = string(c.Response().Body())
        }
        
        // Log response
        if err != nil {
            reqLogger.Error("HTTP Request Failed", logger.Fields{
                "error": err.Error(),
            })
        } else if c.Response().StatusCode() >= 500 {
            reqLogger.Error("HTTP Request Error", respFields)
        } else if c.Response().StatusCode() >= 400 {
            reqLogger.Warn("HTTP Request Warning", respFields)
        } else {
            reqLogger.Info("HTTP Request Success", respFields)
        }
        
        return err
    }
}

func shouldSkip(path string, skipPaths []string) bool {
    for _, skip := range skipPaths {
        if path == skip {
            return true
        }
    }
    return false
}

func headersToMap(c *fiber.Ctx) map[string]string {
    headers := make(map[string]string)
    c.Request().Header.VisitAll(func(key, value []byte) {
        headers[string(key)] = string(value)
    })
    return headers
}
```

### Usage

```go
// pkg/http/server.go
func (s *Server) SetupMiddleware() {
    s.app.Use(middleware.LoggerWithConfig(middleware.LoggingConfig{
        Logger:          s.logger,
        SkipPaths:       []string{"/health", "/metrics"},
        LogRequestBody:  false, // Enable in development only
        LogResponseBody: false, // Enable in development only
        LogHeaders:      true,
    }))
}
```

---

## Request ID Tracking

### Request ID Middleware

```go
// pkg/http/middleware/request_id.go
package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
)

func RequestID() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Get or generate request ID
        requestID := c.Get("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        // Store in context
        c.Locals("request_id", requestID)
        
        // Add to response header
        c.Set("X-Request-ID", requestID)
        
        return c.Next()
    }
}
```

### Using Request ID in Logs

```go
// modules/user/controller.go
func (ctrl *controller) Create(c *fiber.Ctx) error {
    // Get request ID
    requestID := c.Locals("request_id").(string)
    
    // Create logger with request ID
    reqLogger := ctrl.logger.WithFields(logger.Fields{
        "request_id": requestID,
    })
    
    // All logs include request ID
    reqLogger.Info("Creating user")
    
    // ... handle request
    
    return c.JSON(response)
}
```

---

## Performance Monitoring

### Response Time Logging

```go
func (ctrl *controller) SlowEndpoint(c *fiber.Ctx) error {
    start := time.Now()
    
    // Process request
    result, err := ctrl.service.HeavyOperation(c.Context())
    
    duration := time.Since(start)
    
    // Warn on slow requests
    if duration > 1*time.Second {
        ctrl.logger.Warn("Slow request detected", logger.Fields{
            "endpoint": c.Path(),
            "duration": duration,
        })
    }
    
    if err != nil {
        return err
    }
    
    return c.JSON(result)
}
```

### Performance Metrics

```go
// pkg/http/middleware/metrics.go
package middleware

type MetricsLogger struct {
    logger logger.Logger
}

func NewMetricsLogger(log logger.Logger) *MetricsLogger {
    return &MetricsLogger{logger: log}
}

func (m *MetricsLogger) Middleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        
        err := c.Next()
        
        duration := time.Since(start)
        statusCode := c.Response().StatusCode()
        
        // Log metrics
        m.logger.Info("HTTP Metrics", logger.Fields{
            "method":      c.Method(),
            "path":        c.Path(),
            "status":      statusCode,
            "duration_ms": duration.Milliseconds(),
            "size_bytes":  len(c.Response().Body()),
        })
        
        // Warn on slow requests
        if duration > 2*time.Second {
            m.logger.Warn("Slow HTTP request", logger.Fields{
                "path":     c.Path(),
                "duration": duration,
            })
        }
        
        return err
    }
}
```

---

## Error Logging

### Error Handler Middleware

```go
// pkg/http/middleware/error_handler.go
package middleware

import (
    "github.com/gofiber/fiber/v2"
    "neonexcore/pkg/logger"
)

func ErrorHandler(log logger.Logger) fiber.ErrorHandler {
    return func(c *fiber.Ctx, err error) error {
        // Default 500 status code
        code := fiber.StatusInternalServerError
        message := "Internal Server Error"
        
        // Fiber error
        if e, ok := err.(*fiber.Error); ok {
            code = e.Code
            message = e.Message
        }
        
        // Build log fields
        fields := logger.Fields{
            "error":      err.Error(),
            "method":     c.Method(),
            "path":       c.Path(),
            "status":     code,
            "ip":         c.IP(),
            "user_agent": c.Get("User-Agent"),
        }
        
        // Add request ID if available
        if requestID := c.Locals("request_id"); requestID != nil {
            fields["request_id"] = requestID
        }
        
        // Log based on status code
        if code >= 500 {
            log.Error("HTTP Error", fields)
        } else if code >= 400 {
            log.Warn("HTTP Client Error", fields)
        }
        
        // Send error response
        return c.Status(code).JSON(fiber.Map{
            "error": message,
        })
    }
}
```

### Registration

```go
// pkg/http/server.go
func NewServer(config *Config, logger logger.Logger) *Server {
    app := fiber.New(fiber.Config{
        ErrorHandler: middleware.ErrorHandler(logger),
    })
    
    return &Server{
        app:    app,
        logger: logger,
    }
}
```

---

## Filtering and Sanitization

### Skip Health Checks

```go
func LoggerWithConfig(config LoggingConfig) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Skip health check endpoints
        skipPaths := []string{
            "/health",
            "/healthz",
            "/ready",
            "/metrics",
            "/favicon.ico",
        }
        
        for _, path := range skipPaths {
            if c.Path() == path {
                return c.Next()
            }
        }
        
        // ... log request
    }
}
```

### Sanitize Sensitive Data

```go
func sanitizeHeaders(headers map[string]string) map[string]string {
    sensitive := []string{
        "Authorization",
        "Cookie",
        "X-API-Key",
    }
    
    sanitized := make(map[string]string)
    for k, v := range headers {
        if contains(sensitive, k) {
            sanitized[k] = "[REDACTED]"
        } else {
            sanitized[k] = v
        }
    }
    
    return sanitized
}

func LoggerWithConfig(config LoggingConfig) fiber.Handler {
    return func(c *fiber.Ctx) error {
        headers := headersToMap(c)
        sanitizedHeaders := sanitizeHeaders(headers)
        
        config.Logger.Info("Request", logger.Fields{
            "headers": sanitizedHeaders,
        })
        
        return c.Next()
    }
}
```

### Redact Request Body

```go
func sanitizeBody(body []byte, contentType string) string {
    // Don't log binary content
    if strings.Contains(contentType, "multipart") ||
       strings.Contains(contentType, "octet-stream") {
        return "[BINARY]"
    }
    
    // Truncate large bodies
    if len(body) > 1024 {
        return string(body[:1024]) + "... [TRUNCATED]"
    }
    
    // Redact sensitive fields (JSON)
    if strings.Contains(contentType, "json") {
        var data map[string]interface{}
        if err := json.Unmarshal(body, &data); err == nil {
            sensitive := []string{"password", "token", "secret", "ssn"}
            for _, field := range sensitive {
                if _, exists := data[field]; exists {
                    data[field] = "[REDACTED]"
                }
            }
            sanitized, _ := json.Marshal(data)
            return string(sanitized)
        }
    }
    
    return string(body)
}
```

---

## Access Logs

### Apache Combined Log Format

```go
// pkg/http/middleware/access_log.go
func AccessLog(log logger.Logger) fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        
        err := c.Next()
        
        duration := time.Since(start)
        
        // Apache Combined Log Format
        // IP - - [timestamp] "METHOD PATH PROTOCOL" STATUS SIZE "REFERER" "USER-AGENT" DURATION
        log.Info("", logger.Fields{
            "remote_addr": c.IP(),
            "timestamp":   start.Format("02/Jan/2006:15:04:05 -0700"),
            "method":      c.Method(),
            "path":        c.Path(),
            "protocol":    c.Protocol(),
            "status":      c.Response().StatusCode(),
            "size":        len(c.Response().Body()),
            "referer":     c.Get("Referer"),
            "user_agent":  c.Get("User-Agent"),
            "duration_ms": duration.Milliseconds(),
        })
        
        return err
    }
}
```

---

## Custom Loggers per Route

### Route-Specific Logging

```go
// modules/user/routes.go
func RegisterRoutes(router fiber.Router, ctrl *controller, logger logger.Logger) {
    userRouter := router.Group("/users")
    
    // Custom logger for user routes
    userLogger := logger.WithFields(logger.Fields{
        "module": "user",
    })
    
    // Apply custom logging middleware
    userRouter.Use(func(c *fiber.Ctx) error {
        c.Locals("logger", userLogger)
        return c.Next()
    })
    
    userRouter.Post("/", ctrl.Create)
    userRouter.Get("/:id", ctrl.GetByID)
}

// In controller
func (ctrl *controller) Create(c *fiber.Ctx) error {
    // Get route-specific logger
    log := c.Locals("logger").(logger.Logger)
    
    log.Info("Creating user")
    
    // ... handle request
}
```

---

## Structured Logging Example

### Complete Implementation

```go
// pkg/http/middleware/structured_logging.go
package middleware

import (
    "time"
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "neonexcore/pkg/logger"
)

func StructuredLogger(log logger.Logger) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Generate request ID
        requestID := uuid.New().String()
        c.Locals("request_id", requestID)
        c.Set("X-Request-ID", requestID)
        
        // Create request logger
        reqLogger := log.WithFields(logger.Fields{
            "request_id": requestID,
            "method":     c.Method(),
            "path":       c.Path(),
        })
        
        // Store logger in context
        c.Locals("logger", reqLogger)
        
        // Log request start
        reqLogger.Info("Request started", logger.Fields{
            "ip":         c.IP(),
            "user_agent": c.Get("User-Agent"),
        })
        
        // Process request
        start := time.Now()
        err := c.Next()
        duration := time.Since(start)
        
        // Build response fields
        fields := logger.Fields{
            "status":      c.Response().StatusCode(),
            "duration_ms": duration.Milliseconds(),
            "size_bytes":  len(c.Response().Body()),
        }
        
        // Log completion
        if err != nil {
            fields["error"] = err.Error()
            reqLogger.Error("Request failed", fields)
        } else if c.Response().StatusCode() >= 500 {
            reqLogger.Error("Request error", fields)
        } else if c.Response().StatusCode() >= 400 {
            reqLogger.Warn("Request client error", fields)
        } else {
            reqLogger.Info("Request completed", fields)
        }
        
        return err
    }
}
```

**JSON Output:**
```json
{
  "level": "info",
  "timestamp": "2024-01-15T10:30:45.123Z",
  "message": "Request completed",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "POST",
  "path": "/api/users",
  "ip": "192.168.1.1",
  "user_agent": "Mozilla/5.0...",
  "status": 201,
  "duration_ms": 45,
  "size_bytes": 156
}
```

---

## Best Practices

### ✅ DO:

**1. Use Request IDs**
```go
requestID := uuid.New().String()
c.Locals("request_id", requestID)
```

**2. Log Important Events**
```go
logger.Info("User created", logger.Fields{
    "user_id": user.ID,
    "email":   user.Email,
})
```

**3. Skip Health Checks**
```go
if c.Path() == "/health" {
    return c.Next() // Don't log health checks
}
```

### ❌ DON'T:

**1. Log Passwords**
```go
// Bad
logger.Info("Login", logger.Fields{
    "password": password, // ❌
})
```

**2. Log Large Payloads**
```go
// Bad
logger.Info("Request", logger.Fields{
    "body": string(c.Body()), // ❌ Could be huge
})
```

**3. Block on Logging**
```go
// Bad: Synchronous logging slows requests
// Good: Use async logger or buffered writer
```

---

## Next Steps

- [**Usage**](usage.md) - Logger usage in code
- [**Configuration**](configuration.md) - Logger setup
- [**Best Practices**](best-practices.md) - Logging guidelines
- [**Overview**](overview.md) - Logger architecture

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
