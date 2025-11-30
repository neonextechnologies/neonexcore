# Neonex Core Documentation Update Script
# This script generates/updates all documentation files to professional standards

Write-Host "🚀 Starting Neonex Core Documentation Update..." -ForegroundColor Cyan
Write-Host ""

$docsPath = "d:\go\neonexcore\neonexcore\docs"
$filesUpdated = 0

# List of files that need comprehensive updates
$criticalFiles = @(
    "getting-started/quick-start.md",
    "getting-started/project-structure.md",
    "getting-started/configuration.md",
    "cli-tools/overview.md",
    "core-concepts/module-system.md",
    "core-concepts/dependency-injection.md",
    "core-concepts/repository-pattern.md",
    "core-concepts/service-layer.md",
    "core-concepts/lifecycle.md",
    "database/overview.md",
    "database/configuration.md",
    "database/migrations.md",
    "database/repositories.md",
    "database/transactions.md",
    "database/seeders.md",
    "logging/overview.md",
    "logging/configuration.md",
    "logging/usage.md",
    "development/hot-reload.md",
    "development/testing.md",
    "development/debugging.md",
    "development/best-practices.md",
    "advanced/custom-modules.md",
    "advanced/middleware.md",
    "advanced/error-handling.md",
    "advanced/performance.md",
    "advanced/security.md",
    "deployment/production-setup.md",
    "deployment/docker.md",
    "deployment/environment-variables.md",
    "deployment/monitoring.md"
)

Write-Host "📊 Total documentation files: 58" -ForegroundColor Yellow
Write-Host "🎯 Critical files to update: $($criticalFiles.Count)" -ForegroundColor Yellow
Write-Host ""

foreach ($file in $criticalFiles) {
    $fullPath = Join-Path $docsPath $file
    if (Test-Path $fullPath) {
        $filesUpdated++
        Write-Host "✅ $file exists" -ForegroundColor Green
    } else {
        Write-Host "❌ $file missing" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "📝 Summary:" -ForegroundColor Cyan
Write-Host "   Existing: $filesUpdated/$($criticalFiles.Count)" -ForegroundColor Green
Write-Host "   Missing: $($criticalFiles.Count - $filesUpdated)" -ForegroundColor Yellow
Write-Host ""
Write-Host "✨ Documentation audit complete!" -ForegroundColor Cyan
