# Design Evaluation: Surface-Aware Justfile — Iteration 2

**Evaluator**: Staff Architect (Adversary)
**Date**: 2026-05-25
**Document**: `docs/features/surface-aware-justfile/design/tech-design.md`
**PRD**: `docs/features/surface-aware-justfile/prd/prd-spec.md`
**Previous**: `docs/features/surface-aware-justfile/eval/design/iteration-1.md` (540/1000)

---

## Iteration 1 Issues — Resolution Check

| # | Iteration 1 Issue | Resolution Status |
|---|-------------------|-------------------|
| 1 | MatchSurface() returns only surface-type, not surface-key | **RESOLVED** — Interface 1a now defines `SurfaceMatch{Key, Type}` struct and new signature explicitly. `match.go` added to Phase Component Map. |
| 2 | Scope deletion without migration violates PRD Story 4 AC | **PARTIALLY RESOLVED** — Migration Notes section added with explicit rationale, migration validation step, and PRD AC update acknowledgment. However, PRD itself still contains the contradictory AC. |
| 3 | --json flag understated as parenthetical | **RESOLVED** — Interface 1b provides full three-mode specification with concrete JSON examples per mode, including error output format. |
| 4 | Silent data loss on migration | **PARTIALLY RESOLVED** — Changed from "静默赋空值" to "记录 warning 日志，赋空值". Migration validation step added. But "赋空值" still means information loss with no automated remediation. |
| 5 | match.go missing from component list | **RESOLVED** — Phase Component Map now explicitly lists `pkg/forgeconfig/match.go` as first row in Phase 1. |

---

## Dimension 1: Architecture Clarity (170 pts)

### Layer Placement (50/60)

Design states: "Forge CLI 层（Go 数据模型 + prompt 合成）+ Forge Plugin 层（Skill 文档 + 规则文件）". Phase Component Map further clarifies which files belong to which phase.

**Deductions**:
- (-5) Layer placement remains a single bullet. There is no description of the boundary contract between CLI and Plugin layers — what crosses the boundary, what doesn't, and what happens when the boundary is violated. For example, `forge surfaces --json` is the sole CLI→Plugin bridge, but this architectural significance is not called out.
- (-5) No justification for why the boundary falls where it does. A reader unfamiliar with Forge's distribution model (skill docs are consumed by LLM agents, Go code by the CLI binary) would not understand the architectural motivation.

### Component Diagram (50/60)

Component diagram is present with data flow arrows, showing CLI → skills → Task struct → prompt.go flow.

**Deductions**:
- (-5) The diagram shows `MatchSurface(path)` with annotation "(需修改: 返回 SurfaceMatch{Key,Type})" which is correct. However, the diagram does not show the migration validation component mentioned in Migration Notes — "扫描 index.json 中含 scope 字段但不含 surface-key 的任务". This is a new component that should appear in the diagram.
- (-5) The diagram shows `prompt.go renderTemplate()` using `{{SURFACE_KEY}} direct` and `{{TEST_TYPE_ARG}} direct` as if they are two separate template variables, but the document does not clarify whether `TEST_TYPE_ARG` is being renamed or if this is a new variable in addition to `SURFACE_KEY`. The relationship between these two template variables is architecturally ambiguous.

### Dependencies (45/50)

Dependencies section now correctly distinguishes "已存在，不变" from "需修改签名" and "需新增".

**Deductions**:
- (-5) The document claims "无新外部依赖" but `just >= 1.4.0` is listed as an external dependency. While this is technically an existing version requirement, the statement "无新外部依赖" is slightly misleading — a reader might interpret it as "no version changes needed" rather than "no new packages". This is a minor clarity issue.

**Subtotal: 145/170**

---

## Dimension 2: Interface & Model Definitions (170 pts)

Using "no er-diagram" scoring criteria (db-schema: "no").

### Interface Signatures Typed (50/60)

Interface 1a provides exact Go function signatures (old and new) with typed struct. Interface 1b provides cobra flag registration code and three-mode JSON output specification.

**Deductions**:
- (-5) Interface 1a defines `SurfaceMatch` struct but does not specify its package. Is it `forgeconfig.SurfaceMatch`? A new shared type? The placement matters for import paths and downstream consumers.
- (-5) Interface 2 (Surface 规则文件格式) remains as comment-level pseudo-code: "编排序列 // 步骤名 | 退出码 | 语义 | 后续动作". This is a description of section headings, not a typed interface. No concrete example file is provided. An implementer cannot determine the exact markdown structure, field separators, or parsing expectations from this specification.

### Models Concrete (45/60)

Five Go structs defined with field-level detail (SurfaceKey, SurfaceType on each).

**Deductions**:
- (-5) Model 1 (Task) specifies JSON tags `"surface-key,omitempty"` / `"surface-type,omitempty"` but does not specify what happens to existing `index.json` files that contain `"scope": "frontend"`. The Migration Notes section says this is addressed via a validation scan, but the Model definition itself does not document the deserialization behavior when `scope` is present — will Go's `json.Unmarshal` silently ignore unknown fields? (Yes, it will by default, but this behavior should be explicitly stated.)
- (-5) Model 2 (AutoGenTaskDef) says "原 Scope 字段删除后新增" for SurfaceKey, but the existing `AutoGenTaskDef` has `Scope string` with `yaml:"scope"` and `json:"scope,omitempty"` tags. The design does not specify the YAML/JSON tags for the new fields. Are they `yaml:"surface-key"` and `json:"surface-key,omitempty"`? This is a concrete typing gap.
- (-5) Model 5 (TaskState) shows SurfaceKey/SurfaceType but does not document how these fields are populated. Current code populates TaskState during task claim from the Task struct. The design does not state whether the claim path needs modification or whether the fields propagate automatically.

### Directly Implementable (35/50)

**Deductions**:
- (-5) Interface 1a states "所有现有调用方（runSurfacesQuery 等）需适配新签名，从 SurfaceMatch.Type 取值" but does not enumerate which callers exist. The implementer must search the codebase to find all call sites.
- (-5) The `--json` flag implementation pattern says "在 runSurfacesList/runSurfacesQuery/runSurfacesTypes 中增加 if jsonFlag 分支" but does not define error handling within JSON mode. What happens if `ReadSurfaces()` fails in JSON mode? The error output is specified (stderr JSON error + exit 1) but the code path from ReadSurfaces failure to this error output is not specified.
- (-5) Surface 规则文件格式 (Interface 2) has zero concrete content. There is no example rule file for any surface type. An implementer must invent the markdown structure, field format, and parsing rules from scratch. This is the primary interface for 10 files being created but it is the least specified.

**Subtotal: 130/170**

---

## Dimension 3: Error Handling (130 pts)

### Error Types Defined (35/45)

Go layer error table defines 4 scenarios with exit codes. Skill layer references PRD's 7 error scenarios.

**Deductions**:
- (-5) No Go error types or sentinel errors defined. The document specifies exit codes and stderr messages but does not define concrete Go error types. For example, what Go error wraps the "surfaces 配置缺失" condition? Is it a `fmt.Errorf("...")` or a typed error? The existing codebase uses `surfacesPathError` in `surfaces.go` — no equivalent convention is specified for new error paths.
- (-5) The error table says "build.go 解析旧 frontmatter 含 scope 但无 surface-key" results in "记录 warning 日志（提示迁移），赋空值" with exit 0. But the warning message format is not specified. What logging library? What message format? Following TECH-error-handling-001, errors must go to stderr with `<context>: <specific-detail>` format, but "warning 日志" implies a different channel (possibly structured logging, not stderr).

### Propagation Strategy Clear (35/45)

**Deductions**:
- (-5) The propagation strategy section says "Go 层错误通过 exit code + stderr 传播给 skill（LLM agent 读取）" and "Skill 层错误通过 just 命令 exit code 和 SKILL.md HARD-GATE 规则控制". The Go→Skill propagation is well-defined (exit code + stderr). But the Skill layer propagation — "SKILL.md HARD-GATE 规则控制" — is still a documentation constraint on an LLM agent, not a deterministic error propagation mechanism. The design conflates prompt engineering constraints (HARD-GATE rules in SKILL.md that instruct an LLM how to behave) with architectural error handling. There is no verification mechanism that an LLM agent actually follows these rules.
- (-5) The `--json` error output format specifies `stderr: {"error": "no surface configured; run 'forge init' to configure surfaces"}` but this mixes JSON error format on stderr with the convention that stderr should contain `<context>: <specific-detail>` format (per injected error handling conventions). Is JSON-on-stderr acceptable when the caller expects structured output? The document does not reconcile these two conventions.

### HTTP Status Codes (40/40)

N/A — no HTTP surface. Full score.

**Subtotal: 110/130**

---

## Dimension 4: Testing Strategy (130 pts)

### Per-Layer Test Plan (35/45)

Test table lists 5 Go unit test rows and 3 Skill test rows. Key Test Scenarios section adds 6 concrete scenarios.

**Deductions**:
- (-5) No integration test between `forge surfaces --json` CLI output and skill consumption. The document lists `go test` for forge surfaces --json output format, but the critical cross-layer integration (skill Bash tool parses JSON output) has no test coverage. This was flagged in iteration 1 and remains unaddressed.
- (-5) No test plan for the migration validation step described in Migration Notes. The document specifies "扫描 index.json 中含 scope 字段但不含 surface-key 的任务" as a validation check, but no test row covers this scenario.

### Coverage Target Numeric (40/45)

Go code 80% stated explicitly.

**Deductions**:
- (-5) No coverage target for Skill layer (still marked N/A). No definition of what constitutes a passing manual test. For 5 surface types × 2 skills = 10 manual test scenarios, the document should specify minimum pass criteria.

### Test Tooling Named (35/40)

`go test` and `just --dry-run` named. Key Test Scenarios provide concrete examples.

**Deductions**:
- (-5) No test tooling specified for JSON output validation. The `forge surfaces --json` output needs deterministic testing (e.g., golden file comparison, `jq` assertions). No tooling is named for this.

**Subtotal: 110/130**

---

## Dimension 5: Breakdown-Readiness (180 pts)

### Components Enumerable (55/65)

Phase Component Map explicitly lists 7 files in Phase 1 with change descriptions and Interface/Model cross-references. Phases 2 and 3 list component groups.

**Deductions**:
- (-5) Phase 2 and Phase 3 components are listed as group names ("breakdown-tasks/quick-tasks skill", "16 个 prompt 模板") without the file-level enumeration that Phase 1 provides. The Phase 1 table is a model of clarity; the absence of equivalent detail in Phases 2-3 is a step down.
- (-5) "16 个 prompt 模板 SURFACE_KEY 变量值域同步" is listed but no enumeration of which 16 templates or what the variable change looks like. This was flagged in iteration 1 and remains unchanged.

### Tasks Derivable (45/65)

Phase ordering provides task derivation skeleton. Phase Component Map adds file-level granularity for Phase 1.

**Deductions**:
- (-10) Phase 2 and Phase 3 items remain single-sentence descriptions. "breakdown-tasks/quick-tasks：任务生成时调用 forge surfaces CLI 填充双字段" — what are the subtasks? Parse CLI output? Call CLI? Update frontmatter generation? Error handling for CLI failure? An implementer cannot derive a complete task list from this.
- (-10) No task-level acceptance criteria for any phase. The "done when" definition is absent. For Phase 1, the Interface/Model cross-references help, but there are no concrete pass/fail criteria per task.

### PRD AC Coverage (45/50)

PRD Coverage Map maps each PRD AC to a design component with Interface/Model reference.

**Deductions**:
- (-5) Story 4 AC: "旧任务文件含 `scope: frontend` 时，run-tests 能正确读取并按默认编排策略执行" — the design says "不保留兼容层，PRD AC 将更新（见 Migration Notes）". The Migration Notes section provides rationale and a migration validation step. However, the PRD itself has NOT been updated. The design acknowledges a PRD AC violation but does not resolve it — the resolution is deferred to "需在实现前同步更新 PRD". This means the current state has an unresolved PRD AC gap.

**Subtotal: 145/180**

---

## Dimension 6: Security Considerations (80 pts)

### Assessment

PRD has no auth/data/multi-user requirements. Marking N/A.

One security-adjacent concern noted (not scored):
- Surface-key naming constraint `[a-zA-Z0-9_-]` is enforced by "init-justfile 配方生成时强制执行" — i.e., by the LLM agent following skill instructions. No Go-layer validation exists in `Task.SurfaceKey` or `forge surfaces` CLI output. If a user manually edits frontmatter with a malicious surface-key containing shell metacharacters, there is no defense.

**Subtotal: 80/80**

---

## Dimension 7: Implementation Feasibility (140 pts)

### Dependencies Available (40/50)

Dependencies section now correctly marks `MatchSurface()` as "需修改签名" rather than "已存在". `surfaces.go` listed as needing `--json` addition.

**Deductions**:
- (-5) The document says "无新外部依赖" but depends on `just >= 1.4.0` which has `[linux]`/`[windows]` recipe attributes. If existing projects use an older just version, this is effectively a new version requirement. The document lists this under Dependencies but does not reconcile the "无新外部依赖" statement.
- (-5) The `--json` flag implementation says "使用 json.NewEncoder(cmd.OutOrStdout()).Encode()" — this uses Go's standard `encoding/json` package, which is available. However, the error output format `{"error": "..."}` on stderr uses a different serialization path than the success JSON on stdout. The design does not specify whether the error JSON uses the same encoder or raw string formatting, which affects error output consistency.

### Architecture Fits Project Structure (40/50)

**Deductions**:
- (-5) The Task struct change deletes `Scope` field (JSON tag `"scope,omitempty"`) and adds `SurfaceKey`/`SurfaceType`. The design acknowledges this is a breaking change but says the migration window is controlled because "升级后首次 breakdown-tasks/quick-tasks 运行即自动填充新字段". However, this assumes all tasks are regenerated via breakdown-tasks/quick-tasks — what about manually created tasks or tasks from other sources? The migration plan does not cover all task origins.
- (-5) The `renderTemplate()` change from `{{SCOPE}}` to `{{SURFACE_KEY}}` requires updating all 16 prompt templates. The design lists this as a Phase 3 item but does not specify the template engine or how the variable injection mechanism works, making it difficult to assess whether this is a simple find-replace or a more complex template restructuring.

### Technical Claims Grounded (35/40)

**Deductions**:
- (-5) The Migration Notes state "迁移窗口可控：升级后首次 breakdown-tasks/quick-tasks 运行即自动填充新字段". This claim assumes that (a) users will run breakdown-tasks/quick-tasks after upgrade, (b) all tasks will be regenerated, and (c) no tasks persist from other creation paths. None of these assumptions are validated. A project that upgrades but does not re-run task generation will have tasks with empty SurfaceKey/SurfaceType — the migration is not automatic, it requires user action that is not enforced.

**Subtotal: 115/140**

---

## Cross-Cutting Deductions

### Cross-Section Inconsistency (-30)

1. PRD Coverage Map row "Story4: 旧任务 scope 兼容" says "不保留兼容层，PRD AC 将更新（见 Migration Notes）" — but the PRD's actual AC text still reads "旧任务文件含 scope: frontend 时，run-tests 能正确读取并按默认编排策略执行". The design document and PRD are in conflict. The design defers the resolution ("需在实现前同步更新 PRD") but does not resolve it. (-30)

### Vague Language (-20)

1. "Skill 层错误通过 SKILL.md HARD-GATE 规则约束 agent 行为" — this is a documentation constraint, not a deterministic mechanism. The design provides no way to verify that an LLM agent follows these rules. (-20, carried from iteration 1 — the wording is identical and the issue is not addressed)

---

## Blindspot Attacks

### Attack 1 [Interface & Model]: Surface 规则文件格式无法实施

**Quote**: Interface 2 defines the rule file format as three section headings with comment-level descriptions: "编排序列 // 步骤名 | 退出码 | 语义 | 后续动作" and "配方调用契约 // 配方名 | 参数签名 | 退出码语义".

**Analysis**: This is the primary interface for 10 new files being created (5 surface rules per skill × 2 skills). Despite being central to the entire feature — it defines how init-justfile generates recipes and how run-tests orchestrates execution — the specification provides no concrete example, no markdown structure, no parsing rules, and no validation criteria. An implementer must invent the entire file format from a three-line comment template. There is no way to verify that two different implementers would produce compatible rule files. **What must improve**: Provide at least one complete example rule file (e.g., `rules/surfaces/web.md`) with concrete content showing the actual markdown structure, step definitions, and recipe contracts.

### Attack 2 [Error Handling]: JSON-on-stderr conflicts with error handling conventions

**Quote**: Interface 1b specifies surfaces configuration missing error as `stderr: {"error": "no surface configured; run 'forge init' to configure surfaces"}` with exit 1.

**Analysis**: The injected error handling conventions (TECH-error-handling-001) require error messages to follow `<context>: <specific-detail>` format on stderr. The `--json` error format outputs `{"error": "..."}` on stderr — this is neither the `<context>: <specific-detail>` plain text format nor is it documented as an exception to the convention. When `--json` mode produces errors, should stderr be plain text (per convention) or JSON (per Interface 1b)? The document does not reconcile this conflict. A consumer parsing stderr in `--json` mode might expect JSON but get plain text, or vice versa if the error originates from a layer that doesn't know about `--json` mode. **What must improve**: Explicitly state that `--json` mode overrides the stderr format convention with structured JSON errors, and ensure all error paths in `--json` mode produce consistent JSON output.

### Attack 3 [Breakdown-Readiness]: Phase 2-3 task derivation is underspecified

**Quote**: Phase Component Map provides detailed file-level breakdown for Phase 1 (7 rows with change descriptions and cross-references) but Phase 2 lists "breakdown-tasks/quick-tasks skill、forge task add CLI、quality-gate fix-task" and Phase 3 lists "init-justfile skill、run-tests skill、16 个 prompt 模板、死代码清理" — each as bare names.

**Analysis**: Phase 1's table format is excellent — it provides file paths, change descriptions, and Interface/Model cross-references. Phases 2 and 3 provide none of this. For example, "init-justfile skill" in Phase 3 involves: (a) rewriting SKILL.md with surface detection flow, (b) creating 5 rule files, (c) handling mixed-project aggregation logic, (d) implementing `# user-customized` protection. These are at least 4 distinct implementation tasks hidden behind one name. An implementer cannot derive accurate task estimates or dependency ordering from the current Phase 2-3 specification. **What must improve**: Apply the same Phase 1 table format to Phases 2 and 3, listing specific files, change descriptions, and cross-references.

### Attack 4 [Implementation Feasibility]: Migration validation is not automatic

**Quote**: Migration Notes states "迁移窗口可控：升级后首次 breakdown-tasks/quick-tasks 运行即自动填充新字段" and "此检查不阻塞正常流程".

**Analysis**: The migration relies on users re-running task generation after upgrade. The validation step ("扫描 index.json 中含 scope 字段但不含 surface-key 的任务") is a warning-only, non-blocking check. This means: (1) Projects that upgrade and don't run task generation will have tasks with empty SurfaceKey/SurfaceType — run-tests will fail to load execution strategy rules for these tasks. (2) The warning is passive — there is no enforcement mechanism. (3) The "零回归" PRD goal ("无 surface 配置的项目输出与当前完全一致") may hold, but projects WITH surface configuration that upgrade will experience regression until they re-run task generation. **What must improve**: Either (a) make the validation check blocking (exit 2 if migration-needed tasks detected), or (b) provide an automated migration command, or (c) acknowledge in the design that post-upgrade regression is expected for surface-configured projects and update the PRD "零回归" goal accordingly.

---

## Summary

```
SCORE: 650/1000
DIMENSIONS:
  Architecture Clarity: 145/170
  Interface & Model Definitions: 130/170
  Error Handling: 110/130
  Testing Strategy: 110/130
  Breakdown-Readiness: 145/180
  Security Considerations: 80/80
  Implementation Feasibility: 115/140
ATTACKS:
1. [Interface & Model]: Surface 规则文件格式仅用注释级伪代码描述 — "编排序列 // 步骤名 | 退出码 | 语义 | 后续动作" — 10 个核心文件的可实施接口未提供具体示例 — 提供至少一个完整的 web.md 规则文件示例
2. [Error Handling]: --json 模式 stderr JSON 输出与 TECH-error-handling-001 的 <context>: <specific-detail> 格式冲突 — "stderr: {"error": "..."}" 未声明为格式例外 — 明确 --json 模式下 stderr 格式规范
3. [Breakdown-Readiness]: Phase 2-3 缺少 Phase 1 级别的文件级拆解 — "init-justfile skill" 隐藏至少 4 个独立实现任务 — 对 Phase 2-3 应用同样的表格级拆解
4. [Implementation Feasibility]: 迁移依赖用户手动重跑任务生成 — "迁移窗口可控：升级后首次 breakdown-tasks/quick-tasks 运行即自动填充" 假设用户行为但无强制机制 — 提供自动迁移或阻塞检查
```
