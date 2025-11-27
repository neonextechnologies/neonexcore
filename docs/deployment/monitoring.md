# Monitoring

Application monitoring and observability guide.

---

## Health Checks

```go
app.Get(\"/health\", func(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        \"status\": \"ok\",
        \"timestamp\": time.Now(),
    })
})

app.Get(\"/ready\", func(c *fiber.Ctx) error {
    // Check dependencies
    if err := db.Ping(); err != nil {
        return c.Status(503).JSON(fiber.Map{
            \"status\": \"not ready\",
        })
    }
    
    return c.JSON(fiber.Map{
        \"status\": \"ready\",
    })
})
```

---

## Prometheus Metrics

```go
import \"github.com/gofiber/fiber/v2/middleware/monitor\"

app.Get(\"/metrics\", monitor.New())
```

---

## Logging

```go
logger.Info(\"Request metrics\", logger.Fields{
    \"endpoint\": c.Path(),
    \"method\": c.Method(),
    \"status\": c.Response().StatusCode(),
    \"duration_ms\": duration.Milliseconds(),
})
```

---

## Error Tracking

```go
// Sentry integration
import \"github.com/getsentry/sentry-go\"

sentry.Init(sentry.ClientOptions{
    Dsn: os.Getenv(\"SENTRY_DSN\"),
})

// Capture errors
if err != nil {
    sentry.CaptureException(err)
}
```

---

## Next Steps

- [**Production Setup**](production-setup.md)
- [**Docker**](docker.md)
- [**Environment Variables**](environment-variables.md)
