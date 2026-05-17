package profile

import (
	"slices"
	"strings"
	"testing"
)

func TestGetStrategy(t *testing.T) {
	tests := []struct {
		name    string
		kind    string
		wantErr bool
	}{
		{"go", "generate", false},
		{"go", "run", false},
		{"go", "graduate", false},
		{"javascript", "generate", false},
		{"go", "invalid", true},
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
	data, err := GetJustfileRecipes("go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("returned empty data")
	}

	_, err = GetJustfileRecipes("unknown")
	if err == nil {
		t.Error("expected error for unknown language")
	}
}

func TestGetTemplate(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{"go", "test-file.go", false},
		{"go", "helpers.go", false},
		{"javascript", "helpers.ts", false},
		{"javascript", "nonexistent.ts", true},
		{"go", "nonexistent.go", true},
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
		{"go", 2, false},
		{"javascript", 8, false},
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

	expected := []string{"go", "java", "javascript", "mobile", "python", "rust"}
	for _, p := range expected {
		if !slices.Contains(profiles, p) {
			t.Errorf("expected language %q not found in list", p)
		}
	}
}

func TestValidateInterfaces(t *testing.T) {
	t.Run("valid single interface", func(t *testing.T) {
		err := ValidateInterfaces([]string{"web-ui"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("valid multiple interfaces", func(t *testing.T) {
		err := ValidateInterfaces([]string{"web-ui", "tui", "api"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("all valid types accepted", func(t *testing.T) {
		err := ValidateInterfaces([]string{"web-ui", "tui", "mobile-ui", "api", "cli"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("invalid value rejected", func(t *testing.T) {
		err := ValidateInterfaces([]string{"web-ui", "invalid-type"})
		if err == nil {
			t.Fatal("expected error for invalid interface, got nil")
		}
		// Error message should list valid values
		for _, valid := range ValidInterfaceTypes {
			if !containsSubstring(err.Error(), valid) {
				t.Errorf("error message should mention valid type %q, got: %v", valid, err)
			}
		}
	})

	t.Run("empty input passes", func(t *testing.T) {
		err := ValidateInterfaces(nil)
		if err != nil {
			t.Fatalf("unexpected error for empty input: %v", err)
		}
	})

	t.Run("case sensitive - uppercase rejected", func(t *testing.T) {
		err := ValidateInterfaces([]string{"Web-UI"})
		if err == nil {
			t.Fatal("expected error for uppercase interface, got nil")
		}
	})

	t.Run("case sensitive - mixed case rejected", func(t *testing.T) {
		err := ValidateInterfaces([]string{"Api"})
		if err == nil {
			t.Fatal("expected error for mixed-case interface, got nil")
		}
	})

	t.Run("duplicate valid values pass", func(t *testing.T) {
		err := ValidateInterfaces([]string{"api", "api"})
		if err != nil {
			t.Fatalf("unexpected error for duplicate valid values: %v", err)
		}
	})
}

func containsSubstring(s, sub string) bool {
	return strings.Contains(s, sub)
}

func TestValidInterfaceTypes(t *testing.T) {
	expected := []string{"web-ui", "tui", "mobile-ui", "api", "cli"}
	if len(ValidInterfaceTypes) != len(expected) {
		t.Fatalf("ValidInterfaceTypes has %d entries, want %d", len(ValidInterfaceTypes), len(expected))
	}
	for _, want := range expected {
		if !slices.Contains(ValidInterfaceTypes, want) {
			t.Errorf("ValidInterfaceTypes missing %q", want)
		}
	}
}

func TestUnionLanguageInterfaces(t *testing.T) {
	t.Run("single language", func(t *testing.T) {
		ifaces, err := UnionLanguageInterfaces([]string{"go"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(ifaces) != 2 {
			t.Errorf("expected 2 interfaces, got %d: %v", len(ifaces), ifaces)
		}
	})

	t.Run("multiple languages deduplicates", func(t *testing.T) {
		// go: [api, cli], javascript: [web-ui, api]
		// union: [api, cli, web-ui] (sorted)
		ifaces, err := UnionLanguageInterfaces([]string{"go", "javascript"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// api appears in both, should be deduplicated
		expected := []string{"api", "cli", "web-ui"}
		if len(ifaces) != len(expected) {
			t.Errorf("expected %d interfaces, got %d: %v", len(expected), len(ifaces), ifaces)
		}
		for _, want := range expected {
			if !slices.Contains(ifaces, want) {
				t.Errorf("expected interface %q not found in union", want)
			}
		}
	})

	t.Run("empty languages returns empty", func(t *testing.T) {
		ifaces, err := UnionLanguageInterfaces(nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(ifaces) != 0 {
			t.Errorf("expected 0 interfaces, got %d: %v", len(ifaces), ifaces)
		}
	})
}
