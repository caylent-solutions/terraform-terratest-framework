package functional

import (
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/assertions"
	"github.com/caylent-solutions/terraform-terratest-framework/internal/testctx"
)

// Example of a custom test function that can be run on all examples
func TestCustomAssertions(t *testing.T) {
	// Create configs for the examples
	configs := map[string]testctx.TestConfig{
		"example-basic": {
			Name: "example-basic",
		},
		"example-advanced": {
			Name: "example-advanced",
		},
	}

	// Define a custom test function
	customTest := func(t *testing.T, ctx testctx.TestContext) {
		t.Logf("Running custom test on example: %s", ctx.Config.Name)

		// Example of a custom assertion
		if ctx.Config.Name == "example-basic" {
			assertions.AssertOutputEquals(t, ctx, "output_content", "Hello, World!")
		} else if ctx.Config.Name == "example-advanced" {
			assertions.AssertOutputEquals(t, ctx, "output_content", "Advanced Example")
		}
	}

	// Run all examples with the custom test function
	testctx.RunAllExamplesWithTests(t, "../terraform-module-test-fixtures", configs, customTest)
}

// Example of running all examples with multiple custom tests in one go
func TestMultipleCustomTests(t *testing.T) {
	// Define multiple custom test functions
	test1 := func(t *testing.T, ctx testctx.TestContext) {
		t.Logf("Running test 1 on example: %s", ctx.Config.Name)
		// Custom assertions for test 1
	}

	test2 := func(t *testing.T, ctx testctx.TestContext) {
		t.Logf("Running test 2 on example: %s", ctx.Config.Name)
		// Custom assertions for test 2
	}

	// Run all examples and then run both custom tests on each example
	testctx.RunAllExamplesWithTests(t, "../terraform-module-test-fixtures", nil, test1, test2)
}
