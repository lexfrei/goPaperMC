package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
)

const (
	// DefaultBaseURL is the base URL for PaperMC API v3.
	DefaultBaseURL = "https://fill.papermc.io"
	// DefaultTimeout is the default timeout for HTTP requests.
	DefaultTimeout = 30 * time.Second
)

// channelToAPI maps lowercase channel names to API format (uppercase).
var channelToAPI = map[Channel]string{
	ChannelAlpha:       "ALPHA",
	ChannelBeta:        "BETA",
	ChannelStable:      "STABLE",
	ChannelRecommended: "RECOMMENDED",
}

// Client represents the PaperMC API client.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Limit      int     // Limit the number of items to return (0 means no limit)
	Channel    Channel // Filter builds by channel (empty means no filter)
}

// NewClient creates a new instance of the PaperMC API client.
func NewClient() *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		Limit: 0,
	}
}

// WithBaseURL sets a custom base URL for the API.
func (c *Client) WithBaseURL(baseURL string) *Client {
	c.BaseURL = baseURL
	return c
}

// WithTimeout sets a timeout for the HTTP client.
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.HTTPClient.Timeout = timeout
	return c
}

// WithLimit sets a limit for the number of items to return.
func (c *Client) WithLimit(limit int) *Client {
	c.Limit = limit
	return c
}

// WithChannel sets a channel filter for builds.
func (c *Client) WithChannel(channel Channel) *Client {
	c.Channel = channel
	return c
}

// GetProjects returns a list of all available projects.
func (c *Client) GetProjects(ctx context.Context) (*ProjectsV3Response, error) {
	url := fmt.Sprintf("%s/v3/projects", c.BaseURL)

	resp, err := c.makeRequest(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request projects")
	}
	defer func() { _ = resp.Body.Close() }()

	var projectsResp ProjectsV3Response
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

// GetProject returns information about a specific project.
func (c *Client) GetProject(ctx context.Context, projectID string) (*ProjectV3Response, error) {
	url := fmt.Sprintf("%s/v3/projects/%s", c.BaseURL, projectID)

	resp, err := c.makeRequest(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request project")
	}
	defer func() { _ = resp.Body.Close() }()

	var projectResp ProjectV3Response
	if err := json.NewDecoder(resp.Body).Decode(&projectResp); err != nil {
		return nil, errors.Wrap(err, "failed to decode project response")
	}

	return &projectResp, nil
}

// GetVersions returns a list of all versions for a project with full metadata.
func (c *Client) GetVersions(ctx context.Context, projectID string) ([]VersionV3Response, error) {
	url := fmt.Sprintf("%s/v3/projects/%s/versions", c.BaseURL, projectID)

	resp, err := c.makeRequest(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request versions")
	}
	defer func() { _ = resp.Body.Close() }()

	var versions []VersionV3Response
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, errors.Wrap(err, "failed to decode versions response")
	}

	// Apply limit if set
	if c.Limit > 0 && len(versions) > c.Limit {
		start := len(versions) - c.Limit
		versions = versions[start:]
	}

	return versions, nil
}

// GetVersion returns information about a project version.
func (c *Client) GetVersion(ctx context.Context, projectID, version string) (*VersionV3Response, error) {
	url := fmt.Sprintf("%s/v3/projects/%s/versions/%s", c.BaseURL, projectID, version)

	resp, err := c.makeRequest(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request version")
	}
	defer func() { _ = resp.Body.Close() }()

	var versionResp VersionV3Response
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

// GetBuilds returns a list of available builds for a project version.
// Optionally filter by channels (ALPHA, BETA, STABLE, RECOMMENDED).
func (c *Client) GetBuilds(ctx context.Context, projectID, version string, channels ...Channel) ([]BuildV3Response, error) {
	url := fmt.Sprintf("%s/v3/projects/%s/versions/%s/builds", c.BaseURL, projectID, version)

	// Add channel filter if specified
	if len(channels) > 0 {
		channelStrs := make([]string, 0, len(channels))
		for _, ch := range channels {
			if apiCh, ok := channelToAPI[ch]; ok {
				channelStrs = append(channelStrs, "channel="+apiCh)
			}
		}
		if len(channelStrs) > 0 {
			url += "?" + strings.Join(channelStrs, "&")
		}
	}

	resp, err := c.makeRequest(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request builds")
	}
	defer func() { _ = resp.Body.Close() }()

	var builds []BuildV3Response
	if err := json.NewDecoder(resp.Body).Decode(&builds); err != nil {
		return nil, errors.Wrap(err, "failed to decode builds response")
	}

	// Apply limit to builds if set
	if c.Limit > 0 && len(builds) > c.Limit {
		start := len(builds) - c.Limit
		builds = builds[start:]
	}

	return builds, nil
}

// GetLatestBuildV3 returns the latest build for the specified version using v3 API.
// If channel filter is set, uses /builds endpoint with filter and returns the last one.
func (c *Client) GetLatestBuildV3(ctx context.Context, projectID, version string) (*BuildV3Response, error) {
	// If channel filter is set, use /builds endpoint (channel not supported on /builds/latest)
	if c.Channel != "" {
		builds, err := c.GetBuilds(ctx, projectID, version, c.Channel)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get builds with channel filter")
		}
		if len(builds) == 0 {
			return nil, errors.Newf("no builds found for channel %s", c.Channel)
		}
		// Find build with highest ID (API doesn't guarantee order)
		latest := &builds[0]
		for i := range builds {
			if builds[i].ID > latest.ID {
				latest = &builds[i]
			}
		}
		return latest, nil
	}

	url := fmt.Sprintf("%s/v3/projects/%s/versions/%s/builds/latest", c.BaseURL, projectID, version)

	resp, err := c.makeRequest(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request latest build")
	}
	defer func() { _ = resp.Body.Close() }()

	var buildResp BuildV3Response
	if err := json.NewDecoder(resp.Body).Decode(&buildResp); err != nil {
		return nil, errors.Wrap(err, "failed to decode build response")
	}

	return &buildResp, nil
}

// GetBuild returns information about a specific build.
func (c *Client) GetBuild(ctx context.Context, projectID, version string, build int32) (*BuildV3Response, error) {
	url := fmt.Sprintf("%s/v3/projects/%s/versions/%s/builds/%d", c.BaseURL, projectID, version, build)

	resp, err := c.makeRequest(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request build")
	}
	defer func() { _ = resp.Body.Close() }()

	var buildResp BuildV3Response
	if err := json.NewDecoder(resp.Body).Decode(&buildResp); err != nil {
		return nil, errors.Wrap(err, "failed to decode build response")
	}

	return &buildResp, nil
}

// DownloadBuild downloads the specified file from a build.
func (c *Client) DownloadBuild(ctx context.Context, downloadURL string) (io.ReadCloser, error) {
	resp, err := c.makeRequest(ctx, downloadURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download build")
	}

	return resp.Body, nil
}

// makeRequest performs an HTTP request to the API.
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
		_ = resp.Body.Close()
		return nil, errors.Newf("API returned non-OK status: %d, body: %s", resp.StatusCode, body)
	}

	return resp, nil
}
