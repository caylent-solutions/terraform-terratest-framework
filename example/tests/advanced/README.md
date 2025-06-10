# Advanced Example Functional Tests

This directory contains tests specific to the advanced example of the module.

## Test Files

### `module_test.go`

This file contains all tests specific to the advanced example:

1. **TestAdvancedOutput**: Tests the advanced functionality of the module by verifying that the JSON file is created with the correct content
2. **TestAdvancedJSONStructure**: Tests that the JSON file created by the advanced example has the expected structure with required keys
3. **TestCollectionAssertions**: Tests various collection and JSON assertions:
   - `AssertOutputMapContainsKey`: Verifies that the JSON data contains specific keys
   - `AssertOutputMapKeyEquals`: Verifies that keys in the JSON data have expected values
   - `AssertOutputListContains`: Verifies that lists in the JSON data contain expected values
   - `AssertOutputListLength`: Verifies that lists in the JSON data have the expected length
   - `AssertOutputJSONContains`: Verifies that the JSON string contains expected key-value pairs
4. **TestAdvancedJSONFormat**: Verifies that the advanced example creates a valid JSON file with specific required fields and structure
5. **AssertJSONStructure**: A custom assertion that checks if the JSON content has the expected structure

The custom assertions demonstrate how to extend the testing framework with your own assertions for specific requirements, particularly for validating structured data like JSON.

## Example Configuration

The advanced example uses:
- `variables.tf` to define variables with default values (using `any` type for complex JSON)
- `terraform.tfvars` to set actual values for the example
- `main.tf` to reference variables instead of hardcoded values

This approach makes the example more maintainable and allows for dynamic testing of complex data structures.

## Running the Tests

To run these tests:

```bash
# Run tests for the advanced example
tftest run --example-path advanced

# Or using Go test directly
cd /path/to/module
go test ./tests/advanced
```

For more information on the `tftest` CLI tool, see the [CLI Usage Documentation](https://github.com/caylent-solutions/terraform-terratest-framework/blob/v0.4.2/docs/CLI_USAGE.md).