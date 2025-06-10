# Basic Example Functional Tests

This directory contains tests specific to the basic example of the module.

## Test Files

### `module_test.go`

This file contains all tests specific to the basic example:

1. **TestBasicOutput**: Tests the basic functionality of the module by verifying that the file is created with the correct content
2. **TestBasicFilePermissions**: Tests that the file created by the basic example has the expected permissions (0644)
3. **TestBasicContentFormat**: Verifies that the basic example creates a plain text file (not JSON)
4. **AssertFilePermissions**: A custom assertion that checks if the file has the expected permissions

The custom assertion demonstrates how to extend the testing framework with your own assertions for specific requirements.

## Example Configuration

The basic example uses:
- `variables.tf` to define variables with default values
- `terraform.tfvars` to set actual values for the example
- `main.tf` to reference variables instead of hardcoded values

This approach makes the example more maintainable and allows for dynamic testing.

## Running the Tests

To run these tests:

```bash
# Run tests for the basic example
tftest run --example-path basic

# Or using Go test directly
cd /path/to/module
go test ./tests/basic
```

For more information on the `tftest` CLI tool, see the [CLI Usage Documentation](https://github.com/caylent-solutions/terraform-terratest-framework/blob/v0.4.2/docs/CLI_USAGE.md).