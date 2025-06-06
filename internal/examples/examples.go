package examples

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/testctx"
)

// Example represents a single terraform example
type Example struct {
	Name   string
	Path   string
	Config testctx.TestConfig
}

// FindAllExamples discovers all examples in the examples directory
func FindAllExamples(t *testing.T, moduleRootPath string) []Example {
	examplesPath := filepath.Join(moduleRootPath, "examples")
	entries, err := os.ReadDir(examplesPath)
	if err != nil {
		t.Fatalf("Failed to read examples directory: %v", err)
	}

	var examples []Example
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		examples = append(examples, Example{
			Name: entry.Name(),
			Path: filepath.Join(examplesPath, entry.Name()),
			Config: testctx.TestConfig{
				Name:      entry.Name(),
				ExtraVars: map[string]interface{}{},
			},
		})
	}

	return examples
}

// ConfigureExamples allows customizing the configuration for each example
func ConfigureExamples(examples []Example, configurator func(Example) testctx.TestConfig) map[string]testctx.TestConfig {
	configs := make(map[string]testctx.TestConfig)

	for _, example := range examples {
		configs[example.Name] = configurator(example)
	}

	return configs
}
