# Container API

Dependency injection container reference.

---

## Functions

```go
// RegisterSingleton - Single instance
func RegisterSingleton[T any](c *Container, name string, factory func() T)

// RegisterTransient - New instance each time
func RegisterTransient[T any](c *Container, name string, factory func() T)

// Resolve - Get service
func Resolve[T any](c *Container, name string) T
```

## Usage

```go
// Register
RegisterSingleton(container, "database", func() *gorm.DB {
    return setupDB()
})

// Resolve
db := Resolve[*gorm.DB](container, "database")
```

---

## Next Steps

- [**Core**](core.md)
- [**Module Interface**](module-interface.md)
