package assertions

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	testing_terratest "github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTestContextSimple implements a simple test context for testing
type MockTestContextSimple struct {
	terraform *terraform.Options
	outputs   map[string]string
}

func NewMockTestContextSimple() *MockTestContextSimple {
	return &MockTestContextSimple{
		terraform: &terraform.Options{},
		outputs: map[string]string{
			"test_key":         "test_value",
			"output_file_path": "/tmp/test.txt",
			"output_content":   "test content",
		},
	}
}

func (m *MockTestContextSimple) GetOutput(t testing.TB, key string) string {
	if val, ok := m.outputs[key]; ok {
		return val
	}
	return ""
}

func (m *MockTestContextSimple) GetTerraform() *terraform.Options {
	return m.terraform
}

func TestAssertOutputEquals(t *testing.T) {
	// Create a mock test context
	mockCtx := NewMockTestContextSimple()

	// Test successful assertion
	AssertOutputEquals(t, mockCtx, "test_key", "test_value")
}

func TestAssertFileExists(t *testing.T) {
	// Create a temporary file
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(tempFile, []byte("test content"), 0644)
	assert.NoError(t, err, "Failed to create temporary file")

	// Create a mock test context
	mockCtx := NewMockTestContextSimple()
	mockCtx.outputs["output_file_path"] = tempFile

	// Test successful assertion
	AssertFileExists(t, mockCtx)
}

func TestAssertFileContent(t *testing.T) {
	// Create a temporary file
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(tempFile, []byte("test content"), 0644)
	assert.NoError(t, err, "Failed to create temporary file")

	// Create a mock test context
	mockCtx := NewMockTestContextSimple()
	mockCtx.outputs["output_file_path"] = tempFile
	mockCtx.outputs["output_content"] = "test content"

	// Test successful assertion
	AssertFileContent(t, mockCtx)
}

func TestAssertOutputContains(t *testing.T) {
	// Create a mock test context
	mockCtx := NewMockTestContextSimple()
	mockCtx.outputs["test_key"] = "this is a test value"

	// Test successful assertion
	AssertOutputContains(t, mockCtx, "test_key", "test value")
}

func TestAssertOutputMatches(t *testing.T) {
	// Create a mock test context
	mockCtx := NewMockTestContextSimple()
	mockCtx.outputs["test_key"] = "abc-123-xyz"

	// Test successful assertion
	AssertOutputMatches(t, mockCtx, "test_key", "^[a-z]+-\\d+-[a-z]+$")
}

func TestAssertOutputNotEmpty(t *testing.T) {
	// Create a mock test context
	mockCtx := NewMockTestContextSimple()
	mockCtx.outputs["test_key"] = "not empty"

	// Test successful assertion
	AssertOutputNotEmpty(t, mockCtx, "test_key")
}

func TestAssertOutputEmpty(t *testing.T) {
	// Create a mock test context
	mockCtx := NewMockTestContextSimple()
	mockCtx.outputs["test_key"] = ""

	// Test successful assertion
	AssertOutputEmpty(t, mockCtx, "test_key")
}

// Mock for terraform.OutputMap
type MockTerraformOptions struct {
	mock.Mock
}

// Store original functions to restore later
var origOutputMap func(t testing_terratest.TestingT, options *terraform.Options, key string) map[string]string
var origOutputList func(t testing_terratest.TestingT, options *terraform.Options, key string) []string
var origRunTerraformCommand func(t testing_terratest.TestingT, options *terraform.Options, args ...string) string

func init() {
	// Store original functions
	origOutputMap = terraform.OutputMap
	origOutputList = terraform.OutputList
	origRunTerraformCommand = terraform.RunTerraformCommand
}

func TestAssertOutputMapContainsKey(t *testing.T) {
	// Save original function and restore after test
	defer func() { terraform.OutputMap = origOutputMap }()

	// Create a local variable for the mock function
	mockOutputMap := func(t testing_terratest.TestingT, options *terraform.Options, key string) map[string]string {
		return map[string]string{
			"test_key":    "test_value",
			"another_key": "another_value",
		}
	}

	// Assign the mock function to terraform.OutputMap
	terraform.OutputMap = mockOutputMap

	// Create a mock test context
	mockCtx := NewMockTestContextSimple()

	// Test successful assertion
	AssertOutputMapContainsKey(t, mockCtx, "test_map", "test_key")
}

func TestAssertOutputMapKeyEquals(t *testing.T) {
	// Save original function and restore after test
	defer func() { terraform.OutputMap = origOutputMap }()

	// Create a local variable for the mock function
	mockOutputMap := func(t testing_terratest.TestingT, options *terraform.Options, key string) map[string]string {
		return map[string]string{
			"test_key":    "test_value",
			"another_key": "another_value",
		}
	}

	// Assign the mock function to terraform.OutputMap
	terraform.OutputMap = mockOutputMap

	// Create a mock test context
	mockCtx := NewMockTestContextSimple()

	// Test successful assertion
	AssertOutputMapKeyEquals(t, mockCtx, "test_map", "test_key", "test_value")
}

func TestAssertOutputListContains(t *testing.T) {
	// Save original function and restore after test
	defer func() { terraform.OutputList = origOutputList }()

	// Create a local variable for the mock function
	mockOutputList := func(t testing_terratest.TestingT, options *terraform.Options, key string) []string {
		return []string{"value1", "value2", "value3"}
	}

	// Assign the mock function to terraform.OutputList
	terraform.OutputList = mockOutputList

	// Create a mock test context
	mockCtx := NewMockTestContextSimple()

	// Test successful assertion
	AssertOutputListContains(t, mockCtx, "test_list", "value2")
}

func TestAssertOutputListLength(t *testing.T) {
	// Save original function and restore after test
	defer func() { terraform.OutputList = origOutputList }()

	// Create a local variable for the mock function
	mockOutputList := func(t testing_terratest.TestingT, options *terraform.Options, key string) []string {
		return []string{"value1", "value2", "value3"}
	}

	// Assign the mock function to terraform.OutputList
	terraform.OutputList = mockOutputList

	// Create a mock test context
	mockCtx := NewMockTestContextSimple()

	// Test successful assertion
	AssertOutputListLength(t, mockCtx, "test_list", 3)
}

func TestAssertOutputJSONContains(t *testing.T) {
	// Create a mock test context
	mockCtx := NewMockTestContextSimple()
	mockCtx.outputs["test_json"] = `{"key1": "value1", "key2": 42}`

	// Test successful assertion
	AssertOutputJSONContains(t, mockCtx, "test_json", "key1", "value1")
	AssertOutputJSONContains(t, mockCtx, "test_json", "key2", float64(42)) // JSON numbers are parsed as float64
}

func TestAssertResourceExists(t *testing.T) {
	// Save original function and restore after test
	defer func() { terraform.RunTerraformCommand = origRunTerraformCommand }()

	// Create a local variable for the mock function
	mockRunTerraformCommand := func(t testing_terratest.TestingT, options *terraform.Options, args ...string) string {
		if len(args) >= 2 && args[0] == "state" && args[1] == "list" {
			return "aws_s3_bucket.example\naws_dynamodb_table.test\naws_lambda_function.handler"
		}
		return ""
	}

	// Assign the mock function to terraform.RunTerraformCommand
	terraform.RunTerraformCommand = mockRunTerraformCommand

	// Create a mock test context
	mockCtx := NewMockTestContextSimple()

	// Test successful assertion
	AssertResourceExists(t, mockCtx, "aws_s3_bucket", "example")
}

func TestAssertResourceCount(t *testing.T) {
	// Save original function and restore after test
	defer func() { terraform.RunTerraformCommand = origRunTerraformCommand }()

	// Create a local variable for the mock function
	mockRunTerraformCommand := func(t testing_terratest.TestingT, options *terraform.Options, args ...string) string {
		if len(args) >= 2 && args[0] == "state" && args[1] == "list" {
			return "aws_s3_bucket.example\naws_dynamodb_table.test\naws_lambda_function.handler"
		}
		return ""
	}

	// Assign the mock function to terraform.RunTerraformCommand
	terraform.RunTerraformCommand = mockRunTerraformCommand

	// Create a mock test context
	mockCtx := NewMockTestContextSimple()

	// Test successful assertion
	AssertResourceCount(t, mockCtx, "aws_s3_bucket", 1)
}

func TestAssertNoResourcesOfType(t *testing.T) {
	// Save original function and restore after test
	defer func() { terraform.RunTerraformCommand = origRunTerraformCommand }()

	// Create a local variable for the mock function
	mockRunTerraformCommand := func(t testing_terratest.TestingT, options *terraform.Options, args ...string) string {
		if len(args) >= 2 && args[0] == "state" && args[1] == "list" {
			return "aws_s3_bucket.example\naws_dynamodb_table.test\naws_lambda_function.handler"
		}
		return ""
	}

	// Assign the mock function to terraform.RunTerraformCommand
	terraform.RunTerraformCommand = mockRunTerraformCommand

	// Create a mock test context
	mockCtx := NewMockTestContextSimple()

	// Test successful assertion
	AssertNoResourcesOfType(t, mockCtx, "aws_ec2_instance")
}

func TestAssertTerraformVersion(t *testing.T) {
	// Save original function and restore after test
	defer func() { terraform.RunTerraformCommand = origRunTerraformCommand }()

	// Create a local variable for the mock function
	mockRunTerraformCommand := func(t testing_terratest.TestingT, options *terraform.Options, args ...string) string {
		if len(args) >= 1 && args[0] == "version" {
			return "Terraform v1.5.0\non darwin_amd64"
		}
		return ""
	}

	// Assign the mock function to terraform.RunTerraformCommand
	terraform.RunTerraformCommand = mockRunTerraformCommand

	// Create a mock test context
	mockCtx := NewMockTestContextSimple()

	// Test successful assertion
	AssertTerraformVersion(t, mockCtx, "1.4.0")
}
