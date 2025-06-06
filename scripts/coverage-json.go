package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

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

func main() {
	coverageDir := "tmp/coverage"
	if len(os.Args) > 1 {
		coverageDir = os.Args[1]
	}

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
