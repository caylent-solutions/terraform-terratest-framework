package unit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
)

func TestDiscoverExamples(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "discover-examples-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create some example directories
	exampleDirs := []string{
		"example-basic",
		"example-advanced",
		"not-an-example",
		"another-dir",
	}

	for _, dir := range exampleDirs {
		err := os.Mkdir(filepath.Join(tempDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create a file to ensure it's ignored
	file, err := os.Create(filepath.Join(tempDir, "some-file.txt"))
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	file.Close()

	// Test the function
	examples := testctx.DiscoverExamples(t, tempDir)

	// Verify results
	assert.Equal(t, 2, len(examples))
	assert.Contains(t, examples, "example-basic")
	assert.Contains(t, examples, "example-advanced")
	assert.NotContains(t, examples, "not-an-example")
	assert.NotContains(t, examples, "another-dir")
	assert.NotContains(t, examples, "some-file.txt")
}
