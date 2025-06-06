package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	loadAsdf()

	// Get current version
	currentVersion, err := os.ReadFile("VERSION")
	if err != nil {
		fmt.Println("Error reading VERSION file:", err)
		os.Exit(1)
	}
	versionStr := strings.TrimSpace(string(currentVersion))
	fmt.Println("Current version:", versionStr)

	parts := strings.Split(versionStr, ".")
	if len(parts) != 3 {
		fmt.Println("Invalid version format in VERSION file. Expected format: X.Y.Z")
		os.Exit(1)
	}

	major, _ := strconv.Atoi(parts[0])
	minor, _ := strconv.Atoi(parts[1])
	patch, _ := strconv.Atoi(parts[2])

	// Determine bump type
	bumpType := ""
	if len(os.Args) > 1 {
		bumpType = os.Args[1]
	} else {
		bumpType = determineBumpType()
	}
	fmt.Println("Bump type:", bumpType)

	// Apply bump
	switch bumpType {
	case "major":
		major++
		minor = 0
		patch = 0
	case "minor":
		minor++
		patch = 0
	case "patch":
		patch++
	default:
		fmt.Println("Invalid bump type:", bumpType)
		fmt.Println("Usage:", os.Args[0], "[major|minor|patch]")
		os.Exit(1)
	}

	newVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)
	tagVersion := fmt.Sprintf("v%s", newVersion)
	fmt.Println("Version bumped to", tagVersion)

	if checkTagExists(tagVersion) {
		fmt.Printf("Tag %s already exists. Please use a different version.\n", tagVersion)
		os.Exit(1)
	}

	// Write VERSION file
	err = os.WriteFile("VERSION", []byte(newVersion+"\n"), 0644)
	if err != nil {
		fmt.Println("Error writing VERSION file:", err)
		os.Exit(1)
	}

	updateMakefile(newVersion)

	runCommand("git", "add", "VERSION", "Makefile")
	runCommand("git", "commit", "-m", fmt.Sprintf("release: cut %s [skip ci]", tagVersion))
	runCommand("git", "tag", "-a", tagVersion, "-m", fmt.Sprintf("release: %s", tagVersion))

	fmt.Println("Changes committed and tagged as", tagVersion)
	fmt.Println("Run 'git push && git push --tags' to push changes to remote")
}

func loadAsdf() {
	_, err := exec.LookPath("asdf")
	if err == nil {
		return
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Warning: Could not determine home directory")
		return
	}
	asdfPath := fmt.Sprintf("%s/.asdf/asdf.sh", homeDir)
	if _, err := os.Stat(asdfPath); err == nil {
		asdfBin := fmt.Sprintf("%s/.asdf/bin", homeDir)
		asdfShims := fmt.Sprintf("%s/.asdf/shims", homeDir)
		path := os.Getenv("PATH")
		os.Setenv("PATH", fmt.Sprintf("%s:%s:%s", asdfBin, asdfShims, path))
	}
}

func runCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running command: %s %v\n", name, args)
		os.Exit(1)
	}
}

func checkTagExists(version string) bool {
	cmd := exec.Command("git", "tag", "-l", version)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) != ""
}

func determineBumpType() string {
	lastTag, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
	if err != nil {
		return determineFromCommits(getAllCommits())
	}
	return determineFromCommits(getCommitsSinceTag(strings.TrimSpace(string(lastTag))))
}

func getAllCommits() []string {
	output, err := exec.Command("git", "log", "--pretty=format:%s").Output()
	if err != nil {
		return []string{}
	}
	return strings.Split(strings.TrimSpace(string(output)), "\n")
}

func getCommitsSinceTag(tag string) []string {
	output, err := exec.Command("git", "log", fmt.Sprintf("%s..HEAD", tag), "--pretty=format:%s").Output()
	if err != nil {
		return []string{}
	}
	return strings.Split(strings.TrimSpace(string(output)), "\n")
}

func determineFromCommits(commits []string) string {
	breakingPatterns := []string{
		"^BREAKING CHANGE:", "^breaking!:", "!:", "feat!:", "fix!:",
		"refactor!:", "docs!:", "style!:", "test!:", "chore!:",
		"ci!:", "build!:", "perf!:",
	}
	for _, commit := range commits {
		for _, pattern := range breakingPatterns {
			if matched, _ := regexp.MatchString(pattern, commit); matched {
				return "major"
			}
		}
	}

	featurePatterns := []string{"^feat:", "^feature:"}
	for _, commit := range commits {
		for _, pattern := range featurePatterns {
			if matched, _ := regexp.MatchString(pattern, commit); matched {
				return "minor"
			}
		}
	}

	return "patch"
}

func updateMakefile(newVersion string) {
	content, err := os.ReadFile("Makefile")
	if err != nil {
		fmt.Println("Error reading Makefile:", err)
		os.Exit(1)
	}
	re := regexp.MustCompile(`Version=v?[0-9]+\.[0-9]+\.[0-9]+`)
	updated := re.ReplaceAllString(string(content), fmt.Sprintf("Version=v%s", newVersion))
	if err := os.WriteFile("Makefile", []byte(updated), 0644); err != nil {
		fmt.Println("Error writing Makefile:", err)
		os.Exit(1)
	}
}
