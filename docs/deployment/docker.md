# Docker Deployment

Deploy Neonex Core applications using Docker containers.

---

## Overview

Docker benefits:
- **Consistency** - Same environment everywhere
- **Isolation** - No dependency conflicts
- **Scalability** - Easy horizontal scaling
- **Portability** - Run anywhere Docker runs
- **Reproducibility** - Identical builds every time

---

## Quick Start

### Dockerfile

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-s -w' -o main .

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Copy .env file (optional, better use environment variables)
COPY .env.example .env

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run application
CMD ["./main"]
```

### Build Image

```bash
# Build image
docker build -t neonex-app:latest .

# Build with specific version
docker build -t neonex-app:v1.0.0 .

# Build with build args
docker build \
  --build-arg VERSION=v1.0.0 \
  --build-arg BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  -t neonex-app:v1.0.0 .
```

### Run Container

```bash
# Run container
docker run -d \
  --name neonex \
  -p 8080:8080 \
  --env-file .env \
  neonex-app:latest

# Run with environment variables
docker run -d \
  --name neonex \
  -p 8080:8080 \
  -e APP_ENV=production \
  -e DB_HOST=postgres \
  -e DB_PASSWORD=secret \
  neonex-app:latest

# Run with volume mounts
docker run -d \
  --name neonex \
  -p 8080:8080 \
  -v $(pwd)/logs:/root/logs \
  -v $(pwd)/data:/root/data \
  neonex-app:latest
```

---

## Docker Compose

### Basic Setup

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: .
    container_name: neonex-app
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_NAME=neonex
      - DB_USER=neonex
      - DB_PASSWORD=${DB_PASSWORD}
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped
    networks:
      - neonex-network
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/health"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 40s

  postgres:
    image: postgres:15-alpine
    container_name: neonex-postgres
    environment:
      - POSTGRES_DB=neonex
      - POSTGRES_USER=neonex
      - POSTGRES_PASSWORD=${DB_PASSWORD}
    volumes:
      - postgres-data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    restart: unless-stopped
    networks:
      - neonex-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U neonex"]
      interval: 10s
      timeout: 5s
      retries: 5

networks:
  neonex-network:
    driver: bridge

volumes:
  postgres-data:
    driver: local
```

### Run with Docker Compose

```bash
# Start services
docker-compose up -d

# View logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f app

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v

# Rebuild and restart
docker-compose up -d --build

# Scale application
docker-compose up -d --scale app=3
```

---

## Multi-Stage Build (Optimized)

### Advanced Dockerfile

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev

WORKDIR /build

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build with optimizations
ARG VERSION=dev
ARG BUILD_TIME
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" \
    -a -installsuffix cgo \
    -o /build/app \
    .

# Runtime stage
FROM alpine:3.18

# Install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    curl \
    && addgroup -g 1000 neonex \
    && adduser -D -u 1000 -G neonex neonex

# Set timezone
ENV TZ=UTC

WORKDIR /app

# Copy binary from builder
COPY --from=builder --chown=neonex:neonex /build/app .

# Create required directories
RUN mkdir -p /app/logs /app/data \
    && chown -R neonex:neonex /app

# Switch to non-root user
USER neonex

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Run application
ENTRYPOINT ["./app"]
```

---

## Production Docker Compose

### Complete Stack

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  # Application
  app:
    build:
      context: .
      dockerfile: Dockerfile
      args:
        VERSION: ${VERSION:-latest}
        BUILD_TIME: ${BUILD_TIME}
    image: neonex-app:${VERSION:-latest}
    container_name: neonex-app
    ports:
      - "8080:8080"
    environment:
      - APP_NAME=neonex-production
      - APP_ENV=production
      - APP_PORT=8080
      - DEBUG=false
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_NAME=${DB_NAME}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_SSL_MODE=require
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - LOG_LEVEL=warn
      - LOG_FORMAT=json
    volumes:
      - app-logs:/app/logs
      - app-data:/app/data
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    restart: unless-stopped
    networks:
      - backend
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 40s

  # PostgreSQL Database
  postgres:
    image: postgres:15-alpine
    container_name: neonex-postgres
    environment:
      - POSTGRES_DB=${DB_NAME}
      - POSTGRES_USER=${DB_USER}
      - POSTGRES_PASSWORD=${DB_PASSWORD}
      - POSTGRES_INITDB_ARGS=--encoding=UTF-8
    volumes:
      - postgres-data:/var/lib/postgresql/data
      - ./init-db:/docker-entrypoint-initdb.d
    ports:
      - "5432:5432"
    restart: unless-stopped
    networks:
      - backend
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 1G
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis Cache
  redis:
    image: redis:7-alpine
    container_name: neonex-redis
    command: redis-server --appendonly yes --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis-data:/data
    ports:
      - "6379:6379"
    restart: unless-stopped
    networks:
      - backend
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 3s
      retries: 3

  # Nginx Reverse Proxy
  nginx:
    image: nginx:alpine
    container_name: neonex-nginx
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/conf.d:/etc/nginx/conf.d:ro
      - ./certs:/etc/nginx/certs:ro
      - nginx-logs:/var/log/nginx
    depends_on:
      - app
    restart: unless-stopped
    networks:
      - backend
      - frontend
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 256M

  # Prometheus Monitoring
  prometheus:
    image: prom/prometheus:latest
    container_name: neonex-prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--storage.tsdb.retention.time=30d'
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - prometheus-data:/prometheus
    ports:
      - "9090:9090"
    restart: unless-stopped
    networks:
      - backend

  # Grafana Dashboards
  grafana:
    image: grafana/grafana:latest
    container_name: neonex-grafana
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=${GRAFANA_PASSWORD}
      - GF_SERVER_ROOT_URL=http://grafana.example.com
    volumes:
      - grafana-data:/var/lib/grafana
      - ./monitoring/grafana/dashboards:/etc/grafana/provisioning/dashboards:ro
      - ./monitoring/grafana/datasources:/etc/grafana/provisioning/datasources:ro
    ports:
      - "3000:3000"
    depends_on:
      - prometheus
    restart: unless-stopped
    networks:
      - backend

networks:
  backend:
    driver: bridge
  frontend:
    driver: bridge

volumes:
  app-logs:
  app-data:
  postgres-data:
  redis-data:
  nginx-logs:
  prometheus-data:
  grafana-data:
```

---

## Environment Variables

### .env File

```bash
# .env
VERSION=v1.0.0
BUILD_TIME=2024-01-15T10:00:00Z

# Database
DB_NAME=neonex_prod
DB_USER=neonex
DB_PASSWORD=secure-password-here

# Redis
REDIS_PASSWORD=redis-password-here

# Grafana
GRAFANA_PASSWORD=grafana-password-here
```

---

## Docker Commands

### Container Management

```bash
# List running containers
docker ps

# List all containers
docker ps -a

# View container logs
docker logs neonex-app

# Follow logs
docker logs -f neonex-app

# Execute command in container
docker exec -it neonex-app sh

# Inspect container
docker inspect neonex-app

# View container stats
docker stats neonex-app

# Stop container
docker stop neonex-app

# Start container
docker start neonex-app

# Restart container
docker restart neonex-app

# Remove container
docker rm neonex-app

# Force remove running container
docker rm -f neonex-app
```

### Image Management

```bash
# List images
docker images

# Remove image
docker rmi neonex-app:latest

# Remove unused images
docker image prune

# Remove all unused images
docker image prune -a

# Tag image
docker tag neonex-app:latest registry.example.com/neonex-app:v1.0.0

# Push image
docker push registry.example.com/neonex-app:v1.0.0

# Pull image
docker pull registry.example.com/neonex-app:v1.0.0
```

### Volume Management

```bash
# List volumes
docker volume ls

# Inspect volume
docker volume inspect postgres-data

# Remove volume
docker volume rm postgres-data

# Remove unused volumes
docker volume prune

# Backup volume
docker run --rm \
  -v postgres-data:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/postgres-backup.tar.gz /data
```

---

## Docker Registry

### Private Registry

```bash
# Run local registry
docker run -d -p 5000:5000 --name registry registry:2

# Tag image for registry
docker tag neonex-app:latest localhost:5000/neonex-app:latest

# Push to registry
docker push localhost:5000/neonex-app:latest

# Pull from registry
docker pull localhost:5000/neonex-app:latest
```

### Docker Hub

```bash
# Login to Docker Hub
docker login

# Tag image
docker tag neonex-app:latest username/neonex-app:latest

# Push to Docker Hub
docker push username/neonex-app:latest

# Pull from Docker Hub
docker pull username/neonex-app:latest
```

---

## CI/CD Integration

### GitHub Actions

```yaml
# .github/workflows/docker-build.yml
name: Docker Build and Push

on:
  push:
    branches: [ main ]
    tags: [ 'v*' ]

jobs:
  build:
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      
      - name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}
      
      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v4
        with:
          images: username/neonex-app
          tags: |
            type=ref,event=branch
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}
      
      - name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=registry,ref=username/neonex-app:buildcache
          cache-to: type=registry,ref=username/neonex-app:buildcache,mode=max
```

---

## Docker Swarm

### Initialize Swarm

```bash
# Initialize swarm
docker swarm init

# Deploy stack
docker stack deploy -c docker-compose.prod.yml neonex

# List services
docker service ls

# Scale service
docker service scale neonex_app=3

# View service logs
docker service logs -f neonex_app

# Remove stack
docker stack rm neonex
```

---

## Kubernetes

### Deployment

```yaml
# k8s/deployment.yml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: neonex-app
  labels:
    app: neonex
spec:
  replicas: 3
  selector:
    matchLabels:
      app: neonex
  template:
    metadata:
      labels:
        app: neonex
    spec:
      containers:
      - name: neonex
        image: neonex-app:latest
        ports:
        - containerPort: 8080
        env:
        - name: APP_ENV
          value: "production"
        - name: DB_HOST
          value: "postgres-service"
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: password
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "1Gi"
            cpu: "1000m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: neonex-service
spec:
  selector:
    app: neonex
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
```

---

## Best Practices

### ✅ DO: Use Multi-Stage Builds

```dockerfile
# Smaller final image
FROM golang:1.21-alpine AS builder
# Build here

FROM alpine:latest
COPY --from=builder /app/main .
```

### ✅ DO: Use Non-Root User

```dockerfile
RUN adduser -D -u 1000 appuser
USER appuser
```

### ✅ DO: Add Health Checks

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s \
  CMD curl -f http://localhost:8080/health || exit 1
```

### ❌ DON'T: Store Secrets in Images

```dockerfile
# Bad
ENV DB_PASSWORD=secret123

# Good
# Use environment variables at runtime
```

### ❌ DON'T: Run as Root

```dockerfile
# Bad
USER root

# Good
USER appuser
```

---

## Troubleshooting

### Check Container Logs

```bash
docker logs neonex-app --tail 100
docker logs -f neonex-app
```

### Enter Container Shell

```bash
docker exec -it neonex-app sh
```

### Check Container Resources

```bash
docker stats neonex-app
```

### Network Issues

```bash
# List networks
docker network ls

# Inspect network
docker network inspect neonex-network

# Test connectivity
docker exec neonex-app ping postgres
```

---

## Next Steps

- [**Production Setup**](production-setup.md) - Non-Docker deployment
- [**Environment Variables**](environment-variables.md) - Configuration management
- [**Monitoring**](monitoring.md) - Container monitoring
- [**Security**](../advanced/security.md) - Container security

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
