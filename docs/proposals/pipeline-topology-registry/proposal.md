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

**Cost of inaction**: Of the last 5 task type additions, 3 caused at least one pipeline bug. Developer velocity on pipeline changes is constrained by the need to manually verify all 6+ touch points.

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

// ResolveLastBusinessTask returns the ID of the highest-numbered business task.
// Used by T-clean-code which must run after all business tasks complete.
var ResolveLastBusinessTask DepResolveFunc = func(ctx *GenContext) []string {
    if len(ctx.BusinessTasks) == 0 { return nil }
    return []string{ctx.BusinessTasks[len(ctx.BusinessTasks)-1].ID}
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
        ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks, StrategyKind: "interface",
    },
    // --- Eval (breakdown only) ---
    {
        Type: TypeEvalJourney, ID: "T-eval-journey",
        Title: "Evaluate Journey Quality", Priority: "P1", EstimatedTime: "20-30min",
        ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks, Mode: "breakdown", MainSession: true,
        DependsOn: []DepRef{{Ref: "T-test-gen-journeys"}},
    },
    {
        Type: TypeTestGenContracts, ID: "T-test-gen-contracts",
        Title: "Generate Test Contracts", Priority: "P1", EstimatedTime: "30-45min",
        ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks, Mode: "breakdown",
        DependsOn: []DepRef{{Ref: "T-eval-journey"}},
    },
    {
        Type: TypeEvalContract, ID: "T-eval-contract",
        Title: "Evaluate Contract Quality", Priority: "P1", EstimatedTime: "20-30min",
        ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks, Mode: "breakdown", MainSession: true,
        DependsOn: []DepRef{{Ref: "T-test-gen-contracts"}},
    },
    // --- Gen Scripts (per surface type) ---
    {
        Type: TypeTestGenScripts, ID: "T-test-gen-scripts-{surface-type}",
        Title: "Generate {test-type-title} Scripts", Priority: "P1", EstimatedTime: "1-2h",
        ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks, Mode: "breakdown", Expansion: "per-surface-type",
        DependsOn: []DepRef{{Ref: "T-eval-contract"}},
        StrategyKind: "generate",
    },
    // --- Run Tests (per surface key) ---
    {
        Type: TypeTestRun, ID: "T-test-run-{surface-key}",
        Title: "Run {test-type-title}", Priority: "P1", EstimatedTime: "30min-1h",
        ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks,
        DependsOn: []DepRef{{Resolve: ResolveUpstream}},
        Expansion: "per-surface-key", StrategyKind: "run",
    },
    // --- Validation ---
    {
        Type: TypeValidationCode, ID: "T-validate-code",
        Title: "Validate Code Quality", Priority: "P2", EstimatedTime: "15min",
        ConfigGate: GateValidation, GenerateCondition: CondAlways,
        DependsOn: []DepRef{{Resolve: ResolveLastRunTestOrBusiness}},
    },
    {
        Type: TypeValidationUx, ID: "T-validate-ux",
        Title: "Validate User Experience", Priority: "P2", EstimatedTime: "15min",
        ConfigGate: GateValidation, GenerateCondition: CondAlways,
        DependsOn: []DepRef{{Resolve: ResolveLastRunTestOrBusiness}},
        UISurfaceOnly: true, MainSession: true,
    },
    // --- Spec Consolidation/Drift ---
    {
        Type: TypeDocConsolidate, ID: "T-specs-consolidate",
        Title: "Consolidate Specs", Priority: "P2", EstimatedTime: "20min",
        ConfigGate: GateConsolidateSpecs, GenerateCondition: CondAlways, Mode: "breakdown",
        DependsOn: []DepRef{{Resolve: ResolveLastRunTestOrBusiness}},
    },
    {
        Type: TypeDocDrift, ID: "T-quick-doc-drift",
        Title: "Detect Spec Drift", Priority: "P2", EstimatedTime: "15min",
        ConfigGate: GateConsolidateSpecs, GenerateCondition: CondAlways, Mode: "quick",
        DependsOn: []DepRef{{Resolve: ResolveLastRunTestOrBusiness}},
    },
    // --- Clean Code ---
    {
        Type: TypeCleanCode, ID: "T-clean-code",
        Title: "Simplify and Clean Code", Priority: "P2", EstimatedTime: "20min",
        ConfigGate: GateCleanCode, GenerateCondition: CondAlways,
        DependsOn: []DepRef{{Resolve: ResolveLastBusinessTask}},
    },
}
```

**Gating rules (both must pass for a node to be generated)**:

1. **ConfigGate**: `nil` = always passes. Otherwise calls the gate function with `(mode, autoConfig)`.
2. **GenerateCondition**: Every node MUST set `GenerateCondition` explicitly. No nil default. Use `CondAlways` for unconditional generation or `CondHasTestableTasks` / `CondHasDocTasks` for conditional. Init-time validation rejects nodes with nil `GenerateCondition`.

**T-review-doc 的三种场景**：

| 场景 | CondHasDocTasks | CondHasTestableTasks | T-review-doc 行为 |
|------|:-:|:-:|---|
| 纯 doc 任务 | true | false | 生成。依赖 @doc-tasks，无 test pipeline |
| coding + doc 混合 | true | true | 生成。依赖 @doc-tasks，同时被注入为 test pipeline 首任务的前置依赖 |
| 纯 coding | false | true | 不生成（CondHasDocTasks 返回 false） |

**Test pipeline disabled / all ConfigGates off 场景**：当 `GateTest` 返回 false 时，所有 test 节点跳过。非 test 节点（validate/consolidate/clean-code）通过 `ResolveLastRunTestOrBusiness` 自动退化为依赖 lastBusinessTask。当所有 ConfigGates 均为 off 时，仅 T-review-doc（若有 doc 任务）可能生成。

**Empty business tasks scenario**：When `businessTasks` is empty, `CondHasTestableTasks`/`CondHasDocTasks` return false — only `CondAlways` nodes may generate. Their resolvers return nil (no RunTestChain, no BusinessTasks), so nil-handling rule applies: nodes generate with empty `DependsOn` as pipeline roots. This scenario occurs when `forge task index` is invoked with only auto-generated tasks (e.g., running clean-code + validation without user-defined business tasks). Expected result:

| Node | ConfigGate | GenerateCondition | Resolver result | Generated? |
|------|-----------|-------------------|-----------------|------------|
| T-review-doc | nil | CondHasDocTasks → false | — | No |
| T-test-gen-journeys | GateTest | CondHasTestableTasks → false | — | No |
| T-validate-code | GateValidation | CondAlways → true | ResolveLastRunTestOrBusiness → nil | Yes, empty DependsOn |
| T-clean-code | GateCleanCode | CondAlways → true | ResolveLastBusinessTask → nil | Yes, empty DependsOn |

**T-review-doc 反向注入机制**：当 doc 和 coding 任务同时存在时，T-review-doc 需要被注入为 test pipeline 首任务的前置依赖。这在 `GenerateTestTasks` 后处理阶段实现：

```go
// Post-generation injection: when T-review-doc exists AND test pipeline is active,
// prepend T-review-doc as a dependency of the first test task.
func injectReviewDocDep(generated []Task, businessTasks []Task) {
    reviewDocID := "T-review-doc"
    hasReviewDoc := false
    for _, t := range generated {
        if t.ID == reviewDocID { hasReviewDoc = true; break }
    }
    if !hasReviewDoc { return }
    for i, t := range generated {
        if IsTestPipelineTask(t) {
            generated[i].DependsOn = append([]string{reviewDocID}, generated[i].DependsOn...)
            break
        }
    }
}
```

此逻辑不属于 registry 声明式定义，而是生成后协调步骤，因为它需要跨 pipeline 阶段（doc → test）的可见性。

**Stage-gate 与 clean-code 注入处理**：`ResolveFirstTestDep` 执行三项逻辑：(a) stage-gate 依赖解析 — 独立于 pipeline topology，由 `GenerateStageGates` 确保；(b) T-clean-code 中间依赖注入 — 通过 post-generation step 实现；(c) T-review-doc 注入 — 由 `injectReviewDocDep` 覆盖。

```go
func injectCleanCodeDep(generated []Task) {
    cleanCodeID := "T-clean-code"
    for _, t := range generated {
        if t.ID == cleanCodeID {
            for i, t := range generated {
                if IsTestPipelineTask(t) {
                    generated[i].DependsOn = append([]string{cleanCodeID}, generated[i].DependsOn...)
                    return
                }
            }
        }
    }
}
```

后处理阶段按序执行：先 `injectCleanCodeDep`，再 `injectReviewDocDep`。执行顺序不影响正确性——两者各自 prepend 到同一 first-test-task 的 DependsOn，无论先后均产生相同集合。两个 step 共同产生**并行扇入**（fan-in）依赖：`first-test-task.DependsOn = [T-review-doc, T-clean-code, ...upstream deps...]`。其中 T-review-doc 自身依赖 `[doc business tasks]`，T-clean-code 自身依赖 `[last business task]`，二者之间无串行关系——它们可以并行执行，first-test-task 等待两者全部完成后才开始。

**依赖关系示意图**：

```
doc-tasks ──────→ T-review-doc ──┐
                                  ├─→ first-test-task ─→ ... test pipeline
business-tasks ─→ T-clean-code ──┘
```

注意：这不是串行链（T-review-doc 不依赖 T-clean-code，反之亦然），而是 first-test-task 对两个独立分支的并行等待。

**Escape-hatch protocol**：上述 post-generation 注入步骤是 registry 声明式定义的 escape hatch，用于处理跨 pipeline 阶段的依赖。为防止 gradual re-proceduralization：(1) 仅跨阶段可见性依赖允许 escape hatch，纯拓扑依赖必须在 registry 声明；(2) escape hatch 总数上限 5 个，达上限时必须扩展 registry 表达力；(3) 每个 step 需独立单元测试 + snapshot 覆盖；(4) 每个 escape hatch 需记录存在原因及未来消除路径。

#### Derived Functions

All current scattered logic derives from the registry:

**InferType**: Instead of a 15-case switch, iterate the registry and match ID patterns (with wildcard support for `{surface-key}`/`{surface-type}` placeholders). Fall back to prefix-based matching for runtime-only tasks (fix-*, doc-fix-*, disc-*).

InferType 保持两阶段签名以获取 surfaces map：先尝试 registry 模式匹配（纯 ID 字符串匹配，无需 surfaces），当展开后的 ID 包含 surface key 时，用 `surfaces` map 验证后缀是否为有效 surface key。当前签名为 `InferType(taskID string, surfaces map[string]string) string`，保持不变。

**isTestTaskID / IsAutoGenTaskID**: Derive from registry by collecting all IDs (after expansion) into a lookup set.

**isAutoGenForDep**: Derive from registry, same as IsAutoGenTaskID (all registry nodes are auto-gen for dependency purposes).

**Task generation**: `GenerateTestTasks(mode, surfaces, executionOrder, auto, businessTasks)` filters the registry by:
1. `Mode` — node's mode must match or be empty ("both")
2. `ConfigGate` — `nil` passes; otherwise calls `ConfigGate(mode, auto)` must return true
3. `GenerateCondition` — calls the function with `businessTasks`; init-time validation guarantees non-nil
4. `UISurfaceOnly` — skip when no surface has a visual UI

**Breaking change: `businessTasks` parameter**: The current signature is `GenerateTestTasks(mode string, surfaces map[string]string, executionOrder []string, auto forgeconfig.AutoConfig)`. The proposal adds a `businessTasks []Task` parameter. Both callers (`build.go` step 7 and step 7.5) already have `businessTasks` in scope — no additional data fetching required.

Then expands per-surface nodes and resolves dependency references:

```
func GenerateTestTasks(mode, surfaces, executionOrder, auto, businessTasks) []Task {
    ctx := &GenContext{
        Mode: mode, Surfaces: surfaces, ExecutionOrder: executionOrder,
        Auto: auto, BusinessTasks: businessTasks,
    }
    var generated []Task

    for _, node := range PipelineRegistry {
        // Step 1: Filter
        if node.Mode != "" && node.Mode != mode { continue }
        if node.ConfigGate != nil && !node.ConfigGate(mode, auto) { continue }
        if !node.GenerateCondition(businessTasks) { continue }
        if node.UISurfaceOnly && !hasVisualUI(surfaces) { continue }

        // Step 2: Expand — produce concrete task(s) from template
        expanded := expand(node, surfaces)
        // expanded is []Task with IDs/Keys substituted:
        //   ""                 → [node.ID]
        //   "per-surface-key"  → [node.ID per key in surfaces]
        //   "per-surface-type" → [node.ID per unique type in surfaces]

        // Step 3: Resolve dependencies for each expanded task
        for _, t := range expanded {
            for _, dep := range node.DependsOn {
                if dep.Resolve != nil {
                    ids := dep.Resolve(ctx)
                    if ids == nil { continue } // skip this dep entry
                    t.DependsOn = append(t.DependsOn, ids...)
                } else {
                    t.DependsOn = append(t.DependsOn, dep.Ref)
                }
            }
        }

        // Step 4: Update GenContext (progressive population)
        ids := taskIDs(expanded)
        ctx.AllGenerated = append(ctx.AllGenerated, ids...)
        ctx.UpstreamIDs = ids
        if node.Type == TypeTestRun {
            ctx.RunTestChain = append(ctx.RunTestChain, ids...)
        }
        generated = append(generated, expanded...)
    }

    // Step 5: Post-generation injection (escape hatches)
    injectCleanCodeDep(generated)
    injectReviewDocDep(generated, businessTasks)
    return generated
}
```


**Expansion rules**：

1. `per-surface-key`：对 `surfaces` map 中的每个 key 生成一个任务。当 `isSingleSurface(surfaces)` 为 true 时退化——任务 ID 去掉 `-{surface-key}` 后缀（例如 `T-test-run` 而非 `T-test-run-{surface-key}`），但仍使用该唯一 surface 的配置。
2. `per-surface-type`：对 surfaces 中的每个唯一 type 生成一个任务。无特殊退化行为。

**Dependency resolution**: Walk the registry in order. For each DepRef:
- If `Resolve` is set: call `Resolve(ctx)` to get concrete IDs dynamically
- If `Resolve` is nil: use `Ref` as a static concrete ID
- `GenContext` is progressively populated as nodes are generated: `UpstreamIDs`, `RunTestChain`, `AllGenerated`

#### Runtime Task Coordination

The registry defines the **static** pipeline topology. Runtime tasks (fix-*, doc-fix-*, disc-*) are generated by quality-gate and run-tasks dispatcher using prefix-based matching, outside the registry. The coordination contract:

1. Runtime tasks use IDs that do NOT match any registry pattern (prefix `fix-`, `doc-fix-`, or `disc-`)
2. `InferType()` tries registry first, then falls back to prefix rules for runtime tasks (fix-*, doc-fix-*, disc-*)
3. `IsAutoGenTaskID()` returns false for runtime tasks (they are business tasks)
4. Runtime tasks have their own dependency wiring (source task ID)

This ensures the registry is a **closed, immutable** data structure at runtime, while runtime task creation remains flexible.

#### Elimination of findTaskIndexOrPanic

Replace all `findTaskIndexOrPanic` calls with registry-based lookups that:
1. Return an error (not panic) when a referenced ID is not found
2. Are validated via two-phase validation: static `DependsOn.Ref` strings checked at init-time; dynamic resolver results checked at runtime

## Alternatives Considered

### Alternative 1: Table-Driven with Preserved Function Structure
Keep existing function structure, replace hardcoded strings with constant arrays. Lower risk but incomplete: dependency resolution remains procedural, and adding a task type still requires touching multiple arrays.

**Trade-off**: Table-driven is a legitimate incremental step — it solves ~80% of "forgot to update a lookup site" bugs by centralizing string definitions into shared arrays, with minimal risk and no algorithmic changes. However, dependency wiring in `resolveBreakdownDeps()`/`resolveQuickDeps()` remains procedural and must be manually kept in sync. Each new task type still requires 3-4 function updates (down from 6). The full registry is worth the additional complexity because it eliminates the "forgot to wire a dependency" class of bugs that table-driven cannot prevent.

### Alternative 2: Dependency Injection Container (Industry Pattern)
Use a DI container (e.g., [google/wire](https://github.com/google/wire)) to declare task types as providers and deps as injection points. Wire provides compile-time dep graph validation.

**Trade-off**: Tekton Pipelines 和 GitHub Actions 采用**全静态 DAG 调度**模式：Tekton 用 `runAfter` + `params` 在 PipelineRun 创建时确定完整拓扑，GitHub Actions 用 `needs:` + `${{ needs.*.outputs }}` 在 workflow 解析时确定 job 依赖。两者的关键特征是：(1) 拓扑在调度前完全确定，不支持运行时条件性节点生成；(2) 依赖目标在声明时可知——`needs: [build]` 中的 `build` 必须存在于同一 workflow。

我们的场景本质不同：`ConfigGate` 和 `GenerateCondition` 在运行时动态决定哪些节点存在，导致 DAG 形状在调度前不可知。例如 `GateTest` 为 false 时整个 test pipeline 消失，`CondHasDocTasks` 为 false 时 T-review-doc 不生成。这意味着静态 DAG 调度不适用——`DepResolveFunc` 需要在生成时从渐进填充的 `GenContext` 中解析依赖目标。

具体对比：

| 维度 | Tekton/GitHub Actions | Pipeline Topology Registry |
|------|----------------------|---------------------------|
| DAG 确定 | 调度前完全静态 | 生成时动态（ConfigGate/GenerateCondition 裁剪） |
| 依赖解析 | 声明时已知目标 ID | 运行时从 GenContext 解析 |
| 条件执行 | 通过 `when`/`if` 过滤，但节点仍存在于 DAG 中 | 不满足条件的节点完全不生成 |
| 扩展模式 | 无内建 per-X 展开 | per-surface-key/type 模板展开 |

若采用 Tekton 模式，需要引入 "条件节点消除" 阶段 + "模板展开" 阶段作为 DAG 前处理——这实质上是在 DAG 之上重建我们的 GenerateTestTasks 算法，增加抽象层但无实际收益。DI wraps tasks as providers — boilerplate without domain alignment。

### Alternative 3: Do Nothing, Strengthen Tests
Add comprehensive integration tests covering all task type combinations. Future bugs are caught by tests.

**Why rejected**: Tests verify behavior but don't prevent the fundamental issue — adding a task type still requires 6+ file changes. Developer friction remains high.

### Alternative 4: Code Generation
Define pipeline in YAML, generate Go code via `go generate`.

**Why rejected**: `go generate` is standard practice, but here it introduces a YAML-to-Go translation layer that developers must mentally trace during debugging, with no runtime benefit over the native Go registry.

## Scope

### In Scope

1. Define `PipelineNode`, `DepRef`, and `PipelineRegistry` in `pkg/task/pipeline.go`
2. Refactor `autogen.go` — `GetBreakdownTestTasks()`, `GetQuickTestTasks()`, and dependency resolution functions to derive from registry
3. Refactor `infer.go` — `InferType()` to match against registry ID patterns, with prefix fallback for runtime tasks
4. Refactor `build.go` steps 7/7.5/7.6 — task generation and injection to use registry-driven generation
5. Refactor `isTestTaskID()`, `isAutoGenForDep()`, `IsAutoGenTaskID()` to derive from registry
6. Eliminate `findTaskIndexOrPanic` — replace with safe lookups and init-time validation
7. Add two-phase validation for the registry:
   - **Phase 1 (static, init-time)**: Run at CLI startup via `init()`. Validates: all `DependsOn.Ref` strings reference existing node IDs or match templates (placeholders treated as wildcards — a `Ref` of `"T-test-gen-scripts-{surface-type}"` matches the node with `ID: "T-test-gen-scripts-{surface-type}"` by normalizing placeholders before comparison; currently no `DepRef` uses template placeholders in `Ref`, but validation handles them for forward compatibility); all expanded IDs are unique; `GenerateCondition` is non-nil; `Key`/`ID` template placeholders match `Expansion` setting; escape-hatch count <= 5; ordering invariants: `ResolveUpstream` users must appear after their expected upstream nodes. Panics on failure. This replaces `ValidateAutogenTemplates`.
   - **Phase 2 (dynamic, runtime)**: Run at the start of `GenerateTestTasks`. Validates: all resolver-returned IDs exist in generated task set; no circular dependencies. Returns errors (does not panic).
8. Update all existing tests (`autogen_test.go`, `infer_test.go`, `build_test.go`, `claim_test.go`)
9. Add new tests for registry validation and derived function correctness

**Effort estimate**: 2-3 development days. Breakdown: items 1 + 7 (registry definition + validation) = 0.5 day; items 2-6 (refactoring consuming code) = 1-1.5 days; items 8-9 (test updates + new tests) = 0.5-1 day. Single PR, single review cycle.

### Functions Relationship to Registry

以下 `build.go` 函数在本次重构中保留但需文档化其与 registry 的关系：

| 函数 | 与 Registry 的关系 | 处理方式 |
|------|-------------------|---------|
| `resolveTestDepsAndInjectReviewDoc` | 被 `injectReviewDocDep` 替代 | 删除，逻辑迁入 post-generation |
| `ResolveFirstTestDep` | 被 `injectCleanCodeDep` + `injectReviewDocDep` 替代 | 删除，三项逻辑分别由 post-generation steps 和 stage-gate 生成流程覆盖 |
| `findHighestGateOrSummary` | 唯一调用方 `ResolveFirstTestDep` 被删除 | 删除（dead code），stage-gate 逻辑由 `GenerateStageGates` 独立处理 |
| `findMaxBusinessTaskID` | 被 `ResolveLastBusinessTask` resolver 替代 | 删除 |
| `ResolveDriftFallbackDep` | 被 `ResolveLastRunTestOrBusiness` resolver 替代 | 删除 |

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
| init() panic in CI environments | Medium | High | Panic-by-default is intentional. Mitigation: `forge validate-pipeline` CI step. `--no-validate` flag for emergency bypass — **warning**: bypass with broken registry produces silently wrong task generation (missing tasks, wrong dependencies, orphaned nodes) with no error signal. Usage logged at WARN level with stack trace and registry checksum, enabling post-incident audit. Not recommended for production CI; intended for hotfix development only. |
| ResolveFirstTestDep logic regression | High | High | Highest-risk deletion. Mitigation: `injectCleanCodeDep` replaces clean-code injection, `injectReviewDocDep` replaces review-doc injection, stage-gate preserved by `GenerateStageGates`. Exhaustive snapshot comparison for all mode/config/surface combinations. |
| Drift/consolidate fallback when test pipeline disabled | Medium | Medium | `ResolveLastRunTestOrBusiness` falls back to last business task when no run-test tasks exist. Test: "quick mode, test disabled, consolidate enabled". |
| Registry reordering silently breaks GenContext progressive population | Low | High | Init-time ordering invariant validation in `ValidatePipelineRegistry` Phase 1: (1) verify `ResolveUpstream` users appear after at least one non-expansion node that populates UpstreamIDs; (2) verify `ResolveLastRunTest` users appear after at least one `TypeTestRun` node; (3) lint rule: any node using `Resolve` (non-nil) that reads `ctx.AllGenerated` or `ctx.UpstreamIDs` must have at least one prior non-gated node in declaration order. Violations produce compile-time actionable error: `"node %s uses ResolveUpstream but has no guaranteed upstream producer before it in declaration order"`. |
| Empty business tasks produces nodes with no dependencies | Low | Low | `CondAlways` nodes generate with empty `DependsOn` via nil-handling rule. Correct (pipeline roots). |

## Success Criteria

1. **Minimal touch points**: Adding a new auto-generated task type requires changes in at most TWO files: `pipeline.go` for the registry entry and `types.go` for the type constant (when a new type is needed; reusing an existing type requires only `pipeline.go`). No other file changes required.
2. **Zero panics**: No `findTaskIndexOrPanic` calls remain; all lookups return errors or are validated at init time
3. **Init-time validation**: Invalid dependency references in the registry fail at CLI startup, not at runtime during `forge task index`
4. **Behavioral parity**: `forge task index` produces identical output for all existing feature configurations after the refactoring
5. **InferType coverage**: All auto-generated task IDs are correctly typed by the registry-derived `InferType`, including surface-expanded variants (e.g., `T-test-gen-scripts-api`) and T-review-doc
6. **GenerateCondition correctness**: T-review-doc is generated in all three scenarios (docs-only, mixed coding+doc, not for coding-only)
7. **Dependency chain completeness**: For all mode/config/surface combinations, the generated task dependency graph contains no orphaned tasks (every generated task has at least one dependency or is the first task in its pipeline chain), and no dangling dependency references (every `DependsOn` entry resolves to an existing task ID). Verified via snapshot tests that include dependency graph validation.
8. **Test coverage >= 80%**: For the new `pipeline.go` and refactored functions
9. **All existing tests pass**: Zero regressions in `autogen_test.go`, `infer_test.go`, `build_test.go`, `claim_test.go`
10. **Escape-hatch count bounded**: Post-generation injection functions <= 5, verified by `ValidatePipelineRegistry` Phase 1. Current count: 2.
