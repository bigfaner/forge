// Package forgeconfig provides types and functions for reading/writing
// the .forge/config.yaml file. This package extracts only the retained
// config types from the legacy profile package: auto and worktree blocks.
//
//nolint:govet // reflect.Ptr inline warnings are toolchain version mismatches, not code issues
package forgeconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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

// EvalConfig controls which eval skills auto-run after document generation.
// Each field is a bool: true means auto-eval is enabled, false means disabled.
//
//nolint:revive // UiDesign matches YAML key convention (camelCase)
type EvalConfig struct {
	Proposal   bool `yaml:"proposal"`
	Prd        bool `yaml:"prd"`
	UiDesign   bool `yaml:"uiDesign"`
	TechDesign bool `yaml:"techDesign"`
}

// UnmarshalYAML implements custom YAML unmarshaling for EvalConfig.
// Supports backward compatibility: if a field receives a map value (old ModeToggle format),
// it extracts the "full" sub-key as the bool value.
func (e *EvalConfig) UnmarshalYAML(value *yaml.Node) error {
	// Use a temporary map to capture raw YAML
	var raw map[string]interface{}
	if err := value.Decode(&raw); err != nil {
		return err
	}

	//nolint:revive // UiDesign matches YAML key convention (camelCase)
	type aliases struct {
		Proposal   interface{} `yaml:"proposal"`
		Prd        interface{} `yaml:"prd"`
		UiDesign   interface{} `yaml:"uiDesign"`
		TechDesign interface{} `yaml:"techDesign"`
	}
	var a aliases
	if err := value.Decode(&a); err != nil {
		return err
	}

	e.Proposal = toBoolWithCompat(a.Proposal)
	e.Prd = toBoolWithCompat(a.Prd)
	e.UiDesign = toBoolWithCompat(a.UiDesign)
	e.TechDesign = toBoolWithCompat(a.TechDesign)
	return nil
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
		},
	}
}

// SurfacesMap is a map[string]string that supports dual-form YAML serialization.
// Scalar form: `surfaces: api` -> map[string]string{".": "api"}
// Map form: `surfaces: {frontend: web}` -> used as-is
// Empty map serializes as `surfaces: {}` (no omitempty).
type SurfacesMap map[string]string

// UnmarshalYAML implements custom YAML unmarshaling for SurfacesMap.
// Handles three cases:
//   - scalar string "api" -> map[string]string{".": "api"}
//   - map form {frontend: web} -> used as-is
//   - nil -> nil map (distinguished from empty map)
func (s *SurfacesMap) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		// Scalar form: "api" -> {".": "api"}
		if value.Value == "" || value.Value == "null" || value.Value == "~" {
			*s = nil
			return nil
		}
		*s = SurfacesMap{".": strings.ToLower(value.Value)}
		return nil

	case yaml.MappingNode:
		// Map form: {frontend: web, backend: api}
		// Normalize keys (spaces/special chars -> hyphens, uppercase -> lowercase)
		// The "." key is the scalar-form marker and must be preserved as-is.
		result := make(SurfacesMap, len(value.Content)/2)
		for i := 0; i < len(value.Content); i += 2 {
			raw := value.Content[i].Value
			key := raw
			if raw != "." {
				key = normalizeSurfaceKeyValue(raw)
			}
			val := strings.ToLower(value.Content[i+1].Value)
			result[key] = val
		}
		*s = result
		return nil

	case yaml.AliasNode:
		return s.UnmarshalYAML(value.Alias)

	default:
		*s = nil
		return nil
	}
}

// MarshalYAML implements custom YAML serialization for SurfacesMap.
// Single entry with key "." -> scalar form.
// Multiple entries or non-"." key -> map form.
// Nil map -> empty map `surfaces: {}`.
func (s SurfacesMap) MarshalYAML() (interface{}, error) {
	if s == nil {
		// Empty surfaces must serialize as `surfaces: {}`, not omitted.
		return map[string]string{}, nil
	}
	if len(s) == 1 {
		if v, ok := s["."]; ok {
			// Single scalar form: return the value directly as a scalar
			return v, nil
		}
	}
	// Map form: return as-is (map[string]string)
	return map[string]string(s), nil
}

// Config represents the .forge/config.yaml structure.
type Config struct {
	Version        string          `yaml:"version,omitempty"`
	ProjectType    string          `yaml:"project-type,omitempty"`
	Auto           *AutoConfig     `yaml:"auto"`
	Worktree       *WorktreeConfig `yaml:"worktree,omitempty"`
	Logs           *LogsConfig     `yaml:"logs,omitempty"`
	Coverage       *CoverageConfig `yaml:"coverage,omitempty"`
	Eval           *EvalSettings   `yaml:"eval,omitempty"`
	TestFramework  string          `yaml:"test-framework,omitempty"`
	Languages      []string        `yaml:"languages,omitempty"`
	Surfaces       SurfacesMap     `yaml:"surfaces"`
	ExecutionOrder []string        `yaml:"execution-order,omitempty"`
}

// Valid project type values.
const (
	ProjectTypeFullstack = "fullstack"
	ProjectTypeMobile    = "mobile"
	ProjectTypeLibrary   = "library"
	ProjectTypeMixed     = "mixed"
)

// ValidProjectTypes lists all valid project type values.
var ValidProjectTypes = []string{ProjectTypeFullstack, ProjectTypeMobile, ProjectTypeLibrary, ProjectTypeMixed}

// ValidProjectType returns true if the given project type is one of the valid values.
func ValidProjectType(pt string) bool {
	for _, v := range ValidProjectTypes {
		if pt == v {
			return true
		}
	}
	return false
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

	// Default version to "1" for backward compatibility with configs
	// created before the version field was introduced.
	if cfg.Version == "" {
		cfg.Version = "1"
	}

	// Parse auto block with explicit-set tracking for default filling
	if cfg.Auto != nil {
		rawAuto, err := parseAutoRaw(data)
		if err == nil {
			cfg.Auto.raw = rawAuto
		}
		cfg.Auto.applyDefaults()
	}

	// v3.0.x migration: detect old key "e2eTest" and map its value to Test field.
	// To be removed in v3.1.0 — v3.1.0 will error on old key instead.
	migrateOldE2eTestKey(data, &cfg)

	// Validate surface keys and execution-order at config load time (fail fast)
	if err := validateSurfacesConfig(&cfg); err != nil {
		return nil, err
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

// GetConfigValue returns the value for a given key from .forge/config.yaml.
// For scalar values, returns the raw string; for arrays, joins with newline.
// Supports arbitrary-depth dot-notation for nested keys via reflection.
// Returns empty string and errKeyNotFound if the key doesn't exist or has zero value.
func GetConfigValue(projectRoot, key string) (string, error) {
	cfg, err := ReadConfig(projectRoot)
	if err != nil {
		return "", err
	}
	if cfg == nil {
		cfg = &Config{}
	}

	// Ensure Auto block has defaults applied for auto.* key lookups
	if strings.HasPrefix(key, "auto") && cfg.Auto == nil {
		cfg.Auto = &AutoConfig{}
		cfg.Auto.applyDefaults()
	}

	// Try reflect-based routing first
	val, err := getByPath(reflect.ValueOf(cfg).Elem(), strings.Split(key, "."))
	if err == nil {
		if val == "" {
			return "", errKeyNotFound
		}
		return val, nil
	}
	if err != errUnsupportedType && err != errKeyNotFound {
		return "", err
	}

	// Fallback: coverage keys (inline map + custom SurfacesMap type)
	if strings.HasPrefix(key, "coverage.") {
		if val, ok, err := getCoverageKeyValue(projectRoot, key); ok || err != nil {
			if err != nil {
				return "", err
			}
			return val, nil
		}
	}

	return "", errKeyNotFound
}

// SetConfigValue sets a config value for a given dot-notation key in .forge/config.yaml.
// Supports arbitrary-depth keys via reflection.
// Returns an error for unknown keys, invalid values, non-leaf sets, or ModeToggle direct sets.
func SetConfigValue(projectRoot, key, value string) error {
	cfg, err := readOrCreateConfig(projectRoot)
	if err != nil {
		return err
	}

	segments := strings.Split(key, ".")
	err = setByPath(reflect.ValueOf(cfg).Elem(), segments, value, key)
	if err == errUnsupportedType || err == errKeyNotFound {
		// Fallback to coverage set for inline map types with dot-containing keys
		if strings.HasPrefix(key, "coverage.") {
			return setCoverageConfigValue(projectRoot, key, value)
		}
		if err == errUnsupportedType {
			return fmt.Errorf("unknown config key: %s", key)
		}
		return fmt.Errorf("config key %q not found", key)
	}
	if err != nil {
		return err
	}
	return writeConfig(projectRoot, cfg)
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

// validateSurfacesConfig runs all surface-related validations at config load time.
// This ensures fail-fast behavior: errors are caught immediately rather than at build time.
// Validations:
//   - Surface-key format: keys must match [a-z][a-z0-9-]* after normalization
//   - Execution-order references: each key must exist in surfaces map
//   - Same-type conflict: multiple surfaces with the same type require explicit execution-order
func validateSurfacesConfig(cfg *Config) error {
	if len(cfg.Surfaces) == 0 {
		return nil
	}

	// Validate surface-key format (after normalization in UnmarshalYAML)
	if err := ValidateSurfaceKeys(cfg.Surfaces); err != nil {
		return fmt.Errorf("config validation: %w", err)
	}

	// Validate execution-order references and same-type conflicts
	if err := ValidateExecutionOrder(cfg.Surfaces, cfg.ExecutionOrder); err != nil {
		return fmt.Errorf("config validation: %w", err)
	}

	return nil
}
