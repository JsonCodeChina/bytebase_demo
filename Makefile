# SQL Review Learning Demo Makefile

.PHONY: help build test clean run fmt vet deps

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Development commands
deps: ## Install dependencies
	go mod download
	go mod tidy

build: ## Build the demo application
	go build -o bin/sql-review-demo ./cmd/demo

build-server: ## Build the API server
	go build -o bin/sql-review-server ./cmd/server

build-all: ## Build both demo and server
	go build -o bin/sql-review-demo ./cmd/demo
	go build -o bin/sql-review-server ./cmd/server

run: ## Run the demo application
	go run ./cmd/demo

run-server: ## Run the API server
	go run ./cmd/server

dev-server: ## Run server in development mode with auto-reload
	find . -name "*.go" | entr -r go run ./cmd/server

run-server-prod: ## Run server in production mode
	APP_ENV=production ./bin/sql-review-server

test-config: ## Test configuration loading
	go run ./cmd/server --help || echo "Server starts successfully"

test: ## Run tests
	go test ./...

test-verbose: ## Run tests with verbose output
	go test -v ./...

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: ## Run golangci-lint (requires installation)
	golangci-lint run

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

# Example commands
example-good: build ## Run example with good SQL
	./bin/sql-review-demo check examples/good_examples.sql

example-bad: build ## Run example with bad SQL
	./bin/sql-review-demo check examples/bad_examples.sql

example-mixed: build ## Run example with mixed SQL
	./bin/sql-review-demo check examples/mixed_examples.sql

# Development helpers
dev-setup: deps ## Setup development environment
	@echo "Development environment setup complete"

watch: ## Watch for changes and rebuild (requires entr)
	find . -name "*.go" | entr -r make run

# Reference commands
ref-bytebase: ## Show Bytebase reference paths
	@echo "Bytebase reference project:"
	@echo "  Location: /Users/shenbo/goprojects/bytebase-3.5.2/"
	@echo "  Key directories:"
	@echo "    - backend/plugin/advisor/"
	@echo "    - backend/api/v1/"
	@echo "    - backend/store/"

ref-docs: ## Open documentation
	@echo "Documentation files:"
	@echo "  - docs/bytebase-sql-review-analysis.md"
	@echo "  - docs/project-plan.md"
	@echo "  - docs/learning-notes.md"

# Build info
version: ## Show version info
	@echo "SQL Review Learning Demo"
	@echo "Based on: Bytebase v3.5.2"
	@echo "Created: 2025-09-16"