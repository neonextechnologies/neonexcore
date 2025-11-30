# Testing

Comprehensive guide to testing Neonex Core applications.

---

## Overview

Neonex Core supports multiple testing strategies:

- **Unit Tests** - Test individual functions and methods
- **Integration Tests** - Test component interactions
- **Repository Tests** - Test database operations
- **Service Tests** - Test business logic
- **Controller Tests** - Test HTTP handlers
- **End-to-End Tests** - Test complete user flows

**Testing Framework:** Go's built-in `testing` package + [Testify](https://github.com/stretchr/testify)

---

## Quick Start

### Install Testing Dependencies

```bash
# Testify for assertions and mocks
go get github.com/stretchr/testify

# SQLite for test database
go get gorm.io/driver/sqlite

# HTTP testing
go get github.com/gofiber/fiber/v2
```

### Run Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package
go test ./modules/user

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Unit Testing

### Testing Functions

```go
// modules/user/service_test.go
package user

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
    tests := []struct{
        name string
        email string
        want bool
    }{
        {"valid email", "user@example.com", true},
        {"invalid email", "invalid", false},
        {"empty email", "", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := ValidateEmail(tt.email)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

### Table-Driven Tests

```go
func TestCalculateDiscount(t *testing.T) {
    testCases := []struct {
        name     string
        price    float64
        discount float64
        expected float64
    }{
        {"10% discount", 100.0, 0.10, 90.0},
        {"50% discount", 200.0, 0.50, 100.0},
        {"no discount", 150.0, 0.00, 150.0},
        {"100% discount", 80.0, 1.00, 0.0},
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := CalculateDiscount(tc.price, tc.discount)
            assert.Equal(t, tc.expected, result)
        })
    }
}
```

---

## Repository Testing

### Setup Test Database

```go
// modules/user/repository_test.go
package user

import (
    "testing"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
    suite.Suite
    db   *gorm.DB
    repo Repository
}

func (suite *RepositoryTestSuite) SetupTest() {
    // Create in-memory database
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    assert.NoError(suite.T(), err)
    
    // Migrate schema
    err = db.AutoMigrate(&User{})
    assert.NoError(suite.T(), err)
    
    suite.db = db
    suite.repo = NewRepository(db)
}

func (suite *RepositoryTestSuite) TearDownTest() {
    // Clean up after each test
    sqlDB, _ := suite.db.DB()
    sqlDB.Close()
}

func TestRepositoryTestSuite(t *testing.T) {
    suite.Run(t, new(RepositoryTestSuite))
}
```

### Test CRUD Operations

```go
func (suite *RepositoryTestSuite) TestCreate() {
    user := &User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    err := suite.repo.Create(user)
    
    assert.NoError(suite.T(), err)
    assert.NotZero(suite.T(), user.ID)
}

func (suite *RepositoryTestSuite) TestFindByID() {
    // Arrange
    user := &User{Name: "Jane", Email: "jane@example.com"}
    suite.repo.Create(user)
    
    // Act
    found, err := suite.repo.FindByID(user.ID)
    
    // Assert
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), user.Email, found.Email)
}

func (suite *RepositoryTestSuite) TestUpdate() {
    user := &User{Name: "Old Name", Email: "old@example.com"}
    suite.repo.Create(user)
    
    user.Name = "New Name"
    err := suite.repo.Update(user)
    
    assert.NoError(suite.T(), err)
    
    updated, _ := suite.repo.FindByID(user.ID)
    assert.Equal(suite.T(), "New Name", updated.Name)
}

func (suite *RepositoryTestSuite) TestDelete() {
    user := &User{Name: "To Delete", Email: "delete@example.com"}
    suite.repo.Create(user)
    
    err := suite.repo.Delete(user.ID)
    assert.NoError(suite.T(), err)
    
    _, err = suite.repo.FindByID(user.ID)
    assert.Error(suite.T(), err)
}
```

---

## Service Testing

### Mocking Dependencies

```go
// modules/user/mocks/repository.go
package mocks

import (
    "github.com/stretchr/testify/mock"
    "myapp/modules/user"
)

type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Create(user *user.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockRepository) FindByID(id uint) (*user.User, error) {
    args := m.Called(id)
    return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockRepository) Update(user *user.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockRepository) Delete(id uint) error {
    args := m.Called(id)
    return args.Error(0)
}
```

### Test Service Logic

```go
// modules/user/service_test.go
package user

import (
    "errors"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "myapp/modules/user/mocks"
)

func TestCreateUser(t *testing.T) {
    mockRepo := new(mocks.MockRepository)
    service := NewService(mockRepo)
    
    user := &User{Name: "John", Email: "john@example.com"}
    
    mockRepo.On("Create", user).Return(nil)
    
    err := service.CreateUser(user)
    
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}

func TestCreateUser_ValidationError(t *testing.T) {
    mockRepo := new(mocks.MockRepository)
    service := NewService(mockRepo)
    
    user := &User{Name: "", Email: "john@example.com"}
    
    err := service.CreateUser(user)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "name is required")
}

func TestGetUser_NotFound(t *testing.T) {
    mockRepo := new(mocks.MockRepository)
    service := NewService(mockRepo)
    
    mockRepo.On("FindByID", uint(999)).
        Return((*User)(nil), errors.New("not found"))
    
    user, err := service.GetUser(999)
    
    assert.Error(t, err)
    assert.Nil(t, user)
    mockRepo.AssertExpectations(t)
}
```

---

## Controller Testing

### HTTP Request Testing

```go
// modules/user/controller_test.go
package user

import (
    "bytes"
    "encoding/json"
    "net/http/httptest"
    "testing"
    
    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "myapp/modules/user/mocks"
)

func TestGetAllUsers(t *testing.T) {
    // Setup
    app := fiber.New()
    mockService := new(mocks.MockService)
    controller := NewController(mockService)
    
    app.Get("/users", controller.GetAll)
    
    // Mock data
    users := []User{
        {ID: 1, Name: "John", Email: "john@example.com"},
        {ID: 2, Name: "Jane", Email: "jane@example.com"},
    }
    
    mockService.On("GetAllUsers").Return(users, nil)
    
    // Test
    req := httptest.NewRequest("GET", "/users", nil)
    resp, _ := app.Test(req)
    
    // Assert
    assert.Equal(t, 200, resp.StatusCode)
    mockService.AssertExpectations(t)
}

func TestCreateUser(t *testing.T) {
    app := fiber.New()
    mockService := new(mocks.MockService)
    controller := NewController(mockService)
    
    app.Post("/users", controller.Create)
    
    user := User{Name: "New User", Email: "new@example.com"}
    mockService.On("CreateUser", mock.AnythingOfType("*user.User")).
        Return(nil)
    
    body, _ := json.Marshal(user)
    req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    
    resp, _ := app.Test(req)
    
    assert.Equal(t, 201, resp.StatusCode)
    mockService.AssertExpectations(t)
}
```

---

## Integration Testing

### Test Complete Flow

```go
// tests/integration/user_test.go
package integration

import (
    "testing"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
    "myapp/modules/user"
)

func TestUserIntegration(t *testing.T) {
    // Setup database
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&user.User{})
    
    // Setup application
    app := fiber.New()
    repo := user.NewRepository(db)
    service := user.NewService(repo)
    controller := user.NewController(service)
    
    user.RegisterRoutes(app, controller)
    
    // Test create
    body := `{"name":"Integration Test","email":"test@example.com"}`
    req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    resp, _ := app.Test(req)
    
    assert.Equal(t, 201, resp.StatusCode)
    
    // Test get
    req = httptest.NewRequest("GET", "/users/1", nil)
    resp, _ = app.Test(req)
    
    assert.Equal(t, 200, resp.StatusCode)
}
```

---

## Test Coverage

### Generate Coverage Report

```bash
# Run tests with coverage
go test -coverprofile=coverage.out ./...

# View coverage percentage
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# Open in browser
open coverage.html  # Mac
xdg-open coverage.html  # Linux
start coverage.html  # Windows
```

### Coverage Goals

**Target Coverage:**
- **Repository Layer:** 90%+ (data access is critical)
- **Service Layer:** 85%+ (business logic)
- **Controller Layer:** 75%+ (HTTP handlers)
- **Overall Project:** 80%+

### Example Output

```
myapp/modules/user/repository.go:15:    Create          100.0%
myapp/modules/user/repository.go:20:    FindByID        100.0%
myapp/modules/user/repository.go:25:    FindAll         100.0%
myapp/modules/user/repository.go:30:    Update          100.0%
myapp/modules/user/repository.go:35:    Delete          100.0%
myapp/modules/user/service.go:15:       CreateUser      85.7%
myapp/modules/user/service.go:25:       GetUser         100.0%
myapp/modules/user/controller.go:15:    GetAll          75.0%
myapp/modules/user/controller.go:25:    GetByID         75.0%
total:                                  (statements)    87.5%
```

---

## Best Practices

### ✅ DO: Test Behavior, Not Implementation

```go
// Good: Test behavior
func TestUserCreation(t *testing.T) {
    service := NewService(repo)
    user := &User{Name: "John", Email: "john@example.com"}
    
    err := service.CreateUser(user)
    
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)  // Verify outcome
}

// Bad: Test implementation details
func TestUserCreation(t *testing.T) {
    service := NewService(repo)
    // Don't test how it's implemented internally
}
```

### ✅ DO: Use Test Fixtures

```go
// test/fixtures/users.go
package fixtures

func CreateTestUser() *user.User {
    return &user.User{
        Name:  "Test User",
        Email: "test@example.com",
    }
}

func CreateTestUsers(n int) []*user.User {
    users := make([]*user.User, n)
    for i := 0; i < n; i++ {
        users[i] = &user.User{
            Name:  fmt.Sprintf("User %d", i),
            Email: fmt.Sprintf("user%d@example.com", i),
        }
    }
    return users
}
```

### ✅ DO: Isolate Tests

```go
func TestIndependentTest1(t *testing.T) {
    db := setupTestDB()
    defer db.Close()
    // Test in isolation
}

func TestIndependentTest2(t *testing.T) {
    db := setupTestDB()
    defer db.Close()
    // Test in isolation
}
```

### ❌ DON'T: Share State Between Tests

```go
// Bad: Global state
var testDB *gorm.DB

func TestA(t *testing.T) {
    testDB = setupDB()  // Affects other tests
}

// Good: Local state
func TestB(t *testing.T) {
    db := setupTestDB()
    defer db.Close()
}
```

---

## Test Organization

### Directory Structure

```
modules/
└── user/
    ├── controller.go
    ├── controller_test.go
    ├── service.go
    ├── service_test.go
    ├── repository.go
    ├── repository_test.go
    ├── mocks/
    │   ├── repository.go
    │   └── service.go
    └── testdata/
        └── fixtures.json

tests/
├── integration/
│   └── user_integration_test.go
└── e2e/
    └── user_flow_test.go
```

### Naming Conventions

```go
// Test function names
func TestCreateUser(t *testing.T) {}           // Basic test
func TestCreateUser_EmptyName(t *testing.T) {}  // Specific scenario
func TestCreateUser_DuplicateEmail(t *testing.T) {}

// Test suite names
type UserServiceTestSuite struct {}
type UserRepositoryTestSuite struct {}
```

---

## Continuous Integration

### GitHub Actions Workflow

```yaml
# .github/workflows/test.yml
name: Tests

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Install dependencies
        run: go mod download
      
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
```

---

## Next Steps

- [**Debugging**](debugging.md) - Debug failing tests
- [**Best Practices**](best-practices.md) - Code quality guidelines
- [**Hot Reload**](hot-reload.md) - Test during development
- [**CI/CD**](../deployment/production-setup.md) - Automate testing

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
