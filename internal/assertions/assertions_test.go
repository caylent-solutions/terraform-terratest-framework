package assertions

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	testing_terratest "github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/assert"
	"github.com/undefinedlabs/go-mpatch"
)

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

func TestAssertOutputMapContainsKey(t *testing.T) {
	// Create a patch for the OutputMap function
	patch, err := mpatch.PatchMethod(terraform.OutputMap, func(t testing_terratest.TestingT, options *terraform.Options, key string) map[string]string {
		return map[string]string{
			"test_key":    "test_value",
			"another_key": "another_value",
		}
	})
	if err != nil {
		t.Fatalf("Failed to patch OutputMap: %v", err)
	}
	defer patch.Unpatch()

	// Create a mock test context
	mockCtx := NewMockTestContextSimple()

	// Test successful assertion
	AssertOutputMapContainsKey(t, mockCtx, "test_map", "test_key")
}

func TestAssertOutputMapKeyEquals(t *testing.T) {
	// Create a patch for the OutputMap function
	patch, err := mpatch.PatchMethod(terraform.OutputMap, func(t testing_terratest.TestingT, options *terraform.Options, key string) map[string]string {
		return map[string]string{
			"test_key":    "test_value",
			"another_key": "another_value",
		}
	})
	if err != nil {
		t.Fatalf("Failed to patch OutputMap: %v", err)
	}
	defer patch.Unpatch()

	// Create a mock test context
	mockCtx := NewMockTestContextSimple()

	// Test successful assertion
	AssertOutputMapKeyEquals(t, mockCtx, "test_map", "test_key", "test_value")
}

func TestAssertOutputListContains(t *testing.T) {
	// Create a patch for the OutputList function
	patch, err := mpatch.PatchMethod(terraform.OutputList, func(t testing_terratest.TestingT, options *terraform.Options, key string) []string {
		return []string{"value1", "value2", "value3"}
	})
	if err != nil {
		t.Fatalf("Failed to patch OutputList: %v", err)
	}
	defer patch.Unpatch()

	// Create a mock test context
	mockCtx := NewMockTestContextSimple()

	// Test successful assertion
	AssertOutputListContains(t, mockCtx, "test_list", "value2")
}

func TestAssertOutputListLength(t *testing.T) {
	// Create a patch for the OutputList function
	patch, err := mpatch.PatchMethod(terraform.OutputList, func(t testing_terratest.TestingT, options *terraform.Options, key string) []string {
		return []string{"value1", "value2", "value3"}
	})
	if err != nil {
		t.Fatalf("Failed to patch OutputList: %v", err)
	}
	defer patch.Unpatch()

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
	// Create a patch for the RunTerraformCommand function
	patch, err := mpatch.PatchMethod(terraform.RunTerraformCommand, func(t testing_terratest.TestingT, options *terraform.Options, args ...string) string {
		if len(args) >= 2 && args[0] == "state" && args[1] == "list" {
			return "aws_s3_bucket.example\naws_dynamodb_table.test\naws_lambda_function.handler"
		}
		return ""
	})
	if err != nil {
		t.Fatalf("Failed to patch RunTerraformCommand: %v", err)
	}
	defer patch.Unpatch()

	// Create a mock test context
	mockCtx := NewMockTestContextSimple()

	// Test successful assertion
	AssertResourceExists(t, mockCtx, "aws_s3_bucket", "example")
}

func TestAssertResourceCount(t *testing.T) {
	// Create a patch for the RunTerraformCommand function
	patch, err := mpatch.PatchMethod(terraform.RunTerraformCommand, func(t testing_terratest.TestingT, options *terraform.Options, args ...string) string {
		if len(args) >= 2 && args[0] == "state" && args[1] == "list" {
			return "aws_s3_bucket.example\naws_dynamodb_table.test\naws_lambda_function.handler"
		}
		return ""
	})
	if err != nil {
		t.Fatalf("Failed to patch RunTerraformCommand: %v", err)
	}
	defer patch.Unpatch()

	// Create a mock test context
	mockCtx := NewMockTestContextSimple()

	// Test successful assertion
	AssertResourceCount(t, mockCtx, "aws_s3_bucket", 1)
}

func TestAssertNoResourcesOfType(t *testing.T) {
	// Create a patch for the RunTerraformCommand function
	patch, err := mpatch.PatchMethod(terraform.RunTerraformCommand, func(t testing_terratest.TestingT, options *terraform.Options, args ...string) string {
		if len(args) >= 2 && args[0] == "state" && args[1] == "list" {
			return "aws_s3_bucket.example\naws_dynamodb_table.test\naws_lambda_function.handler"
		}
		return ""
	})
	if err != nil {
		t.Fatalf("Failed to patch RunTerraformCommand: %v", err)
	}
	defer patch.Unpatch()

	// Create a mock test context
	mockCtx := NewMockTestContextSimple()

	// Test successful assertion
	AssertNoResourcesOfType(t, mockCtx, "aws_ec2_instance")
}

func TestAssertTerraformVersion(t *testing.T) {
	// Create a patch for the RunTerraformCommand function
	patch, err := mpatch.PatchMethod(terraform.RunTerraformCommand, func(t testing_terratest.TestingT, options *terraform.Options, args ...string) string {
		if len(args) >= 1 && args[0] == "version" {
			return "Terraform v1.5.0\non darwin_amd64"
		}
		return ""
	})
	if err != nil {
		t.Fatalf("Failed to patch RunTerraformCommand: %v", err)
	}
	defer patch.Unpatch()

	// Create a mock test context
	mockCtx := NewMockTestContextSimple()

	// Test successful assertion
	AssertTerraformVersion(t, mockCtx, "1.4.0")
}
