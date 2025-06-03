package functional

import (
	"testing"

	"github.com/caylent-solutions/terraform-test-framework/internal/testctx"
	"github.com/caylent-solutions/terraform-test-framework/internal/assertions"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestFrameworkScaffold(t *testing.T) {
	config := testctx.TestConfig{
		Name: "scaffold",
		ExtraVars: map[string]interface{}{
			"output_content":  "framework test",
			"output_filename": "framework-test.txt",
		},
	}
	ctx := testctx.Run("../../examples/basic", config)
	defer terraform.Destroy(t, ctx.Terraform)
	terraform.InitAndApply(t, ctx.Terraform)

	assertions.AssertFileExists(t, ctx)
	assertions.AssertFileContent(t, ctx, "framework test")
}
