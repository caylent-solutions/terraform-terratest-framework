# 🚀 TFTest CLI

This command-line tool provides a fun and engaging way to run tests for Terraform modules using the Terraform Test Framework.

## 📥 Installation

You have two options for installing the TFTest CLI:

### Option 1: Install the version used by your module

```bash
# This installs the same version specified in your module's go.mod
cd /path/to/your/terraform-module
go install github.com/caylent-solutions/terraform-terratest-framework/cmd/tftest@$(grep terraform-terratest-framework go.mod | awk '{print $2}')
```

### Option 2: Install from a specific branch

```bash
# This installs from a specific branch (e.g., main, develop)
go install github.com/caylent-solutions/terraform-terratest-framework/cmd/tftest@main
```

## 🎮 Usage

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

# Enable verbose logging with different levels
tftest --verbose DEBUG run        # Most detailed logging
tftest -v INFO format --all       # Standard information (default)
tftest -v WARN run --common       # Only warnings and above
tftest -v ERROR format --all      # Only errors and fatal messages
tftest -v FATAL run               # Only fatal errors
```

## 📁 Directory Structure

For details on the required directory structure, see the [Directory Structure Documentation](../docs/DIRECTORY_STRUCTURE.md).

## 🎯 Commands

- `tftest` - Show help and version information
- `tftest version` - Show version information
- `tftest run` - Run tests for a Terraform module
- `tftest format` - Format and verify Go test code

## 🔧 Global Options

- `--help, -h` - Show help for any command or subcommand
- `--version, -V` - Show version information (root command only)
- `--verbose, -v` - Set verbosity level:
  - `DEBUG` - Detailed information for diagnosing problems
  - `INFO` - General information (default)
  - `WARN` - Warning messages for potential issues
  - `ERROR` - Error messages for operation failures
  - `FATAL` - Critical errors that terminate the program

## 🔧 Options for 'run' command

- `--module-root` - Path to the root of the Terraform module (runs all tests)
- `--example-path` - Specific example to test (verifies both example and test directories exist)
- `--common` - Run only common tests (verifies common directory exists)
- `--timeout` - Test timeout duration (e.g., 20m, 1h). Can also be set via GO_TEST_TIMEOUT environment variable
- `--parallel-fixtures` - Run test fixtures in parallel (default: false)
- `--parallel-tests` - Run tests within each fixture in parallel (default: false)
- `--help, -h` - Show help for the run command

### ⏱️ Test Timeout Behavior

**Important**: Go's test framework applies a default 10-minute timeout to the cumulative duration of ALL tests within a fixture, not per individual test. This means:

- If you have multiple tests in a single fixture, they share the 10-minute timeout
- The timeout applies to the entire test suite execution within that fixture
- Individual tests within the fixture do not get their own separate 10-minute timeout

To override this default timeout:

```bash
# Set timeout via command line flag
tftest run --timeout=20m

# Set timeout via environment variable
export GO_TEST_TIMEOUT=20m
tftest run

# Environment variable takes precedence if both are set
```

## 🔧 Options for 'format' command

- `--all, -A` - Format all Go test files (verifies each example has a matching test directory)
- `--example-path` - Format a specific example's test files (verifies both example and test directories exist)
- `--common` - Format only common test files (verifies common directory exists)
- `--module-root` - Path to the root of the Terraform module
- `--help, -h` - Show help for the format command

## 🧩 How It Works

For detailed information on how the commands work, see the [CLI Usage Documentation](../docs/CLI_USAGE.md#how-it-works).

## 📋 Requirements

- Go 1.23 or later
- A properly structured Terraform module with tests
- AWS credentials (if testing AWS resources)

## 🏗️ Test Architecture Notes

When running tests with fixtures:

- **Serial execution** (`--parallel-fixtures=false --parallel-tests=false`): Tests run sequentially, which is slower but produces minimal output volume to prevent overwhelming AI agents when debugging assistance is required
- **Parallel execution**: Tests run concurrently, which is faster but generates extensive interleaved output that can exceed AI agent context limits

For AI-assisted development workflows, implementing one test per testctx instance is recommended to maintain output volumes within AI agent processing capabilities.