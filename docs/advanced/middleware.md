# Middleware

Build custom middleware for Neonex Core applications using Fiber.

---

## Overview

Middleware intercepts HTTP requests/responses for:
- **Authentication** - Verify user identity
- **Authorization** - Check permissions
- **Logging** - Request/response logging
- **Rate Limiting** - Prevent abuse
- **CORS** - Cross-origin resource sharing
- **Validation** - Input validation
- **Compression** - Response compression

**Execution Order:** Middleware runs in registration order.

---

## Built-in Middleware

### Logger Middleware

```go
import "github.com/gofiber/fiber/v2/middleware/logger"

app.Use(logger.New(logger.Config{
    Format: "[${time}] ${status} - ${method} ${path} ${latency}\n",
    TimeFormat: "2006-01-02 15:04:05",
    TimeZone: "Local",
}))
```

### CORS Middleware

```go
import "github.com/gofiber/fiber/v2/middleware/cors"

app.Use(cors.New(cors.Config{
    AllowOrigins: "http://localhost:3000,https://example.com",
    AllowHeaders: "Origin,Content-Type,Accept,Authorization",
    AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
    AllowCredentials: true,
}))
```

### Recover Middleware

```go
import "github.com/gofiber/fiber/v2/middleware/recover"

app.Use(recover.New(recover.Config{
    EnableStackTrace: true,
}))
```

### Compress Middleware

```go
import "github.com/gofiber/fiber/v2/middleware/compress"

app.Use(compress.New(compress.Config{
    Level: compress.LevelBestSpeed,
}))
```

### Rate Limiter

```go
import "github.com/gofiber/fiber/v2/middleware/limiter"

app.Use(limiter.New(limiter.Config{
    Max: 20,
    Expiration: 30 * time.Second,
    KeyGenerator: func(c *fiber.Ctx) string {
        return c.IP()
    },
    LimitReached: func(c *fiber.Ctx) error {
        return c.Status(429).JSON(fiber.Map{
            "error": "Too many requests",
        })
    },
}))
```

---

## Custom Middleware

### Basic Structure

```go
func MyMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Before request
        start := time.Now()
        
        // Continue to next handler
        err := c.Next()
        
        // After request
        duration := time.Since(start)
        log.Printf("Request took %v", duration)
        
        return err
    }
}

// Use middleware
app.Use(MyMiddleware())
```

### Authentication Middleware

```go
func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := c.Get("Authorization")
        
        if token == "" {
            return c.Status(401).JSON(fiber.Map{
                "error": "Missing authorization token",
            })
        }
        
        // Remove "Bearer " prefix
        token = strings.TrimPrefix(token, "Bearer ")
        
        // Validate token
        userID, err := validateToken(token)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{
                "error": "Invalid token",
            })
        }
        
        // Store user ID in context
        c.Locals("userID", userID)
        
        return c.Next()
    }
}

// Use on protected routes
api := app.Group("/api", AuthMiddleware())
api.Get("/profile", getProfile)
```

### Authorization Middleware

```go
func RequireRole(roles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userID := c.Locals("userID").(string)
        
        user, err := getUserByID(userID)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{
                "error": "Unauthorized",
            })
        }
        
        // Check if user has required role
        hasRole := false
        for _, role := range roles {
            if user.Role == role {
                hasRole = true
                break
            }
        }
        
        if !hasRole {
            return c.Status(403).JSON(fiber.Map{
                "error": "Insufficient permissions",
            })
        }
        
        return c.Next()
    }
}

// Usage
app.Get("/admin", AuthMiddleware(), RequireRole("admin"), adminHandler)
```

### Request ID Middleware

```go
import "github.com/google/uuid"

func RequestIDMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        requestID := c.Get("X-Request-ID")
        
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        c.Set("X-Request-ID", requestID)
        c.Locals("requestID", requestID)
        
        return c.Next()
    }
}
```

### Timeout Middleware

```go
func TimeoutMiddleware(timeout time.Duration) fiber.Handler {
    return func(c *fiber.Ctx) error {
        ctx, cancel := context.WithTimeout(c.Context(), timeout)
        defer cancel()
        
        c.SetUserContext(ctx)
        
        done := make(chan error, 1)
        go func() {
            done <- c.Next()
        }()
        
        select {
        case err := <-done:
            return err
        case <-ctx.Done():
            return c.Status(408).JSON(fiber.Map{
                "error": "Request timeout",
            })
        }
    }
}

// Use with 5 second timeout
app.Use(TimeoutMiddleware(5 * time.Second))
```

### Validation Middleware

```go
func ValidateJSON() fiber.Handler {
    return func(c *fiber.Ctx) error {
        if c.Method() == "POST" || c.Method() == "PUT" {
            contentType := c.Get("Content-Type")
            
            if !strings.Contains(contentType, "application/json") {
                return c.Status(400).JSON(fiber.Map{
                    "error": "Content-Type must be application/json",
                })
            }
            
            var data map[string]interface{}
            if err := c.BodyParser(&data); err != nil {
                return c.Status(400).JSON(fiber.Map{
                    "error": "Invalid JSON",
                })
            }
        }
        
        return c.Next()
    }
}
```

---

## Middleware Patterns

### Chain Multiple Middleware

```go
// Method 1: Use()
app.Use(
    RequestIDMiddleware(),
    logger.New(),
    recover.New(),
)

// Method 2: Route-specific
api := app.Group("/api",
    AuthMiddleware(),
    RateLimitMiddleware(),
)

// Method 3: Handler-specific
app.Get("/admin",
    AuthMiddleware(),
    RequireRole("admin"),
    adminHandler,
)
```

### Conditional Middleware

```go
func ConditionalAuth() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Skip auth for public endpoints
        if strings.HasPrefix(c.Path(), "/public") {
            return c.Next()
        }
        
        // Require auth for other endpoints
        return AuthMiddleware()(c)
    }
}
```

### Middleware with Dependencies

```go
func LoggingMiddleware(logger logger.Logger) fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        
        logger.Info("Request started",
            zap.String("method", c.Method()),
            zap.String("path", c.Path()),
            zap.String("ip", c.IP()),
        )
        
        err := c.Next()
        
        logger.Info("Request completed",
            zap.Int("status", c.Response().StatusCode()),
            zap.Duration("duration", time.Since(start)),
        )
        
        return err
    }
}

// Use with injected dependencies
app.Use(LoggingMiddleware(myLogger))
```

---

## Advanced Examples

### API Key Authentication

```go
func APIKeyMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        apiKey := c.Get("X-API-Key")
        
        if apiKey == "" {
            return c.Status(401).JSON(fiber.Map{
                "error": "Missing API key",
            })
        }
        
        // Validate API key
        valid := validateAPIKey(apiKey)
        if !valid {
            return c.Status(401).JSON(fiber.Map{
                "error": "Invalid API key",
            })
        }
        
        return c.Next()
    }
}

func validateAPIKey(key string) bool {
    validKeys := []string{
        os.Getenv("API_KEY_1"),
        os.Getenv("API_KEY_2"),
    }
    
    for _, validKey := range validKeys {
        if key == validKey {
            return true
        }
    }
    
    return false
}
```

### IP Whitelist

```go
func IPWhitelistMiddleware(allowedIPs []string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        clientIP := c.IP()
        
        allowed := false
        for _, ip := range allowedIPs {
            if clientIP == ip {
                allowed = true
                break
            }
        }
        
        if !allowed {
            return c.Status(403).JSON(fiber.Map{
                "error": "IP not allowed",
            })
        }
        
        return c.Next()
    }
}

// Usage
allowedIPs := []string{"127.0.0.1", "192.168.1.100"}
app.Use("/admin", IPWhitelistMiddleware(allowedIPs))
```

### Cache Middleware

```go
var cache = make(map[string][]byte)
var cacheMutex sync.RWMutex

func CacheMiddleware(ttl time.Duration) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Only cache GET requests
        if c.Method() != "GET" {
            return c.Next()
        }
        
        key := c.Path()
        
        // Check cache
        cacheMutex.RLock()
        if data, ok := cache[key]; ok {
            cacheMutex.RUnlock()
            c.Set("X-Cache", "HIT")
            return c.Send(data)
        }
        cacheMutex.RUnlock()
        
        // Continue to handler
        err := c.Next()
        
        // Cache response
        if err == nil && c.Response().StatusCode() == 200 {
            data := c.Response().Body()
            cacheMutex.Lock()
            cache[key] = data
            cacheMutex.Unlock()
            
            // Clear cache after TTL
            go func() {
                time.Sleep(ttl)
                cacheMutex.Lock()
                delete(cache, key)
                cacheMutex.Unlock()
            }()
            
            c.Set("X-Cache", "MISS")
        }
        
        return err
    }
}

// Use with 5 minute cache
app.Use("/api", CacheMiddleware(5*time.Minute))
```

### Request Size Limit

```go
func RequestSizeLimitMiddleware(maxBytes int64) fiber.Handler {
    return func(c *fiber.Ctx) error {
        if c.Request().Header.ContentLength() > int(maxBytes) {
            return c.Status(413).JSON(fiber.Map{
                "error": "Request entity too large",
                "max_bytes": maxBytes,
            })
        }
        
        return c.Next()
    }
}

// Limit to 10MB
app.Use(RequestSizeLimitMiddleware(10 * 1024 * 1024))
```

---

## Testing Middleware

```go
func TestAuthMiddleware(t *testing.T) {
    app := fiber.New()
    app.Use(AuthMiddleware())
    app.Get("/test", func(c *fiber.Ctx) error {
        return c.SendString("OK")
    })
    
    // Test without token
    req := httptest.NewRequest("GET", "/test", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, 401, resp.StatusCode)
    
    // Test with valid token
    req = httptest.NewRequest("GET", "/test", nil)
    req.Header.Set("Authorization", "Bearer valid-token")
    resp, _ = app.Test(req)
    assert.Equal(t, 200, resp.StatusCode)
}
```

---

## Best Practices

### ✅ DO: Keep Middleware Focused

```go
// Good: Single responsibility
func LoggingMiddleware() fiber.Handler { /* ... */ }
func AuthMiddleware() fiber.Handler { /* ... */ }
func RateLimitMiddleware() fiber.Handler { /* ... */ }

// Bad: Does too much
func MegaMiddleware() fiber.Handler {
    // Logging + Auth + Rate Limit + Validation
}
```

### ✅ DO: Use Early Returns

```go
func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := c.Get("Authorization")
        
        if token == "" {
            return c.Status(401).JSON(fiber.Map{
                "error": "Missing token",
            })
        }
        
        // Continue processing
        return c.Next()
    }
}
```

### ❌ DON'T: Block c.Next()

```go
// Bad: Blocks request
func BadMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        time.Sleep(5 * time.Second)  // ❌ Blocks
        return c.Next()
    }
}

// Good: Use goroutines for async work
func GoodMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        go func() {
            // Async work
            time.Sleep(5 * time.Second)
            // Log or process
        }()
        
        return c.Next()
    }
}
```

---

## Middleware Order

```go
app := fiber.New()

// 1. Recovery (catch panics)
app.Use(recover.New())

// 2. Request ID (for tracing)
app.Use(RequestIDMiddleware())

// 3. Logger (log everything)
app.Use(logger.New())

// 4. CORS (handle preflight)
app.Use(cors.New())

// 5. Compression (compress responses)
app.Use(compress.New())

// 6. Rate Limit (prevent abuse)
app.Use(limiter.New())

// 7. Auth (protect routes)
protected := app.Group("/api", AuthMiddleware())

// 8. Business logic
protected.Get("/users", getUsers)
```

---

## Next Steps

- [**Error Handling**](error-handling.md) - Handle errors in middleware
- [**Security**](security.md) - Security middleware patterns
- [**Performance**](performance.md) - Optimize middleware
- [**Testing**](../development/testing.md) - Test middleware

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
