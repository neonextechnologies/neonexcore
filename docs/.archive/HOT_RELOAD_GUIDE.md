# Hot Reload Development Guide

Complete guide for using hot reload in Neonex Core development.

## What is Hot Reload?

Hot reload automatically rebuilds and restarts your application when you make code changes. This speeds up development by eliminating the need to manually stop and restart the server.

## Setup

### 1. Install Air

Air is the hot reload tool we use for Go applications.

**Windows (PowerShell):**
```powershell
go install github.com/cosmtrek/air@latest
```

**Linux/Mac:**
```bash
go install github.com/cosmtrek/air@latest
```

**Verify installation:**
```bash
air -v
```

### 2. Configuration

The `.air.toml` file is already configured in the project root:

```toml
[build]
  cmd = "go build -o ./tmp/main ."
  bin = "./tmp/main"
  include_ext = ["go", "tpl", "tmpl", "html", "json"]
  exclude_dir = ["tmp", "vendor", "logs"]
  delay = 1000
```

## Usage

### Method 1: Using Neonex CLI (Recommended)

```bash
# Start with hot reload
neonex serve --hot
# or
neonex serve -w

# Start without hot reload
neonex serve

# Custom port
neonex serve --hot --port 3000
```

### Method 2: Using Make

```bash
# Start development server with hot reload
make dev

# Other commands
make serve    # Alias for dev
make watch    # Alias for dev
make run      # Run without hot reload
```

### Method 3: Using PowerShell Scripts (Windows)

```powershell
# Load dev tools
. .\dev.ps1

# Start with hot reload
Start-Dev
# or simply
dev

# Start without hot reload
Start-Server
```

### Method 4: Using Bash Scripts (Linux/Mac)

```bash
# Make script executable
chmod +x dev.sh

# Start with hot reload
./dev.sh dev

# Start without hot reload
./dev.sh serve
```

### Method 5: Direct Air Command

```bash
air
```

## What Gets Watched?

Air watches these file types:
- `.go` files
- `.json` files
- `.html`, `.tpl`, `.tmpl` templates

These directories are excluded:
- `tmp/` - Build output
- `vendor/` - Dependencies
- `logs/` - Log files
- `testdata/` - Test data
- `test-project/` - Test projects

## How It Works

1. **Watch**: Air monitors your files for changes
2. **Build**: When changes detected, rebuilds the application
3. **Kill**: Stops the running process
4. **Start**: Runs the newly built binary
5. **Repeat**: Continues watching

## Build Output

- Binary location: `./tmp/main` (Linux/Mac) or `./tmp/main.exe` (Windows)
- Build errors: `build-errors.log`
- Console shows colored output:
  - 🟡 Yellow: Build process
  - 🟣 Magenta: Main app output
  - 🟢 Green: Air runner messages
  - 🔵 Cyan: File watcher

## Development Workflow

### Typical Session

```bash
# 1. Start hot reload
neonex serve --hot

# 2. Make changes to your code
# 3. Save the file
# 4. Air automatically rebuilds and restarts
# 5. Test your changes
# 6. Repeat from step 2
```

### Example Output

```
🔥 Starting server with hot reload...
  
  __    _   ___  
 / /\  | | | |_) 
/_/--\ |_| |_| \_ , built with Go 

watching...
building...
running...

[2024-11-25 10:30:00] INFO  Logger initialized
[2024-11-25 10:30:00] INFO  Database initialized
[2024-11-25 10:30:00] INFO  HTTP server starting | port=8080

┌───────────────────────────────────────────────────┐
│              Neonex Core v0.1-alpha               │
│               http://127.0.0.1:8080               │
└───────────────────────────────────────────────────┘
```

When you save a file:
```
main.go has changed
building...
running...
```

## Tips & Best Practices

### 1. Clean Rebuilds

If you encounter issues, clean the tmp directory:

```bash
# Using make
make clean

# Using PowerShell
Clean-Build

# Manually
rm -rf tmp/
```

### 2. Exclude Directories

Add directories to exclude in `.air.toml`:

```toml
[build]
  exclude_dir = ["tmp", "vendor", "logs", "yourdir"]
```

### 3. Custom Build Commands

Modify the build command in `.air.toml`:

```toml
[build]
  cmd = "go build -tags dev -o ./tmp/main ."
```

### 4. Disable Specific Files

```toml
[build]
  exclude_regex = ["_test.go", "_generated.go"]
```

### 5. Adjust Delay

Change the delay before rebuilding (milliseconds):

```toml
[build]
  delay = 1000  # 1 second
```

## Troubleshooting

### Air not found

```bash
# Reinstall Air
go install github.com/cosmtrek/air@latest

# Check GOPATH/bin is in PATH
echo $PATH
```

### Build errors not showing

Check `build-errors.log`:

```bash
cat build-errors.log
```

### Port already in use

Kill the process using the port:

**Windows:**
```powershell
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

**Linux/Mac:**
```bash
lsof -ti:8080 | xargs kill -9
```

### Changes not detected

1. Save the file properly
2. Check file is not in excluded directories
3. Check file extension is in `include_ext`
4. Restart Air

### Too many rebuilds

Increase the delay:

```toml
[build]
  delay = 2000  # 2 seconds
```

## Advanced Configuration

### Custom Air Config

Create a custom `.air.toml`:

```toml
root = "."
tmp_dir = "tmp"

[build]
  # Custom build command
  cmd = "go build -race -o ./tmp/main ."
  
  # Custom binary location
  bin = "./tmp/main"
  
  # Full command with arguments
  full_bin = "./tmp/main --debug"
  
  # Delay before restart
  delay = 1000
  
  # Stop on build error
  stop_on_error = true
  
  # Send interrupt signal
  send_interrupt = true
  
  # Kill delay
  kill_delay = "500ms"

[log]
  # Show only main output
  main_only = false
  
  # Show timestamp
  time = true

[color]
  main = "magenta"
  watcher = "cyan"
  build = "yellow"
  runner = "green"

[misc]
  # Clean on exit
  clean_on_exit = true

[screen]
  # Clear on rebuild
  clear_on_rebuild = true
```

### Pre/Post Build Commands

Use shell scripts:

```toml
[build]
  # Run before build
  pre_cmd = ["echo Building...", "go generate ./..."]
  
  # Run after build
  post_cmd = ["echo Build complete!"]
```

## Integration with IDEs

### VS Code

Create `.vscode/tasks.json`:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Air: Hot Reload",
      "type": "shell",
      "command": "air",
      "problemMatcher": [],
      "group": {
        "kind": "build",
        "isDefault": true
      }
    }
  ]
}
```

Press `Ctrl+Shift+B` to start.

### GoLand/IntelliJ

1. Run → Edit Configurations
2. Add new "Go Build"
3. Program: `air`
4. Working directory: `$ProjectFileDir$`

## Performance Tips

1. **Use fast disk**: SSD recommended
2. **Exclude unnecessary directories**: Reduce file watching overhead
3. **Increase delay**: Prevent too many rebuilds
4. **Use build cache**: Go build cache speeds up rebuilds
5. **Disable unused features**: Comment out in code during dev

## Comparison: With vs Without Hot Reload

**Without Hot Reload:**
```
1. Make code change
2. Ctrl+C to stop server
3. Run `go run main.go`
4. Wait for build
5. Test changes
Time: ~10-30 seconds per iteration
```

**With Hot Reload:**
```
1. Make code change
2. Save file
3. Wait for auto-rebuild (1-3 seconds)
4. Test changes
Time: ~2-5 seconds per iteration
```

**Productivity gain: 5-10x faster development cycle!**

## Resources

- Air GitHub: https://github.com/cosmtrek/air
- Air Documentation: https://github.com/cosmtrek/air/blob/master/README.md
- Go Build Cache: https://pkg.go.dev/cmd/go#hdr-Build_and_test_caching

---

**Happy Hot Reloading! 🔥**
