package examples

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/testctx"
	"github.com/stretchr/testify/assert"
)

func TestFindAllExamples(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "examples-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create examples directory
	examplesDir := filepath.Join(tempDir, "examples")
	err = os.Mkdir(examplesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create examples directory: %v", err)
	}

	// Create example directories
	exampleDirs := []string{
		"example1",
		"example2",
		"not-an-example",
	}

	for _, dir := range exampleDirs {
		err := os.Mkdir(filepath.Join(examplesDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create a file to ensure it's ignored
	file, err := os.Create(filepath.Join(examplesDir, "some-file.txt"))
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	file.Close()

	// Test the function
	examples := FindAllExamples(t, tempDir)

	// Verify results
	assert.Equal(t, 3, len(examples))

	// Check that all directories were found
	foundNames := make(map[string]bool)
	for _, example := range examples {
		foundNames[example.Name] = true
		assert.Equal(t, example.Name, example.Config.Name)
		assert.Equal(t, filepath.Join(examplesDir, example.Name), example.Path)
		assert.NotNil(t, example.Config.ExtraVars)
	}

	assert.True(t, foundNames["example1"])
	assert.True(t, foundNames["example2"])
	assert.True(t, foundNames["not-an-example"])
}

func TestConfigureExamples(t *testing.T) {
	// Create test examples
	examples := []Example{
		{
			Name: "example1",
			Path: "/path/to/example1",
			Config: testctx.TestConfig{
				Name:      "example1",
				ExtraVars: map[string]interface{}{},
			},
		},
		{
			Name: "example2",
			Path: "/path/to/example2",
			Config: testctx.TestConfig{
				Name:      "example2",
				ExtraVars: map[string]interface{}{},
			},
		},
	}

	// Create configurator function
	configurator := func(example Example) testctx.TestConfig {
		return testctx.TestConfig{
			Name: example.Name + "-configured",
			ExtraVars: map[string]interface{}{
				"configured": true,
				"name":       example.Name,
			},
		}
	}

	// Test the function
	configs := ConfigureExamples(examples, configurator)

	// Verify results
	assert.Equal(t, 2, len(configs))

	// Check example1 config
	config1, exists := configs["example1"]
	assert.True(t, exists)
	assert.Equal(t, "example1-configured", config1.Name)
	assert.Equal(t, true, config1.ExtraVars["configured"])
	assert.Equal(t, "example1", config1.ExtraVars["name"])

	// Check example2 config
	config2, exists := configs["example2"]
	assert.True(t, exists)
	assert.Equal(t, "example2-configured", config2.Name)
	assert.Equal(t, true, config2.ExtraVars["configured"])
	assert.Equal(t, "example2", config2.ExtraVars["name"])
}
