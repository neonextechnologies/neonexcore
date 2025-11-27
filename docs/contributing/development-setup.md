# Development Setup

Set up your development environment for contributing to Neonex Core.

---

## Prerequisites

- **Go 1.21+**
- **Git**
- **Air** (for hot reload)
- **VS Code** (recommended)

---

## Clone Repository

```bash
git clone https://github.com/yourusername/neonexcore.git
cd neonexcore
```

---

## Install Dependencies

```powershell
# Download Go modules
go mod download

# Install Air
go install github.com/air-verse/air@latest

# Install tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

---

## Environment Setup

```bash
# Copy example env
cp .env.example .env

# Edit configuration
# Edit .env with your settings
```

---

## Database Setup

```powershell
# Run with SQLite (default)
# No setup needed

# Or PostgreSQL
createdb neonex_dev
```

---

## Run Development Server

```powershell
# With hot reload
air

# Or
neonex serve

# Or directly
go run ./cmd/neonex
```

---

## Run Tests

```powershell
# All tests
go test ./...

# With coverage
go test -cover ./...

# Verbose
go test -v ./...
```

---

## Code Quality

```powershell
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Fix issues
golangci-lint run --fix
```

---

## VS Code Extensions

Recommended extensions:
- **Go** (golang.go)
- **GitLens**
- **EditorConfig**
- **Better Comments**

---

## Next Steps

- [**How to Contribute**](how-to-contribute.md)
- [**Code of Conduct**](code-of-conduct.md)
- [**Hot Reload**](../development/hot-reload.md)
