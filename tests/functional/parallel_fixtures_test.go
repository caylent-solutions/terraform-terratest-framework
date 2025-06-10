package functional

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestParallelFixturesFlag(t *testing.T) {
	// Skip this test in CI environments
	if os.Getenv("CI") != "" {
		t.Skip("Skipping test in CI environment")
	}

	// Create a temporary directory for the test
	tempDir := t.TempDir()

	// Create a simple test file
	testFile := filepath.Join(tempDir, "simple_test.go")
	content := `
package simple

import (
	"testing"
	"time"
)

func TestExample1(t *testing.T) {
	time.Sleep(100 * time.Millisecond)
}

func TestExample2(t *testing.T) {
	time.Sleep(100 * time.Millisecond)
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test with parallel fixtures enabled
	cmd := exec.Command("go", "test", "-v", tempDir)
	cmd.Env = append(os.Environ(), "TERRATEST_DISABLE_PARALLEL_TESTS=false")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run test with parallel fixtures enabled: %v\nOutput: %s", err, output)
	}

	// Test with parallel fixtures disabled
	cmd = exec.Command("go", "test", "-v", "-p", "1", tempDir)
	cmd.Env = append(os.Environ(), "TERRATEST_DISABLE_PARALLEL_TESTS=false")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run test with parallel fixtures disabled: %v\nOutput: %s", err, output)
	}

	// Verify the output contains the expected test names
	outputStr := string(output)
	if !strings.Contains(outputStr, "TestExample1") || !strings.Contains(outputStr, "TestExample2") {
		t.Errorf("Expected output to contain both test names, got: %s", outputStr)
	}
}
