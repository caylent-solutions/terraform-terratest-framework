package functional

import (
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
)

func TestIdempotency(t *testing.T) {
	// Configure examples with default settings
	configs := map[string]testctx.TestConfig{
		"example-basic": {
			Name: "example-basic",
			ExtraVars: map[string]interface{}{
				"output_content":  "framework test",
				"output_filename": "framework-test.txt",
			},
		},
		"example-advanced": {
			Name: "example-advanced",
			ExtraVars: map[string]interface{}{
				"output_content":  "framework test",
				"output_filename": "framework-test.txt",
			},
		},
	}

	// Run all examples in parallel and test idempotency
	_ = testctx.RunAllExamples(t, "../terraform-module-test-fixtures", configs)

	// No need for additional assertions here since idempotency is already tested in RunExample
}
