# Eval Report: unify-enum-constants tech-design — Iteration 1

**Date**: 2026-05-28
**Scorer Persona**: Staff Architect
**Document**: `docs/features/unify-enum-constants/design/tech-design.md`

---

## Phase 1: Reasoning Audit (Pre-Score Anchors)

### Problem -> Solution Trace

The problem is clear: 250+ magic string literals for enums (Status, SurfaceType, Priority) scattered across the codebase with no compile-time safety. The solution — typed constants in a leaf `pkg/types/` package — directly addresses this. The mapping is sound.

### Solution -> Evidence Trace

The document claims existing constants in `pkg/feature/constants.go` are "almost never referenced" (Status 5, Priority 0). This is verifiable and the design correctly proposes migrating these. However, the design claims `isTerminalStatus` in `statemachine.go` currently only checks `completed` and `rejected`, while the new `IsTerminalStatus` in `pkg/types/` would also include `skipped`. The business rule `BIZ-task-lifecycle-001` explicitly states terminal states are `completed`, `rejected`, and `skipped`. This means the current code has a bug (missing `skipped`), and the design silently "fixes" it without calling it out. This is a **semantic change masquerading as refactoring**.

### Evidence -> Success Criteria Trace

Success criteria focus on zero magic values, build passing, and tests passing. But the test plan for `IsTerminalStatus` (test scenario 4) says `completed`, `skipped`, `rejected` return true. The existing `isTerminalStatus` in `statemachine.go` does NOT include `skipped`. If the design changes this function's behavior, it changes runtime behavior — violating the "zero behavior change" goal.

### Self-Contradiction Check

The document claims "零行为变更" (zero behavior change) in both PRD and design. But the `IsTerminalStatus` definition in the design expands terminal states beyond what the current code implements. This is a contradiction. Also, the design proposes changing `KnownSurfaceTypes` from `map[string]bool` to `map[types.SurfaceType]bool` — this changes the map key type, which could silently affect callers that pass raw strings (e.g., from config YAML parsing).

### Pre-Score Anchors

1. **`IsTerminalStatus` semantic change not acknowledged**: Design adds `skipped` to terminal states without flagging this as a behavior change.
2. **Map key type change risk**: `KnownSurfaceTypes` changes from `map[string]bool` to `map[types.SurfaceType]bool`. Since `type SurfaceType string`, a raw `string` lookup will fail to compile — but the design does not enumerate all lookup sites or explain the conversion strategy for config-sourced strings.
3. **`surfacePriority` map not addressed**: `detect_surface.go` has `surfacePriority = map[string]int` with 5 surface type string keys. The design does not list this as a change target, even though it uses the same magic values.
4. **`transitionTable` fields**: `TransitionRule.From`/`To` are `string`, but the design only shows the `Transition` struct change (Model 1) while `TransitionRule` in the actual code is the one in `transitionTable` with 28+ string literals. The design conflates two different structs.

---

## Phase 2: Rubric Scoring

### Dimension 1: Architecture Clarity (170 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Layer placement explicit | 52/60 | States "单层重构——仅涉及 `pkg/` 和 `internal/` 层" and clearly positions `pkg/types/` as a leaf package. Dependency direction diagram is present and correct. Loses points for not mentioning that `internal/cmd/task/validate_index.go` imports both `pkg/feature` and `pkg/task` (verified from source), which means the re-export in `pkg/feature` creates a dual-import path that could cause confusion. |
| Component diagram present | 55/60 | ASCII diagram present with correct layering. Shows the three consumer packages and `pkg/types/` as leaf. Loses points for not showing `pkg/forgeconfig/detect.go` (which owns `KnownSurfaceTypes` — a key change point) and for omitting `internal/cmd/surfaces.go` (which re-exports `KnownSurfaceTypes`). |
| Dependencies listed | 45/50 | States "无新增外部依赖" and "纯 Go 标准库代码". Correct. Loses points for not listing the existing internal packages that will need import changes (the Migration Plan table lists files but not the import statements that must change). |

**Dimension 1 Total: 152/170**

### Dimension 2: Interface & Model Definitions (170 pts, db-schema: no)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Interface signatures typed | 50/60 | Interface 1-3 show typed params and return values for `AllStatuses()`, `AllSurfaceTypes()`, `AllPriorities()`, `IsTerminalStatus()`. However: (1) `IsTerminalStatus` has no return type annotation in the signature block — it shows `func IsTerminalStatus(s Status) bool` which is correct, but the design does not show the implementation or define what "terminal" means here; (2) Re-export interface (Interface 4) uses `// ... 其余 Status 常量` comment instead of listing all constants — a developer cannot implement from this. |
| Models concrete | 45/60 | Models 1-3 show before/after struct changes. However: Model 1 (`Transition` struct) does not match the actual code — the real `statemachine.go` uses `TransitionRule` (not `Transition`) with `From string` / `To string` fields, and there is no `Transition` struct. The design appears to reference a different struct name. Model 3 lists `KnownSurfaceTypes` but the actual location is `pkg/forgeconfig/detect.go`, not a field in `ForgeConfig` struct. |
| Directly implementable | 40/50 | A developer can create the `pkg/types/` files from this spec. But the Migration Plan phases 3-5 say "魔法值替换 + 签名升级" without specifying which function signatures change to what. The developer would need to grep each file to discover all change points. |

**Dimension 2 Total: 135/170**

### Dimension 3: Error Handling (130 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Error types defined | 35/45 | States "无新增错误类型" and "类型不匹配在编译期被捕获". This is correct for the happy path. But the design does not address what happens when a runtime string (e.g., from JSON deserialization of `index.json`) does not match any typed constant — the `Task.Status` field typed as `types.Status` will accept any string value at runtime (Go does not enforce enum constraints at runtime). The design should acknowledge this limitation. |
| Propagation strategy clear | 35/45 | States "不适用". Partially correct — but the `TransitionError` struct in `statemachine.go` has `From`/`To` fields typed as `string`. If these change to `types.Status`, error formatting changes. The design does not address this. |
| HTTP status codes mapped | N/A | No API. Full credit for this sub-criterion. |

**Dimension 3 Total: 70/130** (35 + 35 = 70, N/A sub-criterion gets 40/40)

Wait — re-reading the rubric: "If no API: N/A" means the HTTP criterion is N/A. Per rubric, N/A means full credit for that criterion. So: 35 + 35 + 40 = 110.

**Dimension 3 Total: 110/130**

### Dimension 4: Testing Strategy (130 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Per-layer test plan | 40/45 | Lists `pkg/types/` unit tests and global build/regression tests. But no test plan for the migration itself — how to verify each phase is complete before moving to the next. The Migration Plan has 6 phases but no per-phase validation gate. Also, no mention of testing the re-export backward compatibility in `pkg/feature/constants.go`. |
| Coverage target numeric | 40/45 | States "100%" for `pkg/types/` new code and "不降覆盖率" for existing code. Good. But the 100% target is only for new code — the 250+ changed lines in existing packages have no coverage target. |
| Test tooling named | 35/40 | Names `go test` and `go build`. Adequate for this scope. No mention of `go vet` or `golangci-lint` which are part of the project's standard CI pipeline (per CLAUDE.md). |

**Dimension 4 Total: 115/130**

### Dimension 5: Breakdown-Readiness (180 pts — critical gate)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Components enumerable | 55/65 | Can enumerate: 3 new files in `pkg/types/`, 1 modified file in `pkg/feature/`, 10 files in `pkg/task/`, 3 files in `pkg/forgeconfig/`, 13 files in `internal/cmd/`. That's 30 files total. However, the design does not mention `internal/cmd/surfaces.go` (which re-exports `KnownSurfaceTypes`) or `pkg/forgeconfig/detect.go` (which defines `KnownSurfaceTypes`). These are verified in the codebase as change targets. Also, `pkg/task/types.go` contains `NewTaskIndex()` with hardcoded status strings (lines 397-411) — listed but not highlighted as a key change point. |
| Tasks derivable | 50/65 | Each interface maps to at least one task. But the tasks would be: (1) create pkg/types/ files, (2) modify pkg/feature/constants.go, (3-5) migrate each package, (6) verify. The design does not break these down further. For a 30-file migration with 250+ changes, a developer would need more granular task decomposition — e.g., per-file tasks for the largest packages. |
| PRD AC coverage | 45/50 | PRD Coverage Map table covers all 4 user stories and the success criteria. US4 (`validate_index.go` uses `AllStatuses()`) is addressed. However, the PRD's "零行为变更" AC is contradicted by the `IsTerminalStatus` expansion (adding `skipped`). Also, PRD's "叶包" AC is correct and addressed. |

**Dimension 5 Total: 150/180**

### Dimension 6: Security Considerations (80 pts)

Per PRD: "不适用（纯代码组织优化，无安全面影响）". No auth/data/multi-user requirements.

**Dimension 6 Total: 80/80** (N/A — full credit)

### Dimension 7: Implementation Feasibility (140 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Dependencies available | 48/50 | No new external dependencies. `pkg/types/` is pure Go. Verified correct. Loses 2 points for not verifying that the type alias pattern (`type Status = types.Status`) works with the existing `pkg/feature/constants.go` untyped constants (currently `StatusPending = "pending"` without a type — these are untyped string constants, and re-typing them requires changing the declaration). |
| Architecture fits project structure | 42/50 | The `pkg/types/` leaf package pattern fits the existing `cmd -> internal -> pkg` dependency direction. However, the design does not address that `pkg/forgeconfig/detect_surface.go` has 100+ string literals in signal maps (`packageJSONSignals`, `goModSignals`, etc.) that map dependency names to surface types — these are NOT enum magic values (they map arbitrary strings to surface types), but the design does not clarify this boundary. A developer might waste time trying to type these map values. |
| Technical claims grounded | 35/40 | Claims `type X string` preserves JSON serialization — correct and grounded in Go spec. Claims "零运行时开销" — correct. Loses points for claiming "零行为变更" when `IsTerminalStatus` changes behavior (adds `skipped` to terminal check), and for the `Transition` vs `TransitionRule` naming discrepancy which suggests the design was not verified against actual code for all model definitions. |

**Dimension 7 Total: 125/140**

---

## Phase 3: Blindspot Hunt

### [blindspot-1] `isTerminalStatus` in `statemachine.go` vs `IsTerminalStatus` in `pkg/types/`

The design defines `IsTerminalStatus` as a public function in `pkg/types/` with terminal states = `{completed, skipped, rejected}`. But the existing private `isTerminalStatus` in `statemachine.go` only checks `{completed, rejected}`. When migrating, should the statemachine's `isTerminalStatus` call `types.IsTerminalStatus()`? If yes, this changes behavior — `skipped` tasks would now be considered terminal by `isActiveFixTask`, changing which fix tasks are considered "active". This is a **silent behavior regression** that no rubric dimension catches.

### [blindspot-2] `surfacePriority` map in `detect_surface.go` not in Migration Plan

The `surfacePriority = map[string]int{"web": 1, "mobile": 2, ...}` in `pkg/forgeconfig/detect_surface.go` uses the same 5 surface type string literals that the design proposes to centralize. But the Migration Plan phase 4 only lists 3 files: `detect_surface.go`, `detect.go`, `execution_order.go`. The `surfacePriority` map is in `detect_surface.go`, but the design does not explicitly say to change its key type to `types.SurfaceType`. The signal maps (`packageJSONSignals`, etc.) map dependency names to surface types — these VALUE strings should become `types.SurfaceType`, but their KEYS are arbitrary dependency names. The design must clarify this distinction.

### [blindspot-3] `NewTaskIndex()` hardcoded status list

In `pkg/task/types.go` lines 397-411, `NewTaskIndex()` hardcodes `StatusEnum` and `PriorityEnum` as string slices. After migration, should these use `types.AllStatuses()` and `types.AllPriorities()`? The design lists `pkg/task/types.go` as a change target but does not call out this specific function. If changed, the enum values' JSON representation must remain `[]string` (not `[]types.Status`) for backward compatibility — the design does not address this.

### [blindspot-4] No rollback plan

The Migration Plan has 6 phases but no rollback strategy. If phase 3 (pkg/task migration) introduces a compilation error that takes hours to debug, the developer cannot easily revert to pre-phase-3 state. The design should recommend committing after each phase (which the migration plan implies but does not state).

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Architecture Clarity | 152 | 170 |
| Interface & Model Definitions | 135 | 170 |
| Error Handling | 110 | 130 |
| Testing Strategy | 115 | 130 |
| Breakdown-Readiness | 150 | 180 |
| Security Considerations | 80 | 80 |
| Implementation Feasibility | 125 | 140 |
| **Total** | **867** | **1000** |

---

## Attack Summary

1. [Interface & Model Definitions]: `Transition` struct name does not exist in codebase — actual struct is `TransitionRule` in `statemachine.go`. Model 1 references a non-existent type. — Quote: "状态转移表中的 `From`/`To` 字段类型从 `string` 升级为 `types.Status`" with `type Transition struct` — must verify all model names against actual code before implementation.

2. [Interface & Model Definitions]: Re-export interface uses `// ... 其余 Status 常量` comment instead of listing all constants. — Quote: Interface 4 shows `StatusPending = types.StatusPending` and `StatusInProgress = types.StatusInProgress` then `// ... 其余 Status 常量` — must enumerate all 7 Status constants and 3 Priority constants explicitly for direct implementability.

3. [Implementation Feasibility]: `IsTerminalStatus` adds `skipped` to terminal states but existing `isTerminalStatus` in `statemachine.go` only checks `completed` and `rejected`. — Quote: "IsTerminalStatus：`completed`、`skipped`、`rejected` 返回 true" vs actual code `return status == "completed" || status == "rejected"` — must acknowledge this as a behavior change, not pure refactoring, and assess impact on `isActiveFixTask`/`canAutoUnblock` logic.

4. [Architecture Clarity]: Component diagram omits `pkg/forgeconfig/detect.go` (owner of `KnownSurfaceTypes`) and `internal/cmd/surfaces.go` (re-exports `KnownSurfaceTypes`). — Both files contain enum magic values that must change but are not shown in the diagram or Migration Plan phase 4/5 file lists.

5. [Breakdown-Readiness]: Migration Plan lists 13 files for `internal/cmd/` but does not mention `internal/cmd/surfaces.go` which re-exports `KnownSurfaceTypes = forgeconfig.KnownSurfaceTypes`. — Quote: Phase 5 lists specific files but omits `surfaces.go` — this file's re-export will break when `KnownSurfaceTypes` changes from `map[string]bool` to `map[types.SurfaceType]bool`.

6. [Error Handling]: `TransitionError` struct in `statemachine.go` has `From string` and `To string` fields. If these change to `types.Status`, the error message formatting changes, which could break error assertion in tests. — Quote: Error Handling section says "不适用" — must analyze all error structs that hold enum values.

7. [Testing Strategy]: No per-phase validation gate in the 6-phase Migration Plan. — Quote: Phase 6 is "验证" but phases 1-5 have no intermediate build/test check — must add "go build + go test" gate after each phase to catch issues early.
