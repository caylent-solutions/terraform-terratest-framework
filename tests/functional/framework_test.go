package functional

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/caylent-solutions/terraform-test-framework/internal/assertions"
	"github.com/caylent-solutions/terraform-test-framework/internal/testctx"
	"github.com/stretchr/testify/assert"
)

func TestRunSingleExample(t *testing.T) {
	// Get the path to the test fixture
	fixtureDir, err := filepath.Abs("../terraform-module-test-fixtures/example-basic")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Run the example
	ctx := testctx.RunSingleExample(t, fixtureDir, ".", testctx.TestConfig{
		Name: "basic-example",
		ExtraVars: map[string]interface{}{
			"output_content": "Test Content",
		},
	})

	// Verify outputs
	assertions.AssertOutputEquals(t, ctx, "output_content", "Test Content")
	
	// Verify file was created
	filePath := assertions.AssertFileExists(t, ctx)
	assert.NotEmpty(t, filePath, "File path should not be empty")
	
	// Verify file content
	assertions.AssertFileContent(t, ctx)
}

func TestRunMultipleExamples(t *testing.T) {
	// This test verifies that multiple examples can be run concurrently
	
	// Get the path to the test fixtures directory
	fixturesDir, err := filepath.Abs("../terraform-module-test-fixtures")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}
	
	// Create a temporary directory for test results
	tempDir := t.TempDir()
	
	// Run the basic example
	basicCtx := testctx.RunSingleExample(t, fixturesDir, "example-basic", testctx.TestConfig{
		Name: "basic-example",
		ExtraVars: map[string]interface{}{
			"output_content": "Basic Content",
			"output_filename": filepath.Join(tempDir, "basic.txt"),
		},
	})
	
	// Run the advanced example
	advancedCtx := testctx.RunSingleExample(t, fixturesDir, "example-advanced", testctx.TestConfig{
		Name: "advanced-example",
		ExtraVars: map[string]interface{}{
			"output_content": "Advanced Content",
			"output_filename": filepath.Join(tempDir, "advanced.txt"),
		},
	})
	
	// Verify basic example outputs
	assertions.AssertOutputEquals(t, basicCtx, "output_content", "Basic Content")
	basicFilePath := assertions.AssertFileExists(t, basicCtx)
	assert.Contains(t, basicFilePath, "basic.txt", "File path should contain basic.txt")
	
	// Verify advanced example outputs
	assertions.AssertOutputEquals(t, advancedCtx, "output_content", "Advanced Content")
	advancedFilePath := assertions.AssertFileExists(t, advancedCtx)
	assert.Contains(t, advancedFilePath, "advanced.txt", "File path should contain advanced.txt")
	
	// Verify both files exist
	_, err = os.Stat(filepath.Join(tempDir, "basic.txt"))
	assert.NoError(t, err, "basic.txt should exist")
	
	_, err = os.Stat(filepath.Join(tempDir, "advanced.txt"))
	assert.NoError(t, err, "advanced.txt should exist")
}

func TestDiscoverAndRunAllTests(t *testing.T) {
	// Skip this test in CI environments as it requires more setup
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping in CI environment")
	}
	
	// Get the path to the test fixtures directory
	fixturesDir, err := filepath.Abs("../terraform-module-test-fixtures")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}
	
	// Run all examples
	testctx.DiscoverAndRunAllTests(t, fixturesDir, func(t *testing.T, ctx testctx.TestContext) {
		// This function will be called for each example
		t.Logf("Running tests for example: %s", ctx.Name)
		
		// Verify that outputs exist
		output := ctx.GetOutput(t, "output_content")
		assert.NotEmpty(t, output, "output_content should not be empty")
		
		// Verify that file exists
		filePath := assertions.AssertFileExists(t, ctx)
		assert.NotEmpty(t, filePath, "File path should not be empty")
	})
}