package profile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectProfiles(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(dir string)
		expected []string
	}{
		{
			name:     "empty directory",
			setup:    func(_ string) {},
			expected: nil,
		},
		{
			name: "go.mod → go",
			setup: func(dir string) {
				mustWrite(dir, "go.mod", "module example.com/test")
			},
			expected: []string{"go"},
		},
		{
			name: "Cargo.toml → rust",
			setup: func(dir string) {
				mustWrite(dir, "Cargo.toml", "[package]\nname = \"test\"")
			},
			expected: []string{"rust"},
		},
		{
			name: "pom.xml → java",
			setup: func(dir string) {
				mustWrite(dir, "pom.xml", "<project></project>")
			},
			expected: []string{"java"},
		},
		{
			name: "build.gradle → java",
			setup: func(dir string) {
				mustWrite(dir, "build.gradle", "plugins { id 'java' }")
			},
			expected: []string{"java"},
		},
		{
			name: "android directory → mobile",
			setup: func(dir string) {
				mustMkdir(dir, "android")
			},
			expected: []string{"mobile"},
		},
		{
			name: "ios directory → mobile",
			setup: func(dir string) {
				mustMkdir(dir, "ios")
			},
			expected: []string{"mobile"},
		},
		{
			name: "package.json with playwright → javascript",
			setup: func(dir string) {
				mustWrite(dir, "package.json", `{"devDependencies":{"@playwright/test":"^1.0.0"}}`)
			},
			expected: []string{"javascript"},
		},
		{
			name: "playwright.config.ts → javascript",
			setup: func(dir string) {
				mustWrite(dir, "playwright.config.ts", "export default {}")
			},
			expected: []string{"javascript"},
		},
		{
			name: "package.json without playwright → javascript fallback",
			setup: func(dir string) {
				mustWrite(dir, "package.json", `{"dependencies":{"express":"^4.0.0"}}`)
			},
			expected: []string{"javascript"},
		},
		{
			name: "requirements.txt with pytest → python",
			setup: func(dir string) {
				mustWrite(dir, "requirements.txt", "pytest>=7.0\nrequests")
			},
			expected: []string{"python"},
		},
		{
			name: "pyproject.toml with pytest → python",
			setup: func(dir string) {
				mustWrite(dir, "pyproject.toml", "[tool.pytest.ini_options]\ntestpaths = [\"tests\"]")
			},
			expected: []string{"python"},
		},
		{
			name: "go.mod + package.json with playwright → javascript + go",
			setup: func(dir string) {
				mustWrite(dir, "go.mod", "module example.com/test")
				mustWrite(dir, "package.json", `{"devDependencies":{"@playwright/test":"^1.0.0"}}`)
			},
			expected: []string{"javascript", "go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			tt.setup(dir)

			got, err := DetectProfiles(dir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
			for i, v := range got {
				if v != tt.expected[i] {
					t.Fatalf("expected %v, got %v", tt.expected, got)
				}
			}
		})
	}
}

func mustWrite(dir, name, content string) {
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		panic(err)
	}
}

func mustMkdir(dir, name string) {
	if err := os.Mkdir(filepath.Join(dir, name), 0o755); err != nil {
		panic(err)
	}
}
