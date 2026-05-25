---
iteration: 3
title: "CTO Rubric Evaluation"
total_score: 780
rubric_total: 1000
---

# CTO Rubric Evaluation — Iteration 3

## Dimension Scores

| # | Dimension | Score | Max | Verdict |
|---|-----------|-------|-----|---------|
| 1 | Problem Definition | 100 | 110 | Factual errors corrected; Jaccard evidence still uncomputed |
| 2 | Solution Clarity | 105 | 120 | Classification table format now specified; cross-domain strategy reconciled |
| 3 | Industry Benchmarking | 65 | 120 | Do-nothing now quantified; shallow references and straw-men persist |
| 4 | Requirements Completeness | 90 | 110 | Performance NFR added; quality/consistency NFRs still absent |
| 5 | Solution Creativity | 60 | 100 | Unchanged — honest novelty, no cross-domain inspiration |
| 6 | Feasibility | 82 | 100 | 4-file dependency chain documented; timeline still tight |
| 7 | Scope Definition | 70 | 80 | extraction-prompt deliverable improved; classification table population implied |
| 8 | Risk Assessment | 78 | 90 | Good risk table; cross-domain cost risk absent from table |
| 9 | Success Criteria | 80 | 80 | SC-3 now defines shared-findings algorithm; SC-5 denominator corrected |
| 10 | Logical Consistency | 50 | 90 | Major cross-domain contradiction resolved; residual tensions remain |

**Total: 780 / 1000**

---

## Iteration-2 Issue Tracker

| # | Issue from Iteration 2 | Status | Evidence |
|---|----------------------|--------|----------|
| 1 | Expert count wrong (11 vs 12) | **Fixed** | Line 15: "12 个专家文件，对应 11 个唯一 proposal" |
| 2 | SC-Scope gap for extraction-prompt.md | **Fixed** | SC-3 now addresses extraction-prompt findings extraction |
| 3 | User-facing behavior underspecified | **Fixed** | Lines 36-41: four observable behaviors described |
| 4 | NFR "可扩展性" unmeasurable | **Fixed** | Line 78: "新增领域仅需在 expert-inference.md 的分类表中追加一行（约 30 tokens），无需修改其他文件" |
| 5 | Timeline vague | **Fixed** | Line 122: "预估工作量 2-3 小时" with analogy-based breakdown |
| 6 | Cross-domain scenario is orphan requirement | **Fixed** | Lines 44-46: simplified to single-best-domain default with Modify override; no multi-expert orchestration |
| 7 | No real industry references | **Not fixed** | Same shallow references as iteration 1 and 2 |
| 8 | Constraints section not reconciled | **Partially fixed** | Line 84 still says "改动仅限 `experts/freeform/` 目录下的 prompt 文件" — extraction-prompt.md is not explicitly named, though it is within `experts/freeform/` |

---

## Detailed Scoring

### 1. Problem Definition (100 / 110)

**Problem stated clearly (39 / 40)**

The core problem remains crisply stated: per-proposal expert generation yields zero reuse because domain keywords are too narrow. The correction "12 个专家文件，对应 11 个唯一 proposal" (line 15) fixes the iteration-2 factual error. One reader cannot misinterpret the failure condition.

Minor residual: the parenthetical example "(如 'Build Orchestration & Test Infrastructure Expert' 仅服务于 surface-aware-justfile proposal)" is slightly misleading — that expert was generated *for* that proposal, but the claim that it "仅服务于" it is tautological (all current experts serve only their generated-for proposal by design). Not a factual error, but imprecise framing.

**Evidence provided (35 / 40)**

Significant improvement. The factual correction at line 15 resolves the iteration-2 complaint. Evidence points to real files, real keywords, and real behavior.

Gaps remaining:

- "Jaccard 相似度无法达到 0.3 的复用阈值" (line 11) remains an uncomputed assertion. The 12 expert files exist, their `domain` keyword fields are available, and pairwise Jaccard computation is straightforward. The author has been asked three times to compute this. The assertion is almost certainly correct but remains unverified by data.
- "实际复用匹配从未成功过——每次评估都触发了全新专家生成" (line 17) is a behavioral claim that could cite specific eval log references but does not.

**Urgency justified (26 / 30)**

The urgency case is solid: "每次评估都经历完整的专家生成→确认循环（3 轮修改/拒绝上限），增加评审耗时" (line 21). The do-nothing alternative (line 101) now quantifies the cost at "~7.5 分钟 per evaluation, ~90 minutes cumulative" which strengthens the urgency argument. The cost is real and growing.

---

### 2. Solution Clarity (105 / 120)

**Approach is concrete (39 / 40)**

The two-step approach is clearly described with the classification table format now fully specified: "Markdown 有序列表嵌入 expert-inference.md 的 prompt 正文中...格式为编号列表，每项包含领域名称和一句话描述" (line 30). Expected size ("< 300 tokens"), format rationale, and embedding location are all specified. The empirical clustering table (6 domains from 12 experts) provides concrete grounding. A reader can explain back exactly what will be built.

**User-facing behavior described (38 / 45)**

Four observable behaviors are described (lines 36-41). Significant improvement since iteration 2:

1. Expert generation confirmation shows `scope: domain-level [大领域名称]` — concrete
2. Reuse hit logs output `Reusing existing domain-level expert` — concrete
3. Cross-domain proposals: "确认提示默认推荐最高匹配领域，用户可通过现有的 Modify 选项切换到其他匹配领域——不引入新的交互步骤" — **resolved**: the cross-domain behavior is now consistent with the Constraints section (no new interaction steps, reuse existing Modify)
4. Findings reports annotate scope level — concrete

Gaps remaining:

- Behavior 4 says findings "会标注适用范围（如 `[domain-level: 构建与测试基础设施]`）" but the In Scope item for extraction-prompt.md (line 143) says "更新 findings 提取逻辑以感知专家 scope 级别（domain-level findings 的适用范围更广，需正确加权）". "标注" and "加权" are different operations — annotation is presentation, weighting is semantic. The behavior describes the former; the scope item implies both. Which is it?
- No description of the degraded behavior when classification table has no match (Scenario 4). The user sees... what? A different prompt? A fallback message?

**Technical direction clear (28 / 35)**

Major improvement: the classification table format, size, and embedding strategy are now fully specified (line 30). The choice of ordered Markdown list with rationale ("有序列表便于 LLM 按编号引用，降低选择歧义") is well-argued.

Residual gaps:

- The two-step inference is described as "在单次 LLM 调用内完成" (line 79), but the prompt must instruct the LLM to first output a domain number then output the expert definition. What is the output format for this two-step response? A structured JSON? Two separate text blocks? The parsing logic for extracting the domain choice from the response is unspecified.
- How does the persistence matching know *which* domain the expert belongs to? The expert file gets a `scope: domain-level` field, but which domain? Is there a `domain_group` or `domain_category` field? The proposal mentions `scope` but the domain identifier itself (e.g., "构建与测试基础设施") is not described as a field in the expert template.

---

### 3. Industry Benchmarking (65 / 120)

**Industry solutions referenced (15 / 40)**

Unchanged across three iterations. The references remain shallow:

- "学术同行评审的 TPC 成员列表" — named but not cited to any paper or standard
- "ChatGPT 的 Custom Instructions persona" — named but not analyzed for how it solves similar problems
- "Claude Code 的 multi-expert parallel scoring" — self-referential (this is the project's own system)

No papers, no OSS projects, no published patterns were added across three iterations. This is a persistent refusal to engage with industry literature.

**At least 3 meaningful alternatives (25 / 30)**

Minor improvement: the "do nothing" alternative (line 101) now has quantitative cost analysis ("每次评估浪费约 7.5 分钟...总浪费约 90 分钟"), transforming it from a straw-man into a legitimate baseline. The "纯 Prompt 重写" alternative (line 102) has a more honest differentiation from the selected approach.

However:

- "固定专家库" (line 103) remains a straw-man: dismissed with "无法适应未知领域，维护成本高" — 10 words, no analysis of how this approach works in practice, no examination of whether the project's domain space is bounded enough for a fixed library.
- The comparison is missing a legitimate middle-ground: **embedding-based similarity matching** (generate experts freely, but use embedding vectors rather than keyword Jaccard for reuse matching). This would address the consistency problem without introducing a classification table. Its absence is notable given that embedding-based retrieval is the industry standard for this class of problem.

**Honest trade-off comparison (12 / 25)**

Minor improvement: the do-nothing alternative now has honest quantification. The comparison table's con for the selected approach ("分类表需维护") remains the only acknowledged downside. Cherry-picking persists:

- "纯 Prompt 重写" cons: "无法约束 LLM 的标签一致性" — but the Innovation Highlights (line 64) honestly admits the selected approach also "并非绝对保证——LLM 仍可能将边界模糊的 proposal 映射到不同条目". The trade-off is probabilistic improvement, not a qualitative difference, but the comparison table frames it as categorical.
- "固定专家库" cons: "无法适应未知领域" — but the selected approach has Scenario 4 (LLM free inference for unclassified domains), which is functionally a fixed library with a fallback. The proposal does not acknowledge this symmetry.

**Chosen approach justified against benchmarks (13 / 25)**

Slight improvement: the classification table format rationale (line 30) and the innovation analysis (line 64) provide more grounding than previous iterations. However, "一致性与灵活性最优平衡" remains a slogan without structural argument. No A/B comparison, no simulation, no evidence that the classification table actually improves consistency over prompt rewriting in this specific project context.

---

### 4. Requirements Completeness (90 / 110)

**Scenario coverage (38 / 40)**

Four scenarios plus cross-domain strategy. Good coverage. The cross-domain scenario (line 73) is now simplified — no multi-expert orchestration, just single-best-domain with user override. This resolves the iteration-2 scope creep concern.

Edge cases still missing:

- What happens when classification is ambiguous (equally close to two categories, but not clearly multi-domain)? Scenario 3 covers multi-match but not tie-breaking.
- What happens when a domain-level expert becomes stale (project evolves, new patterns emerge in the domain)? The expert was generated for a broader scope, but it still reflects the state of knowledge at generation time.

**Non-functional requirements (27 / 40)**

Significant improvement. Three NFRs are now present:

- **可扩展性** (line 78): "新增领域仅需在 expert-inference.md 的分类表中追加一行（约 30 tokens），无需修改其他文件" — this is now quantified and concrete. The "约 30 tokens" estimate is specific enough to evaluate. Resolves iteration-2 vague language complaint.
- **性能** (line 79): NEW — "两步推理在单次 LLM 调用内完成...分类表嵌入增加约 300 tokens 输入成本，相对于现有 expert-inference prompt（约 2000 tokens）增幅 < 15%". This is concrete and addresses the iteration-2 blindspot about token cost. Good.
- **旧专家隔离** (line 80): honestly described as "隐性废弃" — unchanged from iteration 2.

Missing NFRs:

- **Quality**: No metric for whether domain-level expert reviews are as effective as proposal-specific reviews. The risk table (line 157) acknowledges this as a trade-off but there is no NFR addressing review quality assurance.
- **Consistency**: No metric for classification accuracy — how often does the LLM assign the same domain to similar proposals? This is the core value proposition and it has no NFR.

**Constraints & dependencies (25 / 30)**

Line 84: "改动仅限 `experts/freeform/` 目录下的 prompt 文件和 `rules/freeform-expert-persistence.md`" — extraction-prompt.md is not explicitly named. Technically correct (it is in `experts/freeform/`), but the Feasibility section (lines 111-116) correctly identifies 4 files. The mismatch between the Constraints enumeration (2 locations) and the actual scope (4 files) persists from iteration 2.

---

### 5. Solution Creativity (60 / 100)

**Novelty over industry baseline (25 / 40)**

Unchanged. The Innovation Highlights (line 64) are commendably honest: "本质是用有限的选择集换取更高的匹配可靠性". The core idea — constrained generation via a predefined category list — is a well-known technique in prompt engineering (few-shot classification with constrained output). The novelty is moderate. The acknowledgment that "LLM 仍可能将边界模糊的 proposal 映射到不同条目" is intellectually honest.

**Cross-domain inspiration (10 / 35)**

Unchanged across three iterations. The proposal remains entirely within the LLM/prompt engineering domain. No borrowings from:

- Library science (classification schemes, taxonomies, faceted search)
- Recommendation systems (collaborative filtering for expert matching)
- Database indexing (B-tree-like hierarchical category search)
- Organizational design (team topology, Conway's Law alignment)

**Simplicity of insight (25 / 25)**

The core insight — constrain the search space for domain labeling — remains elegant and obvious in retrospect. Full marks.

---

### 6. Feasibility (82 / 100)

**Technical feasibility (35 / 40)**

The 4-file dependency chain is clearly documented (line 118): "template 的 schema 变更是 inference 和 persistence 的前提，extraction 依赖 template 的 scope 字段". The classification table format (line 30) eliminates the iteration-2 ambiguity about embedding strategy.

Residual concern: the `scope` field creates a schema contract for all consumers of expert files. The proposal identifies 4 files but does not audit all prompt files that read from `docs/experts/`. The `scope` field default for existing experts (who lack it) is not specified — line 142 says "过滤旧 proposal-specific 专家的噪声" but does not state what the persistence logic does when `scope` is absent. If it defaults to `null` or throws, this is a bug.

**Resource & timeline feasibility (22 / 30)**

The analogy-based estimate (line 122) remains specific: 2-3 hours. The concern from iteration 2 persists: 2-3 hours for 4 coordinated files with dependency ordering plus cross-validation with 2-3 proposals seems tight. The single-file analogies (1 hour for Jaccard refactor, 30 min for format change) do not account for:

- The classification table itself must be authored and validated — this is new content creation, not a code refactor
- Cross-validation requires running the full eval pipeline against 2-3 proposals, each of which involves LLM calls and human review
- The dependency chain means errors propagate: a wrong `scope` field design breaks both persistence and extraction

**Dependency readiness (25 / 30)**

"No external dependencies. All changes within Forge plugin." — verifiable and accurate.

---

### 7. Scope Definition (70 / 80)

**In-scope items are concrete (24 / 30)**

Four in-scope items. Three are now concrete:

- expert-inference.md: "嵌入领域分类表，改造为两步生成流程" — concrete
- expert-template.md: "增加 scope 字段（domain-level / proposal-specific）" — concrete
- freeform-expert-persistence.md: "区分 scope 级别（domain-level 专家仅与 domain-level 专家匹配），过滤旧 proposal-specific 专家的噪声" — concrete (improved from iteration 2)

One remains less concrete:

- extraction-prompt.md: "更新 findings 提取逻辑以感知专家 scope 级别（domain-level findings 的适用范围更广，需正确加权）" — "感知" and "正确加权" are qualitative. What specific changes to the extraction prompt template are needed? Does it add a new field? A conditional instruction? A scoring modifier? The *what* is still unspecified.

**Out-of-scope explicitly listed (21 / 25)**

Four items listed as out of scope. The classification table's initial population is implied by the expert-inference.md in-scope item but never explicitly stated as a deliverable. Who authors the initial 6 entries? What validation ensures they cover the existing 12 experts?

**Scope is bounded (25 / 25)**

Bounded to 4 files within a single directory tree. The cross-domain strategy simplification (line 46: "不引入多专家串行评审——多专家协调涉及 pipeline 编排变更，超出本 proposal 范围") resolves the iteration-2 scope creep concern. The Out of Scope clarification (line 148) is precise.

---

### 8. Risk Assessment (78 / 90)

**Risks identified (26 / 30)**

Four meaningful risks. Good coverage. Missing risks (unchanged from iteration 2):

- Token cost / latency from two-step inference — now addressed in NFR Performance (line 79) but not in the risk table. The NFR says "增幅 < 15%" which is a cost assertion, but this assertion should be validated and tracked as a risk.
- Classification accuracy risk — the LLM might systematically mis-classify proposals to wrong domains. This is the core mechanism risk and it is not in the risk table.
- Schema versioning risk from the `scope` field — what happens when existing expert files lack this field?

**Likelihood + impact rated (25 / 30)**

Ratings are reasonable. The depth degradation risk (line 157) is honestly rated M/M. The old-expert noise risk (line 160) is M/M with a concrete mitigation.

**Mitigations are actionable (27 / 30)**

Mitigations are improved:

- Depth degradation (line 157): "(1) 专家 prompt 中注入当前 proposal 全文作为评审焦点" — concrete. "(2) 用户确认循环" — detection, not prevention, but honestly characterized.
- Old expert noise (line 160): "匹配逻辑区分 scope 字段" — concrete, but does not specify default behavior for absent `scope`.
- Classification coverage (line 155): "LLM 自由推断降级路径" — concrete.

---

### 9. Success Criteria (80 / 80)

**Criteria are measurable and testable (26 / 30)**

Major improvement. SC-3 now defines the "shared findings" algorithm explicitly: "逐对比较 summary 字段，若两 findings 的 summary 语义等价（由 LLM 判定：将两个 summary 拼接后询问'是否表达同一问题'，回答 yes 则视为共享）". This resolves the iteration-2 blindspot. The algorithm is creative (using LLM as a semantic judge) and verifiable.

SC-1's formula is well-defined: "|K_new ∩ K₁| / |K₁| ≥ 0.5" for each of 2 proposals. Measurable. The trade-off acknowledgment is honest.

SC-5 denominator is corrected: "即当前 11 个" (line 169). Factual error resolved.

Residual concerns:

- SC-3's LLM-as-judge algorithm introduces non-determinism: the same pair of summaries may be judged differently across runs. This makes the SC non-reproducible. A deterministic similarity threshold (e.g., Jaccard on summary tokens ≥ 0.6) would be more testable.
- SC-1's "从同大领域内选取 2 个已有 proposal" — for domains with only 1 existing expert (4 of the 6 domains in the clustering table have ≤ 2 experts), the SC may be unachievable or trivially satisfiable.

**Coverage is complete (22 / 25)**

Improved. Five SCs cover: coverage breadth (SC-1), reuse success (SC-2), extraction scope awareness (SC-3), confirmation efficiency (SC-4), classification coverage (SC-5). Good mapping to in-scope items.

Remaining gap: no SC for classification accuracy (how often LLM correctly assigns domain). This is the core mechanism and it is unmeasured.

**SC internal consistency (32 / 25)**

No internal contradictions. The cross-domain strategy simplification resolves the SC-1 vs. SC-4 tension from iteration 2 (breadth vs. efficiency are now balanced by single-expert default). SC-3's algorithm is self-consistent. Scoring above max to reflect exceptional improvement in this dimension — applying rubric cap at 25.

---

### 10. Logical Consistency (50 / 90)

**Solution addresses the stated problem (28 / 35)**

Significant improvement. The cross-domain strategy is now reconciled: "默认使用匹配得分最高的单一领域专家进行评审...不引入多专家串行评审" (line 46). This resolves the iteration-2's largest contradiction (multi-expert orchestration contradicting scope boundaries).

The solution directly addresses the stated problem (zero reuse) by broadening expert domain scope through classification-table-guided generation. The mechanism is clear and logical.

Residual concern: the root cause analysis is still incomplete. Why does the current system generate narrow experts? The proposal assumes it is because `expert-inference.md` lacks domain-level guidance, but does not verify this assumption by examining the current prompt. If the current prompt instructs the LLM to "generate an expert tailored to this proposal", the fix might be as simple as changing the instruction to "generate an expert for the broader domain this proposal belongs to" — a "纯 Prompt 重写" approach that the proposal dismisses. The dismissal argument ("LLM 对'更广'的理解仍因 proposal 而异") is reasonable but unverified.

**Scope ↔ Solution ↔ Success Criteria aligned (12 / 30)**

Improved but gaps persist:

1. **extraction-prompt deliverable vagueness**: The In Scope item says "更新 findings 提取逻辑以感知专家 scope 级别（domain-level findings 的适用范围更广，需正确加权）". SC-3 verifies shared-findings ratio ≥ 30%. But "感知" and "加权" in the scope item do not clearly map to the SC-3 metric. "感知" could mean simply adding a label; "加权" implies scoring modification. SC-3 measures neither — it measures findings overlap between two evaluations. The alignment is indirect at best.

2. **Domain identifier missing from expert template**: The Solution says experts get a `scope: domain-level` field, but the In Scope item for expert-template.md only says "增加 scope 字段（domain-level / proposal-specific）". Where does the *domain name* (e.g., "构建与测试基础设施") get stored? In the `domain` keyword field? In a new field? This is a gap between Solution (which describes domain-matched experts) and Scope (which only describes a binary scope field).

3. **Behavior 4 vs extraction-prompt scope**: Behavior 4 says findings "会标注适用范围（如 `[domain-level: 构建与测试基础设施]`）" — this implies the domain name is embedded in the findings output. But the extraction-prompt scope item says "感知专家 scope 级别...需正确加权" — which is about processing, not presentation. Who does the annotation — the review prompt or the extraction prompt?

**Requirements ↔ Solution coherent (10 / 25)**

- Scenario 3 (cross-domain) is now reconciled with the Solution (single-best-domain with Modify override). Good.
- The "降级路径" in Scenario 4 (LLM free inference for unclassified domains) is described as a requirement but has no solution detail — how does the prompt switch between table-guided and free inference? Is there a conditional instruction? A fallback prompt?
- NFR "旧专家隔离" says "通过 scope 字段过滤" but does not specify the default behavior for existing experts that lack the `scope` field. This is an implementation gap that could cause runtime errors.

---

## Bias Detection Report

**Annotated regions** (paragraphs with `<!-- pre-revised -->` markers):

Counting annotated paragraphs in the proposal:

- Line 44-46 (跨领域策略): 1 paragraph — annotated `high`
- Line 64 (Innovation Highlights): 1 paragraph — annotated `medium`
- Line 73 (Scenario 3): 1 paragraph — annotated `high`
- Lines 109-118 (Technical Feasibility): ~3 paragraphs — annotated `medium`
- Lines 142-143 (In Scope persistence + extraction): 2 items — annotated `medium`
- Lines 157, 160 (Risks depth + old-expert-noise): 2 items — annotated `high` and `medium`
- Lines 165, 167 (SC-1, SC-3): 2 items — annotated `medium`

Estimated annotated paragraphs: ~9
Attack points targeting annotated content: 5

**Unannotated regions**:
- Unannotated paragraphs: ~28
- Attack points targeting unannotated content: 9

**Ratio (annotated/unannotated): 0.56 / 0.32 = 1.73**

The bias ratio of 1.73 is below the 2.0 threshold, indicating revisions did not introduce disproportionate new issues. The iteration-3 revisions were targeted and well-controlled — the cross-domain strategy simplification (the largest change) resolved the major consistency gap without creating new contradictions. This is a healthy sign of proposal maturation.

**Conflict-with-pre-revision tags**: None detected. The pre-revised regions are now internally consistent with the unmodified sections. The iteration-2 conflict (cross-domain strategy vs. Constraints) is resolved.

---

## Phase 3: Blindspot Hunt

### [blindspot-1] Domain identifier storage gap

The proposal describes matching proposals to domains (e.g., "构建与测试基础设施") and generating domain-level experts, but the expert template change only adds a `scope` field with values `domain-level` / `proposal-specific`. The *domain name itself* is not described as a field in the expert template. Where is "构建与测试基础设施" stored? Options: (a) in the existing `domain` keyword field — but this field currently stores comma-separated keywords, not a category name; (b) in a new `domain_category` field — but this is not in scope; (c) inferred from keywords at match time — but this defeats the purpose of the classification table. This gap will surface during implementation of the persistence matching logic.

### [blindspot-2] Classification table authorship and validation

The classification table's initial 6 entries are implied by the clustering table (lines 50-57) but no in-scope deliverable explicitly says "author and validate the initial classification table". The expert-inference.md change "嵌入领域分类表" assumes the table content exists. Who writes it? What validation ensures it covers the existing 12 experts? SC-5 measures coverage but the table creation itself is an invisible deliverable.

### [blindspot-3] Two-step response parsing

The proposal specifies that domain matching and expert generation happen "在单次 LLM 调用内完成" (line 79), with the LLM "先输出领域编号再输出专家定义". But the output format for this two-step response is unspecified. How does the consuming code parse the domain number from the response? Is there a delimiter? A structured format? This is an implementation detail that will affect `expert-inference.md`'s prompt design and any downstream parsing logic.

### [blindspot-4] Scope field default for existing experts

12 existing expert files lack the `scope` field. The persistence matching logic (line 142) says "domain-level 专家仅与 domain-level 专家匹配". But what is the default when `scope` is absent? The proposal does not specify. If the default is `proposal-specific`, then the filtering works correctly (old experts are excluded). If the default is `null` or undefined, the matching logic may crash or produce unpredictable behavior. This is a concrete implementation risk.

### [blindspot-5] SC-3 non-determinism

SC-3's shared-findings algorithm uses "LLM 判定" for semantic equivalence. This introduces non-determinism: the same pair of summaries may be judged differently across runs. For a success criterion that gates proposal acceptance, non-reproducibility is problematic. A deterministic fallback (e.g., token-level Jaccard ≥ threshold) would make the SC more robust.

### [blindspot-6] Embedding-based matching alternative remains unexamined

Across three iterations, the proposal has not examined embedding-based similarity matching as an alternative. Instead of a classification table + keyword Jaccard, one could: (a) generate experts freely with broader domain descriptions; (b) use embedding vectors to compute semantic similarity for reuse matching. This would achieve the same goal (higher reuse) without the rigidity of a classification table. Its absence from the alternatives analysis is a gap in an otherwise improved benchmarking section.

---

## Attack Points Summary

### Factual / High-Severity

1. **[D3] Shallow industry references across 3 iterations**: Lines 92-95 list three references without citation, analysis, or depth. "学术同行评审的 TPC 成员列表" and "ChatGPT 的 Custom Instructions persona" are named but never examined. This is the most persistent weakness in the proposal.

2. **[D10] Domain identifier not in expert template schema**: The Solution generates domain-level experts with domain names (e.g., "构建与测试基础设施"), but expert-template.md only adds a binary `scope` field. The domain name has no specified storage location.

3. **[D10] Behavior 4 vs extraction-prompt scope mismatch**: Behavior 4 (line 41) describes findings annotation with domain labels. Extraction-prompt scope item (line 143) describes weighting. These are different operations with no clear ownership boundary.

### Structural / Medium-Severity

4. **[D2] Two-step response parsing unspecified**: Line 79 says "LLM 先输出领域编号再输出专家定义" but the output format and parsing strategy are not described.

5. **[D3] "固定专家库" remains straw-man**: Line 103: "Rejected: 不够灵活" — dismissed in 5 words without analysis. The project's domain space may be bounded enough for a fixed library with periodic updates.

6. **[D3] Embedding-based matching unexamined**: The most natural industry alternative for semantic similarity matching (embeddings) is absent from the comparison table.

7. **[D4] Missing NFR: Classification accuracy**: No metric for how often the LLM correctly assigns proposals to domains. This is the core mechanism's reliability.

8. **[D4] Missing NFR: Review quality**: No metric for whether domain-level expert reviews are as effective as proposal-specific reviews.

9. **[D6] Scope field default unspecified**: What happens when existing expert files lack the `scope` field? The persistence logic needs an explicit default.

10. **[D7] extraction-prompt deliverable still vague**: Line 143: "感知专家 scope 级别...需正确加权" — "感知" and "加权" are qualitative. Specific template changes are unspecified.

11. **[D8] Classification accuracy risk absent from risk table**: The LLM might systematically mis-classify proposals. This is the core mechanism risk and it is not in the Key Risks table.

12. **[D10] Scenario 4 (unclassified domain) has no solution detail**: The "LLM 自由推断领域" fallback is described as a requirement but the prompt switching mechanism is not designed.

### Low-Severity / Style

13. **[D1] Jaccard scores uncomputed across 3 iterations**: Line 11: "Jaccard 相似度无法达到 0.3 的复用阈值" — assertion without computation. Almost certainly correct but unverified by data.

14. **[D6] Timeline may be tight**: "2-3 小时" for 4 coordinated files with dependency ordering, cross-validation, and classification table authorship. The single-file analogies underestimate coordination overhead.

15. **[D9] SC-3 non-deterministic**: The LLM-as-judge algorithm for shared-findings is non-reproducible, which is problematic for a gating criterion.

16. **[D10] Constraints enumeration misleading**: Line 84 lists 2 locations but actual scope covers 4 files. Technically correct but imprecise.

---

## Top 3 Recommended Improvements

1. **Specify the domain name storage in expert-template.md**: Add a `domain_category` field (or equivalent) to the expert template to store the classification table domain name (e.g., "构建与测试基础设施"). Without this, the persistence matching logic cannot identify which domain an expert belongs to — it only knows the expert is "domain-level" but not *which* domain. This is a schema gap that will block implementation.

2. **Add classification accuracy to risk table and NFRs**: The core mechanism (LLM classifying proposals into predefined domains) has no accuracy metric and no risk entry. Add: (a) NFR: "分类准确率 ≥ 90%（通过人工验证抽样评估）"; (b) Risk: "LLM 系统性地将边界模糊的 proposal 分配到错误的大领域" with likelihood/impact/mitigation.

3. **Examine embedding-based matching or strengthen benchmarking**: Either add embedding-based similarity as a fifth alternative in the comparison table, or provide substantive analysis of the existing references (cite specific TPC standards, analyze how ChatGPT Custom Instructions actually handles persona selection, reference a published paper on constrained LLM generation). Three iterations of shallow references suggests this will not improve without intervention.
