# Debugging Guide

Learn effective debugging techniques for Neonex Core applications.

---

## Overview

Debugging strategies for Neonex Core:
- **Logging** - Structured debug logs
- **Delve** - Go debugger
- **VS Code** - Integrated debugging
- **pprof** - Performance profiling

---

## Debug Logging

### Enable Debug Level

```bash
# .env
LOG_LEVEL=debug
```

### Debug Logs

```go
// Add debug logs
logger.Debug("Processing request", logger.Fields{
    "user_id": userID,
    "payload":  payload,
})
```

---

## Delve Debugger

### Install Delve

```powershell
go install github.com/go-delve/delve/cmd/dlv@latest
```

### Start Debugging

```powershell
# Debug main package
dlv debug ./cmd/neonex

# Set breakpoint
(dlv) break main.main
(dlv) continue

# Inspect variables
(dlv) print variableName
(dlv) locals
```

---

## VS Code Debugging

### Launch Configuration

**.vscode/launch.json:**
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/neonex",
      "env": {
        "ENV": "development"
      },
      "args": []
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

### Set Breakpoints

1. Click left margin in editor
2. Press F5 to start debugging
3. Step through code

---

## Common Issues

### Database Connection Errors

```go
// Add debug logging
logger.Debug("Connecting to database", logger.Fields{
    "host": config.DBHost,
    "port": config.DBPort,
    "database": config.DBName,
})

db, err := gorm.Open(...)
if err != nil {
    logger.Error("Connection failed", logger.Fields{
        "error": err.Error(),
    })
}
```

### HTTP Request Issues

```go
// Log request details
logger.Debug("HTTP Request", logger.Fields{
    "method": c.Method(),
    "path":   c.Path(),
    "body":   string(c.Body()),
    "headers": c.GetReqHeaders(),
})
```

---

## Performance Profiling

### CPU Profiling

```go
import _ "net/http/pprof"

// Start pprof server
go func() {
    http.ListenAndServe("localhost:6060", nil)
}()
```

**Access:** http://localhost:6060/debug/pprof/

### Memory Profiling

```powershell
# Capture heap profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Analyze
(pprof) top
(pprof) list functionName
```

---

## Error Handling

### Stack Traces

```go
import (
    "github.com/pkg/errors"
    "runtime/debug"
)

// Wrap errors with stack trace
if err != nil {
    return errors.Wrap(err, "operation failed")
}

// Capture stack on panic
defer func() {
    if r := recover(); r != nil {
        logger.Error("Panic", logger.Fields{
            "panic": r,
            "stack": string(debug.Stack()),
        })
    }
}()
```

---

## Best Practices

### ✅ DO:
- Use debug logging liberally
- Set strategic breakpoints
- Inspect variable state
- Profile performance bottlenecks

### ❌ DON'T:
- Leave debug logs in production
- Debug without context
- Ignore error messages

---

## Next Steps

- [**Testing**](testing.md) - Testing strategies
- [**Hot Reload**](hot-reload.md) - Development workflow
- [**Logging**](../logging/overview.md) - Logger configuration

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
