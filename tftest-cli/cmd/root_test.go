package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	// Test that the root command exists and has the expected properties
	assert.NotNil(t, rootCmd, "Root command should not be nil")
	assert.Equal(t, "tftest", rootCmd.Use, "Root command should have correct name")
	assert.NotEmpty(t, rootCmd.Short, "Root command should have a short description")
	assert.NotEmpty(t, rootCmd.Long, "Root command should have a long description")
}

func TestVersionFlag(t *testing.T) {
	// Create a command to test the version flag
	cmd := &cobra.Command{}
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	
	// Add version flag
	cmd.Flags().BoolP("version", "V", false, "Print version information")
	cmd.SetVersionTemplate("TFTest CLI {{.Version}}\n")
	cmd.Version = "test-version"
	
	// Execute with version flag
	cmd.SetArgs([]string{"--version"})
	err := cmd.Execute()
	
	// Check results
	assert.NoError(t, err, "Command execution should not error")
	assert.Contains(t, buf.String(), "TFTest CLI test-version", "Version output should contain version information")
}

func TestVerboseFlag(t *testing.T) {
	// Save original value to restore later
	originalVerboseLevel := verboseLevel
	defer func() { verboseLevel = originalVerboseLevel }()
	
	// Create a command to test the verbose flag
	cmd := &cobra.Command{}
	cmd.PersistentFlags().StringVarP(&verboseLevel, "verbose", "v", "", "Set verbosity level")
	
	// Execute with verbose flag
	cmd.SetArgs([]string{"--verbose", "DEBUG"})
	err := cmd.Parse([]string{"--verbose", "DEBUG"})
	
	// Check results
	assert.NoError(t, err, "Command parsing should not error")
	assert.Equal(t, "DEBUG", verboseLevel, "Verbose level should be set correctly")
}

func TestExecute(t *testing.T) {
	// This is a minimal test to ensure Execute doesn't panic
	// A more comprehensive test would mock cobra's execution
	assert.NotPanics(t, func() {
		// Create a temporary rootCmd that doesn't actually execute anything
		origRootCmd := rootCmd
		rootCmd = &cobra.Command{
			Use:   "test",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		defer func() { rootCmd = origRootCmd }()
		
		Execute()
	}, "Execute should not panic")
}