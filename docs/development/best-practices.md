# Development Best Practices

Guidelines for writing clean, maintainable, and efficient Neonex Core applications.

---

## Code Organization

### Module Structure

✅ **DO: Follow Standard Structure**

```
modules/user/
├── model.go          # Data models
├── repository.go     # Data access
├── service.go        # Business logic
├── controller.go     # HTTP handlers
├── routes.go         # Route definitions
├── di.go            # Dependency injection
├── seeder.go        # Test data
├── user.go          # Module entry point
└── module.json      # Metadata
```

❌ **DON'T: Mix Concerns**

```
modules/user/
├── everything.go     # ❌ All code in one file
└── utils.go          # ❌ Vague naming
```

### Package Organization

```go
// ✅ Good: Clear package purpose
package user

import (
    "github.com/myapp/pkg/database"
    "github.com/myapp/pkg/logger"
)

// ❌ Bad: Importing between modules
package user

import (
    "github.com/myapp/modules/product"  // ❌ Tight coupling
)
```

---

## Naming Conventions

### Variables

```go
// ✅ Good: Clear and descriptive
userID := 123
firstName := "John"
isActive := true
userRepository := NewUserRepository(db)

// ❌ Bad: Unclear abbreviations
uid := 123
fn := "John"
flag := true
ur := NewUserRepository(db)
```

### Functions

```go
// ✅ Good: Verb + noun
func CreateUser(user *User) error
func GetUserByID(id uint) (*User, error)
func ValidateEmail(email string) bool
func CalculateDiscount(price float64) float64

// ❌ Bad: Unclear purpose
func DoStuff(data interface{}) error
func Process(x int) int
func Handle(c *fiber.Ctx) error
```

### Constants

```go
// ✅ Good: Clear context
const (
    MaxLoginAttempts     = 3
    SessionTimeoutMinutes = 30
    DefaultPageSize      = 10
)

// ❌ Bad: Magic numbers
const (
    Three  = 3
    Thirty = 30
    Ten    = 10
)
```

---

## Error Handling

### Always Check Errors

```go
// ✅ Good
user, err := service.GetUser(id)
if err != nil {
    logger.Error("Failed to get user", zap.Error(err))
    return c.Status(500).JSON(fiber.Map{
        "error": "Internal server error",
    })
}

// ❌ Bad
user, _ := service.GetUser(id)  // Ignoring errors
```

### Wrap Errors with Context

```go
// ✅ Good
if err := repo.Create(user); err != nil {
    return fmt.Errorf("failed to create user %s: %w", user.Email, err)
}

// ❌ Bad
if err := repo.Create(user); err != nil {
    return err  // No context
}
```

### Custom Error Types

```go
// ✅ Good: Custom errors
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Usage
if user.Email == "" {
    return &ValidationError{
        Field:   "email",
        Message: "email is required",
    }
}
```

---

## Database Best Practices

### Use Transactions for Multiple Operations

```go
// ✅ Good
err := database.WithTransaction(ctx, db, func(tx *gorm.DB) error {
    if err := tx.Create(&user).Error; err != nil {
        return err  // Auto-rollback
    }
    
    if err := tx.Create(&profile).Error; err != nil {
        return err  // Auto-rollback
    }
    
    return nil  // Auto-commit
})

// ❌ Bad: No transaction
db.Create(&user)
db.Create(&profile)  // If this fails, user still created
```

### Avoid N+1 Queries

```go
// ✅ Good: Preload relationships
var users []User
db.Preload("Profile").Preload("Orders").Find(&users)

// ❌ Bad: N+1 problem
var users []User
db.Find(&users)
for _, user := range users {
    db.Find(&user.Profile)   // N queries
    db.Find(&user.Orders)    // N queries
}
```

### Use Indexes

```go
// ✅ Good
type User struct {
    ID    uint   `gorm:"primarykey"`
    Email string `gorm:"uniqueIndex;size:100"`  // Index for lookups
    Name  string `gorm:"index"`                 // Index for searches
}

// ❌ Bad: No indexes on frequently queried fields
type User struct {
    ID    uint
    Email string  // No index, slow lookups
    Name  string  // No index, slow searches
}
```

---

## API Design

### RESTful Conventions

```go
// ✅ Good: RESTful routes
GET    /users         // List all
GET    /users/:id     // Get one
POST   /users         // Create
PUT    /users/:id     // Update
DELETE /users/:id     // Delete

// ❌ Bad: Non-standard
GET    /getUsers
POST   /createUser
POST   /deleteUser/:id
```

### Consistent Response Format

```go
// ✅ Good: Consistent structure
type Response struct {
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
    Message string      `json:"message,omitempty"`
}

func (c *Controller) GetAll(ctx *fiber.Ctx) error {
    users, err := c.service.GetAllUsers()
    if err != nil {
        return ctx.Status(500).JSON(Response{
            Error: err.Error(),
        })
    }
    
    return ctx.JSON(Response{
        Data: users,
    })
}
```

### Use Proper HTTP Status Codes

```go
// ✅ Good
return ctx.Status(200).JSON(user)    // OK
return ctx.Status(201).JSON(user)    // Created
return ctx.Status(204).Send(nil)     // No Content
return ctx.Status(400).JSON(err)     // Bad Request
return ctx.Status(401).JSON(err)     // Unauthorized
return ctx.Status(404).JSON(err)     // Not Found
return ctx.Status(500).JSON(err)     // Internal Error

// ❌ Bad
return ctx.Status(200).JSON(err)     // Error with 200
return ctx.Status(500).JSON(user)    // Success with 500
```

---

## Dependency Injection

### Register as Interfaces

```go
// ✅ Good: Interface-based
container.Singleton(func() UserRepository {
    return NewUserRepository(db)
})

container.Singleton(func() UserService {
    repo := container.Resolve((*UserRepository)(nil)).(UserRepository)
    return NewUserService(repo)
})

// ❌ Bad: Concrete types
container.Singleton(func() *userRepository {
    return NewUserRepository(db)
})
```

### Avoid Service Locator Pattern

```go
// ✅ Good: Constructor injection
type Service struct {
    repo   Repository
    logger logger.Logger
}

func NewService(repo Repository, log logger.Logger) Service {
    return &Service{repo: repo, logger: log}
}

// ❌ Bad: Service locator
type Service struct {
    container *Container
}

func (s *Service) GetRepo() Repository {
    return s.container.Resolve("repository").(Repository)
}
```

---

## Logging Best Practices

### Structured Logging

```go
// ✅ Good: Structured with context
logger.Info("User created",
    zap.String("user_id", user.ID),
    zap.String("email", user.Email),
    zap.Time("created_at", time.Now()),
)

// ❌ Bad: String concatenation
logger.Info("User created: " + user.Email)
```

### Log Levels

```go
// ✅ Good: Appropriate levels
logger.Debug("Cache miss for key", zap.String("key", key))
logger.Info("Server started", zap.Int("port", 8080))
logger.Warn("Rate limit exceeded", zap.String("ip", ip))
logger.Error("Database connection failed", zap.Error(err))
logger.Fatal("Configuration missing", zap.String("key", "DB_HOST"))

// ❌ Bad: Everything as Info
logger.Info("Cache miss")
logger.Info("Rate limit exceeded")
logger.Info("Database error: " + err.Error())
```

---

## Testing Best Practices

### Test Coverage

```go
// ✅ Good: Comprehensive tests
func TestCreateUser(t *testing.T) {
    t.Run("success", func(t *testing.T) {
        // Test happy path
    })
    
    t.Run("empty name", func(t *testing.T) {
        // Test validation
    })
    
    t.Run("duplicate email", func(t *testing.T) {
        // Test constraints
    })
}

// ❌ Bad: Only happy path
func TestCreateUser(t *testing.T) {
    // Only test success case
}
```

### Use Table-Driven Tests

```go
// ✅ Good
tests := []struct {
    name     string
    input    User
    wantErr  bool
    errMsg   string
}{
    {"valid user", User{Name: "John", Email: "john@example.com"}, false, ""},
    {"empty name", User{Name: "", Email: "john@example.com"}, true, "name required"},
    {"invalid email", User{Name: "John", Email: "invalid"}, true, "invalid email"},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        err := ValidateUser(tt.input)
        if tt.wantErr {
            assert.Error(t, err)
            assert.Contains(t, err.Error(), tt.errMsg)
        } else {
            assert.NoError(t, err)
        }
    })
}
```

---

## Security Best Practices

### Never Log Sensitive Data

```go
// ✅ Good
logger.Info("User login",
    zap.String("user_id", user.ID),
    zap.String("email", user.Email),
)

// ❌ Bad
logger.Info("User login",
    zap.String("password", user.Password),  // ❌ Never log passwords
    zap.String("api_key", user.APIKey),     // ❌ Never log secrets
)
```

### Validate Input

```go
// ✅ Good
func CreateUser(ctx *fiber.Ctx) error {
    var user User
    if err := ctx.BodyParser(&user); err != nil {
        return ctx.Status(400).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    
    if err := ValidateUser(&user); err != nil {
        return ctx.Status(400).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    // Safe to use user
}

// ❌ Bad: No validation
func CreateUser(ctx *fiber.Ctx) error {
    var user User
    ctx.BodyParser(&user)
    db.Create(&user)  // Unsafe
}
```

### Use Parameterized Queries

```go
// ✅ Good: GORM handles parameterization
db.Where("email = ?", email).First(&user)

// ❌ Bad: String concatenation (SQL injection)
query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)
db.Raw(query).Scan(&user)
```

---

## Performance Best Practices

### Use Pagination

```go
// ✅ Good
func GetUsers(page, pageSize int) ([]User, error) {
    var users []User
    offset := (page - 1) * pageSize
    
    err := db.Limit(pageSize).Offset(offset).Find(&users).Error
    return users, err
}

// ❌ Bad: Load all records
func GetUsers() ([]User, error) {
    var users []User
    err := db.Find(&users).Error  // Could be millions
    return users, err
}
```

### Use Connection Pooling

```go
// ✅ Good
sqlDB, err := db.DB()
sqlDB.SetMaxOpenConns(25)
sqlDB.SetMaxIdleConns(5)
sqlDB.SetConnMaxLifetime(5 * time.Minute)

// ❌ Bad: Default settings
// No connection pool configuration
```

### Cache Expensive Operations

```go
// ✅ Good
var cache = make(map[string]interface{})
var mu sync.RWMutex

func GetExpensiveData(key string) (interface{}, error) {
    mu.RLock()
    if data, ok := cache[key]; ok {
        mu.RUnlock()
        return data, nil
    }
    mu.RUnlock()
    
    // Compute expensive result
    data := computeExpensiveResult(key)
    
    mu.Lock()
    cache[key] = data
    mu.Unlock()
    
    return data, nil
}
```

---

## Documentation

### Document Public APIs

```go
// ✅ Good
// CreateUser creates a new user in the database.
// It validates the user input and returns an error if validation fails.
//
// Parameters:
//   - user: User struct containing name and email (required)
//
// Returns:
//   - error: Validation or database error, nil if successful
func CreateUser(user *User) error {
    // Implementation
}

// ❌ Bad: No documentation
func CreateUser(user *User) error {
    // Implementation
}
```

### Keep README Updated

```markdown
# User Module

## Features
- User CRUD operations
- Email validation
- Password hashing

## API Endpoints
- `POST /users` - Create user
- `GET /users/:id` - Get user
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user

## Usage
\`\`\`go
service := user.NewService(repo)
err := service.CreateUser(&user)
\`\`\`
```

---

## Git Practices

### Commit Messages

```bash
# ✅ Good: Clear and descriptive
feat: add user authentication
fix: resolve nil pointer in user service
docs: update API documentation
refactor: extract validation to separate function
test: add tests for user repository

# ❌ Bad: Unclear
update stuff
fix
changes
wip
```

### Branch Naming

```bash
# ✅ Good
feature/user-authentication
bugfix/nil-pointer-user-service
hotfix/security-vulnerability
release/v1.2.0

# ❌ Bad
new-branch
fix
updates
test123
```

---

## Next Steps

- [**Testing**](testing.md) - Write quality tests
- [**Debugging**](debugging.md) - Debug effectively
- [**Security**](../advanced/security.md) - Secure your app
- [**Performance**](../advanced/performance.md) - Optimize performance

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
