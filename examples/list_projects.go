package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/lexfrei/goPaperMC/pkg/api"
)

func main() {
	// Parse command line arguments
	limit := 0
	if len(os.Args) > 1 {
		fmt.Sscanf(os.Args[1], "%d", &limit)
	}

	// Create API client
	client := api.NewClient().WithTimeout(10 * time.Second)
	if limit > 0 {
		client.WithLimit(limit)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
		
		if len(projectInfo.Versions) > 0 {
			count := 3
			if limit > 0 && limit < count {
				count = limit
			}
			
			start := 0
			if len(projectInfo.Versions) > count {
				start = len(projectInfo.Versions) - count
			}
			for _, version := range projectInfo.Versions[start:] {
				fmt.Printf("  %s\n", version)
			}
		}
	}
}
