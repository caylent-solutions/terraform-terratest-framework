package testctx

import (
	"testing"
)

// CustomTestFunc defines a function type for custom tests that can be run on each example
type CustomTestFunc func(t *testing.T, ctx TestContext)

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
