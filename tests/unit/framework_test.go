package unit

import (
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
	"github.com/stretchr/testify/assert"
)

// Test the IdempotencyEnabled function
func TestIdempotencyEnabledFramework(t *testing.T) {
	// Test default behavior
	assert.True(t, testctx.IdempotencyEnabled(), "Idempotency should be enabled by default")
}

// Test the runner functions
func TestInitTerraform(t *testing.T) {
	config := testctx.TestConfig{
		Name: "test-config",
		ExtraVars: map[string]interface{}{
			"var1": "value1",
			"var2": 42,
		},
	}

	options := testctx.InitTerraform("/path/to/example", config)
	assert.Equal(t, "/path/to/example", options.TerraformDir)
	assert.Equal(t, config.ExtraVars, options.Vars)
}

func TestRun(t *testing.T) {
	config := testctx.TestConfig{
		Name: "test-config",
		ExtraVars: map[string]interface{}{
			"var1": "value1",
		},
	}

	ctx := testctx.Run("/path/to/example", config)
	assert.Equal(t, config, ctx.Config)
	assert.Equal(t, "/path/to/example", ctx.ExamplePath)
	assert.Equal(t, "test-config", ctx.Name)
	assert.NotNil(t, ctx.Terraform)
	assert.Equal(t, "/path/to/example", ctx.Terraform.TerraformDir)
	assert.Equal(t, config.ExtraVars, ctx.Terraform.Vars)
}

func TestRunCustomTests(t *testing.T) {
	// Create test contexts
	results := map[string]testctx.TestContext{
		"example1": {
			Name: "example1",
			Config: testctx.TestConfig{
				Name: "example1",
			},
		},
		"example2": {
			Name: "example2",
			Config: testctx.TestConfig{
				Name: "example2",
			},
		},
	}

	// Track which examples were tested
	tested := make(map[string]bool)

	// Create test function
	testFunc := func(t *testing.T, ctx testctx.TestContext) {
		tested[ctx.Name] = true
	}

	// Run custom tests
	testctx.RunCustomTests(t, results, testFunc)

	// Verify all examples were tested
	assert.True(t, tested["example1"])
	assert.True(t, tested["example2"])
	assert.Equal(t, 2, len(tested))
}
