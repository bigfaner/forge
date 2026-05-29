package task

import (
	"fmt"
	"sort"
	"strings"

	"forge-cli/pkg/forgeconfig"
	"forge-cli/pkg/types"
)

// ---------------------------------------------------------------------------
// Core types
// ---------------------------------------------------------------------------

// ConfigGateFunc returns true when the auto config enables this node for the given mode.
// mode is "quick" or "breakdown".
type ConfigGateFunc func(mode string, auto forgeconfig.AutoConfig) bool

// IntentGateFunc returns true when the intent permits this node to generate.
// intent is "new-feature", "refactor", or "cleanup".
type IntentGateFunc func(intent string) bool

// GenerateCondFunc returns true when the business task composition permits this node.
type GenerateCondFunc func(tasks []Task) bool

// DepResolveFunc dynamically resolves dependency IDs at generation time.
// Returns nil when the reference cannot be resolved (e.g., no run-test tasks generated,
// no business tasks present). When nil is returned, GenerateTestTasks skips that
// dependency entry — the node is still generated but with one fewer DependsOn entry.
// If ALL dependencies of a node resolve to nil, the node generates with empty
// DependsOn, becoming a pipeline root (no upstream constraint).
type DepResolveFunc func(ctx *GenContext) []string

// GenContext carries state accumulated during pipeline generation.
// Populated progressively as nodes are processed in declaration order.
type GenContext struct {
	Mode           string
	Intent         string
	Surfaces       map[string]string
	ExecutionOrder []string
	Auto           forgeconfig.AutoConfig
	BusinessTasks  []Task
	ExistingTasks  map[string]Task // full index including gates/summaries (populated by caller)
	// Filled during generation as nodes are expanded:
	UpstreamIDs  []string // IDs of the immediately preceding generated node(s)
	RunTestChain []string // IDs of expanded run-test tasks in serial order
	AllGenerated []string // IDs of all nodes generated so far (in order)
}

// PipelineNode defines a single node in the auto-generated task pipeline.
type PipelineNode struct {
	// Type is the task type constant (e.g., TypeTestGenJourneys).
	Type string
	// Key is the map key used in index.json for this task. For expanded nodes,
	// Key is a template with the same placeholders as ID (e.g., "test-run-{surface-key}").
	// When Key is empty, it is derived from ID by stripping the "T-" prefix and
	// lowercasing. This matches the current AutoGenTaskDef.Key convention.
	Key string
	// ID is the task ID or ID template. Placeholders: {surface-key}, {surface-type}.
	ID string
	// Title is the task title or title template.
	Title string
	// Priority is the task priority (P0, P1, P2).
	Priority string
	// EstimatedTime is the task duration estimate.
	EstimatedTime string
	// ConfigGate returns true when the config enables this node. nil = no config gate.
	ConfigGate ConfigGateFunc
	// IntentGate returns true when the intent permits this node. nil = GateBlockSkipTest (default).
	// Use GateAllowAll for nodes that should generate regardless of intent (e.g., T-review-doc).
	IntentGate IntentGateFunc
	// Mode restricts this node to a specific mode. Empty means both modes.
	// "quick" = quick mode only, "breakdown" = breakdown mode only.
	Mode string
	// GenerateCondition returns true when the business task composition permits this node.
	// MUST be explicitly set for every node. No implicit default.
	GenerateCondition GenerateCondFunc
	// DependsOn defines dependency references.
	DependsOn []DepRef
	// Expansion controls how this node expands into multiple tasks.
	// "" (empty) - single task
	// "per-surface-key" - one task per surface key
	// "per-surface-type" - one task per surface type
	Expansion string
	// MainSession indicates this task runs in the main session.
	MainSession bool
	// StrategyKind for task definition ("generate", "run", "interface", "").
	StrategyKind string
	// UISurfaceOnly indicates this node is only generated when at least one surface has a visual UI.
	UISurfaceOnly bool
}

// DepRef represents a dependency reference.
// Use Ref for static IDs, Resolve for dynamic references.
// If Resolve is non-nil, Ref is ignored.
type DepRef struct {
	Ref     string         // concrete task ID (e.g., "T-test-gen-journeys")
	Resolve DepResolveFunc // dynamic resolver; nil = use Ref as-is
}

// ---------------------------------------------------------------------------
// Config Gate functions
// ---------------------------------------------------------------------------

// GateTest returns true when the Test category is enabled for the given mode.
func GateTest(mode string, auto forgeconfig.AutoConfig) bool {
	if mode == "quick" {
		return auto.Test.Quick
	}
	return auto.Test.Full
}

// GateValidation returns true when the Validation category is enabled for the given mode.
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

// ---------------------------------------------------------------------------
// Intent Gate functions
// ---------------------------------------------------------------------------

// GateAllowAll permits all intents. Used by T-review-doc.
func GateAllowAll(_ string) bool { return true }

// GateBlockSkipTest blocks refactor/cleanup intents. Used by all config-gated nodes.
func GateBlockSkipTest(intent string) bool {
	return !isSkipTestIntent(intent)
}

// ---------------------------------------------------------------------------
// Generate Condition functions
// ---------------------------------------------------------------------------

// CondHasTestableTasks returns true when any business task has a testable type.
// When tasks is nil, returns true (legacy compat: old procedural code did not gate on business tasks).
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
// When tasks is nil, returns false (no doc tasks to trigger T-review-doc).
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

// ---------------------------------------------------------------------------
// Dependency Resolver functions
// ---------------------------------------------------------------------------

// ResolveLastRunTest returns the ID of the last task in the run-test expansion chain.
// Returns nil when no run-test tasks have been generated.
var ResolveLastRunTest DepResolveFunc = func(ctx *GenContext) []string {
	if len(ctx.RunTestChain) == 0 {
		return nil
	}
	return []string{ctx.RunTestChain[len(ctx.RunTestChain)-1]}
}

// ResolveUpstream returns the IDs of the immediately preceding generated node(s).
// For single nodes: one ID. For expanded nodes: all expanded IDs of the previous node.
var ResolveUpstream DepResolveFunc = func(ctx *GenContext) []string {
	if len(ctx.UpstreamIDs) == 0 {
		return nil
	}
	return ctx.UpstreamIDs
}

// ResolveDocTasks returns the IDs of all doc-category business tasks.
var ResolveDocTasks DepResolveFunc = func(ctx *GenContext) []string {
	var ids []string
	for _, t := range ctx.BusinessTasks {
		if CategoryForType(t.Type) == CategoryDoc {
			ids = append(ids, t.ID)
		}
	}
	return ids
}

// ResolveLastBusinessTask returns the ID of the highest-numbered business task.
// Used by T-clean-code which must run after all business tasks complete.
// Note: uses numericID for sorting (extracts leading number from task ID), matching
// current findMaxBusinessTaskID behavior.
var ResolveLastBusinessTask DepResolveFunc = func(ctx *GenContext) []string {
	if len(ctx.BusinessTasks) == 0 {
		return nil
	}
	var maxID string
	var maxNum int
	for _, t := range ctx.BusinessTasks {
		num := numericID(t.ID)
		if num > maxNum {
			maxNum = num
			maxID = t.ID
		}
	}
	if maxID == "" {
		return nil
	}
	return []string{maxID}
}

// ResolveHighestGateOrLastBiz returns the ID of the highest-phase gate/summary,
// or the last business task if its phase is higher. Used by T-clean-code in breakdown
// mode to ensure stage-gates gate the test pipeline. Matches current ResolveFirstTestDep behavior.
// Two-pass logic matches current findHighestGateOrSummary: gate priority over summary.
var ResolveHighestGateOrLastBiz DepResolveFunc = func(ctx *GenContext) []string {
	// Pass 1: find highest-phase gate (gate priority)
	var dep string
	var depPhase int
	for _, t := range ctx.ExistingTasks {
		if strings.HasSuffix(t.ID, IDSuffixGate) {
			p := phaseFromID(t.ID)
			if p > depPhase {
				depPhase = p
				dep = t.ID
			}
		}
	}
	// Pass 2: if no gate found, fall back to highest-phase summary
	if dep == "" {
		for _, t := range ctx.ExistingTasks {
			if strings.HasSuffix(t.ID, IDSuffixSummary) {
				p := phaseFromID(t.ID)
				if p > depPhase {
					depPhase = p
					dep = t.ID
				}
			}
		}
	}
	// Compare with last business task phase
	var maxBizID string
	var maxBizPhase int
	for _, t := range ctx.BusinessTasks {
		p := phaseFromID(t.ID)
		if p > maxBizPhase {
			maxBizPhase = p
			maxBizID = t.ID
		}
	}
	if maxBizID != "" && maxBizPhase > depPhase {
		dep = maxBizID
	}
	if dep == "" {
		return nil
	}
	return []string{dep}
}

// ResolveLastRunTestOrBusiness returns the last run-test task ID when test pipeline
// is active, otherwise falls back to the last business task ID.
// Used by drift/consolidate/validation nodes that currently depend on ResolveDriftFallbackDep.
var ResolveLastRunTestOrBusiness DepResolveFunc = func(ctx *GenContext) []string {
	if len(ctx.RunTestChain) > 0 {
		return []string{ctx.RunTestChain[len(ctx.RunTestChain)-1]}
	}
	if len(ctx.BusinessTasks) > 0 {
		return []string{ctx.BusinessTasks[len(ctx.BusinessTasks)-1].ID}
	}
	return nil
}

// ResolveIfGenerated returns the task ID if it was already generated, nil otherwise.
// Used for cross-stage dependencies where the target node appears earlier in declaration order.
// Init-time validation ensures the referenced ID belongs to a node declared before the caller.
func ResolveIfGenerated(id string) DepResolveFunc {
	return func(ctx *GenContext) []string {
		for _, generated := range ctx.AllGenerated {
			if generated == id {
				return []string{id}
			}
		}
		return nil
	}
}

// ---------------------------------------------------------------------------
// GenerateTestTasks — registry-driven task generation
// ---------------------------------------------------------------------------

// GenerateTestTasks filters PipelineRegistry by mode/config/intent/condition/ui constraints,
// expands per-surface nodes, resolves dependencies via GenContext progressive population,
// and returns the generated AutoGenTaskDef list.
//
// Implements the 5-step algorithm: filter -> expand -> resolve -> update -> return.
func GenerateTestTasks(mode string, surfaces map[string]string, executionOrder []string, auto forgeconfig.AutoConfig, intent string, businessTasks []Task, existingTasks map[string]Task) []AutoGenTaskDef {
	// Do NOT early-return on empty surfaces — non-surface tasks (T-review-doc, T-clean-code,
	// T-validate-code) can still generate when surfaces is empty. Surface-dependent nodes
	// (per-surface-key/type expansion) naturally produce zero expanded tasks.

	ctx := &GenContext{
		Mode:           mode,
		Intent:         intent,
		Surfaces:       surfaces,
		ExecutionOrder: executionOrder,
		Auto:           auto,
		BusinessTasks:  businessTasks,
		ExistingTasks:  existingTasks,
	}

	var generated []AutoGenTaskDef

	for _, node := range PipelineRegistry {
		// Step 1: Filter — apply all 5 gate conditions
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

		// Step 2: Expand — produce concrete task(s) from template
		expanded := expandNode(node, surfaces, executionOrder)

		// Step 3: Resolve dependencies for each expanded task
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
		}

		// Step 4: Update GenContext (progressive population)
		ids := pipelineTaskIDs(expanded)
		ctx.AllGenerated = append(ctx.AllGenerated, ids...)
		ctx.UpstreamIDs = ids
		if node.Type == TypeTestRun {
			ctx.RunTestChain = append(ctx.RunTestChain, ids...)
		}
		generated = append(generated, expanded...)
	}

	// Phase 2: Dynamic validation - verify generated task set.
	// Checks: all resolver-returned IDs exist in generated task set; no circular dependencies.
	// Deliberately not returning errors from GenerateTestTasks to preserve
	// backward compatibility. Phase 2 errors indicate programming bugs in registry configuration.
	_ = validateGeneratedTasks(generated)

	return generated
}

// hasVisualUI returns true when at least one surface has a visual UI type
// (TUI, Web, or Mobile).
func hasVisualUI(surfaces map[string]string) bool {
	for _, typ := range surfaces {
		if uiSurfaceTypes[types.SurfaceType(typ)] {
			return true
		}
	}
	return false
}

// expandNode produces concrete AutoGenTaskDef instances from a PipelineNode template.
//
//	"" (empty)         → single task
//	"per-surface-key"  → one task per surface key (serial chain wiring)
//	"per-surface-type" → one task per unique surface type (parallel)
func expandNode(node PipelineNode, surfaces map[string]string, executionOrder []string) []AutoGenTaskDef {
	singleSurface := isSingleSurface(surfaces)

	switch node.Expansion {
	case "per-surface-key":
		return expandPerSurfaceKey(node, surfaces, singleSurface, executionOrder)
	case "per-surface-type":
		return expandPerSurfaceType(node, surfaces)
	default:
		// Single (no expansion)
		key := deriveKey(node.Key, node.ID)
		return []AutoGenTaskDef{
			{
				ID:            node.ID,
				Key:           key,
				Title:         node.Title,
				Priority:      node.Priority,
				EstimatedTime: node.EstimatedTime,
				Type:          node.Type,
				MainSession:   node.MainSession,
				Breaking:      true,
				StrategyKind:  node.StrategyKind,
			},
		}
	}
}

// expandPerSurfaceKey creates one task per surface key.
// When isSingleSurface is true, the surface-key suffix is stripped from ID.
// Serial chain: expanded[i] depends on expanded[i-1] (after node-level deps).
func expandPerSurfaceKey(node PipelineNode, surfaces map[string]string, singleSurface bool, executionOrder []string) []AutoGenTaskDef {
	if singleSurface {
		// Degenerate: single surface, strip suffix
		for key, typ := range surfaces {
			title := expandTitle(node.Title, typ)
			stripID := strings.ReplaceAll(node.ID, "-{surface-key}", "")
			stripKey := strings.ReplaceAll(node.Key, "-{surface-key}", "")
			return []AutoGenTaskDef{
				{
					ID:            stripID,
					Key:           deriveKey(stripKey, stripID),
					Title:         title,
					Priority:      node.Priority,
					EstimatedTime: node.EstimatedTime,
					Type:          node.Type,
					MainSession:   node.MainSession,
					Breaking:      true,
					SurfaceKey:    key,
					SurfaceType:   typ,
					StrategyKind:  node.StrategyKind,
				},
			}
		}
	}

	// Multi-surface: expand by execution order (provided by caller).
	// Fall back to sorted keys when execution order is not available.
	var keys []string
	if len(executionOrder) > 0 {
		keys = executionOrder
	} else {
		keys = sortedSurfaceKeys(surfaces)
	}
	var tasks []AutoGenTaskDef
	for _, key := range keys {
		typ := surfaces[key]
		title := expandTitle(node.Title, typ)
		id := strings.ReplaceAll(node.ID, "{surface-key}", key)
		keyVal := strings.ReplaceAll(node.Key, "{surface-key}", key)
		if keyVal == "" {
			keyVal = deriveKey("", id)
		}
		tasks = append(tasks, AutoGenTaskDef{
			ID:            id,
			Key:           keyVal,
			Title:         title,
			Priority:      node.Priority,
			EstimatedTime: node.EstimatedTime,
			Type:          node.Type,
			MainSession:   node.MainSession,
			Breaking:      true,
			SurfaceKey:    key,
			SurfaceType:   typ,
			StrategyKind:  node.StrategyKind,
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
		title := expandTitle(node.Title, typ)
		id := strings.ReplaceAll(node.ID, "{surface-type}", typ)
		keyVal := strings.ReplaceAll(node.Key, "{surface-type}", typ)
		if keyVal == "" {
			keyVal = deriveKey("", id)
		}
		tasks = append(tasks, AutoGenTaskDef{
			ID:            id,
			Key:           keyVal,
			Title:         title,
			Priority:      node.Priority,
			EstimatedTime: node.EstimatedTime,
			Type:          node.Type,
			MainSession:   node.MainSession,
			Breaking:      true,
			SurfaceType:   typ,
			StrategyKind:  node.StrategyKind,
		})
	}
	return tasks
}

// deriveKey derives the index.json key from the node's Key field or from the ID.
// When key is empty, it is derived from ID by stripping "T-" prefix and lowercasing.
func deriveKey(key, id string) string {
	if key != "" {
		return key
	}
	if strings.HasPrefix(id, "T-") {
		return strings.TrimPrefix(id, "T-")
	}
	return id
}

// expandTitle substitutes {test-type-title} in the title template with the
// TestTypeTitle for the given surface type.
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

// pipelineTaskIDs extracts IDs from a slice of AutoGenTaskDef.
func pipelineTaskIDs(tasks []AutoGenTaskDef) []string {
	ids := make([]string, len(tasks))
	for i, t := range tasks {
		ids[i] = t.ID
	}
	return ids
}

// ---------------------------------------------------------------------------
// Registry-derived lookup functions
// ---------------------------------------------------------------------------

// matchRegistryID attempts to match a task ID against registry node ID patterns.
// Returns the node's Type if matched, or "" if no match.
// Handles: exact match, per-surface-type suffix match, per-surface-key with
// surfaces validation, and single-surface degenerate IDs.
func matchRegistryID(id string, surfaces map[string]string) string {
	for _, node := range PipelineRegistry {
		if node.Expansion == "" {
			// Exact match
			if id == node.ID {
				return node.Type
			}
			continue
		}

		switch node.Expansion {
		case "per-surface-type":
			// Match prefix + "-" + type suffix (e.g., "T-test-gen-scripts-api")
			// Also matches degenerate form (no suffix) for backward compat.
			if matchTypeSuffixedID(id, node.ID) {
				return node.Type
			}
		case "per-surface-key":
			// Match prefix + "-" + surface-key suffix
			if matched := matchSurfaceKeyID(id, node.ID, surfaces); matched {
				return node.Type
			}
		}
	}
	return ""
}

// matchTypeSuffixedID checks if id matches the pattern "baseTemplate" with
// {surface-type} replaced by a concrete type suffix.
// e.g., matchTypeSuffixedID("T-test-gen-scripts-api", "T-test-gen-scripts-{surface-type}")
// returns true.
// Also accepts the degenerate form (ID without type suffix) to maintain backward
// compatibility with InferType's original exact match for IDs like "T-test-gen-scripts".
func matchTypeSuffixedID(id, idTemplate string) bool {
	// Find the placeholder in the template
	placeholder := "{surface-type}"
	idx := strings.Index(idTemplate, placeholder)
	if idx < 0 {
		return false
	}
	prefix := idTemplate[:idx]

	// Degenerate form: ID equals prefix minus trailing "-" (e.g., "T-test-gen-scripts"
	// matches "T-test-gen-scripts-{surface-type}")
	stripPrefix := strings.TrimSuffix(prefix, "-")
	if id == stripPrefix {
		return true
	}

	if !strings.HasPrefix(id, prefix) {
		return false
	}
	rem := id[len(prefix):]
	if len(rem) == 0 {
		return false // no suffix after prefix
	}
	// Validate: suffix must be lowercase letters and hyphens
	for _, c := range rem {
		if (c < 'a' || c > 'z') && c != '-' {
			return false
		}
	}
	return true
}

// matchSurfaceKeyID checks if id matches the pattern "baseTemplate" with
// {surface-key} replaced by a surface key from the surfaces map.
// Also handles single-surface degenerate case (no suffix) — accepted regardless
// of surfaces to maintain backward compatibility with InferType's original exact match.
func matchSurfaceKeyID(id, idTemplate string, surfaces map[string]string) bool {
	placeholder := "{surface-key}"
	idx := strings.Index(idTemplate, placeholder)
	if idx < 0 {
		return false
	}
	prefix := idTemplate[:idx]

	// Single-surface degenerate: the template minus "-{surface-key}" should equal id.
	// Accepted regardless of whether surfaces is nil/single/multi — matches the
	// original InferType exact-match behavior for IDs like "T-test-run".
	stripTemplate := strings.ReplaceAll(idTemplate, "-{surface-key}", "")
	if id == stripTemplate {
		return true
	}

	if !strings.HasPrefix(id, prefix) {
		return false
	}

	// Multi-surface: extract suffix and check if it's a known surface key
	suffix := id[len(prefix):]
	if suffix == "" {
		return false
	}
	_, ok := surfaces[suffix]
	return ok
}

// ---------------------------------------------------------------------------
// PipelineRegistry — the single source of truth for auto-generated tasks
// ---------------------------------------------------------------------------

// PipelineRegistry defines all auto-generated task nodes in declaration order.
// The order determines generation sequence; execution order is determined by DependsOn.
var PipelineRegistry = []PipelineNode{
	// --- Doc Review (generated whenever business tasks include doc-category types) ---
	{
		Type: TypeDocReview, Key: "review-doc", ID: "T-review-doc",
		Title: "Review Documentation Quality", Priority: string(types.PriorityP1), EstimatedTime: "30min",
		ConfigGate: nil, IntentGate: GateAllowAll,
		GenerateCondition: CondHasDocTasks,
		DependsOn:         []DepRef{{Resolve: ResolveDocTasks}},
	},
	// --- Clean Code (declared early for cross-stage dependency resolution;
	//     execution still occurs after all business tasks via ResolveLastBusinessTask) ---
	{
		Type: TypeCleanCode, Key: "clean-code", ID: "T-clean-code",
		Title: "Simplify and Clean Code", Priority: string(types.PriorityP2), EstimatedTime: "20min",
		ConfigGate: GateCleanCode, GenerateCondition: CondAlways,
		DependsOn: []DepRef{{Resolve: ResolveHighestGateOrLastBiz}},
	},
	// --- Test Generation (first test task depends on T-review-doc and T-clean-code) ---
	{
		Type: TypeTestGenJourneys, Key: "gen-journeys", ID: "T-test-gen-journeys",
		Title: "Generate Test Journeys", Priority: string(types.PriorityP1), EstimatedTime: "20-30min",
		ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks, StrategyKind: "interface",
		DependsOn: []DepRef{
			{Resolve: ResolveIfGenerated("T-review-doc")},
			{Resolve: ResolveIfGenerated("T-clean-code")},
		},
	},
	// --- Eval (breakdown only) ---
	{
		Type: TypeEvalJourney, Key: "eval-journey", ID: "T-eval-journey",
		Title: "Evaluate Journey Quality", Priority: string(types.PriorityP1), EstimatedTime: "20-30min",
		ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks, Mode: "breakdown", MainSession: true,
		DependsOn: []DepRef{{Ref: "T-test-gen-journeys"}},
	},
	{
		Type: TypeTestGenContracts, Key: "gen-contracts", ID: "T-test-gen-contracts",
		Title: "Generate Test Contracts", Priority: string(types.PriorityP1), EstimatedTime: "30-45min",
		ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks, Mode: "breakdown",
		DependsOn: []DepRef{{Ref: "T-eval-journey"}},
	},
	{
		Type: TypeEvalContract, Key: "eval-contract", ID: "T-eval-contract",
		Title: "Evaluate Contract Quality", Priority: string(types.PriorityP1), EstimatedTime: "20-30min",
		ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks, Mode: "breakdown", MainSession: true,
		DependsOn: []DepRef{{Ref: "T-test-gen-contracts"}},
	},
	// --- Gen Scripts (per surface type) ---
	{
		Type: TypeTestGenScripts, Key: "gen-test-scripts-{surface-type}", ID: "T-test-gen-scripts-{surface-type}",
		Title: "Generate {test-type-title} Scripts", Priority: string(types.PriorityP1), EstimatedTime: "1-2h",
		ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks, Mode: "breakdown", Expansion: "per-surface-type",
		DependsOn:    []DepRef{{Ref: "T-eval-contract"}},
		StrategyKind: "generate",
	},
	// --- Run Tests (per surface key) ---
	{
		Type: TypeTestRun, Key: "run-test-{surface-key}", ID: "T-test-run-{surface-key}",
		Title: "Run {test-type-title}", Priority: string(types.PriorityP1), EstimatedTime: "30min-1h",
		ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks,
		DependsOn: []DepRef{{Resolve: ResolveUpstream}},
		Expansion: "per-surface-key", StrategyKind: "run",
	},
	// --- Validation ---
	{
		Type: TypeValidationCode, Key: "validate-code", ID: "T-validate-code",
		Title: "Validate Code Quality", Priority: string(types.PriorityP2), EstimatedTime: "15min",
		ConfigGate: GateValidation, GenerateCondition: CondAlways,
		DependsOn: []DepRef{{Resolve: ResolveLastRunTestOrBusiness}},
	},
	{
		Type: TypeValidationUx, Key: "validate-ux", ID: "T-validate-ux",
		Title: "Validate User Experience", Priority: string(types.PriorityP2), EstimatedTime: "15min",
		ConfigGate: GateValidation, GenerateCondition: CondAlways,
		DependsOn:     []DepRef{{Resolve: ResolveLastRunTestOrBusiness}},
		UISurfaceOnly: true, MainSession: true,
	},
	// --- Spec Consolidation/Drift ---
	{
		Type: TypeDocConsolidate, Key: "consolidate-specs", ID: "T-specs-consolidate",
		Title: "Consolidate Specs", Priority: string(types.PriorityP2), EstimatedTime: "20min",
		ConfigGate: GateConsolidateSpecs, GenerateCondition: CondAlways, Mode: "breakdown",
		DependsOn: []DepRef{{Resolve: ResolveLastRunTestOrBusiness}},
	},
	{
		Type: TypeDocDrift, Key: "quick-drift-detection", ID: "T-quick-doc-drift",
		Title: "Detect Spec Drift", Priority: string(types.PriorityP2), EstimatedTime: "15min",
		ConfigGate: GateConsolidateSpecs, GenerateCondition: CondAlways, Mode: "quick",
		DependsOn: []DepRef{{Resolve: ResolveLastRunTestOrBusiness}},
	},
}

// ---------------------------------------------------------------------------
// Two-phase validation
// ---------------------------------------------------------------------------

// escapeHatchLimit is the maximum allowed number of post-generation injection
// functions. Current count: 0. See proposal "Escape hatch protocol" for rationale.
const escapeHatchLimit = 5

func init() {
	if err := ValidatePipelineRegistry(); err != nil {
		panic("pipeline registry validation failed: " + err.Error())
	}
}

// ValidatePipelineRegistry performs Phase 1 (static) validation of the pipeline
// registry. It runs at init-time and validates structural invariants that can be
// checked without runtime state. Panics on failure with actionable error messages.
func ValidatePipelineRegistry() error {
	// Collect all node IDs (including template IDs) for reference checking.
	nodeIDs := make(map[string]int) // ID -> declaration index
	for i, node := range PipelineRegistry {
		// Check 1: GenerateCondition must be non-nil
		if node.GenerateCondition == nil {
			return fmt.Errorf("node %q (index %d): GenerateCondition must be non-nil, use CondAlways for unconditional generation", node.ID, i)
		}

		// Check 2: Key/ID template placeholders must match Expansion setting
		if err := validatePlaceholders(node, i); err != nil {
			return err
		}

		// Check 3: DependsOn.Ref strings must reference existing node IDs
		for _, dep := range node.DependsOn {
			if dep.Resolve == nil && dep.Ref != "" {
				if _, ok := nodeIDs[dep.Ref]; !ok {
					// Not found among previously-declared nodes — check full registry
					if !idExistsInRegistry(dep.Ref) {
						return fmt.Errorf("node %q (index %d): DependsOn.Ref %q does not match any registry node ID", node.ID, i, dep.Ref)
					}
				}
			}

			// Check 4: ResolveIfGenerated references must point to nodes declared before the caller
			if dep.Resolve != nil && isResolveIfGenerated(dep.Resolve) {
				refID := extractResolveIfGeneratedID(dep.Resolve)
				if refID != "" {
					if _, ok := nodeIDs[refID]; !ok {
						return fmt.Errorf("node %q (index %d): ResolveIfGenerated(%q) references a node not yet declared — must appear before this node in declaration order", node.ID, i, refID)
					}
				}
			}
		}

		nodeIDs[node.ID] = i
	}

	// Check 5: All expanded IDs must be unique
	if err := validateExpandedIDsUnique(); err != nil {
		return err
	}

	// Check 6: Escape hatch count <= 5
	// Currently hardcoded as 0 post-generation injection functions.
	// When escape hatches are added, increment this counter.
	escapeHatchCount := 0
	if escapeHatchCount > escapeHatchLimit {
		return fmt.Errorf("escape hatch count %d exceeds limit %d — extend registry expressiveness instead", escapeHatchCount, escapeHatchLimit)
	}

	// Check 7: Ordering invariants
	if err := validateOrderingInvariants(); err != nil {
		return err
	}

	return nil
}

// validatePlaceholders checks that Key/ID template placeholders match the Expansion setting.
func validatePlaceholders(node PipelineNode, idx int) error {
	hasSurfaceKey := strings.Contains(node.ID, "{surface-key}") || strings.Contains(node.Key, "{surface-key}")
	hasSurfaceType := strings.Contains(node.ID, "{surface-type}") || strings.Contains(node.Key, "{surface-type}")

	switch node.Expansion {
	case "per-surface-key":
		if !hasSurfaceKey {
			return fmt.Errorf("node %q (index %d): Expansion is per-surface-key but ID/Key lacks {surface-key} placeholder", node.ID, idx)
		}
	case "per-surface-type":
		if !hasSurfaceType {
			return fmt.Errorf("node %q (index %d): Expansion is per-surface-type but ID/Key lacks {surface-type} placeholder", node.ID, idx)
		}
	default:
		// Single (no expansion) — must NOT have placeholders
		if hasSurfaceKey || hasSurfaceType {
			return fmt.Errorf("node %q (index %d): Expansion is empty but ID/Key contains placeholders {surface-key} or {surface-type}", node.ID, idx)
		}
	}
	return nil
}

// idExistsInRegistry checks if a concrete ID matches any registry node ID.
func idExistsInRegistry(id string) bool {
	for _, node := range PipelineRegistry {
		if node.ID == id {
			return true
		}
		// For template IDs, check if the query matches the pattern
		if strings.Contains(node.ID, "{") {
			if matchIDToTemplate(id, node.ID) {
				return true
			}
		}
	}
	return false
}

// matchIDToTemplate checks if a concrete ID matches a template with placeholders.
func matchIDToTemplate(id, template string) bool {
	parts := strings.SplitN(template, "{", 2)
	if len(parts) < 2 {
		return id == template
	}
	prefix := parts[0]
	if !strings.HasPrefix(id, prefix) {
		return false
	}
	return len(id) > len(prefix)
}

// isResolveIfGenerated checks if a DepResolveFunc is a ResolveIfGenerated wrapper.
func isResolveIfGenerated(fn DepResolveFunc) bool {
	return extractResolveIfGeneratedID(fn) != ""
}

// extractResolveIfGeneratedID attempts to extract the target ID from a ResolveIfGenerated
// closure by testing it against all registry node IDs.
func extractResolveIfGeneratedID(fn DepResolveFunc) string {
	// Create a context with all registry node IDs in AllGenerated
	allIDs := make([]string, 0, len(PipelineRegistry))
	for _, node := range PipelineRegistry {
		allIDs = append(allIDs, node.ID)
	}

	// Test with all IDs present — ResolveIfGenerated returns [id] if id is in AllGenerated
	ctx := &GenContext{AllGenerated: allIDs}
	result := fn(ctx)
	if len(result) == 1 {
		// Now test with no IDs — should return nil
		ctx2 := &GenContext{AllGenerated: nil}
		result2 := fn(ctx2)
		if result2 == nil {
			return result[0]
		}
	}
	return ""
}

// validateExpandedIDsUnique verifies that all expanded IDs are unique
// using representative surface configurations.
func validateExpandedIDsUnique() error {
	testSurfaces := map[string]string{
		".":      "api",
		"api":    "api",
		"cli":    "cli",
		"tui":    "tui",
		"web":    "web",
		"mobile": "mobile",
	}

	seen := make(map[string]string) // ID -> node Type that produced it
	for _, node := range PipelineRegistry {
		expanded := expandNode(node, testSurfaces, nil)
		for _, t := range expanded {
			if prev, dup := seen[t.ID]; dup {
				return fmt.Errorf("duplicate expanded ID %q produced by both %q and %q", t.ID, prev, node.Type)
			}
			seen[t.ID] = node.Type
		}
	}

	// Also test single-surface degenerate case
	singleSurfaces := map[string]string{".": "api"}
	singleSeen := make(map[string]string)
	for _, node := range PipelineRegistry {
		expanded := expandNode(node, singleSurfaces, nil)
		for _, t := range expanded {
			if prev, dup := singleSeen[t.ID]; dup {
				return fmt.Errorf("duplicate expanded ID %q (single-surface) produced by both %q and %q", t.ID, prev, node.Type)
			}
			singleSeen[t.ID] = node.Type
		}
	}

	return nil
}

// validateOrderingInvariants checks that resolver dependencies match declaration order.
func validateOrderingInvariants() error {
	var hasUpstreamProducer bool
	var hasRunTestProducer bool

	for i, node := range PipelineRegistry {
		for _, dep := range node.DependsOn {
			if dep.Resolve == nil {
				continue
			}

			// ResolveUpstream: needs at least one prior node that populates UpstreamIDs
			if isSameResolver(dep.Resolve, ResolveUpstream) {
				if !hasUpstreamProducer {
					return fmt.Errorf("node %q (index %d): uses ResolveUpstream but has no guaranteed upstream producer before it in declaration order", node.ID, i)
				}
			}

			// ResolveLastRunTest: needs at least one prior TypeTestRun node
			if isSameResolver(dep.Resolve, ResolveLastRunTest) {
				if !hasRunTestProducer {
					return fmt.Errorf("node %q (index %d): uses ResolveLastRunTest but no TypeTestRun node declared before it", node.ID, i)
				}
			}
		}

		// After processing this node, mark it as a potential producer
		if node.Expansion == "" || node.Expansion == "per-surface-key" {
			hasUpstreamProducer = true
		}
		if node.Type == TypeTestRun {
			hasRunTestProducer = true
		}
	}

	return nil
}

// isSameResolver checks if two DepResolveFunc values are the same function pointer.
func isSameResolver(a, b DepResolveFunc) bool {
	return fmt.Sprintf("%p", a) == fmt.Sprintf("%p", b)
}

// validateGeneratedTasks performs Phase 2 (dynamic) validation of generated tasks.
// Runs at the end of GenerateTestTasks to validate runtime invariants.
// Validates: all resolver-returned IDs exist in generated task set; no circular dependencies.
func validateGeneratedTasks(tasks []AutoGenTaskDef) error {
	// Build ID set for lookup
	idSet := make(map[string]bool, len(tasks))
	for _, t := range tasks {
		idSet[t.ID] = true
	}

	// Check 1: All dependency references point to existing IDs
	for _, t := range tasks {
		for _, dep := range t.Dependencies {
			if !idSet[dep] {
				// Dependencies might reference business tasks (not in generated set).
				// Only flag references that look like pipeline task IDs (T- prefix).
				if strings.HasPrefix(dep, "T-") {
					return fmt.Errorf("generated task %q depends on %q which is not in the generated task set", t.ID, dep)
				}
			}
		}
	}

	// Check 2: No circular dependencies (topological sort)
	if err := checkNoCycles(tasks); err != nil {
		return err
	}

	return nil
}

// checkNoCycles performs topological sort to detect circular dependencies.
func checkNoCycles(tasks []AutoGenTaskDef) error {
	idSet := make(map[string]bool, len(tasks))
	for _, t := range tasks {
		idSet[t.ID] = true
	}

	// Build forward adjacency: dependency -> dependent (for topo sort)
	forward := make(map[string][]string)
	inDegree := make(map[string]int)
	for _, t := range tasks {
		inDegree[t.ID] = 0
	}
	for _, t := range tasks {
		for _, dep := range t.Dependencies {
			if idSet[dep] {
				forward[dep] = append(forward[dep], t.ID)
				inDegree[t.ID]++
			}
		}
	}

	// Enqueue nodes with in-degree 0
	queue := make([]string, 0)
	for _, t := range tasks {
		if inDegree[t.ID] == 0 {
			queue = append(queue, t.ID)
		}
	}

	visited := 0
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		visited++
		for _, next := range forward[curr] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	if visited < len(tasks) {
		return fmt.Errorf("circular dependency detected among %d pipeline tasks (topological sort visited %d of %d)", len(tasks)-visited, visited, len(tasks))
	}

	return nil
}
