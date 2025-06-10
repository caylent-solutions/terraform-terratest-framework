package unit

import (
	"os"
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
)

func TestIsParallelTestsEnabled(t *testing.T) {
	// Save original env var value to restore later
	originalValue := os.Getenv("TERRATEST_DISABLE_PARALLEL_TESTS")
	defer os.Setenv("TERRATEST_DISABLE_PARALLEL_TESTS", originalValue)

	// Test when env var is not set
	os.Unsetenv("TERRATEST_DISABLE_PARALLEL_TESTS")
	if !testctx.IsParallelTestsEnabled() {
		t.Error("Expected parallel tests to be enabled when env var is not set")
	}

	// Test when env var is set to "true"
	os.Setenv("TERRATEST_DISABLE_PARALLEL_TESTS", "true")
	if testctx.IsParallelTestsEnabled() {
		t.Error("Expected parallel tests to be disabled when env var is set to 'true'")
	}

	// Test when env var is set to "TRUE" (case insensitive)
	os.Setenv("TERRATEST_DISABLE_PARALLEL_TESTS", "TRUE")
	if testctx.IsParallelTestsEnabled() {
		t.Error("Expected parallel tests to be disabled when env var is set to 'TRUE'")
	}

	// Test when env var is set to "false"
	os.Setenv("TERRATEST_DISABLE_PARALLEL_TESTS", "false")
	if !testctx.IsParallelTestsEnabled() {
		t.Error("Expected parallel tests to be enabled when env var is set to 'false'")
	}
}
