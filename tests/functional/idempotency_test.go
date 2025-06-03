package functional

import (
	"testing"

	"github.com/caylent-solutions/terraform-test-framework/internal/assertions"
	"github.com/caylent-solutions/terraform-test-framework/internal/examples"
	"github.com/caylent-solutions/terraform-test-framework/internal/testctx"
)

func TestIdempotency(t *testing.T) {
	// Find all examples
	allExamples := examples.FindAllExamples(t, "../..")
	
	// Configure examples with default settings
	configs := examples.ConfigureExamples(allExamples, func(ex examples.Example) testctx.TestConfig {
		return testctx.TestConfig{
			Name: ex.Name,
			ExtraVars: map[string]interface{}{
				"output_content":  "framework test",
				"output_filename": "framework-test.txt",
			},
		}
	})
	
	// Run all examples in parallel and test idempotency
	results := testctx.RunAllExamples(t, "../..", configs)
	
	// Additional assertions can be added here for each example
	for name, ctx := range results {
		t.Run("Assertions_"+name, func(t *testing.T) {
			assertions.AssertIdempotent(t, ctx)
		})
	}
}