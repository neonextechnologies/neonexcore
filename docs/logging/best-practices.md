# Logging Best Practices

Guidelines and recommendations for effective logging in Neonex Core applications.

---

## General Principles

### 1. Log Meaningful Events

**✅ DO:**
```go
// Important business events
logger.Info("Order created", logger.Fields{
    "order_id":   order.ID,
    "user_id":    user.ID,
    "total":      order.Total,
})

// Significant state changes
logger.Info("User account activated", logger.Fields{
    "user_id": user.ID,
})

// External service calls
logger.Info("Payment processed", logger.Fields{
    "payment_id": payment.ID,
    "gateway":    "stripe",
    "status":     "completed",
})
```

**❌ DON'T:**
```go
// Trivial operations
logger.Info("Setting variable") // ❌ Too granular

// Noise without value
logger.Debug("Entering function") // ❌ Not useful
logger.Debug("Exiting function")  // ❌ Not useful
```

---

## Log Levels

### Use Appropriate Levels

| Level | When to Use | Example |
|-------|-------------|---------|
| **Debug** | Development diagnostics | Variable values, flow tracing |
| **Info** | Normal operations | Business events, milestones |
| **Warn** | Potential issues | Slow queries, deprecated API usage |
| **Error** | Operation failures | Failed requests, exceptions |
| **Fatal** | Critical failures | Cannot start app, missing config |

**✅ DO:**
```go
// Debug: Detailed diagnostics
logger.Debug("Validating input", logger.Fields{
    "input": input,
})

// Info: Business events
logger.Info("User logged in", logger.Fields{
    "user_id": user.ID,
})

// Warn: Concerning but not critical
logger.Warn("Slow database query", logger.Fields{
    "query":    query,
    "duration": duration,
})

// Error: Operation failed
logger.Error("Payment processing failed", logger.Fields{
    "error":      err.Error(),
    "payment_id": payment.ID,
})

// Fatal: Critical system failure
logger.Fatal("Database connection failed", logger.Fields{
    "error": err.Error(),
})
```

**❌ DON'T:**
```go
// Wrong level usage
logger.Error("User logged in")      // ❌ Not an error
logger.Info("Database crashed")     // ❌ Should be Fatal
logger.Debug("Payment failed")      // ❌ Should be Error
```

---

## Structured Logging

### Always Use Fields

**✅ DO:**
```go
// Structured with fields
logger.Info("User action", logger.Fields{
    "user_id": 123,
    "action":  "profile_update",
    "ip":      "192.168.1.1",
})
```

**❌ DON'T:**
```go
// String interpolation
logger.Info(fmt.Sprintf("User %d performed %s from %s", 
    123, "profile_update", "192.168.1.1")) // ❌

// Concatenation
logger.Info("User " + userID + " logged in") // ❌
```

### Consistent Field Names

**✅ DO:**
```go
// Use consistent naming convention
logger.Info("Event", logger.Fields{
    "user_id":    123,        // snake_case
    "order_id":   456,        // Consistent
    "request_id": "abc-123",  // Consistent
})
```

**❌ DON'T:**
```go
// Inconsistent naming
logger.Info("Event", logger.Fields{
    "userId":     123,        // ❌ camelCase
    "order-id":   456,        // ❌ kebab-case
    "RequestID":  "abc-123",  // ❌ PascalCase
})
```

---

## Security

### Never Log Sensitive Data

**✅ DO:**
```go
// Omit sensitive fields
logger.Info("User created", logger.Fields{
    "user_id": user.ID,
    "email":   user.Email,
})

// Redact sensitive data
logger.Info("Payment processed", logger.Fields{
    "card_last4": payment.CardLast4,
    "amount":     payment.Amount,
})
```

**❌ DON'T:**
```go
// Logging sensitive data
logger.Info("User login", logger.Fields{
    "password": password,        // ❌ NEVER log passwords
    "token":    authToken,       // ❌ NEVER log auth tokens
    "ssn":      user.SSN,        // ❌ NEVER log PII
    "card_number": cardNumber,   // ❌ NEVER log credit cards
})
```

### Sanitize User Input

**✅ DO:**
```go
// Truncate or sanitize
func sanitizeForLog(input string) string {
    if len(input) > 100 {
        return input[:100] + "..."
    }
    return input
}

logger.Info("Search query", logger.Fields{
    "query": sanitizeForLog(userInput),
})
```

**❌ DON'T:**
```go
// Raw user input
logger.Info("Search", logger.Fields{
    "query": userInput, // ❌ Could contain sensitive data
})
```

---

## Performance

### Avoid Expensive Operations

**✅ DO:**
```go
// Conditional logging for expensive operations
if logger.IsDebugEnabled() {
    debugData := gatherExpensiveDebugInfo() // Only if needed
    logger.Debug("Debug info", logger.Fields{
        "data": debugData,
    })
}

// Lazy evaluation
logger.Debug("User data", logger.Fields{
    "user": func() string {
        return user.ToJSON() // Only called if debug enabled
    }(),
})
```

**❌ DON'T:**
```go
// Always computing expensive data
debugData := gatherExpensiveDebugInfo() // ❌ Computed even if not logged
logger.Debug("Debug info", logger.Fields{
    "data": debugData,
})
```

### Batch When Possible

**✅ DO:**
```go
// Log summary instead of individual items
processed := 0
failed := 0

for _, item := range items {
    if processItem(item) {
        processed++
    } else {
        failed++
    }
}

logger.Info("Batch complete", logger.Fields{
    "total":     len(items),
    "processed": processed,
    "failed":    failed,
})
```

**❌ DON'T:**
```go
// Log every iteration
for i, item := range items {
    logger.Info("Processing", logger.Fields{ // ❌ Too many logs
        "index": i,
        "item":  item,
    })
}
```

---

## Context and Traceability

### Use Request IDs

**✅ DO:**
```go
// Create contextual logger with request ID
reqLogger := logger.WithFields(logger.Fields{
    "request_id": requestID,
})

// All logs automatically include request_id
reqLogger.Info("Request started")
reqLogger.Info("Validating input")
reqLogger.Info("Request completed")
```

### Add Contextual Information

**✅ DO:**
```go
// Include relevant context
logger.Error("Order processing failed", logger.Fields{
    "error":     err.Error(),
    "order_id":  order.ID,
    "user_id":   order.UserID,
    "status":    order.Status,
    "timestamp": time.Now(),
})
```

**❌ DON'T:**
```go
// Minimal context
logger.Error("Failed", logger.Fields{ // ❌ What failed?
    "error": err.Error(),
})
```

---

## Error Logging

### Log Errors with Context

**✅ DO:**
```go
if err != nil {
    logger.Error("Failed to create user", logger.Fields{
        "error":  err.Error(),
        "email":  req.Email,
        "name":   req.Name,
    })
    return err
}
```

### Include Stack Traces for Critical Errors

**✅ DO:**
```go
import "github.com/pkg/errors"

if err != nil {
    wrappedErr := errors.Wrap(err, "payment processing failed")
    logger.Error("Critical error", logger.Fields{
        "error": wrappedErr.Error(),
        "stack": fmt.Sprintf("%+v", wrappedErr),
    })
}
```

### Log Recoveries from Panics

**✅ DO:**
```go
func (s *service) SafeOperation(ctx context.Context) (err error) {
    defer func() {
        if r := recover(); r != nil {
            logger.Error("Panic recovered", logger.Fields{
                "panic": r,
                "stack": string(debug.Stack()),
            })
            err = fmt.Errorf("panic: %v", r)
        }
    }()
    
    // ... operation
    return nil
}
```

---

## Message Quality

### Clear and Descriptive Messages

**✅ DO:**
```go
logger.Info("User account activated successfully")
logger.Error("Failed to connect to payment gateway")
logger.Warn("Database query exceeded timeout threshold")
```

**❌ DON'T:**
```go
logger.Info("Done")                    // ❌ Too vague
logger.Error("Error")                  // ❌ Not descriptive
logger.Warn("Something is wrong")      // ❌ Not specific
```

### Use Consistent Tense and Voice

**✅ DO:**
```go
// Past tense for completed actions
logger.Info("User created")
logger.Info("Payment processed")

// Present continuous for ongoing actions
logger.Info("Processing order")
logger.Info("Connecting to database")
```

---

## Environment-Specific Logging

### Development

```bash
# .env.development
LOG_LEVEL=debug
LOG_FORMAT=console
LOG_OUTPUT=stdout
LOG_ENABLE_CALLER=true
```

**Characteristics:**
- Verbose (debug level)
- Human-readable format
- Caller information
- Console output

### Staging

```bash
# .env.staging
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=file
LOG_FILE_PATH=logs/app.log
```

**Characteristics:**
- Moderate verbosity (info level)
- JSON format
- File output with rotation
- Similar to production

### Production

```bash
# .env.production
LOG_LEVEL=warn
LOG_FORMAT=json
LOG_OUTPUT=file
LOG_FILE_PATH=/var/log/neonex/app.log
LOG_FILE_MAX_SIZE=100
LOG_FILE_MAX_BACKUPS=10
LOG_ENABLE_STACKTRACE=true
LOG_SAMPLING_ENABLED=true
```

**Characteristics:**
- Minimal verbosity (warn/error only)
- JSON format (parsing)
- File output with rotation
- Stack traces for errors
- Sampling for high volume

---

## Testing

### Mock Logger for Tests

**✅ DO:**
```go
// Create mock logger
type mockLogger struct {
    logs []LogEntry
}

func (m *mockLogger) Info(msg string, fields ...logger.Fields) {
    m.logs = append(m.logs, LogEntry{
        Level:   "info",
        Message: msg,
        Fields:  fields[0],
    })
}

// Test with mock
func TestService(t *testing.T) {
    mockLog := &mockLogger{}
    service := NewService(mockLog)
    
    service.DoSomething()
    
    assert.Len(t, mockLog.logs, 1)
    assert.Equal(t, "Action completed", mockLog.logs[0].Message)
}
```

### Verify Important Logs

**✅ DO:**
```go
func TestUserCreation(t *testing.T) {
    mockLog := &mockLogger{}
    service := NewService(repo, mockLog)
    
    user, err := service.CreateUser(ctx, req)
    
    require.NoError(t, err)
    
    // Verify success log
    assert.Contains(t, mockLog.logs, LogEntry{
        Level:   "info",
        Message: "User created successfully",
        Fields: logger.Fields{
            "user_id": user.ID,
        },
    })
}
```

---

## Monitoring Integration

### Log for Monitoring Systems

**✅ DO:**
```go
// Structured logs for monitoring
logger.Info("API Response", logger.Fields{
    "endpoint":    "/api/users",
    "status_code": 200,
    "duration_ms": 45,
    "user_id":     123,
})

// Metric-like logs
logger.Info("Cache Hit", logger.Fields{
    "cache_key":  key,
    "hit":        true,
    "ttl_seconds": 300,
})
```

### Use Standard Metric Names

**✅ DO:**
```go
// Consistent metric naming
logger.Info("HTTP Request", logger.Fields{
    "http.method":      "POST",
    "http.status_code": 201,
    "http.duration_ms": 45,
    "http.path":        "/api/users",
})
```

---

## Log Rotation and Retention

### Configure Rotation

**✅ DO:**
```bash
# Rotate logs to prevent disk fill
LOG_FILE_MAX_SIZE=100        # MB
LOG_FILE_MAX_BACKUPS=10      # Keep 10 old files
LOG_FILE_MAX_AGE=30          # Delete after 30 days
LOG_FILE_COMPRESS=true       # Compress old logs
```

### Regular Cleanup

**✅ DO:**
```go
// Periodic log cleanup (via cron or systemd timer)
// Example: Keep logs for 30 days
find /var/log/neonex -name "*.log*" -mtime +30 -delete
```

---

## Compliance and Audit

### Audit Logging

**✅ DO:**
```go
// Log sensitive operations for audit
logger.Info("Audit: User role changed", logger.Fields{
    "actor_id":    adminID,
    "target_id":   userID,
    "old_role":    oldRole,
    "new_role":    newRole,
    "timestamp":   time.Now(),
    "ip_address":  clientIP,
})

// Log data access
logger.Info("Audit: PII accessed", logger.Fields{
    "user_id":     userID,
    "accessed_by": adminID,
    "resource":    "user_profile",
    "action":      "read",
})
```

---

## Quick Reference

### DO's

✅ Use appropriate log levels  
✅ Use structured logging (fields)  
✅ Include request IDs for traceability  
✅ Log errors with context  
✅ Use consistent field names  
✅ Configure log rotation  
✅ Use environment-specific configs  
✅ Batch logs when possible  
✅ Write clear, descriptive messages  
✅ Test important log statements  

### DON'Ts

❌ Log sensitive data (passwords, tokens, PII)  
❌ Log in loops without batching  
❌ Use string concatenation for messages  
❌ Use inconsistent field naming  
❌ Log at wrong levels  
❌ Write vague messages  
❌ Log trivial operations  
❌ Ignore performance impact  
❌ Log without context  
❌ Forget to configure rotation  

---

## Common Patterns

### Pattern 1: Operation Logging

```go
func (s *service) ComplexOperation(ctx context.Context, req *Request) error {
    opLogger := s.logger.WithFields(logger.Fields{
        "operation":  "complex_operation",
        "request_id": req.ID,
    })
    
    opLogger.Info("Operation started")
    
    start := time.Now()
    
    // Step 1
    opLogger.Debug("Step 1: Validation")
    if err := s.validate(req); err != nil {
        opLogger.Error("Validation failed", logger.Fields{
            "error": err.Error(),
        })
        return err
    }
    
    // Step 2
    opLogger.Debug("Step 2: Processing")
    result, err := s.process(req)
    if err != nil {
        opLogger.Error("Processing failed", logger.Fields{
            "error": err.Error(),
        })
        return err
    }
    
    duration := time.Since(start)
    opLogger.Info("Operation completed", logger.Fields{
        "duration": duration,
        "result":   result,
    })
    
    return nil
}
```

### Pattern 2: HTTP Request Logging

```go
func (ctrl *controller) HandleRequest(c *fiber.Ctx) error {
    requestID := c.Locals("request_id").(string)
    
    reqLogger := ctrl.logger.WithFields(logger.Fields{
        "request_id": requestID,
        "endpoint":   c.Path(),
        "method":     c.Method(),
    })
    
    reqLogger.Info("Request received")
    
    result, err := ctrl.service.Process(c.Context())
    if err != nil {
        reqLogger.Error("Request failed", logger.Fields{
            "error": err.Error(),
        })
        return err
    }
    
    reqLogger.Info("Request succeeded")
    return c.JSON(result)
}
```

### Pattern 3: Background Job Logging

```go
func (w *worker) ProcessJob(ctx context.Context, job *Job) error {
    jobLogger := w.logger.WithFields(logger.Fields{
        "job_id":   job.ID,
        "job_type": job.Type,
    })
    
    jobLogger.Info("Job started")
    
    defer func() {
        if r := recover(); r != nil {
            jobLogger.Error("Job panicked", logger.Fields{
                "panic": r,
                "stack": string(debug.Stack()),
            })
        }
    }()
    
    if err := w.execute(ctx, job); err != nil {
        jobLogger.Error("Job failed", logger.Fields{
            "error": err.Error(),
        })
        return err
    }
    
    jobLogger.Info("Job completed")
    return nil
}
```

---

## Next Steps

- [**Overview**](overview.md) - Logger architecture
- [**Configuration**](configuration.md) - Logger setup
- [**Usage**](usage.md) - Logging in your code
- [**Middleware**](middleware.md) - HTTP logging

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
