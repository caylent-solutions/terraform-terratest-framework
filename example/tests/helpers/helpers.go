package helpers

import (
	"os"
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// VerifyFilePermissions checks if a file has the expected permissions
func VerifyFilePermissions(t *testing.T, filePath string, expectedPerm os.FileMode) {
	// Get the file info
	fileInfo, err := os.Stat(filePath)
	assert.NoError(t, err, "Should be able to get file info")

	// Check the file permissions
	actualPerm := fileInfo.Mode().Perm()
	assert.Equal(t, expectedPerm, actualPerm, "File should have the expected permissions")
}

// VerifyFileContent checks if a file has the expected content
func VerifyFileContent(t *testing.T, filePath string, expectedContent string) {
	// Read the file content
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err, "Should be able to read file content")

	// Check the file content
	assert.Equal(t, expectedContent, string(content), "File should have the expected content")
}

// TestContext represents the context for a Terraform test
// This is a wrapper around the framework's TestContext for backward compatibility
type TestContext struct {
	Terraform *terraform.Options
}

// AssertInputMatchesOutput verifies that a specific input variable matches the corresponding output
func AssertInputMatchesOutput(t *testing.T, ctx TestContext, inputName string, outputName string) {
	// Convert our TestContext to the framework's TestContext to use GetVariableAsMap
	frameworkCtx := testctx.TestContext{
		Terraform: ctx.Terraform,
	}

	// Get the input variable value using the framework's GetVariableAsMap method
	inputValue := frameworkCtx.GetVariableAsMap()[inputName]

	// Get the output value
	outputValue := terraform.Output(t, ctx.Terraform, outputName)

	// Verify that the input matches the output
	assert.Equal(t, inputValue, outputValue, "Input '%s' should match output '%s'", inputName, outputName)
}

// AssertAllInputsMatchOutputs verifies that all input variables match their corresponding outputs
// This assumes that for each input variable, there is an output with the same name
func AssertAllInputsMatchOutputs(t *testing.T, ctx TestContext, inputOutputMap map[string]string) {
	// Convert our TestContext to the framework's TestContext to use GetVariableAsMap
	frameworkCtx := testctx.TestContext{
		Terraform: ctx.Terraform,
	}

	// Get all input variables using the framework's GetVariableAsMap method
	inputs := frameworkCtx.GetVariableAsMap()

	// Get all outputs
	outputs := terraform.OutputAll(t, ctx.Terraform)

	// Verify each input-output pair
	for inputName, outputName := range inputOutputMap {
		inputValue, inputExists := inputs[inputName]
		assert.True(t, inputExists, "Input '%s' should exist", inputName)

		outputValue, outputExists := outputs[outputName]
		assert.True(t, outputExists, "Output '%s' should exist", outputName)

		assert.Equal(t, inputValue, outputValue, "Input '%s' should match output '%s'", inputName, outputName)
	}
}
