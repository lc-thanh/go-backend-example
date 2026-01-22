# Jenkins CI/CD Configuration

Hướng dẫn cấu hình Jenkins CI/CD cho dự án Go Backend với 2 kịch bản riêng biệt:

- **CI (Continuous Integration)**: Validate code trên tất cả nhánh
- **CD (Continuous Deployment)**: Deploy khi có thay đổi trên `main` hoặc khi có tag mới

## Yêu cầu

- Jenkins 2.x trở lên
- Docker đã được cài đặt trên Jenkins server
- Plugins cần thiết:
  - **Docker Pipeline** (bắt buộc)
  - **HTML Publisher** (cho coverage reports)
  - **Pipeline** (bắt buộc)
  - **Git** (bắt buộc)
  - **Credentials Binding** (cho Docker registry)
  - **Input Step** (cho deployment approval)

**Lưu ý quan trọng**:

- Pipeline sử dụng Docker containers để chạy tất cả tasks (không cần cài Go trên Jenkins server)
- Go version: `1.23-alpine`
- Environment files phải được đặt sẵn trên Jenkins node tại `/app/env/`

## Cấu hình Jenkins

### 0. Cài đặt Docker trên Jenkins Agent (node)

Đảm bảo Docker đã được cài đặt và Jenkins user có quyền sử dụng:

```bash
# Thêm Jenkins user vào docker group
sudo usermod -aG docker jenkins

# Verify
sudo -u jenkins docker ps
```

### 1. Cấu hình Gitlab Connection

1. **Manage Jenkins** → **System**
2. Tìm phần **GitLab**
3. **GitLab connections**:
   - **Connection name**: `gitlab server`
   - **GitLab host URL**: `https://gitlab.com` (hoặc self-hosted GitLab URL)
   - **Credentials**: Add GitLab Personal Access Token
     - Kind: **GitLab API token**
     - API token: (tạo từ GitLab: Settings → Access Tokens)
     - Scopes: `api`
4. Click **Test Connection**
5. Save

### 2. Cấu hình Git Credentials (nếu private repository)

1. Vào **Manage Jenkins** → **Credentials**
2. Click **(global)** → **Add Credentials**
3. Chọn một trong hai:

**Option A: Username with password**

- Username: Git username
- Password: Personal Access Token hoặc password
- ID: `jenkins-gitlab-user-account`

**Option B: SSH Username with private key**

- Kind: **SSH Username with private key**
- Username: `git`
- Private Key: Paste SSH private key
- ID: `git-ssh-key`

### 3. Cấu hình Docker Credentials

**Bắt buộc** cho việc push Docker images lên registry:

1. Vào **Manage Jenkins** → **Credentials**
2. Click **(global)** → **Add Credentials**
3. Chọn **Username with password**
4. Nhập:
   - Username: Docker Hub username (hoặc registry username)
   - Password: Docker Hub access token/password
   - ID: `docker-hub-credentials` (phải khớp với Jenkinsfile)
5. Save

> **Lưu ý**: ID phải là `docker-hub-credentials` hoặc cập nhật biến `DOCKER_CREDENTIALS_ID` trong Jenkinsfile.

### 4. Cấu hình Environment Files

Pipeline yêu cầu các file .env được đặt sẵn trên Jenkins node:

```bash
# Tạo thư mục env trên node
sudo mkdir -p /app/env
sudo chown jenkins:jenkins /app/env

# Tạo các file environment
# File cho testing (CI stages)
sudo nano /app/env/.env.test

# File cho production deployment
sudo nano /app/env/.env.production

# File cho staging deployment (optional)
sudo nano /app/env/.env.staging
```

**Ví dụ nội dung .env file**:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_backend_db

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Server Configuration
PORT=8080

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Environment
APP_ENV=development
```

### 5. Tạo Pipeline Job

1. Click **New Item**
2. Nhập tên job (vd: `go-backend-ci-cd`)
3. Chọn **Pipeline**
4. Click **OK**

#### Cấu hình Pipeline:

**General:**

- ☐ Discard old builds: Keep last 10 builds (đã cấu hình trong Jenkinsfile)

**Build Triggers:**

Chọn **Build when a change is pushed to GitLab...**, sau đó chọn những tùy chọn sau:

- ✓ **Push Events**:
- ✓ **Opened Merge Request Events**
- ✓ **Build only if new commits were pushed to Merge Request**
- ✓ **Accepted Merge Request Events**

Xuống phần **Advanced**, nhấn **Generate** để tạo Secret token, copy để cấu hình Webhook bên Gitlab

**Pipeline:**

- Definition: **Pipeline script from SCM**
- SCM: **Git**
- Repository URL: `http://gitlab.local/go_backend/backend_app.git`
- Credentials: (chọn Git credentials nếu private repo)
- Branches to build:
  - **`*/*`** - Build tất cả branches (CI sẽ chạy cho tất cả)
  - Hoặc cụ thể: `*/main`, `*/develop`, `*/feature/*`
- Script Path: `Jenkinsfile`

**Additional Behaviours** (Optional):

- Add: **Clean before checkout** (đảm bảo workspace sạch)

### 6. Cấu hình Webhooks (GitLab)

#### Cho GitLab:

1. Vào repository → **Settings** → **Webhooks**
2. Add webhook:
   - **URL**: `http://your-jenkins-server/project/go-backend-ci-cd`
   - **Secret token**: nhập Secret token đã generate trước đó
   - **Trigger**:
     - ✓ Push events
     - ✓ Tag push events
     - ✓ Merge request events (optional)
3. Click **Add webhook**
4. Test webhook: Click **Test** → **Push events**

### 7. Cấu hình Deployment Approvers

Pipeline yêu cầu approval trước khi deploy. Cấu hình users có quyền approve:

1. Vào **Manage Jenkins** → **Configure Global Security**
2. Đảm bảo **Authorization** được bật
3. Trong Jenkinsfile, update `submitter` parameter:
   ```groovy
   submitter: 'admin,deployer,your-username'
   ```

## Cấu hình Jenkinsfile

Cần cập nhật các biến sau trong file `Jenkinsfile`:

```groovy
environment {
    // Docker registry configuration
    DOCKER_REGISTRY = 'docker.io'  // Hoặc: 'gcr.io', 'your-registry.com'
    DOCKER_IMAGE_NAME = 'your-dockerhub-username/go-backend'  // ⚠️ PHẢI ĐỔI
    DOCKER_CREDENTIALS_ID = 'docker-hub-credentials'  // Phải khớp với Jenkins credentials ID

    // Go version
    GO_VERSION = '1.23'  // Có thể thay đổi nếu cần

    // Environment files directory
    ENV_DIR = '/app/env'  // Phải khớp với đường dẫn trên Jenkins node
}
```

**Các biến quan trọng**:

- `DOCKER_IMAGE_NAME`: Tên Docker image, format: `registry-username/image-name`
- `DOCKER_CREDENTIALS_ID`: ID của credentials trong Jenkins (mặc định: `docker-hub-credentials`)
- `ENV_DIR`: Đường dẫn chứa file .env trên Jenkins server (mặc định: `/app/env`)

## Pipeline Stages

### CI Stages (Chạy trên TẤT CẢ branches)

Pipeline tự động chạy các stages sau cho mọi branch:

#### 1. **Checkout**

- Lấy source code từ Git repository
- Hiển thị thông tin branch và commit

#### 2. **Setup Environment**

- Kiểm tra Go environment trong Docker
- Verify Go version và configuration

#### 3. **Download Dependencies**

- Tải các Go modules
- Sử dụng `go mod download`

#### 4. **Code Quality** (Parallel)

Chạy song song 3 checks:

- **Format Check**: Kiểm tra code formatting với `gofmt`
- **Go Vet**: Phân tích code với `go vet`
- **Lint**: Chạy `golangci-lint` với timeout 10 phút

#### 5. **Unit Tests**

- Chạy tests với coverage: `go test -v -race -coverprofile=coverage.out`
- **Yêu cầu coverage tối thiểu: 70%** (build sẽ fail nếu dưới 70%)
- Load environment từ `/app/env/.env.test` nếu có
- Generate HTML coverage report
- Archive coverage artifacts

#### 6. **Security Scan**

- Scan bảo mật với `gosec`
- Generate JSON report
- Archive security report

#### 7. **Build Binary**

- Build Go binary với optimizations (`-ldflags="-s -w"`)
- Inject version và build time vào binary
- Output: `bin/go-backend`
- Archive binary artifact

---

### CD Stages (Chỉ chạy cho `main` branch hoặc tags)

Các stages sau chỉ chạy khi:

- Branch là `main` (nhánh chính)
- Hoặc có tag mới được push (format: `refs/tags/*`)

#### 8. **Build Docker Image**

- Build Docker image từ Dockerfile
- Tag image với: `${BRANCH_NAME}-${COMMIT_SHORT}`
- Nếu là `main` branch: tag thêm `latest`

#### 9. **Test Docker Image**

- Start container với port mapping (8081:8080)
- Load environment từ `/app/env/.env.test`
- Kiểm tra container có chạy được không
- Test health endpoint (nếu có)
- Cleanup container sau khi test

#### 10. **Push Docker Image**

- Push image lên Docker registry
- Sử dụng credentials đã cấu hình
- Push cả tag `latest` nếu là main branch

#### 11. **Deploy Approval** ⚠️ **YÊU CẦU XÁC NHẬN**

**Stage này dừng lại và chờ approval từ user!**

- Timeout: 30 phút
- Chỉ users trong danh sách `submitter` mới approve được
- Có 2 options:
  - **DEPLOY_ENVIRONMENT**: Chọn `production` hoặc `staging`
  - **DEPLOY_NOTES**: Ghi chú về deployment

Để approve:

1. Vào Jenkins build đang chờ
2. Click **Paused for Input**
3. Chọn environment và điền notes
4. Click **Deploy** để tiếp tục hoặc **Abort** để hủy

#### 12. **Deploy**

- Deploy application lên environment đã chọn
- Load environment từ `/app/env/.env.{environment}`
- **Cần cấu hình deployment method** (xem phần Deployment Strategy bên dưới)

#### 13. **Post-Deploy Verification**

- Verify deployment thành công
- Chạy smoke tests (nếu có)
- Kiểm tra health endpoints

## Deployment Strategy

### Branch/Tag Strategy

| Branch/Tag    | CI  | Build Image | Deploy | Approval Required |
| ------------- | --- | ----------- | ------ | ----------------- |
| `feature/*`   | ✅  | ❌          | ❌     | -                 |
| `develop`     | ✅  | ❌          | ❌     | -                 |
| `main`        | ✅  | ✅          | ✅     | **✅ Required**   |
| `refs/tags/*` | ✅  | ✅          | ✅     | **✅ Required**   |

### Cấu hình Deployment Method

Trong stage `Deploy` của Jenkinsfile (dòng ~430), chọn một trong các phương thức sau:

#### Option 1: Docker Compose (Đơn giản nhất)

```groovy
sh """
    # Copy docker-compose file
    scp docker-compose.yml user@server:/app/

    # Deploy
    ssh user@server "
        cd /app
        docker-compose down
        docker-compose pull
        docker-compose up -d
    "
"""
```

#### Option 2: Docker Run (Đơn giản, phù hợp cho 1 container)

```groovy
sh """
    # Stop và remove container cũ
    ssh user@server "
        docker stop go-backend-\${deployEnv} || true
        docker rm go-backend-\${deployEnv} || true

        # Pull image mới
        docker pull ${DOCKER_IMAGE_NAME}:${IMAGE_TAG}

        # Run container mới
        docker run -d \\
            --name go-backend-\${deployEnv} \\
            -p 8080:8080 \\
            --env-file ${envFile} \\
            --restart unless-stopped \\
            ${DOCKER_IMAGE_NAME}:${IMAGE_TAG}
    "
"""
```

#### Option 3: Docker Swarm (Cho production scale)

```groovy
sh """
    docker stack deploy \\
        -c docker-compose.\${deployEnv}.yml \\
        go-backend-stack
"""
```

#### Option 4: Kubernetes (Enterprise)

```groovy
sh """
    # Update image
    kubectl set image deployment/go-backend \\
        go-backend=${DOCKER_IMAGE_NAME}:${IMAGE_TAG} \\
        -n \${deployEnv}

    # Wait for rollout
    kubectl rollout status deployment/go-backend -n \${deployEnv}

    # Verify pods
    kubectl get pods -n \${deployEnv} -l app=go-backend
"""
```

**Lưu ý**: Uncomment và sửa đổi deployment method phù hợp trong Jenkinsfile!

### Post-Deploy Verification

Cập nhật verification checks trong stage `Post-Deploy Verification` (dòng ~466):

```groovy
sh '''
    # Health check
    response=$(curl -s -o /dev/null -w "%{http_code}" https://your-app.com/health)
    if [ $response -ne 200 ]; then
        echo "❌ Health check failed with status: $response"
        exit 1
    fi
    echo "✅ Health check passed"

    # Smoke tests
    # curl -f https://your-app.com/api/version
    # curl -f https://your-app.com/api/status

    # Database migration check (if applicable)
    # kubectl exec -n production deployment/go-backend -- ./bin/migrate status

    echo "✅ All verification checks passed"
'''
```

## Environment Files Structure

Cấu trúc file .env trên Jenkins node (`/app/env/`):

```
/app/env/
├── .env.test          # Cho testing/CI stages
├── .env.production    # Cho production deployment
└── .env.staging       # Cho staging deployment (optional)
```

**Template .env file**:

```env
# Application
APP_ENV=production
APP_PORT=8080
APP_DEBUG=false

# Database
DB_HOST=your-db-host.com
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secure-password
DB_NAME=go_backend
DB_SSL_MODE=require

# Redis
REDIS_HOST=your-redis-host.com
REDIS_PORT=6379
REDIS_PASSWORD=redis-password
REDIS_DB=0

# JWT
JWT_SECRET=your-very-long-secret-key-here
JWT_EXPIRY=24h

# External Services (if any)
# AWS_ACCESS_KEY_ID=xxx
# AWS_SECRET_ACCESS_KEY=xxx
```

## Viewing Results

Sau khi build chạy xong, xem kết quả tại:

### Build Dashboard

- **Build Status**: Màn hình chính của job hiển thị:
  - ✅ Success / ❌ Failed / ⚠️ Unstable
  - Build number và duration
  - Last successful build, last failed build

### Artifacts & Reports

Click vào build number → Xem các artifacts:

1. **Coverage Report** (HTML Publisher)
   - Click **Coverage Report** tab
   - Xem coverage chi tiết theo từng file/function
   - Minimum required: 70%

2. **Build Artifacts**
   - `coverage.out` - Raw coverage data
   - `gosec-report.json` - Security scan results
   - `bin/go-backend` - Binary executable

3. **Console Output**
   - Click **Console Output**
   - Xem logs chi tiết của từng stage
   - Debug errors nếu build failed

### Stage View

- Pipeline visualization hiển thị:
  - Các stages đã chạy
  - Thời gian mỗi stage
  - Stage nào passed/failed
  - Parallel stages (Code Quality)

## Pipeline Behavior Summary

### Tất cả branches (`*/*`):

```
✅ Checkout
✅ Setup Environment
✅ Download Dependencies
✅ Code Quality (Format + Vet + Lint)
✅ Unit Tests (với coverage ≥70%)
✅ Security Scan
✅ Build Binary
```

### Chỉ `main` branch và tags:

```
Tất cả CI stages ở trên +
✅ Build Docker Image
✅ Test Docker Image
✅ Push Docker Image
⏸️ Deploy Approval (CHỜ USER APPROVE)
✅ Deploy (sau khi approve)
✅ Post-Deploy Verification
```

## Best Practices

````

## Troubleshooting

### Issue: GitLab webhook không trigger build

**Solution**:

1. Kiểm tra webhook URL đúng format
2. Check Jenkins firewall cho phép GitLab truy cập
3. Test webhook trong GitLab Settings
4. Xem GitLab webhook logs: Settings → Webhooks → Recent Deliveries
5. Check Jenkins logs: Manage Jenkins → System Log

### Issue: Authentication failed với GitLab

**Solution**:

1. Verify GitLab API token còn valid
2. Check token scopes: `api`, `read_repository`
3. Test connection: Manage Jenkins → Configure System → GitLab → Test Connection

### Issue: Docker image pull bị chậm hoặc timeout

**Solution**:

- Sử dụng Docker registry mirror
- Pre-pull image: `docker pull golang:1.23-alpine`
- Tăng timeout trong Jenkins global configuration

### Issue: Docker permission denied

**Solution**: Thêm Jenkins user vào docker group:

```bash
sudo usermod -aG docker jenkins
sudo systemctl restart jenkins
````

### Issue: Tests fail do thiếu database

**Solution**:

- Cấu hình test database trong Jenkins
- Hoặc dùng Docker để tạo test database trong pipeline
- Hoặc skip integration tests: `go test -short ./...`

### Issue: Golangci-lint quá chậm

**Solution**: Thêm timeout trong Jenkinsfile:

```groovy
golangci-lint run --timeout 10m
```

## Best Practices

1. **Environment Files Security**:
   - Không commit file .env vào Git
   - Chỉ lưu trên Jenkins server với quyền restricted
   - Sử dụng Jenkins Credentials Plugin cho sensitive data
   - Rotate secrets định kỳ

2. **Docker Image Management**:
   - Tag images với commit hash để dễ rollback
   - Cleanup old images định kỳ (đã cấu hình trong pipeline)
   - Sử dụng multi-stage builds để giảm image size
   - Scan images với security tools

3. **Testing Strategy**:
   - Maintain coverage ≥70%
   - Chạy tests nhanh (CI) trước, tests chậm (integration) sau
   - Mock external services trong tests
   - Sử dụng test containers nếu cần database

4. **Deployment Safety**:
   - Luôn yêu cầu approval cho production
   - Test deployment trong staging trước
   - Có rollback plan
   - Monitor metrics sau deploy

5. **Pipeline Optimization**:
   - Parallel stages cho code quality checks
   - Cache dependencies khi có thể
   - Cleanup workspace sau build
   - Keep build history hợp lý (10 builds)

6. **Monitoring & Alerts**:
   - Enable build notifications (Slack/Email)
   - Track build duration trends
   - Alert khi coverage giảm
   - Monitor deployment frequency

## Troubleshooting

### Issue: Docker permission denied

**Solution**: Thêm Jenkins user vào docker group:

```bash
sudo usermod -aG docker jenkins
sudo systemctl restart jenkins

# Verify
sudo -u jenkins docker ps
```

### Issue: Tests fail do thiếu database

**Solution**:

1. **Sử dụng test .env file** với in-memory DB hoặc mock
2. **Docker test database** trong pipeline:

```groovy
sh '''
    docker run -d --name test-db -e POSTGRES_PASSWORD=test postgres:15
    # Run tests
    docker rm -f test-db
'''
```

3. **Skip integration tests**: `go test -short ./...`

### Issue: Coverage không đạt 70%

**Solution**:

- Viết thêm tests cho untested code
- Hoặc tạm thời giảm threshold trong Jenkinsfile:

```groovy
if [ $(echo "$coverage < 70" | bc -l) -eq 1 ]; then
    # Thay 70 thành giá trị thấp hơn
```

### Issue: Golangci-lint quá chậm hoặc timeout

**Solution**: Đã cấu hình timeout 10m trong Jenkinsfile:

```groovy
golangci-lint run --timeout 10m
```

Nếu vẫn chậm, có thể disable một số linters:

```bash
golangci-lint run --timeout 10m --disable=linter1,linter2
```

### Issue: Build Docker image failed

**Solution**:

1. Kiểm tra Dockerfile syntax
2. Verify base image tồn tại: `docker pull golang:1.23-alpine`
3. Check network connectivity từ Jenkins server
4. Xem logs chi tiết trong Console Output

### Issue: Deployment approval timeout

**Solution**:

- Tăng timeout trong Jenkinsfile (hiện tại: 30 phút):

```groovy
timeout(time: 60, unit: 'MINUTES') {  // Tăng lên 60 phút
```

- Hoặc approve nhanh hơn trong 30 phút

### Issue: Environment file không tồn tại

**Error**: `No such file or directory: /app/env/.env.production`

**Solution**:

```bash
# Tạo missing file trên Jenkins server
sudo touch /app/env/.env.production
sudo chown jenkins:jenkins /app/env/.env.production
sudo chmod 600 /app/env/.env.production
# Điền nội dung environment variables
sudo nano /app/env/.env.production
```

### Issue: Webhook không trigger build

**Solution**:

1. **Check webhook URL** đúng format
2. **Test webhook** trong Git repository settings
3. **Firewall**: Đảm bảo Git server access được Jenkins
4. **Logs**: Check Jenkins system logs
5. **Manual trigger**: Thử click "Build Now" để test pipeline

### Issue: Push image to registry failed (authentication)

**Solution**:

1. Verify credentials ID khớp: `docker-hub-credentials`
2. Test login manually:

```bash
docker login -u your-username
```

3. Recreate Jenkins credentials nếu cần
4. Check network từ Jenkins đến Docker registry

## Notifications (Tùy chọn)

Pipeline đã có template notifications trong `post` blocks. Để enable:

### Slack Integration:

1. Cài **Slack Notification Plugin**
2. Cấu hình Slack trong Jenkins: Manage Jenkins → Configure System
3. Uncomment Slack code trong Jenkinsfile:

```groovy
post {
    success {
        slackSend(
            channel: '#ci-cd',
            color: 'good',
            message: "✅ Build #${env.BUILD_NUMBER} succeeded\nBranch: ${env.BRANCH_NAME}"
        )
    }
    failure {
        slackSend(
            channel: '#ci-cd',
            color: 'danger',
            message: "❌ Build #${env.BUILD_NUMBER} failed\nBranch: ${env.BRANCH_NAME}\n${env.BUILD_URL}"
        )
    }
}
```

### Email Notification:

1. Cấu hình SMTP trong Jenkins: Manage Jenkins → Configure System → E-mail Notification
2. Uncomment email code trong Jenkinsfile:

```groovy
post {
    failure {
        emailext(
            subject: "❌ Build Failed: ${env.JOB_NAME} #${env.BUILD_NUMBER}",
            body: """
                Build: ${env.BUILD_URL}
                Branch: ${env.BRANCH_NAME}
                Commit: ${GIT_COMMIT_SHORT}
            """,
            to: "devteam@example.com",
            recipientProviders: [[$class: 'DevelopersRecipientProvider']]
        )
    }
}
```

## Monitoring & Metrics

### Build Metrics (Jenkins Dashboard)

Theo dõi:

- **Build Success Rate**: Tỷ lệ builds thành công/thất bại
- **Build Duration**: Thời gian trung bình của pipeline
- **Test Coverage Trend**: Coverage tăng/giảm theo thời gian
- **Deployment Frequency**: Số lần deploy/tuần

### Application Metrics (Post-Deploy)

Monitor sau deployment:

- Health endpoint status
- Response time
- Error rates
- Resource usage (CPU, Memory)

## Quick Reference

### Common Commands

```bash
# Test local build (không cần Jenkins)
docker build -t go-backend:test .
docker run --rm go-backend:test

# Manual test coverage locally
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run golangci-lint locally
docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run

# Check security locally
docker run --rm -v $(pwd):/app securego/gosec:latest /app/...
```

### Environment Variables trong Pipeline

| Variable             | Description         | Example                          |
| -------------------- | ------------------- | -------------------------------- |
| `BRANCH_NAME`        | Tên branch hiện tại | `main`, `develop`, `feature/xyz` |
| `TAG_NAME`           | Tên tag (nếu có)    | `v1.0.0`, `release-2024`         |
| `BUILD_NUMBER`       | Build number        | `42`                             |
| `GIT_COMMIT_SHORT`   | Short commit hash   | `a1b2c3d`                        |
| `IMAGE_TAG`          | Docker image tag    | `main-a1b2c3d` hoặc `latest`     |
| `DEPLOY_ENVIRONMENT` | Environment đã chọn | `production`, `staging`          |

## Tài liệu tham khảo

- [Jenkins Pipeline Documentation](https://www.jenkins.io/doc/book/pipeline/)
- [Docker Pipeline Plugin](https://plugins.jenkins.io/docker-workflow/)
- [HTML Publisher Plugin](https://plugins.jenkins.io/htmlpublisher/)
- [Golang Official Docker Images](https://hub.docker.com/_/golang)
- [golangci-lint Documentation](https://golangci-lint.run/)
- [gosec - Security Scanner](https://github.com/securego/gosec)

---

**Cập nhật lần cuối**: January 2026  
**Jenkins Version**: 2.x+  
**Go Version**: 1.23  
**Pipeline Type**: Declarative Pipeline
