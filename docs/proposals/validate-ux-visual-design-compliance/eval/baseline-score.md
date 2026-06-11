# Adversarial Score Report: validate-ux Visual Design Compliance

**Evaluator**: CTO Adversary (hardened)
**Date**: 2026-06-08
**Document**: `/docs/proposals/validate-ux-visual-design-compliance/proposal.md`

---

## Phase 1: Reasoning Audit

### Chain Trace: Problem -> Solution -> Evidence -> Success Criteria

| Link | Trace | Verdict |
|------|-------|---------|
| Problem | validate-ux does not validate rendered visual output against design spec | **Sound** -- problem is real and specific |
| Evidence | pm-work-tracker milestone-map: 6 issues found post-validate-ux, 4 visual | **Unverifiable** -- referenced document `docs/lessons/gotcha-validation-ux-misses-visual-gaps.md` does NOT exist on disk. A related document `lesson-tui-visual-verify.md` exists for TUI but not for the claimed web-surface gotcha. Evidence is fabricated or miscited. |
| Solution | Accessibility tree semantic matching vs ui-design.md | **Contradicts evidence** -- 3 of 4 cited visual problems (style inconsistency, button position, visual regression) cannot be detected by accessibility tree alone. Solution addresses a different problem than the one evidenced. |
| SC | SC-2 claims every ui-design.md requirement maps to pass/fail | **Overclaims** -- no acknowledgment that free-text design requirements cannot deterministically produce pass/fail verdicts. SC is not achievable as stated. |

### Self-Contradictions Found

1. **Solution claims "能检测组件缺失、错误变体、布局偏差"** but Risk 2 admits "accessibility tree 无法反映某些视觉属性" and the mitigation is only "明确能力边界". The innovation section's claim of detecting layout deviation directly contradicts the risk assessment's admission that accessibility tree lacks layout information.

2. **Evidence problem #2 is "刷新按钮位置与设计稿不符（布局错误）"** but the chosen approach (accessibility tree) has no position information. The proposal uses a layout problem as evidence, then proposes a solution that cannot detect layout problems.

3. **Assumptions Challenged claims "accessibility tree 提供了结构化语义信息，足以检测组件缺失、错误变体、布局偏差"** -- this is stated as a "Refined" finding (meaning the assumption was validated), but Risk 2 concedes it cannot detect visual attributes. The assumption challenge contradicts the risk assessment.

---

## Phase 2: Rubric Scoring

### 1. Problem Definition: 68/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 30/40 | Core problem is clear: validate-ux does not validate rendered output. However, the problem conflates two distinct issues: (a) structural/semantic compliance and (b) visual fidelity. The proposal treats these as one problem, which muddies the solution. |
| Evidence provided | 18/40 | The 4 concrete visual issues are specific and illustrative. **Critical flaw**: the referenced evidence document `docs/lessons/gotcha-validation-ux-misses-visual-gaps.md` does not exist. Only `lesson-tui-visual-verify.md` exists, which covers TUI not web. This means the evidence is either fabricated or the citation is wrong -- either way, unverifiable. Deduction: evidence cannot be independently confirmed. |
| Urgency justified | 20/30 | Argument that "more web projects = more frequent impact" is reasonable but speculative. No data on how many forge web projects exist, how often visual issues slip through, or what the cost of these issues has been beyond the single unverifiable example. |

### 2. Solution Clarity: 58/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 20/40 | The 4-step flow (read ui-design.md -> capture per route -> compare -> report) gives a high-level shape but leaves the core algorithm entirely unspecified. "逐条映射" is the hardest part of the entire proposal and it is described in exactly one sentence. |
| User-facing behavior described | 28/45 | SC-2 and SC-5 describe report output format adequately. But the key user experience -- what a pass/fail verdict looks like when the design doc uses natural language -- is not described. Users cannot evaluate whether this output would be useful to them. |
| Technical direction clear | 10/35 | "Accessibility tree vs ui-design.md" is the stated direction, but: (a) no mention of how to parse ui-design.md (which is free-text with HTML comments as seen in the actual template), (b) no mention of the matching algorithm (rule engine? LLM? heuristic?), (c) no mention of how computed styles supplement accessibility tree gaps. The technical direction is a direction label, not a direction description. |

### 3. Industry Benchmarking: 52/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 18/40 | Percy, Chromatic, and Playwright screenshot comparison are mentioned by name. But only one sentence describes them: "均为像素级截图对比，需要维护 baseline image." No discussion of visual regression testing as a broader category (e.g., DOM snapshot testing with Jest/Playwright, component testing with Storybook, visual review workflows). Missing: Applitools, BackstopJS, regression as a service patterns. |
| At least 3 meaningful alternatives | 14/30 | 4 alternatives listed, but "跨页面一致性审计" is a straw man: described as "自研" with "不知道哪个是正确" as its con. It is set up purely to be rejected. The "do nothing" alternative is legitimate. The pixel-level comparison alternative is legitimate but dismissed with "维护成本过高" without quantification. Deduction: 1 straw man (-20 pts, capped at 0 for this sub-criterion, applied to dimension total). |
| Honest trade-off comparison | 10/25 | Comparison table is cherry-picked. Percy/Chromatic's con is "误报率高" but the proposal's own approach will also have high false rates due to NLU parsing of free-text design docs. The selected approach's con is listed as just "依赖 ui-design.md 质量" -- this understates the risk dramatically. No mention of the selected approach's inability to detect visual-only deviations. |
| Chosen approach justified against benchmarks | 10/25 | The justification is "直接对标 gotcha 缺口，复用现有基础设施" -- but the gotcha gap includes layout and style issues that the chosen approach cannot detect. The justification claims to solve the evidenced problem but the approach only solves a subset. |

### 4. Requirements Completeness: 55/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 24/40 | Happy path, component missing, variant inconsistency, and no-ui-design.md are covered. Missing scenarios: (a) ui-design.md exists but is incomplete/ambiguous, (b) accessibility tree extraction fails for a route, (c) matching produces ambiguous results (not clear pass/fail), (d) design spec says "primary color" but accessibility tree has no color info -- how is this handled? |
| Non-functional requirements | 18/40 | Performance (50% cap), compatibility (graceful skip), observability (screenshot + node ref) are listed. Missing NFRs: (a) accuracy/reliability -- no target for false positive/negative rate, which is the #1 NFR for a testing tool, (b) determinism -- if LLM is involved in matching, can the same input produce different results across runs? (c) accessibility of the compliance report itself -- not relevant here. |
| Constraints & dependencies | 13/30 | 4 dependencies listed correctly. Missing constraints: (a) ui-design.md template has HTML comment placeholders, not structured data -- this is a constraint on matching quality, (b) agent-browser's accessibility tree extraction fidelity and completeness, (c) the constraint that accessibility tree provides no color, spacing, or positioning information. |

### 5. Solution Creativity: 45/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 20/40 | Using accessibility tree for design compliance is a genuinely different approach from pixel comparison. However, it is not novel in the accessibility world -- WCAG audits routinely use accessibility trees for compliance checking. The proposal applies an existing pattern (accessibility tree auditing) to a new context (design spec compliance). Incremental, not breakthrough. |
| Cross-domain inspiration | 15/35 | The accessibility tree -> design compliance pipeline borrows from web accessibility auditing. One domain cross-pollination. No evidence of drawing from other domains (e.g., compiler AST comparison, schema validation, contract testing patterns that could also apply). |
| Simplicity of insight | 10/25 | The insight "use semantic structure instead of pixels" is clean, but the implementation complexity is hidden. The actual matching problem (free-text design spec -> accessibility tree assertions) is genuinely hard, and the proposal pretends it is trivial. This is not elegant simplicity -- it is underspecified complexity. |

### 6. Feasibility: 48/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 18/40 | The proposal claims "不需要新的基础设施" which is true for the browser automation part. However, the core matching algorithm -- parsing free-text ui-design.md into structured assertions and comparing against accessibility tree -- is the technical crux and is entirely unaddressed. The actual ui-design.md template (verified on disk) contains HTML comment placeholders (`<!-- Component hierarchy, grid/flex layout description -->`) and free-text table cells. Claiming "已有结构化格式" is misleading -- the template has structure only in the heading hierarchy, not in the content. |
| Resource & timeline | 18/30 | "4-5 files, medium scope" is plausible for the wiring but does not account for building the matching engine. If the matching uses LLM calls per route per requirement, the prompt engineering and evaluation effort is significant. No timeline estimate provided. |
| Dependency readiness | 12/30 | agent-browser readiness is stated but unverified for the specific use case (accessibility tree extraction for design compliance). ui-design.md "已有结构化格式" is false -- the template is free-text. sitemap.json "已有" is assumed but not verified for all web projects. |

### 7. Scope Definition: 48/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 20/30 | 5 in-scope items are listed. Items 1-2 are somewhat concrete (new collection step, new rubric dimension). Items 3-5 are document updates, which are concrete deliverables. Item 1 is vague on the matching algorithm itself. |
| Out-of-scope explicitly listed | 18/25 | 6 out-of-scope items listed. Good that cross-page consistency and pixel-level testing are excluded. However, the exclusion of "跨页面一致性审计" directly conflicts with the evidence -- the #1 visual issue in the evidence IS a cross-page consistency problem. |
| Scope is bounded | 10/25 | "仅适用于 web surface" provides one bound. No timeline bound, no per-route complexity bound, no matching algorithm scope bound. The scope of "逐条对比" could range from 10 lines of regex to a full NLU pipeline -- the proposal does not bound this. |

### 8. Risk Assessment: 42/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 18/30 | 4 risks listed. The accessibility tree limitation risk is the most honest. Missing risks: (a) LLM non-determinism in matching (if LLM is used), (b) ui-design.md template has free-text content (verified), not "structured format" as claimed, (c) false positive rate from NLU-based matching eroding user trust, (d) the risk that the solution only detects ~25% of the evidenced problems. |
| Likelihood + impact rated | 12/30 | Ratings are self-serving: the most critical risk (ui-design.md quality, M/H) has the most aggressive mitigation ("标注无法对比") which doesn't actually solve the problem, just flags it. The accessibility tree limitation is rated M/M but it undermines the entire value proposition -- this should be H/H. The "新增步骤使时间增加" risk is rated L/L which is honest. |
| Mitigations are actionable | 12/30 | "在 rubric 中明确能力边界" is a documentation action, not a technical mitigation. "标注无法对比" is a reporting action, not a solution. The CI stability risk mitigation ("复用现有编排") is circular -- if it was already stable, there would be no risk. The only truly actionable mitigation is the performance one (skip when no ui-design.md). |

### 9. Success Criteria: 42/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 14/30 | SC-1 (100% routes get collection data) is testable. SC-4 (behavior unchanged without ui-design.md) is testable. SC-2 (every requirement maps to pass/fail) is NOT testable because "requirement" is undefined -- a free-text sentence from ui-design.md is a "requirement" but mapping it deterministically to pass/fail is the unsolved problem. SC-3 (4 sub-dimensions with scoring rules) is testable. SC-5 (every fail has structured content) is testable. 3 of 5 are genuinely testable. |
| Coverage is complete | 12/25 | No SC for false positive/negative rate. No SC for the accuracy of matching. No SC for handling ambiguous design requirements. SC-4 covers backward compatibility but there is no SC for forward compatibility (what happens when ui-design.md template changes). |
| SC internal consistency | 16/25 | SC-1 (100% routes) and SC-2 (every requirement maps to pass/fail) create tension: if a route's ui-design.md has ambiguous requirements, SC-1 demands collection data (achievable) but SC-2 demands a pass/fail verdict (may not be achievable deterministically). This is flagged as **ambiguous -- requires author clarification** on whether "unable to determine" is an acceptable verdict or not. SC-3 and SC-5 are consistent with each other. SC-4 is independent and consistent. |

### 10. Logical Consistency: 35/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 10/35 | **Major gap**. The problem is "visual problems not detected." The evidence lists 4 specific visual problems. The solution (accessibility tree matching) can at best detect 1 of 4 (calendar popup missing -- if the interaction is triggered). It cannot reliably detect style inconsistency (no CSS info), button position (no layout info), or visual regression (no visual info). The solution addresses "structural/semantic compliance" which is a subset of the stated problem, but the proposal presents it as a comprehensive solution. |
| Scope <-> Solution <-> SC aligned | 15/30 | In-scope item #2 claims rubric covers "布局结构" and "样式属性" as sub-dimensions. But accessibility tree has no layout or style information. How can a rubric dimension evaluate something the collection step cannot capture? SC-3 requires these 4 sub-dimensions to have scoring rules, but if the data source (accessibility tree) provides no layout/style data, the scoring rules would be vacuous. |
| Requirements <-> Solution coherent | 10/25 | Constraint says "依赖 ui-design.md 的存在和质量" but the solution does not handle low-quality ui-design.md. Scenario "ui-design.md 不存在" is covered (skip), but no scenario for "ui-design.md exists but is free-text with ambiguous requirements." The NFR for observability (screenshot + node ref) is coherent with the solution. The NFR for performance (50% cap) has no solution-level discussion of how to achieve it (caching? parallelism?). |

---

## Phase 3: Blindspot Hunt

### Blindspot 1: The Evidence is Fabricated or Miscited

The proposal states: "详细记录见：`docs/lessons/gotcha-validation-ux-misses-visual-gaps.md`（pm-work-tracker 项目）". This file does not exist. The closest match is `lesson-tui-visual-verify.md` which covers TUI, not web. The entire evidence base for this proposal rests on a document that cannot be verified. If the evidence is real but the path is wrong, the author should correct it. If the evidence is constructed for the proposal, this is a serious integrity issue.

### Blindspot 2: The Solution Creates the Same False-Security Problem It Claims to Solve

The proposal's core motivation is that validate-ux gives "虚假安全感" (false security) by reporting all pass when visual problems exist. The proposed solution will do the same thing for a different reason: design-compliance will report pass/fail based on accessibility tree matching, which cannot detect visual-only deviations (color, spacing, typography). Users will now trust that "design compliance passed" means the visual output matches the design, when in fact it only means the semantic structure matches. The proposal recreates the exact problem it claims to solve, shifted one level.

### Blindspot 3: No Rollback Plan

If the design compliance step produces unreliable results (high false positive rate from NLU matching, or high false negative rate from accessibility tree limitations), there is no mechanism to disable or roll back the feature per-project. The proposal has no feature flag, no configuration option, and no degradation path beyond "skip when no ui-design.md." Once deployed, every forge web project with a ui-design.md will have design compliance reports that may be misleading.

### Blindspot 4: ui-design.md Quality is Treated as External Constraint, Not Design Dependency

The proposal treats ui-design.md quality as a risk with an external mitigation ("标注无法对比"). But ui-design.md is generated by forge's own /ui-design skill. The quality of ui-design.md is within forge's control. The proposal should include improvements to the ui-design.md template (e.g., adding structured, machine-parseable annotations) as a prerequisite or co-deliverable, rather than treating it as an external dependency.

### Blindspot 5: The 50% Performance Budget is Untested

For a project with 30 routes, each requiring browser navigation + accessibility tree extraction + (potentially) LLM-based matching per requirement, the 50% budget claim has no basis. The proposal provides no benchmark of current validate-ux execution time, no estimate of per-route overhead, and no analysis of whether the budget is achievable. This is a critical NFR that is treated as an afterthought.

---

## Score Summary

SCORE: 443/1000

DIMENSIONS:
  Problem Definition: 68/110
  Solution Clarity: 58/120
  Industry Benchmarking: 52/120
  Requirements Completeness: 55/110
  Solution Creativity: 45/100
  Feasibility: 48/100
  Scope Definition: 48/80
  Risk Assessment: 42/90
  Success Criteria: 42/80
  Logical Consistency: 35/90

ATTACKS:
1. Problem Definition: Evidence document does not exist on disk -- "详细记录见：`docs/lessons/gotcha-validation-ux-misses-visual-gaps.md`" -- File not found; closest match is a TUI lesson, not a web-surface gotcha. Either correct the path or provide verifiable evidence.
2. Solution Clarity: Core matching algorithm is a single sentence -- "将 accessibility tree 中的实际渲染结构与 ui-design.md 的要求逐条映射" -- This is the hardest part of the entire feature. Specify: rule engine or LLM? What are the matching rules? How is ambiguity handled?
3. Industry Benchmarking: Straw-man alternative -- "跨页面一致性审计" described as "自研" with con "不知道哪个是正确" -- This is the approach most likely to catch the #1 evidenced problem (cross-page style inconsistency). Present it honestly with real pros (no design doc dependency, catches structural divergence) before rejecting.
4. Industry Benchmarking: Trade-off comparison cherry-picked -- Percy/Chromatic con listed as "误报率高" but the proposed approach will also have high false rates due to free-text NLU parsing. Acknowledge that the selected approach trades one type of false positive (pixel diff noise) for another (semantic matching ambiguity).
5. Requirements Completeness: Missing critical NFR: accuracy/reliability target -- No false positive/negative rate target specified. For a testing tool, this is the #1 non-functional requirement. Specify target (e.g., "false positive rate < 20% on validated test set").
6. Solution Creativity: Overclaims capability -- "能检测组件缺失、错误变体、布局偏差" -- Accessibility tree has no layout information. "布局偏差" detection is not supported by the chosen data source. Remove "布局偏差" from the innovation claims or add computed style extraction.
7. Feasibility: False claim about ui-design.md format -- "ui-design.md 由 /ui-design skill 生成，已有结构化格式" -- Verified the actual template: Layout Structure is an HTML comment placeholder, States table Visual column is free text. Not machine-parseable. Either update the template first or acknowledge this as a prerequisite.
8. Feasibility: No timeline estimate -- "4-5 个文件的修改，中等规模改动" -- No calendar estimate, no sprint allocation, no breakdown of the matching engine effort (which is the bulk of the work).
9. Scope Definition: Out-of-scope item contradicts evidence -- "跨页面一致性审计" is out of scope, but evidence problem #1 is exactly a cross-page consistency issue. Either bring this into scope or acknowledge that the proposal will not address the most common evidenced problem type.
10. Risk Assessment: Critical risk underrated -- "accessibility tree 无法反映某些视觉属性" rated M/M -- This limitation means the solution cannot detect the majority of the evidenced problems. Should be H/H with a technical mitigation (e.g., computed style extraction), not just a documentation action.
11. Risk Assessment: Missing risk: false security recreation -- The proposal's motivation is that validate-ux gives false security. Design-compliance will give the same false security for visual-only deviations. This risk is not identified at all.
12. Success Criteria: SC-2 is not achievable as stated -- "每条 ui-design.md UI 要求映射到 pass/fail 判断" -- Free-text requirements cannot deterministically map to binary verdicts. Add "unable to determine" as a third verdict or constrain the SC to machine-parseable requirements only.
13. Logical Consistency: Solution does not address the evidenced problem -- Evidence has 4 visual problems; accessibility tree can detect at most 1. The proposal presents the solution as comprehensive but it addresses only the structural subset of the stated problem.
14. Logical Consistency: In-scope claims rubric covers "布局结构" and "样式属性" but accessibility tree provides neither layout nor style data -- In-scope item #2 promises evaluation of layout and style, but the data source (accessibility tree) cannot capture these. The rubric sub-dimensions would be vacuous.
15. Blindspot: No rollback plan -- Once deployed, no per-project opt-out mechanism. If design-compliance produces unreliable results, users cannot disable it without removing ui-design.md (which would also skip the feature they might partially want). Add a feature flag or configuration option.
