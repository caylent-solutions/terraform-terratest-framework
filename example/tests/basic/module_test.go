package basic_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestBasicOutput tests the basic functionality of the module
// It verifies that the file is created with the correct content
func TestBasicOutput(t *testing.T) {
	// Use the values from terraform.tfvars instead of overriding them
	ctx := testctx.RunSingleExample(t, "../../examples", "basic", testctx.TestConfig{
		Name: "basic",
	})

	// Verify file exists and has correct content
	filePath := terraform.Output(t, ctx.Terraform, "output_file_path")
	assert.NotEmpty(t, filePath, "File path should not be empty")

	content := terraform.Output(t, ctx.Terraform, "output_content")
	assert.Equal(t, "hello from basic", content, "File content should match expected value")
}

// AssertFilePermissions is a custom assertion that checks if the file has the expected permissions
func AssertFilePermissions(t *testing.T, ctx testctx.TestContext, expectedPerm os.FileMode) {
	// Get the file path from the output
	filePath := terraform.Output(t, ctx.Terraform, "output_file_path")
	assert.NotEmpty(t, filePath, "File path should not be empty")

	// Use the correct path by prepending the working directory
	fullPath := filepath.Join(ctx.Terraform.TerraformDir, filePath)

	// Get the file info
	fileInfo, err := os.Stat(fullPath)
	assert.NoError(t, err, "Should be able to get file info")

	// Check the file permissions
	actualPerm := fileInfo.Mode().Perm()
	assert.Equal(t, expectedPerm, actualPerm, "File should have the expected permissions")
}

// TestBasicFilePermissions tests that the file created by the basic example has the expected permissions
// This demonstrates a custom assertion for file permissions
func TestBasicFilePermissions(t *testing.T) {
	// Run the example
	ctx := testctx.RunSingleExample(t, "../../examples", "basic", testctx.TestConfig{
		Name: "basic-permissions-test",
	})

	// Use our custom assertion to check file permissions
	// 0644 is a common permission for files (rw-r--r--)
	AssertFilePermissions(t, ctx, 0644)
}

// TestBasicContentFormat verifies that the basic example creates a plain text file
// This is a unique test specific to the basic example
func TestBasicContentFormat(t *testing.T) {
	// Run the example
	ctx := testctx.RunSingleExample(t, "../../examples", "basic", testctx.TestConfig{
		Name: "basic-content-format-test",
	})

	// Get the output content
	content := terraform.Output(t, ctx.Terraform, "output_content")

	// Verify that the content is plain text (not JSON)
	assert.Equal(t, "hello from basic", content, "Basic example should output plain text")

	// Verify that attempting to parse as JSON would fail
	// This is a unique characteristic of the basic example
	assert.NotContains(t, content, "{", "Basic example should not contain JSON syntax")
	assert.NotContains(t, content, "}", "Basic example should not contain JSON syntax")
}
