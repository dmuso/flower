.PHONY: setup start stop logs clean \
	test lint help \
	frontend-setup frontend-build frontend-dev frontend-serve test-frontend lint-frontend \
	lint-api test-api \
	dev loc migrate migrate-down migrate-force status

SHELL := /bin/bash

-include .env
API_PORT ?= 8180
FRONTEND_PORT ?= 4273

NIX_SHELL ?= nix-shell --run

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*##"; printf "%-20s %s\n", "Target", "Description"} /^[a-zA-Z_-]+:.*?##/ { printf "%-20s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

setup: ## Initialise development environment
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp .env.example .env; fi
	@make -C api setup
	@$(MAKE) frontend-setup

start: ## Start all Docker Compose services
	@echo "Starting all services..."
	@docker-compose up -d --build

stop: ## Stop Docker Compose services and host processes on API/frontend ports
	@echo "Stopping all services..."
	@docker-compose down
	@echo "Killing stray host processes on ports $(API_PORT) and $(FRONTEND_PORT)..."
	@ports="$(API_PORT) $(FRONTEND_PORT)"; for port in $$ports; do \
		pids=$$(lsof -ti tcp:$$port); \
		if [ -n "$$pids" ]; then \
			echo " - Port $$port: terminating $$pids"; \
			kill -9 $$pids >/dev/null 2>&1 || true; \
		fi; \
	done

logs: ## View Docker Compose service logs
	@echo "Viewing service logs..."
	@docker-compose logs -f

clean: ## Clean up containers, volumes, and API build artefacts
	@echo "Cleaning up containers and volumes..."
	@docker-compose down -v --remove-orphans
	@echo "Cleaning API build artefacts..."
	@make -C api clean

test: ## Run API + frontend unit tests
	@echo "Running API and frontend tests..."
	@$(MAKE) test-api
	@$(MAKE) test-frontend
	@echo "All tests passed!"

lint: ## Run linting on all services
	@echo "Running linting..."
	@$(MAKE) lint-api
	@$(MAKE) lint-frontend

frontend-setup:
	@echo "Installing frontend dependencies..."
	@cd frontend && bun install

frontend-build:
	@echo "Building frontend bundle..."
	@if [ -f .env ]; then set -a && . ./.env && set +a; fi && cd frontend && bun run build

frontend-dev:
	@mkdir -p tmp
	@if [ -f .env ]; then set -a && . ./.env && set +a; fi; \
		export FRONTEND_API_URL="$${FRONTEND_API_URL}"; \
		export FRONTEND_PORT="$${FRONTEND_PORT}"; \
		$(NIX_SHELL) "set -o pipefail; overmind start -f Procfile frontend | tee tmp/server.log"

frontend-serve:
	@if [ -f .env ]; then set -a && . ./.env && set +a; fi && cd frontend && bun run serve

test-frontend:
	@cd frontend && bun test --preload ./test-setup.ts --timeout 5000 src scripts

lint-frontend:
	@cd frontend && bun run lint && bun run check

lint-api:
	@cd api && make lint

test-api:
	@if [ -f .env.test ]; then set -a && . ./.env.test && set +a; fi && cd api && make test

dev: ## Start Postgres in Docker and run API + frontend locally with hot reload
	@echo "Killing stray host processes on ports $(API_PORT) and $(FRONTEND_PORT)..."
	@ports="$(API_PORT) $(FRONTEND_PORT)"; for port in $$ports; do \
		pids=$$(lsof -ti tcp:$$port); \
		if [ -n "$$pids" ]; then \
			echo " - Port $$port: terminating $$pids"; \
			kill -9 $$pids >/dev/null 2>&1 || true; \
		fi; \
	done
	@echo "Killing any previous Overmind lock files..."
	@rm -f ./.overmind.sock
	@echo "Ensuring Postgres is running..."
	@docker-compose up --remove-orphans -d db
	@echo "Starting Overmind dev supervisor (Ctrl+C to stop everything)..."
	@mkdir -p tmp
	@rm -f tmp/server.log
	@set -a && . ./.env && set +a; \
		export FRONTEND_API_URL="$${FRONTEND_API_URL}"; \
		export FRONTEND_PORT="$${FRONTEND_PORT}"; \
		$(NIX_SHELL) "set -o pipefail; trap 'overmind quit 2>/dev/null || true' INT TERM EXIT; overmind start -f Procfile 2>&1 | tee tmp/server.log"

migrate: ## Run database migrations
	@echo "Running database migrations..."
	@make -C api migrate

migrate-down: ## Rollback last migration
	@echo "Rolling back last migration..."
	@make -C api migrate-down

migrate-force: ## Force a specific migration version
	@echo "Forcing migration version..."
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION parameter is required"; \
		echo "Example: make migrate-force VERSION=1"; \
		exit 1; \
	fi
	@make -C api migrate-force VERSION=$(VERSION)

status: ## Check Docker Compose service status
	@echo "Checking service status..."
	@docker-compose ps

loc: ## Count lines of code across the codebase
	@scc \
    	--exclude-dir vendor,node_modules,.git,dist,build,tmp \
    	--cocomo-project-type "semi-detached,3.0,1.12,2.5,0.35" \
    	.
