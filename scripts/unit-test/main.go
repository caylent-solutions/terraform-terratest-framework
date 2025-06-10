package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type TestEvent struct {
	Action  string `json:"Action"`
	Package string `json:"Package"`
	Test    string `json:"Test"`
}

func main() {
	// Get test paths from command line arguments
	if len(os.Args) <= 1 {
		fmt.Println("Error: Test path argument is required")
		fmt.Println("Usage: go run scripts/unit-test/main.go \"./path/to/tests/...\"")
		os.Exit(1)
	}
	testPaths := strings.Split(os.Args[1], " ")

	// Clean Go test cache
	fmt.Println("Cleaning Go test cache...")
	cleanCmd := exec.Command("go", "clean", "-testcache")
	cleanCmd.Stdout = os.Stdout
	cleanCmd.Stderr = os.Stderr
	cleanCmd.Run()

	// Run unit tests with verbose output
	fmt.Println("Running unit tests (verbose output)...")
	args := append([]string{"test", "-v"}, testPaths...)
	testCmd := exec.Command("go", args...)
	testCmd.Stdout = os.Stdout
	testCmd.Stderr = os.Stderr
	testCmd.Run()

	// Run tests again in JSON mode and process the output directly
	fmt.Println("\nSummarizing unit test results...")
	jsonArgs := append([]string{"test", "-json"}, testPaths...)
	cmd := exec.Command("go", jsonArgs...)
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
	label := "Unit Test Summary"
	fmt.Printf("📊 %s:\n", label)
	fmt.Printf("✅ Passed: %d\n", passed)
	fmt.Printf("❌ Failed: %d\n", failed)
	fmt.Printf("⚠️ Skipped: %d\n", skipped)
}
