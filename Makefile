.PHONY: help test test-unit test-integration test-e2e test-coverage test-db-setup test-db-teardown build run docker-build docker-run lint

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Variables
BINARY_NAME=task-service
DOCKER_IMAGE=task-service
GO_FILES=$(shell find . -type f -name '*.go' -not -path './vendor/*')
TEST_TIMEOUT=300s

# Build targets
build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BINARY_NAME) ./cmd/api

run: ## Run the application locally
	@echo "Running $(BINARY_NAME)..."
	go run ./cmd/api

# Docker targets
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run application in Docker
	@echo "Running application in Docker..."
	docker compose up -d

docker-stop: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	docker compose down

docker-logs: ## Show Docker container logs
	@echo "Showing Docker logs..."
	docker compose logs -f

# Test targets
test: test-unit test-integration ## Run all tests

#docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
test-unit: ## Run unit tests
	@echo "Running unit tests..."
	#go test -v -race -timeout $(TEST_TIMEOUT) ./internal/service/... ./internal/api/handler/...
	docker compose -f docker-compose.test.yml run --rm unit-tests

test-integration: test-db-setup ## Run integration tests
	@echo "Running integration tests..."
	#go test -v -race -timeout $(TEST_TIMEOUT) ./internal/repository/postgres/...
	docker compose -f docker-compose.test.yml run --rm integration-tests

test-e2e: test-db-setup ## Run end-to-end tests
	@echo "Running E2E tests..."
	#go test -v -race -timeout $(TEST_TIMEOUT) ./tests/e2e/...
	docker compose -f docker-compose.test.yml run --rm e2e-tests

test-coverage: test-db-setup ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-benchmark: ## Run benchmark tests
	@echo "Running benchmark tests..."
	go test -bench=. -benchmem ./...

# Test database setup
test-db-setup: ## Setup test database
	@echo "Setting up test database..."
	@docker run --name test-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=postgres -p 5433:5432 -d postgres:15-alpine || true
	@sleep 5

test-db-teardown: ## Teardown test database
	@echo "Tearing down test database..."
	@docker stop test-postgres || true
	@docker rm test-postgres || true

# Code quality targets
lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	gofmt -s -w $(GO_FILES)

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

# Dependency management
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod verify

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

deps-clean: ## Clean dependency cache
	@echo "Cleaning dependency cache..."
	go clean -modcache

# Development targets
dev-setup: deps test-db-setup ## Setup development environment
	@echo "Development environment setup complete!"

dev-reset: test-db-teardown test-db-setup ## Reset development environment
	@echo "Development environment reset complete!"

# Database migration targets
migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/taskdb?sslmode=disable" up

migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/taskdb?sslmode=disable" down

migrate-create: ## Create new migration (usage: make migrate-create NAME=migration_name)
	@echo "Creating migration $(NAME)..."
	migrate create -ext sql -dir ./migrations $(NAME)

# Security targets
security-scan: ## Run security scan
	@echo "Running security scan..."
	gosec ./...

# Documentation targets
docs: ## Generate documentation
	@echo "Generating documentation..."
	godoc -http=:6060

# Clean targets
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	go clean -cache
	go clean -testcache

# Performance testing
load-test: ## Run load tests (requires wrk)
	@echo "Running load tests..."
	@if ! command -v wrk > /dev/null; then \
		echo "wrk is not installed. Please install it first."; \
		exit 1; \
	fi
	wrk -t12 -c400 -d30s --timeout 30s -H "Authorization: Bearer $(shell curl -s -X POST http://localhost:8080/auth/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin"}' | jq -r '.token')" http://localhost:8080/api/v1/tasks