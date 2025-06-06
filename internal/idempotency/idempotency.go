package idempotency

import (
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// Test runs an idempotency test on the given context
// This test is run by default unless explicitly disabled via TERRATEST_IDEMPOTENCY=false
func Test(t *testing.T, ctx testctx.TestContext) bool {
	if !testctx.IdempotencyEnabled() {
		t.Logf("Idempotency testing disabled for %s via TERRATEST_IDEMPOTENCY=false", ctx.Config.Name)
		return true
	}

	t.Logf("Running idempotency test for %s", ctx.Config.Name)
	planResult := terraform.Plan(t, ctx.Terraform)

	if planResult != "" {
		t.Errorf("Idempotency test failed for %s: Terraform plan would make changes: %s", ctx.Config.Name, planResult)
		return false
	}

	t.Logf("Idempotency test passed for %s", ctx.Config.Name)
	return true
}

// TestAll runs idempotency tests on all contexts
func TestAll(t *testing.T, contexts map[string]testctx.TestContext) {
	for name, ctx := range contexts {
		t.Run("Idempotency_"+name, func(t *testing.T) {
			Test(t, ctx)
		})
	}
}
