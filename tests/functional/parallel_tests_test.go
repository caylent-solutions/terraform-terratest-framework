package functional

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestParallelTestsFlag(t *testing.T) {
	// Skip this test in CI environments
	if os.Getenv("CI") != "" {
		t.Skip("Skipping test in CI environment")
	}

	// Create a temporary directory for the test
	tempDir := t.TempDir()

	// Create a test file that uses the framework
	testFile := filepath.Join(tempDir, "framework_test.go")
	content := `
package framework

import (
	"testing"
	"os"
	"github.com/caylent-solutions/terraform-terratest-framework/pkg/testctx"
)

func TestFrameworkParallelism(t *testing.T) {
	// This test just verifies that the environment variable is correctly read
	if os.Getenv("TERRATEST_DISABLE_PARALLEL_TESTS") == "true" {
		if testctx.IsParallelTestsEnabled() {
			t.Error("Expected parallel tests to be disabled")
		}
	} else {
		if !testctx.IsParallelTestsEnabled() {
			t.Error("Expected parallel tests to be enabled")
		}
	}
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test with parallel tests enabled
	cmd := exec.Command("go", "test", "-v", tempDir)
	cmd.Env = append(os.Environ(), "TERRATEST_DISABLE_PARALLEL_TESTS=false")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run test with parallel tests enabled: %v\nOutput: %s", err, output)
	}

	// Test with parallel tests disabled
	cmd = exec.Command("go", "test", "-v", tempDir)
	cmd.Env = append(os.Environ(), "TERRATEST_DISABLE_PARALLEL_TESTS=true")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run test with parallel tests disabled: %v\nOutput: %s", err, output)
	}

	// Verify the output contains the expected test name
	outputStr := string(output)
	if !strings.Contains(outputStr, "TestFrameworkParallelism") {
		t.Errorf("Expected output to contain test name, got: %s", outputStr)
	}
}
