package functional

import (
	"testing"

	"github.com/caylent-solutions/terraform-test-framework/internal/examples"
	"github.com/caylent-solutions/terraform-test-framework/internal/idempotency"
	"github.com/caylent-solutions/terraform-test-framework/internal/testctx"
)

func TestParallelExamples(t *testing.T) {
	// Find all examples
	allExamples := examples.FindAllExamples(t, "../..")
	
	// Skip test if no examples found
	if len(allExamples) == 0 {
		t.Skip("No examples found to test")
	}
	
	// Configure examples with default settings
	configs := examples.ConfigureExamples(allExamples, func(ex examples.Example) testctx.TestConfig {
		return testctx.TestConfig{
			Name: ex.Name,
			ExtraVars: map[string]interface{}{
				"output_content":  "parallel test",
				"output_filename": "parallel-test.txt",
			},
		}
	})
	
	// Run all examples in parallel
	results := testctx.RunAllExamples(t, "../..", configs)
	
	// Run idempotency tests on all examples
	idempotency.TestAll(t, results)
}