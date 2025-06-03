package testctx

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdempotencyEnabled(t *testing.T) {
	// Test default behavior (should be enabled)
	os.Unsetenv("TERRATEST_IDEMPOTENCY")
	assert.True(t, IdempotencyEnabled(), "Idempotency should be enabled by default")

	// Test with environment variable set to "false"
	os.Setenv("TERRATEST_IDEMPOTENCY", "false")
	assert.False(t, IdempotencyEnabled(), "Idempotency should be disabled when TERRATEST_IDEMPOTENCY=false")

	// Test with environment variable set to something else
	os.Setenv("TERRATEST_IDEMPOTENCY", "anything")
	assert.True(t, IdempotencyEnabled(), "Idempotency should be enabled when TERRATEST_IDEMPOTENCY is not 'false'")

	// Cleanup
	os.Unsetenv("TERRATEST_IDEMPOTENCY")
}

func TestNewTestContext(t *testing.T) {
	// Create a test context
	ctx := NewTestContext("test-example", nil)

	// Verify the context
	assert.Equal(t, "test-example", ctx.Name, "Context name should match")
	assert.NotNil(t, ctx.Terraform, "Terraform options should not be nil")
	assert.NotNil(t, ctx.TerraformVars, "Terraform variables should not be nil")
}

func TestTestConfig(t *testing.T) {
	// Create a test config
	config := TestConfig{
		Name: "test-config",
		ExtraVars: map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		},
	}

	// Verify the config
	assert.Equal(t, "test-config", config.Name, "Config name should match")
	assert.Equal(t, "value1", config.ExtraVars["key1"], "ExtraVars key1 should match")
	assert.Equal(t, 42, config.ExtraVars["key2"], "ExtraVars key2 should match")
}

// Mock functions for testing RunExample and RunAllExamples would be more complex
// and would require mocking the filesystem and terraform execution
// These would be added in a more comprehensive test suite