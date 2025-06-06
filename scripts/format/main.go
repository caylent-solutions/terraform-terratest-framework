package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// shouldIgnoreFile checks if a file should be ignored based on ignoredDirs
func shouldIgnoreFile(filePath string, ignoredDirs []string) bool {
	for _, dir := range ignoredDirs {
		if dir == "" {
			continue
		}
		if strings.HasPrefix(filePath, dir+"/") || filePath == dir {
			return true
		}
	}
	return false
}

func main() {
	ignoreDirs := flag.String("ignore", "", "Comma-separated list of directories to ignore")
	flag.Parse()

	loadAsdf()

	fmt.Println("Formatting Go code...")

	var ignoredDirList []string
	if *ignoreDirs != "" {
		for _, dir := range strings.Split(*ignoreDirs, ",") {
			ignoredDirList = append(ignoredDirList, strings.TrimSpace(dir))
		}

		if len(ignoredDirList) == 1 {
			fmt.Printf("⚠️  Ignoring directory during formatting: %s\n", ignoredDirList[0])
		} else if len(ignoredDirList) > 1 {
			fmt.Printf("⚠️  Ignoring directories during formatting: %s\n", strings.Join(ignoredDirList, ", "))
		}
	}

	fmt.Println("Checking which files need formatting...")
	checkCmd := exec.Command("gofmt", "-l", ".")
	output, err := checkCmd.Output()
	if err != nil {
		fmt.Printf("Error checking files: %v\n", err)
		os.Exit(1)
	}

	files := strings.TrimSpace(string(output))
	if files == "" {
		fmt.Println("✅ All files are already properly formatted")
		return
	}

	fmt.Println("Formatting the following files:")
	fileCount := 0
	for _, file := range strings.Split(files, "\n") {
		if shouldIgnoreFile(file, ignoredDirList) {
			continue
		}
		fmt.Printf("  - %s\n", file)
		fileCount++
	}

	if fileCount == 0 {
		fmt.Println("No files need formatting after applying ignore rules")
		return
	}

	fmt.Println("Running gofmt to fix formatting...")
	cmd := exec.Command("gofmt", "-w", ".")
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running gofmt: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✨ Format complete - fixed %d file(s)\n", fileCount)
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

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			os.Setenv(parts[0], parts[1])
		}
	}
}
