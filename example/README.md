# Example Terraform Module with Terratest Framework

This directory contains a complete working example of how to use the [Terraform Terratest Framework](https://github.com/caylent-solutions/terraform-terratest-framework) to test Terraform modules. It serves as a reference implementation and demonstration of best practices for testing Terraform modules.

## Purpose of This Example

This example exists within the Terraform Terratest Framework repository to:

1. Demonstrate how to structure a Terraform module with comprehensive tests
2. Show real-world usage of the framework's features
3. Provide a template that can be used as a starting point for new modules
4. Illustrate best practices for Terraform module testing

Unlike the examples in the `examples` directory, which are minimal test fixtures used by the framework's own tests, this is a complete, production-ready example of how to use the framework in your own projects.

## Module Structure

```
terraform-module/
├── examples/                # Example implementations of the module
│   ├── basic/              # Basic example with minimal configuration
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── terraform.tfvars
│   └── advanced/           # Advanced example with more complex configuration
│       ├── main.tf
│       ├── variables.tf
│       └── terraform.tfvars
├── tests/                   # Tests for the module
│   ├── common/              # Tests that run on all examples
│   │   ├── module_test.go
│   │   ├── input_validation_test.go
│   │   └── README.md
│   ├── helpers/             # Helper functions for tests
│   │   ├── helpers.go
│   │   └── README.md
│   ├── basic/               # Tests for the basic example
│   │   ├── module_test.go
│   │   └── README.md
│   └── advanced/            # Tests for the advanced example
│       ├── module_test.go
│       └── README.md
├── main.tf                  # Main module code
├── variables.tf             # Input variables
├── outputs.tf               # Output values
├── versions.tf              # Required providers and versions
└── Makefile                 # Automation for common tasks
```

## Writing Tests

### Using the TestCtx Package

The `testctx` package is the core of the Terraform Terratest Framework, providing the essential functionality for running and managing Terraform tests. It offers several key functions:

- **RunSingleExample**: Runs a specific example with the given configuration
  ```go
  ctx := testctx.RunSingleExample(t, "../../examples", "basic", testctx.TestConfig{
      Name: "basic-test",
  })
  ```

- **RunAllExamples**: Runs all examples in parallel
  ```go
  results := testctx.RunAllExamples(t, "../../examples", configs)
  ```

- **DiscoverAndRunAllTests**: Automatically discovers and runs all examples
  ```go
  testctx.DiscoverAndRunAllTests(t, "../../", func(t *testing.T, ctx testctx.TestContext) {
      // Common assertions for all examples
  })
  ```

For more detailed documentation on the testctx package, see the [TestCtx Package Documentation](https://github.com/caylent-solutions/terraform-terratest-framework/blob/v0.9.0/docs/TESTCTX_PACKAGE.md).

### Using Assertions

The framework provides a set of assertions in the `pkg/assertions` package that you can use to verify your Terraform module's behavior:

```go
import (
    "testing"

    "github.com/caylent-solutions/terraform-terratest-framework/pkg/assertions"
    "github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
)

func TestExample(t *testing.T) {
    ctx := testctx.RunSingleExample(t, "../../examples", "example", testctx.TestConfig{
        Name: "example-test",
    })
    
    // Use assertions
    assertions.AssertOutputEquals(t, ctx, "instance_type", "t2.micro")
    assertions.AssertOutputContains(t, ctx, "bucket_name", "my-bucket")
    assertions.AssertOutputMapContainsKey(t, ctx, "tags", "Environment")
}
```

For a complete list of available assertions, see the [Assertions Documentation](https://github.com/caylent-solutions/terraform-terratest-framework/blob/v0.5.0/docs/ASSERTIONS.md).

### Common Tests

Common tests run on all examples and verify basic functionality:

- **Idempotency**: Automatically included by the framework - ensures that applying the same Terraform code multiple times produces the same result
- **Validation**: Verifies that the Terraform code is syntactically valid
- **Formatting**: Checks if the Terraform code is properly formatted
- **Required Outputs**: Ensures that all required outputs are defined
- **Input Validation**: Verifies that inputs from terraform.tfvars match the provisioned outputs

See the [Common Tests README](tests/common/README.md) for more details on these tests and how to control idempotency testing.

### Example-Specific Tests

Each example has its own tests that verify specific functionality:

- **Basic Example**: Tests the basic functionality of the module
- **Advanced Example**: Tests more complex configurations and features

See the [Basic Tests README](tests/basic/README.md) and [Advanced Tests README](tests/advanced/README.md) for more details.

## Variable Management

Each example uses a consistent approach to variable management:

1. **variables.tf**: Defines variables with default values
2. **terraform.tfvars**: Sets actual values for the example
3. **main.tf**: References variables instead of hardcoded values

This approach:
- Makes examples more maintainable
- Allows for dynamic testing
- Follows Terraform best practices

## Running Tests

Tests can be run using the provided Makefile commands:

```bash
# Run all tests (with both parallel flags set to false for maximum stability)
make test

# Run only common tests
make test-common

# Run a specific test
go test ./tests/common -run '^TestInputsMatchProvisioned

## Developer Workflow

1. **Clone the skeleton**: Start by cloning this skeleton to a new directory
2. **Modify the module**: Update the module code to implement your functionality
3. **Update examples**: Modify the examples to demonstrate your module's usage
4. **Write tests**: Update the tests to verify your module's functionality
5. **Run tests**: Run the tests using `make test`
6. **Commit and push**: Commit your changes and push to your repository

## References

- [Terraform Terratest Framework](https://github.com/caylent-solutions/terraform-terratest-framework/blob/v0.5.0/README.md)
- [TestCtx Package Documentation](https://github.com/caylent-solutions/terraform-terratest-framework/blob/v0.5.0/docs/TESTCTX_PACKAGE.md)
- [Assertions Documentation](https://github.com/caylent-solutions/terraform-terratest-framework/blob/v0.5.0/docs/ASSERTIONS.md)
- [Directory Structure Documentation](https://github.com/caylent-solutions/terraform-terratest-framework/blob/v0.5.0/docs/DIRECTORY_STRUCTURE.md)
- [Writing Tests Documentation](https://github.com/caylent-solutions/terraform-terratest-framework/blob/v0.5.0/docs/WRITING_TESTS.md)
- [CLI Usage Documentation](https://github.com/caylent-solutions/terraform-terratest-framework/blob/v0.5.0/docs/CLI_USAGE.md)

# Format all test files
make format

# Clean up temporary files
make clean
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

For more information on the `tftest` CLI tool, see the [CLI Usage Documentation](https://github.com/caylent-solutions/terraform-terratest-framework/blob/main/docs/CLI_USAGE.md).

## Developer Workflow

1. **Clone the skeleton**: Start by cloning this skeleton to a new directory
2. **Modify the module**: Update the module code to implement your functionality
3. **Update examples**: Modify the examples to demonstrate your module's usage
4. **Write tests**: Update the tests to verify your module's functionality
5. **Run tests**: Run the tests using `make test`
6. **Commit and push**: Commit your changes and push to your repository

## References

- [Terraform Terratest Framework](https://github.com/caylent-solutions/terraform-terratest-framework/blob/v0.5.0/README.md)
- [TestCtx Package Documentation](https://github.com/caylent-solutions/terraform-terratest-framework/blob/v0.5.0/docs/TESTCTX_PACKAGE.md)
- [Assertions Documentation](https://github.com/caylent-solutions/terraform-terratest-framework/blob/v0.5.0/docs/ASSERTIONS.md)
- [Directory Structure Documentation](https://github.com/caylent-solutions/terraform-terratest-framework/blob/v0.5.0/docs/DIRECTORY_STRUCTURE.md)
- [Writing Tests Documentation](https://github.com/caylent-solutions/terraform-terratest-framework/blob/v0.5.0/docs/WRITING_TESTS.md)
- [CLI Usage Documentation](https://github.com/caylent-solutions/terraform-terratest-framework/blob/v0.5.0/docs/CLI_USAGE.md)