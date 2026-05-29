package task

import (
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
func CondHasTestableTasks(tasks []Task) bool {
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
// PipelineRegistry — the single source of truth for auto-generated tasks
// ---------------------------------------------------------------------------

// PipelineRegistry defines all auto-generated task nodes in declaration order.
// The order determines generation sequence; execution order is determined by DependsOn.
var PipelineRegistry = []PipelineNode{
	// --- Doc Review (generated whenever business tasks include doc-category types) ---
	{
		Type: TypeDocReview, ID: "T-review-doc",
		Title: "Review Documentation Quality", Priority: string(types.PriorityP1), EstimatedTime: "30min",
		ConfigGate: nil, IntentGate: GateAllowAll,
		GenerateCondition: CondHasDocTasks,
		DependsOn:         []DepRef{{Resolve: ResolveDocTasks}},
	},
	// --- Clean Code (declared early for cross-stage dependency resolution;
	//     execution still occurs after all business tasks via ResolveLastBusinessTask) ---
	{
		Type: TypeCleanCode, ID: "T-clean-code",
		Title: "Simplify and Clean Code", Priority: string(types.PriorityP2), EstimatedTime: "20min",
		ConfigGate: GateCleanCode, GenerateCondition: CondAlways,
		DependsOn: []DepRef{{Resolve: ResolveHighestGateOrLastBiz}},
	},
	// --- Test Generation (first test task depends on T-review-doc and T-clean-code) ---
	{
		Type: TypeTestGenJourneys, ID: "T-test-gen-journeys",
		Title: "Generate Test Journeys", Priority: string(types.PriorityP1), EstimatedTime: "20-30min",
		ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks, StrategyKind: "interface",
		DependsOn: []DepRef{
			{Resolve: ResolveIfGenerated("T-review-doc")},
			{Resolve: ResolveIfGenerated("T-clean-code")},
		},
	},
	// --- Eval (breakdown only) ---
	{
		Type: TypeEvalJourney, ID: "T-eval-journey",
		Title: "Evaluate Journey Quality", Priority: string(types.PriorityP1), EstimatedTime: "20-30min",
		ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks, Mode: "breakdown", MainSession: true,
		DependsOn: []DepRef{{Ref: "T-test-gen-journeys"}},
	},
	{
		Type: TypeTestGenContracts, ID: "T-test-gen-contracts",
		Title: "Generate Test Contracts", Priority: string(types.PriorityP1), EstimatedTime: "30-45min",
		ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks, Mode: "breakdown",
		DependsOn: []DepRef{{Ref: "T-eval-journey"}},
	},
	{
		Type: TypeEvalContract, ID: "T-eval-contract",
		Title: "Evaluate Contract Quality", Priority: string(types.PriorityP1), EstimatedTime: "20-30min",
		ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks, Mode: "breakdown", MainSession: true,
		DependsOn: []DepRef{{Ref: "T-test-gen-contracts"}},
	},
	// --- Gen Scripts (per surface type) ---
	{
		Type: TypeTestGenScripts, ID: "T-test-gen-scripts-{surface-type}",
		Title: "Generate {test-type-title} Scripts", Priority: string(types.PriorityP1), EstimatedTime: "1-2h",
		ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks, Mode: "breakdown", Expansion: "per-surface-type",
		DependsOn:    []DepRef{{Ref: "T-eval-contract"}},
		StrategyKind: "generate",
	},
	// --- Run Tests (per surface key) ---
	{
		Type: TypeTestRun, ID: "T-test-run-{surface-key}",
		Title: "Run {test-type-title}", Priority: string(types.PriorityP1), EstimatedTime: "30min-1h",
		ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks,
		DependsOn: []DepRef{{Resolve: ResolveUpstream}},
		Expansion: "per-surface-key", StrategyKind: "run",
	},
	// --- Validation ---
	{
		Type: TypeValidationCode, ID: "T-validate-code",
		Title: "Validate Code Quality", Priority: string(types.PriorityP2), EstimatedTime: "15min",
		ConfigGate: GateValidation, GenerateCondition: CondAlways,
		DependsOn: []DepRef{{Resolve: ResolveLastRunTestOrBusiness}},
	},
	{
		Type: TypeValidationUx, ID: "T-validate-ux",
		Title: "Validate User Experience", Priority: string(types.PriorityP2), EstimatedTime: "15min",
		ConfigGate: GateValidation, GenerateCondition: CondAlways,
		DependsOn:     []DepRef{{Resolve: ResolveLastRunTestOrBusiness}},
		UISurfaceOnly: true, MainSession: true,
	},
	// --- Spec Consolidation/Drift ---
	{
		Type: TypeDocConsolidate, ID: "T-specs-consolidate",
		Title: "Consolidate Specs", Priority: string(types.PriorityP2), EstimatedTime: "20min",
		ConfigGate: GateConsolidateSpecs, GenerateCondition: CondAlways, Mode: "breakdown",
		DependsOn: []DepRef{{Resolve: ResolveLastRunTestOrBusiness}},
	},
	{
		Type: TypeDocDrift, ID: "T-quick-doc-drift",
		Title: "Detect Spec Drift", Priority: string(types.PriorityP2), EstimatedTime: "15min",
		ConfigGate: GateConsolidateSpecs, GenerateCondition: CondAlways, Mode: "quick",
		DependsOn: []DepRef{{Resolve: ResolveLastRunTestOrBusiness}},
	},
}
