# Example-Specific Tests

This guide explains how to write tests for specific examples in your Terraform module.

## Overview

The Terraform Terratest Framework supports two testing approaches:

1. **Centralized Testing**: Running all examples from a single test file
2. **Example-Specific Testing**: Dedicated test files for each example

This document focuses on the second approach, which is useful when:
- Each example requires unique test logic
- Examples are complex and warrant dedicated test files
- You want to organize tests by example

## Directory Structure

For example-specific testing, your module should follow this structure:

```
terraform-module/
├── examples/
│   ├── example1/
│   │   ├── main.tf
│   │   └── ...
│   ├── example2/
│   │   ├── main.tf
│   │   └── ...
│   └── ...
├── tests/
│   ├── common/
│   │   └── helpers.go
│   ├── example1/
│   │   └── module_test.go
│   ├── example2/
│   │   └── module_test.go
│   └── ...
└── ...
```

## Writing Example-Specific Tests

In each example's test file (e.g., `tests/example1/module_test.go`):

```go
package example1

import (
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/assertions"
	"github.com/caylent-solutions/terraform-terratest-framework/internal/testctx"
)

func TestExample(t *testing.T) {
	// Run just this specific example
	ctx := testctx.RunSingleExample(t, "../../", "example1", testctx.TestConfig{
		Name: "example1",
		ExtraVars: map[string]interface{}{
			"region": "us-west-2",
		},
	})
	
	// Run example-specific assertions
	assertions.AssertOutputEquals(t, ctx, "instance_type", "t2.micro")
	
	// Custom verification for this example
	verifyResources(t, ctx)
}

func verifyResources(t *testing.T, ctx testctx.TestContext) {
	// Example-specific verification logic
	// ...
}
```

## Sharing Common Test Logic

You can create shared test helpers in the `tests/common` directory:

```go
// tests/common/helpers.go
package common

import (
	"testing"
	"github.com/caylent-solutions/terraform-terratest-framework/internal/testctx"
)

func VerifyS3Bucket(t *testing.T, ctx testctx.TestContext) {
	// Common S3 bucket verification logic
	// ...
}
```

Then import and use these helpers in your example-specific tests:

```go
package example1

import (
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/testctx"
	"terraform-module/tests/common"
)

func TestExample(t *testing.T) {
	ctx := testctx.RunSingleExample(t, "../../", "example1", testctx.TestConfig{})
	
	// Use common test helpers
	common.VerifyS3Bucket(t, ctx)
	
	// Example-specific tests
	// ...
}
```

## Best Practices

1. **Use Relative Paths**: Always use relative paths (`"../../"`) to reference the module root

2. **Organize by Resource Type**: Group test functions by the resources they verify

3. **Share Common Logic**: Put reusable test functions in the `tests/common` package

4. **Use Descriptive Test Names**: Name tests based on what they're verifying

5. **Test One Example Per File**: Keep each test file focused on a single example

## Example: Complete Test File

```go
package ecs_public

import (
	"testing"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	
	"github.com/caylent-solutions/terraform-terratest-framework/internal/assertions"
	"github.com/caylent-solutions/terraform-terratest-framework/internal/testctx"
	"terraform-module/tests/common"
)

func TestEcsPublic(t *testing.T) {
	// Run the specific example
	ctx := testctx.RunSingleExample(t, "../../", "ecs-public", testctx.TestConfig{
		Name: "ecs-public",
		ExtraVars: map[string]interface{}{
			"cluster_name": "test-cluster",
			"region":       "us-west-2",
		},
	})
	
	// Verify outputs
	assertions.AssertOutputEquals(t, ctx, "cluster_name", "test-cluster")
	
	// Verify ECS cluster
	verifyEcsCluster(t, ctx)
	
	// Use common test helpers
	common.VerifySecurityGroups(t, ctx)
}

func verifyEcsCluster(t *testing.T, ctx testctx.TestContext) {
	clusterName := terraform.Output(t, ctx.Terraform, "cluster_name")
	
	// AWS SDK verification logic
	// ...
}
```