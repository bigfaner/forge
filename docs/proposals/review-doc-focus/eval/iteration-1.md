---
iteration: 1
evaluator: CTO Adversary (Rubric Scoring)
date: 2026-05-24
document: docs/proposals/review-doc-focus/proposal.md
---

# Iteration-1 Evaluation: Review-Doc Execution Scope Focus

**Evaluation Method**: Annotated blind review — scored against the CTO rubric, tracing Problem -> Solution -> Evidence -> Success Criteria chain. Annotated regions (`<!-- pre-revised -->`) receive focused audit for revision-introduced issues; unannotated regions receive standard scrutiny.

---

## Phase 1: Reasoning Audit

**Chain**: Problem (agent reads irrelevant files -> waste + quality loss + mis-edits) -> Solution (build-time AC embedding + allowlist discovery) -> Evidence (prompt template quotes, execution logs) -> Success Criteria (6 checkbox items).

**Chain integrity**: The reasoning chain is fundamentally sound. Build-time AC embedding directly eliminates the need for runtime task-file scanning (addresses Problem 1 and 2). Allowlist discovery directly prevents reading and modifying non-deliverable files (addresses Problem 3). However, the chain has two stress points: (a) the "allowlist" is still a prompt instruction, partially relying on the mechanism the proposal rejects as unreliable, and (b) the success criteria measure structural changes but do not verify that review quality is maintained.

---

## Phase 2: Rubric Scoring

### 1. Problem Definition: 80/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 32/40 | Three sub-problems enumerated with specificity: token waste, quality degradation, agent mis-operation. The core problem is unambiguous. Minor deduction: "token waste" and "quality degradation" are partially overlapping (quality loss is a downstream consequence of attention dilution from token waste), making them feel like two facets of one problem rather than two independent problems. |
| Evidence provided | 27/40 | Four concrete evidence items cite specific file paths and template steps. The first three are verifiable against source code (confirmed correct on inspection). However, the fourth item — "实际执行记录显示 Referenced Documents 包含任务文件和记录文件" — remains anecdotal. No specific execution record is named, no count is given, no percentage of wasted reads is quantified. This was flagged in iteration-0 and the revision did not address it. Quote: "实际执行记录显示 Referenced Documents 包含任务文件和记录文件" — which records? How many runs? What proportion of referenced docs were irrelevant? |
| Urgency justified | 21/30 | States the pipeline is merged and active. Valid but thin. No estimate of execution frequency or per-execution waste magnitude. The urgency argument is "every execution suffers" — but how often? Daily? Weekly? And what is the per-execution cost? Without this, the cost-of-delay argument is asserted rather than demonstrated. |

**Deductions**:
- Vague language without quantification: "实际执行记录显示" (-20 pts) — still no specific records or numbers cited after revision.

### 2. Solution Clarity: 90/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 35/40 | Two-pronged approach with named implementation targets. The pre-revision added the AC aggregation schema format (lines 122-142) and prompt template restructuring spec (lines 145-155), which resolve the major gap from iteration-0. An implementer can now explain back the data flow: BuildIndex extracts AC per doc-task -> BodyContext carries it -> autogen template renders `## Acceptance Criteria Summary` with per-task subsections -> agent reads from that section. Remaining gap: the `renderBody()` mechanism for the new `AcceptanceCriteria map[string]string` field is described as `{{range .AcceptanceCriteria}}` but the actual template uses `{{ACCEPTANCE_CRITERIA}}` (a simple string placeholder). The proposal does not specify whether a new placeholder token is added or the existing one is repurposed. |
| User-facing behavior described | 35/45 | The Innovation Highlights section (lines 37-41) and the Prompt Template Restructuring Spec (lines 145-155) now describe observable agent behavior: agent reads pre-extracted AC from the task file, discovers documents only under `docs/`, and modifies only `docs/` files. However, there is still no explicit before/after comparison table showing what the agent reads and modifies in the old flow vs. the new flow. The user-facing change is implied across multiple sections rather than stated in one place. |
| Technical direction clear | 20/35 | Four implementation threads named with function-level detail (lines 78-86). However: (1) The proposal says "新增 `AcceptanceCriteria map[string]string` 字段" to BodyContext, but the existing BodyContext already has `AcceptanceCriteria []string` for PRD-level AC. This creates a naming collision that the proposal ignores. (2) The proposal says `renderBody()` uses `{{range .AcceptanceCriteria}}` — but `renderBody()` uses simple string replacement (`strings.ReplaceAll`), not Go template ranges. The actual implementation would need either a new placeholder + rendering logic, or a switch to Go templates. The technical direction is at the right granularity but points at wrong mechanisms. |

**Deductions**: None applied beyond score reduction.

### 3. Industry Benchmarking: 40/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 10/40 | Zero industry references. No citation of agent frameworks (LangGraph, CrewAI, AutoGen), RAG filtering patterns, agent sandboxing approaches, or context window management strategies. This was flagged in iteration-0 and remains unaddressed. |
| At least 3 meaningful alternatives | 15/30 | Three alternatives: selected approach, pure prompt, do-nothing. The pre-revision added the "Challenge Override" note (line 41) clarifying that the approach is "dual-layer defense" — build-time data removal + prompt-level allowlist. This partially addresses the straw-man concern by admitting prompt constraints are complementary rather than rejected. However, no genuinely different alternatives are explored: no execution-time AC extraction, no separate AC manifest file, no filesystem-level sandboxing, no RAG-style retrieval with relevance scoring. The alternatives table is still "our approach" vs. "weak alternative" vs. "nothing." |
| Honest trade-off comparison | 7/25 | The selected approach's con is "build 阶段需解析 AC，两个模板强耦合需同步修改." This is more honest than iteration-0's "复杂度略增." The pre-revision added the Coupling Constraints section (lines 158-161) which explicitly names the template synchronization risk. However, the temporal coupling between BuildIndex and execution time is still understated in the trade-off table — it appears only in the risk table, not in the alternatives comparison. |
| Chosen approach justified against benchmarks | 8/25 | No industry benchmarks to justify against. The justification is purely internal reasoning. The pre-revision's "dual-layer defense" framing (line 41) is a better argument than iteration-0's blanket rejection of prompt constraints, but it is still not benchmarked against how other systems solve similar context-scoping problems. |

**Deductions**:
- Straw-man alternative: The "pure prompt constraint" row still functions as a straw-man despite the Challenge Override note, because the alternatives table itself does not acknowledge that the selected approach also uses prompt constraints. The table says pure prompt is "Rejected: 效果不确定" while the selected approach's prompt filtering component faces the same uncertainty. (-20 pts)

### 4. Requirements Completeness: 80/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 30/40 | Three scenarios (pure-doc, mixed, no-doc). The pre-revision added edge case handling for missing AC (lines 62, 167): heading tolerance (`## Acceptance Criteria`, `## Acceptance criteria`, `## 验收标准`), build-time warning, and `free-review` mode fallback. This substantially improves coverage. Remaining gaps: (1) scenario where multiple doc tasks share deliverable files — how are overlapping AC handled? (2) scenario where `docs/` directory does not exist yet (first doc task still in progress). |
| Non-functional requirements | 25/40 | Token reduction, AC coverage, backward compat listed. The pre-revision did not add quantification. Quote: "review-doc 执行 token 消耗降低" — lower by how much? The NFR is unverifiable as stated. No security consideration (can task-file AC content inject agent instructions via prompt injection?). |
| Constraints & dependencies | 25/30 | The pre-revision added heading tolerance strategy (line 62), build-time validation (line 167), and explicit coupling constraint (lines 158-161). This resolves most gaps from iteration-0. Remaining: the naming collision between `BodyContext.AcceptanceCriteria []string` (existing PRD field) and the proposed `AcceptanceCriteria map[string]string` (new per-task field) is an unstated schema constraint. |

**Deductions**:
- Vague language: "token 消耗降低" without threshold (-20 pts) — still unquantified after revision.

### 5. Solution Creativity: 40/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 10/40 | The proposal honestly states "无特殊创新" and the pre-revision does not change this. Build-time data extraction for runtime consumption is a well-known pattern (compiled configs, AOT optimization). The allowlist approach is standard security practice (principle of least privilege applied to file access). |
| Cross-domain inspiration | 10/35 | No cross-domain borrowing. The "Assumptions Challenged" section (lines 98-103) demonstrates good analytical thinking but does not draw from other domains. The proposal could have referenced: compiler dead-code elimination (removing unreachable code = removing irrelevant files), database query optimization's predicate pushdown (pushing filters closer to data source = pushing filters into build time), or information retrieval's relevance filtering. |
| Simplicity of insight | 20/25 | The core insight — "don't give the agent files it shouldn't read" — remains elegant. The "Assumptions Challenged" table is the strongest intellectual content: correctly identifying that runtime discovery is unnecessary when data is available at build time. The pre-revision's "dual-layer defense" framing adds a useful conceptual model. |

### 6. Feasibility: 65/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 25/40 | The pre-revision significantly improved technical specification. The four-thread breakdown (lines 78-86), BodyContext schema naming (`AcceptanceCriteria map[string]string`), and AC extraction pipeline description provide enough detail for implementation. However, the `renderBody()` mechanism description is incorrect — it says `{{range .AcceptanceCriteria}}` but `renderBody()` uses `strings.ReplaceAll`, not Go templates. This means the implementation approach as described will not work without an architectural change to the template rendering system. Verified against actual code: `autogen.go` line 300-357 uses simple string replacement for all placeholders. The `{{range .AcceptanceCriteria}}` syntax would require switching to Go's `text/template`, which is a larger change than the proposal acknowledges. |
| Resource & timeline feasibility | 20/30 | The pre-revision updated the estimate to 4-5 coding tasks (lines 91-96), which is more realistic. However, if the template rendering system needs to change from string replacement to Go templates (to support `{{range}}`), this adds another task. The estimate remains tight. |
| Dependency readiness | 20/30 | All dependencies are internal. The heading tolerance strategy (line 62) reduces extraction fragility. The `extractCheckboxItems()` function in `extract.go` can be adapted for per-task extraction but will need a new function — the existing one only handles single-file extraction. |

### 7. Scope Definition: 65/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 24/30 | Four in-scope items with file paths and change descriptions. The pre-revision added the AC aggregation schema format section and the prompt template restructuring spec, making each item more implementable. The in-scope items now map to the correct modification surface (BodyContext + renderBody, not GetReviewDocTask). However, the `extract.go` changes (new per-task AC extraction function) are not explicitly listed as in-scope despite being required — they are implied by "修改 `build.go`" but actually live in `extract.go`. |
| Out-of-scope explicitly listed | 21/25 | Four out-of-scope items. Clear exclusion of other task types, AC format changes, eval pipeline, and record templates. Good. |
| Scope is bounded | 20/25 | The 4-5 task estimate bounds the work. The pre-revision added explicit task breakdown (lines 91-96). No calendar timeframe given, but the task count is a reasonable proxy. |

### 8. Risk Assessment: 72/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 24/30 | Three risks. The pre-revision substantially improved all three: (1) AC extraction instability now has heading tolerance, build warnings, and `free-review` fallback mode (line 167); (2) stale AC risk now has timestamp annotation and re-index guidance (line 169); (3) over-aggressive filtering risk now explicitly names allowlist strategy (line 170). Missing risks: (a) `renderBody()` architectural mismatch (string replacement vs. Go templates), (b) BodyContext field naming collision between existing `AcceptanceCriteria []string` and proposed `AcceptanceCriteria map[string]string`. |
| Likelihood + impact rated | 23/30 | Ratings are honest. The stale-AC risk is now L/M (line 169) with a concrete mitigation (timestamp + re-index guidance). The extraction instability risk at M/M is fair. The filtering risk at L/M is fair for an allowlist approach. Remaining concern: the template rendering architectural mismatch is an unidentified risk — if `{{range}}` cannot be implemented with current `renderBody()`, the entire AC rendering pipeline needs redesign, which is a High-impact risk. |
| Mitigations are actionable | 25/30 | All three mitigations are now actionable: (1) heading tolerance with specific variants + warning format + `free-review` behavioral fallback, (2) timestamp annotation + re-index instruction, (3) explicit allowlist. The `free-review` mode with `free-review` flag annotation (line 167) is particularly well-specified. Missing: the proposal does not define what happens if ALL doc tasks have no AC — does the entire review become freeform? |

### 9. Success Criteria: 62/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 42/55 | Six checkbox criteria. The pre-revision added criteria for: `## Acceptance Criteria Summary` section existence (line 175), `> No acceptance criteria defined.` display (line 180), and `free-review` flag (line 180). These are testable. The three "Referenced Documents" criteria (lines 176-178) are directly testable via execution logs. Remaining gap: "所有 doc 任务的 AC 均在 review-doc.md 的汇总区域中可查" (line 179) — how is this verified? By manual inspection? By a build-time assertion? The criterion is structural but the verification mechanism is undefined. |
| Coverage is complete | 20/25 | Success criteria cover: AC summary generation, file exclusion, modification scope, AC completeness, and empty-AC handling. Missing: (1) no criterion for the prompt template restructuring (the autogen template has criteria but the prompt template does not), (2) no criterion for backward compatibility (listed as NFR but no verification), (3) no criterion measuring review quality preservation (Problem 2 is "quality degradation" but no criterion measures that quality is maintained). |

### 10. Logical Consistency: 72/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses stated problem | 28/35 | Build-time AC embedding addresses Problems 1-2 (eliminates need to read irrelevant files). Allowlist discovery addresses Problem 3 (prevents mis-edits). The pre-revision's "dual-layer defense" framing (line 41) resolves the logical inconsistency flagged in iteration-0 — the proposal now explicitly acknowledges that prompt constraints are complementary to structural changes, not rejected outright. Remaining tension: the allowlist is still a prompt instruction, so Problem 3's mitigation is fundamentally the same mechanism the proposal argues is unreliable. The structural change (AC embedding) makes the allowlist *less necessary* but does not *guarantee* the agent won't read irrelevant files — that still depends on prompt compliance. |
| Scope <-> Solution <-> Success Criteria aligned | 22/30 | Scope items map to solution components. Success criteria cover most scope items. Gap: the prompt template restructuring is in scope (line 112) but has no dedicated success criterion — the success criteria focus on the autogen template and execution behavior, not on the prompt template's structural correctness. |
| Requirements <-> Solution coherent | 22/25 | Requirements (token reduction, AC coverage, backward compat) map to solution. The AC coverage requirement now has a concrete verification path via the `## Acceptance Criteria Summary` section. The heading tolerance strategy (line 62) maps to the AC extraction implementation. The coupling constraint (lines 158-161) maps to the template synchronization requirement. Coherent overall. |

---

## Phase 3: Blindspot Hunt

1. **renderBody() architectural mismatch is an unacknowledged showstopper**: The proposal says `renderBody()` uses `{{range .AcceptanceCriteria}}` (line 80), but `renderBody()` in `autogen.go` uses `strings.ReplaceAll` for all placeholder substitution. The `{{range}}` syntax is Go template syntax, not a simple string placeholder. Implementing this requires either (a) switching the template rendering from string replacement to Go's `text/template` package (a significant architectural change affecting all 12+ autogen templates), or (b) using a new simple placeholder like `{{DOC_TASK_AC}}` and building the markdown in Go code before substitution. The proposal does not recognize this architectural decision point.

2. **BodyContext field naming collision**: The proposal says to add `AcceptanceCriteria map[string]string` to BodyContext (line 79). But BodyContext already has `AcceptanceCriteria []string` (autogen.go line 60). This is a direct naming collision in the same struct. The new field needs a different name (e.g., `DocTaskCriteria`) or the existing field needs renaming.

3. **No rollback plan**: If the build-time AC extraction produces incorrect or incomplete results, there is no documented way to revert to the current behavior. For a production pipeline, this is a serious omission. The proposal should specify: (a) a feature flag to disable the new extraction, (b) a manual override mechanism, or (c) a documented revert procedure.

4. **Review quality regression is unmeasured**: Problem 2 is "审校质量下降" (review quality degradation due to attention dilution). The success criteria measure structural changes (agent doesn't read certain files) but do not measure whether review quality actually improves. The proposal assumes quality improvement follows from attention focus, but this is an untested hypothesis. A before/after comparison of review quality on the same feature would validate the core claim.

5. **The free-review fallback creates a perverse incentive**: Line 167 specifies that when AC is missing, the agent enters "free-review mode" using "generic document quality standards." This creates an incentive to skip writing AC in doc tasks — the agent will still review, just with less specific criteria. The proposal does not discuss whether this is acceptable or whether missing AC should be a hard failure.

---

## Bias Detection Report

- Annotated regions: 9 attack points / 12 paragraphs = density 0.75
- Unannotated regions: 8 attack points / 20 paragraphs = density 0.40
- Ratio (annotated/unannotated): 1.88

**Interpretation**: Annotated (pre-revised) regions received ~1.9x more scrutiny than unannotated regions. This is expected per the annotated review protocol — pre-revised regions receive focused audit for revision-introduced issues. The absolute counts are close enough (9 vs 8) that the density difference is driven primarily by paragraph count, not by disproportionate harshness toward revised content. No evidence of systematic bias against revisions.

---

## Summary of Deductions Applied

| Deduction Type | Instances | Points |
|---------------|-----------|--------|
| Vague language without quantification | 2 instances ("实际执行记录显示", "token 消耗降低") | -40 pts |
| Straw-man alternative | 1 instance (pure prompt rejected while selected approach also uses prompt constraints) | -20 pts |

---

## SCORE: 606/1000

## DIMENSIONS:
  Problem Definition: 80/110
  Solution Clarity: 90/120
  Industry Benchmarking: 40/120
  Requirements Completeness: 80/110
  Solution Creativity: 40/100
  Feasibility: 65/100
  Scope Definition: 65/80
  Risk Assessment: 72/90
  Success Criteria: 62/80
  Logical Consistency: 72/90

## ATTACKS:

1. **Problem Definition**: Evidence remains anecdotal — "实际执行记录显示 Referenced Documents 包含任务文件和记录文件" (line 22) names no specific execution record, provides no count, no percentage. After revision, this is still a hand-wave, not data. Must cite at least one specific execution log with concrete file counts.

2. **Solution Clarity**: `renderBody()` mechanism description is architecturally incorrect — "renderBody() 使用 {{range .AcceptanceCriteria}} 展开" (line 80). Verified against `autogen.go`: `renderBody()` uses `strings.ReplaceAll` for all placeholders, not Go template ranges. The `{{range}}` syntax requires switching to `text/template`, which is an architectural change affecting all 12+ autogen templates. Must specify whether the rendering system is being changed or a different placeholder strategy is used.

3. **Solution Clarity**: BodyContext field naming collision — "新增 AcceptanceCriteria map[string]string 字段" (line 79). The existing `BodyContext` already has `AcceptanceCriteria []string` at `autogen.go` line 60. Adding another field with the same name is a compile error. Must use a distinct field name (e.g., `DocTaskCriteria map[string]string`).

4. **Industry Benchmarking**: Zero industry references. No citation of how LangGraph, CrewAI, AutoGen, or any agent framework handles context scoping, task isolation, or context window management. The proposal operates in a vacuum of its own reasoning. Must reference at least one industry-validated pattern for agent context management.

5. **Industry Benchmarking**: Straw-man alternative persists — "纯 prompt 约束 | Rejected: 效果不确定" (line 70). The selected approach's "allowlist discovery" is itself a prompt instruction. The alternatives table contradicts the Challenge Override note (line 41) which admits prompt constraints are complementary. Must either update the table to reflect that the selected approach also uses prompt constraints, or remove the categorical rejection of pure-prompt approaches.

6. **Requirements Completeness**: NFR "token 消耗降低" (line 53) has no measurable threshold. After revision, still no target percentage or token count. Must specify a verifiable threshold (e.g., "reduce Referenced Documents count by 50%+ for typical features with 3+ doc tasks").

7. **Feasibility**: Template rendering architectural mismatch is an unacknowledged risk. The proposal assumes `{{range .AcceptanceCriteria}}` works in the current rendering system, but the system uses `strings.ReplaceAll`. If this requires switching to Go templates, the scope expands to all 12+ autogen template renderings. Must explicitly address this architectural decision.

8. **Risk Assessment**: The `renderBody()` architectural mismatch (attack #2) is an unidentified risk. If the template system cannot support `{{range}}` without refactoring, this is a High-impact, Medium-likelihood risk that the proposal does not acknowledge. Must add this risk to the risk table with an explicit mitigation.

9. **Success Criteria**: "所有 doc 任务的 AC 均在 review-doc.md 的汇总区域中可查" (line 179) — the verification mechanism is undefined. Is this checked by a build-time assertion? By manual inspection? By a test? Must specify how completeness is verified.

10. **Success Criteria**: No criterion for the prompt template restructuring. The prompt template (`prompt/data/doc-review.md`) is explicitly in scope (line 112) but no success criterion verifies that the restructured prompt produces correct agent behavior. Must add at least one criterion for prompt template correctness (e.g., "prompt template Step 1 references AC Summary section, Step 2 uses docs/ allowlist, Step 3 includes modification restriction").

11. **Logical Consistency** (conflict-with-pre-revision): The pre-revision added the "dual-layer defense" Challenge Override (line 41) to resolve the contradiction between rejecting prompt constraints and using them. However, this creates a new inconsistency: the Alternatives table (line 70) still says pure prompt is "Rejected: 效果不确定" while the selected approach explicitly relies on prompt constraints as half its defense. The table and the Challenge Override contradict each other. Must update the Alternatives table to reflect the nuanced position stated in the Challenge Override.

12. **Logical Consistency**: Review quality regression is unmeasured. Problem 2 states "审校质量下降" but no success criterion measures whether quality actually improves after the change. The proposal assumes focus = quality improvement, which is an untested hypothesis. Must add a qualitative or quantitative quality measure.

13. **Scope Definition**: `extract.go` changes are not listed as in-scope. The proposal says "修改 `build.go`" (line 109) for the AC extraction pipeline, but the extraction logic (heading matching, content parsing) would naturally extend `extract.go`, not `build.go`. `BuildIndex()` in `build.go` would call the new extraction function, but the function itself belongs in `extract.go`. Must either add `extract.go` to in-scope or clarify that `build.go` subsumes all extraction logic.

14. **Feasibility**: The `free-review` fallback (line 167) introduces a behavioral risk that is not assessed. When all doc tasks lack AC sections, the agent performs a completely freeform review with no quality baseline. This could produce worse results than the current approach (where at least the agent reads the task files for context). The proposal treats this as a minor degradation path but it could be a major behavioral regression for features with poorly structured doc tasks.

15. **Blindspot**: No rollback plan. For a production pipeline ("已合并并投入使用"), the absence of any revert strategy is a serious omission. If the new AC extraction produces incorrect results or the filtered prompt causes the agent to miss critical quality checks, there is no documented path back to the current behavior.

16. **Blindspot**: No migration consideration for existing `review-doc.md` files. `BuildIndex()` is idempotent and regenerates `review-doc.md` only if the file does not already exist (verified in `build.go` lines 273-274: `if _, err := os.Stat(mdPath); os.IsNotExist(err)`). This means existing `review-doc.md` files will NOT be updated with the new AC Summary section. The proposal does not address whether users need to manually delete existing review-doc task files to get the new format.
