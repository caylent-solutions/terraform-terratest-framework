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
		fmt.Println("Usage: go run scripts/test-coverage-json/main.go path/to/coverage-groups.json")
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

	// Run tests for each group
	for _, group := range groups {
		runTestWithCoverage(group.OutputFile, group.TestPath, group.CoverPkg)
		generateCoverageSummary(group.OutputFile)
	}

	// Merge coverage profiles
	mergeCoverageProfiles(groups)

	// Generate JSON coverage
	generateJSONCoverage(groups)
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

	cmd.Stdout = nil
	cmd.Stderr = nil

	_ = cmd.Run() // intentionally ignore errors
}

func generateCoverageSummary(coverageFile string) {
	coveragePath := filepath.Join(coverageDir, coverageFile)
	summaryPath := filepath.Join(coverageDir, strings.Replace(coverageFile, ".out", "-summary.log", 1))

	cmd := exec.Command("go", "tool", "cover", "-func", coveragePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}

	_ = os.WriteFile(summaryPath, output, 0644)
}

func mergeCoverageProfiles(groups []CoverageGroup) {
	mergedPath := filepath.Join(coverageDir, "coverage.out")
	mergedFile, err := os.Create(mergedPath)
	if err != nil {
		fmt.Printf("Error creating merged coverage file: %v\n", err)
		return
	}
	defer mergedFile.Close()

	mergedFile.WriteString("mode: set\n")

	// Get output files from groups
	var coverageFiles []string
	for _, group := range groups {
		coverageFiles = append(coverageFiles, group.OutputFile)
	}

	for _, file := range coverageFiles {
		filePath := filepath.Join(coverageDir, file)
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		lines := strings.Split(string(data), "\n")
		for i, line := range lines {
			if i > 0 && line != "" {
				_, _ = mergedFile.WriteString(line + "\n")
			}
		}
	}

	summaryPath := filepath.Join(coverageDir, "coverage-summary.log")
	cmd := exec.Command("go", "tool", "cover", "-func", mergedPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}

	_ = os.WriteFile(summaryPath, output, 0644)
}

func generateJSONCoverage(groups []CoverageGroup) {
	// Create a map to store group coverage data
	groupCoverages := make(map[string]GroupCoverage)

	// Process each group
	for _, group := range groups {
		summaryPath := filepath.Join(coverageDir, strings.Replace(group.OutputFile, ".out", "-summary.log", 1))
		files, total := parseCoverageOutput(summaryPath)

		// Create a simplified name for the JSON key
		simpleName := strings.ToLower(strings.Split(group.Name, " ")[0])

		groupCoverages[simpleName] = GroupCoverage{
			Name:  group.Name,
			Total: total,
			Files: files,
		}
	}

	// Get combined coverage
	combinedSummaryPath := filepath.Join(coverageDir, "coverage-summary.log")
	_, combinedTotal := parseCoverageOutput(combinedSummaryPath)

	// Create the final JSON structure
	jsonMap := make(map[string]interface{})

	// Add all groups to the JSON
	for key, coverage := range groupCoverages {
		jsonMap[key] = coverage
	}

	// Add combined total
	jsonMap["combined_total"] = combinedTotal

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(jsonMap, "", "  ")
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

			if file == "mode:" || !strings.HasSuffix(coverageStr, "%") {
				continue
			}

			files = append(files, FileCoverage{
				File:     file,
				Coverage: parsePercentage(coverageStr),
			})
		}
	}

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

// parsePercentage converts a percentage string like "75.0%" to a float64 value
func parsePercentage(s string) float64 {
	s = strings.TrimSuffix(s, "%")
	value := 0.0
	fmt.Sscanf(s, "%f", &value)
	return value
}
