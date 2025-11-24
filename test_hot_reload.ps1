# Test Hot Reload Setup

Write-Host "🔥 Testing Hot Reload Setup" -ForegroundColor Cyan
Write-Host ""

# Check if Air is installed
Write-Host "1. Checking Air installation..." -ForegroundColor Yellow
$airPath = Get-Command air -ErrorAction SilentlyContinue
if ($airPath) {
    Write-Host "   ✅ Air is installed at: $($airPath.Source)" -ForegroundColor Green
    $airVersion = & air -v 2>&1
    Write-Host "   Version: $airVersion" -ForegroundColor Gray
} else {
    Write-Host "   ❌ Air is not installed" -ForegroundColor Red
    Write-Host "   Installing Air..." -ForegroundColor Yellow
    go install github.com/air-verse/air@latest
    Write-Host "   ✅ Air installed" -ForegroundColor Green
}

Write-Host ""

# Check if .air.toml exists
Write-Host "2. Checking .air.toml configuration..." -ForegroundColor Yellow
if (Test-Path ".air.toml") {
    Write-Host "   ✅ .air.toml found" -ForegroundColor Green
} else {
    Write-Host "   ❌ .air.toml not found" -ForegroundColor Red
}

Write-Host ""

# Check if tmp directory exists
Write-Host "3. Checking tmp directory..." -ForegroundColor Yellow
if (Test-Path "tmp") {
    Write-Host "   ⚠️  tmp directory exists (will be used for builds)" -ForegroundColor Yellow
} else {
    Write-Host "   ✅ tmp directory will be created on first run" -ForegroundColor Green
}

Write-Host ""

# Test build
Write-Host "4. Testing build..." -ForegroundColor Yellow
$buildResult = go build -o tmp/test.exe . 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "   ✅ Build successful" -ForegroundColor Green
    Remove-Item tmp/test.exe -ErrorAction SilentlyContinue
} else {
    Write-Host "   ❌ Build failed" -ForegroundColor Red
    Write-Host "   $buildResult" -ForegroundColor Red
}

Write-Host ""

# Summary
Write-Host "═══════════════════════════════════════" -ForegroundColor Cyan
Write-Host "📋 Summary" -ForegroundColor Cyan
Write-Host "═══════════════════════════════════════" -ForegroundColor Cyan
Write-Host ""
Write-Host "Ready to use hot reload!" -ForegroundColor Green
Write-Host ""
Write-Host "Commands to try:" -ForegroundColor White
Write-Host "  .\neonex.exe serve --hot" -ForegroundColor Gray
Write-Host "  air" -ForegroundColor Gray
Write-Host ""
Write-Host "Make changes to any .go file and it will auto-rebuild!" -ForegroundColor Yellow
Write-Host ""
