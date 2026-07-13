package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/lexfrei/goPaperMC/pkg/api"
)

func TestCIMatrixCommand(t *testing.T) {
	// Skip if not in CI environment or if it would make real API calls
	if os.Getenv("CI") == "" {
		t.Skip("Skipping test in non-CI environment")
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute command
	ciMatrixCmd.Run(ciMatrixCmd, []string{"paper"})

	// Restore stdout
	if err := w.Close(); err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("Failed to copy output: %v", err)
	}
	output := buf.String()

	// Check if output is valid JSON
	var buildInfos []BuildInfo
	err := json.Unmarshal([]byte(output), &buildInfos)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check if we have at least one build info
	if len(buildInfos) == 0 {
		t.Fatal("Expected at least one build info")
	}

	// Check if fields are populated
	firstBuild := buildInfos[0]
	if firstBuild.Version == "" {
		t.Fatal("Expected Version to be populated")
	}
	if firstBuild.Build == 0 {
		t.Fatal("Expected Build to be non-zero")
	}
	if firstBuild.URL == "" {
		t.Fatal("Expected URL to be populated")
	}
}

func TestCIGitHubActionsCommand(t *testing.T) {
	// Skip if not in CI environment or if it would make real API calls
	if os.Getenv("CI") == "" {
		t.Skip("Skipping test in non-CI environment")
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute command
	ciActionsCmd.Run(ciActionsCmd, []string{"paper"})

	// Restore stdout
	if err := w.Close(); err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("Failed to copy output: %v", err)
	}
	output := buf.String()

	// Check if output is valid JSON
	var matrixObj map[string][]BuildInfo
	err := json.Unmarshal([]byte(output), &matrixObj)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check if we have the include field
	include, ok := matrixObj["include"]
	if !ok {
		t.Fatal("Expected include field in matrix")
	}

	// Check if we have at least one build info
	if len(include) == 0 {
		t.Fatal("Expected at least one build info in include field")
	}
}

// TestCILatestCommand_ChannelFilter verifies that "ci latest" honors the
// --channel flag. Before the fix, ciLatestCmd never consulted GetChannel(),
// so it always returned the newest version overall even when that version
// had no build in the requested channel.
func TestCILatestCommand_ChannelFilter(t *testing.T) {
	// Skip if not in CI environment or if it would make real API calls
	if os.Getenv("CI") == "" {
		t.Skip("Skipping test in non-CI environment")
	}

	if err := rootCmd.PersistentFlags().Set("channel", "stable"); err != nil {
		t.Fatalf("Failed to set channel flag: %v", err)
	}
	defer func() {
		_ = rootCmd.PersistentFlags().Set("channel", "")
	}()

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute command
	ciLatestCmd.Run(ciLatestCmd, []string{"paper"})

	// Restore stdout
	if err := w.Close(); err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("Failed to copy output: %v", err)
	}
	version := strings.TrimSpace(buf.String())

	if version == "" {
		t.Fatal("Expected a version to be printed")
	}

	// The returned version must actually have a stable build. This is the
	// behavior that was previously broken by the ignored --channel flag.
	client := api.NewClient()
	builds, err := client.GetBuilds(context.Background(), "paper", version, api.ChannelStable)
	if err != nil {
		t.Fatalf("Failed to verify builds for version %q: %v", version, err)
	}
	if len(builds) == 0 {
		t.Fatalf("Expected version %q returned with --channel=stable to have a stable build", version)
	}
}
