# Neonex Core Development Scripts

# Quick start with hot reload
function Start-Dev {
    Write-Host "🔥 Starting development server with hot reload..." -ForegroundColor Cyan
    air
}

# Run without hot reload
function Start-Server {
    Write-Host "🚀 Starting server..." -ForegroundColor Green
    go run main.go
}

# Build the application
function Build-App {
    Write-Host "🔨 Building application..." -ForegroundColor Yellow
    go build -o neonex.exe .
    Write-Host "✅ Build complete: neonex.exe" -ForegroundColor Green
}

# Build CLI
function Build-CLI {
    Write-Host "🔨 Building CLI..." -ForegroundColor Yellow
    go build -o neonex.exe ./cmd/neonex
    Write-Host "✅ CLI built: neonex.exe" -ForegroundColor Green
}

# Install Air
function Install-Air {
    Write-Host "📦 Installing Air..." -ForegroundColor Cyan
    go install github.com/air-verse/air@latest
    Write-Host "✅ Air installed" -ForegroundColor Green
}

# Install dependencies
function Install-Deps {
    Write-Host "📦 Installing dependencies..." -ForegroundColor Cyan
    go mod download
    go mod tidy
    Write-Host "✅ Dependencies installed" -ForegroundColor Green
}

# Run tests
function Run-Tests {
    Write-Host "🧪 Running tests..." -ForegroundColor Cyan
    go test -v ./...
}

# Clean build artifacts
function Clean-Build {
    Write-Host "🧹 Cleaning..." -ForegroundColor Yellow
    Remove-Item -Path "tmp" -Recurse -Force -ErrorAction SilentlyContinue
    Remove-Item -Path "*.exe" -Force -ErrorAction SilentlyContinue
    Remove-Item -Path "build-errors.log" -Force -ErrorAction SilentlyContinue
    Write-Host "✅ Clean complete" -ForegroundColor Green
}

# Format code
function Format-Code {
    Write-Host "🎨 Formatting code..." -ForegroundColor Cyan
    go fmt ./...
    Write-Host "✅ Code formatted" -ForegroundColor Green
}

# Create new module
function New-Module {
    param([string]$Name)
    if (-not $Name) {
        Write-Host "❌ Module name required: New-Module -Name modulename" -ForegroundColor Red
        return
    }
    Write-Host "📦 Creating module: $Name" -ForegroundColor Cyan
    .\neonex.exe module create $Name
}

# List modules
function List-Modules {
    Write-Host "📦 Available modules:" -ForegroundColor Cyan
    .\neonex.exe module list
}

# Show help
function Show-Help {
    Write-Host @"

Neonex Core - Development Commands
===================================

Quick Start:
  Start-Dev          - Start with hot reload (air)
  Start-Server       - Start without hot reload
  
Build:
  Build-App          - Build main application
  Build-CLI          - Build CLI tool
  Clean-Build        - Clean build artifacts
  
Development:
  Install-Air        - Install Air for hot reload
  Install-Deps       - Install Go dependencies
  Format-Code        - Format all Go code
  Run-Tests          - Run all tests
  
Modules:
  New-Module -Name x - Create new module
  List-Modules       - List all modules
  
Examples:
  PS> Start-Dev
  PS> New-Module -Name product
  PS> Build-CLI
  PS> List-Modules

"@ -ForegroundColor White
}

# Aliases
Set-Alias dev Start-Dev
Set-Alias serve Start-Server
Set-Alias build Build-App
Set-Alias clean Clean-Build
Set-Alias test Run-Tests
Set-Alias fmt Format-Code

# Export functions
Export-ModuleMember -Function * -Alias *

# Show welcome message
Write-Host "✨ Neonex Core Dev Tools Loaded!" -ForegroundColor Green
Write-Host "   Run 'Show-Help' for available commands" -ForegroundColor Gray
