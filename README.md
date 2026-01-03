# Go Backend RESTful API

Đây là một ví dụ đầy đủ về RESTful API được viết bằng Golang với cấu trúc chuẩn, sử dụng Gin framework, GORM ORM, và PostgreSQL database.

## 🏗️ Cấu trúc dự án

```
go_backend/
├── config/              # Cấu hình ứng dụng
│   └── config.go
├── controllers/         # Xử lý logic cho các endpoint
│   ├── auth_controller.go
│   └── product_controller.go
├── database/           # Kết nối và migration database
│   └── database.go
├── middleware/         # Middleware cho authentication, logging, CORS
│   ├── auth.go
│   ├── cors.go
│   └── logger.go
├── models/            # Định nghĩa cấu trúc dữ liệu
│   ├── product.go
│   ├── response.go
│   └── user.go
├── routes/            # Định nghĩa các API endpoints
│   └── routes.go
├── utils/             # Các hàm tiện ích (JWT, validation, v.v.)
│   └── jwt.go
├── .env.example       # File mẫu environment variables
├── go.mod            # Quản lý dependencies
└── main.go           # Entry point của ứng dụng
```

## 🚀 Tính năng

### Authentication

- ✅ Đăng ký tài khoản (Register)
- ✅ Đăng nhập (Login)
- ✅ JWT Authentication
- ✅ Xem thông tin profile

### Product Management (CRUD)

- ✅ Lấy danh sách sản phẩm (có phân trang)
- ✅ Lấy chi tiết sản phẩm
- ✅ Tạo sản phẩm mới
- ✅ Cập nhật sản phẩm
- ✅ Xóa sản phẩm (soft delete)

### Middleware

- ✅ CORS middleware
- ✅ Logger middleware
- ✅ Authentication middleware
- ✅ Admin middleware

## 📋 Yêu cầu

- Go 1.21 hoặc cao hơn
- Docker và Docker Compose (để chạy PostgreSQL và pgAdmin)

## ⚙️ Cài đặt

### 1. Clone repository

```bash
cd d:\Coding\go_projects\go_backend
```

> **Lưu ý:** Đảm bảo Docker Desktop đã được cài đặt và đang chạy trên máy của bạn.

### 2. Cài đặt dependencies

```bash
go get .
```

### 3. Khởi động Database với Docker Compose

Chạy PostgreSQL và pgAdmin bằng Docker Compose:

```bash
docker-compose up -d
```

Kiểm tra containers đã chạy:

```bash
docker-compose ps
```

### 4. Cấu hình môi trường

Copy file `.env.example` thành `.env`:

```bash
copy .env.example .env  # Windows
# hoặc
cp .env.example .env    # Linux/Mac
```

File `.env` mặc định đã được cấu hình phù hợp với Docker Compose:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_backend_db
PORT=8080
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
APP_ENV=development
```

⚠️ **Lưu ý:** Đổi `JWT_SECRET` trong môi trường production!

### 5. Chạy ứng dụng

```bash
go run main.go
```

Server sẽ chạy tại: `http://localhost:8080`

### 6. Truy cập pgAdmin (Quản lý Database)

Mở trình duyệt và truy cập: `http://localhost:5050`

**Thông tin đăng nhập pgAdmin:**

- Email: `admin@admin.com`
- Password: `admin`

**Kết nối PostgreSQL trong pgAdmin:**

1. Click **Add New Server**
2. Tab **General**:
   - Name: `Go Backend DB` (tùy chọn)
3. Tab **Connection**:
   - Host: `postgres`
   - Port: `5432`
   - Database: `go_backend_db`
   - Username: `postgres`
   - Password: `postgres`
4. Click **Save**

### 7. Dừng các services

```bash
# Dừng containers
docker-compose down

# Dừng và xóa volumes (XÓA TOÀN BỘ DỮ LIỆU)
docker-compose down -v
```

## 📚 API Endpoints

### Health Check

```
GET /health
```

### Authentication

#### Đăng ký

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "password123",
  "full_name": "John Doe"
}
```

#### Đăng nhập

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

#### Xem profile (Cần authentication)

```http
GET /api/v1/profile
Authorization: Bearer <your_jwt_token>
```

### Products (Cần authentication)

#### Lấy danh sách sản phẩm

```http
GET /api/v1/products?page=1&limit=10&category=electronics
Authorization: Bearer <your_jwt_token>
```

#### Lấy chi tiết sản phẩm

```http
GET /api/v1/products/1
Authorization: Bearer <your_jwt_token>
```

#### Tạo sản phẩm mới

```http
POST /api/v1/products
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
  "name": "iPhone 15 Pro",
  "description": "Latest iPhone model",
  "price": 999.99,
  "stock": 100,
  "category": "electronics",
  "sku": "IPH15PRO",
  "image_url": "https://example.com/image.jpg"
}
```

#### Cập nhật sản phẩm

```http
PUT /api/v1/products/1
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
  "name": "iPhone 15 Pro Max",
  "price": 1199.99,
  "stock": 50
}
```

#### Xóa sản phẩm

```http
DELETE /api/v1/products/1
Authorization: Bearer <your_jwt_token>
```

## 🔐 Authentication

API sử dụng JWT (JSON Web Tokens) cho authentication. Sau khi login hoặc register thành công, bạn sẽ nhận được một token. Sử dụng token này trong header `Authorization`:

```
Authorization: Bearer <your_jwt_token>
```

Token có thời hạn 24 giờ.

## 📖 Giải thích cấu trúc

### 1. **main.go**

- Entry point của ứng dụng
- Load environment variables
- Kết nối database
- Chạy auto migration
- Setup routes
- Khởi động server

### 2. **config/**

- Quản lý cấu hình từ environment variables
- Centralized configuration

### 3. **database/**

- Xử lý kết nối database
- Auto migration cho models
- Singleton pattern cho database instance

### 4. **models/**

- Định nghĩa struct cho database tables
- Request/Response DTOs
- Validation tags

### 5. **controllers/**

- Xử lý business logic
- Validate request data
- Trả về response

### 6. **routes/**

- Định nghĩa API endpoints
- Group routes theo version (v1, v2, ...)
- Apply middleware cho các route groups

### 7. **middleware/**

- **Auth**: Xác thực JWT token
- **CORS**: Xử lý Cross-Origin requests
- **Logger**: Log request details

### 8. **utils/**

- Các helper functions
- JWT generation và validation

## 🛠️ Technologies

- **Gin** - Web framework
- **GORM** - ORM library
- **PostgreSQL** - Database
- **JWT** - Authentication
- **bcrypt** - Password hashing
- **godotenv** - Environment variables
- **Docker & Docker Compose** - Database containerization

## 📝 Best Practices được áp dụng

1. ✅ **Separation of Concerns** - Tách biệt logic thành các layers
2. ✅ **RESTful API Design** - Tuân thủ chuẩn REST
3. ✅ **Environment Variables** - Không hardcode sensitive data
4. ✅ **Password Hashing** - Sử dụng bcrypt
5. ✅ **JWT Authentication** - Stateless authentication
6. ✅ **Soft Delete** - Không xóa vĩnh viễn data quan trọng
7. ✅ **Pagination** - Hỗ trợ phân trang cho danh sách
8. ✅ **Validation** - Validate input data
9. ✅ **Error Handling** - Xử lý lỗi consistent
10. ✅ **Logging** - Log requests và errors

## 🧪 Testing với cURL

### Register

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser",
    "password": "password123",
    "full_name": "Test User"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Tạo Product

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "Test Product",
    "description": "A test product",
    "price": 99.99,
    "stock": 10,
    "category": "test",
    "sku": "TEST001"
  }'
```

## 📚 Học hỏi thêm

### Các concepts quan trọng:

1. **RESTful API** - Chuẩn thiết kế API
2. **MVC Pattern** - Model-View-Controller
3. **Middleware** - Request/Response processing
4. **ORM** - Object-Relational Mapping
5. **JWT** - JSON Web Tokens
6. **CRUD Operations** - Create, Read, Update, Delete

### Tài liệu tham khảo:

- [Gin Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [Go Documentation](https://golang.org/doc/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## 🔧 Mở rộng

Để phát triển thêm, bạn có thể:

- ✅ Docker containerization (Database)
- Thêm unit tests
- Implement refresh token
- Thêm role-based access control
- Upload files/images
- Email verification
- Password reset
- Rate limiting
- API documentation với Swagger
- Dockerize Go application
- CI/CD pipeline

## ⚠️ Lưu ý bảo mật

1. **Đổi JWT_SECRET** trong production
2. **Sử dụng HTTPS** trong production
3. **Không commit file .env** vào git
4. **Validate tất cả input** từ client
5. **Implement rate limiting** để chống brute force
6. **Sử dụng prepared statements** (GORM đã hỗ trợ)

---

**Happy Coding! 🚀**

Nếu có câu hỏi, hãy tạo issue hoặc liên hệ qua email.
