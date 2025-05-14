
APP_NAME := marketplace-api
BINARY := $(APP_NAME)
GO := go
GOLANGCI_LINT := golangci-lint
MIGRATION_DIR := migration
CONFIG_FILE := config.yaml
PORT := 30011
DOCKER_IMAGE := $(APP_NAME):latest
POSTGRES_DSN := host=localhost user=postgres password=postgres dbname=auth sslmode=disable
POSTGRES_DSN_PROD := host=database-1.cm9ewocwci8f.us-east-1.rds.amazonaws.com user=postgres password=Swanhtetaungphyo dbname=postgres port=5432 sslmode=require
ENV := development

# Colors for output
YELLOW := \033[1;33m
GREEN := \033[1;32m
RED := \033[1;31m
NC := \033[0m

# Default target
.PHONY: all
all: deps fmt lint test build ## Run all essential development tasks

# Help documentation
.PHONY: help
help: ## Display this help message
	@echo "$(YELLOW)Makefile for $(APP_NAME)$(NC)"
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

# Dependency management
.PHONY: deps
deps: ## Install Go dependencies
	@echo "$(YELLOW)📦 Installing dependencies...$(NC)"
	@$(GO) mod tidy
	@$(GO) mod vendor

.PHONY: deps-update
deps-update: ## Update Go dependencies to latest versions
	@echo "$(YELLOW)🔄 Updating dependencies...$(NC)"
	@$(GO) get -u ./...
	@$(GO) mod tidy
	@$(GO) mod vendor

# Code quality
.PHONY: fmt
fmt: ## Format Go code
	@echo "$(YELLOW)🖌️ Formatting code...$(NC)"
	@$(GO) fmt ./...

.PHONY: lint
lint: ## Run golangci-lint
	@echo "$(YELLOW)🔍 Running linter...$(NC)"
	@if ! command -v $(GOLANGCI_LINT) >/dev/null; then \
		echo "$(RED)Error: golangci-lint not installed. Run 'go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest'$(NC)"; \
		exit 1; \
	fi
	@$(GOLANGCI_LINT) run --timeout=5m

.PHONY: vet
vet: ## Run go vet
	@echo "$(YELLOW)🔎 Vetting code...$(NC)"
	@$(GO) vet ./...

# Testing
.PHONY: test
test: ## Run tests with coverage
	@echo "$(YELLOW)🧪 Running tests...$(NC)"
	@$(GO) test -v -coverprofile=coverage.out ./...
	@$(GO) tool cover -func=coverage.out

.PHONY: test-html
test-html: ## Run tests and generate HTML coverage report
	@echo "$(YELLOW)🧪 Running tests with HTML coverage...$(NC)"
	@$(GO) test -v -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Coverage report generated: coverage.html$(NC)"

# Building
.PHONY: build
build: ## Build the binary
	@echo "$(YELLOW)🏗️ Building binary...$(NC)"
	@CGO_ENABLED=0 $(GO) build -o $(BINARY) main.go
	@echo "$(GREEN)✅ Binary built: $(BINARY)$(NC)"

.PHONY: build-prod
build-prod: ## Build the binary for production
	@echo "$(YELLOW)🏗️ Building production binary...$(NC)"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags="-s -w" -o $(BINARY) main.go
	@echo "$(GREEN)✅ Production binary built: $(BINARY)$(NC)"

# Running
.PHONY: run
run: ## Run the application in development mode
	@echo "$(YELLOW)🚀 Starting application (development)...$(NC)"
	@APP_ENV=$(ENV) APP_POSTGRES_DSN=$(POSTGRES_DSN) $(GO) run main.go

.PHONY: run-prod
run-prod: ## Run the application in production mode
	@echo "$(YELLOW)🚀 Starting application (production)...$(NC)"
	@APP_ENV=production APP_POSTGRES_DSN=$(POSTGRES_DSN_PROD) ./$(BINARY)

# Database migration
.PHONY: migrate
migrate: ## Run database migrations
	@echo "$(YELLOW)🗄️ Migrating database tables...$(NC)"
	@if [ ! -d "$(MIGRATION_DIR)" ]; then \
		echo "$(RED)Error: Migration directory $(MIGRATION_DIR) not found$(NC)"; \
		exit 1; \
	fi
	@APP_ENV=$(ENV) APP_POSTGRES_DSN=$(POSTGRES_DSN) $(GO) run $(MIGRATION_DIR)/models.go $(MIGRATION_DIR)/mock_data.go $(MIGRATION_DIR)/migrate.go
	@echo "$(GREEN)✅ Database migration completed$(NC)"

# Git operations
.PHONY: git-cycle
git-cycle: ## Run git cycle (add, commit, pull, push) with unrelated histories handling
	@echo "$(YELLOW)🔄 Running git cycle...$(NC)"
	@git add .
	@if git diff --staged --quiet; then \
		echo "$(YELLOW)⚠️ No changes to commit$(NC)"; \
	else \
		git commit -m "update: new feature or fix" || { echo "$(RED)❌ Commit failed$(NC)"; exit 1; }; \
		echo "$(GREEN)✅ Changes committed$(NC)"; \
	fi
	@git fetch origin
	@git pull origin main --rebase --allow-unrelated-histories || { echo "$(RED)❌ Pull failed. Resolve conflicts manually$(NC)"; exit 1; }
	@git push origin main || { echo "$(RED)❌ Push failed$(NC)"; exit 1; }
	@echo "$(GREEN)✅ Git cycle completed$(NC)"

# Docker support
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "$(YELLOW)🐳 Building Docker image...$(NC)"
	@docker build -t $(DOCKER_IMAGE) .
	@echo "$(GREEN)✅ Docker image built: $(DOCKER_IMAGE)$(NC)"

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "$(YELLOW)🐳 Running Docker container...$(NC)"
	@docker run -p $(PORT):$(PORT) \
		-e APP_ENV=$(ENV) \
		-e APP_POSTGRES_DSN=$(POSTGRES_DSN) \
		--name $(APP_NAME) $(DOCKER_IMAGE)
	@echo "$(GREEN)✅ Docker container running on port $(PORT)$(NC)"

.PHONY: docker-stop
docker-stop: ## Stop and remove Docker container
	@echo "$(YELLOW)🐳 Stopping Docker container...$(NC)"
	@docker stop $(APP_NAME) || true
	@docker rm $(APP_NAME) || true
	@echo "$(GREEN)✅ Docker container stopped$(NC)"

# Cleaning
.PHONY: clean
clean: ## Clean build artifacts and coverage files
	@echo "$(YELLOW)🧹 Cleaning up...$(NC)"
	@rm -f $(BINARY) coverage.out coverage.html
	@$(GO) clean
	@echo "$(GREEN)✅ Cleaned$(NC)"

# Mock generation
.PHONY: mock
mock: ## Generate mocks for testing
	@echo "$(YELLOW)🤖 Generating mocks...$(NC)"
	@if ! command -v mockgen >/dev/null; then \
		echo "$(RED)Error: mockgen not installed. Run 'go install github.com/golang/mock/mockgen@latest'$(NC)"; \
		exit 1; \
	fi
	@mockgen -source=main.go -destination=mocks/mock_main.go -package=mocks
	@echo "$(GREEN)✅ Mocks generated$(NC)"

# Security scanning
.PHONY: scan
scan: ## Run security scan with trivy
	@echo "$(YELLOW)🔒 Running security scan...$(NC)"
	@if ! command -v trivy >/dev/null; then \
		echo "$(RED)Error: trivy not installed. Install via 'curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh'$(NC)"; \
		exit 1; \
	fi
	@trivy fs .
	@echo "$(GREEN)✅ Security scan completed$(NC)"

# Database utilities
.PHONY: db-connect
db-connect: ## Connect to local PostgreSQL database
	@echo "$(YELLOW)🗄️ Connecting to PostgreSQL...$(NC)"
	@psql "$(POSTGRES_DSN)"

.PHONY: db-connect-prod
db-connect-prod: ## Connect to production RDS database
	@echo "$(YELLOW)🗄️ Connecting to production RDS...$(NC)"
	@psql "$(POSTGRES_DSN_PROD)"

# Development utilities
.PHONY: watch
watch: ## Run application with hot reload using air
	@echo "$(YELLOW)👀 Starting application with hot reload...$(NC)"
	@if ! command -v air >/dev/null; then \
		echo "$(RED)Error: air not installed. Run 'go install github.com/cosmtrek/air@latest'$(NC)"; \
		exit 1; \
	fi
	@APP_ENV=$(ENV) APP_POSTGRES_DSN=$(POSTGRES_DSN) air

.PHONY: generate
generate: ## Run go generate
	@echo "$(YELLOW)⚙️ Running go generate...$(NC)"
	@$(GO) generate ./...
	@echo "$(GREEN)✅ Generation completed$(NC)"