package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const maxAsdfVersion = "v0.15.0"

func main() {
	updateOnly := false
	asdfVersion := maxAsdfVersion

	for _, arg := range os.Args[1:] {
		if arg == "--update" {
			updateOnly = true
		} else if strings.HasPrefix(arg, "--asdf-version=") {
			requestedVersion := strings.TrimPrefix(arg, "--asdf-version=")
			if compareVersions(requestedVersion, maxAsdfVersion) <= 0 {
				asdfVersion = requestedVersion
			} else {
				fmt.Printf("Warning: Requested asdf version %s is higher than maximum allowed %s. Using %s instead.\n",
					requestedVersion, maxAsdfVersion, maxAsdfVersion)
			}
		}
	}

	if updateOnly {
		updateTools()
		return
	}

	if os.Getenv("DEVCONTAINER") == "true" {
		fmt.Println("Caylent Devcontainer detected. Tools already installed.")
		fmt.Println("Running update-tools to ensure everything is up to date...")
		updateTools()
		return
	}

	if _, err := exec.LookPath("asdf"); err != nil {
		fmt.Printf("Installing asdf version %s...\n", asdfVersion)
		installAsdf(asdfVersion)
	} else {
		fmt.Println("asdf already installed.")
	}

	installPlugins()
}

func compareVersions(v1, v2 string) int {
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	for i := 0; i < len(parts1) && i < len(parts2); i++ {
		if parts1[i] < parts2[i] {
			return -1
		}
		if parts1[i] > parts2[i] {
			return 1
		}
	}

	if len(parts1) < len(parts2) {
		return -1
	}
	if len(parts1) > len(parts2) {
		return 1
	}
	return 0
}

func installAsdf(version string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}

	asdfDir := filepath.Join(homeDir, ".asdf")
	cmd := exec.Command("git", "clone", "https://github.com/asdf-vm/asdf.git", asdfDir, "--branch", version)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error cloning asdf repository:", err)
		os.Exit(1)
	}

	asdfBin := filepath.Join(asdfDir, "bin")
	asdfShims := filepath.Join(asdfDir, "shims")
	path := os.Getenv("PATH")
	os.Setenv("PATH", fmt.Sprintf("%s:%s:%s", asdfBin, asdfShims, path))
}

func installPlugins() {
	content, err := os.ReadFile(".tool-versions")
	if err != nil {
		fmt.Println("Error reading .tool-versions file:", err)
		os.Exit(1)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) > 0 {
			plugin := parts[0]
			fmt.Printf("Adding plugin: %s\n", plugin)
			cmd := exec.Command("asdf", "plugin", "add", plugin)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		}
	}

	fmt.Println("Installing tools from .tool-versions...")
	cmd := exec.Command("asdf", "install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error installing tools:", err)
		os.Exit(1)
	}

	cmd = exec.Command("asdf", "reshim")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error reshimming:", err)
		os.Exit(1)
	}
}

func updateTools() {
	if _, err := exec.LookPath("asdf"); err != nil {
		fmt.Println("asdf not found. Please run 'make install-tools' first.")
		os.Exit(1)
	}

	fmt.Println("Checking and updating asdf tools...")

	content, err := os.ReadFile(".tool-versions")
	if err != nil {
		fmt.Println("Error reading .tool-versions file:", err)
		os.Exit(1)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) > 0 {
			plugin := parts[0]
			fmt.Printf("Ensuring plugin %s is installed...\n", plugin)
			cmd := exec.Command("asdf", "plugin", "add", plugin)
			cmd.Run()
		}
	}

	fmt.Println("Installing/updating tools from .tool-versions...")
	cmd := exec.Command("asdf", "install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error installing tools:", err)
		os.Exit(1)
	}

	cmd = exec.Command("asdf", "reshim")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error reshimming:", err)
		os.Exit(1)
	}

	fmt.Println("All tools are up to date.")
}
