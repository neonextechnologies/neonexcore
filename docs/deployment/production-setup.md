# Production Setup

Guide to deploying Neonex Core applications to production.

---

## Build for Production

```powershell
# Build binary
go build -o neonex.exe ./cmd/neonex

# Linux/Mac
GOOS=linux GOARCH=amd64 go build -o neonex ./cmd/neonex
```

---

## Environment Configuration

```bash
# .env.production
ENV=production
PORT=3000

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=neonex_prod
DB_USER=neonex
DB_PASSWORD=your-secure-password

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=file
LOG_FILE_PATH=/var/log/neonex/app.log

# Security
JWT_SECRET=your-secret-key-here
```

---

## Systemd Service

```ini
# /etc/systemd/system/neonex.service
[Unit]
Description=Neonex Application
After=network.target

[Service]
Type=simple
User=neonex
WorkingDirectory=/opt/neonex
ExecStart=/opt/neonex/neonex
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

**Commands:**
```bash
sudo systemctl daemon-reload
sudo systemctl enable neonex
sudo systemctl start neonex
sudo systemctl status neonex
```

---

## Nginx Reverse Proxy

```nginx
# /etc/nginx/sites-available/neonex
server {
    listen 80;
    server_name example.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

---

## Process Manager (PM2 Alternative)

```bash
# Using supervisord
[program:neonex]
command=/opt/neonex/neonex
directory=/opt/neonex
autostart=true
autorestart=true
stderr_logfile=/var/log/neonex/err.log
stdout_logfile=/var/log/neonex/out.log
```

---

## Best Practices

### ✅ DO:
- Use environment variables for config
- Enable HTTPS
- Set up log rotation
- Monitor application health
- Use reverse proxy
- Implement graceful shutdown

### ❌ DON'T:
- Run as root user
- Hardcode secrets
- Ignore error logs
- Skip backup strategy

---

## Next Steps

- [**Docker**](docker.md)
- [**Environment Variables**](environment-variables.md)
- [**Monitoring**](monitoring.md)
