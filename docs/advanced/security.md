# Security

Secure your Neonex Core applications against common vulnerabilities.

---

## Overview

Security best practices for:
- **Authentication** - Verify user identity
- **Authorization** - Control access
- **Input Validation** - Prevent injection attacks
- **Data Protection** - Encrypt sensitive data
- **HTTPS** - Secure communication
- **CORS** - Control cross-origin requests
- **Rate Limiting** - Prevent abuse

---

## Authentication

### Password Hashing

```go
import "golang.org/x/crypto/bcrypt"

// Hash password
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// Verify password
func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// Usage
func CreateUser(user *User) error {
    hashedPassword, err := HashPassword(user.Password)
    if err != nil {
        return err
    }
    
    user.Password = hashedPassword
    return repo.Create(user)
}

func Login(email, password string) (*User, error) {
    user, err := repo.FindByEmail(email)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }
    
    if !CheckPassword(password, user.Password) {
        return nil, errors.New("invalid credentials")
    }
    
    return user, nil
}
```

### JWT Authentication

```go
import "github.com/golang-jwt/jwt/v5"

type Claims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Generate token
func GenerateToken(user *User) (string, error) {
    claims := Claims{
        UserID: user.ID,
        Email:  user.Email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// Validate token
func ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("invalid token")
}

// Middleware
func JWTMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(401).JSON(fiber.Map{
                "error": "Missing authorization header",
            })
        }
        
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        
        claims, err := ValidateToken(tokenString)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{
                "error": "Invalid token",
            })
        }
        
        c.Locals("userID", claims.UserID)
        c.Locals("email", claims.Email)
        
        return c.Next()
    }
}
```

---

## Authorization

### Role-Based Access Control (RBAC)

```go
type Role string

const (
    RoleAdmin     Role = "admin"
    RoleModerator Role = "moderator"
    RoleUser      Role = "user"
)

type User struct {
    ID    uint   `gorm:"primarykey"`
    Email string `gorm:"uniqueIndex"`
    Role  Role   `gorm:"type:varchar(20)"`
}

// Middleware
func RequireRole(allowedRoles ...Role) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userID := c.Locals("userID").(uint)
        
        user, err := userService.GetUser(userID)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{
                "error": "Unauthorized",
            })
        }
        
        for _, role := range allowedRoles {
            if user.Role == role {
                return c.Next()
            }
        }
        
        return c.Status(403).JSON(fiber.Map{
            "error": "Insufficient permissions",
        })
    }
}

// Usage
app.Delete("/users/:id",
    JWTMiddleware(),
    RequireRole(RoleAdmin),
    deleteUser,
)
```

### Permission-Based Access Control

```go
type Permission string

const (
    PermissionReadUsers   Permission = "read:users"
    PermissionWriteUsers  Permission = "write:users"
    PermissionDeleteUsers Permission = "delete:users"
)

type User struct {
    ID          uint         `gorm:"primarykey"`
    Permissions []Permission `gorm:"type:json"`
}

func HasPermission(user *User, perm Permission) bool {
    for _, p := range user.Permissions {
        if p == perm {
            return true
        }
    }
    return false
}

func RequirePermission(perm Permission) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userID := c.Locals("userID").(uint)
        user, _ := userService.GetUser(userID)
        
        if !HasPermission(user, perm) {
            return c.Status(403).JSON(fiber.Map{
                "error": "Permission denied",
            })
        }
        
        return c.Next()
    }
}
```

---

## Input Validation

### Validate Request Data

```go
import "github.com/go-playground/validator/v10"

type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8,max=100"`
    Age      int    `json:"age" validate:"required,gte=18,lte=120"`
}

var validate = validator.New()

func CreateUser(c *fiber.Ctx) error {
    var req CreateUserRequest
    
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    
    if err := validate.Struct(req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    // Proceed with validated data
    return nil
}
```

### Sanitize Input

```go
import (
    "html"
    "strings"
)

func SanitizeString(s string) string {
    // Remove leading/trailing whitespace
    s = strings.TrimSpace(s)
    
    // Escape HTML
    s = html.EscapeString(s)
    
    return s
}

func CreateUser(c *fiber.Ctx) error {
    var user User
    c.BodyParser(&user)
    
    // Sanitize inputs
    user.Name = SanitizeString(user.Name)
    user.Bio = SanitizeString(user.Bio)
    
    return service.CreateUser(&user)
}
```

### Prevent SQL Injection

```go
// ✅ Good: Parameterized queries (GORM does this automatically)
db.Where("email = ?", email).First(&user)

// ❌ Bad: String concatenation
query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)
db.Raw(query).Scan(&user)  // SQL injection vulnerability!

// ✅ Good: Named parameters
db.Where("email = @email AND status = @status", map[string]interface{}{
    "email":  email,
    "status": "active",
}).Find(&users)
```

---

## Data Protection

### Encrypt Sensitive Data

```go
import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "io"
)

var encryptionKey = []byte(os.Getenv("ENCRYPTION_KEY")) // 32 bytes for AES-256

func Encrypt(plaintext string) (string, error) {
    block, err := aes.NewCipher(encryptionKey)
    if err != nil {
        return "", err
    }
    
    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }
    
    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))
    
    return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(ciphertext string) (string, error) {
    data, err := base64.URLEncoding.DecodeString(ciphertext)
    if err != nil {
        return "", err
    }
    
    block, err := aes.NewCipher(encryptionKey)
    if err != nil {
        return "", err
    }
    
    if len(data) < aes.BlockSize {
        return "", errors.New("ciphertext too short")
    }
    
    iv := data[:aes.BlockSize]
    data = data[aes.BlockSize:]
    
    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(data, data)
    
    return string(data), nil
}

// Usage for sensitive fields
type User struct {
    ID          uint
    Email       string
    SSN         string `gorm:"-"` // Don't store directly
    SSNEncrypted string `gorm:"column:ssn_encrypted"`
}

func (u *User) BeforeSave(tx *gorm.DB) error {
    if u.SSN != "" {
        encrypted, err := Encrypt(u.SSN)
        if err != nil {
            return err
        }
        u.SSNEncrypted = encrypted
    }
    return nil
}

func (u *User) AfterFind(tx *gorm.DB) error {
    if u.SSNEncrypted != "" {
        decrypted, err := Decrypt(u.SSNEncrypted)
        if err != nil {
            return err
        }
        u.SSN = decrypted
    }
    return nil
}
```

---

## HTTPS

### Enable TLS

```go
func main() {
    app := fiber.New()
    
    // Routes
    app.Get("/", handler)
    
    // Start with TLS
    log.Fatal(app.ListenTLS(":443",
        "/path/to/cert.pem",
        "/path/to/key.pem",
    ))
}
```

### Force HTTPS Redirect

```go
func RedirectToHTTPS() fiber.Handler {
    return func(c *fiber.Ctx) error {
        if c.Protocol() != "https" {
            return c.Redirect("https://" + c.Hostname() + c.OriginalURL())
        }
        return c.Next()
    }
}

app.Use(RedirectToHTTPS())
```

---

## CORS

### Configure CORS

```go
import "github.com/gofiber/fiber/v2/middleware/cors"

app.Use(cors.New(cors.Config{
    AllowOrigins: "https://example.com,https://app.example.com",
    AllowMethods: "GET,POST,PUT,DELETE",
    AllowHeaders: "Origin,Content-Type,Accept,Authorization",
    AllowCredentials: true,
    MaxAge: 3600,
}))

// Development: Allow all origins
if os.Getenv("APP_ENV") == "development" {
    app.Use(cors.New(cors.Config{
        AllowOrigins: "*",
    }))
}
```

---

## Rate Limiting

### Basic Rate Limiter

```go
import "github.com/gofiber/fiber/v2/middleware/limiter"

app.Use(limiter.New(limiter.Config{
    Max:        20,
    Expiration: 30 * time.Second,
    KeyGenerator: func(c *fiber.Ctx) string {
        return c.IP()
    },
    LimitReached: func(c *fiber.Ctx) error {
        return c.Status(429).JSON(fiber.Map{
            "error": "Too many requests",
        })
    },
}))
```

### Per-Route Rate Limiting

```go
// Login endpoint: 5 attempts per minute
loginLimiter := limiter.New(limiter.Config{
    Max:        5,
    Expiration: 1 * time.Minute,
})

app.Post("/login", loginLimiter, loginHandler)

// API endpoint: 100 requests per minute
apiLimiter := limiter.New(limiter.Config{
    Max:        100,
    Expiration: 1 * time.Minute,
})

api := app.Group("/api", apiLimiter)
```

---

## Security Headers

### Set Security Headers

```go
func SecurityHeaders() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Prevent clickjacking
        c.Set("X-Frame-Options", "DENY")
        
        // XSS Protection
        c.Set("X-XSS-Protection", "1; mode=block")
        
        // Prevent MIME sniffing
        c.Set("X-Content-Type-Options", "nosniff")
        
        // Referrer Policy
        c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
        
        // Content Security Policy
        c.Set("Content-Security-Policy", "default-src 'self'")
        
        // HSTS (HTTPS only)
        if c.Protocol() == "https" {
            c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        }
        
        return c.Next()
    }
}

app.Use(SecurityHeaders())
```

---

## Secure File Uploads

```go
func UploadFile(c *fiber.Ctx) error {
    file, err := c.FormFile("file")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "No file uploaded",
        })
    }
    
    // Validate file size (10MB max)
    if file.Size > 10*1024*1024 {
        return c.Status(400).JSON(fiber.Map{
            "error": "File too large (max 10MB)",
        })
    }
    
    // Validate file type
    allowedTypes := map[string]bool{
        "image/jpeg": true,
        "image/png":  true,
        "image/gif":  true,
    }
    
    contentType := file.Header.Get("Content-Type")
    if !allowedTypes[contentType] {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid file type",
        })
    }
    
    // Generate safe filename
    ext := filepath.Ext(file.Filename)
    filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
    
    // Save file
    uploadPath := filepath.Join("uploads", filename)
    if err := c.SaveFile(file, uploadPath); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to save file",
        })
    }
    
    return c.JSON(fiber.Map{
        "filename": filename,
    })
}
```

---

## Environment Variables

### Secure Configuration

```go
// ❌ Bad: Hardcoded secrets
const apiKey = "sk_live_12345"
const dbPassword = "admin123"

// ✅ Good: Environment variables
apiKey := os.Getenv("API_KEY")
dbPassword := os.Getenv("DB_PASSWORD")

if apiKey == "" {
    log.Fatal("API_KEY is required")
}
```

### Use .env.example

```bash
# .env.example (commit this)
API_KEY=your_api_key_here
DB_PASSWORD=your_password_here

# .env (don't commit)
API_KEY=sk_live_actual_key
DB_PASSWORD=actual_password
```

---

## Security Checklist

- [ ] Passwords hashed with bcrypt
- [ ] JWT tokens for authentication
- [ ] Role-based authorization
- [ ] Input validation on all endpoints
- [ ] SQL injection prevention (parameterized queries)
- [ ] XSS prevention (escape output)
- [ ] CSRF protection (tokens)
- [ ] HTTPS enabled
- [ ] Security headers set
- [ ] CORS configured properly
- [ ] Rate limiting enabled
- [ ] File upload validation
- [ ] Secrets in environment variables
- [ ] Error messages don't leak info
- [ ] Dependencies updated regularly

---

## Common Vulnerabilities

### SQL Injection

```go
// ❌ Vulnerable
email := c.Query("email")
db.Raw("SELECT * FROM users WHERE email = '" + email + "'").Scan(&user)

// ✅ Safe
db.Where("email = ?", email).First(&user)
```

### XSS (Cross-Site Scripting)

```go
// ❌ Vulnerable
userInput := c.Query("name")
return c.SendString("<h1>Hello " + userInput + "</h1>")

// ✅ Safe
import "html"
userInput := html.EscapeString(c.Query("name"))
return c.SendString("<h1>Hello " + userInput + "</h1>")
```

### CSRF (Cross-Site Request Forgery)

```go
import "github.com/gofiber/fiber/v2/middleware/csrf"

app.Use(csrf.New(csrf.Config{
    KeyLookup:      "header:X-CSRF-Token",
    CookieName:     "csrf_",
    CookieSameSite: "Strict",
    Expiration:     1 * time.Hour,
}))
```

---

## Next Steps

- [**Error Handling**](error-handling.md) - Secure error messages
- [**Middleware**](middleware.md) - Security middleware
- [**Deployment**](../deployment/production-setup.md) - Production security
- [**Monitoring**](../deployment/monitoring.md) - Security monitoring

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
