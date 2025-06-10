# Terraform Terratest Framework

A framework for testing Terraform modules using Terratest with a focus on functional testing and parallel execution of examples.

## Overview

This framework provides a structured way to test Terraform modules by:

1. Automatically discovering and running all examples in parallel
2. Testing idempotency of Terraform code
3. Providing common assertions for Terraform outputs and resources
4. Supporting custom test functions for advanced verification
5. Simplifying test setup and execution

## Features

- **Parallel Example Testing**: Run all examples in the `examples/` directory in parallel
- **Flexible Parallelism Control**: Control parallelism at both test fixture and individual test levels
- **Idempotency Testing**: Verify that Terraform code is idempotent by running a plan after apply
- **Common Assertions**: Pre-built assertions for common testing scenarios
- **Custom Tests**: Support for custom test functions to verify specific resource behaviors
- **Automatic Discovery**: Automatically finds and runs tests on all examples without manual configuration
- **Configurable**: Easily customize test configurations for each example
- **Environment Control**: Disable idempotency testing with the `TERRATEST_IDEMPOTENCY=false` environment variable

## Development Environment Setup

### Using Caylent Devcontainer (Recommended)

The easiest way to set up your development environment is using the Caylent Devcontainer:

1. **Install the Caylent Devcontainer CLI**:
   ```bash
   pip install caylent-devcontainer-cli
   ```

2. **Set up the devcontainer in your project**:
   ```bash
   cdevcontainer setup-devcontainer .
   ```
   Follow the interactive prompts to configure your environment.

3. **Launch VS Code with the devcontainer**:
   ```bash
   cdevcontainer code
   ```
   Accept the prompt to reopen in container when VS Code launches.

For more detailed instructions on using the Caylent Devcontainer, run `cdevcontainer --help` or refer to the [Caylent Devcontainer documentation](https://github.com/caylent-solutions/devcontainer).

## Usage

### Directory Structure

For details on the required directory structure, see the [Directory Structure Documentation](docs/DIRECTORY_STRUCTURE.md).

### Running Tests

For details on installing and using the tftest CLI, see the [CLI Usage Documentation](docs/CLI_USAGE.md).

#### Required Environment Variables

- **AWS Authentication**: Tests require AWS credentials to be available
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

- **Test Control** (optional):
  ```bash
  # To disable idempotency testing
  export TERRATEST_IDEMPOTENCY=false
  
  # To control parallelism of tests within fixtures
  export TERRATEST_DISABLE_PARALLEL_TESTS=true  # Disable parallel tests within fixtures
  ```

### Writing Tests

For details on writing tests, see the [Writing Tests Documentation](docs/WRITING_TESTS.md).

### Idempotency Testing

By default, the framework automatically runs idempotency tests for all Terraform examples. This ensures your Terraform code is idempotent (running it multiple times produces the same result).

The idempotency test:
- Runs automatically when using `RunExample` or `RunAllExamples` functions
- Verifies that running `terraform plan` after `terraform apply` shows no changes
- Is enabled by default

To disable idempotency testing (useful when there are known issues with providers):

```bash
TERRATEST_IDEMPOTENCY=false tftest run
```

The test will run if:
- `TERRATEST_IDEMPOTENCY` environment variable doesn't exist
- `TERRATEST_IDEMPOTENCY` is set to any value other than "false"

The test will be skipped if:
- `TERRATEST_IDEMPOTENCY` is set to "false"

### Controlling Parallel Test Execution

You can control parallelism at two levels:

1. **Test Fixtures**: Control whether different test fixtures run in parallel
2. **Tests Within Fixtures**: Control whether tests within a single fixture run in parallel

By default, both levels of parallelism are disabled for maximum stability. You can enable them as needed:

```bash
# Run test fixtures in parallel
tftest run --parallel-fixtures=true

# Run tests within fixtures in parallel
tftest run --parallel-tests=true

# Enable both levels of parallelism
tftest run --parallel-fixtures=true --parallel-tests=true
```

Controlling parallelism is particularly useful when:
- Tests interact with shared state files
- Tests modify the same AWS resources
- You're experiencing intermittent failures due to race conditions

You can also control test-level parallelism using an environment variable:
```bash
# Disable parallel tests within fixtures
export TERRATEST_DISABLE_PARALLEL_TESTS=true
```

## Available Assertions

### Basic Assertions
- `AssertIdempotent`: Verifies that a Terraform plan shows no changes after apply
- `AssertOutputEquals`: Checks if a specified Terraform output matches an expected value
- `AssertOutputContains`: Checks if a specified Terraform output contains an expected substring
- `AssertOutputMatches`: Checks if a specified Terraform output matches a regular expression
- `AssertOutputNotEmpty`: Checks if a specified Terraform output is not empty
- `AssertOutputEmpty`: Checks if a specified Terraform output is empty

### File Assertions
- `AssertFileExists`: Checks if a file exists at the path specified by the `output_file_path` Terraform output
- `AssertFileContent`: Checks if the `output_content` Terraform output matches the expected value

### Collection Assertions
- `AssertOutputMapContainsKey`: Checks if a Terraform map output contains a specific key
- `AssertOutputMapKeyEquals`: Checks if a key in a Terraform map output equals an expected value
- `AssertOutputListContains`: Checks if a Terraform list output contains an expected value
- `AssertOutputListLength`: Checks if a Terraform list output has the expected length

### JSON Assertions
- `AssertOutputJSONContains`: Checks if a JSON string output contains an expected key-value pair

### Resource Assertions
- `AssertResourceExists`: Checks if a specific resource exists in the Terraform state
- `AssertResourceCount`: Checks if the number of resources of a specific type matches the expected count
- `AssertNoResourcesOfType`: Checks that no resources of a specific type exist in the Terraform state

### Environment Assertions
- `AssertTerraformVersion`: Checks if the Terraform version meets the minimum required version

For detailed documentation on all assertions, including usage examples and requirements, see the [Assertions Documentation](docs/ASSERTIONS.md).

## Documentation

- [Directory Structure](docs/DIRECTORY_STRUCTURE.md) - Required directory structure for using the framework
- [CLI Usage](docs/CLI_USAGE.md) - How to use the TFTest CLI
- [Writing Tests](docs/WRITING_TESTS.md) - Guide to writing tests with the framework
- [Standard Tests](docs/STANDARD_TESTS.md) - Standard tests to run on every Terraform module
- [Assertions Documentation](docs/ASSERTIONS.md) - Detailed guide to all built-in assertions
- [Testing Guide](docs/TESTING_GUIDE.md) - Comprehensive guide to writing and organizing tests
- [Error Handling](docs/ERROR_HANDLING.md) - Error handling strategy for the framework
- [Logging](docs/LOGGING.md) - Logging framework for the framework
- [Benchmarking](docs/BENCHMARKING.md) - Benchmarking and performance optimization
- [CI/CD Pipeline](docs/CI_CD_PIPELINE.md) - Overview of the CI/CD pipeline and release process
- [Contributing Guide](docs/CONTRIBUTING.md) - Guidelines for contributing to the project

## Examples

For a complete working example of how to use this framework in a real Terraform module, see the [Example Directory](example/README.md). This example demonstrates best practices for structuring and testing Terraform modules using this framework.

For simple test fixtures used by the framework's own tests, see the `tests/terraform-module-test-fixtures` directory.

## Installation

To install the framework as a dependency in your Terraform module project:

```bash
go get github.com/caylent-solutions/terraform-terratest-framework@v1.0.0

# To install the CLI tool
go install github.com/caylent-solutions/terraform-terratest-framework/cmd/tftest@v1.0.0
```

This adds the framework to your module's `go.mod` file, allowing you to use its testing capabilities in your Go test files.

## Development

### Prerequisites

- Go 1.23 or later - Required to run the Go tests
- Terraform CLI - Required to execute the `terraform` commands used by the framework
- AWS credentials (if testing AWS resources) - Required for deploying and testing AWS resources

### Setup

```bash
make install-tools
make install
```

### Running Tests

```bash
# Run all tests (unit tests and functional tests)
make test

# Run all tests including example tests
make test-all

# List all functional tests with descriptions
make list-functional-tests

# Run a specific functional test
FUNCTIONAL_TEST=TestName make run-specific-functional-test
```

The `list-functional-tests` command will show all available functional tests in the project. This is useful when you need to run a specific test or understand what tests are available.

The `run-specific-functional-test` command allows you to run a single functional test by name. You must set the `FUNCTIONAL_TEST` environment variable to the name of the test you want to run. For example:

```bash
FUNCTIONAL_TEST=TestCliVersion make run-specific-functional-test
```

### Releasing

Releases are typically managed through the GitHub Actions pipeline. However, for manual releases outside of the pipeline, you can use the `make release-manual` command:

```bash
# Create a manual release based on commit messages
make release-manual

# Explicitly specify the version bump type
make release-manual TYPE=major
make release-manual TYPE=minor
make release-manual TYPE=patch
```

Note: This manual release process should only be used in exceptional circumstances when the automated pipeline cannot be used.

The release process automatically determines the version bump type based on conventional commit messages:
- **Major version bump** is triggered by:
  - `BREAKING CHANGE:` prefix
  - Any prefix with an exclamation mark (`!`): `feat!:`, `fix!:`, `refactor!:`, etc.
- **Minor version bump** is triggered by:
  - `feat:` or `feature:` prefix
- **Patch version bump** is triggered by all other prefixes:
  - `fix:` - Bug fixes
  - `docs:` - Documentation changes
  - `style:` - Code style changes
  - `refactor:` - Code refactoring
  - `test:` - Test changes
  - `chore:` - Routine tasks
  - `ci:` - CI/CD changes
  - `build:` - Build system changes
  - `perf:` - Performance improvements

## Contributing

We welcome contributions from both external contributors and Caylent team members. Please see our [Contributing Guide](docs/CONTRIBUTING.md) for details on how to contribute to this project.

## License

This project is licensed under the Apache License, Version 2.0. See the [LICENSE](LICENSE) file for details.