# Error Handling

Comprehensive error handling patterns in Neonex Core.

---

## Custom Errors

```go
var (
    ErrNotFound = errors.New("resource not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrValidation = errors.New("validation error")
)
```

## Error Responses

```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
    Code    string `json:"code"`
}

func HandleError(c *fiber.Ctx, err error) error {
    if errors.Is(err, ErrNotFound) {
        return c.Status(404).JSON(ErrorResponse{
            Error: "Not Found",
            Message: err.Error(),
            Code: "NOT_FOUND",
        })
    }
    
    return c.Status(500).JSON(ErrorResponse{
        Error: "Internal Server Error",
        Message: err.Error(),
        Code: "INTERNAL_ERROR",
    })
}
```

## Validation Errors

```go
func ValidateUser(user *User) error {
    if user.Email == "" {
        return fmt.Errorf("%w: email is required", ErrValidation)
    }
    return nil
}
```

---

## Next Steps

- [**Middleware**](middleware.md)
- [**Performance**](performance.md)
- [**Security**](security.md)