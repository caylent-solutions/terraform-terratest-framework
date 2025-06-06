package functional

import (
	"strings"
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
	"github.com/gruntwork-io/terratest/modules/terraform"
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

	// Create a map to store test contexts
	results := make(map[string]testctx.TestContext)

	// Process each example manually to control when resources are destroyed
	for name, config := range configs {
		t.Run(name, func(t *testing.T) {
			// Initialize terraform in the example directory
			examplePath := "../terraform-module-test-fixtures/" + name
			ctx := testctx.Run(examplePath, config)

			// Apply the terraform configuration
			terraform.InitAndApply(t, ctx.Terraform)

			// Store the context
			results[name] = ctx

			// Run assertions on the outputs before destroying resources
			if name == "example-basic" {
				output := terraform.Output(t, ctx.Terraform, "output_content")
				if output != "Hello, World!" {
					t.Errorf("Expected output_content to be 'Hello, World!', but got '%s'", output)
				}
			} else if name == "example-advanced" {
				output := terraform.Output(t, ctx.Terraform, "output_content")
				if output != "Advanced Example" {
					t.Errorf("Expected output_content to be 'Advanced Example', but got '%s'", output)
				}
			}

			// Run idempotency test if enabled
			if testctx.IdempotencyEnabled() {
				t.Log("Running idempotency test...")
				planOutput := terraform.Plan(t, ctx.Terraform)
				// Check if the plan output contains "No changes" or "no changes"
				if strings.Contains(planOutput, "No changes") || strings.Contains(planOutput, "no changes") {
					t.Log("Idempotency test passed")
				} else {
					t.Fatalf("Idempotency test failed: Terraform plan would make changes: %s", planOutput)
				}
			}

			// Clean up resources after assertions
			terraform.Destroy(t, ctx.Terraform)
		})
	}
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
