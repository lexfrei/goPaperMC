package api

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/cockroachdb/errors"
)

// DownloadLatestBuild downloads the latest build of the specified project version
func (c *Client) DownloadLatestBuild(ctx context.Context, projectID, version, destDir string) (*DownloadResult, error) {
	// Get the latest build number
	buildNum, err := c.GetLatestBuild(ctx, projectID, version)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get latest build number")
	}

	// Get the default file name
	downloadName, err := c.GetDefaultDownloadName(ctx, projectID, version, buildNum)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get default download name")
	}

	// Form the save path
	destPath := filepath.Join(destDir, downloadName)

	// Download the file
	result, err := c.DownloadFile(ctx, projectID, version, buildNum, downloadName, destPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download file")
	}

	return result, nil
}

// DownloadLatestStableVersion downloads the latest stable version of the project
func (c *Client) DownloadLatestStableVersion(ctx context.Context, projectID, destDir string) (*DownloadResult, error) {
	// Get the latest version
	version, err := c.GetLatestVersion(ctx, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get latest version")
	}

	// Download the latest build of this version
	result, err := c.DownloadLatestBuild(ctx, projectID, version, destDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download latest build")
	}

	return result, nil
}

// FindPromotedBuild finds a recommended (promoted) build for the specified version
func (c *Client) FindPromotedBuild(ctx context.Context, projectID, version string) (int32, error) {
	builds, err := c.GetBuilds(ctx, projectID, version)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get builds")
	}

	// Look for promoted builds, starting from the end (from new to old)
	for i := len(builds.Builds) - 1; i >= 0; i-- {
		if builds.Builds[i].Promoted {
			return builds.Builds[i].Build, nil
		}
	}

	// If a promoted build is not found, return the latest
	if len(builds.Builds) > 0 {
		return builds.Builds[len(builds.Builds)-1].Build, nil
	}

	return 0, errors.New("no builds found for this version")
}

// DownloadPromotedBuild downloads the recommended build of the specified version
func (c *Client) DownloadPromotedBuild(ctx context.Context, projectID, version, destDir string) (*DownloadResult, error) {
	// Find the promoted build
	buildNum, err := c.FindPromotedBuild(ctx, projectID, version)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find promoted build")
	}

	// Get the default file name
	downloadName, err := c.GetDefaultDownloadName(ctx, projectID, version, buildNum)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get default download name")
	}

	// Form the save path
	destPath := filepath.Join(destDir, downloadName)

	// Download the file
	result, err := c.DownloadFile(ctx, projectID, version, buildNum, downloadName, destPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download file")
	}

	return result, nil
}

// GetRecommendedVersion returns the recommended version for the project
// Usually it's the latest stable (not SNAPSHOT and not pre) version
func (c *Client) GetRecommendedVersion(ctx context.Context, projectID string) (string, error) {
	projectInfo, err := c.GetProject(ctx, projectID)
	if err != nil {
		return "", errors.Wrap(err, "failed to get project info")
	}

	if len(projectInfo.Versions) == 0 {
		return "", errors.New("no versions found for this project")
	}

	// Look for versions without SNAPSHOT and pre, starting from the end (from new to old)
	for i := len(projectInfo.Versions) - 1; i >= 0; i-- {
		version := projectInfo.Versions[i]
		if !isSnapshotOrPreRelease(version) {
			return version, nil
		}
	}

	// If a stable version is not found, return the latest
	return projectInfo.Versions[len(projectInfo.Versions)-1], nil
}

// isSnapshotOrPreRelease checks if a version is pre-release or SNAPSHOT
func isSnapshotOrPreRelease(version string) bool {
	return contains(version, "SNAPSHOT") || contains(version, "pre")
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}

// FormatDownloadURL returns a URL for direct file download
func (c *Client) FormatDownloadURL(projectID, version string, build int32, downloadName string) string {
	return fmt.Sprintf("%s/v2/projects/%s/versions/%s/builds/%d/downloads/%s", 
		c.BaseURL, projectID, version, build, downloadName)
}

// GetLatestBuildURL returns the download URL for the latest build of a version
func (c *Client) GetLatestBuildURL(ctx context.Context, projectID, version string) (string, error) {
	// Get the latest build number
	buildNum, err := c.GetLatestBuild(ctx, projectID, version)
	if err != nil {
		return "", errors.Wrap(err, "failed to get latest build number")
	}

	// Get the default file name
	downloadName, err := c.GetDefaultDownloadName(ctx, projectID, version, buildNum)
	if err != nil {
		return "", errors.Wrap(err, "failed to get default download name")
	}

	// Format and return the URL
	return c.FormatDownloadURL(projectID, version, buildNum, downloadName), nil
}

// GetLatestVersionURL returns the download URL for the latest version of a project
func (c *Client) GetLatestVersionURL(ctx context.Context, projectID string) (string, error) {
	// Get the latest version
	version, err := c.GetLatestVersion(ctx, projectID)
	if err != nil {
		return "", errors.Wrap(err, "failed to get latest version")
	}

	// Get the URL for the latest build of this version
	url, err := c.GetLatestBuildURL(ctx, projectID, version)
	if err != nil {
		return "", errors.Wrap(err, "failed to get latest build URL")
	}

	return url, nil
}

// GetPromotedBuildURL returns the download URL for the promoted build of a version
func (c *Client) GetPromotedBuildURL(ctx context.Context, projectID, version string) (string, error) {
	// Find the promoted build
	buildNum, err := c.FindPromotedBuild(ctx, projectID, version)
	if err != nil {
		return "", errors.Wrap(err, "failed to find promoted build")
	}

	// Get the default file name
	downloadName, err := c.GetDefaultDownloadName(ctx, projectID, version, buildNum)
	if err != nil {
		return "", errors.Wrap(err, "failed to get default download name")
	}

	// Format and return the URL
	return c.FormatDownloadURL(projectID, version, buildNum, downloadName), nil
}

// GetBuildURL returns the download URL for a specific build
func (c *Client) GetBuildURL(ctx context.Context, projectID, version string, build int32) (string, error) {
	// Get the default file name
	downloadName, err := c.GetDefaultDownloadName(ctx, projectID, version, build)
	if err != nil {
		return "", errors.Wrap(err, "failed to get default download name")
	}

	// Format and return the URL
	return c.FormatDownloadURL(projectID, version, build, downloadName), nil
}
