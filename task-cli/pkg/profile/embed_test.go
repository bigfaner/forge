package profile

import (
	"slices"
	"testing"
)

func TestGetManifest(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"web-playwright", false},
		{"go-test", false},
		{"maestro", false},
		{"java-junit", false},
		{"rust-test", false},
		{"pytest", false},
		{"unknown-profile", true},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := GetManifest(tt.name)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetManifest(%q) expected error, got nil", tt.name)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetManifest(%q) unexpected error: %v", tt.name, err)
			}
			if len(data) == 0 {
				t.Errorf("GetManifest(%q) returned empty data", tt.name)
			}
		})
	}
}

func TestGetStrategy(t *testing.T) {
	tests := []struct {
		name    string
		kind    string
		wantErr bool
	}{
		{"go-test", "generate", false},
		{"go-test", "run", false},
		{"go-test", "graduate", false},
		{"web-playwright", "generate", false},
		{"go-test", "invalid", true},
		{"unknown", "generate", true},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/"+tt.kind, func(t *testing.T) {
			data, err := GetStrategy(tt.name, tt.kind)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(data) == 0 {
				t.Error("returned empty data")
			}
		})
	}
}

func TestGetJustfileRecipes(t *testing.T) {
	data, err := GetJustfileRecipes("go-test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("returned empty data")
	}

	_, err = GetJustfileRecipes("unknown")
	if err == nil {
		t.Error("expected error for unknown profile")
	}
}

func TestGetTemplate(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{"go-test", "test-file.go", false},
		{"go-test", "helpers.go", false},
		{"web-playwright", "helpers.ts", false},
		{"web-playwright", "nonexistent.ts", true},
		{"go-test", "nonexistent.go", true},
		{"unknown", "test-file.go", true},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/"+tt.filename, func(t *testing.T) {
			data, err := GetTemplate(tt.name, tt.filename)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(data) == 0 {
				t.Error("returned empty data")
			}
		})
	}
}

func TestListProfileTemplates(t *testing.T) {
	tests := []struct {
		name      string
		wantCount int
		wantErr   bool
	}{
		{"go-test", 2, false},
		{"web-playwright", 8, false},
		{"unknown", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templates, err := ListProfileTemplates(tt.name)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(templates) != tt.wantCount {
				t.Errorf("ListProfileTemplates(%q) = %d templates, want %d", tt.name, len(templates), tt.wantCount)
			}
		})
	}
}

func TestListEmbeddedProfiles(t *testing.T) {
	profiles := ListEmbeddedProfiles()
	if len(profiles) != 6 {
		t.Errorf("ListEmbeddedProfiles() = %d profiles, want 6", len(profiles))
	}

	expected := []string{"go-test", "java-junit", "maestro", "pytest", "rust-test", "web-playwright"}
	for _, p := range expected {
		if !slices.Contains(profiles, p) {
			t.Errorf("expected profile %q not found in list", p)
		}
	}
}
