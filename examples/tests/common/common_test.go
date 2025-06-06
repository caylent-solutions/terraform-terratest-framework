package common

import (
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
)

// TestAllExamples runs all examples and their tests in parallel
// This is the main entry point for running all tests
func TestAllExamples(t *testing.T) {
	// This will discover and run all examples and their tests in parallel
	testctx.DiscoverAndRunAllTests(t, "../..", func(t *testing.T, ctx testctx.TestContext) {
		// Run standard tests on each example
		RunStandardTests(t, ctx, []string{}, []string{})
	})
}

// Example of a common test that would run on all examples
func TestCommonS3BucketEncryption(t *testing.T) {
	// This test would be run for each example
	// The actual implementation would be provided by the user
	// based on their specific requirements
}

// Example of another common test
func TestCommonTagCompliance(t *testing.T) {
	// This test would verify that all resources have required tags
	// The actual implementation would be provided by the user
}
