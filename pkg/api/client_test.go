package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProjectsV3(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/projects" {
			t.Errorf("Expected path /v3/projects, got %s", r.URL.Path)
		}

		resp := ProjectsV3Response{
			Projects: []ProjectV3Info{
				{
					Project: ProjectMeta{ID: "paper", Name: "Paper"},
					Versions: map[string][]string{
						"1.21": {"1.21.11", "1.21.10"},
					},
				},
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient().WithBaseURL(server.URL)
	projects, err := client.GetProjects(context.Background())
	if err != nil {
		t.Fatalf("GetProjects failed: %v", err)
	}

	if len(projects.Projects) != 1 {
		t.Errorf("Expected 1 project, got %d", len(projects.Projects))
	}

	if projects.Projects[0].Project.ID != "paper" {
		t.Errorf("Expected project ID 'paper', got '%s'", projects.Projects[0].Project.ID)
	}
}

func TestGetProjectV3(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/projects/paper" {
			t.Errorf("Expected path /v3/projects/paper, got %s", r.URL.Path)
		}

		resp := ProjectV3Response{
			Project: ProjectMeta{ID: "paper", Name: "Paper"},
			Versions: map[string][]string{
				"1.21": {"1.21.11", "1.21.11-rc3", "1.21.10"},
				"1.20": {"1.20.6", "1.20.5"},
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient().WithBaseURL(server.URL)
	project, err := client.GetProject(context.Background(), "paper")
	if err != nil {
		t.Fatalf("GetProject failed: %v", err)
	}

	if project.Project.ID != "paper" {
		t.Errorf("Expected project ID 'paper', got '%s'", project.Project.ID)
	}

	if len(project.Versions) != 2 {
		t.Errorf("Expected 2 version groups, got %d", len(project.Versions))
	}
}

func TestGetProjectV3_FlattenVersions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ProjectV3Response{
			Project: ProjectMeta{ID: "paper", Name: "Paper"},
			Versions: map[string][]string{
				"1.21": {"1.21.11", "1.21.10"},
				"1.20": {"1.20.6"},
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient().WithBaseURL(server.URL)
	project, err := client.GetProject(context.Background(), "paper")
	if err != nil {
		t.Fatalf("GetProject failed: %v", err)
	}

	flattened := project.FlattenVersions()
	if len(flattened) != 3 {
		t.Errorf("Expected 3 flattened versions, got %d", len(flattened))
	}

	// Verify all versions are present (order depends on map iteration and sort)
	versionSet := make(map[string]bool)
	for _, v := range flattened {
		versionSet[v] = true
	}
	if !versionSet["1.21.11"] || !versionSet["1.21.10"] || !versionSet["1.20.6"] {
		t.Errorf("Expected all versions to be present, got %v", flattened)
	}
}

func TestGetVersionV3(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/projects/paper/versions/1.21.11" {
			t.Errorf("Expected path /v3/projects/paper/versions/1.21.11, got %s", r.URL.Path)
		}

		resp := VersionV3Response{
			Version: VersionMeta{
				ID: "1.21.11",
				Support: SupportInfo{
					Status: "SUPPORTED",
				},
			},
			Builds: []int32{74, 73, 72},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient().WithBaseURL(server.URL)
	version, err := client.GetVersion(context.Background(), "paper", "1.21.11")
	if err != nil {
		t.Fatalf("GetVersion failed: %v", err)
	}

	if version.Version.ID != "1.21.11" {
		t.Errorf("Expected version ID '1.21.11', got '%s'", version.Version.ID)
	}

	if len(version.Builds) != 3 {
		t.Errorf("Expected 3 builds, got %d", len(version.Builds))
	}
}

func TestGetVersionV3_RCVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/projects/paper/versions/1.21.11-rc3" {
			t.Errorf("Expected path /v3/projects/paper/versions/1.21.11-rc3, got %s", r.URL.Path)
		}

		resp := VersionV3Response{
			Version: VersionMeta{
				ID: "1.21.11-rc3",
				Support: SupportInfo{
					Status: "UNSUPPORTED",
				},
			},
			Builds: []int32{31, 30},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient().WithBaseURL(server.URL)
	version, err := client.GetVersion(context.Background(), "paper", "1.21.11-rc3")
	if err != nil {
		t.Fatalf("GetVersion for RC failed: %v", err)
	}

	if len(version.Builds) != 2 {
		t.Errorf("Expected 2 builds for RC version, got %d", len(version.Builds))
	}
}

func TestGetLatestBuildV3(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/projects/paper/versions/1.21.11/builds/latest" {
			t.Errorf("Expected path /v3/projects/paper/versions/1.21.11/builds/latest, got %s", r.URL.Path)
		}

		resp := BuildV3Response{
			ID:      74,
			Channel: "STABLE",
			Downloads: map[string]DownloadV3{
				"server:default": {
					Name: "paper-1.21.11-74.jar",
					URL:  "https://example.com/paper-1.21.11-74.jar",
					Checksums: ChecksumsV3{
						SHA256: "abc123",
					},
					Size: 54819307,
				},
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient().WithBaseURL(server.URL)
	build, err := client.GetLatestBuildV3(context.Background(), "paper", "1.21.11")
	if err != nil {
		t.Fatalf("GetLatestBuildV3 failed: %v", err)
	}

	if build.ID != 74 {
		t.Errorf("Expected build ID 74, got %d", build.ID)
	}

	if build.Channel != "STABLE" {
		t.Errorf("Expected channel 'STABLE', got '%s'", build.Channel)
	}

	download, ok := build.Downloads["server:default"]
	if !ok {
		t.Fatal("Expected 'server:default' download")
	}

	if download.URL == "" {
		t.Error("Expected direct download URL in response")
	}
}

func TestGetLatestBuildV3_RCVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/projects/paper/versions/1.21.11-rc3/builds/latest" {
			t.Errorf("Expected path for RC version, got %s", r.URL.Path)
		}

		resp := BuildV3Response{
			ID:      31,
			Channel: "STABLE",
			Downloads: map[string]DownloadV3{
				"server:default": {
					Name: "paper-1.21.11-rc3-31.jar",
					URL:  "https://example.com/paper-1.21.11-rc3-31.jar",
					Checksums: ChecksumsV3{
						SHA256: "def456",
					},
				},
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient().WithBaseURL(server.URL)
	build, err := client.GetLatestBuildV3(context.Background(), "paper", "1.21.11-rc3")
	if err != nil {
		t.Fatalf("GetLatestBuildV3 for RC version should work: %v", err)
	}

	if build.ID != 31 {
		t.Errorf("Expected build ID 31 for RC version, got %d", build.ID)
	}
}

func TestGetBuildV3(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/projects/paper/versions/1.21.11/builds/74" {
			t.Errorf("Expected path /v3/projects/paper/versions/1.21.11/builds/74, got %s", r.URL.Path)
		}

		resp := BuildV3Response{
			ID:      74,
			Channel: "STABLE",
			Downloads: map[string]DownloadV3{
				"server:default": {
					Name: "paper-1.21.11-74.jar",
					URL:  "https://example.com/paper-1.21.11-74.jar",
					Checksums: ChecksumsV3{
						SHA256: "abc123",
					},
					Size: 54819307,
				},
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient().WithBaseURL(server.URL)
	build, err := client.GetBuild(context.Background(), "paper", "1.21.11", 74)
	if err != nil {
		t.Fatalf("GetBuild failed: %v", err)
	}

	if build.ID != 74 {
		t.Errorf("Expected build ID 74, got %d", build.ID)
	}
}

func TestFlattenVersions_SemverSort(t *testing.T) {
	// Test that versions are sorted correctly using semver, not lexicographically
	// This prevents the bug where 1.7, 1.8, 1.9 would appear after 1.10-1.21
	project := &ProjectV3Response{
		Project: ProjectMeta{ID: "paper", Name: "Paper"},
		Versions: map[string][]string{
			// Intentionally unordered to test sorting
			"1.21": {"1.21.11", "1.21.11-rc3", "1.21.11-rc2", "1.21.10"},
			"1.7":  {"1.7.10"},
			"1.20": {"1.20.6", "1.20.4"},
			"1.10": {"1.10.2"},
			"1.9":  {"1.9.4"},
		},
	}

	flattened := project.FlattenVersions()

	// Expected order: 1.7.10 < 1.9.4 < 1.10.2 < 1.20.4 < 1.20.6 < 1.21.10 < 1.21.11-rc2 < 1.21.11-rc3 < 1.21.11
	expected := []string{
		"1.7.10",
		"1.9.4",
		"1.10.2",
		"1.20.4",
		"1.20.6",
		"1.21.10",
		"1.21.11-rc2",
		"1.21.11-rc3",
		"1.21.11",
	}

	if len(flattened) != len(expected) {
		t.Fatalf("Expected %d versions, got %d: %v", len(expected), len(flattened), flattened)
	}

	for i, v := range expected {
		if flattened[i] != v {
			t.Errorf("Position %d: expected %s, got %s. Full list: %v", i, v, flattened[i], flattened)
		}
	}

	// Verify that taking last 3 gives us the newest versions (the original bug fix)
	last3 := flattened[len(flattened)-3:]
	expectedLast3 := []string{"1.21.11-rc2", "1.21.11-rc3", "1.21.11"}
	for i, v := range expectedLast3 {
		if last3[i] != v {
			t.Errorf("Last 3 position %d: expected %s, got %s", i, v, last3[i])
		}
	}
}

func TestGetBuildDownloadURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := BuildV3Response{
			ID:      74,
			Channel: "STABLE",
			Downloads: map[string]DownloadV3{
				"server:default": {
					Name: "paper-1.21.11-74.jar",
					URL:  "https://fill-data.papermc.io/v1/objects/abc/paper-1.21.11-74.jar",
					Checksums: ChecksumsV3{
						SHA256: "abc123",
					},
				},
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient().WithBaseURL(server.URL)
	build, err := client.GetLatestBuildV3(context.Background(), "paper", "1.21.11")
	if err != nil {
		t.Fatalf("GetLatestBuildV3 failed: %v", err)
	}

	url := build.GetDownloadURL()
	if url != "https://fill-data.papermc.io/v1/objects/abc/paper-1.21.11-74.jar" {
		t.Errorf("Expected direct download URL, got '%s'", url)
	}
}
