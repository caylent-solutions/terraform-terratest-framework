package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var ignoredDirs []string

func main() {
	ignoreFlag := flag.String("ignore", "", "Comma-separated list of directories to ignore during linting")
	flag.Parse()

	if *ignoreFlag != "" {
		ignoredDirs = strings.Split(*ignoreFlag, ",")
		for i, dir := range ignoredDirs {
			ignoredDirs[i] = strings.TrimSpace(dir)
		}
		if len(ignoredDirs) == 1 {
			fmt.Printf("⚠️  Ignoring directory during linting: %s\n", ignoredDirs[0])
		} else if len(ignoredDirs) > 1 {
			fmt.Printf("⚠️  Ignoring directories during linting: %s\n", strings.Join(ignoredDirs, ", "))
		}
	}

	loadAsdf()
	os.Setenv("GOGC", "off")

	exitCode := 0

	fmt.Println("Step 1: Running gofmt checks...")
	gofmtResult := runGofmtChecks()
	if gofmtResult != 0 {
		exitCode = 1
	}

	fmt.Println("Step 2: Running go vet checks...")
	goVetResult := runGoVetChecks()
	if goVetResult != 0 {
		exitCode = 1
	}

	fmt.Println("\n=== Lint Summary ===")
	fmt.Printf("gofmt checks: %s\n", formatStatus(gofmtResult == 0, "violates code formatting policy"))
	fmt.Printf("go vet checks: %s\n", formatStatus(goVetResult == 0, "violates code correctness policy"))

	os.Exit(exitCode)
}

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

func runGoVetChecks() int {
	govetExit := 0

	// Step 1: vet all Go packages (except scripts)
	pkgCmd := exec.Command("go", "list", "./...")
	pkgOutput, err := pkgCmd.Output()
	if err != nil {
		fmt.Println("Failed to list Go packages:", err)
		return 1
	}
	packages := strings.Split(strings.TrimSpace(string(pkgOutput)), "\n")

	for _, pkg := range packages {
		if strings.HasPrefix(pkg, "github.com/caylent-solutions/terraform-terratest-framework/scripts") {
			continue // skip scripts package
		}
		if shouldIgnoreFile(pkg) {
			continue
		}
		cmd := exec.Command("go", "vet", pkg)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("❌ %s (violates go vet policy)\n", pkg)
			printLines(output)
			govetExit = 1
		} else {
			fmt.Printf("✅ %s\n", pkg)
		}
	}

	// Step 2: vet individual .go files under scripts/
	err = filepath.WalkDir("scripts", func(path string, d os.DirEntry, err error) error {
		if d == nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		if shouldIgnoreFile(path) {
			return nil
		}

		cmd := exec.Command("go", "vet", path)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("❌ %s (violates go vet policy)\n", path)
			printLines(output)
			govetExit = 1
		} else {
			fmt.Printf("✅ %s\n", path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking scripts directory: %v\n", err)
		govetExit = 1
	}

	return govetExit
}

func shouldIgnoreFile(filePath string) bool {
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

func printLines(output []byte) {
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line != "" {
			fmt.Printf("   %s\n", line)
		}
	}
}

func formatStatus(success bool, failMessage string) string {
	if success {
		return "PASS ✅"
	}
	return fmt.Sprintf("FAIL ❌ (%s)", failMessage)
}

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

	for _, line := range strings.Split(string(output), "\n") {
		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			os.Setenv(parts[0], parts[1])
		}
	}
}
