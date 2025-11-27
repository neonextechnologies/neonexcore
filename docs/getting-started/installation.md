# Installation

## Prerequisites

Before installing Neonex Core, ensure you have:

* **Go 1.21 or higher** - [Download Go](https://golang.org/dl/)
* **Git** - [Download Git](https://git-scm.com/downloads)
* **(Optional) Air** - For hot reload development

### Verify Go Installation

```bash
go version
# Should output: go version go1.21.x or higher
```

## Installation Methods

### Method 1: Clone from GitHub (Recommended)

```bash
# Clone the repository
git clone https://github.com/neonextechnologies/neonexcore.git
cd neonexcore

# Install dependencies
go mod download

# Build CLI tool
go build -o neonex.exe ./cmd/neonex

# Verify installation
./neonex --version
```

### Method 2: Using Go Install

```bash
# Install CLI globally
go install github.com/neonextechnologies/neonexcore/cmd/neonex@latest

# Verify installation
neonex --version
```

### Method 3: Download Binary

Download pre-built binaries from the [Releases page](https://github.com/neonextechnologies/neonexcore/releases).

## Optional: Install Air (Hot Reload)

For development with hot reload:

```bash
go install github.com/air-verse/air@latest

# Verify Air installation
air -v
```

## Post-Installation

### Create Your First Project

```bash
# Create a new project
neonex new my-first-app

# Navigate to project
cd my-first-app

# Install dependencies
go mod download

# Run the application
neonex serve
```

The server will start at `http://localhost:8080`.

### Verify Installation

Test the default endpoint:

```bash
curl http://localhost:8080
```

Expected response:

```json
{
  "framework": "Neonex Core",
  "version": "0.1-alpha",
  "status": "running",
  "engine": "Fiber (fasthttp)"
}
```

## Troubleshooting

### Go Command Not Found

If `go` command is not recognized:

1. Ensure Go is installed correctly
2. Add Go to your PATH:
   * **Windows**: Add `C:\Go\bin` to PATH
   * **Linux/Mac**: Add `/usr/local/go/bin` to PATH

### Permission Denied (Linux/Mac)

If you get permission errors:

```bash
chmod +x neonex
```

### Module Download Issues

If `go mod download` fails:

```bash
# Set Go proxy (if behind firewall)
go env -w GOPROXY=https://goproxy.io,direct

# Try again
go mod download
```

## Next Steps

* [Quick Start Guide](quick-start.md)
* [Project Structure](project-structure.md)
* [Configuration](configuration.md)

## System Requirements

### Minimum

* **CPU**: 1 core
* **RAM**: 512 MB
* **Disk**: 100 MB

### Recommended

* **CPU**: 2+ cores
* **RAM**: 2 GB+
* **Disk**: 1 GB+
* **OS**: Windows 10+, Ubuntu 20.04+, macOS 11+

## Supported Platforms

Neonex Core supports all platforms that Go supports:

* ✅ Windows (amd64, arm64)
* ✅ Linux (amd64, arm64, arm)
* ✅ macOS (amd64, arm64)
* ✅ FreeBSD, OpenBSD
