.PHONY: install install-tools update-tools configure \
        unit-test functional-test test test-coverage test-coverage-json \
        lint format clean build-cli release

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
	go test -v ./internal/... ./tftest-cli/...

functional-test:
	go test -v ./tests/functional/...

test: unit-test functional-test
	@echo "All tests passed! 🎉"

test-coverage:
	@echo "Running tests with coverage..."
	@echo "Framework coverage:"
	@go test -coverprofile=coverage.out ./internal/...
	@go tool cover -func=coverage.out
	@echo "\nCLI coverage:"
	@go test -coverprofile=coverage-cli.out ./tftest-cli/...
	@go tool cover -func=coverage-cli.out
	@echo "\nTest coverage report complete! 📊"

test-coverage-json:
	@echo "Running tests with JSON coverage output..."
	@go test -coverprofile=coverage.out ./internal/...
	@go test -coverprofile=coverage-cli.out ./tftest-cli/...
	@echo "{\"framework\": \"$(shell go tool cover -func=coverage.out | grep total | awk '{print $$3}')\", \"cli\": \"$(shell go tool cover -func=coverage-cli.out | grep total | awk '{print $$3}')\"}"

lint:
	@echo "Checking code for linting issues..."
	golangci-lint run
	@echo "Linting check complete! ✅"

format:
	@echo "Formatting code..."
	golangci-lint run --fix
	gofmt -w ./internal ./tftest-cli ./tests
	@echo "Code formatting complete! ✨"

clean:
	rm -rf .terraform terraform.tfstate* *.txt *.json bin/ coverage*.out

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