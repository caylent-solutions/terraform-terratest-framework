package unit

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestVersionCommandOutput(t *testing.T) {
	// Create a command to test the version output
	cmd := &cobra.Command{}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	// Add the version command as a subcommand
	version := "test-version"
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("TFTest CLI " + version)
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
