package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/caylent-solutions/terraform-terratest-framework/cmd/tftest/logger"
	"github.com/spf13/cobra"
)

var (
	// Run command flags
	moduleRoot  string
	examplePath string
	commonOnly  bool
	parallel    bool
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run tests for a Terraform module",
	Long: `Run tests for a Terraform module using the Terraform Test Framework.

Examples:
  tftest run                     # Run all tests in the current directory
  tftest run --example-path vpc  # Run tests for the vpc example
  tftest run --common            # Run only common tests
  tftest run --module-root /path/to/terraform-module  # Run all tests in the specified module
  tftest run --parallel=false    # Run tests sequentially (disables parallel execution)

This command expects a specific directory structure:
- Examples in the 'examples/' directory
- Tests in the 'tests/' directory with the same name as the example
- Common tests in 'tests/common/'
- Helper functions in 'tests/helpers/'`,
	Run: func(cmd *cobra.Command, args []string) {
		runTests()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Add flags to run command
	runCmd.Flags().StringVar(&moduleRoot, "module-root", ".", "Path to the root of the Terraform module (runs all tests)")
	runCmd.Flags().StringVar(&examplePath, "example-path", "", "Specific example to test (leave empty to test all)")
	runCmd.Flags().BoolVar(&commonOnly, "common", false, "Run only common tests")
	runCmd.Flags().BoolVar(&parallel, "parallel", true, "Run tests in parallel (set to false to run sequentially)")
}

// runTests executes the tests based on the provided flags
func runTests() {
	// Get absolute path to module root
	absPath, err := filepath.Abs(moduleRoot)
	if err != nil {
		logger.Fatal("Error resolving path: %v", err)
	}

	// Verify directory structure
	if !verifyDirectoryStructure(absPath) {
		logger.Error("Invalid directory structure at %s", absPath)
		logger.Info("Expected structure:")
		logger.Info("  - examples/")
		logger.Info("  - tests/")
		logger.Info("  - tests/common/ (optional)")
		logger.Info("  - tests/helpers/ (optional)")
		os.Exit(1)
	}

	// If specific example, verify it exists
	if examplePath != "" {
		exampleDir := filepath.Join(absPath, "examples", examplePath)
		testDir := filepath.Join(absPath, "tests", examplePath)

		if _, err := os.Stat(exampleDir); os.IsNotExist(err) {
			logger.Fatal("Example directory not found: %s", exampleDir)
		}

		if _, err := os.Stat(testDir); os.IsNotExist(err) {
			logger.Fatal("Test directory for example not found: %s", testDir)
		}
	}

	// If common only, verify common directory exists
	if commonOnly {
		commonDir := filepath.Join(absPath, "tests", "common")
		if _, err := os.Stat(commonDir); os.IsNotExist(err) {
			logger.Fatal("Common test directory not found: %s", commonDir)
		}
	}

	// Build the test command
	testPath := "./tests/..."
	if examplePath != "" {
		testPath = fmt.Sprintf("./tests/%s/...", examplePath)
		logger.Info("Running tests for example: %s", examplePath)
	} else if commonOnly {
		testPath = "./tests/common/..."
		logger.Info("Running common tests")
	} else {
		logger.Info("Running all tests")
	}

	logger.Info("Module root: %s", absPath)
	if !parallel {
		logger.Info("Running tests sequentially (parallel execution disabled)")
	} else {
		logger.Info("Running tests in parallel")
	}
	logger.Info("Starting tests...")

	// Run the tests
	args := []string{"test", testPath, "-v"}

	// Add -p 1 flag if parallel is false to disable parallel execution
	if !parallel {
		args = append(args, "-p", "1")
	}

	cmd := exec.Command("go", args...)
	cmd.Dir = absPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		logger.Error("Tests failed: %v", err)
		os.Exit(1)
	}

	logger.Info("All tests passed! 🎉")
}

// verifyDirectoryStructure checks if the directory structure is as expected
func verifyDirectoryStructure(path string) bool {
	// Check if examples directory exists
	examplesPath := filepath.Join(path, "examples")
	if _, err := os.Stat(examplesPath); os.IsNotExist(err) {
		logger.Error("Examples directory not found at: %s", examplesPath)
		return false
	}

	// Check if tests directory exists
	testsPath := filepath.Join(path, "tests")
	if _, err := os.Stat(testsPath); os.IsNotExist(err) {
		logger.Error("Tests directory not found at: %s", testsPath)
		return false
	}

	return true
}
