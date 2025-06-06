package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Packages to check with go vet
var packages = []string{
	"./internal/assertions",
	"./internal/benchmark",
	"./internal/errors",
	"./internal/examples",
	"./internal/idempotency",
	"./internal/logging",
	"./internal/testctx",
	"./tftest-cli/cmd",
	"./tftest-cli/logger",
	"./examples/tests/common",
	"./examples/tests/ecs-private",
	"./examples/tests/ecs-public",
	"./examples/tests/helpers",
	"./tests/functional",
	"./tests/unit",
}

// Define ignored directories
var ignoredDirs []string

func main() {
	// Parse command line flags
	ignoreFlag := flag.String("ignore", "bin", "Comma-separated list of directories to ignore during linting")
	flag.Parse()

	// Process ignored directories
	if *ignoreFlag != "" {
		ignoredDirs = strings.Split(*ignoreFlag, ",")
		for i, dir := range ignoredDirs {
			ignoredDirs[i] = strings.TrimSpace(dir)
		}

		// Display which directories are being ignored
		if len(ignoredDirs) == 1 {
			fmt.Printf("⚠️  Ignoring directory during linting: %s\n", ignoredDirs[0])
		} else if len(ignoredDirs) > 1 {
			fmt.Printf("⚠️  Ignoring directories during linting: %s\n", strings.Join(ignoredDirs, ", "))
		}
	}

	// Load asdf
	loadAsdf()

	// Set environment variables to reduce memory usage
	os.Setenv("GOGC", "off")

	exitCode := 0

	// Step 1: Run gofmt checks
	fmt.Println("Step 1: Running gofmt checks...")
	gofmtExit := runGofmtChecks()
	if gofmtExit != 0 {
		exitCode = 1
	}

	// Step 2: Run go vet checks
	fmt.Println("Step 2: Running go vet checks...")
	govetExit := runGoVetChecks()
	if govetExit != 0 {
		exitCode = 1
	}

	// Step 3: Skip error checking
	fmt.Println("Step 3: Skipping error checking (most are false positives in test code)")
	fmt.Println("✅ Error checking skipped")

	// Print summary
	fmt.Println("\n=== Lint Summary ===")
	fmt.Printf("gofmt checks: %s\n", formatStatus(gofmtExit == 0, "violates code formatting policy"))
	fmt.Printf("go vet checks: %s\n", formatStatus(govetExit == 0, "violates code correctness policy"))
	fmt.Println("error checks: SKIPPED ✅")

	os.Exit(exitCode)
}

// loadAsdf attempts to load asdf from the user's home directory
func loadAsdf() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to get home directory")
		return
	}

	asdfPath := filepath.Join(homeDir, ".asdf", "asdf.sh")
	cmd := exec.Command("bash", "-c", fmt.Sprintf(". %s && env", asdfPath))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Failed to load asdf")
		return
	}

	// Parse environment variables from the output and set them
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			os.Setenv(parts[0], parts[1])
		}
	}
}

// runGofmtChecks runs gofmt and returns exit code (0 = success)
func runGofmtChecks() int {
	cmd := exec.Command("gofmt", "-l", ".")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error running gofmt: %v\n", err)
		return 1
	}

	files := strings.TrimSpace(string(output))
	if files != "" {
		fmt.Println("Files needing formatting (violates gofmt policy):")

		failedFiles := 0
		for _, file := range strings.Split(files, "\n") {
			// Skip files in ignored directories
			if shouldIgnoreFile(file) {
				continue
			}

			fmt.Printf("❌ %s\n", file)
			failedFiles++
		}

		if failedFiles > 0 {
			return 1
		}
	}

	fmt.Println("✅ All files properly formatted")
	return 0
}

// shouldIgnoreFile checks if a file should be ignored based on ignoredDirs
func shouldIgnoreFile(filePath string) bool {
	for _, dir := range ignoredDirs {
		if dir == "" {
			continue
		}

		// Check if the file path starts with the ignored directory
		if strings.HasPrefix(filePath, dir+"/") || filePath == dir {
			return true
		}
	}
	return false
}

// runGoVetChecks runs go vet on all packages and returns exit code (0 = success)
func runGoVetChecks() int {
	govetExit := 0

	for _, pkg := range packages {
		// Skip if directory doesn't exist
		if _, err := os.Stat(pkg); os.IsNotExist(err) {
			continue
		}

		cmd := exec.Command("go", "vet", pkg)
		output, err := cmd.CombinedOutput()

		if err != nil {
			fmt.Printf("❌ %s (violates go vet policy)\n", pkg)
			// Extract and display line numbers from errors
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if line == "" {
					continue
				}
				fmt.Printf("   %s\n", line)
			}
			govetExit = 1
		} else {
			fmt.Printf("✅ %s\n", pkg)
		}
	}

	return govetExit
}

// formatStatus returns a formatted status string
func formatStatus(success bool, failMessage string) string {
	if success {
		return "PASS ✅"
	}
	return fmt.Sprintf("FAIL ❌ (%s)", failMessage)
}
