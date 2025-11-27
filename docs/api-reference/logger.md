# Logger API

Logging system reference.

---

## Interface

```go
type Logger interface {
    Debug(msg string, fields ...Fields)
    Info(msg string, fields ...Fields)
    Warn(msg string, fields ...Fields)
    Error(msg string, fields ...Fields)
    Fatal(msg string, fields ...Fields)
    WithFields(fields Fields) Logger
}

type Fields map[string]interface{}
```

## Usage

```go
logger.Info(\"User created\", logger.Fields{
    \"user_id\": 123,
    \"email\": \"user@example.com\",
})
```

## Contextual Logger

```go
userLogger := logger.WithFields(logger.Fields{
    \"user_id\": userID,
})

userLogger.Info(\"Action performed\")
```

---

## Next Steps

- [**Core**](core.md)
- [**Logging**](../logging/overview.md)
