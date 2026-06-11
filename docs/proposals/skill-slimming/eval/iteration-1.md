# Proposal Evaluation: Iteration 1

**Document**: `skill-slimming/proposal.md`
**Date**: 2026-05-20
**Scorer**: CTO persona, adversarial mode
**Iteration**: 1

---

## Phase 1: Reasoning Audit (Pre-Score Anchors)

### Anchor 1: Problem-Solution Mismatch (Partial)
The problem statement conflates two distinct issues: (a) large files causing "LLM context waste" and (b) "instruction ambiguity" (noTest/doc* confusion). The solution addresses (a) thoroughly via splitting, but addresses (b) only as a side-effect tag-along ("消歧" is listed as a goal but never scoped with specific ambiguous items to fix). The solution is primarily a file-splitting exercise dressed up as a comprehensive "three-layer" slimming method.

### Anchor 2: Evidence-Quality Gap
Top-3 average is stated as 495 lines. Actual: (607+472+407)/3 = 495.3. This is accurate. However, the "~6700 lines total" claim is wrong — actual total is 6394 lines. The "21 skills" count is wrong — there are 22 skill directories. These are verifiable factual errors in the opening paragraph that undermine confidence.

### Anchor 3: Success Criteria Test Easy Proxies
The success criteria measure line counts and path existence — easy-to-measure mechanical properties. None measure whether the slimming actually improves LLM instruction-following accuracy, reduces agent execution errors, or speeds up development. The core problem ("LLM context waste" and "execution deviation risk") is never tested by the success criteria.

### Anchor 4: Self-Contradiction Risk
The proposal requires obeying `skill-self-containment.md` (SKILL.md must contain complete flow steps) while simultaneously splitting content out to rules/templates. The boundary between "complete flow steps" and "rules detail" is undefined and subjective — a potential source of inconsistency across 9 independent tasks.

---

## Phase 2: Dimension Scoring

### 1. Problem Definition (110 pts)

**Problem stated clearly: 30/40**
The core problem is stated in the opening paragraph: 21 SKILL.md files totaling ~6700 lines, large files mixing flow instructions with business rules, and instruction ambiguity. However, the problem conflates two issues (size and ambiguity) without establishing their relative priority or causal relationship. Quote: "导致 LLM 上下文浪费且维护困难" — "LLM context waste" is asserted but never quantified. How much context is wasted? What does "waste" mean — token cost, instruction dilution, or something else?

Deduction: -10 for conflating two problems without priority, -0 for ambiguity about "waste" quantification (partial; the concept is understandable).

**Evidence provided: 28/40**
Four bullet points of evidence are given. Top-3 average line count (495) is verifiable and correct. The consolidate-specs 607-line claim is correct. However: (1) The "~6700 total lines" claim is inaccurate — actual is 6394, a 4.8% overstatement. (2) The "21 skills" count is wrong — there are 22 directories. (3) The noTest/doc* ambiguity is mentioned but no specific example or quote is provided. (4) "多个 skill 内嵌大量模板文本和解释性段落" is vague — which skills, how many lines of template text?

Deduction: -8 for factual errors in data claims, -4 for vague evidence items without specifics.

**Urgency justified: 22/30**
Quote: "v3.0.0 重构窗口期。已有 5 个瘦身相关提案均未执行——方向分散、范围过大是主因。" This provides context but not urgency. What happens if we delay past v3.0.0? What concrete cost does the 6394 lines impose per day/week? The "5 failed proposals" argument actually weakens urgency — it suggests the problem has been lived with for a while without catastrophic consequences.

Deduction: -8 for no concrete cost of delay, -0 for weak urgency signal.

**Dimension Total: 80/110**

---

### 2. Solution Clarity (120 pts)

**Approach is concrete: 32/40**
The three-tier approach (by file size) is clearly described. Nine task groups are enumerated with specific skill names and line counts. A reader can explain back what will happen. However, the actual transformation for each skill is vaguely described as "拆分/精简/消歧" — it's unclear which specific content moves to which file for any given skill. The proposal lacks even one worked example showing "before → after" for a single skill.

Deduction: -8 for no concrete transformation example.

**User-facing behavior described: 30/45**
The Requirements Analysis section lists three scenarios, but they describe developer/maintainer experience rather than end-user (agent) behavior. What does the agent experience differently after slimming? The proposal says "Agent 加载 skill 后获得精简、无歧义的指令" but does not describe how agent behavior changes — does it follow instructions more accurately? Does it produce different outputs? Does it skip steps less often?

Deduction: -15 for absent agent-facing behavioral description.

**Technical direction clear: 28/35**
The general approach (split SKILL.md into SKILL.md + rules/ or templates/) is stated. The constraint about `skill-self-containment.md` is referenced. However, there is no description of the splitting heuristic — how does the implementer decide what stays in SKILL.md vs. what moves to rules/? The proposal says "SKILL.md 保留流程骨架 + 关键约束" but "流程骨架" and "关键约束" are undefined.

Deduction: -7 for undefined splitting heuristic.

**Dimension Total: 90/120**

---

### 3. Industry Benchmarking (120 pts)

**Industry solutions referenced: 18/40**
Quote: "OpenAI 的 GPT best practices 建议将 system prompt 控制在关键指令内，详细规则外置。" This is a single reference to a general best-practices guideline, not a concrete industry solution or pattern. No open-source projects, no published architectures, no product names (beyond OpenAI as a company). No reference to prompt engineering frameworks (e.g., DSPy, LangChain prompt templates, Claude's own tool-use patterns).

Deduction: -22 for thin, single-source industry reference.

**At least 3 meaningful alternatives: 18/30**
Three alternatives are listed: "Do nothing", "模式审计 + 批量清理", "按大小分层逐组处理". The "模式审计 + 批量清理" alternative references a prior proposal ("skill-slim-down") but provides no details — it's presented only to be rejected as "太理论化". This is a straw man. Missing: alternative approaches like (1) automated prompt compression via LLM, (2) hierarchical prompt loading (lazy-load rules only when needed), (3) shared rule files across skills, (4) template inheritance.

Deduction: -12 for straw-man alternative (per rubric rule: -20 for straw man, but partial credit for having three rows).

**Honest trade-off comparison: 15/25**
The comparison table has Pros/Cons columns but they are shallow. "小组内 skill 可能需不同策略" is the only con for the selected approach — a minor concern. The cons for rejected alternatives are dismissive ("已搁置", "债务积累") rather than analytical.

Deduction: -10 for shallow trade-off analysis.

**Chosen approach justified against benchmarks: 12/25**
The proposal does not explain why the selected approach outperforms industry patterns. It claims "增量重构最佳实践" as the source but does not cite any specific methodology or reference. No benchmark comparison with actual industry practices.

Deduction: -13 for absent justification against benchmarks.

**Dimension Total: 63/120**

---

### 4. Requirements Completeness (110 pts)

**Scenario coverage: 28/40**
Three scenarios are listed (agent loading, developer maintenance, large skill splitting). Edge cases and error scenarios are absent: What if a skill's flow instructions are tightly interleaved with rules and cannot be cleanly separated? What if splitting a skill introduces circular references between SKILL.md and rules/? What if the 350-line target is unreachable for a complex skill without losing critical instructions?

Deduction: -12 for missing edge cases and error scenarios.

**Non-functional requirements: 30/40**
Three NFRs are stated: 350-line cap, rules/templates subdirectory placement, and no I/O contract change. These are reasonable but incomplete. Missing: (1) Performance — does file splitting affect skill loading time or token efficiency? (2) Compatibility — will existing scripts or hooks that reference SKILL.md content break? (3) Backward compatibility for any external tooling.

Deduction: -10 for incomplete NFR coverage.

**Constraints & dependencies: 24/30**
Four constraints are listed: forge-distribution.md, skill-self-containment.md, no Go source changes, no skill merging. Good. However, the dependency on `skill-self-containment.md` is listed without acknowledging the tension: the document's principle of "SKILL.md must contain complete flow steps" may conflict with aggressive splitting. This unstated tension is a constraint risk.

Deduction: -6 for unaddressed constraint tension.

**Dimension Total: 82/110**

---

### 5. Solution Creativity (100 pts)

**Novelty over industry baseline: 20/40**
The "三层瘦身法" (split, slim, disambiguate) is presented as an innovation but is standard refactoring practice: extract to modules, remove dead code, clarify naming. There is no differentiation from basic file-splitting work. The tier-based grouping (by file size) is a simple organizational heuristic, not a creative insight.

Deduction: -20 for lack of genuine novelty.

**Cross-domain inspiration: 10/35**
No evidence of borrowing from other domains. The proposal could have drawn from: microservices decomposition patterns, database normalization principles, progressive loading in frontend architecture, or information architecture principles. None are referenced.

Deduction: -25 for absent cross-domain thinking.

**Simplicity of insight: 18/25**
The approach is straightforward and not overengineered — that's a strength. The insight that "split by file size, process in groups" is simple and practical. However, it's also unremarkable — it's the most obvious approach one would take.

Deduction: -7 for being obvious rather than elegantly insightful.

**Dimension Total: 48/100**

---

### 6. Feasibility (100 pts)

**Technical feasibility: 35/40**
Pure text file operations, git for rollback, existing patterns for rules/templates subdirectories. The approach is technically straightforward. Minor concern: the proposal assumes all skills can be cleanly split without structural issues, which may not hold for tightly-coupled skill files.

Deduction: -5 for unverified assumption that clean splitting is always possible.

**Resource & timeline feasibility: 22/30**
Nine tasks are estimated but no timeline is given — no days, no sprints, no deadline. "预计 9 个任务" is a count, not a schedule. There is no estimate of effort per task or total project duration. Without a timeline, feasibility cannot be assessed.

Deduction: -8 for absent timeline.

**Dependency readiness: 25/30**
Quote: "无外部依赖。所有文件已在本地。" This is accurate — the work is self-contained. However, the dependency on `skill-self-containment.md`'s principles (which define what can/cannot be split out) is acknowledged but not analyzed for readiness — is the document up to date? Does it cover splitting scenarios?

Deduction: -5 for unverified convention readiness.

**Dimension Total: 82/100**

---

### 7. Scope Definition (80 pts)

**In-scope items are concrete: 24/30**
"21 个 SKILL.md 文件的拆分、精简、消歧" is clear. Creating rules/templates subdirectories is concrete. However, "清理过时标签、路径引用和歧义描述" is vague — which tags, which paths, which descriptions?

Deduction: -6 for vague cleanup scope.

**Out-of-scope explicitly listed: 22/25**
Five out-of-scope items are listed: Go source, I/O contracts, skill merging, commands/agents, hooks/references/scripts. Good coverage. Minor gap: are test files in scope? What about documentation outside SKILL.md (e.g., README.md in skill directories)?

Deduction: -3 for incomplete exclusion boundaries.

**Scope is bounded: 18/25**
Nine task groups are enumerated, which provides structure. However, the scope says "21 个 skills" when there are 22 directories. This factual error means the scope is literally miscounted. Additionally, no timeline bounds the work — without a deadline or timebox, the scope is open-ended despite being grouped.

Deduction: -5 for incorrect skill count, -2 for no time boundary.

**Dimension Total: 64/80**

---

### 8. Risk Assessment (90 pts)

**Risks identified: 22/30**
Three risks are listed. The first (key instruction loss) is meaningful. The second (inconsistent naming) is low-impact and arguably trivial. The third (introducing new ambiguity) is relevant. Missing risks: (1) Splitting may violate skill-self-containment principle — a stated constraint. (2) Nine independent tasks may diverge in splitting style, creating inconsistency. (3) No automated regression test for "does the agent still follow instructions correctly after splitting."

Deduction: -8 for missing high-impact risks.

**Likelihood + impact rated: 24/30**
Ratings are provided (M/H, M/L, L/M). The first risk (key instruction loss, M/H) is honest. However, the second risk (naming inconsistency, M/L) is rated as Medium likelihood but the mitigation (a convention) is trivial — this risk should be Low likelihood given the mitigation. Assessment is mostly honest but slightly padded.

Deduction: -6 for slight rating inflation.

**Mitigations are actionable: 18/30**
First mitigation: "SKILL.md 保留流程骨架 + 关键约束，辅助文件仅放规则细节" — this is a design principle, not an actionable mitigation. How do you verify no key instruction was lost? Second mitigation: "约定 rules/ 放规则、templates/ 放模板" — a convention, not a verification mechanism. Third mitigation: "commit message 中注明原文和修改理由" — this is documentation, not prevention. None of the mitigations include verification or testing steps.

Deduction: -12 for non-actionable mitigations.

**Dimension Total: 64/90**

---

### 9. Success Criteria (80 pts)

**Criteria are measurable and testable: 35/55**
Four criteria are listed. Three are mechanically testable: (1) line count per file ≤ 350, (2) total line reduction ≥ 25%, (3) no broken internal references. The fourth (each commit touches only one group) is also testable. However, the most important criterion — "agent instruction-following accuracy is maintained or improved" — is absent. The criteria measure structural properties, not functional correctness. A 300-line SKILL.md that lost critical instructions would pass all four criteria.

Deduction: -20 for missing functional correctness criteria.

**Coverage is complete: 15/25**
The criteria cover the structural goals (line counts, file organization) but do not cover the disambiguation goal at all. The proposal promises "消歧" as a core operation, but no success criterion measures whether ambiguity was actually resolved. Additionally, the "不改变 skill 的输入/输出契约" NFR has no corresponding verification criterion.

Deduction: -10 for missing disambiguation criteria.

**Dimension Total: 50/80**

---

### 10. Logical Consistency (90 pts)

**Solution addresses the stated problem: 25/35**
The stated problem has two parts: (a) large files causing context waste, (b) instruction ambiguity. The solution addresses (a) well via splitting. It addresses (b) only nominally — "消歧" is mentioned as a task step but never scoped with specific ambiguous items. The proposal does not identify which skills have ambiguity issues or what the ambiguous instructions are. The noTest/doc* example is mentioned in the evidence section but never appears in the solution, scope, or success criteria.

Deduction: -10 for solution only partially addressing the stated problem.

**Scope ↔ Solution ↔ Success Criteria aligned: 22/30**
The scope lists 21 skills (actually 22), the solution proposes 9 task groups covering those skills, and the success criteria measure structural outcomes. However: the scope includes "消歧" as an activity, but no success criterion tests for it. The scope includes "清理过时标签、路径引用", but no success criterion verifies tag cleanup. The solution's "三层瘦身法" claims three operations but success criteria only test two (split and slim).

Deduction: -8 for misalignment on disambiguation coverage.

**Requirements ↔ Solution coherent: 20/25**
The requirements (350-line cap, subdirectory structure, no I/O contract change) map to the solution approach. However, the requirement "SKILL.md 必须包含完整流程步骤" (from the constraint referencing skill-self-containment.md) creates tension with the splitting approach — where is the line drawn? The proposal does not address this coherence gap.

Deduction: -5 for unaddressed requirement-solution tension.

**Dimension Total: 67/90**

---

## Phase 3: Blindspot Hunt

### [blindspot-1] Factually Incorrect Baseline Data Undermines Entire Proposal

The proposal states "21 个 SKILL.md 文件总计 ~6700 行" in its opening sentence. Actual: 22 files totaling 6394 lines. This is not a rounding error — it is a miscount (22 vs 21) and a 4.8% inflation of total lines. Since the entire proposal's justification (task grouping, success criteria of "25% reduction to 5000 lines") rests on these baseline numbers, the error propagates through the success criteria. The 25% reduction target from 6700 yields 5025, but from the actual 6394, a 25% reduction yields 4796. The target "5000 行以下" would require only a 21.8% reduction from actual baseline. The proposal's own success criterion may be easier to achieve than presented, which inflates the perceived ambition.

Quote: "Forge plugin 的 21 个 SKILL.md 文件总计 ~6700 行"

### [blindspot-2] No Rollback Plan Beyond "Git"

The CTO failure pattern checklist explicitly flags "Missing rollback plans for infrastructure or architecture changes." The proposal says "git 提供完整回滚能力" and "每个任务独立 commit，可逐个验证回滚." But git revert only restores file content — it does not address: (1) What if the agent was already used with a split SKILL.md and produced incorrect results? (2) How do you detect that a split degraded agent behavior? (3) What is the rollback trigger — who decides and when? The proposal has no monitoring or detection mechanism for regression.

Quote: "纯文本修改 + 文件拆分。git 提供完整回滚能力。"

### [blindspot-3] "三层瘦身法" Masks Absence of Disambiguation Strategy

The proposal's "Innovation Highlight" claims a "三层瘦身法：拆分、精简、消歧" but the document provides zero detail on how disambiguation will be performed. There is no list of ambiguous terms, no examples of ambiguous instructions, no methodology for identifying ambiguity, and no success criterion for disambiguation. The noTest/doc* issue is mentioned once in the Evidence section and never appears again. The "three-layer method" branding disguises that two layers (split, slim) are well-defined while the third (disambiguate) is a placeholder.

Quote: "三层瘦身法：对每个 skill 按需施以拆分（大文件）、精简（冗余文本）、消歧（模糊指令）三种操作"

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 80 | 110 |
| 2. Solution Clarity | 90 | 120 |
| 3. Industry Benchmarking | 63 | 120 |
| 4. Requirements Completeness | 82 | 110 |
| 5. Solution Creativity | 48 | 100 |
| 6. Feasibility | 82 | 100 |
| 7. Scope Definition | 64 | 80 |
| 8. Risk Assessment | 64 | 90 |
| 9. Success Criteria | 50 | 80 |
| 10. Logical Consistency | 67 | 90 |
| **Total** | **690** | **1000** |

---

## Attack Points

1. [Industry Benchmarking]: Single vague reference, no real benchmarks — "OpenAI 的 GPT best practices 建议将 system prompt 控制在关键指令内" — Cite at least 2-3 concrete industry patterns or open-source projects with specific URL references, and compare the chosen approach against them.

2. [Industry Benchmarking]: Straw-man alternative — "模式审计 + 批量清理...已搁置 Rejected: 太理论化" — Provide substantive alternatives with real analysis, not dismissed prior proposals. Include at least one industry-validated approach (e.g., prompt compression, lazy loading).

3. [Solution Creativity]: No cross-domain inspiration — the entire Innovation Highlights section offers only "三层瘦身法" which is standard refactoring — Reference decomposition patterns from other domains (microservices, information architecture, progressive loading) and articulate what this proposal borrows or adapts.

4. [Success Criteria]: Functional correctness is untested — criteria only measure line counts and file existence — Add at least one criterion that verifies agent behavior after splitting (e.g., "each modified skill passes its existing eval suite" or "agent instruction-following accuracy ≥ baseline").

5. [Success Criteria]: Disambiguation has no measurable criterion — the proposal promises 消歧 as a core operation but no criterion tests for it — Add a criterion such as "all identified ambiguous terms (e.g., noTest, doc*) are resolved with documented before/after definitions."

6. [Risk Assessment]: Mitigations are design principles, not actionable steps — "SKILL.md 保留流程骨架 + 关键约束，辅助文件仅放规则细节" — Provide concrete verification mechanisms: diff review checklists, automated reference-integrity tests, or pairwise comparison of pre/post agent outputs.

7. [Problem Definition]: Factual errors in baseline data — "21 个 SKILL.md 文件总计 ~6700 行" — actual count is 22 files totaling 6394 lines. Correct all baseline numbers and recalculate the 25% reduction target.

8. [Scope Definition]: Incorrect skill count propagates through scope — "21 个 skills" appears in In Scope, Scope boundaries, and Task Grouping — Verify and correct the count to 22 throughout the document.

9. [Logical Consistency]: Solution only partially addresses the stated problem — ambiguity (noTest/doc*) is listed as a core problem but never scoped in the solution — Either remove ambiguity from the problem statement or add concrete disambiguation items to the task descriptions.

10. [blindspot]: No regression detection mechanism — "git 提供完整回滚能力" is not a rollback plan — Define what constitutes a regression, how it will be detected, and what the rollback trigger is. Git can revert files but cannot tell you when to revert.

11. [blindspot]: Disambiguation is a labeled but undefined operation — "三层瘦身法：拆分（大文件）、精简（冗余文本）、消歧（模糊指令）" — The third layer has no methodology, no examples, and no success criterion. Either define it concretely (with specific ambiguous items to resolve) or remove it from the proposal's claims.

12. [blindspot]: Missing timeline makes feasibility unverifiable — "21 个 skill 分 9 组，预计 9 个任务" gives a count but no schedule — Add estimated duration per task group and a total project timeline, so reviewers can assess whether the scope is realistic.

---

## Improvement Priority (Top 3)

1. **Fix baseline data**: Correct "21 skills, ~6700 lines" to "22 skills, 6394 lines" and recalculate targets.
2. **Add functional success criteria**: At minimum, a criterion verifying agent behavior is preserved after splitting.
3. **Define or remove disambiguation**: Either scope specific ambiguous items with a resolution plan, or remove the "消歧" claim from the solution.
