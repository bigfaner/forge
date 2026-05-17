package profile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Language represents a detected programming language key.
type Language string

// Language constants for type-safe language keys.
const (
	LanguageGo         Language = "go"
	LanguageJavaScript Language = "javascript"
	LanguagePython     Language = "python"
	LanguageJava       Language = "java"
	LanguageRust       Language = "rust"
	LanguageMobile     Language = "mobile"
)

// DetectLanguages scans the project root for file signals to infer languages.
// Returns nil (not error) if no signals match.
// Multiple signals may return multiple languages.
func DetectLanguages(projectRoot string) ([]Language, error) {
	var detected []Language
	seen := make(map[Language]bool)

	add := func(lang Language) {
		if !seen[lang] {
			seen[lang] = true
			detected = append(detected, lang)
		}
	}

	// Check package.json for @playwright/test dependency
	pkgData, err := os.ReadFile(filepath.Join(projectRoot, "package.json"))
	if err == nil && detectPlaywrightInPackageJSON(pkgData) {
		add(LanguageJavaScript)
	}

	// Go
	if fileExists(projectRoot, "go.mod") {
		add(LanguageGo)
	}

	// Mobile (Maestro)
	if dirExists(projectRoot, "android") || dirExists(projectRoot, "ios") {
		add(LanguageMobile)
	}

	// Java (Maven or Gradle)
	if fileExists(projectRoot, "pom.xml") || fileExists(projectRoot, "build.gradle") || fileExists(projectRoot, "build.gradle.kts") {
		add(LanguageJava)
	}

	// Rust
	if fileExists(projectRoot, "Cargo.toml") {
		add(LanguageRust)
	}

	// Python (pytest)
	if detectPytest(projectRoot) {
		add(LanguagePython)
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
