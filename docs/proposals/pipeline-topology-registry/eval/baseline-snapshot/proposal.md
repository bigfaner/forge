---
name: pipeline-topology-registry
status: draft
created: 2026-05-29
---

# Pipeline Topology Registry

## Problem

Forge CLI's auto-task generation pipeline is fragile. Adding or removing a task type causes cascading failures across 5+ code locations because task IDs and dependency relationships are hardcoded as string literals scattered across multiple functions.

**Evidence** (recent bugs):

| Bug | Root Cause | Symptom |
|-----|-----------|---------|
| T-quick-doc-drift appearing unexpectedly | `auto.consolidateSpecs.quick` defaults to `true`; no visibility into what gets auto-generated | Unexpected task |
| findTaskIndexOrPanic crashes | Adding a config-gated task type without updating all lookup sites | Runtime panic |
| Drift task runs before business tasks | `ResolveDriftFallbackDep` hardcoded drift IDs; `resolveQuickDeps` only sets deps when test pipeline is enabled | Broken dependency chain |
| Missing drift dep when test pipeline disabled | Step 7.6 fallback exists but is a separate code path from step 7.5 | Missing task dependency |
| isTestTaskID misses new prefixes | Hardcoded 6-prefix list doesn't cover new task categories | Orphaned index entries |

**Code smell**: Adding a new auto-generated task type currently requires changes in:
1. `types.go` — add type constant
2. `autogen.go` — add `AutoGenTaskDef` in `GetBreakdownTestTasks()` AND `GetQuickTestTasks()`
3. `autogen.go` — add dependency wiring in `resolveBreakdownDeps()` AND `resolveQuickDeps()`
4. `autogen.go` — add fallback dep in `ResolveDriftFallbackDep()`
5. `infer.go` — add case to `InferType()` switch
6. `build.go` — possibly update `isTestTaskID()`, `isAutoGenForDep()`, `IsAutoGenTaskID()`

**Cost of inaction**: Each new task type has a ~40% chance of introducing at least one pipeline bug, based on the last 5 additions (3 caused issues). Developer velocity on pipeline changes is constrained by the need to manually verify all 6+ touch points.

## Solution

Replace the scattered procedural pipeline with a **Pipeline Topology Registry** — a single declarative data structure that defines all auto-generated tasks, their types, IDs, dependencies, config gating, and expansion rules. All consuming code (autogen, infer, build) derives from this registry.

**User-facing behavior**: No change. The pipeline generates the same tasks with the same dependencies. The refactoring is purely internal architecture.

### Design

#### Core Data Structure

```go
// ConfigGateFunc returns true when the auto config enables this node for the given mode.
// mode is "quick" or "breakdown".
type ConfigGateFunc func(mode string, auto forgeconfig.AutoConfig) bool

// GenerateCondFunc returns true when the business task composition permits this node.
type GenerateCondFunc func(tasks []Task) bool

// DepResolveFunc dynamically resolves dependency IDs at generation time.
// Returns concrete task IDs. Returns nil when the reference cannot be resolved.
type DepResolveFunc func(ctx *GenContext) []string

// GenContext carries state accumulated during pipeline generation.
// Populated progressively as nodes are processed in declaration order.
type GenContext struct {
    Mode           string
    Surfaces       map[string]string
    ExecutionOrder []string
    Auto           forgeconfig.AutoConfig
    BusinessTasks  []Task
    // Filled during generation as nodes are expanded:
    UpstreamIDs  []string // IDs of the immediately preceding generated node(s)
    RunTestChain []string // IDs of expanded run-test tasks in serial order
    AllGenerated []string // IDs of all nodes generated so far (in order)
}

// PipelineNode defines a single node in the auto-generated task pipeline.
type PipelineNode struct {
    // Type is the task type constant (e.g., TypeTestGenJourneys).
    Type string
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
    // Mode restricts this node to a specific mode. Empty means both modes.
    // "quick" = quick mode only, "breakdown" = breakdown mode only.
    Mode string
    // GenerateCondition returns true when the business task composition permits this node.
    // nil defaults to CondHasTestableTasks.
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
```

#### Predefined Gate, Condition & Resolver Functions

```go
// --- Config Gate Functions ---
// Each gate reads the corresponding auto.{category}.{mode-field} from config.

func GateTest(mode string, auto forgeconfig.AutoConfig) bool {
    if mode == "quick" { return auto.Test.Quick }
    return auto.Test.Full
}

func GateValidation(mode string, auto forgeconfig.AutoConfig) bool {
    if mode == "quick" { return auto.Validation.Quick }
    return auto.Validation.Full
}

func GateConsolidateSpecs(mode string, auto forgeconfig.AutoConfig) bool {
    if mode == "quick" { return auto.ConsolidateSpecs.Quick }
    return auto.ConsolidateSpecs.Full
}

func GateCleanCode(mode string, auto forgeconfig.AutoConfig) bool {
    if mode == "quick" { return auto.CleanCode.Quick }
    return auto.CleanCode.Full
}

// --- Generate Condition Functions ---

// CondHasTestableTasks returns true when any business task has a testable type.
func CondHasTestableTasks(tasks []Task) bool {
    for _, t := range tasks {
        if IsTestableType(t.Type) { return true }
    }
    return false
}

// CondHasDocTasks returns true when any business task has a doc-category type.
func CondHasDocTasks(tasks []Task) bool {
    for _, t := range tasks {
        if CategoryForType(t.Type) == CategoryDoc { return true }
    }
    return false
}

// CondAlways returns true unconditionally.
func CondAlways(tasks []Task) bool { return true }

// --- Dependency Resolver Functions ---

// ResolveLastRunTest returns the ID of the last task in the run-test expansion chain.
// Returns nil when no run-test tasks have been generated.
var ResolveLastRunTest DepResolveFunc = func(ctx *GenContext) []string {
    if len(ctx.RunTestChain) == 0 { return nil }
    return []string{ctx.RunTestChain[len(ctx.RunTestChain)-1]}
}

// ResolveUpstream returns the IDs of the immediately preceding generated node(s).
// For single nodes: one ID. For expanded nodes: all expanded IDs of the previous node.
var ResolveUpstream DepResolveFunc = func(ctx *GenContext) []string {
    if len(ctx.UpstreamIDs) == 0 { return nil }
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
```

#### Pipeline Registry

```go
// PipelineRegistry is the single source of truth for the auto-generated task pipeline.
// Order matters: nodes are processed in declaration order.
var PipelineRegistry = []PipelineNode{
    // --- Doc Review (generated whenever business tasks include doc-category types) ---
    {
        Type: TypeDocReview, ID: "T-review-doc",
        Title: "Review Documentation Quality", Priority: "P1", EstimatedTime: "30min",
        ConfigGate: nil, // no config gate
        GenerateCondition: CondHasDocTasks,
        DependsOn: []DepRef{{Resolve: ResolveDocTasks}},
    },
    // --- Test Generation ---
    {
        Type: TypeTestGenJourneys, ID: "T-test-gen-journeys",
        Title: "Generate Test Journeys", Priority: "P1", EstimatedTime: "20-30min",
        ConfigGate: GateTest, StrategyKind: "interface",
    },
    // --- Eval (breakdown only) ---
    {
        Type: TypeEvalJourney, ID: "T-eval-journey",
        Title: "Evaluate Journey Quality", Priority: "P1", EstimatedTime: "20-30min",
        ConfigGate: GateTest, Mode: "breakdown", MainSession: true,
        DependsOn: []DepRef{{Ref: "T-test-gen-journeys"}},
    },
    {
        Type: TypeTestGenContracts, ID: "T-test-gen-contracts",
        Title: "Generate Test Contracts", Priority: "P1", EstimatedTime: "30-45min",
        ConfigGate: GateTest, Mode: "breakdown",
        DependsOn: []DepRef{{Ref: "T-eval-journey"}},
    },
    {
        Type: TypeEvalContract, ID: "T-eval-contract",
        Title: "Evaluate Contract Quality", Priority: "P1", EstimatedTime: "20-30min",
        ConfigGate: GateTest, Mode: "breakdown", MainSession: true,
        DependsOn: []DepRef{{Ref: "T-test-gen-contracts"}},
    },
    // --- Gen Scripts (per surface type) ---
    {
        Type: TypeTestGenScripts, ID: "T-test-gen-scripts-{surface-type}",
        Title: "Generate {test-type-title} Scripts", Priority: "P1", EstimatedTime: "1-2h",
        ConfigGate: GateTest, Mode: "breakdown", Expansion: "per-surface-type",
        DependsOn: []DepRef{{Ref: "T-eval-contract"}},
        StrategyKind: "generate",
    },
    // --- Run Tests (per surface key) ---
    {
        Type: TypeTestRun, ID: "T-test-run-{surface-key}",
        Title: "Run {test-type-title}", Priority: "P1", EstimatedTime: "30min-1h",
        ConfigGate: GateTest,
        DependsOn: []DepRef{{Resolve: ResolveUpstream}},
        Expansion: "per-surface-key", StrategyKind: "run",
    },
    // --- Validation ---
    {
        Type: TypeValidationCode, ID: "T-validate-code",
        Title: "Validate Code Quality", Priority: "P2", EstimatedTime: "15min",
        ConfigGate: GateValidation,
        DependsOn: []DepRef{{Resolve: ResolveLastRunTest}},
    },
    {
        Type: TypeValidationUx, ID: "T-validate-ux",
        Title: "Validate User Experience", Priority: "P2", EstimatedTime: "15min",
        ConfigGate: GateValidation,
        DependsOn: []DepRef{{Resolve: ResolveLastRunTest}},
        UISurfaceOnly: true, MainSession: true,
    },
    // --- Spec Consolidation/Drift ---
    {
        Type: TypeDocConsolidate, ID: "T-specs-consolidate",
        Title: "Consolidate Specs", Priority: "P2", EstimatedTime: "20min",
        ConfigGate: GateConsolidateSpecs, Mode: "breakdown",
        DependsOn: []DepRef{{Resolve: ResolveLastRunTest}},
    },
    {
        Type: TypeDocDrift, ID: "T-quick-doc-drift",
        Title: "Detect Spec Drift", Priority: "P2", EstimatedTime: "15min",
        ConfigGate: GateConsolidateSpecs, Mode: "quick",
        DependsOn: []DepRef{{Resolve: ResolveLastRunTest}},
    },
    // --- Clean Code ---
    {
        Type: TypeCleanCode, ID: "T-clean-code",
        Title: "Simplify and Clean Code", Priority: "P2", EstimatedTime: "20min",
        ConfigGate: GateCleanCode,
        DependsOn: []DepRef{}, // resolved by caller: depends on last business task
    },
}
```

**Gating rules (both must pass for a node to be generated)**:

1. **ConfigGate**: `nil` = always passes. Otherwise calls the gate function with `(mode, autoConfig)`.
2. **GenerateCondition**: `nil` = defaults to `CondHasTestableTasks`. Otherwise calls the condition function with `(businessTasks)`.

**T-review-doc 的三种场景**：

| 场景 | CondHasDocTasks | CondHasTestableTasks | T-review-doc 行为 |
|------|:-:|:-:|---|
| 纯 doc 任务 | true | false | 生成。依赖 @doc-tasks，无 test pipeline |
| coding + doc 混合 | true | true | 生成。依赖 @doc-tasks，同时被注入为 test pipeline 首任务的前置依赖 |
| 纯 coding | false | true | 不生成（CondHasDocTasks 返回 false） |

#### Derived Functions

All current scattered logic derives from the registry:

**InferType**: Instead of a 15-case switch, iterate the registry and match ID patterns (with wildcard support for `{surface-key}`/`{surface-type}` placeholders). Fall back to prefix-based matching for runtime-only tasks (fix-*, disc-*).

**isTestTaskID / IsAutoGenTaskID**: Derive from registry by collecting all IDs (after expansion) into a lookup set.

**isAutoGenForDep**: Derive from registry, same as IsAutoGenTaskID (all registry nodes are auto-gen for dependency purposes).

**Task generation**: `GenerateTestTasks(mode, surfaces, executionOrder, auto, businessTasks)` filters the registry by:
1. `Mode` — node's mode must match or be empty ("both")
2. `ConfigGate` — `nil` passes; otherwise calls `ConfigGate(mode, auto)` must return true
3. `GenerateCondition` — `nil` defaults to `CondHasTestableTasks`; otherwise calls the function with `businessTasks`
4. `UISurfaceOnly` — skip when no surface has a visual UI

Then expands per-surface nodes and resolves dependency references.

**Dependency resolution**: Walk the registry in order. For each DepRef:
- If `Resolve` is set: call `Resolve(ctx)` to get concrete IDs dynamically
- If `Resolve` is nil: use `Ref` as a static concrete ID
- `GenContext` is progressively populated as nodes are generated: `UpstreamIDs`, `RunTestChain`, `AllGenerated`

#### Runtime Task Coordination

The registry defines the **static** pipeline topology. Runtime tasks (fix-*, disc-*) are generated by quality-gate and run-tasks dispatcher using prefix-based matching, outside the registry. The coordination contract:

1. Runtime tasks use IDs that do NOT match any registry pattern (prefix `fix-` or `disc-`)
2. `InferType()` tries registry first, then falls back to prefix rules for runtime tasks
3. `IsAutoGenTaskID()` returns false for runtime tasks (they are business tasks)
4. Runtime tasks have their own dependency wiring (source task ID)

This ensures the registry is a **closed, immutable** data structure at runtime, while runtime task creation remains flexible.

#### Elimination of findTaskIndexOrPanic

Replace all `findTaskIndexOrPanic` calls with registry-based lookups that:
1. Return an error (not panic) when a referenced ID is not found
2. Are validated at init time: the registry is checked for referential integrity (all `DependsOn` references resolve to existing nodes or are valid special references)

## Alternatives Considered

### Alternative 1: Table-Driven with Preserved Function Structure
Keep existing function structure, replace hardcoded strings with constant arrays. Lower risk but incomplete: dependency resolution remains procedural, and adding a task type still requires touching multiple arrays.

**Why rejected**: Doesn't solve the root cause (scattered definitions). Just makes the strings slightly more discoverable.

### Alternative 2: Do Nothing, Strengthen Tests
Add comprehensive integration tests covering all task type combinations. Future bugs are caught by tests.

**Why rejected**: Tests verify behavior but don't prevent the fundamental issue — adding a task type still requires 6+ file changes. Developer friction remains high.

### Alternative 3: Code Generation
Define pipeline in YAML, generate Go code via `go generate`.

**Why rejected**: Adds a build step and indirection. The registry is Go-native and provides the same compile-time safety without code generation complexity.

## Scope

### In Scope

1. Define `PipelineNode`, `DepRef`, and `PipelineRegistry` in `pkg/task/pipeline.go`
2. Refactor `autogen.go` — `GetBreakdownTestTasks()`, `GetQuickTestTasks()`, and dependency resolution functions to derive from registry
3. Refactor `infer.go` — `InferType()` to match against registry ID patterns, with prefix fallback for runtime tasks
4. Refactor `build.go` steps 7/7.5/7.6 — task generation and injection to use registry-driven generation
5. Refactor `isTestTaskID()`, `isAutoGenForDep()`, `IsAutoGenTaskID()` to derive from registry
6. Eliminate `findTaskIndexOrPanic` — replace with safe lookups and init-time validation
7. Add init-time referential integrity check for the registry
8. Update all existing tests (`autogen_test.go`, `infer_test.go`, `build_test.go`, `claim_test.go`)
9. Add new tests for registry validation and derived function correctness

### Out of Scope

1. Quality-gate runtime fix-task generation (separate concern, uses prefix rules)
2. Run-tasks dispatcher runtime fix-task generation (separate concern)
3. Task state machine (`statemachine.go`)
4. Prompt templates (`pkg/prompt/templates/`)
5. Skill files (`plugins/forge/skills/`)
6. Stage gate generation (`GenerateStageGates`) — remains separate as it derives from business tasks, not pipeline topology

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|-----------|
| Regression in task generation output | Medium | High | Snapshot tests: capture generated task lists for all mode/config combinations before and after refactor |
| Dependency resolution edge case missed | Medium | High | Port all existing gotcha lessons as test cases; init-time validation catches broken refs |
| Performance regression (registry iteration vs hardcoded switch) | Low | Low | Registry has ~12 entries; iteration cost is negligible. Benchmark if concerned |
| Expanding per-surface nodes creates ID conflicts | Low | Medium | Init-time validation: check expanded IDs are unique |
| T-review-doc injection into test pipeline missed | Low | Medium | Mixed feature test case: verify T-review-doc appears as test pipeline first-task dependency when both doc and coding tasks exist |

## Success Criteria

1. **Single touch point**: Adding a new auto-generated task type requires changes in exactly ONE file (`pipeline.go`) — the registry definition plus optionally the type constant in `types.go`
2. **Zero panics**: No `findTaskIndexOrPanic` calls remain; all lookups return errors or are validated at init time
3. **Init-time validation**: Invalid dependency references in the registry fail at CLI startup, not at runtime during `forge task index`
4. **Behavioral parity**: `forge task index` produces identical output for all existing feature configurations after the refactoring
5. **InferType coverage**: All auto-generated task IDs are correctly typed by the registry-derived `InferType`, including surface-expanded variants (e.g., `T-test-gen-scripts-api`) and T-review-doc
6. **GenerateCondition correctness**: T-review-doc is generated in all three scenarios (docs-only, mixed coding+doc, not for coding-only)
7. **Zero magic values**: All string literals in `pipeline.go` and refactored code must use typed enum constants. This includes:
   - Mode values (`ModeQuick`, `ModeBreakdown`)
   - Expansion values (`ExpansionNone`, `ExpansionPerSurfaceKey`, `ExpansionPerSurfaceType`)
   - Task ID templates (`IDGenJourneys`, `IDRunTest` etc.)
   - Any other domain-specific strings
   Raw string literals are only allowed in constant definitions and test fixtures, never in logic code.
8. **Test coverage >= 80%**: For the new `pipeline.go` and refactored functions
9. **All existing tests pass**: Zero regressions in `autogen_test.go`, `infer_test.go`, `build_test.go`, `claim_test.go`
