package profile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// DetectProfiles scans the project root for file signals to infer test profiles.
// Returns nil (not error) if no signals match.
// Multiple signals may return multiple profiles.
func DetectProfiles(projectRoot string) ([]string, error) {
	var detected []string
	seen := make(map[string]bool)

	add := func(name string) {
		if !seen[name] {
			seen[name] = true
			detected = append(detected, name)
		}
	}

	hasPlaywright := false
	hasPackageJSON := false

	// Check package.json for Playwright
	pkgData, err := os.ReadFile(filepath.Join(projectRoot, "package.json"))
	if err == nil {
		hasPackageJSON = true
		hasPlaywright = detectPlaywrightInPackageJSON(pkgData)
	}

	// Check for playwright.config.* files
	if !hasPlaywright {
		matches, _ := filepath.Glob(filepath.Join(projectRoot, "playwright.config.*"))
		if len(matches) > 0 {
			hasPlaywright = true
		}
	}

	if hasPlaywright {
		add("javascript")
	}

	// Go
	if fileExists(projectRoot, "go.mod") {
		add("go")
	}

	// Mobile (Maestro)
	if dirExists(projectRoot, "android") || dirExists(projectRoot, "ios") {
		add("mobile")
	}

	// Java (Maven or Gradle)
	if fileExists(projectRoot, "pom.xml") || fileExists(projectRoot, "build.gradle") || fileExists(projectRoot, "build.gradle.kts") {
		add("java")
	}

	// Rust
	if fileExists(projectRoot, "Cargo.toml") {
		add("rust")
	}

	// Python (pytest)
	if detectPytest(projectRoot) {
		add("python")
	}

	// Fallback: package.json without Playwright → javascript
	if hasPackageJSON && !hasPlaywright {
		add("javascript")
	}

	return detected, nil
}

func detectPlaywrightInPackageJSON(data []byte) bool {
	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return false
	}
	_, ok := pkg.DevDependencies["@playwright/test"]
	if !ok {
		_, ok = pkg.Dependencies["@playwright/test"]
	}
	return ok
}

func detectPytest(root string) bool {
	// Check requirements.txt
	if data, err := os.ReadFile(filepath.Join(root, "requirements.txt")); err == nil {
		if strings.Contains(string(data), "pytest") {
			return true
		}
	}

	// Check pyproject.toml
	if data, err := os.ReadFile(filepath.Join(root, "pyproject.toml")); err == nil {
		content := string(data)
		if strings.Contains(content, "pytest") {
			return true
		}
	}

	return false
}

func fileExists(root, name string) bool {
	_, err := os.Stat(filepath.Join(root, name))
	return err == nil
}

func dirExists(root, name string) bool {
	info, err := os.Stat(filepath.Join(root, name))
	return err == nil && info.IsDir()
}
