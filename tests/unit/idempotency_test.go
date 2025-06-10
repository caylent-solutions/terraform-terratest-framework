package unit

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/idempotency"
	"github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
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

func TestIdempotencyPackage(t *testing.T) {
	// This is a basic test to ensure the package can be imported and compiled
	assert.NotPanics(t, func() {
		_ = idempotency.Test
		_ = idempotency.TestAll
	})
}
