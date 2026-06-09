package task

import (
	"sort"
	"strings"

	"forge-cli/pkg/forgeconfig"
	"forge-cli/pkg/types"
)

const (
	estTimeQuickTask = "15min"
	estTimeMedium    = "1-2h"
)

// ConfigGateFunc returns true when the auto config enables this node for the given mode.
type ConfigGateFunc func(mode string, auto forgeconfig.AutoConfig) bool

// IntentGateFunc returns true when the intent permits this node to generate.
type IntentGateFunc func(intent string) bool

// GenerateCondFunc returns true when the business task composition permits this node.
type GenerateCondFunc func(tasks []Task) bool

// DepResolveFunc dynamically resolves dependency IDs at generation time.
// Returns nil when the reference cannot be resolved. If ALL dependencies of a node
// resolve to nil, the node generates with empty DependsOn (pipeline root).
type DepResolveFunc func(ctx *GenContext) []string

// GenContext carries state accumulated during pipeline generation.
type GenContext struct {
	Mode           string
	Intent         string
	Surfaces       map[string]string
	ExecutionOrder []string
	Auto           forgeconfig.AutoConfig
	BusinessTasks  []Task
	ExistingTasks  map[string]Task
	UpstreamIDs    []string
	RunTestChain   []string
	AllGenerated   []string
	GenScriptsMap  map[string]string // surface-key -> gen-scripts-task-ID
	RunTestMap     map[string]string // surface-key -> test-run-task-ID
}

// PipelineNode defines a single node in the auto-generated task pipeline.
type PipelineNode struct {
	Type              string
	Key               string
	ID                string
	Title             string
	Priority          string
	EstimatedTime     string
	ConfigGate        ConfigGateFunc
	IntentGate        IntentGateFunc
	Mode              string
	GenerateCondition GenerateCondFunc
	DependsOn         []DepRef
	Expansion         string // "", "per-surface-key", "per-surface-type"
	MainSession       bool
	Breaking          bool // when true, submit quality gate includes unit-test; default false for verification/validation tasks
	StrategyKind      string
	UISurfaceOnly     bool
}

// DepRef represents a dependency reference. If Resolve is non-nil, Ref is ignored.
type DepRef struct {
	Ref     string
	Resolve DepResolveFunc
}

// GateTest returns true when the Test category is enabled for the given mode.
func GateTest(mode string, auto forgeconfig.AutoConfig) bool {
	if mode == "quick" {
		return auto.Test.Quick
	}
	return auto.Test.Full
}

// GateValidation returns true when Validation is enabled for the given mode.
func GateValidation(mode string, auto forgeconfig.AutoConfig) bool {
	if mode == "quick" {
		return auto.Validation.Quick
	}
	return auto.Validation.Full
}

// GateConsolidateSpecs returns true when ConsolidateSpecs is enabled for the given mode.
func GateConsolidateSpecs(mode string, auto forgeconfig.AutoConfig) bool {
	if mode == "quick" {
		return auto.ConsolidateSpecs.Quick
	}
	return auto.ConsolidateSpecs.Full
}

// GateCleanCode returns true when CleanCode is enabled for the given mode.
func GateCleanCode(mode string, auto forgeconfig.AutoConfig) bool {
	if mode == "quick" {
		return auto.CleanCode.Quick
	}
	return auto.CleanCode.Full
}

// GateAllowAll permits all intents. Used by T-review-doc.
func GateAllowAll(_ string) bool { return true }

// GateBlockSkipTest blocks refactor/cleanup intents. Used by all config-gated nodes.
func GateBlockSkipTest(intent string) bool {
	return !isSkipTestIntent(intent)
}

// CondHasTestableTasks returns true when any business task has a testable type.
// When tasks is nil, returns true (legacy compat).
func CondHasTestableTasks(tasks []Task) bool {
	if tasks == nil {
		return true
	}
	for _, t := range tasks {
		if IsTestableType(t.Type) {
			return true
		}
	}
	return false
}

// CondHasDocTasks returns true when any business task has a doc-category type.
func CondHasDocTasks(tasks []Task) bool {
	for _, t := range tasks {
		if CategoryForType(t.Type) == CategoryDoc {
			return true
		}
	}
	return false
}

// CondAlways returns true unconditionally.
func CondAlways(_ []Task) bool { return true }

// protocolSurfaceTypes is the set of surface types that use protocol-level
// contracts (API, CLI, TUI). Web and Mobile are interaction-level surfaces
// that do not produce meaningful contract files.
var protocolSurfaceTypes = map[types.SurfaceType]bool{
	types.SurfaceAPI: true,
	types.SurfaceCLI: true,
	types.SurfaceTUI: true,
}

// CondHasProtocolSurfaceTask returns true when any business task has a
// protocol-level surface type (api/cli/tui). Returns true conservatively
// when tasks is nil/empty or when any task has an unknown/empty surface type.
// Returns false only when ALL tasks have explicitly interaction-level surfaces
// (web/mobile).
func CondHasProtocolSurfaceTask(tasks []Task) bool {
	if len(tasks) == 0 {
		return true // conservative: no tasks = don't skip
	}
	for _, t := range tasks {
		st := types.SurfaceType(t.SurfaceType)
		if st == "" {
			return true // conservative: unknown = don't skip
		}
		if _, ok := protocolSurfaceTypes[st]; ok {
			return true // has protocol surface
		}
		// If it's not a protocol type and not a known interaction type, conservative
		if st != types.SurfaceWeb && st != types.SurfaceMobile {
			return true // unknown surface type = don't skip
		}
	}
	return false // all tasks are interaction-level (web/mobile only)
}

// ---------------------------------------------------------------------------
// GenerateTestTasks — registry-driven task generation
// ---------------------------------------------------------------------------

// GenerateTestTasks filters PipelineRegistry by mode/config/intent/condition/ui constraints,
// expands per-surface nodes, resolves dependencies via GenContext progressive population,
// and returns the generated AutoGenTaskDef list.
func GenerateTestTasks(mode string, surfaces map[string]string, executionOrder []string, auto forgeconfig.AutoConfig, intent string, businessTasks []Task, existingTasks map[string]Task) []AutoGenTaskDef {
	ctx := &GenContext{
		Mode: mode, Intent: intent, Surfaces: surfaces, ExecutionOrder: executionOrder,
		Auto: auto, BusinessTasks: businessTasks, ExistingTasks: existingTasks,
	}

	var generated []AutoGenTaskDef

	// Phase 1: Expand all nodes and apply default dependency wiring.
	// Track gen-scripts and test-run expansions for interleaved rewiring.
	type nodeExpansion struct {
		nodeType string
		indices  []int // indices into generated slice
	}
	var expansions []nodeExpansion

	for _, node := range PipelineRegistry {
		if node.Mode != "" && node.Mode != mode {
			continue
		}
		if node.ConfigGate != nil && !node.ConfigGate(mode, auto) {
			continue
		}
		intentGate := node.IntentGate
		if intentGate == nil {
			intentGate = GateBlockSkipTest
		}
		if !intentGate(intent) {
			continue
		}
		if !node.GenerateCondition(businessTasks) {
			continue
		}
		if node.UISurfaceOnly && !hasVisualUI(surfaces) {
			continue
		}

		expanded := expandNode(node, surfaces, executionOrder)

		startIdx := len(generated)
		var indices []int
		for i := range expanded {
			if node.Expansion == "per-surface-key" && i > 0 {
				expanded[i].Dependencies = []string{expanded[i-1].ID}
			} else {
				for _, dep := range node.DependsOn {
					if dep.Resolve != nil {
						ids := dep.Resolve(ctx)
						if ids == nil {
							continue
						}
						expanded[i].Dependencies = append(expanded[i].Dependencies, ids...)
					} else {
						expanded[i].Dependencies = append(expanded[i].Dependencies, dep.Ref)
					}
				}
			}
			indices = append(indices, startIdx+i)
		}

		ids := pipelineTaskIDs(expanded)
		ctx.AllGenerated = append(ctx.AllGenerated, ids...)
		ctx.UpstreamIDs = ids
		if node.Type == TypeTestRun {
			ctx.RunTestChain = append(ctx.RunTestChain, ids...)
		}

		// Track gen-scripts and test-run expansions for interleaved rewiring.
		if (node.Type == TypeTestGenScripts || node.Type == TypeTestRun) && node.Expansion == "per-surface-key" {
			expansions = append(expansions, nodeExpansion{nodeType: node.Type, indices: indices})
		}

		generated = append(generated, expanded...)
	}

	// Phase 2: Interleaved rewiring for multi-surface gen-scripts and test-run.
	// After all nodes are expanded, we have full surface-key -> task-ID mappings.
	if len(expansions) == 2 && !isSingleSurface(surfaces) {
		// Build maps from the expanded tasks.
		genScriptsMap := make(map[string]string)
		runTestMap := make(map[string]string)
		for _, exp := range expansions {
			for _, idx := range exp.indices {
				t := &generated[idx]
				if exp.nodeType == TypeTestGenScripts && t.SurfaceKey != "" {
					genScriptsMap[t.SurfaceKey] = t.ID
				}
				if exp.nodeType == TypeTestRun && t.SurfaceKey != "" {
					runTestMap[t.SurfaceKey] = t.ID
				}
			}
		}

		keys := executionOrderKeys(executionOrder, surfaces)
		for _, exp := range expansions {
			for i, idx := range exp.indices {
				t := &generated[idx]
				switch {
				// gen-scripts i>0: depend on previous surface's test-run
				case exp.nodeType == TypeTestGenScripts && i > 0 && i-1 < len(keys):
					prevSurfaceKey := keys[i-1]
					if runID, ok := runTestMap[prevSurfaceKey]; ok {
						t.Dependencies = []string{runID}
					}
				// test-run: depend on corresponding surface's gen-scripts
				case exp.nodeType == TypeTestRun && t.SurfaceKey != "":
					if genID, ok := genScriptsMap[t.SurfaceKey]; ok {
						t.Dependencies = []string{genID}
					}
				}
			}
		}
	}

	// Phase 3: Dynamic validation (errors indicate programming bugs).
	_ = validateGeneratedTasks(generated)
	return generated
}

// hasVisualUI returns true when at least one surface has a visual UI type (TUI, Web, Mobile).
func hasVisualUI(surfaces map[string]string) bool {
	for _, typ := range surfaces {
		if uiSurfaceTypes[types.SurfaceType(typ)] {
			return true
		}
	}
	return false
}

// expandNode produces concrete AutoGenTaskDef instances from a PipelineNode template.
func expandNode(node PipelineNode, surfaces map[string]string, executionOrder []string) []AutoGenTaskDef {
	singleSurface := isSingleSurface(surfaces)
	switch node.Expansion {
	case "per-surface-key":
		return expandPerSurfaceKey(node, surfaces, singleSurface, executionOrder)
	case "per-surface-type":
		return expandPerSurfaceType(node, surfaces)
	default:
		key := deriveKey(node.Key, node.ID)
		return []AutoGenTaskDef{
			{
				ID: node.ID, Key: key, Title: node.Title, Priority: node.Priority,
				EstimatedTime: node.EstimatedTime, Type: node.Type,
				MainSession: node.MainSession, Breaking: node.Breaking, StrategyKind: node.StrategyKind,
			},
		}
	}
}

// expandPerSurfaceKey creates one task per surface key. Serial chain wiring.
func expandPerSurfaceKey(node PipelineNode, surfaces map[string]string, singleSurface bool, executionOrder []string) []AutoGenTaskDef {
	if singleSurface {
		for key, typ := range surfaces {
			stripID := strings.ReplaceAll(node.ID, "-{surface-key}", "")
			stripKey := strings.ReplaceAll(node.Key, "-{surface-key}", "")
			return []AutoGenTaskDef{{
				ID: stripID, Key: deriveKey(stripKey, stripID), Title: expandTitle(node.Title, typ),
				Priority: node.Priority, EstimatedTime: node.EstimatedTime, Type: node.Type,
				MainSession: node.MainSession, Breaking: node.Breaking, SurfaceKey: key, SurfaceType: typ,
				StrategyKind: node.StrategyKind,
			}}
		}
	}

	var keys []string
	if len(executionOrder) > 0 {
		keys = executionOrder
	} else {
		keys = sortedSurfaceKeys(surfaces)
	}
	var tasks []AutoGenTaskDef
	for _, key := range keys {
		typ := surfaces[key]
		id := strings.ReplaceAll(node.ID, "{surface-key}", key)
		keyVal := strings.ReplaceAll(node.Key, "{surface-key}", key)
		if keyVal == "" {
			keyVal = deriveKey("", id)
		}
		tasks = append(tasks, AutoGenTaskDef{
			ID: id, Key: keyVal, Title: expandTitle(node.Title, typ),
			Priority: node.Priority, EstimatedTime: node.EstimatedTime, Type: node.Type,
			MainSession: node.MainSession, Breaking: node.Breaking, SurfaceKey: key, SurfaceType: typ,
			StrategyKind: node.StrategyKind,
		})
	}
	return tasks
}

// expandPerSurfaceType creates one task per unique surface type (parallel).
func expandPerSurfaceType(node PipelineNode, surfaces map[string]string) []AutoGenTaskDef {
	seen := make(map[string]bool)
	var tasks []AutoGenTaskDef
	for _, key := range sortedSurfaceKeys(surfaces) {
		typ := surfaces[key]
		if seen[typ] {
			continue
		}
		seen[typ] = true
		id := strings.ReplaceAll(node.ID, "{surface-type}", typ)
		keyVal := strings.ReplaceAll(node.Key, "{surface-type}", typ)
		if keyVal == "" {
			keyVal = deriveKey("", id)
		}
		tasks = append(tasks, AutoGenTaskDef{
			ID: id, Key: keyVal, Title: expandTitle(node.Title, typ),
			Priority: node.Priority, EstimatedTime: node.EstimatedTime, Type: node.Type,
			MainSession: node.MainSession, Breaking: node.Breaking, SurfaceType: typ,
			StrategyKind: node.StrategyKind,
		})
	}
	return tasks
}

// deriveKey derives the index.json key from the node's Key field or from the ID.
func deriveKey(key, id string) string {
	if key != "" {
		return key
	}
	if strings.HasPrefix(id, "T-") {
		return strings.TrimPrefix(id, "T-")
	}
	return id
}

// expandTitle substitutes {test-type-title} in the title template.
func expandTitle(titleTemplate, surfaceType string) string {
	return strings.ReplaceAll(titleTemplate, "{test-type-title}", TestTypeTitle(surfaceType))
}

// sortedSurfaceKeys returns surface map keys in sorted order for deterministic output.
func sortedSurfaceKeys(surfaces map[string]string) []string {
	keys := make([]string, 0, len(surfaces))
	for k := range surfaces {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// executionOrderKeys returns the ordered surface keys for multi-surface expansion.
func executionOrderKeys(executionOrder []string, surfaces map[string]string) []string {
	if len(executionOrder) > 0 {
		return executionOrder
	}
	return sortedSurfaceKeys(surfaces)
}

// pipelineTaskIDs extracts IDs from a slice of AutoGenTaskDef.
func pipelineTaskIDs(tasks []AutoGenTaskDef) []string {
	ids := make([]string, len(tasks))
	for i, t := range tasks {
		ids[i] = t.ID
	}
	return ids
}
