# TFTest CLI Usage

This document describes how to use the TFTest CLI tool.

## Installation

You have two options for installing the TFTest CLI:

### Option 1: Install the version used by your module

```bash
# This installs the same version specified in your module's go.mod
cd /path/to/your/terraform-module
go install github.com/caylent-solutions/terraform-terratest-framework/tftest-cli@$(grep terraform-terratest-framework go.mod | awk '{print $2}')
```

### Option 2: Install from a specific branch

```bash
# This installs from a specific branch (e.g., main, develop)
go install github.com/caylent-solutions/terraform-terratest-framework/tftest-cli@main
```

## Basic Usage

```bash
# Show help
tftest --help
tftest -h

# Show version
tftest --version
tftest version

# Run all tests
tftest run

# Run tests for a specific example
tftest run --example-path vpc

# Run common tests only
tftest run --common

# Run all tests in a specific module
tftest run --module-root /path/to/terraform-module

# Format and verify all Go test files
tftest format --all

# Format and verify a specific example's test files
tftest format --example-path vpc

# Format and verify common test files
tftest format --common
```

## Logging Levels

Enable verbose logging with different levels:

```bash
tftest --verbose DEBUG run        # Most detailed logging
tftest -v INFO format --all       # Standard information (default)
tftest -v WARN run --common       # Only warnings and above
tftest -v ERROR format --all      # Only errors and fatal messages
tftest -v FATAL run               # Only fatal errors
```

## Commands

- `tftest` - Show help and version information
- `tftest version` - Show version information
- `tftest run` - Run tests for a Terraform module
- `tftest format` - Format and verify Go test code

## Global Options

- `--help, -h` - Show help for any command or subcommand
- `--version, -V` - Show version information (root command only)
- `--verbose, -v` - Set verbosity level:
  - `DEBUG` - Detailed information for diagnosing problems
  - `INFO` - General information (default)
  - `WARN` - Warning messages for potential issues
  - `ERROR` - Error messages for operation failures
  - `FATAL` - Critical errors that terminate the program

## Options for 'run' command

- `--module-root` - Path to the root of the Terraform module (runs all tests)
- `--example-path` - Specific example to test (verifies both example and test directories exist)
- `--common` - Run only common tests (verifies common directory exists)
- `--help, -h` - Show help for the run command

## Options for 'format' command

- `--all, -A` - Format all Go test files (verifies each example has a matching test directory)
- `--example-path` - Format a specific example's test files (verifies both example and test directories exist)
- `--common` - Format only common test files (verifies common directory exists)
- `--module-root` - Path to the root of the Terraform module
- `--help, -h` - Show help for the format command

## How It Works

### Run Command

1. Verifies your module follows the expected directory structure
2. When using `--example-path`, verifies both the example and its test directory exist
3. When using `--common`, verifies the common test directory exists
4. Runs the appropriate tests using the Go test command
5. Displays the test results in real-time with colorful output

### Format Command

When run with `--all`:
1. Verifies each example in `examples/` has a corresponding test directory in `tests/`
2. Includes optional directories (`common/`, `helpers/`) if they exist
3. Exits with error if any required test directory is missing

When run with `--example-path`:
1. Verifies both the example directory and its test directory exist
2. Exits with error if either directory is missing

When run with `--common`:
1. Verifies the common test directory exists
2. Exits with error if the directory is missing

For all paths:
1. Checks Go test files for formatting issues using `gofmt`
2. Automatically fixes formatting issues when possible
3. Runs `go vet` to check for code issues
4. Exits with non-zero status if any issues couldn't be fixed automatically

## Required Environment Variables

- **AWS Authentication**:
  ```bash
  export AWS_ACCESS_KEY_ID=your_access_key
  export AWS_SECRET_ACCESS_KEY=your_secret_key
  export AWS_REGION=us-west-2
  ```
  
  Or using AWS SSO:
  ```bash
  aws sso login --profile your-profile
  export AWS_PROFILE=your-profile
  ```

- **Idempotency Control** (optional):
  ```bash
  # To disable idempotency testing
  export TERRATEST_IDEMPOTENCY=false
  ```