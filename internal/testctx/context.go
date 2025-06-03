package testctx

import (
	"os"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

// TestConfig holds configuration for a single test
type TestConfig struct {
	Name      string
	ExtraVars map[string]interface{}
}

// TestContext combines test configuration with terraform options
type TestContext struct {
	Config    TestConfig
	Terraform *terraform.Options
	ExamplePath string
}

// IdempotencyEnabled checks if idempotency testing is enabled
// Returns true (enabled) by default unless TERRATEST_IDEMPOTENCY=false
// The test will run if:
// - TERRATEST_IDEMPOTENCY environment variable doesn't exist
// - TERRATEST_IDEMPOTENCY is set to any value other than "false"
// The test will be skipped if:
// - TERRATEST_IDEMPOTENCY is set to "false"
func IdempotencyEnabled() bool {
	val := os.Getenv("TERRATEST_IDEMPOTENCY")
	return val != "false"
}
