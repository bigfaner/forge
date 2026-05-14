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
			name: "go.mod → go-test",
			setup: func(dir string) {
				mustWrite(dir, "go.mod", "module example.com/test")
			},
			expected: []string{"go-test"},
		},
		{
			name: "Cargo.toml → rust-test",
			setup: func(dir string) {
				mustWrite(dir, "Cargo.toml", "[package]\nname = \"test\"")
			},
			expected: []string{"rust-test"},
		},
		{
			name: "pom.xml → java-junit",
			setup: func(dir string) {
				mustWrite(dir, "pom.xml", "<project></project>")
			},
			expected: []string{"java-junit"},
		},
		{
			name: "build.gradle → java-junit",
			setup: func(dir string) {
				mustWrite(dir, "build.gradle", "plugins { id 'java' }")
			},
			expected: []string{"java-junit"},
		},
		{
			name: "android directory → maestro",
			setup: func(dir string) {
				mustMkdir(dir, "android")
			},
			expected: []string{"maestro"},
		},
		{
			name: "ios directory → maestro",
			setup: func(dir string) {
				mustMkdir(dir, "ios")
			},
			expected: []string{"maestro"},
		},
		{
			name: "package.json with playwright → web-playwright",
			setup: func(dir string) {
				mustWrite(dir, "package.json", `{"devDependencies":{"@playwright/test":"^1.0.0"}}`)
			},
			expected: []string{"web-playwright"},
		},
		{
			name: "playwright.config.ts → web-playwright",
			setup: func(dir string) {
				mustWrite(dir, "playwright.config.ts", "export default {}")
			},
			expected: []string{"web-playwright"},
		},
		{
			name: "package.json without playwright → web-playwright fallback",
			setup: func(dir string) {
				mustWrite(dir, "package.json", `{"dependencies":{"express":"^4.0.0"}}`)
			},
			expected: []string{"web-playwright"},
		},
		{
			name: "requirements.txt with pytest → pytest",
			setup: func(dir string) {
				mustWrite(dir, "requirements.txt", "pytest>=7.0\nrequests")
			},
			expected: []string{"pytest"},
		},
		{
			name: "pyproject.toml with pytest → pytest",
			setup: func(dir string) {
				mustWrite(dir, "pyproject.toml", "[tool.pytest.ini_options]\ntestpaths = [\"tests\"]")
			},
			expected: []string{"pytest"},
		},
		{
			name: "go.mod + package.json with playwright → go-test + web-playwright",
			setup: func(dir string) {
				mustWrite(dir, "go.mod", "module example.com/test")
				mustWrite(dir, "package.json", `{"devDependencies":{"@playwright/test":"^1.0.0"}}`)
			},
			expected: []string{"web-playwright", "go-test"},
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
