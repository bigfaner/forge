# Design Evaluation: Surface-Aware Justfile — Iteration 1

**Evaluator**: Staff Architect (Adversary)
**Date**: 2026-05-25
**Document**: `docs/features/surface-aware-justfile/design/tech-design.md`
**PRD**: `docs/features/surface-aware-justfile/prd/prd-spec.md`

---

## Dimension 1: Architecture Clarity (170 pts)

### Layer Placement (45/60)

Design explicitly states two layers: "Forge CLI 层（Go 数据模型 + prompt 合成）+ Forge Plugin 层（Skill 文档 + 规则文件）"。

**Deductions**:
- (-10) Layer placement is a single bullet sentence. No explanation of *why* the boundary falls where it does, no data-flow across boundary constraints, and no explicit statement about which layer owns what invariants. The distinction between "CLI" and "Plugin" is assumed knowledge.
- (-5) The three-phase implementation ordering is stated but not justified architecturally — why must Go data model precede skill docs? The dependency chain is self-evident but not architecturally motivated.

### Component Diagram (50/60)

A text-based component diagram is present with data flow arrows.

**Deductions**:
- (-5) The diagram shows `ReadSurfaces()` and `MatchSurface()` as methods inside "forge surfaces CLI" but the actual codebase has these in `pkg/forgeconfig/detect.go` and `pkg/forgeconfig/match.go`, not in the CLI command package. This is a misrepresentation of actual architecture.
- (-5) The diagram shows `prompt.go renderTemplate()` using `{{SURFACE_KEY}}` directly, but current code uses `{{SCOPE}}` and `resolveScope()`. The diagram shows the *target* state without marking it as such, creating ambiguity about what exists vs what is new.

### Dependencies (40/50)

Dependencies section lists internal deps (`ReadSurfaces()`, `MatchSurface()`, task package) and external dep (just >= 1.4.0).

**Deductions**:
- (-10) The document claims "forge surfaces CLI 命令已存在" and states `--json` flag is "新增". But code verification shows `--json` flag does **not exist** in `surfaces.go` at all — no `jsonFlag`, no JSON output code. The document says `(已存在，新增 --json)` which conflates two different things. This is a **factual inaccuracy** in the dependencies section.

**Subtotal: 135/170**

---

## Dimension 2: Interface & Model Definitions (170 pts)

### Interface Signatures Typed (35/60)

Two interfaces are defined: `forge surfaces <path> --json` (CLI) and Surface 规则文件格式 (markdown convention).

**Deductions**:
- (-15) Interface 1 (`forge surfaces --json`) defines JSON output as `[{"surface-key": "admin-panel", "surface-type": "web"}]`. However, current `MatchSurface()` returns only the surface-type (the map value), not the surface-key (the map key). The interface spec requires a **breaking API change** to `MatchSurface()` to return both key and type, but this change is never mentioned in the Data Models or Dependencies sections. The interface assumes a function signature that does not exist.
- (-5) Interface 2 (surface rules file format) uses pseudo-code comments rather than a concrete schema or example file. "编排序列 // 步骤名 | 退出码 | 语义 | 后续动作" is a description of what should be in the file, not a typed interface.
- (-5) No function signatures for the Go-layer changes. The design says `resolveScope()` will be deleted and `renderTemplate()` will use `{{SURFACE_KEY}}` directly, but never provides the new function signatures or call signatures.

### Models Concrete (40/60)

Five Go structs are defined with field-level detail: Task, AutoGenTaskDef, FrontmatterData, AddTaskOpts, TaskState.

**Deductions**:
- (-10) Model 1 (Task) shows `SurfaceKey string` and `SurfaceType string` as new additions with `Scope` deleted. But the current `Task` struct (verified in `types.go`) has `Scope string` with JSON tag `"scope,omitempty"`. The design says `"surface-key,omitempty"` / `"surface-type,omitempty"` for JSON tags but does not address the **breaking JSON deserialization** for existing `index.json` files that contain `"scope": "frontend"`. The `// Scope 字段删除，不保留兼容层` comment acknowledges this but provides no migration strategy.
- (-5) Model 2 (AutoGenTaskDef) claims `Scope` and `TestType` fields will be deleted and replaced with `SurfaceKey`/`SurfaceType`. Current code shows `AutoGenTaskDef` has `Scope string` and `TestType string` but the design doesn't specify how `TaskFromFile()` (which maps AutoGenTaskDef → Task) will propagate the new fields — it currently maps `Scope: d.Scope` (line 841 of autogen.go) but there's no mapping for SurfaceKey/SurfaceType in the design.
- (-5) Model 5 (TaskState) shows new SurfaceKey/SurfaceType fields but current `TaskState` (line 256-273 of types.go) has `Scope string` with comment "mirrors Task.Scope". No explanation of how `TaskState` is populated (it's derived from `Task` during claim) or whether the claim path needs modification.

### Directly Implementable (25/50)

**Deductions**:
- (-15) The design states `MatchSurface()` returns surface-key + surface-type in Interface 1, but the actual `MatchSurface()` function (verified in `match.go`) returns only a single `string` (the surface type). Making it return both key and type requires changing the function signature from `(string, error)` to a struct or dual return. This is a **required but undocumented API change** to an existing function.
- (-10) The `--json` flag for `forge surfaces` is specified as new but has zero implementation detail: no cobra flag registration code, no JSON serialization path, no output format specification beyond the example. The implementer must reverse-engineer the full CLI change from a 3-line interface spec.

**Subtotal: 100/170**

---

## Dimension 3: Error Handling (130 pts)

### Error Types Defined (30/45)

The design references PRD's Error Handling Paths table (7 scenarios) and defines Go-layer error behavior.

**Deductions**:
- (-10) Go layer errors: the design says `build.go` parsing frontmatter without surface-key results in "静默赋空值" (silent empty value, exit 0). This silently degrades data quality — a task without surface-key will be invisible to surface-aware dispatching, but there's no logging or metric to flag this. This contradicts the principle of observable failures.
- (-5) No Go error types are defined. The design references exit codes but doesn't define concrete Go error types or sentinel errors for the new code paths. Compare with existing `surfacesPathError` in `surfaces.go` — no equivalent is specified for the new error scenarios.

### Propagation Strategy Clear (25/45)

**Deductions**:
- (-10) The design says "Go 层错误通过 exit code + stderr 传播给 skill（LLM agent 读取）" and "Skill 层错误通过 just 命令 exit code 和 SKILL.md HARD-GATE 规则控制". But the Go layer `MatchSurface()` error propagation uses cobra's `SilenceErrors: true` pattern with raw stderr writes — the design doesn't explain how new error paths (e.g., `--json` output when config is missing) integrate with this existing pattern.
- (-10) The "Skill 层" error handling is described as "通过 SKILL.md 中 HARD-GATE 规则约束 agent 行为" — this is a documentation constraint enforced on an LLM agent, not a deterministic error propagation mechanism. The design conflates prompt engineering constraints with error handling architecture.

### HTTP Status Codes (40/40)

N/A — no HTTP surface. Auto-full-score.

**Subtotal: 95/130**

---

## Dimension 4: Testing Strategy (130 pts)

### Per-Layer Test Plan (30/45)

A table lists Go CLI unit tests (5 rows) and Skill manual/dry-run tests (3 rows).

**Deductions**:
- (-10) No integration test plan. The critical integration between `forge surfaces --json` CLI output and skill consumption (breakdown-tasks, run-tests parsing JSON via Bash tool) has no test coverage. This is the primary data bridge between Go and Skill layers.
- (-5) The Go unit test rows are all "80%" coverage target but don't specify test scenarios for the **migration path** — e.g., loading an old `index.json` with `"scope": "frontend"` into the new struct without `Scope` field. This is a regression risk.

### Coverage Target Numeric (40/45)

Go 代码 80% is stated explicitly.

**Deductions**:
- (-5) No coverage target for Skill layer (marked N/A). While skill tests are manual, there's no definition of what "passing" means — how many surface types must be verified? What counts as a successful manual test run?

### Test Tooling Named (35/40)

`go test` and `just --dry-run` are named.

**Deductions**:
- (-5) No test tooling for JSON output validation (e.g., `jq` for parsing, or `go test` golden file comparison). The `forge surfaces --json` output format needs deterministic testing but no tooling is specified for this.

**Subtotal: 105/130**

---

## Dimension 5: Breakdown-Readiness (180 pts)

### Components Enumerable (45/65)

The document lists specific files to change: `task/types.go`, `autogen.go`, `frontmatter.go`, `build.go`, `prompt.go`, 10 surface rules files, 2 SKILL.md files.

**Deductions**:
- (-10) Missing `forgeconfig/match.go` from component list. The design requires `MatchSurface()` to return surface-key + surface-type, but `match.go` is never listed as a file to modify. This is a required change that's invisible in the breakdown.
- (-5) Missing `surfaces.go` from component list. Adding `--json` flag requires modifying the cobra command registration and adding JSON output logic, but `surfaces.go` is not listed.
- (-5) The "16 个 prompt 模板 SURFACE_KEY 变量值域同步" is listed as a single item but no breakdown of which 16 templates or what changes each needs.

### Tasks Derivable (40/65)

The three-phase ordering provides a task derivation skeleton.

**Deductions**:
- (-15) The design's three phases don't map to PRD's implementation items cleanly. PRD lists 11 "In Scope" items with checkbox syntax; the design has 3 phases with ~10 items. Cross-referencing reveals gaps:
  - PRD item "死代码清理：extractTestTypeArg()、genScriptBases" has no corresponding design task.
  - PRD item "16 个 prompt 模板 SURFACE_KEY 变量值域同步" is in Phase 3 but no enumeration of which templates or what the variable change looks like.
- (-10) No task-level acceptance criteria. Each phase item is a single sentence with no "done when" definition. An implementer cannot determine when a task is complete without referring back to the PRD.

### PRD AC Coverage (40/50)

The PRD Coverage Map table maps each PRD AC to a design component.

**Deductions**:
- (-5) Story 3 AC: "prompt.go resolveScope() 基于 surfaces map 集合查询而非 projectType 硬编码" — the design says "**删除 resolveScope**，直接读 SurfaceKey" which technically satisfies the AC by removing the function entirely. But the PRD AC says "基于 surfaces map 集合查询", implying a rewrite, not a deletion. The design's approach (deleting resolveScope and reading SurfaceKey directly from the task struct) means surface resolution moves upstream to breakdown-tasks. This is an architectural decision that should be explicitly flagged as a deviation from PRD wording.
- (-5) Story 4 AC: "旧任务文件含 `scope: frontend` 时，run-tests 能正确读取并按默认编排策略执行" — the design says "Scope 字段删除，不保留兼容层" which directly contradicts this PRD AC. Old tasks with `scope: frontend` will have their Scope field silently dropped during deserialization, and no fallback behavior is specified.

**Subtotal: 125/180**

---

## Dimension 6: Security Considerations (80 pts)

### Assessment

PRD has no auth/data/multi-user requirements. Marking N/A.

However, one security-adjacent concern is noted (not scored, but flagged):
- The design mentions surface-key naming constraint `[a-zA-Z0-9_-]` to prevent just recipe name injection, but this constraint is specified only as "由 init-justfile 配方生成时强制执行" — i.e., enforced by the skill's LLM agent, not by Go code. There is no server-side validation in `Task.SurfaceKey` or in the `forge surfaces` CLI output.

**Subtotal: 80/80**

---

## Dimension 7: Implementation Feasibility (140 pts)

### Dependencies Available (30/50)

**Deductions**:
- (-10) The design claims `ReadSurfaces()` and `MatchSurface()` are "已存在" (already exist). While `ReadSurfaces()` exists and works as described, `MatchSurface()` **does not return surface-key** — it returns only the surface-type string. The design assumes a `MatchSurface()` that returns both key and type, which does not exist. This is a **factual error** in the dependency assessment.
- (-10) The `--json` flag for `forge surfaces` is described as "新增" but no existing JSON output infrastructure exists in the CLI. The cobra command setup in `surfaces.go` uses plain `write()` calls — adding JSON output requires either a new output format abstraction or a parallel code path, neither of which is acknowledged.

### Architecture Fits Project Structure (40/50)

**Deductions**:
- (-10) The design modifies `Task` struct to delete `Scope` and add `SurfaceKey`/`SurfaceType`. This is a breaking change to the JSON serialization format of `index.json`. The current `Task` struct has `"scope,omitempty"` JSON tag. Existing `index.json` files with `"scope": "frontend"` will silently lose this data on round-trip (load → save). PRD Story 4 AC explicitly requires backward compatibility: "旧任务文件含 scope: frontend 时，run-tests 能正确读取并按默认编排策略执行". The design directly violates this AC.

### Technical Claims Grounded (25/40)

**Deductions**:
- (-10) "forge surfaces CLI 命令已存在" is misleading. The base command exists but the `--json` output and surface-key-in-output capabilities required by Interface 1 do not. The parenthetical "(已存在，新增 --json)" downplays the significance of this change.
- (-5) The design claims "无新外部依赖" but the Skill layer depends on `just >= 1.4.0` for `[linux]`/`[windows]` recipe attributes. This is listed in Dependencies but not counted as an external dependency, which is inconsistent.

**Subtotal: 95/140**

---

## Cross-Cutting Deductions

### Vague Language (-40)

1. "Skill 层错误通过 SKILL.md HARD-GATE 规则约束 agent 行为" — "约束 agent 行为" is not a deterministic specification. LLM agents do not have guaranteed compliance with documentation rules. (-20)

2. "静默赋空值" for missing surface-key in frontmatter parsing — "静默" means no logging, no metric, no observable signal. This is a silent failure path. (-20)

### PRD AC Gap (-30)

- Story 4 AC: "旧任务文件含 `scope: frontend` 时，run-tests 能正确读取并按默认编排策略执行" — the design's "Scope 字段删除，不保留兼容层" directly violates this acceptance criterion. (-30)

---

## Blindspot Attacks

### Attack 1 [Interface & Model]: MatchSurface() return value mismatch

**Quote**: Interface 1 specifies output `{"surface-key": "admin-panel", "surface-type": "web"}`, and Dependencies section says `MatchSurface()` is "已存在".

**Reality**: `MatchSurface()` in `forgeconfig/match.go` returns `(string, error)` where the string is the surface-type only (the map value). The surface-key (the map key, e.g., "admin-panel") is **never extracted or returned**. The entire surface-key propagation chain from CLI to task frontmatter depends on a function that doesn't provide surface-key. **What must improve**: Define the new `MatchSurface()` signature (or new function) that returns both key and type, and list `match.go` as a file to modify.

### Attack 2 [Implementation Feasibility]: Backward compatibility violation

**Quote**: Design says "Scope 字段删除，不保留兼容层" for Model 1 (Task struct). PRD Story 4 AC says "旧任务文件含 `scope: frontend` 时，run-tests 能正确读取并按默认编排策略执行".

**Analysis**: These are directly contradictory. Deleting the `Scope` field means existing `index.json` files with `"scope": "frontend"` will lose this data on deserialization. There is no migration strategy, no compatibility layer, and no transition plan for existing task files. **What must improve**: Either (a) add a migration path that converts `scope` → `surface-key` on load, or (b) acknowledge the breaking change and update the PRD AC to reflect the new contract, or (c) retain a read-only `Scope` field with a deprecation path.

### Attack 3 [Architecture Clarity]: --json flag not a minor addition

**Quote**: "(已存在，新增 --json)" in the Component Diagram.

**Analysis**: Adding JSON output to `forge surfaces` is not a trivial flag addition. The current command has three sub-invocations (list, query, --types) each with different output formats. The `--json` flag must produce structured output for all three, including the new surface-key field that `MatchSurface()` doesn't currently return. This is a non-trivial CLI change that needs its own interface specification, not a parenthetical note.

### Attack 4 [Error Handling]: Silent data loss on migration

**Quote**: "build.go 解析任务 frontmatter 缺 surface-key | 静默赋空值 | 0"

**Analysis**: When a project migrates to the new model, all existing tasks will have empty `SurfaceKey`. Tasks created before the migration will silently lose their scope information (because `Scope` field is deleted). There is no migration tooling, no validation step, and no logging when this happens. This means `run-tests` will be unable to determine the surface type for pre-migration tasks, violating the "zero regression" goal from PRD ("无 surface 配置的项目输出与当前完全一致").

### Attack 5 [Breakdown-Readiness]: Missing component: match.go

**Quote**: Dependencies says "内部：forgeconfig.ReadSurfaces() + MatchSurface()（已存在）".

**Analysis**: `match.go` must be modified to return surface-key alongside surface-type. This file is never listed in any component enumeration, phase plan, or PRD coverage map. An implementer following the breakdown would miss this required change, leading to compilation errors or incorrect CLI output.

---

## Summary

```
SCORE: 540/1000
DIMENSIONS:
  Architecture Clarity: 135/170
  Interface & Model Definitions: 100/170
  Error Handling: 95/130
  Testing Strategy: 105/130
  Breakdown-Readiness: 125/180
  Security Considerations: 80/80
  Implementation Feasibility: 95/140
ATTACKS:
1. [Interface & Model]: MatchSurface() returns only surface-type, not surface-key — Interface 1 output `{"surface-key": "...", "surface-type": "..."}` requires new return signature — Define new MatchSurface signature and list match.go in component changes.
2. [Implementation Feasibility]: Design deletes Scope field without migration — "Scope 字段删除，不保留兼容层" violates PRD Story 4 AC "旧任务文件含 scope: frontend 时，run-tests 能正确读取" — Add migration path or update PRD AC.
3. [Architecture Clarity]: --json flag understated — "(已存在，新增 --json)" hides non-trivial CLI change requiring structured output for 3 sub-invocations — Provide full CLI interface spec for --json mode.
4. [Error Handling]: Silent data loss on migration — "静默赋空值" for missing surface-key means pre-migration tasks lose scope with no signal — Add migration validation and logging.
5. [Breakdown-Readiness]: match.go missing from component list — "ReadSurfaces() + MatchSurface()（已存在）" implies no changes needed, but MatchSurface must be modified to return surface-key — Add match.go to Phase 1 component changes.
```
