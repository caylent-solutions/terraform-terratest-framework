# Standard Terraform Module Tests

This document describes standard tests that should be run on every Terraform module.

## Overview

The framework provides a set of standard tests that can be run on any Terraform module. These tests check for common issues and best practices.

## Standard Tests

### 1. Terraform Validate

Runs `terraform validate` to ensure the module is syntactically valid and internally consistent.

```go
func TestTerraformValidate(t *testing.T, ctx testctx.TestContext) {
    terraform.Validate(t, ctx.Terraform)
}
```

### 2. Terraform Format

Checks if the Terraform code is properly formatted according to the HashiCorp style guide.

```go
func TestTerraformFormat(t *testing.T, ctx testctx.TestContext) {
    stdout, stderr, err := terraform.RunTerraformCommandE(t, ctx.Terraform, "fmt", "-check", "-recursive")
    assert.Empty(t, stdout, "Terraform code should be properly formatted")
    assert.Empty(t, stderr, "Terraform fmt should not produce errors")
    assert.NoError(t, err, "Terraform fmt should not fail")
}
```

### 3. No Hardcoded Credentials

Checks for hardcoded credentials in Terraform code.

```go
func TestNoHardcodedCredentials(t *testing.T, ctx testctx.TestContext) {
    // Check for common credential patterns
    // ...
}
```

### 4. Required Outputs

Checks that required outputs are defined.

```go
func TestRequiredOutputs(t *testing.T, ctx testctx.TestContext, requiredOutputs []string) {
    outputs := terraform.OutputAll(t, ctx.Terraform)
    
    for _, output := range requiredOutputs {
        _, exists := outputs[output]
        assert.True(t, exists, "Required output '%s' should be defined", output)
    }
}
```

### 5. Required Tags

Checks that all resources have required tags.

```go
func TestRequiredTags(t *testing.T, ctx testctx.TestContext, requiredTags []string) {
    // Check for required tags
    // ...
}
```

### 6. Idempotency

Checks that the Terraform code is idempotent (running it multiple times produces the same result).

```go
func TestIdempotency(t *testing.T, ctx testctx.TestContext) {
    assertions.AssertIdempotent(t, ctx)
}
```

## Using Standard Tests

You can use the standard tests in your common test file:

```go
// In tests/common/common_test.go
package common

import (
    "testing"

    "github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
    "github.com/caylent-solutions/terraform-terratest-framework/examples/tests/common"
)

func TestStandardTests(t *testing.T) {
    // Run all examples and apply standard tests to each
    testctx.DiscoverAndRunAllTests(t, "../..", func(t *testing.T, ctx testctx.TestContext) {
        // Define required outputs and tags for your module
        requiredOutputs := []string{"id", "arn"}
        requiredTags := []string{"Name", "Environment", "Owner"}
        
        // Run all standard tests
        common.RunStandardTests(t, ctx, requiredOutputs, requiredTags)
    })
}
```

## Customizing Standard Tests

You can customize the standard tests by:

1. Creating your own version of the standard tests
2. Adding additional tests specific to your module
3. Skipping certain tests that don't apply to your module

Example:

```go
func TestCustomStandardTests(t *testing.T) {
    testctx.DiscoverAndRunAllTests(t, "../..", func(t *testing.T, ctx testctx.TestContext) {
        // Run only specific standard tests
        common.TestTerraformValidate(t, ctx)
        common.TestIdempotency(t, ctx)
        
        // Add custom tests
        TestCustomRequirement(t, ctx)
    })
}
```