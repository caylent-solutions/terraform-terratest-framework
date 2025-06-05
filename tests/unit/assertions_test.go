package unit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAssertIdempotent tests the AssertIdempotent function
func TestAssertIdempotent(t *testing.T) {
	// This test requires mocking terraform.Plan which is difficult in a unit test
	// We'll test this in a more comprehensive integration test
	t.Skip("Skipping test that requires mocking terraform.Plan function")
}

// TestAssertFileExists tests the AssertFileExists function
func TestAssertFileExists(t *testing.T) {
	// Create a temporary file for testing
	tempDir, err := os.MkdirTemp("", "test-assertions")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tempFile := filepath.Join(tempDir, "test-file.txt")
	err = os.WriteFile(tempFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a mock test context
	ctx := NewMockTestContextSimple()
	ctx.SetOutput("output_file_path", tempFile)

	// Test the function
	filePath := AssertFileExists(t, ctx)
	assert.Equal(t, tempFile, filePath)
}

// TestAssertFileContent tests the AssertFileContent function
func TestAssertFileContent(t *testing.T) {
	// Create a temporary file for testing
	tempDir, err := os.MkdirTemp("", "test-assertions")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tempFile := filepath.Join(tempDir, "test-file.txt")
	err = os.WriteFile(tempFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a mock test context
	ctx := NewMockTestContextSimple()
	ctx.SetOutput("output_file_path", tempFile)
	ctx.SetOutput("output_content", "test content")

	// Test the function with content from output
	AssertFileContent(t, ctx)

	// Test the function with explicit content
	AssertFileContent(t, ctx, "test content")
}

// TestAssertOutputEquals tests the AssertOutputEquals function
func TestAssertOutputEquals(t *testing.T) {
	// Create a mock test context
	ctx := NewMockTestContextSimple()
	ctx.SetOutput("test_output", "test value")

	// Test the function
	AssertOutputEquals(t, ctx, "test_output", "test value")
}

// TestAssertOutputContains tests the AssertOutputContains function
func TestAssertOutputContains(t *testing.T) {
	// Create a mock test context
	ctx := NewMockTestContextSimple()
	ctx.SetOutput("test_output", "this is a test value")

	// Test the function
	AssertOutputContains(t, ctx, "test_output", "test value")
}

// TestAssertOutputMatches tests the AssertOutputMatches function
func TestAssertOutputMatches(t *testing.T) {
	// Create a mock test context
	ctx := NewMockTestContextSimple()
	ctx.SetOutput("test_output", "abc123")

	// Test the function
	AssertOutputMatches(t, ctx, "test_output", "^[a-z]+[0-9]+$")
}

// TestAssertOutputNotEmpty tests the AssertOutputNotEmpty function
func TestAssertOutputNotEmpty(t *testing.T) {
	// Create a mock test context
	ctx := NewMockTestContextSimple()
	ctx.SetOutput("test_output", "not empty")

	// Test the function
	AssertOutputNotEmpty(t, ctx, "test_output")
}

// TestAssertOutputEmpty tests the AssertOutputEmpty function
func TestAssertOutputEmpty(t *testing.T) {
	// Create a mock test context
	ctx := NewMockTestContextSimple()
	ctx.SetOutput("test_output", "")

	// Test the function
	AssertOutputEmpty(t, ctx, "test_output")
}

// TestGetTerraform tests the GetTerraform method
func TestGetTerraform(t *testing.T) {
	// Create a mock test context
	ctx := NewMockTestContextSimple()

	// Test the function
	result := ctx.GetTerraform()
	assert.NotNil(t, result)
}
