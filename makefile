
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

YELLOW := \033[1;33m
GREEN := \033[1;32m
RED := \033[1;31m
NC := \033[0m

.PHONY: all
all: deps fmt lint test build

.PHONY: help
help:
	@echo "$(YELLOW)Makefile for $(APP_NAME)$(NC)"
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

.PHONY: deps
deps:
	@echo "$(YELLOW)📦 Installing dependencies...$(NC)"
	@$(GO) mod tidy
	@$(GO) mod vendor

.PHONY: deps-update
deps-update:
	@echo "$(YELLOW)🔄 Updating dependencies...$(NC)"
	@$(GO) get -u ./...
	@$(GO) mod tidy
	@$(GO) mod vendor


.PHONY: fmt
fmt:
	@echo "$(YELLOW)🖌️ Formatting code...$(NC)"
	@$(GO) fmt ./...

.PHONY: lint
lint:
	@echo "$(YELLOW)🔍 Running linter...$(NC)"
	@if ! command -v $(GOLANGCI_LINT) >/dev/null; then \
		echo "$(RED)Error: golangci-lint not installed. Run 'go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest'$(NC)"; \
		exit 1; \
	fi
	@$(GOLANGCI_LINT) run --timeout=5m

.PHONY: vet
vet:
	@echo "$(YELLOW)🔎 Vetting code...$(NC)"
	@$(GO) vet ./...

.PHONY: test
test:
	@echo "$(YELLOW)🧪 Running tests...$(NC)"
	@$(GO) test -v -coverprofile=coverage.out ./...
	@$(GO) tool cover -func=coverage.out

.PHONY: test-html
test-html:
	@echo "$(YELLOW)🧪 Running tests with HTML coverage...$(NC)"
	@$(GO) test -v -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Coverage report generated: coverage.html$(NC)"

.PHONY: build
build:
	@echo "$(YELLOW)🏗️ Building binary...$(NC)"
	@CGO_ENABLED=0 $(GO) build -o $(BINARY) main.go
	@echo "$(GREEN)✅ Binary built: $(BINARY)$(NC)"

.PHONY: build-prod
build-prod:
	@echo "$(YELLOW)🏗️ Building production binary...$(NC)"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags="-s -w" -o $(BINARY) main.go
	@echo "$(GREEN)✅ Production binary built: $(BINARY)$(NC)"

.PHONY: run
run:
	@echo "$(YELLOW)🚀 Starting application (development)...$(NC)"
	@APP_ENV=$(ENV) APP_POSTGRES_DSN=$(POSTGRES_DSN) $(GO) run main.go

.PHONY: run-prod
run-prod:
	@echo "$(YELLOW)🚀 Starting application (production)...$(NC)"
	@APP_ENV=production APP_POSTGRES_DSN=$(POSTGRES_DSN_PROD) ./$(BINARY)


.PHONY: migrate
migrate:
	@echo "$(YELLOW)🗄️ Migrating database tables...$(NC)"
	@if [ ! -d "$(MIGRATION_DIR)" ]; then \
		echo "$(RED)Error: Migration directory $(MIGRATION_DIR) not found$(NC)"; \
		exit 1; \
	fi
	@APP_ENV=$(ENV) APP_POSTGRES_DSN=$(POSTGRES_DSN) $(GO) run $(MIGRATION_DIR)/models.go $(MIGRATION_DIR)/mock_data.go $(MIGRATION_DIR)/migrate.go
	@echo "$(GREEN)✅ Database migration completed$(NC)"

.PHONY: git-cycle
git-cycle:
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


.PHONY: docker-build
docker-build:
	@echo "$(YELLOW)🐳 Building Docker image...$(NC)"
	@docker build -t $(DOCKER_IMAGE) .
	@echo "$(GREEN)✅ Docker image built: $(DOCKER_IMAGE)$(NC)"

.PHONY: docker-run
docker-run:
	@echo "$(YELLOW)🐳 Running Docker container...$(NC)"
	@docker run -p $(PORT):$(PORT) \
		-e APP_ENV=$(ENV) \
		-e APP_POSTGRES_DSN=$(POSTGRES_DSN) \
		--name $(APP_NAME) $(DOCKER_IMAGE)
	@echo "$(GREEN)✅ Docker container running on port $(PORT)$(NC)"

.PHONY: docker-stop
docker-stop:
	@echo "$(YELLOW)🐳 Stopping Docker container...$(NC)"
	@docker stop $(APP_NAME) || true
	@docker rm $(APP_NAME) || true
	@echo "$(GREEN)✅ Docker container stopped$(NC)"


.PHONY: clean
clean:
	@echo "$(YELLOW)🧹 Cleaning up...$(NC)"
	@rm -f $(BINARY) coverage.out coverage.html
	@$(GO) clean
	@echo "$(GREEN)✅ Cleaned$(NC)"

.PHONY: mock
mock:
	@echo "$(YELLOW)🤖 Generating mocks...$(NC)"
	@if ! command -v mockgen >/dev/null; then \
		echo "$(RED)Error: mockgen not installed. Run 'go install github.com/golang/mock/mockgen@latest'$(NC)"; \
		exit 1; \
	fi
	@mockgen -source=main.go -destination=mocks/mock_main.go -package=mocks
	@echo "$(GREEN)✅ Mocks generated$(NC)"

.PHONY: scan
scan:
	@echo "$(YELLOW)🔒 Running security scan...$(NC)"
	@if ! command -v trivy >/dev/null; then \
		echo "$(RED)Error: trivy not installed. Install via 'curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh'$(NC)"; \
		exit 1; \
	fi
	@trivy fs .
	@echo "$(GREEN)✅ Security scan completed$(NC)"


.PHONY: db-connect
db-connect:
	@echo "$(YELLOW)🗄️ Connecting to PostgreSQL...$(NC)"
	@psql "$(POSTGRES_DSN)"

.PHONY: db-connect-prod
db-connect-prod:
	@echo "$(YELLOW)🗄️ Connecting to production RDS...$(NC)"
	@psql "$(POSTGRES_DSN_PROD)"


.PHONY: watch
watch:
	@echo "$(YELLOW)👀 Starting application with hot reload...$(NC)"
	@if ! command -v air >/dev/null; then \
		echo "$(RED)Error: air not installed. Run 'go install github.com/cosmtrek/air@latest'$(NC)"; \
		exit 1; \
	fi
	@APP_ENV=$(ENV) APP_POSTGRES_DSN=$(POSTGRES_DSN) air

.PHONY: generate
generate:
	@echo "$(YELLOW)⚙️ Running go generate...$(NC)"
	@$(GO) generate ./...
	@echo "$(GREEN)✅ Generation completed$(NC)"

.PHONY: user-pool
user-pool:
	@aws cognito-idp create-user-pool \
		--pool-name my-user-pool \
		--schema '[ \
			{ "Name": "email", "AttributeDataType": "String", "Mutable": true, "Required": true }, \
			{ "Name": "given_name", "AttributeDataType": "String", "Mutable": true, "Required": false }, \
			{ "Name": "family_name", "AttributeDataType": "String", "Mutable": true, "Required": false }, \
			{ "Name": "custom:country", "AttributeDataType": "String", "StringAttributeConstraints": { "MinLength": "2", "MaxLength": "56" }, "Mutable": true, "Required": false }, \
			{ "Name": "custom:bio_hash", "AttributeDataType": "String", "StringAttributeConstraints": { "MinLength": "1", "MaxLength": "128" }, "Mutable": true, "Required": false } \
		]'
