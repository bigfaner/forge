---
iteration: 2
title: "CTO Rubric Evaluation"
total_score: 690
rubric_total: 1000
---

# CTO Rubric Evaluation — Iteration 2

## Dimension Scores

| # | Dimension | Score | Max | Verdict |
|---|-----------|-------|-----|---------|
| 1 | Problem Definition | 90 | 110 | Strong problem statement, evidence still has one inaccuracy |
| 2 | Solution Clarity | 88 | 120 | Substantial improvement in user-facing behavior; cross-domain strategy added |
| 3 | Industry Benchmarking | 55 | 120 | Unchanged — shallow benchmarks, straw-man alternatives persist |
| 4 | Requirements Completeness | 82 | 110 | Scenarios improved, NFRs still thin |
| 5 | Solution Creativity | 60 | 100 | Honest novelty assessment, no cross-domain inspiration |
| 6 | Feasibility | 78 | 100 | Timeline now quantified; 4-file coordination acknowledged |
| 7 | Scope Definition | 65 | 80 | In-scope items more specific; extraction-prompt deliverable still vague |
| 8 | Risk Assessment | 78 | 90 | Good risk table; cross-domain cost acknowledged |
| 9 | Success Criteria | 65 | 80 | Extraction-prompt SC added; SC-1 formula improved but denominator ambiguity remains |
| 10 | Logical Consistency | 29 | 90 | Cross-domain strategy contradicts original single-expert constraint; scope-solution gaps persist |

**Total: 690 / 1000**

---

## Iteration-1 Issue Tracker

| # | Issue from Iteration 1 | Status | Evidence |
|---|----------------------|--------|----------|
| 1 | Expert count wrong (11 vs 12) | **Fixed** | Line 15: "12 个专家文件" |
| 2 | SC-Scope gap for extraction-prompt.md | **Fixed** | SC-3 now addresses extraction-prompt findings extraction |
| 3 | User-facing behavior underspecified | **Fixed** | New "用户可观测行为" subsection (lines 33-39) |
| 4 | NFR "可扩展性" unmeasurable | **Not fixed** | Line 76: "应易于扩展" remains vague |
| 5 | Timeline vague ("工作量中等") | **Fixed** | Lines 117-119: "预估工作量 2-3 小时" with analogy-based breakdown |
| 6 | Cross-domain scenario is orphan requirement | **Partially fixed** | New "跨领域多专家策略" section (lines 42-44), but creates new consistency issues (see D10) |
| 7 | No real industry references | **Not fixed** | Same shallow references as iteration 1 |
| 8 | Constraints section not reconciled | **Partially fixed** | Scope section updated, but Constraints line 81 still says "改动仅限 `experts/freeform/` 目录下的 prompt 文件和 `rules/freeform-expert-persistence.md`" — missing extraction-prompt.md from the constraint enumeration |

---

## Detailed Scoring

### 1. Problem Definition (90 / 110)

**Problem stated clearly (38 / 40)**

The core problem is unambiguous: freeform experts generated per-proposal have near-zero reuse, measured by Jaccard similarity failing the 0.3 threshold. The statement is precise and a reader cannot misinterpret the failure condition.

One residual inaccuracy: line 15 claims "每个的 `generated_for` 都指向唯一的 proposal" — but `build-orchestration-test-infra.md` and `surface-aware-dispatcher-orchestrator.md` both have `generated_for: "docs/proposals/surface-aware-justfile/proposal.md"`. The 12 experts map to 11 unique proposals, not 12. This does not undermine the argument (zero reuse is still zero reuse), but it is a factual error in the evidence section.

**Evidence provided (27 / 40)**

The evidence points to real files, real keywords, and real behavior. The iteration-2 fix of "12 个专家文件" corrects the previous count error. However:

- "Jaccard 相似度无法达到 0.3 的复用阈值" remains a qualitative assertion — no actual pairwise Jaccard scores are computed from the 12 experts' domain keywords. The data is available; the computation is straightforward; the author chose not to do it.
- "实际复用匹配从未成功过——每次评估都触发了全新专家生成" is a behavioral claim that could be verified from eval logs but isn't cited with specific log references.
- The "generated_for 指向唯一 proposal" claim is factually wrong (see above).

**Urgency justified (25 / 30)**

"每次评估都经历完整的专家生成→确认循环（3 轮修改/拒绝上限），增加评审耗时" — the cost is real and concrete. The claim "显著降低评审启动成本" is reasonable. Still no quantification of actual time cost per evaluation cycle.

---

### 2. Solution Clarity (88 / 120)

**Approach is concrete (38 / 40)**

The two-step approach (domain classification via lookup table → expert generation within domain scope) is clearly described. The empirical clustering table (6 domains from 12 experts, projected ceiling 8-10) provides concrete grounding. The dependency chain between 4 files is well-documented. A reader can explain back what will be built.

**User-facing behavior described (30 / 45)**

The new "用户可观测行为" subsection (lines 33-39) is a significant improvement over iteration 1. Four observable behaviors are described:

1. Expert generation confirmation prompt shows `scope: domain-level [大领域名称]`
2. Reuse hit logs output `Reusing existing domain-level expert`
3. Cross-domain proposals show interactive domain selection
4. Findings reports annotate scope level

Gaps remaining:

- Behavior 3 ("交互式提示呈现") references a UI/interaction pattern ("使用全部匹配领域" or "仅使用最相关领域") but the "跨领域多专家策略" section (line 44) describes sequential review, not interactive selection. These two descriptions of the same cross-domain workflow are inconsistent.
- No description of what happens when the user selects "Modify" on a domain-level expert — does the modification scope stay within the matched domain, or can the user redirect to a different domain?

**Technical direction clear (20 / 35)**

The general approach remains clear but the most critical technical question is unanswered: how is the classification table embedded in `expert-inference.md`? The proposal says "嵌入领域分类表" (line 136) but never describes the table format (markdown table? YAML? structured list?), its expected size, or how it interacts with the LLM inference prompt. This is the architectural centerpiece of the proposal, and its implementation format is left to the imagination.

---

### 3. Industry Benchmarking (55 / 120)

**Industry solutions referenced (15 / 40)**

Unchanged from iteration 1. The references remain shallow:

- "学术同行评审的 TPC 成员列表" — named but not cited
- "ChatGPT 的 Custom Instructions persona" — named but not analyzed
- "Claude Code 的 multi-expert parallel scoring" — self-referential

No papers, no OSS projects, no published patterns were added.

**At least 3 meaningful alternatives (20 / 30)**

Unchanged. Four alternatives presented, but:

- "Do nothing" dismissed in one line: "核心问题不会自行消失" — this is a tautology, not an engagement with the alternative.
- "纯 Prompt 重写" — the comparison table now provides a longer explanation (line 99) distinguishing it from the selected approach ("从开放式生成约束为从有限候选项中选取"), which is a genuine attempt to differentiate. However, the distinction is about *degree* of prompt engineering, not fundamentally different approaches. The selected approach is a specialized form of prompt rewriting.
- "固定专家库" — dismissed as "不够灵活" without analysis.

**Honest trade-off comparison (10 / 25)**

Unchanged. The comparison table acknowledges "分类表需维护" as a con. Cherry-picking remains:

- "纯 Prompt 重写" cons: "无法约束 LLM 的标签一致性" — but the selected approach also cannot *guarantee* consistency; it only improves the probability. The Innovation Highlights section (line 62) honestly admits this: "LLM 仍可能将边界模糊的 proposal 映射到不同条目".
- "固定专家库" cons: "无法适应未知领域" — but the selected approach has a degradation path (Scenario 4: LLM free inference for unclassified domains), which is functionally equivalent to a fixed library with a fallback. The proposal doesn't acknowledge this symmetry.

**Chosen approach justified against benchmarks (10 / 25)**

Unchanged. "一致性与灵活性最优平衡" remains a slogan. No quantitative or structural argument for why the classification table approach outperforms alternatives in this specific project context.

---

### 4. Requirements Completeness (82 / 110)

**Scenario coverage (37 / 40)**

Four scenarios plus cross-domain strategy. Good coverage of main paths. The cross-domain scenario (Scenario 3) now explicitly acknowledges the cost trade-off ("评审成本与领域数成正比"). Edge cases still missing:

- What happens when a domain-level expert's quality degrades over time (stale expertise as the project evolves)?
- What happens when classification is ambiguous (equally close to two categories) — Scenario 3 addresses multi-match but not ambiguous single-match.

**Non-functional requirements (20 / 40)**

Two NFRs: extensibility and old-expert isolation. Both remain thin:

- "可扩展性：领域分类表应易于扩展，新增领域只需修改 prompt 文件" — "应易于扩展" is vague language. Deduct per rubric vague language rule (-20 pts applied to this sub-criterion).
- "旧专家隔离" — honestly described as "隐性废弃：文件仍存在但永远不会被新的 domain-level 评估选中" — this is a refreshingly honest characterization but it means backward compatibility is nominal, not substantive.

Missing NFRs (unchanged from iteration 1):
- **Performance**: No analysis of latency or token cost from two-step inference.
- **Quality**: No metric for whether domain-level expert reviews are as effective.
- **Consistency**: No metric for classification accuracy across similar proposals.

**Constraints & dependencies (25 / 30)**

Line 81: "改动仅限 `experts/freeform/` 目录下的 prompt 文件和 `rules/freeform-expert-persistence.md`" — this enumeration does not explicitly name `extraction-prompt.md`, though it is within `experts/freeform/`. The Feasibility section (line 108) correctly identifies 4 files. The mismatch between the Constraints enumeration (2 locations) and the actual scope (4 files in those locations) creates confusion about what exactly is being changed.

---

### 5. Solution Creativity (60 / 100)

**Novelty over industry baseline (25 / 40)**

The Innovation Highlights section (lines 59-62) is commendably honest: "本质是用有限的选择集换取更高的匹配可靠性". The novelty is moderate — constrained generation applied to LLM expert profiling. The acknowledgment that "LLM 仍可能将边界模糊的 proposal 映射到不同条目" is intellectually honest.

**Cross-domain inspiration (10 / 35)**

No cross-domain inspiration. The proposal remains entirely within the LLM/prompt engineering domain. No borrowings from library science, recommendation systems, database indexing, or organizational design.

**Simplicity of insight (25 / 25)**

The core insight — constrain the search space for domain labeling — is elegant and seems obvious in retrospect. Full marks.

---

### 6. Feasibility (78 / 100)

**Technical feasibility (33 / 40)**

The 4-file coordination complexity is honestly assessed with a dependency chain: "template 的 schema 变更是 inference 和 persistence 的前提，extraction 依赖 template 的 scope 字段" (line 115). However, the assessment still assumes only 4 files are affected without auditing all consumers of `docs/experts/` files. The `scope` field creates a schema contract that any code or prompt reading expert files must handle — the proposal identifies 4 files but the actual surface area may be larger.

**Resource & timeline feasibility (20 / 30)**

Significantly improved. Line 117-119: "预估工作量 2-3 小时" with analogy-based breakdown ("两步流程改造类似...Jaccard 匹配逻辑重构（单文件约 1 小时）"). The estimate is specific enough to evaluate. Concern: 2-3 hours for coordinated changes to 4 prompt files with dependency ordering, plus cross-validation with 2-3 proposals, seems tight. The analogy to past single-file changes may not account for the coordination overhead of multi-file dependent changes.

**Dependency readiness (25 / 30)**

"No external dependencies. All changes within Forge plugin." — verifiable and accurate.

---

### 7. Scope Definition (65 / 80)

**In-scope items are concrete (22 / 30)**

Four in-scope items. Two are concrete:
- "expert-inference.md 嵌入领域分类表，改造为两步生成流程" — concrete
- "expert-template.md 增加 scope 字段" — concrete

Two are less concrete:
- "freeform-expert-persistence.md 更新复用匹配逻辑：区分 scope 级别" — the pre-revised expansion (line 139) adds specificity about *what* changes ("domain-level 专家仅与 domain-level 专家匹配，过滤旧 proposal-specific 专家的噪声"), which helps.
- "extraction-prompt.md 更新 findings 提取逻辑以感知专家 scope 级别" — the parenthetical (line 140) explains the *why* but not the *what*. What specific changes to the extraction prompt template are needed? The current extraction-prompt.md has no concept of scope — what fields, conditions, or output format changes are required?

**Out-of-scope explicitly listed (20 / 25)**

Four items listed as out of scope. Adequate. Notable: the classification table's initial population is implied to be part of the expert-inference.md change but never explicitly stated.

**Scope is bounded (23 / 25)**

Bounded to 4 files within a single directory tree. The boundary note (line 145) explicitly clarifies that the two-step inference is "单个 prompt 内部推理流程的改造，不涉及 pipeline 步骤间的编排变更" — this is a fair clarification.

---

### 8. Risk Assessment (78 / 90)

**Risks identified (26 / 30)**

Four meaningful risks. Good coverage. Missing risks:
- Token cost / latency from two-step inference (identified as a blindspot in iteration 1, still missing)
- Schema versioning risk from the `scope` field
- Classification table staleness as project evolves

**Likelihood + impact rated (25 / 30)**

Ratings are reasonable. The new cross-domain cost risk ("评审成本与领域数成正比") is acknowledged in Scenario 3 and the strategy section rather than in the risk table — it should be in both places.

**Mitigations are actionable (27 / 30)**

Mitigations improved since iteration 1:

- Depth degradation mitigation (line 154): "(1) 专家 prompt 中注入当前 proposal 全文作为评审焦点" — concrete and actionable. "(2) 用户确认循环作为最终质量保障" — still detection rather than prevention, but honestly characterized as a trade-off.
- Old expert noise mitigation (line 157): "匹配逻辑区分 scope 字段" — concrete and actionable.
- Cross-domain cost mitigation: "用户可通过确认循环选择仅使用最相关的一个大领域" — actionable.

---

### 9. Success Criteria (65 / 80)

**Criteria are measurable and testable (20 / 30)**

- SC-1 (coverage breadth): The formula is now explicitly defined — "|K_new ∩ K₁| / |K₁| ≥ 0.5" for each of 2 proposals. This is a significant improvement in measurability. However, the selection of which 2 proposals to test against is arbitrary — for domains with only 1 existing expert, this SC cannot be verified.
- SC-2 (reuse match success): Clear binary criterion. Measurable.
- SC-3 (extraction scope awareness): New — "domain-level 专家产出的 findings 在 freeform review 中被标记为适用范围更广...两者共享的 findings 比例 ≥ 30%". This is measurable but introduces a new concept ("共享 findings") that is not defined in the solution or requirements. What counts as a "shared finding"? Exact match? Semantic overlap?
- SC-4 (confirmation rounds ≤ 2): Measurable, but conflation with user standards persists.
- SC-5 (classification table coverage ≥ 80%): The denominator is now defined ("docs/experts/ 下已有专家文件对应的 proposal 数量，即当前 12 个"), which resolves the iteration-1 ambiguity. However, the denominator counts proposals (12) not expert files (12) — but as noted in D1, 2 experts share 1 proposal, so the actual unique proposal count is 11, not 12. The "即当前 12 个" parenthetical is wrong.

**Coverage is complete (20 / 25)**

Improved since iteration 1. SC-3 now covers extraction-prompt changes. Remaining gap:
- No SC for classification accuracy (how often LLM correctly assigns domain).
- No SC for the quality/effectiveness of domain-level expert reviews vs. proposal-specific.

**SC internal consistency (25 / 25)**

No internal contradictions detected. SC-1 (breadth) and SC-4 (confirmation rounds ≤ 2) create tension but the trade-off is explicitly acknowledged. SC-3's "shared findings ≥ 30%" is a novel metric but not internally contradictory.

---

### 10. Logical Consistency (29 / 90)

**Solution addresses the stated problem (18 / 35)**

The solution addresses the core problem (zero reuse) by broadening expert domain scope through classification-table-guided generation. However:

- **Gap persists from iteration 1**: The root cause analysis is incomplete. Why does the current system generate narrow experts? If `expert-inference.md` instructs the LLM to focus on the specific proposal, the fix might be simpler (adjust the prompt to generate broader keywords) than introducing a full classification table. The proposal doesn't rule out this simpler alternative — the "纯 Prompt 重写" alternative (line 99) argues it fails because "LLM 对'更广'的理解仍因 proposal 而异", but the selected approach still relies on LLM judgment within a domain, which could also vary. The proposal doesn't provide evidence that constrained-domain LLM inference is more consistent than instructed-breadth LLM inference.

- **New contradiction from cross-domain strategy**: The original Constraints section (line 83) states "用户确认循环（Accept / Modify / Regenerate）保持不变". But the cross-domain strategy (lines 42-44, 71) introduces a new interactive choice: "使用全部匹配领域" or "仅使用最相关领域" — this is a *new* decision point within the confirmation flow, which contradicts "保持不变". The existing Accept/Modify/Regenerate options do not include domain scope selection.

**Scope ↔ Solution ↔ Success Criteria aligned (5 / 30)**

Multiple misalignments persist:

1. **Cross-domain strategy scope creep**: The "跨领域多专家策略" section (lines 42-44) describes sequential multi-expert review with findings merging. This is a significant new capability: generating/instantiating multiple experts per evaluation, merging findings from multiple sources. But this capability is not listed in In Scope — In Scope lists 4 file modifications, not a multi-expert coordination feature. The Out of Scope section (line 145) explicitly says "修改 freeform-pipeline 的跨步骤编排逻辑" is out of scope, but multi-expert sequential review with findings merging *is* cross-step orchestration.

2. **extraction-prompt SC alignment**: SC-3 verifies extraction-prompt scope awareness, which now maps to the In Scope item. This is fixed from iteration 1. However, SC-3's "共享 findings 比例 ≥ 30%" introduces a concept not defined anywhere in the solution — what is a "shared finding" between two evaluations?

3. **用户可观测行为 vs Solution**: Behavior 3 (line 38) describes an interactive choice between "全部匹配领域" and "仅最相关领域", but the cross-domain strategy (line 44) says sequential review is the default, with user choice as an override. The Solution section doesn't describe the default — is sequential multi-expert review the default, or is single-best-domain the default with opt-in for multi-domain?

4. **Constraints enumeration vs actual scope**: Line 81 says "改动仅限 `experts/freeform/` 目录下的 prompt 文件和 `rules/freeform-expert-persistence.md`" — this lists 2 locations but the actual scope covers 4 files. While technically correct (all 4 prompt files are in `experts/freeform/`), the enumeration is misleading because it doesn't convey the breadth of the change.

**Requirements ↔ Solution coherent (6 / 25)**

- NFR "可扩展性" (easy to add domains) has no corresponding solution mechanism — who modifies the table, what's the process?
- Scenario 3 / cross-domain strategy describes multi-expert sequential review with merging. The Solution only describes single-domain expert generation. Multi-expert coordination (generation sequencing, findings merging) has no solution design — it's described as a behavior and a requirement but the *how* is unspecified.
- The cross-domain strategy section (line 44) says "findings 合并逻辑在当前 freeform review 中已有先例" — but no evidence is provided for this claim. If findings merging already exists, why is it listed as a new behavior in "用户可观测行为" point 4?

---

## Bias Detection Report

**Annotated regions** (paragraphs with `<!-- pre-revised -->` markers):

Counting annotated paragraphs in the proposal:
- Line 41-44 (跨领域多专家策略): 1 paragraph — annotated `high`
- Line 61-62 (Innovation Highlights): 1 paragraph — annotated `medium`
- Lines 70-71 (Scenario 3): 1 paragraph — annotated `high`
- Lines 105-115 (Technical Feasibility): ~3 paragraphs — annotated `medium`
- Lines 138-140 (In Scope persistence + extraction): 2 items — annotated `medium`
- Lines 153-157 (Risks depth + old-expert-noise): 2 items — annotated `high` and `medium`
- Lines 161-164 (SC-1, SC-3): 2 items — annotated `medium`

Estimated annotated paragraphs: ~8
Attack points targeting annotated content: 7

**Unannotated regions**:
- Unannotated paragraphs: ~30
- Attack points targeting unannotated content: 11

**Ratio (annotated/unannotated): 0.88 / 0.37 = 2.38**

The bias ratio of 2.38 is elevated, indicating higher attack density on pre-revised regions. This is expected — pre-revision addressed the most problematic areas, and the revisions introduced new consistency issues (particularly the cross-domain strategy's tension with the original constraints). The ratio is above the 2.0 threshold, suggesting the scoring should be interpreted with this bias in mind.

**Conflict-with-pre-revision tags**:

1. `conflict-with-pre-revision`: The cross-domain strategy (pre-revised, annotated `high`) introduces multi-expert sequential review. This conflicts with the original (unannotated) Constraints section that says "用户确认循环保持不变" — the cross-domain choice is a new decision point.

2. `conflict-with-pre-revision`: Behavior 3 (pre-revised, line 38) describes interactive domain selection. The cross-domain strategy section (also pre-revised, line 44) describes sequential review as the default. These two pre-revised passages describe the same feature inconsistently.

---

## Phase 3: Blindspot Hunt

### [blindspot-1] Jaccard threshold untested for domain-level keywords

The current Jaccard 0.3 threshold was calibrated for narrow, proposal-specific keyword sets. Domain-level experts will have broader keyword sets, which increases the denominator in Jaccard computation. A broader keyword set matching a narrower proposal could yield a *lower* Jaccard score even when the expert is appropriate, because the intersection doesn't grow proportionally to the domain expert's keyword set. The proposal does not analyze whether the 0.3 threshold needs adjustment. This was identified in iteration-1 blindspot-3 and remains unaddressed.

### [blindspot-2] "Shared findings" in SC-3 is undefined

SC-3 introduces "两者共享的 findings 比例 ≥ 30%" as a success criterion, but "shared findings" is never defined. Is it exact string match of summary fields? Semantic overlap of quote fields? Matching severity? This undefined concept is the lynchpin of the extraction-prompt SC, and its ambiguity makes the SC unverifiable in practice.

### [blindspot-3] Classification table is a hidden coupling point

The classification table embedded in `expert-inference.md` becomes a shared dependency for expert generation, reuse matching, and scope-based filtering. Any change to the table (adding, renaming, or removing a domain) propagates to all domain-level experts already generated. The proposal treats the table as a configuration item but it is actually a schema — changes to it require migration of existing expert files. No migration strategy is described.

### [blindspot-4] "generated_for 指向唯一 proposal" is still wrong

Line 15 claims "每个的 `generated_for` 都指向唯一的 proposal". In reality, `build-orchestration-test-infra.md` and `surface-aware-dispatcher-orchestrator.md` both point to `surface-aware-justfile/proposal.md`. This means the 12 experts map to 11 unique proposals, not 12. This was not caught by the iteration-1 review (which focused on the count, not the uniqueness) and was not fixed in iteration 2 because it was not flagged.

### [blindspot-5] Sequential review cost is underestimated

The cross-domain strategy (line 44) says "顺序串行评审" with cost proportional to matched domain count. For a proposal matching 3 domains, this triples the evaluation time and token cost. The proposal acknowledges this ("复杂 proposal 成本更高") but the risk is not in the Key Risks table — it's buried in a scenario description and a strategy section. A 3x cost increase for the proposals that most need careful evaluation is a significant trade-off that deserves explicit risk analysis.

### [blindspot-6] Scope field migration is unaddressed

The proposal adds a `scope` field to `expert-template.md` and uses it for filtering in persistence matching. But what happens to the existing 12 experts that lack this field? The proposal says "通过 scope 字段过滤" (line 77), implying absent-scope experts are filtered out. But this is an implicit behavior — the filtering logic must explicitly handle the absence case. If a consumer reads an expert file and the `scope` field is missing, what is the default? `proposal-specific`? `null`? The proposal doesn't specify, and this ambiguity could cause bugs.

---

## Attack Points Summary

### Factual / High-Severity

1. **[D1] Evidence inaccuracy persists**: Line 15: "每个的 `generated_for` 都指向唯一的 proposal" — false. Two experts share one proposal. 12 experts map to 11 unique proposals.

2. **[D10] Cross-domain strategy contradicts Constraints**: Lines 42-44 describe sequential multi-expert review with findings merging, which is functionally cross-step orchestration — contradicting Out of Scope item "修改 freeform-pipeline 的跨步骤编排逻辑" (line 145).

3. **[D10] Cross-domain strategy introduces new decision point**: Line 38 describes "使用全部匹配领域" / "仅使用最相关领域" interactive choice, contradicting Constraints line 83: "用户确认循环（Accept / Modify / Regenerate）保持不变".

4. **[D9] SC-3 defines unverifiable metric**: "共享的 findings 比例 ≥ 30%" — "shared findings" is undefined. No algorithm or criteria for what constitutes a shared finding is provided.

5. **[D9] SC-5 denominator is wrong**: "分母为...当前 12 个" — but 2 experts share 1 proposal, so the unique proposal count is 11, not 12.

### Structural / Medium-Severity

6. **[D3] Straw-man "do nothing" alternative**: Line 98: "Rejected: 核心问题不会自行消失" — tautological dismissal. Does not engage with whether the problem's cost justifies the solution's complexity.

7. **[D3] No real industry references**: Lines 89-92 list shallow references without citation, analysis, or depth. Unchanged from iteration 1.

8. **[D2] Technical direction unclear on table format**: The classification table is the architectural centerpiece but its format, embedding strategy, and expected size are never described.

9. **[D4] NFR "可扩展性" remains vague**: Line 76: "应易于扩展" — vague language without quantification. Deduct per rubric rule (-20 pts applied to NFR sub-criterion).

10. **[D4] Missing NFR: Performance**: No analysis of latency or token cost from two-step inference. Identified in iteration-1, still missing.

11. **[D10] Multi-expert coordination has no solution design**: Scenario 3 and cross-domain strategy describe multi-expert review with findings merging. The Solution section only describes single-domain generation. The *how* of multi-expert coordination is unspecified.

12. **[D10] Behavior 3 and cross-domain strategy inconsistent**: Behavior 3 (line 38) implies interactive domain selection as the primary interface. Cross-domain strategy (line 44) implies sequential review is the default with user choice as override. The default behavior is unclear.

13. **[D5] No cross-domain inspiration**: The proposal stays entirely within LLM/prompt engineering. Missing opportunities from classification science, recommendation systems, or organizational design.

14. **[D4] Backward compatibility is nominal**: Line 77 honestly admits "本质是隐性废弃" — old experts are preserved as files but will never be selected by the new system. This is not backward compatibility.

15. **[D7] extraction-prompt.md deliverable is vague**: "更新 findings 提取逻辑以感知专家 scope 级别" (line 140) — the *what* (specific changes to extraction prompt template) is unspecified.

### Low-Severity / Style

16. **[D6] Timeline estimate is tight**: "2-3 小时" for 4 coordinated files with dependency ordering plus cross-validation may underestimate coordination overhead.

17. **[D8] Depth mitigation (2) is detection, not prevention**: Line 154: "用户确认循环作为最终质量保障" — the existing system, not a new mitigation.

18. **[D10] Constraints enumeration misleading**: Line 81 lists 2 locations but actual scope covers 4 files. Technically correct but misleading.

---

## Top 3 Recommended Improvements

1. **Reconcile cross-domain strategy with Constraints and Scope**: Either (a) add multi-expert sequential review to In Scope and remove the "pipeline 编排" exclusion from Out of Scope, or (b) simplify the cross-domain strategy to a single-best-domain approach with optional user override, keeping the pipeline unchanged. The current state has the cross-domain feature half-in and half-out of scope, creating the largest consistency gap in the proposal.

2. **Define "shared findings" in SC-3 or replace with a verifiable metric**: The extraction-prompt SC is currently unverifiable because its key concept is undefined. Either provide an algorithm for computing "shared findings" (e.g., "findings with semantically equivalent summary fields, verified by...") or replace with a more directly measurable criterion.

3. **Fix factual errors and provide computed evidence**: (a) Correct "每个的 generated_for 都指向唯一的 proposal" — 2 experts share 1 proposal, yielding 11 unique. (b) Compute actual pairwise Jaccard scores for the 12 experts' domain keywords to substantiate the "无法达到 0.3" claim. (c) Correct SC-5 denominator from 12 to 11.
