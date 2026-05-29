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

// ConfigGateFunc determines whether a PipelineNode should be included based on
// the auto-configuration for the current mode (quick vs full).
// Returns true when the corresponding category is enabled for the given mode.
type ConfigGateFunc func(auto forgeconfig.AutoConfig, mode string) bool

// IntentGateFunc determines whether a PipelineNode is allowed for the given
// pipeline intent (e.g. "feature", "refactor", "cleanup").
type IntentGateFunc func(intent string) bool

// GenerateCondFunc determines whether a node should be generated based on the
// generation context (surfaces, types, etc.).
type GenerateCondFunc func(ctx GenContext) bool

// DepResolveFunc resolves a dependency reference into a concrete task ID using
// the full task index. Returns empty string when the dependency target is not
// found (meaning this dependency link is skipped).
type DepResolveFunc func(ref DepRef, ctx GenContext) string

// DepRef represents a single dependency reference within a PipelineNode.
// Resolver is the strategy for turning this reference into a concrete task ID.
type DepRef struct {
	// TaskTemplate is the ID pattern of the dependency (e.g. "T-review-doc").
	// For per-surface-key tasks, use the base prefix (e.g. "T-test-run").
	TaskTemplate string
	// Resolver resolves this dependency reference to a concrete task ID.
	Resolver DepResolveFunc
}

// GenContext carries the runtime context needed for dependency resolution and
// generation-condition evaluation.
type GenContext struct {
	// ExistingTasks is the full task index (including gates/summaries/auto-gen).
	// Populated by the caller before any resolution.
	ExistingTasks map[string]Task
	// Auto is the auto-configuration from .forge/config.yaml.
	Auto forgeconfig.AutoConfig
	// Mode is the pipeline mode: "quick" or "breakdown".
	Mode string
	// Intent is the pipeline intent: "feature", "refactor", "cleanup", etc.
	Intent string
	// SurfaceTypes is the deduplicated list of surface types in the project.
	SurfaceTypes []string
	// Surfaces is the surface-key-to-type map from config.
	Surfaces map[string]string
	// ExecutionOrder is the resolved execution order of surface keys.
	ExecutionOrder []string
	// GeneratedTasks is the list of tasks generated so far in the current run.
	// Used by ResolveIfGenerated to check whether a dependency target was generated.
	GeneratedTasks []AutoGenTaskDef
}

// PipelineNode describes a single auto-generated task in the pipeline topology.
// The PipelineRegistry is the single source of truth for all auto-generated task
// definitions; BuildIndex consumes it to produce the final task set.
type PipelineNode struct {
	// ID is the task identifier (e.g. "T-review-doc").
	// For per-surface-key nodes, this is the base ID (e.g. "T-test-run").
	ID string
	// Key is the index map key for single-instance tasks (e.g. "review-doc").
	// Empty for per-surface-key nodes (key is derived from surface key).
	Key string
	// Title is the task title. May contain format verbs for per-surface-key tasks.
	Title string
	// Type is the task type constant (e.g. TypeDocReview).
	Type string
	// Priority is the task priority (P0/P1/P2).
	Priority string
	// EstimatedTime is the estimated task duration.
	EstimatedTime string
	// MainSession controls whether the task runs in the main session.
	MainSession bool
	// ConfigGate controls whether this node is included based on auto config.
	ConfigGate ConfigGateFunc
	// IntentGate controls whether this node is blocked by the pipeline intent.
	IntentGate IntentGateFunc
	// GenerateCond controls whether this node is generated based on context.
	// When nil, the node is always generated (if ConfigGate and IntentGate pass).
	GenerateCond GenerateCondFunc
	// DependsOn lists the dependency references for this node.
	// Resolvers are called in order; each produces one concrete dependency ID.
	DependsOn []DepRef
	// PerSurfaceKey indicates whether this node is expanded per surface key.
	// When true, one task instance is created per surface key in execution order.
	// The first instance uses DependsOn resolvers; subsequent instances depend on
	// the previous expanded task (serial chain).
	PerSurfaceKey bool
	// StrategyKind is the strategy kind for auto-gen task definitions
	// ("generate", "run", "interface", or "" for generic).
	StrategyKind string
}

// ---------------------------------------------------------------------------
// ConfigGate functions
// ---------------------------------------------------------------------------

// GateTest returns true when the Test category is enabled for the given mode.
func GateTest(auto forgeconfig.AutoConfig, mode string) bool {
	switch mode {
	case "quick":
		return auto.Test.Quick
	case "breakdown":
		return auto.Test.Full
	default:
		return false
	}
}

// GateValidation returns true when the Validation category is enabled for the given mode.
func GateValidation(auto forgeconfig.AutoConfig, mode string) bool {
	switch mode {
	case "quick":
		return auto.Validation.Quick
	case "breakdown":
		return auto.Validation.Full
	default:
		return false
	}
}

// GateConsolidateSpecs returns true when ConsolidateSpecs is enabled for the given mode.
func GateConsolidateSpecs(auto forgeconfig.AutoConfig, mode string) bool {
	switch mode {
	case "quick":
		return auto.ConsolidateSpecs.Quick
	case "breakdown":
		return auto.ConsolidateSpecs.Full
	default:
		return false
	}
}

// GateCleanCode returns true when CleanCode is enabled for the given mode.
func GateCleanCode(auto forgeconfig.AutoConfig, mode string) bool {
	switch mode {
	case "quick":
		return auto.CleanCode.Quick
	case "breakdown":
		return auto.CleanCode.Full
	default:
		return false
	}
}

// ---------------------------------------------------------------------------
// IntentGate functions
// ---------------------------------------------------------------------------

// GateAllowAll allows all intents — the node is never blocked by intent.
func GateAllowAll(string) bool { return true }

// GateBlockSkipTest blocks the node when the intent is "refactor" or "cleanup"
// (intents that skip the test pipeline).
func GateBlockSkipTest(intent string) bool {
	return !isSkipTestIntent(intent)
}

// ---------------------------------------------------------------------------
// GenerateCondition functions
// ---------------------------------------------------------------------------

// CondHasTestableTasks returns true when there are test pipeline tasks to generate.
// This is true when the test config gate is enabled and the intent does not skip tests.
func CondHasTestableTasks(ctx GenContext) bool {
	return GateTest(ctx.Auto, ctx.Mode) && !isSkipTestIntent(ctx.Intent)
}

// CondHasDocTasks returns true when there are doc-category business tasks
// in the existing task index.
func CondHasDocTasks(ctx GenContext) bool {
	for _, t := range ctx.ExistingTasks {
		if isAutoGenForDep(t.ID) {
			continue
		}
		if CategoryForType(t.Type) == CategoryDoc {
			return true
		}
	}
	return false
}

// CondAlways always returns true — the node is always generated (subject to
// ConfigGate and IntentGate).
func CondAlways(GenContext) bool { return true }

// ---------------------------------------------------------------------------
// Dependency resolver functions
// ---------------------------------------------------------------------------

// ResolveIfGenerated resolves to the task ID only if a task with that ID
// (or matching the template for per-surface-key nodes) exists in the
// generated tasks list. Returns empty string otherwise.
func ResolveIfGenerated(ref DepRef, ctx GenContext) string {
	for _, t := range ctx.GeneratedTasks {
		if t.ID == ref.TaskTemplate || strings.HasPrefix(t.ID, ref.TaskTemplate+"-") {
			return t.ID
		}
	}
	// Also check existing tasks (for tasks from prior runs or external generation)
	for id := range ctx.ExistingTasks {
		if id == ref.TaskTemplate || strings.HasPrefix(id, ref.TaskTemplate+"-") {
			return id
		}
	}
	return ""
}

// ResolveHighestGateOrLastBiz implements two-pass gate-priority resolution:
//   - Pass 1: find the highest-phase gate task in the index
//   - Pass 2 (fallback): find the highest-phase summary task
//   - Finally: compare the result against the last business task; prefer whichever
//     is in the higher phase (handles final-phase single-task case with no gate)
//
// This matches the current findHighestGateOrSummary + ResolveFirstTestDep logic.
func ResolveHighestGateOrLastBiz(_ DepRef, ctx GenContext) string {
	// Pass 1: find highest-phase gate
	bestID := ""
	bestPhase := 0
	for _, t := range ctx.ExistingTasks {
		if strings.HasSuffix(t.ID, IDSuffixGate) {
			phase := phaseFromID(t.ID)
			if phase > bestPhase {
				bestPhase = phase
				bestID = t.ID
			}
		}
	}
	if bestID != "" {
		// Check if last business task is in a higher phase
		lastBiz := findMaxBusinessTaskID(ctx.ExistingTasks)
		if lastBiz != "" && phaseFromID(lastBiz) > phaseFromID(bestID) {
			return lastBiz
		}
		return bestID
	}

	// Pass 2: fallback to highest-phase summary
	bestPhase = 0
	for _, t := range ctx.ExistingTasks {
		if strings.HasSuffix(t.ID, IDSuffixSummary) {
			phase := phaseFromID(t.ID)
			if phase > bestPhase {
				bestPhase = phase
				bestID = t.ID
			}
		}
	}

	// Compare with last business task
	lastBiz := findMaxBusinessTaskID(ctx.ExistingTasks)
	if lastBiz != "" && (bestID == "" || phaseFromID(lastBiz) > phaseFromID(bestID)) {
		return lastBiz
	}
	return bestID
}

// ResolveLastBusinessTask resolves to the business task with the highest numeric ID.
// Uses numericID for sorting (not slice order), matching findMaxBusinessTaskID.
func ResolveLastBusinessTask(_ DepRef, ctx GenContext) string {
	return findMaxBusinessTaskID(ctx.ExistingTasks)
}

// ResolveLastRunTestOrBusiness resolves to the last run-test task in the
// generated list; if no run-test exists (e.g. refactor/cleanup intent),
// falls back to the last business task.
func ResolveLastRunTestOrBusiness(_ DepRef, ctx GenContext) string {
	// Find the last run-test task from generated tasks
	var lastRunTest string
	for _, t := range ctx.GeneratedTasks {
		if strings.HasPrefix(t.ID, "T-test-run") {
			lastRunTest = t.ID
		}
	}
	if lastRunTest != "" {
		return lastRunTest
	}

	// Fallback to last business task
	return findMaxBusinessTaskID(ctx.ExistingTasks)
}

// ResolveUpstream resolves to the previous task in a per-surface-key serial chain.
// If idx is 0, it returns the resolved dependency from DependsOn.
// If idx > 0, it returns the ID of the previous surface-key task in the chain.
//
// Note: This resolver is used for per-surface-key chain wiring and requires
// the caller to pass context about the current expansion index. It resolves
// to the TaskTemplate as-is, which the expansion loop overrides for non-first
// instances.
func ResolveUpstream(ref DepRef, _ GenContext) string {
	return ref.TaskTemplate
}

// ResolveDocTasks resolves to all doc-category business task IDs, sorted
// deterministically. Returns empty string (signaling no dependency) when
// there are no doc tasks.
func ResolveDocTasks(_ DepRef, ctx GenContext) string {
	var deps []string
	for _, t := range ctx.ExistingTasks {
		if isAutoGenForDep(t.ID) {
			continue
		}
		if CategoryForType(t.Type) == CategoryDoc {
			deps = append(deps, t.ID)
		}
	}
	sort.Strings(deps)
	if len(deps) == 0 {
		return ""
	}
	return strings.Join(deps, ",")
}

// ---------------------------------------------------------------------------
// PipelineRegistry — the single source of truth for auto-generated tasks
// ---------------------------------------------------------------------------

// PipelineRegistry defines all auto-generated task nodes in declaration order.
// The order determines generation sequence: tasks are emitted in this order,
// with per-surface-key tasks expanded inline.
//
// Declaration order (12 base nodes):
//
//  1. T-review-doc          — doc review (gated by intent: allow-all, gated by doc tasks)
//  2. T-clean-code          — code cleanup (gated by clean-code config)
//  3. T-test-gen-journeys   — generate test journeys (gated by test config + intent)
//  4. T-eval-journey        — evaluate journey quality (breakdown only)
//  5. T-test-gen-contracts  — generate test contracts (breakdown only)
//  6. T-eval-contract       — evaluate contract quality (breakdown only)
//  7. T-test-gen-scripts    — generate test scripts (breakdown only, per-type)
//  8. T-test-run            — run tests (per-surface-key)
//  9. T-validate-code       — validate code quality (gated by validation config)
//  10. T-validate-ux         — validate UX (gated by validation config + UI surface)
//  11. T-specs-consolidate   — consolidate specs (breakdown) / T-quick-doc-drift (quick)
//  12. (node 11 is mode-dependent; quick uses drift, breakdown uses consolidate)
var PipelineRegistry = []PipelineNode{
	// 1. Review documentation quality
	{
		ID:            "T-review-doc",
		Key:           "review-doc",
		Title:         "Review Documentation Quality",
		Type:          TypeDocReview,
		Priority:      string(types.PriorityP1),
		EstimatedTime: "30min",
		IntentGate:    GateAllowAll,
		GenerateCond:  CondHasDocTasks,
		DependsOn: []DepRef{
			{TaskTemplate: "", Resolver: ResolveDocTasks},
		},
	},
	// 2. Clean code
	{
		ID:            "T-clean-code",
		Key:           "clean-code",
		Title:         "Simplify and Clean Code",
		Type:          TypeCleanCode,
		Priority:      string(types.PriorityP2),
		EstimatedTime: "20min",
		ConfigGate:    GateCleanCode,
		IntentGate:    GateAllowAll,
		DependsOn: []DepRef{
			{TaskTemplate: "", Resolver: ResolveHighestGateOrLastBiz},
		},
	},
	// 3. Generate test journeys
	{
		ID:            "T-test-gen-journeys",
		Key:           "gen-journeys",
		Title:         "Generate Test Journeys",
		Type:          TypeTestGenJourneys,
		Priority:      string(types.PriorityP1),
		EstimatedTime: "20-30min",
		ConfigGate:    GateTest,
		IntentGate:    GateBlockSkipTest,
		StrategyKind:  "interface",
		DependsOn: []DepRef{
			{TaskTemplate: "T-review-doc", Resolver: ResolveIfGenerated},
			{TaskTemplate: "T-clean-code", Resolver: ResolveIfGenerated},
		},
	},
	// 4. Evaluate journey quality (breakdown only)
	{
		ID:            "T-eval-journey",
		Key:           "eval-journey",
		Title:         "Evaluate Journey Quality",
		Type:          TypeEvalJourney,
		Priority:      string(types.PriorityP1),
		EstimatedTime: "20-30min",
		MainSession:   true,
		ConfigGate:    GateTest,
		IntentGate:    GateBlockSkipTest,
		DependsOn: []DepRef{
			{TaskTemplate: "T-test-gen-journeys", Resolver: ResolveIfGenerated},
		},
	},
	// 5. Generate test contracts (breakdown only)
	{
		ID:            "T-test-gen-contracts",
		Key:           "gen-contracts",
		Title:         "Generate Test Contracts",
		Type:          TypeTestGenContracts,
		Priority:      string(types.PriorityP1),
		EstimatedTime: "30-45min",
		ConfigGate:    GateTest,
		IntentGate:    GateBlockSkipTest,
		DependsOn: []DepRef{
			{TaskTemplate: "T-eval-journey", Resolver: ResolveIfGenerated},
		},
	},
	// 6. Evaluate contract quality (breakdown only)
	{
		ID:            "T-eval-contract",
		Key:           "eval-contract",
		Title:         "Evaluate Contract Quality",
		Type:          TypeEvalContract,
		Priority:      string(types.PriorityP1),
		EstimatedTime: "20-30min",
		MainSession:   true,
		ConfigGate:    GateTest,
		IntentGate:    GateBlockSkipTest,
		DependsOn: []DepRef{
			{TaskTemplate: "T-test-gen-contracts", Resolver: ResolveIfGenerated},
		},
	},
	// 7. Generate test scripts (per surface type, breakdown only)
	{
		ID:            "T-test-gen-scripts",
		Key:           "gen-test-scripts",
		Title:         fmt.Sprintf("Generate %%s Scripts"),
		Type:          TypeTestGenScripts,
		Priority:      string(types.PriorityP1),
		EstimatedTime: "1-2h",
		ConfigGate:    GateTest,
		IntentGate:    GateBlockSkipTest,
		StrategyKind:  "generate",
		DependsOn: []DepRef{
			{TaskTemplate: "T-eval-contract", Resolver: ResolveIfGenerated},
		},
	},
	// 8. Run tests (per surface key)
	{
		ID:            "T-test-run",
		Key:           "run-test",
		Title:         "Run %ss",
		Type:          TypeTestRun,
		Priority:      string(types.PriorityP1),
		EstimatedTime: "30min-1h",
		PerSurfaceKey: true,
		ConfigGate:    GateTest,
		IntentGate:    GateBlockSkipTest,
		StrategyKind:  "run",
		DependsOn: []DepRef{
			{TaskTemplate: "T-test-gen-scripts", Resolver: ResolveIfGenerated},
		},
	},
	// 9. Validate code quality
	{
		ID:            "T-validate-code",
		Key:           "validate-code",
		Title:         "Validate Code Quality",
		Type:          TypeValidationCode,
		Priority:      string(types.PriorityP2),
		EstimatedTime: "15min",
		ConfigGate:    GateValidation,
		IntentGate:    GateAllowAll,
		DependsOn: []DepRef{
			{TaskTemplate: "T-test-run", Resolver: ResolveLastRunTestOrBusiness},
		},
	},
	// 10. Validate user experience (only when UI surfaces exist)
	{
		ID:            "T-validate-ux",
		Key:           "validate-ux",
		Title:         "Validate User Experience",
		Type:          TypeValidationUx,
		Priority:      string(types.PriorityP2),
		EstimatedTime: "15min",
		MainSession:   true,
		ConfigGate:    GateValidation,
		IntentGate:    GateAllowAll,
		GenerateCond: func(ctx GenContext) bool {
			return hasUISurface(ctx.SurfaceTypes)
		},
		DependsOn: []DepRef{
			{TaskTemplate: "T-test-run", Resolver: ResolveLastRunTestOrBusiness},
		},
	},
	// 11. Consolidate specs (breakdown mode)
	{
		ID:            "T-specs-consolidate",
		Key:           "consolidate-specs",
		Title:         "Consolidate Specs",
		Type:          TypeDocConsolidate,
		Priority:      string(types.PriorityP2),
		EstimatedTime: "20min",
		ConfigGate:    GateConsolidateSpecs,
		IntentGate:    GateAllowAll,
		DependsOn: []DepRef{
			{TaskTemplate: "T-test-run", Resolver: ResolveLastRunTestOrBusiness},
		},
	},
	// 12. Detect spec drift (quick mode)
	{
		ID:            "T-quick-doc-drift",
		Key:           "quick-drift-detection",
		Title:         "Detect Spec Drift",
		Type:          TypeDocDrift,
		Priority:      string(types.PriorityP2),
		EstimatedTime: "15min",
		ConfigGate:    GateConsolidateSpecs,
		IntentGate:    GateAllowAll,
		DependsOn: []DepRef{
			{TaskTemplate: "T-test-run", Resolver: ResolveLastRunTestOrBusiness},
		},
	},
}
