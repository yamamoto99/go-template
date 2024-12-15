.PHONY: all up down seed prune fmt migrate help

# Default target
.DEFAULT_GOAL := help

# Colors for output
GREEN  := $(shell tput setaf 2)
RESET  := $(shell tput sgr0)

all: up migrate seed ## Start the application, run migrations and seed the database

up: ## Start the application
	@docker-compose up -d --build

down: ## Stop the application
	@docker-compose down -v

migrate: ## Run database migrations
	@docker-compose run --rm api go run ../migrate/main.go

seed: ## Seed the database
	@docker-compose run --rm -e PGPASSWORD=postgres api sh -c "psql -h db -U postgres -d postgres -f ../../../seed_data.sql"

prune: ## Remove dangling images
	@docker image prune -f

fmt: ## Format all Go code files
	@go fmt ./...

help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-8s$(RESET) %s\n", $$1, $$2}'
