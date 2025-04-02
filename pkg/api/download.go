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

// DownloadResult contains the result of downloading a file
type DownloadResult struct {
	Filename       string
	ExpectedSHA256 string
	ActualSHA256   string
	Valid          bool
}

// DownloadFile downloads a file from a build and verifies its hash
func (c *Client) DownloadFile(ctx context.Context, projectID, version string, build int32, downloadName, destPath string) (*DownloadResult, error) {
	// Get build information for hash verification
	buildInfo, err := c.GetBuild(ctx, projectID, version, build)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get build info")
	}

	download, ok := buildInfo.Downloads[downloadName]
	if !ok {
		return nil, errors.Newf("download %s not found in build %d", downloadName, build)
	}

	// Download the file
	reader, err := c.DownloadBuild(ctx, projectID, version, build, downloadName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download build")
	}
	defer reader.Close()

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return nil, errors.Wrap(err, "failed to create destination directory")
	}

	// Open file for writing
	file, err := os.Create(destPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create destination file")
	}
	defer file.Close()

	// Create hash function
	hasher := sha256.New()
	writer := io.MultiWriter(file, hasher)

	// Copy data
	if _, err := io.Copy(writer, reader); err != nil {
		return nil, errors.Wrap(err, "failed to copy data")
	}

	// Verify hash
	actualSHA256 := hex.EncodeToString(hasher.Sum(nil))
	valid := actualSHA256 == download.SHA256

	result := &DownloadResult{
		Filename:       destPath,
		ExpectedSHA256: download.SHA256,
		ActualSHA256:   actualSHA256,
		Valid:          valid,
	}

	if !valid {
		return result, errors.Newf("SHA256 mismatch: expected %s, got %s", download.SHA256, actualSHA256)
	}

	return result, nil
}

// GetLatestBuild returns the number of the latest build for the specified version
func (c *Client) GetLatestBuild(ctx context.Context, projectID, version string) (int32, error) {
	versionInfo, err := c.GetVersion(ctx, projectID, version)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get version info")
	}

	if len(versionInfo.Builds) == 0 {
		return 0, errors.New("no builds found for this version")
	}

	// The latest build is usually the last one in the list
	return versionInfo.Builds[len(versionInfo.Builds)-1], nil
}

// GetLatestVersion returns the latest available version for a project
func (c *Client) GetLatestVersion(ctx context.Context, projectID string) (string, error) {
	projectInfo, err := c.GetProject(ctx, projectID)
	if err != nil {
		return "", errors.Wrap(err, "failed to get project info")
	}

	if len(projectInfo.Versions) == 0 {
		return "", errors.New("no versions found for this project")
	}

	// The latest version is usually the last one in the list
	return projectInfo.Versions[len(projectInfo.Versions)-1], nil
}

// GetDefaultDownloadName returns the name of the main downloadable file for a build
func (c *Client) GetDefaultDownloadName(ctx context.Context, projectID, version string, build int32) (string, error) {
	buildInfo, err := c.GetBuild(ctx, projectID, version, build)
	if err != nil {
		return "", errors.Wrap(err, "failed to get build info")
	}

	// Check if there's an "application" download (which is common for PaperMC)
	if download, ok := buildInfo.Downloads["application"]; ok {
		return download.Name, nil
	}

	// As a fallback, look for any download with a .jar extension
	for _, download := range buildInfo.Downloads {
		if filepath.Ext(download.Name) == ".jar" {
			return download.Name, nil
		}
	}

	// If no jar file, get the first file's name
	for _, download := range buildInfo.Downloads {
		return download.Name, nil
	}

	return "", errors.New("no downloads found for this build")
}
