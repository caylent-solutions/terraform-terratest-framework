package assertions

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

// Helper functions for mocking terraform functions in tests

// MockTerraformFunctions provides mock implementations for terraform functions
type MockTerraformFunctions struct {
	OutputMapFunc           func(t testing.TB, options *terraform.Options, key string) map[string]string
	OutputListFunc          func(t testing.TB, options *terraform.Options, key string) []string
	RunTerraformCommandFunc func(t testing.TB, options *terraform.Options, args ...string) string
	PlanFunc                func(t testing.TB, options *terraform.Options) string
}

// DefaultMockTerraformFunctions returns a MockTerraformFunctions with default implementations
func DefaultMockTerraformFunctions() *MockTerraformFunctions {
	return &MockTerraformFunctions{
		OutputMapFunc: func(t testing.TB, options *terraform.Options, key string) map[string]string {
			return map[string]string{
				"test_key":    "test_value",
				"another_key": "another_value",
			}
		},
		OutputListFunc: func(t testing.TB, options *terraform.Options, key string) []string {
			return []string{"value1", "value2", "value3"}
		},
		RunTerraformCommandFunc: func(t testing.TB, options *terraform.Options, args ...string) string {
			if len(args) >= 1 {
				if args[0] == "state" && args[1] == "list" {
					return "aws_s3_bucket.example\naws_dynamodb_table.test\naws_lambda_function.handler"
				} else if args[0] == "version" {
					return "Terraform v1.5.0\non darwin_amd64"
				}
			}
			return ""
		},
		PlanFunc: func(t testing.TB, options *terraform.Options) string {
			return "No changes. Your infrastructure matches the configuration."
		},
	}
}

// SetupTestWithMocks sets up a test with mocked terraform functions
func SetupTestWithMocks(t *testing.T) *MockTestContext {
	// Create a mock test context
	mockCtx := new(MockTestContext)
	mockCtx.terraform = &terraform.Options{}

	// Set up mock functions
	// Note: We're not actually replacing the terraform functions anymore
	// as that would require modifying the package-level functions

	return mockCtx
}
