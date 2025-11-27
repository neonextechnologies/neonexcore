# Module Interface

Module system API reference.

---

## Module Interface

```go
type Module interface {
    Name() string
    Version() string
    RegisterDependencies(*Container)
    RegisterRoutes(fiber.Router)
}
```

## Implementation

```go
type ProductModule struct{}

func (m *ProductModule) Name() string {
    return "product"
}

func (m *ProductModule) Version() string {
    return "1.0.0"
}

func (m *ProductModule) RegisterDependencies(c *Container) {
    // Register services
}

func (m *ProductModule) RegisterRoutes(router fiber.Router) {
    // Register routes
}
```

---

## Next Steps

- [**Core**](core.md)
- [**Container**](container.md)
