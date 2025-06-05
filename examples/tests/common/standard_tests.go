package common

import (
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/assertions"
	"github.com/caylent-solutions/terraform-terratest-framework/internal/testctx"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestTerraformValidate runs 'terraform validate' on the module
func TestTerraformValidate(t *testing.T, ctx testctx.TestContext) {
	// Run terraform validate
	terraform.Validate(t, ctx.Terraform)
}

// TestTerraformFormat checks if the Terraform code is properly formatted
func TestTerraformFormat(t *testing.T, ctx testctx.TestContext) {
	// Check if terraform code is formatted
	output, err := terraform.RunTerraformCommandE(t, ctx.Terraform, "fmt", "-check", "-recursive")
	assert.Empty(t, output, "Terraform code should be properly formatted")
	assert.NoError(t, err, "Terraform fmt should not fail")
}

// TestNoHardcodedCredentials checks for hardcoded credentials in Terraform code
func TestNoHardcodedCredentials(t *testing.T, ctx testctx.TestContext) {
	// This is a simplified check - in a real implementation, you would use a more robust method
	// like static code analysis or regex patterns to check for common credential patterns
	output, err := terraform.RunTerraformCommandE(t, ctx.Terraform, "show", "-json")
	assert.NoError(t, err, "Terraform show should not fail")

	// Check for common credential patterns
	credentialPatterns := []string{
		"AKIA", // AWS Access Key ID prefix
		"-----BEGIN RSA PRIVATE KEY-----",
		"-----BEGIN PRIVATE KEY-----",
		"-----BEGIN OPENSSH PRIVATE KEY-----",
		"-----BEGIN PGP PRIVATE KEY BLOCK-----",
		"AccessKeyId",
		"SecretAccessKey",
		"aws_access_key_id",
		"aws_secret_access_key",
	}

	for _, pattern := range credentialPatterns {
		assert.NotContains(t, output, pattern, "Terraform code should not contain hardcoded credentials")
	}
}

// TestRequiredOutputs checks that required outputs are defined
func TestRequiredOutputs(t *testing.T, ctx testctx.TestContext, requiredOutputs []string) {
	outputs := terraform.OutputAll(t, ctx.Terraform)

	for _, output := range requiredOutputs {
		_, exists := outputs[output]
		assert.True(t, exists, "Required output '%s' should be defined", output)
	}
}

// TestRequiredTags checks that all resources have required tags
func TestRequiredTags(t *testing.T, ctx testctx.TestContext, requiredTags []string) {
	// This is a simplified check - in a real implementation, you would parse the state
	// and check each resource that supports tags
	output, err := terraform.RunTerraformCommandE(t, ctx.Terraform, "show", "-json")
	assert.NoError(t, err, "Terraform show should not fail")

	// Check for required tags
	for _, tag := range requiredTags {
		assert.Contains(t, output, tag, "Resources should have required tag '%s'", tag)
	}
}

// TestIdempotency checks that the Terraform code is idempotent
func TestIdempotency(t *testing.T, ctx testctx.TestContext) {
	assertions.AssertIdempotent(t, ctx)
}

// RunStandardTests runs all standard tests on a module
func RunStandardTests(t *testing.T, ctx testctx.TestContext, requiredOutputs []string, requiredTags []string) {
	TestTerraformValidate(t, ctx)
	TestTerraformFormat(t, ctx)
	TestNoHardcodedCredentials(t, ctx)
	TestRequiredOutputs(t, ctx, requiredOutputs)
	TestRequiredTags(t, ctx, requiredTags)
	TestIdempotency(t, ctx)
}
