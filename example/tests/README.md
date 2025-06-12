# Terraform Module Tests

This directory contains tests for the Terraform module.

## Directory Structure

```
tests/
├── common/              # Tests that run on all examples
│   ├── module_test.go   # Common tests for all examples
│   └── README.md        # Documentation for common tests
├── helpers/             # Helper functions for tests
│   └── helpers.go       # Common helper functions
├── basic/               # Tests for the basic example
│   └── module_test.go   # Tests specific to the basic example
└── advanced/            # Tests for the advanced example
    └── module_test.go   # Tests specific to the advanced example
```

## Test Categories

### Common Tests

The `common/` directory contains tests that run on all examples:
- Terraform validation
- Terraform formatting
- Required outputs
- File creation
- Various assertion types
- Input validation
- Benchmarking

### Example-Specific Tests

Each example has its own test directory with the same name:
- `basic/`: Tests for the basic example
- `advanced/`: Tests for the advanced example

### Helper Functions

The `helpers/` directory contains reusable helper functions for tests:
- File validation helpers
- Input validation helpers

## Running Tests

Tests can be run using the provided Makefile commands from the root directory:

```bash
# Install dependencies (only needed once)
make install

# Run all tests (with both parallel flags set to false for maximum stability)
make test

# Run only common tests
make test-common

# Run a specific test
go test ./tests/common -run '^TestInputsMatchProvisioned$'

# Format all test files
make format
```

### Using the CLI Directly

You can also run tests directly using the `tftest` CLI:

```bash
# Run all tests with both parallel flags set to false (recommended)
tftest run --parallel-fixtures=false --parallel-tests=false

# Run tests for a specific example
tftest run --example-path basic --parallel-fixtures=false --parallel-tests=false

# Run only common tests
tftest run --common --parallel-fixtures=false --parallel-tests=false
```

The `--parallel-fixtures=false` and `--parallel-tests=false` flags ensure that:
1. Test fixtures run sequentially (not in parallel)
2. Tests within each fixture run sequentially (not in parallel)

This provides maximum stability and is recommended for most use cases.

For more information on the `tftest` CLI tool, see the [CLI Usage Documentation](../../docs/CLI_USAGE.md).