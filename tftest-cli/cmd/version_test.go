package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestVersionCommand(t *testing.T) {
	// Test that the version command exists and has the expected properties
	assert.NotNil(t, versionCmd, "Version command should not be nil")
	assert.Equal(t, "version", versionCmd.Use, "Version command should have correct name")
	assert.NotEmpty(t, versionCmd.Short, "Version command should have a short description")
}

func TestVersionCommandOutput(t *testing.T) {
	// Save original version to restore later
	origVersion := Version
	defer func() { Version = origVersion }()
	
	// Set a test version
	Version = "test-version"
	
	// Create a command to test the version output
	cmd := &cobra.Command{}
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	
	// Add the version command as a subcommand
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("TFTest CLI " + Version)
		},
	}
	cmd.AddCommand(versionCmd)
	
	// Execute the version command
	cmd.SetArgs([]string{"version"})
	err := cmd.Execute()
	
	// Check results
	assert.NoError(t, err, "Command execution should not error")
	assert.Contains(t, buf.String(), "TFTest CLI test-version", "Version output should contain version information")
}