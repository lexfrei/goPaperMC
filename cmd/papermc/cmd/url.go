package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/cockroachdb/errors"
	"github.com/lexfrei/goPaperMC/pkg/api"
	"github.com/spf13/cobra"
)

// urlCmd represents the get-url command
var urlCmd = &cobra.Command{
	Use:   "get-url PROJECT_ID [VERSION] [BUILD]",
	Short: "Get download URL without downloading",
	Long: `Get the download URL for a build without actually downloading the file.
If only PROJECT_ID is provided, the URL for the latest stable version and build will be returned.
If PROJECT_ID and VERSION are provided, the URL for the latest build for that version will be returned.
If PROJECT_ID, VERSION, and BUILD are provided, the URL for that specific build will be returned.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectID := args[0]
		
		var url string
		var err error
		client := api.NewClient()

		// Create context
		ctx := context.Background()

		// Process arguments
		switch len(args) {
		case 1: // Only project_id - get URL for latest version
			url, err = client.GetLatestVersionURL(ctx, projectID)
			if err != nil {
				fmt.Printf("Error getting URL: %v\n", errors.UnwrapAll(err))
				os.Exit(1)
			}
		
		case 2: // project_id and version - get URL for latest build of version
			version := args[1]
			url, err = client.GetLatestBuildURL(ctx, projectID, version)
			if err != nil {
				fmt.Printf("Error getting URL: %v\n", errors.UnwrapAll(err))
				os.Exit(1)
			}
		
		case 3: // project_id, version and build - get URL for specific build
			version := args[1]
			
			build, err := strconv.ParseInt(args[2], 10, 32)
			if err != nil {
				fmt.Printf("Error parsing build number: %v\n", errors.UnwrapAll(err))
				os.Exit(1)
			}
			
			url, err = client.GetBuildURL(ctx, projectID, version, int32(build))
			if err != nil {
				fmt.Printf("Error getting URL: %v\n", errors.UnwrapAll(err))
				os.Exit(1)
			}
		}

		// Print only the URL without any additional text
		fmt.Println(url)
	},
	Aliases: []string{"url"},
}

func init() {
	rootCmd.AddCommand(urlCmd)
}
