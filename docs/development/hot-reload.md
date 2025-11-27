# Hot Reload Development

Learn how to use hot reload for rapid development in Neonex Core with Air.

---

## Overview

**Hot reload** automatically restarts your application when code changes are detected, enabling instant feedback during development.

**Benefits:**
- ⚡ **Instant feedback** - See changes immediately
- 🔄 **Auto-restart** - No manual server restarts
- 💻 **Better DX** - Smooth development experience
- 🎯 **Focus** - Stay in your editor

**Tool:** [Air](https://github.com/air-verse/air) - Live reload for Go apps

---

## Quick Start

### Install Air

```powershell
# Install Air globally
go install github.com/air-verse/air@latest

# Verify installation
air -v
```

### Start with Hot Reload

```powershell
# Using air command
air

# Or using neonex CLI
neonex serve

# Or using PowerShell script
.\run.ps1
```

**Output:**
```
  __    _   ___  
 / /\  | | | |_) 
/_/--\ |_| |_| \_ v1.52.0

watching .
building...
running...

[INFO] Server started on :3000
```

---

## Configuration

### Air Configuration File

**Location:** `.air.toml` (project root)

```toml
# .air.toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  # Binary output
  bin = "./tmp/main.exe"
  
  # Build command
  cmd = "go build -o ./tmp/main.exe ./cmd/neonex"
  
  # Delay before restart (ms)
  delay = 1000
  
  # Exclude directories
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "docs"]
  
  # Exclude files
  exclude_file = []
  
  # Watch these extensions
  include_ext = ["go", "tpl", "tmpl", "html", "env"]
  
  # Exclude patterns
  exclude_regex = ["_test.go"]
  
  # Stop on error
  stop_on_error = false
  
  # Send interrupt signal before kill
  send_interrupt = true
  
  # Delay after sending interrupt signal
  kill_delay = 500

[color]
  main = "magenta"
  watcher = "cyan"
  build = "yellow"
  runner = "green"

[log]
  time = false

[misc]
  clean_on_exit = true
```

---

## Project Structure

### Required Files

```
your-project/
├── .air.toml          # Air configuration
├── cmd/
│   └── neonex/
│       └── main.go    # Entry point
├── tmp/               # Build output (auto-created)
│   └── main.exe
└── ...
```

### Temporary Directory

```toml
# .air.toml
tmp_dir = "tmp"
```

**Purpose:**
- Stores compiled binary
- Auto-created by Air
- Cleaned on exit (if configured)
- Should be in `.gitignore`

**.gitignore:**
```
tmp/
*.exe
```

---

## Watch Configuration

### File Extensions

```toml
# Watch these file types
include_ext = ["go", "tpl", "tmpl", "html", "env"]
```

**Common extensions:**
- `.go` - Go source files
- `.env` - Environment variables
- `.html`, `.tmpl` - Templates
- `.json` - Configuration files

### Exclude Directories

```toml
# Don't watch these directories
exclude_dir = [
    "assets",
    "tmp",
    "vendor",
    "testdata",
    "docs",
    "node_modules"
]
```

### Exclude Files

```toml
# Don't watch test files
exclude_regex = ["_test.go"]

# Or specific files
exclude_file = ["README.md"]
```

---

## Build Configuration

### Basic Build

```toml
[build]
  cmd = "go build -o ./tmp/main.exe ./cmd/neonex"
  bin = "./tmp/main.exe"
```

### Build with Tags

```toml
[build]
  cmd = "go build -tags dev -o ./tmp/main.exe ./cmd/neonex"
```

### Build with Multiple Commands

```toml
[build]
  # Run multiple commands
  cmd = "swag init && go build -o ./tmp/main.exe ./cmd/neonex"
```

### Custom Arguments

```toml
[build]
  # Pass arguments to binary
  args_bin = ["--port", "3000", "--env", "development"]
```

---

## Development Workflow

### Typical Development Session

```powershell
# 1. Start air
air

# 2. Edit code
# File: modules/user/service.go

# 3. Save file
# Air automatically:
#   - Detects change
#   - Rebuilds binary
#   - Restarts server

# 4. Test changes
# http://localhost:3000
```

### Console Output

```
[INFO] 14:30:45 | file changed: modules/user/service.go
[INFO] 14:30:45 | Building...
[INFO] 14:30:46 | Build finished
[INFO] 14:30:46 | Restarting...
[INFO] 14:30:46 | Server started on :3000
```

---

## Common Use Cases

### Watch Environment Files

```toml
include_ext = ["go", "env"]
```

**Use case:** Reload when `.env` changes

### Exclude Test Files

```toml
exclude_regex = ["_test.go"]
```

**Use case:** Don't rebuild on test file changes

### Watch Templates

```toml
include_ext = ["go", "html", "tmpl"]
```

**Use case:** Reload when templates change

### Fast Builds

```toml
[build]
  delay = 100  # Shorter delay for fast rebuilds
```

**Use case:** Rapid iteration

---

## Integration with Neonex CLI

### neonex serve Command

```go
// internal/commands/serve.go
var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start development server with hot reload",
    Run: func(cmd *cobra.Command, args []string) {
        // Check if air is installed
        if !isAirInstalled() {
            fmt.Println("Air not found. Install with:")
            fmt.Println("  go install github.com/air-verse/air@latest")
            os.Exit(1)
        }
        
        // Run air
        airCmd := exec.Command("air")
        airCmd.Stdout = os.Stdout
        airCmd.Stderr = os.Stderr
        
        if err := airCmd.Run(); err != nil {
            fmt.Printf("Error: %v\n", err)
            os.Exit(1)
        }
    },
}
```

### Usage

```powershell
# Start with hot reload
neonex serve

# Equivalent to
air
```

---

## PowerShell Integration

### run.ps1 Script

```powershell
# run.ps1
Write-Host "Starting Neonex development server..." -ForegroundColor Green

# Check if air is installed
$airPath = Get-Command air -ErrorAction SilentlyContinue
if (-not $airPath) {
    Write-Host "Air not installed. Installing..." -ForegroundColor Yellow
    go install github.com/air-verse/air@latest
}

# Start air
air
```

### Usage

```powershell
# Run script
.\run.ps1

# Or with execution policy
powershell -ExecutionPolicy Bypass -File .\run.ps1
```

---

## Troubleshooting

### Air Not Found

**Problem:**
```
'air' is not recognized as an internal or external command
```

**Solution:**
```powershell
# Install air
go install github.com/air-verse/air@latest

# Add Go bin to PATH
$env:PATH += ";$env:USERPROFILE\go\bin"

# Or permanently in system settings
[System.Environment]::SetEnvironmentVariable(
    "Path",
    $env:Path + ";$env:USERPROFILE\go\bin",
    [System.EnvironmentVariableTarget]::User
)
```

### Build Errors

**Problem:**
```
[ERROR] Build failed
```

**Solution:**
```powershell
# Check build manually
go build -o ./tmp/main.exe ./cmd/neonex

# Fix any compilation errors
# Then restart air
```

### Port Already in Use

**Problem:**
```
Error: address already in use
```

**Solution:**
```powershell
# Find process using port 3000
netstat -ano | findstr :3000

# Kill process
taskkill /PID <PID> /F

# Or use different port
$env:PORT = "3001"
air
```

### Too Many Rebuilds

**Problem:**
- Air rebuilding constantly
- Watching too many files

**Solution:**
```toml
# Exclude more directories
exclude_dir = [
    "tmp",
    "vendor",
    "node_modules",
    "docs",
    ".git"
]

# Increase delay
[build]
  delay = 2000  # Wait 2 seconds before rebuild
```

### Slow Rebuilds

**Problem:**
- Builds taking too long

**Solution:**
```toml
# Disable CGO for faster builds
[build]
  cmd = "CGO_ENABLED=0 go build -o ./tmp/main.exe ./cmd/neonex"
```

---

## Advanced Configuration

### Multiple Build Commands

```toml
[build]
  # Generate code, then build
  cmd = """
    go generate ./... &&
    swag init &&
    go build -o ./tmp/main.exe ./cmd/neonex
  """
```

### Custom Pre/Post Commands

```toml
[build]
  # Pre-build command
  pre_cmd = ["echo Building..."]
  
  # Post-build command
  post_cmd = ["echo Build complete!"]
```

### Environment Variables

```toml
[build]
  # Pass environment variables
  args_bin = []
  
  # In Go code, load from .env
  # godotenv.Load(".env")
```

---

## VS Code Integration

### Launch Configuration

**.vscode/launch.json:**
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch with Air",
      "type": "go",
      "request": "launch",
      "mode": "exec",
      "program": "${workspaceFolder}/tmp/main.exe",
      "preLaunchTask": "air"
    }
  ]
}
```

### Task Configuration

**.vscode/tasks.json:**
```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "air",
      "type": "shell",
      "command": "air",
      "isBackground": true,
      "problemMatcher": {
        "pattern": {
          "regexp": "^(.*):(\\d+):(\\d+):\\s+(warning|error):\\s+(.*)$",
          "file": 1,
          "line": 2,
          "column": 3,
          "severity": 4,
          "message": 5
        },
        "background": {
          "activeOnStart": true,
          "beginsPattern": "building...",
          "endsPattern": "running..."
        }
      }
    }
  ]
}
```

---

## Docker Development

### Dockerfile with Air

```dockerfile
# Dockerfile.dev
FROM golang:1.21-alpine

WORKDIR /app

# Install air
RUN go install github.com/air-verse/air@latest

# Copy air config
COPY .air.toml .

# Copy source
COPY . .

# Run with air
CMD ["air"]
```

### Docker Compose

```yaml
# docker-compose.dev.yml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile.dev
    volumes:
      - .:/app
    ports:
      - "3000:3000"
    environment:
      - ENV=development
```

**Usage:**
```powershell
docker-compose -f docker-compose.dev.yml up
```

---

## Best Practices

### ✅ DO:

**1. Exclude Unnecessary Directories**
```toml
exclude_dir = ["vendor", "tmp", "docs", "node_modules"]
```

**2. Use Appropriate Delay**
```toml
delay = 1000  # 1 second is usually good
```

**3. Clean on Exit**
```toml
[misc]
  clean_on_exit = true
```

### ❌ DON'T:

**1. Watch Too Many Files**
```toml
# Bad: Watching everything
exclude_dir = []  # ❌
```

**2. No Delay**
```toml
# Bad: Rebuilds too often
delay = 0  # ❌
```

**3. Commit Binaries**
```
# Always add to .gitignore
tmp/
*.exe
```

---

## Comparison: With vs Without Air

### Without Air

```powershell
# Manual workflow
go build -o app.exe ./cmd/neonex
./app.exe

# Code change...
Ctrl+C  # Stop server
go build -o app.exe ./cmd/neonex
./app.exe
# Repeat...
```

### With Air

```powershell
# Automated workflow
air

# Code change...
# Auto rebuild & restart
# Keep coding!
```

**Time saved:** ~10-30 seconds per change

---

## Next Steps

- [**Testing**](testing.md) - Testing strategies
- [**Debugging**](debugging.md) - Debug techniques
- [**Best Practices**](best-practices.md) - Development tips
- [**Quick Start**](../getting-started/quick-start.md) - Get started

---

**Need help?** Check [FAQ](../resources/faq.md) or [get support](../resources/support.md).
