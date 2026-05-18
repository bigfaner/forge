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
// e2eTest: quick=false, full=true; consolidateSpecs: quick=true, full=true;
// cleanCode=false, gitPush=false.
func AutoConfigDefaults() AutoConfig {
	return AutoConfig{
		E2eTest:          ModeToggle{Quick: false, Full: true},
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
func (a AutoConfig) WithDefaults() AutoConfig {
	if a.IsZero() {
		return AutoConfigDefaults()
	}
	d := AutoConfigDefaults()
	if a.E2eTest == (ModeToggle{}) {
		a.E2eTest = d.E2eTest
	}
	if a.ConsolidateSpecs == (ModeToggle{}) {
		a.ConsolidateSpecs = d.ConsolidateSpecs
	}
	if a.CleanCode == (ModeToggle{}) {
		a.CleanCode = d.CleanCode
	}
	return a
}

// WorktreeConfig controls worktree creation behavior.
type WorktreeConfig struct {
	SourceBranch string   `yaml:"source-branch"`
	CopyFiles    []string `yaml:"copy-files"`
}

// ForgeConfig represents the .forge/config.yaml structure.
type ForgeConfig struct {
	ProjectType   string          `yaml:"project-type"`
	Interfaces    []string        `yaml:"interfaces"`
	Languages     []string        `yaml:"languages"`
	TestFramework string          `yaml:"test-framework,omitempty"`
	TestCommand   string          `yaml:"test-command,omitempty"`
	Auto          *AutoConfig     `yaml:"auto,omitempty"`
	Worktree      *WorktreeConfig `yaml:"worktree,omitempty"`
}

// KnownLanguages is the set of valid language keys.
var KnownLanguages = []string{
	"go",
	"javascript",
	"mobile",
	"java",
	"rust",
	"python",
}

// IsKnownLanguage checks whether a language key is valid.
func IsKnownLanguage(name string) bool {
	return slices.Contains(KnownLanguages, name)
}

// ReadLanguages reads languages from .forge/config.yaml.
// Returns config.Languages if set, otherwise auto-detects via DetectLanguages.
// Returns empty slice (not error) if file doesn't exist or key is missing and detection finds nothing.
func ReadLanguages(projectRoot string) ([]string, error) {
	path := configPath(projectRoot)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return detectLanguagesAsStrings(projectRoot)
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg ForgeConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if len(cfg.Languages) > 0 {
		return cfg.Languages, nil
	}

	return detectLanguagesAsStrings(projectRoot)
}

// ReadInterfaces reads interfaces from .forge/config.yaml.
// Returns config.Interfaces if set, otherwise defaults to union of all
// detected languages' capabilities via the languageCapabilities map in embed.go.
func ReadInterfaces(projectRoot string) ([]string, error) {
	path := configPath(projectRoot)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultInterfaces(projectRoot)
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg ForgeConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if len(cfg.Interfaces) > 0 {
		return cfg.Interfaces, nil
	}

	return defaultInterfaces(projectRoot)
}

// defaultInterfaces detects profiles and returns the union of their interfaces.
func defaultInterfaces(projectRoot string) ([]string, error) {
	langs, err := DetectLanguages(projectRoot)
	if err != nil {
		return nil, err
	}
	if len(langs) == 0 {
		return nil, nil
	}
	strs, err := languagesToStrings(langs)
	if err != nil {
		return nil, err
	}
	return UnionLanguageInterfaces(strs)
}

// configPath returns the path to .forge/config.yaml.
func configPath(projectRoot string) string {
	return filepath.Join(projectRoot, forgeDir, forgeConfigFile)
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
func (a *AutoConfig) applyDefaults() {
	d := AutoConfigDefaults()
	if a.raw == nil {
		a.E2eTest = d.E2eTest
		a.ConsolidateSpecs = d.ConsolidateSpecs
		a.CleanCode = d.CleanCode
		return
	}

	applyModeDefault(&a.E2eTest, a.raw, "e2eTest", d.E2eTest)
	applyModeDefault(&a.ConsolidateSpecs, a.raw, "consolidateSpecs", d.ConsolidateSpecs)
	applyModeDefault(&a.CleanCode, a.raw, "cleanCode", d.CleanCode)
}

// applyModeDefault sets default values for a ModeToggle field using per-mode defaults.
func applyModeDefault(mt *ModeToggle, raw map[string]map[string]bool, field string, defaults ModeToggle) {
	fieldRaw, exists := raw[field]
	if !exists {
		mt.Quick = defaults.Quick
		mt.Full = defaults.Full
		return
	}
	if !fieldRaw["quick"] {
		mt.Quick = defaults.Quick
	}
	if !fieldRaw["full"] {
		mt.Full = defaults.Full
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
	"interfaces": {
		slice: func(c *ForgeConfig) ([]string, bool) { return c.Interfaces, len(c.Interfaces) > 0 },
	},
	"languages": {
		slice: func(c *ForgeConfig) ([]string, bool) { return c.Languages, len(c.Languages) > 0 },
	},
	"test-framework": {
		scalar: func(c *ForgeConfig) (string, bool) { return c.TestFramework, c.TestFramework != "" },
	},
	"test-command": {
		scalar: func(c *ForgeConfig) (string, bool) { return c.TestCommand, c.TestCommand != "" },
	},
}

// GetConfigValue returns the value for a given key from .forge/config.yaml.
// For scalar values, returns the raw string; for arrays, joins with newline.
// Supports dot-notation for nested keys (e.g. "auto.gitPush", "worktree.source-branch").
// Returns empty string and ErrKeyNotFound if the key doesn't exist or has zero value.
func GetConfigValue(projectRoot, key string) (string, error) {
	// Handle dot-notation auto keys
	if val, ok, err := getAutoKeyValue(projectRoot, key); ok || err != nil {
		if err != nil {
			return "", err
		}
		return val, nil
	}

	// Handle dot-notation worktree keys
	if val, ok, err := getWorktreeKeyValue(projectRoot, key); ok || err != nil {
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

// getWorktreeKeyValue handles dot-notation keys for the worktree config block.
// Returns (value, true, nil) if the key was handled, ("", false, nil) if not a worktree key.
func getWorktreeKeyValue(projectRoot, key string) (string, bool, error) {
	if key != "worktree.source-branch" && key != "worktree.copy-files" {
		return "", false, nil
	}

	cfg, err := ReadConfig(projectRoot)
	if err != nil {
		return "", true, err
	}
	if cfg == nil || cfg.Worktree == nil {
		return "", true, ErrKeyNotFound
	}

	switch key {
	case "worktree.source-branch":
		if cfg.Worktree.SourceBranch == "" {
			return "", true, ErrKeyNotFound
		}
		return cfg.Worktree.SourceBranch, true, nil
	case "worktree.copy-files":
		if len(cfg.Worktree.CopyFiles) == 0 {
			return "", true, ErrKeyNotFound
		}
		return joinSlice(cfg.Worktree.CopyFiles), true, nil
	}

	return "", true, ErrKeyNotFound
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

// WriteLanguages writes languages to .forge/config.yaml.
// Creates the file if it doesn't exist. Preserves other keys if the file exists.
func WriteLanguages(projectRoot string, languages []string) error {
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

	cfg.Languages = languages

	out, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, out, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

// languagesToStrings converts []Language to []string for compatibility with
// config and interface resolution functions that operate on string keys.
// Preserves nil: if langs is nil, returns nil (not empty slice).
func languagesToStrings(langs []Language) ([]string, error) {
	if langs == nil {
		return nil, nil
	}
	result := make([]string, len(langs))
	for i, l := range langs {
		result[i] = string(l)
	}
	return result, nil
}

// detectLanguagesAsStrings calls DetectLanguages and converts the result to []string.
func detectLanguagesAsStrings(projectRoot string) ([]string, error) {
	langs, err := DetectLanguages(projectRoot)
	if err != nil {
		return nil, err
	}
	return languagesToStrings(langs)
}
