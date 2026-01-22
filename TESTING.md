# Testing Documentation

Dự án này đã được thiết lập với một bộ test toàn diện để đảm bảo chất lượng code và tích hợp CI/CD.

## Cấu trúc Test

Các test được tổ chức theo từng package:

```
├── config/
│   └── config_test.go          # Tests cho configuration loading
├── controllers/
│   └── (integration tests)     # Tests cho API controllers
├── middleware/
│   ├── auth_test.go           # Tests cho authentication middleware
│   ├── cors_test.go           # Tests cho CORS middleware
│   └── logger_test.go         # Tests cho logging middleware
├── models/
│   ├── user_test.go           # Tests cho User model
│   ├── product_test.go        # Tests cho Product model
│   └── response_test.go       # Tests cho API response models
└── utils/
    └── jwt_test.go            # Tests cho JWT utilities
```

## Chạy Tests

### Chạy tất cả tests

```bash
make test
# hoặc
go test -v -race ./...
```

### Chạy tests với coverage

```bash
make test-coverage
# hoặc
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -html=coverage.out -o coverage.html
```

### Chạy tests cho một package cụ thể

```bash
go test -v ./utils
go test -v ./middleware
go test -v ./models
```

### Chạy một test cụ thể

```bash
go test -v -run TestGenerateJWT ./utils
go test -v -run TestAuthMiddleware ./middleware
```

### Chạy tests với verbose output

```bash
make test-verbose
```

### Chạy benchmarks

```bash
make bench
```

## CI/CD Pipeline

Dự án sử dụng GitHub Actions để tự động chạy tests và build code.

### Workflow CI/CD

Pipeline CI/CD được định nghĩa trong `.github/workflows/ci.yml` và bao gồm:

1. **Test Job**
   - Chạy trên Ubuntu latest
   - Khởi động PostgreSQL và Redis containers
   - Chạy tất cả tests với race detector
   - Tạo coverage report
   - Upload coverage lên Codecov

2. **Lint Job**
   - Chạy golangci-lint để kiểm tra code quality
   - Thực thi các rules được định nghĩa trong `.golangci.yml`

3. **Build Job**
   - Build binary từ source code
   - Upload artifact để lưu trữ

4. **Docker Job** (chỉ chạy khi push lên main)
   - Build Docker image
   - Push image lên Docker Hub

### Trigger CI/CD

CI/CD pipeline sẽ tự động chạy khi:

- Push code lên branch `main` hoặc `develop`
- Tạo Pull Request vào `main` hoặc `develop`

## Dependencies

Các dependencies cần thiết cho testing:

- `github.com/stretchr/testify` - Assertion library
- `github.com/gin-gonic/gin` - Web framework
- Testing database: PostgreSQL
- Testing cache: Redis

## Environment Variables cho Testing

Khi chạy tests local, các biến môi trường sau được sử dụng:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_backend_test
JWT_SECRET=test-secret-key
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

## Code Coverage

Để xem coverage report:

1. Chạy tests với coverage:

   ```bash
   make test-coverage
   ```

2. Mở file `coverage.html` trong browser:

   ```bash
   # Windows
   start coverage.html

   # Linux/Mac
   open coverage.html
   ```

## Linting

### Chạy linter local

```bash
make lint
# hoặc
golangci-lint run ./...
```

### Cài đặt golangci-lint

```bash
make install-tools
# hoặc
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Cấu hình Linter

Linter được cấu hình trong file `.golangci.yml` với các rules:

- `errcheck` - Kiểm tra unchecked errors
- `gosimple` - Đơn giản hóa code
- `govet` - Phát hiện suspicious constructs
- `staticcheck` - Static analysis
- `gosec` - Security checks
- `gocritic` - Opinionated linter
- Và nhiều rules khác...

## Best Practices

1. **Viết tests cho code mới**
   - Mọi function/method mới nên có unit test
   - Integration tests cho API endpoints

2. **Chạy tests trước khi commit**

   ```bash
   make ci  # Chạy full CI pipeline local
   ```

3. **Maintain code coverage**
   - Mục tiêu: ít nhất 70% coverage
   - Review coverage report thường xuyên

4. **Fix linter warnings**
   - Chạy linter và fix warnings trước khi push
   - Không disable linter rules trừ khi cần thiết

## Troubleshooting

### Tests fail do database connection

Đảm bảo PostgreSQL đang chạy:

```bash
# Với Docker
docker-compose up -d postgres

# Kiểm tra connection
psql -h localhost -U postgres -d go_backend_test
```

### Tests fail do Redis connection

Đảm bảo Redis đang chạy:

```bash
# Với Docker
docker-compose up -d redis

# Kiểm tra connection
redis-cli ping
```

### Race condition detected

Nếu phát hiện race condition:

1. Xem output chi tiết từ `-race` flag
2. Fix data races bằng cách sử dụng sync primitives
3. Re-run tests để verify fix

## Contributing

Khi contribute code mới:

1. Viết tests cho code mới
2. Đảm bảo tất cả tests pass:
   ```bash
   make test
   ```
3. Đảm bảo code pass linter:
   ```bash
   make lint
   ```
4. Kiểm tra coverage không giảm
5. Tạo Pull Request

CI/CD pipeline sẽ tự động verify tất cả requirements trên.
