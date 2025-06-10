package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const coverageDir = "tmp/coverage"

type TestEvent struct {
	Action  string `json:"Action"`
	Package string `json:"Package"`
	Test    string `json:"Test"`
}

func main() {
	// Get test path from command line
	if len(os.Args) < 2 {
		fmt.Println("Error: Test path argument is required")
		fmt.Println("Usage: go run scripts/functional-test/main.go <test-path>")
		os.Exit(1)
	}
	testPath := os.Args[1]

	// Read version from VERSION file
	version, err := os.ReadFile("VERSION")
	if err != nil {
		fmt.Printf("Error reading VERSION file: %v\n", err)
		os.Exit(1)
	}
	versionStr := fmt.Sprintf("v%s", strings.TrimSpace(string(version)))

	// Build CLI
	fmt.Println("Running functional tests (verbose output)...")
	buildCmd := exec.Command("go", "build", "-o", "bin/tftest",
		fmt.Sprintf("-ldflags=-X 'github.com/caylent-solutions/terraform-terratest-framework/cmd/tftest.Version=%s'", versionStr),
		"./cmd/tftest")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		fmt.Printf("Error building CLI: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("🎉 TFTest CLI built at bin/tftest")

	// Create coverage directory
	if err := os.MkdirAll(coverageDir, 0755); err != nil {
		fmt.Printf("Error creating coverage directory: %v\n", err)
		os.Exit(1)
	}

	// Run tests with coverage
	coverageFile := filepath.Join(coverageDir, "functional.out")
	testCmd := exec.Command("go", "test", "-v", "-covermode=atomic",
		fmt.Sprintf("-coverprofile=%s", coverageFile),
		"-coverpkg=./pkg/...", testPath)
	testCmd.Stdout = os.Stdout
	testCmd.Stderr = os.Stderr
	testCmd.Run() // Ignore error as tests might fail but we still want to generate coverage

	// Print coverage summary
	fmt.Println("\nFunctional Test Coverage of Packages:")
	coverCmd := exec.Command("go", "tool", "cover", "-func", coverageFile)
	coverOutput, err := coverCmd.CombinedOutput()
	if err == nil {
		fmt.Println(string(coverOutput))
	}

	// Run tests in JSON mode and process the output
	fmt.Println("\nSummarizing functional test results...")
	cmd := exec.Command("go", "test", "-json", testPath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error creating pipe: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting test command: %v\n", err)
		os.Exit(1)
	}

	// Process the JSON output
	passed := 0
	failed := 0
	skipped := 0

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		var event TestEvent
		err := json.Unmarshal(scanner.Bytes(), &event)
		if err != nil {
			continue
		}

		switch event.Action {
		case "pass":
			if event.Test != "" {
				passed++
			}
		case "fail":
			if event.Test != "" {
				failed++
			}
		case "skip":
			if event.Test != "" {
				skipped++
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		// Test failures are expected and shouldn't stop the summary
	}

	// Print summary
	fmt.Printf("📊 Functional Test Summary:\n")
	fmt.Printf("✅ Passed: %d\n", passed)
	fmt.Printf("❌ Failed: %d\n", failed)
	fmt.Printf("⚠️ Skipped: %d\n", skipped)
}
