package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/lexfrei/goPaperMC/pkg/api"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: version_groups <project_id> [limit]")
		os.Exit(1)
	}

	projectID := os.Args[1]

	// Parse limit if provided
	limit := 0
	if len(os.Args) > 2 {
		limit, _ = strconv.Atoi(os.Args[2])
	}

	// Create API client
	client := api.NewClient().WithTimeout(10 * time.Second)
	if limit > 0 {
		client.WithLimit(limit)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get project information
	projectInfo, err := client.GetProject(ctx, projectID)
	if err != nil {
		fmt.Printf("Error: %v\n", errors.UnwrapAll(err))
		os.Exit(1)
	}

	// Get sorted version groups
	groups := make([]string, 0, len(projectInfo.Versions))
	for group := range projectInfo.Versions {
		groups = append(groups, group)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(groups)))

	for _, groupName := range groups {
		fmt.Printf("%s\n", groupName)

		versions := projectInfo.Versions[groupName]

		// Display versions in the group
		count := 3
		if limit > 0 && limit < count {
			count = limit
		}

		if len(versions) < count {
			count = len(versions)
		}

		for i := len(versions) - count; i < len(versions); i++ {
			version := versions[i]

			// Get latest build for this version
			build, err := client.GetLatestBuildV3(ctx, projectID, version)
			if err != nil {
				fmt.Printf("  %s (Error: %v)\n", version, errors.UnwrapAll(err))
				continue
			}

			fmt.Printf("  %s build %d (%s)\n", version, build.ID, build.Channel)
		}
	}
}
