# Core API Reference

API documentation for Neonex Core framework.

---

## App

```go
type App struct {
    Container *Container
    Router    fiber.Router
}

// NewApp creates new application instance
func NewApp() *App

// Run starts the application
func (a *App) Run() error
```

---

## Container

```go
type Container struct{}

// RegisterSingleton registers singleton service
func RegisterSingleton[T any](c *Container, name string, factory func() T)

// RegisterTransient registers transient service  
func RegisterTransient[T any](c *Container, name string, factory func() T)

// Resolve retrieves service
func Resolve[T any](c *Container, name string) T
```

---

## Module

```go
type Module interface {
    Name() string
    Version() string
    RegisterDependencies(c *Container)
    RegisterRoutes(router fiber.Router)
}
```

---

## Repository

```go
type Repository[T any] interface {
    Create(ctx context.Context, entity *T) error
    FindByID(ctx context.Context, id uint) (*T, error)
    FindAll(ctx context.Context) ([]*T, error)
    Update(ctx context.Context, entity *T) error
    Delete(ctx context.Context, id uint) error
}
```

---

## Logger

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

---

## Next Steps

- [**Module Interface**](module-interface.md)
- [**Container**](container.md)
- [**Repository**](repository.md)
- [**Logger**](logger.md)
