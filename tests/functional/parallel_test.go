package functional

import (
	"path/filepath"
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
)

func TestParallelExamples(t *testing.T) {
	// Create a temporary directory for test results
	tempDir := t.TempDir()

	// Configure examples with default settings
	configs := map[string]testctx.TestConfig{
		"example-basic": {
			Name: "example-basic",
			ExtraVars: map[string]interface{}{
				"output_content":  "parallel test",
				"output_filename": filepath.Join(tempDir, "basic-parallel-test.txt"),
			},
		},
		"example-advanced": {
			Name: "example-advanced",
			ExtraVars: map[string]interface{}{
				"output_content":  "parallel test",
				"output_filename": filepath.Join(tempDir, "advanced-parallel-test.txt"),
			},
		},
	}

	// Run all examples in parallel
	_ = testctx.RunAllExamples(t, "../terraform-module-test-fixtures", configs)

	// No need for additional assertions here since we're just testing parallel execution
}
