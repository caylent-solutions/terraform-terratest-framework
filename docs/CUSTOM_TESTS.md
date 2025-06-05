# Writing Custom Tests

This guide explains how to write custom tests for your Terraform modules using the Terraform Terratest Framework.

## Overview

While the framework provides common assertions for basic testing, real-world Terraform modules often require custom tests specific to their functionality. The framework supports running custom test functions on all examples in parallel.

## Writing Custom Test Functions

A custom test function is any function that matches the `CustomTestFunc` type:

```go
type CustomTestFunc func(t *testing.T, ctx testctx.TestContext)
```

Inside this function, you can:
- Access the Terraform outputs
- Make assertions about resources
- Call AWS APIs to verify resources
- Check for specific conditions
- Use any of the built-in assertions

Example of a custom test function:

```go
func verifyS3BucketExists(t *testing.T, ctx testctx.TestContext) {
    // Get the bucket name from Terraform outputs
    bucketName := terraform.Output(t, ctx.Terraform, "bucket_name")
    
    // Use AWS SDK to verify the bucket exists
    awsRegion := "us-west-2"
    s3Client, err := s3.NewFromConfig(aws.Config{
        Region: awsRegion,
    })
    require.NoError(t, err)
    
    _, err = s3Client.HeadBucket(context.Background(), &s3.HeadBucketInput{
        Bucket: aws.String(bucketName),
    })
    assert.NoError(t, err, "S3 bucket should exist")
}
```

## Running Custom Tests on All Examples

There are two main approaches to running custom tests on all examples:

### Approach 1: Run Examples First, Then Custom Tests

```go
func TestModule(t *testing.T) {
    // Run all examples with default configs
    results := testctx.RunAllExamples(t, "../..", nil)
    
    // Define custom test functions
    verifyS3 := func(t *testing.T, ctx testctx.TestContext) {
        // Custom S3 verification logic
    }
    
    verifyIAM := func(t *testing.T, ctx testctx.TestContext) {
        // Custom IAM verification logic
    }
    
    // Run custom tests on all examples
    testctx.RunCustomTests(t, results, verifyS3, verifyIAM)
}
```

### Approach 2: Run Examples and Custom Tests in One Go

```go
func TestModule(t *testing.T) {
    // Define custom test functions
    verifyS3 := func(t *testing.T, ctx testctx.TestContext) {
        // Custom S3 verification logic
    }
    
    verifyIAM := func(t *testing.T, ctx testctx.TestContext) {
        // Custom IAM verification logic
    }
    
    // Run all examples and then run custom tests on each
    testctx.RunAllExamplesWithTests(t, "../..", nil, verifyS3, verifyIAM)
}
```

## Example-Specific Tests

You can run different tests for different examples:

```go
func TestModule(t *testing.T) {
    // Run all examples
    results := testctx.RunAllExamples(t, "../..", nil)
    
    // Define a custom test that behaves differently based on the example
    customTest := func(t *testing.T, ctx testctx.TestContext) {
        switch ctx.Config.Name {
        case "basic":
            // Test logic for basic example
            assertions.AssertOutputEquals(t, ctx, "instance_type", "t2.micro")
        
        case "advanced":
            // Test logic for advanced example
            assertions.AssertOutputEquals(t, ctx, "instance_type", "t2.large")
            
            // Additional advanced-specific tests
            bucketName := terraform.Output(t, ctx.Terraform, "bucket_name")
            assert.Contains(t, bucketName, "advanced")
        }
    }
    
    // Run the custom test on all examples
    testctx.RunCustomTests(t, results, customTest)
}
```

## Best Practices

1. **Organize Tests by Resource Type**: Group tests by the AWS resource they're testing (S3, EC2, IAM, etc.)

2. **Use Helper Functions**: Create reusable helper functions for common verification tasks

3. **Handle Example-Specific Logic**: Use conditionals or switch statements to handle differences between examples

4. **Clean Up Resources**: The framework automatically cleans up Terraform resources, but if your tests create additional resources, clean them up

5. **Use Descriptive Test Names**: When defining subtests, use descriptive names that indicate what's being tested

6. **Error Handling**: Include proper error handling and descriptive assertion messages

## Example: Complete Custom Test

```go
package functional

import (
    "testing"
    "context"
    
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/gruntwork-io/terratest/modules/terraform"
    
    "github.com/caylent-solutions/terraform-terratest-framework/internal/assertions"
    "github.com/caylent-solutions/terraform-terratest-framework/internal/testctx"
)

func TestS3Module(t *testing.T) {
    // Define custom test functions
    verifyBucketExists := func(t *testing.T, ctx testctx.TestContext) {
        bucketName := terraform.Output(t, ctx.Terraform, "bucket_name")
        
        // AWS SDK verification logic
        // ...
    }
    
    verifyBucketPolicy := func(t *testing.T, ctx testctx.TestContext) {
        bucketName := terraform.Output(t, ctx.Terraform, "bucket_name")
        
        // AWS SDK verification logic
        // ...
    }
    
    // Run all examples and custom tests
    testctx.RunAllExamplesWithTests(t, "../..", nil, 
        verifyBucketExists, 
        verifyBucketPolicy,
    )
}
```

## Advanced: Test Fixtures

For complex tests, you might want to create test fixtures:

```go
type S3TestFixture struct {
    BucketName string
    Client     *s3.Client
}

func setupS3Fixture(t *testing.T, ctx testctx.TestContext) *S3TestFixture {
    bucketName := terraform.Output(t, ctx.Terraform, "bucket_name")
    
    client, err := s3.NewFromConfig(aws.Config{
        Region: "us-west-2",
    })
    require.NoError(t, err)
    
    return &S3TestFixture{
        BucketName: bucketName,
        Client:     client,
    }
}

func TestS3ModuleWithFixture(t *testing.T) {
    results := testctx.RunAllExamples(t, "../..", nil)
    
    for name, ctx := range results {
        t.Run(fmt.Sprintf("S3Tests_%s", name), func(t *testing.T) {
            fixture := setupS3Fixture(t, ctx)
            
            t.Run("BucketExists", func(t *testing.T) {
                // Use fixture to test bucket existence
            })
            
            t.Run("BucketPolicy", func(t *testing.T) {
                // Use fixture to test bucket policy
            })
        })
    }
}
```

This approach allows you to organize complex tests with shared setup logic.