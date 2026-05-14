.PHONY: dev fmt migrate-up migrate-down migrate-create help test test-setup test-repository test-usecase test-handler test-all

# Default target
.DEFAULT_GOAL := help

# Colors for output
GREEN  := $(shell tput setaf 2)
RESET  := $(shell tput sgr0)

dev: ## Run application with hot reload (air)
	@air

migrate-up: ## Run database migrations
	@echo "--- Applying database migrations ---"
	@migrate -path db/migrations -database "postgres://postgres:postgres@db:5432/postgres?sslmode=disable" up

migrate-down: ## Rollback last database migration (override steps with N=2)
	@echo "--- Rolling back database migrations ---"
	@migrate -path db/migrations -database "postgres://postgres:postgres@db:5432/postgres?sslmode=disable" down $(or $(N),1)

migrate-create: ## Create new migration file pair (usage: make migrate-create NAME=add_xxx_table)
	@if [ -z "$(NAME)" ]; then echo "Error: NAME is required. Usage: make migrate-create NAME=add_xxx_table"; exit 1; fi
	@migrate create -ext sql -dir db/migrations -seq $(NAME)

fmt: ## Format all Go code files
	@go fmt ./...

test-setup: ## Setup test environment
	@echo "--- Applying test database migrations ---"
	@migrate -path db/migrations -database "postgres://postgres:postgres@test-db:5432/testdb?sslmode=disable" up

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
