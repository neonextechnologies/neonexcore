# Frequently Asked Questions

Common questions about Neonex Core.

---

## General

### What is Neonex Core?

Neonex Core is a modular Go framework for building scalable web applications with clean architecture, dependency injection, and modern development practices.

### Why Neonex Core?

- **Modular** - Organize code by features
- **Type-safe** - Generic repository and DI
- **Modern** - Latest Go patterns
- **Fast** - Built on Fiber (fasthttp)
- **Developer-friendly** - Hot reload, CLI tools

---

## Installation

### How do I install Neonex Core?

```bash
go install github.com/neophp/neonexcore/cmd/neonex@latest
neonex new my-project
cd my-project
go mod download
```

### What are the requirements?

- Go 1.21 or higher
- 64-bit operating system
- Git (for version control)

---

## Development

### How do I create a new module?

```bash
neonex module create <module-name>
```

### How do I enable hot reload?

```bash
air
# or
neonex serve
```

### How do I run tests?

```bash
go test ./...
```

---

## Database

### What databases are supported?

- SQLite (default)
- PostgreSQL
- MySQL
- Turso (libSQL)

### How do I change database?

Update `.env`:
```bash
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=mydb
DB_USER=user
DB_PASSWORD=password
```

### How do migrations work?

Neonex uses GORM AutoMigrate:
```go
db.AutoMigrate(&User{}, &Product{})
```

---

## Deployment

### How do I build for production?

```bash
go build -o app ./cmd/neonex
```

### Can I use Docker?

Yes! See [Docker guide](../deployment/docker.md)

### What about Kubernetes?

Coming in v0.3 - see [Roadmap](roadmap.md)

---

## Troubleshooting

### Port already in use?

```powershell
# Windows
netstat -ano | findstr :3000
taskkill /PID <PID> /F

# Or change port
$env:PORT = \"3001\"
```

### Database connection failed?

Check:
1. Database is running
2. Credentials in `.env`
3. Network connectivity
4. Firewall settings

### Module not found?

```bash
go mod tidy
go mod download
```

---

## Contributing

### How can I contribute?

See [How to Contribute](../contributing/how-to-contribute.md)

### Where do I report bugs?

Create an issue on GitHub with:
- Description
- Steps to reproduce
- Expected vs actual behavior
- Environment details

---

## Next Steps

- [**Roadmap**](roadmap.md)
- [**Changelog**](changelog.md)
- [**Support**](support.md)
