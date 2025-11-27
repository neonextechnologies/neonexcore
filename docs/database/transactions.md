# Database Transactions

Learn how to manage database transactions in Neonex Core for ACID-compliant operations.

---

## Overview

**Transactions** ensure data integrity by grouping multiple database operations into an atomic unit. Either all operations succeed (commit) or all fail (rollback).

**ACID Properties:**
- **Atomicity** - All or nothing
- **Consistency** - Valid state transitions
- **Isolation** - Concurrent transactions don't interfere
- **Durability** - Committed changes persist

---

## Transaction Manager

### Basic Usage

```go
// pkg/database/transaction.go
type TxManager struct {
    db *gorm.DB
}

func NewTxManager(db *gorm.DB) *TxManager {
    return &TxManager{db: db}
}
```

### Automatic Transaction

Recommended approach with automatic commit/rollback:

```go
func (s *service) TransferMoney(ctx context.Context, fromID, toID uint, amount float64) error {
    txManager := database.NewTxManager(s.db)
    
    return txManager.WithTransaction(ctx, func(tx *gorm.DB) error {
        // All operations in this function are part of the transaction
        
        // 1. Debit from sender
        var sender Account
        if err := tx.First(&sender, fromID).Error; err != nil {
            return err
        }
        
        if sender.Balance < amount {
            return errors.New("insufficient funds")
        }
        
        sender.Balance -= amount
        if err := tx.Save(&sender).Error; err != nil {
            return err // Rollback automatically
        }
        
        // 2. Credit to receiver
        var receiver Account
        if err := tx.First(&receiver, toID).Error; err != nil {
            return err // Rollback automatically
        }
        
        receiver.Balance += amount
        if err := tx.Save(&receiver).Error; err != nil {
            return err // Rollback automatically
        }
        
        // 3. Create transaction record
        record := &Transaction{
            FromID: fromID,
            ToID:   toID,
            Amount: amount,
        }
        if err := tx.Create(record).Error; err != nil {
            return err // Rollback automatically
        }
        
        // If we reach here without errors, transaction commits automatically
        return nil
    })
}
```

**How it works:**
- Returns `nil` → Transaction commits
- Returns error → Transaction rolls back
- Panic → Transaction rolls back

---

## Manual Transaction Control

### Begin, Commit, Rollback

```go
func (s *service) ManualTransaction(ctx context.Context) error {
    txManager := database.NewTxManager(s.db)
    
    // Begin transaction
    tx := txManager.BeginTx(ctx)
    
    // Ensure rollback on panic
    defer func() {
        if r := recover(); r != nil {
            txManager.Rollback(tx)
            panic(r)
        }
    }()
    
    // Operation 1
    if err := tx.Create(&user).Error; err != nil {
        txManager.Rollback(tx)
        return err
    }
    
    // Operation 2
    if err := tx.Create(&profile).Error; err != nil {
        txManager.Rollback(tx)
        return err
    }
    
    // Commit
    return txManager.Commit(tx)
}
```

### Savepoints

```go
func (s *service) WithSavepoint(ctx context.Context) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Operation 1
        if err := tx.Create(&user).Error; err != nil {
            return err
        }
        
        // Create savepoint
        tx.SavePoint("sp1")
        
        // Operation 2 (might fail)
        if err := tx.Create(&optional).Error; err != nil {
            // Rollback to savepoint
            tx.RollbackTo("sp1")
            // Continue without failing
        }
        
        // Operation 3
        return tx.Create(&required).Error
    })
}
```

---

## Common Patterns

### Pattern 1: Create with Relationships

```go
func (s *OrderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 1. Create order
        order := &Order{
            UserID: req.UserID,
            Total:  req.Total,
            Status: "pending",
        }
        if err := tx.Create(order).Error; err != nil {
            return err
        }
        
        // 2. Create order items
        for _, item := range req.Items {
            orderItem := &OrderItem{
                OrderID:   order.ID,
                ProductID: item.ProductID,
                Quantity:  item.Quantity,
                Price:     item.Price,
            }
            if err := tx.Create(orderItem).Error; err != nil {
                return err
            }
            
            // 3. Update product stock
            if err := tx.Model(&Product{}).
                Where("id = ?", item.ProductID).
                Update("stock", gorm.Expr("stock - ?", item.Quantity)).
                Error; err != nil {
                return err
            }
        }
        
        return nil
    })
}
```

### Pattern 2: Update Multiple Records

```go
func (s *UserService) DeactivateUsers(ctx context.Context, userIDs []uint) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Update users
        if err := tx.Model(&User{}).
            Where("id IN ?", userIDs).
            Update("active", false).
            Error; err != nil {
            return err
        }
        
        // Log deactivation
        for _, id := range userIDs {
            log := &AuditLog{
                UserID: id,
                Action: "deactivated",
            }
            if err := tx.Create(log).Error; err != nil {
                return err
            }
        }
        
        return nil
    })
}
```

### Pattern 3: Conditional Operations

```go
func (s *InventoryService) ProcessOrder(ctx context.Context, orderID uint) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Get order
        var order Order
        if err := tx.Preload("Items").First(&order, orderID).Error; err != nil {
            return err
        }
        
        // Check inventory for all items
        for _, item := range order.Items {
            var product Product
            if err := tx.First(&product, item.ProductID).Error; err != nil {
                return err
            }
            
            if product.Stock < item.Quantity {
                return fmt.Errorf("insufficient stock for product %d", product.ID)
            }
        }
        
        // All checks passed, update inventory
        for _, item := range order.Items {
            if err := tx.Model(&Product{}).
                Where("id = ?", item.ProductID).
                Update("stock", gorm.Expr("stock - ?", item.Quantity)).
                Error; err != nil {
                return err
            }
        }
        
        // Update order status
        order.Status = "confirmed"
        return tx.Save(&order).Error
    })
}
```

---

## Transaction with Repository

### Repository with Transaction

```go
// Repository method that accepts transaction
func (r *repository) WithTx(tx *gorm.DB) Repository {
    return &repository{
        BaseRepository: database.NewBaseRepository[User](tx),
    }
}

// Service using transactional repository
func (s *service) ComplexOperation(ctx context.Context) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Create transactional repository
        txRepo := s.repo.WithTx(tx)
        
        // Use repository with transaction
        user := &User{Name: "John"}
        if err := txRepo.Create(ctx, user); err != nil {
            return err
        }
        
        profile := &Profile{UserID: user.ID}
        if err := txRepo.Create(ctx, profile); err != nil {
            return err
        }
        
        return nil
    })
}
```

### Multiple Repositories

```go
func (s *OrderService) CreateOrderWithPayment(ctx context.Context, req *OrderRequest) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Create transactional repositories
        orderRepo := s.orderRepo.WithTx(tx)
        paymentRepo := s.paymentRepo.WithTx(tx)
        productRepo := s.productRepo.WithTx(tx)
        
        // 1. Create order
        order := &Order{
            UserID: req.UserID,
            Total:  req.Total,
        }
        if err := orderRepo.Create(ctx, order); err != nil {
            return err
        }
        
        // 2. Process payment
        payment := &Payment{
            OrderID: order.ID,
            Amount:  req.Total,
            Status:  "completed",
        }
        if err := paymentRepo.Create(ctx, payment); err != nil {
            return err
        }
        
        // 3. Update inventory
        for _, item := range req.Items {
            if err := productRepo.UpdateStock(ctx, item.ProductID, -item.Quantity); err != nil {
                return err
            }
        }
        
        return nil
    })
}
```

---

## Nested Transactions

### Using Savepoints

```go
func (s *service) NestedOperation(ctx context.Context) error {
    return s.db.Transaction(func(tx1 *gorm.DB) error {
        // Outer transaction operations
        if err := tx1.Create(&user).Error; err != nil {
            return err
        }
        
        // Nested transaction (uses savepoint)
        err := tx1.Transaction(func(tx2 *gorm.DB) error {
            // Inner transaction operations
            if err := tx2.Create(&profile).Error; err != nil {
                return err // Rollback inner only
            }
            return nil
        })
        
        if err != nil {
            // Handle inner transaction failure
            // Outer transaction continues
        }
        
        return nil
    })
}
```

---

## Isolation Levels

### Setting Isolation Level

```go
// Read Uncommitted
tx := db.Begin(&sql.TxOptions{
    Isolation: sql.LevelReadUncommitted,
})

// Read Committed (default for most databases)
tx := db.Begin(&sql.TxOptions{
    Isolation: sql.LevelReadCommitted,
})

// Repeatable Read
tx := db.Begin(&sql.TxOptions{
    Isolation: sql.LevelRepeatableRead,
})

// Serializable (highest isolation)
tx := db.Begin(&sql.TxOptions{
    Isolation: sql.LevelSerializable,
})
```

### Example

```go
func (s *service) WithIsolation(ctx context.Context) error {
    sqlDB, _ := s.db.DB()
    
    // Begin with specific isolation level
    tx, err := sqlDB.BeginTx(ctx, &sql.TxOptions{
        Isolation: sql.LevelSerializable,
    })
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // Use GORM with this transaction
    gormTx := s.db.WithContext(ctx).Begin(&sql.TxOptions{
        Isolation: sql.LevelSerializable,
    })
    defer gormTx.Rollback()
    
    // Operations...
    if err := gormTx.Create(&user).Error; err != nil {
        return err
    }
    
    return gormTx.Commit().Error
}
```

---

## Error Handling

### Rollback on Error

```go
func (s *service) SafeTransaction(ctx context.Context) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Operation 1
        if err := tx.Create(&user).Error; err != nil {
            // Log error
            s.logger.Error("User creation failed", logger.Fields{
                "error": err.Error(),
            })
            return err // Automatic rollback
        }
        
        // Operation 2
        if err := tx.Create(&profile).Error; err != nil {
            s.logger.Error("Profile creation failed", logger.Fields{
                "user_id": user.ID,
                "error":   err.Error(),
            })
            return err // Automatic rollback
        }
        
        s.logger.Info("Transaction completed", logger.Fields{
            "user_id": user.ID,
        })
        
        return nil // Automatic commit
    })
}
```

### Custom Error Types

```go
var (
    ErrInsufficientFunds = errors.New("insufficient funds")
    ErrInvalidAmount     = errors.New("invalid amount")
    ErrAccountNotFound   = errors.New("account not found")
)

func (s *service) Transfer(ctx context.Context, fromID, toID uint, amount float64) error {
    if amount <= 0 {
        return ErrInvalidAmount
    }
    
    return s.db.Transaction(func(tx *gorm.DB) error {
        var from, to Account
        
        if err := tx.First(&from, fromID).Error; err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                return ErrAccountNotFound
            }
            return err
        }
        
        if from.Balance < amount {
            return ErrInsufficientFunds
        }
        
        // Process transfer...
        return nil
    })
}
```

---

## Performance Considerations

### Batch Operations

```go
func (s *service) BatchInsert(ctx context.Context, users []User) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Insert in batches of 1000
        return tx.CreateInBatches(users, 1000).Error
    })
}
```

### Disable Hooks in Transaction

```go
func (s *service) BulkUpdate(ctx context.Context) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Disable hooks for performance
        return tx.Session(&gorm.Session{
            SkipHooks: true,
        }).Model(&User{}).
            Where("active = ?", false).
            Update("status", "inactive").
            Error
    })
}
```

### Lock Records

```go
func (s *service) UpdateWithLock(ctx context.Context, userID uint) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        var user User
        
        // SELECT ... FOR UPDATE
        if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
            First(&user, userID).Error; err != nil {
            return err
        }
        
        // Update locked record
        user.Balance += 100
        return tx.Save(&user).Error
    })
}
```

---

## Testing Transactions

### Test Rollback

```go
func TestTransactionRollback(t *testing.T) {
    db := setupTestDB(t)
    defer teardownTestDB(t, db)
    
    service := NewService(db)
    
    // This should rollback
    err := service.FailingTransaction(context.Background())
    assert.Error(t, err)
    
    // Verify rollback
    var count int64
    db.Model(&User{}).Count(&count)
    assert.Equal(t, int64(0), count)
}
```

### Test Commit

```go
func TestTransactionCommit(t *testing.T) {
    db := setupTestDB(t)
    defer teardownTestDB(t, db)
    
    service := NewService(db)
    
    // This should commit
    err := service.SuccessfulTransaction(context.Background())
    assert.NoError(t, err)
    
    // Verify commit
    var users []User
    db.Find(&users)
    assert.Len(t, users, 2)
}
```

---

## Best Practices

### ✅ DO:

**1. Keep Transactions Short**
```go
// Good: Quick transaction
return db.Transaction(func(tx *gorm.DB) error {
    return tx.Create(&user).Error
})
```

**2. Handle Errors Properly**
```go
// Good: Check all errors
if err := tx.Create(&user).Error; err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}
```

**3. Use Context**
```go
// Good: Respects cancellation
return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    // Operations...
})
```

**4. Log Transaction Events**
```go
// Good: Track transactions
logger.Info("Starting transaction")
err := db.Transaction(func(tx *gorm.DB) error {
    // ...
})
if err != nil {
    logger.Error("Transaction failed", logger.Fields{"error": err})
}
```

### ❌ DON'T:

**1. Long-Running Transactions**
```go
// Bad: Holds locks too long
return db.Transaction(func(tx *gorm.DB) error {
    time.Sleep(10 * time.Second) // ❌
    return tx.Create(&user).Error
})
```

**2. External API Calls in Transaction**
```go
// Bad: Network calls in transaction
return db.Transaction(func(tx *gorm.DB) error {
    tx.Create(&order)
    http.Post("https://api.example.com/notify", ...) // ❌
    return nil
})
```

**3. Ignore Transaction Errors**
```go
// Bad: No error handling
db.Transaction(func(tx *gorm.DB) error {
    tx.Create(&user)
    return nil
})
```

---

## Troubleshooting

### Deadlocks

**Problem:**
```
Error: deadlock detected
```

**Solution:**
- Access tables in same order
- Use shorter transactions
- Reduce isolation level if appropriate

### Transaction Timeout

**Problem:**
```
Error: transaction timeout
```

**Solution:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    // Operations...
})
```

---

## Next Steps

- [**Repositories**](repositories.md) - Repository pattern with transactions
- [**Migrations**](migrations.md) - Schema management
- [**Performance**](../advanced/performance.md) - Optimization tips
- [**Testing**](../development/testing.md) - Testing transactions

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
