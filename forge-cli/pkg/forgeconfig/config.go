// Package forgeconfig provides types and functions for reading/writing
// the .forge/config.yaml file. This package extracts only the retained
// config types from the legacy profile package: auto and worktree blocks.
package forgeconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	forgeDir        = ".forge"
	forgeConfigFile = "config.yaml"
)

// ModeToggle holds per-mode (quick/full) boolean flags.
// The zero-value defaults to true for both modes (backward compat).
type ModeToggle struct {
	Quick bool `yaml:"quick"`
	Full  bool `yaml:"full"`
}

// AutoConfig controls which auto-generated tasks are produced by `forge task index`.
// When the `auto` block is missing from config, all fields use defaults that match
// pre-auto-behavior.
type AutoConfig struct {
	E2eTest          ModeToggle `yaml:"e2eTest"`
	ConsolidateSpecs ModeToggle `yaml:"consolidateSpecs"`
	CleanCode        ModeToggle `yaml:"cleanCode"`
	Validation       ModeToggle `yaml:"validation"`
	RunTasks         ModeToggle `yaml:"runTasks"`
	GitPush          bool       `yaml:"gitPush"`
	KnowledgeSave    ModeToggle `yaml:"knowledgeSave"`
	// raw tracks which sub-fields were explicitly present in the YAML.
	// Used by applyDefaults to distinguish "false" from "missing".
	raw map[string]map[string]bool
}

// AutoConfigDefaults returns an AutoConfig with backward-compatible defaults:
// e2eTest: quick=false, full=true; consolidateSpecs: quick=true, full=true;
// cleanCode=false, validation=false, gitPush=false.
func AutoConfigDefaults() AutoConfig {
	return AutoConfig{
		E2eTest:          ModeToggle{Quick: false, Full: true},
		ConsolidateSpecs: ModeToggle{Quick: true, Full: true},
		CleanCode:        ModeToggle{Quick: false, Full: false},
		Validation:       ModeToggle{Quick: false, Full: false},
		RunTasks:         ModeToggle{Quick: true, Full: false},
		GitPush:          false,
		KnowledgeSave:    ModeToggle{Quick: true, Full: false},
	}
}

// IsZero returns true if the AutoConfig has all zero-value fields.
func (a AutoConfig) IsZero() bool {
	return a.E2eTest == ModeToggle{} &&
		a.ConsolidateSpecs == ModeToggle{} &&
		a.CleanCode == ModeToggle{} &&
		a.Validation == ModeToggle{} &&
		a.RunTasks == ModeToggle{} &&
		!a.GitPush &&
		a.KnowledgeSave == ModeToggle{}
}

// WithDefaults returns an AutoConfig with defaults applied for any zero-value fields.
// IMPORTANT: This cannot distinguish "user explicitly set ModeToggle{false, false}" from
// "field was never set" because both equal ModeToggle{}. Use ReadConfig -> applyDefaults()
// (which tracks raw YAML fields) for proper per-field defaults. This function only handles
// the all-zero case (no config loaded at all).
func (a AutoConfig) WithDefaults() AutoConfig {
	if a.IsZero() {
		return AutoConfigDefaults()
	}
	return a
}

// WorktreeConfig controls worktree creation behavior.
type WorktreeConfig struct {
	SourceBranch string   `yaml:"source-branch"`
	CopyFiles    []string `yaml:"copy-files"`
}

// CoverageStrategy defines the coverage strategy for a single task type.
// Two strategy types are supported:
//   - "percentage": target a specific coverage percentage
//   - "maintain": keep existing coverage, don't add new tests
type CoverageStrategy struct {
	Type       string `yaml:"type"`
	Percentage *int   `yaml:"percentage,omitempty"`
}

// CoverageConfig holds per-task-type coverage strategies.
type CoverageConfig struct {
	ByType map[string]CoverageStrategy `yaml:",inline"`
}

// CoverageConfigDefaults returns built-in default coverage strategies.
// Returns a fresh map each time to prevent mutation issues.
func CoverageConfigDefaults() CoverageConfig {
	feature := 80
	enhancement := 60
	fix := 60
	return CoverageConfig{
		ByType: map[string]CoverageStrategy{
			"coding.feature":     {Type: "percentage", Percentage: &feature},
			"coding.enhancement": {Type: "percentage", Percentage: &enhancement},
			"coding.fix":         {Type: "percentage", Percentage: &fix},
			"coding.refactor":    {Type: "maintain"},
			"coding.cleanup":     {Type: "maintain"},
			"coding.clean":       {Type: "maintain"},
		},
	}
}

// ReadCoverageConfig reads the coverage config block from .forge/config.yaml.
// Returns defaults when the block is missing or the file doesn't exist.
// User-provided values are merged on top of defaults.
func ReadCoverageConfig(projectRoot string) (CoverageConfig, error) {
	cfg, err := ReadConfig(projectRoot)
	if err != nil {
		return CoverageConfigDefaults(), err
	}
	if cfg == nil || cfg.Coverage == nil {
		return CoverageConfigDefaults(), nil
	}
	// Merge: start with defaults, overlay user-provided
	result := CoverageConfigDefaults()
	for k, v := range cfg.Coverage.ByType {
		result.ByType[k] = v
	}
	return result, nil
}

// Config represents the .forge/config.yaml structure.
type Config struct {
	Auto          *AutoConfig     `yaml:"auto,omitempty"`
	Worktree      *WorktreeConfig `yaml:"worktree,omitempty"`
	Coverage      *CoverageConfig `yaml:"coverage,omitempty"`
	TestFramework string          `yaml:"test-framework,omitempty"`
	Languages     []string        `yaml:"languages,omitempty"`
	Interfaces    []string        `yaml:"interfaces,omitempty"`
}

// configPath returns the path to .forge/config.yaml.
func configPath(projectRoot string) string {
	return filepath.Join(projectRoot, forgeDir, forgeConfigFile)
}

// ReadConfig reads the Config from .forge/config.yaml.
// Returns nil, nil if file doesn't exist.
// Unknown fields in the YAML are silently ignored.
func ReadConfig(projectRoot string) (*Config, error) {
	path := configPath(projectRoot)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
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
// Returns value type (AutoConfig), not pointer.
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

	autoNode := findMappingKey(&root, "auto")
	if autoNode == nil {
		return nil, fmt.Errorf("auto block not found")
	}

	result := make(map[string]map[string]bool)

	modeFields := []string{"e2eTest", "consolidateSpecs", "cleanCode", "validation", "runTasks", "knowledgeSave"}
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
		a.RunTasks = d.RunTasks
		a.KnowledgeSave = d.KnowledgeSave
		return
	}

	applyModeDefault(&a.E2eTest, a.raw, "e2eTest", d.E2eTest)
	applyModeDefault(&a.ConsolidateSpecs, a.raw, "consolidateSpecs", d.ConsolidateSpecs)
	applyModeDefault(&a.CleanCode, a.raw, "cleanCode", d.CleanCode)
	applyModeDefault(&a.Validation, a.raw, "validation", d.Validation)
	applyModeDefault(&a.RunTasks, a.raw, "runTasks", d.RunTasks)
	applyModeDefault(&a.KnowledgeSave, a.raw, "knowledgeSave", d.KnowledgeSave)
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

// errKeyNotFound is returned when a config key does not exist or has a zero value.
var errKeyNotFound = fmt.Errorf("config key not found")

// GetConfigValue returns the value for a given key from .forge/config.yaml.
// For scalar values, returns the raw string; for arrays, joins with newline.
// Supports dot-notation for nested keys (e.g. "auto.gitPush", "worktree.source-branch", "coverage.coding.feature").
// Also supports top-level keys: "test-framework".
// Returns empty string and errKeyNotFound if the key doesn't exist or has zero value.
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

	// Handle coverage.* keys
	if val, ok, err := getCoverageKeyValue(projectRoot, key); ok || err != nil {
		if err != nil {
			return "", err
		}
		return val, nil
	}

	// Handle top-level scalar keys
	if key == "test-framework" {
		cfg, err := ReadConfig(projectRoot)
		if err != nil {
			return "", err
		}
		if cfg == nil {
			return "", errKeyNotFound
		}
		if cfg.TestFramework == "" {
			return "", errKeyNotFound
		}
		return cfg.TestFramework, nil
	}

	return "", errKeyNotFound
}

// autoModeField returns the ModeToggle pointer for a given auto field name.
func autoModeField(a *AutoConfig, field string) *ModeToggle {
	switch field {
	case "e2eTest":
		return &a.E2eTest
	case "consolidateSpecs":
		return &a.ConsolidateSpecs
	case "cleanCode":
		return &a.CleanCode
	case "validation":
		return &a.Validation
	case "runTasks":
		return &a.RunTasks
	case "knowledgeSave":
		return &a.KnowledgeSave
	}
	return nil
}

// SetConfigValue sets a config value for a given dot-notation key in .forge/config.yaml.
// Supported key patterns:
//   - auto.{field}            (sets both quick and full of a ModeToggle)
//   - auto.{field}.quick      (sets quick field of a ModeToggle)
//   - auto.{field}.full       (sets full field of a ModeToggle)
//   - auto.gitPush            (sets bool)
//   - worktree.source-branch  (sets string)
//   - test-framework          (sets string)
//   - coverage.{task-type}    (sets coverage strategy)
//
// Returns an error for unknown keys or invalid values.
func SetConfigValue(projectRoot, key, value string) error {
	switch {
	case strings.HasPrefix(key, "auto."):
		return setAutoConfigValue(projectRoot, key, value)
	case strings.HasPrefix(key, "worktree."):
		return setWorktreeConfigValue(projectRoot, key, value)
	case key == "test-framework":
		cfg, err := readOrCreateConfig(projectRoot)
		if err != nil {
			return err
		}
		cfg.TestFramework = value
		return writeConfig(projectRoot, cfg)
	case strings.HasPrefix(key, "coverage."):
		return setCoverageConfigValue(projectRoot, key, value)
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}
}

// readOrCreateConfig reads config or returns an empty Config if file doesn't exist.
func readOrCreateConfig(projectRoot string) (*Config, error) {
	cfg, err := ReadConfig(projectRoot)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = &Config{}
	}
	if cfg.Auto == nil {
		cfg.Auto = &AutoConfig{}
	}
	return cfg, nil
}

// setAutoConfigValue sets a value in the auto config block by dot-notation key.
func setAutoConfigValue(projectRoot, key, value string) error {
	cfg, err := readOrCreateConfig(projectRoot)
	if err != nil {
		return err
	}

	// auto.gitPush
	if key == "auto.gitPush" {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid bool value for %s: %s", key, value)
		}
		cfg.Auto.GitPush = b
		return writeConfig(projectRoot, cfg)
	}

	// auto.{field} or auto.{field}.{subfield}
	rest := strings.TrimPrefix(key, "auto.")
	parts := strings.SplitN(rest, ".", 2)

	field := parts[0]
	mt := autoModeField(cfg.Auto, field)
	if mt == nil {
		return fmt.Errorf("unknown config key: %s", key)
	}

	if len(parts) == 1 {
		// auto.{field} — set both quick and full to the same bool value
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid bool value for %s: %s", key, value)
		}
		mt.Quick = b
		mt.Full = b
		return writeConfig(projectRoot, cfg)
	}

	// auto.{field}.{subfield}
	subfield := parts[1]
	b, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("invalid bool value for %s: %s", key, value)
	}
	switch subfield {
	case "quick":
		mt.Quick = b
	case "full":
		mt.Full = b
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	return writeConfig(projectRoot, cfg)
}

// setWorktreeConfigValue sets a value in the worktree config block by dot-notation key.
func setWorktreeConfigValue(projectRoot, key, value string) error {
	cfg, err := readOrCreateConfig(projectRoot)
	if err != nil {
		return err
	}

	switch key {
	case "worktree.source-branch":
		if cfg.Worktree == nil {
			cfg.Worktree = &WorktreeConfig{}
		}
		cfg.Worktree.SourceBranch = value
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	return writeConfig(projectRoot, cfg)
}

// setCoverageConfigValue sets a coverage strategy.
func setCoverageConfigValue(projectRoot, key, value string) error {
	cfg, err := readOrCreateConfig(projectRoot)
	if err != nil {
		return err
	}

	taskType := strings.TrimPrefix(key, "coverage.")
	if taskType == "" {
		return fmt.Errorf("unknown config key: %s", key)
	}

	if cfg.Coverage == nil {
		cfg.Coverage = &CoverageConfig{ByType: make(map[string]CoverageStrategy)}
	}

	// Parse value: either "maintain" or a percentage number
	pct, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid coverage value for %s: %s (expected percentage number)", key, value)
	}

	cfg.Coverage.ByType[taskType] = CoverageStrategy{
		Type:       "percentage",
		Percentage: &pct,
	}

	return writeConfig(projectRoot, cfg)
}

// getAutoKeyValue handles dot-notation keys for the auto config block.
func getAutoKeyValue(projectRoot, key string) (string, bool, error) {
	// auto.gitPush is a simple bool
	if key == "auto.gitPush" {
		auto, err := ReadAutoConfig(projectRoot)
		if err != nil {
			return "", true, err
		}
		return strconv.FormatBool(auto.GitPush), true, nil
	}

	// auto.{field} or auto.{field}.{subfield}
	rest, ok := strings.CutPrefix(key, "auto.")
	if !ok {
		return "", false, nil
	}

	var field, subfield string
	if idx := strings.Index(rest, "."); idx >= 0 {
		field = rest[:idx]
		subfield = rest[idx+1:]
	} else {
		field = rest
	}

	auto, err := ReadAutoConfig(projectRoot)
	if err != nil {
		return "", true, err
	}

	mt := autoModeField(&auto, field)
	if mt == nil {
		return "", false, nil
	}

	if subfield == "" {
		return fmt.Sprintf("quick:%v full:%v", mt.Quick, mt.Full), true, nil
	}

	switch subfield {
	case "quick":
		return strconv.FormatBool(mt.Quick), true, nil
	case "full":
		return strconv.FormatBool(mt.Full), true, nil
	}

	return "", false, nil
}

// getWorktreeKeyValue handles dot-notation keys for the worktree config block.
func getWorktreeKeyValue(projectRoot, key string) (string, bool, error) {
	if key != "worktree.source-branch" && key != "worktree.copy-files" {
		return "", false, nil
	}

	cfg, err := ReadConfig(projectRoot)
	if err != nil {
		return "", true, err
	}
	if cfg == nil || cfg.Worktree == nil {
		return "", true, errKeyNotFound
	}

	switch key {
	case "worktree.source-branch":
		if cfg.Worktree.SourceBranch == "" {
			return "", true, errKeyNotFound
		}
		return cfg.Worktree.SourceBranch, true, nil
	case "worktree.copy-files":
		if len(cfg.Worktree.CopyFiles) == 0 {
			return "", true, errKeyNotFound
		}
		return joinSlice(cfg.Worktree.CopyFiles), true, nil
	}

	return "", true, errKeyNotFound
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

// getCoverageKeyValue handles dot-notation keys for the coverage config block.
// Key format: "coverage.<task-type>" (e.g. "coverage.coding.feature").
// Returns the strategy type or percentage value as a string.
func getCoverageKeyValue(projectRoot, key string) (string, bool, error) {
	if !strings.HasPrefix(key, "coverage.") {
		return "", false, nil
	}

	taskType := strings.TrimPrefix(key, "coverage.")
	if taskType == "" {
		return "", false, nil
	}

	coverage, err := ReadCoverageConfig(projectRoot)
	if err != nil {
		return "", true, err
	}

	strategy, ok := coverage.ByType[taskType]
	if !ok {
		return "", true, errKeyNotFound
	}

	switch strategy.Type {
	case "maintain":
		return "maintain", true, nil
	case "percentage":
		if strategy.Percentage != nil {
			return strconv.Itoa(*strategy.Percentage), true, nil
		}
		return "", true, errKeyNotFound
	default:
		return "", true, errKeyNotFound
	}
}

// writeConfig writes a Config to .forge/config.yaml.
// Creates the file and directory if they don't exist.
func writeConfig(projectRoot string, cfg *Config) error {
	path := configPath(projectRoot)

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create .forge dir: %w", err)
	}

	out, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, out, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}
