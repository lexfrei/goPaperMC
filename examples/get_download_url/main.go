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
		fmt.Println("Usage: get_download_url <project_id> [version] [build]")
		os.Exit(1)
	}
	
	projectID := os.Args[1]
	
	// Create API client
	client := api.NewClient().WithTimeout(10 * time.Second)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	var url string
	var err error

	// Process arguments
	switch len(os.Args) {
	case 2: // Only project_id - get URL for latest version
		url, err = client.GetLatestVersionURL(ctx, projectID)
		if err != nil {
			fmt.Printf("Error: %v\n", errors.UnwrapAll(err))
			os.Exit(1)
		}
	
	case 3: // project_id and version - get URL for latest build of version
		version := os.Args[2]
		url, err = client.GetLatestBuildURL(ctx, projectID, version)
		if err != nil {
			fmt.Printf("Error: %v\n", errors.UnwrapAll(err))
			os.Exit(1)
		}
	
	case 4: // project_id, version and build - get URL for specific build
		version := os.Args[2]
		
		build, err := strconv.ParseInt(os.Args[3], 10, 32)
		if err != nil {
			fmt.Printf("Error parsing build number: %v\n", errors.UnwrapAll(err))
			os.Exit(1)
		}
		
		url, err = client.GetBuildURL(ctx, projectID, version, int32(build))
		if err != nil {
			fmt.Printf("Error: %v\n", errors.UnwrapAll(err))
			os.Exit(1)
		}
	}

	// Print only the URL
	fmt.Println(url)
}
