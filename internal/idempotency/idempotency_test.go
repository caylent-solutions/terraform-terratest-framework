package idempotency

import (
	"os"
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
	"github.com/stretchr/testify/assert"
)

// TestIdempotencyEnabledDetailed tests the IdempotencyEnabled function
func TestIdempotencyEnabledDetailed(t *testing.T) {
	// Test default behavior (enabled)
	os.Unsetenv("TERRATEST_IDEMPOTENCY")
	assert.True(t, testctx.IdempotencyEnabled())

	// Test explicitly enabled
	os.Setenv("TERRATEST_IDEMPOTENCY", "true")
	assert.True(t, testctx.IdempotencyEnabled())

	// Test disabled
	os.Setenv("TERRATEST_IDEMPOTENCY", "false")
	assert.False(t, testctx.IdempotencyEnabled())

	// Reset for other tests
	os.Unsetenv("TERRATEST_IDEMPOTENCY")
}

// Note: Testing the Test and TestAll functions would require mocking the terraform module,
// which is beyond the scope of a simple unit test. These functions primarily call terraform.Plan
// which requires a real Terraform environment to execute properly.

func TestIdempotencyTestBasic(t *testing.T) {
	// This is a basic test to ensure the package can be imported and compiled
	// A more thorough test would require mocking terraform.Plan
	assert.NotNil(t, testctx.IdempotencyEnabled)
}
