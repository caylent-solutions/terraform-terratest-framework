package unit

import (
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/testctx"
	"github.com/stretchr/testify/assert"
)

func TestConfigInit(t *testing.T) {
	cfg := testctx.TestConfig{
		Name: "test",
		ExtraVars: map[string]interface{}{
			"output_content": "example",
		},
	}
	assert.Equal(t, "test", cfg.Name)
	assert.Equal(t, "example", cfg.ExtraVars["output_content"])
}
