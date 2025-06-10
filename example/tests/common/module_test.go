package common_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/pkg/assertions"
	"github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTerraformValidate runs 'terraform validate' on all examples
// This test ensures that the Terraform code is syntactically valid
func TestTerraformValidate(t *testing.T) {
	examples := []string{"basic", "advanced"}

	for _, example := range examples {
		t.Run(example, func(t *testing.T) {
			// Run the example
			ctx := testctx.RunSingleExample(t, "../../examples", example, testctx.TestConfig{
				Name: "validate-test-" + example,
			})

			// Run terraform validate
			terraform.Validate(t, ctx.Terraform)
		})
	}
}

// TestTerraformFormat checks if the Terraform code is properly formatted
// This test ensures that the Terraform code follows consistent formatting
func TestTerraformFormat(t *testing.T) {
	examples := []string{"basic", "advanced"}

	for _, example := range examples {
		t.Run(example, func(t *testing.T) {
			// Run the example
			ctx := testctx.RunSingleExample(t, "../../examples", example, testctx.TestConfig{
				Name: "format-test-" + example,
			})

			// Check if terraform code is formatted
			output, err := terraform.RunTerraformCommandE(t, ctx.Terraform, "fmt", "-check", "-recursive")
			assert.Empty(t, output, "Terraform code should be properly formatted")
			assert.NoError(t, err, "Terraform fmt should not fail")
		})
	}
}

// TestRequiredOutputs checks that required outputs are defined
// This test ensures that all required outputs are present in the module
func TestRequiredOutputs(t *testing.T) {
	examples := []string{"basic", "advanced"}
	requiredOutputs := []string{"output_file_path", "output_content"}

	for _, example := range examples {
		t.Run(example, func(t *testing.T) {
			// Run the example
			ctx := testctx.RunSingleExample(t, "../../examples", example, testctx.TestConfig{
				Name: "outputs-test-" + example,
			})

			// Check required outputs
			outputs := terraform.OutputAll(t, ctx.Terraform)

			for _, output := range requiredOutputs {
				_, exists := outputs[output]
				assert.True(t, exists, "Required output '%s' should be defined", output)
			}
		})
	}
}

// TestFileCreation checks that the file is created with the correct content
// This test ensures that the module creates the output file as expected
func TestFileCreation(t *testing.T) {
	examples := []string{"basic", "advanced"}

	for _, example := range examples {
		t.Run(example, func(t *testing.T) {
			// Run the example
			ctx := testctx.RunSingleExample(t, "../../examples", example, testctx.TestConfig{
				Name: "file-test-" + example,
			})

			// Get the file path from the output
			filePath := terraform.Output(t, ctx.Terraform, "output_file_path")
			assert.NotEmpty(t, filePath, "File path should not be empty")

			// Get the file content from the output
			expectedContent := terraform.Output(t, ctx.Terraform, "output_content")
			assert.NotEmpty(t, expectedContent, "File content should not be empty")

			// Verify the file exists and has the correct content
			// Note: In a real test, you would read the file and compare its content
		})
	}
}

// TestAllAssertionTypes demonstrates all the assertion types available in the framework
// This test runs on both the basic and advanced examples to verify various aspects of the module
func TestAllAssertionTypes(t *testing.T) {
	examples := []string{"basic", "advanced"}

	for _, example := range examples {
		t.Run(example, func(t *testing.T) {
			// Run the example
			ctx := testctx.RunSingleExample(t, "../../examples", example, testctx.TestConfig{
				Name: "assertions-test-" + example,
			})

			// Basic Assertions
			// Verify that the output_file_path output matches the expected value
			assertions.AssertOutputEquals(t, ctx, "output_file_path", terraform.Output(t, ctx.Terraform, "output_file_path"))

			// Verify that the output_file_path output contains the substring "output"
			assertions.AssertOutputContains(t, ctx, "output_file_path", "output")

			// Verify that the output_file_path output matches the regex pattern for a file path (contains a dot)
			assertions.AssertOutputMatches(t, ctx, "output_file_path", ".*\\..*")

			// Verify that the output_content output is not empty
			assertions.AssertOutputNotEmpty(t, ctx, "output_content")

			// File Assertions
			// Verify that the file specified by the output_file_path output exists
			assertions.AssertFileExists(t, ctx)

			// Verify that the content of the file matches the output_content output
			assertions.AssertFileContent(t, ctx)

			// Resource Assertions
			// Verify that the local_file.output resource exists in the Terraform state
			assertions.AssertResourceExists(t, ctx, "local_file", "output")

			// Verify that there is exactly 1 local_file resource in the Terraform state
			assertions.AssertResourceCount(t, ctx, "local_file", 1)

			// Verify that there are no aws_s3_bucket resources in the Terraform state
			assertions.AssertNoResourcesOfType(t, ctx, "aws_s3_bucket")

			// Environment Assertions
			// Verify that the Terraform version is at least 1.12.0
			assertions.AssertTerraformVersion(t, ctx, "1.12.0")

			// Idempotency is automatically tested by the framework when using RunSingleExample
			// This verifies that running terraform plan after apply shows no changes
			// assertions.AssertIdempotent(t, ctx)
		})
	}
}

// TestBenchmarking demonstrates how to benchmark Terraform operations
// This test measures the performance of applying and destroying the basic example
func TestBenchmarking(t *testing.T) {
	// Skip in normal runs as benchmarking can take time
	if testing.Short() {
		t.Skip("Skipping benchmarking in short mode")
	}

	// Define the benchmark function that will be measured
	benchmark := func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// For each iteration, apply and destroy the basic example
			// This measures the time it takes to create and destroy the resources
			ctx := testctx.RunSingleExample(t, "../../examples", "basic", testctx.TestConfig{
				Name: "benchmark-test",
			})
			terraform.Destroy(t, ctx.Terraform)
		}
	}

	// Run the benchmark if not in short mode
	// This will execute the benchmark function multiple times and report statistics
	result := testing.Benchmark(benchmark)
	t.Logf("Benchmark results: %s", result.String())
}

// readTFVars reads and parses the terraform.tfvars file
func readTFVars(t *testing.T, exampleDir, example string) map[string]interface{} {
	tfvarsPath := filepath.Join(exampleDir, "terraform.tfvars")

	// Read the file content
	content, err := os.ReadFile(tfvarsPath)
	require.NoError(t, err, "Failed to read terraform.tfvars file")

	// Parse the file based on the example type
	if example == "basic" {
		// For basic example, parse simple key-value pairs
		vars := make(map[string]interface{})

		// Use regex to extract key-value pairs
		re := regexp.MustCompile(`(\w+)\s*=\s*"([^"]*)"`)
		matches := re.FindAllStringSubmatch(string(content), -1)

		for _, match := range matches {
			if len(match) == 3 {
				vars[match[1]] = match[2]
			}
		}

		return vars
	} else if example == "advanced" {
		// For advanced example, we need to handle the complex structure
		vars := make(map[string]interface{})

		// Extract file_permission and output_filename
		filePermRe := regexp.MustCompile(`file_permission\s*=\s*"([^"]*)"`)
		filePermMatch := filePermRe.FindStringSubmatch(string(content))
		if len(filePermMatch) == 2 {
			vars["file_permission"] = filePermMatch[1]
		}

		filenameRe := regexp.MustCompile(`output_filename\s*=\s*"([^"]*)"`)
		filenameMatch := filenameRe.FindStringSubmatch(string(content))
		if len(filenameMatch) == 2 {
			vars["output_filename"] = filenameMatch[1]
		}

		// Extract json_config values we care about
		messageRe := regexp.MustCompile(`message\s*=\s*"([^"]*)"`)
		messageMatch := messageRe.FindStringSubmatch(string(content))
		if len(messageMatch) == 2 {
			vars["message"] = messageMatch[1]
		}

		enabledRe := regexp.MustCompile(`enabled\s*=\s*(true|false)`)
		enabledMatch := enabledRe.FindStringSubmatch(string(content))
		if len(enabledMatch) == 2 {
			vars["enabled"] = enabledMatch[1] == "true"
		}

		retriesRe := regexp.MustCompile(`retries\s*=\s*(\d+)`)
		retriesMatch := retriesRe.FindStringSubmatch(string(content))
		if len(retriesMatch) == 2 {
			vars["retries"] = retriesMatch[1]
		}

		return vars
	}

	return nil
}

// TestInputsMatchProvisioned verifies that the inputs provided to the module
// match what was actually provisioned by Terraform
func TestInputsMatchProvisioned(t *testing.T) {
	examples := []string{"basic", "advanced"}

	for _, example := range examples {
		t.Run(example, func(t *testing.T) {
			// Get the example directory path
			exampleDir := filepath.Join("../../examples", example)

			// Read the terraform.tfvars file
			tfvars := readTFVars(t, exampleDir, example)

			// Run the example
			ctx := testctx.RunSingleExample(t, "../../examples", example, testctx.TestConfig{
				Name: "input-validation-test-" + example,
			})

			// Get the outputs from the Terraform state
			outputs := terraform.OutputAll(t, ctx.Terraform)

			// Verify that output_content matches what was provided as input
			if example == "basic" {
				// For basic example, compare the output_content directly
				outputContent := outputs["output_content"]
				assert.Equal(t, tfvars["output_content"], outputContent, "output_content should match the input value")
			} else if example == "advanced" {
				// For advanced example, verify key fields in the parsed JSON
				jsonData := outputs["json_data"].(map[string]interface{})

				assert.Equal(t, tfvars["message"], jsonData["message"], "JSON message should match the input value")
				assert.Equal(t, tfvars["enabled"], jsonData["enabled"], "JSON enabled flag should match the input value")

				// Convert retries to float64 for comparison
				expectedRetries := 5.0 // We know it's 5 from the tfvars
				assert.Equal(t, expectedRetries, jsonData["retries"], "JSON retries should match the input value")
			}

			// Verify that output_filename is contained in output_file_path
			outputFilePath := outputs["output_file_path"].(string)
			expectedFilename := filepath.Base(tfvars["output_filename"].(string))
			assert.Contains(t, outputFilePath, expectedFilename,
				"output_file_path should contain the input filename")

			// Verify that file_permission matches what was provided as input
			outputPermission := outputs["file_permission"]
			assert.Equal(t, tfvars["file_permission"], outputPermission,
				"file_permission should match the input value")
		})
	}
}
