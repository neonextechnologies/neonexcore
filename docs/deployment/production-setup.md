# Production Setup

Deploy Neonex Core applications to production environments.

---

## Pre-Deployment Checklist

### Code Quality
- [ ] All tests passing
- [ ] Code coverage > 80%
- [ ] No security vulnerabilities
- [ ] Dependencies updated
- [ ] Linter warnings resolved

### Configuration
- [ ] Environment variables configured
- [ ] Database migrations tested
- [ ] Secrets managed securely
- [ ] HTTPS certificates ready
- [ ] CORS configured properly

### Performance
- [ ] Database indexed
- [ ] Connection pooling configured
- [ ] Caching strategy implemented
- [ ] Rate limiting enabled
- [ ] Load testing completed

### Monitoring
- [ ] Logging configured
- [ ] Error tracking setup
- [ ] Performance monitoring ready
- [ ] Health checks implemented
- [ ] Alerts configured

---

## Build for Production

### Compile Binary

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o bin/app-linux-amd64 main.go

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o bin/app-windows-amd64.exe main.go

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o bin/app-darwin-arm64 main.go

# With optimizations
go build -ldflags="-s -w" -o bin/app main.go
```

### Build Script

```bash
#!/bin/bash
# build.sh

VERSION=$(git describe --tags --always)
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

go build \
  -ldflags="-s -w -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME" \
  -o bin/app \
  main.go

echo "Built version $VERSION at $BUILD_TIME"
```

---

## Environment Configuration

### Production .env

```bash
# Application
APP_NAME=neonex-production
APP_ENV=production
APP_PORT=8080
DEBUG=false

# Database
DB_DRIVER=postgres
DB_HOST=prod-db.example.com
DB_PORT=5432
DB_NAME=neonex_prod
DB_USER=neonex
DB_PASSWORD=${DB_PASSWORD}  # From secrets
DB_SSL_MODE=require
DB_MAX_CONNECTIONS=25
DB_MAX_IDLE=5

# Logging
LOG_LEVEL=warn
LOG_FORMAT=json
LOG_FILE=/var/log/neonex/app.log
LOG_MAX_SIZE=100
LOG_MAX_AGE=30
LOG_MAX_BACKUPS=10

# Security
JWT_SECRET=${JWT_SECRET}
ENCRYPTION_KEY=${ENCRYPTION_KEY}

# CORS
CORS_ALLOW_ORIGINS=https://app.example.com
CORS_ALLOW_CREDENTIALS=true

# Rate Limiting
RATE_LIMIT_MAX=100
RATE_LIMIT_WINDOW=60s
```

---

## Systemd Service

### Create Service File

```ini
# /etc/systemd/system/neonex.service
[Unit]
Description=Neonex Core Application
After=network.target postgresql.service

[Service]
Type=simple
User=neonex
Group=neonex
WorkingDirectory=/opt/neonex
Environment="APP_ENV=production"
EnvironmentFile=/opt/neonex/.env
ExecStart=/opt/neonex/bin/app
Restart=always
RestartSec=10
StandardOutput=append:/var/log/neonex/app.log
StandardError=append:/var/log/neonex/error.log

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/log/neonex /opt/neonex/data

[Install]
WantedBy=multi-user.target
```

### Manage Service

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service
sudo systemctl enable neonex

# Start service
sudo systemctl start neonex

# Check status
sudo systemctl status neonex

# View logs
sudo journalctl -u neonex -f

# Restart service
sudo systemctl restart neonex

# Stop service
sudo systemctl stop neonex
```

---

## Nginx Reverse Proxy

### Configuration

```nginx
# /etc/nginx/sites-available/neonex
upstream neonex {
    server 127.0.0.1:8080;
    keepalive 32;
}

server {
    listen 80;
    server_name api.example.com;
    
    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.example.com;
    
    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/api.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.example.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    
    # Security Headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "DENY" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    
    # Logging
    access_log /var/log/nginx/neonex-access.log;
    error_log /var/log/nginx/neonex-error.log;
    
    # Gzip Compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml;
    
    # Rate Limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req zone=api burst=20 nodelay;
    
    location / {
        proxy_pass http://neonex;
        proxy_http_version 1.1;
        
        # Headers
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Connection "";
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
        
        # Buffering
        proxy_buffering off;
    }
    
    # Health check endpoint
    location /health {
        proxy_pass http://neonex;
        access_log off;
    }
}
```

### Enable Site

```bash
# Create symlink
sudo ln -s /etc/nginx/sites-available/neonex /etc/nginx/sites-enabled/

# Test configuration
sudo nginx -t

# Reload Nginx
sudo systemctl reload nginx
```

---

## SSL/TLS Certificates

### Let's Encrypt (Certbot)

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx

# Obtain certificate
sudo certbot --nginx -d api.example.com

# Auto-renewal (cron)
sudo crontab -e
0 0 * * * certbot renew --quiet
```

---

## Database Setup

### PostgreSQL Production

```bash
# Create database and user
sudo -u postgres psql

CREATE DATABASE neonex_prod;
CREATE USER neonex WITH ENCRYPTED PASSWORD 'secure-password';
GRANT ALL PRIVILEGES ON DATABASE neonex_prod TO neonex;
\q

# Configure PostgreSQL
sudo nano /etc/postgresql/14/main/postgresql.conf

# Recommended settings
max_connections = 100
shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 64MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
effective_io_concurrency = 200

# Restart PostgreSQL
sudo systemctl restart postgresql
```

---

## Deployment Strategies

### Blue-Green Deployment

```bash
#!/bin/bash
# deploy.sh

# Build new version
go build -o bin/app-new main.go

# Health check function
health_check() {
    curl -f http://localhost:$1/health || return 1
}

# Start new instance on different port
PORT=8081 ./bin/app-new &
NEW_PID=$!

# Wait for new instance to be ready
sleep 5
if health_check 8081; then
    echo "New instance healthy"
    
    # Update Nginx to point to new instance
    sudo sed -i 's/8080/8081/' /etc/nginx/sites-available/neonex
    sudo nginx -s reload
    
    # Stop old instance
    kill $OLD_PID
    
    # Move new binary to production
    mv bin/app-new bin/app
else
    echo "New instance failed health check"
    kill $NEW_PID
    exit 1
fi
```

### Rolling Deployment

```bash
#!/bin/bash
# rolling-deploy.sh

INSTANCES=(8080 8081 8082)

for PORT in "${INSTANCES[@]}"; do
    echo "Deploying to instance on port $PORT"
    
    # Remove from load balancer
    # (Nginx upstream configuration)
    
    # Stop instance
    kill $(lsof -ti:$PORT)
    
    # Start new version
    PORT=$PORT ./bin/app &
    
    # Wait for health check
    sleep 5
    curl -f http://localhost:$PORT/health
    
    # Add back to load balancer
    
    echo "Instance $PORT deployed successfully"
done
```

---

## Monitoring & Health Checks

### Health Check Endpoint

```go
func HealthCheck(c *fiber.Ctx) error {
    // Check database
    if err := db.DB().Ping(); err != nil {
        return c.Status(503).JSON(fiber.Map{
            "status": "unhealthy",
            "database": "down",
        })
    }
    
    return c.JSON(fiber.Map{
        "status": "healthy",
        "database": "up",
        "uptime": time.Since(startTime).String(),
        "version": version,
    })
}

app.Get("/health", HealthCheck)
```

### Readiness Check

```go
func ReadinessCheck(c *fiber.Ctx) error {
    // Check if app is ready to serve traffic
    if !app.Ready {
        return c.Status(503).JSON(fiber.Map{
            "ready": false,
        })
    }
    
    return c.JSON(fiber.Map{
        "ready": true,
    })
}

app.Get("/ready", ReadinessCheck)
```

---

## Logging

### Production Logging

```go
import "go.uber.org/zap"

func setupLogger() *zap.Logger {
    config := zap.NewProductionConfig()
    config.OutputPaths = []string{
        "stdout",
        "/var/log/neonex/app.log",
    }
    config.ErrorOutputPaths = []string{
        "stderr",
        "/var/log/neonex/error.log",
    }
    
    logger, _ := config.Build()
    return logger
}
```

### Log Rotation

```yaml
# /etc/logrotate.d/neonex
/var/log/neonex/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 0640 neonex neonex
    sharedscripts
    postrotate
        systemctl reload neonex
    endscript
}
```

---

## Backup Strategy

### Database Backups

```bash
#!/bin/bash
# backup.sh

BACKUP_DIR="/var/backups/neonex"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/neonex_$DATE.sql"

# Create backup
pg_dump -U neonex -h localhost neonex_prod > $BACKUP_FILE

# Compress
gzip $BACKUP_FILE

# Upload to S3 (optional)
aws s3 cp $BACKUP_FILE.gz s3://my-backups/neonex/

# Keep only last 30 days
find $BACKUP_DIR -name "*.sql.gz" -mtime +30 -delete

echo "Backup completed: $BACKUP_FILE.gz"
```

### Automated Backups

```bash
# Cron job (daily at 2 AM)
0 2 * * * /opt/neonex/scripts/backup.sh
```

---

## Performance Tuning

### System Limits

```bash
# /etc/security/limits.conf
neonex soft nofile 65536
neonex hard nofile 65536
neonex soft nproc 4096
neonex hard nproc 4096
```

### Kernel Parameters

```bash
# /etc/sysctl.conf
net.core.somaxconn = 65535
net.ipv4.tcp_max_syn_backlog = 8192
net.ipv4.ip_local_port_range = 1024 65535
net.ipv4.tcp_tw_reuse = 1
net.ipv4.tcp_fin_timeout = 15
```

---

## Troubleshooting

### Check Application Logs

```bash
# Systemd logs
sudo journalctl -u neonex -f

# Application logs
tail -f /var/log/neonex/app.log

# Error logs
tail -f /var/log/neonex/error.log
```

### Check Resource Usage

```bash
# CPU and Memory
top -p $(pgrep app)

# Network connections
netstat -an | grep :8080

# Open files
lsof -p $(pgrep app)
```

### Database Issues

```bash
# Check connections
psql -U neonex -c "SELECT count(*) FROM pg_stat_activity WHERE datname='neonex_prod';"

# Check slow queries
psql -U neonex -c "SELECT pid, query, state FROM pg_stat_activity WHERE state != 'idle';"
```

---

## Next Steps

- [**Docker Deployment**](docker.md) - Containerize your app
- [**Environment Variables**](environment-variables.md) - Manage configuration
- [**Monitoring**](monitoring.md) - Set up monitoring
- [**Security**](../advanced/security.md) - Production security

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
