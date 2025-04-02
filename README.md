# goPaperMC

goPaperMC is a Go client for the [PaperMC API](https://api.papermc.io), which allows you to interact with the PaperMC API to retrieve information about projects, versions, builds, and download files.

## Features

- Get a list of available projects
- Get information about a project, its versions, and version groups
- Get a list of available builds for a specific version
- Get information about a specific build
- Download files from builds or get just the download URL
- Verify SHA256 hashes of downloaded files
- Get the latest and recommended versions and builds
- Limit query results to N latest items
- CLI utility with Cobra for powerful command handling
- Configuration with Viper for config files and environment variables
- Shell completions (bash, zsh, fish, powershell)
- Robust error handling with github.com/cockroachdb/errors
- Minimalistic output with information only when necessary

## Installation

```bash
go install github.com/lexfrei/goPaperMC/cmd/papermc@latest
```

Or clone and build:

```bash
git clone https://github.com/lexfrei/goPaperMC.git
cd goPaperMC
make build
```

## CLI Usage

```bash
# Get list of projects
papermc list projects

# Get the 5 latest versions for a project
papermc --limit=5 list versions paper

# Get list of builds for a version
papermc list builds paper 1.19.4

# Get download URL for the latest version without downloading
papermc get-url paper

# Get download URL for the latest build of a specific version
papermc get-url paper 1.19.4

# Get download URL for a specific build
papermc get-url paper 1.19.4 100

# Download latest version of a project
papermc download paper

# Download latest build of a specific version
papermc download paper 1.19.4

# Download specific build
papermc download paper 1.19.4 100

# Download to specific directory
papermc download paper 1.19.4 100 -d ./server

# Generate shell completions
papermc completion bash > ~/.bash_completion.d/papermc

# Get version information
papermc version
```

## Configuration

goPaperMC supports configuration through:

1. Command-line flags
2. Environment variables (prefixed with `PAPERMC_`)
3. Configuration files (`.papermc.yaml`, `.papermc.json`, etc.)

### Config File Example

```yaml
# ~/.papermc.yaml or ./.papermc.yaml
limit: 10
verbose: false
default_project: "paper"
```

### Environment Variables

```bash
# Set limit for all commands
export PAPERMC_LIMIT=5

# Set default project
export PAPERMC_DEFAULT_PROJECT=paper
```

## Library Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/lexfrei/goPaperMC/pkg/api"
)

func main() {
	// Create API client
	client := api.NewClient().
		WithTimeout(30 * time.Second).
		WithLimit(5) // Show only the 5 latest items

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Get list of projects
	projects, err := client.GetProjects(ctx)
	if err != nil {
		log.Fatalf("Error: %v", errors.UnwrapAll(err))
	}

	for _, project := range projects.Projects {
		fmt.Printf("%s\n", project)
	}

	// Get download URL for the latest version of Paper (without downloading)
	url, err := client.GetLatestVersionURL(ctx, "paper")
	if err != nil {
		log.Fatalf("Error: %v", errors.UnwrapAll(err))
	}
	fmt.Printf("Download URL: %s\n", url)

	// Download the latest stable Paper version
	result, err := client.DownloadLatestStableVersion(ctx, "paper", "./server")
	if err != nil {
		log.Fatalf("Error: %v", errors.UnwrapAll(err))
	}

	fmt.Printf("Downloaded %s\n", result.Filename)
	
	if !result.Valid {
		log.Fatalf("Checksum verification FAILED!")
	}
}
```

## API URL Methods

These methods allow getting download URLs without actually downloading the files:

```go
// Get URL for the latest version of a project
url, err := client.GetLatestVersionURL(ctx, "paper")

// Get URL for the latest build of a specific version
url, err := client.GetLatestBuildURL(ctx, "paper", "1.19.4")

// Get URL for a specific build
url, err := client.GetBuildURL(ctx, "paper", "1.19.4", 100)

// Get URL for the promoted (recommended) build of a version
url, err := client.GetPromotedBuildURL(ctx, "paper", "1.19.4")

// Format a download URL directly (if you already know all parameters)
url := client.FormatDownloadURL("paper", "1.19.4", 100, "paper-1.19.4-100.jar")
```

## License

BSD 3-Clause
