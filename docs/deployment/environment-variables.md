# Environment Variables

Complete environment variable reference.

---

## Application

```bash
ENV=production              # development, staging, production
PORT=3000                   # Server port
HOST=0.0.0.0               # Server host
```

---

## Database

```bash
DB_DRIVER=postgres          # sqlite, postgres, mysql, turso
DB_HOST=localhost
DB_PORT=5432
DB_NAME=neonex
DB_USER=neonex
DB_PASSWORD=your-password
DB_SSL_MODE=disable         # disable, require
```

---

## Logging

```bash
LOG_LEVEL=info              # debug, info, warn, error, fatal
LOG_FORMAT=json             # json, text, console
LOG_OUTPUT=stdout           # stdout, stderr, file
LOG_FILE_PATH=logs/app.log
LOG_FILE_MAX_SIZE=100       # MB
LOG_FILE_MAX_BACKUPS=5
LOG_FILE_MAX_AGE=30         # days
LOG_FILE_COMPRESS=true
```

---

## Security

```bash
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h
API_KEY=your-api-key
ENCRYPTION_KEY=your-encryption-key
```

---

## CORS

```bash
CORS_ALLOW_ORIGINS=*
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOW_HEADERS=Content-Type,Authorization
```

---

## Loading .env Files

```go
import \"github.com/joho/godotenv\"

func init() {
    if err := godotenv.Load(); err != nil {
        log.Fatal(\"Error loading .env file\")
    }
}
```

---

## Next Steps

- [**Production Setup**](production-setup.md)
- [**Docker**](docker.md)
- [**Monitoring**](monitoring.md)
