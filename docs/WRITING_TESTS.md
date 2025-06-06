# Writing Tests

This document explains how to write tests using the Terraform Test Framework.

## Test Organization

The framework supports a structured approach to testing:

1. **Example-Specific Tests**: Each example has its own test directory with the same name
2. **Common Tests**: Tests in the `common` directory run on all examples
3. **Helper Functions**: Reusable test helpers in the `helpers` directory

## Example-Specific Tests

Each example should have its own test directory with the same name:

```go
// In tests/example1/module_test.go
package example1

import (
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/assertions"
	"github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
	"github.com/gruntwork-io/terratest/modules/terraform"
	
	"terraform-module/tests/helpers"  // Import custom helpers
)

func TestExample(t *testing.T) {
	// Run this specific example
	ctx := testctx.RunSingleExample(t, "../../", "example1", testctx.TestConfig{
		Name: "example1",
		ExtraVars: map[string]interface{}{
			"region": "us-west-2",
		},
	})
	
	// Example-specific assertions
	assertions.AssertOutputEquals(t, ctx, "instance_type", "t2.micro")
	
	// Use helper functions
	helpers.VerifyS3Bucket(t, ctx, "AES256")
	
	// Custom verification
	verifyResources(t, ctx)
}

func verifyResources(t *testing.T, ctx testctx.TestContext) {
	// Example-specific verification logic
	// ...
}
```

## Common Tests

Tests in the `common` directory will run on all examples:

```go
// In tests/common/common_test.go
package common

import (
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
)

// TestAllExamples runs all examples and their tests in parallel
func TestAllExamples(t *testing.T) {
	// This will discover and run all examples and their tests in parallel
	testctx.DiscoverAndRunAllTests(t, "../..")
}

// Common test that runs on all examples
func TestCommonS3BucketEncryption(t *testing.T) {
	// This test would verify that all S3 buckets are encrypted
	// ...
}
```

## Helper Functions

Create reusable helper functions in the `helpers` directory:

```go
// In tests/helpers/aws_helpers.go
package helpers

import (
	"testing"
	"github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
)

func VerifyS3Bucket(t *testing.T, ctx testctx.TestContext, expectedEncryption string) {
	// Helper logic to verify S3 bucket
	// ...
}
```

## Idempotency Testing

By default, the framework automatically runs idempotency tests for all Terraform examples, whether running all examples, a specific example, or common tests. This ensures your Terraform code is idempotent (running it multiple times produces the same result).

The idempotency test:
- Runs automatically when using `RunExample`, `RunSingleExample`, or `RunAllExamples` functions
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

## Best Practices

1. **Organize Tests by Example**: Keep tests for each example in their own directory

2. **Use Common Tests for Shared Requirements**: Put tests that should run on all examples in the `common` directory

3. **Create Helper Functions**: Extract reusable test logic into helper functions in the `helpers` directory

4. **Clean Up Resources**: The framework automatically cleans up Terraform resources, but if your tests create additional resources, clean them up

5. **Use Descriptive Test Names**: Name tests based on what they're verifying

6. **Error Handling**: Include proper error handling and descriptive assertion messages

7. **AWS Authentication**: Ensure AWS credentials are properly configured before running tests