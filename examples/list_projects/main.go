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
		_, _ = fmt.Sscanf(os.Args[1], "%d", &limit)
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

	for _, projectInfo := range projects.Projects {
		fmt.Printf("%s (%s)\n", projectInfo.Project.ID, projectInfo.Project.Name)

		versions := make([]string, 0)
		for _, groupVersions := range projectInfo.Versions {
			versions = append(versions, groupVersions...)
		}

		if len(versions) > 0 {
			count := 3
			if limit > 0 && limit < count {
				count = limit
			}

			start := 0
			if len(versions) > count {
				start = len(versions) - count
			}
			for _, version := range versions[start:] {
				fmt.Printf("  %s\n", version)
			}
		}
	}
}
