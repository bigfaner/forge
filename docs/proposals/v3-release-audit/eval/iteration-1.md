---
iteration: 1
reviewer: CTO Adversary
date: 2026-05-24
status: completed
---

# Iteration 1: CTO Adversary Rubric Scoring

## Phase 1: Reasoning Audit

The proposal argues: Forge v3.0.0 docs are systemically drifted from implementation → 50 audit findings across 5 dimensions → tiered remediation (P0/P1/P2) → ~5h effort → no runtime code changes.

**Chain trace:**

1. Problem claim: docs-impl drift is systematic and severe. Evidence: 50 items with severity breakdown. **Chain holds** — the table is specific and auditable.
2. Urgency: v3.0.0 is a major release, users will be misled. **Chain holds** — version 2.16.1 in README for a 3.0.0 release is objectively wrong.
3. Solution: tiered fix by severity. **Chain mostly holds**, but the scope boundary claim ("no runtime code changes") is contradicted by P0.4 (SKILL.md splitting changes agent context loading), P0.5 option (a) (rubric creation), and P1.8 (adding Load directives). Pre-revision partially addressed this by broadening the scope statement to include "插件内容调整（rules/SKILL.md 结构重组）" — an improvement, but still underplays the runtime-adjacent nature of SKILL.md restructuring.
4. Alternatives: 4 options, including do-nothing. **Chain holds** but alternatives are mostly scope variations of the same approach, not genuinely different strategies (e.g., automated doc generation, contract testing, CI doc linting).
5. Feasibility: ~5h, no external deps. **Chain is optimistic** — eval/SKILL.md splitting involves Mermaid flowcharts, Phase 0 conditional branching, and 7+ cross-references. The 1.5h estimate for this alone is aggressive given the stated complexity.
6. Success criteria: 6 checkboxes. **Chain is incomplete** — P1 items 6-12 have no corresponding success criteria. Item 12 (ARCHITECTURE.md subsystem documentation) is explicitly P1 but has no measurable criterion.

## Phase 2: Dimension Scoring

### 1. Problem Definition: 88/110

**Problem stated clearly (35/40):** The core problem — systematic documentation-implementation drift — is unambiguous. Two readers would interpret it the same way. Minor deduction: "系统性偏差" is stated but the structural pattern of the drift (v2→v3 naming migration fallout, component removal without doc cleanup, feature additions without doc updates) is not analyzed. Understanding *why* the drift accumulated would strengthen the problem frame.

**Evidence provided (35/40):** The 5-dimension audit table with severity breakdown is strong quantitative evidence. The "Assumptions Challenged" table provides excellent falsification evidence. Deduction: the headline originally said "27 个偏差项" but was corrected to "50 个偏差项" in pre-revision. The pre-revision fix resolved the factual error. However, the audit methodology ("5维度交叉审计") is described in one sentence — no explanation of how each dimension was measured, what the unit of analysis was, or how severity thresholds were defined. An external reader cannot reproduce the audit.

**Urgency justified (18/30):** "v3.0.0 主版本发布" and "README 版本号仍停留在 2.16.1" are concrete urgency signals. But the urgency section does not quantify the cost of delay: How many users will encounter the wrong docs? Are there beta users already affected? Is there a release date deadline? Without these, urgency rests on general principle rather than project-specific pressure.

### 2. Solution Clarity: 95/120

**Approach is concrete (38/40):** The P0/P1/P2 tiering with execution ordering is highly concrete. Each item specifies what changes, where, and why. The pre-revision addition of explicit execution order (P0.4→P0.5→P0.2→P0.3→P0.1) is a significant structural improvement. Deduction: P1 items 10 ("if applicable") and P2 items 13-21 are progressively less specific — some are single-line descriptions that could mean different things to different implementers.

**User-facing behavior described (30/45):** This is the weakest sub-dimension. The proposal describes *what files change* but not *what the user will experience differently after the fix*. For README: what will the new README look like? What sections will it have? For ARCHITECTURE.md: what will the corrected agent architecture section say? The proposal lists error categories but does not show the target state. The reader cannot visualize the end product.

**Technical direction clear (27/35):** SKILL.md splitting direction is clear (extract to rules/). CLI reference fix direction is clear (rename commands). ARCHITECTURE.md fix direction is clear (correct factual errors). Deduction: the README rewrite direction is vague — "全面重写" with a laundry list of corrections, but no indication of structural approach (reorganize sections? rewrite from scratch? patch individual claims?). The pre-revision improved P0.1 by adding "完整的命名体系替换——旧命名与 dot-notation 新命名零重叠，非单纯数量修正", which is helpful but still does not describe the target README structure.

### 3. Industry Benchmarking: 55/120

**Industry solutions referenced (15/40):** One sentence: "开源项目的发布审计通常通过 CHANGELOG + Breaking Changes 文档完成。" No specific project cited, no tool named, no published pattern referenced. This is a hand-waving reference to an industry practice, not evidence of actual benchmarking. No mention of doc-linting tools (e.g., `markdown-link-check`, `lychee`), automated doc testing patterns (e.g., test fixtures that assert doc accuracy), or how mature open-source projects (Kubernetes, React, Rust) handle doc-impl alignment at major version boundaries.

**At least 3 meaningful alternatives (18/30):** The comparison table has 4 rows including do-nothing. However, "仅修复 README", "仅修复 Critical", and "分层全量修复" are all scope variations of the same manual-fix approach. None represents a genuinely different strategy. Missing alternatives: automated doc generation from code (e.g., generating CLI reference from Go command definitions), contract testing between docs and code (e.g., CI checks that validate doc claims against codebase state), or partial automation (e.g., scripts that verify counts/paths). At least one alternative should be an industry-validated approach to preventing doc drift, not just fixing it.

**Honest trade-off comparison (12/25):** Trade-offs are present but shallow. "工作量较大" for the selected approach is the only stated con — no analysis of risk of overcorrection, no comparison of manual vs. automated approaches, no consideration of maintenance cost. The "do-nothing" alternative is presented as a straw man with only "用户和 agent 都会被误导" as a con.

**Chosen approach justified against benchmarks (10/25):** The justification is "问题间存在依赖，分批修复不如一次性对齐". This is a reasonable argument for the selected tier, but since no actual industry benchmarks were cited, the justification is self-referential. There is no argument for why this approach is better or worse than what comparable projects do.

### 4. Requirements Completeness: 85/110

**Scenario coverage (30/40):** Four key scenarios are identified with current-state failures. These cover the primary user types (reader, contributor, agent). Deduction: missing scenarios include: (a) what happens when a user follows the incorrect task type names in README to create a task — does it silently fail or produce wrong results? (b) what happens when an agent loads eval SKILL.md and encounters the `harness` type with no rubric — does it error or hang? (c) regression scenario — after the fix, how do we verify no new errors were introduced? The third scenario is partially addressed by NFR #1 but not as an explicit test scenario.

**Non-functional requirements (30/40):** Three NFRs are identified, all relevant. Deduction: missing NFRs include: (a) performance — does adding Load directives for orphaned rules increase agent context consumption measurably? (b) compatibility — are there downstream consumers of the current (wrong) documentation format that might break? (c) maintainability — what prevents this drift from reoccurring? No mention of prevention mechanisms.

**Constraints & dependencies (25/30):** Three constraints are well-specified, all referencing actual project conventions by filename. This is strong. Deduction: the constraint that "ARCHITECTURE.md must not grow beyond a readable length" is implied by item 12's scope but not stated as an explicit constraint. Adding all missing subsystem documentation to ARCHITECTURE.md could make it unwieldy.

### 5. Solution Creativity: 40/100

**Novelty over industry baseline (10/40):** The proposal explicitly states "这是一次标准的技术债清理，无特殊创新." Honest but self-scoring at zero. The 5-dimension audit framework has minor novelty but is essentially a structured checklist. No differentiation from standard doc audit approaches.

**Cross-domain inspiration (5/35):** No evidence of cross-domain thinking. The proposal does not reference how other domains handle doc-impl drift (API contract testing, schema validation, generated documentation patterns). No inspiration from unrelated fields.

**Simplicity of insight (25/25):** The core insight is dead-simple: audit docs against code, fix what's wrong, prioritize by severity. This is appropriately simple for a tech-debt cleanup. No overengineering. The execution order dependency analysis (P0.4→P0.5→P0.2→P0.3→P0.1) is an elegant insight that prevents rework.

### 6. Feasibility: 75/100

**Technical feasibility (30/40):** All changes are docs and file cleanup. No showstoppers. The pre-revision added detailed risk analysis for SKILL.md splitting ((a) explicit Load instruction, (b) Phase 0 cross-file consistency, (c) Mermaid sync), which improves confidence. Deduction: the proposal does not address how the eval SKILL.md Mermaid flowchart will be updated after splitting — Mermaid diagrams are fragile to restructure and a visual diagram that references non-existent file boundaries is worse than no diagram.

**Resource & timeline feasibility (25/30):** ~5h total is plausible for an experienced developer who already has the audit data. However, the 1.5h estimate for eval SKILL.md splitting is aggressive given the stated complexity (7 cross-refs, 4 expert refs, Mermaid, 3-way conditional branching, distribution model compliance, post-split verification). A more realistic estimate would be 2.5-3h for this item alone, pushing total to ~6.5h.

**Dependency readiness (20/30):** "无外部依赖。所有信息已在审计中收集完毕。" This is partially true — the audit data is collected, but several P0 items depend on *decisions* not yet made: (a) P0.5 harness rubric decision affects eval SKILL.md content; (b) P0.4 splitting strategy affects README counts; (c) P1.12 subsystem documentation scope is undefined. The pre-revision improved this by making the harness decision explicit (recommending option (c)), but the decision is not yet confirmed.

### 7. Scope Definition: 68/80

**In-scope items are concrete (26/30):** P0 items are highly concrete with specific file names, line numbers, and change descriptions. The pre-revision significantly improved P0 items 1, 4, and 5 with additional detail. Deduction: P2 items 13-21 are one-liners that lack specificity. Item 10 ("if applicable") explicitly introduces uncertainty about scope.

**Out-of-scope explicitly listed (22/25):** Five items listed. Clear and relevant. Deduction: the boundary between "documentation update" and "runtime-adjacent change" is blurred. SKILL.md splitting changes agent behavior (what gets loaded into context). Adding Load directives changes runtime behavior. The pre-revision broadened the scope description to "插件内容调整" which is more honest, but the Out of Scope section still says "运行时代码重构" — the distinction between "runtime code" and "runtime-affecting plugin content" is not made explicit.

**Scope is bounded (20/25):** The P0/P1/P2 tiering with item counts provides clear boundaries. The ~5h timeline estimate bounds effort. Deduction: P1 item 12 (ARCHITECTURE.md subsystem documentation for 9+ new subsystems) is potentially unbounded. Writing architecture documentation for 9 subsystems could take anywhere from 2h to 20h depending on depth. No depth specification is given.

### 8. Risk Assessment: 65/90

**Risks identified (22/30):** Four risks identified. The SKILL.md splitting risk is well-detailed with three sub-risks in the pre-revision. Deduction: missing risks: (a) **Cascade dependency risk** — the execution order P0.4→P0.5→P0.2→P0.3→P0.1 means any delay in P0.4 blocks all subsequent items; (b) **Scope creep risk** — P1.12 (9 subsystem docs) could expand significantly; (c) **Validation regression risk** — no automated validation means success criteria are checked manually and subjectively. The pre-revision improved the risk section by making P0.4 and P0.5 risks explicit, but the overall risk table has not been expanded.

**Likelihood + impact rated (18/30):** Ratings are present (L/M/H) but not justified. Why is "SKILL.md 拆分破坏已有引用链" Medium likelihood when it's the most complex operation? Why is "README 重写引入新的不准确描述" Medium likelihood when the README has the most changes? The assessment appears to average out to "Medium" rather than being honestly calibrated. A truly honest assessment would have at least one High-likelihood risk.

**Mitigations are actionable (25/30):** Mitigations are specific and actionable: "逐条与代码交叉验证", "拆分前用 grep 确认所有引用", "拆分后验证 agent 上下文加载行为". The pre-revision added specific sub-points to the SKILL.md splitting mitigation. Deduction: the mitigation for harness rubric risk ("采用方案 (c) 纯删除，规避此风险") is excellent — but this is a mitigation that eliminates the risk rather than managing it, which suggests the risk should not be in the table at all (it's already decided).

### 9. Success Criteria: 55/80

**Criteria are measurable and testable (40/55):** The pre-revision improved the ARCHITECTURE.md criterion by scoping it to "已有内容" rather than all components. The CLI cross-reference criterion specifies exact grep patterns. The SKILL.md line count criterion is objectively testable. Deduction: (a) "100% 一致" for README and ARCHITECTURE.md is aspirational but practically untestable without a verification script — who decides what constitutes a "事实性声明"? (b) "所有 rules/ 文件至少被其父 SKILL.md 引用一次" — what about parameterized references (surface-<type>.md)? Does a pattern reference count? (c) The init-justfile template criterion is now "已评估（保留为参考实现或清理）" which is a process criterion, not an outcome criterion — it can be checked off by doing the evaluation regardless of the outcome.

**Coverage is complete (15/25):** The 6 criteria cover P0 items well. Deduction: **No success criteria exist for P1 items** (items 6-12). Specifically: (a) P1.8 (orphaned rules Load directives) — no criterion checks that Load instructions work correctly; (b) P1.12 (9 subsystem docs) — no criterion defines adequate documentation depth; (c) P1.6 (CLI reference flags) — no criterion defines "补全" completeness. The 13 Major items represent ~26% of the total audit findings but have 0% of the success criteria coverage.

### 10. Logical Consistency: 72/90

**Solution addresses the stated problem (30/35):** The tiered remediation directly addresses the 5-dimension audit findings. P0 items map to Critical findings. P1 items map to Major findings. The solution is proportionate to the problem. Deduction: the problem statement focuses on README.md and ARCHITECTURE.md as "核心文档", but the solution extends significantly into plugin internals (SKILL.md splitting, rules/ restructuring, dead code cleanup). This is not inconsistent — the audit expanded the scope beyond the initial problem frame — but the problem statement should be updated to reflect the full scope of what was discovered.

**Scope ↔ Solution ↔ Success Criteria aligned (20/30):** P0 scope → solution → criteria are well-aligned. P1 scope → solution → criteria have gaps: P1 items 6-12 are in scope and have solution descriptions but no success criteria. P2 items are in scope with neither solution details nor criteria. The pre-revision improved P0 alignment significantly by adding explicit execution ordering and detailed risk analysis.

**Requirements ↔ Solution coherent (22/25):** The NFRs (no new errors, no path resolution breakage, no reference chain breakage) map cleanly to the solution's cautious approach. The constraints map to actual project conventions. Deduction: NFR #2 ("死代码清理不能影响分发后的路径解析") is listed but the solution for P1.9 (dead code cleanup) does not describe how path resolution impact will be verified. The constraint exists but the solution does not demonstrate compliance with it.

## Phase 3: Blindspot Hunt

### CTO Failure Pattern 1: Overstated Value

The proposal is appropriately restrained — it does not claim the fixes will dramatically improve user experience or adoption. The value proposition is accurate: "users will stop being misled." No overstatement detected.

### CTO Failure Pattern 2: Hidden Costs

**Found.** The proposal estimates ~5h but underestimates:
1. The eval SKILL.md splitting requires post-split regression testing (agent context loading verification, end-to-end eval pipeline test). This is not in the estimate.
2. P1.12 (9 subsystem docs) is a significant authoring effort with no time estimate. If each subsystem takes 30-60 minutes to document, that's 4.5-9h additional — potentially doubling the total effort.
3. The README rewrite is listed last because it depends on all other items. If any upstream item uncovers additional discrepancies, the README may need multiple revision passes.

### CTO Failure Pattern 3: Solution Reintroducing the Problem

**Found.** The proposal's P0.1 (README rewrite) introduces risk of *new* factual inaccuracies — the exact problem the proposal exists to fix. The mitigation is "逐条与代码交叉验证" but no mechanism is proposed to prevent future drift. After this fix, the next release could drift again. There is no prevention, only remediation. A CI check or automated doc validation script would break this cycle but is not proposed (even as a P2 item or follow-up recommendation).

### CTO Failure Pattern 4: Unstated Assumptions

1. **Assumption: the audit is complete.** The 5-dimension framework covers existing doc content but not missing content. The pre-revision added P1.12 (missing subsystem docs), but the success criterion for P1.12 does not exist, so completeness cannot be verified.
2. **Assumption: one person does all the work.** The ~5h estimate assumes a single executor with full context. If work is split across team members, the dependency chain (P0.4→P0.5→P0.2→P0.3→P0.1) creates serialization bottlenecks.
3. **Assumption: the codebase is stable during remediation.** If other team members merge changes during the ~5h window, the audit findings may become stale. No mention of branching strategy or merge freeze.

### CTO Failure Pattern 5: Missing Rollback Plan

**Found.** No rollback plan exists. If the SKILL.md splitting breaks agent behavior in production, how do we revert? The proposal should specify: (a) work on a feature branch, (b) tag the pre-split state, (c) define a go/no-go checkpoint after P0.4 before proceeding to P0.3/P0.1. The execution order dependency makes this especially critical — if P0.4 introduces a regression, P0.2/P0.3/P0.1 cannot proceed.

---

## Bias Detection Report

**Annotated regions (pre-revised):**

Paragraphs with `<!-- pre-revised -->` markers: 8 paragraphs (lines 15-16, 33-34, 108-109, 114-115, 116-117, 127-128, 162-163)

Attack points targeting annotated regions:
1. [Problem Definition] Evidence table methodology not reproducible — "5维度交叉审计" described in one sentence
2. [Solution Clarity] User-facing behavior of corrected docs not shown — "按优先级分层修复"
3. [Scope Definition] P0 execution order creates serialization bottleneck — "P0.4→P0.5→P0.2→P0.3→P0.1"
4. [Scope Definition] SKILL.md splitting runtime impact understated in scope boundary — "不修改任何 Go 运行时代码"
5. [Scope Definition] P0.5 harness decision still unconfirmed — "推荐方案 (c)"
6. [Scope Definition] P1.12 unbounded scope for 9 subsystem docs
7. [Success Criteria] ARCHITECTURE.md criterion scoped to "已有内容" but P1.12 adds new content with no criterion
8. [Risk Assessment] Missing cascade dependency risk for serialized P0 execution

**Annotated attack density: 8 / 8 paragraphs = 1.00**

**Unannotated regions:**

Paragraphs without markers: ~33 paragraphs (remaining content)

Attack points targeting unannotated regions:
1. [Problem Definition] Urgency lacks project-specific pressure metrics
2. [Solution Clarity] README target state not described
3. [Industry Benchmarking] No real industry solutions cited — only one generic sentence
4. [Industry Benchmarking] Alternatives are scope variations, not different strategies
5. [Industry Benchmarking] Straw-man do-nothing alternative
6. [Requirements Completeness] Missing regression/verification scenario
7. [Requirements Completeness] Missing maintainability NFR (drift prevention)
8. [Solution Creativity] Self-admitted zero innovation
9. [Feasibility] eval SKILL.md splitting 1.5h estimate is aggressive
10. [Feasibility] P0.5 decision dependency not resolved
11. [Scope Definition] P2 items are one-liners lacking specificity
12. [Scope Definition] Out-of-scope boundary between "runtime code" and "runtime-affecting content" unclear
13. [Risk Assessment] Likelihood ratings not justified — tend toward "Medium"
14. [Risk Assessment] Missing scope creep risk for P1.12
15. [Risk Assessment] Missing validation regression risk
16. [Success Criteria] "100% 一致" not objectively testable without verification script
17. [Success Criteria] Parameterized references not addressed in rules criterion
18. [Success Criteria] Zero criteria for P1 items (13 Major findings)
19. [Success Criteria] init-justfile criterion is process, not outcome
20. [Logical Consistency] Problem statement narrower than solution scope
21. [Logical Consistency] NFR #2 compliance not demonstrated in P1.9 solution
22. [Blindspot] No rollback plan
23. [Blindspot] Hidden cost of P1.12 (9 subsystem docs) doubling total effort
24. [Blindspot] No drift prevention mechanism proposed

**Unannotated attack density: 24 / 33 paragraphs = 0.73**

**Ratio (annotated/unannotated): 1.00 / 0.73 = 1.37**

The ratio of 1.37 indicates a moderate bias toward scrutinizing pre-revised regions more heavily. This is expected — pre-revised regions received attention specifically because they had known issues. The bias is within acceptable range (< 2.0) and does not materially affect scoring validity. I have flagged one attack as `conflict-with-pre-revision` below.

---

## ATTACK POINTS (Tagged by Dimension)

### Annotated Region Attacks

1. **[Problem Definition]** Evidence methodology not reproducible — "5维度交叉审计——文档自检、文档-实现比对、接口契约验证、架构规范合规、依赖图分析" — Each dimension needs a brief description of measurement method; without it, audit completeness is an unverifiable claim.

2. **[Solution Clarity]** User-facing target state absent — "按优先级分层修复：Critical（发布阻塞）→ Major（发布前建议修复）→ Minor/Advisory（发布后迭代）" — The reader knows what gets fixed but not what the fixed state looks like. For a doc-remediation proposal, showing the target state (even as a section outline for the new README) is table stakes.

3. **[Scope Definition]** Serialized P0 execution creates single-point-of-failure — "P0.4（SKILL.md 拆分）→ P0.5（harness 决策）→ P0.2（ARCHITECTURE.md）→ P0.3（CLI 引用）→ P0.1（README 重写最后）" — The pre-revision correctly identified the dependency order, but did not add the corresponding risk (P0.4 failure blocks everything). The risk table must include cascade dependency risk.

4. **[Scope Definition]** Scope boundary claim contradicted by solution — "不修改任何 Go 运行时代码" — The pre-revision broadened the description to "插件内容调整（rules/SKILL.md 结构重组）" which is more honest, but the Out of Scope section still says only "运行时代码重构." SKILL.md splitting and Load directive additions change agent runtime behavior. The Out of Scope should explicitly exclude "Go source code" rather than "runtime behavior changes."

5. **[Scope Definition]** P0.5 decision listed as recommendation, not commitment — "推荐方案 (c)" — A P0 blocking item should have a decided solution, not a recommendation. If the decision is deferred to execution time, it should be explicitly flagged as a pre-requisite decision point.

6. **[Scope Definition]** `conflict-with-pre-revision` P1.12 scope is unbounded — "ARCHITECTURE.md 补充缺失的 v3.0.0 子系统文档：surface detection、worktree 管理、Convention（consolidate-specs）、forensic、deep-research、clean-code、extract-design-md、test-guide、learn" — The pre-revision added this as P1, but 9 subsystem documentation entries with no depth specification and no success criterion is an unbounded commitment. The pre-revision addressed the "missing subsystems" gap by listing them, but introduced a new scope control problem.

7. **[Success Criteria]** ARCHITECTURE.md criterion covers "已有内容" but P1.12 adds new content — "ARCHITECTURE.md 中已有内容的所有事实性声明与代码库 100% 一致；缺失子系统文档列为 P1 后续任务" — The success criterion correctly scopes P0 work, but P1.12 has no corresponding criterion. The proposal claims P1 items are "发布前建议修复" but provides no way to verify they succeeded.

8. **[Risk Assessment]** Missing cascade dependency risk — "注意执行顺序：P0.4→P0.5→P0.2→P0.3→P0.1" — The serialized execution order means P0.4 (the most complex operation) failure blocks all 17 P0 items. This is the highest-impact risk in the entire proposal and it is not in the risk table.

### Unannotated Region Attacks

9. **[Problem Definition]** Urgency lacks project-specific pressure — "v3.0.0 是主版本发布" — No release date, no beta user count, no open issues from confused users. Urgency rests on general principle alone.

10. **[Solution Clarity]** README target structure undefined — "README.md 全面重写：版本号、技能/命令/Agent 计数、任务类型表..." — A "全面重写" needs at minimum a section outline or structural mockup. The current description is a list of corrections, not a design.

11. **[Industry Benchmarking]** No real-world solutions cited — "开源项目的发布审计通常通过 CHANGELOG + Breaking Changes 文档完成" — This is one generic sentence. No specific projects, tools, or patterns referenced. A proposal at this quality level should cite at least 2-3 real projects or tools.

12. **[Industry Benchmarking]** Alternatives are scope variations — The 4 alternatives (do-nothing, README-only, Critical-only, full) are all "manual fix with different scope." Missing: automated doc generation from code, CI doc-validation, contract testing between docs and implementation.

13. **[Industry Benchmarking]** Straw-man alternative — "Do nothing | 零成本 | 用户和 agent 都会被误导 | Rejected" — The do-nothing alternative is presented with only cons, making it a clear straw man.

14. **[Requirements Completeness]** Missing regression verification scenario — "Agent 执行 skill 调用 CLI 命令 → 当前会在 gen-test-scripts 和 run-tests 中调用不存在的 CLI 命令" — There is no scenario for "after fix, verify the fix did not introduce new errors." NFR #1 mentions this but no test scenario describes how.

15. **[Requirements Completeness]** Missing maintainability NFR — No requirement or follow-up for preventing future drift. The proposal fixes current drift but has no mechanism to prevent recurrence.

16. **[Feasibility]** eval SKILL.md splitting time estimate is aggressive — "SKILL.md 拆分（gen-test-scripts + eval）：~1.5h（含拆分后引用链验证）" — The proposal itself describes 7 cross-refs, 4 expert refs, Mermaid flowchart, 3-way conditional branching, and distribution model compliance. 1.5h for splitting + verification is unrealistic.

17. **[Risk Assessment]** Likelihood ratings not calibrated — "README 重写引入新的不准确描述 | M" and "SKILL.md 拆分破坏已有引用链 | M" — Both are rated M. But SKILL.md splitting is objectively more complex (the proposal says so) and should be rated H likelihood. No risk is rated H likelihood, suggesting anchoring to Medium.

18. **[Success Criteria]** "100% 一致" is not independently testable — "README.md 所有事实性声明（版本号、计数、路径、命令名）与代码库 100% 一致" — Who enumerates "所有事实性声明"? Without a verification script or explicit checklist, this criterion is subjective.

19. **[Success Criteria]** Zero criteria for 13 P1 items — Items 6-12 are in scope as "发布前建议修复" but have no success criteria. 26% of audit findings have 0% criteria coverage.

20. **[Success Criteria]** init-justfile criterion is a process check, not outcome — "init-justfile/templates/ 下的 .just 模板文件已评估（保留为参考实现或清理）" — This can be checked off by saying "evaluated" regardless of the outcome. Should specify the expected outcome.

21. **[Logical Consistency]** Problem statement narrower than solution scope — "核心文档（README.md、ARCHITECTURE.md）与实际实现存在系统性偏差" — The solution extends to plugin internals (SKILL.md splitting, rules/ restructuring). The problem statement should mention that the audit revealed drift extends beyond core docs into plugin content.

22. **[Logical Consistency]** NFR compliance not demonstrated — "死代码清理不能影响分发后的路径解析" — The P1.9 solution (dead code cleanup) does not describe how path resolution impact will be verified after cleanup.

23. **[Blindspot]** No rollback plan — If P0.4 (SKILL.md splitting) breaks agent behavior, there is no documented rollback strategy. The serialized execution order makes this critical.

24. **[Blindspot]** No drift prevention mechanism — The proposal fixes current drift but proposes nothing (even as P2 or follow-up) to prevent recurrence. This is a maintenance anti-pattern: fix-once, break-again.

25. **[Blindspot]** P1.12 hidden cost — Writing architecture documentation for 9 subsystems is a significant authoring effort with no time estimate. Could double the ~5h total if executed as part of this proposal.
