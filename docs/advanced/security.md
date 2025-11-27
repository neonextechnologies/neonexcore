# Security Best Practices

Security guidelines for Neonex Core applications.

---

## Authentication

### JWT Authentication

```go
import "github.com/golang-jwt/jwt/v5"

func GenerateToken(userID uint) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    })
    
    return token.SignedString([]byte(jwtSecret))
}

func ValidateToken(tokenString string) (uint, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte(jwtSecret), nil
    })
    
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return uint(claims["user_id"].(float64)), nil
    }
    
    return 0, err
}
```

---

## Password Hashing

```go
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

---

## Input Validation

```go
import "github.com/go-playground/validator/v10"

type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Name     string `json:"name" validate:"required"`
}

func ValidateRequest(req interface{}) error {
    validate := validator.New()
    return validate.Struct(req)
}
```

---

## SQL Injection Prevention

```go
// Use parameterized queries
db.Where("email = ?", email).First(&user)

// Never concatenate strings
// ❌ BAD: db.Where("email = '" + email + "'").First(&user)
```

---

## XSS Prevention

```go
import "html"

func SanitizeInput(input string) string {
    return html.EscapeString(input)
}
```

---

## CSRF Protection

```go
import "github.com/gofiber/fiber/v2/middleware/csrf"

app.Use(csrf.New())
```

---

## Rate Limiting

```go
import "github.com/gofiber/fiber/v2/middleware/limiter"

app.Use(limiter.New(limiter.Config{
    Max: 100,
    Expiration: 1 * time.Minute,
}))
```

---

## HTTPS Configuration

```go
// Use HTTPS in production
app.ListenTLS(":443", "./cert.pem", "./key.pem")

// Redirect HTTP to HTTPS
app.Use(func(c *fiber.Ctx) error {
    if c.Protocol() != "https" {
        return c.Redirect("https://" + c.Hostname() + c.OriginalURL())
    }
    return c.Next()
})
```

---

## Environment Variables

```bash
# Never commit secrets
JWT_SECRET=your-secret-key
DB_PASSWORD=your-db-password
API_KEY=your-api-key
```

---

## Best Practices

### ✅ DO:
- Hash passwords with bcrypt
- Use HTTPS in production
- Validate all input
- Use parameterized queries
- Implement rate limiting
- Keep dependencies updated

### ❌ DON'T:
- Store passwords in plaintext
- Trust user input
- Expose error details
- Hardcode secrets
- Use weak encryption

---

## Next Steps

- [**Middleware**](middleware.md)
- [**Error Handling**](error-handling.md)
- [**Performance**](performance.md)

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
