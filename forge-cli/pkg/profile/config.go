// Package profile provides test profile resolution utilities.
package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"forge-cli/pkg/feature"

	"gopkg.in/yaml.v3"
)

// ForgeConfig represents the .forge/config.yaml structure.
type ForgeConfig struct {
	TestProfiles []string `yaml:"test-profiles"`
}

// KnownProfiles is the set of valid profile names.
var KnownProfiles = []string{
	"web-playwright",
	"go-test",
	"maestro",
	"java-junit",
	"rust-test",
	"pytest",
}

// IsKnownProfile checks whether a profile name is valid.
func IsKnownProfile(name string) bool {
	return slices.Contains(KnownProfiles, name)
}

// configPath returns the path to .forge/config.yaml.
func configPath(projectRoot string) string {
	return filepath.Join(projectRoot, feature.ForgeDir, feature.ForgeConfigFileName)
}

// ReadTestProfiles reads test-profiles from .forge/config.yaml.
// Returns empty slice (not error) if file doesn't exist or key is missing.
func ReadTestProfiles(projectRoot string) ([]string, error) {
	path := configPath(projectRoot)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg ForgeConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return cfg.TestProfiles, nil
}

// WriteTestProfiles writes test-profiles to .forge/config.yaml.
// Creates the file if it doesn't exist. Preserves other keys if the file exists.
func WriteTestProfiles(projectRoot string, profiles []string) error {
	path := configPath(projectRoot)

	// Ensure .forge/ directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create .forge dir: %w", err)
	}

	// Read existing config to preserve other keys
	var cfg ForgeConfig
	data, err := os.ReadFile(path)
	if err == nil {
		_ = yaml.Unmarshal(data, &cfg)
	}

	cfg.TestProfiles = profiles

	out, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, out, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}
