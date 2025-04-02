package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information set at build time
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Print the version, commit, and build date information for the PaperMC CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("PaperMC CLI version %s\n", Version)
		fmt.Printf("Commit: %s\n", Commit)
		fmt.Printf("Built: %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
