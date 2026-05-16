// Package profile provides test profile resolution utilities.
package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config path constants (mirrored from feature package to avoid import cycle).
const (
	forgeDir        = ".forge"
	forgeConfigFile = "config.yaml"
)

// ModeToggle holds per-mode (quick/full) boolean flags.
// The zero-value defaults to true for both modes (backward compat).
// Use pointer types or explicit default-filling for fields that default to false.
type ModeToggle struct {
	Quick bool `yaml:"quick"`
	Full  bool `yaml:"full"`
}

// AutoConfig controls which auto-generated tasks are produced by `forge task index`.
// When the `auto` block is missing from config, all fields use defaults that match
// pre-auto-behavior behavior.
type AutoConfig struct {
	E2eTest          ModeToggle `yaml:"e2eTest"`
	ConsolidateSpecs ModeToggle `yaml:"consolidateSpecs"`
	CleanCode        ModeToggle `yaml:"cleanCode"`
	GitPush          bool       `yaml:"gitPush"`
	// raw tracks which sub-fields were explicitly present in the YAML.
	// Used by applyDefaults to distinguish "false" from "missing".
	raw map[string]map[string]bool
}

// AutoConfigDefaults returns an AutoConfig with backward-compatible defaults:
// e2eTest=true, consolidateSpecs=true, cleanCode=false, gitPush=false.
func AutoConfigDefaults() AutoConfig {
	return AutoConfig{
		E2eTest:          ModeToggle{Quick: true, Full: true},
		ConsolidateSpecs: ModeToggle{Quick: true, Full: true},
		CleanCode:        ModeToggle{Quick: false, Full: false},
		GitPush:          false,
	}
}

// IsZero returns true if the AutoConfig has all zero-value fields.
func (a AutoConfig) IsZero() bool {
	return a.E2eTest == ModeToggle{} &&
		a.ConsolidateSpecs == ModeToggle{} &&
		a.CleanCode == ModeToggle{} &&
		!a.GitPush
}

// WithDefaults returns an AutoConfig with defaults applied for any zero-value fields.
// This handles the case where BuildIndexOpts.AutoConfig was not explicitly set
// (Go zero-value for bool is false, but our defaults for e2eTest and consolidateSpecs are true).
func (a AutoConfig) WithDefaults() AutoConfig {
	if a.IsZero() {
		return AutoConfigDefaults()
	}
	// If not fully zero but some fields might need defaults
	if a.E2eTest == (ModeToggle{}) {
		a.E2eTest = ModeToggle{Quick: true, Full: true}
	}
	if a.ConsolidateSpecs == (ModeToggle{}) {
		a.ConsolidateSpecs = ModeToggle{Quick: true, Full: true}
	}
	if a.CleanCode == (ModeToggle{}) {
		a.CleanCode = ModeToggle{Quick: false, Full: false}
	}
	return a
}

// ForgeConfig represents the .forge/config.yaml structure.
type ForgeConfig struct {
	ProjectType  string      `yaml:"project-type"`
	TestProfiles []string    `yaml:"test-profiles"`
	Capabilities []string    `yaml:"capabilities"`
	Auto         *AutoConfig `yaml:"auto,omitempty"`
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
	return filepath.Join(projectRoot, forgeDir, forgeConfigFile)
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

	// Parse auto block with explicit-set tracking for default filling
	if cfg.Auto != nil {
		rawAuto, err := parseAutoRaw(data)
		if err == nil {
			cfg.Auto.raw = rawAuto
		}
		cfg.Auto.applyDefaults()
	}

	return &cfg, nil
}

// ReadAutoConfig reads the auto config block from .forge/config.yaml.
// Returns defaults when the block is missing or the file doesn't exist.
func ReadAutoConfig(projectRoot string) (AutoConfig, error) {
	cfg, err := ReadConfig(projectRoot)
	if err != nil {
		return AutoConfigDefaults(), err
	}
	if cfg == nil || cfg.Auto == nil {
		return AutoConfigDefaults(), nil
	}
	return *cfg.Auto, nil
}

// parseAutoRaw parses the raw YAML to detect which auto fields and sub-fields were present.
func parseAutoRaw(data []byte) (map[string]map[string]bool, error) {
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	// Find the "auto" mapping node
	autoNode := findMappingKey(&root, "auto")
	if autoNode == nil {
		return nil, fmt.Errorf("auto block not found")
	}

	result := make(map[string]map[string]bool)

	modeFields := []string{"e2eTest", "consolidateSpecs", "cleanCode"}
	for _, field := range modeFields {
		node := findMappingKey(autoNode, field)
		if node == nil {
			continue
		}
		result[field] = make(map[string]bool)
		if node.Kind == yaml.MappingNode {
			for i := 0; i < len(node.Content); i += 2 {
				key := node.Content[i].Value
				if key == "quick" || key == "full" {
					result[field][key] = true
				}
			}
		}
	}

	return result, nil
}

// findMappingKey finds a mapping node value by key within a YAML node tree.
func findMappingKey(node *yaml.Node, key string) *yaml.Node {
	if node == nil {
		return nil
	}
	// If the node itself is a document, look at its content
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		return findMappingKey(node.Content[0], key)
	}
	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			if node.Content[i].Value == key {
				return node.Content[i+1]
			}
		}
	}
	return nil
}

// applyDefaults fills in defaults for fields that were not explicitly set in YAML.
// The Hard Rule requires: e2eTest defaults true, consolidateSpecs defaults true, cleanCode defaults false.
func (a *AutoConfig) applyDefaults() {
	if a.raw == nil {
		// No raw tracking means all fields get defaults
		a.E2eTest = ModeToggle{Quick: true, Full: true}
		a.ConsolidateSpecs = ModeToggle{Quick: true, Full: true}
		a.CleanCode = ModeToggle{Quick: false, Full: false}
		return
	}

	applyModeDefault(&a.E2eTest, a.raw, "e2eTest", true)
	applyModeDefault(&a.ConsolidateSpecs, a.raw, "consolidateSpecs", true)
	applyModeDefault(&a.CleanCode, a.raw, "cleanCode", false)
}

// applyModeDefault sets default values for a ModeToggle field.
// If the entire field block is missing from raw, both Quick and Full get the default.
// If only one sub-key is missing, only that one gets the default.
func applyModeDefault(mt *ModeToggle, raw map[string]map[string]bool, field string, defaultVal bool) {
	fieldRaw, exists := raw[field]
	if !exists {
		// Entire field block missing: set both to default
		mt.Quick = defaultVal
		mt.Full = defaultVal
		return
	}
	// Field block exists: check individual sub-keys
	if !fieldRaw["quick"] {
		mt.Quick = defaultVal
	}
	if !fieldRaw["full"] {
		mt.Full = defaultVal
	}
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
// Supports dot-notation for nested keys (e.g. "auto.gitPush").
// Returns empty string and ErrKeyNotFound if the key doesn't exist or has zero value.
func GetConfigValue(projectRoot, key string) (string, error) {
	// Handle dot-notation auto keys
	if val, ok, err := getAutoKeyValue(projectRoot, key); ok || err != nil {
		if err != nil {
			return "", err
		}
		return val, nil
	}

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

// getAutoKeyValue handles dot-notation keys for the auto config block.
// Returns (value, true, nil) if the key was handled, ("", false, nil) if not an auto key.
func getAutoKeyValue(projectRoot, key string) (string, bool, error) {
	if key != "auto.gitPush" {
		return "", false, nil
	}

	auto, err := ReadAutoConfig(projectRoot)
	if err != nil {
		return "", true, err
	}

	return strconv.FormatBool(auto.GitPush), true, nil
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
