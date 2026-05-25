# Design Evaluation: Surface-Aware Justfile — Iteration 3

**Evaluator**: Staff Architect (Adversary)
**Date**: 2026-05-25
**Document**: `docs/features/surface-aware-justfile/design/tech-design.md`
**PRD**: `docs/features/surface-aware-justfile/prd/prd-spec.md`
**Previous**: `docs/features/surface-aware-justfile/eval/design/iteration-2.md` (650/1000)

---

## Iteration 2 Issues — Resolution Check

| # | Iteration 2 Issue | Resolution Status |
|---|-------------------|-------------------|
| 1 | PRD AC conflict: design says "不保留兼容层" but PRD Story 4 AC requires backward compat | **RESOLVED** — Migration Notes now specifies blocking exit 2 check when old `scope` field detected, forcing users to regenerate tasks. Acknowledges PRD AC must be updated. |
| 2 | Surface 规则文件格式 lacks concrete example (10 files) | **RESOLVED** — Interface 2 now provides complete `web.md` example with编排序列 table, 配方调用契约 table, journey 过滤策略, and implementation constraints. |
| 3 | --json stderr JSON conflicts with TECH-error-handling-001 | **RESOLVED** — Explicit "stderr 格式覆盖声明" added: `--json` mode overrides plain-text stderr convention with structured JSON, justified by machine-consumption requirement. |
| 4 | Phase 2-3 lack Phase 1-level file-level breakdown | **RESOLVED** — Phase 2 now has 5-row table with specific file paths and cross-references. Phase 3 has 12-row table enumerating all files including individual rule files. |
| 5 | Migration validation is passive (warning-only) | **RESOLVED** — Changed from warning to blocking exit 2: `build.go` detects old `scope` field → stderr error → exit 2. |
| 6 | 16 prompt templates not enumerated | **NOT RESOLVED** — Phase 3 lists "skills/*/templates/*.md（共 16 个 prompt 模板）" as a single glob row. No enumeration of which 16 templates. |
| 7 | Model 2 (AutoGenTaskDef) missing YAML/JSON tags | **NOT RESOLVED** — Model 2 still does not specify YAML/JSON tags for the new SurfaceKey/SurfaceType fields. |
| 8 | Model 5 (TaskState) population path unspecified | **NOT RESOLVED** — No description of how TaskState.SurfaceKey/SurfaceType get populated from Task. |

---

## Dimension 1: Architecture Clarity (170 pts)

### Layer Placement (55/60)

Design states: "Forge CLI 层（Go 数据模型 + prompt 合成）+ Forge Plugin 层（Skill 文档 + 规则文件）". Phase Component Map provides file-level assignment to layers.

**Deductions**:
- (-5) The boundary contract between CLI and Plugin layers is still implicit. The sole bridge is `forge surfaces --json` CLI output consumed by skills via Bash tool. This architectural significance — that the CLI is the *only* data channel from Go layer to skill layer — is not called out as a constraint. A future contributor could inadvertently add a direct file-read bridge without understanding the boundary semantics.

### Component Diagram (55/60)

Component diagram shows the full data flow with `MatchSurface()` marked as "(需修改: 返回 SurfaceMatch{Key,Type})". Interface 1a now provides the exact struct definition.

**Deductions**:
- (-5) The diagram shows `prompt.go renderTemplate()` with `{{SURFACE_KEY}} direct` and `{{TEST_TYPE_ARG}} direct` as two separate template variable substitutions. The document does not clarify what `TEST_TYPE_ARG` is — is it a new variable? A renamed variable? The relationship between these two template variables and the deleted `{{SCOPE}}` / `{{TEST_TYPE}}` is stated in the Phase 3 table ("`{{SCOPE}}` → `{{SURFACE_KEY}}`，`{{TEST_TYPE}}` → `{{TEST_TYPE_ARG}}`") but not reflected in the diagram, creating a disconnect.

### Dependencies (45/50)

Dependencies section correctly distinguishes "已存在" / "需修改签名" / "需新增".

**Deductions**:
- (-5) "无新外部依赖" statement while listing `just >= 1.4.0` as external dep remains slightly misleading. The statement should say "无新外部包依赖" to avoid confusion.

**Subtotal: 155/170**

---

## Dimension 2: Interface & Model Definitions (170 pts)

Using "no er-diagram" scoring criteria (db-schema: "no").

### Interface Signatures Typed (55/60)

Interface 1a provides exact Go function signatures (old and new) with typed `SurfaceMatch` struct including field semantics. Interface 1b provides cobra flag registration, three-mode JSON output table with concrete examples, and explicit error format specification.

**Deductions**:
- (-5) Interface 1a defines `SurfaceMatch` struct but does not specify its package. The struct is referenced in `match.go` which lives in `pkg/forgeconfig/` — is it `forgeconfig.SurfaceMatch`? A shared type? The import path matters for downstream consumers (surfaces.go, autogen.go, etc.).

### Models Concrete (50/60)

Five Go structs defined with field-level detail. JSON tags specified for Model 1.

**Deductions**:
- (-5) Model 2 (AutoGenTaskDef) does not specify JSON/YAML tags for SurfaceKey/SurfaceType. Current code has `yaml:"scope"` and `json:"scope,omitempty"` for the Scope field. The new fields need equivalent tags. This was flagged in iteration 2 and remains unfixed.
- (-5) Model 5 (TaskState) shows SurfaceKey/SurfaceType fields but does not document how these are populated. Current TaskState is populated during task claim from the Task struct — does the claim path propagate these fields automatically or does it need modification? This was flagged in iteration 2 and remains unfixed.

### Directly Implementable (40/50)

**Deductions**:
- (-5) Interface 1a says "所有现有调用方（runSurfacesQuery 等）需适配新签名，从 SurfaceMatch.Type 取值" but does not enumerate which callers exist. An implementer must search the codebase to find all call sites of `MatchSurface()`.
- (-5) The `build.go` migration check ("扫描含 scope 字段但不含 surface-key 的任务") specifies the detection logic but not the implementation location. Does this check live in `build.go`'s existing frontmatter parsing function? A new validation function? The error message format is specified but the code structure is not.

**Subtotal: 145/170**

---

## Dimension 3: Error Handling (130 pts)

### Error Types Defined (40/45)

Go layer error table now defines 4 scenarios with clear exit codes. The migration check now uses blocking exit 2 instead of silent empty-value assignment.

**Deductions**:
- (-5) No Go error types or sentinel errors defined. The document specifies exit codes and stderr messages but does not define concrete Go error types. For example, what Go error wraps the "surfaces 配置缺失" condition? Is it `fmt.Errorf("...")` or a typed error? The existing codebase uses `surfacesPathError` in `surfaces.go` — no equivalent convention is specified for new error paths.

### Propagation Strategy Clear (40/45)

Go→Skill propagation via exit code + stderr is well-defined. The `--json` stderr format override is now explicitly declared.

**Deductions**:
- (-5) The Skill layer propagation — "Skill 层错误通过 SKILL.md HARD-GATE 规则控制" — remains a documentation constraint on an LLM agent, not a deterministic error propagation mechanism. This is an architectural limitation of the Forge distribution model (skills are consumed by LLM agents, not deterministic code), but the design does not acknowledge this limitation or propose any verification mechanism.

### HTTP Status Codes (40/40)

N/A — no HTTP surface. Full score.

**Subtotal: 120/130**

---

## Dimension 4: Testing Strategy (130 pts)

### Per-Layer Test Plan (40/45)

Test table lists 5 Go unit test rows and 3 Skill test rows. Key Test Scenarios section adds 6 concrete scenarios.

**Deductions**:
- (-5) No integration test between `forge surfaces --json` CLI output and skill consumption. The document lists `go test` for forge surfaces --json output format, but the critical cross-layer integration (skill Bash tool parses JSON output) has no test coverage. This was flagged in iteration 1 and remains unaddressed across all three iterations.

### Coverage Target Numeric (40/45)

Go code 80% stated explicitly.

**Deductions**:
- (-5) No coverage target for Skill layer (still marked N/A). No definition of what constitutes a passing manual test. For 5 surface types × 2 skills = 10 manual test scenarios, the document should specify minimum pass criteria.

### Test Tooling Named (35/40)

`go test` and `just --dry-run` named. Key Test Scenarios provide concrete examples.

**Deductions**:
- (-5) No test tooling specified for JSON output validation. The `forge surfaces --json` output needs deterministic testing (e.g., golden file comparison, `jq` assertions). No tooling is named for this. This was flagged in iteration 1 and remains unchanged.

**Subtotal: 115/130**

---

## Dimension 5: Breakdown-Readiness (180 pts)

### Components Enumerable (60/65)

Phase Component Map now covers all three phases with file-level tables. Phase 1 has 7 rows, Phase 2 has 5 rows, Phase 3 has 12 rows, each with change descriptions and Interface/Model cross-references.

**Deductions**:
- (-5) "skills/*/templates/*.md（共 16 个 prompt 模板）" is a glob pattern, not an enumeration. An implementer cannot derive the exact file list from this. This was flagged in iteration 2 and remains unresolved.

### Tasks Derivable (55/65)

Phase ordering provides task derivation skeleton with file-level granularity across all three phases.

**Deductions**:
- (-5) Phase 2 items like "quality-gate 相关逻辑 | fix-task 流程中调用 forge surfaces <path> --json 推断 surface-key" are still somewhat vague — which specific file(s) in quality-gate need modification? The quality-gate system spans multiple files (`quality_gate.go`, skill docs, rules).
- (-5) No task-level acceptance criteria. Each phase item has a change description but no "done when" definition. An implementer cannot determine when a task is complete without cross-referencing multiple sections.

### PRD AC Coverage (45/50)

PRD Coverage Map maps each PRD AC to a design component with Interface/Model reference.

**Deductions**:
- (-5) Story 4 AC: "旧任务文件含 `scope: frontend` 时，run-tests 能正确读取并按默认编排策略执行" — the design says the PRD AC will be updated. The Migration Notes section provides a complete rationale and blocking migration check. However, the PRD itself has NOT been updated. The design explicitly acknowledges the PRD AC will change but defers the actual change to "需在实现前同步更新 PRD". This is a deferred resolution, not a resolved one.

**Subtotal: 160/180**

---

## Dimension 6: Security Considerations (80 pts)

### Assessment

PRD has no auth/data/multi-user requirements. Marking N/A.

One security-adjacent concern noted (not scored):
- Surface-key naming constraint `[a-zA-Z0-9_-]` is enforced by "init-justfile 配方生成时强制执行" — i.e., by the LLM agent following skill instructions. No Go-layer validation exists in `Task.SurfaceKey` or `forge surfaces` CLI output. If a user manually edits frontmatter with a malicious surface-key containing shell metacharacters, there is no defense.

**Subtotal: 80/80**

---

## Dimension 7: Implementation Feasibility (140 pts)

### Dependencies Available (45/50)

Dependencies section correctly marks `MatchSurface()` as "需修改签名" and `surfaces.go` as "需新增 --json flag".

**Deductions**:
- (-5) The `--json` flag error output specifies `json.NewEncoder(cmd.ErrOrStderr()).Encode()` for JSON errors on stderr. However, the success paths in the three modes use different JSON structures (array for query, object for list/types). The implementation must handle serialization for each mode separately, and the design does not specify whether a shared response struct or per-mode inline serialization is preferred.

### Architecture Fits Project Structure (45/50)

**Deductions**:
- (-5) The Task struct change deletes `Scope` and adds `SurfaceKey`/`SurfaceType`. The migration is now blocking (exit 2 when old `scope` detected), which is a significant improvement. However, the blocking check only fires when `build.go` loads `index.json`. Other code paths that read tasks (e.g., `forge task list`, `forge task show`) may not go through `build.go` — will they also detect old-format tasks? The migration check scope is not fully specified.
- (-5) The `renderTemplate()` change from `{{SCOPE}}` to `{{SURFACE_KEY}}` and `{{TEST_TYPE}}` to `{{TEST_TYPE_ARG}}` is listed in Phase 3 but the template engine mechanism is not described. Is this a string replacement? A Go `text/template` variable? A custom tokenizer? The implementation pattern affects the scope of change.

### Technical Claims Grounded (35/40)

**Deductions**:
- (-5) The migration blocking check claims to solve the backward compatibility issue. But consider: a team upgrades Forge CLI, and CI runs `forge run-tests` before anyone has regenerated tasks. `run-tests` reads task frontmatter, loads `build.go`, which returns exit 2. The CI pipeline is now broken until someone manually runs `breakdown-tasks`. This is a disruptive migration experience. The claim "此检查确保升级后不存在数据不一致的任务" is technically true but the user experience impact is not acknowledged.

**Subtotal: 125/140**

---

## Cross-Cutting Deductions

### Cross-Section Inconsistency (-0)

No cross-section inconsistencies found. The major inconsistency from iteration 2 (PRD AC conflict) is now acknowledged with a documented mitigation strategy (blocking migration + PRD update plan). The SurfaceMatch struct package omission is a gap but not an inconsistency.

### Vague Language (-20)

1. "Skill 层错误通过 SKILL.md HARD-GATE 规则约束 agent 行为" — this is a documentation constraint, not a deterministic mechanism. The design provides no way to verify that an LLM agent follows these rules. (-20, persistent from iterations 1-2 — the wording is identical and the fundamental limitation is not acknowledged)

---

## Blindspot Attacks

### Attack 1 [Interface & Model]: SurfaceMatch package placement unspecified

**Quote**: Interface 1a defines `type SurfaceMatch struct { Key string; Type string }` without specifying the package.

**Analysis**: `MatchSurface()` lives in `pkg/forgeconfig/match.go`. If `SurfaceMatch` is also in `pkg/forgeconfig`, then `surfaces.go` (in `internal/cmd/`) must import `pkg/forgeconfig` to access the type — which it already does. But `build.go` and `autogen.go` (in the `task` package) may also need to reference `SurfaceMatch` if they call `MatchSurface()` directly. If they don't, then who converts the CLI's JSON output into task fields? The data flow from `SurfaceMatch` to `Task.SurfaceKey`/`Task.SurfaceType` has a gap in the middle — the skill layer writes frontmatter, `build.go` reads frontmatter, but there is no Go-layer function that produces a `SurfaceMatch` and directly populates a `Task`. The entire surface-key resolution happens in the skill layer via CLI call. The `SurfaceMatch` type exists only for `surfaces.go` and `match.go` — but this is never stated, leaving the reader to wonder about the type's scope. **What must improve**: Specify the package for `SurfaceMatch` and clarify which Go packages consume it.

### Attack 2 [Breakdown-Readiness]: 16 prompt templates remain a glob

**Quote**: Phase 3 lists "skills/*/templates/*.md（共 16 个 prompt 模板）" as a single row with change description "`{{SCOPE}}` → `{{SURFACE_KEY}}`，`{{TEST_TYPE}}` → `{{TEST_TYPE_ARG}}` 变量值域同步".

**Analysis**: This has been flagged in iterations 1 and 2 without resolution. Sixteen files need modification but the design treats them as a single glob. The change description says two template variables are being renamed, but does not specify: (a) which templates use `{{SCOPE}}` vs which use `{{TEST_TYPE}}` — they may not all use both; (b) whether some templates have conditional logic around these variables; (c) whether the replacement is a simple find-replace or requires structural changes. For an implementer, this single row could represent anywhere from 16 trivial changes to 16 non-trivial refactorings. **What must improve**: Enumerate the 16 template files and specify which variable(s) each uses.

### Attack 3 [Error Handling]: Migration blocking check scope is narrow

**Quote**: Migration Notes states "build.go 加载 index.json 时执行阻塞式迁移检查：扫描含 scope 字段但不含 surface-key 的任务".

**Analysis**: The migration check is located in `build.go` during `index.json` loading. But `build.go` is only one consumer of task data. Other code paths that access tasks — `forge task list`, `forge task show`, `forge task add` — may read task data through different code paths that don't go through `build.go`'s frontmatter parsing. If these paths don't perform the same migration check, they will encounter tasks with empty `SurfaceKey`/`SurfaceType` and silently proceed, potentially causing incorrect behavior downstream (e.g., `run-tests` selecting the wrong execution strategy). The migration check is correctly placed for the primary flow but its scope is not comprehensive. **What must improve**: Either (a) extract the migration check into a shared function called by all task-loading code paths, or (b) document that the migration check only covers the `build.go` frontmatter parsing path and identify which other paths exist.

### Attack 4 [Implementation Feasibility]: CI-breaking migration experience

**Quote**: Migration Notes states "build.go 检测到旧 scope 字段时返回 exit 2，用户必须重跑 breakdown-tasks/quick-tasks 后才能继续操作".

**Analysis**: The blocking migration strategy means that any CI pipeline running `forge run-tests` or `forge task list` after a Forge CLI upgrade (but before task regeneration) will fail with exit 2. For teams with multiple developers or automated CI, this creates a coordination problem: the first person to upgrade Forge CLI must also regenerate all tasks before anyone else can proceed. There is no graceful degradation path. Compare with the PRD's "零回归保证" goal: "无 surface 配置的项目输出与当前完全一致". This goal holds for projects without surface configuration, but the design does not acknowledge that projects WITH surface configuration face a mandatory disruption window. **What must improve**: Acknowledge the CI disruption impact in the design document and consider providing a `forge task migrate` command or flag that can be run independently of `breakdown-tasks`/`quick-tasks`.

---

## Summary

```
SCORE: 760/1000
DIMENSIONS:
  Architecture Clarity: 155/170
  Interface & Model Definitions: 145/170
  Error Handling: 120/130
  Testing Strategy: 115/130
  Breakdown-Readiness: 160/180
  Security Considerations: 80/80
  Implementation Feasibility: 125/140
ATTACKS:
1. [Interface & Model]: SurfaceMatch struct package 未指定 — Interface 1a 定义 type SurfaceMatch struct 但未声明所在 package — 明确 SurfaceMatch 所在 package 及消费方
2. [Breakdown-Readiness]: 16 个 prompt 模板仍为 glob 模式 — "skills/*/templates/*.md（共 16 个）" 未列举具体文件名及各模板使用的变量 — 枚举 16 个模板文件并标注各模板涉及的变量替换
3. [Error Handling]: 迁移阻塞检查仅覆盖 build.go — "build.go 加载 index.json 时执行阻塞式迁移检查" 未考虑 forge task list/show/add 等其他 task 读取路径 — 提取迁移检查为共享函数或明确检查覆盖范围
4. [Implementation Feasibility]: 阻塞式迁移导致 CI 中断 — "用户必须重跑 breakdown-tasks/quick-tasks 后才能继续操作" 意味着升级后首次 CI 必定失败 — 提供 forge task migrate 独立命令或承认迁移窗口的 CI 中断影响
```
