package cmd

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/caylent-solutions/terraform-terratest-framework/tftest-cli/logger"
	"github.com/spf13/cobra"
)

var (
	// Format command flags
	formatModuleRoot  string
	formatExamplePath string
	formatCommonOnly  bool
	allFlag           bool
)

// formatCmd represents the format command
var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Format and verify Go test code",
	Long: `Format and verify Go test code in the tests directory.

This command formats Go test files using 'gofmt' and verifies the formatting.
If formatting issues are found, it will attempt to fix them and exit with a non-zero
status code if any issues couldn't be automatically fixed.

Examples:
  tftest format --all            # Format all Go test files, verifies each example has a matching test directory
  tftest format --example-path vpc    # Format only the vpc example test files, verifies the example exists
  tftest format --common         # Format only common test files, verifies the common directory exists

The command verifies the directory structure follows the expected pattern:
- With --all: Checks each example has a matching test directory
- With --example-path: Checks both the example and its test directory exist
- With --common: Checks the common test directory exists`,
	Run: func(cmd *cobra.Command, args []string) {
		formatTests()
	},
}

func init() {
	rootCmd.AddCommand(formatCmd)

	// Add flags to format command
	formatCmd.Flags().StringVar(&formatModuleRoot, "module-root", ".", "Path to the root of the Terraform module")
	formatCmd.Flags().StringVar(&formatExamplePath, "example-path", "", "Specific example test files to format")
	formatCmd.Flags().BoolVar(&formatCommonOnly, "common", false, "Format only common test files")
	formatCmd.Flags().BoolVarP(&allFlag, "all", "A", false, "Format all Go test files")
}

// formatTests formats and verifies the Go test files
func formatTests() {
	// Get absolute path to module root
	absPath, err := filepath.Abs(formatModuleRoot)
	if err != nil {
		logger.Fatal("Error resolving path: %v", err)
	}

	// Verify basic directory structure
	if !verifyDirectoryStructure(absPath) {
		logger.Error("Invalid directory structure at %s", absPath)
		logger.Info("Expected structure:")
		logger.Info("  - examples/")
		logger.Info("  - tests/")
		os.Exit(1)
	}

	// Determine which paths to format
	var paths []string
	hasErrors := false

	if allFlag {
		// Format all test directories
		examplesPath := filepath.Join(absPath, "examples")
		testsPath := filepath.Join(absPath, "tests")

		// Get all examples
		examples, err := os.ReadDir(examplesPath)
		if err != nil {
			logger.Fatal("Error reading examples directory: %v", err)
		}

		// Check that each example has a corresponding test directory
		for _, example := range examples {
			if !example.IsDir() {
				continue
			}

			exampleName := example.Name()
			exampleTestPath := filepath.Join(testsPath, exampleName)

			if _, err := os.Stat(exampleTestPath); os.IsNotExist(err) {
				logger.Error("Missing test directory for example %s: %s", exampleName, exampleTestPath)
				hasErrors = true
				continue
			}

			paths = append(paths, exampleTestPath)
		}

		// Add optional directories if they exist
		commonPath := filepath.Join(testsPath, "common")
		if _, err := os.Stat(commonPath); !os.IsNotExist(err) {
			paths = append(paths, commonPath)
		}

		helpersPath := filepath.Join(testsPath, "helpers")
		if _, err := os.Stat(helpersPath); !os.IsNotExist(err) {
			paths = append(paths, helpersPath)
		}

		logger.Info("Formatting all Go test files")
	} else if formatExamplePath != "" {
		// Verify both example and test directories exist
		exampleDir := filepath.Join(absPath, "examples", formatExamplePath)
		if _, err := os.Stat(exampleDir); os.IsNotExist(err) {
			logger.Fatal("Example directory not found: %s", exampleDir)
		}

		// Format a specific example's test directory
		exampleTestPath := filepath.Join(absPath, "tests", formatExamplePath)

		if _, err := os.Stat(exampleTestPath); os.IsNotExist(err) {
			logger.Fatal("Test directory for example %s not found: %s", formatExamplePath, exampleTestPath)
		}

		paths = append(paths, exampleTestPath)
		logger.Info("Formatting example test files: %s", formatExamplePath)
	} else if formatCommonOnly {
		// Format common test directory
		commonPath := filepath.Join(absPath, "tests", "common")

		if _, err := os.Stat(commonPath); os.IsNotExist(err) {
			logger.Fatal("Common test directory not found: %s", commonPath)
		}

		paths = append(paths, commonPath)
		logger.Info("Formatting common test files")
	} else {
		logger.Fatal("Please specify what to format: --all, --example-path, or --common")
	}

	if hasErrors {
		logger.Fatal("Directory structure validation failed")
	}

	logger.Info("Module root: %s", absPath)
	logger.Info("Starting Go code formatting...")

	formatErrors := false

	// Run gofmt on each path
	for _, path := range paths {
		logger.Info("Formatting Go files in: %s", path)

		// First check if there are any formatting issues
		checkCmd := exec.Command("gofmt", "-l", ".")
		checkCmd.Dir = path
		output, err := checkCmd.Output()

		if err != nil {
			logger.Error("Error checking format: %v", err)
			formatErrors = true
			continue
		}

		if len(output) > 0 {
			logger.Warn("Found formatting issues in:\n%s", string(output))

			// Try to fix the formatting issues
			fixCmd := exec.Command("gofmt", "-w", ".")
			fixCmd.Dir = path
			fixErr := fixCmd.Run()

			if fixErr != nil {
				logger.Error("Failed to fix formatting issues: %v", fixErr)
				formatErrors = true
			} else {
				logger.Info("Fixed formatting issues")
			}
		} else {
			logger.Info("No formatting issues found")
		}

		// Run go vet to check for code issues
		vetCmd := exec.Command("go", "vet", "./...")
		vetCmd.Dir = path
		vetOutput, vetErr := vetCmd.CombinedOutput()

		if vetErr != nil {
			logger.Error("Go vet found issues:\n%s", string(vetOutput))
			formatErrors = true
		} else {
			logger.Info("Go vet passed")
		}
	}

	if formatErrors {
		logger.Fatal("Format verification failed! Some issues could not be automatically fixed.")
	} else {
		logger.Info("Format verification passed! 🎉")
	}
}
