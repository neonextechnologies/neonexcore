# Logger Usage

Practical guide to using the logger in Neonex Core applications.

---

## Getting Started

### Inject Logger

```go
// modules/user/service.go
package user

import "neonexcore/pkg/logger"

type service struct {
    repo   Repository
    logger logger.Logger
}

func NewService(repo Repository, logger logger.Logger) Service {
    return &service{
        repo:   repo,
        logger: logger,
    }
}
```

### Basic Logging

```go
func (s *service) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    s.logger.Info("Creating user", logger.Fields{
        "email": req.Email,
    })
    
    user := &User{
        Name:  req.Name,
        Email: req.Email,
    }
    
    if err := s.repo.Create(ctx, user); err != nil {
        s.logger.Error("Failed to create user", logger.Fields{
            "error": err.Error(),
            "email": req.Email,
        })
        return nil, err
    }
    
    s.logger.Info("User created successfully", logger.Fields{
        "user_id": user.ID,
    })
    
    return user, nil
}
```

---

## Log Levels

### Debug Level

**Use for:** Detailed diagnostic information

```go
func (s *service) ProcessOrder(ctx context.Context, orderID uint) error {
    s.logger.Debug("Processing order", logger.Fields{
        "order_id": orderID,
        "step":     "validation",
    })
    
    order, err := s.repo.FindByID(ctx, orderID)
    s.logger.Debug("Order loaded", logger.Fields{
        "order_id": orderID,
        "status":   order.Status,
        "items":    len(order.Items),
    })
    
    // ... processing
    
    s.logger.Debug("Order processing complete", logger.Fields{
        "order_id": orderID,
        "duration": time.Since(start),
    })
    
    return nil
}
```

### Info Level

**Use for:** Important business events

```go
func (s *service) Login(ctx context.Context, email, password string) (*Session, error) {
    s.logger.Info("User login attempt", logger.Fields{
        "email": email,
    })
    
    user, err := s.repo.FindByEmail(ctx, email)
    if err != nil {
        s.logger.Info("Login failed: user not found", logger.Fields{
            "email": email,
        })
        return nil, ErrInvalidCredentials
    }
    
    if !user.VerifyPassword(password) {
        s.logger.Info("Login failed: invalid password", logger.Fields{
            "user_id": user.ID,
            "email":   email,
        })
        return nil, ErrInvalidCredentials
    }
    
    session := s.createSession(user)
    
    s.logger.Info("User logged in successfully", logger.Fields{
        "user_id":    user.ID,
        "session_id": session.ID,
    })
    
    return session, nil
}
```

### Warn Level

**Use for:** Potentially harmful situations

```go
func (s *service) ProcessPayment(ctx context.Context, payment *Payment) error {
    // Slow query warning
    start := time.Now()
    result, err := s.gateway.Process(payment)
    duration := time.Since(start)
    
    if duration > 5*time.Second {
        s.logger.Warn("Slow payment processing", logger.Fields{
            "payment_id": payment.ID,
            "duration":   duration,
            "gateway":    payment.Gateway,
        })
    }
    
    // Retry warning
    if payment.RetryCount > 3 {
        s.logger.Warn("Multiple payment retries", logger.Fields{
            "payment_id":  payment.ID,
            "retry_count": payment.RetryCount,
        })
    }
    
    // Deprecated feature warning
    if payment.UseLegacyAPI {
        s.logger.Warn("Using deprecated payment API", logger.Fields{
            "payment_id": payment.ID,
        })
    }
    
    return nil
}
```

### Error Level

**Use for:** Error conditions that need attention

```go
func (s *service) SendEmail(ctx context.Context, to, subject, body string) error {
    err := s.mailer.Send(to, subject, body)
    if err != nil {
        s.logger.Error("Failed to send email", logger.Fields{
            "error":   err.Error(),
            "to":      to,
            "subject": subject,
        })
        return err
    }
    return nil
}

func (s *service) ProcessBatch(ctx context.Context, items []Item) error {
    var errors []error
    
    for _, item := range items {
        if err := s.processItem(ctx, item); err != nil {
            s.logger.Error("Item processing failed", logger.Fields{
                "error":   err.Error(),
                "item_id": item.ID,
            })
            errors = append(errors, err)
        }
    }
    
    if len(errors) > 0 {
        s.logger.Error("Batch processing completed with errors", logger.Fields{
            "total_items":  len(items),
            "failed_items": len(errors),
        })
        return fmt.Errorf("%d items failed", len(errors))
    }
    
    return nil
}
```

### Fatal Level

**Use for:** Critical errors requiring immediate termination

```go
func main() {
    logger, err := logger.New(config)
    if err != nil {
        panic("Failed to initialize logger")
    }
    
    // Fatal database connection error
    db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
    if err != nil {
        logger.Fatal("Failed to connect to database", logger.Fields{
            "error": err.Error(),
        })
        // Application exits here
    }
    
    // Fatal configuration error
    if config.JWTSecret == "" {
        logger.Fatal("JWT secret not configured", logger.Fields{
            "config_file": configPath,
        })
    }
}
```

---

## Structured Fields

### Basic Fields

```go
// Simple fields
s.logger.Info("User created", logger.Fields{
    "user_id": 123,
    "email":   "user@example.com",
    "name":    "John Doe",
})

// Multiple types
s.logger.Info("Order processed", logger.Fields{
    "order_id":  12345,              // int
    "amount":    99.99,               // float
    "currency":  "USD",               // string
    "completed": true,                // bool
    "items":     3,                   // int
    "timestamp": time.Now(),          // time.Time
})
```

### Complex Fields

```go
// Nested objects
s.logger.Info("Payment processed", logger.Fields{
    "payment": Payment{
        ID:     123,
        Amount: 99.99,
        Status: "completed",
    },
})

// Arrays
s.logger.Info("Batch processed", logger.Fields{
    "ids":      []int{1, 2, 3, 4, 5},
    "statuses": []string{"ok", "ok", "failed", "ok", "ok"},
})

// Maps
s.logger.Info("Metadata", logger.Fields{
    "headers": map[string]string{
        "User-Agent": "Mozilla/5.0",
        "Accept":     "application/json",
    },
})
```

### Error Fields

```go
// Log errors with context
if err != nil {
    s.logger.Error("Operation failed", logger.Fields{
        "error":      err.Error(),
        "error_type": fmt.Sprintf("%T", err),
        "user_id":    userID,
        "action":     "create_order",
    })
}

// Wrap errors with fields
import "github.com/pkg/errors"

if err != nil {
    wrappedErr := errors.Wrap(err, "failed to process payment")
    s.logger.Error("Payment processing error", logger.Fields{
        "error":       wrappedErr.Error(),
        "payment_id":  paymentID,
        "stack_trace": fmt.Sprintf("%+v", wrappedErr),
    })
}
```

---

## Contextual Logging

### Logger with Fields

```go
// Create contextual logger
func (s *service) HandleRequest(ctx context.Context, req *Request) error {
    // Logger with request context
    reqLogger := s.logger.WithFields(logger.Fields{
        "request_id": req.ID,
        "user_id":    req.UserID,
        "ip":         req.IP,
    })
    
    // All logs include request context
    reqLogger.Info("Request started")
    
    if err := s.validateRequest(req); err != nil {
        reqLogger.Error("Validation failed", logger.Fields{
            "error": err.Error(),
        })
        return err
    }
    
    reqLogger.Info("Request validated")
    
    result, err := s.processRequest(req)
    if err != nil {
        reqLogger.Error("Processing failed", logger.Fields{
            "error": err.Error(),
        })
        return err
    }
    
    reqLogger.Info("Request completed", logger.Fields{
        "result": result,
    })
    
    return nil
}
```

### Per-User Logger

```go
// Create user-specific logger
func (s *service) GetUserLogger(userID uint) logger.Logger {
    return s.logger.WithFields(logger.Fields{
        "user_id": userID,
    })
}

// Usage
func (s *service) UpdateProfile(ctx context.Context, userID uint, updates *ProfileUpdates) error {
    userLogger := s.GetUserLogger(userID)
    
    userLogger.Info("Updating profile")
    
    if err := s.repo.Update(ctx, userID, updates); err != nil {
        userLogger.Error("Profile update failed", logger.Fields{
            "error": err.Error(),
        })
        return err
    }
    
    userLogger.Info("Profile updated successfully")
    return nil
}
```

---

## Performance Optimization

### Conditional Logging

```go
// Check if level enabled before expensive operations
func (s *service) ProcessData(ctx context.Context, data *Data) error {
    if s.logger.IsDebugEnabled() {
        // Only compute debug info if debug is enabled
        debugInfo := s.gatherDebugInfo(data) // Expensive
        s.logger.Debug("Processing data", logger.Fields{
            "debug_info": debugInfo,
        })
    }
    
    // ... process data
    
    return nil
}
```

### Lazy Field Evaluation

```go
// Fields are only evaluated if log level is active
s.logger.Debug("Complex data", logger.Fields{
    "data": func() string {
        // This expensive operation only runs if debug is enabled
        return s.serializeComplexData(data)
    }(),
})
```

### Batching

```go
// Batch log events
func (s *service) ProcessBatch(ctx context.Context, items []Item) error {
    processed := 0
    failed := 0
    
    for _, item := range items {
        if err := s.processItem(ctx, item); err != nil {
            failed++
        } else {
            processed++
        }
    }
    
    // Single summary log instead of logging each item
    s.logger.Info("Batch processing complete", logger.Fields{
        "total":     len(items),
        "processed": processed,
        "failed":    failed,
    })
    
    return nil
}
```

---

## Common Patterns

### Request/Response Logging

```go
func (s *service) HandleAPIRequest(ctx context.Context, req *APIRequest) (*APIResponse, error) {
    start := time.Now()
    
    // Log request
    s.logger.Info("API request", logger.Fields{
        "method":     req.Method,
        "path":       req.Path,
        "request_id": req.ID,
    })
    
    // Process
    resp, err := s.process(ctx, req)
    duration := time.Since(start)
    
    // Log response
    if err != nil {
        s.logger.Error("API request failed", logger.Fields{
            "request_id": req.ID,
            "error":      err.Error(),
            "duration":   duration,
        })
        return nil, err
    }
    
    s.logger.Info("API request completed", logger.Fields{
        "request_id":  req.ID,
        "status_code": resp.StatusCode,
        "duration":    duration,
    })
    
    return resp, nil
}
```

### Database Operation Logging

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
            "email":    user.Email,
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

### Background Job Logging

```go
func (s *service) ProcessJob(ctx context.Context, job *Job) error {
    jobLogger := s.logger.WithFields(logger.Fields{
        "job_id":   job.ID,
        "job_type": job.Type,
    })
    
    jobLogger.Info("Job started")
    
    defer func() {
        if r := recover(); r != nil {
            jobLogger.Error("Job panicked", logger.Fields{
                "panic": r,
            })
        }
    }()
    
    if err := s.executeJob(ctx, job); err != nil {
        jobLogger.Error("Job failed", logger.Fields{
            "error": err.Error(),
        })
        return err
    }
    
    jobLogger.Info("Job completed successfully")
    return nil
}
```

### Audit Logging

```go
func (s *service) AuditAction(ctx context.Context, action *Action) {
    s.logger.Info("Audit event", logger.Fields{
        "action":     action.Type,
        "user_id":    action.UserID,
        "resource":   action.Resource,
        "resource_id": action.ResourceID,
        "ip":         action.IP,
        "user_agent": action.UserAgent,
        "timestamp":  action.Timestamp,
    })
}

// Example usage
func (s *service) DeleteUser(ctx context.Context, userID uint) error {
    // ... delete user
    
    s.AuditAction(ctx, &Action{
        Type:       "user.delete",
        UserID:     currentUserID,
        Resource:   "user",
        ResourceID: userID,
        IP:         clientIP,
        Timestamp:  time.Now(),
    })
    
    return nil
}
```

---

## Error Handling

### Logging Errors

```go
// Basic error logging
if err != nil {
    s.logger.Error("Operation failed", logger.Fields{
        "error": err.Error(),
    })
    return err
}

// Error with context
if err != nil {
    s.logger.Error("Failed to create order", logger.Fields{
        "error":   err.Error(),
        "user_id": userID,
        "items":   len(items),
    })
    return fmt.Errorf("create order: %w", err)
}

// Error with stack trace
import "github.com/pkg/errors"

if err != nil {
    wrappedErr := errors.Wrap(err, "payment processing failed")
    s.logger.Error("Payment error", logger.Fields{
        "error": wrappedErr.Error(),
        "stack": fmt.Sprintf("%+v", wrappedErr),
    })
    return wrappedErr
}
```

### Panic Recovery

```go
func (s *service) SafeOperation(ctx context.Context) (err error) {
    defer func() {
        if r := recover(); r != nil {
            s.logger.Error("Panic recovered", logger.Fields{
                "panic": r,
                "stack": string(debug.Stack()),
            })
            err = fmt.Errorf("panic: %v", r)
        }
    }()
    
    // ... operation that might panic
    
    return nil
}
```

---

## Testing

### Mock Logger

```go
// test/mocks/logger.go
package mocks

type MockLogger struct {
    InfoCalls  []LogCall
    ErrorCalls []LogCall
}

type LogCall struct {
    Message string
    Fields  logger.Fields
}

func (m *MockLogger) Info(msg string, fields ...logger.Fields) {
    m.InfoCalls = append(m.InfoCalls, LogCall{
        Message: msg,
        Fields:  mergeFields(fields...),
    })
}

func (m *MockLogger) Error(msg string, fields ...logger.Fields) {
    m.ErrorCalls = append(m.ErrorCalls, LogCall{
        Message: msg,
        Fields:  mergeFields(fields...),
    })
}
```

### Testing with Mock Logger

```go
func TestServiceCreateUser(t *testing.T) {
    mockLogger := &mocks.MockLogger{}
    mockRepo := &mocks.MockRepository{}
    
    service := NewService(mockRepo, mockLogger)
    
    _, err := service.CreateUser(context.Background(), &CreateUserRequest{
        Name:  "John",
        Email: "john@example.com",
    })
    
    require.NoError(t, err)
    
    // Verify logging
    assert.Len(t, mockLogger.InfoCalls, 2)
    assert.Equal(t, "Creating user", mockLogger.InfoCalls[0].Message)
    assert.Equal(t, "User created successfully", mockLogger.InfoCalls[1].Message)
}
```

---

## Best Practices

### ✅ DO:

**1. Use Appropriate Levels**
```go
// Debug for diagnostics
s.logger.Debug("Validating input", logger.Fields{"data": data})

// Info for business events
s.logger.Info("Order created", logger.Fields{"order_id": id})

// Error for failures
s.logger.Error("Payment failed", logger.Fields{"error": err.Error()})
```

**2. Add Context**
```go
s.logger.Info("User action", logger.Fields{
    "user_id": userID,
    "action":  "update_profile",
    "ip":      clientIP,
})
```

**3. Log Start and End**
```go
s.logger.Info("Job started", logger.Fields{"job_id": id})
// ... work
s.logger.Info("Job completed", logger.Fields{"job_id": id, "duration": duration})
```

### ❌ DON'T:

**1. Log Sensitive Data**
```go
// Bad
s.logger.Info("User login", logger.Fields{
    "password": password, // ❌ Never log passwords
    "ssn":      ssn,      // ❌ Never log PII
})

// Good
s.logger.Info("User login", logger.Fields{
    "user_id": userID,
})
```

**2. Excessive Logging**
```go
// Bad
for _, item := range items {
    s.logger.Debug("Processing", logger.Fields{"item": item}) // ❌
}

// Good
s.logger.Debug("Processing batch", logger.Fields{"count": len(items)})
```

**3. String Concatenation**
```go
// Bad
s.logger.Info("User " + userID + " created") // ❌

// Good
s.logger.Info("User created", logger.Fields{"user_id": userID})
```

---

## Next Steps

- [**Configuration**](configuration.md) - Logger setup
- [**Middleware**](middleware.md) - HTTP logging
- [**Best Practices**](best-practices.md) - Logging guidelines
- [**Overview**](overview.md) - Logger architecture

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
