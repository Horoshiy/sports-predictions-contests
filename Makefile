.PHONY: help setup build test clean docker-up docker-down proto frontend backend

# Default target
help: ## Show this help message
	@echo "Sports Prediction Contests - Development Commands"
	@echo "================================================="
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  %-15s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

setup: ## Setup development environment
	@echo "Setting up development environment..."
	@cp .env.example .env 2>/dev/null || true
	@chmod +x scripts/setup.sh
	@./scripts/setup.sh

build: ## Build all services
	@echo "Building all services..."
	@$(MAKE) backend
	@$(MAKE) frontend

backend: ## Build backend services
	@echo "Building backend services..."
	@cd backend && go work sync
	@cd backend/shared && go mod tidy
	@echo "Backend services ready for development"

frontend: ## Build frontend application
	@echo "Building frontend application..."
	@cd frontend && npm install
	@cd frontend && npm run build

test: ## Run all tests
	@echo "Running tests..."
	@cd backend && go test ./...
	@cd frontend && npm test

proto: ## Generate Protocol Buffers code
	@echo "Generating Protocol Buffers code..."
	@protoc --proto_path=backend/proto \
		--go_out=backend/shared \
		--go-grpc_out=backend/shared \
		backend/proto/*.proto

docker-up: ## Start development environment with Docker
	@echo "Starting development environment..."
	@docker-compose up -d postgres redis

docker-down: ## Stop development environment
	@echo "Stopping development environment..."
	@docker-compose down

docker-services: ## Start all services with Docker
	@echo "Starting all services..."
	@docker-compose --profile services up -d

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf backend/*/bin
	@rm -rf frontend/build
	@rm -rf frontend/node_modules
	@docker-compose down -v

dev: ## Start development servers
	@echo "Starting development servers..."
	@$(MAKE) docker-up
	@echo "Development environment ready!"
	@echo "PostgreSQL: localhost:5432"
	@echo "Redis: localhost:6379"

logs: ## Show Docker logs
	@docker-compose logs -f

status: ## Show service status
	@docker-compose ps

e2e-test: ## Run end-to-end tests with Docker services
	@echo "Running E2E tests..."
	@./scripts/e2e-test.sh

e2e-test-only: ## Run E2E tests (assumes services are running)
	@echo "Running E2E tests against running services..."
	@cd tests/e2e && go test -tags=e2e -v -timeout 5m ./...

seed-small: ## Seed database with small dataset (20 users, 8 contests)
	@echo "Seeding database with small dataset..."
	@./scripts/seed-data.sh --size small

seed-medium: ## Seed database with medium dataset (100 users, 25 contests)
	@echo "Seeding database with medium dataset..."
	@./scripts/seed-data.sh --size medium

seed-large: ## Seed database with large dataset (500 users, 50 contests)
	@echo "Seeding database with large dataset..."
	@./scripts/seed-data.sh --size large

seed-test: ## Test seeding configuration without adding data
	@echo "Testing seeding configuration..."
	@./scripts/seed-data.sh --test

seed-custom: ## Seed with custom parameters (use SEED_SIZE and SEED_VALUE env vars)
	@echo "Seeding with custom parameters..."
	@./scripts/seed-data.sh --size $(or $(SEED_SIZE),small) --seed $(or $(SEED_VALUE),42)
