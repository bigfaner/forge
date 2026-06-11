---
iteration: 1
title: "CTO Rubric Evaluation"
total_score: 640
rubric_total: 1000
---

# CTO Rubric Evaluation — Iteration 1

## Dimension Scores

| # | Dimension | Score | Max | Verdict |
|---|-----------|-------|-----|---------|
| 1 | Problem Definition | 85 | 110 | Strong problem, evidence has factual error |
| 2 | Solution Clarity | 75 | 120 | Approach clear, user-facing behavior underspecified |
| 3 | Industry Benchmarking | 55 | 120 | Shallow benchmarks, straw-man alternatives |
| 4 | Requirements Completeness | 80 | 110 | Good scenarios, NFRs thin |
| 5 | Solution Creativity | 60 | 100 | Modest innovation, no cross-domain inspiration |
| 6 | Feasibility | 70 | 100 | Acknowledges complexity but underestimates it |
| 7 | Scope Definition | 60 | 80 | In-scope items lack deliverable specificity |
| 8 | Risk Assessment | 75 | 90 | Good risk identification, mitigations partially actionable |
| 9 | Success Criteria | 50 | 80 | Measurability gaps, incomplete coverage, consistency issues |
| 10 | Logical Consistency | 30 | 90 | SC-Scope-Solution misalignment, unacknowledged scope creep |

**Total: 640 / 1000**

---

## Detailed Scoring

### 1. Problem Definition (85 / 110)

**Problem stated clearly (35 / 40)**

The core problem is unambiguous: freeform experts generated per-proposal have near-zero reuse. The Jaccard 0.3 threshold provides a concrete, measurable failure condition. One minor ambiguity: "Jaccard 相似度无法达到 0.3 的复用阈值" conflates pairwise similarity with the reuse threshold — the threshold applies to the best-match candidate, not to pairwise averages. This is a presentation clarity issue, not a logical flaw.

**Evidence provided (25 / 40)**

Evidence points to real files (`docs/experts/`), real keywords, and real behavior ("每次评估都触发了全新专家生成"). However, there is a factual error: the document states "11 个专家文件" but `docs/experts/` contains 12 files. The 12th (`expert-system-design-prompt-architecture.md`) was generated for this very proposal, which actually strengthens the argument (even the meta-proposal generates a new expert) — but the inaccurate count raises credibility concerns.

The domain keyword evidence ("5 个关键词中任何一个出现在其他 proposal 的概率极低") is qualitative assertion without computation. The author could have computed actual pairwise Jaccard scores for the 12 experts to make this concrete — the data is available.

**Urgency justified (25 / 30)**

The urgency argument is clear and honest: "每次评估都经历完整的专家生成→确认循环（3 轮修改/拒绝上限），增加评审耗时". The cost is real (time per evaluation). The claim "显著降低评审启动成本" is reasonable for the domain-level approach. One missing element: no quantification of the actual time cost (e.g., "each expert generation adds ~5 minutes to evaluation").

---

### 2. Solution Clarity (75 / 120)

**Approach is concrete (35 / 40)**

The two-step approach (domain classification → expert generation within domain) is clearly described. The pre-revised addition of the empirical clustering analysis (Section "领域分布实证分析") substantially improves concreteness by showing the 12 experts map to 6 natural domains with a projected ceiling of 8-10. A reader can explain back what will be built.

**User-facing behavior described (20 / 45)**

This is the weakest aspect. The proposal describes internal mechanism changes (classification table, scope field, persistence logic) but provides minimal description of what the end user (the person running evaluations) actually experiences differently. Key gaps:

- What does the user see when a domain match is found vs. not found?
- How does the "Accept / Modify / Regenerate" flow change for domain-level experts?
- What does the user confirmation prompt look like with the new `scope` field?
- Scenario 3 mentions "用户可选择仅使用最相关的一个大领域专家" — how is this choice presented?

The proposal optimizes for internal architecture clarity at the expense of observable behavior specification. This is a significant gap for a proposal that modifies a user-facing workflow.

**Technical direction clear (20 / 35)**

The general approach is clear (classification table in prompt, scope field in template, modified persistence logic), but the technical details are thin on the most critical part: how does the classification table actually get embedded in `expert-inference.md`? Is it a hard-coded list in the prompt? A referenced external file? The proposal says "只需修改 prompt 文件" (Constraints) but never shows or describes the classification table format or embedding strategy.

The dependency chain between the 4 files is well-documented (template schema is prerequisite for inference and persistence, extraction depends on scope field), which partially compensates.

---

### 3. Industry Benchmarking (55 / 120)

**Industry solutions referenced (15 / 40)**

The references are shallow. "学术同行评审的 TPC 成员列表" and "ChatGPT 的 Custom Instructions persona" are named but not cited with any depth — no papers, no OSS projects, no published patterns. The comparison to "Claude Code 的 multi-expert parallel scoring" is self-referential (this is the system being modified, not an industry benchmark).

Missing references that would be relevant:
- Academic literature on expert system knowledge organization (hierarchical vs. flat taxonomies)
- Industry classification/taxonomy systems (e.g., ACM Computing Classification System, Stack Overflow tag taxonomy)
- Published research on LLM role assignment consistency

**At least 3 meaningful alternatives (20 / 30)**

Four alternatives are presented including "do nothing". However, two of the four are effectively straw men:

1. "Do nothing" — dismissed in one line ("核心问题不会自行消失"). This is not a genuine engagement with the alternative; it's a token mention.
2. "纯 Prompt 重写" — described as "改动最小" then dismissed as "没解决一致性问题" without explaining what "prompt rewrite" means or why it fails. The proposal's own solution is essentially a prompt rewrite with a classification table embedded — making this alternative suspiciously close to the chosen approach but presented as inadequate.

The "固定专家库" alternative is more genuine but still lacks depth — no analysis of what a fixed expert library would look like for this project's domain.

**Honest trade-off comparison (10 / 25)**

The comparison table is honest in acknowledging "分类表需维护" as a con for the selected approach. However, the pros/cons for other approaches are cherry-picked:

- "纯 Prompt 重写" cons list is just "领域标签不一致" — but the selected approach also relies on LLM inference within a domain, which could produce inconsistent results within that domain.
- "固定专家库" is dismissed as "无法适应未知领域" — but the proposal itself includes a degradation path (LLM free inference) for unclassified domains, which is functionally equivalent to having a fixed library with a fallback.

**Chosen approach justified against benchmarks (10 / 25)**

The justification is "一致性与灵活性最优平衡" — a slogan rather than a reasoned argument. The proposal does not explain why "classification table + LLM inference within domain" is superior to "improved prompt instructions for consistent keyword generation" or "hierarchical keyword taxonomy without classification table." The innovation claim is moderation between two extremes, but the middle ground is not uniquely justified.

---

### 4. Requirements Completeness (80 / 110)

**Scenario coverage (35 / 40)**

Four scenarios are identified: new domain first evaluation, same-domain reuse, cross-domain proposal, and unclassified domain. This covers the main paths well. Edge cases that are missing:

- What happens when a domain-level expert is auto-deprecated (per existing persistence rules)?
- What happens when a proposal's domain classification is ambiguous (equally close to two categories)?
- What happens when a domain-level expert was generated for a different proposal's needs and doesn't cover the current proposal's specific concerns?

The pre-revised addition of Scenario 3 (cross-domain) with explicit trade-off acknowledgment is valuable.

**Non-functional requirements (20 / 40)**

Two NFRs are listed: extensibility and backward compatibility. Both are relevant but thin:

- "可扩展性：领域分类表应易于扩展" — no metric for how easy is "easy". One new domain requires modifying the prompt file — is that "easy"?
- "向后兼容" — stated as "现有 11 个 proposal-specific 专家文件保留不变" but the `scope` field isolation means old experts will never match new domain-level experts, which is de facto deprecation without explicit deprecation. This is backward compatibility in letter but not in spirit.

Missing NFRs:
- **Performance**: Will the two-step inference add latency to evaluations?
- **Quality**: How to measure whether domain-level expert reviews are as effective as proposal-specific ones?
- **Consistency**: How to verify that the same domain classification is produced for similar proposals?

**Constraints & dependencies (25 / 30)**

Constraints are well-specified: "改动仅限 `experts/freeform/` 目录下的 prompt 文件和 `rules/freeform-expert-persistence.md`". The dependencies are clear (no external deps, all within Forge plugin). However, the pre-revised Feasibility section added `extraction-prompt.md` as a 4th file — this expands the constraint boundary beyond what was stated here. The Scope section was also updated to include `extraction-prompt.md`, but the Constraints section was not reconciled to reflect the expanded scope.

---

### 5. Solution Creativity (60 / 100)

**Novelty over industry baseline (25 / 40)**

The pre-revised Innovation Highlights section is commendably honest: "大多数专家系统要么用固定专家库（无灵活性），要么完全依赖 LLM 自由推断（无一致性）——分层方案在两者之间取得平衡". The novelty is moderate — it's a known pattern (constrained generation space) applied to a new domain (LLM expert profile generation). The acknowledgment that "LLM 仍可能将边界模糊的 proposal 映射到不同条目" shows intellectual honesty but also undermines the uniqueness claim.

**Cross-domain inspiration (10 / 35)**

No cross-domain inspiration is identified. The proposal stays entirely within the LLM/prompt engineering domain. No borrowings from:
- Library science (classification systems, faceted search)
- Recommendation systems (collaborative filtering for expert matching)
- Database indexing (hierarchical clustering for retrieval)
- Organization design (role-based vs. skill-based team assignment)

**Simplicity of insight (25 / 25)**

The core insight — "constrain the search space for domain labeling to improve matching consistency" — is genuinely simple and elegant. It's the kind of insight that seems obvious in retrospect, which is a hallmark of good design. The pre-revised analysis that re-frames this from "absolute guarantee" to "probability improvement" is appropriately scoped.

---

### 6. Feasibility (70 / 100)

**Technical feasibility (30 / 40)**

The pre-revised section is honest about the 4-file coordination complexity and the dependency chain. The analysis that "template 的 schema 变更是 inference 和 persistence 的前提" shows understanding of the change propagation. However, the feasibility assessment underestimates one risk: the `scope` field creates a de facto schema versioning problem. All downstream consumers of expert files (not just the 4 named files — any skill that reads `docs/experts/`) must handle both `scope`-present and `scope`-absent experts. The proposal assumes only 4 files are affected without auditing all consumers.

**Resource & timeline feasibility (15 / 30)**

"单次 prompt 链协调修改 + 交叉验证，工作量中等" — this is too vague to evaluate. No timeline estimate, no team size assumption, no comparison to similar past changes. "中等" is subjective and unsubstantiated. Deducted for lack of specificity.

**Dependency readiness (25 / 30)**

"No external dependencies. All changes within Forge plugin." This is verifiable and accurate based on the codebase structure. The only dependency is internal consistency between the 4 prompt files, which the proposal acknowledges.

---

### 7. Scope Definition (60 / 80)

**In-scope items are concrete (20 / 30)**

The 4 in-scope items are partially concrete:

- "`expert-inference.md` 嵌入领域分类表，改造为两步生成流程" — concrete enough to be a deliverable.
- "`expert-template.md` 增加 `scope` 字段" — concrete.
- "`freeform-expert-persistence.md` 更新复用匹配逻辑" — less concrete. What exactly changes? "区分 scope 级别" is a design goal, not a deliverable.
- "`extraction-prompt.md` 更新 findings 提取逻辑以感知专家 scope 级别" — the pre-revised addition is specific about the *why* ("domain-level findings 的适用范围更广，需正确加权") but not the *what* — what specific changes to the extraction prompt are needed?

**Out-of-scope explicitly listed (20 / 25)**

Three items are explicitly listed as out of scope. This is adequate. However, one notable omission: the classification table itself — is its initial population in scope or out of scope? The proposal implies it's part of the `expert-inference.md` change but never explicitly states who creates the initial table or what domains it contains.

**Scope is bounded (20 / 25)**

The scope is bounded to 4 files within a single directory tree. However, the boundary is slightly porous — `extraction-prompt.md` was added via pre-revision, and the claim "不涉及代码变更" in Constraints may not hold if the `scope` field requires changes to any Go code that reads expert files. The proposal assumes only prompt files are affected but hasn't audited Go consumers.

---

### 8. Risk Assessment (75 / 90)

**Risks identified (25 / 30)**

Four risks are identified, all meaningful:

1. Classification table coverage gaps (M/L)
2. Domain-level expert review depth degradation (M/M)
3. Classification table bloat over time (L/L)
4. Old proposal-specific expert noise in Jaccard matching (M/M)

This is good coverage. Missing risks:
- The `scope` field creating a schema contract that constrains future evolution
- Classification table entries becoming stale as the project's domain evolves
- The two-step inference increasing LLM token usage / latency

**Likelihood + impact rated (25 / 30)**

Ratings are reasonable. The M/L for classification gaps is honest (not everything is low). The L/L for table bloat is defensible given the 8-12 domain ceiling. One questionable rating: "旧 proposal-specific 专家持续参与 Jaccard 匹配产生噪声" is rated M/M — but the pre-revised mitigation (scope-based filtering) effectively eliminates this risk, making the M likelihood rating inconsistent with the mitigation's thoroughness.

**Mitigations are actionable (25 / 30)**

Mitigations are generally actionable:

- "提供 LLM 自由推断降级路径" — actionable, and degradation path is described in Scenario 4.
- The depth degradation mitigation is the weakest: "(1) 专家 prompt 中注入当前 proposal 全文作为评审焦点" and "(2) 用户确认循环作为最终质量保障". Mitigation (1) is concrete. Mitigation (2) is the existing system, not a new mitigation — it's "the problem will be caught by existing safeguards," which is a detection strategy, not a prevention strategy.
- "分类表控制在大领域粒度（8-12 个）" — actionable but unenforced. What mechanism prevents bloat beyond 12?
- The pre-revised mitigation for old expert noise (scope-based filtering) is concrete and actionable.

---

### 9. Success Criteria (50 / 80)

**Criteria are measurable and testable (15 / 30)**

- SC-1: "domain 关键词覆盖范围 >= 2 个 proposal 的领域交集" — partially measurable. "覆盖范围" is not precisely defined. Does it mean the expert's domain keywords must include terms from at least 2 proposals? How is "领域交集" computed? The parenthetical "用已有 proposal 交叉比对" suggests manual checking, not automated verification. **Pre-revised trade-off acknowledgment is honest but doesn't fix measurability.**
- SC-2: "复用匹配成功" — measurable and testable. Clear binary criterion.
- SC-3: "用户确认的轮次 <= 2" — measurable. However, this SC conflates the quality of domain-level experts with the user's personal threshold for acceptance. A domain-level expert might be accepted in 1 round not because it's good, but because the user has lower standards for breadth vs. depth.
- SC-4: "分类表覆盖 >= 80% 的已有 proposal" — measurable. But "已有 proposal" is undefined — does this mean the 12 proposals that already have experts? Or all proposals in `docs/proposals/`? The denominator matters.

**Coverage is complete (15 / 25)**

SC coverage gaps:

- No SC for `extraction-prompt.md` changes (added to In Scope but no corresponding SC).
- No SC for the quality/effectiveness of domain-level expert reviews compared to proposal-specific ones. The proposal acknowledges this trade-off in Key Risks but has no SC to measure it.
- No SC for classification accuracy (how often the LLM correctly assigns a proposal to its domain).

The SC-InScope alignment gap for `extraction-prompt.md` is a concrete omission.

**SC internal consistency (20 / 25)**

SC set is internally satisfiable — no logical contradictions. SC-1 (coverage breadth) and the acknowledged depth trade-off create tension but not contradiction because SC-3 (confirmation rounds) serves as a quality backstop.

One ambiguity: SC-1 says "覆盖范围 >= 2 个 proposal" but the domain clustering analysis shows 6 domains from 12 experts (2 per domain on average). For domains with only 1 existing expert (e.g., "Prompt 与 Agent 协议" with just 1 expert), SC-1 cannot be verified against 2 proposals. This is not a contradiction but a verification gap for niche domains.

---

### 10. Logical Consistency (30 / 90)

**Solution addresses the stated problem (20 / 35)**

The solution partially addresses the problem:

- **Addressed**: Domain-level experts with broader keywords will have higher Jaccard overlap with future proposals in the same domain, improving reuse matching.
- **Partially addressed**: The classification table improves labeling consistency, but the proposal itself acknowledges LLM may still produce inconsistent mappings for boundary cases.
- **Gap**: The root cause of zero reuse is that expert domain keywords are too narrow (proposal-specific). The solution generates broader keywords, which should help — but the proposal never analyzes *why* the current system generates narrow keywords. Is it because `expert-inference.md` instructs the LLM to focus on the specific proposal? If so, the fix might be simpler (adjust the inference prompt to generate broader keywords) than introducing a full classification table. The proposal doesn't rule out this simpler alternative.

**Scope ↔ Solution ↔ Success Criteria aligned (5 / 30)**

Significant misalignment:

1. `extraction-prompt.md` is listed In Scope but has no corresponding SC. The SC set does not cover all in-scope items.
2. The Constraints section says "改动仅限 `experts/freeform/` 目录下的 prompt 文件和 `rules/freeform-expert-persistence.md`" — but In Scope includes `extraction-prompt.md` which is in `experts/freeform/`, so technically within the constraint. However, the Feasibility section says "改动涉及 4 个 Markdown prompt 文件" while Constraints says "改动仅限...prompt 文件和...persistence.md" — the count and scope description are inconsistent across sections.
3. The Solution describes updating `expert-template.md` with a `scope` field, and In Scope confirms this. But there's no SC verifying that the `scope` field correctly differentiates domain-level from proposal-specific experts in persistence matching.
4. Out of Scope says "修改 freeform-pipeline 编排流程" — but the two-step inference (domain match → expert generate) is a pipeline modification. The proposal argues it's only a prompt change within `expert-inference.md`, but adding a two-step decision flow to a single prompt is functionally a pipeline change, even if it stays within one file.

**Requirements ↔ Solution coherent (5 / 25)**

Requirements (scenarios, NFRs) generally map to the solution, but:

- NFR "可扩展性" (easy to add new domains) has no corresponding solution mechanism. The proposal says "只需修改 prompt 文件" but doesn't describe the modification process or who performs it.
- Scenario 3 (cross-domain) describes multi-expert parallel review with merge, but the Solution section only describes single-domain expert generation. Multi-expert coordination is an orphan requirement — described in scenarios but not in the solution.
- The NFR "向后兼容" claims old experts are preserved, but the scope-based filtering effectively isolates them from the new system. This is "compatible" in the sense that old experts won't break anything, but it's not compatible in the sense that old experts remain useful — they become dead weight in `docs/experts/`.

---

## Bias Detection Report

**Annotated regions** (paragraphs with `<!-- pre-revised -->` markers):
- Attack points targeting annotated content: 8
- Annotated paragraphs: 10
- Density: 0.80

**Unannotated regions**:
- Attack points targeting unannotated content: 12
- Unannotated paragraphs: 28
- Density: 0.43

**Ratio (annotated/unannotated): 1.86**

The bias ratio of 1.86 indicates higher attack density on pre-revised regions. This is expected — pre-revision addressed the most problematic areas, which remain the densest sources of issues even after revision. However, the ratio is below 2.0, suggesting the scoring is not excessively biased toward pre-revised content. The unannotated regions still receive substantial scrutiny.

**Conflict-with-pre-revision tags**: None detected. All attack points on annotated regions align with or extend the pre-revision direction rather than contradicting it.

---

## Phase 3: Blindspot Hunt

### [blindspot-1] Classification table authority and evolution

The proposal never specifies who owns and maintains the classification table. Is it the Forge plugin maintainers? Is it auto-generated? Is it user-configurable? The assumption that "8-12 domains" is stable implies central maintenance, but no governance model is described. This is an architectural decision with long-term implications that the proposal treats as an implementation detail.

### [blindspot-2] Token cost and latency impact

The two-step inference (domain match → expert generate) will increase LLM token consumption. The current system does one inference pass; the proposed system does two (first to extract features and match domain, then to generate expert within domain). No analysis of the cost/latency impact. For a system that runs evaluations interactively, this could be significant.

### [blindspot-3] Jaccard threshold validity for domain-level experts

The current Jaccard threshold of 0.3 was calibrated for proposal-specific experts with narrow keyword sets. Domain-level experts will have broader keyword sets, which could make the Jaccard denominator larger and *reduce* match scores even when the expert is appropriate. The proposal does not analyze whether the 0.3 threshold needs adjustment for the new expert format.

### [blindspot-4] Single-point-of-failure in classification table

The classification table embedded in `expert-inference.md` becomes a single point of failure. If the table has a wrong or missing domain, *every* proposal in that domain will be misclassified or degraded. The current system distributes this risk across individual inference runs (each can produce different keywords). The proposal trades distributed risk for concentrated risk without acknowledgment.

### [blindspot-5] "12 experts" factual error compounds

The document states "11 个专家" but there are 12. The 12th expert (`expert-system-design-prompt-architecture.md`) was generated for this very proposal, which means:
1. The evidence section's claim of "11 个专家文件" is wrong
2. The clustering analysis (which lists 11 experts) missed one
3. The "6 大领域" conclusion may be inaccurate with the 12th expert included

This factual error, while minor in impact (the 12th expert fits the "评估管线架构" cluster), undermines the empirical analysis's credibility.

### [blindspot-6] Expert count as denominator for SC-4

SC-4 says "分类表覆盖 >= 80% 的已有 proposal." If the denominator is all proposals in `docs/proposals/`, the pass threshold depends on how many proposals exist — a moving target. If it's the 12 that already have experts, it's fixed but arbitrary. The proposal doesn't define this.

---

## Attack Points Summary

### Factual / High-Severity

1. **[D1] Evidence count is wrong**: "11 个专家文件" — actual count is 12. The clustering table also lists 11, missing `expert-system-design-prompt-architecture.md`.

2. **[D2] SC-Scope coverage gap**: `extraction-prompt.md` is listed In Scope but has no corresponding Success Criterion. In Scope has 4 items; SC has 4 items but none addresses extraction-prompt changes.

3. **[D3] NFR "可扩展性" is unmeasurable**: "领域分类表应易于扩展，新增领域只需修改 prompt 文件" — "应易于扩展" is vague language without quantification. What constitutes "easy"? Deducted per rubric vague language rule (-20).

4. **[D10] Scope-Solution inconsistency on pipeline modification**: Out of Scope says "修改 freeform-pipeline 编排流程" but the two-step inference flow in `expert-inference.md` is functionally a pipeline modification, even if confined to one file.

### Structural / Medium-Severity

5. **[D2] User-facing behavior underspecified**: No description of what the evaluation runner sees differently. The proposal optimizes for internal architecture at the expense of observable behavior specification.

6. **[D2] Technical direction unclear on classification table embedding**: The proposal never describes how the classification table is formatted or embedded in `expert-inference.md`. Is it a markdown table? A structured list? How large will it be?

7. **[D3] Straw-man alternatives**: "纯 Prompt 重写" is dismissed as "没解决一致性问题" but the selected approach is also a prompt rewrite with a classification table. The distinction between these two alternatives is not clearly articulated.

8. **[D3] No real industry references**: "学术同行评审的 TPC 成员列表" and "ChatGPT 的 Custom Instructions" are mentioned without depth, citation, or analysis of how they handle the specific problem of expert reuse.

9. **[D4] Missing NFR: Performance**: No analysis of latency or token cost impact from two-step inference.

10. **[D4] Backward compatibility is nominal**: "现有 11 个 proposal-specific 专家文件保留不变" — but scope-based filtering means they will never match new domain-level proposals and will never be selected. They are preserved as files but de facto deprecated.

11. **[D6] Timeline is vague**: "工作量中等" with no specific estimate. A proposal that modifies 4 coordinated files with dependency chains deserves more than "中等."

12. **[D7] Classification table ownership undefined**: No governance model for who maintains the table, how new domains are proposed/validated, or how stale domains are pruned.

13. **[D8] Depth mitigation (2) is detection, not prevention**: "用户确认循环（Accept / Modify / Regenerate）作为最终质量保障" — this is the existing system, not a new mitigation. It detects quality problems after generation rather than preventing them.

14. **[D9] SC-1 measurability**: "domain 关键词覆盖范围 >= 2 个 proposal 的领域交集" — "覆盖范围" and "领域交集" are not precisely defined operations.

15. **[D9] SC-4 denominator undefined**: "覆盖 >= 80% 的已有 proposal" — "已有 proposal" is ambiguous.

16. **[D10] Cross-domain scenario is an orphan requirement**: Scenario 3 describes multi-expert parallel review with merge, but the Solution only describes single-domain generation. Multi-expert coordination has no solution design.

### Low-Severity / Style

17. **[D3] Comparison table verdict column is editorializing**: "Rejected: 核心问题不会自行消失" for "do nothing" — while reasonable, the one-line dismissal doesn't engage with the legitimate question of whether the problem's cost justifies the solution's complexity.

18. **[D5] No cross-domain inspiration**: The proposal stays entirely within LLM/prompt engineering. Missing opportunities to cite classification science, recommendation systems, or organizational design.

19. **[D10] Constraints section not reconciled after pre-revision**: Original constraints said "改动仅限...prompt 文件和...persistence.md" but pre-revision added `extraction-prompt.md` to scope without updating constraints language.

---

## Top 3 Recommended Improvements

1. **Add a missing SC for `extraction-prompt.md` changes** and define a clear deliverable for what changes to that file are needed. This fixes the SC-Scope coverage gap.

2. **Define user-facing behavior explicitly**: Add a subsection describing what the evaluation runner experiences at each step — what prompts they see, what choices they make, what feedback they receive. This is the largest gap in Solution Clarity.

3. **Fix the factual error (12 experts, not 11) and compute actual pairwise Jaccard scores** for the existing experts as concrete evidence. This strengthens the problem statement and eliminates the credibility concern.
