BUILD_DIR=build

.PHONY: all clean test dep lint coverage-report field-align sec-check help

##@ Commands

all: clean dep lint test ## Run all commands

clean: ## Clean build directory
	@echo "Cleaning modcache..."
	@go clean -modcache -i -r
	@echo "[DONE]: Cleaned modcache"

test: ## Run unit tests
	@echo "Running tests..."
	@mkdir -p ${BUILD_DIR}
	@go test -coverprofile=./${BUILD_DIR}/coverage.out -cover -v ./...
	@echo "[DONE]: Tests completed, coverage report generated in coverage.out"
	@echo "Validating race conditions..."
	@go test -race -short ./...
	@echo "[DONE]: Testing completed"

dep: ## Download and tidy dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@echo "[DONE]: Dependencies downloaded"
	@echo "Tidying dependencies..."
	@go mod tidy
	@echo "[DONE]: Dependencies tidied"

lint: ## Run linter
	@echo "Running linter..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo >&2 "golangci-lint is required but not installed. Aborting."; exit 1; }
	@golangci-lint run --timeout 5m
	@echo "[DONE]: Linter completed"

coverage-report: test ## Serve coverage report in browser
	@echo "Serving coverage report..."
	@go tool cover -html=./${BUILD_DIR}/coverage.out  
	@go tool cover -html=./${BUILD_DIR}/coverage.out -o ./${BUILD_DIR}/coverage.html

field-align: ## Run field analysis
	@echo "Running field analysis..."
	@command -v fieldalignment >/dev/null 2>&1 || { echo >&2 "fieldalignment is required but not installed. Aborting."; exit 1; }
	@fieldalignment -fix ./instructions
	@fieldalignment -fix ./types
	@fieldalignment -fix ./
	@echo "[DONE]: Field analysis completed"

sec-check: ## Automatic Static Code Security Analysis
	@echo "Static Code Checking for Security Vulnerabilities..."
	@gosec ./...
	@echo "[DONE]: Security Check completed"

help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <command> \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
