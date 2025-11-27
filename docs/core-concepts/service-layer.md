# Service Layer

Master the Service Layer pattern for organizing business logic in Neonex Core applications.

---

## Overview

The **Service Layer** contains your application's business logic, sitting between Controllers (presentation) and Repositories (data access). It orchestrates operations, enforces business rules, and coordinates between multiple repositories.

**Responsibilities:**
- ✅ Business logic and validation
- ✅ Transaction management
- ✅ Coordination between repositories
- ✅ Data transformation
- ✅ Error handling

**NOT Responsible For:**
- ❌ HTTP request handling (Controller's job)
- ❌ Database queries (Repository's job)
- ❌ Data modeling (Model's job)

---

## Architecture

### Layered Flow

```
Controller → Service → Repository → Database
   ↓          ↓          ↓
 HTTP      Business    Data
Request     Logic     Access
```

**Example Flow:**
1. Controller receives HTTP request
2. Controller calls Service method
3. Service validates and processes
4. Service calls Repository
5. Repository queries database
6. Results flow back up

---

## Creating a Service

### Service Interface

```go
// modules/user/service.go
package user

import "context"

type Service interface {
    CreateUser(ctx context.Context, user *User) error
    GetUser(ctx context.Context, id uint) (*User, error)
    GetAllUsers(ctx context.Context) ([]*User, error)
    UpdateUser(ctx context.Context, id uint, user *User) error
    DeleteUser(ctx context.Context, id uint) error
    
    // Business operations
    ActivateUser(ctx context.Context, id uint) error
    DeactivateUser(ctx context.Context, id uint) error
    ChangePassword(ctx context.Context, id uint, oldPass, newPass string) error
}
```

### Service Implementation

```go
type service struct {
    repo   Repository
    logger logger.Logger
}

func NewService(repo Repository, log logger.Logger) Service {
    return &service{
        repo:   repo,
        logger: log,
    }
}

func (s *service) CreateUser(ctx context.Context, user *User) error {
    // 1. Validation
    if user.Email == "" {
        return errors.New("email is required")
    }
    
    // 2. Business logic
    existing, _ := s.repo.FindByEmail(ctx, user.Email)
    if existing != nil {
        return errors.New("email already exists")
    }
    
    // 3. Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(user.Password), 
        bcrypt.DefaultCost,
    )
    if err != nil {
        return err
    }
    user.Password = string(hashedPassword)
    
    // 4. Set defaults
    user.Active = true
    user.Role = "user"
    
    // 5. Call repository
    if err := s.repo.Create(ctx, user); err != nil {
        s.logger.Error("Failed to create user", logger.Fields{
            "email": user.Email,
            "error": err.Error(),
        })
        return err
    }
    
    // 6. Log success
    s.logger.Info("User created", logger.Fields{
        "id":    user.ID,
        "email": user.Email,
    })
    
    return nil
}
```

---

## Common Patterns

### Pattern 1: Validation

```go
func (s *service) CreateProduct(ctx context.Context, product *Product) error {
    // Input validation
    if err := s.validateProduct(product); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    // Business rules
    if product.Price < 0 {
        return errors.New("price cannot be negative")
    }
    
    if product.Stock < 0 {
        return errors.New("stock cannot be negative")
    }
    
    return s.repo.Create(ctx, product)
}

func (s *service) validateProduct(product *Product) error {
    if product.Name == "" {
        return errors.New("name is required")
    }
    if len(product.Name) < 3 {
        return errors.New("name must be at least 3 characters")
    }
    return nil
}
```

### Pattern 2: Data Transformation

```go
func (s *service) GetUserProfile(ctx context.Context, id uint) (*UserProfile, error) {
    // Get data
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Transform to DTO
    profile := &UserProfile{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        JoinedAt:  user.CreatedAt,
        IsActive:  user.Active,
    }
    
    // Don't expose password!
    return profile, nil
}
```

### Pattern 3: Multiple Repositories

```go
type OrderService struct {
    orderRepo    OrderRepository
    productRepo  ProductRepository
    userRepo     UserRepository
    paymentRepo  PaymentRepository
    logger       logger.Logger
}

func (s *OrderService) CreateOrder(ctx context.Context, order *Order) error {
    // 1. Validate user
    user, err := s.userRepo.FindByID(ctx, order.UserID)
    if err != nil {
        return fmt.Errorf("user not found: %w", err)
    }
    
    // 2. Validate products
    for _, item := range order.Items {
        product, err := s.productRepo.FindByID(ctx, item.ProductID)
        if err != nil {
            return fmt.Errorf("product %d not found", item.ProductID)
        }
        
        // Check stock
        if product.Stock < item.Quantity {
            return fmt.Errorf("insufficient stock for product %s", product.Name)
        }
    }
    
    // 3. Calculate total
    order.Total = s.calculateTotal(order.Items)
    
    // 4. Create order
    if err := s.orderRepo.Create(ctx, order); err != nil {
        return err
    }
    
    // 5. Update stock
    for _, item := range order.Items {
        if err := s.productRepo.UpdateStock(ctx, item.ProductID, -item.Quantity); err != nil {
            // Rollback needed!
            return err
        }
    }
    
    return nil
}
```

### Pattern 4: Transactions

```go
func (s *service) TransferBalance(ctx context.Context, fromID, toID uint, amount float64) error {
    // Use transaction
    return database.WithTransaction(ctx, s.db, func(tx *gorm.DB) error {
        // Create transactional repository
        txRepo := s.repo.WithTx(tx)
        
        // Debit from sender
        from, err := txRepo.FindByID(ctx, fromID)
        if err != nil {
            return err
        }
        
        if from.Balance < amount {
            return errors.New("insufficient balance")
        }
        
        from.Balance -= amount
        if err := txRepo.Update(ctx, from); err != nil {
            return err
        }
        
        // Credit to receiver
        to, err := txRepo.FindByID(ctx, toID)
        if err != nil {
            return err
        }
        
        to.Balance += amount
        if err := txRepo.Update(ctx, to); err != nil {
            return err
        }
        
        // Create transaction record
        txRecord := &Transaction{
            FromUserID: fromID,
            ToUserID:   toID,
            Amount:     amount,
        }
        return s.txRepo.Create(ctx, txRecord)
    })
}
```

### Pattern 5: Error Handling

```go
// Define custom errors
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrInvalidPassword  = errors.New("invalid password")
    ErrEmailTaken       = errors.New("email already taken")
    ErrInsufficientFunds = errors.New("insufficient funds")
)

func (s *service) GetUser(ctx context.Context, id uint) (*User, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrUserNotFound
        }
        return nil, fmt.Errorf("database error: %w", err)
    }
    
    if user == nil {
        return nil, ErrUserNotFound
    }
    
    return user, nil
}
```

---

## Business Logic Examples

### Example 1: E-commerce Order

```go
func (s *OrderService) PlaceOrder(ctx context.Context, req *PlaceOrderRequest) (*Order, error) {
    // 1. Validate user
    user, err := s.userRepo.FindByID(ctx, req.UserID)
    if err != nil {
        return nil, ErrUserNotFound
    }
    
    if !user.Active {
        return nil, errors.New("user account is inactive")
    }
    
    // 2. Validate cart
    if len(req.Items) == 0 {
        return nil, errors.New("cart is empty")
    }
    
    // 3. Calculate totals
    var subtotal float64
    for _, item := range req.Items {
        product, err := s.productRepo.FindByID(ctx, item.ProductID)
        if err != nil {
            return nil, fmt.Errorf("product not found: %w", err)
        }
        
        if product.Stock < item.Quantity {
            return nil, fmt.Errorf("insufficient stock for %s", product.Name)
        }
        
        subtotal += product.Price * float64(item.Quantity)
    }
    
    // 4. Apply discount
    discount := s.calculateDiscount(user, subtotal)
    tax := subtotal * 0.1
    total := subtotal - discount + tax
    
    // 5. Validate payment
    if user.Balance < total {
        return nil, ErrInsufficientFunds
    }
    
    // 6. Create order in transaction
    var order *Order
    err = database.WithTransaction(ctx, s.db, func(tx *gorm.DB) error {
        // Create order
        order = &Order{
            UserID:   req.UserID,
            Items:    req.Items,
            Subtotal: subtotal,
            Discount: discount,
            Tax:      tax,
            Total:    total,
            Status:   "pending",
        }
        
        if err := s.orderRepo.Create(ctx, order); err != nil {
            return err
        }
        
        // Deduct stock
        for _, item := range req.Items {
            if err := s.productRepo.UpdateStock(ctx, item.ProductID, -item.Quantity); err != nil {
                return err
            }
        }
        
        // Deduct balance
        if err := s.userRepo.UpdateBalance(ctx, req.UserID, -total); err != nil {
            return err
        }
        
        return nil
    })
    
    if err != nil {
        return nil, err
    }
    
    // 7. Send notifications (async)
    go s.notifyOrderCreated(order)
    
    return order, nil
}

func (s *OrderService) calculateDiscount(user *User, subtotal float64) float64 {
    if user.Role == "premium" {
        return subtotal * 0.1 // 10% discount
    }
    if subtotal > 1000 {
        return subtotal * 0.05 // 5% discount for large orders
    }
    return 0
}
```

### Example 2: User Authentication

```go
func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
    // 1. Find user
    user, err := s.userRepo.FindByEmail(ctx, email)
    if err != nil {
        return nil, ErrInvalidCredentials
    }
    
    if user == nil {
        return nil, ErrInvalidCredentials
    }
    
    // 2. Check if active
    if !user.Active {
        return nil, errors.New("account is deactivated")
    }
    
    // 3. Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        // Log failed attempt
        s.logger.Warn("Failed login attempt", logger.Fields{
            "email": email,
        })
        return nil, ErrInvalidCredentials
    }
    
    // 4. Generate token
    token, err := s.generateJWT(user)
    if err != nil {
        return nil, err
    }
    
    // 5. Update last login
    if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
        s.logger.Error("Failed to update last login", logger.Fields{
            "user_id": user.ID,
            "error":   err.Error(),
        })
    }
    
    // 6. Log success
    s.logger.Info("User logged in", logger.Fields{
        "user_id": user.ID,
        "email":   email,
    })
    
    return &AuthResponse{
        Token: token,
        User:  user,
    }, nil
}
```

---

## Best Practices

### ✅ DO:

**1. Keep Services Focused**
```go
// Good: Focused on user operations
type UserService struct {
    repo Repository
}

// Bad: Too many responsibilities
type MegaService struct {
    userRepo    UserRepository
    orderRepo   OrderRepository
    productRepo ProductRepository
    paymentRepo PaymentRepository
    // ...10 more repos
}
```

**2. Use Dependency Injection**
```go
// Good: Dependencies injected
func NewService(repo Repository, logger logger.Logger) Service {
    return &service{
        repo:   repo,
        logger: logger,
    }
}
```

**3. Return Custom Errors**
```go
// Good: Meaningful errors
if user == nil {
    return ErrUserNotFound
}

// Bad: Generic errors
if user == nil {
    return errors.New("error")
}
```

**4. Log Important Events**
```go
// Good: Log business events
s.logger.Info("Order placed", logger.Fields{
    "order_id": order.ID,
    "user_id":  order.UserID,
    "total":    order.Total,
})
```

### ❌ DON'T:

**1. Put HTTP Logic in Service**
```go
// Bad: HTTP concerns in service
func (s *service) CreateUser(ctx *fiber.Ctx) error {
    var user User
    ctx.BodyParser(&user)  // ❌ HTTP parsing
    // ...
}

// Good: Service is HTTP-agnostic
func (s *service) CreateUser(ctx context.Context, user *User) error {
    // Pure business logic
}
```

**2. Expose Repository Directly**
```go
// Bad: Expose repository
type Service struct {
    Repo Repository  // ❌ Public field
}

// Good: Keep it private
type service struct {
    repo Repository  // ✅ Private field
}
```

**3. Ignore Errors**
```go
// Bad: Silently ignore errors
s.logger.Info("User created")  // What if it failed?

// Good: Handle errors
if err := s.repo.Create(ctx, user); err != nil {
    return err
}
```

---

## Testing

### Unit Test with Mock Repository

```go
func TestUserService_CreateUser(t *testing.T) {
    // Setup
    mockRepo := &MockUserRepository{
        users: make(map[uint]*User),
    }
    mockLogger := &MockLogger{}
    
    service := NewService(mockRepo, mockLogger)
    
    // Test
    user := &User{
        Name:     "John Doe",
        Email:    "john@example.com",
        Password: "secret123",
    }
    
    err := service.CreateUser(context.Background(), user)
    
    // Assert
    assert.NoError(t, err)
    assert.NotEqual(t, "secret123", user.Password) // Password hashed
    assert.True(t, user.Active)
    assert.Equal(t, "user", user.Role)
}

func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
    // Setup with existing user
    mockRepo := &MockUserRepository{
        users: map[uint]*User{
            1: {ID: 1, Email: "existing@example.com"},
        },
    }
    
    service := NewService(mockRepo, &MockLogger{})
    
    // Test
    user := &User{Email: "existing@example.com"}
    err := service.CreateUser(context.Background(), user)
    
    // Assert
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "already exists")
}
```

---

## Next Steps

- [**Repository Pattern**](repository-pattern.md) - Data access layer
- [**Dependency Injection**](dependency-injection.md) - Wire services together
- [**Module System**](module-system.md) - Organize services in modules
- [**Testing**](../development/testing.md) - Test services thoroughly

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
