package api

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
)

// DownloadLatestBuild downloads the latest build of the specified project version.
func (c *Client) DownloadLatestBuild(ctx context.Context, projectID, version, destDir string) (*DownloadResult, error) {
	// Get the latest build
	build, err := c.GetLatestBuildV3(ctx, projectID, version)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get latest build")
	}

	// Get download name
	downloadName := build.GetDownloadName()
	if downloadName == "" {
		return nil, errors.New("no download found for this build")
	}

	// Form the save path
	destPath := filepath.Join(destDir, downloadName)

	// Download the file
	result, err := c.DownloadFile(ctx, projectID, version, build.ID, destPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download file")
	}

	return result, nil
}

// DownloadLatestStableVersion downloads the latest stable version of the project.
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

// FindPromotedBuild finds a recommended (promoted) build for the specified version.
// In v3 API, looks for builds with RECOMMENDED channel first, then STABLE.
func (c *Client) FindPromotedBuild(ctx context.Context, projectID, version string) (int32, error) {
	// Try to find RECOMMENDED builds first
	builds, err := c.GetBuilds(ctx, projectID, version, ChannelRecommended)
	if err == nil && len(builds) > 0 {
		return builds[len(builds)-1].ID, nil
	}

	// Fallback to STABLE builds
	builds, err = c.GetBuilds(ctx, projectID, version, ChannelStable)
	if err == nil && len(builds) > 0 {
		return builds[len(builds)-1].ID, nil
	}

	// Fallback to latest build
	latestBuild, err := c.GetLatestBuildV3(ctx, projectID, version)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get latest build")
	}

	return latestBuild.ID, nil
}

// DownloadPromotedBuild downloads the recommended build of the specified version.
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
	result, err := c.DownloadFile(ctx, projectID, version, buildNum, destPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download file")
	}

	return result, nil
}

// GetRecommendedVersion returns the recommended version for the project.
// Usually it's the latest stable (not SNAPSHOT and not pre/rc) version.
func (c *Client) GetRecommendedVersion(ctx context.Context, projectID string) (string, error) {
	projectInfo, err := c.GetProject(ctx, projectID)
	if err != nil {
		return "", errors.Wrap(err, "failed to get project info")
	}

	versions := projectInfo.FlattenVersions()
	if len(versions) == 0 {
		return "", errors.New("no versions found for this project")
	}

	// Look for versions without SNAPSHOT, pre, and rc, starting from the end (from new to old)
	for i := len(versions) - 1; i >= 0; i-- {
		version := versions[i]
		if !isSnapshotOrPreRelease(version) {
			return version, nil
		}
	}

	// If a stable version is not found, return the latest
	return versions[len(versions)-1], nil
}

// isSnapshotOrPreRelease checks if a version is pre-release, RC, or SNAPSHOT.
func isSnapshotOrPreRelease(version string) bool {
	lower := strings.ToLower(version)
	return strings.Contains(lower, "snapshot") ||
		strings.Contains(lower, "-pre") ||
		strings.Contains(lower, "-rc")
}

// GetLatestBuildURL returns the download URL for the latest build of a version.
func (c *Client) GetLatestBuildURL(ctx context.Context, projectID, version string) (string, error) {
	build, err := c.GetLatestBuildV3(ctx, projectID, version)
	if err != nil {
		return "", errors.Wrap(err, "failed to get latest build")
	}

	url := build.GetDownloadURL()
	if url == "" {
		return "", errors.New("no download URL found for this build")
	}

	return url, nil
}

// GetLatestVersionURL returns the download URL for the latest version of a project.
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

// GetPromotedBuildURL returns the download URL for the promoted build of a version.
func (c *Client) GetPromotedBuildURL(ctx context.Context, projectID, version string) (string, error) {
	// Find the promoted build
	buildNum, err := c.FindPromotedBuild(ctx, projectID, version)
	if err != nil {
		return "", errors.Wrap(err, "failed to find promoted build")
	}

	build, err := c.GetBuild(ctx, projectID, version, buildNum)
	if err != nil {
		return "", errors.Wrap(err, "failed to get build info")
	}

	url := build.GetDownloadURL()
	if url == "" {
		return "", errors.New("no download URL found for this build")
	}

	return url, nil
}

// GetBuildURL returns the download URL for a specific build.
func (c *Client) GetBuildURL(ctx context.Context, projectID, version string, build int32) (string, error) {
	buildInfo, err := c.GetBuild(ctx, projectID, version, build)
	if err != nil {
		return "", errors.Wrap(err, "failed to get build info")
	}

	url := buildInfo.GetDownloadURL()
	if url == "" {
		return "", errors.New("no download URL found for this build")
	}

	return url, nil
}
