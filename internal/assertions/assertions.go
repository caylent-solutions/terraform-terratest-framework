package assertions

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// AssertIdempotent verifies that a Terraform plan shows no changes after apply
func AssertIdempotent(t *testing.T, ctx interface{ GetTerraform() *terraform.Options }) {
	t.Helper()
	
	// Run terraform plan to check for changes
	planOutput := terraform.Plan(t, ctx.GetTerraform())
	
	// Check if the plan output contains "No changes"
	if !strings.Contains(planOutput, "No changes") && !strings.Contains(planOutput, "no changes") {
		t.Fatalf("Expected Terraform plan to show no changes, but got:\n%s", planOutput)
	}
}

// AssertFileExists checks if a file exists at the path specified by the output_file_path Terraform output
func AssertFileExists(t *testing.T, ctx interface{ GetOutput(t testing.TB, key string) string }) string {
	t.Helper()
	
	// Get the file path from the output
	filePath := ctx.GetOutput(t, "output_file_path")
	
	// Check if the file exists
	_, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Expected file to exist at %s, but got error: %v", filePath, err)
	}
	
	return filePath
}

// AssertFileContent checks if the output_content Terraform output matches the expected value
func AssertFileContent(t *testing.T, ctx interface{ GetOutput(t testing.TB, key string) string }) {
	t.Helper()
	
	// Get the file path from the output
	filePath := ctx.GetOutput(t, "output_file_path")
	
	// Get the expected content from the output
	expectedContent := ctx.GetOutput(t, "output_content")
	
	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file at %s: %v", filePath, err)
	}
	
	// Check if the content matches
	if string(content) != expectedContent {
		t.Fatalf("Expected file content to be %q, but got %q", expectedContent, string(content))
	}
}

// AssertOutputEquals checks if a specified Terraform output matches an expected value
func AssertOutputEquals(t *testing.T, ctx interface{ GetOutput(t testing.TB, key string) string }, outputName string, expectedValue string) {
	t.Helper()
	
	// Get the output value
	outputValue := ctx.GetOutput(t, outputName)
	
	// Check if the output matches
	if outputValue != expectedValue {
		t.Fatalf("Expected output %s to be %q, but got %q", outputName, expectedValue, outputValue)
	}
}

// AssertOutputContains checks if a specified Terraform output contains an expected substring
func AssertOutputContains(t *testing.T, ctx interface{ GetOutput(t testing.TB, key string) string }, outputName string, expectedSubstring string) {
	t.Helper()
	
	// Get the output value
	outputValue := ctx.GetOutput(t, outputName)
	
	// Check if the output contains the substring
	if !strings.Contains(outputValue, expectedSubstring) {
		t.Fatalf("Expected output %s to contain %q, but got %q", outputName, expectedSubstring, outputValue)
	}
}

// AssertOutputMatches checks if a specified Terraform output matches a regular expression
func AssertOutputMatches(t *testing.T, ctx interface{ GetOutput(t testing.TB, key string) string }, outputName string, pattern string) {
	t.Helper()
	
	// Get the output value
	outputValue := ctx.GetOutput(t, outputName)
	
	// Compile the regular expression
	regex, err := regexp.Compile(pattern)
	if err != nil {
		t.Fatalf("Invalid regular expression pattern %q: %v", pattern, err)
	}
	
	// Check if the output matches the pattern
	if !regex.MatchString(outputValue) {
		t.Fatalf("Expected output %s to match pattern %q, but got %q", outputName, pattern, outputValue)
	}
}

// AssertOutputNotEmpty checks if a specified Terraform output is not empty
func AssertOutputNotEmpty(t *testing.T, ctx interface{ GetOutput(t testing.TB, key string) string }, outputName string) {
	t.Helper()
	
	// Get the output value
	outputValue := ctx.GetOutput(t, outputName)
	
	// Check if the output is not empty
	if outputValue == "" {
		t.Fatalf("Expected output %s to not be empty", outputName)
	}
}

// AssertOutputEmpty checks if a specified Terraform output is empty
func AssertOutputEmpty(t *testing.T, ctx interface{ GetOutput(t testing.TB, key string) string }, outputName string) {
	t.Helper()
	
	// Get the output value
	outputValue := ctx.GetOutput(t, outputName)
	
	// Check if the output is empty
	if outputValue != "" {
		t.Fatalf("Expected output %s to be empty, but got %q", outputName, outputValue)
	}
}

// AssertOutputMapContainsKey checks if a Terraform map output contains a specific key
func AssertOutputMapContainsKey(t *testing.T, ctx interface{ GetTerraform() *terraform.Options }, outputName string, key string) {
	t.Helper()
	
	// Get the output map
	outputMap := terraform.OutputMap(t, ctx.GetTerraform(), outputName)
	
	// Check if the map contains the key
	if _, ok := outputMap[key]; !ok {
		t.Fatalf("Expected output map %s to contain key %q, but it doesn't", outputName, key)
	}
}

// AssertOutputMapKeyEquals checks if a key in a Terraform map output equals an expected value
func AssertOutputMapKeyEquals(t *testing.T, ctx interface{ GetTerraform() *terraform.Options }, outputName string, key string, expectedValue string) {
	t.Helper()
	
	// Get the output map
	outputMap := terraform.OutputMap(t, ctx.GetTerraform(), outputName)
	
	// Check if the map contains the key
	value, ok := outputMap[key]
	if !ok {
		t.Fatalf("Expected output map %s to contain key %q, but it doesn't", outputName, key)
	}
	
	// Check if the value matches
	if value != expectedValue {
		t.Fatalf("Expected output map %s key %q to be %q, but got %q", outputName, key, expectedValue, value)
	}
}

// AssertOutputListContains checks if a Terraform list output contains an expected value
func AssertOutputListContains(t *testing.T, ctx interface{ GetTerraform() *terraform.Options }, outputName string, expectedValue string) {
	t.Helper()
	
	// Get the output list
	outputList := terraform.OutputList(t, ctx.GetTerraform(), outputName)
	
	// Check if the list contains the value
	found := false
	for _, value := range outputList {
		if value == expectedValue {
			found = true
			break
		}
	}
	
	if !found {
		t.Fatalf("Expected output list %s to contain %q, but it doesn't", outputName, expectedValue)
	}
}

// AssertOutputListLength checks if a Terraform list output has the expected length
func AssertOutputListLength(t *testing.T, ctx interface{ GetTerraform() *terraform.Options }, outputName string, expectedLength int) {
	t.Helper()
	
	// Get the output list
	outputList := terraform.OutputList(t, ctx.GetTerraform(), outputName)
	
	// Check if the list has the expected length
	if len(outputList) != expectedLength {
		t.Fatalf("Expected output list %s to have length %d, but got %d", outputName, expectedLength, len(outputList))
	}
}

// AssertOutputJSONContains checks if a JSON string output contains an expected key-value pair
func AssertOutputJSONContains(t *testing.T, ctx interface{ GetOutput(t testing.TB, key string) string }, outputName string, key string, expectedValue interface{}) {
	t.Helper()
	
	// Get the output value
	outputValue := ctx.GetOutput(t, outputName)
	
	// Parse the JSON
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(outputValue), &jsonData)
	if err != nil {
		t.Fatalf("Failed to parse JSON from output %s: %v", outputName, err)
	}
	
	// Check if the JSON contains the key
	value, ok := jsonData[key]
	if !ok {
		t.Fatalf("Expected JSON output %s to contain key %q, but it doesn't", outputName, key)
	}
	
	// Check if the value matches
	assert.Equal(t, expectedValue, value, fmt.Sprintf("Expected JSON output %s key %q to be %v, but got %v", outputName, key, expectedValue, value))
}

// AssertResourceExists checks if a specific resource exists in the Terraform state
func AssertResourceExists(t *testing.T, ctx interface{ GetTerraform() *terraform.Options }, resourceType string, resourceName string) {
	t.Helper()
	
	// Run terraform state list to get all resources
	stateOutput := terraform.RunTerraformCommand(t, ctx.GetTerraform(), "state", "list")
	
	// Check if the resource exists
	resourceAddress := fmt.Sprintf("%s.%s", resourceType, resourceName)
	if !strings.Contains(stateOutput, resourceAddress) {
		t.Fatalf("Expected resource %s to exist in Terraform state, but it doesn't", resourceAddress)
	}
}

// AssertResourceCount checks if the number of resources of a specific type matches the expected count
func AssertResourceCount(t *testing.T, ctx interface{ GetTerraform() *terraform.Options }, resourceType string, expectedCount int) {
	t.Helper()
	
	// Run terraform state list to get all resources
	stateOutput := terraform.RunTerraformCommand(t, ctx.GetTerraform(), "state", "list")
	
	// Count the resources of the specified type
	count := 0
	for _, line := range strings.Split(stateOutput, "\n") {
		if strings.HasPrefix(line, resourceType+".") {
			count++
		}
	}
	
	// Check if the count matches
	if count != expectedCount {
		t.Fatalf("Expected %d resources of type %s, but found %d", expectedCount, resourceType, count)
	}
}

// AssertNoResourcesOfType checks that no resources of a specific type exist in the Terraform state
func AssertNoResourcesOfType(t *testing.T, ctx interface{ GetTerraform() *terraform.Options }, resourceType string) {
	t.Helper()
	
	// Run terraform state list to get all resources
	stateOutput := terraform.RunTerraformCommand(t, ctx.GetTerraform(), "state", "list")
	
	// Check if any resources of the specified type exist
	for _, line := range strings.Split(stateOutput, "\n") {
		if strings.HasPrefix(line, resourceType+".") {
			t.Fatalf("Expected no resources of type %s, but found %s", resourceType, line)
		}
	}
}

// AssertTerraformVersion checks if the Terraform version meets the minimum required version
func AssertTerraformVersion(t *testing.T, ctx interface{ GetTerraform() *terraform.Options }, minVersion string) {
	t.Helper()
	
	// Run terraform version to get the current version
	versionOutput := terraform.RunTerraformCommand(t, ctx.GetTerraform(), "version")
	
	// Extract the version number
	versionRegex := regexp.MustCompile(`Terraform v(\d+\.\d+\.\d+)`)
	matches := versionRegex.FindStringSubmatch(versionOutput)
	if len(matches) < 2 {
		t.Fatalf("Failed to extract Terraform version from output: %s", versionOutput)
	}
	
	currentVersion := matches[1]
	
	// Compare versions (simple string comparison works for semantic versions)
	if currentVersion < minVersion {
		t.Fatalf("Terraform version %s is less than required minimum version %s", currentVersion, minVersion)
	}
}