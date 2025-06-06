.PHONY: build-cli clean clean-coverage configure format functional-test \
        install install-tools lint list-functional-tests pre-commit pre-commit-install \
        release run-specific-functional-test test test-coverage test-coverage-json \
        test-html-coverage unit-test update-tools

COVERAGE_DIR := tmp/coverage

build-cli:
	go build -o bin/tftest -ldflags="-X 'github.com/caylent-solutions/terraform-terratest-framework/tftest-cli.Version=0.2.0'" ./tftest-cli
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

format:
	@echo "Fixing code formatting and lint issues..."
	@if [ ! -f ./bin/format-safely ]; then \
		echo "Building format tool..."; \
		go build -o ./bin/format-safely ./scripts/format-safely.go; \
	fi
	./bin/format-safely
	@rm -f ./bin/format-safely

functional-test:
	@echo "Running functional tests (verbose output)..."
	@$(MAKE) build-cli
	@mkdir -p $(COVERAGE_DIR)
	@go test -v -covermode=atomic -coverprofile=$(COVERAGE_DIR)/functional.out -coverpkg=./internal/... ./tests/functional/... ; \
	echo "\nFunctional Test Coverage of Internal Packages:" ; \
	go tool cover -func=$(COVERAGE_DIR)/functional.out | grep total: | awk '{print "  " $1 " " $2 " " $3}' ; \
	echo "\nSummarizing functional test results..." ; \
	go test -json ./tests/functional/... | go run scripts/test-summary.go "Functional Test Summary" || true

install:
	go mod tidy

install-tools:
	@echo "Installing asdf and required development tools..."
	@if [ ! -f ./bin/install-tools ]; then \
		echo "Building install-tools..."; \
		go build -o ./bin/install-tools ./scripts/install-tools.go; \
	fi
	@./bin/install-tools --asdf-version=v0.15.0

lint:
	@echo "Checking code for linting issues..."
	@if [ ! -f ./bin/lint-all ]; then \
		echo "Building lint tool..."; \
		go build -o ./bin/lint-all ./scripts/lint-all.go; \
	fi
	./bin/lint-all || echo "Lint check failed ❌"
	@rm -f ./bin/lint-all
	@echo "Lint check complete"

list-functional-tests:
	@echo "Listing functional tests:"
	@echo "------------------------"
	@echo "Top-level tests (use these with FUNCTIONAL_TEST=TestName):"
	@cd tests/functional && go test -list "^Test" ./... | grep -v "^ok" | sort | while read -r test_name; do \
		if [ ! -z "$test_name" ]; then \
			echo "  - $test_name"; \
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
	@if [ ! -f ./bin/bump-version ]; then \
		echo "Building bump-version tool..."; \
		go build -o ./bin/bump-version ./scripts/bump-version.go; \
	fi
	@if [ -z "$(TYPE)" ]; then \
		./bin/bump-version; \
	else \
		./bin/bump-version $(TYPE); \
	fi
	@echo "Release created! 🚀"
	@echo "Run 'git push && git push --tags' to publish the release"

run-specific-functional-test:
	@if [ -z "$(FUNCTIONAL_TEST)" ]; then \
		echo "Error: FUNCTIONAL_TEST environment variable must be set."; \
		echo "Usage: FUNCTIONAL_TEST=TestName make run-specific-functional-test"; \
		echo "Run 'make list-functional-tests' to see available tests."; \
		exit 1; \
	fi
	@echo "Running functional test: $(FUNCTIONAL_TEST)"
	@$(MAKE) build-cli
	@cd tests/functional && go test -v -run "^$(FUNCTIONAL_TEST)$"

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

	@echo "\n🧪 Functional test coverage:"
	@go test -covermode=atomic -coverprofile=$(COVERAGE_DIR)/coverage-functional.out -coverpkg=./internal/...,./tests/functional/... ./tests/functional/... || true
	@go tool cover -func=$(COVERAGE_DIR)/coverage-functional.out | tee $(COVERAGE_DIR)/coverage-functional-summary.log

	@echo "\n🧪 Unit test helpers coverage:"
	@go test -covermode=atomic -coverprofile=$(COVERAGE_DIR)/coverage-unit-helpers.out ./tests/unit/... || true
	@go tool cover -func=$(COVERAGE_DIR)/coverage-unit-helpers.out | tee $(COVERAGE_DIR)/coverage-unit-helpers-summary.log

	@echo "\n🔗 Merging coverage profiles..."
	@echo "mode: atomic" > $(COVERAGE_DIR)/coverage.out
	@tail -n +2 $(COVERAGE_DIR)/coverage-framework.out >> $(COVERAGE_DIR)/coverage.out
	@tail -n +2 $(COVERAGE_DIR)/coverage-cli.out >> $(COVERAGE_DIR)/coverage.out
	@tail -n +2 $(COVERAGE_DIR)/coverage-functional.out >> $(COVERAGE_DIR)/coverage.out
	@tail -n +2 $(COVERAGE_DIR)/coverage-unit-helpers.out >> $(COVERAGE_DIR)/coverage.out
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out > $(COVERAGE_DIR)/coverage-summary.log

	@echo "\n📊 Test Coverage Summary:"
	@echo "🧱 Framework Total Coverage:"
	@grep total: $(COVERAGE_DIR)/coverage-framework-summary.log | awk '{print "  " $$1 " " $$2 " " $$3}'
	@echo "\n🧪 CLI Total Coverage:"
	@grep total: $(COVERAGE_DIR)/coverage-cli-summary.log | awk '{print "  " $$1 " " $$2 " " $$3}'
	@echo "\n🧪 Functional Test Coverage:"
	@grep total: $(COVERAGE_DIR)/coverage-functional-summary.log | awk '{print "  " $$1 " " $$2 " " $$3}'
	@echo "\n🧪 Unit Test Helpers Coverage:"
	@grep total: $(COVERAGE_DIR)/coverage-unit-helpers-summary.log | awk '{print "  " $$1 " " $$2 " " $$3}'
	@echo "\n🧩 Combined Total Coverage (All Components):"
	@grep total: $(COVERAGE_DIR)/coverage-summary.log | awk '{print "  " $$1 " " $$2 " " $$3}'

test-coverage-json:
	@mkdir -p $(COVERAGE_DIR)
	@go test -coverprofile=$(COVERAGE_DIR)/coverage-framework.out ./internal/... > /dev/null 2>&1
	@go test -coverprofile=$(COVERAGE_DIR)/coverage-cli.out ./tftest-cli/... > /dev/null 2>&1
	@go test -coverprofile=$(COVERAGE_DIR)/coverage-functional.out -coverpkg=./internal/...,./tests/functional/... ./tests/functional/... > /dev/null 2>&1
	@go test -coverprofile=$(COVERAGE_DIR)/coverage-unit-helpers.out ./tests/unit/... > /dev/null 2>&1
	
	@go tool cover -func=$(COVERAGE_DIR)/coverage-framework.out > $(COVERAGE_DIR)/coverage-framework-summary.log 2>/dev/null
	@go tool cover -func=$(COVERAGE_DIR)/coverage-cli.out > $(COVERAGE_DIR)/coverage-cli-summary.log 2>/dev/null
	@go tool cover -func=$(COVERAGE_DIR)/coverage-functional.out > $(COVERAGE_DIR)/coverage-functional-summary.log 2>/dev/null
	@go tool cover -func=$(COVERAGE_DIR)/coverage-unit-helpers.out > $(COVERAGE_DIR)/coverage-unit-helpers-summary.log 2>/dev/null
	
	@echo "mode: set" > $(COVERAGE_DIR)/coverage.out
	@tail -n +2 $(COVERAGE_DIR)/coverage-framework.out >> $(COVERAGE_DIR)/coverage.out 2>/dev/null || true
	@tail -n +2 $(COVERAGE_DIR)/coverage-cli.out >> $(COVERAGE_DIR)/coverage.out 2>/dev/null || true
	@tail -n +2 $(COVERAGE_DIR)/coverage-functional.out >> $(COVERAGE_DIR)/coverage.out 2>/dev/null || true
	@tail -n +2 $(COVERAGE_DIR)/coverage-unit-helpers.out >> $(COVERAGE_DIR)/coverage.out 2>/dev/null || true
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out > $(COVERAGE_DIR)/coverage-summary.log 2>/dev/null
	
	@go run scripts/coverage-json.go $(COVERAGE_DIR)

test-coverage-html:
	@mkdir -p $(COVERAGE_DIR)
	@echo "🔍 Generating HTML coverage report..."
	@go test -covermode=atomic -coverprofile=$(COVERAGE_DIR)/coverage.out ./internal/... ./tftest-cli/... ./tests/functional/... ./tests/unit/... || true
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "🌐 HTML coverage report generated at $(COVERAGE_DIR)/coverage.html"

unit-test:
	@echo "Cleaning Go test cache..."
	@go clean -testcache
	@echo "Running unit tests (verbose output)..."
	@go test -v ./internal/... ./tftest-cli/... ; \
	echo "\nSummarizing unit test results..." ; \
	go test -json ./internal/... ./tftest-cli/... | go run scripts/test-summary.go "Unit Test Summary" || true

update-tools:
	@if [ ! -f ./bin/install-tools ]; then \
		echo "Building install-tools..."; \
		go build -o ./bin/install-tools ./scripts/install-tools.go; \
	fi
	@./bin/install-tools --update