# Test Helpers

This directory contains helper functions that can be used across all tests.

## Helper Functions

### `helpers.go`

This file contains all helper functions for tests:

#### File Validation

- **VerifyFilePermissions**: Checks if a file has the expected permissions
- **VerifyFileContent**: Checks if a file has the expected content

#### Input Validation

- **AssertInputMatchesOutput**: Verifies that a specific input variable matches the corresponding output
- **AssertAllInputsMatchOutputs**: Verifies that all input variables match their corresponding outputs based on a provided mapping

## Usage

To use these helpers in your tests:

```go
import (
    "testing"
    
    "github.com/your-org/terraform-modules/skeletons/generic-skeleton/tests/helpers"
    "github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
    "github.com/gruntwork-io/terratest/modules/terraform"
)

func TestExample(t *testing.T) {
    // Run your test
    ctx := testctx.RunSingleExample(t, "../../examples", "example", testctx.TestConfig{
        Name: "example-test",
    })
    
    // Use the helpers
    helpers.VerifyFilePermissions(t, "/path/to/file", 0644)
    helpers.VerifyFileContent(t, "/path/to/file", "expected content")
    
    // Create a TestContext wrapper
    testCtx := helpers.TestContext{
        Terraform: ctx.Terraform,
    }
    
    // Validate inputs match outputs
    inputOutputMap := map[string]string{
        "input_var1": "output_var1",
        "input_var2": "output_var2",
    }
    helpers.AssertAllInputsMatchOutputs(t, testCtx, inputOutputMap)
}
```

## Note on Input Validation

For more robust input validation, consider using the approach in `input_validation_test.go` which reads inputs directly from terraform.tfvars files and compares them with outputs. This approach is more reliable than trying to access input variables through the Terraform context.