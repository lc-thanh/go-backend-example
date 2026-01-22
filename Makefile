.PHONY: help test test-coverage test-verbose build run clean docker-build docker-up docker-down lint fmt vet

# Variables
BINARY_NAME=server
DOCKER_COMPOSE=docker-compose
GO=go
GOTEST=$(GO) test
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOVET=$(GO) vet
GOFMT=$(GO) fmt

# Colors for terminal output
GREEN=\033[0;32m
NC=\033[0m # No Color

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  ${GREEN}%-15s${NC} %s\n", $$1, $$2}' $(MAKEFILE_LIST)

test: ## Run all tests
	@echo "Running tests..."
	$(GOTEST) -v -race ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	$(GOTEST) -v -race -cover ./...

test-short: ## Run short tests only
	@echo "Running short tests..."
	$(GOTEST) -short ./...

bench: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

build: ## Build the application
	@echo "Building..."
	$(GOBUILD) -o bin/$(BINARY_NAME) -v .
	@echo "Build complete: bin/$(BINARY_NAME)"

run: ## Run the application
	@echo "Running application..."
	$(GO) run .

clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f bin/$(BINARY_NAME)
	rm -f coverage.out coverage.html
	@echo "Clean complete"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t go-backend:latest .

docker-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	$(DOCKER_COMPOSE) up -d

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) down

docker-logs: ## View Docker logs
	$(DOCKER_COMPOSE) logs -f

docker-restart: ## Restart Docker containers
	@echo "Restarting Docker containers..."
	$(DOCKER_COMPOSE) restart

lint: ## Run golangci-lint
	@echo "Running linter..."
	golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	$(GOFMT) ./...

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod verify

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	$(GO) mod tidy

update-deps: ## Update all dependencies
	@echo "Updating dependencies..."
	$(GO) get -u ./...
	$(GO) mod tidy

install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed successfully"

ci: deps vet lint test ## Run CI pipeline locally

all: clean deps build test ## Run all tasks
