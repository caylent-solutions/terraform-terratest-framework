package main

import (
	"github.com/caylent-solutions/terraform-terratest-framework/tftest-cli/cmd"
)

// Version will be set during build
var Version = "dev"

func main() {
	// Set version
	cmd.Version = Version

	// Execute the root command
	cmd.Execute()
}
