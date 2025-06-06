package unit

import (
	"os"
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/testctx"
	"github.com/stretchr/testify/assert"
)

func TestConfigInit(t *testing.T) {
	cfg := testctx.TestConfig{
		Name: "test",
		ExtraVars: map[string]interface{}{
			"output_content": "example",
		},
	}
	assert.Equal(t, "test", cfg.Name)
	assert.Equal(t, "example", cfg.ExtraVars["output_content"])
}

// Test the TestConfig structure
func TestTestConfig(t *testing.T) {
	config := testctx.TestConfig{
		Name: "test",
		ExtraVars: map[string]interface{}{
			"key": "value",
		},
	}

	assert.Equal(t, "test", config.Name)
	assert.Equal(t, "value", config.ExtraVars["key"])
}

func TestNewTestContext(t *testing.T) {
	// Test with nil vars
	ctx := testctx.NewTestContext("example-path", nil)
	assert.Equal(t, "example-path", ctx.Name)
	assert.Equal(t, "example-path", ctx.ExamplePath)
	assert.NotNil(t, ctx.Terraform)
	assert.NotNil(t, ctx.TerraformVars)

	// Test with provided vars
	vars := map[string]interface{}{
		"test_var": "test_value",
	}
	ctx = testctx.NewTestContext("another-path", vars)
	assert.Equal(t, "another-path", ctx.Name)
	assert.Equal(t, "another-path", ctx.ExamplePath)
	assert.Equal(t, vars, ctx.TerraformVars)
}

// Test the IdempotencyEnabled function
func TestIdempotencyEnabled(t *testing.T) {
	// Test default behavior (enabled)
	os.Unsetenv("TERRATEST_IDEMPOTENCY")
	assert.True(t, testctx.IdempotencyEnabled())

	// Test explicitly enabled
	os.Setenv("TERRATEST_IDEMPOTENCY", "true")
	assert.True(t, testctx.IdempotencyEnabled())

	// Test disabled
	os.Setenv("TERRATEST_IDEMPOTENCY", "false")
	assert.False(t, testctx.IdempotencyEnabled())

	// Reset for other tests
	os.Unsetenv("TERRATEST_IDEMPOTENCY")
}
