.PHONY: help test test-coverage test-race test-integration bench lint fmt vet clean build install tidy deps example docs

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Testing
test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race: ## Run tests with race detector
	go test -race ./...

test-integration: ## Run integration tests (requires API key)
	@echo "Running integration tests..."
	go test -tags=integration -v ./...

bench: ## Run benchmarks
	go test -bench=. -benchmem ./...

# Code quality
lint: ## Run linter
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

# Dependencies
deps: ## Download dependencies
	go mod download

tidy: ## Tidy dependencies
	go mod tidy

# Build and install
build: ## Build the package
	go build ./...

install: ## Install the package
	go install ./...

# Examples
example: ## Run the example
	cd examples && go run main.go

# Documentation
docs: ## Generate documentation
	@echo "Generating documentation..."
	godoc -http=:6060 &
	@echo "Documentation server started at http://localhost:6060/pkg/github.com/checkhim/go-sdk/"
	@echo "Press Ctrl+C to stop"

# Development
dev-setup: deps ## Setup development environment
	@echo "Setting up development environment..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development environment ready!"

# Clean
clean: ## Clean build artifacts
	go clean
	rm -f coverage.out coverage.html

# Git hooks
setup-hooks: ## Setup git hooks
	@echo "Setting up git hooks..."
	@mkdir -p .git/hooks
	@echo '#!/bin/sh\nmake test' > .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Git hooks installed!"

# Release preparation
pre-release: clean fmt vet lint test-coverage ## Prepare for release
	@echo "Pre-release checks completed successfully!"

# CI targets
ci-test: ## Run CI tests
	go test -race -coverprofile=coverage.out ./...

ci-lint: ## Run CI linting
	golangci-lint run --timeout=5m

# Security
security: ## Run security checks
	@which gosec > /dev/null || (echo "Installing gosec..." && go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest)
	gosec ./...

# All quality checks
check-all: fmt vet lint test-coverage security ## Run all quality checks
	@echo "All checks passed!"
