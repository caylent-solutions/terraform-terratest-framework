package common

import (
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/assertions"
	"github.com/caylent-solutions/terraform-terratest-framework/internal/testctx"
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
	assertions.AssertOutputEquals(t, ctx, "region", "us-west-2")
	assertions.AssertOutputContains(t, ctx, "instance_id", "i-")
	assertions.AssertOutputMatches(t, ctx, "instance_id", "^i-[a-f0-9]{17}$")
	assertions.AssertOutputNotEmpty(t, ctx, "vpc_id")
	assertions.AssertOutputEmpty(t, ctx, "error_message")
	
	// File assertions
	filePath := assertions.AssertFileExists(t, ctx)
	t.Logf("File exists at: %s", filePath)
	assertions.AssertFileContent(t, ctx)
	
	// Collection assertions
	assertions.AssertOutputMapContainsKey(t, ctx, "tags", "Name")
	assertions.AssertOutputMapKeyEquals(t, ctx, "tags", "Environment", "test")
	assertions.AssertOutputListContains(t, ctx, "subnet_ids", "subnet-12345")
	assertions.AssertOutputListLength(t, ctx, "subnet_ids", 3)
	
	// JSON assertions
	assertions.AssertOutputJSONContains(t, ctx, "config_json", "region", "us-west-2")
	
	// Resource assertions
	assertions.AssertResourceExists(t, ctx, "aws_vpc", "main")
	assertions.AssertResourceCount(t, ctx, "aws_subnet", 3)
	assertions.AssertNoResourcesOfType(t, ctx, "aws_db_instance")
	
	// Environment assertions
	assertions.AssertTerraformVersion(t, ctx, "1.0.0")
	
	// Idempotency assertion
	assertions.AssertIdempotent(t, ctx)
}