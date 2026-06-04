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

// toBoolWithCompat converts a YAML value to bool, handling backward compatibility
// with the old ModeToggle map format. If the value is a map, it extracts the "full"
// sub-key. Otherwise, it interprets the value as a bool.
func toBoolWithCompat(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case map[string]interface{}:
		// Old ModeToggle format: extract "full" sub-key
		if full, ok := val["full"]; ok {
			if b, ok := full.(bool); ok {
				return b
			}
		}
		return false
	default:
		return false
	}
}

// AutoConfig controls which auto-generated tasks are produced by `forge task index`.
// When the `auto` block is missing from config, all fields use defaults that match
// pre-auto-behavior.
type AutoConfig struct {
	Test             ModeToggle `yaml:"test"`
	ConsolidateSpecs ModeToggle `yaml:"consolidateSpecs"`
	CleanCode        ModeToggle `yaml:"cleanCode"`
	Validation       ModeToggle `yaml:"validation"`
	RunTasks         ModeToggle `yaml:"runTasks"`
	GitPush          bool       `yaml:"gitPush"`
	KnowledgeSave    ModeToggle `yaml:"knowledgeSave"`
	Eval             EvalConfig `yaml:"eval"`
	// raw tracks which sub-fields were explicitly present in the YAML.
	// Used by applyDefaults to distinguish "false" from "missing".
	raw map[string]map[string]bool
}

// AutoConfigDefaults returns an AutoConfig with backward-compatible defaults:
// test: quick=false, full=true; consolidateSpecs: quick=true, full=true;
// cleanCode=false, validation=false, gitPush=false.
// Eval: proposal=true, prd=false, uiDesign=true, techDesign=false.
func AutoConfigDefaults() AutoConfig {
	return AutoConfig{
		Test:             ModeToggle{Quick: false, Full: true},
		ConsolidateSpecs: ModeToggle{Quick: true, Full: true},
		CleanCode:        ModeToggle{Quick: false, Full: false},
		Validation:       ModeToggle{Quick: false, Full: false},
		RunTasks:         ModeToggle{Quick: true, Full: false},
		GitPush:          false,
		KnowledgeSave:    ModeToggle{Quick: true, Full: false},
		Eval: EvalConfig{
			Proposal:   true,
			Prd:        false,
			UiDesign:   true,
			TechDesign: false,
		},
	}
}

// IsZero returns true if the AutoConfig has all zero-value fields.
func (a AutoConfig) IsZero() bool {
	return a.Test == ModeToggle{} &&
		a.ConsolidateSpecs == ModeToggle{} &&
		a.CleanCode == ModeToggle{} &&
		a.Validation == ModeToggle{} &&
		a.RunTasks == ModeToggle{} &&
		!a.GitPush &&
		a.KnowledgeSave == ModeToggle{} &&
		a.Eval == EvalConfig{}
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

// SurfacesMap is a map[string]string that supports dual-form YAML serialization.
// Scalar form: `surfaces: api` → map[string]string{".": "api"}
// Map form: `surfaces: {frontend: web}` → used as-is
// Empty map serializes as `surfaces: {}` (no omitempty).
type SurfacesMap map[string]string

// UnmarshalYAML implements custom YAML unmarshaling for SurfacesMap.
// Handles three cases:
//   - scalar string "api" → map[string]string{".": "api"}
//   - map form {frontend: web} → used as-is
//   - nil → nil map (distinguished from empty map)
func (s *SurfacesMap) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		// Scalar form: "api" → {".": "api"}
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
// Single entry with key "." → scalar form.
// Multiple entries or non-"." key → map form.
// Nil map → empty map `surfaces: {}`.
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

// EvalTypeSettings holds per-eval-type target score and iteration count.
// Pointer fields: nil means not configured (fallback to rubric defaults),
// non-nil overrides. The reflection router's derefPointer returns
// errKeyNotFound for nil pointers, providing correct "not configured" semantics.
type EvalTypeSettings struct {
	Target     *int `yaml:"target,omitempty"`
	Iterations *int `yaml:"iterations,omitempty"`
}

// EvalSettings holds eval configuration for all 7 eval types.
// Each type maps to a rubric file: proposal, prd, design, ui, journey, contract, consistency.
type EvalSettings struct {
	Proposal EvalTypeSettings `yaml:"proposal"`
	Prd      EvalTypeSettings `yaml:"prd"`
	Design   EvalTypeSettings `yaml:"design"`
	//nolint:revive // Ui matches YAML key convention (lowercase)
	Ui          EvalTypeSettings `yaml:"ui"`
	Journey     EvalTypeSettings `yaml:"journey"`
	Contract    EvalTypeSettings `yaml:"contract"`
	Consistency EvalTypeSettings `yaml:"consistency"`
}

// Config represents the .forge/config.yaml structure.
type Config struct {
	Version        string          `yaml:"version,omitempty"`
	ProjectType    string          `yaml:"project-type,omitempty"`
	Auto           *AutoConfig     `yaml:"auto"`
	Worktree       *WorktreeConfig `yaml:"worktree,omitempty"`
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

// migrateOldE2eTestKey detects the old "auto.e2eTest" key in raw YAML,
// outputs a migration hint to stderr, and maps the value to the new Test field.
// This is a v3.0.x transitional helper; remove in v3.1.0.
func migrateOldE2eTestKey(data []byte, cfg *Config) {
	if cfg == nil {
		return
	}
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return
	}
	autoNode := findMappingKey(&root, "auto")
	if autoNode == nil {
		return
	}
	oldNode := findMappingKey(autoNode, "e2eTest")
	if oldNode == nil {
		return
	}

	fmt.Fprintln(os.Stderr, "config key 'auto.e2eTest' is renamed to 'auto.test' in v3.0.0; please update your config.yaml")

	if cfg.Auto == nil {
		cfg.Auto = &AutoConfig{}
	}

	// Parse the old e2eTest block into ModeToggle
	var mt ModeToggle
	if oldNode.Kind == yaml.MappingNode {
		for i := 0; i < len(oldNode.Content); i += 2 {
			key := oldNode.Content[i].Value
			val := oldNode.Content[i+1].Value
			switch key {
			case "quick":
				mt.Quick = val == "true"
			case "full":
				mt.Full = val == "true"
			}
		}
	}
	cfg.Auto.Test = mt

	// Re-apply defaults with the raw tracking
	rawAuto, err := parseAutoRaw(data)
	if err == nil {
		cfg.Auto.raw = rawAuto
	}
	cfg.Auto.applyDefaults()
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
// For the old key name "e2eTest" (renamed to "test" in v3.0.0), it maps the value to
// the new "test" key in the result. Migration hint output is handled by migrateOldE2eTestKey.
// Uses recursive scanning to support arbitrary nesting (e.g. "eval.proposal").
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

	// Detect old key name "e2eTest" and map to "test" (v3.0.x migration, to be removed in v3.1.0)
	if oldNode := findMappingKey(autoNode, "e2eTest"); oldNode != nil {
		result["test"] = make(map[string]bool)
		if oldNode.Kind == yaml.MappingNode {
			for i := 0; i < len(oldNode.Content); i += 2 {
				key := oldNode.Content[i].Value
				if key == "quick" || key == "full" {
					result["test"][key] = true
				}
			}
		}
	}

	// Recursive scan of auto block
	scanMappingNode(autoNode, "", result)

	return result, nil
}

// scanMappingNode recursively scans a YAML mapping node for ModeToggle-like
// sub-fields (quick/full), building flat-path keys.
// For eval sub-fields (under "eval" prefix), it also tracks scalar (bool) values
// to support the new flat bool format alongside the old ModeToggle map format.
func scanMappingNode(node *yaml.Node, prefix string, result map[string]map[string]bool) {
	if node.Kind != yaml.MappingNode {
		return
	}
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i].Value
		// Skip old e2eTest key (already handled above)
		if key == "e2eTest" {
			continue
		}
		valNode := node.Content[i+1]
		flatKey := key
		if prefix != "" {
			flatKey = prefix + "." + key
		}

		if valNode.Kind == yaml.MappingNode {
			// Check if this is a ModeToggle node (has quick/full sub-keys)
			if isModeToggleNode(valNode) {
				if _, exists := result[flatKey]; !exists {
					result[flatKey] = make(map[string]bool)
				}
				for j := 0; j < len(valNode.Content); j += 2 {
					subKey := valNode.Content[j].Value
					if subKey == "quick" || subKey == "full" {
						result[flatKey][subKey] = true
					}
				}
			} else {
				// Recurse into nested struct (e.g. "eval")
				scanMappingNode(valNode, flatKey, result)
			}
		} else if valNode.Kind == yaml.ScalarNode && prefix == "eval" {
			// Eval sub-field in new bool format: track its presence
			if _, exists := result[flatKey]; !exists {
				result[flatKey] = make(map[string]bool)
			}
			result[flatKey]["set"] = true
		}
	}
}

// isModeToggleNode checks if a YAML mapping node looks like a ModeToggle
// (has "quick" and/or "full" keys).
func isModeToggleNode(node *yaml.Node) bool {
	if node.Kind != yaml.MappingNode {
		return false
	}
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i].Value
		if key == "quick" || key == "full" {
			return true
		}
	}
	return false
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
		a.Test = d.Test
		a.ConsolidateSpecs = d.ConsolidateSpecs
		a.CleanCode = d.CleanCode
		a.Validation = d.Validation
		a.RunTasks = d.RunTasks
		a.KnowledgeSave = d.KnowledgeSave
		a.Eval = d.Eval
		return
	}

	applyModeDefault(&a.Test, a.raw, "test", d.Test)
	applyModeDefault(&a.ConsolidateSpecs, a.raw, "consolidateSpecs", d.ConsolidateSpecs)
	applyModeDefault(&a.CleanCode, a.raw, "cleanCode", d.CleanCode)
	applyModeDefault(&a.Validation, a.raw, "validation", d.Validation)
	applyModeDefault(&a.RunTasks, a.raw, "runTasks", d.RunTasks)
	applyModeDefault(&a.KnowledgeSave, a.raw, "knowledgeSave", d.KnowledgeSave)

	// Eval sub-fields: simple bool defaults (no ModeToggle)
	applyBoolDefault(&a.Eval.Proposal, a.raw, "eval.proposal", d.Eval.Proposal)
	applyBoolDefault(&a.Eval.Prd, a.raw, "eval.prd", d.Eval.Prd)
	applyBoolDefault(&a.Eval.UiDesign, a.raw, "eval.uiDesign", d.Eval.UiDesign)
	applyBoolDefault(&a.Eval.TechDesign, a.raw, "eval.techDesign", d.Eval.TechDesign)
}

// applyBoolDefault sets a default bool value for a field that was not explicitly set in YAML.
// The raw map tracks whether the field was present in the YAML; if absent, the default is applied.
func applyBoolDefault(field *bool, raw map[string]map[string]bool, key string, defaults bool) {
	_, exists := raw[key]
	if !exists {
		*field = defaults
		return
	}
	// Field was explicitly set in YAML — keep the value set by UnmarshalYAML
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

// errUnsupportedType is returned when a config field implements yaml.Unmarshaler
// and the generic reflect router cannot handle it (e.g. SurfacesMap).
var errUnsupportedType = fmt.Errorf("unsupported type for reflect routing")

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

// getByPath traverses a reflect.Value by path segments, returning the formatted value.
func getByPath(v reflect.Value, segments []string) (string, error) {
	var err error
	for i, seg := range segments {
		v, err = navigateToSegment(v, seg)
		if err == errKeyNotFound {
			// Try inline map with dot-joined remaining segments
			v2 := derefPointer(v)
			if v2.IsValid() && v2.Kind() == reflect.Struct {
				if inlineField, ok := findInlineMapField(v2); ok {
					mapVal := derefPointer(inlineField)
					if mapVal.IsValid() && mapVal.Kind() == reflect.Map {
						dotKey := strings.Join(segments[i:], ".")
						entry := mapVal.MapIndex(reflect.ValueOf(dotKey))
						if entry.IsValid() {
							return formatValue(derefPointer(entry))
						}
					}
				}
			}
			return "", err
		}
		if err != nil {
			return "", err
		}

		// Reached the target segment
		if i == len(segments)-1 {
			return formatValue(v)
		}

		// More segments to go — check if current value is navigable
		if isLeafType(v) {
			return "", errKeyNotFound
		}
		// Continue descending
	}
	return "", errKeyNotFound
}

// navigateToSegment resolves one path segment within the given reflect.Value.
func navigateToSegment(v reflect.Value, seg string) (reflect.Value, error) {
	// Dereference pointers
	v = derefPointer(v)
	if !v.IsValid() {
		return reflect.Value{}, errKeyNotFound
	}

	kind := v.Kind()

	switch kind {
	case reflect.Struct:
		field, found := findFieldByYAMLTag(v, seg)
		if !found {
			// Check for yaml:",inline" map fields
			if inlineField, ok := findInlineMapField(v); ok {
				mapVal := derefPointer(inlineField)
				if mapVal.IsValid() && mapVal.Kind() == reflect.Map {
					entry := mapVal.MapIndex(reflect.ValueOf(seg))
					if entry.IsValid() {
						return derefPointer(entry), nil
					}
				}
			}
			return reflect.Value{}, errKeyNotFound
		}
		return derefPointer(field), nil

	case reflect.Map:
		keyVal := reflect.ValueOf(seg)
		entry := v.MapIndex(keyVal)
		if !entry.IsValid() {
			return reflect.Value{}, errKeyNotFound
		}
		return derefPointer(entry), nil

	default:
		return reflect.Value{}, errKeyNotFound
	}
}

// findFieldByYAMLTag finds a struct field matching the segment by YAML tag (priority)
// or Go field name.
func findFieldByYAMLTag(v reflect.Value, seg string) (reflect.Value, bool) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		tag := field.Tag.Get("yaml")
		tagName := parseYAMLTagName(tag, field.Name)
		if tagName == seg {
			return v.Field(i), true
		}
	}
	return reflect.Value{}, false
}

// parseYAMLTagName extracts the YAML key name from a yaml tag.
// Priority: yaml:"name" → name; yaml:",inline" → "" (skip); no tag → GoFieldName.
// Returns empty string for ",inline" and "omitempty" only tags.
func parseYAMLTagName(tag, goName string) string {
	if tag == "" {
		return goName
	}
	// Split by comma: first part is name, rest are options
	parts := strings.Split(tag, ",")
	name := parts[0]
	if name == "" {
		// ",inline" or ",omitempty" — not a key match target
		return ""
	}
	return name
}

// findInlineMapField finds a struct field tagged with yaml:",inline" that is a map type.
func findInlineMapField(v reflect.Value) (reflect.Value, bool) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		tag := field.Tag.Get("yaml")
		if strings.Contains(tag, ",inline") {
			return v.Field(i), true
		}
	}
	return reflect.Value{}, false
}

// derefPointer dereferences a pointer, returning the zero Value if nil.
func derefPointer(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return reflect.Value{}
		}
		v = v.Elem()
	}
	return v
}

// isLeafType returns true if the value cannot be further navigated.
func isLeafType(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	kind := v.Kind()
	switch kind {
	case reflect.Struct, reflect.Map:
		return false
	default:
		return true
	}
}

// formatValue formats a reflect.Value for CLI output.
// For leaf types, returns the scalar value.
// For non-leaf types (struct, map), returns a multi-line summary.
func formatValue(v reflect.Value) (string, error) {
	if !v.IsValid() {
		return "", errKeyNotFound
	}

	kind := v.Kind()

	// Check for custom YAML types that reflect routing cannot handle.
	// Only applies to non-struct types (e.g. SurfacesMap as map type).
	// Structs that implement yaml.Unmarshaler (like EvalConfig for compat)
	// are handled by the struct formatting path below.
	if kind != reflect.Struct && implementsYAMLUnmarshaler(v) {
		return "", errUnsupportedType
	}

	switch kind {
	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), nil
	case reflect.String:
		return v.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), nil
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.String {
			slice := make([]string, v.Len())
			for i := 0; i < v.Len(); i++ {
				slice[i] = v.Index(i).String()
			}
			return joinSlice(slice), nil
		}
		return "", errUnsupportedType
	case reflect.Struct:
		if isModeToggle(v.Type()) {
			q := v.FieldByName("Quick").Bool()
			f := v.FieldByName("Full").Bool()
			return fmt.Sprintf("quick:%v full:%v", q, f), nil
		}
		return formatStructSummary(v, "")
	case reflect.Map:
		return formatMapSummary(v, "")
	default:
		return "", errUnsupportedType
	}
}

// implementsYAMLUnmarshaler checks if the value's type implements yaml.Unmarshaler.
//
//nolint:govet // reflect.PtrTo inline warning is a toolchain version mismatch, not a code issue
func implementsYAMLUnmarshaler(v reflect.Value) bool {
	t := v.Type()
	ptr := reflect.PointerTo(t)
	return ptr.Implements(reflect.TypeOf((*yaml.Unmarshaler)(nil)).Elem())
}

// formatStructSummary formats a struct's exported fields as a multi-line summary.
func formatStructSummary(v reflect.Value, indent string) (string, error) {
	var lines []string
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		// Skip unexported internal fields (like 'raw')
		tag := field.Tag.Get("yaml")
		if tag == "" && field.Name == "raw" {
			continue
		}
		if strings.Contains(tag, ",inline") {
			continue
		}

		fieldName := parseYAMLTagName(tag, field.Name)
		if fieldName == "" {
			continue
		}

		fv := derefPointer(v.Field(i))
		if !fv.IsValid() {
			continue
		}

		line, err := formatFieldLine(fieldName, fv, indent)
		if err != nil {
			continue
		}
		lines = append(lines, line)
	}

	if len(lines) == 0 {
		return "", errKeyNotFound
	}
	return strings.Join(lines, "\n"), nil
}

// formatMapSummary formats a map's entries as a multi-line summary.
func formatMapSummary(v reflect.Value, indent string) (string, error) {
	var lines []string
	iter := v.MapRange()
	for iter.Next() {
		key := iter.Key().String()
		entry := derefPointer(iter.Value())
		if !entry.IsValid() {
			continue
		}
		line, err := formatFieldLine(key, entry, indent)
		if err != nil {
			continue
		}
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		return "", errKeyNotFound
	}
	return strings.Join(lines, "\n"), nil
}

// formatFieldLine formats a single field for summary output.
func formatFieldLine(name string, v reflect.Value, indent string) (string, error) {
	kind := v.Kind()
	switch kind {
	case reflect.Struct:
		if isModeToggle(v.Type()) {
			// ModeToggle → "name: quick:X full:Y"
			q := v.FieldByName("Quick").Bool()
			f := v.FieldByName("Full").Bool()
			return fmt.Sprintf("%s%s: quick:%v full:%v", indent, name, q, f), nil
		}
		// Nested struct → "name:\n" + recursive lines with +2 indent
		sub, err := formatStructSummary(v, indent+"  ")
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s%s:\n%s", indent, name, sub), nil
	case reflect.Bool:
		return fmt.Sprintf("%s%s: %v", indent, name, v.Bool()), nil
	case reflect.String:
		return fmt.Sprintf("%s%s: %s", indent, name, v.String()), nil
	case reflect.Map:
		sub, err := formatMapSummary(v, indent+"  ")
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s%s:\n%s", indent, name, sub), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%s%s: %d", indent, name, v.Int()), nil
	default:
		return "", errUnsupportedType
	}
}

// isModeToggle checks if a type is ModeToggle.
func isModeToggle(t reflect.Type) bool {
	return t.Name() == "ModeToggle" && t.Kind() == reflect.Struct &&
		t.NumField() == 2
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

// setByPath traverses a reflect.Value by segments and sets the leaf value.
func setByPath(v reflect.Value, segments []string, value string, fullKey string) error {
	for i, seg := range segments {
		v = ensureAddressable(v)

		// Dereference pointers, initializing nil pointers as needed
		for v.Kind() == reflect.Ptr {
			if v.IsNil() {
				newVal := reflect.New(v.Type().Elem())
				v.Set(newVal)
			}
			v = v.Elem()
		}

		if v.Kind() == reflect.Struct { //nolint:gocritic // ifElseChain
			field, found := findSettableField(v, seg)
			if !found {
				// Check for inline map - try joining remaining segments as dot-separated key
				if inlineField, ok := findInlineMapField(v); ok {
					mapVal := ensureAddressable(inlineField)
					for mapVal.Kind() == reflect.Ptr {
						if mapVal.IsNil() {
							newVal := reflect.New(mapVal.Type().Elem())
							mapVal.Set(newVal)
						}
						mapVal = mapVal.Elem()
					}
					if mapVal.Kind() == reflect.Map {
						if mapVal.IsNil() {
							mapVal.Set(reflect.MakeMap(mapVal.Type()))
						}
						if i == len(segments)-1 {
							return fmt.Errorf("cannot set non-leaf key, use %s.<field>", fullKey)
						}
						// Join remaining segments as a single dot-separated key for inline maps
						dotKey := strings.Join(segments[i:], ".")
						return setMapEntry(mapVal, []string{dotKey}, value, fullKey)
					}
				}
				return fmt.Errorf("config key %q not found", fullKey)
			}

			// Last segment - set the value
			if i == len(segments)-1 {
				return setFieldValue(field, value, fullKey)
			}

			// Intermediate segment - allow descending into ModeToggle
			if isLeafType(field) && !isModeToggle(field.Type()) {
				return fmt.Errorf("cannot set non-leaf key, use %s.<field>", fullKey)
			}
			v = field
		} else if v.Kind() == reflect.Map {
			if i == len(segments)-1 {
				return fmt.Errorf("cannot set non-leaf key, use %s.<field>", fullKey)
			}
			return setMapEntry(v, segments[i+1:], value, fullKey)
		} else {
			return errKeyNotFound
		}
	}
	return fmt.Errorf("cannot set non-leaf key, use %s.<field>", fullKey)
}

// ensureAddressable returns an addressable reflect.Value.
func ensureAddressable(v reflect.Value) reflect.Value {
	if v.CanAddr() {
		return v
	}
	// For non-addressable values (from reflect.ValueOf), try to get a pointer
	return v
}

// findSettableField finds a struct field matching the segment and returns a settable Value.
func findSettableField(v reflect.Value, seg string) (reflect.Value, bool) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		tag := field.Tag.Get("yaml")
		tagName := parseYAMLTagName(tag, field.Name)
		if tagName == seg {
			fv := v.Field(i)
			// For pointer fields, initialize nil and dereference
			if fv.Kind() == reflect.Ptr && fv.IsNil() {
				newVal := reflect.New(fv.Type().Elem())
				fv.Set(newVal)
			}
			if fv.Kind() == reflect.Ptr {
				fv = fv.Elem()
			}
			return fv, true
		}
	}
	return reflect.Value{}, false
}

// setFieldValue sets a leaf field's value from a string.
func setFieldValue(field reflect.Value, value string, fullKey string) error {
	if isModeToggle(field.Type()) {
		return fmt.Errorf("cannot set ModeToggle directly, use %s.quick or %s.full", fullKey, fullKey)
	}

	switch field.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid value %q for bool field %s: expected true or false", value, fullKey)
		}
		field.SetBool(b)
		return nil
	case reflect.String:
		field.SetString(value)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid value %q for int field %s: expected integer", value, fullKey)
		}
		field.SetInt(int64(n))
		return nil
	default:
		return fmt.Errorf("cannot set non-leaf key, use %s.<field>", fullKey)
	}
}

// setMapEntry sets a value in a map for the given remaining segments.
func setMapEntry(mapVal reflect.Value, segments []string, value string, fullKey string) error {
	if len(segments) != 1 {
		return errUnsupportedType
	}
	key := segments[0]

	// For CoverageConfig.ByType: value is a percentage number
	pct, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid coverage value for %s: %s (expected percentage number)", fullKey, value)
	}

	strategyType := reflect.TypeOf(CoverageStrategy{})
	strategyVal := reflect.New(strategyType).Elem()
	strategyVal.FieldByName("Type").SetString("percentage")
	pctField := strategyVal.FieldByName("Percentage")
	pctVal := reflect.New(pctField.Type().Elem())
	pctVal.Elem().SetInt(int64(pct))
	strategyVal.FieldByName("Percentage").Set(pctVal)

	mapVal.SetMapIndex(reflect.ValueOf(key), strategyVal)
	return nil
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

// joinSlice joins slice values with newline for plain-text output.
func joinSlice(vals []string) string {
	return strings.Join(vals, "\n")
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
