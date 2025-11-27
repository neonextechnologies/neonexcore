# Docker Deployment

Containerize Neonex Core applications with Docker.

---

## Dockerfile

```dockerfile
# Multi-stage build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o neonex ./cmd/neonex

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary
COPY --from=builder /app/neonex .

# Copy config
COPY .env.production .env

EXPOSE 3000

CMD ["./neonex"]
```

---

## Docker Compose

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - \"3000:3000\"
    environment:
      - ENV=production
      - DB_HOST=postgres
    depends_on:
      - postgres
    restart: unless-stopped

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: neonex
      POSTGRES_USER: neonex
      POSTGRES_PASSWORD: password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  postgres_data:
```

---

## Build and Run

```powershell
# Build image
docker build -t neonex:latest .

# Run container
docker run -p 3000:3000 neonex:latest

# Using docker-compose
docker-compose up -d
```

---

## Next Steps

- [**Production Setup**](production-setup.md)
- [**Environment Variables**](environment-variables.md)
- [**Monitoring**](monitoring.md)
