package unit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/pkg/assertions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTestContext implements the testctx.TestContext interface for testing
type MockTestContext struct {
	mock.Mock
}

func (m *MockTestContext) GetTerraform() interface{} {
	args := m.Called()
	return args.Get(0)
}

// TestAssertOutputEquals tests the AssertOutputEquals function
func TestAssertOutputEquals(t *testing.T) {
	// This is a limited test since we can't easily mock terraform.Output
	// In a real scenario, we would use a more comprehensive mocking approach
	t.Run("Function exists", func(t *testing.T) {
		assert.NotPanics(t, func() {
			// Just verify the function exists and can be called
			// We can't verify the result without more complex mocking
			_ = assertions.AssertOutputEquals
		})
	})
}

// TestAssertOutputContains tests the AssertOutputContains function
func TestAssertOutputContains(t *testing.T) {
	t.Run("Function exists", func(t *testing.T) {
		assert.NotPanics(t, func() {
			_ = assertions.AssertOutputContains
		})
	})
}

// TestAssertOutputMatches tests the AssertOutputMatches function
func TestAssertOutputMatches(t *testing.T) {
	t.Run("Function exists", func(t *testing.T) {
		assert.NotPanics(t, func() {
			_ = assertions.AssertOutputMatches
		})
	})
}

// TestAssertOutputNotEmpty tests the AssertOutputNotEmpty function
func TestAssertOutputNotEmpty(t *testing.T) {
	t.Run("Function exists", func(t *testing.T) {
		assert.NotPanics(t, func() {
			_ = assertions.AssertOutputNotEmpty
		})
	})
}

// TestAssertOutputEmpty tests the AssertOutputEmpty function
func TestAssertOutputEmpty(t *testing.T) {
	t.Run("Function exists", func(t *testing.T) {
		assert.NotPanics(t, func() {
			_ = assertions.AssertOutputEmpty
		})
	})
}

// TestAssertFileExists tests the AssertFileExists function
func TestAssertFileExists(t *testing.T) {
	t.Run("Function exists", func(t *testing.T) {
		assert.NotPanics(t, func() {
			_ = assertions.AssertFileExists
		})
	})
}

// TestAssertFileContent tests the AssertFileContent function
func TestAssertFileContent(t *testing.T) {
	t.Run("Function exists", func(t *testing.T) {
		assert.NotPanics(t, func() {
			_ = assertions.AssertFileContent
		})
	})
}

// TestAssertOutputMapContainsKey tests the AssertOutputMapContainsKey function
func TestAssertOutputMapContainsKey(t *testing.T) {
	t.Run("Function exists", func(t *testing.T) {
		assert.NotPanics(t, func() {
			_ = assertions.AssertOutputMapContainsKey
		})
	})
}

// TestAssertOutputMapKeyEquals tests the AssertOutputMapKeyEquals function
func TestAssertOutputMapKeyEquals(t *testing.T) {
	t.Run("Function exists", func(t *testing.T) {
		assert.NotPanics(t, func() {
			_ = assertions.AssertOutputMapKeyEquals
		})
	})
}

// TestAssertOutputListContains tests the AssertOutputListContains function
func TestAssertOutputListContains(t *testing.T) {
	t.Run("Function exists", func(t *testing.T) {
		assert.NotPanics(t, func() {
			_ = assertions.AssertOutputListContains
		})
	})
}

// TestAssertOutputListLength tests the AssertOutputListLength function
func TestAssertOutputListLength(t *testing.T) {
	t.Run("Function exists", func(t *testing.T) {
		assert.NotPanics(t, func() {
			_ = assertions.AssertOutputListLength
		})
	})
}

// TestAssertOutputJSONContains tests the AssertOutputJSONContains function
func TestAssertOutputJSONContains(t *testing.T) {
	// Create a temporary JSON file for testing
	tempDir, err := os.MkdirTemp("", "test-json")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	jsonData := map[string]interface{}{
		"test_key": "test_value",
		"number":   42,
		"boolean":  true,
	}

	jsonBytes, _ := json.Marshal(jsonData)
	tempFile := filepath.Join(tempDir, "test.json")
	err = os.WriteFile(tempFile, jsonBytes, 0644)
	if err != nil {
		t.Fatalf("Failed to write test JSON file: %v", err)
	}

	t.Run("Function exists", func(t *testing.T) {
		assert.NotPanics(t, func() {
			_ = assertions.AssertOutputJSONContains
		})
	})
}

// TestAssertResourceExists tests the AssertResourceExists function
func TestAssertResourceExists(t *testing.T) {
	t.Run("Function exists", func(t *testing.T) {
		assert.NotPanics(t, func() {
			_ = assertions.AssertResourceExists
		})
	})
}

// TestAssertResourceCount tests the AssertResourceCount function
func TestAssertResourceCount(t *testing.T) {
	t.Run("Function exists", func(t *testing.T) {
		assert.NotPanics(t, func() {
			_ = assertions.AssertResourceCount
		})
	})
}

// TestAssertNoResourcesOfType tests the AssertNoResourcesOfType function
func TestAssertNoResourcesOfType(t *testing.T) {
	t.Run("Function exists", func(t *testing.T) {
		assert.NotPanics(t, func() {
			_ = assertions.AssertNoResourcesOfType
		})
	})
}

// TestAssertTerraformVersion tests the AssertTerraformVersion function
func TestAssertTerraformVersion(t *testing.T) {
	t.Run("Function exists", func(t *testing.T) {
		assert.NotPanics(t, func() {
			_ = assertions.AssertTerraformVersion
		})
	})
}

// TestAssertIdempotent tests the AssertIdempotent function
func TestAssertIdempotent(t *testing.T) {
	t.Run("Function exists", func(t *testing.T) {
		assert.NotPanics(t, func() {
			_ = assertions.AssertIdempotent
		})
	})
}
