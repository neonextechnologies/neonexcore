# Installation & Setup Guide

## Prerequisites

### 1. Install Go
Download and install Go 1.21 or later:
- **Windows**: https://go.dev/dl/go1.21.windows-amd64.msi
- **macOS**: https://go.dev/dl/go1.21.darwin-amd64.pkg
- **Linux**: https://go.dev/dl/go1.21.linux-amd64.tar.gz

Verify installation:
```bash
go version
```

### 2. Install Dependencies

```bash
cd /e/go/neonexcore
go mod download
go mod tidy
```

## Required External Services

### 1. Redis (for Caching)
**Windows:**
- Download: https://github.com/microsoftarchive/redis/releases
- Or use Docker:
```bash
docker run -d -p 6379:6379 redis:alpine
```

**Linux/macOS:**
```bash
# Ubuntu/Debian
sudo apt-get install redis-server

# macOS
brew install redis
```

### 2. PostgreSQL or MySQL (for Database)
**PostgreSQL:**
```bash
# Docker
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=password postgres:15

# Windows: Download from https://www.postgresql.org/download/windows/
# macOS: brew install postgresql
```

**MySQL:**
```bash
# Docker
docker run -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password mysql:8

# Windows: Download from https://dev.mysql.com/downloads/installer/
# macOS: brew install mysql
```

## Component-Specific Requirements

### Web3/Blockchain Components
No additional installation needed. Uses go-ethereum library.

**Note:** For production use, you'll need:
- Infura API Key: https://infura.io/
- Alchemy API Key: https://www.alchemy.com/
- Or your own Ethereum node

### AI/ML Components
No additional installation needed for the framework.

**Note:** For OpenAI integration:
- OpenAI API Key: https://platform.openai.com/api-keys

### Metrics Dashboard
Uses built-in metrics collector. For Prometheus integration:
```bash
# Docker
docker run -d -p 9090:9090 prom/prometheus
```

### gRPC Components
Protocol Buffers compiler (optional, for generating new .proto files):
```bash
# Install protoc
# Windows: Download from https://github.com/protocolbuffers/protobuf/releases
# macOS: brew install protobuf
# Linux: sudo apt-get install protobuf-compiler

# Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Quick Start

### 1. Clone Repository
```bash
git clone https://github.com/neonextechnologies/neonexcore.git
cd neonexcore
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Set Environment Variables
Create `.env` file:
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=neonexcore

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Server
PORT=3000

# Web3 (Optional)
INFURA_API_KEY=your_infura_key
ETHEREUM_RPC_URL=https://mainnet.infura.io/v3/YOUR_KEY

# OpenAI (Optional)
OPENAI_API_KEY=your_openai_key
```

### 4. Run Application
```bash
# Development mode with hot reload
go run main.go

# Or build and run
go build -o neonexcore
./neonexcore
```

### 5. Run Examples
```bash
# WebSocket example
go run examples/websocket_example.go

# GraphQL example
go run examples/graphql_example.go

# Cache example
go run examples/cache_example.go

# Metrics example
go run examples/metrics_example.go

# AI/ML example
go run examples/ai_example.go

# Workflow example
go run examples/workflow_example.go

# Web3 example
go run examples/web3_example.go
```

## Build for Production

### Build Binary
```bash
# Current platform
go build -o neonexcore

# Windows
GOOS=windows GOARCH=amd64 go build -o neonexcore.exe

# Linux
GOOS=linux GOARCH=amd64 go build -o neonexcore

# macOS
GOOS=darwin GOARCH=amd64 go build -o neonexcore
```

### Docker Deployment
```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o neonexcore

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/neonexcore .
COPY --from=builder /app/.env .
EXPOSE 3000
CMD ["./neonexcore"]
```

Build and run:
```bash
docker build -t neonexcore .
docker run -p 3000:3000 neonexcore
```

## Troubleshooting

### Go Not Found
Add Go to PATH:
```bash
# Windows (PowerShell as Admin)
$env:Path += ";C:\Go\bin"
[Environment]::SetEnvironmentVariable("Path", $env:Path, [System.EnvironmentVariableTarget]::Machine)

# Linux/macOS
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
```

### Module Dependencies Error
```bash
go clean -modcache
go mod download
go mod tidy
```

### Redis Connection Error
Check if Redis is running:
```bash
redis-cli ping
# Should return: PONG
```

### Database Connection Error
Verify database is running and credentials are correct:
```bash
# PostgreSQL
psql -h localhost -U postgres -d neonexcore

# MySQL
mysql -h localhost -u root -p neonexcore
```

### Port Already in Use
Change port in `.env` or kill process:
```bash
# Windows
netstat -ano | findstr :3000
taskkill /PID <PID> /F

# Linux/macOS
lsof -ti:3000 | xargs kill -9
```

## Development Tools

### Recommended VS Code Extensions
- Go (golang.go)
- Docker (ms-azuretools.vscode-docker)
- REST Client (humao.rest-client)
- Thunder Client (rangav.vscode-thunder-client)

### Useful Commands
```bash
# Format code
go fmt ./...

# Run tests
go test ./...

# Check for issues
go vet ./...

# Generate documentation
godoc -http=:6060

# View dependencies
go mod graph

# Update dependencies
go get -u ./...
```

## Next Steps

1. ✅ Install Go and required services
2. ✅ Download dependencies: `go mod download`
3. ✅ Configure environment variables
4. ✅ Run examples to test components
5. ✅ Start building your application!

For more information, see component-specific README files in `pkg/` directories.
