.PHONY: setup build run test lint clean migrate migrate-down migrate-force help

SHELL := /bin/bash

GO := go
GOFLAGS ?=
LDFLAGS := -s -w
BINARY_NAME := server

help: ## Show this help message
	@echo "Available commands for Flower API:"
	@awk 'BEGIN {FS = ":.*##"; printf "%-20s %s\n", "Target", "Description"} /^[a-zA-Z_-]+:.*?##/ { printf "%-20s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

setup: ## Initialise Go module and dependencies
	@echo "Initializing Go module..."
	@mkdir -p cmd/server internal/handlers internal/platform/config internal/platform/db internal/migrations tmp
	@if [ ! -f go.mod ]; then \
		$(GO) mod init flower/api; \
	fi
	@$(GO) mod tidy
	@echo "Setup complete!"

build: ## Build the API binary
	@echo "Building API binary..."
	@$(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BINARY_NAME) ./cmd/server
	@echo "Build complete!"

run: ## Run the API locally
	@echo "Running API locally..."
	@$(GO) run ./cmd/server

test: ## Run tests
	@echo "Running tests..."
	@$(GO) test $(GOFLAGS) -p 1 -count=1 ./...

lint: ## Run linting
	@echo "Running linting..."
	@golangci-lint run ./cmd/... ./internal/...

clean: ## Clean up build artefacts
	@echo "Cleaning up build artefacts..."
	@rm -f $(BINARY_NAME)
	@echo "Clean complete!"

deps: ## Update dependencies
	@echo "Updating dependencies..."
	@$(GO) mod tidy
	@$(GO) mod download

migrate: ## Run database migrations
	@echo "Running database migrations..."
	@$(GO) run -ldflags '$(LDFLAGS)' ./cmd/server migrate

migrate-down: ## Rollback last migration
	@echo "Rolling back last migration..."
	@$(GO) run -ldflags '$(LDFLAGS)' ./cmd/server rollback

migrate-force: ## Force a specific migration version
	@echo "Forcing migration version..."
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION parameter is required"; \
		echo "Example: make migrate-force VERSION=1"; \
		exit 1; \
	fi
	@$(GO) run -ldflags '$(LDFLAGS)' ./cmd/server force $(VERSION)
