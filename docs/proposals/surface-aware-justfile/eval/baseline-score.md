---
iteration: 1
scorer: CTO
date: 2026-05-25
---

# Baseline Score Report (Iteration 1)

## Overall Score: 810/1000

---

## Dimension Breakdown

### 1. Problem Definition: 90/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 35/40 | Two problems clearly articulated: surface-agnostic init-justfile and redundant test.execution delegation layer. Minor deduction: Problem 2's framing as "redundancy" assumes justfile is sufficient without fully proving that the non-just paths (go test, npx vitest, make test) are trivially wrappable. Quote: "绝大多数模板变量最终都解析为 just 命令" — "绝大多数" is not "全部", leaving ambiguity about the real cost. |
| Evidence provided | 35/40 | Concrete evidence cited: "Web UI 的 e2e 测试必须先启动应用，但当前配方没有 surface 特定的启动逻辑"; config-schema.md examples enumerated. The evidence is observational, not quantitative — no data on how many projects use non-just paths, no user feedback quoted verbatim. |
| Urgency justified | 20/30 | Quote: "随着 v3.0.0 test profile 的引入，surface 成为测试流程的核心维度。" This justifies timing but not urgency — no cost-of-delay analysis, no explanation of what breaks if deferred to v3.1. The "紧迫性" section is two sentences and reads more like a scheduling justification than an urgency argument. |

### 2. Solution Clarity: 105/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 38/40 | Very specific: surface-aware recipe generation, test.execution removal, surface-orchestration.yaml as interface contract. The fallback chain (4 levels) is precisely defined. Quote: "init-justfile 在生成 justfile 的同时，写入 `.forge/surface-orchestration.yaml` 文件" — actionable and unambiguous. |
| User-facing behavior described | 42/45 | Excellent table of "Surface 测试编排模式" shows exactly what happens per surface type. Recipe signatures with consumer and semantic expectations are comprehensive. Deduction: The mixed-project experience (web+api) is described in depth but the "无 surface 配置" fallback experience is only one bullet: "回退到当前行为（纯语言配方，run-tests 保持原有逻辑）" — what does the user actually see? |
| Technical direction clear | 25/35 | The just platform attributes `[linux]`/`[windows]`, PID tracking, probe retry logic are well-specified. However, the surface-orchestration.yaml schema is shown via example but never formally specified (required fields, optional fields, validation rules). The run-tests execution flow is listed as pseudo-code steps but not as a state machine diagram or formal grammar, despite the proposal mentioning "状态机驱动". |

### 3. Industry Benchmarking: 95/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 35/40 | Docker Compose healthcheck, Kubernetes readinessProbe, Cypress start-server-and-test, Makefile target dependency, GitHub Actions service containers — five real-world solutions cited with specific feature comparison (编排模型, 就绪检测, 进程管理). Good coverage. |
| At least 3 meaningful alternatives | 22/30 | Four alternatives listed including "do nothing". However, "Go 代码直接管理进程生命周期" is not a genuinely different approach for the current scope — the proposal itself acknowledges it's "v3.0.0 范围过大". The "仅 surface 感知，保留 test.execution" alternative is weakly motivated — described in one line as "治标不治本" without demonstrating why partial improvement is inferior. This reads as a straw man. |
| Honest trade-off comparison | 20/25 | The trade-off analysis section on "justfile 作为唯一抽象层的 trade-off" is good — acknowledges inability to switch strategies without editing justfile, acknowledges new surface types require run-tests updates. However, the "已知局限" only lists two items. The non-just path sacrifice (go test, npx vitest users must wrap) is discussed once but not compared quantitatively against how many users this affects. |
| Chosen approach justified against benchmarks | 18/25 | The "从行业方案借鉴的设计" table maps K8s → probe retry, Cypress → teardown, Docker Compose → declarative sequence. But the key differentiator — "run-tests 是 LLM agent 执行的 SKILL" — is stated but not rigorously justified against, say, Cypress's deterministic approach. The proposal acknowledges LLM unreliability and builds mitigations, but doesn't explain why the deterministic Go CLI alternative is deferred rather than prioritized, given that every industry solution uses deterministic code. |

### 4. Requirements Completeness: 90/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 35/40 | All five surface types covered with concrete scenarios. Mixed project (web+api) scenario is thorough. "无 surface 配置" fallback included. Missing: What happens when a surface is detected but the language has no test framework? What happens when probe succeeds but test fails mid-run — is there partial teardown? |
| Non-functional requirements | 35/40 | Excellent NFR table with 6 items: cross-platform, backward compatibility, observability, performance, reliability, just version. Each has a verification method. Deduction: The "性能" NFR targets "不超过 1 秒" for init-justfile overhead but no baseline measurement is provided — is 1 second tight or generous? The "可靠性" NFR mentions "故障注入测试" as verification but doesn't specify which faults to inject. |
| Constraints & dependencies | 20/30 | Surface info sources listed with priority rules. Forge-distribution.md path conventions referenced. Missing: No mention of just version compatibility testing matrix, no discussion of Windows CI availability, no dependency on specific test runner versions (playwright, maestro). The dependency readiness section claims "直接可行" in one sentence — insufficient for a proposal that touches 15-20 tasks. |

### 5. Solution Creativity: 75/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 30/40 | The surface-orchestration.yaml as a shared contract between producer (init-justfile) and consumer (run-tests) is a clean idea. The fallback chain (4 levels) is well-structured. However, the core pattern — start service, wait, test, teardown — is directly from industry (Cypress/K8s). The novelty is in the LLM-agent-aware mitigations (HARD-GATE, state machine, exit code enforcement), not in the orchestration itself. |
| Cross-domain inspiration | 25/35 | PID file management from Unix daemon conventions, probe from K8s health checks, teardown from Cypress. The state machine pattern for LLM agent control is borrowed from software engineering but applied in a novel context (constraining LLM behavior). No inspiration from outside software engineering — e.g., no reference to how test orchestration works in game development, embedded systems, or hardware verification, which face similar "start hardware → verify ready → run tests → cleanup" patterns. |
| Simplicity of insight | 20/25 | "justfile 作为唯一抽象层" is a clean simplification. Removing the delegation layer is elegant. However, the proposal's implementation is not simple — 15-20 tasks, 5 surface rule files, PID tracking, cross-platform branches, state machine enforcement, scope migration across 6+ components. The insight is simple but the solution is not. |

### 6. Feasibility: 80/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 35/40 | Just platform attributes exist. PID tracking is standard. The probe retry logic is straightforward. However, the Windows compatibility story is incomplete — quote: "`curl` 可能不可用——`just probe` 需要考虑使用 PowerShell `Invoke-WebRequest` 作为 fallback" — this is acknowledged but not resolved in the proposal. The Windows `[windows]` attribute example uses CMD syntax but most Forge users are likely on macOS/Linux; Windows support adds significant complexity. |
| Resource & timeline feasibility | 20/30 | "预计 15-20 个编码任务" — this is a substantial scope. The proposal lists at least 8 distinct sub-areas (5 surface rules, init-justfile update, run-tests simplification, config schema change, scope migration across 6+ components). No timeline estimate is provided (days? weeks? sprint count?). Quote: "中等范围" — but 15-20 tasks with cross-component scope migration is not "中等" by any reasonable definition. |
| Dependency readiness | 25/30 | Surface detection is in place. The discovery that `test.execution` is not mapped in Go struct but is used by LLM agent is valuable. However, the `GetConfigValue` extension is listed as a dependency with "需独立评审" — this is a blocking dependency that is not yet approved. |

### 7. Scope Definition: 65/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 25/30 | Each in-scope item lists specific files and changes. The config schema sub-proposal is well-bounded. However, the "scope 统一迁移" section lists changes across 6+ components (prompt.go, scope-assignment.md, quick-tasks SKILL.md, db-schema.md, init-justfile, 16 prompt templates) — this is a significant scope expansion that is presented as if it were a minor migration. |
| Out-of-scope explicitly listed | 20/25 | Clear list: no language template changes, no forge-cli/pkg/just gate sequence changes, no quality_gate.go or testrunner changes, no new CLI commands, no feature flags. However, "回滚基础设施" is listed as out-of-scope with "回滚通过 git revert 实现" — this is a rollback strategy, not an out-of-scope item. The distinction matters because git revert on 15-20 tasks is not trivial. |
| Scope is bounded | 20/25 | The proposal mentions "15-20 个编码任务" which provides a number, but no time boundary. The scope-统一迁移 sub-section was clearly added in response to earlier iteration feedback and inflates the original scope significantly. Quote: "scope 值域统一迁移：混合项目所有配方的 scope 参数值从 frontend/backend 迁移为 surfaces map key" — this touches fundamental infrastructure (prompt.go, 16 prompt templates) and could easily be a separate proposal. |

### 8. Risk Assessment: 70/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 25/30 | 6 risks listed. The "test [journey] 过滤与原生运行器不兼容" risk is specific and high-impact. However, some risks are missing: (1) no risk for PID file corruption (stale PID pointing to a different process), (2) no risk for Windows-specific failures in CI, (3) no risk for the scope migration breaking existing mixed projects during transition. |
| Likelihood + impact rated | 20/30 | The ratings show a suspicious pattern: 3/6 risks are "中可能性", 2 are "低可能性", none are "高可能性". This is optimistic for a 15-20 task scope that touches fundamental infrastructure. The "run-tests 简化导致已配置 test.execution 的项目不兼容" is rated "低可能性/低影响" with justification "v3.0.0 未发布，无存量用户" — this is honest but the risk should also consider v2 users who might upgrade. |
| Mitigations are actionable | 25/30 | Most mitigations are concrete: "回退到当前行为", "forge surfaces 基于路径检测", "Surface 规则记录映射关系". However, "Surface 规则过于泛化" mitigation is "LLM 组合语言模板 + surface 规则" — this is a description of how the system works, not a mitigation for when it fails. The rollback plan is clear but relies entirely on git revert. |

### 9. Success Criteria: 60/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 40/55 | Most criteria are checkbox items that can be verified: "CLI/TUI 不生成 run 配方" (checkable), "run-tests 不再依赖 test.execution.run" (checkable), "无 surface 配置的项目输出与当前一致" (diffable). However, several criteria are vague: "语言模板与 surface 规则的配方职责边界清晰" — what constitutes "清晰"? "端到端验证" — what constitutes "通过"? The criteria lack quantitative thresholds for the NFRs (e.g., no "init-justfile overhead < 1s" as a success criterion despite it being an NFR). |
| Coverage is complete | 20/25 | 11 criteria cover the main in-scope items. Gaps: (1) No success criterion for cross-platform compatibility despite it being an NFR. (2) No success criterion for PID file corruption handling. (3) No success criterion for the LLM agent determinism mitigations (HARD-GATE, exit code enforcement). (4) No success criterion for `.forge/surface-orchestration.yaml` validation. |

### 10. Logical Consistency: 80/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 32/35 | Surface-aware init-justfile directly addresses Problem 1. Removing test.execution delegation directly addresses Problem 2. The surface-orchestration.yaml contract connects both solutions. Minor gap: Problem 2 claims the delegation layer "增加了复杂度但没有增加灵活性", but the proposal's own trade-off analysis acknowledges that the new approach cannot support CI-specific orchestration without editing justfile — this is a flexibility loss that contradicts the original claim of "没有增加灵活性". |
| Scope ↔ Solution ↔ Success Criteria aligned | 25/30 | Mostly aligned: scope items map to solution components, success criteria cover scope items. However, the scope-统一迁移 touches prompt.go and 16 prompt templates — these are not reflected in the success criteria's "scope 值域统一迁移" item which only mentions "混合项目所有配方的 scope 参数值". The config schema sub-proposal is in scope but has a vague success criterion: "`test.execution` 节点从 config-schema 中完全删除". |
| Requirements ↔ Solution coherent | 23/25 | Requirements (5 surface scenarios + fallback + mixed project) map cleanly to the solution. The recipe signature table ensures no orphan requirements. Minor issue: the "下游集成契约" section defines `unit-test` recipe as consumed by "forge task submit, clean-code, fix-bug, testrunner" — but the proposal doesn't change unit-test behavior at all. Including it in the contract table is correct but slightly misleading about the scope of change. |

---

## ATTACKS

1. **[Problem Definition - Urgency]**: Cost of delay is entirely missing — Quote: "随着 v3.0.0 test profile 的引入，surface 成为测试流程的核心维度。init-justfile 和 run-tests 需要协同工作" — This explains timing, not urgency. What happens if this is deferred to v3.1? Does v3.0 ship with broken surface-aware testing? The proposal must quantify the cost of delay or acknowledge this is a scheduling preference, not an urgent need.

2. **[Industry Benchmarking - Straw-man alternative]**: "仅 surface 感知，保留 test.execution" is dismissed in one line as "治标不治本" — This is the textbook definition of a straw man. The proposal does not demonstrate why incremental improvement (surface awareness first, delegation cleanup later) is architecturally inferior. For a 15-20 task scope, the incremental approach would reduce risk significantly. The proposal must provide a structural argument for why partial implementation creates worse technical debt than the current state.

3. **[Feasibility - Resource]**: No timeline estimate provided — Quote: "预计 15-20 个编码任务" and "中等范围" — 15-20 tasks touching init-justfile, run-tests, config schema (Go CLI), scope-assignment rules, prompt.go, 16 prompt templates, 5 new surface rule files is not "中等范围". This is a major cross-cutting change. The proposal must either provide a sprint/week estimate or acknowledge this is a "large" scope and justify why it must be a single release.

4. **[Feasibility - Technical]**: Windows story is acknowledged but unresolved — Quote: "`curl` 可能不可用——`just probe` 需要考虑使用 PowerShell `Invoke-WebRequest` 作为 fallback" and the `[windows]` recipe example. The proposal sets a cross-platform NFR but does not commit to a specific Windows testing plan. Without Windows CI, the cross-platform claim is aspirational. Either commit to a Windows testing matrix or downgrade the NFR to "best-effort Windows support."

5. **[Requirements - Coverage]**: PID corruption scenario is unaddressed — The proposal details PID file management extensively but never discusses what happens when a stale PID points to an unrelated process (e.g., dev server crashed, OS reused the PID for a different process). The teardown section says "先校验 PID 有效性" by checking `/proc/<pid>` or `ps -p <pid>`, but process existence does not prove it's the dev server. The proposal must add a PID-validity check that verifies the process identity (e.g., checking the process command line matches expected dev server command).

6. **[Solution Clarity]**: surface-orchestration.yaml schema is example-only, not formally specified — The YAML example shows `version`, `surfaces` (with `type`, `dev_recipe`, `probe_target`, `teardown_order`), and `orchestration` (with `startup_order`, `probe_after`, `teardown_reverse`). But there is no specification of: required vs optional fields, validation rules, what happens when fields are missing, version migration strategy. For a file that is the "统一接口" between two skills, this informality is a gap.

7. **[Risk Assessment]**: Optimistic likelihood ratings — 6 risks, zero rated "高可能性". For a proposal that modifies 6+ infrastructure components and touches prompt templates used across the system, this is surprisingly optimistic. The scope migration alone (changing `frontend`/`backend` to surfaces map keys in prompt.go, 16 prompt templates, scope-assignment rules) has high likelihood of introducing regressions in existing workflows.

8. **[Success Criteria]**: NFR "性能" (init-justfile overhead < 1s) has no corresponding success criterion — The NFR table specifies a measurable threshold but the success criteria checklist does not include it. Either add it as a success criterion or acknowledge the NFR is aspirational guidance, not a gate.

9. **[Success Criteria]**: Cross-platform NFR has no success criterion — Quote from NFR: "Windows/macOS/Linux 三平台均可运行编排序列". But no success criterion verifies Windows behavior. The 11 criteria can all be met on macOS alone.

10. **[Scope]**: scope-统一迁移 is scope creep within the proposal — The original problem statement mentions nothing about `frontend`/`backend` scope values. This entire sub-proposal was added in response to earlier iteration feedback. While the migration is logical, it doubles the scope of change and touches fundamental infrastructure (prompt.go, 16 templates). The proposal should either (a) split this into a separate proposal, or (b) justify why it must be in the same release despite the increased risk.

11. **[Logical Consistency]**: Flexibility claim contradiction — Problem 2 claims test.execution "增加了复杂度但没有增加灵活性" (no added flexibility). But the trade-off analysis later acknowledges: "若需要在不修改 justfile 的情况下切换编排策略（如 CI 环境用不同启动命令），当前方案无法支持" (current approach cannot support this). This is a flexibility loss. The problem statement overclaims to strengthen the case for removal.

12. **[Requirements]**: No requirement for `.forge/surface-orchestration.yaml` migration — What happens when a user runs init-justfile again after editing surface-orchestration.yaml? Does it overwrite their changes? Preserve them? Merge? The file is user-editable (quote: "高级用户可直接编辑此文件调整编排行为") but no requirement addresses the regeneration scenario.

13. **[Solution Clarity]**: run-tests execution flow uses pseudo-code, not the state machine it claims — Quote: "run-tests 的编排步骤本质上是状态机（init → dev → probe → test → teardown）。SKILL.md 中显式声明当前步骤和下一步骤的映射关系" — But the actual execution flow is listed as sequential steps, not as a state machine with transitions, guards, and error states. The discrepancy between the claimed abstraction (state machine) and the actual specification (sequential steps) creates a gap for the LLM agent that must implement it.
