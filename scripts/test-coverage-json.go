package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const coverageDir = "tmp/coverage"

// JSON output structures
type FileCoverage struct {
	File     string  `json:"file"`
	Coverage float64 `json:"coverage"`
}

type GroupCoverage struct {
	Name  string         `json:"name"`
	Total string         `json:"total"`
	Files []FileCoverage `json:"files"`
}

type CoverageSummary struct {
	Framework   GroupCoverage `json:"framework"`
	CLI         GroupCoverage `json:"cli"`
	Functional  GroupCoverage `json:"functional"`
	UnitHelpers GroupCoverage `json:"unit_helpers"`
	Combined    string        `json:"combined_total"`
}

func main() {
	// Create coverage directory if it doesn't exist
	if err := os.MkdirAll(coverageDir, 0755); err != nil {
		fmt.Printf("Error creating coverage directory: %v\n", err)
		os.Exit(1)
	}

	// Run tests with coverage
	runTestWithCoverage("coverage-framework.out", "./internal/...", "")
	runTestWithCoverage("coverage-pkg.out", "./pkg/...", "")
	runTestWithCoverage("coverage-cli.out", "./cmd/tftest/...", "")
	runTestWithCoverage("coverage-functional.out", "./tests/functional/...", "./pkg/...")
	runTestWithCoverage("coverage-unit-helpers.out", "./tests/unit/...", "")

	// Generate coverage summaries
	generateCoverageSummary("coverage-framework.out")
	generateCoverageSummary("coverage-pkg.out")
	generateCoverageSummary("coverage-cli.out")
	generateCoverageSummary("coverage-functional.out")
	generateCoverageSummary("coverage-unit-helpers.out")

	// Merge coverage profiles
	mergeCoverageProfiles()

	// Generate JSON coverage report
	generateJSONCoverage()
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

	// Redirect output to /dev/null
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		// Silently continue on errors
	}
}

func generateCoverageSummary(coverageFile string) {
	coveragePath := filepath.Join(coverageDir, coverageFile)
	summaryPath := filepath.Join(coverageDir, strings.Replace(coverageFile, ".out", "-summary.log", 1))

	cmd := exec.Command("go", "tool", "cover", "-func", coveragePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Silently continue on errors
		return
	}

	// Write to the summary file
	os.WriteFile(summaryPath, output, 0644)
}

func mergeCoverageProfiles() {
	// Create the merged coverage file with mode header
	mergedPath := filepath.Join(coverageDir, "coverage.out")
	mergedFile, err := os.Create(mergedPath)
	if err != nil {
		fmt.Printf("Error creating merged coverage file: %v\n", err)
		return
	}
	defer mergedFile.Close()

	// Write the mode line
	mergedFile.WriteString("mode: set\n")

	// List of coverage files to merge
	coverageFiles := []string{
		"coverage-framework.out",
		"coverage-pkg.out",
		"coverage-cli.out",
		"coverage-functional.out",
		"coverage-unit-helpers.out",
	}

	// Append all coverage data (skipping the mode line)
	for _, file := range coverageFiles {
		filePath := filepath.Join(coverageDir, file)
		data, err := os.ReadFile(filePath)
		if err != nil {
			// Skip files that don't exist
			continue
		}

		lines := strings.Split(string(data), "\n")
		for i, line := range lines {
			// Skip the mode line (first line) and empty lines
			if i > 0 && line != "" {
				mergedFile.WriteString(line + "\n")
			}
		}
	}

	// Generate the summary
	summaryPath := filepath.Join(coverageDir, "coverage-summary.log")
	cmd := exec.Command("go", "tool", "cover", "-func", mergedPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Silently continue on errors
		return
	}

	os.WriteFile(summaryPath, output, 0644)
}

func generateJSONCoverage() {
	frameworkSummaryPath := filepath.Join(coverageDir, "coverage-framework-summary.log")
	cliSummaryPath := filepath.Join(coverageDir, "coverage-cli-summary.log")
	functionalSummaryPath := filepath.Join(coverageDir, "coverage-functional-summary.log")
	unitHelpersSummaryPath := filepath.Join(coverageDir, "coverage-unit-helpers-summary.log")
	combinedSummaryPath := filepath.Join(coverageDir, "coverage-summary.log")

	frameworkFiles, frameworkTotal := parseCoverageOutput(frameworkSummaryPath)
	cliFiles, cliTotal := parseCoverageOutput(cliSummaryPath)
	functionalFiles, functionalTotal := parseCoverageOutput(functionalSummaryPath)
	unitHelpersFiles, unitHelpersTotal := parseCoverageOutput(unitHelpersSummaryPath)
	_, combinedTotal := parseCoverageOutput(combinedSummaryPath)

	summary := CoverageSummary{
		Framework: GroupCoverage{
			Name:  "framework",
			Total: frameworkTotal,
			Files: frameworkFiles,
		},
		CLI: GroupCoverage{
			Name:  "cli",
			Total: cliTotal,
			Files: cliFiles,
		},
		Functional: GroupCoverage{
			Name:  "functional",
			Total: functionalTotal,
			Files: functionalFiles,
		},
		UnitHelpers: GroupCoverage{
			Name:  "unit_helpers",
			Total: unitHelpersTotal,
			Files: unitHelpersFiles,
		},
		Combined: combinedTotal,
	}

	jsonData, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
}

func parseCoverageOutput(filePath string) ([]FileCoverage, string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading coverage file %s: %v\n", filePath, err)
		return nil, "0.0%"
	}

	lines := strings.Split(string(data), "\n")
	var files []FileCoverage
	var totalLine string

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		if strings.Contains(line, "total:") {
			totalLine = line
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 3 {
			file := parts[0]
			coverageStr := parts[len(parts)-1]

			// Skip mode line and non-coverage lines
			if file == "mode:" || !strings.HasSuffix(coverageStr, "%") {
				continue
			}

			files = append(files, FileCoverage{
				File:     file,
				Coverage: parsePercentage(coverageStr),
			})
		}
	}

	// Extract total coverage
	totalCoverage := "0.0%"
	if totalLine != "" {
		re := regexp.MustCompile(`(\d+\.\d+%)`)
		matches := re.FindStringSubmatch(totalLine)
		if len(matches) > 0 {
			totalCoverage = matches[0]
		}
	}

	return files, totalCoverage
}

func parsePercentage(s string) float64 {
	s = strings.TrimSuffix(s, "%")
	value := 0.0
	fmt.Sscanf(s, "%f", &value)
	return value
}
