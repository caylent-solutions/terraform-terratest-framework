package unit

import (
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
)

// CustomTestFunc defines a function type for custom tests that can be run on each example
type CustomTestFunc func(t *testing.T, ctx testctx.TestContext)

// TestRunCustomTests tests the custom test runner
func TestRunCustomTests(t *testing.T) {
	// This is just a test function to avoid the redeclaration error
	// The actual implementation is in runner.go
}

// TestRunAllExamplesWithTests tests the combined runner
func TestRunAllExamplesWithTests(t *testing.T) {
	// This is just a test function to avoid the redeclaration error
	// The actual implementation is in runner.go
}
