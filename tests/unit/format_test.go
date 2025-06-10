package unit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatDirectoryStructure(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "TestFormatDirectoryStructure")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the required directory structure
	err = os.MkdirAll(filepath.Join(tempDir, "examples"), 0755)
	assert.NoError(t, err, "Failed to create examples directory")

	err = os.MkdirAll(filepath.Join(tempDir, "tests"), 0755)
	assert.NoError(t, err, "Failed to create tests directory")

	// Create a dummy example directory
	exampleName := "example-test"
	err = os.MkdirAll(filepath.Join(tempDir, "examples", exampleName), 0755)
	assert.NoError(t, err, "Failed to create example directory")

	// Create the corresponding test directory
	err = os.MkdirAll(filepath.Join(tempDir, "tests", exampleName), 0755)
	assert.NoError(t, err, "Failed to create test directory")

	// Create common and helpers directories
	err = os.MkdirAll(filepath.Join(tempDir, "tests", "common"), 0755)
	assert.NoError(t, err, "Failed to create common directory")

	err = os.MkdirAll(filepath.Join(tempDir, "tests", "helpers"), 0755)
	assert.NoError(t, err, "Failed to create helpers directory")

	// Verify the directory structure exists as expected
	_, err = os.Stat(filepath.Join(tempDir, "examples"))
	assert.NoError(t, err, "Examples directory should exist")

	_, err = os.Stat(filepath.Join(tempDir, "tests"))
	assert.NoError(t, err, "Tests directory should exist")

	_, err = os.Stat(filepath.Join(tempDir, "examples", exampleName))
	assert.NoError(t, err, "Example directory should exist")

	_, err = os.Stat(filepath.Join(tempDir, "tests", exampleName))
	assert.NoError(t, err, "Test directory should exist")

	_, err = os.Stat(filepath.Join(tempDir, "tests", "common"))
	assert.NoError(t, err, "Common directory should exist")

	_, err = os.Stat(filepath.Join(tempDir, "tests", "helpers"))
	assert.NoError(t, err, "Helpers directory should exist")
}
