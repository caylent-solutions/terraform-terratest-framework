package testctx

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

// TestConfig holds configuration for a single test
type TestConfig struct {
	Name      string
	ExtraVars map[string]interface{}
}

// TestContext combines test configuration with terraform options
type TestContext struct {
	Config        TestConfig
	Terraform     *terraform.Options
	ExamplePath   string
	Name          string
	TerraformVars map[string]interface{}
}

// GetOutput retrieves a terraform output value by key
func (ctx TestContext) GetOutput(t testing.TB, key string) string {
	return terraform.Output(t, ctx.Terraform, key)
}

// GetTerraform returns the terraform options
func (ctx TestContext) GetTerraform() *terraform.Options {
	return ctx.Terraform
}

// NewTestContext creates a new test context
func NewTestContext(examplePath string, vars map[string]interface{}) TestContext {
	if vars == nil {
		vars = make(map[string]interface{})
	}
	return TestContext{
		Name:          examplePath,
		Terraform:     &terraform.Options{},
		TerraformVars: vars,
		ExamplePath:   examplePath,
	}
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
