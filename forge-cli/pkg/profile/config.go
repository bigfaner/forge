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
	ProjectType  string   `yaml:"project-type"`
	TestProfiles []string `yaml:"test-profiles"`
	Capabilities []string `yaml:"capabilities"`
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

// ReadConfig reads the full ForgeConfig from .forge/config.yaml.
// Returns nil, nil if file doesn't exist.
func ReadConfig(projectRoot string) (*ForgeConfig, error) {
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

	return &cfg, nil
}

// ErrKeyNotFound is returned when a config key does not exist or has a zero value.
var ErrKeyNotFound = fmt.Errorf("config key not found")

// configKeyMap maps CLI key names to ForgeConfig struct field accessors.
type configKeyAccessor struct {
	scalar func(*ForgeConfig) (string, bool)
	slice  func(*ForgeConfig) ([]string, bool)
}

var configKeyAccessors = map[string]configKeyAccessor{
	"project-type": {
		scalar: func(c *ForgeConfig) (string, bool) { return c.ProjectType, c.ProjectType != "" },
	},
	"test-profiles": {
		slice: func(c *ForgeConfig) ([]string, bool) { return c.TestProfiles, len(c.TestProfiles) > 0 },
	},
	"capabilities": {
		slice: func(c *ForgeConfig) ([]string, bool) { return c.Capabilities, len(c.Capabilities) > 0 },
	},
}

// GetConfigValue returns the value for a given key from .forge/config.yaml.
// For scalar values, returns the raw string; for arrays, joins with newline.
// Returns empty string and ErrKeyNotFound if the key doesn't exist or has zero value.
func GetConfigValue(projectRoot, key string) (string, error) {
	accessor, ok := configKeyAccessors[key]
	if !ok {
		return "", ErrKeyNotFound
	}

	cfg, err := ReadConfig(projectRoot)
	if err != nil {
		return "", err
	}
	if cfg == nil {
		return "", ErrKeyNotFound
	}

	if accessor.scalar != nil {
		val, found := accessor.scalar(cfg)
		if !found {
			return "", ErrKeyNotFound
		}
		return val, nil
	}

	if accessor.slice != nil {
		vals, found := accessor.slice(cfg)
		if !found {
			return "", ErrKeyNotFound
		}
		return joinSlice(vals), nil
	}

	return "", ErrKeyNotFound
}

// joinSlice joins slice values with newline for plain-text output.
func joinSlice(vals []string) string {
	result := ""
	for i, v := range vals {
		if i > 0 {
			result += "\n"
		}
		result += v
	}
	return result
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
