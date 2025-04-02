package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
)

const (
	// DefaultBaseURL - base URL for PaperMC API
	DefaultBaseURL = "https://api.papermc.io"
	// DefaultTimeout - default timeout for HTTP requests
	DefaultTimeout = 30 * time.Second
)

// Client represents the PaperMC API client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Limit      int // Limit the number of items to return (0 means no limit)
}

// NewClient creates a new instance of the PaperMC API client
func NewClient() *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		Limit: 0,
	}
}

// WithBaseURL sets a custom base URL for the API
func (c *Client) WithBaseURL(baseURL string) *Client {
	c.BaseURL = baseURL
	return c
}

// WithTimeout sets a timeout for the HTTP client
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.HTTPClient.Timeout = timeout
	return c
}

// WithLimit sets a limit for the number of items to return
func (c *Client) WithLimit(limit int) *Client {
	c.Limit = limit
	return c
}

// GetProjects returns a list of all available projects
func (c *Client) GetProjects(ctx context.Context) (*ProjectsResponse, error) {
	url := fmt.Sprintf("%s/v2/projects", c.BaseURL)
	
	resp, err := c.makeRequest(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request projects")
	}
	defer resp.Body.Close()
	
	var projectsResp ProjectsResponse
	if err := json.NewDecoder(resp.Body).Decode(&projectsResp); err != nil {
		return nil, errors.Wrap(err, "failed to decode projects response")
	}
	
	// Apply limit if set
	if c.Limit > 0 && len(projectsResp.Projects) > c.Limit {
		start := len(projectsResp.Projects) - c.Limit
		projectsResp.Projects = projectsResp.Projects[start:]
	}
	
	return &projectsResp, nil
}

// GetProject returns information about a specific project
func (c *Client) GetProject(ctx context.Context, projectID string) (*ProjectResponse, error) {
	url := fmt.Sprintf("%s/v2/projects/%s", c.BaseURL, projectID)
	
	resp, err := c.makeRequest(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request project")
	}
	defer resp.Body.Close()
	
	var projectResp ProjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&projectResp); err != nil {
		return nil, errors.Wrap(err, "failed to decode project response")
	}
	
	// Apply limit to versions if set
	if c.Limit > 0 && len(projectResp.Versions) > c.Limit {
		start := len(projectResp.Versions) - c.Limit
		projectResp.Versions = projectResp.Versions[start:]
	}
	
	// Apply limit to version groups if set
	if c.Limit > 0 && len(projectResp.VersionGroups) > c.Limit {
		start := len(projectResp.VersionGroups) - c.Limit
		projectResp.VersionGroups = projectResp.VersionGroups[start:]
	}
	
	return &projectResp, nil
}

// GetVersion returns information about a project version
func (c *Client) GetVersion(ctx context.Context, projectID, version string) (*VersionResponse, error) {
	url := fmt.Sprintf("%s/v2/projects/%s/versions/%s", c.BaseURL, projectID, version)
	
	resp, err := c.makeRequest(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request version")
	}
	defer resp.Body.Close()
	
	var versionResp VersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&versionResp); err != nil {
		return nil, errors.Wrap(err, "failed to decode version response")
	}
	
	// Apply limit to builds if set
	if c.Limit > 0 && len(versionResp.Builds) > c.Limit {
		start := len(versionResp.Builds) - c.Limit
		versionResp.Builds = versionResp.Builds[start:]
	}
	
	return &versionResp, nil
}

// GetBuilds returns a list of available builds for a project version
func (c *Client) GetBuilds(ctx context.Context, projectID, version string) (*BuildsResponse, error) {
	url := fmt.Sprintf("%s/v2/projects/%s/versions/%s/builds", c.BaseURL, projectID, version)
	
	resp, err := c.makeRequest(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request builds")
	}
	defer resp.Body.Close()
	
	var buildsResp BuildsResponse
	if err := json.NewDecoder(resp.Body).Decode(&buildsResp); err != nil {
		return nil, errors.Wrap(err, "failed to decode builds response")
	}
	
	// Apply limit to builds if set
	if c.Limit > 0 && len(buildsResp.Builds) > c.Limit {
		start := len(buildsResp.Builds) - c.Limit
		buildsResp.Builds = buildsResp.Builds[start:]
	}
	
	return &buildsResp, nil
}

// GetBuild returns information about a specific build
func (c *Client) GetBuild(ctx context.Context, projectID, version string, build int32) (*BuildResponse, error) {
	url := fmt.Sprintf("%s/v2/projects/%s/versions/%s/builds/%d", c.BaseURL, projectID, version, build)
	
	resp, err := c.makeRequest(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request build")
	}
	defer resp.Body.Close()
	
	var buildResp BuildResponse
	if err := json.NewDecoder(resp.Body).Decode(&buildResp); err != nil {
		return nil, errors.Wrap(err, "failed to decode build response")
	}
	
	// Apply limit to changes if set
	if c.Limit > 0 && len(buildResp.Changes) > c.Limit {
		start := len(buildResp.Changes) - c.Limit
		buildResp.Changes = buildResp.Changes[start:]
	}
	
	return &buildResp, nil
}

// DownloadBuild downloads the specified file from a build
func (c *Client) DownloadBuild(ctx context.Context, projectID, version string, build int32, download string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s/v2/projects/%s/versions/%s/builds/%d/downloads/%s", c.BaseURL, projectID, version, build, download)
	
	resp, err := c.makeRequest(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download build")
	}
	
	return resp.Body, nil
}

// GetVersionGroup returns information about a project's version group
func (c *Client) GetVersionGroup(ctx context.Context, projectID, family string) (*VersionFamilyResponse, error) {
	url := fmt.Sprintf("%s/v2/projects/%s/version_group/%s", c.BaseURL, projectID, family)
	
	resp, err := c.makeRequest(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request version group")
	}
	defer resp.Body.Close()
	
	var versionFamilyResp VersionFamilyResponse
	if err := json.NewDecoder(resp.Body).Decode(&versionFamilyResp); err != nil {
		return nil, errors.Wrap(err, "failed to decode version family response")
	}
	
	// Apply limit to versions if set
	if c.Limit > 0 && len(versionFamilyResp.Versions) > c.Limit {
		start := len(versionFamilyResp.Versions) - c.Limit
		versionFamilyResp.Versions = versionFamilyResp.Versions[start:]
	}
	
	return &versionFamilyResp, nil
}

// GetVersionGroupBuilds returns a list of available builds for a version group
func (c *Client) GetVersionGroupBuilds(ctx context.Context, projectID, family string) (*VersionFamilyBuildsResponse, error) {
	url := fmt.Sprintf("%s/v2/projects/%s/version_group/%s/builds", c.BaseURL, projectID, family)
	
	resp, err := c.makeRequest(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request version group builds")
	}
	defer resp.Body.Close()
	
	var versionFamilyBuildsResp VersionFamilyBuildsResponse
	if err := json.NewDecoder(resp.Body).Decode(&versionFamilyBuildsResp); err != nil {
		return nil, errors.Wrap(err, "failed to decode version family builds response")
	}
	
	// Apply limit to versions if set
	if c.Limit > 0 && len(versionFamilyBuildsResp.Versions) > c.Limit {
		start := len(versionFamilyBuildsResp.Versions) - c.Limit
		versionFamilyBuildsResp.Versions = versionFamilyBuildsResp.Versions[start:]
	}
	
	// Apply limit to builds if set
	if c.Limit > 0 && len(versionFamilyBuildsResp.Builds) > c.Limit {
		start := len(versionFamilyBuildsResp.Builds) - c.Limit
		versionFamilyBuildsResp.Builds = versionFamilyBuildsResp.Builds[start:]
	}
	
	return &versionFamilyBuildsResp, nil
}

// makeRequest performs an HTTP request to the API
func (c *Client) makeRequest(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute request")
	}
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, errors.Newf("API returned non-OK status: %d, body: %s", resp.StatusCode, body)
	}
	
	return resp, nil
}
