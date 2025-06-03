package functional

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCliVersion(t *testing.T) {
	// Skip if CLI is not built
	cliPath := filepath.Join("..", "..", "bin", "tftest")
	if _, err := os.Stat(cliPath); os.IsNotExist(err) {
		t.Skip("CLI not built, skipping test. Run 'make build-cli' first.")
	}

	// Run the CLI with version flag
	cmd := exec.Command(cliPath, "--version")
	output, err := cmd.CombinedOutput()
	
	// Check results
	assert.NoError(t, err, "CLI should execute without error")
	assert.Contains(t, string(output), "TFTest CLI", "Version output should contain CLI name")
}

func TestCliHelp(t *testing.T) {
	// Skip if CLI is not built
	cliPath := filepath.Join("..", "..", "bin", "tftest")
	if _, err := os.Stat(cliPath); os.IsNotExist(err) {
		t.Skip("CLI not built, skipping test. Run 'make build-cli' first.")
	}

	// Run the CLI with help flag
	cmd := exec.Command(cliPath, "--help")
	output, err := cmd.CombinedOutput()
	
	// Check results
	assert.NoError(t, err, "CLI should execute without error")
	assert.Contains(t, string(output), "Usage:", "Help output should contain usage information")
	assert.Contains(t, string(output), "Available Commands:", "Help output should list available commands")
	assert.Contains(t, string(output), "run", "Help output should mention run command")
	assert.Contains(t, string(output), "format", "Help output should mention format command")
}

func TestCliRunHelp(t *testing.T) {
	// Skip if CLI is not built
	cliPath := filepath.Join("..", "..", "bin", "tftest")
	if _, err := os.Stat(cliPath); os.IsNotExist(err) {
		t.Skip("CLI not built, skipping test. Run 'make build-cli' first.")
	}

	// Run the CLI with run help
	cmd := exec.Command(cliPath, "run", "--help")
	output, err := cmd.CombinedOutput()
	
	// Check results
	assert.NoError(t, err, "CLI should execute without error")
	assert.Contains(t, string(output), "Usage:", "Help output should contain usage information")
	assert.Contains(t, string(output), "--example-path", "Help output should mention example-path flag")
	assert.Contains(t, string(output), "--common", "Help output should mention common flag")
	assert.Contains(t, string(output), "--module-root", "Help output should mention module-root flag")
}

func TestCliFormatHelp(t *testing.T) {
	// Skip if CLI is not built
	cliPath := filepath.Join("..", "..", "bin", "tftest")
	if _, err := os.Stat(cliPath); os.IsNotExist(err) {
		t.Skip("CLI not built, skipping test. Run 'make build-cli' first.")
	}

	// Run the CLI with format help
	cmd := exec.Command(cliPath, "format", "--help")
	output, err := cmd.CombinedOutput()
	
	// Check results
	assert.NoError(t, err, "CLI should execute without error")
	assert.Contains(t, string(output), "Usage:", "Help output should contain usage information")
	assert.Contains(t, string(output), "--all", "Help output should mention all flag")
	assert.Contains(t, string(output), "--example-path", "Help output should mention example-path flag")
	assert.Contains(t, string(output), "--common", "Help output should mention common flag")
}

func TestCliVerboseFlag(t *testing.T) {
	// Skip if CLI is not built
	cliPath := filepath.Join("..", "..", "bin", "tftest")
	if _, err := os.Stat(cliPath); os.IsNotExist(err) {
		t.Skip("CLI not built, skipping test. Run 'make build-cli' first.")
	}

	// Run the CLI with verbose flag
	cmd := exec.Command(cliPath, "--verbose", "DEBUG", "version")
	output, err := cmd.CombinedOutput()
	
	// Check results
	assert.NoError(t, err, "CLI should execute without error")
	assert.Contains(t, string(output), "TFTest CLI", "Output should contain CLI name")
}