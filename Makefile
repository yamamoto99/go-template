.PHONY: all up down prune fmt migrate-up help test test-setup test-repository test-usecase test-handler test-all test-cleanup

# Default target
.DEFAULT_GOAL := help

# Colors for output
GREEN  := $(shell tput setaf 2)
RESET  := $(shell tput sgr0)

all: up migrate-up ## Start the application, run migrations

up: ## Start the application
	docker compose up -d --build

down: ## Stop the application
	docker compose down -v

migrate-up: ## Run database migrations
	@echo "--- Applying database migrations ---"
	@migrate -path app/infrastructure/db/migrations -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up

fmt: ## Format all Go code files
	@go fmt ./...

test-setup: ## Setup test environment
	docker compose up -d test-db
	@echo "--- Waiting for test-db to be ready ---"
	@until docker exec test-db pg_isready -U postgres >/dev/null 2>&1; do sleep 1; done
	@echo "--- Applying test database migrations ---"
	@migrate -path app/infrastructure/db/migrations -database "postgres://postgres:postgres@localhost:5433/testdb?sslmode=disable" up

test-repository: test-setup ## Run repository tests
	@go test -v ./app/internal/repository/...

test-usecase: ## Run usecase tests
	@go test -v ./app/internal/usecase/...

test-handler: ## Run handler tests
	@go test -v ./app/internal/handler/...

test-all: ## Run all tests
	@go test -v ./app/internal/repository/... ./app/internal/usecase/... ./app/internal/handler/...

help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(RESET) %s\n", $$1, $$2}'
