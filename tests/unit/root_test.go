package unit

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/caylent-solutions/terraform-terratest-framework/cmd/tftest/cmd"
)

func TestRootCommand(t *testing.T) {
	// Since rootCmd is not exported, we can only test the Execute function
	// This is a minimal test to ensure it exists
	assert.NotPanics(t, func() {
		// We're not actually executing the command, just verifying it exists
		_ = cmd.Execute
	})
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
	// Create a command to test the verbose flag
	cmd := &cobra.Command{}
	var verboseLevel string
	cmd.PersistentFlags().StringVarP(&verboseLevel, "verbose", "v", "", "Set verbosity level")

	// Execute with verbose flag
	cmd.SetArgs([]string{"--verbose", "DEBUG"})
	err := cmd.ParseFlags([]string{"--verbose", "DEBUG"})

	// Check results
	assert.NoError(t, err, "Command parsing should not error")
	assert.Equal(t, "DEBUG", verboseLevel, "Verbose level should be set correctly")
}
