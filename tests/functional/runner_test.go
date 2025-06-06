package functional

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/testctx"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestFrameworkScaffold(t *testing.T) {
	// Get the absolute path to the test fixture
	fixtureDir, err := filepath.Abs("../terraform-module-test-fixtures/example-basic")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Create a temporary directory for test output
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "framework-test.txt")

	config := testctx.TestConfig{
		Name: "scaffold",
		ExtraVars: map[string]interface{}{
			"output_content":  "framework test",
			"output_filename": outputFile,
		},
	}
	ctx := testctx.Run(fixtureDir, config)

	// Apply Terraform configuration
	terraform.InitAndApply(t, ctx.Terraform)

	// Verify the file exists directly BEFORE destroying resources
	fileFound := false

	// Try the direct output file path first
	_, err = os.Stat(outputFile)
	if err == nil {
		fileFound = true
		content, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("Failed to read file at %s: %v", outputFile, err)
		}

		if string(content) != "framework test" {
			t.Fatalf("Expected file content to be %q, but got %q", "framework test", string(content))
		}
	}

	// If not found, try with the path from the output
	if !fileFound {
		outputPath := terraform.Output(t, ctx.Terraform, "output_file_path")

		// Try the original output path
		_, err1 := os.Stat(outputPath)
		if err1 == nil {
			fileFound = true
			content, err := os.ReadFile(outputPath)
			if err != nil {
				t.Fatalf("Failed to read file at %s: %v", outputPath, err)
			}

			if string(content) != "framework test" {
				t.Fatalf("Expected file content to be %q, but got %q", "framework test", string(content))
			}
		} else {
			// Try without the ./ prefix if it exists
			var cleanPath string
			if strings.HasPrefix(outputPath, "./") {
				cleanPath = outputPath[2:] // Remove the ./ prefix
				_, err2 := os.Stat(cleanPath)
				if err2 == nil {
					fileFound = true
					content, err := os.ReadFile(cleanPath)
					if err != nil {
						t.Fatalf("Failed to read file at %s: %v", cleanPath, err)
					}

					if string(content) != "framework test" {
						t.Fatalf("Expected file content to be %q, but got %q", "framework test", string(content))
					}
				}
			}

			// Try with the ./ prefix if it doesn't exist
			if !fileFound && !strings.HasPrefix(outputPath, "./") {
				withPrefix := "./" + outputPath
				_, err3 := os.Stat(withPrefix)
				if err3 == nil {
					fileFound = true
					content, err := os.ReadFile(withPrefix)
					if err != nil {
						t.Fatalf("Failed to read file at %s: %v", withPrefix, err)
					}

					if string(content) != "framework test" {
						t.Fatalf("Expected file content to be %q, but got %q", "framework test", string(content))
					}
				}
			}

			// If still not found, report error
			if !fileFound {
				if cleanPath != "" {
					t.Fatalf("Expected file to exist at %s, %s, or %s, but got errors", outputFile, outputPath, cleanPath)
				} else {
					t.Fatalf("Expected file to exist at %s or %s, but got errors", outputFile, outputPath)
				}
			}
		}
	}

	// Clean up resources after verification
	terraform.Destroy(t, ctx.Terraform)
}
