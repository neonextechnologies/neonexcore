# Environment Variables

Comprehensive guide to environment variable configuration in Neonex Core.

---

## Overview

Environment variables provide:
- **Configuration** - External configuration without code changes
- **Security** - Keep secrets out of version control
- **Flexibility** - Different settings per environment
- **12-Factor** - Follow best practices

---

## Core Variables

### Application Settings

```bash
# Application Name
APP_NAME=neonex-core
# Description: Application identifier used in logs and responses
# Default: neonex-core
# Required: No

# Environment
APP_ENV=production
# Description: Application environment (development, staging, production)
# Default: development
# Required: Yes
# Values: development | staging | production

# Port
APP_PORT=8080
# Description: HTTP server port
# Default: 8080
# Required: No

# Debug Mode
DEBUG=false
# Description: Enable debug mode (verbose logging, detailed errors)
# Default: false
# Required: No
# Values: true | false
```

### Database Configuration

```bash
# Driver
DB_DRIVER=postgres
# Description: Database driver
# Default: sqlite
# Required: Yes
# Values: sqlite | postgres | mysql | turso

# Host
DB_HOST=localhost
# Description: Database server hostname or IP
# Default: localhost
# Required: Yes (except for SQLite)

# Port
DB_PORT=5432
# Description: Database server port
# Default: 5432 (postgres), 3306 (mysql)
# Required: Yes (except for SQLite)

# Database Name
DB_NAME=neonex
# Description: Database name or file path (for SQLite)
# Default: neonex.db
# Required: Yes

# Username
DB_USER=neonex
# Description: Database username
# Default: (empty)
# Required: Yes (except for SQLite)

# Password
DB_PASSWORD=secure-password
# Description: Database password
# Default: (empty)
# Required: Yes (except for SQLite)
# Security: Store securely, never commit

# SSL Mode
DB_SSL_MODE=require
# Description: SSL/TLS mode for database connection
# Default: disable
# Required: No
# Values: disable | require | verify-ca | verify-full

# Connection Pool
DB_MAX_OPEN_CONNS=25
# Description: Maximum number of open connections
# Default: 25
# Required: No

DB_MAX_IDLE_CONNS=5
# Description: Maximum number of idle connections
# Default: 5
# Required: No

DB_CONN_MAX_LIFETIME=5m
# Description: Maximum lifetime of connections
# Default: 5m
# Required: No
# Format: Duration (e.g., 5m, 1h, 30s)
```

### Logging Configuration

```bash
# Log Level
LOG_LEVEL=info
# Description: Minimum log level
# Default: info
# Required: No
# Values: debug | info | warn | error | fatal

# Log Format
LOG_FORMAT=json
# Description: Log output format
# Default: json
# Required: No
# Values: json | text

# Log Output
LOG_OUTPUT=both
# Description: Log output destination
# Default: console
# Required: No
# Values: console | file | both

# Log File Path
LOG_FILE=/var/log/neonex/app.log
# Description: Log file path (when output is file or both)
# Default: logs/app.log
# Required: No

# Log Rotation
LOG_MAX_SIZE=100
# Description: Maximum log file size in MB before rotation
# Default: 100
# Required: No

LOG_MAX_AGE=30
# Description: Maximum days to retain old log files
# Default: 30
# Required: No

LOG_MAX_BACKUPS=10
# Description: Maximum number of old log files to keep
# Default: 10
# Required: No

LOG_COMPRESS=true
# Description: Compress rotated log files
# Default: true
# Required: No
# Values: true | false
```

### Security Variables

```bash
# JWT Secret
JWT_SECRET=your-jwt-secret-key-here-min-32-chars
# Description: Secret key for JWT token signing
# Default: (none)
# Required: Yes in production
# Security: Generate with: openssl rand -base64 32

# JWT Expiration
JWT_EXPIRATION=24h
# Description: JWT token expiration time
# Default: 24h
# Required: No
# Format: Duration (e.g., 24h, 7d, 30m)

# Encryption Key
ENCRYPTION_KEY=your-32-byte-encryption-key-here
# Description: AES-256 encryption key (32 bytes)
# Default: (none)
# Required: Yes if using encryption
# Security: Generate with: openssl rand -hex 32

# API Keys
API_KEY_PUBLIC=pk_live_xxxxxxxxxxxxx
# Description: Public API key for external services
# Default: (none)
# Required: No

API_KEY_SECRET=sk_live_xxxxxxxxxxxxx
# Description: Secret API key for external services
# Default: (none)
# Required: No
# Security: Never expose publicly
```

### CORS Configuration

```bash
# Allow Origins
CORS_ALLOW_ORIGINS=https://app.example.com,https://admin.example.com
# Description: Comma-separated list of allowed origins
# Default: *
# Required: No
# Note: Use * only in development

# Allow Methods
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS
# Description: Comma-separated list of allowed HTTP methods
# Default: GET,POST,PUT,DELETE,OPTIONS
# Required: No

# Allow Headers
CORS_ALLOW_HEADERS=Origin,Content-Type,Accept,Authorization
# Description: Comma-separated list of allowed headers
# Default: Origin,Content-Type,Accept
# Required: No

# Allow Credentials
CORS_ALLOW_CREDENTIALS=true
# Description: Allow credentials (cookies, authorization headers)
# Default: false
# Required: No
# Values: true | false

# Max Age
CORS_MAX_AGE=3600
# Description: Preflight cache duration in seconds
# Default: 3600
# Required: No
```

### Rate Limiting

```bash
# Rate Limit Max
RATE_LIMIT_MAX=100
# Description: Maximum requests per window
# Default: 100
# Required: No

# Rate Limit Window
RATE_LIMIT_WINDOW=60s
# Description: Time window for rate limiting
# Default: 60s
# Required: No
# Format: Duration (e.g., 60s, 1m, 5m)

# Rate Limit Strategy
RATE_LIMIT_STRATEGY=ip
# Description: Rate limiting strategy
# Default: ip
# Required: No
# Values: ip | user | api_key
```

### External Services

```bash
# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=redis-password
REDIS_DB=0
# Description: Redis connection settings
# Required: No (only if using Redis)

# Email (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=noreply@example.com
SMTP_PASSWORD=email-password
SMTP_FROM=Neonex <noreply@example.com>
# Description: SMTP settings for sending emails
# Required: No (only if using email)

# AWS S3
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
AWS_S3_BUCKET=my-bucket
# Description: AWS S3 storage configuration
# Required: No (only if using S3)

# Sentry (Error Tracking)
SENTRY_DSN=https://xxx@sentry.io/xxx
SENTRY_ENVIRONMENT=production
SENTRY_TRACES_SAMPLE_RATE=0.1
# Description: Sentry error tracking configuration
# Required: No (only if using Sentry)
```

---

## Environment-Specific Configurations

### Development (.env.development)

```bash
APP_ENV=development
DEBUG=true
LOG_LEVEL=debug
LOG_FORMAT=text

DB_DRIVER=sqlite
DB_NAME=dev.db

CORS_ALLOW_ORIGINS=*

RATE_LIMIT_MAX=1000
```

### Staging (.env.staging)

```bash
APP_ENV=staging
DEBUG=false
LOG_LEVEL=info
LOG_FORMAT=json

DB_DRIVER=postgres
DB_HOST=staging-db.example.com
DB_NAME=neonex_staging
DB_SSL_MODE=require

CORS_ALLOW_ORIGINS=https://staging.example.com

RATE_LIMIT_MAX=500
```

### Production (.env.production)

```bash
APP_ENV=production
DEBUG=false
LOG_LEVEL=warn
LOG_FORMAT=json
LOG_FILE=/var/log/neonex/app.log

DB_DRIVER=postgres
DB_HOST=prod-db.example.com
DB_NAME=neonex_prod
DB_SSL_MODE=require
DB_MAX_OPEN_CONNS=25

JWT_SECRET=${JWT_SECRET}
ENCRYPTION_KEY=${ENCRYPTION_KEY}

CORS_ALLOW_ORIGINS=https://app.example.com
CORS_ALLOW_CREDENTIALS=true

RATE_LIMIT_MAX=100

SENTRY_DSN=${SENTRY_DSN}
```

---

## Loading Environment Variables

### Using godotenv

```go
import "github.com/joho/godotenv"

func main() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
    
    // Load specific file
    godotenv.Load(".env.production")
    
    // Load multiple files (first found wins)
    godotenv.Load(".env.local", ".env")
    
    // Your application code
    app := core.NewApp()
    app.Run()
}
```

### Environment-Based Loading

```go
func loadEnv() {
    env := os.Getenv("GO_ENV")
    if env == "" {
        env = "development"
    }
    
    // Load environment-specific file
    envFile := fmt.Sprintf(".env.%s", env)
    if err := godotenv.Load(envFile); err != nil {
        log.Printf("No %s file found, using defaults", envFile)
        godotenv.Load(".env")
    }
}
```

---

## Validation

### Required Variables Check

```go
func validateEnv() error {
    required := []string{
        "APP_ENV",
        "DB_DRIVER",
        "DB_NAME",
    }
    
    if os.Getenv("APP_ENV") == "production" {
        required = append(required,
            "DB_HOST",
            "DB_USER",
            "DB_PASSWORD",
            "JWT_SECRET",
        )
    }
    
    var missing []string
    for _, key := range required {
        if os.Getenv(key) == "" {
            missing = append(missing, key)
        }
    }
    
    if len(missing) > 0 {
        return fmt.Errorf("missing required environment variables: %v", missing)
    }
    
    return nil
}
```

### Type Validation

```go
func getEnvInt(key string, defaultValue int) int {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    
    intValue, err := strconv.Atoi(value)
    if err != nil {
        log.Printf("Invalid integer for %s, using default: %d", key, defaultValue)
        return defaultValue
    }
    
    return intValue
}

func getEnvBool(key string, defaultValue bool) bool {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    
    boolValue, err := strconv.ParseBool(value)
    if err != nil {
        log.Printf("Invalid boolean for %s, using default: %v", key, defaultValue)
        return defaultValue
    }
    
    return boolValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    
    duration, err := time.ParseDuration(value)
    if err != nil {
        log.Printf("Invalid duration for %s, using default: %v", key, defaultValue)
        return defaultValue
    }
    
    return duration
}
```

---

## Secrets Management

### AWS Secrets Manager

```go
import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/secretsmanager"
)

func loadSecrets() error {
    sess := session.Must(session.NewSession())
    svc := secretsmanager.New(sess)
    
    input := &secretsmanager.GetSecretValueInput{
        SecretId: aws.String("neonex/production"),
    }
    
    result, err := svc.GetSecretValue(input)
    if err != nil {
        return err
    }
    
    // Parse secrets JSON
    var secrets map[string]string
    json.Unmarshal([]byte(*result.SecretString), &secrets)
    
    // Set environment variables
    for key, value := range secrets {
        os.Setenv(key, value)
    }
    
    return nil
}
```

### HashiCorp Vault

```go
import vault "github.com/hashicorp/vault/api"

func loadVaultSecrets() error {
    config := vault.DefaultConfig()
    config.Address = os.Getenv("VAULT_ADDR")
    
    client, err := vault.NewClient(config)
    if err != nil {
        return err
    }
    
    client.SetToken(os.Getenv("VAULT_TOKEN"))
    
    secret, err := client.Logical().Read("secret/data/neonex/production")
    if err != nil {
        return err
    }
    
    data := secret.Data["data"].(map[string]interface{})
    
    for key, value := range data {
        os.Setenv(key, value.(string))
    }
    
    return nil
}
```

---

## Docker Integration

### Docker Compose

```yaml
services:
  app:
    image: neonex-app:latest
    environment:
      - APP_ENV=production
      - DB_HOST=postgres
      - DB_PASSWORD=${DB_PASSWORD}
    env_file:
      - .env
      - .env.production
```

### Docker Run

```bash
# With env file
docker run --env-file .env neonex-app:latest

# With individual variables
docker run \
  -e APP_ENV=production \
  -e DB_HOST=postgres \
  -e DB_PASSWORD=secret \
  neonex-app:latest

# From host environment
docker run \
  -e DB_PASSWORD \
  neonex-app:latest
```

---

## Kubernetes

### ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: neonex-config
data:
  APP_ENV: "production"
  LOG_LEVEL: "info"
  DB_HOST: "postgres-service"
```

### Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: neonex-secrets
type: Opaque
stringData:
  DB_PASSWORD: "secure-password"
  JWT_SECRET: "jwt-secret-key"
```

### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: neonex-app
spec:
  template:
    spec:
      containers:
      - name: neonex
        image: neonex-app:latest
        envFrom:
        - configMapRef:
            name: neonex-config
        env:
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: neonex-secrets
              key: DB_PASSWORD
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: neonex-secrets
              key: JWT_SECRET
```

---

## Security Best Practices

### ✅ DO: Use Strong Secrets

```bash
# Generate JWT secret (32+ characters)
openssl rand -base64 32

# Generate encryption key (32 bytes for AES-256)
openssl rand -hex 32

# Generate random password
openssl rand -base64 24
```

### ✅ DO: Never Commit Secrets

```gitignore
# .gitignore
.env
.env.local
.env.*.local
.env.production
*.key
*.pem
secrets/
```

### ✅ DO: Use .env.example

```bash
# .env.example (commit this)
APP_ENV=development
DB_PASSWORD=change-me
JWT_SECRET=generate-with-openssl-rand

# .env (never commit)
APP_ENV=production
DB_PASSWORD=actual-secure-password
JWT_SECRET=actual-jwt-secret
```

### ❌ DON'T: Hardcode Secrets

```go
// Bad
const dbPassword = "admin123"
const jwtSecret = "my-secret"

// Good
dbPassword := os.Getenv("DB_PASSWORD")
jwtSecret := os.Getenv("JWT_SECRET")
```

### ❌ DON'T: Log Sensitive Data

```go
// Bad
log.Printf("DB Password: %s", os.Getenv("DB_PASSWORD"))

// Good
log.Printf("DB connection established")
```

---

## Testing with Environment Variables

```go
func TestWithEnv(t *testing.T) {
    // Save original
    originalEnv := os.Getenv("APP_ENV")
    defer os.Setenv("APP_ENV", originalEnv)
    
    // Set test environment
    os.Setenv("APP_ENV", "test")
    os.Setenv("DB_DRIVER", "sqlite")
    os.Setenv("DB_NAME", ":memory:")
    
    // Run test
    app := core.NewApp()
    assert.NotNil(t, app)
}
```

---

## Reference Table

| Variable | Type | Default | Required | Description |
|----------|------|---------|----------|-------------|
| APP_ENV | string | development | Yes | Environment (development/staging/production) |
| APP_PORT | int | 8080 | No | HTTP server port |
| DEBUG | bool | false | No | Debug mode |
| DB_DRIVER | string | sqlite | Yes | Database driver |
| DB_HOST | string | localhost | Conditional | Database host |
| DB_PASSWORD | string | - | Conditional | Database password |
| LOG_LEVEL | string | info | No | Log level |
| JWT_SECRET | string | - | Production | JWT signing key |

---

## Next Steps

- [**Docker Deployment**](docker.md) - Use with containers
- [**Production Setup**](production-setup.md) - Production configuration
- [**Security**](../advanced/security.md) - Secure configuration
- [**Configuration Guide**](../getting-started/configuration.md) - Basic setup

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
