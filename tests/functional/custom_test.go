package functional

import (
	"testing"

	"github.com/caylent-solutions/terraform-test-framework/internal/assertions"
	"github.com/caylent-solutions/terraform-test-framework/internal/testctx"
)

// Example of a custom test function that can be run on all examples
func TestCustomAssertions(t *testing.T) {
	// Run all examples with default configs
	results := testctx.RunAllExamples(t, "../..", nil)
	
	// Define a custom test function
	customTest := func(t *testing.T, ctx testctx.TestContext) {
		t.Logf("Running custom test on example: %s", ctx.Config.Name)
		
		// Example of a custom assertion
		if ctx.Config.Name == "basic" {
			assertions.AssertOutputEquals(t, ctx, "output_content", "framework test")
		}
	}
	
	// Run the custom test on all examples
	testctx.RunCustomTests(t, results, customTest)
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
	testctx.RunAllExamplesWithTests(t, "../..", nil, test1, test2)
}