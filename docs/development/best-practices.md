# Development Best Practices

Guidelines and recommendations for developing with Neonex Core.

---

## Project Organization

### Module Structure

```
modules/
└── user/
    ├── model.go          # Data models
    ├── repository.go     # Data access
    ├── service.go        # Business logic
    ├── controller.go     # HTTP handlers
    ├── routes.go         # Route registration
    ├── di.go            # Dependency injection
    └── module.json      # Module metadata
```

### Clean Architecture

- **Models** - Data structures only
- **Repository** - Database operations
- **Service** - Business logic
- **Controller** - HTTP/API layer

---

## Code Quality

### Formatting

```powershell
# Format code
go fmt ./...

# Imports organization
goimports -w .
```

### Linting

```powershell
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

---

## Error Handling

### Wrap Errors

```go
if err != nil {
    return fmt.Errorf("create user: %w", err)
}
```

### Custom Errors

```go
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidEmail = errors.New("invalid email")
)

// Usage
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return ErrUserNotFound
    }
    return err
}
```

---

## Configuration

### Environment Variables

```bash
# Development
ENV=development
LOG_LEVEL=debug

# Production
ENV=production
LOG_LEVEL=info
```

### Configuration Struct

```go
type Config struct {
    Env      string `env:"ENV" envDefault:"development"`
    Port     string `env:"PORT" envDefault:"3000"`
    LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}
```

---

## Testing

### Test Coverage

```powershell
# Aim for 80%+ coverage
go test -cover ./...
```

### Test Organization

```
modules/user/
├── service.go
├── service_test.go       # Unit tests
└── integration_test.go   # Integration tests
```

---

## Documentation

### Code Comments

```go
// CreateUser creates a new user account.
// Returns ErrInvalidEmail if email format is invalid.
func (s *service) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // Implementation
}
```

### README Files

```
modules/user/
└── README.md  # Module documentation
```

---

## Version Control

### Commit Messages

```
feat: add user authentication
fix: resolve database connection issue
docs: update API documentation
refactor: simplify service layer
test: add integration tests
```

### .gitignore

```
# Binaries
*.exe
*.dll
*.so

# Temp files
tmp/
.env
.env.local

# IDE
.vscode/
.idea/
```

---

## Performance

### Database Optimization

- Use indexes
- Preload associations
- Batch operations
- Connection pooling

### HTTP Performance

- Use compression
- Implement caching
- Rate limiting
- Connection reuse

---

## Security

### Best Practices

- Validate all input
- Sanitize user data
- Use prepared statements
- Implement rate limiting
- Hash passwords (bcrypt)
- Use HTTPS in production

---

## Monitoring

### Health Checks

```go
app.Get("/health", func(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status": "ok",
    })
})
```

### Metrics

- Request latency
- Error rates
- Database connections
- Memory usage

---

## Next Steps

- [**Testing**](testing.md) - Testing guide
- [**Debugging**](debugging.md) - Debug techniques
- [**Hot Reload**](hot-reload.md) - Development workflow

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
