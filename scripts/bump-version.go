package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	// Load asdf if available
	loadAsdf()

	// Get the current version from VERSION file
	currentVersion, err := ioutil.ReadFile("VERSION")
	if err != nil {
		fmt.Println("Error reading VERSION file:", err)
		os.Exit(1)
	}

	versionStr := strings.TrimSpace(string(currentVersion))
	fmt.Println("Current version:", versionStr)

	// Parse the version
	parts := strings.Split(versionStr, ".")
	if len(parts) != 3 {
		fmt.Println("Invalid version format in VERSION file. Expected format: X.Y.Z")
		os.Exit(1)
	}

	major, _ := strconv.Atoi(parts[0])
	minor, _ := strconv.Atoi(parts[1])
	patch, _ := strconv.Atoi(parts[2])

	// Determine version bump type
	bumpType := ""
	if len(os.Args) > 1 {
		bumpType = os.Args[1]
	} else {
		bumpType = determineBumpType()
	}

	fmt.Println("Bump type:", bumpType)

	// Bump version
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

	// Create new version
	newVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)
	fmt.Println("New version:", newVersion)

	// Update VERSION file
	err = ioutil.WriteFile("VERSION", []byte(newVersion+"\n"), 0644)
	if err != nil {
		fmt.Println("Error writing VERSION file:", err)
		os.Exit(1)
	}

	// Update version in Makefile
	updateMakefile(newVersion)

	fmt.Println("Version bumped to", newVersion)

	// Stage the changes
	runCommand("git", "add", "VERSION", "Makefile")

	// Commit the changes
	runCommand("git", "commit", "-m", fmt.Sprintf("chore: bump version to %s", newVersion))

	// Create a tag
	runCommand("git", "tag", "-a", fmt.Sprintf("v%s", newVersion), "-m", fmt.Sprintf("Version %s", newVersion))

	fmt.Println("Changes committed and tagged as v" + newVersion)
	fmt.Println("Run 'git push && git push --tags' to push changes to remote")
}

func loadAsdf() {
	// Check if asdf is available in the path
	_, err := exec.LookPath("asdf")
	if err == nil {
		return // asdf is already in PATH
	}

	// Try to load asdf from common locations
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Warning: Could not determine home directory, asdf might not be loaded")
		return
	}

	asdfPath := fmt.Sprintf("%s/.asdf/asdf.sh", homeDir)
	_, err = os.Stat(asdfPath)
	if err == nil {
		fmt.Println("Loading asdf from", asdfPath)
		// We can't directly source the shell script in Go, but we can set the PATH
		// to include the asdf bin directory
		asdfBinPath := fmt.Sprintf("%s/.asdf/bin", homeDir)
		asdfShimsPath := fmt.Sprintf("%s/.asdf/shims", homeDir)
		
		path := os.Getenv("PATH")
		if !strings.Contains(path, asdfBinPath) {
			os.Setenv("PATH", fmt.Sprintf("%s:%s:%s", asdfBinPath, asdfShimsPath, path))
		}
	} else {
		fmt.Println("Warning: asdf not found at", asdfPath)
	}
}

func determineBumpType() string {
	// Get the last tag
	lastTag, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
	if err != nil {
		// If no tags exist, get all commits
		return determineFromCommits(getAllCommits())
	}

	// Get commits since last tag
	commits := getCommitsSinceTag(strings.TrimSpace(string(lastTag)))
	return determineFromCommits(commits)
}

func getAllCommits() []string {
	output, err := exec.Command("git", "log", "--pretty=format:%s").Output()
	if err != nil {
		fmt.Println("Error getting commit messages:", err)
		return []string{}
	}
	return strings.Split(strings.TrimSpace(string(output)), "\n")
}

func getCommitsSinceTag(tag string) []string {
	output, err := exec.Command("git", "log", fmt.Sprintf("%s..HEAD", tag), "--pretty=format:%s").Output()
	if err != nil {
		fmt.Println("Error getting commit messages since tag:", err)
		return []string{}
	}
	return strings.Split(strings.TrimSpace(string(output)), "\n")
}

func determineFromCommits(commits []string) string {
	// Check for breaking changes
	breakingChangePatterns := []string{
		"^BREAKING CHANGE:", "^breaking!:", "!:", "feat!:", "fix!:",
		"refactor!:", "docs!:", "style!:", "test!:", "chore!:",
		"ci!:", "build!:", "perf!:",
	}

	for _, commit := range commits {
		for _, pattern := range breakingChangePatterns {
			matched, _ := regexp.MatchString(pattern, commit)
			if matched {
				return "major"
			}
		}
	}

	// Check for feature commits
	featurePatterns := []string{"^feat:", "^feature:"}
	for _, commit := range commits {
		for _, pattern := range featurePatterns {
			matched, _ := regexp.MatchString(pattern, commit)
			if matched {
				return "minor"
			}
		}
	}

	// All other commit types result in a patch bump
	return "patch"
}

func updateMakefile(newVersion string) {
	// Read Makefile
	content, err := ioutil.ReadFile("Makefile")
	if err != nil {
		fmt.Println("Error reading Makefile:", err)
		os.Exit(1)
	}

	// Update version in Makefile
	re := regexp.MustCompile(`Version=v[0-9]*\.[0-9]*\.[0-9]*`)
	updatedContent := re.ReplaceAllString(string(content), fmt.Sprintf("Version=v%s", newVersion))

	// Write updated content back to Makefile
	err = ioutil.WriteFile("Makefile", []byte(updatedContent), 0644)
	if err != nil {
		fmt.Println("Error writing Makefile:", err)
		os.Exit(1)
	}
}

func runCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running command '%s %s': %v\n", name, strings.Join(args, " "), err)
		os.Exit(1)
	}
}