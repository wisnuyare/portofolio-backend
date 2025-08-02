# Makefile for Portfolio Backend
.PHONY: help build run test clean tidy lint docker-build docker-run migrate-up migrate-down dev deps

# Variables
BINARY_NAME := portfolio-backend
BUILD_DIR := build
VERSION := 1.0.0
COMMIT_SHA := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Help target
help: ## Show this help message
	@echo "Portfolio Backend Makefile"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development targets
dev: tidy build run ## Run development workflow (tidy, build, run)

run: ## Run the application
	go run cmd/api/main.go

build: ## Build the application
	@./scripts/build.sh

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests and show coverage
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linters
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, skipping linting"; \
	fi

tidy: ## Tidy go modules
	go mod tidy

deps: ## Download dependencies
	go mod download

clean: ## Clean build artifacts
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Docker targets
docker-build: ## Build Docker image
	docker build -f deployments/docker/Dockerfile -t $(BINARY_NAME):latest .

docker-run: ## Run application in Docker
	docker run --rm -p 8080:8080 \
		-e DB_HOST=host.docker.internal \
		-e DB_USER=portfolio_user \
		-e DB_PASSWORD=your_password \
		-e DB_NAME=portfolio_db \
		$(BINARY_NAME):latest

# Database targets
migrate-up: ## Run database migrations up
	@./scripts/migrate.sh up

migrate-down: ## Roll back last migration
	@./scripts/migrate.sh down

migrate-version: ## Show current migration version
	@./scripts/migrate.sh version

migrate-create: ## Create new migration (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "Please provide a migration name: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@./scripts/migrate.sh create $(NAME)

# Deployment targets
deploy: ## Deploy to GCP (requires PROJECT_ID)
	@./scripts/deploy.sh

deploy-terraform: ## Deploy infrastructure with Terraform
	@cd deployments/gcp/terraform && \
		terraform init && \
		terraform plan && \
		terraform apply

# Configuration targets
config-example: ## Generate example configuration file
	@echo "# Portfolio Backend Configuration" > config.example.yaml
	@echo "server:" >> config.example.yaml
	@echo "  host: 0.0.0.0" >> config.example.yaml
	@echo "  port: 8080" >> config.example.yaml
	@echo "  read_timeout: 30s" >> config.example.yaml
	@echo "  write_timeout: 30s" >> config.example.yaml
	@echo "" >> config.example.yaml
	@echo "database:" >> config.example.yaml
	@echo "  host: localhost" >> config.example.yaml
	@echo "  port: 3306" >> config.example.yaml
	@echo "  user: portfolio_user" >> config.example.yaml
	@echo "  password: your_secure_password" >> config.example.yaml
	@echo "  database: portfolio_db" >> config.example.yaml
	@echo "  max_open_conns: 25" >> config.example.yaml
	@echo "  max_idle_conns: 5" >> config.example.yaml
	@echo "  conn_max_lifetime: 5m" >> config.example.yaml
	@echo "" >> config.example.yaml
	@echo "cors:" >> config.example.yaml
	@echo "  allowed_origins:" >> config.example.yaml
	@echo "    - https://your-frontend-domain.com" >> config.example.yaml
	@echo "    - http://localhost:3000" >> config.example.yaml
	@echo "  allowed_methods:" >> config.example.yaml
	@echo "    - GET" >> config.example.yaml
	@echo "    - POST" >> config.example.yaml
	@echo "    - PUT" >> config.example.yaml
	@echo "    - DELETE" >> config.example.yaml
	@echo "    - OPTIONS" >> config.example.yaml
	@echo "  allowed_headers:" >> config.example.yaml
	@echo "    - Content-Type" >> config.example.yaml
	@echo "    - Authorization" >> config.example.yaml
	@echo "" >> config.example.yaml
	@echo "logging:" >> config.example.yaml
	@echo "  level: info" >> config.example.yaml
	@echo "  format: json" >> config.example.yaml
	@echo "Example configuration file created: config.example.yaml"

# Environment file
.env: ## Create .env file from example
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo ".env file created from .env.example"; \
		echo "Please update the values in .env file"; \
	else \
		echo ".env file already exists"; \
	fi

# Quality targets
check: lint test ## Run quality checks (lint and test)

ci: deps tidy lint test build ## Run CI pipeline locally

# Installation targets
install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "Development tools installed successfully"

# Default target
all: deps tidy lint test build ## Run all tasks