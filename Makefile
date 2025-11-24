# Neonex Core Makefile

.PHONY: help build run dev test clean install air-install

# Variables
BINARY_NAME=neonex
MAIN_PATH=.
BUILD_DIR=tmp
GO=go

help: ## Show this help message
	@echo "Neonex Core - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "🔨 Building Neonex Core..."
	@$(GO) build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Build complete: ./$(BINARY_NAME)"

build-cli: ## Build the CLI tool
	@echo "🔨 Building Neonex CLI..."
	@$(GO) build -o $(BINARY_NAME).exe ./cmd/neonex
	@echo "✅ CLI built: ./$(BINARY_NAME).exe"

run: ## Run the application
	@echo "🚀 Running Neonex Core..."
	@$(GO) run $(MAIN_PATH)

dev: ## Run with hot reload (requires Air)
	@echo "🔥 Starting development server with hot reload..."
	@air

serve: ## Start development server (alias for dev)
	@$(MAKE) dev

watch: ## Watch and reload (alias for dev)
	@$(MAKE) dev

test: ## Run tests
	@echo "🧪 Running tests..."
	@$(GO) test -v ./...

test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	@$(GO) test -v -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME) $(BINARY_NAME).exe
	@rm -f build-errors.log
	@rm -f coverage.out coverage.html
	@echo "✅ Clean complete"

install: ## Install dependencies
	@echo "📦 Installing dependencies..."
	@$(GO) mod download
	@$(GO) mod tidy
	@echo "✅ Dependencies installed"

air-install: ## Install Air for hot reload
	@echo "📦 Installing Air..."
	@$(GO) install github.com/air-verse/air@latest
	@echo "✅ Air installed"

migrate: ## Run database migrations
	@echo "🗄️  Running migrations..."
	@$(GO) run $(MAIN_PATH) migrate

seed: ## Seed the database
	@echo "🌱 Seeding database..."
	@$(GO) run $(MAIN_PATH) seed

fresh: clean install ## Fresh install (clean + install)
	@echo "✨ Fresh install complete"

fmt: ## Format code
	@echo "🎨 Formatting code..."
	@$(GO) fmt ./...
	@echo "✅ Code formatted"

lint: ## Run linter
	@echo "🔍 Running linter..."
	@golangci-lint run ./...

tidy: ## Tidy go modules
	@echo "🧹 Tidying modules..."
	@$(GO) mod tidy
	@echo "✅ Modules tidied"

update: ## Update dependencies
	@echo "⬆️  Updating dependencies..."
	@$(GO) get -u ./...
	@$(GO) mod tidy
	@echo "✅ Dependencies updated"

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	@docker build -t neonexcore:latest .

docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	@docker run -p 8080:8080 neonexcore:latest

# Development helpers
new-module: ## Create a new module (usage: make new-module name=modulename)
	@./$(BINARY_NAME).exe module create $(name)

module-list: ## List all modules
	@./$(BINARY_NAME).exe module list

# Quick start
start: install build run ## Quick start (install + build + run)

# Production build
prod-build: ## Build for production
	@echo "🏗️  Building for production..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags="-w -s" -o $(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -ldflags="-w -s" -o $(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build -ldflags="-w -s" -o $(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@echo "✅ Production builds complete"

.DEFAULT_GOAL := help
