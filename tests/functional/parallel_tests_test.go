package functional

import (
	"os"
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
)

func TestParallelTestsFlag(t *testing.T) {
	// Skip this test in CI environments
	if os.Getenv("CI") != "" {
		t.Skip("Skipping test in CI environment")
	}

	// Test with parallel tests disabled
	originalValue := os.Getenv("TERRATEST_DISABLE_PARALLEL_TESTS")
	defer os.Setenv("TERRATEST_DISABLE_PARALLEL_TESTS", originalValue)

	// Test when env var is set to "true"
	os.Setenv("TERRATEST_DISABLE_PARALLEL_TESTS", "true")
	if testctx.IsParallelTestsEnabled() {
		t.Error("Expected parallel tests to be disabled when env var is set to 'true'")
	}

	// Test when env var is set to "false"
	os.Setenv("TERRATEST_DISABLE_PARALLEL_TESTS", "false")
	if !testctx.IsParallelTestsEnabled() {
		t.Error("Expected parallel tests to be enabled when env var is set to 'false'")
	}
}
