# terraform-terratest-framework

A reusable Go-based testing framework for Terraform modules, built with [Terratest](https://github.com/gruntwork-io/terratest). This library provides standardized interfaces, helper assertions, and test context scaffolding to make testing Terraform modules simple, consistent, and composable across Caylent’s Terraform ecosystem.

---

## 📁 Folder Structure

```
terraform-terratest-framework/
├── .tool-versions             # ASDF tool versions
├── .golangci.yml              # Go linter config
├── .gitignore                 # Standard excludes
├── Makefile                   # Local dev and test automation
├── go.mod / go.sum            # Go module definition
├── README.md                  # Project documentation
├── internal/                  # Go test support packages
│   ├── assertions/            # Assertion helpers (file content, structure, etc.)
│   └── testctx/               # Test context and config scaffolding
└── tests/
    ├── unit/                  # Unit tests for Go framework logic
    └── functional/            # Example functional usage tests
```

---

## ⚙️ Dev Environment Setup

This project uses [ASDF](https://asdf-vm.com) for consistent tool versioning and supports full Devcontainer-based development.

### 🧱 Using the Caylent Devcontainer

To get started with the Caylent standard devcontainer:

1. **Clone the devcontainer repo and copy it in:**

   ```bash
   git clone https://github.com/caylent-solutions/devcontainer.git
   cp -r devcontainer/.devcontainer ./  # from the root of this repo
   ```

2. **Customize your environment:**

   ```bash
   cp .devcontainer/example-container-env-values.json devcontainer-environment-variables.json
   cp .devcontainer/example-aws-profile-map.json .devcontainer/aws-profile-map.json
   ```

3. **Generate and persist shell exports:**

   ```bash
   python .devcontainer/generate-shell-exports.py export devcontainer-environment-variables.json --output shell.env
   echo "source $(pwd)/shell.env" >> ~/.zshrc  # or .bashrc
   ```

4. **Launch VS Code:**

   ```bash
   code .
   ```

   Then accept the prompt to reopen in container.

> For full setup instructions, see the [Caylent Devcontainer README](https://github.com/caylent-solutions/devcontainer#readme)

---

## 🛠 Makefile Tasks

| Command                 | Description                                      |
|-------------------------|--------------------------------------------------|
| `make configure`        | Full setup: tools, lint, pre-commit, install     |
| `make install-tools`    | Installs ASDF + tools from .tool-versions        |
| `make install`          | Runs `go mod tidy`                               |
| `make pre-commit-install` | Installs Git hooks via pre-commit             |
| `make pre-commit`       | Runs all pre-commit checks across all files     |
| `make lint`             | Runs `golangci-lint` with auto-fix              |
| `make unit-test`        | Runs unit tests in `tests/unit`                 |
| `make functional-test`  | Runs functional tests in `tests/functional`     |
| `make all-test`         | Runs unit and functional tests                  |
| `make clean`            | Cleans Terraform state and generated files      |

---

## ✅ .tool-versions (Managed by ASDF)

```
golang 1.23.9
golangci-lint 2.1.6
pre-commit 4.2.0
```

These tools are automatically installed via `make install-tools`.

---

## 🚀 Contributing

This repo uses **trunk-based development**:

1. Create a short-lived branch from `main`
2. Run `make configure` before committing
3. Add your code and tests
4. Push and open a PR against `main`
5. Ensure all tests, lints, and pre-commit checks pass

---

## 🔐 License

This project is licensed under the Apache 2.0 License.
