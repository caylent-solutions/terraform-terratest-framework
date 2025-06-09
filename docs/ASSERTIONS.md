# Assertions Documentation

This document provides detailed information about the assertions available in the Terraform Terratest Framework.

## Using Assertions

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

## Available Assertions

### Basic Assertions

- **AssertOutputEquals**: Checks if a specified Terraform output matches an expected value
  ```go
  assertions.AssertOutputEquals(t, ctx, "instance_type", "t2.micro")
  ```

- **AssertOutputContains**: Checks if a specified Terraform output contains an expected substring
  ```go
  assertions.AssertOutputContains(t, ctx, "bucket_name", "my-bucket")
  ```

- **AssertOutputMatches**: Checks if a specified Terraform output matches a regular expression
  ```go
  assertions.AssertOutputMatches(t, ctx, "instance_id", "^i-[a-f0-9]{17}$")
  ```

- **AssertOutputNotEmpty**: Checks if a specified Terraform output is not empty
  ```go
  assertions.AssertOutputNotEmpty(t, ctx, "instance_id")
  ```

- **AssertOutputEmpty**: Checks if a specified Terraform output is empty
  ```go
  assertions.AssertOutputEmpty(t, ctx, "error_message")
  ```

### File Assertions

- **AssertFileExists**: Checks if a file exists at the path specified by the `output_file_path` Terraform output
  ```go
  assertions.AssertFileExists(t, ctx)
  ```

- **AssertFileContent**: Checks if the `output_content` Terraform output matches the expected value
  ```go
  assertions.AssertFileContent(t, ctx)
  ```

### Collection Assertions

- **AssertOutputMapContainsKey**: Checks if a Terraform map output contains a specific key
  ```go
  assertions.AssertOutputMapContainsKey(t, ctx, "tags", "Environment")
  ```

- **AssertOutputMapKeyEquals**: Checks if a key in a Terraform map output equals an expected value
  ```go
  assertions.AssertOutputMapKeyEquals(t, ctx, "tags", "Environment", "dev")
  ```

- **AssertOutputListContains**: Checks if a Terraform list output contains an expected value
  ```go
  assertions.AssertOutputListContains(t, ctx, "availability_zones", "us-west-2a")
  ```

- **AssertOutputListLength**: Checks if a Terraform list output has the expected length
  ```go
  assertions.AssertOutputListLength(t, ctx, "subnet_ids", 3)
  ```

### JSON Assertions

- **AssertOutputJSONContains**: Checks if a JSON string output contains an expected key-value pair
  ```go
  assertions.AssertOutputJSONContains(t, ctx, "json_output", "enabled", true)
  ```

### Resource Assertions

- **AssertResourceExists**: Checks if a specific resource exists in the Terraform state
  ```go
  assertions.AssertResourceExists(t, ctx, "aws_instance", "example")
  ```

- **AssertResourceCount**: Checks if the number of resources of a specific type matches the expected count
  ```go
  assertions.AssertResourceCount(t, ctx, "aws_subnet", 3)
  ```

- **AssertNoResourcesOfType**: Checks that no resources of a specific type exist in the Terraform state
  ```go
  assertions.AssertNoResourcesOfType(t, ctx, "aws_db_instance")
  ```

### Environment Assertions

- **AssertTerraformVersion**: Checks if the Terraform version meets the minimum required version
  ```go
  assertions.AssertTerraformVersion(t, ctx, "1.12.1")
  ```

- **AssertIdempotent**: Verifies that a Terraform plan shows no changes after apply
  ```go
  assertions.AssertIdempotent(t, ctx)
  ```

## Creating Custom Assertions

You can create your own custom assertions by building on top of the provided assertions:

```go
// In your tests/helpers/custom_assertions.go file
package helpers

import (
    "testing"
    
    "github.com/caylent-solutions/terraform-terratest-framework/pkg/assertions"
    "github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
)

// AssertS3BucketEncryption checks if an S3 bucket has the expected encryption
func AssertS3BucketEncryption(t testing.TB, ctx testctx.TestContext, expectedEncryption string) {
    // Get bucket name from outputs
    bucketName := terraform.Output(t, ctx.Terraform, "bucket_name")
    
    // Get encryption configuration from outputs
    encryption := terraform.Output(t, ctx.Terraform, "bucket_encryption")
    
    // Use the built-in assertions
    assertions.AssertOutputEquals(t, ctx, "bucket_encryption", expectedEncryption)
}
```

## Best Practices

1. **Use Descriptive Error Messages**: Include descriptive error messages in your assertions to make test failures easier to understand.

2. **Group Related Assertions**: Group related assertions together to make your tests more readable and maintainable.

3. **Create Custom Assertions for Complex Checks**: If you find yourself repeating the same set of assertions, create a custom assertion function.

4. **Test Both Happy and Error Paths**: Verify that your module works correctly and handles errors gracefully.

5. **Use Assertions to Document Requirements**: Assertions serve as documentation for your module's requirements and expected behavior.