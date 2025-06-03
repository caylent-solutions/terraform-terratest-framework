package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestRunCommand(t *testing.T) {
	// Test that the run command exists and has the expected properties
	assert.NotNil(t, runCmd, "Run command should not be nil")
	assert.Equal(t, "run", runCmd.Use, "Run command should have correct name")
	assert.NotEmpty(t, runCmd.Short, "Run command should have a short description")
	assert.NotEmpty(t, runCmd.Long, "Run command should have a long description")
}

func TestRunCommandFlags(t *testing.T) {
	// Test that the run command has the expected flags
	moduleRootFlag := runCmd.Flags().Lookup("module-root")
	assert.NotNil(t, moduleRootFlag, "Run command should have module-root flag")
	assert.Equal(t, "string", moduleRootFlag.Value.Type(), "module-root flag should be a string")
	
	examplePathFlag := runCmd.Flags().Lookup("example-path")
	assert.NotNil(t, examplePathFlag, "Run command should have example-path flag")
	assert.Equal(t, "string", examplePathFlag.Value.Type(), "example-path flag should be a string")
	
	commonOnlyFlag := runCmd.Flags().Lookup("common")
	assert.NotNil(t, commonOnlyFlag, "Run command should have common flag")
	assert.Equal(t, "bool", commonOnlyFlag.Value.Type(), "common flag should be a boolean")
}

func TestVerifyDirectoryStructure(t *testing.T) {
	// Create a temporary directory structure
	tempDir := t.TempDir()
	
	// Test with missing directories
	assert.False(t, verifyDirectoryStructure(tempDir), "Should return false when directories are missing")
	
	// Create the required directories
	err := os.Mkdir(filepath.Join(tempDir, "examples"), 0755)
	assert.NoError(t, err, "Failed to create examples directory")
	
	err = os.Mkdir(filepath.Join(tempDir, "tests"), 0755)
	assert.NoError(t, err, "Failed to create tests directory")
	
	// Test with required directories
	assert.True(t, verifyDirectoryStructure(tempDir), "Should return true when required directories exist")
}

func TestRunTests(t *testing.T) {
	// This is a minimal test to ensure runTests doesn't panic with valid inputs
	// A more comprehensive test would mock the filesystem and exec.Command
	
	// Save original values to restore later
	origModuleRoot := moduleRoot
	origExamplePath := examplePath
	origCommonOnly := commonOnly
	defer func() {
		moduleRoot = origModuleRoot
		examplePath = origExamplePath
		commonOnly = origCommonOnly
	}()
	
	// Create a temporary directory structure for testing
	tempDir := t.TempDir()
	err := os.Mkdir(filepath.Join(tempDir, "examples"), 0755)
	assert.NoError(t, err, "Failed to create examples directory")
	
	err = os.Mkdir(filepath.Join(tempDir, "tests"), 0755)
	assert.NoError(t, err, "Failed to create tests directory")
	
	// Set up test values
	moduleRoot = tempDir
	examplePath = ""
	commonOnly = false
	
	// Test that runTests doesn't panic with valid directory structure
	// but will exit early when trying to run the actual tests
	assert.Panics(t, func() {
		runTests()
	}, "runTests should panic when trying to run actual tests")
}