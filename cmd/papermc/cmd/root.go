package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	limit   int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "papermc",
	Short: "PaperMC CLI - Command line interface for PaperMC API",
	Long: `A command line tool to interact with the PaperMC API. 
It allows you to list projects, versions, and builds, as well as
download or get download URLs for PaperMC artifacts.`,
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.papermc.yaml)")
	rootCmd.PersistentFlags().IntVar(&limit, "limit", 0, "limit the number of items to show (0 means no limit)")

	// Bind the limit flag to viper
	viper.BindPFlag("limit", rootCmd.PersistentFlags().Lookup("limit"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".papermc" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".papermc")
	}

	viper.SetEnvPrefix("PAPERMC")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}

// GetLimit returns the limit set from flags or config
func GetLimit() int {
	return viper.GetInt("limit")
}
