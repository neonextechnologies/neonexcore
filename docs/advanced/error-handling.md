# Error Handling

Comprehensive error handling strategies for Neonex Core applications.

---

## Overview

Effective error handling improves:
- **User Experience** - Clear error messages
- **Debugging** - Easier troubleshooting
- **Reliability** - Graceful failure recovery
- **Security** - No sensitive data leaks

---

## Error Types

### Standard Errors

```go
import "errors"

// Simple error
err := errors.New("user not found")

// Formatted error
err := fmt.Errorf("user %d not found", id)

// Wrapped error
err := fmt.Errorf("failed to create user: %w", originalErr)
```

### Custom Error Types

```go
// ValidationError represents input validation failures
type ValidationError struct {
    Field   string
    Message string
    Value   interface{}
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// DatabaseError represents database operation failures
type DatabaseError struct {
    Operation string
    Table     string
    Err       error
}

func (e *DatabaseError) Error() string {
    return fmt.Sprintf("database %s failed on %s: %v", 
        e.Operation, e.Table, e.Err)
}

// NotFoundError for resource not found
type NotFoundError struct {
    Resource string
    ID       interface{}
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s with ID %v not found", e.Resource, e.ID)
}
```

---

## Error Handling Patterns

### Repository Layer

```go
// repository.go
func (r *repository) FindByID(id uint) (*User, error) {
    var user User
    err := r.DB.First(&user, id).Error
    
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, &NotFoundError{
                Resource: "User",
                ID:       id,
            }
        }
        return nil, &DatabaseError{
            Operation: "FindByID",
            Table:     "users",
            Err:       err,
        }
    }
    
    return &user, nil
}

func (r *repository) Create(user *User) error {
    err := r.DB.Create(user).Error
    if err != nil {
        return &DatabaseError{
            Operation: "Create",
            Table:     "users",
            Err:       err,
        }
    }
    return nil
}
```

### Service Layer

```go
// service.go
func (s *service) CreateUser(user *User) error {
    // Validate input
    if err := validateUser(user); err != nil {
        return err  // ValidationError
    }
    
    // Check duplicates
    existing, err := s.repo.FindByEmail(user.Email)
    if err == nil && existing != nil {
        return &ValidationError{
            Field:   "email",
            Message: "email already exists",
            Value:   user.Email,
        }
    }
    
    // Create user
    if err := s.repo.Create(user); err != nil {
        logger.Error("Failed to create user",
            zap.Error(err),
            zap.String("email", user.Email),
        )
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    return nil
}

func validateUser(user *User) error {
    if user.Name == "" {
        return &ValidationError{
            Field:   "name",
            Message: "name is required",
        }
    }
    
    if user.Email == "" {
        return &ValidationError{
            Field:   "email",
            Message: "email is required",
        }
    }
    
    if !isValidEmail(user.Email) {
        return &ValidationError{
            Field:   "email",
            Message: "invalid email format",
            Value:   user.Email,
        }
    }
    
    return nil
}
```

### Controller Layer

```go
// controller.go
func (c *Controller) Create(ctx *fiber.Ctx) error {
    var user User
    
    // Parse request body
    if err := ctx.BodyParser(&user); err != nil {
        return ctx.Status(400).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    
    // Create user
    err := c.service.CreateUser(&user)
    if err != nil {
        return handleError(ctx, err)
    }
    
    return ctx.Status(201).JSON(user)
}

func (c *Controller) GetByID(ctx *fiber.Ctx) error {
    id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
    if err != nil {
        return ctx.Status(400).JSON(fiber.Map{
            "error": "Invalid ID format",
        })
    }
    
    user, err := c.service.GetUser(uint(id))
    if err != nil {
        return handleError(ctx, err)
    }
    
    return ctx.JSON(user)
}
```

---

## Error Response Handler

### Centralized Error Handler

```go
// pkg/http/errors.go
func HandleError(ctx *fiber.Ctx, err error) error {
    // Validation errors -> 400
    var validationErr *ValidationError
    if errors.As(err, &validationErr) {
        return ctx.Status(400).JSON(fiber.Map{
            "error": "Validation failed",
            "field": validationErr.Field,
            "message": validationErr.Message,
        })
    }
    
    // Not found errors -> 404
    var notFoundErr *NotFoundError
    if errors.As(err, &notFoundErr) {
        return ctx.Status(404).JSON(fiber.Map{
            "error": notFoundErr.Error(),
        })
    }
    
    // Database errors -> 500
    var dbErr *DatabaseError
    if errors.As(err, &dbErr) {
        logger.Error("Database error", zap.Error(dbErr))
        return ctx.Status(500).JSON(fiber.Map{
            "error": "Internal server error",
        })
    }
    
    // Default -> 500
    logger.Error("Unhandled error", zap.Error(err))
    return ctx.Status(500).JSON(fiber.Map{
        "error": "Internal server error",
    })
}
```

### Usage in Controllers

```go
func (c *Controller) Create(ctx *fiber.Ctx) error {
    var user User
    if err := ctx.BodyParser(&user); err != nil {
        return HandleError(ctx, err)
    }
    
    if err := c.service.CreateUser(&user); err != nil {
        return HandleError(ctx, err)
    }
    
    return ctx.Status(201).JSON(user)
}
```

---

## Error Middleware

### Global Error Handler

```go
// middleware/error.go
func ErrorHandler() fiber.ErrorHandler {
    return func(ctx *fiber.Ctx, err error) error {
        code := fiber.StatusInternalServerError
        message := "Internal Server Error"
        
        // Fiber errors
        if e, ok := err.(*fiber.Error); ok {
            code = e.Code
            message = e.Message
        }
        
        // Custom errors
        var validationErr *ValidationError
        if errors.As(err, &validationErr) {
            code = fiber.StatusBadRequest
            message = validationErr.Error()
        }
        
        var notFoundErr *NotFoundError
        if errors.As(err, &notFoundErr) {
            code = fiber.StatusNotFound
            message = notFoundErr.Error()
        }
        
        // Log error
        logger.Error("Request failed",
            zap.Error(err),
            zap.Int("status", code),
            zap.String("method", ctx.Method()),
            zap.String("path", ctx.Path()),
        )
        
        // Send response
        return ctx.Status(code).JSON(fiber.Map{
            "error": message,
            "code":  code,
        })
    }
}

// Use in app
app := fiber.New(fiber.Config{
    ErrorHandler: ErrorHandler(),
})
```

---

## Panic Recovery

### Recover Middleware

```go
func RecoverMiddleware() fiber.Handler {
    return func(ctx *fiber.Ctx) error {
        defer func() {
            if r := recover(); r != nil {
                // Log panic with stack trace
                logger.Error("Panic recovered",
                    zap.Any("panic", r),
                    zap.String("stack", string(debug.Stack())),
                )
                
                // Send error response
                ctx.Status(500).JSON(fiber.Map{
                    "error": "Internal server error",
                })
            }
        }()
        
        return ctx.Next()
    }
}

// Use in app
app.Use(RecoverMiddleware())
```

---

## Error Wrapping

### Context-Rich Errors

```go
func (s *service) CreateUser(user *User) error {
    if err := s.repo.Create(user); err != nil {
        return fmt.Errorf("create user %s: %w", user.Email, err)
    }
    return nil
}

func (r *repository) Create(user *User) error {
    if err := r.DB.Create(user).Error; err != nil {
        return fmt.Errorf("insert into users table: %w", err)
    }
    return nil
}

// Error chain:
// create user john@example.com: insert into users table: UNIQUE constraint failed
```

### Unwrap Errors

```go
func GetRootCause(err error) error {
    for {
        unwrapped := errors.Unwrap(err)
        if unwrapped == nil {
            return err
        }
        err = unwrapped
    }
}

// Usage
err := service.CreateUser(user)
if err != nil {
    rootCause := GetRootCause(err)
    logger.Error("Root cause", zap.Error(rootCause))
}
```

---

## Best Practices

### ✅ DO: Log Errors with Context

```go
logger.Error("Failed to create user",
    zap.Error(err),
    zap.String("user_id", user.ID),
    zap.String("email", user.Email),
    zap.String("operation", "CreateUser"),
)
```

### ✅ DO: Return Meaningful Messages

```go
// Good
return &ValidationError{
    Field:   "email",
    Message: "email must be valid format",
    Value:   user.Email,
}

// Bad
return errors.New("error")
```

### ❌ DON'T: Expose Internal Errors

```go
// Bad: Exposes internal details
return ctx.Status(500).JSON(fiber.Map{
    "error": err.Error(),  // "sql: no rows in result set"
})

// Good: Generic message
return ctx.Status(500).JSON(fiber.Map{
    "error": "Internal server error",
})
```

### ❌ DON'T: Ignore Errors

```go
// Bad
user, _ := service.GetUser(id)  // Silent failure

// Good
user, err := service.GetUser(id)
if err != nil {
    return HandleError(ctx, err)
}
```

---

## Testing Error Handling

```go
func TestCreateUser_ValidationError(t *testing.T) {
    service := NewService(mockRepo)
    
    user := &User{Name: "", Email: "test@example.com"}
    err := service.CreateUser(user)
    
    assert.Error(t, err)
    
    var validationErr *ValidationError
    assert.True(t, errors.As(err, &validationErr))
    assert.Equal(t, "name", validationErr.Field)
}

func TestCreateUser_DatabaseError(t *testing.T) {
    mockRepo := new(MockRepository)
    mockRepo.On("Create", mock.Anything).
        Return(&DatabaseError{
            Operation: "Create",
            Table:     "users",
            Err:       errors.New("connection failed"),
        })
    
    service := NewService(mockRepo)
    err := service.CreateUser(&User{Name: "John", Email: "john@example.com"})
    
    assert.Error(t, err)
    
    var dbErr *DatabaseError
    assert.True(t, errors.As(err, &dbErr))
}
```

---

## Next Steps

- [**Middleware**](middleware.md) - Custom middleware patterns
- [**Security**](security.md) - Secure error handling
- [**Logging**](../logging/overview.md) - Error logging strategies
- [**Testing**](../development/testing.md) - Test error scenarios

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
