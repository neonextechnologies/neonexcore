# Development Setup

Complete guide to setting up your development environment for contributing to Neonex Core.

---

## Prerequisites

### Required Software

| Software | Minimum Version | Purpose |
|----------|----------------|---------|
| **Go** | 1.21+ | Runtime and compiler |
| **Git** | 2.30+ | Version control |
| **Make** | 4.0+ | Build automation |
| **Docker** | 20.10+ | Container testing |
| **PostgreSQL** | 14+ | Database (optional) |
| **MySQL** | 8.0+ | Database (optional) |

### Optional Tools

| Tool | Purpose |
|------|---------|
| **VS Code** | Recommended IDE |
| **Delve** | Go debugger |
| **golangci-lint** | Code linter |
| **Air** | Hot reload |
| **gosec** | Security scanner |

---

## Installation

### 1. Install Go

#### Windows

```powershell
# Using Chocolatey
choco install golang

# Verify installation
go version
```

#### macOS

```bash
# Using Homebrew
brew install go

# Verify installation
go version
```

#### Linux

```bash
# Download and install
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# Add to PATH (~/.bashrc or ~/.zshrc)
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# Reload shell
source ~/.bashrc

# Verify installation
go version
```

### 2. Install Git

#### Windows

```powershell
# Using Chocolatey
choco install git

# Or download from https://git-scm.com/
```

#### macOS

```bash
# Using Homebrew
brew install git
```

#### Linux

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install git

# Fedora/RHEL
sudo dnf install git
```

### 3. Install Development Tools

```bash
# Go development tools
go install github.com/go-delve/delve/cmd/dlv@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install github.com/cosmtrek/air@latest

# Verify installations
dlv version
golangci-lint version
gosec --version
air -v
```

### 4. Install Docker (Optional)

#### Windows

Download and install [Docker Desktop](https://www.docker.com/products/docker-desktop/)

#### macOS

```bash
# Using Homebrew
brew install --cask docker
```

#### Linux

```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Add user to docker group
sudo usermod -aG docker $USER
newgrp docker

# Verify
docker --version
docker-compose --version
```

---

## Project Setup

### 1. Fork and Clone

```bash
# Fork repository on GitHub
# https://github.com/neonexcore/neonexcore

# Clone your fork
git clone https://github.com/YOUR_USERNAME/neonexcore.git
cd neonexcore

# Add upstream remote
git remote add upstream https://github.com/neonexcore/neonexcore.git

# Verify remotes
git remote -v
# origin    https://github.com/YOUR_USERNAME/neonexcore.git (fetch)
# origin    https://github.com/YOUR_USERNAME/neonexcore.git (push)
# upstream  https://github.com/neonexcore/neonexcore.git (fetch)
# upstream  https://github.com/neonexcore/neonexcore.git (push)
```

### 2. Install Dependencies

```bash
# Download Go modules
go mod download

# Verify dependencies
go mod verify

# Tidy up (optional)
go mod tidy
```

### 3. Configure Environment

```bash
# Copy example environment file
cp .env.example .env

# Edit with your settings
# Linux/macOS
nano .env

# Windows
notepad .env
```

**Example .env:**

```bash
# Application
APP_ENV=development
APP_PORT=8080
DEBUG=true

# Database
DB_DRIVER=sqlite
DB_NAME=dev.db

# Logging
LOG_LEVEL=debug
LOG_FORMAT=text
```

---

## Database Setup

### SQLite (Default)

No setup required - file-based database.

```bash
# Database file will be created automatically
# Default: dev.db
```

### PostgreSQL

```bash
# Start PostgreSQL with Docker
docker run -d \
  --name postgres-dev \
  -e POSTGRES_USER=neonex \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=neonex_dev \
  -p 5432:5432 \
  postgres:15-alpine

# Or install locally
# Ubuntu/Debian
sudo apt-get install postgresql postgresql-contrib

# macOS
brew install postgresql
brew services start postgresql

# Create database
createdb neonex_dev
```

**Update .env:**

```bash
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=neonex_dev
DB_USER=neonex
DB_PASSWORD=password
DB_SSL_MODE=disable
```

### MySQL

```bash
# Start MySQL with Docker
docker run -d \
  --name mysql-dev \
  -e MYSQL_ROOT_PASSWORD=rootpass \
  -e MYSQL_DATABASE=neonex_dev \
  -e MYSQL_USER=neonex \
  -e MYSQL_PASSWORD=password \
  -p 3306:3306 \
  mysql:8

# Or install locally
# Ubuntu/Debian
sudo apt-get install mysql-server

# macOS
brew install mysql
brew services start mysql

# Create database
mysql -u root -p
CREATE DATABASE neonex_dev;
CREATE USER 'neonex'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON neonex_dev.* TO 'neonex'@'localhost';
FLUSH PRIVILEGES;
```

**Update .env:**

```bash
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_NAME=neonex_dev
DB_USER=neonex
DB_PASSWORD=password
```

---

## IDE Setup

### VS Code

#### Install Extensions

```bash
# Essential extensions
code --install-extension golang.go
code --install-extension ms-vscode.makefile-tools

# Recommended extensions
code --install-extension eamodio.gitlens
code --install-extension streetsidesoftware.code-spell-checker
code --install-extension ms-azuretools.vscode-docker
```

#### Configure Settings

Create `.vscode/settings.json`:

```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "workspace",
  "go.formatTool": "gofmt",
  "go.formatOnSave": true,
  "go.testOnSave": true,
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  },
  "go.testFlags": ["-v"],
  "go.coverOnSave": true,
  "go.coverageDecorator": {
    "type": "gutter",
    "coveredHighlightColor": "rgba(64,128,64,0.5)",
    "uncoveredHighlightColor": "rgba(128,64,64,0.5)"
  }
}
```

#### Configure Debugging

Create `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Server",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/main.go",
      "env": {
        "APP_ENV": "development"
      },
      "args": []
    },
    {
      "name": "Test Current File",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${file}"
    },
    {
      "name": "Test Current Package",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${fileDirname}"
    }
  ]
}
```

#### Configure Tasks

Create `.vscode/tasks.json`:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Run Server",
      "type": "shell",
      "command": "go run main.go",
      "group": {
        "kind": "build",
        "isDefault": true
      },
      "presentation": {
        "reveal": "always",
        "panel": "new"
      }
    },
    {
      "label": "Run Tests",
      "type": "shell",
      "command": "go test ./...",
      "group": "test"
    },
    {
      "label": "Run Linter",
      "type": "shell",
      "command": "golangci-lint run",
      "group": "test"
    }
  ]
}
```

### GoLand / IntelliJ IDEA

1. **Open Project**: File → Open → Select neonexcore directory
2. **Configure Go SDK**: Settings → Languages & Frameworks → Go → Add SDK
3. **Enable Go Modules**: Settings → Go → Go Modules → Enable
4. **Configure Run Configuration**: 
   - Run → Edit Configurations → Add New → Go Build
   - Files: main.go
   - Working directory: Project root

---

## Running the Application

### Standard Run

```bash
# Run directly
go run main.go

# With environment variable
APP_ENV=development go run main.go

# Build and run
go build -o neonex main.go
./neonex
```

### Hot Reload Development

```bash
# Install Air (if not already)
go install github.com/cosmtrek/air@latest

# Run with hot reload
air

# Or use neonex CLI
neonex serve
```

### Docker Development

```bash
# Build image
docker build -t neonex-dev .

# Run container
docker run -p 8080:8080 \
  -v $(pwd):/app \
  -e APP_ENV=development \
  neonex-dev

# Or use Docker Compose
docker-compose up
```

---

## Testing

### Run All Tests

```bash
# All tests
go test ./...

# With verbose output
go test -v ./...

# With coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Specific Tests

```bash
# Test specific package
go test ./modules/user/...

# Test specific function
go test -run TestUserRepository_Create ./modules/user/

# Test with race detector
go test -race ./...
```

### Benchmark Tests

```bash
# Run benchmarks
go test -bench=. ./...

# With memory stats
go test -bench=. -benchmem ./...

# Specific benchmark
go test -bench=BenchmarkUserCreate ./modules/user/
```

---

## Code Quality

### Format Code

```bash
# Format all Go files
go fmt ./...

# Or use gofmt directly
gofmt -w .
```

### Run Linter

```bash
# Run golangci-lint
golangci-lint run

# With specific linters
golangci-lint run --enable=gosec,gofmt

# Fix auto-fixable issues
golangci-lint run --fix
```

### Security Scan

```bash
# Run gosec
gosec ./...

# With JSON output
gosec -fmt=json -out=results.json ./...
```

### Dependency Check

```bash
# List outdated dependencies
go list -u -m all

# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify

# Vendor dependencies (optional)
go mod vendor
```

---

## Troubleshooting

### Common Issues

#### Issue: Module Download Fails

```bash
# Set Go proxy
go env -w GOPROXY=https://proxy.golang.org,direct

# Or use module mirror
go env -w GOPROXY=https://goproxy.io,direct
```

#### Issue: Port Already in Use

```bash
# Find process using port
# Linux/macOS
lsof -i :8080

# Windows
netstat -ano | findstr :8080

# Kill process
# Linux/macOS
kill -9 <PID>

# Windows
taskkill /PID <PID> /F

# Or change port in .env
APP_PORT=8081
```

#### Issue: Database Connection Failed

```bash
# Check database is running
# PostgreSQL
pg_isready

# MySQL
mysqladmin ping

# Check credentials in .env
DB_USER=correct_username
DB_PASSWORD=correct_password
```

#### Issue: Tests Failing

```bash
# Clean test cache
go clean -testcache

# Run tests with verbose output
go test -v ./...

# Check for race conditions
go test -race ./...
```

---

## Development Workflow

### Daily Workflow

```bash
# 1. Sync with upstream
git fetch upstream
git checkout main
git merge upstream/main

# 2. Create feature branch
git checkout -b feature/my-feature

# 3. Make changes
# ... edit files ...

# 4. Run tests
go test ./...

# 5. Format and lint
go fmt ./...
golangci-lint run

# 6. Commit changes
git add .
git commit -m "feat: add my feature"

# 7. Push to your fork
git push origin feature/my-feature

# 8. Create pull request on GitHub
```

### Keeping Fork Updated

```bash
# Fetch upstream changes
git fetch upstream

# Merge into main
git checkout main
git merge upstream/main

# Push to your fork
git push origin main

# Rebase feature branch (if needed)
git checkout feature/my-feature
git rebase main
```

---

## Quick Reference

### Useful Commands

```bash
# Development
go run main.go              # Run application
air                         # Hot reload
go test ./...              # Run tests
go fmt ./...               # Format code
golangci-lint run          # Lint code

# Building
go build                   # Build binary
go build -o neonex         # Build with custom name
go build -ldflags="-s -w"  # Build optimized

# Modules
go mod download            # Download dependencies
go mod tidy               # Clean dependencies
go mod verify             # Verify dependencies

# Testing
go test -v ./...          # Verbose tests
go test -cover ./...      # With coverage
go test -race ./...       # Race detection
go test -bench=. ./...    # Benchmarks

# Tools
dlv debug                 # Debug with Delve
gosec ./...              # Security scan
go tool pprof            # Profiling
```

---

## Next Steps

- [**How to Contribute**](how-to-contribute.md) - Contribution guidelines
- [**Code of Conduct**](code-of-conduct.md) - Community standards
- [**Architecture**](../introduction/architecture.md) - System design
- [**Testing Guide**](../development/testing.md) - Writing tests

---

**Need help?** Ask in [GitHub Discussions](https://github.com/neonexcore/neonexcore/discussions) or [Discord](https://discord.gg/neonexcore).
