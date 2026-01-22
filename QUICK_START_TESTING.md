# Quick Start Guide - Testing & CI/CD

## ⚡ Chạy Tests Ngay

### Bước 1: Install Dependencies

```bash
go mod download
go mod tidy
```

### Bước 2: Start Test Services (PostgreSQL & Redis)

```bash
docker-compose up -d postgres redis
```

### Bước 3: Run Tests

```bash
# Tất cả tests
go test -v ./...

# Hoặc dùng Makefile
make test
```

## 📊 Coverage Report

```bash
# Generate coverage
make test-coverage

# Mở coverage.html trong browser
# Windows:
start coverage.html

# Linux/Mac:
open coverage.html
```

## 🔍 Code Quality Check

```bash
# Install linter (one-time)
make install-tools

# Run linter
make lint

# Format code
make fmt
```

## 🐳 Docker Commands

```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs -f

# Restart services
docker-compose restart
```

## 🚀 Run Application

```bash
# Development mode
go run main.go

# Or using Makefile
make run

# Build binary
make build
./bin/server
```

## ✅ Pre-Commit Checklist

Trước khi commit code:

```bash
# 1. Format code
make fmt

# 2. Run linter
make lint

# 3. Run tests
make test

# Or run all at once:
make ci
```

## 🔄 CI/CD Workflow

### Automatic Triggers:

- ✅ Push to `main` or `develop` → Full CI pipeline
- ✅ Pull Request → Tests + Linting + Security scan

### Manual Check:

```bash
# Check if code will pass CI
make ci
```

## 📁 Test Files Location

```
├── config/config_test.go
├── middleware/
│   ├── auth_test.go
│   ├── cors_test.go
│   └── logger_test.go
├── models/
│   ├── user_test.go
│   ├── product_test.go
│   └── response_test.go
└── utils/jwt_test.go
```

## 🆘 Troubleshooting

### "Cannot connect to database"

```bash
# Ensure PostgreSQL is running
docker-compose up -d postgres

# Check status
docker-compose ps
```

### "Redis connection failed"

```bash
# Ensure Redis is running
docker-compose up -d redis

# Test connection
docker exec -it go_backend_redis redis-cli ping
# Should return: PONG
```

### "Tests fail with race condition"

```bash
# Run tests without race detector
go test ./...

# Fix race conditions in code
# Then re-enable: go test -race ./...
```

### "Linter errors"

```bash
# Auto-fix what can be fixed
golangci-lint run --fix ./...

# Format code
gofmt -w .

# Fix imports
goimports -w .
```

## 💡 Tips

1. **Fast testing during development:**

   ```bash
   # Test specific package
   go test -v ./utils
   go test -v ./middleware
   ```

2. **Watch mode (using external tool):**

   ```bash
   # Install gotestsum
   go install gotest.tools/gotestsum@latest

   # Watch and re-run tests
   gotestsum --watch
   ```

3. **Verbose output:**

   ```bash
   make test-verbose
   ```

4. **Check what CI will do:**
   ```bash
   # Local CI simulation
   make ci
   ```

## 📚 More Info

- Full testing guide: [TESTING.md](TESTING.md)
- Implementation summary: [TEST_IMPLEMENTATION_SUMMARY.md](TEST_IMPLEMENTATION_SUMMARY.md)
- Redis caching guide: [REDIS_CACHING.md](REDIS_CACHING.md)

---

**Need help?** Create an issue or check the documentation!
