# Assertions Documentation

This document provides detailed information about the assertions available in the Terraform Test Framework.

## Overview

Assertions are functions that verify specific conditions in your Terraform code. They help you ensure that your Terraform modules behave as expected.

## Available Assertions

### Basic Assertions

#### AssertIdempotent

Verifies that a Terraform plan shows no changes after apply, indicating that the Terraform code is idempotent.

```go
func AssertIdempotent(t *testing.T, ctx testctx.TestContext)
```

**Usage:**
```go
// Verify that the Terraform code is idempotent
assertions.AssertIdempotent(t, ctx)
```

**When to use:** Use this assertion to ensure that running Terraform multiple times produces the same result, which is a key property of well-designed infrastructure as code.

#### AssertOutputEquals

Checks if a specified Terraform output matches an expected value.

```go
func AssertOutputEquals(t *testing.T, ctx testctx.TestContext, outputName string, expectedValue string)
```

**Usage:**
```go
// Verify that the instance_type output is t2.micro
assertions.AssertOutputEquals(t, ctx, "instance_type", "t2.micro")
```

**When to use:** Use this assertion to verify that specific outputs have the expected values.

#### AssertOutputContains

Checks if a specified Terraform output contains an expected substring.

```go
func AssertOutputContains(t *testing.T, ctx testctx.TestContext, outputName string, expectedSubstring string)
```

**Usage:**
```go
// Verify that the instance_id output contains "i-"
assertions.AssertOutputContains(t, ctx, "instance_id", "i-")
```

**When to use:** Use this assertion when you need to verify that an output contains a specific substring, but you don't need to match the entire string.

#### AssertOutputMatches

Checks if a specified Terraform output matches a regular expression.

```go
func AssertOutputMatches(t *testing.T, ctx testctx.TestContext, outputName string, pattern string)
```

**Usage:**
```go
// Verify that the instance_id output matches the pattern for AWS instance IDs
assertions.AssertOutputMatches(t, ctx, "instance_id", "^i-[a-f0-9]{17}$")
```

**When to use:** Use this assertion when you need to verify that an output matches a specific pattern, such as an ID format.

#### AssertOutputNotEmpty

Checks if a specified Terraform output is not empty.

```go
func AssertOutputNotEmpty(t *testing.T, ctx testctx.TestContext, outputName string)
```

**Usage:**
```go
// Verify that the instance_id output is not empty
assertions.AssertOutputNotEmpty(t, ctx, "instance_id")
```

**When to use:** Use this assertion when you need to verify that an output exists and has a non-empty value.

#### AssertOutputEmpty

Checks if a specified Terraform output is empty.

```go
func AssertOutputEmpty(t *testing.T, ctx testctx.TestContext, outputName string)
```

**Usage:**
```go
// Verify that the error_message output is empty
assertions.AssertOutputEmpty(t, ctx, "error_message")
```

**When to use:** Use this assertion when you need to verify that an output exists but has an empty value.

### File Assertions

#### AssertFileExists

Checks if a file exists at the path specified by the `output_file_path` Terraform output.

```go
func AssertFileExists(t *testing.T, ctx testctx.TestContext) string
```

**Usage:**
```go
// Verify that the file specified by output_file_path exists
filePath := assertions.AssertFileExists(t, ctx)
```

**When to use:** Use this assertion when your Terraform module creates a file and outputs its path.

#### AssertFileContent

Checks if the `output_content` Terraform output matches the content of the file specified by the `output_file_path` output.

```go
func AssertFileContent(t *testing.T, ctx testctx.TestContext)
```

**Usage:**
```go
// Verify that the file content matches the output_content
assertions.AssertFileContent(t, ctx)
```

**When to use:** Use this assertion when your Terraform module creates a file with specific content and outputs both the file path and the expected content.

### Collection Assertions

#### AssertOutputMapContainsKey

Checks if a Terraform map output contains a specific key.

```go
func AssertOutputMapContainsKey(t *testing.T, ctx testctx.TestContext, outputName string, key string)
```

**Usage:**
```go
// Verify that the tags output contains the Name key
assertions.AssertOutputMapContainsKey(t, ctx, "tags", "Name")
```

**When to use:** Use this assertion when you need to verify that a map output contains a specific key.

#### AssertOutputMapKeyEquals

Checks if a key in a Terraform map output equals an expected value.

```go
func AssertOutputMapKeyEquals(t *testing.T, ctx testctx.TestContext, outputName string, key string, expectedValue string)
```

**Usage:**
```go
// Verify that the Name tag is set to "example"
assertions.AssertOutputMapKeyEquals(t, ctx, "tags", "Name", "example")
```

**When to use:** Use this assertion when you need to verify that a specific key in a map output has the expected value.

#### AssertOutputListContains

Checks if a Terraform list output contains an expected value.

```go
func AssertOutputListContains(t *testing.T, ctx testctx.TestContext, outputName string, expectedValue string)
```

**Usage:**
```go
// Verify that the subnet_ids output contains a specific subnet ID
assertions.AssertOutputListContains(t, ctx, "subnet_ids", "subnet-12345")
```

**When to use:** Use this assertion when you need to verify that a list output contains a specific value.

#### AssertOutputListLength

Checks if a Terraform list output has the expected length.

```go
func AssertOutputListLength(t *testing.T, ctx testctx.TestContext, outputName string, expectedLength int)
```

**Usage:**
```go
// Verify that the subnet_ids output contains exactly 3 subnet IDs
assertions.AssertOutputListLength(t, ctx, "subnet_ids", 3)
```

**When to use:** Use this assertion when you need to verify that a list output has a specific number of elements.

### JSON Assertions

#### AssertOutputJSONContains

Checks if a JSON string output contains an expected key-value pair.

```go
func AssertOutputJSONContains(t *testing.T, ctx testctx.TestContext, outputName string, key string, expectedValue interface{})
```

**Usage:**
```go
// Verify that the config_json output contains the correct region
assertions.AssertOutputJSONContains(t, ctx, "config_json", "region", "us-west-2")
```

**When to use:** Use this assertion when your Terraform module outputs a JSON string and you need to verify specific values within it.

### Resource Assertions

#### AssertResourceExists

Checks if a specific resource exists in the Terraform state.

```go
func AssertResourceExists(t *testing.T, ctx testctx.TestContext, resourceType string, resourceName string)
```

**Usage:**
```go
// Verify that the aws_s3_bucket.example resource exists
assertions.AssertResourceExists(t, ctx, "aws_s3_bucket", "example")
```

**When to use:** Use this assertion when you need to verify that a specific resource was created.

#### AssertResourceCount

Checks if the number of resources of a specific type matches the expected count.

```go
func AssertResourceCount(t *testing.T, ctx testctx.TestContext, resourceType string, expectedCount int)
```

**Usage:**
```go
// Verify that exactly 2 S3 buckets were created
assertions.AssertResourceCount(t, ctx, "aws_s3_bucket", 2)
```

**When to use:** Use this assertion when you need to verify that a specific number of resources of a certain type were created.

#### AssertNoResourcesOfType

Checks that no resources of a specific type exist in the Terraform state.

```go
func AssertNoResourcesOfType(t *testing.T, ctx testctx.TestContext, resourceType string)
```

**Usage:**
```go
// Verify that no S3 buckets were created
assertions.AssertNoResourcesOfType(t, ctx, "aws_s3_bucket")
```

**When to use:** Use this assertion when you need to verify that no resources of a certain type were created.

### Environment Assertions

#### AssertTerraformVersion

Checks if the Terraform version meets the minimum required version.

```go
func AssertTerraformVersion(t *testing.T, ctx testctx.TestContext, minVersion string)
```

**Usage:**
```go
// Verify that Terraform version is at least 1.0.0
assertions.AssertTerraformVersion(t, ctx, "1.0.0")
```

**When to use:** Use this assertion when your module requires a specific minimum version of Terraform.

## Best Practices

1. **Use Descriptive Error Messages**: The assertions provide clear error messages, but you can add additional context in your test functions.

2. **Combine Assertions**: Use multiple assertions to verify different aspects of your Terraform module.

3. **Group Related Assertions**: Group related assertions into helper functions to make your tests more readable.

4. **Test Edge Cases**: Use assertions to verify that your module handles edge cases correctly.

5. **Test Failure Cases**: Use assertions to verify that your module fails gracefully when given invalid inputs.

## Examples

Here's an example of using multiple assertions to verify a Terraform module:

```go
func TestExample(t *testing.T) {
	// Run the example
	ctx := testctx.RunSingleExample(t, "../../", "example1", testctx.TestConfig{
		Name: "example1",
		ExtraVars: map[string]interface{}{
			"region": "us-west-2",
		},
	})
	
	// Verify outputs
	assertions.AssertOutputEquals(t, ctx, "region", "us-west-2")
	assertions.AssertOutputNotEmpty(t, ctx, "instance_id")
	assertions.AssertOutputMatches(t, ctx, "instance_id", "^i-[a-f0-9]{17}$")
	
	// Verify resources
	assertions.AssertResourceExists(t, ctx, "aws_instance", "example")
	assertions.AssertResourceCount(t, ctx, "aws_security_group", 1)
	
	// Verify idempotency
	assertions.AssertIdempotent(t, ctx)
}
```