package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/cockroachdb/errors"
	"github.com/lexfrei/goPaperMC/pkg/api"
	"github.com/spf13/cobra"
)

var destination string

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download PROJECT_ID [VERSION] [BUILD] [DESTINATION]",
	Short: "Download a build file",
	Long: `Download a build file from PaperMC API. 
If only PROJECT_ID is provided, the latest stable version and build will be downloaded.
If PROJECT_ID and VERSION are provided, the latest build for that version will be downloaded.
If PROJECT_ID, VERSION, and BUILD are provided, that specific build will be downloaded.
If DESTINATION is provided, the file will be saved to that location.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectID := args[0]
		
		var version string
		var buildNum int32
		var destDir string
		var err error

		// Set default values
		destDir = "."
		if destination != "" {
			destDir = destination
		}

		// Process arguments
		switch len(args) {
		case 1: // Only project_id
			// Download the latest stable version
			client := api.NewClient()
			
			version, err = client.GetRecommendedVersion(context.Background(), projectID)
			if err != nil {
				fmt.Printf("Error finding recommended version: %v\n", errors.UnwrapAll(err))
				os.Exit(1)
			}
			
			buildNum, err = client.GetLatestBuild(context.Background(), projectID, version)
			if err != nil {
				fmt.Printf("Error finding latest build: %v\n", errors.UnwrapAll(err))
				os.Exit(1)
			}
		
		case 2: // project_id and version
			version = args[1]
			
			client := api.NewClient()
			buildNum, err = client.GetLatestBuild(context.Background(), projectID, version)
			if err != nil {
				fmt.Printf("Error finding latest build: %v\n", errors.UnwrapAll(err))
				os.Exit(1)
			}
		
		case 3, 4: // project_id, version, build, and optionally destination
			version = args[1]
			
			build, err := strconv.ParseInt(args[2], 10, 32)
			if err != nil {
				fmt.Printf("Error parsing build number: %v\n", errors.UnwrapAll(err))
				os.Exit(1)
			}
			buildNum = int32(build)
			
			if len(args) >= 4 {
				destDir = args[3]
			}
		}

		client := api.NewClient()
		
		// Get the file name for download
		downloadName, err := client.GetDefaultDownloadName(context.Background(), projectID, version, buildNum)
		if err != nil {
			fmt.Printf("Error getting download name: %v\n", errors.UnwrapAll(err))
			os.Exit(1)
		}

		// Form the full path
		destPath := filepath.Join(destDir, downloadName)

		// Download the file
		result, err := client.DownloadFile(context.Background(), projectID, version, buildNum, downloadName, destPath)
		if err != nil {
			fmt.Printf("Error downloading file: %v\n", errors.UnwrapAll(err))
			os.Exit(1)
		}

		fmt.Printf("Downloaded %s\n", result.Filename)
		
		if !result.Valid {
			fmt.Printf("Checksum verification FAILED! Expected: %s, got: %s\n", 
				result.ExpectedSHA256, result.ActualSHA256)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(&destination, "destination", "d", "", "Destination directory for the downloaded file")
}
