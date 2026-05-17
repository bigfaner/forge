package profile

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed all:profiles
var profileFS embed.FS

const profilesDir = "profiles"

// GetManifest returns the manifest.yaml content for the given profile.
func GetManifest(name string) ([]byte, error) {
	if err := validateProfileName(name); err != nil {
		return nil, err
	}
	return profileFS.ReadFile(path.Join(profilesDir, name, "manifest.yaml"))
}

// GetStrategy returns a strategy file content for the given profile.
// kind must be one of: "generate", "run", "graduate".
func GetStrategy(name, kind string) ([]byte, error) {
	if err := validateProfileName(name); err != nil {
		return nil, err
	}
	if !slices.Contains([]string{"generate", "run", "graduate"}, kind) {
		return nil, fmt.Errorf("invalid strategy kind: %s (must be generate, run, or graduate)", kind)
	}
	return profileFS.ReadFile(path.Join(profilesDir, name, kind+".md"))
}

// GetJustfileRecipes returns the justfile-recipes content for the given profile.
func GetJustfileRecipes(name string) ([]byte, error) {
	if err := validateProfileName(name); err != nil {
		return nil, err
	}
	return profileFS.ReadFile(path.Join(profilesDir, name, "justfile-recipes"))
}

// GetTemplate returns a specific template file content for the given profile.
func GetTemplate(name, filename string) ([]byte, error) {
	if err := validateProfileName(name); err != nil {
		return nil, err
	}
	return profileFS.ReadFile(path.Join(profilesDir, name, "templates", filename))
}

// ListProfileTemplates returns the template filenames available for the given profile.
func ListProfileTemplates(name string) ([]string, error) {
	if err := validateProfileName(name); err != nil {
		return nil, err
	}
	entries, err := fs.ReadDir(profileFS, path.Join(profilesDir, name, "templates"))
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

// ListEmbeddedProfiles returns the names of all embedded profiles.
func ListEmbeddedProfiles() []string {
	entries, err := fs.ReadDir(profileFS, profilesDir)
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

// validateProfileName checks that the profile name is known.
func validateProfileName(name string) error {
	if !slices.Contains(KnownProfiles, name) {
		return fmt.Errorf("unknown profile: %s (known: %s)", name, strings.Join(KnownProfiles, ", "))
	}
	return nil
}

// profileManifest represents the parsed manifest.yaml structure.
type profileManifest struct {
	Capabilities []string `yaml:"capabilities"`
}

// GetProfileInterfaces returns the interfaces for a given profile.
func GetProfileInterfaces(name string) ([]string, error) {
	data, err := GetManifest(name)
	if err != nil {
		return nil, err
	}
	var manifest profileManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("parse manifest for %s: %w", name, err)
	}
	return manifest.Capabilities, nil
}

// ValidInterfaceTypes is the closed set of valid interface types.
// Sourced from all profile manifests under pkg/profile/profiles/.
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

// languageCapabilities maps each profile name to its supported interface types.
// Used by ReadInterfaces as the default when config.Interfaces is empty.
var languageCapabilities = map[string][]string{
	"go-test":        {"api", "cli", "tui"},
	"web-playwright": {"web-ui", "api", "cli"},
	"pytest":         {"api", "cli"},
	"java-junit":     {"api", "cli"},
	"rust-test":      {"api", "cli"},
	"maestro":        {"mobile-ui"},
}

// UnionLanguageInterfaces returns the union of interfaces for the given profiles.
func UnionLanguageInterfaces(profiles []string) ([]string, error) {
	seen := make(map[string]bool)
	var result []string
	for _, name := range profiles {
		caps, ok := languageCapabilities[name]
		if !ok {
			// Fallback to manifest for unknown profiles
			var err error
			caps, err = GetProfileInterfaces(name)
			if err != nil {
				return nil, err
			}
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
