# Redis Caching Implementation Guide

## 📌 Tổng quan

Redis caching đã được tích hợp vào Go Backend để tăng hiệu suất và giảm tải cho database PostgreSQL.

## 🎯 Các tính năng đã implement

### 1. Redis Client Setup

- File: `cache/redis.go`
- Kết nối Redis với cấu hình linh hoạt
- Hỗ trợ các thao tác: Set, Get, Delete, DeletePattern, Exists, SetExpiration
- Auto-reconnect và health check

### 2. Configuration Management

- File: `config/config.go`
- Thêm các biến môi trường:
  - `REDIS_HOST`: Địa chỉ Redis server (mặc định: localhost)
  - `REDIS_PORT`: Port Redis (mặc định: 6379)
  - `REDIS_PASSWORD`: Password (mặc định: trống)
  - `REDIS_DB`: Database number (mặc định: 0)

### 3. Caching Strategy

#### Products List Caching

- **Endpoint**: `GET /api/v1/products`
- **Cache Key Format**: `products:page:{page}:limit:{limit}:category:{category}`
- **TTL**: 5 phút
- **Cache Invalidation**: Khi create/update/delete product

#### Single Product Caching

- **Endpoint**: `GET /api/v1/products/:id`
- **Cache Key Format**: `product:{id}`
- **TTL**: 10 phút
- **Cache Invalidation**: Khi update/delete product cụ thể

### 4. Cache Headers

- `X-Cache: HIT` - Dữ liệu từ cache
- `X-Cache: MISS` - Dữ liệu từ database

## 🚀 Cách sử dụng

### Khởi động Redis

```bash
# Sử dụng Docker Compose
docker-compose up -d redis

# Kiểm tra Redis đang chạy
docker ps | grep redis
```

### Cấu hình Environment Variables

Cập nhật file `.env`:

```env
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### Chạy Application

```bash
go run main.go
```

Nếu Redis kết nối thành công, bạn sẽ thấy:

```
✅ Redis connected successfully
```

## 📊 Kiểm tra Cache

### Test Cache HIT/MISS

```bash
# Lần 1: MISS (lấy từ database)
curl -I http://localhost:8080/api/v1/products
# X-Cache: MISS

# Lần 2: HIT (lấy từ cache - trong vòng 5 phút)
curl -I http://localhost:8080/api/v1/products
# X-Cache: HIT
```

### Xem Cache trong Redis

```bash
# Vào Redis CLI
docker exec -it go_backend_redis redis-cli

# Xem tất cả keys
KEYS *

# Xem nội dung một key
GET product:1

# Xem TTL còn lại
TTL product:1

# Xóa một key
DEL product:1

# Xóa tất cả cache (cẩn thận!)
FLUSHALL
```

## 🔄 Cache Invalidation Flow

### Khi tạo product mới:

1. Product được lưu vào database
2. Tất cả cache có pattern `products:*` bị xóa
3. Response trả về không có cache

### Khi update product:

1. Product được cập nhật trong database
2. Cache `product:{id}` cụ thể bị xóa
3. Tất cả cache có pattern `products:*` bị xóa
4. Response trả về không có cache

### Khi delete product:

1. Product bị soft delete trong database
2. Cache `product:{id}` cụ thể bị xóa
3. Tất cả cache có pattern `products:*` bị xóa
4. Response trả về không có cache

## 📈 Performance Benefits

### Trước khi có Cache:

- Mỗi request đều query database
- Response time: ~50-100ms (tùy độ phức tạp query)
- Database load cao khi có nhiều concurrent requests

### Sau khi có Cache:

- Request đầu tiên query database và lưu cache
- Các request tiếp theo lấy từ cache
- Response time với cache: ~5-10ms
- Giảm 80-90% load cho database
- Tăng khả năng xử lý concurrent requests

## ⚙️ Cấu hình nâng cao

### Thay đổi TTL

Trong `controllers/product_controller.go`:

```go
// Thay đổi TTL cho products list (hiện tại: 5 phút)
cache.Set(cacheKey, response, 10*time.Minute) // Đổi thành 10 phút

// Thay đổi TTL cho single product (hiện tại: 10 phút)
cache.Set(cacheKey, product, 15*time.Minute) // Đổi thành 15 phút
```

### Redis Password Protection

Trong file `.env`:

```env
REDIS_PASSWORD=your-strong-password
```

Trong `docker-compose.yml`:

```yaml
redis:
  image: redis:7-alpine
  command: redis-server --requirepass your-strong-password
```

### Redis Persistence

Redis đã được cấu hình với AOF (Append Only File):

```yaml
redis:
  command: redis-server --appendonly yes
  volumes:
    - redis_data:/data
```

## 🛠️ Troubleshooting

### Redis không kết nối được

```bash
# Kiểm tra Redis container
docker ps | grep redis

# Xem logs
docker logs go_backend_redis

# Restart Redis
docker-compose restart redis
```

### Cache không được invalidate

```bash
# Xóa toàn bộ cache manually
docker exec -it go_backend_redis redis-cli FLUSHALL

# Hoặc xóa pattern cụ thể
docker exec -it go_backend_redis redis-cli
> KEYS products:*
> DEL products:page:1:limit:10:category:
```

### Application chạy không có Redis

Application sẽ vẫn chạy bình thường nếu Redis không available. Chỉ có warning:

```
⚠️ Warning: Failed to connect to Redis: ...
Continuing without cache...
```

## 📚 Best Practices

1. **TTL Selection**: Chọn TTL phù hợp với tần suất thay đổi dữ liệu
2. **Cache Key Naming**: Sử dụng naming convention rõ ràng và có pattern
3. **Cache Invalidation**: Luôn invalidate cache khi dữ liệu thay đổi
4. **Error Handling**: Graceful degradation khi Redis không available
5. **Monitoring**: Theo dõi hit rate và performance metrics

## 🔐 Security

- Sử dụng password cho Redis trong production
- Không expose Redis port ra ngoài (bind to localhost)
- Sử dụng Redis ACL để giới hạn quyền truy cập
- Regular backup Redis data
- Monitor Redis memory usage

## 📝 Next Steps

Có thể mở rộng thêm:

- [ ] Implement caching cho user authentication
- [ ] Add Redis Pub/Sub cho real-time updates
- [ ] Implement distributed caching cho multi-server
- [ ] Add cache warming strategy
- [ ] Implement cache versioning
- [ ] Add Redis Sentinel cho high availability
- [ ] Monitor cache metrics với Prometheus
