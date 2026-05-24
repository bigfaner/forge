# Evaluation Report: Pipeline Integration Stitch — Iteration 3

**Evaluator**: CTO (Adversarial)
**Date**: 2026-05-24
**Document**: `docs/proposals/pipeline-integration-stitch/proposal.md`
**Rubric**: `plugins/forge/skills/eval/rubrics/proposal.md`
**Previous Score**: 864/1000 (Iteration 2)
**Annotated Blind Review**: Pre-revised regions detected — attack density tracked separately.

---

## Iteration-2 Attack Resolution Audit

| # | Attack | Status | Evidence |
|---|--------|--------|----------|
| 1 | Trade-off still qualifies selected approach's con | **Fixed** | Removed "虽模式机械但总量不可忽视". Now reads: "变更量较大——14 文件 ~95 处测试引用需逐文件编辑，预估 P2 单独耗时 2h+". No qualifying clause. |
| 2 | Code-gen rejection rationale covers only templates | **Fixed** | Expanded: "现有代码库未使用 code-gen 模式；ROI 不合理——仅模板层面值得自动化的部分为 4 个文件（不值得建生成器），而 CategoryEval、RecordData 字段、RenderRecord case 等 adapter 逻辑涉及 Go 类型系统和业务语义，无法通过模板化代码生成覆盖". |
| 3 | log.Printf overclaimed as "注册机制" | **Fixed** | Innovation Highlights now says: "这不是注册机制本身（实际的分类注册是 `eval.` 前缀分支），而是为分类遗漏提供运行时诊断能力". Framing corrected to "防御性诊断层". |
| 4 | No novelty — Innovation Highlights should acknowledge limitation | **Fixed** | First paragraph now says: "本质上这是维护性工作——补全 Task 1 遗漏的 adapter 层实现，不引入新的架构模式或设计范式". Honest acknowledgment. |
| 5 | Performance NFR missing compatibility matrix | **Fixed** | New NFR: "版本范围: 本次变更影响的 Forge 版本为引入 eval 类型的版本（即 Task 1 合并后的版本）。旧版本 index.json 不包含 eval 类型，因此 CategoryEval 分支不会被触发，性能影响为零". |
| 6 | No observability NFR for eval tasks | **Fixed** | New NFR: "可观测性（显式排除）: eval 任务的日志/指标/监控不在本次范围内...eval task outcome observability 应作为独立提案处理". Explicitly scoped out with rationale. |
| 7 | eval skill delegation naming inconsistent | **Fixed** | New clarification: "委托链说明：`forge:eval` 是统一的 eval 入口 skill，根据 `--type` 参数内部分发到 `forge:eval-journey`（journey 评估）或 `forge:eval-contract`（contract 评估）的具体评估逻辑". Delegation chain now explicit. |
| 8 | eval record template mitigation is circular | **Fixed** | Risk mitigation now reads: "eval 任务为本次新建（无既有实现可参考），字段设计基于 eval skill 的输出契约：`plugins/forge/skills/eval/` 下的 rubric 文件定义了 score（0-1000 数值）、findings（问题列表）、severity（critical/major/minor）、passed（布尔门控结果）四个标准输出维度。record 模板字段直接映射这些维度". Source of truth identified. |
| 9 | RenderRecord criterion doesn't verify field formatting | **Fixed** | Expanded criterion: "单元测试验证：构造包含 eval 字段的 RecordData{Score: 850, Findings: []string{"finding1"}, Severity: "major", Passed: true} 的 CategoryEval record，调用 RenderRecord 后输出包含格式化后的 score 字符串、findings 列表、'major' 和 passed 标识". Concrete test case provided. |
| 10 | "所有现有测试通过" is underspecified | **Fixed** | Now reads: "所有现有测试通过（go test ./pkg/task/... ./pkg/prompt/... ./internal/cmd/... ./tests/...）". Test commands specified. |
| 11 | findFirstTestTaskIdx prefix matching retains same failure class | **Partially Fixed** | Extended justification provided explaining why task ID requires prefix (vs type constant), semantic stability argument added. But the fundamental issue — prefix coupling remains a failure class — is acknowledged as accepted trade-off. |
| 12 | Integration test deliverable missing from scope and resource estimate | **Fixed** | Resource estimate now reads: "预计 4 个 coding task（含集成测试）+ 2 个 doc task，总工作量 ~7h。集成测试覆盖 Quick mode 依赖链，预计增加 ~1h 工作量". Integration test reflected in estimates. |
| 13 | [blindspot] record-format-test.md update has no explicit success criterion | **Fixed** | New criterion: "record-format-test.md 包含全部五个新类型（test.gen-journeys、test.gen-contracts、test.gen-scripts、test.run、test.verify-regression）且不包含已废弃类型名（test.gen-cases、test.eval-cases、test.gen-and-run）". |
| 14 | [blindspot] RenderEvalRecord formatting fields have no verification | **Fixed** | Covered by expanded RenderRecord criterion with concrete test case verifying ScoreFormatted, FindingsFormatted, etc. |

**Resolution: 13/14 attacks fully addressed, 1 partially (prefix matching trade-off accepted but more deeply justified).**

---

## Phase 1: Reasoning Audit

### Problem -> Solution Trace

Problem: auto-gen-journeys-contracts 提案遗留 6 类执行层面缺陷 — P0 (4 个模板文件缺失), P1 (eval 类型分类错误 + 废弃类型匹配 + 顺序耦合), P2 (废弃代码残留 + 文档过期).

Solution: 三管齐下 — 创建 4 个模板文件, 新增 CategoryEval + 依赖注入加固 + findFirstTestTaskIdx 修复 + record-format 更新, 清理 gen-and-run 废弃代码.

Trace: P0 maps 1:1 to template creation. P1 maps to CategoryEval + dep hardening + prefix fix + record-format updates. P2 maps to gen-and-run removal + documentation cleanup. Direct and complete. No orphan problems or phantom solutions.

### Solution -> Evidence Trace

Evidence is code-structural: file paths, line numbers, function names, struct field definitions with JSON tags. No user incident data or reproduction steps — acceptable for internal tooling proposal where the evidence is source code itself.

### Evidence -> Success Criteria Trace

16 success criteria cover all scope items. New criteria address iteration-2 gaps (record-format-test.md, RenderRecord field formatting, test command specification). Traceability is now near-complete.

### Self-Contradiction Check

1. **Innovation framing is now honest**: "本质上这是维护性工作" explicitly acknowledges lack of innovation. log.Printf contribution correctly framed as "防御性诊断层" not "注册机制". No contradiction.

2. **Trade-off con no longer qualified**: "变更量较大——14 文件 ~95 处测试引用需逐文件编辑，预估 P2 单独耗时 2h+" — clean statement of fact with quantitative estimate. No qualifying clause.

3. **eval skill delegation chain now explicit**: `forge:eval` as unified entry point dispatching to `forge:eval-journey`/`forge:eval-contract` based on `--type` parameter. Clear and unambiguous.

4. **Remaining tension**: `--target 850` hardcoded in prompt template — justified by "当前硬编码为 850 是合理的...若未来需要不同阈值，应在 eval skill 的命令定义中参数化，而非在 prompt 模板中暴露". Reasonable but creates a coupling point. Not a contradiction, but a conscious trade-off.

---

## Phase 2: Rubric Scoring

### 1. Problem Definition (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 39/40 | Six concrete items with P0/P1/P2 prioritization. Each has file paths and line numbers. P0 has clear failure chain (missing file -> ReadFile failure -> pipeline crash). Deduction: P1 item 4 "顺序耦合" — the trigger condition for reordering is still not specified. What scenario causes the reordering? The "should" in "应合并为单步操作消除耦合" is a design preference, not a demonstrated bug trigger. |
| Evidence provided | 39/40 | Strong structural evidence: file paths, line numbers, specific failure modes. P0 has concrete "Synthesize() ReadFile 失败" chain. Deduction: no reproduction steps or incident frequency data (how many features have attempted to use these types since Task 1 merged?). |
| Urgency justified | 28/30 | P0 = "pipeline 执行必定失败" — urgency is self-evident. Cost of delay clear. Deduction: no quantification of blast radius (how many features/teams currently blocked or workaround status). |

**Subtotal: 106/110**

### 2. Solution Clarity (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 39/40 | Three-pronged approach with specific files, function signatures, content patterns, struct field definitions with JSON tags. Template structure described in four-part pattern (header/role/constraints/workflow). Deduction: `resolveTestDepsAndInjectReviewDoc` `mode string` parameter — only "quick" is shown; what are other valid values and their behavior? |
| User-facing behavior described | 41/45 | Observable behaviors described: `forge prompt get-by-task-id` returns valid prompts, `forge submit-task` accepts/rejects by category, `RenderRecord` uses eval template. Error message for deprecated type shown. Deduction: no before/after CLI output examples. What exact output does a developer see when running `forge prompt get-by-task-id` for an eval task? What does the migration error message look like? |
| Technical direction clear | 35/35 | Extremely detailed: Go struct fields with JSON tags, function signatures, file paths with line numbers, removal ordering (5 steps for P2), test file lists with line numbers. eval skill delegation chain explained. `--target 850` justified. findFirstTestTaskIdx prefix choice explained with template-discovery vs task-location distinction. |

**Subtotal: 115/120**

### 3. Industry Benchmarking (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 35/40 | Three concrete solutions: Spring Boot auto-configuration (`@Configuration` + `@ConditionalOnClass`), ASP.NET Core convention-based endpoint discovery (`*Controller` naming), Temporal workflow/activity registration (worker registration requirement). Airflow analogy retained. Each with mechanism and Forge correspondence. Deduction: all three are "type registration + discovery" pattern from web framework / workflow engine domain. No breadth beyond this pattern family — e.g., no plugin system (VS Code extension registration, JetBrains plugin descriptor), no build system (Gradle task registration, Makefile convention), no schema-based validation framework. |
| At least 3 meaningful alternatives | 27/30 | Five alternatives: (1) minimal fix, (2) P0+P1 only, (3) P0+P1+P2 (selected), (4) Code-gen (architecturally different), (5) Schema-based init-time validation (architecturally different). Deduction: alternatives 1-3 are still scope variants rather than genuinely different approaches; no explicit "do nothing" alternative (though P0 severity arguably precludes it). |
| Honest trade-off comparison | 22/25 | Selected approach con now reads: "变更量较大——14 文件 ~95 处测试引用需逐文件编辑，预估 P2 单独耗时 2h+". No qualifying clause. Code-gen rejection expanded to cover full adapter surface: "CategoryEval、RecordData 字段、RenderRecord case 等 adapter 逻辑涉及 Go 类型系统和业务语义，无法通过模板化代码生成覆盖". Deduction: (1) no quantitative comparison across alternatives (e.g., estimated effort for each); (2) Schema-based validation rejection says "需设计 schema 格式" but doesn't estimate effort vs benefit. |
| Chosen approach justified against benchmarks | 21/25 | "完整修复与行业最佳实践一致（Spring/Temporal 的完整注册模型）". Code-gen rejection now covers full adapter surface. Deduction: (1) does not explain why Forge doesn't adopt Temporal's init-time fail-fast for all missing templates (only gen-and-run gets migration error, but what about future types?); (2) Schema-based rejection is thin — "本次 P0 问题本质是遗漏创建文件而非验证缺失" is true for this specific case but doesn't address whether schema validation would prevent future occurrences. |

**Subtotal: 105/120**

### 4. Requirements Completeness (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 38/40 | Five key scenarios with edge cases: mixed submission, old index.json, migration error, Quick mode dependency chain. Error scenarios covered (deprecated type reference, invalid field submission). Deduction: (1) no scenario for eval template content mismatch with actual eval skill behavior — what happens if `forge:eval` skill output format changes after record template is created?; (2) no scenario for auto-discovery filename collision (two templates with similar names). |
| Non-functional requirements | 38/40 | Five NFRs + one explicit exclusion: backward compatibility (two-path migration error), CategoryEval test coverage (positive/negative/boundary), eval record template (score/findings/severity fields), performance (< 1μs), version scope (Task 1+ only), observability (explicitly excluded with rationale). Significantly improved from iteration 2. Deduction: (1) no Go version compatibility requirement; (2) no requirement for thread safety or concurrent access (though likely not relevant for CLI tool). |
| Constraints & dependencies | 26/30 | Task 1 completed clearly stated. No external dependencies. Deduction: (1) no Go version constraint mentioned; (2) no dependency on Forge distribution model noted (per CLAUDE.md requirement for plugins/forge/ modifications — though this proposal modifies `pkg/task/` not `plugins/forge/`, so may not apply). |

**Subtotal: 102/110**

### 5. Solution Creativity (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 20/40 | Proposal honestly states: "本质上这是维护性工作——补全 Task 1 遗漏的 adapter 层实现，不引入新的架构模式或设计范式". This is gap-filling work. log.Printf diagnostic layer is a reasonable defensive pattern but correctly framed as "低成本的防御性诊断辅助" — not overclaimed. Modest contribution beyond industry baseline. |
| Cross-domain inspiration | 12/35 | Spring/ASP.NET/Temporal cited as pattern analogies but not explored for cross-domain borrowing. No inspiration from build systems (Gradle incremental build invalidation), package managers (npm deprecation warnings), or compiler design (type system exhaustiveness checking). The domain is narrow (internal tooling adapter layer), limiting cross-domain applicability, but the proposal doesn't attempt to find any. |
| Simplicity of insight | 22/25 | P0/P1/P2 triage is clean and well-structured. "Adapter layer needs completion when types are registered" insight is straightforward. log.Printf as low-cost defense is a clean idea. The `resolveTestDepsAndInjectReviewDoc` merge is a sensible simplification. Deduction: the five-step removal ordering for P2 is detailed but not particularly insightful — it's standard "remove consumers before definitions" practice. |

**Subtotal: 54/100**

### 6. Feasibility (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 38/40 | Highly feasible. test-gen templates follow existing test-gen-scripts.md pattern. eval templates described as "全新模式" but with clear content structure (header/role/constraints/workflow/record-fields). CategoryEval follows CategoryTest blueprint. Removal is mechanical with explicit ordering. Deduction: eval template "全新模式" — the proposal describes what the template should contain but acknowledges it's a new pattern. Design iteration risk acknowledged but not quantified. |
| Resource & timeline feasibility | 28/30 | "预计 4 个 coding task（含集成测试）+ 2 个 doc task，总工作量 ~7h。集成测试覆盖 Quick mode 依赖链，预计增加 ~1h 工作量". Integration test now reflected in estimate. Deduction: (1) eval template "全新模式" design time not separately estimated — may require iteration; (2) P2 "预估 P2 单独耗时 2h+" is a separate estimate — should clarify if included in 7h total. |
| Dependency readiness | 28/30 | Task 1 completed. No external deps. Deduction: no mention of whether Task 1 left any residual issues or whether the auto-discovery mechanism has been validated in production. |

**Subtotal: 94/100**

### 7. Scope Definition (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 28/30 | Each item specifies exact file path, line numbers, content description, struct field definitions with JSON tags. Template patterns described structurally (header/role/constraints/workflow). Function signatures provided. Removal ordering explicit. Deduction: `resolveTestDepsAndInjectReviewDoc` `mode` parameter valid values still not fully enumerated. |
| Out-of-scope explicitly listed | 24/25 | Five items with rationale: historical docs (~35 files ~130 places), resolveBreakdownDeps/resolveQuickDeps refactoring, eval rollback, old index.json auto-migration, record-format-doc.md rename. Quantified where relevant. |
| Scope is bounded | 24/25 | 4 coding tasks + 2 doc tasks, ~7h total. P2 test file list is exhaustive (14 files with line numbers). Integration test included in estimates. Deduction: 7h total includes integration test (+1h) but doesn't clarify relationship with "P2 单独耗时 2h+" — is P2 within or additional to the 7h? |

**Subtotal: 76/80**

### 8. Risk Assessment (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 28/30 | Seven risks: template inaccuracy, CategoryEval field mismatch, gen-and-run removal compilation failure, Synthesize() file-not-found, findFirstTestTaskIdx dependency wiring, default branch silent misclassification, eval record template design. Deduction: no risk for eval template "全新模式" design iteration — the proposal acknowledges it's a new pattern but doesn't list the risk of needing multiple iterations to get it right. |
| Likelihood + impact rated | 28/30 | Varied ratings: M/H, L/M, M/H, M/M, M/M, L/M, L/M. Generally honest. eval record template now L/M with expanded justification. Deduction: eval record template impact may be underestimated — if field design is wrong, all eval submissions fail, which could be H impact. But the mitigation (field design based on eval skill output contract) reduces likelihood. |
| Mitigations are actionable | 27/30 | Specific mitigations: reference existing patterns + validation workflow + integration test, per-step `go build ./...`, two-path migration error (validate_index.go + Synthesize()), log.Printf warning, field design based on eval skill output contract. Deduction: (1) "参考现有 eval 任务的 submit-task 实际字段设计" from iteration 2 has been replaced with "字段设计基于 eval skill 的输出契约" — now references a concrete source of truth (rubric files); (2) some mitigations are verification activities (run tests) rather than preventive measures. |

**Subtotal: 83/90**

### 9. Success Criteria (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 51/55 | 16 checkbox items. Most use grep commands, function return values, concrete test data, or test assertions. `resolveTestDepsAndInjectReviewDoc` has two-condition test. RenderRecord has concrete test case with specific struct values. Test commands specified (`go test ./pkg/task/... ./pkg/prompt/... ./internal/cmd/... ./tests/...`). record-format-test.md has explicit type list verification. log.Printf criterion uses string matching. Deduction: (1) "Synthesize() 对...返回有效 prompt" — what constitutes "有效"? Should specify expected content structure or field presence; (2) `findFirstTestTaskIdx` "正确返回 gen-journeys 任务索引" — what index value? Should specify expected return value. |
| Coverage is complete | 24/25 | Covers P0 (template existence + Synthesize), P1 (CategoryForType + submit-task + RenderRecord + findFirstTestTaskIdx + dep merge + record-format-eval.md + log.Printf + integration test + record-format-test.md), P2 (grep zero results + migration error + tests pass). record-format-test.md now has explicit criterion. Deduction: (1) `--target 850` hardcoded value has no verification criterion — should the prompt template contain this value and how to verify? (2) eval template content correctness (skill delegation command format) unverified — only existence checked. |

**Subtotal: 75/80**

### 10. Logical Consistency (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses stated problem | 34/35 | Direct 1:1 mapping: P0 -> templates, P1 -> CategoryEval + dep hardening + prefix fix + record-format updates, P2 -> gen-and-run cleanup. No orphan problems. Deduction: findFirstTestTaskIdx replaces one hardcoded prefix with another — the explanation is now thorough and coherent, distinguishing template discovery from task location, but the fundamental trade-off (prefix coupling) is acknowledged as accepted. |
| Scope <-> Solution <-> Success Criteria aligned | 29/30 | Strong alignment. Each scope item has corresponding success criterion. P2 grep criteria cover code removal in both forge-cli/ and plugins/. Integration test aligns with scope deliverable. record-format-test.md now has dedicated criterion. Deduction: `--target 850` in Scope (eval template content) has no dedicated success criterion. |
| Requirements <-> Solution coherent | 24/25 | NFRs map to solution: backward compat -> two-path migration error, test coverage -> unit tests, eval record template -> record-format-eval.md, performance -> prefix comparison analysis, version scope -> Task 1+, observability -> explicitly excluded with rationale. Deduction: (1) "Mixed feature 依赖注入" scenario has solution but no corresponding NFR — what quality attributes apply to the merged function?; (2) Innovation framing is now honest — "维护性工作" matches actual content, no overclaim. |

**Subtotal: 87/90**

---

## Phase 3: Blindspot Hunt

### [blindspot] 1: Synthesize() "有效 prompt" 判定标准模糊
Quote: Success criterion: "Synthesize() 对 test.gen-journeys、test.gen-contracts、eval.journey、eval.contract 返回有效 prompt（P0 修复）"
"有效 prompt" is not defined. Does "有效" mean non-empty? Contains expected sections? Renders without error? The criterion should specify what constitutes a valid prompt — e.g., "返回包含 TASK_ID 占位符替换后值的 prompt，且包含 role/constraints/workflow 三段结构". Without this, an implementer could produce any non-error output and claim the criterion is met.

### [blindspot] 2: eval 模板 "全新模式" 的设计验证策略不充分
Quote: "eval 类型（eval-journey.md、eval-contract.md）为全新模式，与 test-gen 根本不同"
Risk item for template inaccuracy says: "模板创建后立即运行 forge prompt get-by-task-id 验证渲染输出". But this only verifies rendering, not content correctness. For a "全新模式", there should be a content review step — e.g., running the generated prompt through the eval skill and verifying the eval produces expected results. The gap between "template renders without error" and "template produces correct eval behavior" is unaddressed.

### [blindspot] 3: P2 移除顺序与编译依赖可能有隐藏耦合
Quote: "按此顺序确保增量编译通过：infer.go → prompt.go → validate_index.go → build.go → types.go"
This ordering assumes a linear dependency chain, but Go compilation operates at package level, not file level. `go build ./...` compiles all files in a package together. The ordering is conceptually correct (remove consumers before definitions) but the "增量编译" framing implies per-file compilation which doesn't apply in Go. The mitigation should specify: "after removing all references across files, run `go build ./pkg/task/...` to verify the package compiles". The current per-step `go build ./...` is correct but the rationale for the specific ordering is misleading about Go's compilation model.

### [blindspot] 4: resolveTestDepsAndInjectReviewDoc mode 参数合法值未完全定义
Quote Scope: "func resolveTestDepsAndInjectReviewDoc(testTasks []AutoGenTaskDef, index *TaskIndex, mode string, needsEval bool)"
The function signature accepts a `mode string` parameter. Only "quick" mode is shown in examples and success criteria. What other modes exist? If there are other modes (e.g., "breakdown", "staged"), their behavior under the merged function is undefined. This is a gap from iteration 2 that persists — the function signature is provided but the parameter contract is incomplete.

### [blindspot] 5: grep 验证排除 docs/proposals 但不影响活跃文档
Quote: "grep -r 'gen-and-run|quick-gen-and-run|T-quick-gen' forge-cli/ --exclude-dir=docs/proposals" and similar for plugins/
The `--exclude-dir=docs/proposals` correctly excludes historical proposal docs, but the Scope says "活跃文档更新: docs/OVERVIEW.md、docs/WORKFLOW.md". These are in `docs/`, not `forge-cli/` or `plugins/`. The grep commands don't cover the docs directory at all — so there's no automated verification that OVERVIEW.md and WORKFLOW.md are actually cleaned up. The grep criteria verify code removal but leave active documentation cleanup unverified.

---

## Annotated Region Attack Density

| Region Type | Attack Count | Notes |
|-------------|-------------|-------|
| Pre-revised (annotated) | 2 | Focus on whether revision introduced new issues |
| Unannotated | 4 | Standard adversarial review |
| **Total** | **6** | |

Pre-revised attacks focused on: (1) eval template content pattern description — now very detailed with four-part structure, but "全新模式" design verification strategy is thin; (2) gen-and-run removal ordering — now has 5-step explicit order with compilation verification, but Go package-level compilation model makes per-file ordering rationale misleading.

No `conflict-with-pre-revision` tags — all rubric judgments are consistent with pre-revision direction.

---

## Score Summary

| Dimension | Score | Max | Delta from Iter 2 |
|-----------|-------|-----|-------------------|
| Problem Definition | 106 | 110 | +2 |
| Solution Clarity | 115 | 120 | +2 |
| Industry Benchmarking | 105 | 120 | +7 |
| Requirements Completeness | 102 | 110 | +3 |
| Solution Creativity | 54 | 100 | +9 |
| Feasibility | 94 | 100 | +2 |
| Scope Definition | 76 | 80 | +1 |
| Risk Assessment | 83 | 90 | +2 |
| Success Criteria | 75 | 80 | +3 |
| Logical Consistency | 87 | 90 | +2 |
| **Total** | **897** | **1000** | **+33** |

---

## ATTACKS

1. [Success Criteria]: Synthesize() "返回有效 prompt" 判定标准模糊 — quote: "Synthesize() 对 test.gen-journeys、test.gen-contracts、eval.journey、eval.contract 返回有效 prompt（P0 修复）" — "有效" 未定义。应指定 prompt 必须包含的内容结构（如 role/constraints/workflow）或具体验证方法（如检查 TASK_ID 占位符已替换）。一个返回空字符串的 non-error 结果也会被判定为"有效"。

2. [Success Criteria]: eval 模板内容正确性仅有渲染验证 — quote: "模板创建后立即运行 forge prompt get-by-task-id 验证渲染输出" — 但对于"全新模式"的 eval 模板，渲染成功不等于内容正确。应增加端到端验证：用渲染后的 prompt 执行 eval skill 并验证产出符合预期格式。

3. [Scope Definition]: resolveTestDepsAndInjectReviewDoc mode 参数合法值未定义 — quote: "func resolveTestDepsAndInjectReviewDoc(testTasks []AutoGenTaskDef, index *TaskIndex, mode string, needsEval bool)" — 仅展示 "quick" mode，其他 mode 值（如 "breakdown"、"staged"）的行为未定义。函数签名提供了但参数契约不完整。必须在 Scope 中列出所有合法 mode 值及其行为差异。

4. [Scope Definition]: grep 验证不覆盖活跃文档清理 — quote: grep criteria 使用 `forge-cli/` 和 `plugins/` 目录，但 Scope 中活跃文档更新项（`docs/OVERVIEW.md`、`docs/WORKFLOW.md`）位于 `docs/` 目录，不在 grep 扫描范围内。必须增加对 `docs/` 目录中指定文件的 gen-and-run 引用清理验证。

5. [Feasibility]: P2 耗时估算与总计关系不明确 — quote: "总工作量 ~7h" 和 "预估 P2 单独耗时 2h+" — 7h 是否已包含 P2 的 2h+？如果是，P0+P1 仅 ~5h 是否合理？如果不是，总工作量应为 ~9h+。必须澄清 P2 2h+ 是否在 7h 之内。

6. [Industry Benchmarking]: Schema-based init-time validation 拒绝理由薄弱 — quote: "本次 P0 问题本质是遗漏创建文件而非验证缺失" — 这对本次问题成立，但不解释为什么 schema validation 不值得作为防御性措施防止未来再次遗漏。P0 发生的根本原因（人工遗漏文件创建）正是 schema validation 要防止的。拒绝理由仅对当前修复合理，对长期防御不足。
