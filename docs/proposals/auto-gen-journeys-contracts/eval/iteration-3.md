# Proposal Evaluation Report — Iteration 3 (Final)

**Document**: `docs/proposals/auto-gen-journeys-contracts/proposal.md`
**Date**: 2026-05-23
**Scorer**: CTO-level adversarial review

---

## Pre-Score Anchors (Phase 1 — Independent Reasoning Audit)

1. **Problem -> Solution link**: Sound. Making gen-journeys/gen-contracts auto-generated tasks directly addresses the pipeline gap. The core argument chain is valid.

2. **Iteration-2 critical gaps addressed**: The proposal now specifies the Quick mode dependency topology ("staged across types" — all gen-journeys parallel, then all gen-contracts parallel, etc.) with explicit rationale. Template injection execution model is explained in detail (prompt-level behavior switching, not code-level variable resolution). Six key scenarios now include edge cases (zero Journey output, index calculation failure). NFRs expanded with execution time, token consumption, and error handling. resolveQuickDeps rewrite has a dedicated risk item with detailed mitigation. Success criteria expanded from 6 to 12, covering dependency chain correctness, Quick mode task count formula, AUTO_COMMIT/SKIP_EVAL_GATE injection, resolveQuickDeps unit tests, backward compatibility, and zero-Journey abort.

3. **Quick mode quality model corrected**: The proposal no longer relies on "downstream code reconnaissance compensates." Instead, it defines explicit upstream quality gates: scope + success criteria as hard prerequisites, key scenarios as optional with smoke-level fallback, and no reliance on downstream compensation for uncovered scenarios.

4. **Resource estimate now justified**: 8 concrete task breakdowns with individual effort estimates, totaling 6-7 equivalent tasks within the stated 5-8 range.

5. **Remaining tension — Industry Benchmarking**: This dimension saw the least improvement across all three iterations. Citations remain shallow (name-plus-one-sentence), the Graphwalker analogy is still questionable, and no industry-validated solution is offered as a genuine alternative.

---

## Rubric Scoring (Phase 2)

### Dimension 1: Problem Definition — 85/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 35/40 | Core problem is clear: gen-journeys/gen-contracts require manual invocation, creating a pipeline gap. The distinction between auto-generated task types and skill invocation is clarified in the first two paragraphs. Minor residual ambiguity: "自动生成的测试流水线任务" could benefit from a parenthetical noting this refers to `autogen.go` task generation, not CI/CD. |
| Evidence provided | 25/40 | Three evidence points: code structure analysis, user behavior observation, Quick mode quality risk. No user feedback, no issue tracking references, no concrete failure scenarios with reproduction steps. Evidence remains "we inspected the code" rather than "this caused user-visible harm X times." This has been consistent across all three iterations and is the main drag on this dimension. |
| Urgency justified | 25/30 | v3.0.0 dependency is stated. Test profile system as prerequisite is mentioned. Cost of delay is implied (pipeline cannot run end-to-end) but not quantified — what specifically is blocked in the release without this? |

**Deductions**: -5 for minor framing ambiguity; -15 for evidence limited to code inspection; -5 for unquantified delay cost.

### Dimension 2: Solution Clarity — 100/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 38/40 | Pipeline diagrams for both modes are clear. Task type names, template files, mode-specific behavior are specified. Quick mode topology (staged across types) is now explicitly chosen with rationale. A reader can explain back what will be built. |
| User-facing behavior described | 35/45 | `forge task index` output described with specific task count formula (nTypes * 3 + 2). Quick mode task list growth is quantifiable. But the actual CLI experience — what the user sees when running `forge task index` before vs after — is not illustrated with a concrete example (e.g., "before: T-quick-gen-and-run-playwright, T-quick-run, T-quick-verify; after: T-quick-gen-journeys-playwright, T-quick-gen-contracts-playwright, ..."). |
| Technical direction clear | 27/35 | resolveQuickDeps full rewrite explicitly stated with topology (staged across types) and rationale. findTaskIndexByPrefix named as the replacement approach. ResolveFirstTestDep adaptation specified. Template injection execution model explained. Gap: resolveBreakdownDeps rewrite scope — the proposal mentions replacing hardcoded indices but doesn't enumerate all the dependency pairs that need to be built after the rewrite (e.g., gen-journeys -> eval-journey for each profile type, eval-journey -> gen-contracts for each profile type, etc.). |

**Deductions**: -2 for approach; -10 for missing CLI output example; -8 for resolveBreakdownDeps rewrite dependency pair enumeration.

### Dimension 3: Industry Benchmarking — 78/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 30/40 | Specific products cited: Cypress Studio, Playwright Codegen, Selenium IDE (record-replay); Testim, Mabl, Reflection (AI-assisted); Graphwalker, Cucumber/SpecFlow (model-based/BDD). But citations are name-plus-one-sentence — no links, versions, papers, or published patterns. |
| At least 3 meaningful alternatives | 18/30 | Three alternatives listed. "Merge into gen-and-run" has Testim analogy. "Do nothing" is genuine. "Independent auto-tasks" is the proposal. None of the alternatives represent adopting an industry-validated solution verbatim — all are Forge-internal architectural choices. |
| Honest trade-off comparison | 18/25 | Selected approach includes "开发量增加约 30%". The "Industry pattern adoption" paragraph maps Cucumber feature -> Journey, scenario -> Contract, step -> test-script. However, the Graphwalker analogy remains questionable — Graphwalker is a graph-based test path generator, not a pipeline framework. The "stage-per-artifact" pattern is generic enough to apply to any multi-stage system. |
| Chosen approach justified against benchmarks | 12/25 | Forge's unique positioning (PRD -> Journey -> Contract -> Test Script full chain) is articulated. The "only Testim has a comparable end-to-end pipeline" claim is interesting but unsubstantiated. The Cucumber mapping adds some justification but the mapping itself is somewhat forced — Cucumber's layering is designed for human collaboration (non-technical stakeholders write Gherkin), while Forge's layering is AI-generated. The fundamental constraint difference is not acknowledged. |

**Deductions**: -10 for shallow citations; -12 for no industry-validated alternative; -7 for questionable Graphwalker analogy; -13 for forced Cucumber mapping and subjective justification.

### Dimension 4: Requirements Completeness — 95/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 36/40 | Six key scenarios: Breakdown full pipeline, Quick simplified pipeline, Quick proposal input, backward compatibility, zero Journey abort, resolveQuickDeps index error. Good coverage of edge cases. Remaining gap: auto-task execution failure mid-pipeline is mentioned in NFR ("遵循现有 autogen 错误传播机制") but not as a standalone scenario with explicit behavior description. |
| Non-functional requirements | 33/40 | Significant improvement. Now includes: execution time impact (10-30 min + 5-15 min, Quick total 15-45 min), token consumption (template injection adds ~50 tokens), error handling (follows existing autogen error propagation — task marked failed, downstream blocked, pipeline halted, no new error paths). Remaining gaps: (1) Concurrent execution resource limits — "all gen-journeys parallel" implies N simultaneous LLM sessions per profile type count, but no discussion of API rate limits or concurrent session caps. (2) Template injection backward compatibility verification. |
| Constraints & dependencies | 26/30 | gen-contracts eval-journey Blocker conflict resolved (SKIP_EVAL_GATE=true). gen-journeys HARD-RULE interactive conflict resolved (AUTO_COMMIT=true). Template injection execution model explained. Test profile system dependency stated. Minor gap: it's unclear whether SKILL.md modifications are prerequisites for or concurrent with the autogen.go changes. |

**Deductions**: -4 for missing mid-pipeline failure scenario; -7 for missing concurrent execution limits; -4 for dependency sequencing ambiguity.

### Dimension 5: Solution Creativity — 42/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 18/40 | Self-admitted: "无特殊创新。" Forge's unique chain is a positioning statement about the existing system, not innovation in this proposal. |
| Cross-domain inspiration | 8/35 | BDD pipeline stage-splitting and Cucumber/SpecFlow patterns are within-domain (testing). No cross-domain borrowing from unrelated domains. |
| Simplicity of insight | 16/25 | Core insight (two task types into existing framework) is simple and good. Template injection is an elegant minimal-invasion solution. The staged-across-types Quick topology choice is a clean design decision. But the overall implementation complexity (two full function rewrites, SKILL.md adaptations, template injection) is non-trivial. |

**Deductions**: -22 for no meaningful innovation; -27 for no cross-domain inspiration; -9 for implementation complexity understated.

### Dimension 6: Feasibility — 85/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 35/40 | resolveBreakdownDeps hardcoded indices (lines 426-427) accurately described. Full rewrite with findTaskIndexByPrefix specified. resolveQuickDeps rewrite has explicit topology (staged across types). ResolveFirstTestDep Quick branch fix specified. Template injection execution model explained. Gap: resolveBreakdownDeps rewrite — the proposal says "将整个 resolveBreakdownDeps 重写为基于 findTaskIndexByPrefix 的 ID 查找" but doesn't enumerate all the dependency pairs that must be built after the rewrite. |
| Resource & timeline | 25/30 | 8 concrete task breakdowns with individual effort estimates (0.5-2 tasks each), totaling 6-7 equivalent tasks within 5-8 range. Much better justified. Minor tension: resolveBreakdownDeps rewrite is estimated at ~1 task but noted as having "核心复杂度在于多 profile type 依赖链构建" — this could expand. SKILL.md adaptation at 1-2 tasks has uncertainty. |
| Dependency readiness | 25/30 | SKILL.md interactive conflicts resolved via template injection. Test profile system merged. SKILL.md files stable. Minor gap: SKILL.md modifications (adding conditional behavior for AUTO_COMMIT/SKIP_EVAL_GATE) must be tested to not break existing manual invocation paths — this testing is implied but not explicitly listed in task breakdown. |

**Deductions**: -5 for incomplete resolveBreakdownDeps rewrite specification; -5 for task estimate tension; -5 for missing manual-path regression testing in task breakdown.

### Dimension 7: Scope Definition — 75/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 29/30 | 14 concrete deliverables, most specifying file names, function names, or task type names. "更新 ARCHITECTURE.md 和相关文档" still slightly vague — what related documents? |
| Out-of-scope explicitly listed | 24/25 | Five out-of-scope items including a disposal strategy for the test-gen-and-run.md template. Good boundary setting. |
| Scope is bounded | 22/25 | "5-8 coding tasks" with 8-item task breakdown supporting the estimate. 14 in-scope items mapped to 6-7 equivalent tasks is reasonable. Minor uncertainty: task 7 (SKILL.md adaptation) at "1-2 tasks" leaves scope boundary slightly soft. |

**Deductions**: -1 for vague documentation item; -1 for minor out-of-scope gap; -3 for soft scope boundary on SKILL.md task.

### Dimension 8: Risk Assessment — 82/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 28/30 | Five risks: proposal.md quality (M/M), Breakdown index offset (H/H), resolveQuickDeps rewrite topology error (H/H), gen-and-run replacement backward compatibility (L/M), ResolveFirstTestDep Quick branch mismatch (H/H). The critical addition is the resolveQuickDeps risk item. Minor gap: no regression risk for resolveBreakdownDeps full rewrite — this is a core function affecting all Breakdown mode pipelines. |
| Likelihood + impact rated | 25/30 | Ratings are honest and mostly well-calibrated. resolveQuickDeps at H/H with specific code references is appropriate. resolveBreakdownDeps at H/H with line numbers is accurate. Minor gap: the resolveQuickDeps mitigation says "使用与 resolveBreakdownDeps 相同的 findTaskIndexByPrefix 模式，复用已验证的查找逻辑" — this assumes resolveBreakdownDeps is rewritten and verified before resolveQuickDeps. If both are developed in parallel (which the task breakdown doesn't rule out), the "已验证" assumption may not hold. |
| Mitigations are actionable | 29/30 | Mitigations are concrete and executable: findTaskIndexByPrefix replacement, unit tests for 2+ profile types, per-type dependency pair construction, downgrade boundaries with abort/degrade conditions. Near-excellent. |

**Deductions**: -2 for missing resolveBreakdownDeps regression risk; -5 for parallel-development assumption in mitigation; -1 for minor mitigation detail.

### Dimension 9: Success Criteria — 75/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 50/55 | 12 criteria, most highly testable. Standouts: dependency chain correctness (depends-on field verification), Quick mode task count formula (nTypes * 3 + 2), AUTO_COMMIT injection (log verification), SKIP_EVAL_GATE injection (log verification), resolveQuickDeps unit test requirement, backward compatibility (old index.json loads without error). Deduction: "gen-journeys 支持 proposal.md 作为输入" verification includes "生成的 Journey 文件内容与 proposal.md 的 scope/success criteria 对应" — the "内容对应" part is subjective and cannot be automated. |
| Coverage is complete | 25/25 | 12 criteria cover the majority of 14 in-scope items. Coverage includes: task generation, pipeline correctness, mode-specific behavior, injection mechanisms, unit testing, backward compatibility, edge cases (zero Journey), and infer.go recognition. Only "更新 ARCHITECTURE.md" lacks a specific criterion, which is acceptable for documentation updates. |

**Deductions**: -5 for partially subjective criterion 4.

### Dimension 10: Logical Consistency — 78/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses stated problem | 32/35 | Yes. Making gen-journeys/gen-contracts auto-generated tasks directly solves the pipeline gap. Quick mode gen-and-run replacement is scope expansion beyond the stated problem, but is now explicitly framed as a necessary co-change with rationale. |
| Scope <-> Solution <-> Success Criteria aligned | 24/30 | Significantly improved. Most in-scope items now have corresponding success criteria. Remaining gaps: (1) "重写 resolveBreakdownDeps" has only indirect coverage through "依赖链正确性" criterion — no criterion specifically verifies the Breakdown mode rewrite. (2) "更新 ARCHITECTURE.md" has no criterion. (3) "修改 gen-journeys SKILL.md 支持非交互模式" has AUTO_COMMIT injection criterion but doesn't verify the default behavior (manual invocation) is preserved. |
| Requirements <-> Solution coherent | 22/25 | Requirements-to-solution mapping is now clear. Quick mode quality degradation boundary (scope + success criteria as hard prerequisites) is consistent with the solution's abort/degrade behavior. Backward compatibility requirement has corresponding criterion. Minor gap: "resolveQuickDeps 索引计算错误" scenario requires panic with diagnostic output, but the solution section doesn't specify the implementation mechanism (where in the code does the panic occur? In findTaskIndexByPrefix? In the caller?). |

**Deductions**: -3 for scope expansion beyond problem; -6 for scope/criteria alignment gaps; -3 for panic mechanism unspecified.

---

## Cross-Dimension Coherence Check

The most improved areas since iteration 2 are:

1. **Requirements Completeness** (70 -> 80 -> 95): Edge cases, NFRs, and constraint resolution have been systematically addressed.
2. **Risk Assessment** (55 -> 72 -> 82): resolveQuickDeps now has a dedicated risk item with detailed mitigation.
3. **Success Criteria** (58 -> 63 -> 75): Expanded from 6 to 12 criteria with high testability.
4. **Feasibility** (60 -> 75 -> 85): Task breakdown with effort estimates now supports the timeline claim.

The persistent weakness across all three iterations is **Industry Benchmarking** (45 -> 78 -> 78): citations remain shallow, no industry-validated solution is a genuine alternative, and the Cucumber/Graphwalker analogies are forced. This dimension saw improvement from iteration 1 to 2 but no meaningful change from iteration 2 to 3.

---

## Phase 3 — Blindspot Hunt

### [blindspot-1] SKILL.md conditional branch default behavior not verified
The proposal requires gen-journeys SKILL.md to add conditional behavior: "if AUTO_COMMIT=true is present in context, skip human approval." But the proposal never explicitly states that the default behavior (no AUTO_COMMIT flag, i.e., manual `/gen-journeys` invocation) must remain unchanged. If the SKILL.md modification inadvertently inverts the logic, manual invocation would also skip human approval. No success criterion verifies that manual invocation still requires user approval. Quote: "gen-journeys 据此跳过人工审批直接提交生成的 Journey 文件" — but what is the default when AUTO_COMMIT is absent?

### [blindspot-2] Staged across types topology has no concurrency limit discussion
The proposal chooses "all gen-journeys parallel" as the Quick mode topology. With N profile types, this means N simultaneous LLM sessions. The proposal provides no discussion of: (1) LLM API rate limits under concurrent sessions; (2) Whether Forge's autogen executor has a concurrency cap; (3) Peak token consumption during parallel execution (not total, but instantaneous). Quote: "所有 profile type 的 gen-journeys 并行执行 → 所有 gen-contracts 并行 → 所有 gen-scripts 并行 → run → verify。"

### [blindspot-3] Deprecated gen-and-run type visible in CLI help
The proposal keeps TypeTestGenAndRun for backward compatibility. Success criterion says "`forge -h` 输出的类型列表包含 test.gen-journeys 和 test.gen-contracts" but does not address whether test.gen-and-run should be hidden from the help output. Users will see a deprecated type that can never be generated again, causing confusion. No criterion or scope item addresses this UX issue. Quote: "保留 TypeTestGenAndRun 类型定义和 infer 逻辑以防有历史 index.json 引用" combined with "`forge -h` 输出的类型列表包含 test.gen-journeys 和 test.gen-contracts."

### [blindspot-4] No rollback plan
If resolveBreakdownDeps or resolveQuickDeps rewrite causes regression in existing Breakdown/Quick mode pipelines, the proposal provides no rollback strategy. While gen-and-run type is preserved for backward compat, the code changes to resolveBreakdownDeps/resolveQuickDeps are not reversible without a git revert. The proposal does not discuss: (1) Feature flag to toggle old vs new dependency resolution; (2) Incremental rollout strategy; (3) Smoke test requirements before merge. Quote from Key Risks: no entry for rollback or regression recovery.

### [blindspot-5] Task breakdown does not include integration/smoke testing task
The 8-item task breakdown covers implementation tasks but has no dedicated integration testing task. The success criteria include "resolveQuickDeps 单元测试" but no criterion or task covers end-to-end pipeline verification (running `forge task index` and `forge task run` on a real feature with both Breakdown and Quick modes). Quote: "合计约 6-7 个等效任务，在 5-8 范围内" — no integration testing task in the list.

---

## Freeform Finding Integration

| Finding | Resolution |
|---------|------------|
| [high] resolveBreakdownDeps hardcoded indices | **Fully addressed**: Lines 426-427 accurately described. Full rewrite with findTaskIndexByPrefix specified as both solution and risk mitigation. |
| [high] resolveQuickDeps rewrite needed | **Fully addressed**: Complete rewrite explicitly stated. New topology (staged across types) chosen with rationale. Dedicated risk item (H/H) with detailed mitigation. Unit test requirement in success criteria. |
| [high] gen-contracts eval-journey Blocker conflict | **Fully addressed**: SKIP_EVAL_GATE=true template injection proposed. Execution model explained. Success criterion verifies injection via log output. |
| [medium] gen-journeys HARD-RULE interactive conflict | **Fully addressed**: AUTO_COMMIT=true template injection proposed. Execution model explained. Success criterion verifies injection via log output. |
| [medium] Quick mode proposal.md input quality | **Largely addressed**: Degradation boundary defined (scope + success criteria as hard prerequisites, key scenarios as optional with smoke-level fallback). No longer relies on downstream compensation. Remaining: blindspot-1 (default behavior verification). |
| [medium] ResolveFirstTestDep Quick branch | **Fully addressed**: Explicitly stated in scope with specific new prefix `T-quick-gen-journeys`. Listed as risk item with H/H rating. |

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 85 | 110 |
| 2. Solution Clarity | 100 | 120 |
| 3. Industry Benchmarking | 78 | 120 |
| 4. Requirements Completeness | 95 | 110 |
| 5. Solution Creativity | 42 | 100 |
| 6. Feasibility | 85 | 100 |
| 7. Scope Definition | 75 | 80 |
| 8. Risk Assessment | 82 | 90 |
| 9. Success Criteria | 75 | 80 |
| 10. Logical Consistency | 78 | 90 |
| **Total** | **795** | **1000** |
