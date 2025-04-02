package main

import (
	"context"
	"fmt"
	"os"
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
	
	for _, groupName := range projectInfo.VersionGroups {
		fmt.Printf("%s\n", groupName)
		
		// Get version group information
		group, err := client.GetVersionGroup(ctx, projectID, groupName)
		if err != nil {
			fmt.Printf("  Error: %v\n", errors.UnwrapAll(err))
			continue
		}
		
		// Get builds for the version group
		builds, err := client.GetVersionGroupBuilds(ctx, projectID, groupName)
		if err != nil {
			fmt.Printf("  Error: %v\n", errors.UnwrapAll(err))
			continue
		}
		
		// Display recent builds (limited by client or default 3)
		count := 3
		if limit > 0 && limit < count {
			count = limit
		}
		
		if len(builds.Builds) < count {
			count = len(builds.Builds)
		}
		
		for i := len(builds.Builds) - count; i < len(builds.Builds); i++ {
			build := builds.Builds[i]
			promoted := ""
			if build.Promoted {
				promoted = " (promoted)"
			}
			
			fmt.Printf("  %s %d%s\n", build.Version, build.Build, promoted)
		}
	}
}
