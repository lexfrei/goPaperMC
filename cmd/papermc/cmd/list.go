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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List various resources from PaperMC API",
	Long: `The list command allows you to retrieve various 
resources from the PaperMC API.`,
}

// listProjectsCmd represents the list projects command
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

		// Create context
		ctx := context.Background()

		// Get list of projects
		projects, err := client.GetProjects(ctx)
		if err != nil {
			fmt.Printf("Error: %v\n", errors.UnwrapAll(err))
			os.Exit(1)
		}

		for _, project := range projects.Projects {
			projectInfo, err := client.GetProject(ctx, project)
			if err != nil {
				fmt.Printf("%s (Error: %v)\n", project, errors.UnwrapAll(err))
				continue
			}
			fmt.Printf("%s (%s)\n", project, projectInfo.ProjectName)
		}
	},
}

// listVersionsCmd represents the list versions command
var listVersionsCmd = &cobra.Command{
	Use:     "versions PROJECT_ID",
	Aliases: []string{"version"},
	Short:   "List all versions for a project",
	Long:    `List all available versions for a specific project.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectID := args[0]

		client := api.NewClient()
		if limit := GetLimit(); limit > 0 {
			client.WithLimit(limit)
		}

		// Create context
		ctx := context.Background()

		projectInfo, err := client.GetProject(ctx, projectID)
		if err != nil {
			fmt.Printf("Error: %v\n", errors.UnwrapAll(err))
			os.Exit(1)
		}

		// Sort versions in reverse order so newer ones appear at the top
		sort.Sort(sort.Reverse(sort.StringSlice(projectInfo.Versions)))

		for _, version := range projectInfo.Versions {
			fmt.Println(version)
		}
	},
}

// listBuildsCmd represents the list builds command
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

		// Create context
		ctx := context.Background()

		builds, err := client.GetBuilds(ctx, projectID, version)
		if err != nil {
			fmt.Printf("Error: %v\n", errors.UnwrapAll(err))
			os.Exit(1)
		}

		for _, build := range builds.Builds {
			promoted := ""
			if build.Promoted {
				promoted = " (promoted)"
			}
			
			fmt.Printf("%d%s\n", build.Build, promoted)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listProjectsCmd)
	listCmd.AddCommand(listVersionsCmd)
	listCmd.AddCommand(listBuildsCmd)
}
