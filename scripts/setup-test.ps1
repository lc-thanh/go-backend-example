# PowerShell script to setup test environment
# This script will:
# 1. Install dependencies
# 2. Setup test database
# 3. Run tests

$ErrorActionPreference = "Stop"

Write-Host "🔧 Setting up test environment..." -ForegroundColor Cyan

# Install dependencies
Write-Host "📦 Installing dependencies..." -ForegroundColor Yellow
go mod download
go mod verify

# Install development tools
Write-Host "🛠️  Installing development tools..." -ForegroundColor Yellow
$golangciLintInstalled = Get-Command golangci-lint -ErrorAction SilentlyContinue
if (-not $golangciLintInstalled) {
    Write-Host "Installing golangci-lint..." -ForegroundColor Yellow
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
} else {
    Write-Host "golangci-lint already installed" -ForegroundColor Green
}

# Check if Docker is running
try {
    docker info | Out-Null
} catch {
    Write-Host "⚠️  Docker is not running. Please start Docker to run tests with database." -ForegroundColor Red
    exit 1
}

# Start test services (PostgreSQL and Redis)
Write-Host "🐳 Starting test services..." -ForegroundColor Yellow
docker-compose up -d postgres redis

# Wait for PostgreSQL to be ready
Write-Host "⏳ Waiting for PostgreSQL to be ready..." -ForegroundColor Yellow
$maxRetries = 30
$retries = 0
while ($retries -lt $maxRetries) {
    try {
        docker-compose exec -T postgres pg_isready -U postgres 2>&1 | Out-Null
        if ($LASTEXITCODE -eq 0) {
            break
        }
    } catch {
        # Continue waiting
    }
    Start-Sleep -Seconds 1
    $retries++
}

if ($retries -eq $maxRetries) {
    Write-Host "❌ PostgreSQL failed to start" -ForegroundColor Red
    exit 1
}

# Wait for Redis to be ready
Write-Host "⏳ Waiting for Redis to be ready..." -ForegroundColor Yellow
$maxRetries = 30
$retries = 0
while ($retries -lt $maxRetries) {
    try {
        $result = docker-compose exec -T redis redis-cli ping 2>&1
        if ($result -match "PONG") {
            break
        }
    } catch {
        # Continue waiting
    }
    Start-Sleep -Seconds 1
    $retries++
}

if ($retries -eq $maxRetries) {
    Write-Host "❌ Redis failed to start" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "✅ Test environment is ready!" -ForegroundColor Green
Write-Host ""
Write-Host "You can now run tests with:" -ForegroundColor Cyan
Write-Host "  make test          - Run all tests" -ForegroundColor White
Write-Host "  make test-coverage - Run tests with coverage" -ForegroundColor White
Write-Host "  make lint          - Run linter" -ForegroundColor White
Write-Host ""
