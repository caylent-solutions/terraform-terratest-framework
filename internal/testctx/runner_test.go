package testctx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitTerraform(t *testing.T) {
	config := TestConfig{
		Name: "test-config",
		ExtraVars: map[string]interface{}{
			"var1": "value1",
			"var2": 42,
		},
	}

	options := InitTerraform("/path/to/example", config)
	assert.Equal(t, "/path/to/example", options.TerraformDir)
	assert.Equal(t, config.ExtraVars, options.Vars)
}

func TestRun(t *testing.T) {
	config := TestConfig{
		Name: "test-config",
		ExtraVars: map[string]interface{}{
			"var1": "value1",
		},
	}

	ctx := Run("/path/to/example", config)
	assert.Equal(t, config, ctx.Config)
	assert.Equal(t, "/path/to/example", ctx.ExamplePath)
	assert.Equal(t, "test-config", ctx.Name)
	assert.NotNil(t, ctx.Terraform)
	assert.Equal(t, "/path/to/example", ctx.Terraform.TerraformDir)
	assert.Equal(t, config.ExtraVars, ctx.Terraform.Vars)
}

// TestRunCustomTestsImplementation tests the RunCustomTests function
// Renamed to avoid conflict with the function in custom_test.go
func TestRunCustomTestsImplementation(t *testing.T) {
	// Create test contexts
	results := map[string]TestContext{
		"example1": {
			Name: "example1",
			Config: TestConfig{
				Name: "example1",
			},
		},
		"example2": {
			Name: "example2",
			Config: TestConfig{
				Name: "example2",
			},
		},
	}

	// Track which examples were tested
	tested := make(map[string]bool)

	// Create test function
	testFunc := func(t *testing.T, ctx TestContext) {
		tested[ctx.Name] = true
	}

	// Run custom tests
	RunCustomTests(t, results, testFunc)

	// Verify all examples were tested
	assert.True(t, tested["example1"])
	assert.True(t, tested["example2"])
	assert.Equal(t, 2, len(tested))
}

// Note: The following functions are difficult to unit test without mocking:
// - RunExample (requires terraform.InitAndApply)
// - RunAllExamplesWithTests (depends on RunAllExamples)
// - DiscoverAndRunAllTests (depends on os.ReadDir and RunAllExamples)
// - RunAllExamples (depends on os.ReadDir and RunExample)
// - RunSingleExample (depends on os.Stat and RunExample)
//
// These would typically be tested with integration tests or with mocking.
