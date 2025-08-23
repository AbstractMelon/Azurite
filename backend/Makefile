.PHONY: build run test clean dev migrate seed docker-build docker-run help

# Variables
BINARY_NAME=azurite-server
MAIN_PATH=./cmd/server/main.go
BUILD_DIR=./bin

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run the application in development mode
dev:
	@echo "Starting development server..."
	@go run $(MAIN_PATH)

# Run the built binary
run: build
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	@go test -v -race ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p coverage
	@go test -v -coverprofile=coverage/coverage.out ./... || true
	@go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@go tool cover -func=coverage/coverage.out
	@echo "Coverage report generated: coverage/coverage.html"

# Run tests with detailed coverage by package
test-coverage-detailed:
	@echo "Running detailed coverage analysis..."
	@mkdir -p coverage
	@echo "mode: atomic" > coverage/coverage.out
	@for pkg in $$(go list ./... | grep -v /vendor/); do \
		go test -v -covermode=atomic -coverprofile=coverage/profile.out $$pkg; \
		if [ -f coverage/profile.out ]; then \
			grep -v "mode: atomic" coverage/profile.out >> coverage/coverage.out; \
			rm coverage/profile.out; \
		fi; \
	done
	@go tool cover -html=coverage/coverage.out -o coverage/detailed_coverage.html
	@go tool cover -func=coverage/coverage.out | tail -1
	@echo "Detailed coverage report generated: coverage/detailed_coverage.html"

# Run benchmark tests
test-bench:
	@echo "Running benchmark tests..."
	@go test -bench=. -benchmem ./...

# Run tests and check coverage threshold
test-coverage-check:
	@echo "Running tests with coverage threshold check..."
	@mkdir -p coverage
	@go test -v -coverprofile=coverage/coverage.out ./...
	@coverage=$$(go tool cover -func=coverage/coverage.out | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$coverage < 70" | bc -l) -eq 1 ]; then \
		echo "Coverage $$coverage% is below 70% threshold!"; \
		exit 1; \
	else \
		echo "Coverage $$coverage% meets threshold"; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -rf coverage/
	@rm -f coverage.out coverage.html
	@rm -f azurite.db azurite.db-shm azurite.db-wal
	@rm -f *.test

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod verify

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run

# Vet code
vet:
	@echo "Vetting code..."
	@go vet ./...

# Check code quality (fmt, vet, lint)
check: fmt vet lint

# Create database and run migrations (happens automatically on startup)
migrate:
	@echo "Running database migrations..."
	@go run $(MAIN_PATH) --migrate-only

# Seed database with sample data
seed:
	@echo "Seeding database..."
	@go run scripts/seed.go

# Watch for changes and restart (requires air: go install github.com/cosmtrek/air@latest)
watch:
	@echo "Starting development server with hot reload..."
	@air

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t azurite-backend .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env azurite-backend

# Docker compose up
docker-up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d

# Docker compose down
docker-down:
	@echo "Stopping services..."
	@docker-compose down

# Create a new migration file
migration:
	@read -p "Enter migration name: " name; \
	timestamp=$$(date +%Y%m%d_%H%M%S); \
	filename="migrations/$${timestamp}_$${name}.sql"; \
	echo "-- Migration: $$name" > $$filename; \
	echo "-- Created: $$(date)" >> $$filename; \
	echo "" >> $$filename; \
	echo "Migration file created: $$filename"

# Setup development environment
setup-dev:
	@echo "Setting up development environment..."
	@cp .env.example .env
	@mkdir -p storage/mods storage/images
	@make deps
	@echo "Development environment setup complete!"
	@echo "Please edit .env with your configuration before running the server."

# Build for production
build-prod:
	@echo "Building for production..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Production build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Cross-compile for different platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "Cross-compilation complete!"

# Generate API documentation (requires swag: go install github.com/swaggo/swag/cmd/swag@latest)
docs:
	@echo "Generating API documentation..."
	@swag init -g $(MAIN_PATH) -o ./docs

# Backup database
backup:
	@echo "Creating database backup..."
	@timestamp=$$(date +%Y%m%d_%H%M%S); \
	cp azurite.db "backups/azurite_backup_$$timestamp.db"; \
	echo "Backup created: backups/azurite_backup_$$timestamp.db"

# Restore database from backup
restore:
	@read -p "Enter backup filename (in backups/): " filename; \
	cp "backups/$$filename" azurite.db; \
	echo "Database restored from: backups/$$filename"

# Show application logs (if running with systemd)
logs:
	@journalctl -u azurite-backend -f

# Check application health
health:
	@curl -s http://localhost:8080/health | jq '.' || echo "Service not responding"

# Show help
help:
	@echo "Azurite Backend - Available Commands:"
	@echo ""
	@echo "Development:"
	@echo "  dev            - Run the application in development mode"
	@echo "  watch          - Run with hot reload (requires air)"
	@echo "  setup-dev      - Setup development environment"
	@echo ""
	@echo "Building:"
	@echo "  build          - Build the application"
	@echo "  build-prod     - Build for production (static binary)"
	@echo "  build-all      - Cross-compile for multiple platforms"
	@echo ""
	@echo "Testing:"
	@echo "  test           - Run tests"
	@echo "  test-race      - Run tests with race detection"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-coverage-detailed - Run detailed coverage analysis by package"
	@echo "  test-coverage-check - Run tests and check coverage threshold (70%)"
	@echo "  test-bench     - Run benchmark tests"
	@echo ""
	@echo "Code Quality:"
	@echo "  fmt            - Format code"
	@echo "  vet            - Vet code"
	@echo "  lint           - Lint code (requires golangci-lint)"
	@echo "  check          - Run fmt, vet, and lint"
	@echo ""
	@echo "Dependencies:"
	@echo "  deps           - Install dependencies"
	@echo "  deps-update    - Update dependencies"
	@echo ""
	@echo "Database:"
	@echo "  migrate        - Run database migrations"
	@echo "  seed           - Seed database with sample data"
	@echo "  backup         - Backup database"
	@echo "  restore        - Restore database from backup"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-up      - Start with Docker Compose"
	@echo "  docker-down    - Stop Docker Compose services"
	@echo ""
	@echo "Utilities:"
	@echo "  migration      - Create new migration file"
	@echo "  docs           - Generate API documentation"
	@echo "  logs           - Show application logs"
	@echo "  health         - Check application health"
	@echo "  clean          - Clean build artifacts and database"
	@echo "  help           - Show this help message"
