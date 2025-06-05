package assertions

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestAssertIdempotent(t *testing.T) {
	// Save original function and restore after test
	origPlan := terraform.Plan
	defer func() { terraform.Plan = origPlan }()

	// Replace with mock implementation
	terraform.Plan = func(t testing.TB, options *terraform.Options) string {
		return "No changes. Your infrastructure matches the configuration."
	}

	// Create a mock test context
	mockCtx := NewMockTestContextSimple()

	// Test successful assertion
	AssertIdempotent(t, mockCtx)
}
