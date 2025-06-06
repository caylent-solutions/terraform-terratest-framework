package unit

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

// TestAssertIdempotentMock tests the AssertIdempotent function with mocking
func TestAssertIdempotentMock(t *testing.T) {
	// Skip this test since we can't easily mock terraform.Plan
	t.Skip("Skipping test that requires mocking terraform.Plan function")
}
