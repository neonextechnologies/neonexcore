# Testing Guide

Complete guide to testing Neonex Core applications.

---

## Overview

Neonex Core supports comprehensive testing strategies:
- **Unit Tests** - Test individual components
- **Integration Tests** - Test component interactions  
- **Repository Tests** - Test data layer
- **Service Tests** - Test business logic
- **HTTP Tests** - Test API endpoints

---

## Quick Start

### Run Tests

```powershell
# All tests
go test ./...

# Specific module
go test ./modules/user/...

# With coverage
go test -cover ./...

# Verbose
go test -v ./...
```

---

## Unit Testing

### Service Tests

```go
// modules/user/service_test.go
package user

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestService_CreateUser(t *testing.T) {
    // Setup
    mockRepo := &MockRepository{}
    mockLogger := &MockLogger{}
    service := NewService(mockRepo, mockLogger)
    
    ctx := context.Background()
    req := &CreateUserRequest{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    // Execute
    user, err := service.CreateUser(ctx, req)
    
    // Assert
    require.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "John Doe", user.Name)
    assert.Equal(t, "john@example.com", user.Email)
    
    // Verify mock calls
    assert.Equal(t, 1, mockRepo.CreateCallCount)
}
```

### Table-Driven Tests

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "user@example.com", false},
        {"invalid format", "invalid", true},
        {"empty email", "", true},
        {"no domain", "user@", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

---

## Mock Objects

### Mock Repository

```go
// test/mocks/repository.go
package mocks

type MockRepository struct {
    Users          []*User
    CreateCallCount int
}

func (m *MockRepository) Create(ctx context.Context, user *User) error {
    m.CreateCallCount++
    user.ID = uint(len(m.Users) + 1)
    m.Users = append(m.Users, user)
    return nil
}

func (m *MockRepository) FindByID(ctx context.Context, id uint) (*User, error) {
    for _, u := range m.Users {
        if u.ID == id {
            return u, nil
        }
    }
    return nil, gorm.ErrRecordNotFound
}
```

### Mock Logger

```go
// test/mocks/logger.go
package mocks

type MockLogger struct {
    InfoCalls  []LogCall
    ErrorCalls []LogCall
}

type LogCall struct {
    Message string
    Fields  map[string]interface{}
}

func (m *MockLogger) Info(msg string, fields ...logger.Fields) {
    m.InfoCalls = append(m.InfoCalls, LogCall{
        Message: msg,
        Fields:  fields[0],
    })
}
```

---

## Integration Tests

### Database Integration Tests

```go
// modules/user/integration_test.go
//go:build integration

package user

import (
    "context"
    "testing"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)
    
    err = db.AutoMigrate(&User{})
    require.NoError(t, err)
    
    return db
}

func TestRepository_Integration(t *testing.T) {
    db := setupTestDB(t)
    repo := NewRepository(db)
    
    ctx := context.Background()
    
    // Create
    user := &User{Name: "John", Email: "john@example.com"}
    err := repo.Create(ctx, user)
    require.NoError(t, err)
    assert.NotZero(t, user.ID)
    
    // Find
    found, err := repo.FindByID(ctx, user.ID)
    require.NoError(t, err)
    assert.Equal(t, user.Email, found.Email)
}
```

**Run integration tests:**
```powershell
go test -tags=integration ./...
```

---

## HTTP Testing

### Controller Tests

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
)

func TestController_Create(t *testing.T) {
    // Setup
    app := fiber.New()
    mockService := &MockService{}
    mockLogger := &MockLogger{}
    ctrl := NewController(mockService, mockLogger)
    
    app.Post("/users", ctrl.Create)
    
    // Request
    req := CreateUserRequest{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    body, _ := json.Marshal(req)
    
    httpReq := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
    httpReq.Header.Set("Content-Type", "application/json")
    
    // Execute
    resp, err := app.Test(httpReq)
    
    // Assert
    require.NoError(t, err)
    assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}
```

---

## Test Helpers

### Common Setup

```go
// test/helpers/setup.go
package helpers

func SetupTestApp(t *testing.T) *fiber.App {
    app := fiber.New(fiber.Config{
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            return c.Status(500).JSON(fiber.Map{
                "error": err.Error(),
            })
        },
    })
    return app
}

func SetupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)
    return db
}
```

---

## Test Coverage

### Generate Coverage Report

```powershell
# Generate coverage
go test -coverprofile=coverage.out ./...

# View coverage
go tool cover -html=coverage.out

# Coverage per package
go test -cover ./...
```

### Coverage Thresholds

```powershell
# Fail if coverage < 80%
go test -cover ./... | grep -E "coverage: [0-7][0-9]\."
```

---

## Benchmarking

### Benchmark Tests

```go
func BenchmarkCreateUser(b *testing.B) {
    service := NewService(repo, logger)
    ctx := context.Background()
    req := &CreateUserRequest{
        Name:  "Test",
        Email: "test@example.com",
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        service.CreateUser(ctx, req)
    }
}
```

**Run benchmarks:**
```powershell
go test -bench=. ./...
```

---

## Best Practices

### ✅ DO:
- Use table-driven tests
- Test error cases
- Mock external dependencies
- Use meaningful test names
- Clean up test data

### ❌ DON'T:
- Test implementation details
- Share state between tests
- Ignore error handling
- Skip edge cases

---

## Next Steps

- [**Hot Reload**](hot-reload.md) - Development workflow
- [**Debugging**](debugging.md) - Debug techniques
- [**Best Practices**](best-practices.md) - Development tips

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
