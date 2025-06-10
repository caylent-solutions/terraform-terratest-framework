.PHONY: install test test-clean test-common format clean clean-all

# Get the version from the root VERSION file
VERSION := $(shell cat ../VERSION | tr -d '\n')

install:
	@echo "Installing dependencies..."
	@go mod edit -replace github.com/caylent-solutions/terraform-terratest-framework=../
	@go mod tidy
	@echo "Dependencies installed successfully"

# Run all tests with clean cache
test:
	@echo "Cleaning Go test cache..."
	@go clean -cache -testcache
	@echo "Running all tests with clean cache..."
	@cd .. && go run ./cmd/tftest/main.go run --parallel-fixtures=false --parallel-tests=false

# Format all test files
format:
	@echo "Formatting test files..."
	@cd .. && go run ./cmd/tftest/main.go format --all

# Clean up temporary files
clean:
	@echo "Cleaning up temporary files..."
	@rm -rf .terraform
	@find . -name ".terraform.lock.hcl" -delete
	@find . -name "terraform.tfstate*" -delete
	@find . -name ".terraform.tfstate.lock.info" -delete

# Clean up everything including Go cache
clean-all: clean
	@echo "Cleaning Go cache..."
	@go clean -cache -testcache
	@echo "All caches cleaned"