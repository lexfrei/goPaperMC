package api

import (
	"sort"
	"time"

	"golang.org/x/mod/semver"
)

// ProjectMeta represents basic project information in v3 API.
type ProjectMeta struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SupportInfo represents version support status in v3 API.
type SupportInfo struct {
	Status string     `json:"status"` // SUPPORTED, DEPRECATED, UNSUPPORTED
	End    *time.Time `json:"end,omitempty"`
}

// JavaVersion represents Java version requirements in v3 API.
type JavaVersion struct {
	Minimum int `json:"minimum"`
}

// JavaFlags represents recommended JVM flags in v3 API.
type JavaFlags struct {
	Recommended []string `json:"recommended"`
}

// JavaInfo represents Java configuration in v3 API.
type JavaInfo struct {
	Version JavaVersion `json:"version"`
	Flags   JavaFlags   `json:"flags"`
}

// VersionMeta represents version metadata in v3 API.
type VersionMeta struct {
	ID      string      `json:"id"`
	Support SupportInfo `json:"support"`
	Java    JavaInfo    `json:"java"`
}

// ChecksumsV3 represents file checksums in v3 API.
type ChecksumsV3 struct {
	SHA256 string `json:"sha256"`
}

// DownloadV3 represents a downloadable file in v3 API.
type DownloadV3 struct {
	Name      string      `json:"name"`
	URL       string      `json:"url"`
	Checksums ChecksumsV3 `json:"checksums"`
	Size      int64       `json:"size"`
}

// CommitV3 represents a commit in v3 API.
type CommitV3 struct {
	SHA     string    `json:"sha"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

// ProjectV3Info represents a project with its versions in v3 projects list.
type ProjectV3Info struct {
	Project  ProjectMeta         `json:"project"`
	Versions map[string][]string `json:"versions"` // e.g., {"1.21": ["1.21.11", "1.21.10"]}
}

// ProjectsV3Response represents the v3 API response for projects list.
type ProjectsV3Response struct {
	Projects []ProjectV3Info `json:"projects"`
}

// ProjectV3Response represents the v3 API response for a single project.
type ProjectV3Response struct {
	Project  ProjectMeta         `json:"project"`
	Versions map[string][]string `json:"versions"`
}

// FlattenVersions returns all versions as a flat slice, sorted from oldest to newest using semver.
func (p *ProjectV3Response) FlattenVersions() []string {
	var allVersions []string

	// Collect all versions from all groups
	for _, versions := range p.Versions {
		allVersions = append(allVersions, versions...)
	}

	// Sort using semver (requires "v" prefix)
	sort.Slice(allVersions, func(i, j int) bool {
		vi := "v" + allVersions[i]
		vj := "v" + allVersions[j]
		return semver.Compare(vi, vj) < 0
	})

	return allVersions
}

// VersionV3Response represents the v3 API response for a version.
type VersionV3Response struct {
	Version VersionMeta `json:"version"`
	Builds  []int32     `json:"builds"`
}

// BuildV3Response represents the v3 API response for a build.
type BuildV3Response struct {
	ID        int32                 `json:"id"`
	Time      time.Time             `json:"time"`
	Channel   string                `json:"channel"`
	Commits   []CommitV3            `json:"commits"`
	Downloads map[string]DownloadV3 `json:"downloads"`
}

// GetDownloadURL returns the download URL for the default server application.
func (b *BuildV3Response) GetDownloadURL() string {
	if download, ok := b.Downloads["server:default"]; ok {
		return download.URL
	}
	// Fallback: return first available download URL
	for _, download := range b.Downloads {
		return download.URL
	}
	return ""
}

// GetDownloadName returns the filename of the default server application.
func (b *BuildV3Response) GetDownloadName() string {
	if download, ok := b.Downloads["server:default"]; ok {
		return download.Name
	}
	for _, download := range b.Downloads {
		return download.Name
	}
	return ""
}

// GetDownloadSHA256 returns the SHA256 checksum of the default server application.
func (b *BuildV3Response) GetDownloadSHA256() string {
	if download, ok := b.Downloads["server:default"]; ok {
		return download.Checksums.SHA256
	}
	for _, download := range b.Downloads {
		return download.Checksums.SHA256
	}
	return ""
}

// BuildsV3Response represents the v3 API response for builds list.
type BuildsV3Response struct {
	Builds []BuildV3Response `json:"builds"`
}

// Channel represents build stability channels in v3 API.
type Channel string

const (
	ChannelAlpha       Channel = "alpha"
	ChannelBeta        Channel = "beta"
	ChannelStable      Channel = "stable"
	ChannelRecommended Channel = "recommended"
)

// SupportStatus represents version support status in v3 API.
type SupportStatus string

const (
	SupportStatusSupported   SupportStatus = "SUPPORTED"
	SupportStatusDeprecated  SupportStatus = "DEPRECATED"
	SupportStatusUnsupported SupportStatus = "UNSUPPORTED"
)
