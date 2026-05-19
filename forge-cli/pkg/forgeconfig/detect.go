package forgeconfig

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// KnownLanguages is the set of valid language keys.
var KnownLanguages = []string{
	"go",
	"javascript",
	"python",
	"java",
	"rust",
	"mobile",
}

// IsKnownLanguage checks whether a language name is recognized.
func IsKnownLanguage(name string) bool {
	return slices.Contains(KnownLanguages, name)
}

// ReadLanguages resolves the effective languages for a project.
// Priority: config.yaml languages > auto-detect from file signals.
func ReadLanguages(projectRoot string) ([]string, error) {
	cfg, err := ReadConfig(projectRoot)
	if err != nil {
		return DetectLanguages(projectRoot)
	}
	if cfg != nil && len(cfg.Languages) > 0 {
		return cfg.Languages, nil
	}
	return DetectLanguages(projectRoot)
}

// DetectLanguages scans the project root for file signals to infer languages.
func DetectLanguages(projectRoot string) ([]string, error) {
	var detected []string
	seen := make(map[string]bool)

	add := func(lang string) {
		if !seen[lang] {
			seen[lang] = true
			detected = append(detected, lang)
		}
	}

	pkgData, err := os.ReadFile(filepath.Join(projectRoot, "package.json"))
	if err == nil {
		var pkg struct {
			Dependencies    map[string]string `json:"dependencies"`
			DevDependencies map[string]string `json:"devDependencies"`
		}
		if json.Unmarshal(pkgData, &pkg) == nil {
			_, ok := pkg.DevDependencies["@playwright/test"]
			if !ok {
				_, ok = pkg.Dependencies["@playwright/test"]
			}
			if ok {
				add("javascript")
			}
		}
	}

	if fileExists(projectRoot, "go.mod") {
		add("go")
	}

	if dirExists(projectRoot, "android") || dirExists(projectRoot, "ios") {
		add("mobile")
	}

	if fileExists(projectRoot, "pom.xml") || fileExists(projectRoot, "build.gradle") || fileExists(projectRoot, "build.gradle.kts") {
		add("java")
	}

	if fileExists(projectRoot, "Cargo.toml") {
		add("rust")
	}

	if detectPytest(projectRoot) {
		add("python")
	}

	return detected, nil
}

// ReadInterfaces resolves the effective interface types for a project.
// Priority: config.yaml interfaces > union of language capabilities.
func ReadInterfaces(projectRoot string) ([]string, error) {
	cfg, err := ReadConfig(projectRoot)
	if err != nil {
		return defaultInterfaces(projectRoot)
	}
	if cfg != nil && len(cfg.Interfaces) > 0 {
		return cfg.Interfaces, nil
	}
	return defaultInterfaces(projectRoot)
}

func defaultInterfaces(projectRoot string) ([]string, error) {
	langs, err := ReadLanguages(projectRoot)
	if err != nil {
		return nil, err
	}
	return UnionLanguageInterfaces(langs)
}

// languageCapabilities maps each language key to its supported interface types.
var languageCapabilities = map[string][]string{
	"go":         {"api", "cli"},
	"javascript": {"web-ui", "api"},
	"python":     {"api", "cli"},
	"java":       {"api", "cli"},
	"rust":       {"api", "cli"},
	"mobile":     {"mobile-ui"},
}

// UnionLanguageInterfaces returns the union of interfaces for the given languages.
func UnionLanguageInterfaces(languages []string) ([]string, error) {
	seen := make(map[string]bool)
	var result []string
	for _, name := range languages {
		caps, ok := languageCapabilities[name]
		if !ok {
			continue
		}
		for _, c := range caps {
			if !seen[c] {
				seen[c] = true
				result = append(result, c)
			}
		}
	}
	slices.Sort(result)
	return result, nil
}

func detectPytest(root string) bool {
	if data, err := os.ReadFile(filepath.Join(root, "requirements.txt")); err == nil {
		if strings.Contains(string(data), "pytest") {
			return true
		}
	}
	if data, err := os.ReadFile(filepath.Join(root, "pyproject.toml")); err == nil {
		if strings.Contains(string(data), "pytest") {
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
