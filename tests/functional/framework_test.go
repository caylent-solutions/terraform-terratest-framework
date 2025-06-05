package functional

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/assertions"
	"github.com/caylent-solutions/terraform-terratest-framework/internal/testctx"
	"github.com/stretchr/testify/assert"
)

func TestRunSingleExample(t *testing.T) {
	// Get the path to the test fixture
	fixtureDir, err := filepath.Abs("../terraform-module-test-fixtures/example-basic")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Create a temporary directory for test output
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "test-output.txt")

	// Run the example
	ctx := testctx.RunSingleExample(t, fixtureDir, ".", testctx.TestConfig{
		Name: "basic-example",
		ExtraVars: map[string]interface{}{
			"output_content":  "Test Content",
			"output_filename": outputFile,
		},
	})

	// Verify outputs
	assertions.AssertOutputEquals(t, ctx, "output_content", "Test Content")

	// Get the output path from terraform
	outputPath := ctx.GetOutput(t, "output_file_path")

	// Handle paths with ./ prefix
	var cleanPath string
	if strings.HasPrefix(outputPath, "./") {
		cleanPath = outputPath[2:] // Remove the ./ prefix
	} else {
		cleanPath = outputPath
	}

	// Try the original output path first
	_, err = os.Stat(outputPath)
	if err == nil {
		// File exists at the original path
		content, err := os.ReadFile(outputPath)
		if err != nil {
			t.Fatalf("Failed to read file at %s: %v", outputPath, err)
		}

		if string(content) != "Test Content" {
			t.Fatalf("Expected file content to be %q, but got %q", "Test Content", string(content))
		}
		return
	}

	// Try the cleaned path
	_, err = os.Stat(cleanPath)
	if err == nil {
		// File exists at the cleaned path
		content, err := os.ReadFile(cleanPath)
		if err != nil {
			t.Fatalf("Failed to read file at %s: %v", cleanPath, err)
		}

		if string(content) != "Test Content" {
			t.Fatalf("Expected file content to be %q, but got %q", "Test Content", string(content))
		}
		return
	}

	// Try the original output file path
	_, err = os.Stat(outputFile)
	if err == nil {
		// File exists at the original output file path
		content, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("Failed to read file at %s: %v", outputFile, err)
		}

		if string(content) != "Test Content" {
			t.Fatalf("Expected file content to be %q, but got %q", "Test Content", string(content))
		}
		return
	}

	// None of the paths worked
	t.Fatalf("Expected file to exist at %s, %s, or %s, but got errors", outputPath, cleanPath, outputFile)

	// Read from the output path
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read file at %s: %v", outputPath, err)
	}

	if string(content) != "Test Content" {
		t.Fatalf("Expected file content to be %q, but got %q", "Test Content", string(content))
	}
}

func TestRunMultipleExamples(t *testing.T) {
	// This test verifies that multiple examples can be run concurrently

	// Get the path to the test fixtures directory
	fixturesDir, err := filepath.Abs("../terraform-module-test-fixtures")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Create a temporary directory for test results
	tempDir := t.TempDir()

	// Run the basic example
	basicCtx := testctx.RunSingleExample(t, fixturesDir, "example-basic", testctx.TestConfig{
		Name: "basic-example",
		ExtraVars: map[string]interface{}{
			"output_content":  "Basic Content",
			"output_filename": filepath.Join(tempDir, "basic.txt"),
		},
	})

	// Run the advanced example
	advancedCtx := testctx.RunSingleExample(t, fixturesDir, "example-advanced", testctx.TestConfig{
		Name: "advanced-example",
		ExtraVars: map[string]interface{}{
			"output_content":  "Advanced Content",
			"output_filename": filepath.Join(tempDir, "advanced.txt"),
		},
	})

	// Verify basic example outputs
	assertions.AssertOutputEquals(t, basicCtx, "output_content", "Basic Content")

	// Verify advanced example outputs
	assertions.AssertOutputEquals(t, advancedCtx, "output_content", "Advanced Content")

	// Get output paths
	basicFilePath := basicCtx.GetOutput(t, "output_file_path")
	advancedFilePath := advancedCtx.GetOutput(t, "output_file_path")

	// Handle paths with ./ prefix
	var basicCleanPath, advancedCleanPath string

	// For basic file path
	if strings.HasPrefix(basicFilePath, "./") {
		basicCleanPath = basicFilePath[2:] // Remove the ./ prefix
	} else {
		basicCleanPath = basicFilePath
	}

	// For advanced file path
	if strings.HasPrefix(advancedFilePath, "./") {
		advancedCleanPath = advancedFilePath[2:] // Remove the ./ prefix
	} else {
		advancedCleanPath = advancedFilePath
	}

	// Verify basic file exists - try both paths
	basicExists := false
	_, err1 := os.Stat(basicFilePath)
	if err1 == nil {
		basicExists = true
	} else {
		_, err2 := os.Stat(basicCleanPath)
		if err2 == nil {
			basicExists = true
		}
	}
	assert.True(t, basicExists, "basic.txt should exist at either %s or %s", basicFilePath, basicCleanPath)

	// Verify advanced file exists - try both paths
	advancedExists := false
	_, err3 := os.Stat(advancedFilePath)
	if err3 == nil {
		advancedExists = true
	} else {
		_, err4 := os.Stat(advancedCleanPath)
		if err4 == nil {
			advancedExists = true
		}
	}
	assert.True(t, advancedExists, "advanced.txt should exist at either %s or %s", advancedFilePath, advancedCleanPath)
}

func TestDiscoverAndRunAllTests(t *testing.T) {
	// Skip this test in CI environments as it requires more setup
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping in CI environment")
	}

	// Get the path to the test fixtures directory
	fixturesDir, err := filepath.Abs("../terraform-module-test-fixtures")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Create configs for the examples
	configs := make(map[string]testctx.TestConfig)
	entries, err := os.ReadDir(fixturesDir)
	if err != nil {
		t.Fatalf("Failed to read examples directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Only include directories that start with "example-"
		name := entry.Name()
		if len(name) >= 8 && name[:8] == "example-" {
			configs[name] = testctx.TestConfig{
				Name:      name,
				ExtraVars: map[string]interface{}{},
			}
		}
	}

	// Run all examples
	for name, config := range configs {
		examplePath := filepath.Join(fixturesDir, name)
		t.Run(fmt.Sprintf("Example_%s", name), func(t *testing.T) {
			// Run the example
			ctx := testctx.RunExample(t, examplePath, config)

			// Verify outputs before cleanup
			output := ctx.GetOutput(t, "output_content")
			assert.NotEmpty(t, output, "output_content should not be empty")

			// Get the output path
			filePath := ctx.GetOutput(t, "output_file_path")

			// Handle paths with ./ prefix
			var cleanPath string
			if strings.HasPrefix(filePath, "./") {
				cleanPath = filePath[2:] // Remove the ./ prefix
			} else {
				cleanPath = filePath
			}

			// Get the absolute path if it's not already absolute
			if !filepath.IsAbs(filePath) {
				absPath, err := filepath.Abs(filepath.Join(examplePath, filePath))
				if err == nil {
					filePath = absPath
				}
			}

			if !filepath.IsAbs(cleanPath) {
				absPath, err := filepath.Abs(filepath.Join(examplePath, cleanPath))
				if err == nil {
					cleanPath = absPath
				}
			}

			// Verify that file exists - try both paths
			fileExists := false
			_, err1 := os.Stat(filePath)
			if err1 == nil {
				fileExists = true
			} else {
				_, err2 := os.Stat(cleanPath)
				if err2 == nil {
					fileExists = true
				}
			}
			assert.True(t, fileExists, "File should exist at either %s or %s", filePath, cleanPath)
		})
	}
}
