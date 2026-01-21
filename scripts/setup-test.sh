#!/bin/bash

# Script to setup test environment
# This script will:
# 1. Install dependencies
# 2. Setup test database
# 3. Run tests

set -e

echo "🔧 Setting up test environment..."

# Install dependencies
echo "📦 Installing dependencies..."
go mod download
go mod verify

# Install development tools
echo "🛠️  Installing development tools..."
if ! command -v golangci-lint &> /dev/null; then
    echo "Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
else
    echo "golangci-lint already installed"
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "⚠️  Docker is not running. Please start Docker to run tests with database."
    exit 1
fi

# Start test services (PostgreSQL and Redis)
echo "🐳 Starting test services..."
docker-compose up -d postgres redis

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
until docker-compose exec -T postgres pg_isready -U postgres; do
    sleep 1
done

# Wait for Redis to be ready
echo "⏳ Waiting for Redis to be ready..."
until docker-compose exec -T redis redis-cli ping | grep -q PONG; do
    sleep 1
done

echo "✅ Test environment is ready!"
echo ""
echo "You can now run tests with:"
echo "  make test          - Run all tests"
echo "  make test-coverage - Run tests with coverage"
echo "  make lint          - Run linter"
echo ""
