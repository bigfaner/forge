package feature

import (
	"path/filepath"
	"testing"
)

func TestGetFeatureDir(t *testing.T) {
	tests := []struct {
		feature string
		want    string
	}{
		{"my-feature", filepath.Join("docs/features", "my-feature")},
		{"", filepath.Join("docs/features", "")},
		{"feature_123", filepath.Join("docs/features", "feature_123")},
	}

	for _, tt := range tests {
		t.Run(tt.feature, func(t *testing.T) {
			if got := GetFeatureDir(tt.feature); got != tt.want {
				t.Errorf("GetFeatureDir(%q) = %q, want %q", tt.feature, got, tt.want)
			}
		})
	}
}

func TestGetFeatureManifest(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "manifest.md")
	if got := GetFeatureManifest(feature); got != want {
		t.Errorf("GetFeatureManifest(%q) = %q, want %q", feature, got, want)
	}
}

func TestGetFeaturePRDDir(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "prd")
	if got := GetFeaturePRDDir(feature); got != want {
		t.Errorf("GetFeaturePRDDir(%q) = %q, want %q", feature, got, want)
	}
}

func TestGetFeaturePRDFile(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "prd", "prd-spec.md")
	if got := GetFeaturePRDFile(feature); got != want {
		t.Errorf("GetFeaturePRDFile(%q) = %q, want %q", feature, got, want)
	}
}

func TestGetFeatureUserStoriesFile(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "prd", "prd-user-stories.md")
	if got := GetFeatureUserStoriesFile(feature); got != want {
		t.Errorf("GetFeatureUserStoriesFile(%q) = %q, want %q", feature, got, want)
	}
}

func TestGetFeatureUIFunctionsFile(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "prd", "prd-ui-functions.md")
	if got := GetFeatureUIFunctionsFile(feature); got != want {
		t.Errorf("GetFeatureUIFunctionsFile(%q) = %q, want %q", feature, got, want)
	}
}

func TestGetFeatureDesignDir(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "design")
	if got := GetFeatureDesignDir(feature); got != want {
		t.Errorf("GetFeatureDesignDir(%q) = %q, want %q", feature, got, want)
	}
}

func TestGetFeatureDesignFile(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "design", "tech-design.md")
	if got := GetFeatureDesignFile(feature); got != want {
		t.Errorf("GetFeatureDesignFile(%q) = %q, want %q", feature, got, want)
	}
}

func TestGetFeatureAPIHandbookFile(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "design", "api-handbook.md")
	if got := GetFeatureAPIHandbookFile(feature); got != want {
		t.Errorf("GetFeatureAPIHandbookFile(%q) = %q, want %q", feature, got, want)
	}
}

func TestGetFeatureUIDesignDir(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "ui")
	if got := GetFeatureUIDesignDir(feature); got != want {
		t.Errorf("GetFeatureUIDesignDir(%q) = %q, want %q", feature, got, want)
	}
}

func TestGetFeatureUIDesignFile(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "ui", "ui-design.md")
	if got := GetFeatureUIDesignFile(feature); got != want {
		t.Errorf("GetFeatureUIDesignFile(%q) = %q, want %q", feature, got, want)
	}
}

func TestGetFeatureIndexFile(t *testing.T) {
	tests := []struct {
		feature string
		want    string
	}{
		{"my-feature", filepath.Join("docs/features", "my-feature", "tasks", "index.json")},
		{"test", filepath.Join("docs/features", "test", "tasks", "index.json")},
	}

	for _, tt := range tests {
		t.Run(tt.feature, func(t *testing.T) {
			if got := GetFeatureIndexFile(tt.feature); got != tt.want {
				t.Errorf("GetFeatureIndexFile(%q) = %q, want %q", tt.feature, got, tt.want)
			}
		})
	}
}

func TestGetFeatureTasksDir(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "tasks")
	if got := GetFeatureTasksDir(feature); got != want {
		t.Errorf("GetFeatureTasksDir(%q) = %q, want %q", feature, got, want)
	}
}

func TestGetFeatureRecordsDir(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "tasks", "records")
	if got := GetFeatureRecordsDir(feature); got != want {
		t.Errorf("GetFeatureRecordsDir(%q) = %q, want %q", feature, got, want)
	}
}

func TestGetTaskFile(t *testing.T) {
	feature := "test-feature"
	filename := "1.1.md"
	want := filepath.Join("docs/features", "test-feature", "tasks", "1.1.md")
	if got := GetTaskFile(feature, filename); got != want {
		t.Errorf("GetTaskFile(%q, %q) = %q, want %q", feature, filename, got, want)
	}
}

func TestGetRecordFile(t *testing.T) {
	feature := "test-feature"
	filename := "1.1.md"
	want := filepath.Join("docs/features", "test-feature", "tasks", "records", "1.1.md")
	if got := GetRecordFile(feature, filename); got != want {
		t.Errorf("GetRecordFile(%q, %q) = %q, want %q", feature, filename, got, want)
	}
}

func TestGetTaskStatePath(t *testing.T) {
	tests := []struct {
		name        string
		projectRoot string
		feature     string
		want        string
	}{
		{"basic", "/project", "my-feature", filepath.Join("/project", "docs", "features", "my-feature", "tasks", "process", "state.json")},
		{"nested path", "/path/to/project", "test", filepath.Join("/path/to/project", "docs", "features", "test", "tasks", "process", "state.json")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTaskStatePath(tt.projectRoot, tt.feature); got != tt.want {
				t.Errorf("GetTaskStatePath(%q, %q) = %q, want %q", tt.projectRoot, tt.feature, got, tt.want)
			}
		})
	}
}

func TestGetProcessRecordPath(t *testing.T) {
	tests := []struct {
		name        string
		projectRoot string
		feature     string
		want        string
	}{
		{"basic", "/project", "my-feature", filepath.Join("/project", "docs", "features", "my-feature", "tasks", "process", "record.json")},
		{"nested path", "/path/to/project", "test", filepath.Join("/path/to/project", "docs", "features", "test", "tasks", "process", "record.json")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetProcessRecordPath(tt.projectRoot, tt.feature); got != tt.want {
				t.Errorf("GetProcessRecordPath(%q, %q) = %q, want %q", tt.projectRoot, tt.feature, got, tt.want)
			}
		})
	}
}

func TestGetProcessDir(t *testing.T) {
	tests := []struct {
		name        string
		projectRoot string
		feature     string
		want        string
	}{
		{"basic", "/project", "my-feature", filepath.Join("/project", "docs", "features", "my-feature", "tasks", "process")},
		{"nested path", "/path/to/project", "test", filepath.Join("/path/to/project", "docs", "features", "test", "tasks", "process")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetProcessDir(tt.projectRoot, tt.feature); got != tt.want {
				t.Errorf("GetProcessDir(%q, %q) = %q, want %q", tt.projectRoot, tt.feature, got, tt.want)
			}
		})
	}
}

func TestGetProposalDir(t *testing.T) {
	slug := "my-idea"
	want := filepath.Join("docs", "proposals", "my-idea")
	if got := GetProposalDir(slug); got != want {
		t.Errorf("GetProposalDir(%q) = %q, want %q", slug, got, want)
	}
}

func TestGetProposalFile(t *testing.T) {
	slug := "my-idea"
	want := filepath.Join("docs", "proposals", "my-idea", "proposal.md")
	if got := GetProposalFile(slug); got != want {
		t.Errorf("GetProposalFile(%q) = %q, want %q", slug, got, want)
	}
}

func TestGetFeatureTestingResultsDir(t *testing.T) {
	tests := []struct {
		name        string
		featureSlug string
		want        string
	}{
		{"basic", "my-feature", filepath.Join("docs/features", "my-feature", "testing", "results")},
		{"empty", "", filepath.Join("docs/features", "", "testing", "results")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFeatureTestingResultsDir(tt.featureSlug); got != tt.want {
				t.Errorf("GetFeatureTestingResultsDir(%q) = %q, want %q", tt.featureSlug, got, tt.want)
			}
		})
	}
}

func TestGetFeatureTestCasesFile(t *testing.T) {
	tests := []struct {
		name        string
		featureSlug string
		want        string
	}{
		{"basic", "my-feature", filepath.Join("docs/features", "my-feature", "testing", "test-cases.md")},
		{"empty", "", filepath.Join("docs/features", "", "testing", "test-cases.md")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFeatureTestCasesFile(tt.featureSlug); got != tt.want {
				t.Errorf("GetFeatureTestCasesFile(%q) = %q, want %q", tt.featureSlug, got, tt.want)
			}
		})
	}
}

func TestGetE2EGraduatedMarker(t *testing.T) {
	tests := []struct {
		name        string
		projectRoot string
		featureSlug string
		want        string
	}{
		{"basic", "/project", "login", filepath.Join("/project", "tests", "e2e", ".graduated", "login")},
		{"nested", "/path/to/project", "signup-flow", filepath.Join("/path/to/project", "tests", "e2e", ".graduated", "signup-flow")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetE2EGraduatedMarker(tt.projectRoot, tt.featureSlug); got != tt.want {
				t.Errorf("GetE2EGraduatedMarker(%q, %q) = %q, want %q", tt.projectRoot, tt.featureSlug, got, tt.want)
			}
		})
	}
}

func TestGetE2ETargetDir(t *testing.T) {
	tests := []struct {
		name        string
		projectRoot string
		target      string
		want        string
	}{
		{"basic", "/project", "ui/login", filepath.Join("/project", "tests", "e2e", "ui", "login")},
		{"api target", "/project", "api/health", filepath.Join("/project", "tests", "e2e", "api", "health")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetE2ETargetDir(tt.projectRoot, tt.target); got != tt.want {
				t.Errorf("GetE2ETargetDir(%q, %q) = %q, want %q", tt.projectRoot, tt.target, got, tt.want)
			}
		})
	}
}

func TestGetE2EStagingDir(t *testing.T) {
	tests := []struct {
		name         string
		projectRoot  string
		featureSlug  string
		want         string
	}{
		{"basic", "/project", "decision-log", filepath.Join("/project", "tests", "e2e", "features", "decision-log")},
		{"nested slug", "/project", "rbac-permissions", filepath.Join("/project", "tests", "e2e", "features", "rbac-permissions")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetE2EStagingDir(tt.projectRoot, tt.featureSlug); got != tt.want {
				t.Errorf("GetE2EStagingDir(%q, %q) = %q, want %q", tt.projectRoot, tt.featureSlug, got, tt.want)
			}
		})
	}
}

func TestGetForgeStatePath(t *testing.T) {
	tests := []struct {
		name        string
		projectRoot string
		want        string
	}{
		{"basic", "/project", filepath.Join("/project", ".forge", "state.json")},
		{"windows style", "C:\\project", filepath.Join("C:\\project", ".forge", "state.json")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetForgeStatePath(tt.projectRoot); got != tt.want {
				t.Errorf("GetForgeStatePath(%q) = %q, want %q", tt.projectRoot, got, tt.want)
			}
		})
	}
}
