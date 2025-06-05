package assertions

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

// MockTestContext is a mock implementation of the TestContext interface for testing
type MockTestContext struct {
	terraform *terraform.Options
	outputs   map[string]string
}

// GetOutput returns a mock output value
func (m *MockTestContext) GetOutput(t testing.TB, key string) string {
	if m.outputs == nil {
		m.outputs = map[string]string{
			"output_content":   "test content",
			"output_file_path": "/tmp/test-output.txt",
		}
	}
	return m.outputs[key]
}

// GetTerraform returns the terraform options
func (m *MockTestContext) GetTerraform() *terraform.Options {
	return m.terraform
}
