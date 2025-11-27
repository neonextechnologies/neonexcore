# neonex serve

Start the Neonex Core development server with optional hot reload support.

---

## Synopsis

```bash
neonex serve [flags]
```

Starts your application in development mode with:
- 🔥 Hot reload (optional) - Auto-rebuild on file changes
- 🚀 Fast startup - Quick development feedback
- 📝 Request logging - See all HTTP requests
- 🛡️ Graceful shutdown - Clean process termination
- ⚡ Auto-migration - Database schema sync

---

## Usage

### Basic Command

```bash
# Start without hot reload
neonex serve
```

**Output:**
```
Starting Neonex Core...
✓ Connected to database (sqlite)
✓ Auto-migrated 3 models
✓ Registered 2 modules
✓ Started HTTP server on :8080

Press Ctrl+C to stop
```

### With Hot Reload (Recommended)

```bash
# Start with Air hot reload
neonex serve --hot
```

**Output:**
```
Air is running... Press Ctrl+C to stop

  __    _   ___
 / /\  | | | |_)
/_/--\ |_| |_| \_ v1.52.0

Watching: .go .html .json
Building...
✓ Built successfully (2.3s)
✓ Server started on :8080
```

**What happens:**
1. Air monitors file changes
2. Rebuilds on `.go` file modifications
3. Restarts server automatically
4. Preserves terminal output

---

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--hot` | `-h` | bool | false | Enable hot reload with Air |
| `--port` | `-p` | string | 8080 | HTTP server port |
| `--config` | `-c` | string | .env | Config file path |
| `--no-migrate` | | bool | false | Skip auto-migration |
| `--verbose` | `-v` | bool | false | Enable verbose logging |

### Examples with Flags

```bash
# Hot reload with custom port
neonex serve --hot --port 3000

# Different config file
neonex serve --config .env.production

# Skip migration
neonex serve --no-migrate

# Verbose logging
neonex serve --hot --verbose

# Combine multiple flags
neonex serve --hot --port 9000 --verbose
```

---

## Hot Reload Mode

### How It Works

**Air** (github.com/air-verse/air) watches your files and rebuilds automatically:

```
File Changed → Build → Restart → Ready
   (0.5s)      (2s)     (0.3s)    ✓
```

**Workflow:**
1. Edit `modules/user/service.go`
2. Save file
3. Air detects change
4. Rebuilds project
5. Restarts server
6. You're ready to test!

### Watched Files

By default, Air watches:
- `*.go` - Go source files
- `*.html` - HTML templates
- `*.json` - JSON configuration

**Configuration:** `.air.toml` in project root

### Configuration Example

```toml
# .air.toml
[build]
  cmd = "go build -o ./tmp/main.exe ."
  bin = "tmp/main.exe"
  include_ext = ["go", "html", "json", "yaml"]
  exclude_dir = ["tmp", "vendor", "node_modules"]
  delay = 1000  # ms to wait before rebuild
```

### First-Time Setup

If Air is not installed:

```bash
neonex serve --hot
```

**Output:**
```
Air not found. Install Air for hot reload? (y/n): y
Installing Air...
✓ Installed github.com/air-verse/air@latest
Starting with hot reload...
```

**Manual installation:**
```bash
go install github.com/air-verse/air@latest
```

---

## Normal Mode

### Without Hot Reload

```bash
neonex serve
```

**Use when:**
- ❌ Air not installed
- ❌ Don't need auto-reload
- ✅ Quick testing
- ✅ CI/CD environments
- ✅ Production-like testing

**Restart manually:**
- Stop: `Ctrl+C`
- Start: `neonex serve`

---

## Port Configuration

### Custom Port

```bash
# Via flag
neonex serve --port 3000

# Via environment variable
APP_PORT=3000 neonex serve

# Via .env file
echo "APP_PORT=3000" >> .env
neonex serve
```

**Priority:**
1. Command flag (`--port`)
2. Environment variable (`APP_PORT`)
3. .env file
4. Default (8080)

### Port Conflicts

If port is already in use:

```bash
Error: listen tcp :8080: bind: address already in use
```

**Solutions:**

```bash
# Use different port
neonex serve --port 3001

# Find and kill process (Linux/Mac)
lsof -ti:8080 | xargs kill

# Find process (Windows)
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

---

## Database Migration

### Auto-Migration

By default, Neonex auto-migrates all models:

```bash
neonex serve
```

**Output:**
```
✓ Auto-migrated: User, Product, Order
```

### Skip Migration

```bash
neonex serve --no-migrate
```

**Use when:**
- Database schema is stable
- Running multiple instances
- Testing without DB changes

### Manual Migration

```go
// In your code
db := database.GetDB()
db.AutoMigrate(&models.User{}, &models.Product{})
```

---

## Logging

### Default Logging

```bash
neonex serve
```

**Output:**
```
2024/11/28 10:30:00 [INFO] Starting application
2024/11/28 10:30:01 [INFO] Database connected
2024/11/28 10:30:02 [INFO] Server listening on :8080
```

### Verbose Logging

```bash
neonex serve --verbose
```

**Output:**
```
2024/11/28 10:30:00 [DEBUG] Loading environment variables
2024/11/28 10:30:00 [DEBUG] Connecting to database: sqlite://neonex.db
2024/11/28 10:30:01 [DEBUG] Discovering modules in: ./modules
2024/11/28 10:30:01 [INFO] Registered module: user
2024/11/28 10:30:01 [INFO] Registered module: product
2024/11/28 10:30:02 [DEBUG] Registering routes: /user/
2024/11/28 10:30:02 [DEBUG] Registering routes: /product/
2024/11/28 10:30:02 [INFO] Server started successfully
```

### Request Logging

HTTP requests are automatically logged:

```
[2024-11-28 10:31:15] GET /user/ - 200 OK (12ms)
[2024-11-28 10:31:20] POST /user/ - 201 Created (45ms)
[2024-11-28 10:31:25] GET /user/1 - 200 OK (8ms)
```

---

## Development Workflow

### Typical Session

```bash
# 1. Start server with hot reload
neonex serve --hot

# 2. Open in browser
# http://localhost:8080

# 3. Edit code
# modules/user/service.go

# 4. Save file
# → Air auto-rebuilds

# 5. Test changes
# curl http://localhost:8080/user/

# 6. Repeat steps 3-5
```

### Multiple Terminals

**Terminal 1: Server**
```bash
neonex serve --hot
```

**Terminal 2: Testing**
```bash
# Test endpoints
curl http://localhost:8080/user/

# Run tests
go test ./...

# Generate modules
neonex module create product
```

**Terminal 3: Database**
```bash
# Check database
sqlite3 neonex.db "SELECT * FROM users;"

# Or use database client
```

---

## Production Mode

### Build for Production

```bash
# Build binary
go build -o bin/app main.go

# Run without neonex CLI
./bin/app

# Or with environment
APP_ENV=production ./bin/app
```

### Production Checklist

Before deploying:

- [ ] Set `APP_ENV=production`
- [ ] Use production database
- [ ] Configure proper logging (`LOG_LEVEL=warn`)
- [ ] Set secure passwords
- [ ] Enable HTTPS
- [ ] Configure CORS properly
- [ ] Set resource limits

---

## Examples

### Example 1: Basic Development

```bash
# Start server
cd my-app
neonex serve --hot

# Edit file
echo 'func NewFeature() {}' >> modules/user/service.go

# Air rebuilds automatically
# Test: curl http://localhost:8080/user/
```

### Example 2: Multiple Ports

```bash
# Service 1 (port 8080)
cd service1
neonex serve --hot --port 8080

# Service 2 (port 8081)
cd service2
neonex serve --hot --port 8081

# Service 3 (port 8082)
cd service3
neonex serve --hot --port 8082
```

### Example 3: Custom Environment

```bash
# Development
neonex serve --hot --config .env.development

# Staging
neonex serve --config .env.staging

# Test production setup
neonex serve --config .env.production --no-migrate
```

### Example 4: API Gateway + Services

```bash
# Terminal 1: API Gateway
cd api-gateway
neonex serve --hot --port 8080

# Terminal 2: User Service
cd user-service
neonex serve --hot --port 8081

# Terminal 3: Product Service
cd product-service
neonex serve --hot --port 8082
```

---

## Troubleshooting

### Air Not Working

**Problem:**
```bash
Error: Air command not found
```

**Solution:**
```bash
# Install Air
go install github.com/air-verse/air@latest

# Verify installation
air -v

# Add to PATH if needed
export PATH=$PATH:$GOPATH/bin
```

### Build Failures

**Problem:**
```bash
Build failed: syntax error
```

**Solution:**
1. Check error message in terminal
2. Fix the syntax error
3. Air will auto-rebuild on save
4. No need to restart `neonex serve`

### Port Already in Use

**Problem:**
```bash
Error: bind: address already in use
```

**Solution:**
```bash
# Option 1: Use different port
neonex serve --port 3000

# Option 2: Kill existing process
# Linux/Mac
lsof -ti:8080 | xargs kill -9

# Windows
netstat -ano | findstr :8080
taskkill /F /PID <PID>
```

### Database Locked

**Problem:**
```bash
Error: database is locked
```

**Solution:**
```bash
# Close other connections
# Stop other instances of the app

# Use different database file
DB_NAME=dev-copy.db neonex serve
```

### Too Many Rebuilds

**Problem:**
Air rebuilds constantly even without changes.

**Solution:**
Edit `.air.toml`:
```toml
[build]
  delay = 3000  # Increase delay (ms)
  exclude_dir = ["tmp", "vendor", "node_modules", ".git"]
  exclude_regex = ["_test.go", ".*_mock.go"]
```

---

## Advanced Usage

### Custom Build Command

Edit `.air.toml`:
```toml
[build]
  cmd = "go build -tags=dev -o ./tmp/main ."
  full_bin = "./tmp/main --dev-mode"
```

### Environment-Specific Air Config

```bash
# Development
air -c .air.dev.toml

# Production testing
air -c .air.prod.toml
```

### Pre/Post Commands

```toml
[build]
  cmd = "go generate ./... && go build -o ./tmp/main ."
  full_bin = "npm run build:assets && ./tmp/main"
```

---

## Performance Tips

### 1. Exclude Unnecessary Directories

```toml
exclude_dir = [
  "tmp", "vendor", "node_modules", 
  ".git", "docs", "scripts"
]
```

### 2. Limit File Extensions

```toml
include_ext = ["go"]  # Only watch .go files
```

### 3. Increase Build Delay

```toml
delay = 2000  # Wait 2s before rebuilding
```

### 4. Use Build Cache

```bash
# Enable Go build cache
export GOCACHE=$(go env GOCACHE)
```

---

## Next Steps

- [**neonex module**](neonex-module.md) - Create modules while server runs
- [**Hot Reload Guide**](../development/hot-reload.md) - Deep dive into Air
- [**Debugging**](../development/debugging.md) - Debug running applications
- [**Production Deployment**](../deployment/production-setup.md) - Deploy your app

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
