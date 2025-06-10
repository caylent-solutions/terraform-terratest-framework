# Common Functional Tests

This directory contains common tests that run on all examples in the module. These tests verify the basic functionality and quality of the Terraform module.

## Test Files

### `module_test.go`

This file contains tests that should run on all examples:

1. **TestTerraformValidate**: Verifies that the Terraform code is syntactically valid
2. **TestTerraformFormat**: Checks if the Terraform code is properly formatted
3. **TestRequiredOutputs**: Ensures that all required outputs are defined
4. **TestFileCreation**: Verifies that the module creates the output file as expected
5. **TestAllAssertionTypes**: Demonstrates all the assertion types available in the framework:
   - Basic Assertions: `AssertOutputEquals`, `AssertOutputContains`, `AssertOutputMatches`, `AssertOutputNotEmpty`
   - File Assertions: `AssertFileExists`, `AssertFileContent`
   - Resource Assertions: `AssertResourceExists`, `AssertResourceCount`, `AssertNoResourcesOfType`
   - Environment Assertions: `AssertTerraformVersion`
6. **TestBenchmarking**: Demonstrates how to benchmark Terraform operations

### `input_validation_test.go`

This file contains tests that verify input variables match the provisioned resources:

1. **TestInputsMatchProvisioned**: Verifies that the inputs provided in terraform.tfvars files match what was actually provisioned by Terraform
   - Reads input values dynamically from terraform.tfvars files
   - Compares these inputs with the actual outputs from Terraform
   - Handles both simple string values and complex JSON structures

## Idempotency Testing

Idempotency testing is automatically included when using the Terraform Terratest Framework. The framework runs idempotency tests whenever you use `testctx.RunSingleExample` or `testctx.RunAllExamples` functions.

The idempotency test:
- Verifies that running `terraform plan` after `terraform apply` shows no changes
- Is enabled by default

To disable idempotency testing (useful when there are known issues with providers):

```bash
TERRATEST_IDEMPOTENCY=false tftest run
```

For more details on idempotency testing, see the [Writing Tests Documentation](https://github.com/caylent-solutions/terraform-terratest-framework/blob/v0.4.2/docs/WRITING_TESTS.md#idempotency-testing).

## Running the Tests

To run these tests:

```bash
# Run only common tests
tftest run --common

# Run a specific test
go test ./tests/common -run '^TestInputsMatchProvisioned$'
```

For more information on the `tftest` CLI tool, see the [CLI Usage Documentation](https://github.com/caylent-solutions/terraform-terratest-framework/blob/v0.4.2/docs/CLI_USAGE.md).