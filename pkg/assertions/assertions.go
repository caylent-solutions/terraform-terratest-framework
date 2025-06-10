package assertions

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// AssertOutputEquals checks if a Terraform output matches an expected value
func AssertOutputEquals(t testing.TB, ctx testctx.TestContext, outputName string, expectedValue interface{}) {
	output := terraform.Output(t, ctx.Terraform, outputName)
	assert.Equal(t, expectedValue, output, "Output %s should match expected value", outputName)
}

// AssertOutputContains checks if a Terraform output contains an expected substring
func AssertOutputContains(t testing.TB, ctx testctx.TestContext, outputName string, expectedSubstring string) {
	output := terraform.Output(t, ctx.Terraform, outputName)
	assert.Contains(t, output, expectedSubstring, "Output %s should contain expected substring", outputName)
}

// AssertOutputMatches checks if a Terraform output matches a regular expression
func AssertOutputMatches(t testing.TB, ctx testctx.TestContext, outputName string, regex string) {
	output := terraform.Output(t, ctx.Terraform, outputName)
	matched, err := regexp.MatchString(regex, output)
	assert.NoError(t, err, "Regex should be valid")
	assert.True(t, matched, "Output %s should match regex %s", outputName, regex)
}

// AssertOutputNotEmpty checks if a Terraform output is not empty
func AssertOutputNotEmpty(t testing.TB, ctx testctx.TestContext, outputName string) {
	output := terraform.Output(t, ctx.Terraform, outputName)
	assert.NotEmpty(t, output, "Output %s should not be empty", outputName)
}

// AssertOutputEmpty checks if a Terraform output is empty
func AssertOutputEmpty(t testing.TB, ctx testctx.TestContext, outputName string) {
	output := terraform.Output(t, ctx.Terraform, outputName)
	assert.Empty(t, output, "Output %s should be empty", outputName)
}

// AssertFileExists checks if a file exists at the path specified by the output_file_path Terraform output
func AssertFileExists(t testing.TB, ctx testctx.TestContext) {
	filePath := terraform.Output(t, ctx.Terraform, "output_file_path")
	fullPath := filepath.Join(ctx.Terraform.TerraformDir, filePath)
	_, err := os.Stat(fullPath)
	assert.NoError(t, err, "File should exist at path: %s", fullPath)
}

// AssertFileContent checks if the output_content Terraform output matches the expected value
func AssertFileContent(t testing.TB, ctx testctx.TestContext) {
	expectedContent := terraform.Output(t, ctx.Terraform, "output_content")
	filePath := terraform.Output(t, ctx.Terraform, "output_file_path")
	fullPath := filepath.Join(ctx.Terraform.TerraformDir, filePath)

	content, err := os.ReadFile(fullPath)
	assert.NoError(t, err, "Should be able to read file: %s", fullPath)
	assert.Equal(t, expectedContent, string(content), "File content should match expected value")
}

// AssertOutputMapContainsKey checks if a Terraform map output contains a specific key
func AssertOutputMapContainsKey(t testing.TB, ctx testctx.TestContext, outputName string, key string) {
	outputMap := terraform.OutputMap(t, ctx.Terraform, outputName)
	_, exists := outputMap[key]
	assert.True(t, exists, "Output map %s should contain key %s", outputName, key)
}

// AssertOutputMapKeyEquals checks if a key in a Terraform map output equals an expected value
func AssertOutputMapKeyEquals(t testing.TB, ctx testctx.TestContext, outputName string, key string, expectedValue interface{}) {
	outputMap := terraform.OutputMap(t, ctx.Terraform, outputName)
	value, exists := outputMap[key]
	assert.True(t, exists, "Output map %s should contain key %s", outputName, key)
	assert.Equal(t, expectedValue, value, "Output map %s key %s should equal expected value", outputName, key)
}

// AssertOutputListContains checks if a Terraform list output contains an expected value
func AssertOutputListContains(t testing.TB, ctx testctx.TestContext, outputName string, expectedValue string) {
	outputList := terraform.OutputList(t, ctx.Terraform, outputName)
	assert.Contains(t, outputList, expectedValue, "Output list %s should contain %s", outputName, expectedValue)
}

// AssertOutputListLength checks if a Terraform list output has the expected length
func AssertOutputListLength(t testing.TB, ctx testctx.TestContext, outputName string, expectedLength int) {
	outputList := terraform.OutputList(t, ctx.Terraform, outputName)
	assert.Len(t, outputList, expectedLength, "Output list %s should have length %d", outputName, expectedLength)
}

// AssertOutputJSONContains checks if a JSON string output contains an expected key-value pair
func AssertOutputJSONContains(t testing.TB, ctx testctx.TestContext, outputName string, key string, expectedValue interface{}) {
	jsonString := terraform.Output(t, ctx.Terraform, outputName)
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &jsonData)
	assert.NoError(t, err, "Output %s should be valid JSON", outputName)

	value, exists := jsonData[key]
	assert.True(t, exists, "JSON output %s should contain key %s", outputName, key)
	assert.Equal(t, expectedValue, value, "JSON output %s key %s should equal expected value", outputName, key)
}

// AssertResourceExists checks if a specific resource exists in the Terraform state
func AssertResourceExists(t testing.TB, ctx testctx.TestContext, resourceType string, resourceName string) {
	output, err := terraform.RunTerraformCommandE(t, ctx.Terraform, "state", "list")
	assert.NoError(t, err, "Terraform state list should not fail")

	resourceAddress := fmt.Sprintf("module.example.%s.%s", resourceType, resourceName)
	assert.Contains(t, output, resourceAddress, "Resource %s should exist in Terraform state", resourceAddress)
}

// AssertResourceCount checks if the number of resources of a specific type matches the expected count
func AssertResourceCount(t testing.TB, ctx testctx.TestContext, resourceType string, expectedCount int) {
	output, err := terraform.RunTerraformCommandE(t, ctx.Terraform, "state", "list")
	assert.NoError(t, err, "Terraform state list should not fail")

	count := 0
	for _, line := range regexp.MustCompile(`\r?\n`).Split(output, -1) {
		if regexp.MustCompile(fmt.Sprintf(`module\.example\.%s\.`, resourceType)).MatchString(line) {
			count++
		}
	}

	assert.Equal(t, expectedCount, count, "Resource count for %s should match expected count", resourceType)
}

// AssertNoResourcesOfType checks that no resources of a specific type exist in the Terraform state
func AssertNoResourcesOfType(t testing.TB, ctx testctx.TestContext, resourceType string) {
	AssertResourceCount(t, ctx, resourceType, 0)
}

// AssertTerraformVersion checks if the Terraform version meets the minimum required version
func AssertTerraformVersion(t testing.TB, ctx testctx.TestContext, minVersion string) {
	output, err := terraform.RunTerraformCommandE(t, ctx.Terraform, "version")
	assert.NoError(t, err, "Terraform version should not fail")

	// Extract version from output (e.g., "Terraform v1.12.1")
	versionRegex := regexp.MustCompile(`Terraform v(\d+\.\d+\.\d+)`)
	matches := versionRegex.FindStringSubmatch(output)
	assert.True(t, len(matches) > 1, "Should be able to extract Terraform version")

	// Compare versions (simplified, assumes semantic versioning)
	assert.True(t, matches[1] >= minVersion, "Terraform version should be at least %s", minVersion)
}

// AssertIdempotent verifies that a Terraform plan shows no changes after apply
func AssertIdempotent(t testing.TB, ctx testctx.TestContext) {
	planOutput := terraform.Plan(t, ctx.Terraform)
	assert.True(t,
		regexp.MustCompile(`No changes|no changes`).MatchString(planOutput),
		"Terraform plan should show no changes after apply")
}
