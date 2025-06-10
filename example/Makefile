.PHONY: install test test-clean test-common format clean clean-all

# Get the version from the root VERSION file
VERSION := $(shell cat ../VERSION | tr -d '\n')

# Install dependencies
install:
	@echo "Installing dependencies..."
	@go mod edit -require github.com/caylent-solutions/terraform-terratest-framework@$(VERSION)
	@go mod tidy
	@go install github.com/caylent-solutions/terraform-terratest-framework/cmd/tftest@$(VERSION)
	@asdf reshim golang || true
	@echo "Dependencies installed successfully"


# Run all tests with clean cache
test:
	@echo "Cleaning Go test cache..."
	@go clean -cache -testcache
	@echo "Running all tests with clean cache..."
	@GOFLAGS="-count=1" tftest run --parallel-fixtures=false --parallel-tests=false

# Run common tests with clean cache
test-common:
	@echo "Cleaning Go test cache..."
	@go clean -cache -testcache
	@echo "Running common tests with clean cache..."
	@. ~/.asdf/asdf.sh && TERRATEST_DISABLE_PARALLEL_TESTS=true go test ./tests/common -v -count=1

# Format all test files
format:
	@echo "Formatting test files..."
	@tftest format --all

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