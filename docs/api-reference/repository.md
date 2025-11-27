# Repository API

Generic repository pattern reference.

---

## Interface

```go
type Repository[T any] interface {
    Create(ctx context.Context, entity *T) error
    FindByID(ctx context.Context, id uint) (*T, error)
    FindAll(ctx context.Context) ([]*T, error)
    Update(ctx context.Context, entity *T) error
    Delete(ctx context.Context, id uint) error
}
```

## Base Implementation

```go
type BaseRepository[T any] struct {
    db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
    return BaseRepository[T]{db: db}
}
```

## Custom Repository

```go
type UserRepository interface {
    Repository[User]
    FindByEmail(ctx context.Context, email string) (*User, error)
}
```

---

## Next Steps

- [**Core**](core.md)
- [**Repository Pattern**](../core-concepts/repository-pattern.md)
