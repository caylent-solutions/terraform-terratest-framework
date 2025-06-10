package testctx

import (
	"os"
	"strings"
)

// IsParallelTestsEnabled checks if parallel tests within fixtures are enabled
// Returns false if TERRATEST_DISABLE_PARALLEL_TESTS is set to "true"
func IsParallelTestsEnabled() bool {
	val := os.Getenv("TERRATEST_DISABLE_PARALLEL_TESTS")
	return !strings.EqualFold(val, "true")
}
