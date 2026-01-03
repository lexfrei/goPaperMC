package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/lexfrei/goPaperMC/pkg/api"
	"github.com/spf13/cobra"
)

// BuildInfo contains information about a build for CI output.
type BuildInfo struct {
	Version string `json:"version"`
	Build   int32  `json:"build"`
	URL     string `json:"url"`
}

// ciCmd represents the ci command.
var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "Commands specifically for CI environments",
	Long:  `Commands designed to work well in Continuous Integration environments like GitHub Actions.`,
}

// ciMatrixCmd represents the ci matrix command.
var ciMatrixCmd = &cobra.Command{
	Use:   "matrix PROJECT_ID",
	Short: "Generate a JSON matrix for CI builds",
	Long: `Generate a JSON array of the latest builds for the last N versions of a project.
This is designed to be used in CI environments to generate a build matrix.

Example:
  papermc ci matrix paper --limit=3

This will output a JSON array of objects with version, build, and URL information
for the latest builds of the last 3 versions of the paper project.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectID := args[0]
		client := api.NewClient()

		// Apply channel filter if set
		if ch := GetChannel(); ch != "" {
			client.WithChannel(api.Channel(ch))
		}

		ctx := context.Background()

		// Get project info to get versions
		projectInfo, err := client.GetProject(ctx, projectID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting project info: %v\n", errors.UnwrapAll(err))
			os.Exit(1)
		}

		versions := projectInfo.FlattenVersions()
		limit := GetLimit()

		// Build array of builds, iterating from newest to oldest
		var buildInfos []BuildInfo
		for i := len(versions) - 1; i >= 0; i-- {
			if limit > 0 && len(buildInfos) >= limit {
				break
			}

			version := versions[i]
			build, err := client.GetLatestBuildV3(ctx, projectID, version)
			if err != nil {
				// Skip versions without matching builds
				continue
			}

			url := build.GetDownloadURL()
			if url == "" {
				continue
			}

			buildInfos = append(buildInfos, BuildInfo{
				Version: version,
				Build:   build.ID,
				URL:     url,
			})
		}

		// Reverse to get oldest-to-newest order
		for i, j := 0, len(buildInfos)-1; i < j; i, j = i+1, j-1 {
			buildInfos[i], buildInfos[j] = buildInfos[j], buildInfos[i]
		}

		// Output as JSON
		jsonOutput, err := json.Marshal(buildInfos)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating JSON: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(jsonOutput))
	},
}

// ciActionsCmd represents the ci github-actions command.
var ciActionsCmd = &cobra.Command{
	Use:   "github-actions PROJECT_ID",
	Short: "Output GitHub Actions compatible JSON matrix",
	Long: `Generate JSON specifically formatted for GitHub Actions matrix strategy.

Example:
  papermc ci github-actions paper --limit=3

This will output JSON that can be directly used in a GitHub Actions workflow:

  matrix=$(papermc ci github-actions paper --limit=3)
  echo "matrix=$matrix" >> $GITHUB_OUTPUT`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectID := args[0]
		client := api.NewClient()

		// Apply channel filter if set
		if ch := GetChannel(); ch != "" {
			client.WithChannel(api.Channel(ch))
		}

		ctx := context.Background()

		// Get project info to get versions
		projectInfo, err := client.GetProject(ctx, projectID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting project info: %v\n", errors.UnwrapAll(err))
			os.Exit(1)
		}

		versions := projectInfo.FlattenVersions()
		limit := GetLimit()

		// Build array of builds, iterating from newest to oldest
		var buildInfos []BuildInfo
		for i := len(versions) - 1; i >= 0; i-- {
			if limit > 0 && len(buildInfos) >= limit {
				break
			}

			version := versions[i]
			build, err := client.GetLatestBuildV3(ctx, projectID, version)
			if err != nil {
				// Skip versions without matching builds
				continue
			}

			url := build.GetDownloadURL()
			if url == "" {
				continue
			}

			buildInfos = append(buildInfos, BuildInfo{
				Version: version,
				Build:   build.ID,
				URL:     url,
			})
		}

		// Reverse to get oldest-to-newest order
		for i, j := 0, len(buildInfos)-1; i < j; i, j = i+1, j-1 {
			buildInfos[i], buildInfos[j] = buildInfos[j], buildInfos[i]
		}

		// Format in the way GitHub Actions expects
		matrixObj := map[string][]BuildInfo{
			"include": buildInfos,
		}

		// Output as JSON
		jsonOutput, err := json.Marshal(matrixObj)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating JSON: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(jsonOutput))
	},
}

var ciLatestCmd = &cobra.Command{
	Use:   "latest PROJECT_ID",
	Short: "Get the latest version",
	Long: `Get the latest version of a project.

Example:
  papermc ci latest paper

This will output just the latest version string, which can be used in scripts:

  latest_version=$(papermc ci latest paper)`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectID := args[0]
		client := api.NewClient()

		ctx := context.Background()

		// Get latest version
		version, err := client.GetLatestVersion(ctx, projectID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting latest version: %v\n", errors.UnwrapAll(err))
			os.Exit(1)
		}

		// Output just the version string
		fmt.Println(version)
	},
}

func init() {
	rootCmd.AddCommand(ciCmd)
	ciCmd.AddCommand(ciMatrixCmd)
	ciCmd.AddCommand(ciActionsCmd)
	ciCmd.AddCommand(ciLatestCmd)
}
