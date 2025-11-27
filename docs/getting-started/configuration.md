# Configuration

Learn how to configure Neonex Core applications using environment variables and configuration files.

---

## Overview

Neonex Core uses **environment variables** for configuration management, following the [12-Factor App](https://12factor.net/config) methodology.

**Benefits:**
- ✅ Easy to change between deployments without code changes
- ✅ Keep secrets out of version control
- ✅ Environment-specific configuration (dev, staging, prod)
- ✅ Standard practice in cloud deployments

---

## Environment Files

### `.env` - Local Development

Create a `.env` file in your project root:

```bash
# .env
APP_NAME=my-app
APP_ENV=development
APP_PORT=8080

DB_DRIVER=sqlite
DB_NAME=app.db
```

**Important:**
- 🚫 **Never commit `.env` to git**
- ✅ Add to `.gitignore`
- ✅ Use `.env.example` as a template

### `.env.example` - Template for Team

Commit this file to show required variables:

```bash
# .env.example
APP_NAME=neonex-core
APP_ENV=development
APP_PORT=8080

DB_DRIVER=sqlite
DB_HOST=
DB_PORT=
DB_NAME=neonex.db
DB_USER=
DB_PASSWORD=

LOG_LEVEL=info
LOG_FORMAT=json
LOG_FILE=logs/app.log
```

Team members copy this to `.env` and fill in their values.

---

## Configuration Variables

### Application Settings

```bash
# Application name (used in logs and responses)
APP_NAME=my-app

# Environment: development, staging, production
APP_ENV=development

# HTTP server port
APP_PORT=8080

# Enable debug mode (more verbose logging)
DEBUG=true
```

**Usage in code:**
```go
appName := os.Getenv("APP_NAME")
port := os.Getenv("APP_PORT")
if port == "" {
    port = "8080" // Default
}
```

### Database Configuration

#### SQLite (Default)

```bash
DB_DRIVER=sqlite
DB_NAME=neonex.db

# Optional: custom path
# DB_NAME=/data/neonex.db
```

#### PostgreSQL

```bash
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=neonex
DB_USER=postgres
DB_PASSWORD=secret
DB_SSLMODE=disable

# Optional: connection string
# DB_DSN=postgres://user:pass@host:5432/dbname?sslmode=disable
```

#### MySQL

```bash
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_NAME=neonex
DB_USER=root
DB_PASSWORD=secret
DB_CHARSET=utf8mb4
DB_PARSE_TIME=true

# Optional: connection string
# DB_DSN=root:pass@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True
```

#### Turso (LibSQL)

```bash
DB_DRIVER=turso
DB_URL=libsql://your-db.turso.io
DB_TOKEN=your-auth-token

# Optional: local replica
# DB_NAME=local.db
```

### Logging Configuration

```bash
# Log level: debug, info, warn, error, fatal
LOG_LEVEL=info

# Output format: json, text
LOG_FORMAT=json

# Log file path (leave empty for stdout only)
LOG_FILE=logs/app.log

# File rotation settings
LOG_MAX_SIZE=100        # MB per file
LOG_MAX_AGE=30          # days to keep
LOG_MAX_BACKUPS=10      # number of old files
LOG_COMPRESS=true       # compress old files
```

**Log Levels:**
- `debug` - Everything including debug messages
- `info` - General information + warnings + errors
- `warn` - Warnings + errors
- `error` - Errors + fatal
- `fatal` - Only fatal errors

### CORS Configuration

```bash
# Allow all origins (use carefully)
CORS_ALLOW_ORIGINS=*

# Specific origins
# CORS_ALLOW_ORIGINS=http://localhost:3000,https://example.com

# Allow credentials
CORS_ALLOW_CREDENTIALS=true

# Allowed methods
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS

# Allowed headers
CORS_ALLOW_HEADERS=Origin,Content-Type,Accept,Authorization

# Max age for preflight cache
CORS_MAX_AGE=3600
```

---

## Loading Configuration

### Automatic Loading

Neonex Core automatically loads `.env` on startup:

```go
// internal/core/app.go
func NewApp() *App {
    // Load .env file if exists
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment")
    }
    
    // Initialize components
    // ...
}
```

### Manual Loading

Load specific files:

```go
import "github.com/joho/godotenv"

// Load specific file
godotenv.Load(".env.production")

// Load multiple files (first found wins)
godotenv.Load(".env.local", ".env")
```

### Accessing Variables

```go
import "os"

// Get with default
port := os.Getenv("APP_PORT")
if port == "" {
    port = "8080"
}

// Must have
dbHost := os.Getenv("DB_HOST")
if dbHost == "" {
    log.Fatal("DB_HOST is required")
}

// Parse to int
maxConns := 10
if v := os.Getenv("DB_MAX_CONNECTIONS"); v != "" {
    maxConns, _ = strconv.Atoi(v)
}

// Parse to bool
debug := os.Getenv("DEBUG") == "true"
```

---

## Configuration Struct Pattern

### Create Config Type

```go
// internal/config/config.go
package config

import (
    "os"
    "strconv"
)

type Config struct {
    App      AppConfig
    Database DatabaseConfig
    Logger   LoggerConfig
}

type AppConfig struct {
    Name  string
    Env   string
    Port  string
    Debug bool
}

type DatabaseConfig struct {
    Driver   string
    Host     string
    Port     string
    Name     string
    User     string
    Password string
}

type LoggerConfig struct {
    Level      string
    Format     string
    File       string
    MaxSize    int
    MaxAge     int
    MaxBackups int
    Compress   bool
}

func Load() *Config {
    return &Config{
        App: AppConfig{
            Name:  getEnv("APP_NAME", "neonex-core"),
            Env:   getEnv("APP_ENV", "development"),
            Port:  getEnv("APP_PORT", "8080"),
            Debug: getEnv("DEBUG", "false") == "true",
        },
        Database: DatabaseConfig{
            Driver:   getEnv("DB_DRIVER", "sqlite"),
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "5432"),
            Name:     getEnv("DB_NAME", "neonex.db"),
            User:     getEnv("DB_USER", ""),
            Password: getEnv("DB_PASSWORD", ""),
        },
        Logger: LoggerConfig{
            Level:      getEnv("LOG_LEVEL", "info"),
            Format:     getEnv("LOG_FORMAT", "json"),
            File:       getEnv("LOG_FILE", ""),
            MaxSize:    getEnvInt("LOG_MAX_SIZE", 100),
            MaxAge:     getEnvInt("LOG_MAX_AGE", 30),
            MaxBackups: getEnvInt("LOG_MAX_BACKUPS", 10),
            Compress:   getEnv("LOG_COMPRESS", "true") == "true",
        },
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}
```

### Use Config

```go
// main.go
package main

import (
    "github.com/YOUR_USERNAME/neonexcore/internal/config"
    "github.com/YOUR_USERNAME/neonexcore/internal/core"
)

func main() {
    cfg := config.Load()
    
    app := core.NewAppWithConfig(cfg)
    app.Run()
}
```

---

## Environment-Specific Configuration

### Development

```bash
# .env.development
APP_ENV=development
DEBUG=true
LOG_LEVEL=debug
LOG_FORMAT=text

DB_DRIVER=sqlite
DB_NAME=dev.db
```

### Staging

```bash
# .env.staging
APP_ENV=staging
DEBUG=false
LOG_LEVEL=info
LOG_FORMAT=json

DB_DRIVER=postgres
DB_HOST=staging-db.example.com
DB_NAME=neonex_staging
```

### Production

```bash
# .env.production
APP_ENV=production
DEBUG=false
LOG_LEVEL=warn
LOG_FORMAT=json
LOG_FILE=/var/log/neonex/app.log

DB_DRIVER=postgres
DB_HOST=prod-db.example.com
DB_NAME=neonex_prod
DB_SSLMODE=require
```

### Load Based on Environment

```go
env := os.Getenv("GO_ENV")
if env == "" {
    env = "development"
}

envFile := fmt.Sprintf(".env.%s", env)
if err := godotenv.Load(envFile); err != nil {
    godotenv.Load(".env")
}
```

---

## Security Best Practices

### 1. Never Commit Secrets

Add to `.gitignore`:
```gitignore
.env
.env.*
!.env.example
*.db
```

### 2. Use Environment Variables in Production

Don't use `.env` files in production. Set variables directly:

**Docker:**
```dockerfile
ENV DB_HOST=prod-db.example.com
ENV DB_PASSWORD=secure-password
```

**Kubernetes:**
```yaml
env:
  - name: DB_HOST
    value: prod-db.example.com
  - name: DB_PASSWORD
    valueFrom:
      secretKeyRef:
        name: db-secret
        key: password
```

**Cloud Platforms:**
- AWS: Use Systems Manager Parameter Store or Secrets Manager
- Google Cloud: Use Secret Manager
- Azure: Use Key Vault

### 3. Validate Required Variables

```go
func validateConfig() error {
    required := []string{
        "DB_HOST",
        "DB_NAME",
        "DB_USER",
        "DB_PASSWORD",
    }
    
    for _, key := range required {
        if os.Getenv(key) == "" {
            return fmt.Errorf("required environment variable %s is not set", key)
        }
    }
    return nil
}
```

### 4. Use Strong Passwords

```bash
# Generate secure password
openssl rand -base64 32

# Use in .env
DB_PASSWORD=8KfJ3mP9qL2nR5tY7wX0zA4bC6dE1fG
```

---

## Configuration Tips

### Use Defaults Wisely

```go
// Good: Sensible defaults for development
port := getEnv("APP_PORT", "8080")

// Bad: Default secrets (security risk)
// password := getEnv("DB_PASSWORD", "admin123") ❌
```

### Group Related Variables

```bash
# Database config together
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=neonex

# Redis config together
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
```

### Document Variables

Add comments in `.env.example`:
```bash
# Database driver: sqlite, postgres, mysql, turso
DB_DRIVER=sqlite

# For SQLite: file path
# For others: hostname
DB_HOST=localhost

# Server port (default: 8080)
APP_PORT=8080
```

---

## Next Steps

- [CLI Tools](../cli-tools/overview.md) - Learn about `neonex` commands
- [Module System](../core-concepts/module-system.md) - Understand module configuration
- [Database Configuration](../database/configuration.md) - Deep dive into database setup
- [Logging](../logging/configuration.md) - Configure logging behavior

---

**Need help?** Check the [FAQ](../resources/faq.md) or [get support](../resources/support.md).
