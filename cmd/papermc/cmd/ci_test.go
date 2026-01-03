package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
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
