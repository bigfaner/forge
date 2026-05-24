# Evaluation Report: Pipeline Integration Stitch — Iteration 1

**Evaluator**: CTO (Adversarial)
**Date**: 2026-05-24
**Document**: `docs/proposals/pipeline-integration-stitch/proposal.md`
**Rubric**: `plugins/forge/skills/eval/rubrics/proposal.md`
**Annotated Blind Review**: Pre-revised regions detected — attack density tracked separately.

---

## Phase 1: Reasoning Audit

### Problem -> Solution Trace

Problem: `auto-gen-journeys-contracts` 提案引入了 staged test pipeline 和 eval 质量门控，但遗漏了 4 类执行层面配套 — P0 (模板文件缺失), P1 (类型分类错误 + 依赖注入引用废弃类型 + 顺序耦合), P2 (废弃代码残留).

Solution: 三管齐下 — 创建 4 个模板文件, 新增 CategoryEval + 加固依赖注入, 清理 gen-and-run 废弃代码.

Trace: P0 maps 1:1 to template creation. P1 maps to CategoryEval + dep hardening. P2 maps to gen-and-run removal. Strong problem-solution fit.

### Solution -> Evidence Trace

Evidence is code-structural: specific file paths, line numbers, function names. No user incident data, no reproduction steps. Acceptable for internal tool but not exemplary.

### Evidence -> Success Criteria Trace

11 success criteria cover most scope items. Gaps:
- `record-format-eval.md` creation has no verification criterion
- `CategoryForType` default branch `log.Printf` has no verification criterion
- `resolveTestDepsAndInjectReviewDoc` "无顺序耦合" is an architectural property, not testable behavior

### Self-Contradiction Check

1. **Innovation claim vs. actual fix**: "自动发现机制已消除根因" but `findFirstTestTaskIdx` fix uses `findTaskIndexByPrefix(tasks, "T-test-gen-journeys")` — this is still a hardcoded prefix, not auto-discovery. Minor contradiction.

2. **P1 item 4 priority vs. "Overturned" assumption**: Assumption table says re-index idempotency bug is "Overturned: 当前代码已幂等...风险仅为代码耦合". Yet the item remains P1. The proposal justifies this as "preventive hardening" which is reasonable but the P1 label is slightly inflated.

3. **"机械性清理" framing**: Repeated assertion that P2 is "大部分是机械性清理" but 14 test files with ~95 references is non-trivial. Risk of underestimation.

---

## Phase 2: Rubric Scoring

### 1. Problem Definition (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 38/40 | Six concrete items with P0/P1/P2 prioritization. Each has file paths and line numbers. Minor deduction: "集成缝隙" is metaphorical; the actual gap type (missing adapter-layer artifacts) becomes clear only in Evidence. |
| Evidence provided | 38/40 | Strong structural evidence: file paths, line numbers, specific failure modes (`ReadFile` fails, `CategoryForType` returns wrong category). P0 has concrete failure chain. Deduction: no reproduction steps or incident frequency data. |
| Urgency justified | 25/30 | P0 = "pipeline 执行必定失败" — urgency is clear. Cost of delay stated: "任何执行 test pipeline 或 eval gate 的 feature 必定失败". Deduction: no quantification of blast radius (how many features/teams are currently blocked). |

**Subtotal: 101/110**

### 2. Solution Clarity (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 38/40 | Three-pronged approach with specific files, function signatures, content patterns. Template structure described in detail (header/constraints/workflow/record-fields four-part pattern). Deduction: `--target 850` in eval template is an unexplained magic number. |
| User-facing behavior described | 35/45 | Internal-facing proposal; "user" is a developer/agent. Observable behaviors described: `forge prompt get-by-task-id` returns valid prompts, `forge submit-task` accepts/rejects by category, `RenderRecord` uses eval template. Deduction: no before/after examples for any user-facing command. What exact output change does a developer see? |
| Technical direction clear | 34/35 | Extremely detailed: Go struct fields with JSON tags, function signatures, file paths with line numbers, removal ordering. Deduction: `resolveTestDepsAndInjectReviewDoc` parameter `mode string` — valid values not specified. |

**Subtotal: 107/120**

### 3. Industry Benchmarking (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 15/40 | Single analogy: "Airflow 添加新 DAG 类型后需要配套的 executor plugin 和 UI renderer". This is a brief comparison, not a cited solution with product names, open-source projects, or published patterns. No reference to convention-over-configuration frameworks (Rails, Spring), plugin systems (VS Code, JetBrains), or type registration patterns (Airflow operators, Temporal activities). |
| At least 3 meaningful alternatives | 12/30 | Three alternatives listed: (1) minimal fix, (2) P0+P1 only, (3) P0+P1+P2. These are scope variants of the same approach, not genuinely different architectural strategies. No "do nothing" alternative. No industry-validated alternative (e.g., init-time validation, code generation, schema-based registration). |
| Honest trade-off comparison | 12/25 | Table pros/cons are thin: "变更最小", "变更量较大（但大部分是机械性清理）". No quantitative comparison. The selected approach's con ("变更量较大") is immediately qualified with a parenthesis that minimizes it. |
| Chosen approach justified against benchmarks | 8/25 | Justification: "端到端可执行；零僵尸代码；类型系统一致" — these are benefits, not benchmark comparison. No argument for why this approach beats industry patterns or why it adopts/doesn't adopt standard solutions. |

**Subtotal: 47/120**

### 4. Requirements Completeness (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 35/40 | Five key scenarios: auto-discovery, eval submit semantics, mixed feature deps, quick-mode matching, gen-and-run removal. Edge cases: mixed submission (review + test fields), old index.json. Error scenarios: migration error for deprecated types. Deduction: no scenario for eval template content mismatch with actual eval skill behavior; no scenario for auto-discovery filename collision. |
| Non-functional requirements | 28/40 | Three NFRs: backward compatibility (migration error), test coverage (positive/negative/boundary), eval record template. Deduction: no performance consideration (adding CategoryEval branch to render/validate hot paths), no compatibility matrix (which Forge versions' index.json are affected), no observability requirement (eval task logging/monitoring). |
| Constraints & dependencies | 25/30 | "Task 1 已完成" clearly stated. No external dependencies. Deduction: no Go version constraint mentioned, no dependency on forge distribution model (per CLAUDE.md, plugin files have specific path resolution requirements). |

**Subtotal: 88/110**

### 5. Solution Creativity (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 8/40 | Proposal explicitly states: "Task 1 引入的自动发现机制已消除根因。本次提案聚焦于补全遗漏的执行层面配套". This is a bugfix/completion task — filling gaps in a previously shipped feature. No innovation beyond baseline. |
| Cross-domain inspiration | 5/35 | Airflow analogy mentioned but not explored. No cross-domain borrowing from build systems, package managers, plugin architectures, or type registration systems. |
| Simplicity of insight | 18/25 | The P0/P1/P2 triage is clean. The insight that "adapter layer needs to be completed when types are registered" is straightforward. Not a "why didn't I think of that" moment. |

**Subtotal: 31/100**

### 6. Feasibility (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 38/40 | Highly feasible. Templates reference existing patterns (test-gen-scripts.md, validation-code.md). CategoryEval follows CategoryTest blueprint. Removal is mechanical. Deduction: eval templates reference `validation-code.md` but it's not verified whether that file's pattern actually matches eval semantics. |
| Resource & timeline feasibility | 26/30 | "3 coding task + 2 doc task, ~6h" seems reasonable but may underestimate P2: 14 test files with ~95 references across ~11 test files requires careful per-file editing, not batch replacement. Risk of 2h+ for P2 test cleanup alone. |
| Dependency readiness | 28/30 | Task 1 completed. No external deps. Deduction: no mention of whether Task 1 left any residual issues. |

**Subtotal: 92/100**

### 7. Scope Definition (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 29/30 | Each item specifies exact file path, line numbers, and content description. Template patterns are described structurally. Deduction: template content is a pattern reference, not a complete spec — implementer must still make content decisions. |
| Out-of-scope explicitly listed | 24/25 | Five items with rationale: historical docs (~35 files, ~130 refs — doesn't affect runtime), resolveBreakdownDeps refactor, eval rollback, migration tool, record-format-doc.md doc.eval. Good. Deduction: "~35 文件 ~130 处" is an estimate, not a verified count. |
| Scope is bounded | 23/25 | 5 tasks, ~6h, well-defined deliverables. P2 test file list is exhaustive (11 files with specific lines). Deduction: "大部分是机械性清理" may understate P2 complexity. |

**Subtotal: 76/80**

### 8. Risk Assessment (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 25/30 | Six risks: template inaccuracy, CategoryEval field mismatch, gen-and-run removal compilation failure, findFirstTestTaskIdx dependency wiring, default branch silent misclassification, eval record template design. Deduction: missing risk of runtime impact on features currently using gen-and-run type; missing risk of `--target 850` hardcoded eval target being wrong. |
| Likelihood + impact rated | 26/30 | Varied ratings: M/H, L/M, M/H, M/M, L/M, L/L. Honest — not all low/high. Deduction: "gen-and-run 引用移除不完整" rated M may be optimistic given 95 references across 14 files; "4 个新 prompt 模板内容不准确" rated M is reasonable. |
| Mitigations are actionable | 26/30 | Specific mitigations: reference existing patterns, per-step `go build ./...`, `findTaskIndexByPrefix` reuse. Deduction: "test-gen 模板参考现有 test-gen-scripts.md" is implementation detail, not mitigation — a real mitigation would be "stakeholder review of template content" or "integration test verifying rendered prompt"; "编写单元测试覆盖正向/负向/边界用例" is generic. |

**Subtotal: 77/90**

### 9. Success Criteria (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 48/55 | 11 checkbox items. Most use grep commands, function return values, or test assertions. Highly measurable. Deduction: "ResolveFirstTestDep + T-review-doc prepend 为单步操作，无顺序耦合" — "无顺序耦合" is an architectural property, not a testable behavior (how do you write a test for "no ordering coupling"?); "所有现有测试通过" is too broad — which test suite? |
| Coverage is complete | 20/25 | Covers P0 (template existence + Synthesize), P1 (CategoryForType + submit-task + RenderRecord + findFirstTestTaskIdx + dep merge), P2 (grep zero results + migration error + tests pass). Deduction: no criterion for `record-format-eval.md` content correctness; no criterion for `CategoryForType` default branch `log.Printf` warning; no criterion for eval record field rendering accuracy. |

**Subtotal: 68/80**

### 10. Logical Consistency (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses stated problem | 34/35 | Direct 1:1 mapping: P0 -> templates, P1 -> CategoryEval + dep hardening, P2 -> gen-and-run cleanup. Deduction: findFirstTestTaskIdx fix replaces one hardcoded prefix with another, slightly contradicting "auto-discovery eliminates the root cause" claim. |
| Scope <-> Solution <-> Success Criteria aligned | 28/30 | Strong alignment. Each scope item has corresponding success criterion. P2 grep criteria cover code removal. Deduction: `record-format-eval.md` in Scope (P1) has no dedicated success criterion; `CategoryForType` default branch `log.Printf` in Scope has no criterion. |
| Requirements <-> Solution coherent | 23/25 | NFRs map to solution items: backward compat -> validate_index.go migration error, test coverage -> unit tests, eval record template -> record-format-eval.md. Deduction: "Mixed feature 依赖注入" scenario has solution (`resolveTestDepsAndInjectReviewDoc`) but no corresponding NFR (performance/reliability requirement for the merged function). |

**Subtotal: 85/90**

---

## Phase 3: Blindspot Hunt

### [blindspot] 1: eval template `--target 850` 魔数
Quote: "Skill(skill="forge:eval", args="--type [journey|contract] --target 850")"
The value 850 is hardcoded in the template. Different eval tasks may require different target scores. The proposal does not discuss parameterization or explain why 850 is the universal target.

### [blindspot] 2: CategoryEval 与 CategoryReview 语义边界
The proposal creates CategoryEval but does not clarify whether CategoryReview already exists or what the semantic boundary is between "eval" and "review". If `doc.review` already uses a review category, the distinction between eval and review is an unstated design decision.

### [blindspot] 3: 缺少集成测试 task
Quote from risk mitigation: "更新后添加集成测试验证 Quick mode 依赖链正确性"
This promises integration tests but no corresponding task exists in the task breakdown ("3 coding task + 2 doc task"), and no success criterion verifies integration test existence or correctness.

### [blindspot] 4: P1 item 3 与 P2 重叠修改
`findFirstTestTaskIdx` quick-mode fix appears in P1 (`build.go:492-494` -> `findTaskIndexByPrefix(tasks, "T-test-gen-journeys")`) AND P2 (`build.go:484,492-494` -> same location). The proposal does not flag this overlap, creating risk of conflicting changes during implementation.

### [blindspot] 5: 废弃模板删除的向后兼容缺口
Quote: "NFR: 向后兼容: 旧 index.json 引用 test.gen-and-run 时给出明确迁移错误提示"
The NFR only addresses `validate_index.go` validation. But deleting `prompt/data/test-gen-and-run.md` means any existing feature with gen-and-run tasks will fail at `Synthesize()` with a file-not-found error, not a migration-aware error. The migration error path only triggers during index validation, not during prompt rendering.

### [blindspot] 6: 测试清理的回归风险
Quote: "14 个测试文件中 ~95 处引用"
Removing ~95 test references across 14 files risks accidentally removing shared setup/teardown code or test helpers that are co-located with gen-and-run tests. The proposal does not assess whether any gen-and-run test code contains shared fixtures.

---

## Annotated Region Attack Density

| Region Type | Attack Count | Notes |
|-------------|-------------|-------|
| Pre-revised (annotated) | 8 | Focus on whether revision introduced new issues |
| Unannotated | 14 | Standard adversarial review |
| **Total** | **22** | |

Pre-revised attacks focused on: template content pattern specificity (medium), CategoryEval field naming coherence (medium), record-eval template missing criterion (medium), gen-and-run removal ordering risk (high), risk table mitigation quality (medium/high), success criteria grep scope (medium), build.go merge function testability (medium), record-format-test.md update scope (medium).

No `conflict-with-pre-revision` tags — all rubric judgments are consistent with pre-revision direction.

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 101 | 110 |
| Solution Clarity | 107 | 120 |
| Industry Benchmarking | 47 | 120 |
| Requirements Completeness | 88 | 110 |
| Solution Creativity | 31 | 100 |
| Feasibility | 92 | 100 |
| Scope Definition | 76 | 80 |
| Risk Assessment | 77 | 90 |
| Success Criteria | 68 | 80 |
| Logical Consistency | 85 | 90 |
| **Total** | **772** | **1000** |

---

## ATTACKS

1. [Industry Benchmarking]: No industry solutions, open-source projects, or published patterns cited — quote: "这是典型的补全遗漏的 adapter 层问题...类比：Airflow 添加新 DAG 类型后需要配套的 executor plugin 和 UI renderer" — this single analogy is insufficient. Must reference at least one concrete industry-validated approach (convention-over-configuration frameworks like Rails/Spring, plugin systems like VS Code extensions, type registration patterns).

2. [Industry Benchmarking]: Three alternatives are scope variants of the same approach — quote: "| 手动补全模板 + 最小化修复 |...| 仅修 P0+P1，保留 gen-and-run |...| 完整修复 P0+P1+P2 |" — these differ only in scope, not in architectural strategy. Must include genuinely different approaches (e.g., init-time schema validation, code generation from type definitions, migration tool for existing features).

3. [Industry Benchmarking]: Trade-off analysis understates selected approach's risk — quote: "变更量较大（但大部分是机械性清理）" — the parenthesis immediately minimizes the con. With 14 test files and ~95 references, this is non-trivial risk that deserves honest weighting.

4. [Solution Creativity]: Proposal explicitly admits it's gap-filling, not innovation — quote: "Task 1 引入的自动发现机制已消除根因。本次提案聚焦于补全遗漏的执行层面配套" — no novelty, no cross-domain inspiration. This is a cleanup task framed as a proposal.

5. [Requirements Completeness]: NFR gap — no performance consideration for CategoryEval branch addition in hot-path functions (`RenderRecord`, `validateRecordData`). No compatibility matrix specifying which Forge versions' index.json files are affected.

6. [Requirements Completeness]: Missing scenario: eval template content must match actual eval skill behavior — quote: "eval 模板参考 validation-code.md（评估 + pass/fail 判定模式，~70 行）" — the proposal does not verify that `validation-code.md`'s pattern is semantically correct for eval tasks, only that it exists.

7. [Risk Assessment]: Missing risk: features currently using gen-and-run type will encounter file-not-found error (not migration-aware error) when `prompt/data/test-gen-and-run.md` is deleted. Quote: "向后兼容: 旧 index.json 引用 test.gen-and-run 时给出明确迁移错误提示" — this only covers `validate_index.go`, not the `Synthesize()` code path.

8. [Risk Assessment]: Mitigation "test-gen 模板参考现有 test-gen-scripts.md" is implementation detail, not risk mitigation — quote: "test-gen 模板参考现有 test-gen-scripts.md（skill 委托模式）" — a real mitigation would specify review/validation of template content correctness.

9. [Success Criteria]: "无顺序耦合" is not testable — quote: "ResolveFirstTestDep + T-review-doc prepend 为单步操作，无顺序耦合" — how do you write a test that verifies "no ordering coupling"? Must replace with concrete behavioral test (e.g., "calling the function with needsEval=false produces identical output to calling with needsEval=true after removing T-review-doc from results").

10. [Success Criteria]: Missing criterion for `record-format-eval.md` — this file is in Scope (P1) but no success criterion verifies its existence or content correctness. Must add: "`plugins/forge/skills/submit-task/data/record-format-eval.md` exists and contains score/findings/severity/passed field definitions".

11. [Success Criteria]: Missing criterion for `CategoryForType` default branch `log.Printf` — quote: "default 分支改为返回 CategoryCoding 并通过 log.Printf(...) 记录警告" — no success criterion verifies this warning is emitted for unknown types.

12. [Logical Consistency]: findFirstTestTaskIdx fix contradicts "auto-discovery eliminates root cause" — quote from Innovation Highlights: "自动发现机制已消除'忘记更新映射'的根因" vs. Scope P1: "替换为 findTaskIndexByPrefix(tasks, 'T-test-gen-journeys')" — replacing one hardcoded prefix with another does not use auto-discovery and retains the same failure mode class.

13. [blindspot]: `--target 850` in eval template is an unexplained hardcoded value — quote: "Skill(skill='forge:eval', args='--type [journey|contract] --target 850')" — the proposal does not explain why 850 is the target, whether it should be parameterized, or what happens when different eval tasks need different targets.

14. [blindspot]: P1 item 3 and P2 modify the same code location (`build.go:484,492-494`) — the proposal does not flag this overlap, creating implementation conflict risk.

15. [blindspot]: Integration test promised in risk mitigation but absent from task count and success criteria — quote: "更新后添加集成测试验证 Quick mode 依赖链正确性" — no task, no criterion.

16. [blindspot]: Deleting `prompt/data/test-gen-and-run.md` causes `Synthesize()` file-not-found for existing features, bypassing the migration-aware error path in `validate_index.go`.
