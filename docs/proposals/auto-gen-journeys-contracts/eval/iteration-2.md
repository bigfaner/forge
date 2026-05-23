# Proposal Evaluation Report — Iteration 2

**Document**: `docs/proposals/auto-gen-journeys-contracts/proposal.md`
**Date**: 2026-05-23
**Scorer**: CTO-level adversarial review

---

## Pre-Score Anchors (Phase 1 — Independent Reasoning Audit)

1. **Problem -> Solution link**: Sound. Making gen-journeys/gen-contracts auto-generated tasks directly addresses the pipeline gap. The core argument chain is valid.

2. **Iteration-1 critical fixes addressed**: The proposal now accurately describes the hardcoded indices in `resolveBreakdownDeps` (lines 426-427, `evalJourneyIdx:=0`, `evalContractIdx:=1`). It correctly proposes a full rewrite using `findTaskIndexByPrefix` instead of the previous false claim that the code already uses that approach. The `ResolveFirstTestDep` Quick branch issue is now specified (`T-quick-gen-journeys` prefix). Both SKILL.md semantic conflicts now have explicit mitigation proposals (`AUTO_COMMIT=true`, `SKIP_EVAL_GATE=true` template injection).

3. **Remaining gap — resolveQuickDeps specificity**: While the proposal now acknowledges "resolveQuickDeps 需要完全重写" and explains why (3 tasks per profile type vs 1 gen-and-run), it does not specify the new dependency topology within Quick mode. Are the 3 tasks per profile type sequential (gen-journeys -> gen-contracts -> gen-scripts per type) or staged (all gen-journeys, then all gen-contracts, then all gen-scripts)? This is an architectural decision that affects correctness.

4. **Template injection mechanism under-specified**: `AUTO_COMMIT=true` and `SKIP_EVAL_GATE=true` are proposed as template variables, but SKILL.md files are markdown prompts — they don't have a variable parsing mechanism. The execution model for reading these injected values is unclear.

5. **Scope estimate tension**: 14 in-scope deliverables mapped to "5-8 coding tasks" with no breakdown or justification. Given two full function rewrites plus SKILL.md modifications, this estimate lacks supporting evidence.

---

## Rubric Scoring (Phase 2)

### Dimension 1: Problem Definition — 85/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 35/40 | Core problem is clear: gen-journeys/gen-contracts require manual invocation, creating a pipeline gap. Minor ambiguity: "自动生成的测试流水线任务" vs "手动调用的 skills" — the distinction between auto-generated task types and skill invocation could benefit from a one-sentence clarification for readers unfamiliar with Forge's autogen system. |
| Evidence provided | 25/40 | Three evidence points: code structure analysis (`autogen.go` `GetBreakdownTestTasks`), user behavior (manual `/gen-journeys` calls), and Quick mode quality risk. Still no user feedback, no issue tracking references, no concrete failure scenarios with reproduction steps. The evidence remains "we inspected the code" rather than "this caused user-visible harm X times." |
| Urgency justified | 25/30 | v3.0.0 dependency is stated. Test profile system as prerequisite is mentioned. Cost of delay is implied (pipeline cannot run end-to-end) but not quantified. What specifically is blocked in the release without this? |

**Deductions**: -5 for minor ambiguity in problem framing; -15 for evidence limited to code inspection with no user-facing impact data; -5 for unquantified delay cost.

### Dimension 2: Solution Clarity — 96/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 38/40 | Pipeline diagrams for both modes are clear. Task type names, template files, and mode-specific behavior are specified. A reader can explain back what will be built. |
| User-facing behavior described | 30/45 | `forge task index` output described (new tasks appear). But the Quick mode UX change is underspecified: replacing 1 gen-and-run task per profile type with 3 split tasks means the task list grows by 2x per type. The document says "替换" but doesn't describe what the user sees differently — task count change, task ordering, or naming conventions in the CLI output. |
| Technical direction clear | 28/35 | Significant improvement over iteration 1. `resolveQuickDeps` full rewrite is now explicitly stated with rationale. `findTaskIndexByPrefix` named as the replacement approach. `ResolveFirstTestDep` adaptation specified (`T-quick-gen-journeys` prefix). Gap: the new dependency topology within Quick mode is not specified — sequential per-type or staged across types? This affects how `resolveQuickDeps` must be implemented. |

**Deductions**: -2 for approach missing the Quick mode UX detail; -15 for user-facing behavior not describing Quick mode task list changes; -7 for unspecified Quick mode intra-task dependency topology.

### Dimension 3: Industry Benchmarking — 78/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 30/40 | Now cites specific products: Cypress Studio, Playwright Codegen, Selenium IDE (record-replay); Testim, Mabl, Reflection (AI-assisted); Graphwalker, Cucumber/SpecFlow (model-based/BDD). Significant improvement over iteration 1. However, citations are name-plus-one-sentence only — no links, no versions, no published patterns or papers. |
| At least 3 meaningful alternatives | 18/30 | Three alternatives listed. "Merge into gen-and-run" is now framed with a Testim analogy ("类似 Testim 单一管道"), which elevates it above straw-man. "Do nothing" is genuine. "Independent auto-tasks" is the proposal. However, none of the alternatives represent adopting an industry-validated solution verbatim — all are Forge-internal architectural choices. |
| Honest trade-off comparison | 18/25 | Selected approach includes a quantitative cost estimate: "开发量增加约 30%". Pros/cons relate to actual project constraints. However, the claim "与 Graphwalker 等工具的 stage-per-artifact 模式一致" is a questionable analogy — Graphwalker is a graph-based test path generator, not a pipeline framework. The "stage-per-artifact" pattern is generic enough to apply to any multi-stage system. |
| Chosen approach justified against benchmarks | 12/25 | Forge's unique positioning (PRD -> Journey -> Contract -> Test Script full chain) is articulated and compared to industry. The "only Testim has a comparable end-to-end pipeline" claim is interesting but unsubstantiated. Why is this specific task-splitting approach better than other structures? The justification is "most aligned with Forge design philosophy" — a subjective criterion, not benchmarked against alternatives. |

**Deductions**: -10 for shallow citations; -12 for no industry-validated solution as a genuine alternative; -7 for questionable Graphwalker analogy; -13 for subjective justification criterion.

### Dimension 4: Requirements Completeness — 80/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 30/40 | Four key scenarios listed. The proposal.md quality degradation boundary is now partially addressed (in Assumptions Challenged). Missing: (1) Error scenarios during auto-task execution — what happens if template injection fails? (2) Migration scenario — what about existing features with gen-and-run tasks in their index.json? (3) What if gen-journeys produces zero Journey files in Quick mode? |
| Non-functional requirements | 22/40 | Still minimal: MainSession=false for both tasks, embed.FS consistency. Missing: (1) Execution time impact — gen-journeys can take 10-30 minutes as an auto-task blocking the pipeline. (2) Token/resource consumption for auto-generated tasks. (3) Error handling behavior when auto-tasks fail mid-pipeline. (4) Backward compatibility verification for the template injection mechanism. |
| Constraints & dependencies | 28/30 | Now explicitly addresses the gen-contracts eval-journey Blocker conflict with a concrete solution (SKIP_EVAL_GATE=true template injection). gen-journeys HARD-RULE interactive conflict addressed via AUTO_COMMIT=true. Dependencies listed (PRD/proposal.md input, test profile system). Well done. |

**Deductions**: -10 for missing edge/error cases; -18 for thin NFRs; -2 for minor dependency detail.

### Dimension 5: Solution Creativity — 41/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 18/40 | The document is honest: "无特殊创新。对已有自动化框架的自然延伸." The differentiation section describes Forge's unique chain but that's a positioning statement about the existing system, not innovation in this proposal. The proposal connects two existing systems. |
| Cross-domain inspiration | 8/35 | BDD pipeline stage-splitting and Graphwalker's stage-per-artifact pattern are cited. These are within-domain references (testing), not cross-domain inspiration. No borrowing from unrelated domains. |
| Simplicity of insight | 15/25 | The core insight (add two task types to an existing framework) is simple, which is good. The template injection mechanism (AUTO_COMMIT, SKIP_EVAL_GATE) is a minimally invasive solution to the SKILL.md adaptation problem — elegant. But the overall implementation is not as simple as presented given the two full function rewrites. |

**Deductions**: -22 for no meaningful innovation (self-admitted); -27 for no cross-domain inspiration; -10 for implementation simplicity overstated.

### Dimension 6: Feasibility — 75/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 32/40 | Now accurately describes the resolveBreakdownDeps issue (hardcoded indices 0 and 1). Proposes a complete rewrite using findTaskIndexByPrefix. The ResolveFirstTestDep Quick branch fix is specified. Remaining gap: the resolveBreakdownDeps rewrite is complex — the function contains index arithmetic for multiple downstream tasks, and the full scope of what "rewrite" entails is not detailed. |
| Resource & timeline | 18/30 | "5-8 coding tasks" stated but not justified. With 14 in-scope deliverables including two full function rewrites (resolveBreakdownDeps, resolveQuickDeps), two SKILL.md modifications, two new task types with templates, and documentation updates, the estimate may be tight. No evidence or task breakdown provided. |
| Dependency readiness | 25/30 | SKILL.md interactive conflicts now resolved via template injection proposal. Test profile system merged. gen-journeys/gen-contracts SKILL.md stated as stable. However, the template injection approach (AUTO_COMMIT=true, SKIP_EVAL_GATE=true) requires SKILL.md modifications to recognize these flags — this work is mentioned but not detailed in terms of readiness. |

**Deductions**: -8 for incomplete technical detail on rewrite scope; -12 for unjustified timeline estimate; -5 for unspecified SKILL.md flag parsing readiness.

### Dimension 7: Scope Definition — 69/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 28/30 | 14 concrete deliverables. Most specify file names, function names, or task type names. "更新 ARCHITECTURE.md 和相关文档" remains vague — what related documents? |
| Out-of-scope explicitly listed | 23/25 | Five out-of-scope items including a disposal strategy for the test-gen-and-run.md template. Good boundary setting. |
| Scope is bounded | 18/25 | "5-8 coding tasks" provides a numerical bound. Scope now includes resolveBreakdownDeps rewrite, resolveQuickDeps rewrite, and SKILL.md adaptations — items that were missing in iteration 1. However, 14 in-scope items for 5-8 tasks averages ~2 deliverables per task, which is aggressive given the complexity of the two full function rewrites. |

**Deductions**: -2 for vague documentation item; -2 for minor out-of-scope gap; -7 for scope-to-task ratio tension.

### Dimension 8: Risk Assessment — 72/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 25/30 | Four risks, including the ResolveFirstTestDep Quick branch mismatch (an iteration-1 blindspot). Missing: (1) Risk of regression during resolveBreakdownDeps full rewrite — this is a core function affecting all Breakdown mode pipelines. (2) Risk of resolveQuickDeps rewrite miscalculating the new multi-task topology. The resolveQuickDeps rewrite is the most architecturally significant change but has no dedicated risk item. |
| Likelihood + impact rated | 22/30 | Ratings are honest: "Breakdown 索引偏移" at H/H with accurate description of lines 426-427. "ResolveFirstTestDep 不匹配" at H/H with specific code reference. "proposal.md 质量" at M/M with clear degradation strategy. Gap: resolveQuickDeps rewrite — the most significant architectural change — is not listed as a standalone risk. |
| Mitigations are actionable | 25/30 | Mitigations are concrete and executable: "重写为基于 findTaskIndexByPrefix" (executable), "更新为 T-quick-gen-journeys 前缀" (executable), "定义降级边界" with specific abort/degrade conditions (executable). Gap: the resolveQuickDeps rewrite is mentioned in the Solution section but not linked to any risk item with a corresponding mitigation. |

**Deductions**: -5 for missing regression risk; -3 for missing resolveQuickDeps-specific risk; -8 for unassigned mitigation; -5 for resolveQuickDeps rewrite having no risk item.

### Dimension 9: Success Criteria — 63/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 45/55 | Six criteria, mostly testable. The iteration-1 typo (`test.test-gen-contracts`) is now correctly stated as `T-test-gen-journeys` and `T-test-gen-contracts`. Gap: Criterion 4 ("gen-journeys 支持 proposal.md 作为输入") lacks a testable output — how do we verify it processed proposal.md correctly? Criterion 3 ("validate-index 通过") is regression testing, not functional correctness. |
| Coverage is complete | 18/25 | Covers main scenarios. Missing: (1) No success criterion verifies dependency chain correctness after resolveBreakdownDeps rewrite. (2) No criterion verifies Quick mode pipeline produces the correct task count/order. (3) No criterion verifies AUTO_COMMIT or SKIP_EVAL_GATE injection works. (4) No criterion verifies backward compatibility (existing index.json unaffected). With 14 in-scope items and 6 success criteria, coverage ratio is low (43%). |

**Deductions**: -10 for untestable criterion 4 and regression-only criterion 3; -7 for coverage gaps across multiple in-scope items.

### Dimension 10: Logical Consistency — 66/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses stated problem | 30/35 | Yes, making gen-journeys/gen-contracts auto-generated tasks directly solves the pipeline gap. The Quick mode gen-and-run replacement is scope expansion beyond the stated problem (which is about automation, not restructuring), but the document now explicitly frames this as a necessary co-change (line 37-39). |
| Scope <-> Solution <-> Success Criteria aligned | 18/30 | Better than iteration 1. SKILL.md adaptations are now in scope with corresponding solution elements (template injection). However: (1) "重写 resolveQuickDeps" is in scope but has no success criterion. (2) "修改 gen-contracts SKILL.md" is in scope but has no success criterion. (3) "修改 gen-journeys SKILL.md 支持非交互模式" is in scope but has no success criterion. 14 in-scope items vs 6 success criteria leaves 8 items unverified. |
| Requirements <-> Solution coherent | 18/25 | Requirements now map more clearly to solution elements. SKILL.md conflicts (identified in requirements) have corresponding solutions (template injection). However, the "向后兼容" requirement (scenario 4) has no corresponding success criterion and no concrete solution element verifying that existing index.json files are unaffected. |

**Deductions**: -5 for scope expansion beyond problem statement; -12 for scope/criteria misalignment (8 of 14 items uncovered); -7 for orphan backward-compatibility requirement.

---

## Cross-Dimension Coherence Check

The most improved area is the **accuracy of technical claims**. The proposal now correctly identifies the hardcoded indices in `resolveBreakdownDeps`, proposes a concrete rewrite strategy, and addresses both SKILL.md semantic conflicts. This eliminates the cascading factual-error problem from iteration 1 that affected Dimensions 2, 6, 8, and 10.

The remaining cross-cutting issue is **scope-to-evidence ratio**. The proposal claims "5-8 coding tasks" for 14 deliverables including two full function rewrites, but provides no breakdown or justification. This affects Dimensions 6, 7, and 10.

---

## Phase 3 — Blindspot Hunt

### [blindspot-1] resolveQuickDeps new dependency topology unspecified
The proposal states "resolveQuickDeps 需要完全重写" and explains that each profile type now produces 3 tasks instead of 1. But it does not specify the new topology: are tasks sequential per type (`T-quick-gen-journeys-{type}` -> `T-quick-gen-contracts-{type}` -> `T-quick-gen-scripts-{type}` for each type) or staged across types (all gen-journeys first, then all gen-contracts)? This architectural decision determines correctness of the rewrite. Quote: "每个 profile type 产生 3 个任务（gen-journeys + gen-contracts + gen-scripts），加上 run 和 verify，索引布局完全不同。"

### [blindspot-2] Template injection execution model unclear
`AUTO_COMMIT=true` and `SKIP_EVAL_GATE=true` are proposed as template-injected variables. But SKILL.md files are markdown prompt files — they don't have a variable parsing mechanism. The execution model for how these values are read and acted upon is not explained. The proposal assumes SKILL.md files can conditionally alter behavior based on injected template context, but this capability is not demonstrated as existing in the current architecture. Quote: "在任务模板中注入 `AUTO_COMMIT=true` 标记，gen-journeys 据此跳过人工审批直接提交生成的 Journey 文件。"

### [blindspot-3] resolveQuickDeps rewrite has no dedicated risk item
The scope explicitly includes "重写 resolveQuickDeps（适配新的多任务拓扑）" and the solution section describes it as a full rewrite. This is the most architecturally significant change in the proposal (replacing the entire Quick mode task topology). Yet the risk table has 4 items and none of them is "resolveQuickDeps rewrite introduces regression or miscalculation." The resolveBreakdownDeps rewrite has a risk item (row 2), but resolveQuickDeps does not. Quote from Key Risks table: no entry for resolveQuickDeps.

### [blindspot-4] "5-8 coding tasks" estimate lacks evidence
14 in-scope deliverables are mapped to "5-8 coding tasks" with no task breakdown or justification. Given the complexity of two full function rewrites (each involving index arithmetic across multiple test profile types) plus two SKILL.md modifications and two new template files, the estimate should include a brief task breakdown. Quote: "预计 5-8 个 coding 任务，在 quick mode 范围内。"

### [blindspot-5] Failure scenarios absent from Key Scenarios
The Requirements Analysis Key Scenarios section lists 4 scenarios, all happy-path or nominal-path. No scenario covers: (1) auto-task execution failure mid-pipeline, (2) gen-journeys producing zero Journey files, (3) template variable injection failure. The Assumptions Challenged section partially addresses proposal.md quality degradation, but error/failure scenarios should be in the requirements section. Quote from Key Scenarios: four scenarios, all describing successful completion.

### [blindspot-6] Quick mode quality relies on downstream compensation without mechanism
The proposal claims proposal.md content "足以提取精简版 Journeys" and that "质量由下游 gen-contracts 和 gen-test-scripts 的代码侦察来补偿." This inverts the quality model: instead of upstream verification (eval), it relies on downstream correction. But the mechanism for how "代码侦察 compensates" is completely unexplained. If proposal.md describes a critical scenario not in the codebase yet, code reconnaissance cannot discover it. The degradation boundary ("scope + success criteria 为硬性前提") mitigates this partially but doesn't address the fundamental gap. Quote: "质量由下游 gen-contracts 和 gen-test-scripts 的代码侦察来补偿。"

---

## Freeform Finding Integration

| Finding | Rubric Dimension | Resolution |
|---------|-----------------|------------|
| [high] resolveBreakdownDeps hardcoded indices | Dim 2, 6, 8, 10 | **Addressed**: Proposal now accurately describes lines 426-427 hardcoded indices and proposes full rewrite. No longer a deduction driver. |
| [high] resolveQuickDeps rewrite needed | Dim 2, 6, 8 | **Partially addressed**: Proposal now explicitly states full rewrite needed. But topology unspecified (blindspot-1) and no risk item (blindspot-3). Moderate deductions remain. |
| [high] gen-contracts eval-journey Blocker conflict | Dim 4, 10 | **Addressed**: SKIP_EVAL_GATE=true template injection proposed as concrete solution. Constraint now in both requirements and scope. |
| [medium] gen-journeys HARD-RULE interactive conflict | Dim 4, 6 | **Addressed**: AUTO_COMMIT=true template injection proposed. Listed in scope. |
| [medium] Quick mode proposal.md input quality | Dim 4 | **Partially addressed**: Degradation boundary defined (scope + success criteria as hard prerequisites, key scenarios as optional with smoke-level fallback). But downstream compensation mechanism unexplained (blindspot-6). |
| [medium] ResolveFirstTestDep Quick branch | Dim 2, 6, 8 | **Addressed**: Explicitly stated in scope with specific new prefix `T-quick-gen-journeys`. Listed as risk item with H/H rating. |

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 85 | 110 |
| 2. Solution Clarity | 96 | 120 |
| 3. Industry Benchmarking | 78 | 120 |
| 4. Requirements Completeness | 80 | 110 |
| 5. Solution Creativity | 41 | 100 |
| 6. Feasibility | 75 | 100 |
| 7. Scope Definition | 69 | 80 |
| 8. Risk Assessment | 72 | 90 |
| 9. Success Criteria | 63 | 80 |
| 10. Logical Consistency | 66 | 90 |
| **Total** | **725** | **1000** |
