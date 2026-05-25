# PRD Evaluation Report — Iteration 1

**Feature**: Surface-Aware Justfile
**Evaluator Persona**: Senior Product Manager
**Date**: 2026-05-25
**Documents Reviewed**: `prd/prd-spec.md`, `prd/prd-user-stories.md`

---

## Phase 1: Reasoning Audit (Pre-Score Anchors)

### Anchor 1: Problem → Solution alignment
The stated problem (init-justfile ignores surface types, leading to incorrect test orchestration) is directly addressed by the solution (surface-aware rule files + dispatcher mode). This is a strong alignment. However, the PRD bundles four distinct changes (surface-aware generation, test.execution removal, surface-key migration, Task model extension) into one feature. The coupling is justified by shared dependencies but creates risk: any one sub-feature's blocker delays all four.

### Anchor 2: Unsubstantiated evidence
Quote: *"75% 的实际示例已通过 just 命令调用"* — this statistic has no source. Is it from codebase analysis? User surveys? Internal testing? As a PM, I cannot assess whether removing test.execution is safe based on an unattributed percentage.

### Anchor 3: Success criteria measure proxies
The "编排链路从 4 层降至 2 层" metric measures architectural simplification, not user outcome. The user outcome should be "developers spend less time configuring test execution" or "zero manual test.execution configuration needed". The current metric could be achieved without improving the developer experience.

### Anchor 4: Scope boundary ambiguity
The Out of Scope list excludes "变更 forge-cli/pkg/just/ 门控序列" and "变更 forge-cli/internal/cmd/quality_gate.go 或 testrunner 的 Go 代码", but Related Changes #9 describes quality-gate fix-task modifications and #3-5 describe Go code changes to prompt.go, task/types.go, and autogen.go. The boundary between "which Go code changes" and "which Go code is excluded" is drawn per-file rather than per-concern, creating confusion about what constitutes a boundary violation.

---

## Phase 2: Rubric Scoring

### Dimension 1: Background & Goals (100 pts)

**Background has three elements (Reason/Target/Users): 28/30**

All three elements are present and structured clearly:
- Reason: init-justfile ignores surface types, different surfaces need different orchestration
- Target: 4 concrete deliverables listed
- Users: Two user types identified (Forge plugin developers, Forge users/project developers)

Minor deduction: The "Who" section describes users by role but does not specify their skill level or typical workflow context (e.g., "developers who configure .forge/config.yaml and run CLI commands" is implied but not stated).

**Goals are quantified: 25/30**

Quantified metrics present:
- "编排链路从 4 层降至 2 层" (layer count)
- "覆盖 5 种 surface 类型" (count)
- "7+ 组件的 surface-key 值域从固定枚举迁移为用户自定义" (count)
- "Task 新增 surface-key + surface-type 双字段" (count)

Deduction: "无 surface 配置的项目输出与当前完全一致" is a binary pass/fail criterion but lacks a specific verification method. "diff 输出对比验证" is mentioned in Notes but should be in the Goals metric column. Also, "不超过 1 秒的额外耗时" in Performance Requirements is a numeric target but is not in the Goals table.

**Background and goals are logically consistent: 35/40**

Goals follow from the stated problem in most cases. The reasoning chain holds: surface-aware generation → correct orchestration → no manual configuration → simplified architecture.

Deduction (-5): Goal "Surface-key 值域统一" is a refactoring/migration goal that doesn't directly address the user-facing problem stated in the Background (which is about test orchestration correctness). This is an internal architecture goal presented as a feature goal. The Background should include a reason why the fixed enumeration causes user pain, not just "消除硬编码约束".

Deduction: The "零回归保证" goal is stated but the Background section never identifies regression as a risk or explains why it matters specifically for this change. The connection is implicit.

**Dimension 1 Total: 88/100**

---

### Dimension 2: Flow Diagrams (150 pts)

**Mermaid diagram exists: 50/50**

Yes, a Mermaid flowchart is present in the "Business Flow Diagram" section covering both init-justfile and run-tests flows.

**Main path complete (start → end): 45/50**

The init-justfile flow covers: start → detect language → load template → generate base recipes → detect surface → load surface rules → generate surface recipes → arbitrate conflicts → protect user customizations → assemble justfile → validate → end. Complete.

The run-tests flow covers: start → get surface info → load rule → branch by surface type → execute sequence → teardown → end. Complete.

Minor deduction: The diagram shows `Probe -->|exit 1/2/3| Teardown1 --> Abort1` for web/api, but `Test1 -->|exit 1| Teardown3 --> RunEnd` — the probe failure leads to "Abort" while test failure leads to "RunEnd" (completion). The semantics of this difference are unclear in the diagram alone: does exit 1 from test mean "passed with test failures" or "completed with failures"? According to BIZ-error-reporting-001, exit code 1 means retryable failure. The diagram does not clarify whether test exit 1 triggers a retry or is treated as a final state.

**Decision points + error branches covered: 40/50**

Decision points present:
- `CheckSurf{检测 surface 类型?}` — diamond node with yes/no branches
- `Arbitrate{Surface 规则与语言模板冲突?}` — diamond node
- `Protect{目标配方有 user-customized?}` — diamond node
- `ExecSeq{编排序列类型?}` — diamond node with 3 branches

Error branches present:
- Probe failure → teardown → abort
- Test exit 1 → teardown → complete

Deductions:
- (-5) No error branch for "surface detection fails" or "surface type unknown" — what happens when a surface type is not one of web/api/cli/tui/mobile?
- (-5) The init-justfile flow has no error branch at all. What happens when language detection fails? When surface rules file is missing?

**Dimension 2 Total: 135/150**

---

### Dimension 3: Flow Completeness (200 pts) — Mode B

**Flow steps describe complete business process: 60/70**

The text describes two complete business processes:
1. init-justfile: detect language → load template → generate base → detect surface → load rules → generate surface recipes → arbitrate → protect → assemble → validate. Steps include state transitions.
2. run-tests: get surface info → load rule → orchestrate just recipe sequence → check exit codes → abort/teardown.

The Surface orchestration table documents 5 surface types with their specific sequences.

Deductions:
- (-5) The init-justfile flow step 4 says "跨平台变体验证" but does not specify what happens when validation fails. Is it a warning? An error? Does it block generation?
- (-5) The run-tests flow step 1 says "优先任务文档 frontmatter → forge surfaces CLI" but does not specify what happens when neither source provides surface information. This is a gap in the complete process description.

**Data flow documented (if multi-system): 60/70**

This is a multi-system feature spanning: config.yaml, Go structs (prompt.go, types.go, autogen.go), SKILL.md documents, prompt templates, and just CLI. The Related Changes table (#1-11) effectively serves as a data flow map showing which modules change and how.

Deductions:
- (-10) No explicit data flow table showing data movement between systems. The Related Changes table shows change points but not data flow direction (what data flows from config.yaml → prompt.go → task template → justfile). For a feature involving 7+ components and 16 prompt templates, an explicit data flow diagram or table would significantly reduce implementation ambiguity.

**Exception handling and edge cases covered: 45/60**

Documented exception handling:
- Probe timeout → teardown → abort (web/api)
- Dev server crash → probe timeout → teardown (in Other Notes)
- Session interruption → .forge/test-state.json recovery
- Teardown idempotency (PID not found → skip)

Deductions:
- (-5) No documented behavior for "just recipe not found" — what happens when init-justfile didn't generate a recipe that run-tests expects?
- (-5) No documented retry logic for probe. The diagram shows probe as a single step, but the text mentions "重试轮询". The retry count, interval, and total timeout are not specified in the Flow Description (only in Observability notes showing `[retry 3/30]`).
- (-5) The "exit 1" from test in the flow diagram goes to teardown → complete, but per BIZ-error-reporting-001, exit code 1 means "retryable failure". The document does not specify whether test failures are retried or always treated as final.

**Dimension 3 Total: 165/200**

---

### Dimension 4: User Stories (200 pts)

**Coverage: one story per target user: 45/50**

Background identifies two user types:
- Forge plugin developers — covered by Story 3 and Story 4
- Forge users (project developers) — covered by Story 1 and Story 2

Deduction (-5): No story for the "mixed project" scenario where a project has multiple surfaces. Story 1 mentions "混合项目 dev 配方接受 surface-key 参数" in scope but no user story describes this multi-surface interaction from the user's perspective.

**Format correct (As a / I want / So that): 50/50**

All four stories follow the correct format:
- Story 1: As a Forge 用户 / I want to 配置 surfaces 字段后... / So that web/api 项目获得...
- Story 2: As a Forge 用户 / I want to 运行 run-tests 时... / So that 我不需要手动配置...
- Story 3: As a Forge 插件开发者 / I want to 将 surface-key 值域... / So that 混合项目的所有配方...
- Story 4: As a Forge 插件开发者 / I want to Task 数据模型... / So that 下游 skill...

Actions are concrete ("配置 surfaces 字段后运行 init-justfile", "运行 run-tests 时自动检测", etc.), not vague like "manage" or "handle".

**AC per story (Given/When/Then): 40/50**

All stories have acceptance criteria in Given/When/Then format:
- Story 1: Given/When/Then + two And clauses
- Story 2: Given/When/Then + two And clauses
- Story 3: Given/When/Then + two And clauses
- Story 4: Given/When/Then + three And clauses

Deduction (-10): Story 4's AC mixes implementation verification with behavioral verification. AC "Then 任务包含 surface-key: 'admin-panel' 和 surface-type: 'web'" checks implementation detail (JSON field values) rather than user-visible behavior. Story 4 targets "Forge 插件开发者" — the behavioral AC should describe what the developer can now do with the surface-key/surface-type data, not just that the fields exist.

**AC verifiability & boundary coverage: 35/50**

Verifiability analysis:
- Story 1 AC: Verifiable. Can check generated justfile contents. Edge case covered: CLI/TUI surface not generating `run` recipe.
- Story 2 AC: Partially verifiable. "probe 失败时（exit 1/2/3）执行 teardown 后中止" is testable. However, "HARD-GATE: probe 失败后禁止重试 probe 或重试 dev" — the term "HARD-GATE" is undefined in the PRD. Is this an architectural constraint? A runtime check? How is it enforced?
- Story 3 AC: "prompt.go resolveScope() 基于 surfaces map 集合查询而非 projectType 硬编码" — this is an implementation check, not a behavioral verification.
- Story 4 AC: All ACs are implementation checks (JSON field values, function existence).

Deductions:
- (-5) No AC covers the error case: what happens when init-justfile is run on a project with an unrecognized surface type?
- (-5) No AC covers the migration path: existing projects with scope field — how does GetSurfaceKey() handle the transition?
- (-5) Story 2 does not cover the case where surface info is unavailable (neither in frontmatter nor via CLI).

**Dimension 4 Total: 170/200**

---

### Dimension 5: Scenario Completeness (150 pts)

**End-to-end scenario coverage: 50/60**

Scenarios covered:
1. init-justfile with surface → generates correct recipes (Story 1)
2. run-tests with surface → correct orchestration sequence (Story 2)
3. surface-key migration from fixed to user-defined (Story 3)
4. Task data model extension with new fields (Story 4)

Missing scenarios:
- (-5) Multi-surface project: A project with both `admin-panel: web` and `payment-service: api` — how does run-tests decide which surface's orchestration to use? The Scope section mentions "混合项目 dev 配方接受 surface-key 参数" but no scenario describes the end-to-end flow.
- (-5) No-surface project: What does init-justfile generate when no surfaces are defined? The Goals say "零回归" but no scenario traces this path explicitly.

**Implicit assumptions surfaced: 25/40**

Implicit assumptions detected:
1. **Assumption: just >= 1.4.0 is installed** — mentioned in Performance Requirements but not surfaced as a prerequisite or pre-condition in any user story or scenario.
2. **Assumption: forge surfaces CLI exists** — run-tests flow step 1 references "forge surfaces CLI" but this command's existence is assumed. Related Changes #11 describes creating it, but no scenario assumes this dependency.
3. **Assumption: `.forge/test-state.json` mechanism exists** — Reliability section mentions it but no scenario describes its creation, format, or recovery flow.
4. **Assumption: Surface types are limited to 5** — the document lists 5 types but does not state whether this is an exhaustive set or whether custom surface types are supported. If a user defines `surfaces: {foo: bar}` where `bar` is not web/api/cli/tui/mobile, what happens?
5. **Assumption: tasks have frontmatter** — run-tests flow step 1 says "优先任务文档 frontmatter" but does not assume the case where tasks lack frontmatter.

Deduction (-15): At least 3 of these assumptions are significant enough to affect implementation and are not surfaced in the scenarios or as explicit pre-conditions.

**Business-rules consistency: 40/50**

Cross-checking with injected business rules:
- BIZ-task-lifecycle-001: Task state transitions — the PRD does not describe how surface-key/surface-type interact with task state transitions. This is acceptable since the new fields are metadata, not state.
- BIZ-error-reporting-001: Exit code semantics — the flow diagram uses exit 0/1/2/3 but does not consistently map to the rule's semantics (0=success, 1=retryable, 2=blocking). The diagram shows "exit 1/2/3" for probe failure all going to teardown+abort, which treats all non-zero exits identically. Per the rule, exit 2 (blocking) and exit 1 (retryable) should have different handling.
- BIZ-error-reporting-002: Actionable error messages — the Observability section shows a status message format but does not specify that error messages include recovery hints as required by this rule.
- BIZ-quality-gate-001: Quality gate pipeline — the PRD mentions quality-gate fix-task (#9 in Related Changes) but does not describe how surface-aware changes affect the quality gate pipeline sequence.

Deduction (-10): Exit code handling in the flow diagram and run-tests description does not align with BIZ-error-reporting-001's exit code semantics.

**Dimension 5 Total: 115/150**

---

### Dimension 6: Edge Case Coverage (100 pts)

**Error paths documented: 30/40**

Documented error paths:
- Probe failure → teardown → abort
- Dev server crash → probe timeout → teardown
- Session interruption → .forge/test-state.json recovery
- Test failure → teardown → complete

Missing error paths:
- (-5) Surface rule file not found (e.g., user defines surface type "custom" with no rule file)
- (-5) just command not found or version < 1.4.0

**Boundary conditions covered: 25/35**

Addressed boundary conditions:
- Just version requirement: >= 1.4.0 (stated)
- Cross-platform: Windows/macOS/Linux (stated)
- Probe retry: implied by `[retry 3/30]` in observability (but specific count/interval not documented in flow)
- Surface count: 5 types enumerated

Missing boundary conditions:
- (-5) Maximum number of surfaces per project (can a project have 20 surfaces?)
- (-5) Surface key naming constraints (characters, length, uniqueness rules)

**Failure recovery described: 20/25**

Recovery mechanisms described:
- `.forge/test-state.json` for session interruption recovery
- Teardown idempotency
- git revert as rollback method

Deduction (-5): The test-state.json recovery mechanism is mentioned but not described in detail. What does recovery look like? Does the user run a command? Does it happen automatically?

**Dimension 6 Total: 75/100**

---

### Dimension 7: Scope Clarity (100 pts)

**In-scope items are concrete deliverables: 30/35**

In-scope items are specific and include:
- 5 surface rule files (with paths)
- SKILL.md changes (with specific changes)
- Config schema changes
- Go code changes (with specific functions/files)
- 16 prompt template updates

Most items are concrete deliverables with file paths.

Deduction (-5): "16 个 prompt 模板 SURFACE_KEY 变量值域同步" — which 16 templates? No list or reference is provided. This is a concrete count attached to a vague scope item.

**Out-of-scope explicitly lists deferred items: 25/30**

Out-of-scope items listed:
- 语言模板变更
- forge-cli/pkg/just/ gate sequence
- quality_gate.go or testrunner Go code
- New forge CLI commands
- Go subprocess management (long-term)
- Feature flag rollback infrastructure

Deduction (-5): "变更 forge-cli/internal/cmd/quality_gate.go 或 testrunner 的 Go 代码" is out of scope, but Related Changes #9 describes "quality-gate fix-task: 从失败文件路径推断 surface-key/type". If quality_gate.go is out of scope, where does this change live? The scope boundary is contradictory.

**Scope consistent with functional specs and user stories: 25/35**

Cross-referencing:
- User Stories 1-4 map to In-Scope items 1-2 (init-justfile + run-tests), but Stories 3-4 (surface-key migration + Task model) map to In-Scope items that are Go code changes. These are present in both.
- Related Changes table has 11 items. In-Scope has ~16 checkbox items. The mapping is roughly consistent but:
  - (-5) Related Changes #6 (config-schema.md: 移除 test.execution 文档) is in In-Scope but has no corresponding user story. Who does this affect? What's the user impact?
  - (-5) Related Changes #10-11 (16 prompt templates, surface-key-assignment rule) are in In-Scope but have no corresponding user story. These affect Story 3's behavior but are not explicitly tested by any AC.

**Dimension 7 Total: 80/100**

---

## Phase 3: Blindspot Hunt

### [blindspot] Attack 1: Undefined behavior for unrecognized surface types

Quote: *"覆盖 5 种 surface 类型（web/api/cli/tui/mobile）"*

The document enumerates 5 surface types and all flows, stories, and scenarios assume only these types. But the surface-key migration goal explicitly moves toward "用户自定义 surface-key". If a user defines `surfaces: {my-app: desktop}` where `desktop` is not one of the 5 types, the init-justfile flow has no branch for this case, run-tests has no rule file to load, and no error behavior is defined. The gap between "user-defined surface-key" and "5 fixed surface types" is a fundamental contradiction that no rubric dimension captures.

### [blindspot] Attack 2: Probe retry parameters are fragmented across sections

The flow description mentions probe as "just probe 重试轮询", the diagram shows a single probe step, and the Observability section shows `[retry 3/30]`. But the actual retry count, interval, timeout, and backoff strategy are never specified in any single location. This is not just a documentation issue — it means the acceptance criteria for Story 2 cannot be objectively verified because "probe 失败时（exit 1/2/3）执行 teardown 后中止" does not define how many retries occur before "failure" is declared.

### [blindspot] Attack 3: run-tests dispatcher mode has no fallback strategy

Quote: *"检测 surface type 后加载对应执行策略规则"*

The dispatcher pattern assumes surface type is always determinable. But the flow step 1 says "优先任务文档 frontmatter → forge surfaces CLI". What happens when: (a) the task has no frontmatter AND surfaces CLI returns nothing? (b) the task frontmatter says `surface-type: web` but the project has no web surface defined? These are not error paths in the diagram and not covered by any edge case. A product that "just works" needs to define its "just fails" behavior too.

### [blindspot] Attack 4: Coupling risk between four sub-features

The PRD packages four distinct changes: (1) surface-aware recipe generation, (2) test.execution removal, (3) surface-key migration, (4) Task model extension. Each has independent failure modes but they share dependencies (surface type detection, surface-key values). If sub-feature 3 (migration) encounters a blocking issue, sub-features 1 and 2 cannot ship because they depend on the new surface-key semantics. No risk mitigation or phased delivery strategy is documented. This is a PM-level concern about delivery risk that no rubric dimension captures.

---

## Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Background & Goals | 88 | 100 |
| Flow Diagrams | 135 | 150 |
| Flow Completeness | 165 | 200 |
| User Stories | 170 | 200 |
| Scenario Completeness | 115 | 150 |
| Edge Case Coverage | 75 | 100 |
| Scope Clarity | 80 | 100 |
| **Total** | **828** | **1000** |

### Key Weaknesses to Address

1. **Unrecognized surface types**: Define behavior when surface type is not one of the 5 known types
2. **Probe retry parameters**: Consolidate retry count, interval, timeout into one spec location
3. **Dispatcher fallback**: Define behavior when surface info is unavailable
4. **Exit code alignment**: Map exit codes to BIZ-error-reporting-001 semantics consistently
5. **16 prompt templates**: List which templates are affected, don't leave as a vague count
6. **Scope boundary contradiction**: Clarify how quality-gate fix-task changes are in-scope when quality_gate.go is out-of-scope
7. **Story AC verifiability**: Replace implementation checks with behavioral verifications, especially in Stories 3-4
8. **Multi-surface scenario**: Add end-to-end scenario for projects with multiple surfaces
9. **No-surface scenario**: Explicitly trace the zero-regression path as a scenario
10. **Surface assumptions**: Surface the assumptions about just version, forge surfaces CLI existence, and test-state.json as explicit pre-conditions
