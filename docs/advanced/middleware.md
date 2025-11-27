# Custom Middleware

Guide to creating custom middleware in Neonex Core.

---

## Overview

Middleware functions execute before/after HTTP handlers, enabling cross-cutting concerns like authentication, logging, and validation.

---

## Basic Middleware

```go
// pkg/http/middleware/custom.go
package middleware

import "github.com/gofiber/fiber/v2"

func Custom() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Before handler
        
        // Call next middleware/handler
        err := c.Next()
        
        // After handler
        return err
    }
}
```

---

## Authentication Middleware

```go
func Auth() fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := c.Get("Authorization")
        
        if token == "" {
            return c.Status(401).JSON(fiber.Map{
                "error": "Unauthorized",
            })
        }
        
        // Validate token
        userID, err := validateToken(token)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{
                "error": "Invalid token",
            })
        }
        
        // Store in context
        c.Locals("user_id", userID)
        
        return c.Next()
    }
}
```

---

## Validation Middleware

```go
func ValidateRequest[T any]() fiber.Handler {
    return func(c *fiber.Ctx) error {
        var req T
        
        if err := c.BodyParser(&req); err != nil {
            return c.Status(400).JSON(fiber.Map{
                "error": "Invalid request body",
            })
        }
        
        // Validate
        if err := validator.Validate(req); err != nil {
            return c.Status(422).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
        
        c.Locals("validated_request", req)
        return c.Next()
    }
}
```

---

## Rate Limiting

```go
func RateLimit(maxRequests int, window time.Duration) fiber.Handler {
    limiter := rate.NewLimiter(rate.Every(window), maxRequests)
    
    return func(c *fiber.Ctx) error {
        if !limiter.Allow() {
            return c.Status(429).JSON(fiber.Map{
                "error": "Too many requests",
            })
        }
        return c.Next()
    }
}
```

---

## CORS Middleware

```go
func CORS() fiber.Handler {
    return cors.New(cors.Config{
        AllowOrigins: "*",
        AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    })
}
```

---

## Next Steps

- [**Custom Modules**](custom-modules.md)
- [**Error Handling**](error-handling.md)
- [**Security**](security.md)
