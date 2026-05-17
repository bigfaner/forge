package profile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectLanguages(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(dir string)
		expected []Language
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
			expected: []Language{LanguageGo},
		},
		{
			name: "Cargo.toml → rust",
			setup: func(dir string) {
				mustWrite(dir, "Cargo.toml", "[package]\nname = \"test\"")
			},
			expected: []Language{LanguageRust},
		},
		{
			name: "pom.xml → java",
			setup: func(dir string) {
				mustWrite(dir, "pom.xml", "<project></project>")
			},
			expected: []Language{LanguageJava},
		},
		{
			name: "build.gradle → java",
			setup: func(dir string) {
				mustWrite(dir, "build.gradle", "plugins { id 'java' }")
			},
			expected: []Language{LanguageJava},
		},
		{
			name: "android directory → mobile",
			setup: func(dir string) {
				mustMkdir(dir, "android")
			},
			expected: []Language{LanguageMobile},
		},
		{
			name: "ios directory → mobile",
			setup: func(dir string) {
				mustMkdir(dir, "ios")
			},
			expected: []Language{LanguageMobile},
		},
		{
			name: "package.json with playwright devDep → javascript",
			setup: func(dir string) {
				mustWrite(dir, "package.json", `{"devDependencies":{"@playwright/test":"^1.0.0"}}`)
			},
			expected: []Language{LanguageJavaScript},
		},
		{
			name: "package.json with playwright dep → javascript",
			setup: func(dir string) {
				mustWrite(dir, "package.json", `{"dependencies":{"@playwright/test":"^1.0.0"}}`)
			},
			expected: []Language{LanguageJavaScript},
		},
		{
			name: "package.json without playwright → no detection",
			setup: func(dir string) {
				mustWrite(dir, "package.json", `{"dependencies":{"express":"^4.0.0"}}`)
			},
			expected: nil,
		},
		{
			name: "requirements.txt with pytest → python",
			setup: func(dir string) {
				mustWrite(dir, "requirements.txt", "pytest>=7.0\nrequests")
			},
			expected: []Language{LanguagePython},
		},
		{
			name: "pyproject.toml with pytest → python",
			setup: func(dir string) {
				mustWrite(dir, "pyproject.toml", "[tool.pytest.ini_options]\ntestpaths = [\"tests\"]")
			},
			expected: []Language{LanguagePython},
		},
		{
			name: "go.mod + package.json with playwright → javascript + go",
			setup: func(dir string) {
				mustWrite(dir, "go.mod", "module example.com/test")
				mustWrite(dir, "package.json", `{"devDependencies":{"@playwright/test":"^1.0.0"}}`)
			},
			expected: []Language{LanguageJavaScript, LanguageGo},
		},
		{
			name: "playwright.config.ts without package.json playwright dep → no javascript",
			setup: func(dir string) {
				mustWrite(dir, "playwright.config.ts", "export default {}")
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			tt.setup(dir)

			got, err := DetectLanguages(dir)
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
