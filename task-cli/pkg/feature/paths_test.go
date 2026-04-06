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

func TestGetFeaturePRDFile(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "prd.md")
	if got := GetFeaturePRDFile(feature); got != want {
		t.Errorf("GetFeaturePRDFile(%q) = %q, want %q", feature, got, want)
	}
}

func TestGetFeatureDesignFile(t *testing.T) {
	feature := "test-feature"
	want := filepath.Join("docs/features", "test-feature", "design.md")
	if got := GetFeatureDesignFile(feature); got != want {
		t.Errorf("GetFeatureDesignFile(%q) = %q, want %q", feature, got, want)
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
