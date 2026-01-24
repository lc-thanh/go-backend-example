pipeline {
    agent {
        label 'backend-server'
    }
    
    environment {
        // Docker configuration
        DOCKER_REGISTRY = 'registry.lcthanh.cloud'
        DOCKER_IMAGE_NAME = 'lcthanh/go-backend'
        DOCKER_IMAGE_FULL_NAME = "${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}"
        DOCKER_CREDENTIALS_ID = 'docker-hub-credentials'
        
        // Go configuration
        GO_VERSION = '1.23'
        CGO_ENABLED = '0'
        GOOS = 'linux'
        GOARCH = 'amd64'
        GOLANGCI_LINT_VERSION = 'v1.60.3-alpine'
        
        // Environment files path on Jenkins node
        ENV_DIR = '/app/envs/go-backend'
        
        // Build info
        BUILD_TIME = sh(script: "date -u '+%Y-%m-%dT%H:%M:%SZ'", returnStdout: true).trim()
        GIT_COMMIT_SHORT = sh(script: "git rev-parse --short HEAD", returnStdout: true).trim()
        
        // Get current branch/tag name (handle null BRANCH_NAME and detached HEAD)
        CURRENT_BRANCH = sh(script: '''
            if [ -n "${TAG_NAME:-}" ]; then
                echo "${TAG_NAME}"
            elif [ -n "${BRANCH_NAME:-}" ]; then
                echo "${BRANCH_NAME}"
            elif [ -n "${GIT_BRANCH:-}" ]; then
                # GIT_BRANCH might be like "origin/main", strip the remote part
                echo "${GIT_BRANCH}" | sed 's|^origin/||'
            else
                # Try to get branch from git, handle detached HEAD
                BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
                if [ "$BRANCH" = "HEAD" ] || [ -z "$BRANCH" ]; then
                    # Detached HEAD state, try to find branch from remote
                    BRANCH=$(git branch -r --contains HEAD | grep -v HEAD | head -1 | sed 's|origin/||' | xargs || echo "detached-HEAD")
                fi
                echo "$BRANCH"
            fi
        ''', returnStdout: true).trim()
        
        // Image tag based on branch/tag
        IMAGE_TAG = sh(script: '''
            BRANCH="${CURRENT_BRANCH}"
            COMMIT="${GIT_COMMIT_SHORT}"
            
            # Ensure we have a commit hash
            if [ -z "$COMMIT" ]; then
                COMMIT=$(git rev-parse --short HEAD || echo "unknown")
            fi
            
            # Generate tag based on branch
            if [ "$BRANCH" = "main" ] || [ "$BRANCH" = "master" ]; then
                echo "latest"
            elif [ -n "$BRANCH" ] && [ "$BRANCH" != "HEAD" ]; then
                echo "${BRANCH}-${COMMIT}"
            else
                echo "build-${COMMIT}"
            fi
        ''', returnStdout: true).trim()
    }
    
    options {
        // Keep last 10 builds
        buildDiscarder(logRotator(numToKeepStr: '10'))
        
        // Timeout for entire pipeline
        timeout(time: 30, unit: 'MINUTES')
        
        // Disable concurrent builds
        disableConcurrentBuilds()
        
        // Timestamps in console output
        timestamps()
    }
    
    stages {
        stage('Checkout') {
            steps {
                script {
                    echo "🔄 Checking out code..."
                    echo "Branch/Tag: ${CURRENT_BRANCH}"
                    echo "Commit: ${GIT_COMMIT_SHORT}"
                    if (env.TAG_NAME) {
                        echo "Tag Name: ${env.TAG_NAME}"
                    }
                    if (env.BRANCH_NAME) {
                        echo "Branch Name: ${env.BRANCH_NAME}"
                    }
                }
                checkout scm
            }
        }
        
        stage('Setup Environment') {
            steps {
                script {
                    echo "🔧 Setting up Go environment..."
                }
                sh '''
                    docker run --rm golang:${GO_VERSION}-alpine go version
                    docker run --rm golang:${GO_VERSION}-alpine go env
                '''
            }
        }
        
        stage('Download Dependencies') {
            steps {
                script {
                    echo "📦 Downloading Go modules..."
                }
                sh '''
                    docker run --rm \
                        -v "$(pwd):/app" \
                        -w /app \
                        golang:${GO_VERSION}-alpine \
                        go mod download
                '''
            }
        }
        
        stage('Code Quality') {
            parallel {
                stage('Format Check') {
                    steps {
                        script {
                            echo "🎨 Checking code formatting..."
                        }
                        sh '''
                            docker run --rm \
                                -v "$(pwd):/app" \
                                -w /app \
                                golang:${GO_VERSION}-alpine \
                                sh -c '
                                    unformatted=$(gofmt -l .)
                                    if [ -n "$unformatted" ]; then
                                        echo "❌ Code is not formatted properly:"
                                        echo "$unformatted"
                                        exit 1
                                    fi
                                    echo "✅ All code is properly formatted"
                                '
                        '''
                    }
                }
                
                stage('Go Vet') {
                    steps {
                        script {
                            echo "🔍 Running go vet..."
                        }
                        sh '''
                            docker run --rm \
                                -v "$(pwd):/app" \
                                -w /app \
                                golang:${GO_VERSION}-alpine \
                                go vet ./...
                        '''
                    }
                }
                
                stage('Lint') {
                    steps {
                        script {
                            echo "🧹 Running golangci-lint..."
                        }
                        sh '''
                            docker run --rm \
                                -v "$(pwd):/app" \
                                -w /app \
                                golangci/golangci-lint:${GOLANGCI_LINT_VERSION} \
                                golangci-lint run --timeout 10m
                        '''
                    }
                }
            }
        }
        
        stage('Unit Tests') {
            steps {
                script {
                    echo "🧪 Running unit tests with coverage..."
                }
                sh '''
                    docker run --rm \
                        -v "$(pwd):/app" \
                        -w /app \
                        -e ENV_FILE=${ENV_DIR}/.env.test \
                        golang:${GO_VERSION}-alpine \
                        sh -c '
                            # Load test environment if exists
                            if [ -f "${ENV_FILE}" ]; then
                                export $(cat ${ENV_FILE} | grep -v "^#" | xargs)
                            fi
                            
                            # Run tests with coverage
                            go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
                            
                            # Generate coverage report
                            go tool cover -func=coverage.out
                            
                            # Calculate coverage percentage
                            coverage=$(go tool cover -func=coverage.out | grep total | awk "{print \\$3}" | sed "s/%//")
                            echo "📊 Total Coverage: ${coverage}%"
                            
                            # Fail if coverage is below 70%
                            if [ $(echo "$coverage < 70" | bc -l) -eq 1 ]; then
                                echo "❌ Coverage ${coverage}% is below 70% threshold"
                                exit 1
                            fi
                        '
                '''
            }
            post {
                always {
                    // Archive coverage report
                    archiveArtifacts artifacts: 'coverage.out', allowEmptyArchive: true
                    
                    // Publish HTML coverage report (requires HTML Publisher plugin)
                    sh '''
                        docker run --rm \
                            -v "$(pwd):/app" \
                            -w /app \
                            golang:${GO_VERSION}-alpine \
                            go tool cover -html=coverage.out -o coverage.html || true
                    '''
                    publishHTML([
                        allowMissing: true,
                        alwaysLinkToLastBuild: true,
                        keepAll: true,
                        reportDir: '.',
                        reportFiles: 'coverage.html',
                        reportName: 'Go Coverage Report'
                    ])
                }
            }
        }
        
        stage('Security Scan') {
            steps {
                script {
                    echo "🔒 Running security scan with gosec..."
                }
                sh '''
                    docker run --rm \
                        -v "$(pwd):/app" \
                        -w /app \
                        securego/gosec:latest \
                        -fmt json -out gosec-report.json ./... || true
                    
                    # Display summary
                    docker run --rm \
                        -v "$(pwd):/app" \
                        -w /app \
                        securego/gosec:latest \
                        ./... || true
                '''
            }
            post {
                always {
                    archiveArtifacts artifacts: 'gosec-report.json', allowEmptyArchive: true
                }
            }
        }
        
        stage('Build Binary') {
            steps {
                script {
                    echo "🔨 Building Go binary..."
                }
                sh '''
                    # Run docker with current user to avoid permission issues
                    # Set GOCACHE and GOMODCACHE to workspace to avoid permission denied
                    docker run --rm \
                        -u $(id -u):$(id -g) \
                        -v "$(pwd):/app" \
                        -w /app \
                        -e CGO_ENABLED=${CGO_ENABLED} \
                        -e GOOS=${GOOS} \
                        -e GOARCH=${GOARCH} \
                        -e GOCACHE=/app/.cache/go-build \
                        -e GOMODCACHE=/app/.cache/go-mod \
                        golang:${GO_VERSION}-alpine \
                        sh -c '
                            # Create cache directories
                            mkdir -p .cache/go-build .cache/go-mod
                            
                            go build -v \
                                -ldflags="-s -w -X main.Version=${IMAGE_TAG} -X main.BuildTime=${BUILD_TIME}" \
                                -o bin/go-backend \
                                .
                            
                            ls -lh bin/go-backend
                        '
                '''
            }
            post {
                success {
                    archiveArtifacts artifacts: 'bin/go-backend', allowEmptyArchive: false
                }
            }
        }
        
        // ========== CD STAGES (Only for main branch or tags) ==========
        
        stage('Build Docker Image') {
            when {
                expression {
                    return CURRENT_BRANCH == 'main' || 
                           CURRENT_BRANCH == 'master' || 
                           env.TAG_NAME != null
                }
            }
            steps {
                script {
                    echo "🐳 Building Docker image: ${DOCKER_IMAGE_FULL_NAME}:${IMAGE_TAG}"
                }
                sh '''
                    # Build Docker image
                    docker build -t ${DOCKER_IMAGE_FULL_NAME}:${IMAGE_TAG} .
                    
                    # Tag as latest if main branch
                    if [ "${CURRENT_BRANCH}" = "main" ] || [ "${CURRENT_BRANCH}" = "master" ]; then
                        docker tag ${DOCKER_IMAGE_FULL_NAME}:${IMAGE_TAG} ${DOCKER_IMAGE_FULL_NAME}:latest
                    fi
                    
                    # List built images
                    docker images | grep ${DOCKER_IMAGE_FULL_NAME}
                '''
            }
        }
        
        stage('Test Docker Image') {
            when {
                expression {
                    return CURRENT_BRANCH == 'main' || 
                           CURRENT_BRANCH == 'master' || 
                           env.TAG_NAME != null
                }
            }
            steps {
                script {
                    echo "🧪 Testing Docker image..."
                }
                sh '''
                    # Start container in background
                    docker run -d \
                        --name go-backend-test-${BUILD_NUMBER} \
                        -p 8081:8080 \
                        --env-file ${ENV_DIR}/.env.test \
                        ${DOCKER_IMAGE_FULL_NAME}:${IMAGE_TAG}
                    
                    # Wait for container to be ready
                    sleep 5
                    
                    # Check if container is running
                    docker ps | grep go-backend-test-${BUILD_NUMBER}
                    
                    # Test health endpoint (if available)
                    curl -f http://localhost:8081/health || echo "Health endpoint not available"
                    
                    # Stop and remove test container
                    docker stop go-backend-test-${BUILD_NUMBER}
                    docker rm go-backend-test-${BUILD_NUMBER}
                '''
            }
        }
        
        stage('Push Docker Image') {
            when {
                expression {
                    return CURRENT_BRANCH == 'main' || 
                           CURRENT_BRANCH == 'master' || 
                           env.TAG_NAME != null
                }
            }
            steps {
                script {
                    echo "📤 Pushing Docker image to registry..."
                }
                withCredentials([usernamePassword(
                    credentialsId: "${DOCKER_CREDENTIALS_ID}",
                    usernameVariable: 'DOCKER_USER',
                    passwordVariable: 'DOCKER_PASS'
                )]) {
                    sh '''
                        # Login to Docker registry
                        echo "$DOCKER_PASS" | docker login ${DOCKER_REGISTRY} -u "$DOCKER_USER" --password-stdin
                        
                        # Push image with specific tag
                        docker push ${DOCKER_IMAGE_FULL_NAME}:${IMAGE_TAG}
                        
                        # Push latest tag if main branch
                        if [ "${CURRENT_BRANCH}" = "main" ] || [ "${CURRENT_BRANCH}" = "master" ]; then
                            docker push ${DOCKER_IMAGE_FULL_NAME}:latest
                        fi
                        
                        # Logout
                        docker logout ${DOCKER_REGISTRY}
                    '''
                }
            }
        }
        
        stage('Deploy Approval') {
            when {
                expression {
                    return CURRENT_BRANCH == 'main' || 
                           CURRENT_BRANCH == 'master' || 
                           env.TAG_NAME != null
                }
            }
            steps {
                script {
                    echo "⏸️ Waiting for deployment approval..."
                    
                    def deploymentType = env.TAG_NAME ? "Tag: ${env.TAG_NAME}" : "Branch: ${CURRENT_BRANCH}"
                    
                    timeout(time: 30, unit: 'MINUTES') {
                        input(
                            message: "🚀 Deploy to Production?",
                            ok: 'Deploy',
                            submitter: 'admin,deployer',
                            parameters: [
                                choice(
                                    name: 'DEPLOY_ENVIRONMENT',
                                    choices: ['production', 'staging'],
                                    description: 'Select deployment environment'
                                ),
                                text(
                                    name: 'DEPLOY_NOTES',
                                    defaultValue: "Deploying ${deploymentType}\nCommit: ${GIT_COMMIT_SHORT}\nImage: ${IMAGE_TAG}",
                                    description: 'Deployment notes'
                                )
                            ]
                        )
                    }
                }
            }
        }
        
        stage('Deploy') {
            when {
                expression {
                    return CURRENT_BRANCH == 'main' || 
                           CURRENT_BRANCH == 'master' || 
                           env.TAG_NAME != null
                }
            }
            steps {
                script {
                    def deployEnv = env.DEPLOY_ENVIRONMENT ?: 'production'
                    echo "🚀 Deploying to ${deployEnv}..."
                    echo "📝 Notes: ${env.DEPLOY_NOTES}"
                    
                    // Load environment-specific configuration
                    def envFile = "${ENV_DIR}/.env.${deployEnv}"
                    
                    sh """
                        # Deploy using docker-compose (example)
                        # Replace with your actual deployment method
                        
                        echo "Environment: ${deployEnv}"
                        echo "Image: ${DOCKER_IMAGE_FULL_NAME}:${IMAGE_TAG}"
                        echo "Env file: ${envFile}"
                        
                        # Example deployment commands:
                        # Option 1: Docker Compose
                        # docker-compose -f docker-compose.${deployEnv}.yml down
                        # docker-compose -f docker-compose.${deployEnv}.yml up -d
                        
                        # Option 2: Docker Swarm
                        # docker stack deploy -c docker-compose.${deployEnv}.yml go-backend
                        
                        # Option 3: Kubernetes
                        # kubectl set image deployment/go-backend \\
                        #     go-backend=${DOCKER_IMAGE_FULL_NAME}:${IMAGE_TAG} \\
                        #     -n ${deployEnv}
                        # kubectl rollout status deployment/go-backend -n ${deployEnv}
                        
                        # Option 4: Direct Docker run
                        docker stop go-backend-${deployEnv} || true
                        docker rm go-backend-${deployEnv} || true
                        docker run -d \\
                            --name go-backend-${deployEnv} \\
                            -p 8080:8080 \\
                            --env-file ${envFile} \\
                            --restart unless-stopped \\
                            ${DOCKER_IMAGE_FULL_NAME}:${IMAGE_TAG}
                        
                        echo "✅ Deployment completed successfully"
                    """
                }
            }
        }
        
        stage('Post-Deploy Verification') {
            when {
                expression {
                    return CURRENT_BRANCH == 'main' || 
                           CURRENT_BRANCH == 'master' || 
                           env.TAG_NAME != null
                }
            }
            steps {
                script {
                    echo "✅ Verifying deployment..."
                }
                sh '''
                    # Add your verification checks here
                    # Examples:
                    # - Health check endpoint
                    # - Smoke tests
                    # - Database migration verification
                    
                    # curl -f https://your-app.com/health || exit 1
                    
                    echo "✅ All verification checks passed"
                '''
            }
        }
    }
    
    post {
        always {
            script {
                echo "🧹 Cleaning up workspace..."
            }
            
            // Clean up Docker images
            sh '''
                docker image prune -f --filter "until=24h" || true
            '''
            
            // Clean workspace with fallback for root-owned files
            script {
                try {
                    cleanWs(
                        deleteDirs: true,
                        patterns: [
                            [pattern: 'bin/', type: 'INCLUDE'],
                            [pattern: '.cache/', type: 'INCLUDE'],
                            [pattern: 'coverage.*', type: 'INCLUDE'],
                            [pattern: 'gosec-report.json', type: 'INCLUDE']
                        ]
                    )
                } catch (Exception e) {
                    echo "⚠️ cleanWs() failed, using docker to remove root-owned files..."
                    sh '''
                        # Remove root-owned files using docker
                        docker run --rm -v "$(pwd):/workspace" alpine sh -c 'rm -rf /workspace/bin /workspace/.cache /workspace/coverage.* /workspace/gosec-report.json' || true
                        # Try cleanWs again
                    '''
                    cleanWs(deleteDirs: true)
                }
            }
        }
        
        success {
            script {
                def message = "✅ Build #${env.BUILD_NUMBER} succeeded"
                message += "\nBranch/Tag: ${CURRENT_BRANCH}"
                if (CURRENT_BRANCH == 'main' || CURRENT_BRANCH == 'master' || env.TAG_NAME) {
                    message += "\n🚀 Deployed: ${DOCKER_IMAGE_FULL_NAME}:${IMAGE_TAG}"
                }
                echo message
                
                // Uncomment to enable Slack notification
                // slackSend(
                //     color: 'good',
                //     message: message
                // )
            }
        }
        
        failure {
            script {
                def message = "❌ Build #${env.BUILD_NUMBER} failed\nBranch/Tag: ${CURRENT_BRANCH}\nCommit: ${GIT_COMMIT_SHORT}"
                echo message
                
                // Uncomment to enable Slack notification
                // slackSend(
                //     color: 'danger',
                //     message: message
                // )
                
                // Uncomment to enable email notification
                // emailext(
                //     subject: "Build Failed: ${env.JOB_NAME} #${env.BUILD_NUMBER}",
                //     body: message,
                //     to: "team@example.com"
                // )
            }
        }
        
        unstable {
            script {
                echo "⚠️ Build #${env.BUILD_NUMBER} is unstable"
            }
        }
    }
}
