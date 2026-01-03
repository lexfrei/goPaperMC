package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/cockroachdb/errors"
	"github.com/lexfrei/goPaperMC/pkg/api"
	"github.com/spf13/cobra"
)

// listCmd represents the list command.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List various resources from PaperMC API",
	Long: `The list command allows you to retrieve various
resources from the PaperMC API.`,
}

// listProjectsCmd represents the list projects command.
var listProjectsCmd = &cobra.Command{
	Use:     "projects",
	Aliases: []string{"project"},
	Short:   "List all available projects",
	Long:    `List all available projects from the PaperMC API.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient()
		if limit := GetLimit(); limit > 0 {
			client.WithLimit(limit)
		}

		ctx := context.Background()

		projects, err := client.GetProjects(ctx)
		if err != nil {
			fmt.Printf("Error: %v\n", errors.UnwrapAll(err))
			os.Exit(1)
		}

		for _, projectInfo := range projects.Projects {
			fmt.Printf("%s (%s)\n", projectInfo.Project.ID, projectInfo.Project.Name)
		}
	},
}

// listVersionsCmd represents the list versions command.
var listVersionsCmd = &cobra.Command{
	Use:     "versions PROJECT_ID",
	Aliases: []string{"version"},
	Short:   "List all versions for a project",
	Long:    `List all available versions for a specific project.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectID := args[0]

		client := api.NewClient()

		ctx := context.Background()

		projectInfo, err := client.GetProject(ctx, projectID)
		if err != nil {
			fmt.Printf("Error: %v\n", errors.UnwrapAll(err))
			os.Exit(1)
		}

		versions := projectInfo.FlattenVersions()

		// Apply limit if set
		limit := GetLimit()
		if limit > 0 && len(versions) > limit {
			start := len(versions) - limit
			versions = versions[start:]
		}

		// Sort versions in reverse order so newer ones appear at the top
		sort.Sort(sort.Reverse(sort.StringSlice(versions)))

		for _, version := range versions {
			fmt.Println(version)
		}
	},
}

// listBuildsCmd represents the list builds command.
var listBuildsCmd = &cobra.Command{
	Use:     "builds PROJECT_ID VERSION",
	Aliases: []string{"build"},
	Short:   "List all builds for a version",
	Long:    `List all available builds for a specific project version.`,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		projectID := args[0]
		version := args[1]

		client := api.NewClient()
		if limit := GetLimit(); limit > 0 {
			client.WithLimit(limit)
		}

		ctx := context.Background()

		builds, err := client.GetBuilds(ctx, projectID, version)
		if err != nil {
			fmt.Printf("Error: %v\n", errors.UnwrapAll(err))
			os.Exit(1)
		}

		for _, build := range builds {
			channel := ""
			if build.Channel != "" {
				channel = fmt.Sprintf(" (%s)", build.Channel)
			}

			fmt.Printf("%d%s\n", build.ID, channel)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listProjectsCmd)
	listCmd.AddCommand(listVersionsCmd)
	listCmd.AddCommand(listBuildsCmd)
}
