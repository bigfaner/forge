---
created: "2026-05-25"
evaluator: CTO-review
iteration: baseline
score: 650
target: 900
---

# Proposal Evaluation: Baseline Iteration

## Phase 1: Reasoning Audit

### Argument Chain Trace

**Problem -> Solution**: The problem (experts generated per-proposal, zero reuse) is real. The proposed solution (domain classification table) does address it by grouping proposals under broader domain labels. However, the argument assumes the 12 existing experts naturally cluster into a small number of domains. Inspection of the actual expert `domain` fields reveals heterogeneous, highly specific keywords (e.g., "golden-dataset, snapshot-testing, go-testing, schema-regression, type-dispatch" vs "prompt template engineering & agent protocol design"). These do NOT obviously cluster into 8-12 broad domains. The solution assumes LLM will bridge the gap, but the mechanism for generating appropriately broad keywords is underspecified.

**Solution -> Evidence**: The proposal cites "Jaccard similarity never reaches 0.3 threshold" as the core evidence. However, this claim is **factually contradicted** by the existing codebase: `docs/experts/config-schema-surface-detection.md` has a `review_history` entry showing it was successfully reused for `docs/proposals/test-recipe-unification/proposal.md` with `rubric_delta: 126` and `substantive_change: true`. This means reuse matching has succeeded at least once. The "0 successful matches" claim is false. Additionally, two experts (`build-orchestration-test-infra.md` and `surface-aware-dispatcher-orchestrator.md`) were generated for the same proposal (`surface-aware-justfile`), which further undermines the "one expert per proposal" characterization.

**Evidence -> Success Criteria**: SC-2 ("reuse match succeeds for second proposal in same domain") is already partially achieved by the `config-schema-surface-detection.md` reuse case. The SC does not acknowledge this baseline, making it uncalibrated. SC-1 ("domain keywords cover >= 2 proposals' domain intersection") is vague on how "domain intersection" is computed and tested.

**Self-contradiction check**: The proposal says "改动仅限 `experts/freeform/` 目录下的 prompt 文件" but In Scope lists `freeform-expert-persistence.md` which is in `rules/` not `experts/freeform/`. Minor path error but indicates imprecise scope analysis.

### SC Consistency Deep-Dive

**Cluster A: expert-inference.md changes**
- In Scope: "嵌入领域分类表，改造为两步生成流程"
- SC-1: Requires domain keywords to cover >= 2 proposals (depends on classification table quality)
- SC-2: Reuse match succeeds (depends on persistence logic update)
- SC-4: Classification table covers >= 80% of existing proposals

**Bidirectional check SC-1 <-> SC-4**: SC-4 measures classification table coverage of proposals. SC-1 measures keyword breadth per expert. These are related but not derived from each other. A classification table covering 80% of proposals does not guarantee each generated expert's keywords span 2+ proposals. **Gap: no mechanism ensures SC-1 follows from SC-4.**

**Bidirectional check SC-2 <-> In Scope item 3**: In Scope says "更新复用匹配逻辑以适配 domain-level 专家". SC-2 says "复用匹配成功". The In Scope item is necessary but not sufficient for SC-2 -- the matching logic update must also produce correct results. **Acceptable: SC-2 is the verification of the In Scope item.**

**SC-3 vs solution**: SC-3 ("确认轮次 <= 2") assumes domain-level experts are closer to user expectations. But the proposal provides no mechanism to reduce confirmation rounds -- the Accept/Modify/Regenerate loop remains unchanged. This SC measures a hypothesis, not a verifiable outcome of the stated scope.

## Phase 2: Rubric Scoring

### D1. Problem Definition: 70/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 30/40 | The core problem (low expert reuse) is identifiable but the "11 个专家" count is wrong (there are 12), and the claim "实际复用匹配从未成功过" is factually false (`config-schema-surface-detection.md` was successfully reused). Two factual errors in the problem statement undermine clarity. |
| Evidence provided | 15/40 | The domain keyword examples are cherry-picked -- the specific 5-keyword example from `go-pipeline-integration-type-system-engineer.md` is real, but the broader pattern is overstated. Some experts have reasonably broad domains (e.g., "Go refactoring, technical debt, code quality, CLI architecture"). The claim that "任何" keyword has near-zero probability of appearing in other proposals is unverified and likely false for keywords like "go-backend" or "prompt-template-architecture" which could match multiple proposals. |
| Urgency justified | 25/30 | The urgency argument (each evaluation triggers full generation cycle, 3 rounds max) is concrete. However, it overstates the cost -- "增加评审耗时" without quantifying how much time. If each generation cycle takes 30 seconds, the total cost is 90 seconds per evaluation, which may not justify a system redesign. |

### D2. Solution Clarity: 75/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 28/40 | The two-step flow (domain matching -> expert generation) is described, but the classification table itself is not shown. What are the 8-12 domains? Without this, a reader cannot explain back what will be built with specificity. The proposal says "预定义领域分类表" but never provides or samples it. |
| User-facing behavior described | 32/45 | The 4 scenarios describe user-facing outcomes reasonably well. Scenario 3 (cross-domain proposal) is notably thin -- "匹配最相关的一个大领域" raises the question: what if the best match is wrong? Scenario 4 (fallback to LLM free inference) exists but the user experience of this fallback path is not described. |
| Technical direction clear | 15/35 | "改动仅涉及 prompt 文件（Markdown），不涉及代码变更" is clear on surface, but the actual technical mechanism is vague. How does the classification table get embedded? As a section in `expert-inference.md`? As a separate reference file? How does the persistence logic change to handle `scope` field? The proposal lists files to change but not how they change. |

### D3. Industry Benchmarking: 55/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 20/40 | Three categories are listed (fixed expert pool, LLM free inference, hybrid) but no specific product or open-source project is named by name except "ChatGPT Custom Instructions" and "Claude Code multi-expert parallel scoring" which are vague references, not formal benchmarks. No academic papers, no open-source repos. |
| At least 3 meaningful alternatives | 20/30 | Four alternatives are presented including "do nothing". However, "纯 Prompt 重写" is a straw man -- the proposal presents it only to dismiss it with "没解决一致性问题" without evidence that pure prompt rewriting cannot improve consistency. "固定专家库" is also a straw man with "不够灵活" dismissal. |
| Honest trade-off comparison | 10/25 | The comparison table lists Pros/Cons but they are generic. "分类表需维护" is the only con for the selected approach -- what about the risk of misclassification? What about the cognitive overhead of maintaining the classification table as the project evolves? |
| Chosen approach justified against benchmarks | 5/25 | The justification is "一致性与灵活性最优平衡" which is a conclusion, not a justification. No analysis of why this specific balance point is optimal for this project's constraints. |

### D4. Requirements Completeness: 55/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 22/40 | The 4 scenarios cover basic flows but miss critical edge cases: (1) What happens when two experts from the same domain exist but with different focus areas? (2) What if a proposal matches multiple domains equally? (3) What happens when the classification table is updated -- do existing domain-level experts become stale? |
| Non-functional requirements | 18/40 | Only two NFRs listed: extensibility and backward compatibility. Missing: performance (does classification table lookup add latency?), reliability (what if LLM misclassifies?), observability (how do we know the system is working?). The "可扩展性" NFR is vague -- "只需修改 prompt 文件" is a claim, not a requirement. |
| Constraints & dependencies | 15/30 | Constraints are listed but the "改动仅限 `experts/freeform/`" constraint is inaccurate -- `freeform-expert-persistence.md` is in `rules/`, not `experts/freeform/`. The dependency on "不影响 freeform-review-protocol" is stated as a constraint but not verified. |

### D5. Solution Creativity: 55/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 20/40 | The "classification table + LLM within domain" approach is a standard hybrid pattern. The proposal acknowledges this: "在两者之间取得平衡". There is minimal innovation beyond the standard hybrid approach. |
| Cross-domain inspiration | 15/35 | No evidence of cross-domain inspiration. The proposal stays entirely within the "expert system" domain. No references to taxonomy systems, ontology matching, clustering algorithms, or any other domain's approach to similar problems. |
| Simplicity of insight | 20/25 | The core insight ("use a fixed taxonomy to ensure consistent domain labeling, then let LLM be flexible within domains") is clean and understandable. It is a reasonable and practical solution. |

### D6. Feasibility: 70/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 35/40 | The change is indeed limited to Markdown prompt files, which makes it technically straightforward. No code changes required is a fair claim. |
| Resource & timeline | 20/30 | "单次 prompt 改写 + 测试验证，工作量小" is vague. No timeline estimate. What does "测试验证" entail? Running the full eval pipeline on how many proposals? |
| Dependency readiness | 15/30 | "无外部依赖" is correct, but the proposal does not address the dependency on the existing Jaccard/weighted scoring system in `freeform-expert-persistence.md`. The scoring threshold (0.3 Jaccard or weighted score >= 5) was designed for proposal-specific keywords. With broader domain-level keywords, the threshold may need recalibration -- this dependency is not acknowledged. |

### D7. Scope Definition: 50/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 18/30 | Three items listed, each referencing a specific file. However, "嵌入领域分类表" and "改造为两步生成流程" describe the approach, not the deliverable. A deliverable would be "expert-inference.md updated with classification table section X and two-step protocol Y". |
| Out-of-scope explicitly listed | 17/25 | Good: 5 out-of-scope items are listed, including the critical "复用匹配对旧专家的兼容". However, missing from out-of-scope: what happens to the Jaccard threshold? Is threshold recalibration in or out of scope? |
| Scope is bounded | 15/25 | The scope is bounded to 3 files, but the "domain classification table" itself is an unspecified deliverable. How many domains? How are they defined? This is the core artifact but its scope is undefined. |

### D8. Risk Assessment: 55/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 20/30 | Three risks listed. Missing risks: (1) Classification table becomes a bottleneck for new domains (the "cold start" problem). (2) Domain-level experts may be too generic for specialized proposals. (3) Backward compatibility of matching scores with existing threshold. |
| Likelihood + impact rated | 18/30 | Ratings are reasonable but the "领域级专家的评审深度不如 proposal-specific 专家" risk is rated M/M -- this is the core quality trade-off and arguably should be H impact. If the domain-level expert is too shallow, the entire system fails its purpose. |
| Mitigations are actionable | 17/30 | Mitigation for Risk 1 ("提供 LLM 自由推断降级路径") is actionable. Mitigation for Risk 2 ("LLM 在领域内细化专业方向时参考 proposal 内容") is vague -- how? Mitigation for Risk 3 ("分类表控制在大领域粒度（8-12 个）") is a design constraint, not a mitigation. |

### D9. Success Criteria: 40/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 15/30 | SC-1: "domain 关键词覆盖范围 >= 2 个 proposal 的领域交集" -- "领域交集" is undefined. How do you compute the intersection of two proposals' domains? SC-3: "确认轮次 <= 2" is testable. SC-4: "覆盖 >= 80% 的已有 proposal" is testable. SC-2: "复用匹配成功" is binary and testable. But SC-1 is ambiguous. |
| Coverage is complete | 10/25 | SC covers the core goal (reuse) but does not cover: (1) classification table quality (are the domains correct?), (2) expert review quality post-change (does domain-level expertise still add value?), (3) backward compatibility (do existing 12 experts still work with unchanged matching logic?). |
| SC internal consistency | 15/25 | SC-1 and SC-4 are loosely coupled but not contradictory. SC-3 (confirmation rounds <= 2) is independent and unverified by any in-scope mechanism. **Ambiguity flag**: SC-1 uses "领域交集" which is undefined. If it means keyword overlap between proposals mapped to the same domain, then it depends on the classification table which is not defined. **Ambiguous -- requires author clarification.** |

### D10. Logical Consistency: 65/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 28/35 | The classification table approach directly addresses the reuse problem by ensuring consistent domain labeling. The causal chain is sound: consistent labels -> higher Jaccard scores -> more reuse matches. However, the problem statement contains factual errors that weaken the logical foundation. |
| Scope <-> Solution <-> SC aligned | 20/30 | In Scope items map to the solution. SC items mostly map to In Scope items. But SC-3 (confirmation rounds <= 2) has no corresponding In Scope mechanism. The solution does not include any change to the confirmation flow, so SC-3 is an aspirational metric, not a verifiable outcome. |
| Requirements <-> Solution coherent | 17/25 | The 4 scenarios map to the solution, but Scenario 3 (cross-domain proposal) and Scenario 4 (classification miss) have no corresponding requirements or success criteria. The solution handles them (match best domain / fallback to LLM) but the requirements section does not specify what constitutes success for these scenarios. |

## Phase 3: Blindspot Hunt

1. **[blindspot] Classification table authorship and governance**: The proposal never addresses who creates and maintains the classification table. Is it embedded in the prompt file? Is it a separate configuration? Who decides when a new domain is needed? This is the core artifact and its governance is completely unspecified.

2. **[blindspot] Metric gaming risk**: If domain-level experts are rewarded for broad coverage, the LLM may generate overly generic experts that match everything but add little value. The proposal does not include a quality floor for domain-level experts (e.g., "expert must identify at least 3 proposal-specific attack points").

3. **[blindspot] The 12th expert file**: The proposal itself (`domain-level-freeform-experts`) has already generated an expert (`expert-system-design-prompt-architecture.md`) that is directly relevant to this proposal's domain. This expert was generated by the current system and has domain keywords like "expert-systems, prompt-engineering, classification-taxonomy, reuse-matching, forge-eval-pipeline". These keywords overlap significantly with this proposal's topic. The current system may already be capable of some reuse, and the proposal does not acknowledge this.

4. **[blindspot] Threshold recalibration**: The existing matching threshold (Jaccard >= 0.3 or weighted >= 5) was calibrated for proposal-specific keywords. With broader domain-level keywords, the threshold will produce more matches by design -- but some of these matches may be false positives (broad keyword overlap but wrong focus). The proposal does not address precision vs recall trade-off.

5. **[blindspot] Migration strategy for existing experts**: The proposal says existing experts are "not disturbed" but also says domain-level experts will naturally outcompete them via higher Jaccard scores. This means existing experts will be gradually replaced by domain-level ones. But what happens to the institutional knowledge captured in proposal-specific expert profiles? No migration or knowledge preservation strategy exists.

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| D1. Problem Definition | 70 | 110 |
| D2. Solution Clarity | 75 | 120 |
| D3. Industry Benchmarking | 55 | 120 |
| D4. Requirements Completeness | 55 | 110 |
| D5. Solution Creativity | 55 | 100 |
| D6. Feasibility | 70 | 100 |
| D7. Scope Definition | 50 | 80 |
| D8. Risk Assessment | 55 | 90 |
| D9. Success Criteria | 40 | 80 |
| D10. Logical Consistency | 65 | 90 |
| **Total** | **590** | **1000** |

### Deduction Details

- **Vague language**: "显著降低评审启动成本" (no quantification) -20 pts from D1
- **Factual error**: "11 个专家" (there are 12) and "复用匹配从未成功过" (one success exists) -20 pts from D1
- **Vague language**: "工作量小" (no quantification) -20 pts from D6
- **Vague language**: "领域交集" (undefined term) -20 pts from D9
- **Straw-man alternative**: "纯 Prompt 重写" dismissed without evidence -20 pts from D3
