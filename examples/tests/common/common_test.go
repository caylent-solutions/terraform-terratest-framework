package common

import (
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
	"github.com/caylent-solutions/terraform-terratest-framework/tests/unit"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// ==================== Main Test Functions ====================

// TestAllExamples runs all examples and their tests in parallel
// This is the main entry point for running all tests
func TestAllExamples(t *testing.T) {
	// This will discover and run all examples and their tests in parallel
	testctx.DiscoverAndRunAllTests(t, "../..", func(t *testing.T, ctx testctx.TestContext) {
		// Run standard tests on each example
		RunStandardTests(t, ctx, []string{}, []string{})
	})
}

// Example of a common test that would run on all examples
func TestCommonS3BucketEncryption(t *testing.T) {
	// This test would be run for each example
	// The actual implementation would be provided by the user
	// based on their specific requirements
}

// Example of another common test
func TestCommonTagCompliance(t *testing.T) {
	// This test would verify that all resources have required tags
	// The actual implementation would be provided by the user
}

// ==================== Assertions Example ====================

// Example of using various assertions in a test
func TestAssertionsExample(t *testing.T) {
	// Skip this test in normal runs as it's just an example
	t.Skip("This is an example test that demonstrates how to use assertions")

	// Run an example
	ctx := testctx.RunSingleExample(t, "../../", "example", testctx.TestConfig{
		Name: "assertions-example",
		ExtraVars: map[string]interface{}{
			"region": "us-west-2",
		},
	})

	// Basic assertions
	unit.AssertOutputEquals(t, ctx, "region", "us-west-2")
	unit.AssertOutputContains(t, ctx, "instance_id", "i-")
	unit.AssertOutputMatches(t, ctx, "instance_id", "^i-[a-f0-9]{17}$")
	unit.AssertOutputNotEmpty(t, ctx, "vpc_id")
	unit.AssertOutputEmpty(t, ctx, "error_message")

	// File assertions
	filePath := unit.AssertFileExists(t, ctx)
	t.Logf("File exists at: %s", filePath)
	unit.AssertFileContent(t, ctx)

	// Collection assertions
	unit.AssertOutputMapContainsKey(t, ctx, "tags", "Name")
	unit.AssertOutputMapKeyEquals(t, ctx, "tags", "Environment", "test")
	unit.AssertOutputListContains(t, ctx, "subnet_ids", "subnet-12345")
	unit.AssertOutputListLength(t, ctx, "subnet_ids", 3)

	// JSON assertions
	unit.AssertOutputJSONContains(t, ctx, "config_json", "region", "us-west-2")

	// Resource assertions
	unit.AssertResourceExists(t, ctx, "aws_vpc", "main")
	unit.AssertResourceCount(t, ctx, "aws_subnet", 3)
	unit.AssertNoResourcesOfType(t, ctx, "aws_db_instance")

	// Environment assertions
	unit.AssertTerraformVersion(t, ctx, "1.0.0")

	// Idempotency assertion
	unit.AssertIdempotent(t, ctx)
}

// ==================== Standard Tests Implementation ====================

// ValidateTerraform runs 'terraform validate' on the module
func ValidateTerraform(t *testing.T, ctx testctx.TestContext) {
	// Run terraform validate
	terraform.Validate(t, ctx.Terraform)
}

// CheckTerraformFormat checks if the Terraform code is properly formatted
func CheckTerraformFormat(t *testing.T, ctx testctx.TestContext) {
	// Check if terraform code is formatted
	output, err := terraform.RunTerraformCommandE(t, ctx.Terraform, "fmt", "-check", "-recursive")
	assert.Empty(t, output, "Terraform code should be properly formatted")
	assert.NoError(t, err, "Terraform fmt should not fail")
}

// CheckNoHardcodedCredentials checks for hardcoded credentials in Terraform code
func CheckNoHardcodedCredentials(t *testing.T, ctx testctx.TestContext) {
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

// VerifyRequiredOutputs checks that required outputs are defined
func VerifyRequiredOutputs(t *testing.T, ctx testctx.TestContext, requiredOutputs []string) {
	outputs := terraform.OutputAll(t, ctx.Terraform)

	for _, output := range requiredOutputs {
		_, exists := outputs[output]
		assert.True(t, exists, "Required output '%s' should be defined", output)
	}
}

// VerifyRequiredTags checks that all resources have required tags
func VerifyRequiredTags(t *testing.T, ctx testctx.TestContext, requiredTags []string) {
	// This is a simplified check - in a real implementation, you would parse the state
	// and check each resource that supports tags
	output, err := terraform.RunTerraformCommandE(t, ctx.Terraform, "show", "-json")
	assert.NoError(t, err, "Terraform show should not fail")

	// Check for required tags
	for _, tag := range requiredTags {
		assert.Contains(t, output, tag, "Resources should have required tag '%s'", tag)
	}
}

// VerifyIdempotency checks that the Terraform code is idempotent
func VerifyIdempotency(t *testing.T, ctx testctx.TestContext) {
	unit.AssertIdempotent(t, ctx)
}

// RunStandardTests runs all standard tests on a module
func RunStandardTests(t *testing.T, ctx testctx.TestContext, requiredOutputs []string, requiredTags []string) {
	ValidateTerraform(t, ctx)
	CheckTerraformFormat(t, ctx)
	CheckNoHardcodedCredentials(t, ctx)
	VerifyRequiredOutputs(t, ctx, requiredOutputs)
	VerifyRequiredTags(t, ctx, requiredTags)
	VerifyIdempotency(t, ctx)
}

// ==================== Test Helpers and Mocks ====================

// Setup mocks for testing
func setupMocks(t *testing.T) (*testing.T, testctx.TestContext, *terraform.Options) {
	mockT := new(testing.T)
	mockTerraformOptions := &terraform.Options{}

	// Create a testctx.TestContext directly
	ctx := testctx.TestContext{
		Terraform:     mockTerraformOptions,
		TerraformVars: map[string]interface{}{},
		Name:          "test-example",
		ExamplePath:   "test-example",
		Config:        testctx.TestConfig{Name: "test-example"},
	}

	return mockT, ctx, mockTerraformOptions
}
