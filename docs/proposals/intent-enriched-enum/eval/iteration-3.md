# Eval Report: Intent Enriched Enum — Iteration 3

**Reviewer**: CTO Adversary
**Date**: 2026-05-31
**Mode**: Adversarial Re-Evaluation (targeting remaining weaknesses from Iteration 2)

---

## Iteration 2 Issue Tracker

| # | Attack | Status | Evidence |
|---|--------|--------|----------|
| 1 | Urgency lacks frequency data | FIXED | Added explicit "关于频率数据的说明" paragraph acknowledging the gap, explaining why data is unavailable (no pipeline execution logs), and providing severity-based justification instead of frequency-based. Honest framing. |
| 2 | Override detection temporal ordering misleading | FIXED | Reworded from "在 Pipeline Configuration 应用后，扫描 PRD 正文段落" to "LLM 在生成 PRD/tech-design 内容的过程中，同步（而非先后）完成信号检测" with explicit note "不存在'先生成再扫描'的时序关系". |
| 3 | `doc` Minimal PRD format underspecified | FIXED | Now contains full description: "标题：一句话描述文档变更对象和目的；目标：列出要更新/新增的文档文件和预期变更点；scope：界定变更涉及的文档范围和不涉及的边界". Parity with enhancement format. |
| 4 | CI lint gate analogy ignores determinism gap | FIXED | Added dedicated "与 CI lint gate 的关键差异" paragraph with 3-point justification: (1) Forge entire chain depends on LLM anyway, (2) override is additive-only so worst case is benign, (3) structured table reduces ambiguity. |
| 5 | "完全内容驱动" rejection rationale contradicts chosen approach | FIXED | Reworked to: "被拒绝的核心原因不是'LLM 判断不可靠'（本方案的 override 同样依赖 LLM），而是'没有基线'" — explicitly names the real differentiator and avoids false implication. |
| 6 | No testability NFR | FIXED | Added "可测试性：每个 override signal 需至少一个 PRD 输入→pipeline 输出的测试用例" with cross-reference to SC. |
| 7 | No scenario for invalid intent fallback | FIXED | Added Scenario 8 with full description of AskUserQuestion structured options, best-effort fallback mapping, and root-cause note about brainstorm output format constraints. |
| 8 | Override annotation protocol has no owner | FIXED | Scope now explicitly lists write-prd/SKILL.md and tech-design/SKILL.md as "实现 `<!-- Override: ... -->` 注释行的生成逻辑". Implementation owner identified. |
| 9 | Eval skills not mentioned in scope | FIXED | Added to Out of Scope: "eval 系列 skill（eval-prd、eval-design 等）的 rubric/contract 更新" with justification "这些 skill 评估的是产物质量而非 pipeline 行为". |
| 10 | Risk #3 likelihood/mitigation contradiction (关键词误触发) | FIXED | Re-rated from M/L to L/L. Mitigation strengthened: "同一 LLM 在处理 Pipeline Configuration 时已具备足够的上下文推理能力（否则整个 pipeline 配置逻辑都不可靠）". Now internally consistent. |
| 11 | Risk #5 mitigation is circular | FIXED | Replaced with: "验证策略：对旧 3 个 intent 分别运行变更前的 write-prd/tech-design 测试用例，对比变更后输出中的 pipeline 相关产物是否与变更前一致". Concrete testing strategy. |
| 12 | SC #7 conflates table structure with behavior | FIXED | Reworded to: "对旧 3 个 intent 分别执行变更后的 write-prd 和 tech-design，验证生成的产物集合（PRD sections、checklist items）与变更前一致". Now behavioral, not structural. |
| 13 | No SC for Architecture Decision | FIXED | SC now includes: "breakdown-tasks 将 fix intent 映射到 coding.fix task type（验证 Type Assignment 表更新）；brainstorm 将 coding.feature → new-feature 和 coding.enhancement → enhancement 作为独立路径（验证 intent mapping 表 split）". |
| 14 | Effective pipeline count is 4, not 6 | FIXED | SC #3 now explicitly states "产生 4 种功能不同的 pipeline 配置：Full、Simplified、Spec-only、Minimal". Honest framing. |
| 15 | brainstorm UX change not acknowledged | FIXED | NFR now has two separate bullets: "Pipeline 向后兼容" and "交互向前演进" with explicit note about brainstorm 3→6 options and parenthetical "非破坏性——旧 3 值仍在选项中". |
| 16 | doc intent override no-op is incidental, not designed | FIXED | Scenario 7 now contains: "设计约束：这是显式的设计决策而非偶然属性——doc pipeline 的 Minimal 格式刻意不包含任何可被 override 开启的检查项...如果未来需要为 doc pipeline 添加可被 override 的检查项，需重新评估此约束". |

**Summary**: All 16 issues from Iteration 2 have been addressed. 16/16 fully fixed, 0 partially fixed, 0 not fixed. The proposal has been substantially strengthened. Iteration 3 attacks focus on residual weaknesses and new blindspots exposed by the revisions.

---

## Phase 1: Reasoning Audit

**Problem → Solution → Evidence → SC chain**:
- Problem: 3-value intent → 5/8 task types unmapped → pipeline squashes 5 types into 1 branch
- Solution: 6-value enum + hybrid pipeline (intent baseline + content override) + keyword-based signals
- Evidence: Codebase grep confirms 5 unmapped types, heuristic at brainstorm/SKILL.md:96, identical treatment in write-prd/tech-design
- SC: 9 criteria, most verifiable with concrete test procedures

**Chain integrity**: The chain is now sound. All major gaps from Iterations 1-2 have been closed. The frequency-data gap is honestly acknowledged and mitigated with severity-based reasoning. The temporal ordering confusion is resolved. The override annotation protocol has an owner. The SC set covers the Architecture Decision. Remaining weaknesses are second-order issues.

---

## Phase 2: Rubric Scoring

### 1. Problem Definition (96/110)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Problem stated clearly | 40/40 | "8 个 task type 中有 5 个被压缩进同一条分支" — quantified, specific, unambiguous. Two readers would interpret this identically. |
| Evidence provided | 35/40 | Four concrete evidence items verified against codebase. "关于频率数据的说明" is a model of honest gap acknowledgment — explains why data is unavailable (no pipeline execution logs) and substitutes severity argument. Deduction: the severity argument ("修复成本和遗漏成本都显著高于一次性修复映射的成本") is strong but still lacks a single concrete example of "I ran brainstorm, got heuristic miss, had to manually fix the PRD" from a real session. The API handbook example is hypothetical, not retrospective. |
| Urgency justified | 21/30 | Substantially improved. The severity argument (cost of each miss > cost of fix) is valid. The honest admission about missing frequency data is a strength. Remaining gap: the urgency is framed entirely in abstract cost terms. One retrospective anecdote ("In the last project, 2 out of 5 brainstorms produced wrong pipeline outputs that required manual correction") would have been more persuasive than the theoretical cost analysis. |

**Attacks**:
- The "关于频率数据的说明" paragraph is excellent practice but reveals a meta-weakness: Forge has no pipeline execution telemetry. This means post-deployment, the proposal's success cannot be measured by "reduction in heuristic miss rate" either — there's no baseline AND no post-deployment measurement capability. SC #8 verifies implementation correctness, not operational impact.

### 2. Solution Clarity (112/120)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Approach is concrete | 40/40 | 6-value enum, explicit mapping, Pipeline Configuration table, 5 override signals, negation handling, multi-signal stacking, override annotations — all fully specified. |
| User-facing behavior described | 43/45 | Override annotation gives user visibility. Enhancement gets Simplified PRD with explicit section list. Doc gets Minimal PRD with 3 described sections. Fix gets spec-only with reproduce→fix→verify test pipeline. Gap: the `<!-- Override: ... -->` annotation is described as being generated for user review, but the user's action on seeing an incorrect override is not specified. Can the user delete the annotation line and re-run? Edit it? The review workflow is incomplete. |
| Technical direction clear | 29/35 | Markdown editing approach clear. LLM-based negation handling explained with "同步（而非先后）" temporal model. Gap: the Scenario 8 invalid intent fallback says "LLM 将回退到最接近的合法值" — who implements this fallback? The brainstorm SKILL.md? The downstream skill's Pipeline Configuration step? There's no scope item for implementing fallback logic. If it's "the LLM just figures it out," that's not a technical direction. |

**Attacks**:
- Scenario 8's fallback mechanism is a critical robustness feature but has no implementation path. The proposal says "根本解决方案依赖 brainstorm SKILL.md 中对输出格式的明确约束" — but scope item #1 only says "更新 Step 4.5 intent mapping 表（6 值）, 移除 fix 启发式, 更新 AskUserQuestion 选项". It doesn't say "add output format constraint to prevent invalid intent values." The fallback and its root fix are described but not scoped.

### 3. Industry Benchmarking (98/120)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Industry solutions referenced | 36/40 | GitHub label timeline with specific dates (2013/2016/2020), TypeScript 2.0 diagnostic categories with behavioral differentiation. CI lint gate analogy with ESLint overrides and GitHub Actions path-based triggers. The "与 CI lint gate 的关键差异" paragraph is a model of honest analogy analysis — names the determinism gap and justifies why it's acceptable. Remaining gap: no citation of pipeline routing systems that are closer analogs (Buildkite step overrides, CircleCI dynamic config, Tekton trigger bindings). These would be more relevant than CI linting. |
| At least 3 meaningful alternatives | 25/30 | 4 alternatives. "只扩枚举" con is now technically accurate and specific ("pipeline 分支仍只有 2 条，5 个 intent 挤在 Spec-only 分支里"). "完全内容驱动" now has honest rejection rationale focused on "no baseline" rather than "LLM unreliable". Remaining gap: no analysis of when "只扩枚举" might be acceptable (e.g., team that only needs brainstorm fix and doesn't care about pipeline granularity). |
| Honest trade-off comparison | 18/25 | Improved: selected approach's cons now mention two-copy table sync. The keyword table governance question from Iteration 2 remains unanswered — who adds/removes keywords? What's the review process? The Override Signals table is effectively a growing ruleset. |
| Chosen approach justified | 19/25 | CI lint gate analogy is well-justified with explicit determinism acknowledgment. The 3-point justification for accepting the determinism gap is strong. Gap: no discussion of alternative override signal mechanisms that might be more deterministic (e.g., structured PRD metadata fields instead of natural language keywords). The jump from "keywords" to "LLM interprets" has no intermediate option explored. |

**Attacks**:
- The keyword table governance gap is now the most significant remaining weakness in Industry Benchmarking. The Override Signals table has 5 rows with specific keywords. Who decides when to add a 6th? When to change existing keywords? What's the review/approval process? This is a maintenance governance question that the proposal's "可维护性" NFR acknowledges for the Pipeline Configuration table but not for the Override Signals table.

### 4. Requirements Completeness (98/110)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Scenario coverage | 38/40 | 8 scenarios including multi-signal, intent-content mismatch, and invalid intent fallback. Coverage is comprehensive. Gap: no scenario for "refactor intent with NO override signals" — the default case. All refactor scenarios in the proposal involve override signals (API, performance), but the most common case (a refactor that doesn't touch APIs or performance) is assumed rather than tested. |
| Non-functional requirements | 38/40 | Excellent: pipeline backward compatibility, interaction forward evolution (distinguished separately), consistency, maintainability (with two-copy acknowledgment), testability (with SC cross-reference). Gap: no NFR for the Override Signals keyword table's governance — covered above in Industry Benchmarking. |
| Constraints & dependencies | 22/30 | Good: 8 files, no Go code, task-doc.md false match excluded. Remaining gap from Iteration 2 not fully addressed: the "grep intent 匹配的 skill 文件" constraint lists 8 files, but the actual codebase shows `gen-contracts/SKILL.md` and `deep-research/SKILL.md` also reference "intent" (albeit in different semantic contexts). The proposal correctly excludes these from scope but doesn't document why they're excluded. |

**Attacks**:
- Codebase verification reveals 7 skill files reference "intent": brainstorm, write-prd, tech-design, breakdown-tasks, quick-tasks, gen-contracts, deep-research. The proposal's scope lists 5 skills (8 files including rules). gen-contracts and deep-research are excluded without justification. gen-contracts uses "intent" in "Parameters serve as prefill hints... do not override the user's confirmed intent" — this is semantically relevant to the proposal. deep-research uses "intent" in "expressing business intent" — this is unrelated. The proposal should explicitly justify excluding gen-contracts.

### 5. Solution Creativity (42/100)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Novelty over industry baseline | 16/40 | Still honest: "无特别创新". The hybrid baseline+override pattern is standard. The Assumptions Challenged section shows deeper analytical thinking but doesn't translate into solution novelty. |
| Cross-domain inspiration | 14/35 | CI lint gate analogy is cross-domain. The "determinism gap" analysis is thoughtful. Still no borrowing from feature flag systems (progressive rollout patterns), test impact analysis, or other adjacent domains. |
| Simplicity of insight | 12/25 | The core insight is still "expand the enum to match the task types." The most elegant element is the additive-only override principle, which is borrowed from CI lint gates. The Assumptions Challenged table is the most thoughtful piece but is analytical, not creative. |

**Attacks**:
- The Assumptions Challenged section is well-crafted but is better classified as "rigorous analysis" than "creative insight." This dimension is inherently weak for this proposal because the problem (enum mismatch) has an obvious solution (expand the enum). The creativity ceiling is low by nature of the problem domain.

### 6. Feasibility (92/100)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Technical feasibility | 39/40 | All markdown, no Go code. Verified against codebase. The LLM-dependent override detection is explicitly acknowledged as relying on the same capability the entire skill chain depends on. One gap: the invalid intent fallback in Scenario 8 relies on "LLM 将回退到最接近的合法值" — this is undefined behavior masquerading as a feature. |
| Resource & timeline | 26/30 | 8 files, 2-3 tasks. The override annotation mechanism adds complexity not accounted for in the original estimate. Testing across 6 intents × 5 signals = 30 override combinations (not all meaningful, but the test matrix is larger than "2-3 tasks" suggests). |
| Dependency readiness | 27/30 | No external dependencies. The override annotation consumer question from Iteration 2 is resolved — the scope now identifies write-prd and tech-design as the generators. Gap: no consumer beyond "user review" is specified. If no downstream process reads the annotations, they're documentation-only, which is fine but should be explicit. |

**Attacks**:
- The test matrix for override signals is understated. The proposal defines 6 intents × 5 signals, and claims additive-only means "worst case is extra output." But the override detection runs within the LLM's generation step — each intent+signal combination could theoretically produce different behavior. The "2-3 tasks" estimate doesn't account for verifying this matrix.

### 7. Scope Definition (74/80)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| In-scope items are concrete | 29/30 | 8 specific files with specific changes. Override annotation implementation assigned to write-prd and tech-design. breakdown-tasks Type Assignment update included. Comprehensive. |
| Out-of-scope explicitly listed | 23/25 | 5 items including eval skills with justification. Gap: gen-contracts/SKILL.md references "intent" semantically relevant to the proposal but is not mentioned in either in-scope or out-of-scope. |
| Scope is bounded | 22/25 | "2-3 tasks" is bounded. Two-copy sync acknowledged. Override Signals keyword table is open-ended maintenance surface. The proposal says "修改 SKILL.md 中的结构化表格即可" for keyword changes, which is bounded technically but not governance-wise. |

### 8. Risk Assessment (82/90)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Risks identified | 27/30 | 6 risks. Invalid intent fallback now addressed via Scenario 8. Gap: no risk for gen-contracts/SKILL.md containing intent references that might need updating. |
| Likelihood + impact rated | 26/30 | "关键词误触发" re-rated to L/L — now consistent with mitigation (LLM handles negation) and impact (additive-only = benign). "6 值枚举仍不够" at M/L is honest. Well-calibrated overall. |
| Mitigations are actionable | 29/30 | Risk #5 now has concrete testing strategy ("对旧 3 个 intent 分别运行变更前后的 write-prd/tech-design"). Risk #3 has explicit fallback argument. Risk #4 has diff-based verification. Near-excellent. |

**Attacks**:
- Risk assessment is substantially improved. The remaining weakness is minor: Risk #1 ("6 值枚举仍不够") mitigation says "如未来出现新 task type，可追加 intent 值——Pipeline Configuration 表新增一行即可". This understates the cascading cost: adding an intent requires updating all 8 scope files, the Override Signals table, brainstorm AskUserQuestion options, and all downstream mapping tables. "新增一行即可" is misleading about the true maintenance cost.

### 9. Success Criteria (76/80)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Criteria are measurable and testable | 28/30 | SC #8 uses specific grep command. SC #5 specifies input→output test pairs. SC #7 now verifies actual output parity, not table structure. SC #9 (new) covers Architecture Decision. Strong improvement. Gap: SC #5 "对每个 override signal 各提供一个包含对应关键词的 PRD 输入" — this defines 5 test cases but doesn't specify what "验证输出包含对应产物" means concretely. What constitutes "API handbook section" in the output? A heading? A paragraph? An artifact file? |
| Coverage is complete | 24/25 | SC #3 now explicitly names "4 种功能不同的 pipeline 配置". SC #6 covers fix→coding.fix mapping. SC #7 covers backward compatibility. SC #9 covers Architecture Decision split. Gap: no SC for the override annotation format — the `<!-- Override: ... -->` protocol is part of the solution (SC #5 mentions it) but no SC verifies the annotation format is correct, only that it exists. |
| SC internal consistency | 24/25 | SC set is internally consistent. SC #3 (4 distinct configurations) is honest and doesn't inflate. SC #7 (backward compat) and SC #9 (Architecture Decision) don't conflict. One concern: SC #5 defines override trigger testing, but Scenario 8 defines invalid intent fallback — there's no SC for fallback behavior. If brainstorm outputs `bug-fix`, the proposal says it should fall back to `fix`, but no SC verifies this. |

### 10. Logical Consistency (84/90)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Solution addresses stated problem | 32/35 | Strong alignment. Effective pipeline count is honestly stated as 4. The refactor/cleanup differentiation is honestly grounded in override probability and measurement rather than default pipeline difference. Gap: Scenario 8's invalid intent fallback is described but not scoped — the fallback logic ("LLM 将回退到最接近的合法值") has no implementation vehicle. If it's emergent LLM behavior, it should be documented as such; if it requires implementation, it's missing from scope. |
| Scope ↔ Solution ↔ SC aligned | 27/30 | Mostly aligned. Override annotation now in scope and in SC. Architecture Decision now in SC. Gap: the invalid intent fallback described in Scenario 8 is in the solution but not in scope (no file is listed as needing fallback implementation) and not in SC (no criterion for fallback correctness). |
| Requirements ↔ Solution coherent | 25/25 | Good coherence. NFR distinguish pipeline backward compat from interaction evolution. Maintainability NFR acknowledges two-copy burden. Testability NFR cross-references SC. No orphan requirements or solution features without requirements. |

---

## Phase 3: Blindspot Hunt

### Residual Weaknesses

1. **Invalid intent fallback is unimplemented**: Scenario 8 describes fallback behavior ("LLM 将回退到最接近的合法值") but no scope item implements this. The proposal says "根本解决方案依赖 brainstorm SKILL.md 中对输出格式的明确约束" — but no scope item adds this constraint either. The fallback is described as both a feature and a non-goal in the same paragraph. Either scope the brainstorm output format constraint or remove the fallback promise.

2. **gen-contracts/SKILL.md intent reference unaddressed**: Codebase grep shows gen-contracts/SKILL.md line 38 references "intent" semantically: "Parameters serve as prefill hints — they reduce the questions asked in Phase 1 but do not override the user's confirmed intent." This is relevant to the proposal's intent system but is not in scope or out-of-scope. The proposal should either justify exclusion or include it.

3. **Override Signals keyword table governance gap**: The keyword table has 5 rows. The "可维护性" NFR mentions Pipeline Configuration table changes ("修改 SKILL.md 中的结构化表格即可") but not Override Signals table changes. Adding a new keyword or signal type requires the same two-copy sync as the Pipeline Configuration table. This is a governance surface that's acknowledged for one table but not the other.

4. **No retrospective anecdote despite honest frequency gap**: The proposal excellently acknowledges the lack of frequency data but compensates with theoretical severity analysis. A single real-world anecdote ("Last week's refactor proposal generated a Full PRD because the heuristic misclassified it as new-feature, requiring 10 minutes of manual cleanup") would be more persuasive than the abstract cost argument. The proposal author likely has such an anecdote but didn't include it.

5. **Risk #1 mitigation understates cascading cost**: "Pipeline Configuration 表新增一行即可" for adding a new intent value understates the actual cost: 8 files, 2 Pipeline Configuration tables, brainstorm AskUserQuestion options, breakdown-tasks Type Assignment table, quick-tasks mapping, and override signal evaluation for the new intent. The mitigation should be honest about the expansion cost.

6. **SC #5 override verification lacks output specificity**: "验证输出包含对应产物" is still vague about what constitutes the expected output. For "API handbook" — is it a section in the PRD? A separate file? A checklist item? The SC should specify the expected output artifact for each signal.

---

## Deductions

- **Vague language**: SC #5 "验证输出包含对应产物" — "对应产物" is vague. What specific artifact? -10 pts from Success Criteria.
- **Missing scope coverage**: gen-contracts/SKILL.md contains semantically relevant "intent" reference not addressed in scope or out-of-scope. -5 pts from Scope Definition.
- **Unimplemented fallback**: Scenario 8 describes invalid intent fallback but no scope item implements it. The fallback is a solution promise without an implementation path. -5 pts from Logical Consistency.

---

SCORE: 784/1000
DIMENSIONS:
  Problem Definition: 96/110
  Solution Clarity: 112/120
  Industry Benchmarking: 98/120
  Requirements Completeness: 98/110
  Solution Creativity: 42/100
  Feasibility: 92/100
  Scope Definition: 74/80
  Risk Assessment: 82/90
  Success Criteria: 76/80
  Logical Consistency: 84/90
ATTACKS:
1. [Problem Definition]: No retrospective anecdote despite honest frequency gap — the "关于频率数据的说明" is excellent but substitutes theoretical severity analysis for a single concrete real-world example. Quote: "bug fix 被当作 new-feature 会生成不必要的 User Stories 和 Full PRD（浪费 token + 用户困惑）" — this is hypothetical. Adding one real session example ("In session X, brainstorm misclassified a fix as new-feature, resulting in...") would strengthen the severity argument measurably.
2. [Solution Clarity]: Override annotation review workflow incomplete — "供用户 review 时确认" describes user review but not the user's action on seeing an incorrect override. Can the user delete the annotation and re-run? Edit the PRD? The review workflow stops at "show annotation" without defining the remediation path.
3. [Solution Clarity]: Invalid intent fallback has no implementation path — Scenario 8 says "LLM 将回退到最接近的合法值（如 bug-fix → fix）" but no scope item implements this fallback logic. Quote: "根本解决方案依赖 brainstorm SKILL.md 中对输出格式的明确约束" — but no scope item adds this constraint. Either scope the brainstorm output format constraint or remove the fallback promise.
4. [Industry Benchmarking]: Override Signals keyword table governance gap — the "可维护性" NFR acknowledges Pipeline Configuration table changes but not Override Signals table changes. Adding a keyword requires the same two-copy sync. Quote: "修改 SKILL.md 中的结构化表格即可" — this covers Pipeline Configuration but not the parallel Override Signals table. Governance for keyword addition/removal is undefined.
5. [Industry Benchmarking]: No intermediate determinism option explored — the proposal jumps from "natural language keywords" to "LLM interprets" without exploring structured metadata alternatives (e.g., PRD frontmatter fields like `overrides: [api-handbook]` that are deterministic and machine-parseable). This intermediate option would address the determinism gap while staying within the baseline+override architecture.
6. [Requirements Completeness]: gen-contracts/SKILL.md intent reference unaddressed — codebase grep reveals gen-contracts/SKILL.md line 38: "Parameters serve as prefill hints... do not override the user's confirmed intent." This semantically references the intent system but is not in scope or out-of-scope. Quote from constraints: "变更限于 plugins/forge/ 目录下的 8 个文件（grep intent 匹配的 skill 文件）" — the grep would match gen-contracts.
7. [Requirements Completeness]: No scenario for "refactor with no override signals" — all refactor scenarios in the proposal involve override triggers (API changes, performance). The most common case — a simple code reorganization that triggers no override signals — is assumed rather than explicitly tested. Add: "brainstorm infers refactor, PRD contains no override keywords → pipeline stays at default Spec-only, no override annotations generated."
8. [Feasibility]: Test matrix understated — 6 intents × 5 override signals creates a combinatorial test space. The "2-3 tasks" estimate doesn't account for verifying the override behavior across the full matrix. Not all combinations are meaningful (doc + any signal = no-op), but the meaningful combinations (refactor + API, refactor + performance, enhancement + security, etc.) still exceed what "2-3 tasks" implies.
9. [Scope Definition]: gen-contracts/SKILL.md not in scope or out-of-scope — contains semantically relevant intent reference. Either add to out-of-scope with justification, or add to in-scope if the intent reference needs updating.
10. [Risk Assessment]: Risk #1 mitigation understates expansion cost — Quote: "如未来出现新 task type，可追加 intent 值——Pipeline Configuration 表新增一行即可". Reality: adding an intent requires updates to 8 scope files, 2 Pipeline Configuration tables, brainstorm AskUserQuestion, breakdown-tasks Type Assignment, quick-tasks mapping, and override signal evaluation for the new intent. "新增一行即可" is misleading.
11. [Success Criteria]: SC #5 output specificity insufficient — "验证输出包含对应产物" doesn't define what the output artifact is. For API handbook: is it a PRD section? A checklist item? A generated file? The SC should specify expected output format per signal.
12. [Success Criteria]: No SC for invalid intent fallback — Scenario 8 describes fallback behavior but no SC verifies it. If the fallback is a documented feature, it needs a success criterion; if it's emergent behavior, it should be explicitly labeled as such and excluded from SC.
13. [Logical Consistency]: Scenario 8 fallback is both a feature and a non-goal — Quote: "这是 best-effort 容错，不保证正确性——根本解决方案依赖 brainstorm SKILL.md 中对输出格式的明确约束". This is simultaneously promising a fallback AND disclaiming its reliability AND identifying a root fix that's not scoped. The paragraph should commit to one position: either scope the root fix, or remove the fallback promise and document it as a known limitation.
