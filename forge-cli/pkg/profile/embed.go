package profile

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"slices"
	"strings"
)

//go:embed all:languages
var profileFS embed.FS

const languagesDir = "languages"

// GetStrategy returns a strategy file content for the given language.
// kind must be one of: "generate", "run", "graduate".
func GetStrategy(name, kind string) ([]byte, error) {
	if err := validateLanguageName(name); err != nil {
		return nil, err
	}
	if !slices.Contains([]string{"generate", "run", "graduate"}, kind) {
		return nil, fmt.Errorf("invalid strategy kind: %s (must be generate, run, or graduate)", kind)
	}
	return profileFS.ReadFile(path.Join(languagesDir, name, kind+".md"))
}

// GetJustfileRecipes returns the justfile-recipes content for the given language.
func GetJustfileRecipes(name string) ([]byte, error) {
	if err := validateLanguageName(name); err != nil {
		return nil, err
	}
	return profileFS.ReadFile(path.Join(languagesDir, name, "justfile-recipes"))
}

// GetTemplate returns a specific template file content for the given language.
func GetTemplate(name, filename string) ([]byte, error) {
	if err := validateLanguageName(name); err != nil {
		return nil, err
	}
	return profileFS.ReadFile(path.Join(languagesDir, name, "templates", filename))
}

// ListProfileTemplates returns the template filenames available for the given language.
func ListProfileTemplates(name string) ([]string, error) {
	if err := validateLanguageName(name); err != nil {
		return nil, err
	}
	entries, err := fs.ReadDir(profileFS, path.Join(languagesDir, name, "templates"))
	if err != nil {
		return nil, fmt.Errorf("read templates for %s: %w", name, err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	slices.Sort(names)
	return names, nil
}

// ListEmbeddedProfiles returns the names of all embedded language directories.
func ListEmbeddedProfiles() []string {
	entries, err := fs.ReadDir(profileFS, languagesDir)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	slices.Sort(names)
	return names
}

// validateLanguageName checks that the language name is known.
func validateLanguageName(name string) error {
	if !slices.Contains(KnownLanguages, name) {
		return fmt.Errorf("unknown language: %s (known: %s)", name, strings.Join(KnownLanguages, ", "))
	}
	return nil
}

// ValidInterfaceTypes is the closed set of valid interface types.
var ValidInterfaceTypes = []string{
	"web-ui",
	"tui",
	"mobile-ui",
	"api",
	"cli",
}

// ValidateInterfaces checks that every value in ifaces is a known interface type.
// Returns an error listing valid values if any unknown interface is found.
func ValidateInterfaces(ifaces []string) error {
	for _, c := range ifaces {
		if !slices.Contains(ValidInterfaceTypes, c) {
			return fmt.Errorf("invalid interface: %s (valid types: %s)", c, strings.Join(ValidInterfaceTypes, ", "))
		}
	}
	return nil
}

// languageCapabilities maps each language key to its supported interface types.
// Used by ReadInterfaces as the default when config.Interfaces is empty.
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
			return nil, fmt.Errorf("unknown language: %s", name)
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
