# Database Configuration

Complete guide to configuring database connections in Neonex Core.

---

## Configuration Files

### Environment Variables

Primary configuration method using `.env` file:

```bash
# Database Driver
DB_DRIVER=sqlite|postgres|mysql|turso

# Connection Details
DB_HOST=localhost
DB_PORT=5432
DB_NAME=neonex
DB_USER=username
DB_PASSWORD=secret

# SQLite specific
DB_NAME=neonex.db

# PostgreSQL specific
DB_SSLMODE=disable|require|verify-full

# MySQL specific
DB_CHARSET=utf8mb4
DB_PARSE_TIME=true

# Connection Pool
DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=100
DB_CONN_MAX_LIFETIME=3600
```

---

## Driver-Specific Configuration

### SQLite

**Minimal Configuration:**
```bash
DB_DRIVER=sqlite
DB_NAME=app.db
```

**Advanced Options:**
```bash
DB_DRIVER=sqlite
DB_NAME=app.db
DB_SQLITE_CACHE=shared
DB_SQLITE_MODE=rwc
DB_SQLITE_JOURNAL=WAL
```

**Connection String:**
```
file:app.db?cache=shared&mode=rwc&_journal_mode=WAL
```

**Use Cases:**
- Development
- Testing
- Small applications
- Embedded systems

### PostgreSQL

**Basic Configuration:**
```bash
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=neonex
DB_USER=postgres
DB_PASSWORD=secret
DB_SSLMODE=disable
```

**With SSL:**
```bash
DB_DRIVER=postgres
DB_HOST=prod-db.example.com
DB_PORT=5432
DB_NAME=neonex_prod
DB_USER=app_user
DB_PASSWORD=secure_password
DB_SSLMODE=require
DB_SSLCERT=/path/to/client-cert.pem
DB_SSLKEY=/path/to/client-key.pem
DB_SSLROOTCERT=/path/to/ca-cert.pem
```

**Connection String:**
```bash
DB_DSN=postgres://user:pass@localhost:5432/dbname?sslmode=disable
```

**Connection Pooling:**
```bash
DB_MAX_IDLE_CONNS=25
DB_MAX_OPEN_CONNS=100
DB_CONN_MAX_LIFETIME=3600
```

### MySQL

**Basic Configuration:**
```bash
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_NAME=neonex
DB_USER=root
DB_PASSWORD=secret
DB_CHARSET=utf8mb4
DB_PARSE_TIME=true
DB_LOC=Local
```

**Connection String:**
```bash
DB_DSN=root:pass@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
```

**With Connection Pool:**
```bash
DB_DRIVER=mysql
DB_HOST=mysql.example.com
DB_PORT=3306
DB_NAME=production
DB_USER=app_user
DB_PASSWORD=secure_pass
DB_CHARSET=utf8mb4
DB_PARSE_TIME=true
DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=50
```

### Turso (LibSQL)

**Configuration:**
```bash
DB_DRIVER=turso
DB_URL=libsql://your-database.turso.io
DB_TOKEN=your_auth_token
```

**With Local Replica:**
```bash
DB_DRIVER=turso
DB_URL=libsql://your-database.turso.io
DB_TOKEN=your_auth_token
DB_NAME=local-replica.db
```

**Features:**
- Edge deployment
- Global replication
- SQLite compatibility
- Low latency reads

---

## Configuration Loading

### Automatic Loading

Framework loads configuration automatically:

```go
// internal/config/database.go
func LoadDatabaseConfig() *DatabaseConfig {
    return &DatabaseConfig{
        Driver:   os.Getenv("DB_DRIVER"),
        Host:     os.Getenv("DB_HOST"),
        Port:     os.Getenv("DB_PORT"),
        Name:     os.Getenv("DB_NAME"),
        User:     os.Getenv("DB_USER"),
        Password: os.Getenv("DB_PASSWORD"),
    }
}
```

### Manual Configuration

Override with code:

```go
config := &config.DatabaseConfig{
    Driver:   "postgres",
    Host:     "localhost",
    Port:     "5432",
    Name:     "mydb",
    User:     "myuser",
    Password: "mypass",
}

db, err := config.InitDatabaseWithConfig(config)
```

---

## Connection Pooling

### Configuration

```bash
# Maximum idle connections
DB_MAX_IDLE_CONNS=10

# Maximum open connections
DB_MAX_OPEN_CONNS=100

# Maximum connection lifetime (seconds)
DB_CONN_MAX_LIFETIME=3600
```

### Implementation

```go
sqlDB, err := db.DB()
if err != nil {
    return err
}

// Set pool settings
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### Recommendations

| Environment | Max Idle | Max Open | Lifetime |
|-------------|----------|----------|----------|
| Development | 5 | 20 | 1 hour |
| Staging | 10 | 50 | 1 hour |
| Production | 25 | 100 | 30 mins |

---

## SSL/TLS Configuration

### PostgreSQL SSL

```bash
# Disable SSL (development)
DB_SSLMODE=disable

# Require SSL (production)
DB_SSLMODE=require

# Verify certificate
DB_SSLMODE=verify-full
DB_SSLROOTCERT=/path/to/ca.pem
```

### MySQL SSL

```bash
DB_DRIVER=mysql
DB_HOST=secure-mysql.example.com
DB_TLS=true
DB_TLS_CERT=/path/to/client-cert.pem
DB_TLS_KEY=/path/to/client-key.pem
DB_TLS_CA=/path/to/ca.pem
```

---

## Environment-Specific Config

### Development

```bash
# .env.development
DB_DRIVER=sqlite
DB_NAME=dev.db
DB_LOG_LEVEL=debug
```

### Staging

```bash
# .env.staging
DB_DRIVER=postgres
DB_HOST=staging-db.example.com
DB_PORT=5432
DB_NAME=neonex_staging
DB_USER=staging_user
DB_PASSWORD=${STAGING_DB_PASSWORD}
DB_SSLMODE=require
DB_MAX_OPEN_CONNS=50
```

### Production

```bash
# .env.production
DB_DRIVER=postgres
DB_HOST=prod-db.example.com
DB_PORT=5432
DB_NAME=neonex_prod
DB_USER=prod_user
DB_PASSWORD=${PROD_DB_PASSWORD}
DB_SSLMODE=verify-full
DB_MAX_IDLE_CONNS=25
DB_MAX_OPEN_CONNS=100
DB_CONN_MAX_LIFETIME=1800
```

---

## Advanced Configuration

### Read Replicas

```go
// Master for writes
masterDB, _ := gorm.Open(postgres.Open(masterDSN))

// Replica for reads
replicaDB, _ := gorm.Open(postgres.Open(replicaDSN))

// Use replicas plugin
db.Use(dbresolver.Register(dbresolver.Config{
    Replicas: []gorm.Dialector{
        postgres.Open(replica1DSN),
        postgres.Open(replica2DSN),
    },
    Policy: dbresolver.RandomPolicy{},
}))
```

### Sharding

```go
// Shard by user ID
func GetShardDB(userID uint) *gorm.DB {
    shardNum := userID % 4
    return shards[shardNum]
}
```

### Custom Logger

```go
import (
    "gorm.io/gorm/logger"
)

newLogger := logger.New(
    log.New(os.Stdout, "\r\n", log.LstdFlags),
    logger.Config{
        SlowThreshold:             time.Second,
        LogLevel:                  logger.Info,
        IgnoreRecordNotFoundError: true,
        Colorful:                  true,
    },
)

db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
    Logger: newLogger,
})
```

---

## Security Best Practices

### ✅ DO:

**1. Use Environment Variables**
```bash
# Good: From environment
DB_PASSWORD=${DATABASE_PASSWORD}
```

**2. Restrict Database User**
```sql
-- Good: Minimal permissions
GRANT SELECT, INSERT, UPDATE, DELETE ON neonex.* TO 'app_user'@'%';
```

**3. Use SSL in Production**
```bash
# Good: Encrypted connection
DB_SSLMODE=require
```

**4. Rotate Credentials**
```bash
# Good: Regular rotation
# Update password monthly
```

### ❌ DON'T:

**1. Commit Secrets**
```bash
# Bad: Hardcoded password
DB_PASSWORD=mysecretpass123
```

**2. Use Root User**
```bash
# Bad: Too permissive
DB_USER=root
```

**3. Disable SSL**
```bash
# Bad: Unencrypted in production
DB_SSLMODE=disable
```

---

## Validation

### Connection Test

```go
func TestConnection(db *gorm.DB) error {
    sqlDB, err := db.DB()
    if err != nil {
        return err
    }
    
    if err := sqlDB.Ping(); err != nil {
        return fmt.Errorf("ping failed: %w", err)
    }
    
    return nil
}
```

### Health Check Endpoint

```go
app.Get("/health/db", func(c *fiber.Ctx) error {
    sqlDB, _ := db.DB()
    
    if err := sqlDB.Ping(); err != nil {
        return c.Status(503).JSON(fiber.Map{
            "status": "unhealthy",
            "error":  err.Error(),
        })
    }
    
    return c.JSON(fiber.Map{
        "status": "healthy",
    })
})
```

---

## Troubleshooting

### Connection Failed

**Problem:**
```
failed to connect to database
```

**Solutions:**
1. Check credentials
2. Verify host/port
3. Check firewall rules
4. Test with database client

### Too Many Connections

**Problem:**
```
too many connections
```

**Solutions:**
```bash
# Reduce max connections
DB_MAX_OPEN_CONNS=50

# Or increase database limit
# PostgreSQL: max_connections = 200
# MySQL: max_connections = 200
```

### SSL Required

**Problem:**
```
SSL connection required
```

**Solution:**
```bash
DB_SSLMODE=require
# Or provide certificates
```

---

## Next Steps

- [**Migrations**](migrations.md) - Schema management
- [**Transactions**](transactions.md) - Transaction handling
- [**Repositories**](repositories.md) - Data access patterns
- [**Performance**](../advanced/performance.md) - Optimization tips

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
