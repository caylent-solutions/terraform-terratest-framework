package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const coverageDir = "tmp/coverage"

func main() {
	// Create coverage directory if it doesn't exist
	if err := os.MkdirAll(coverageDir, 0755); err != nil {
		fmt.Printf("Error creating coverage directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🔍 Generating HTML coverage report...")

	// Run tests with coverage for all packages
	outputPath := filepath.Join(coverageDir, "coverage.out")
	cmd := exec.Command("go", "test",
		"-covermode=atomic",
		fmt.Sprintf("-coverprofile=%s", outputPath),
		"-coverpkg=./internal/...,./pkg/...",
		"./internal/...", "./pkg/...", "./cmd/tftest/...", "./tests/functional/...", "./tests/unit/...")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: Test command returned error: %v\n", err)
	}

	// Generate HTML report
	htmlPath := filepath.Join(coverageDir, "coverage.html")
	htmlCmd := exec.Command("go", "tool", "cover",
		"-html", outputPath,
		"-o", htmlPath)

	htmlCmd.Stdout = os.Stdout
	htmlCmd.Stderr = os.Stderr

	if err := htmlCmd.Run(); err != nil {
		fmt.Printf("Error generating HTML report: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🌐 HTML coverage report generated at %s\n", htmlPath)
}
