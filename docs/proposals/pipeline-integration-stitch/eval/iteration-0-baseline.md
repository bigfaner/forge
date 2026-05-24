---
created: 2026-05-24
scorer: CTO-adversary
iteration: 0
status: baseline
---

# Eval Report: Pipeline Integration Stitch (Iteration 0 — Baseline)

## Scoring Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 103 | 110 |
| 2. Solution Clarity | 108 | 120 |
| 3. Industry Benchmarking | 60 | 120 |
| 4. Requirements Completeness | 92 | 110 |
| 5. Solution Creativity | 50 | 100 |
| 6. Feasibility | 95 | 100 |
| 7. Scope Definition | 75 | 80 |
| 8. Risk Assessment | 78 | 90 |
| 9. Success Criteria | 72 | 80 |
| 10. Logical Consistency | 85 | 90 |
| **Total** | **818** | **1000** |

---

## Phase 1: Reasoning Audit

### Problem -> Solution

The problem chain is tight: the `auto-gen-journeys-contracts` proposal added new types (gen-journeys, gen-contracts, eval.journey, eval.contract) to the type registry, but missed the execution-layer adapters. The three-pronged solution (templates + CategoryEval + gen-and-run removal) directly addresses each identified gap. No phantom problems detected — all 6 evidence items were verified against the codebase:

- **P0-1 confirmed**: `prompt/data/` is missing `test-gen-journeys.md`, `test-gen-contracts.md`, `eval-journey.md`, `eval-contract.md`. The `templatePath()` convention (`typeName` with `.` -> `-` + `.md`) means `Synthesize()` will attempt to `ReadFile` a non-existent path.
- **P1-2 confirmed**: `category.go` has no `eval.` prefix branch; `CategoryForType("eval.journey")` falls through to `default` -> `CategoryCoding`.
- **P1-3 confirmed**: `build.go:492-494` matches `T-quick-gen-and-run*`.
- **P1-4 confirmed**: `build.go:329-347` shows the two-step operation.
- **P2-5 confirmed**: `TypeTestGenAndRun`, `test-gen-and-run.md`, `genScriptBases` entry, `validate_index.go:224` check all present.
- **P2-6 confirmed**: `record-format-test.md` lists `test.gen-cases`, `test.eval-cases`, `test.gen-and-run` — all deprecated. No `record-format-eval.md` exists.

### Solution -> Evidence

Evidence is code-level and specific. File paths, line numbers, function names are all accurate. This is a strong evidence chain.

### Evidence -> Success Criteria

Success criteria are testable and map to the evidence items. Some gaps noted below.

### Self-contradiction check

One minor tension: the proposal claims P1-4 is a "顺序耦合" (ordering coupling) risk, but also admits "当前幂等" (currently idempotent). The Assumptions Challenged table correctly flags this as "Overturned: 当前代码已幂等... 风险仅为代码耦合". This is honest but slightly inflates the severity from P1 to "nice-to-have refactor".

---

## Phase 2: Rubric Scoring

### 1. Problem Definition (103/110)

**Problem stated clearly (38/40)**: The problem is unambiguous — a prior proposal left execution-layer gaps. The P0/P1/P2 priority scheme is effective. One reader might interpret "Pipeline 执行必定失败" as hyperbolic since it only affects 4 specific type strings, not the entire pipeline. Minor deduction for lack of clarity on whether "all features using test pipeline" means "all features that have test pipeline tasks in their index" or "all features period".

**Evidence provided (40/40)**: Excellent. Every claim is backed by specific file paths, line numbers, and observable behavior. Verified against codebase — all claims check out.

**Urgency justified (25/30)**: P0 is convincingly argued. However, the cost of delay is stated only for P0. No urgency argument for P1 or P2. The proposal bundles P1+P2 into the solution without explaining why they cannot be deferred.

### 2. Solution Clarity (108/120)

**Approach is concrete (40/40)**: The three-pronged approach is crystal clear. A reader could explain it back immediately: create 4 templates, add CategoryEval, remove gen-and-run.

**User-facing behavior described (38/45)**: The proposal describes what `forge submit-task` and `Synthesize()` will do, which are developer-facing CLI behaviors. However, it does not describe what the *end user* (agent executing a task) will experience differently. For instance: what does the eval prompt template tell the agent to do? What fields does the agent need to fill in? The "user" here is the AI agent, and its experience is only partially specified.

**Technical direction clear (30/35)**: Sufficient for implementation. References to existing patterns (CategoryTest, code-quality-simplify.md) give concrete anchors. However, the eval record template fields are listed but not designed — `evalScore`, `evalFindings`, `evalSeverity`, `evalPassed` are named but their types, validation rules, and relationships are unspecified.

### 3. Industry Benchmarking (60/120)

**Industry solutions referenced (15/40)**: The proposal contains one analogy: "类比：Airflow 添加新 DAG 类型后需要配套的 executor plugin 和 UI renderer". This is a single analogy, not an industry benchmark with real-world references. No product names, open-source projects, or published patterns are cited beyond this. The rubric asks for "real-world solutions/patterns for this type of problem" — the Airflow mention is superficial and does not constitute research into how established systems handle type registration vs. execution adapter synchronization.

**At least 3 meaningful alternatives (20/30)**: Three alternatives are presented in the comparison table: (1) minimal fix only, (2) fix P0+P1 keep gen-and-run, (3) full fix. However, alternatives 1 and 2 are effectively "do less" variants, not genuinely different approaches. A meaningful alternative could be: lazy template generation (generate templates on first use), or a validation gate that catches missing templates at index-load time. None of these are explored. At least the "do nothing" is implicitly covered by alternative 1 (which is rejected).

**Honest trade-off comparison (15/25)**: The cons are somewhat cherry-picked. For the selected approach, the con is "变更量较大（但大部分是机械性清理）" — the parenthetical immediately downplays the con. For alternative 1, the con is "eval 分类仍错" which is accurate but doesn't explore whether the eval misclassification actually causes failures in practice (the proposal notes it only affects submit validation, which may not be in the critical path for most features).

**Chosen approach justified against benchmarks (10/25)**: No benchmarking was performed. The Airflow analogy was not used to derive a design decision. The choice is justified purely by internal logic (fix everything > fix partially > fix minimally), not by comparison to industry practices.

### 4. Requirements Completeness (92/110)

**Scenario coverage (35/40)**: Key scenarios are well-identified. Happy path (new types work), mixed feature dependencies, quick-mode routing, and gen-and-run removal are covered. Edge cases partially covered (old index.json with deprecated types). Missing: what happens if a feature has only eval tasks but no test tasks? What happens if both eval and test tasks submit simultaneously?

**Non-functional requirements (35/40)**: Backward compatibility is mentioned (migration error message). Test coverage requirements are specified. Missing: no performance NFR (will the extra category check add latency? negligible, but not stated). No security NFR (is there a risk of eval submissions being used to bypass test evidence requirements?).

**Constraints & dependencies (22/30)**: Only one dependency is listed: "Task 1 已完成". Missing constraints: Go version compatibility, CLI version compatibility for users with older forge-cli, timing constraints relative to the current sprint.

### 5. Solution Creativity (50/100)

**Novelty over industry baseline (20/40)**: The proposal acknowledges in Innovation Highlights that "Task 1 已引入的自动发现机制已消除根因" and "本次提案聚焦于补全". This is explicitly a catch-up/patch proposal, not an innovation. The proposal is honest about this but scores low on novelty by definition.

**Cross-domain inspiration (10/35)**: No cross-domain inspiration is demonstrated. The Airflow analogy is mentioned once but not explored or leveraged for design insights.

**Simplicity of insight (20/25)**: The "just create the missing files and fix the category" insight is elegant in its simplicity. The proposal correctly identifies that the root cause was a prior proposal's incompleteness, not a design flaw.

### 6. Feasibility (95/100)

**Technical feasibility (38/40)**: All technical steps are verified feasible against the codebase. The template files follow existing conventions. CategoryEval follows the CategoryTest pattern. One minor concern: `RecordData` struct changes (adding eval fields) could break existing JSON deserialization if not handled carefully — the proposal doesn't mention backward compatibility for the struct itself.

**Resource & timeline feasibility (28/30)**: "3 coding + 2 doc tasks, ~6h" is specific and realistic for the described scope. The gen-and-run removal affecting 14 test files with ~95 references is the main risk to the estimate, but the proposal correctly identifies this as "mechanical".

**Dependency readiness (29/30)**: Task 1 is complete. No external dependencies. Strong.

### 7. Scope Definition (75/80)

**In-scope items are concrete (28/30)**: Each P0/P1/P2 item maps to specific files and specific changes. The file-level granularity is excellent.

**Out-of-scope explicitly listed (22/25)**: Five out-of-scope items are named. Good. However, one out-of-scope item ("旧 index.json 自动迁移工具") is actually partially in scope (validate_index.go will return a migration error), creating a gray area.

**Scope is bounded (25/25)**: Clear deliverables, clear timeline, clear boundary. The proposal can be executed within the stated timeframe.

### 8. Risk Assessment (78/90)

**Risks identified (25/30)**: Five risks are identified. Missing risks: (1) what if the eval record template design is wrong and needs to change after tasks are submitted? (no backward compatibility plan for stored records); (2) what if gen-and-run removal breaks features that are currently in-progress with gen-and-run tasks?

**Likelihood + impact rated (25/30)**: Ratings are reasonable. However, the first risk ("4 个新 prompt 模板内容不准确") is rated M/H — given that these are instruction templates for AI agents, "inaccurate" is vague. What constitutes accuracy for a prompt template? This should be more precisely defined.

**Mitigations are actionable (28/30)**: Mitigations are concrete: "参考现有模板结构", "编写单元测试覆盖正向/负向/边界用例", "每项后执行 go build 验证". These are actionable. Minor gap: no mitigation for the risk of incomplete gen-and-run removal (the "grep 返回零结果" success criterion is itself the mitigation, but no incremental check strategy is described).

### 9. Success Criteria (72/80)

**Criteria are measurable and testable (50/55)**: Most criteria are excellent: "grep 返回零结果", "CategoryForType returns X", "Synthesize returns valid prompt". However:
- "ResolveFirstTestDep + T-review-doc prepend 为单步操作，无顺序耦合" — how do you test "无顺序耦合"? This is a design property, not an observable behavior. A better criterion would be: "reordering the calls does not change the output" or "single function call encapsulates both operations".
- "所有现有测试通过" — this is a necessary condition, not a success criterion. It should be a prerequisite, not a deliverable.

**Coverage is complete (22/25)**: Most in-scope items have corresponding success criteria. Missing: no success criterion for `record-format-eval.md` creation (it's listed in scope but not checked). No criterion for the "向后兼容" NFR (old index.json migration error).

### 10. Logical Consistency (85/90)

**Solution addresses the stated problem (33/35)**: Tight mapping. P0 -> create templates. P1 -> add CategoryEval + harden deps. P2 -> remove gen-and-run. Each problem item has a corresponding solution item. One gap: P1-4 (顺序耦合) is presented as a bug but the proposal itself admits it's not actually a bug (Assumptions Challenged table says "Overturned"). This is honest but creates a logical inconsistency: why is a non-bug listed as P1 evidence?

**Scope <-> Solution <-> Success Criteria aligned (27/30)**: Generally well-aligned. One gap: the Scope section mentions updating `record-format-test.md` as P1 work, but the Success Criteria section does not verify this update. The criterion `grep -r "gen-and-run" forge-cli/` would not catch an outdated `record-format-test.md` in `plugins/forge/`.

**Requirements <-> Solution coherent (25/25)**: Requirements map cleanly to solution items. No orphan requirements or solution features without requirements.

---

## Phase 3: Blindspot Hunt

### [blindspot-1] Eval record immutability not addressed

The proposal adds new eval-specific fields to `RecordData` (`evalScore`, `evalFindings`, `evalSeverity`, `evalPassed`) but does not discuss what happens when existing submitted records (stored as markdown files) are read back. If the record format changes, old records may fail to parse or display correctly. This is a data migration concern that the rubric's Risk Assessment dimension should have caught.

**Quote**: "types.go: RecordData 结构添加 eval 特有字段（evalScore、evalFindings、evalSeverity、evalPassed）"

### [blindspot-2] Prompt template content is entirely unspecified

The proposal identifies 4 missing prompt templates as P0 (highest priority) and lists them as in-scope deliverables, but provides zero guidance on what these templates should contain beyond "任务上下文说明、输入格式、期望输出格式、质量标准" (a generic four-part structure). For P0 items that block all test pipeline execution, this is a significant design gap. The templates are the actual user-facing deliverable, and their content determines whether the eval quality gate works correctly.

**Quote**: "每个模板包含：任务上下文说明、输入格式（从 index.json task configuration 读取）、期望输出格式、质量标准"

### [blindspot-3] No rollback plan

The proposal removes `TypeTestGenAndRun` from the type system entirely, which means any feature currently using gen-and-run tasks in their index.json will break. While `validate_index.go` will return a migration error, there is no rollback plan if the migration error message is insufficient or if users need to continue using gen-and-run during a transition period. The "Out of Scope" explicitly excludes "旧 index.json 自动迁移工具", leaving users without a migration path.

**Quote**: "旧 index.json 自动迁移工具" (listed as Out of Scope)

### [blindspot-4] Assumptions Challenged section is unusual placement

The "Assumptions Challenged" section is a meta-cognitive artifact (it documents the proposal author's own thought process) rather than a standard proposal section. While intellectually honest, it doesn't map to any rubric dimension and creates ambiguity about whether the "overturned" assumptions should still be listed as evidence (they are). P1-4 remains in the Evidence section despite being acknowledged as non-bug in Assumptions Challenged.

**Quote**: "re-index 幂等性是实际 bug | Code Audit | Overturned: 当前代码已幂等"

### [blindspot-5] Test count estimate may understate complexity

The proposal estimates "14 个测试文件中 ~95 处引用" for gen-and-run cleanup. However, the grep results show gen-and-run references in 108 files across the entire codebase. While most of these are in historical docs (correctly scoped out), the boundary between "active docs" and "historical docs" is not clearly defined. OVERVIEW.md and task-lifecycle.md are mentioned, but other potentially active docs (like expert files, convention docs) are not enumerated.

**Quote**: "活跃文档：更新 OVERVIEW.md、task-lifecycle.md 中的 gen-and-run 引用"

---

## Improvement Priority

1. **[Critical] Industry Benchmarking (60/120)**: Research how established task pipeline systems (Airflow, Temporal, Prefect, Jenkins) handle type registration vs. execution adapter synchronization. Cite specific patterns (e.g., Airflow's "plugin discovery" mechanism, Temporal's "activity interface registration").
2. **[Important] Prompt template design**: Specify the actual content or at least the key design decisions for the 4 P0 templates. What instructions do eval agents receive? What output format do they produce?
3. **[Important] Success criteria gaps**: Add criteria for `record-format-eval.md` creation and old index.json migration error behavior.
4. **[Minor] P1-4 reclassification**: Move the ordering coupling from P1 evidence to a "nice-to-have refactoring" note, since it's acknowledged as non-bug.
