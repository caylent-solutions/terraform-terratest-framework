package cmd

import (
	"fmt"
	"os"

	"github.com/caylent-solutions/terraform-terratest-framework/tftest-cli/logger"
	"github.com/spf13/cobra"
)

var (
	// Version will be set during build
	Version = "dev"

	// Verbose flag
	verboseLevel string

	// Root command
	rootCmd = &cobra.Command{
		Use:   "tftest",
		Short: "TFTest CLI - A command-line tool for testing Terraform modules",
		Long: `🚀 TFTest CLI 🚀
A command-line tool for testing Terraform modules using the Terraform Test Framework.

This framework is opinionated about directory structure and expects:
- Examples in the 'examples/' directory
- Tests in the 'tests/' directory with the same name as the example
- Common tests in 'tests/common/'
- Helper functions in 'tests/helpers/'

Run 'tftest run' to execute tests for your Terraform module.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Set log level based on verbose flag
			if verboseLevel != "" {
				level, err := logger.ParseLogLevel(verboseLevel)
				if err != nil {
					fmt.Printf("Warning: %v, using INFO level instead\n", err)
					level = logger.INFO
				}
				logger.SetDefaultLogLevel(level)
				logger.Debug("Log level set to %s", level.String())
			}
		},
	}
)

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error("Command failed: %v", err)
		os.Exit(1)
	}
}

func init() {
	// Add persistent flags that work across all subcommands
	rootCmd.PersistentFlags().StringVarP(&verboseLevel, "verbose", "v", "", "Set verbosity level (DEBUG, INFO, WARN, ERROR, FATAL)")

	// Add version flag
	rootCmd.Flags().BoolP("version", "V", false, "Print version information")
	rootCmd.SetVersionTemplate("TFTest CLI {{.Version}}\n")
	rootCmd.Version = Version
}
