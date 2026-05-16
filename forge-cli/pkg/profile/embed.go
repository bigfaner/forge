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

// GetProfileCapabilities returns the capabilities for a given profile.
func GetProfileCapabilities(name string) ([]string, error) {
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

// ValidTestTypes is the closed set of valid test-type capabilities.
// Sourced from all profile manifests under pkg/profile/profiles/.
var ValidTestTypes = []string{
	"web-ui",
	"tui",
	"mobile-ui",
	"api",
	"cli",
}

// ValidateCapabilities checks that every value in caps is a known test-type capability.
// Returns an error listing valid values if any unknown capability is found.
func ValidateCapabilities(caps []string) error {
	for _, c := range caps {
		if !slices.Contains(ValidTestTypes, c) {
			return fmt.Errorf("invalid capability: %s (valid types: %s)", c, strings.Join(ValidTestTypes, ", "))
		}
	}
	return nil
}

// UnionCapabilities returns the union of capabilities from the given profiles.
func UnionCapabilities(profileNames []string) ([]string, error) {
	seen := make(map[string]bool)
	var result []string
	for _, name := range profileNames {
		caps, err := GetProfileCapabilities(name)
		if err != nil {
			return nil, err
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
