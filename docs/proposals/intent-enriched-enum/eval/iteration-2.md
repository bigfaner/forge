# Eval Report: Intent Enriched Enum — Iteration 2

**Reviewer**: CTO Adversary
**Date**: 2026-05-31
**Mode**: Adversarial Re-Evaluation (targeting remaining weaknesses from Iteration 1)

---

## Iteration 1 Issue Tracker

| # | Attack | Status | Evidence |
|---|--------|--------|----------|
| 1 | Urgency lacks frequency data | PARTIALLY FIXED | "已有 8 个 task type 中 5 个...无对应 intent 值" added, but still no count of actual heuristic failures in practice |
| 2 | Vague "pipeline 分支过粗" | PARTIALLY FIXED | Problem section now says "5 个被压缩进同一条分支" — quantified |
| 3 | Override UX: silent or visible? | FIXED | Added `<!-- Override: ... -->` annotation line and user review mechanism |
| 4 | Keyword negation handling | FIXED | Added explicit negation handling paragraph with LLM context reasoning |
| 5 | Industry examples too shallow | FIXED | Added GitHub 2013/2016/2020 timeline, TypeScript 2.0 specifics, behavioral differentiation pattern |
| 6 | "CI lint gate" unexplained | FIXED | Full paragraph explaining ESLint overrides and GitHub Actions path-based triggers as analog |
| 7 | "内容驱动" rejection contradicts hybrid | NOT FIXED | Still says "不提供稳定基线...同一 PRD 可能产出不同配置" — rejection rationale now slightly more honest ("无基线意味着不可复现，但本方案保留了 LLM 覆盖能力") but the core contradiction remains: the chosen approach ALSO depends on LLM judgment for override signals |
| 8 | "只扩枚举" straw man | PARTIALLY FIXED | Con column expanded to "pipeline 分支仍只有 2 条，5 个 intent 挤在 Spec-only 分支里" — more specific, but still presented as obviously inferior without genuine analysis of why someone might accept this tradeoff |
| 9 | Multi-signal scenario missing | FIXED | Scenario 6 added with explicit independent stacking rule |
| 10 | Intent-content mismatch | FIXED | Scenario 7 added with no-op behavior for doc intent |
| 11 | Missing maintainability NFR | FIXED | Added "可维护性" NFR with explicit mention of two-copy sync burden |
| 12 | task-doc.md scope error | FIXED | Scope now lists 8 files, task-doc.md removed. Constraints section explicitly notes "排除 task-doc.md 的误匹配" |
| 13 | Missing scope: breakdown-tasks Type Assignment | FIXED | Scope item now says "更新 Type Assignment 表中 coding.fix 的约束描述" |
| 14 | "6 值仍不够" risk underrated L/L | FIXED | Re-rated to M/L with honest justification |
| 15 | Missing risk: keyword false-positives | FIXED | Added as explicit risk row (关键词误触发) |
| 16 | Missing risk: two-copy sync drift | FIXED | Added as explicit risk row (Pipeline Configuration 表同步漂移) |
| 17 | Missing SC for enhancement PRD format | FIXED | SC #4 now explicitly verifies "Simplified PRD 格式（Background + Goals + Test Pipeline）" |
| 18 | "9 个文件无遗漏" untestable | FIXED | Changed to specific grep verification command (SC #8) |
| 19 | Missing backward compatibility SC | FIXED | SC #7 now verifies "现有 new-feature、refactor、cleanup 值的 pipeline 行为不变" |
| 20 | refactor/cleanup differentiation gap | PARTIALLY FIXED | Assumptions Challenged table now addresses this explicitly: "两者的 override 概率不同...区分 intent 让后续度量成为可能" — honest admission that default pipeline is identical, value is in measurement and override probability |
| 21 | enhancement mapping ambiguity | FIXED | Architecture Decision section explicitly states "enhancement → coding.enhancement（不再是 coding.feature）" |
| 22 | doc umbrella vs 1:1 mapping contradiction | FIXED | Explicit parenthetical: "breakdown-tasks 的 Intent Propagation 将 doc intent 解析为 doc task type，不区分子类型" |

**Summary**: 22 issues raised, 11 fully fixed, 9 partially fixed, 2 not fixed. The proposal has been substantially strengthened since Iteration 1. Remaining attacks focus on partially-fixed and unfixed items, plus new weaknesses introduced by revisions.

---

## Phase 1: Reasoning Audit

**Problem → Solution → Evidence → SC chain**:
- Problem: 3-value intent → 5/8 task types unmapped → pipeline squashes 5 types into 1 branch
- Solution: 6-value enum + hybrid pipeline (intent baseline + content override) + keyword-based signals
- Evidence: Codebase grep confirms 5 unmapped types, heuristic at brainstorm/SKILL.md:96, identical treatment in write-prd/tech-design
- SC: 8 criteria, most verifiable

**Chain integrity**: The chain remains sound. Revisions strengthened weak links (multi-signal, negation, mapping gaps). One new gap introduced: the override annotation mechanism creates an implicit user review step that isn't acknowledged as a workflow change.

---

## Phase 2: Rubric Scoring

### 1. Problem Definition (88/110)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Problem stated clearly | 38/40 | "8 个 task type 中有 5 个被压缩进同一条分支" — now quantified. The core problem is unambiguous. |
| Evidence provided | 32/40 | Four concrete evidence items, two verified against codebase. Still lacks frequency data — how many pipeline runs actually produced wrong artifacts? The evidence proves the gap exists, not that it causes pain frequently. |
| Urgency justified | 18/30 | Improved: "已有 8 个 task type 中 5 个无对应 intent 值" quantifies the gap. "每次 heuristic miss 都产生 pipeline 错配" still lacks count. The argument remains "it'll get worse" without data on current pain frequency. |

**Attacks on revised regions**:
- Urgency now quantifies the gap (5/8) but still conflates gap existence with pain frequency. Having unmapped types is a completeness problem; producing wrong artifacts is a pain problem. The proposal conflates the two without evidence that the 5 unmapped types frequently produce wrong artifacts.

### 2. Solution Clarity (108/120)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Approach is concrete | 40/40 | 6-value enum with explicit mapping, Pipeline Configuration table, 5 override signals, negation handling, multi-signal stacking — all concrete. |
| User-facing behavior described | 40/45 | Override annotation (`<!-- Override: ... -->`) gives user visibility. Enhancement gets simplified PRD format described. Gap: the `doc` Minimal PRD format is defined as "标题 + 目标 + scope" but what does the user actually see? A 3-section document? A single paragraph? The format is underspecified relative to the other formats. |
| Technical direction clear | 28/35 | Markdown editing approach clear. LLM-based negation handling explained. Gap: the override detection "在 Pipeline Configuration 应用后，扫描 PRD 正文段落" — but write-prd generates the PRD, so it's scanning content it hasn't generated yet? This is a temporal ordering issue. The proposal seems to mean the LLM applies overrides DURING PRD generation, not after. The wording is misleading. |

**Attacks on revised regions**:
- The negation handling paragraph is well-reasoned ("依赖 LLM 的上下文理解能力...这是合理的，因为 Pipeline Configuration 步骤本身就是 LLM 执行的") — this is a strong argument that closes the keyword-fragility concern. However, it introduces a new issue: if the LLM is doing context-aware signal detection, the structured conditional table is not the actual mechanism — the LLM's interpretation is. The table is a suggestion, not a rule. This is fine, but the proposal should be honest about this.

### 3. Industry Benchmarking (82/120)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Industry solutions referenced | 30/40 | Much improved: GitHub label timeline (2013/2016/2020), TypeScript 2.0 diagnostic categories with specific behavioral differences (Suggestion vs Error). The "CI lint gate 模式" paragraph now cites ESLint overrides and GitHub Actions path-based triggers. However: no citation of actual pipeline routing systems (e.g., Buildkite step overrides, CircleCI dynamic config, Tekton trigger bindings) which are the closest industry analogs to intent→pipeline routing. |
| At least 3 meaningful alternatives | 22/30 | 4 alternatives. "只扩枚举" improved but still framed to be rejected — the con is technically accurate ("5 个 intent 挤在 Spec-only 分支") but there's no analysis of when this tradeoff might be acceptable (e.g., small team, low pipeline diversity needs). "完全内容驱动" improved with "无基线意味着不可复现" — more honest framing. |
| Honest trade-off comparison | 15/25 | Improved: the cons column for the selected approach now mentions the two-copy table sync. However: still no analysis of the maintenance cost of the keyword table. The keyword list is effectively a growing ruleset — what's the governance for adding/removing keywords? |
| Chosen approach justified | 15/25 | "CI lint gate 模式" now properly explained. The analogy is valid (baseline + conditional overrides, only-add-no-remove). However, the analogy breaks down in one key way: CI lint gates use machine-parseable conditions (glob patterns, file paths), while this system uses LLM-interpreted natural language. The proposal doesn't acknowledge this difference. |

**Attacks on revised regions**:
- The CI lint gate analogy is the strongest addition, but it creates a false sense of determinism. ESLint overrides match file paths — deterministic, testable, version-controlled. Override signals match PRD content via LLM interpretation — probabilistic, untestable in CI, no version control on the interpretation. The analogy papered over this fundamental difference.

### 4. Requirements Completeness (88/110)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Scenario coverage | 35/40 | 7 scenarios now including multi-signal and intent-content mismatch. Gap: no scenario for the Architecture Decision — what happens when a user selects `fix` intent in brainstorm but the bug turns out to require API changes? The fix pipeline defaults to spec-only, but the override mechanism should catch "API" — is this tested? Also: no error scenario for when the LLM hallucinates an intent not in the 6-value enum. |
| Non-functional requirements | 30/40 | Added maintainability NFR with explicit two-copy acknowledgment. Gap: no testability NFR — how do you verify the override signals work correctly? No mention of test strategy for the pipeline table changes. |
| Constraints & dependencies | 23/30 | Good: 8 files, no Go code, task-doc.md false match explicitly excluded. Gap: the constraint "grep intent 匹配的 skill 文件" still doesn't prove the list is exhaustive — what if a skill file uses "intent" in a different semantic context (e.g., "design intent")? |

**Attacks on revised regions**:
- Scenario 7 (intent-content mismatch with doc) is a clever edge case. But it reveals an assumption: that override signals are no-op for doc intent because "doc intent 的 pipeline 没有可被'开启'的检查项". What if a future iteration adds checks to the doc pipeline? The no-op behavior becomes a bug. This should be documented as an explicit design decision, not an incidental property.

### 5. Solution Creativity (38/100)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Novelty over industry baseline | 15/40 | Still honest: "无特别创新". The hybrid baseline+override pattern is standard in routing/middleware systems. The Assumptions Challenged section adds depth but doesn't make the solution itself more creative. |
| Cross-domain inspiration | 10/35 | Improved: CI lint gate analogy is cross-domain. Still no borrowing from adjacent domains like feature flag systems (progressive rollout → progressive pipeline activation), test impact analysis (coverage-based test selection → content-based pipeline selection), or A/B testing frameworks (traffic splitting → pipeline splitting). |
| Simplicity of insight | 13/25 | The Assumptions Challenged section shows deeper thinking, particularly the "refactor/cleanup differentiation" stress test. But the core insight remains obvious: expand the enum to match the task types. |

### 6. Feasibility (88/100)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Technical feasibility | 38/40 | All markdown, no Go code. Verified against codebase. The LLM-dependent negation handling is a quality concern, not a feasibility blocker. |
| Resource & timeline | 26/30 | 8 files, 2-3 tasks. Write-prd and tech-design are the heaviest changes. The estimate is reasonable but doesn't account for the override annotation mechanism — this is a new behavior that needs testing across all 6 intents × 5 signals. |
| Dependency readiness | 24/30 | No external dependencies. LLM instruction-following dependency acknowledged through the structured table argument. Gap: the override annotation (`<!-- Override: ... -->`) depends on downstream skills/consumers understanding and respecting these annotations. What consumes these annotations? Is there a verification step? |

### 7. Scope Definition (72/80)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| In-scope items are concrete | 28/30 | 8 specific files with specific changes. task-doc.md removed. breakdown-tasks Type Assignment update added. Improvement over Iteration 1. |
| Out-of-scope explicitly listed | 22/25 | 4 items. Gap: eval skills not mentioned. The proposal changes pipeline behavior — do eval rubrics or contracts that reference intent need updating? Also: the `<!-- Override: ... -->` annotation format is a new protocol — is any consumer of PRD/tech-design output in scope to handle it? |
| Scope is bounded | 22/25 | "2-3 tasks" is bounded. The two-copy sync is now acknowledged as a risk. But: the override signal keyword table is an open-ended maintenance surface (who adds new signal types?) and this isn't bounded. |

### 8. Risk Assessment (72/90)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Risks identified | 25/30 | 6 risks now (up from 4). Added keyword false-positives, two-copy sync drift, branch rewrite inconsistency. Gap: no risk for LLM hallucinating an intent not in the 6-value enum (e.g., a brainstorm that outputs `bug-fix` instead of `fix`). |
| Likelihood + impact rated | 22/30 | "6 值枚举仍不够" re-rated to M/L — honest. However: "关键词误触发" rated M/L with mitigation "依赖 LLM 上下文理解" — but if the LLM is already doing context understanding, why is the likelihood M? Either the LLM handles negation well (likelihood should be L) or it doesn't (mitigation is weak). The rating and mitigation contradict. |
| Mitigations are actionable | 25/30 | Improved: two-copy sync mitigation now has "diff 检查成本极低" as argument. Override annotation enables user review. But: risk #5 ("分支重写引入不一致") mitigation is "Pipeline Configuration 表统一两处逻辑，减少不一致可能性" — this is a description of the solution, not a mitigation for the risk. The risk is that the REWRITE introduces bugs; the mitigation should be a testing strategy, not a restatement of the approach. |

**Attacks on revised regions**:
- Risk #3 (关键词误触发) is the most interesting addition. The mitigation says "LLM 可识别否定语境" and "最坏情况是开启了一个不必要的检查". This is correct — the override is additive-only, so false positives are safe. But then why is likelihood M? If the worst case is benign, shouldn't the impact be L (which it is) and the likelihood irrelevant? The risk assessment is internally inconsistent in its urgency signaling.

### 9. Success Criteria (72/80)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Criteria are measurable and testable | 26/30 | SC #8 now uses a specific grep command — excellent. SC #4 specifies enhancement PRD format sections. SC #7 has specific behavioral invariants. Gap: SC #5 "可被 PRD 内容触发" — still vague. Which specific PRD content triggers which specific signal? Need input→output test pairs. |
| Coverage is complete | 22/25 | Improved: enhancement PRD format covered (SC #4), backward compatibility covered (SC #7). Gap: no SC for the Architecture Decision — "fix intent 映射到 coding.fix" is in scope but not in SC. Also: no SC for the override annotation format — the `<!-- Override: ... -->` protocol is part of the solution but has no verification criterion. |
| SC internal consistency | 24/25 | SC set is internally consistent. The grep verification (SC #8) combined with behavioral invariants (SC #7) provides good coverage. One concern: SC #8 grep command checks specific directories — if intent references exist in other skill directories not listed, the SC would pass despite incomplete coverage. |

**Attacks on revised regions**:
- SC #7 "现有 new-feature、refactor、cleanup 值的 pipeline 行为不变（旧 3 值在 Pipeline Configuration 表中对应行与当前行为一致）" — the parenthetical redefines the SC from "behavioral test" to "table row comparison". These are not equivalent. The old behavior might not be perfectly captured by the new table rows. The SC should verify actual output, not table structure.

### 10. Logical Consistency (78/90)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Solution addresses stated problem | 30/35 | Strong alignment. The refactor/cleanup differentiation is now honestly addressed in Assumptions Challenged: value is in measurement and override probability, not default pipeline difference. One remaining gap: the problem says "pipeline 分支过粗" with 5 types in 1 branch, but the solution gives refactor/cleanup the SAME default pipeline (Spec-only). So 6 intents produce only 4 distinct pipeline configurations (Full, Simplified, Spec-only × 3, Minimal). The "refinement" is less dramatic than claimed — 2 effective pipelines become 4, not 2 become 6. |
| Scope ↔ Solution ↔ SC aligned | 24/30 | Mostly aligned. Gap: the override annotation protocol is in the solution but not in scope (no file is listed as needing annotation-related changes) and not in SC. Also: the Architecture Decision's Type Assignment table update is in scope but the mapping change (coding.feature/coding.enhancement split) isn't reflected in the brainstorm scope item — brainstorm/SKILL.md line 93 currently maps both to `new-feature`, and the scope says "更新 Step 4.5 intent mapping 表（6 值）" which should cover it, but doesn't explicitly call out the split. |
| Requirements ↔ Solution coherent | 24/25 | Good coherence. The NFR maintainability section now honestly acknowledges the two-copy burden. The Architecture Decision resolves the fix mapping conflict. One gap: the NFR says "向后兼容：现有的 new-feature、refactor、cleanup 值行为不变" but the brainstorm AskUserQuestion changes from 3 to 6 options — this is a behavioral change for users running brainstorm. Not a pipeline behavior change, but a UX change. The NFR should distinguish pipeline backward compatibility from interaction backward compatibility. |

---

## Phase 3: Blindspot Hunt

### New Issues Introduced by Revisions

1. **Override annotation has no consumer**: The `<!-- Override: ... -->` annotation is written into PRD/tech-design output, but no downstream skill or process is specified to read, validate, or act on these annotations. They exist for "user review" but users have no structured workflow for reviewing them. This is a documentation feature with no integration.

2. **Temporal ordering confusion in override detection**: "在 Pipeline Configuration 应用后，扫描 PRD 正文段落中是否命中关键词" — the PRD is being generated by write-prd. The LLM is simultaneously generating content and detecting signals in that content. The override isn't applied "after" — it's applied "during". The sequential framing is misleading about the actual LLM execution model.

3. **Pipeline configuration count inflation**: The proposal implies 6 distinct pipeline configurations, but refactor/cleanup/fix all default to Spec-only PRD. The actual distinct configurations are: Full (new-feature), Simplified (enhancement), Spec-only (refactor/cleanup/fix), Minimal (doc) = 4 configurations. The expansion from 2 to 4 effective pipelines is meaningful but less than the "6 行" framing suggests.

4. **Assumptions Challenged section is well-crafted but buried**: The "refactor/cleanup" stress test is the most honest piece of reasoning in the document. It should be in the Solution section, not hidden in a section that most reviewers might skip.

### Persistent Weaknesses from Iteration 1

5. **"完全内容驱动" rejection still partially contradictory**: The revision adds "无基线意味着不可复现，但本方案保留了 LLM 覆盖能力" — acknowledging the contradiction partially. But the core issue remains: the chosen approach's override mechanism relies on the same LLM judgment that "完全内容驱动" was rejected for. The difference is the presence of a baseline, which is a valid distinction, but the rejection should be framed as "no baseline" not "LLM judgment unreliable".

6. **Frequency data still missing**: The proposal proves the mapping gap exists (5/8 types unmapped) but not that it causes frequent pain. One concrete example (API handbook skip) is provided, but no data on how often heuristic misses produce wrong artifacts. This is the difference between a completeness problem and a pain problem.

---

## Deductions

- **Vague language**: "pipeline 错配" in urgency — used without specific examples of what wrong artifacts were produced. -10 pts from Problem Definition (reduced from -20 in Iteration 1 because the quantification improvement partially offsets).
- **Straw-man alternative**: "只扩枚举" con column improved but still frames it as obviously wrong without genuine analysis of when the tradeoff might be acceptable. -10 pts from Industry Benchmarking (reduced from -20 because the con is now technically accurate).

---

SCORE: 706/1000
DIMENSIONS:
  Problem Definition: 88/110
  Solution Clarity: 108/120
  Industry Benchmarking: 82/120
  Requirements Completeness: 88/110
  Solution Creativity: 38/100
  Feasibility: 88/100
  Scope Definition: 72/80
  Risk Assessment: 72/90
  Success Criteria: 72/80
  Logical Consistency: 78/90
ATTACKS:
1. [Problem Definition]: Frequency data still missing — "每次 heuristic miss 都产生 pipeline 错配" is a claim about frequency with no frequency data. One concrete example (API handbook skip) proves the gap exists, not that it hurts frequently. Provide: count of brainstorm runs that used the heuristic in the last N sessions, or count of PRDs where wrong pipeline output was manually corrected.
2. [Solution Clarity]: Override detection temporal ordering is misleading — "在 Pipeline Configuration 应用后，扫描 PRD 正文段落" implies sequential execution, but the LLM generates PRD content and detects signals simultaneously. The sequential framing obscures the actual execution model. Reword to describe the LLM's simultaneous generation+detection process.
3. [Solution Clarity]: `doc` Minimal PRD format is underspecified — "标题 + 目标 + scope" describes 3 section names but not what goes in them. Compare to the enhancement format which explicitly lists "Background（说明增强什么）、Goals（增强目标）、Test Pipeline（确保增强有测试覆盖）". Doc format needs equivalent specificity.
4. [Industry Benchmarking]: CI lint gate analogy breaks down on determinism — ESLint overrides use machine-parseable glob patterns (deterministic, testable in CI, version-controlled). Override signals use LLM-interpreted natural language (probabilistic, untestable in CI, interpretation not version-controlled). The analogy papered over this fundamental difference. Acknowledge the determinism gap.
5. [Industry Benchmarking]: "完全内容驱动" rejection rationale contradicts chosen approach — rejected for "不提供稳定基线...同一 PRD 可能产出不同配置" but the justification should be "no baseline", not "unreliable LLM judgment", since the chosen approach also depends on LLM for override signal interpretation. The revision partially addresses this but the framing still implies LLM judgment is the problem, when it's the lack of baseline.
6. [Requirements Completeness]: No testability NFR — how do you verify override signals work correctly? No mention of test strategy. The SC says "可被 PRD 内容触发" but doesn't specify test pairs. Add: "For each override signal, at least one specific PRD input→pipeline output test case must be defined."
7. [Requirements Completeness]: No scenario for LLM hallucinating an invalid intent — brainstorm could output `bug-fix` instead of `fix`, or `documentation` instead of `doc`. The 6-value enum has no fallback or validation described. Add: "What happens when brainstorm outputs an intent not in the 6-value set?"
8. [Scope Definition]: Override annotation protocol has no owner — `<!-- Override: ... -->` is a new output format. Which file/skill change implements the annotation writing? Not listed in scope. Which downstream process consumes it? Not in scope or out-of-scope. The annotation exists in limbo.
9. [Scope Definition]: Eval skills not mentioned in scope or out-of-scope — the proposal changes pipeline behavior for write-prd and tech-design. Do eval rubrics, contracts, or journey tests that reference intent need updating? Either add to in-scope or explicitly add to out-of-scope.
10. [Risk Assessment]: Risk #3 likelihood/mitigation contradiction — "关键词误触发" rated M likelihood with mitigation "LLM 可识别否定语境". If the LLM handles negation well enough to be a mitigation, why is likelihood M? Either the LLM handles it (likelihood L, mitigation validated) or it doesn't (mitigation weak). Resolve the contradiction.
11. [Risk Assessment]: Risk #5 mitigation is circular — "Pipeline Configuration 表统一两处逻辑，减少不一致可能性" describes the solution design, not a mitigation for the risk of introducing bugs during the rewrite. Add: specific testing strategy (e.g., "run existing write-prd and tech-design test contracts before and after change, verify output parity for old 3 values").
12. [Success Criteria]: SC #7 conflates table structure with behavior — "旧 3 值在 Pipeline Configuration 表中对应行与当前行为一致" verifies table rows, not actual output. A table row that looks correct can still produce different output if the LLM interprets it differently. Replace with: "For proposals with old 3 intents, the generated PRD/tech-design output matches pre-change output."
13. [Success Criteria]: No SC for Architecture Decision — "fix intent 映射到 coding.fix" and "coding.feature/coding.enhancement split" are key decisions in scope but have no success criterion. Add: "breakdown-tasks maps fix intent to coding.fix task type" and "brainstorm maps coding.feature → new-feature and coding.enhancement → enhancement as separate paths."
14. [Logical Consistency]: Effective pipeline count is 4, not 6 — refactor/cleanup/fix all default to Spec-only PRD. The "6 行" framing implies 6 distinct configurations, but only 4 are functionally different. This is a framing gap, not a logical error, but it inflates the perceived impact of the change.
15. [Logical Consistency]: brainstorm AskUserQuestion UX change not acknowledged as breaking — NFR says "向后兼容" but brainstorm goes from 3 to 6 options. Pipeline backward compatibility holds, but interaction backward compatibility changes. Distinguish: "Pipeline output backward compatible; brainstorm interaction changed (3→6 options)."
16. [Logical Consistency]: Override annotation no-op for doc intent is incidental, not designed — Scenario 7 says "override 信号为 no-op，不改变 pipeline 行为" because doc pipeline has no check items to override. This is a coincidental property, not a design guarantee. If future doc pipeline adds checks, the no-op becomes a bug. Document this as an explicit design constraint, not an incidental observation.
