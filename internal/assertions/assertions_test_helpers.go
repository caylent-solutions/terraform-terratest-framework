package assertions

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
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

// Mock for terraform.OutputMap
type MockTerraformOptions struct {
	mock.Mock
}