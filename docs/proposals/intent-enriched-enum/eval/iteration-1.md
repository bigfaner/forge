# Eval Report: Intent Enriched Enum — Iteration 1

**Reviewer**: CTO Adversary
**Date**: 2026-05-31
**Mode**: Annotated Blind Review (Pre-Revised markers present)

---

## Phase 1: Reasoning Audit

**Problem → Solution chain**:
- Problem: 3-value intent enum → poor mapping to 8-value task type → pipeline branching too coarse
- Solution: Expand to 6 values + hybrid pipeline (intent baseline + content override signals)
- Evidence: Concrete — grep confirms `coding.fix` heuristic exists in brainstorm/SKILL.md line 96, refactor/cleanup indeed treated identically in write-prd and tech-design
- SC: 7 criteria, mostly verifiable

**Chain integrity**: The problem-solution chain is sound. The core thesis (enum expansion + hybrid pipeline) directly addresses both stated motivations.

---

## Phase 2: Rubric Scoring

### 1. Problem Definition (82/110)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Problem stated clearly | 35/40 | Core problem unambiguous: 3-value intent can't map to 8-value task type. One interpretation issue: "pipeline 分支过粗" is asserted but the exact consequence (wrong artifacts produced) is only illustrated by one example (API handbook skip). |
| Evidence provided | 30/40 | Four concrete evidence items, two verified against codebase (coding.fix heuristic at brainstorm/SKILL.md:96, refactor/cleanup identical treatment in write-prd self-check.md). However, no frequency data — how often does this mismatch actually cause problems? One concrete case does not prove systemic pain. |
| Urgency justified | 17/30 | "随着 Forge 处理的场景增多" is trend-based justification but lacks quantification. No deadline pressure. No cost-of-delay metric. The argument is "it'll get worse" — a valid but weak urgency case. |

**Attacks on annotated regions**:
- `<!-- pre-revised: high -->` in Problem section: No pre-revised markers in the Problem section. Clean region.

**Attacks on unannotated regions**: The urgency section is the weakest part of the entire proposal. "每次遇到非标准场景都需要 LLM 做启发式判断" is a claim about frequency with zero frequency data.

### 2. Solution Clarity (98/120)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Approach is concrete | 38/40 | Four-part solution with a concrete enum (6 values named), a concrete pipeline table (6 rows × 6 columns), and concrete override signals (5 signal types with keywords). Reader can explain back exactly what will be built. |
| User-facing behavior described | 35/45 | User sees 6 options instead of 3 in AskUserQuestion. PRD/tech-design produce different artifacts based on intent + content signals. But: what does the user *experience* when an override signal fires? Is there a prompt? A silent switch? The proposal doesn't say. |
| Technical direction clear | 25/35 | "Markdown editing" is the technical approach — clear but shallow. The override signal detection mechanism ("扫描 PRD 正文段落中是否命中关键词") is described as keyword matching, which is fragile. No discussion of how to handle multi-language keywords (the codebase uses Chinese terms like "接口变更" alongside English "API"). |

**Attacks on annotated regions**:
- `<!-- pre-revised: high -->` on point 1: The 6-value enum justification now includes explicit rationale for doc.consolidate/doc.drift exclusion. Good revision — eliminates the "why not 8?" question.
- `<!-- pre-revised: medium -->` on point 2: The Override Signals table is well-structured. Revision strengthened the section by replacing prose with a conditional table. No new issues introduced.

**Attacks on unannotated regions**: The keyword-based override detection is the weakest technical element. "API" as a keyword will false-positive on any PRD that mentions "this does NOT change the API". The proposal has no negation handling. Also: "命中任意一个信号即触发对应覆盖" means a single keyword hit triggers the override — no confidence threshold, no context window.

### 3. Industry Benchmarking (58/120)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Industry solutions referenced | 18/40 | Two examples (GitHub Issues labels, TypeScript diagnostic categories) cited but analyzed at hand-wave level. No specifics: which GitHub label evolution? What TypeScript version added what category? The comparison is "they also grew categories" — trivially true for any taxonomy. No reference to pipeline routing systems (e.g., CI/CD gating, feature flag management, test classification frameworks) which are the actual industry analogs of intent-based pipeline branching. |
| At least 3 meaningful alternatives | 22/30 | 4 alternatives listed including "do nothing". However, "只扩枚举" is a straw man — it's presented only to be rejected. The "完全内容驱动 pipeline" alternative is poorly characterized ("依赖 LLM 判断力，不稳定") when the chosen hybrid approach ALSO relies on LLM judgment for override signal detection. |
| Honest trade-off comparison | 8/25 | Cons column for selected approach says only "9 个文件变更" — this is a cost, not a con. No trade-off analysis of the keyword-based override fragility. No discussion of maintenance burden for the pipeline table. |
| Chosen approach justified | 10/25 | "CI lint gate 模式" is cited as Source but never explained. What CI lint gate? What specific pattern? The justification is "双重改进" which is circular — the approach is good because it solves both problems, but that's the definition of the approach, not a justification against benchmarks. |

**Attacks on annotated regions**:
- `<!-- pre-revised: medium -->` on Industry Solutions: The revision added the GitHub/TypeScript examples. The examples are relevant but shallow. The common pattern identified ("分类粒度必须匹配行为差异") is a tautology, not an insight.

### 4. Requirements Completeness (72/110)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Scenario coverage | 28/40 | 5 key scenarios listed. Missing: (1) What happens when multiple override signals fire simultaneously? (2) What happens when intent is `doc` but PRD content contains "API" keyword? (3) Error scenario: brainstorm infers `fix` but user overrides to `new-feature` — does the override signal table still apply? These gaps matter because the override mechanism is the novel part of the solution. |
| Non-functional requirements | 22/40 | Backward compatibility mentioned. Consistency mentioned. Missing: (1) Performance — keyword scanning cost on long PRDs? (2) Maintainability — who updates the keyword list when new terms emerge? (3) Testability — how do you test that the override signals fire correctly? (4) Internationalization — keywords are mixed Chinese/English, what about future languages? |
| Constraints & dependencies | 22/30 | Good: no Go code dependency, limited to 9 files. However, the "9 files" claim needs verification (see Logical Consistency section). Missing constraint: the proposal assumes LLM compliance with the pipeline table — what if the LLM hallucinates a different format? |

**Attacks on annotated regions**:
- `<!-- pre-revised: high -->` on Architecture Decision: The fix mapping strategy is well-reasoned. The distinction between "fix intent via brainstorm" vs "coding.fix via CLI" is clear. Good addition that addresses a real conflict. However, it's buried in Requirements Analysis rather than being in Solution — this is a solution-level decision, not a requirement.

### 5. Solution Creativity (35/100)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Novelty over industry baseline | 15/40 | The proposal itself says "无特别创新". The hybrid intent+signals approach is a straightforward two-layer dispatch pattern — common in routing systems, middleware chains, and rule engines. Not novel, but honest about it. |
| Cross-domain inspiration | 5/35 | No cross-domain borrowing. The "CI lint gate 模式" is mentioned as inspiration source but never elaborated. No reference to rule engines, middleware pipelines, or feature flag systems which are the actual cross-domain analogs. |
| Simplicity of insight | 15/25 | The core insight ("expand enum + add content-based overrides") is simple and appropriate. Not elegant enough for "why didn't I think of that" — it's the obvious solution once you identify the problem. |

**Attacks on unannotated regions**: The Innovation Highlights section is admirably honest ("无特别创新") but then claims "混合模式 pipeline 是对当前二元分支的自然细化" — "natural refinement" is not an innovation highlight, it's a description.

### 6. Feasibility (85/100)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Technical feasibility | 38/40 | Verified: all changes are markdown files. brainstorm/SKILL.md, write-prd/SKILL.md, tech-design/SKILL.md all confirmed to use intent branching. No Go code involved. The only risk is the keyword-based override detection being fragile, but that's a quality concern, not a feasibility blocker. |
| Resource & timeline | 25/30 | "9 markdown files, 2-3 tasks" is reasonable. However, write-prd and tech-design need full pipeline rewrites (the current binary branching is deeply embedded in step-by-step instructions). The estimate may undercount the effort for these two files. |
| Dependency readiness | 22/30 | No external dependencies. But: the proposal depends on the LLM correctly following a 6-row pipeline table + 5 override signal rules + consistent keyword matching. This is a dependency on LLM instruction-following reliability, which is not an external API but is a real dependency. |

### 7. Scope Definition (62/80)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| In-scope items are concrete | 24/30 | 9 specific files named with specific changes described. However, "更新 AskUserQuestion 选项" is vague — what exactly are the 6 options and their prompts? The quick-tasks/templates/task-doc.md entry claims "更新 intent 有效值注释" but the file contains NO intent field at all (verified: task-doc.md has no intent reference). This is a scope item that doesn't match reality. |
| Out-of-scope explicitly listed | 20/25 | 4 items listed. Good. Missing: what about eval skills? The proposal changes pipeline behavior but doesn't mention updating eval rubrics or contracts that reference intent. |
| Scope is bounded | 18/25 | "2-3 tasks" is bounded. But the override signal table introduces an open-ended maintenance surface — who adds new signal types? This isn't bounded. |

**Attacks on annotated regions**:
- `<!-- pre-revised: medium -->` on quick-tasks/templates/task-doc.md: The claim is "更新 intent 有效值注释（grep 扫描发现的遗漏文件）". Verified: this file has NO intent field and NO intent-related comments. The scope item is based on a false grep finding (likely `grep -r "intent"` matched some other content or the search was imprecise). This is a factual error in scope.

### 8. Risk Assessment (60/90)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Risks identified | 22/30 | 4 risks identified. Missing: (1) Keyword false-positive risk (mentioned in my solution clarity analysis) — "API" matching in negation contexts. (2) Maintenance drift risk — the pipeline table in write-prd and tech-design must stay synchronized manually, a classic two-copy problem. (3) Override signal language coverage — the keywords are bilingual, future contributors may add terms in only one language. |
| Likelihood + impact rated | 18/30 | The ratings are reasonable but not deeply justified. "6 值枚举仍不够" is rated L/L — but the proposal already acknowledges doc.consolidate and doc.drift are excluded, which means the enum doesn't cover all task types. This is a current gap, not a future risk. Should be M/L. |
| Mitigations are actionable | 20/30 | The mitigations for risk #2 (override signals ignored by LLM) is well-reasoned — "结构化条件表" and "只开启不关闭" are good arguments. But risk #3 mitigation ("Pipeline Configuration 表统一两处逻辑") doesn't address the root cause — it just claims the new approach is better. And risk #4 mitigation ("旧 3 值在表中仍有对应行") is a property of the solution, not an actionable mitigation step. |

**Attacks on annotated regions**:
- `<!-- pre-revised: medium -->` on risk #2 mitigation: The revision improved the argument significantly by adding the "只开启不关闭" safety property and the "表格每行是原子化的 if-then 规则" justification. This is the strongest piece of reasoning in the entire proposal. No new issues introduced by revision.

### 9. Success Criteria (60/80)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Criteria are measurable and testable | 22/30 | Most criteria are verifiable: "6 值之一" is countable, "不再使用启发式" is grep-verifiable, "Pipeline Configuration 表（6 行）" is countable. However, "Override Signals 规则存在且可被 PRD 内容触发" is weak — "可被触发" is not specific about WHICH content triggers WHICH signal. A test would need specific input-output pairs. |
| Coverage is complete | 18/25 | Missing: (1) No SC for the enhancement simplified PRD format — the pipeline table describes a custom format but no SC verifies it. (2) No SC for backward compatibility of existing proposals — "行为不变" is asserted in NFR but not in SC. (3) No SC for the fix mapping strategy (Architecture Decision section) — the key conflict resolution isn't verified. |
| SC internal consistency | 20/25 | SC set is internally consistent — no contradictions. However, SC #3 ("统一的 Pipeline Configuration 表") and SC #4 ("Override Signals 规则") together imply that both write-prd AND tech-design must implement the same table AND the same override logic. The proposal doesn't verify they share implementation — they're separate markdown files that must be kept in sync manually. |

**Attacks on unannotated regions**: SC #7 "9 个文件全部更新，无遗漏" is problematic because: (a) the task-doc.md scope item appears to be based on a false grep finding, and (b) "无遗漏" is untestable without an exhaustive list of all files that reference intent — the proposal provides a list but doesn't prove it's exhaustive.

### 10. Logical Consistency (70/90)

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Solution addresses stated problem | 28/35 | The 6-value enum directly addresses the mapping gap. The hybrid pipeline directly addresses the "pipeline 分支过粗" problem. The fix heuristic elimination directly addresses the inconsistency risk. Strong alignment. One gap: the problem states "refactor 和 cleanup 在 write-prd/tech-design 中被完全等同对待" but the solution keeps them with the same PRD format (Spec-only) for both — the only differentiation is the override signals, which are content-based, not intent-based. So refactor and cleanup ARE still treated identically by intent, just with a safety net. |
| Scope ↔ Solution ↔ SC aligned | 22/30 | Mostly aligned. Two misalignments: (1) Scope lists quick-tasks/templates/task-doc.md but this file has no intent reference to update — scope item is phantom. (2) Solution describes enhancement simplified PRD format in detail but no SC verifies it and no scope item explicitly creates/enforces it. |
| Requirements ↔ Solution coherent | 20/25 | The Requirements section identifies the fix mapping conflict and resolves it coherently. However, the NFR "向后兼容" is stated as a requirement but the solution changes brainstorm's AskUserQuestion from 3 to 6 options — is that backward compatible? Yes for existing proposals, but it changes the brainstorm interaction, which isn't acknowledged. |

---

## Phase 3: Blindspot Hunt

### Rubric Missed

1. **Override signal false-positive handling**: The keyword-based approach has no negation detection. A PRD saying "this does NOT change the API" would still trigger the API signal. The rubric doesn't have a criterion for "mechanism robustness" — this is a design flaw that falls between Solution Clarity and Risk Assessment.

2. **Two-copy synchronization burden**: The Pipeline Configuration table must be maintained identically in write-prd/SKILL.md AND tech-design/SKILL.md. This is a classic single-source-of-truth violation. The proposal doesn't address it. The rubric's Feasibility dimension doesn't cover maintenance cost.

3. **Grep-based scope verification**: The scope claims "grep 扫描发现" task-doc.md needs updating, but the file has no intent references. This suggests the grep search was either too broad or matched false positives. The scope verification methodology is unreliable.

4. **Missing `doc` task type mapping**: The proposal maps `doc` intent to... what task type? The breakdown-tasks Type Assignment table has `doc`, `doc.consolidate`, `doc.drift` as task types, but the proposal's 1:1 mapping doesn't specify which doc sub-type `doc` intent maps to. This is an ambiguity in the "严格 1:1 映射" claim.

5. **`enhancement` vs `coding.enhancement` mapping gap**: The brainstorm currently maps `coding.feature` AND `coding.enhancement` to `new-feature`. The proposal introduces `enhancement` as a separate intent but doesn't explicitly address whether `coding.feature` and `coding.enhancement` now map to different intents (`new-feature` vs `enhancement`). The pipeline table gives them different treatments, so this mapping matters.

---

## Summary Statistics

- **Annotated regions attack density**: 7 annotations reviewed, 2 substantive findings (task-doc.md scope error, industry benchmark shallowness)
- **Unannotated regions attack density**: 12 substantive attacks on unmarked content
- **conflict-with-pre-revision**: 0 instances — no disagreements with pre-revision direction

---

## Deductions

- **Vague language**: "pipeline 分支过粗" (multiple uses) — no quantification of what "粗" means. -20 pts from Problem Definition.
- **Straw-man alternative**: "只扩枚举" is presented only to be rejected with the weakest possible argument. -20 pts from Industry Benchmarking.
- **Vague language**: "自然细化" in Innovation Highlights — no quantification. Already scored low in creativity, no additional deduction.
- **Vague language**: "覆盖缺口随场景增长扩大" in urgency — no quantification. Already reflected in urgency score.

---

SCORE: 622/1000
DIMENSIONS:
  Problem Definition: 82/110
  Solution Clarity: 98/120
  Industry Benchmarking: 58/120
  Requirements Completeness: 72/110
  Solution Creativity: 35/100
  Feasibility: 85/100
  Scope Definition: 62/80
  Risk Assessment: 60/90
  Success Criteria: 60/80
  Logical Consistency: 70/90
ATTACKS:
1. [Problem Definition]: Urgency is asserted without quantification — "随着 Forge 处理的场景增多" — how many scenarios? How many heuristic failures? One concrete example does not establish systemic pain. Must provide frequency data or at minimum a count of affected pipeline runs.
2. [Problem Definition]: Vague language "pipeline 分支过粗" — no quantification of what "粗" means or how many downstream decisions are affected. Define a metric for branching granularity.
3. [Solution Clarity]: Override signal firing has no user-visible feedback — "命中任意一个信号即触发对应覆盖" — the user doesn't know an override happened. Is this silent? If so, how does the user verify the pipeline produced the right artifacts? Specify the UX of override detection.
4. [Solution Clarity]: Keyword-based detection has no negation handling — "API" matches "this does NOT change the API". The mechanism is fragile. Add negation context handling or at minimum document this as a known limitation.
5. [Industry Benchmarking]: Two industry examples (GitHub labels, TypeScript diagnostics) are analyzed at surface level — no specific version, no specific evolution, no pattern extraction beyond "they grew categories". Cite specific GitHub label changes or TypeScript diagnostic category additions with version numbers and motivations.
6. [Industry Benchmarking]: "CI lint gate 模式" is cited as source for the chosen approach but never explained — what specific CI lint gate pattern? What project? This is an unsubstantiated reference. Either explain the reference or remove it.
7. [Industry Benchmarking]: "完全内容驱动 pipeline" alternative is rejected for "依赖 LLM 判断力，不稳定" but the chosen approach ALSO relies on LLM for override signal detection — the rejection rationale applies equally to the selected approach. Provide an honest comparison.
8. [Industry Benchmarking]: "只扩枚举" is a straw-man alternative — presented with "pipeline 分支仍然过粗" as the only con, which is the problem statement restated. Provide a genuine analysis of why someone might choose this approach.
9. [Requirements Completeness]: No scenario covers multiple override signals firing simultaneously — what happens when PRD content mentions both "API" and "性能"? Both triggers fire independently? Is there a conflict resolution? Add a multi-signal scenario.
10. [Requirements Completeness]: No scenario covers intent-content mismatch — what if intent is `doc` but PRD contains "API变更" keyword? Does the override still fire? Define the interaction between intent baseline and content override for edge cases.
11. [Requirements Completeness]: Missing NFR for maintainability — who updates the keyword list when new terms emerge? Who ensures the pipeline table in write-prd and tech-design stay synchronized? Add maintainability as an NFR.
12. [Scope Definition]: quick-tasks/templates/task-doc.md is listed as needing "更新 intent 有效值注释" but the file contains NO intent field or comment. This scope item is based on a false grep finding. Remove or correct this scope item.
13. [Scope Definition]: Missing scope item for the Architecture Decision — "允许 fix intent 映射到 coding.fix" requires updating the Type Assignment table in breakdown-tasks/SKILL.md, but this isn't listed as an in-scope change to that file. Add the Type Assignment table update to scope.
14. [Risk Assessment]: "6 值枚举仍不够" is rated L/L but the proposal already acknowledges 2 excluded task types (doc.consolidate, doc.drift) — the gap exists NOW, not in the future. Should be M/L. Re-rate this risk honestly.
15. [Risk Assessment]: Missing risk: keyword false-positives in negation contexts — a PRD saying "this does NOT change the API" would trigger the API signal. This is a design flaw, not a low-probability event.
16. [Risk Assessment]: Missing risk: two-copy synchronization — the Pipeline Configuration table must be maintained identically in write-prd/SKILL.md AND tech-design/SKILL.md with no single source of truth. Divergence is a matter of when, not if.
17. [Success Criteria]: No SC for the enhancement simplified PRD format — the pipeline table defines a custom format (Background/Goals/Test Pipeline) but no SC verifies write-prd produces this format. Add a specific SC for enhancement PRD format.
18. [Success Criteria]: "9 个文件全部更新，无遗漏" is untestable — the proposal provides a list but doesn't prove it's exhaustive. "无遗漏" requires proof of completeness, not just a list. Replace with a grep-verification criterion.
19. [Success Criteria]: No SC for backward compatibility — NFR section claims "现有的 new-feature、refactor、cleanup 值行为不变" but no SC verifies this. Add backward compatibility SC with specific behavioral invariants.
20. [Logical Consistency]: Solution claims refactor and cleanup will be differentiated, but both get Spec-only PRD format and the only differentiation is content-based override signals — the INTENT itself doesn't distinguish them. The problem statement says they're "等同对待" and the solution... still treats them identically by intent, just with a safety net. This is a gap between stated problem and actual solution.
21. [Logical Consistency]: The 1:1 mapping claim is incomplete — `enhancement` intent maps to which task type? `coding.enhancement` or `coding.feature`? The current brainstorm maps BOTH to `new-feature`. The proposal must specify the new mapping explicitly.
22. [Logical Consistency]: `doc` intent maps to which task type? The task types include `doc`, `doc.consolidate`, `doc.drift` — but the proposal says `doc` umbrella. "严格 1:1 映射" contradicts the umbrella approach. Resolve this contradiction.
