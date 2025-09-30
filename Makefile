# Makefile for Arbak Portfolio Backend

# Variables
APP_NAME=arbak-portfolio-backend
DOCKER_IMAGE=portfolio-api
DOCKER_TAG=latest
GO_VERSION=1.21

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: help build run test clean docker-build docker-run docker-compose-up docker-compose-down migrate-up migrate-down lint format

# Default target
help: ## Show this help message
	@echo "$(BLUE)Arbak Portfolio Backend - Available Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

# Development commands
build: ## Build the application
	@echo "$(BLUE)Building $(APP_NAME)...$(NC)"
	@go build -o bin/$(APP_NAME) main.go
	@echo "$(GREEN)Build completed!$(NC)"

run: ## Run the application locally
	@echo "$(BLUE)Starting $(APP_NAME)...$(NC)"
	@go run main.go

test: ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)Clean completed!$(NC)"

# Code quality
lint: ## Run linter
	@echo "$(BLUE)Running linter...$(NC)"
	@golangci-lint run

format: ## Format code
	@echo "$(BLUE)Formatting code...$(NC)"
	@go fmt ./...
	@goimports -w .

# Database commands
migrate-up: ## Run database migrations
	@echo "$(BLUE)Running database migrations...$(NC)"
	@go run main.go migrate up

migrate-down: ## Rollback database migrations
	@echo "$(BLUE)Rolling back database migrations...$(NC)"
	@go run main.go migrate down

# Docker commands
docker-build: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(NC)"
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "$(GREEN)Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)$(NC)"

docker-run: ## Run Docker container
	@echo "$(BLUE)Running Docker container...$(NC)"
	@docker run -d \
		--name $(APP_NAME) \
		-p 8080:8080 \
		-e ENVIRONMENT=development \
		-e DATABASE_URL=postgres://portfolio_user:portfolio_password@localhost:5432/portfolio_db?sslmode=disable \
		-e REDIS_URL=redis://localhost:6379 \
		-e JWT_SECRET=dev-jwt-secret \
		$(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo "$(GREEN)Container started: $(APP_NAME)$(NC)"

docker-stop: ## Stop Docker container
	@echo "$(BLUE)Stopping Docker container...$(NC)"
	@docker stop $(APP_NAME) || true
	@docker rm $(APP_NAME) || true
	@echo "$(GREEN)Container stopped and removed$(NC)"

# Docker Compose commands
docker-compose-up: ## Start all services with Docker Compose
	@echo "$(BLUE)Starting all services with Docker Compose...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)All services started!$(NC)"
	@echo "$(YELLOW)API: http://localhost:8080$(NC)"
	@echo "$(YELLOW)Health: http://localhost:8080/health$(NC)"

docker-compose-down: ## Stop all services with Docker Compose
	@echo "$(BLUE)Stopping all services with Docker Compose...$(NC)"
	@docker-compose down
	@echo "$(GREEN)All services stopped!$(NC)"

docker-compose-logs: ## View Docker Compose logs
	@echo "$(BLUE)Viewing Docker Compose logs...$(NC)"
	@docker-compose logs -f

# Development setup
setup: ## Setup development environment
	@echo "$(BLUE)Setting up development environment...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)Development environment ready!$(NC)"

install-tools: ## Install development tools
	@echo "$(BLUE)Installing development tools...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "$(GREEN)Development tools installed!$(NC)"

# API documentation
docs: ## Generate API documentation
	@echo "$(BLUE)Generating API documentation...$(NC)"
	@swag init -g main.go -o ./docs
	@echo "$(GREEN)API documentation generated!$(NC)"

# Production deployment
deploy-staging: ## Deploy to staging environment
	@echo "$(BLUE)Deploying to staging...$(NC)"
	@docker build -t $(DOCKER_IMAGE):staging .
	@echo "$(GREEN)Staging deployment completed!$(NC)"

deploy-production: ## Deploy to production environment
	@echo "$(BLUE)Deploying to production...$(NC)"
	@docker build -t $(DOCKER_IMAGE):production .
	@echo "$(GREEN)Production deployment completed!$(NC)"

# Monitoring and health checks
health: ## Check application health
	@echo "$(BLUE)Checking application health...$(NC)"
	@curl -f http://localhost:8080/health || echo "$(RED)Health check failed!$(NC)"

logs: ## View application logs
	@echo "$(BLUE)Viewing application logs...$(NC)"
	@docker logs -f $(APP_NAME) 2>/dev/null || echo "$(RED)Container not running$(NC)"

# Database utilities
db-shell: ## Connect to database shell
	@echo "$(BLUE)Connecting to database shell...$(NC)"
	@docker exec -it portfolio_postgres psql -U portfolio_user -d portfolio_db

redis-cli: ## Connect to Redis CLI
	@echo "$(BLUE)Connecting to Redis CLI...$(NC)"
	@docker exec -it portfolio_redis redis-cli

# Backup and restore
backup-db: ## Backup database
	@echo "$(BLUE)Backing up database...$(NC)"
	@docker exec portfolio_postgres pg_dump -U portfolio_user portfolio_db > backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "$(GREEN)Database backup completed!$(NC)"

restore-db: ## Restore database from backup
	@echo "$(BLUE)Restoring database...$(NC)"
	@echo "$(YELLOW)Usage: make restore-db BACKUP_FILE=backup_20231201_120000.sql$(NC)"
	@if [ -z "$(BACKUP_FILE)" ]; then echo "$(RED)Please specify BACKUP_FILE$(NC)"; exit 1; fi
	@docker exec -i portfolio_postgres psql -U portfolio_user -d portfolio_db < $(BACKUP_FILE)
	@echo "$(GREEN)Database restore completed!$(NC)"

# Performance testing
benchmark: ## Run performance benchmarks
	@echo "$(BLUE)Running performance benchmarks...$(NC)"
	@go test -bench=. -benchmem ./...

load-test: ## Run load tests (requires hey tool)
	@echo "$(BLUE)Running load tests...$(NC)"
	@hey -n 1000 -c 10 http://localhost:8080/health

# Security
security-scan: ## Run security scan
	@echo "$(BLUE)Running security scan...$(NC)"
	@gosec ./...

# All-in-one commands
dev: setup docker-compose-up ## Setup and start development environment
	@echo "$(GREEN)Development environment ready!$(NC)"
	@echo "$(YELLOW)API: http://localhost:8080$(NC)"
	@echo "$(YELLOW)Health: http://localhost:8080/health$(NC)"

prod: docker-build docker-compose-up ## Build and start production environment
	@echo "$(GREEN)Production environment ready!$(NC)"

# Cleanup
cleanup: docker-compose-down clean ## Clean up everything
	@echo "$(BLUE)Cleaning up everything...$(NC)"
	@docker system prune -f
	@echo "$(GREEN)Cleanup completed!$(NC)"
