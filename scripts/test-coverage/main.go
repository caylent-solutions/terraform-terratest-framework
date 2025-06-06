package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const coverageDir = "tmp/coverage"

func main() {
	// Create coverage directory if it doesn't exist
	if err := os.MkdirAll(coverageDir, 0755); err != nil {
		fmt.Printf("Error creating coverage directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🔍 Running tests with coverage...")

	// Run framework coverage
	fmt.Println("\n🧱 Framework coverage:")
	runTestWithCoverage("coverage-framework.out", "./internal/...", "")
	printCoverageDetails("coverage-framework.out")

	// Run package coverage
	fmt.Println("\n📦 Package coverage:")
	runTestWithCoverage("coverage-pkg.out", "./pkg/...", "")
	printCoverageDetails("coverage-pkg.out")

	// Run CLI coverage
	fmt.Println("\n🧪 CLI coverage:")
	runTestWithCoverage("coverage-cli.out", "./cmd/tftest/...", "")
	printCoverageDetails("coverage-cli.out")

	// Run functional test coverage
	fmt.Println("\n🧪 Functional test coverage:")
	runTestWithCoverage("coverage-functional.out", "./tests/functional/...", "./pkg/...")
	printCoverageDetails("coverage-functional.out")

	// Run unit test helpers coverage
	fmt.Println("\n🧪 Unit test helpers coverage:")
	runTestWithCoverage("coverage-unit-helpers.out", "./tests/unit/...", "")
	printCoverageDetails("coverage-unit-helpers.out")

	// Merge coverage profiles
	fmt.Println("\n🔗 Merging coverage profiles...")
	mergeCoverageProfiles()

	// Print summary
	fmt.Println("\n📊 Test Coverage Summary:")
	printCoverageSummary()
}

func runTestWithCoverage(outputFile, testPackages, coverpkg string) {
	outputPath := filepath.Join(coverageDir, outputFile)

	var cmd *exec.Cmd
	if coverpkg != "" {
		cmd = exec.Command("go", "test", "-covermode=atomic",
			fmt.Sprintf("-coverprofile=%s", outputPath),
			fmt.Sprintf("-coverpkg=%s", coverpkg),
			testPackages)
	} else {
		cmd = exec.Command("go", "test", "-covermode=atomic",
			fmt.Sprintf("-coverprofile=%s", outputPath),
			testPackages)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: Test command returned error: %v\n", err)
	}
}

func printCoverageDetails(coverageFile string) {
	coveragePath := filepath.Join(coverageDir, coverageFile)
	summaryPath := filepath.Join(coverageDir, strings.Replace(coverageFile, ".out", "-summary.log", 1))

	summaryFile, err := os.Create(summaryPath)
	if err != nil {
		fmt.Printf("Error creating summary file: %v\n", err)
		return
	}
	defer summaryFile.Close()

	cmd := exec.Command("go", "tool", "cover", "-func", coveragePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Warning: Error getting coverage details: %v\n", err)
		return
	}

	fmt.Print(string(output))
	summaryFile.Write(output)
}

func mergeCoverageProfiles() {
	mergedPath := filepath.Join(coverageDir, "coverage.out")
	mergedFile, err := os.Create(mergedPath)
	if err != nil {
		fmt.Printf("Error creating merged coverage file: %v\n", err)
		return
	}
	defer mergedFile.Close()

	mergedFile.WriteString("mode: atomic\n")

	coverageFiles := []string{
		"coverage-framework.out",
		"coverage-pkg.out",
		"coverage-cli.out",
		"coverage-functional.out",
		"coverage-unit-helpers.out",
	}

	for _, file := range coverageFiles {
		filePath := filepath.Join(coverageDir, file)
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Warning: Could not read %s: %v\n", file, err)
			continue
		}

		lines := strings.Split(string(data), "\n")
		for i, line := range lines {
			if i > 0 && line != "" {
				mergedFile.WriteString(line + "\n")
			}
		}
	}

	summaryPath := filepath.Join(coverageDir, "coverage-summary.log")
	summaryFile, err := os.Create(summaryPath)
	if err != nil {
		fmt.Printf("Error creating summary file: %v\n", err)
		return
	}
	defer summaryFile.Close()

	cmd := exec.Command("go", "tool", "cover", "-func", mergedPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Warning: Error generating coverage summary: %v\n", err)
		return
	}

	summaryFile.Write(output)
}

func printCoverageSummary() {
	coverageTypes := []struct {
		name       string
		emoji      string
		summaryLog string
	}{
		{"Framework Total Coverage", "🧱", "coverage-framework-summary.log"},
		{"Package Total Coverage", "📦", "coverage-pkg-summary.log"},
		{"CLI Total Coverage", "🧪", "coverage-cli-summary.log"},
		{"Functional Test Coverage", "🧪", "coverage-functional-summary.log"},
		{"Unit Test Helpers Coverage", "🧪", "coverage-unit-helpers-summary.log"},
		{"Combined Total Coverage (All Components)", "🧩", "coverage-summary.log"},
	}

	for _, ct := range coverageTypes {
		fmt.Printf("%s %s:\n", ct.emoji, ct.name)

		summaryPath := filepath.Join(coverageDir, ct.summaryLog)
		data, err := os.ReadFile(summaryPath)
		if err != nil {
			fmt.Printf("  Warning: Could not read %s: %v\n", ct.summaryLog, err)
			continue
		}

		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.Contains(line, "total:") {
				parts := strings.Fields(line)
				if len(parts) >= 3 {
					fmt.Printf("  %s %s %s\n", parts[0], parts[1], parts[2])
					break
				}
			}
		}

		fmt.Println()
	}
}
