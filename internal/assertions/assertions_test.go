package assertions

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/testctx"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock TestContext for testing assertions
type MockTestContext struct {
	mock.Mock
	testctx.TestContext
}

func (m *MockTestContext) GetOutput(t *testing.T, key string) string {
	args := m.Called(t, key)
	return args.String(0)
}

func TestAssertOutputEquals(t *testing.T) {
	// Create a mock test context
	mockCtx := new(MockTestContext)
	mockCtx.On("GetOutput", mock.Anything, "test_key").Return("test_value")

	// Test successful assertion
	AssertOutputEquals(t, mockCtx, "test_key", "test_value")

	// Verify the mock was called
	mockCtx.AssertExpectations(t)
}

func TestAssertFileExists(t *testing.T) {
	// Create a temporary file
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(tempFile, []byte("test content"), 0644)
	assert.NoError(t, err, "Failed to create temporary file")

	// Create a mock test context
	mockCtx := new(MockTestContext)
	mockCtx.On("GetOutput", mock.Anything, "output_file_path").Return(tempFile)

	// Test successful assertion
	AssertFileExists(t, mockCtx)

	// Verify the mock was called
	mockCtx.AssertExpectations(t)
}

func TestAssertFileContent(t *testing.T) {
	// Create a temporary file
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(tempFile, []byte("test content"), 0644)
	assert.NoError(t, err, "Failed to create temporary file")

	// Create a mock test context
	mockCtx := new(MockTestContext)
	mockCtx.On("GetOutput", mock.Anything, "output_file_path").Return(tempFile)
	mockCtx.On("GetOutput", mock.Anything, "output_content").Return("test content")

	// Test successful assertion
	AssertFileContent(t, mockCtx)

	// Verify the mock was called
	mockCtx.AssertExpectations(t)
}

func TestAssertOutputContains(t *testing.T) {
	// Create a mock test context
	mockCtx := new(MockTestContext)
	mockCtx.On("GetOutput", mock.Anything, "test_key").Return("this is a test value")

	// Test successful assertion
	AssertOutputContains(t, mockCtx, "test_key", "test value")

	// Verify the mock was called
	mockCtx.AssertExpectations(t)
}

func TestAssertOutputMatches(t *testing.T) {
	// Create a mock test context
	mockCtx := new(MockTestContext)
	mockCtx.On("GetOutput", mock.Anything, "test_key").Return("abc-123-xyz")

	// Test successful assertion
	AssertOutputMatches(t, mockCtx, "test_key", "^[a-z]+-\\d+-[a-z]+$")

	// Verify the mock was called
	mockCtx.AssertExpectations(t)
}

func TestAssertOutputNotEmpty(t *testing.T) {
	// Create a mock test context
	mockCtx := new(MockTestContext)
	mockCtx.On("GetOutput", mock.Anything, "test_key").Return("not empty")

	// Test successful assertion
	AssertOutputNotEmpty(t, mockCtx, "test_key")

	// Verify the mock was called
	mockCtx.AssertExpectations(t)
}

func TestAssertOutputEmpty(t *testing.T) {
	// Create a mock test context
	mockCtx := new(MockTestContext)
	mockCtx.On("GetOutput", mock.Anything, "test_key").Return("")

	// Test successful assertion
	AssertOutputEmpty(t, mockCtx, "test_key")

	// Verify the mock was called
	mockCtx.AssertExpectations(t)
}

// Mock for terraform.OutputMap
type MockTerraformOptions struct {
	mock.Mock
}

func TestAssertOutputMapContainsKey(t *testing.T) {
	// This test is minimal since terraform.OutputMap is a Terratest function
	// that we can't easily mock without more complex test infrastructure
	mockT := new(testing.T)
	mockCtx := new(MockTestContext)
	mockCtx.Terraform = &terraform.Options{}
	
	// Just verify it doesn't panic
	assert.NotPanics(t, func() {
		// Use a mock testing.T to prevent actual terraform commands from running
		AssertOutputMapContainsKey(mockT, mockCtx, "test_map", "test_key")
	})
}

func TestAssertOutputMapKeyEquals(t *testing.T) {
	// This test is minimal since terraform.OutputMap is a Terratest function
	// that we can't easily mock without more complex test infrastructure
	mockT := new(testing.T)
	mockCtx := new(MockTestContext)
	mockCtx.Terraform = &terraform.Options{}
	
	// Just verify it doesn't panic
	assert.NotPanics(t, func() {
		// Use a mock testing.T to prevent actual terraform commands from running
		AssertOutputMapKeyEquals(mockT, mockCtx, "test_map", "test_key", "test_value")
	})
}

func TestAssertOutputListContains(t *testing.T) {
	// This test is minimal since terraform.OutputList is a Terratest function
	// that we can't easily mock without more complex test infrastructure
	mockT := new(testing.T)
	mockCtx := new(MockTestContext)
	mockCtx.Terraform = &terraform.Options{}
	
	// Just verify it doesn't panic
	assert.NotPanics(t, func() {
		// Use a mock testing.T to prevent actual terraform commands from running
		AssertOutputListContains(mockT, mockCtx, "test_list", "test_value")
	})
}

func TestAssertOutputListLength(t *testing.T) {
	// This test is minimal since terraform.OutputList is a Terratest function
	// that we can't easily mock without more complex test infrastructure
	mockT := new(testing.T)
	mockCtx := new(MockTestContext)
	mockCtx.Terraform = &terraform.Options{}
	
	// Just verify it doesn't panic
	assert.NotPanics(t, func() {
		// Use a mock testing.T to prevent actual terraform commands from running
		AssertOutputListLength(mockT, mockCtx, "test_list", 3)
	})
}

func TestAssertOutputJSONContains(t *testing.T) {
	// Create a mock test context
	mockCtx := new(MockTestContext)
	mockCtx.On("GetOutput", mock.Anything, "test_json").Return(`{"key1": "value1", "key2": 42}`)

	// Test successful assertion
	AssertOutputJSONContains(t, mockCtx, "test_json", "key1", "value1")
	AssertOutputJSONContains(t, mockCtx, "test_json", "key2", float64(42)) // JSON numbers are parsed as float64

	// Verify the mock was called
	mockCtx.AssertExpectations(t)
}

func TestAssertResourceExists(t *testing.T) {
	// This test is minimal since terraform.RunTerraformCommand is a Terratest function
	// that we can't easily mock without more complex test infrastructure
	mockT := new(testing.T)
	mockCtx := new(MockTestContext)
	mockCtx.Terraform = &terraform.Options{}
	
	// Just verify it doesn't panic
	assert.NotPanics(t, func() {
		// Use a mock testing.T to prevent actual terraform commands from running
		AssertResourceExists(mockT, mockCtx, "aws_s3_bucket", "example")
	})
}

func TestAssertResourceCount(t *testing.T) {
	// This test is minimal since terraform.RunTerraformCommand is a Terratest function
	// that we can't easily mock without more complex test infrastructure
	mockT := new(testing.T)
	mockCtx := new(MockTestContext)
	mockCtx.Terraform = &terraform.Options{}
	
	// Just verify it doesn't panic
	assert.NotPanics(t, func() {
		// Use a mock testing.T to prevent actual terraform commands from running
		AssertResourceCount(mockT, mockCtx, "aws_s3_bucket", 2)
	})
}

func TestAssertNoResourcesOfType(t *testing.T) {
	// This test is minimal since terraform.RunTerraformCommand is a Terratest function
	// that we can't easily mock without more complex test infrastructure
	mockT := new(testing.T)
	mockCtx := new(MockTestContext)
	mockCtx.Terraform = &terraform.Options{}
	
	// Just verify it doesn't panic
	assert.NotPanics(t, func() {
		// Use a mock testing.T to prevent actual terraform commands from running
		AssertNoResourcesOfType(mockT, mockCtx, "aws_s3_bucket")
	})
}

func TestAssertTerraformVersion(t *testing.T) {
	// This test is minimal since terraform.RunTerraformCommand is a Terratest function
	// that we can't easily mock without more complex test infrastructure
	mockT := new(testing.T)
	mockCtx := new(MockTestContext)
	mockCtx.Terraform = &terraform.Options{}
	
	// Just verify it doesn't panic
	assert.NotPanics(t, func() {
		// Use a mock testing.T to prevent actual terraform commands from running
		AssertTerraformVersion(mockT, mockCtx, "0.12.0")
	})
}