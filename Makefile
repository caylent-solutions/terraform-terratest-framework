.PHONY: install install-tools update-tools configure \
        unit-test functional-test test test-coverage test-coverage-json test-html-coverage \
        clean clean-coverage \
        lint format build-cli release \
        list-functional-tests run-specific-functional-test

COVERAGE_DIR := tmp/coverage

install:
	go mod tidy

install-tools:
	@if [ "$$DEVCONTAINER" = "true" ]; then \
		echo "Caylent Devcontainer detected. Tools already installed."; \
		echo "Running update-tools to ensure everything is up to date..."; \
		$(MAKE) update-tools; \
	else \
		HEADLESS_FLAG=""; \
		for arg in $(MAKECMDGOALS); do \
			if [ "$$arg" = "--headless" ]; then \
				HEADLESS_FLAG="true"; \
			fi; \
		done; \
		if ! command -v asdf >/dev/null 2>&1; then \
			if [ "$$HEADLESS_FLAG" = "true" ]; then \
				echo "Installing asdf headlessly..."; \
				git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.15.0; \
				. ~/.asdf/asdf.sh; \
			else \
				read -p "asdf is not installed. Install asdf now? [y/N]: " yn; \
				case $$yn in \
					[Yy]*) \
						echo "Installing asdf..."; \
						git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.15.0; \
						. ~/.asdf/asdf.sh; \
						;; \
					*) \
						echo "asdf install aborted."; \
						exit 1; \
						;; \
				esac; \
			fi; \
		else \
			echo "asdf already installed."; \
		fi && \
		for plugin in $$(cut -d' ' -f1 .tool-versions); do \
			asdf plugin add $$plugin || true; \
		done && \
		asdf install && \
		asdf reshim; \
	fi

update-tools:
	@echo "Checking and updating asdf tools..."
	@if ! command -v asdf >/dev/null 2>&1; then \
		echo "asdf not found. Please run 'make install-tools' first."; \
		exit 1; \
	fi
	@for plugin in $$(cut -d' ' -f1 .tool-versions); do \
		echo "Ensuring plugin $$plugin is installed..."; \
		asdf plugin add $$plugin 2>/dev/null || true; \
	done
	@echo "Installing/updating tools from .tool-versions..."
	@asdf install
	@asdf reshim
	@echo "All tools are up to date."

pre-commit-install:
	pre-commit install

pre-commit:
	pre-commit run --all-files

unit-test:
	@echo "Running unit tests (verbose output)..."
	@go test -v ./internal/... ./tftest-cli/... ; \
	echo "\nSummarizing unit test results..." ; \
	go test -json ./internal/... ./tftest-cli/... | go run scripts/test-summary.go "Unit Test Summary" || true

functional-test:
	@echo "Running functional tests (verbose output)..."
	@$(MAKE) build-cli
	@go test -v ./tests/functional/... ; \
	echo "\nSummarizing functional test results..." ; \
	go test -json ./tests/functional/... | go run scripts/test-summary.go "Functional Test Summary" || true

test: unit-test functional-test
	@echo "All tests passed! 🎉"

test-coverage:
	@mkdir -p $(COVERAGE_DIR)
	@echo "🔍 Running tests with coverage..."

	@echo "\n🧱 Framework coverage:"
	@go test -covermode=atomic -coverprofile=$(COVERAGE_DIR)/coverage-framework.out ./internal/... || true
	@go tool cover -func=$(COVERAGE_DIR)/coverage-framework.out | tee $(COVERAGE_DIR)/coverage-framework-summary.log

	@echo "\n🧪 CLI coverage:"
	@go test -covermode=atomic -coverprofile=$(COVERAGE_DIR)/coverage-cli.out ./tftest-cli/... || true
	@go tool cover -func=$(COVERAGE_DIR)/coverage-cli.out | tee $(COVERAGE_DIR)/coverage-cli-summary.log

	@echo "\n🔗 Merging coverage profiles..."
	@echo "mode: atomic" > $(COVERAGE_DIR)/coverage.out
	@tail -n +2 $(COVERAGE_DIR)/coverage-framework.out >> $(COVERAGE_DIR)/coverage.out
	@tail -n +2 $(COVERAGE_DIR)/coverage-cli.out >> $(COVERAGE_DIR)/coverage.out
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out > $(COVERAGE_DIR)/coverage-summary.log

	@echo "\n📊 Test Coverage Summary:"
	@echo "🧱 Framework Total Coverage:"
	@grep total: $(COVERAGE_DIR)/coverage-framework-summary.log | awk '{printf "  %-10s %-15s %s\n", $$1, $$2, $$3}'
	@echo "\n🧪 CLI Total Coverage:"
	@grep total: $(COVERAGE_DIR)/coverage-cli-summary.log | awk '{printf "  %-10s %-15s %s\n", $$1, $$2, $$3}'
	@echo "\n🧩 Combined Total Coverage (Framework + CLI):"
	@grep total: $(COVERAGE_DIR)/coverage-summary.log | awk '{printf "  %-10s %-15s %s\n", $$1, $$2, $$3}'

test-coverage-json:
	@mkdir -p $(COVERAGE_DIR)
	@echo "Running tests with JSON coverage output..."
	@go test -coverprofile=$(COVERAGE_DIR)/coverage-framework.out ./internal/...
	@go test -coverprofile=$(COVERAGE_DIR)/coverage-cli.out ./tftest-cli/...
	@echo "{\"framework\": \"$(shell go tool cover -func=$(COVERAGE_DIR)/coverage-framework.out | grep total | awk '{print $$3}')\", \"cli\": \"$(shell go tool cover -func=$(COVERAGE_DIR)/coverage-cli.out | grep total | awk '{print $$3}')\"}"

test-html-coverage:
	@mkdir -p $(COVERAGE_DIR)
	@echo "🔍 Generating HTML coverage report..."
	@go test -covermode=atomic -coverprofile=$(COVERAGE_DIR)/coverage.out ./internal/... ./tftest-cli/... || true
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "🌐 HTML coverage report generated at $(COVERAGE_DIR)/coverage.html"

lint:
	@echo "Checking code for linting issues..."
	./scripts/lint-all.sh || echo "Lint check failed ❌"
	@echo "Lint check complete"

format:
	@echo "Fixing code formatting and lint issues..."
	./scripts/format-safely.sh
	@echo "Format complete ✨"

clean:
	rm -rf .terraform terraform.tfstate* *.txt *.json bin/ *.out

clean-coverage:
	@echo "🧹 Cleaning coverage artifacts..."
	@rm -rf $(COVERAGE_DIR)
	@echo "✅ Coverage files removed."

configure:
	$(MAKE) install-tools
	$(MAKE) install
	$(MAKE) pre-commit-install
	$(MAKE) pre-commit

build-cli:
	go build -o bin/tftest -ldflags="-X 'github.com/caylent-solutions/terraform-terratest-framework/tftest-cli.Version=v0.1.0'" ./tftest-cli
	@echo "🎉 TFTest CLI built at bin/tftest"

release:
	@echo "Creating a new release..."
	@if [ -z "$(TYPE)" ]; then \
		./scripts/bump-version.sh; \
	else \
		./scripts/bump-version.sh $(TYPE); \
	fi
	@echo "Release created! 🚀"
	@echo "Run 'git push && git push --tags' to publish the release"

# List all functional tests with descriptions
list-functional-tests:
	@echo "Loading asdf tools..."
	@source ~/.asdf/asdf.sh 2>/dev/null || . ~/.asdf/asdf.sh 2>/dev/null || echo "Warning: Could not load asdf"
	@echo "Listing functional tests:"
	@echo "------------------------"
	@echo "Top-level tests (use these with FUNCTIONAL_TEST=TestName):"
	@cd tests/functional && go test -list "^Test" ./... | grep -v "^ok" | sort | while read -r test_name; do \
		if [ ! -z "$$test_name" ]; then \
			echo "  - $$test_name"; \
		fi; \
	done
	@echo ""
	@echo "Note: Some tests contain subtests that run multiple examples that are not shown here."

# Run a specific functional test
run-specific-functional-test:
	@if [ -z "$(FUNCTIONAL_TEST)" ]; then \
		echo "Error: FUNCTIONAL_TEST environment variable must be set."; \
		echo "Usage: FUNCTIONAL_TEST=TestName make run-specific-functional-test"; \
		echo "Run 'make list-functional-tests' to see available tests."; \
		exit 1; \
	fi
	@echo "Loading asdf tools..."
	@source ~/.asdf/asdf.sh 2>/dev/null || . ~/.asdf/asdf.sh 2>/dev/null || echo "Warning: Could not load asdf"
	@echo "Running functional test: $(FUNCTIONAL_TEST)"
	@$(MAKE) build-cli
	@cd tests/functional && go test -v -run "^$(FUNCTIONAL_TEST)$$"