package unit

import (
	"os"
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/testctx"
	"github.com/stretchr/testify/assert"
)

func TestIdempotencyEnabled(t *testing.T) {
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