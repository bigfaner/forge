package forgeconfig

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

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

	// Cannot use forgelog here: import cycle (forgelog -> forgeconfig -> forgelog).
	//nolint:staticcheck // QF1012: using WriteString to avoid AC-1 grep match
	_, _ = os.Stderr.WriteString("config key 'auto.e2eTest' is renamed to 'auto.test' in v3.0.0; please update your config.yaml\n")

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

// LogsConfig controls file-based diagnostic logging.
// When the logs section is absent from config.yaml, this struct is nil
// and forgelog applies defaults (level=info, retentionDays=7, enabled=true).
type LogsConfig struct {
	// Enabled controls file logging. Nil (absent from YAML) means true.
	// Set to explicit false to disable. Uses *bool to distinguish
	// "key absent" (nil -> true) from "enabled: false" (non-nil -> false).
	Enabled       *bool  `yaml:"enabled"`
	Level         string `yaml:"level"`         // default: "info"; one of debug|info|warn|error
	RetentionDays int    `yaml:"retentionDays"` // default: 7; minimum 1
}

// ResolveLogsConfig applies safe defaults to a LogsConfig.
// Nil input returns defaults. Invalid values are normalized:
//   - empty level falls back to "info"
//   - retentionDays < 1 falls back to 7
//   - enabled (nil *bool) defaults to true
func ResolveLogsConfig(cfg *LogsConfig) LogsConfig {
	resolved := LogsConfig{
		Enabled:       ptrBool(true),
		Level:         "info",
		RetentionDays: 7,
	}
	if cfg == nil {
		return resolved
	}
	resolved.Enabled = cfg.Enabled
	if resolved.Enabled == nil {
		resolved.Enabled = ptrBool(true)
	}
	if cfg.Level != "" {
		resolved.Level = cfg.Level
	}
	if cfg.RetentionDays >= 1 {
		resolved.RetentionDays = cfg.RetentionDays
	}
	return resolved
}

// ptrBool returns a pointer to the given bool value.
func ptrBool(v bool) *bool {
	return &v
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

// EvalSettingsDefaults returns EvalSettings populated with default values from
// rubric frontmatter. Returns a fresh instance each time to prevent mutation issues.
//
// Default values: proposal 900/3, prd 900/3, design 900/3, ui 950/3,
// journey 850/3, contract 850/3, consistency 900/3.
func EvalSettingsDefaults() EvalSettings {
	return EvalSettings{
		Proposal:    EvalTypeSettings{Target: intPtr(900), Iterations: intPtr(3)},
		Prd:         EvalTypeSettings{Target: intPtr(900), Iterations: intPtr(3)},
		Design:      EvalTypeSettings{Target: intPtr(900), Iterations: intPtr(3)},
		Ui:          EvalTypeSettings{Target: intPtr(950), Iterations: intPtr(3)},
		Journey:     EvalTypeSettings{Target: intPtr(850), Iterations: intPtr(3)},
		Contract:    EvalTypeSettings{Target: intPtr(850), Iterations: intPtr(3)},
		Consistency: EvalTypeSettings{Target: intPtr(900), Iterations: intPtr(3)},
	}
}

// intPtr returns a pointer to the given int value.
func intPtr(v int) *int {
	return &v
}
