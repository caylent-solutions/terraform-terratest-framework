package common

import (
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/testctx"
	"github.com/caylent-solutions/terraform-terratest-framework/tests/unit"
)

// Example of using various assertions in a test
func TestAssertionsExample(t *testing.T) {
	// Skip this test in normal runs as it's just an example
	t.Skip("This is an example test that demonstrates how to use assertions")

	// Run an example
	ctx := testctx.RunSingleExample(t, "../../", "example", testctx.TestConfig{
		Name: "assertions-example",
		ExtraVars: map[string]interface{}{
			"region": "us-west-2",
		},
	})

	// Basic assertions
	unit.AssertOutputEquals(t, ctx, "region", "us-west-2")
	unit.AssertOutputContains(t, ctx, "instance_id", "i-")
	unit.AssertOutputMatches(t, ctx, "instance_id", "^i-[a-f0-9]{17}$")
	unit.AssertOutputNotEmpty(t, ctx, "vpc_id")
	unit.AssertOutputEmpty(t, ctx, "error_message")

	// File assertions
	filePath := unit.AssertFileExists(t, ctx)
	t.Logf("File exists at: %s", filePath)
	unit.AssertFileContent(t, ctx)

	// Collection assertions
	unit.AssertOutputMapContainsKey(t, ctx, "tags", "Name")
	unit.AssertOutputMapKeyEquals(t, ctx, "tags", "Environment", "test")
	unit.AssertOutputListContains(t, ctx, "subnet_ids", "subnet-12345")
	unit.AssertOutputListLength(t, ctx, "subnet_ids", 3)

	// JSON assertions
	unit.AssertOutputJSONContains(t, ctx, "config_json", "region", "us-west-2")

	// Resource assertions
	unit.AssertResourceExists(t, ctx, "aws_vpc", "main")
	unit.AssertResourceCount(t, ctx, "aws_subnet", 3)
	unit.AssertNoResourcesOfType(t, ctx, "aws_db_instance")

	// Environment assertions
	unit.AssertTerraformVersion(t, ctx, "1.0.0")

	// Idempotency assertion
	unit.AssertIdempotent(t, ctx)
}
