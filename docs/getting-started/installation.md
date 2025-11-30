# Installation

Learn how to install Neonex Core and get started building high-performance Go applications.

---

## Prerequisites

Before installing Neonex Core, ensure you have the following installed on your system:

| Requirement | Version | Purpose | Download |
|------------|---------|---------|----------|
| **Go** | 1.21+ | Core runtime | [golang.org/dl](https://golang.org/dl/) |
| **Git** | Latest | Version control | [git-scm.com](https://git-scm.com/downloads) |
| **PostgreSQL** | 14+ (Optional) | Database | [postgresql.org](https://www.postgresql.org/download/) |
| **Air** | Latest (Optional) | Hot reload | `go install github.com/air-verse/air@latest` |

### Verify Prerequisites

Check if you have the required tools installed:

```bash
# Check Go version
go version
# Expected: go version go1.21.x or higher

# Check Git
git --version
# Expected: git version 2.x.x or higher

# Check PostgreSQL (if installed)
psql --version
# Expected: psql (PostgreSQL) 14.x or higher
```

> **💡 Tip:** Neonex Core uses Go generics, which requires Go 1.21 or higher. Earlier versions will not work.

---

## Installation Methods

Choose the installation method that best fits your workflow:

### 🚀 Method 1: Clone from GitHub (Recommended)

This is the recommended approach for development and contribution:

```bash
# Clone the repository
git clone https://github.com/neonextechnologies/neonexcore.git
cd neonexcore

# Install dependencies
go mod download

# Build the CLI tool
go build -o neonex.exe ./cmd/neonex

# (Optional) Add to PATH for global access
# Windows PowerShell:
$env:Path += ";$(pwd)"
# Windows CMD:
set PATH=%PATH%;%cd%
# Linux/Mac:
export PATH=$PATH:$(pwd)

# Verify installation
./neonex --version
```

**Advantages:**
- ✅ Full source code access
- ✅ Easy to contribute and customize
- ✅ Latest features and updates
- ✅ Can debug framework internals

### 📦 Method 2: Using Go Install

Quick installation for using Neonex Core as a CLI tool:

```bash
# Install CLI globally
go install github.com/neonextechnologies/neonexcore/cmd/neonex@latest

# Verify installation
neonex --version
```

**Advantages:**
- ✅ Quick and simple
- ✅ Automatically added to PATH
- ✅ Easy to update with same command

> **📌 Note:** This installs the CLI tool only. To work with the framework, you'll still need to create a new project.

### 💾 Method 3: Download Pre-Built Binary

Download pre-built binaries for your platform:

1. Visit [Releases page](https://github.com/neonextechnologies/neonexcore/releases)
2. Download the binary for your OS:
   - `neonex-windows-amd64.exe` (Windows 64-bit)
   - `neonex-linux-amd64` (Linux 64-bit)
   - `neonex-darwin-amd64` (macOS Intel)
   - `neonex-darwin-arm64` (macOS Apple Silicon)
3. Extract and add to PATH

**Advantages:**
- ✅ No compilation needed
- ✅ Works on systems without Go installed
- ✅ Smaller download size

---

## Optional Tools

### Install Air for Hot Reload

Air automatically rebuilds and restarts your application when files change:

```bash
# Install Air globally
go install github.com/air-verse/air@latest

# Verify installation
air -v
# Expected: air version x.x.x

# Create Air configuration (optional)
air init
```

**Benefits:**
- ⚡ Instant feedback during development
- 🔄 Automatic rebuild on file changes
- 🎯 No manual server restarts
- ⏱️ Saves development time

### Install Database Tools (Optional)

#### PostgreSQL

```bash
# Ubuntu/Debian
sudo apt-get install postgresql-client

# macOS
brew install postgresql

# Windows
# Download from https://www.postgresql.org/download/windows/
```

#### Database Migration Tools

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Verify
migrate -version
```

---

## Quick Start

### Create Your First Project

Once Neonex Core is installed, create your first application:

```bash
# Create a new project
neonex new my-first-app

# Output:
# ✨ Creating new Neonex Core project...
# 📁 Created directory: my-first-app
# 📝 Generated go.mod
# 📦 Created project structure
# ✅ Project created successfully!

# Navigate to project directory
cd my-first-app

# Install project dependencies
go mod download

# Run the application
neonex serve
```

**Expected output:**

```
🚀 Starting Neonex Core...
📦 Loading modules...
✅ User module loaded
✅ Auth module loaded
🌐 HTTP server starting on :8080
⚡ Server is ready to handle requests
```

The server is now running at **`http://localhost:8080`**

### Verify Installation

Test the default API endpoint:

```bash
# Using curl
curl http://localhost:8080

# Or using PowerShell
Invoke-WebRequest -Uri http://localhost:8080 | Select-Object -Expand Content
```

**Expected response:**

```json
{
  "framework": "Neonex Core",
  "version": "0.1-alpha",
  "status": "running",
  "engine": "Fiber (fasthttp)",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Test with Hot Reload

For development with automatic reload:

```bash
# Start with Air
air

# Or use the built-in watch mode
neonex serve --watch
```

Now any changes to your Go files will automatically rebuild and restart the server! 🔥

---

## Troubleshooting

### Common Issues and Solutions

#### ❌ "Go command not found"

**Problem:** The `go` command is not recognized in your terminal.

**Solution:**

1. Ensure Go is installed correctly:
   - Download from [golang.org/dl](https://golang.org/dl/)
   - Follow the installation wizard

2. Add Go to your PATH:

**Windows:**
```powershell
# Check current Go location
where go

# Add to PATH permanently
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Go\bin", "User")

# Restart terminal and verify
go version
```

**Linux/Mac:**
```bash
# Add to .bashrc or .zshrc
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify
go version
```

#### ❌ "Permission denied" (Linux/Mac)

**Problem:** Cannot execute the `neonex` binary.

**Solution:**

```bash
# Make the binary executable
chmod +x neonex

# Verify
./neonex --version
```

#### ❌ "Module download fails"

**Problem:** Cannot download Go modules (network/firewall issues).

**Solution:**

```bash
# Option 1: Use Go proxy
go env -w GOPROXY=https://goproxy.io,direct

# Option 2: Use Chinese mirror (if in China)
go env -w GOPROXY=https://goproxy.cn,direct

# Option 3: Disable checksum validation (not recommended for production)
go env -w GOSUMDB=off

# Try again
go mod download
```

#### ❌ "Port 8080 already in use"

**Problem:** Another application is using port 8080.

**Solution:**

```bash
# Option 1: Stop the other application

# Windows - Find process using port
netstat -ano | findstr :8080
# Kill process by PID
taskkill /PID <PID> /F

# Linux/Mac - Find and kill process
lsof -ti:8080 | xargs kill -9

# Option 2: Use a different port
neonex serve --port 3000
```

#### ❌ "Database connection fails"

**Problem:** Cannot connect to the database.

**Solution:**

1. Check database is running:
```bash
# PostgreSQL
sudo systemctl status postgresql  # Linux
brew services list                # Mac
```

2. Verify connection string in `.env`:
```env
DB_CONNECTION=postgres
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=neonex
DB_USERNAME=postgres
DB_PASSWORD=your_password
```

3. Test connection manually:
```bash
psql -h localhost -U postgres -d neonex
```

#### ❌ "Build fails with 'generics not supported'"

**Problem:** Your Go version is too old.

**Solution:**

```bash
# Check Go version
go version

# If < 1.21, upgrade Go
# Download from https://golang.org/dl/

# Verify upgrade
go version
# Should show: go1.21.x or higher
```

### Getting Help

Still stuck? Here's how to get help:

1. **📖 Read the docs** - Check our [FAQ](../resources/faq.md)
2. **🐛 GitHub Issues** - [Report a bug](https://github.com/neonextechnologies/neonexcore/issues)
3. **💬 Discussions** - [Ask the community](https://github.com/neonextechnologies/neonexcore/discussions)
4. **📧 Email** - contact@neonextechnologies.com

---

## System Requirements

### Minimum Requirements

For basic development and testing:

| Component | Requirement |
|-----------|------------|
| **CPU** | 1 core @ 1.0 GHz |
| **RAM** | 512 MB |
| **Disk** | 100 MB free |
| **Network** | Internet (for initial setup) |

### Recommended Requirements

For optimal development experience:

| Component | Requirement |
|-----------|------------|
| **CPU** | 2+ cores @ 2.0 GHz |
| **RAM** | 2 GB+ |
| **Disk** | 1 GB+ free (SSD preferred) |
| **Network** | Stable internet connection |

### Production Requirements

For production deployments:

| Component | Requirement |
|-----------|------------|
| **CPU** | 4+ cores |
| **RAM** | 4 GB+ |
| **Disk** | 10 GB+ (SSD required) |
| **Network** | High-speed, redundant |
| **Database** | PostgreSQL 14+, MySQL 8+, or SQLite |

> **💡 Performance Tip:** Neonex Core is built on Fiber (fasthttp), which is highly optimized for performance. Even with minimal resources, it can handle thousands of requests per second.

---

## Supported Platforms

Neonex Core supports all platforms that Go supports:

| Platform | Architecture | Status | Notes |
|----------|--------------|--------|-------|
| **Windows** | amd64 | ✅ Fully Supported | Windows 10+ |
| **Windows** | arm64 | ✅ Fully Supported | Surface Pro X, etc. |
| **Linux** | amd64 | ✅ Fully Supported | All major distros |
| **Linux** | arm64 | ✅ Fully Supported | Raspberry Pi 4+, ARM servers |
| **Linux** | arm (32-bit) | ✅ Supported | Raspberry Pi 3 |
| **macOS** | amd64 | ✅ Fully Supported | Intel Macs |
| **macOS** | arm64 | ✅ Fully Supported | M1/M2/M3 Macs |
| **FreeBSD** | amd64 | ⚠️ Experimental | Community tested |
| **OpenBSD** | amd64 | ⚠️ Experimental | Community tested |

### Tested Environments

We actively test on:
- **Windows**: 10, 11, Server 2019, Server 2022
- **Linux**: Ubuntu 20.04+, Debian 11+, CentOS 8+, Alpine 3.15+
- **macOS**: 11 (Big Sur), 12 (Monterey), 13 (Ventura), 14 (Sonoma)

### Container Support

Neonex Core works great in containers:
- ✅ **Docker** - Full support with official images
- ✅ **Kubernetes** - Tested with K8s 1.25+
- ✅ **Podman** - Compatible
- ✅ **Cloud Run** - GCP Cloud Run compatible
- ✅ **AWS ECS/Fargate** - Tested and working
- ✅ **Azure Container Apps** - Compatible

---

## Next Steps

Now that Neonex Core is installed, continue your journey:

<div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 1rem; margin: 2rem 0;">
  <div style="border: 1px solid #e0e0e0; border-radius: 8px; padding: 1.5rem;">
    <h3>🚀 Quick Start</h3>
    <p>Build your first API in 5 minutes</p>
    <a href="quick-start.md">Start Tutorial →</a>
  </div>
  
  <div style="border: 1px solid #e0e0e0; border-radius: 8px; padding: 1.5rem;">
    <h3>📁 Project Structure</h3>
    <p>Understand the project layout</p>
    <a href="project-structure.md">View Structure →</a>
  </div>
  
  <div style="border: 1px solid #e0e0e0; border-radius: 8px; padding: 1.5rem;">
    <h3>⚙️ Configuration</h3>
    <p>Configure your application</p>
    <a href="configuration.md">Learn Configuration →</a>
  </div>
  
  <div style="border: 1px solid #e0e0e0; border-radius: 8px; padding: 1.5rem;">
    <h3>🧩 Module System</h3>
    <p>Learn about modular architecture</p>
    <a href="../core-concepts/module-system.md">Explore Modules →</a>
  </div>
</div>

---

## Additional Resources

### Video Tutorials

- 📺 [Installing Neonex Core (5 min)](https://youtube.com/watch?v=xxx)
- 📺 [Setting up Development Environment (10 min)](https://youtube.com/watch?v=xxx)
- 📺 [First API with Neonex Core (15 min)](https://youtube.com/watch?v=xxx)

### Community

- 💬 [Discord Community](https://discord.gg/neonexcore)
- 🐦 [Twitter Updates](https://twitter.com/neonexcore)
- 📰 [Dev Blog](https://blog.neonextechnologies.com)

### Related Documentation

- [CLI Tools Reference](../cli-tools/overview.md)
- [Database Setup](../database/configuration.md)
- [Deployment Guide](../deployment/production-setup.md)

---

**Need help?** Check our [Support page](../resources/support.md) or [FAQ](../resources/faq.md).

