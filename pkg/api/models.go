package api

import "time"

// ProjectsResponse represents a list of all available projects
type ProjectsResponse struct {
	Projects []string `json:"projects"`
}

// ProjectResponse represents information about a project
type ProjectResponse struct {
	ProjectID     string   `json:"project_id"`
	ProjectName   string   `json:"project_name"`
	VersionGroups []string `json:"version_groups"`
	Versions      []string `json:"versions"`
}

// VersionResponse represents information about a version
type VersionResponse struct {
	ProjectID   string  `json:"project_id"`
	ProjectName string  `json:"project_name"`
	Version     string  `json:"version"`
	Builds      []int32 `json:"builds"`
}

// BuildsResponse represents a list of available builds for a project version
type BuildsResponse struct {
	ProjectID   string        `json:"project_id"`
	ProjectName string        `json:"project_name"`
	Version     string        `json:"version"`
	Builds      []VersionBuild `json:"builds"`
}

// BuildResponse represents information about a specific build
type BuildResponse struct {
	ProjectID   string               `json:"project_id"`
	ProjectName string               `json:"project_name"`
	Version     string               `json:"version"`
	Build       int32                `json:"build"`
	Time        time.Time            `json:"time"`
	Channel     string               `json:"channel"`
	Promoted    bool                 `json:"promoted"`
	Changes     []Change             `json:"changes"`
	Downloads   map[string]Download  `json:"downloads"`
}

// VersionBuild represents information about a build in the version context
type VersionBuild struct {
	Build     int32                `json:"build"`
	Time      time.Time            `json:"time"`
	Channel   string               `json:"channel"`
	Promoted  bool                 `json:"promoted"`
	Changes   []Change             `json:"changes"`
	Downloads map[string]Download  `json:"downloads"`
}

// Change represents information about changes in a build
type Change struct {
	Commit  string `json:"commit"`
	Summary string `json:"summary"`
	Message string `json:"message"`
}

// Download represents information about a downloadable file
type Download struct {
	Name   string `json:"name"`
	SHA256 string `json:"sha256"`
}

// VersionFamilyResponse represents information about a project's version group
type VersionFamilyResponse struct {
	ProjectID    string   `json:"project_id"`
	ProjectName  string   `json:"project_name"`
	VersionGroup string   `json:"version_group"`
	Versions     []string `json:"versions"`
}

// VersionFamilyBuildsResponse represents a list of available builds for a version group
type VersionFamilyBuildsResponse struct {
	ProjectID    string               `json:"project_id"`
	ProjectName  string               `json:"project_name"`
	VersionGroup string               `json:"version_group"`
	Versions     []string             `json:"versions"`
	Builds       []VersionFamilyBuild `json:"builds"`
}

// VersionFamilyBuild represents information about a build in the version group context
type VersionFamilyBuild struct {
	Version   string               `json:"version"`
	Build     int32                `json:"build"`
	Time      time.Time            `json:"time"`
	Channel   string               `json:"channel"`
	Promoted  bool                 `json:"promoted"`
	Changes   []Change             `json:"changes"`
	Downloads map[string]Download  `json:"downloads"`
}
