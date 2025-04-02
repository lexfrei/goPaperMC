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
	if len(os.Args) < 2 {
		fmt.Println("Usage: download_latest <destination_directory>")
		os.Exit(1)
	}
	
	destDir := os.Args[1]
	
	// Create API client
	client := api.NewClient().WithTimeout(30 * time.Second)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	
	// Download latest stable Paper version
	result, err := client.DownloadLatestStableVersion(ctx, "paper", destDir)
	if err != nil {
		fmt.Printf("Error: %v\n", errors.UnwrapAll(err))
		
		// Print full error details in verbose mode
		if len(os.Args) > 2 && os.Args[2] == "-v" {
			fmt.Printf("\nDetailed error:\n%+v\n", err)
		}
		
		os.Exit(1)
	}

	fmt.Printf("Downloaded %s\n", result.Filename)
	
	if !result.Valid {
		fmt.Printf("Checksum verification FAILED! Expected: %s, got: %s\n", 
			result.ExpectedSHA256, result.ActualSHA256)
		os.Exit(1)
	}
}
