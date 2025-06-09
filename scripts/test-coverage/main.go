package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const coverageDir = "tmp/coverage"

// CoverageGroup represents a group of tests to run coverage on
type CoverageGroup struct {
	Name       string `json:"name"`
	Emoji      string `json:"emoji"`
	OutputFile string `json:"outputFile"`
	TestPath   string `json:"testPath"`
	CoverPkg   string `json:"coverPkg"`
}

func main() {
	// Get JSON file path from command line
	if len(os.Args) < 2 {
		fmt.Println("Error: JSON file path must be provided")
		fmt.Println("Usage: go run scripts/test-coverage/main.go path/to/coverage-groups.json")
		os.Exit(1)
	}

	jsonFilePath := os.Args[1]

	// Read JSON from file
	jsonData, err := os.ReadFile(jsonFilePath)
	if err != nil {
		fmt.Printf("Error reading coverage groups file: %v\n", err)
		os.Exit(1)
	}

	var groups []CoverageGroup
	if err := json.Unmarshal(jsonData, &groups); err != nil {
		fmt.Printf("Error parsing coverage groups: %v\n", err)
		os.Exit(1)
	}

	if len(groups) == 0 {
		fmt.Println("Error: No coverage groups provided in the JSON file")
		os.Exit(1)
	}

	// Create coverage directory if it doesn't exist
	if err := os.MkdirAll(coverageDir, 0755); err != nil {
		fmt.Printf("Error creating coverage directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🔍 Running tests with coverage...")

	// Run coverage for each group
	for _, group := range groups {
		fmt.Printf("\n%s %s:\n", group.Emoji, group.Name)
		runTestWithCoverage(group.OutputFile, group.TestPath, group.CoverPkg)
		printCoverageDetails(group.OutputFile)
	}

	// Merge coverage profiles
	fmt.Println("\n🔗 Merging coverage profiles...")
	mergeCoverageProfiles(groups)

	// Print summary
	fmt.Println("\n📊 Test Coverage Summary:")
	printCoverageSummary(groups)
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

func mergeCoverageProfiles(groups []CoverageGroup) {
	mergedPath := filepath.Join(coverageDir, "coverage.out")
	mergedFile, err := os.Create(mergedPath)
	if err != nil {
		fmt.Printf("Error creating merged coverage file: %v\n", err)
		return
	}
	defer mergedFile.Close()

	mergedFile.WriteString("mode: atomic\n")

	// Get output files from groups
	var coverageFiles []string
	for _, group := range groups {
		coverageFiles = append(coverageFiles, group.OutputFile)
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

func printCoverageSummary(groups []CoverageGroup) {
	// Add combined coverage to the groups
	coverageTypes := make([]struct {
		name       string
		emoji      string
		summaryLog string
	}, 0, len(groups)+1)

	for _, group := range groups {
		coverageTypes = append(coverageTypes, struct {
			name       string
			emoji      string
			summaryLog string
		}{
			name:       group.Name,
			emoji:      group.Emoji,
			summaryLog: strings.Replace(group.OutputFile, ".out", "-summary.log", 1),
		})
	}

	// Add combined coverage
	coverageTypes = append(coverageTypes, struct {
		name       string
		emoji      string
		summaryLog string
	}{
		name:       "Combined Total Coverage (All Components)",
		emoji:      "🧩",
		summaryLog: "coverage-summary.log",
	})

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
