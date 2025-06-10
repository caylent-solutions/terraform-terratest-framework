.PHONY: build-cli clean clean-coverage configure format functional-test \
        install install-tools lint list-functional-tests pre-commit pre-commit-install \
        release run-specific-functional-test test test-coverage test-coverage-json \
        test-html-coverage unit-test update-tools test-sequential test-parallel-fixtures test-parallel-tests

COVERAGE_DIR := tmp/coverage

build-cli:
	@go build -o bin/tftest -ldflags="-X 'github.com/caylent-solutions/terraform-terratest-framework/cmd/tftest.Version=$(cat VERSION)'" ./cmd/tftest
	@echo "🎉 TFTest CLI built at bin/tftest"

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

.PHONY: format

format:
	@echo "Fixing code formatting and lint issues..."
	@mkdir -p ./bin
	@echo "Building format tool..."
	@go build -o ./bin/format ./scripts/format/main.go
	@./bin/format --ignore="bin"
	@rm -f ./bin/format


functional-test:
	@go run scripts/functional-test/main.go "./tests/functional/..."

install:
	go mod tidy

install-tools:
	@echo "Installing asdf and required development tools..."
	@mkdir -p ./bin
	@echo "Building install-tools..."
	@go build -o ./bin/install-tools ./scripts/install-tools/main.go
	@./bin/install-tools --asdf-version=v0.15.0
	@rm -f ./bin/install-tools

lint:
	@echo "Checking code for linting issues..."
	@mkdir -p ./bin
	@echo "Building lint tool..."
	@go build -o ./bin/lint ./scripts/lint/main.go
	@./bin/lint --ignore="bin" --skip-prefix="github.com/caylent-solutions/terraform-terratest-framework/scripts" || echo "Lint check failed ❌"
	@rm -f ./bin/lint
	@echo "Lint check complete"

list-functional-tests:
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

pre-commit:
	pre-commit run --all-files

pre-commit-install:
	pre-commit install

release:
	@echo "Creating a new release..."
	@mkdir -p ./bin
	@go build -o ./bin/release ./scripts/release/main.go
	@./bin/release $(TYPE)
	@rm -f ./bin/release

run-specific-functional-test:
	@if [ -z "$(FUNCTIONAL_TEST)" ]; then \
		echo "Error: FUNCTIONAL_TEST environment variable must be set."; \
		echo "Usage: FUNCTIONAL_TEST=TestName make run-specific-functional-test"; \
		echo "Run 'make list-functional-tests' to see available tests."; \
		exit 1; \
	fi
	@echo "Running functional test: $(FUNCTIONAL_TEST)"
	@$(MAKE) build-cli
	@cd tests/functional && go test -v -run "^$(FUNCTIONAL_TEST)$$"

test: unit-test functional-test
	@echo "All tests passed! 🎉"

test-coverage:
	@go run scripts/test-coverage/main.go scripts/coverage-groups.json

test-coverage-json:
	@go run scripts/test-coverage-json/main.go scripts/coverage-groups.json

test-coverage-html:
	@go run scripts/test-coverage-html/main.go scripts/coverage-groups.json

unit-test:
	@go run scripts/unit-test/main.go "./tests/unit/..."

update-tools:
	@if [ ! -f ./bin/install-tools ]; then \
		echo "Building install-tools..."; \
		go build -o ./bin/install-tools ./scripts/install-tools.go; \
	fi
	@./bin/install-tools --update