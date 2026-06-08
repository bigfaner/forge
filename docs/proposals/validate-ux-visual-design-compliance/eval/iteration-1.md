---
iteration: 1
scorer: adversary-cto
model: glm-5.1
date: 2026-06-08
---

# Iteration 1: Adversarial Scoring Report

## Document Under Review

`/docs/proposals/validate-ux-visual-design-compliance/proposal.md`

## Phase 1: Reasoning Audit

### Problem -> Solution Chain

The problem is clearly stated: validate-ux verifies code structure but not rendered visual output, creating a false-confidence blind spot. The evidence from pm-work-tracker's milestone-map feature is concrete and verifiable (gotcha document confirmed at `/Users/fanhuifeng/Projects/ai/pm-work-tracker/docs/lessons/gotcha-validation-ux-misses-visual-gaps.md`).

The solution proposes adding a "Design Compliance" collection phase using accessibility tree + computed style matching against ui-design.md requirements. The chain is: read design spec -> render pages -> capture accessibility tree + computed style -> match against spec -> generate pass/fail report.

**Gap in the chain**: The evidence contains 4 visual problems. The Capability Boundary section honestly admits the solution can only cover 2-3 of those 4 types. This is a positive transparency signal but also means the solution only partially addresses the stated problem. The proposal acknowledges this but does not quantify what percentage of *future* visual issues it expects to catch beyond this single sample.

### Solution -> Evidence Chain

The solution (accessibility tree + computed style matching) was clearly shaped by the pre-revision feedback. The two-layer matching strategy (rule-based high confidence + LLM low confidence) addresses the core algorithm gap identified in iteration 0. The Capability Boundary table is honest about detection limits.

**Remaining gap**: The solution's primary value proposition is "design compliance vs ui-design.md," but ui-design.md itself is acknowledged as unstructured free text. The constraint section admits "结构化提取精度取决于文档质量" — this is a critical dependency that undermines the entire approach but is treated as a minor constraint.

### Evidence -> Success Criteria Chain

SC-6 requires covering at least 2 of 4 gotcha problem types. This is a low bar — it means the solution can fail to detect half the known problem types and still meet success criteria. Furthermore, SC-6 only measures against the single pm-work-tracker incident, not against a broader taxonomy of visual defects.

### Self-Contradiction Check

1. **Solution claims "pass/fail" structured output but depends on free-text ui-design.md parsing.** The two-layer matching strategy mitigates this partially, but the rule layer requires the design doc to contain machine-parseable component type information. Without a schema, rule layer coverage will be narrow.

2. **Innovation section claims "不依赖 golden image 维护" as a differentiator.** This is true, but the solution trades one maintenance burden for another: it depends on ui-design.md quality instead. This is acknowledged in constraints but not in the innovation comparison.

3. **NFR says "不应使 validate-ux 的执行时间增加超过 50%."** But the solution adds per-route screenshot + accessibility tree extraction + computed style collection + LLM analysis. No evidence or estimation is provided that 50% is achievable. The feasibility section mentions "4-5 files" but does not address runtime complexity at all.

### SC Consistency Deep-Dive

Clustered by affected area:

**validate-ux pipeline (SC-1, SC-4)**: SC-1 requires 100% route coverage for projects with ui-design.md. SC-4 requires zero behavioral change for projects without ui-design.md. These are internally consistent — the pipeline branches on ui-design.md presence.

**Report format (SC-2, SC-5)**: SC-2 requires every ui-design.md requirement maps to pass/fail with screenshot path, accessibility tree node reference, and confidence level. SC-5 requires every fail entry includes failure reason, expected value, actual value, screenshot path, and confidence level. These are consistent — SC-5 is a stricter subset of SC-2's fail entries.

**Rubric update (SC-3)**: Requires 4 sub-dimensions with scoring criteria and deduction rules. This is standalone and consistent with other SCs.

**Coverage claim (SC-6)**: Requires at least 2 of 4 gotcha types. This is achievable per the Capability Boundary analysis. However, "并在报告中显式列出未覆盖的问题类型" creates an obligation that should have a corresponding SC entry for the report template — minor gap but not a contradiction.

**One potential tension**: SC-1 says "100% 的路由" but what if a route has no corresponding ui-design.md section? This edge case is not addressed. If ui-design.md describes 5 pages but sitemap.json has 8 routes, the remaining 3 routes are in a grey zone — not covered by SC-1's 100% claim, not covered by SC-4's backward compatibility (since ui-design.md does exist).

## Phase 2: Rubric Scoring

### Dimension 1: Problem Definition — 82/110

**Problem stated clearly (35/40)**: The problem is specific and unambiguous — validate-ux reports "all pass" but visual defects exist. Two readers would interpret this the same way. Minor deduction: the problem scope conflates "visual design compliance" with "rendering validation" — the title says "visual design compliance" but the problem is really about "rendered output validation," which is broader.

**Evidence provided (32/40)**: Concrete evidence from pm-work-tracker with specific issues listed and reference to gotcha document. Verified that the gotcha document exists and corroborates the claims. However, the evidence is from a *single* project and a *single* feature. The urgency section claims "随着 forge 承担更多 web 项目的自动化开发，这个问题的影响频率会持续增加" — this is a projection without supporting data (how many web projects? how frequently do visual issues occur across them?).

**Urgency justified (15/30)**: The urgency argument is hand-waving. "每次 web 功能迭代都可能产生同样的视觉盲区" is speculative. No data on: how many web projects forge handles, what percentage have ui-design.md, how often visual issues occur. The only hard data point is one project, one feature. The cost of delay is not quantified — what is the actual impact of a visual defect reaching production in forge's target use case?

### Dimension 2: Solution Clarity — 85/120

**Approach is concrete (35/40)**: The four-step process is clearly described. A reader can explain back what will be built: read ui-design.md -> render + capture -> compare -> report. The two-layer matching strategy adds specificity. Minor gap: step 2 says "复用现有 agent-browser 基础设施" but agent-browser is not found as a named component in the forge plugin — the pipeline rules reference "agent-browser" as an execution method but there is no dedicated agent-browser module. This is a terminology precision issue.

**User-facing behavior described (30/45)**: The report format is described (pass/fail with confidence levels, screenshot paths, accessibility tree references). However, the actual user experience is underspecified: what does the user *see* when design compliance fails? How does it integrate into the existing validate-ux output? Is it a separate section, a merged report? SC-2 says "ux-snapshot.md 中新增 design-compliance section" — this is about the artifact format, not the user experience of reading the results.

**Technical direction clear (20/35)**: The two-layer strategy (rule + LLM) is described but implementation details are thin. The rule layer depends on a "组件 role 映射表" that does not exist yet — who creates it? How many rules does it need to be useful? The LLM layer is described as fallback but no detail on prompt design, context window management, or how to prevent hallucinated matches. The computed style collection is mentioned in Innovation Highlights but not in the main solution steps — it appears only in the Capability Boundary table.

### Dimension 3: Industry Benchmarking — 72/120

**Industry solutions referenced (30/40)**: Percy, Chromatic, and Playwright screenshot comparison are mentioned. These are real tools. However, the references are superficial — no version numbers, no specific features analyzed, no pricing/licensing constraints discussed. The comparison is at a high level ("pixel-level comparison needs baselines").

**At least 3 meaningful alternatives (18/30)**: Four alternatives are listed including "do nothing." However, "跨页面一致性审计" is presented as a self-invented option without industry precedent — is this a novel idea or does it exist elsewhere? The Percy/Chromatic option is a genuine industry-validated solution. "Do nothing" is appropriate. But the fourth option (the selected one) is also self-invented. So only 1 out of 4 alternatives is industry-validated. Missing: visual snapshot testing tools like BackstopJS, Applitools, or browser-level assertion libraries like jest-image-snapshot. Missing: accessibility testing tools (axe-core, pa11y) which already work with accessibility trees and could be partially relevant.

**Honest trade-off comparison (14/25)**: The comparison table has clear verdicts but the pros/cons are cherry-picked. Percy's con is "任何细微变化触发 diff，误报率高" — but Percy has configurable diff thresholds, and Chromatic supports change approval workflows. The cons are overstated to favor the selected approach. The selected approach's con ("依赖 ui-design.md 质量") is understated — it should be a showstopper-level concern given the free-text format.

**Chosen approach justified (10/25)**: The selection rationale is "直接对标 gotcha 缺口，复用现有基础设施." This is pragmatic but not a rigorous justification. The proposal does not explain why a hybrid approach (e.g., accessibility tree checks + lightweight screenshot diff for visual properties) was not considered. The selected approach leaves the hardest problem (visual detail detection) unsolved while adding complexity for partial coverage.

### Dimension 4: Requirements Completeness — 72/110

**Scenario coverage (28/40)**: Five scenarios are listed covering happy path, component missing, variant inconsistency, layout deviation, and no-ui-design.md fallback. Missing scenarios: (1) ui-design.md exists but is empty or malformed, (2) sitemap.json has routes not mentioned in ui-design.md (the grey zone identified in SC consistency check), (3) accessibility tree extraction fails (browser crash, timeout), (4) computed style collection fails for specific elements, (5) rule layer and LLM layer produce conflicting results for the same component.

**Non-functional requirements (25/40)**: Three NFRs are stated: compatibility, performance, observability. Missing NFRs: (1) Reliability — what happens when LLM returns inconsistent results across runs? (2) Accuracy — no target for false positive/negative rates. (3) Security — screenshots may contain sensitive data from development environments. The 50% performance constraint is stated without justification or measurement methodology.

**Constraints & dependencies (19/30)**: Four dependencies listed. The ui-design.md quality dependency is acknowledged but its impact is understated. Missing constraints: (1) browser/headless Chrome version compatibility, (2) accessibility tree format differences across browsers, (3) computed style values that differ between rendering engines, (4) LLM API availability and latency as a dependency for the low-confidence layer.

### Dimension 5: Solution Creativity — 62/100

**Novelty over industry baseline (28/40)**: The accessibility tree + computed style approach is genuinely different from pixel-level screenshot comparison. Borrowing from web accessibility domain is a creative insight. However, the innovation is primarily about *what to compare* (semantic structure instead of pixels), not about *how to compare* — the comparison itself relies on standard rule matching and LLM fallback.

**Cross-domain inspiration (22/35)**: The accessibility tree insight is borrowed from the accessibility testing domain, applied to design compliance. This is valid cross-domain inspiration. However, no other cross-domain ideas are explored — e.g., compiler AST diff techniques for structural comparison, or contract testing patterns from API testing applied to UI components.

**Simplicity of insight (12/25)**: The solution is *almost* elegant — "use accessibility tree as structured representation of visual output" is a clean idea. But the implementation requires two layers (rule + LLM), computed style supplementation, and still cannot cover 40-50% of the evidence problems. The elegance is undermined by the complexity needed to make it work partially.

### Dimension 6: Feasibility — 68/100

**Technical feasibility (30/40)**: The claim that "validate-ux 管线已有 web surface 支持" is verifiable — the pipeline rules file confirms agent-browser, sitemap.json, and accessibility tree extraction are existing capabilities. The computed style collection via `getComputedStyle()` is technically straightforward. The main risk is LLM reliability for the low-confidence layer — the proposal does not address how to handle LLM non-determinism in a CI/CD pipeline.

**Resource & timeline (18/30)**: "涉及 4-5 个文件的修改" is plausible but vague. No timeline estimate. The claim "属于中等规模改动" is self-assessed without evidence. The LLM layer requires prompt engineering, testing, and iteration — this is non-trivial work not reflected in the "4-5 files" estimate.

**Dependency readiness (20/30)**: agent-browser and sitemap.json are confirmed ready. ui-design.md "由现有 /ui-design skill 生成" — but the ui-design.md template is free-text, which is the core dependency risk. The proposal treats this as "ready" when it is actually "ready but inadequate."

### Dimension 7: Scope Definition — 60/80

**In-scope items are concrete (22/30)**: Five in-scope items are specific deliverables (pipeline step, rubric dimension, rule file update, snapshot format update, expert update). Item 1 ("读取 ui-design.md → 提取 UI 要求 → 逐路由对比渲染结果 → 生成 design-compliance 报告") is actually a multi-step workflow described as a single item — could be decomposed further.

**Out-of-scope explicitly listed (22/25)**: Six out-of-scope items are clearly named. The exclusion of "跨页面一致性审计" is well-argued in the alternatives section. Item 6 ("task 层面的 validation.ux 改动——仅增强 eval 层面的 validate-ux") is a useful boundary.

**Scope is bounded (16/25)**: The scope is bounded to validate-ux pipeline enhancement. However, the "未来可引入可选的 machine-readable 标注格式（如 @design-spec JSON 块）" in the constraints section is scope creep speculation that should either be in-scope or explicitly out-of-scope. Ambiguous scope boundaries invite implementation divergence.

### Dimension 8: Risk Assessment — 62/90

**Risks identified (20/30)**: Five risks listed. The most critical risk — LLM non-determinism affecting pass/fail reproducibility — is missing. Also missing: the risk that the rule mapping table proves insufficient for real-world component variety, requiring constant maintenance. The risks listed are real but skewed toward external factors rather than internal solution weaknesses.

**Likelihood + impact rated (20/30)**: Ratings are reasonable. The "ui-design.md 质量不足" risk is rated M/H — arguably should be H/H given it is the foundational dependency. The "evidence 中部分问题类型超出检测能力" is rated H/M — the pre-revised annotation adds honesty but the likelihood should perhaps be H given the Capability Boundary already acknowledges this.

**Mitigations are actionable (22/30)**: Most mitigations are actionable: "标注'无法对比'的条目", "显式列出未覆盖的问题类型." However, "补充 computed style 采集关键样式属性" is a design choice, not a mitigation — it is already part of the solution. The agent-browser instability mitigation ("复用现有 validate-ux 的 dev→probe→test 编排，已有重试机制") is valid but does not address the new failure modes introduced by design compliance.

### Dimension 9: Success Criteria — 52/80

**Criteria are measurable and testable (20/30)**: SC-1 ("100% 的路由生成 design compliance 采集数据") is measurable. SC-4 ("行为与当前完全一致") is testable via comparison. SC-5 specifies exact fields required in fail entries — testable. However, SC-6 ("至少覆盖 2 类") is a low bar that is easily met but proves little. SC-3 ("包含至少 4 个子维度") is about counting, not about quality of those sub-dimensions.

**Coverage is complete (17/25)**: SCs cover report format, pipeline behavior, rubric update, and coverage claims. Missing: no SC for the rule mapping table quality (how many rules? what coverage?), no SC for false positive rate, no SC for the LLM layer accuracy, no SC for performance (the NFR says 50% but there is no corresponding SC enforcing it).

**SC internal consistency (15/25)**: As identified in Phase 1, SC-1's "100% 路由" claim has an edge case when routes in sitemap.json exceed pages described in ui-design.md. SC-2 requires "每条 ui-design.md UI 要求映射到 pass/fail" — but if LLM layer produces a result, is that a "pass/fail" or a "low-confidence pass/fail"? The SC does not distinguish. SC-4 says "无额外耗时" for projects without ui-design.md — but the pipeline still needs to *check* for ui-design.md existence, which adds minimal but non-zero overhead.

### Dimension 10: Logical Consistency — 62/90

**Solution addresses the stated problem (22/35)**: The problem is "validate-ux does not verify rendered visual output." The solution adds rendered output verification — but only partial. The Capability Boundary admits colors, fonts, shadows, and precise layout positions cannot be detected. The solution addresses the problem *directionally* but not *completely*. The gotcha evidence shows 4 visual problems; the solution covers 2-3. This is a significant gap between problem statement and solution coverage.

**Scope <-> Solution <-> Success Criteria aligned (20/30)**: Generally aligned. In-scope items map to SCs. However, In-Scope item 2 mentions "布局结构" as a rubric sub-dimension, but the Capability Boundary says layout position deviation has "低" confidence. Scoring layout structure in the rubric when the detection capability is low-confidence creates a misalignment — the rubric may deduct points based on unreliable data.

**Requirements <-> Solution coherent (20/25)**: Requirements and solution are mostly coherent. The "组件变体不一致" scenario maps to the rule + computed style approach. The "布局偏差" scenario maps to computed style. One orphan: the "刷新按钮无反馈" problem from the gotcha is neither in the problem scope (correctly excluded as non-visual) nor in the solution, but the proposal does not explicitly address why this type of issue is excluded from the problem redefinition.

## Phase 3: Blindspot Hunt

### Blindspot 1: No rollback plan

The proposal has no rollback or kill-switch mechanism. If design compliance produces excessive false positives in production use, there is no documented way to disable it per-project or globally. The SC-4 backward compatibility only covers "no ui-design.md" projects — but what about projects *with* ui-design.md where the feature causes problems?

### Blindspot 2: LLM cost and latency unaddressed

The LLM layer is treated as a minor fallback, but for projects with complex ui-design.md files (many components, many routes), LLM calls could be numerous. No cost estimation, no latency budget, no fallback when LLM API is unavailable. In a CI pipeline, LLM API failures would block validation.

### Blindspot 3: The "at least 2 of 4" success criterion normalizes mediocrity

SC-6 sets the bar at 50% coverage of known problem types. This means the feature can ship while failing to detect half the problems that motivated its creation. There is no path defined for improving from 2 to 4 types.

### Blindspot 4: No incremental rollout strategy

The proposal jumps directly to "add to validate-ux pipeline." No mention of: opt-in alpha period, shadow mode (run but don't affect scoring), or gradual rollout. Given the solution's acknowledged limitations, a phased approach would be prudent.

### Blindspot 5: Accessibility tree browser inconsistency

The proposal does not address that accessibility trees differ across browsers (Chrome vs Firefox vs Safari). If agent-browser uses headless Chrome (likely), the accessibility tree may differ from what developers see in their browser. This is not mentioned as a constraint.

### Blindspot 6: The rule mapping table is an uncosted deliverable

The rule layer requires "组件 role 映射表" — this is a new artifact that needs creation, maintenance, and testing. It is not listed as an in-scope deliverable, not estimated in the resource assessment, and has no corresponding SC.

## Bias Detection Report

- Annotated regions: 5 attack points / 8 paragraphs = density 0.625
- Unannotated regions: 18 attack points / 19 paragraphs = density 0.947
- Ratio (annotated/unannotated): 0.66

The annotated (pre-revised) regions show lower attack density than unannotated regions. This is expected — pre-revision addressed the most severe issues (accessibility tree expressiveness, matching algorithm, capability boundary). The remaining attacks on annotated regions focus on whether revisions introduced new issues: the two-layer matching strategy (annotated) is structurally sound but under-specified; the computed style addition (annotated) is mentioned in Capability Boundary but not in the main solution steps.

## Conflict-with-Pre-Revision Tags

None identified. The scorer's findings align with or are orthogonal to the pre-revision direction — no cases where the scorer recommends deleting what pre-revision added, or where the scorer disagrees with the pre-revision assessment's severity.
