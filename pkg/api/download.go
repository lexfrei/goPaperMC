package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
)

// DownloadResult contains the result of downloading a file.
type DownloadResult struct {
	Filename       string
	ExpectedSHA256 string
	ActualSHA256   string
	Valid          bool
}

// DownloadFile downloads a file from a build and verifies its hash.
func (c *Client) DownloadFile(ctx context.Context, projectID, version string, build int32, destPath string) (*DownloadResult, error) {
	// Get build information for hash verification
	buildInfo, err := c.GetBuild(ctx, projectID, version, build)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get build info")
	}

	downloadURL := buildInfo.GetDownloadURL()
	if downloadURL == "" {
		return nil, errors.Newf("no download URL found for build %d", build)
	}

	expectedSHA256 := buildInfo.GetDownloadSHA256()

	// Download the file
	reader, err := c.DownloadBuild(ctx, downloadURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download build")
	}
	defer func() { _ = reader.Close() }()

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return nil, errors.Wrap(err, "failed to create destination directory")
	}

	// Open file for writing
	file, err := os.Create(destPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create destination file")
	}
	defer func() { _ = file.Close() }()

	// Create hash function
	hasher := sha256.New()
	writer := io.MultiWriter(file, hasher)

	// Copy data
	if _, err := io.Copy(writer, reader); err != nil {
		return nil, errors.Wrap(err, "failed to copy data")
	}

	// Verify hash
	actualSHA256 := hex.EncodeToString(hasher.Sum(nil))
	valid := actualSHA256 == expectedSHA256

	result := &DownloadResult{
		Filename:       destPath,
		ExpectedSHA256: expectedSHA256,
		ActualSHA256:   actualSHA256,
		Valid:          valid,
	}

	if !valid && expectedSHA256 != "" {
		return result, errors.Newf("SHA256 mismatch: expected %s, got %s", expectedSHA256, actualSHA256)
	}

	return result, nil
}

// GetLatestBuild returns the number of the latest build for the specified version.
func (c *Client) GetLatestBuild(ctx context.Context, projectID, version string) (int32, error) {
	build, err := c.GetLatestBuildV3(ctx, projectID, version)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get latest build")
	}

	return build.ID, nil
}

// GetLatestVersion returns the latest available version for a project.
func (c *Client) GetLatestVersion(ctx context.Context, projectID string) (string, error) {
	projectInfo, err := c.GetProject(ctx, projectID)
	if err != nil {
		return "", errors.Wrap(err, "failed to get project info")
	}

	versions := projectInfo.FlattenVersions()
	if len(versions) == 0 {
		return "", errors.New("no versions found for this project")
	}

	// The latest version is the last one in the flattened list
	return versions[len(versions)-1], nil
}

// GetDefaultDownloadName returns the name of the main downloadable file for a build.
func (c *Client) GetDefaultDownloadName(ctx context.Context, projectID, version string, build int32) (string, error) {
	buildInfo, err := c.GetBuild(ctx, projectID, version, build)
	if err != nil {
		return "", errors.Wrap(err, "failed to get build info")
	}

	name := buildInfo.GetDownloadName()
	if name == "" {
		return "", errors.New("no downloads found for this build")
	}

	return name, nil
}
