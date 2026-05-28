# Eval Report: unify-enum-constants tech-design — Iteration 2

**Date**: 2026-05-28
**Scorer Persona**: Staff Architect
**Document**: `docs/features/unify-enum-constants/design/tech-design.md`
**Previous Score**: 867/1000 (Iteration 1)

---

## Iteration-1 Issues Remediation Check

| # | Iteration-1 Attack | Status | Evidence |
|---|---------------------|--------|----------|
| 1 | `Transition` struct name does not exist — actual is `TransitionRule` | **FIXED** | Model 1 now reads "`statemachine.go` 中 `TransitionRule` 结构体" with correct struct name |
| 2 | Re-export uses `// ... 其余 Status 常量` instead of listing all constants | **FIXED** | Interface 4 now enumerates all 7 Status constants and all 3 Priority constants explicitly |
| 3 | `IsTerminalStatus` adds `skipped` silently — behavior change not acknowledged | **FIXED** | Interface 1 now has explicit "行为变更声明" blockquote acknowledging the semantic change, calling it "对齐业务规则的 bug 修复" and specifying impact on `isActiveFixTask`/`canAutoUnblock` |
| 4 | Component diagram omits `pkg/forgeconfig/detect.go` and `internal/cmd/surfaces.go` | **FIXED** | Component diagram now shows `surfaces.go (重导出 KnownSurfaceTypes)` in `internal/cmd/*` box, and `detect.go` in `pkg/forgeconfig` box |
| 5 | `surfaces.go` re-export not in Migration Plan | **FIXED** | Phase 5 now explicitly includes `surfaces.go` with detailed description of map key type change and `runSurfacesTypes` lookup logic impact |
| 6 | `TransitionError` fields not addressed in Error Handling | **FIXED** | Error Handling now includes "Error Structs Holding Enum Values" table analyzing `TransitionError` and `ActiveFixExistsError`, plus explanation that `Error()` format is preserved via `%s` formatting |
| 7 | No per-phase validation gate in Migration Plan | **FIXED** | Migration Plan now has explicit Gate rows after each phase with `go build + go test` commands |

**All 7 iteration-1 attacks addressed.** This is strong remediation.

---

## Phase 1: Reasoning Audit (Pre-Score Anchors)

### Problem -> Solution Trace

Problem: 250+ magic string literals across enums with no compile-time safety. Solution: typed constants in `pkg/types/` leaf package, per-package batch migration. Mapping remains sound and direct.

### Solution -> Evidence Trace

1. **`IsTerminalStatus` behavioral change**: Now explicitly declared in Interface 1 with full impact analysis. The design correctly identifies that `add.go`'s `terminalStatuses` map already includes `skipped` (verified in source: `"skipped": true` at line 37 of `add.go`), while `statemachine.go`'s `isTerminalStatus` does not. The unification aligns two divergent implementations. However, the design says "将 skipped 纳入终态，影响 `isActiveFixTask`/`canAutoUnblock` 逻辑——skipped 的 fix-task 将被视为非活跃" — but `isActiveFixTask` calls `!isTerminalStatus(t.Status)`, meaning adding `skipped` to terminal states means a fix-task with status `skipped` will be considered **inactive** (not active). This is the **correct** behavior per business rules but the design description is slightly misleading: it says "skipped 的 fix-task 将被视为非活跃" as if this is new, when in fact the current `hasActiveFixTasks` in `add.go` already treats skipped as terminal via its own `terminalStatuses` map. The inconsistency is between two files in the same package, and the design unifies them correctly.

2. **`TransitionError` field type change**: Design states `Error()` output format is unchanged because `%s` formatting still outputs the raw string value. Verified correct: `type Status string` means `fmt.Sprintf("... %s ...", types.StatusCompleted)` outputs `"completed"` identically.

3. **`surfaces.go` lookup logic**: Phase 5 describes that `KnownSurfaceTypes[typ]` lookup needs `typ` converted to `types.SurfaceType`. Verified: in `runSurfacesTypes` line 201, `KnownSurfaceTypes[typ]` where `typ` is a `string` (from `surfaces map[string]string` iteration). After migration, `KnownSurfaceTypes` becomes `map[types.SurfaceType]bool`, so `typ` must be cast: `KnownSurfaceTypes[types.SurfaceType(typ)]`. The design notes this but does not address the **broader pattern**: `surfaces map[string]string` itself holds surface type values as `string`. When surface types come from YAML config parsing, they arrive as `string`. The design does not specify where the `string -> types.SurfaceType` conversion boundary sits for the entire `ReadSurfaces` -> `ValidateSurfaceTypes` pipeline.

### Self-Contradiction Check

1. The "零行为变更" claim in Overview (line 11) now partially contradicts the explicit "行为变更声明" in Interface 1. The Overview still says "保持零行为变更" but the design itself acknowledges a behavioral change in `IsTerminalStatus`. This is a **minor inconsistency** — the Overview should say "最小行为变更" or "除 IsTerminalStatus 对齐外零行为变更".

2. The wildcard `"*"` in `TransitionRule.From`/`To` is typed as `types.Status` but `*` is not a valid Status value. The design acknowledges this in Model 1's note about `types.Status("*")` or `StatusAny`, but does not resolve the Open Question — should a constant `StatusAny = Status("*")` be defined? This leaves ambiguity: if a developer uses `types.Status("*")` literal, it's still a magic value. If a constant is defined, it adds a value to the Status type that is not a real status. The design should state a clear decision.

### Pre-Score Anchors

1. **Overview "零行为变更" vs Interface 1 "行为变更声明"**: Contradictory framing in the same document.
2. **Wildcard `"*"` handling unresolved**: Model 1 acknowledges the issue but defers the decision.
3. **`string -> types.SurfaceType` conversion boundary**: The design addresses `surfaces.go` lookup but not the upstream `ReadSurfaces` / `ValidateSurfaceTypes` pipeline where YAML-parsed strings first enter the system.

---

## Phase 2: Rubric Scoring

### Dimension 1: Architecture Clarity (170 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Layer placement explicit | 57/60 | Clearly states "单层重构" and positions `pkg/types/` as leaf package. Dependency direction diagram is explicit with correct arrows. Loses 3 points for the Overview's "零行为变更" claim contradicting the actual behavioral change acknowledged later — this mischaracterization could lead implementers to underestimate risk. |
| Component diagram present | 58/60 | ASCII diagram now includes `surfaces.go` re-export and `detect.go` / `detect_surface` / `execution_order` in `pkg/forgeconfig`. Comprehensive. Loses 2 points for not showing the data flow between `ReadSurfaces` (returns `map[string]string`) and the typed lookup boundary — this is where the `string -> types.SurfaceType` conversion happens and is architecturally significant. |
| Dependencies listed | 48/50 | States "无新增外部依赖" and "纯 Go 标准库代码". Correct. Migration Plan per-phase file lists are now complete. Loses 2 points for not listing the specific import statements that change per phase (e.g., Phase 3 adds `"forge-cli/pkg/types"` to 10 files in `pkg/task/`) — this would make implementation more mechanical. |

**Dimension 1 Total: 163/170**

### Dimension 2: Interface & Model Definitions (170 pts, db-schema: no)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Interface signatures typed | 57/60 | All four interfaces show typed params, return values, and full constant enumerations. `AllStatuses() []Status`, `AllSurfaceTypes() []SurfaceType`, `AllPriorities() []Priority`, `IsTerminalStatus(s Status) bool` — all clear. Loses 3 points for not specifying the return behavior of `AllXxx()` functions: should they return a new slice each call (safe) or a shared slice (efficient but mutable)? This matters for correctness. |
| Models concrete | 55/60 | Model 1 now correctly names `TransitionRule` (not `Transition`). Model 2 shows `Task` struct field changes. Model 3 covers `KnownSurfaceTypes`, `surfacePriority`, `defaultExecutionOrder`. Loses 5 points for Model 1's unresolved wildcard handling: the `After` struct shows `From types.Status` but the transition table contains `"*"` wildcards. The note says "需保留为 `types.Status("*")` 或定义常量 `StatusAny`" — this is a design decision left open, not a concrete model. Also loses points for not modeling the `TaskIndex.StatusEnum` / `TaskIndex.PriorityEnum` fields (`[]string` → should they become `[]types.Status`/`[]types.Priority`? Or stay `[]string` for JSON compat?). |
| Directly implementable | 46/50 | A developer can create `pkg/types/` files and the re-export from these specs. The per-phase Migration Plan with gates makes execution plan clear. Loses 4 points for: (1) wildcard `"*"` handling decision deferred; (2) `StatusEnum`/`PriorityEnum` field migration strategy not specified; (3) `ValidateTransition` and `matchRule` function signature changes not shown — these are the most complex signature upgrades (they compare `string` against `TransitionRule.From`/`To`, and the comparison logic changes when these become `types.Status`). |

**Dimension 2 Total: 158/170**

### Dimension 3: Error Handling (130 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Error types defined | 42/45 | States "无新增错误类型" — correct. The "Error Structs Holding Enum Values" table is a strong addition, listing `TransitionError` and `ActiveFixExistsError` with field-level impact. Loses 3 points for not addressing the runtime string-to-typed-constant boundary: when `index.json` is deserialized, `Task.Status` is `types.Status` (which is `string`), and Go allows any string value at runtime. If a corrupt JSON file contains `"status": "unknown"`, `types.IsTerminalStatus(types.Status("unknown"))` returns false — which is safe but undocumented. The design should acknowledge this limitation. |
| Propagation strategy clear | 42/45 | `TransitionError` field type change is now fully analyzed. `Error()` format preservation via `%s` is correctly explained. Phase 3 gate recommends running tests after `statemachine.go` migration. Loses 3 points for not specifying whether `TransitionError.From`/`To` should be typed as `types.Status` (which means callers constructing this error must pass `types.Status` values) — the design implies yes but does not explicitly state this, and the `ValidateTransition` function signature (`current, target string`) would also need to change to `types.Status`. |
| HTTP status codes mapped | N/A | No API. Full credit. |

**Dimension 3 Total: 124/130** (42 + 42 + 40)

### Dimension 4: Testing Strategy (130 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Per-layer test plan | 43/45 | Per-layer test table with specific test types and tools. 5 key test scenarios specified. Per-phase gates in Migration Plan serve as intermediate validation. Loses 2 points for not specifying a test for the `IsTerminalStatus` behavioral change: the design acknowledges this is a bug fix, but Test Scenario 4 only says "completed、skipped、rejected 返回 true" — it should also test that the behavioral change propagates correctly to `canAutoUnblock` (e.g., "source task with status=skipped allows auto-unblock"). |
| Coverage target numeric | 43/45 | "100% for new code, 不降覆盖率 for existing". Clear. Loses 2 points for not specifying coverage for the behavioral change: if `isTerminalStatus` gains `skipped`, existing tests that verify `canAutoUnblock` behavior may need updating. The design should state whether existing tests cover the `skipped` case or if new tests are needed. |
| Test tooling named | 38/40 | Names `go test`, `go build`, `go vet` (Phase 6 now includes `go vet`). Adequate. Loses 2 points for not mentioning `golangci-lint` which the project's CLAUDE.md lists as a standard CI tool (`golangci-lint run ./...`). |

**Dimension 4 Total: 124/130**

### Dimension 5: Breakdown-Readiness (180 pts — critical gate)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Components enumerable | 60/65 | All change targets now identified: 3 new files in `pkg/types/`, `pkg/feature/constants.go`, 10 files in `pkg/task/`, 3 files in `pkg/forgeconfig/`, 15 files in `internal/cmd/` (including `surfaces.go`). Loses 5 points for not calling out `pkg/task/types.go`'s `NewTaskIndex()` function specifically — it hardcodes `StatusEnum: []string{"pending", "in_progress", ...}` and `PriorityEnum: []string{"P0", "P1", "P2"}`. After migration, should these use `types.AllStatuses()` and `types.AllPriorities()`? And should `StatusEnum` field type change from `[]string` to `[]types.Status`? The JSON serialization must remain `[]string` for backward compatibility with existing `index.json` files. This is a non-trivial design decision hidden in a file that is listed but not analyzed. |
| Tasks derivable | 58/65 | Each interface maps to tasks. Migration phases provide structure. But for a 30+ file migration, the task granularity could be finer — e.g., Phase 3 (`pkg/task/`) touches 10 files and is the largest phase. This could be split into: (3a) `statemachine.go` + `add.go` (core state logic), (3b) `state.go` + `deps.go` + `build.go` (state-dependent operations), (3c) `types.go` + `record.go` + `index.go` + `tasktemplate.go` + `autogen.go` (data layer). The design does not suggest this sub-decomposition. |
| PRD AC coverage | 48/50 | PRD Coverage Map covers all user stories and success criteria. The behavioral change in `IsTerminalStatus` is now acknowledged with explicit linkage to `BIZ-task-lifecycle-001`. Loses 2 points for the Overview's "零行为变更" claim which slightly misaligns with the PRD AC — the PRD likely expects zero behavior change, and the design should have called this out as a PRD amendment rather than silently fixing it. |

**Dimension 5 Total: 166/180**

### Dimension 6: Security Considerations (80 pts)

No auth/data/multi-user requirements. Pure code reorganization.

**Dimension 6 Total: 80/80** (N/A — full credit)

### Dimension 7: Implementation Feasibility (140 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Dependencies available | 49/50 | No new external dependencies. `pkg/types/` is pure Go. Type alias pattern (`type Status = types.Status`) is correctly specified for re-export. Loses 1 point for not verifying that untyped string constants in the current `pkg/feature/constants.go` (e.g., `StatusPending = "pending"` without an explicit type) will work correctly with the type alias pattern — the design shows the re-export correctly using typed constants from `types`, but does not note that the current constants are untyped and will need re-typing. |
| Architecture fits project structure | 47/50 | `pkg/types/` leaf package fits `cmd -> internal -> pkg` direction. Per-phase migration with gates is practical. Loses 3 points for the `detect_surface.go` signal maps: `packageJSONSignals` maps dependency names like `"react"` to surface types like `"web"`. The design's Model 3 implies `surfacePriority` key type changes to `types.SurfaceType`, but does not address the signal maps. Should their VALUES become `types.SurfaceType`? Their KEYS are arbitrary dependency names (not enum values). The design does not clarify this boundary, and a developer might waste time or miss the distinction. |
| Technical claims grounded | 38/40 | `type X string` preserves JSON serialization — grounded. `Error()` format preservation — grounded and verified. `IsTerminalStatus` behavioral change now acknowledged with correct analysis. Loses 2 points for the unresolved wildcard `"*"` question — `types.Status("*")` is a valid expression but `"*"` is semantically different from a status value. The design acknowledges this but does not make a decision. |

**Dimension 7 Total: 134/140**

---

## Phase 3: Blindspot Hunt

### [blindspot-1] `StatusEnum`/`PriorityEnum` field type in `TaskIndex` — JSON backward compatibility

**Quote**: Model 2 shows `Status types.Status` and `Priority types.Priority` for the `Task` struct, but does not address `TaskIndex.StatusEnum []string` and `TaskIndex.PriorityEnum []string` fields in `pkg/task/types.go` (lines 269-270, 282-283).

**Impact**: `NewTaskIndex()` hardcodes `StatusEnum: []string{"pending", "in_progress", ...}` and `PriorityEnum: []string{"P0", "P1", "P2"}`. These are serialized to JSON as `"statusEnum": ["pending", ...]`. After migration, there are three options: (a) change field type to `[]types.Status` (JSON serialization is identical since `type Status string`), (b) keep `[]string` and populate from `types.AllStatuses()` with conversion, or (c) replace `StatusEnum`/`PriorityEnum` fields entirely with calls to `types.AllStatuses()`/`types.AllPriorities()` at consumption sites. Each option has different implications for the `MarshalJSON`/`UnmarshalJSON` methods. The design must specify which approach to take.

### [blindspot-2] `ReadSurfaces` returns `map[string]string` — conversion boundary undefined

**Quote**: Phase 5 describes `surfaces.go`'s `KnownSurfaceTypes[typ]` lookup needing `types.SurfaceType(typ)` conversion, but does not address the upstream data source.

**Impact**: `forgeconfig.ReadSurfaces()` returns `map[string]string` (path -> surface type). This map is populated from YAML config parsing, where all values are raw strings. The design changes `KnownSurfaceTypes` to `map[types.SurfaceType]bool` and `surfacePriority` to `map[types.SurfaceType]int`. When `ReadSurfaces` returns a `string` value for a surface type, every lookup like `surfacePriority[surfaceType]` or `KnownSurfaceTypes[surfaceType]` will fail to compile unless `surfaceType` is explicitly cast. The design should state where the canonical conversion boundary is: at `ReadSurfaces` return (changing its signature to return typed values), at each consumption site, or via a validation function like `ValidateSurfaceTypes`.

### [blindspot-3] `matchRule` function comparison logic

**Quote**: Model 1 shows `TransitionRule.From types.Status` and notes wildcard `"*"` handling, but does not show the `matchRule` function that performs the comparison.

**Impact**: The current `matchRule` function in `statemachine.go` compares `string` values. After `From`/`To` become `types.Status`, the comparison `rule.From == current` where `current` is `types.Status` and `rule.From` is also `types.Status` works naturally. But the wildcard comparison `rule.From == "*"` becomes `rule.From == types.Status("*")` — or should it be compared against a `StatusAny` constant? The design notes this in Model 1 but leaves the decision open. More critically, `ValidateTransition` takes `current, target string` — these parameters must become `types.Status`, which changes the function's public API and every call site. The design does not show this signature change.

### [blindspot-4] Overview "零行为变更" vs Interface 1 "行为变更声明" — framing contradiction

**Quote**: Overview line 11 states "保持零行为变更" while Interface 1 line 88 has explicit "行为变更声明" acknowledging the `IsTerminalStatus` change.

**Impact**: This is a documentation inconsistency that could cause confusion. An implementer reading the Overview first may proceed assuming zero risk, then discover the behavioral change mid-implementation. The Overview should be updated to reflect the actual scope: "除 IsTerminalStatus 对齐业务规则外，保持零行为变更".

---

## Score Summary

| Dimension | Score | Max | Delta from Iter 1 |
|-----------|-------|-----|--------------------|
| Architecture Clarity | 163 | 170 | +11 |
| Interface & Model Definitions | 158 | 170 | +23 |
| Error Handling | 124 | 130 | +14 |
| Testing Strategy | 124 | 130 | +9 |
| Breakdown-Readiness | 166 | 180 | +16 |
| Security Considerations | 80 | 80 | 0 |
| Implementation Feasibility | 134 | 140 | +9 |
| **Total** | **949** | **1000** | **+82** |

---

## Attack Summary

1. [Interface & Model Definitions]: `StatusEnum`/`PriorityEnum` field migration strategy not specified. — Quote: Model 2 shows `Status types.Status` but does not address `TaskIndex.StatusEnum []string` (lines 269-270 in `types.go`). `NewTaskIndex()` hardcodes string slices that must either become typed or be replaced with `types.AllStatuses()` calls. — Must specify which approach for these fields and whether JSON output format is preserved.

2. [Architecture Clarity]: Overview claims "保持零行为变更" but Interface 1 has explicit "行为变更声明". — Quote: line 11 "保持零行为变更" vs line 88 "行为变更声明：IsTerminalStatus 包含 completed、skipped、rejected 三种终态...这是对齐业务规则的 bug 修复". — Must update Overview to acknowledge the exception.

3. [Interface & Model Definitions]: Wildcard `"*"` handling in `TransitionRule` unresolved — a design decision deferred, not made. — Quote: Model 1 note "需保留为 types.Status(\"*\") 或定义常量 StatusAny". — Must decide: either define `StatusAny` constant or document that `types.Status("*")` literal is acceptable. Leaving this open means the implementer must make an architectural choice.

4. [Architecture Clarity]: `string -> types.SurfaceType` conversion boundary for `ReadSurfaces` pipeline not defined. — Quote: Phase 5 describes `surfaces.go` lookup conversion but `ReadSurfaces()` returns `map[string]string` from YAML parsing. Every downstream lookup (`surfacePriority[surfaceType]`) requires explicit cast. — Must specify where the canonical conversion happens.

5. [Interface & Model Definitions]: `ValidateTransition` and `matchRule` function signature changes not shown. — Quote: Model 1 changes `TransitionRule.From`/`To` to `types.Status` but `ValidateTransition(current, target string, role TransitionRole)` takes `string` params that must also become `types.Status`. — Must show the signature change for these key functions.

6. [Testing Strategy]: No test specified for `canAutoUnblock` behavioral change propagation. — Quote: Test Scenario 4 only tests `IsTerminalStatus` directly: "completed、skipped、rejected 返回 true". The behavioral change (skipped fix-tasks now inactive) should have an integration-level test in `pkg/task/`. — Must add test scenario for `canAutoUnblock` with source task status=skipped.

7. [Testing Strategy]: `golangci-lint` not included in Phase 6 verification despite being a project-standard CI tool. — Quote: Phase 6 gate is "go build ./... && go test -race -cover ./... && go vet ./..." but project CLAUDE.md lists `golangci-lint run ./...` as standard. — Must add `golangci-lint run ./...` to Phase 6 gate.
