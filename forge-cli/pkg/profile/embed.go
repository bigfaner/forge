package profile

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"slices"
	"strings"
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
