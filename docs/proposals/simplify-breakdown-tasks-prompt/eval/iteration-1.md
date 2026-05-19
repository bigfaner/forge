# Evaluation Report: Simplify breakdown-tasks Prompt

**Proposal**: `docs/proposals/simplify-breakdown-tasks-prompt/proposal.md`
**Iteration**: 1
**Date**: 2026-05-19
**Evaluator**: CTO Expert Persona

---

## Phase 1: Reasoning Audit (Pre-Score Anchors)

### Argument Chain Trace

1. **Problem -> Solution**: The problem lists 4 pain points (execution instability, maintenance cost, token waste, learning curve). The solution (on-demand decomposition) directly addresses token waste through size reduction and indirectly addresses learning curve through smaller files. However, "execution instability" — the claim that LLMs lose coherence navigating dense branching logic — is addressed by moving the branching to a load-time decision rather than a read-time decision. This assumes the LLM will reliably execute the condition-rule matrix, which is itself a branching construct. The solution may simply relocate the complexity rather than eliminate it.

2. **Solution -> Evidence**: Token savings are estimated but not empirically validated. The table shows "~8KB", "~2KB" etc. — these are projections, not measurements from a prototype. No evidence is provided that reducing prompt size from 23KB to 8-15KB actually improves LLM output quality or consistency.

3. **Evidence -> Success Criteria**: Success criteria measure file sizes, conditional behavior, and functional equivalence. They do NOT measure execution stability improvement (pain point #1) or learning curve reduction (pain point #4). The criteria test the mechanism (smaller files) rather than the outcome (better task generation).

4. **Self-contradiction check**: The proposal replaces 6 conditional tags with a condition-rule matrix + on-demand loading. The conditional logic is not eliminated — it is relocated from inline tags to a dispatch table. This is fine if the goal is token savings, but the problem statement frames conditional complexity as a cause of "execution instability," and the solution preserves conditional complexity in a different form.

### Pre-Score Anchors

- The proposal is technically sound for token reduction but overclaims on execution stability
- Success criteria are well-structured but incomplete relative to the problem statement
- The alternatives section is honest but thin — "incremental trim" vs "rewrite" vs "do nothing" are not deeply analyzed
- No industry benchmarking is present
- Risk assessment is present but narrowly scoped

---

## Phase 2: Rubric Scoring

### 1. Problem Definition (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 35/40 | The four pain points are listed with specificity. Line count, file size, and tag count are concrete. Slight ambiguity: "execution instability" is asserted but never demonstrated — what does "inconsistent task quality" look like concretely? No example of a bad output is provided. |
| Evidence provided | 28/40 | The file metrics (421 lines, 23KB, 6 tags, 2^3=8 variants) are verifiable facts. However, the causal link between these facts and the pain points is asserted rather than demonstrated. Quote: "LLMs lose coherence navigating dense branching logic, producing inconsistent task quality" — no specific failure example or user report is cited. |
| Urgency justified | 20/30 | Quote: "growing worse with each addition" suggests urgency, but no concrete cost of delay is quantified. What is the cost of the current 23KB per execution? How many executions per day? What is the actual dollar or time cost? |

**Dimension Total: 83/110**

### 2. Solution Clarity (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 38/40 | The on-demand decomposition model is clearly described. The before/after structure diagram, file sizes, and load conditions are specific. A reader could explain this back accurately. |
| User-facing behavior described | 30/45 | The proposal describes internal restructuring but does not clearly state what the end user (a developer running `breakdown-tasks`) experiences differently. Will the output look the same? Will it run faster? Will errors change? Quote: "Generated tasks are functionally identical to current output" — this is stated in success criteria but not in the solution description. The observable behavior from the user's perspective is "nothing changes" which is important but not explicitly framed. |
| Technical direction clear | 33/35 | The condition-rule matrix, file extraction plan, and load model provide sufficient technical direction. The dispatch table in Step 2/3/4a is specific enough to implement. |

**Dimension Total: 101/120**

### 3. Industry Benchmarking (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 0/40 | No industry solutions, patterns, or prior art are cited. The proposal does not reference how other prompt engineering frameworks, LLM orchestration tools, or similar codebases handle conditional prompt assembly. |
| At least 3 meaningful alternatives | 22/30 | Three alternatives are listed (do nothing, rewrite, incremental trim). However, the alternatives are shallow — each gets a single row in a table with 1-2 sentences. "Do nothing" is included, which is good. But there are no genuinely different technical approaches explored (e.g., prompt compression, RAG-based rule retrieval, hierarchical skill delegation, dynamic prompt assembly via CLI pre-processing). |
| Honest trade-off comparison | 15/25 | The cons listed are real (regression risk for rewrite, complexity from split files for incremental trim). But the comparison is thin — no quantitative trade-off analysis. Quote for "Incremental trim" cons: "Rules split across files, slightly less 'one-file' readable" — this understates the real risk of LLMs failing to load rule files correctly. |
| Chosen approach justified against benchmarks | 10/25 | Quote: "Selected — best risk/reward ratio" — but no analysis supports this claim. Why is the risk/reward better than, say, having the CLI pre-assemble the prompt based on artifact detection? |

**Dimension Total: 47/120**

### 4. Requirements Completeness (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 30/40 | The token savings table covers 5 scenarios (greenfield backend, backend+phases+DB, backend+existing code, full-stack, UI-only). Edge cases like "PRD has phase keywords but no actual phase sections" or "ui-design.md exists but is empty" are not addressed. The scenario where multiple rule files interact (e.g., UI placement + DB schema in the same task) is implied but not explicitly discussed. |
| Non-functional requirements | 20/40 | Token cost reduction is the primary NFR and is well-analyzed. However: no mention of execution latency impact (loading multiple files may be slower than one large file), no backward compatibility requirements stated (what if other skills reference `breakdown-tasks/SKILL.md` directly?), no error handling requirements (what happens if a rule file is missing or corrupted?). |
| Constraints & dependencies | 22/30 | The proposal references the forge distribution model and skill-relative paths. It correctly identifies the templates directory as unchanged. However, it does not explicitly state the dependency on the CLI's artifact detection mechanism or the LLM's ability to reliably check file existence and conditionally read files. |

**Dimension Total: 72/110**

### 5. Solution Creativity (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 15/40 | On-demand loading / lazy evaluation is a well-established pattern in software engineering. Applying it to prompt assembly is straightforward, not novel. The proposal does not introduce any technique beyond standard conditional file inclusion. |
| Cross-domain inspiration | 5/35 | No evidence of cross-domain inspiration. The proposal applies a basic code organization technique (module splitting with lazy loading) to prompts. No reference to how game engines handle asset loading, how compilers handle conditional compilation, or how other LLM frameworks manage prompt modularity. |
| Simplicity of insight | 18/25 | The core insight — "only load rules you need" — is elegantly simple and directly addresses the token waste problem. It is a good application of a simple principle. |

**Dimension Total: 38/100**

### 6. Feasibility (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 35/40 | The proposal is technically feasible. File splitting, conditional loading via condition-rule matrix, and skill-relative paths are all supported by the current forge distribution model. The main concern is whether LLMs reliably execute the conditional loading, which the proposal acknowledges as a medium-likelihood risk. |
| Resource & timeline feasibility | 22/30 | No explicit timeline or resource estimate is provided. The extraction plan lists 4 files with content sources, which gives a rough scope, but no estimate of how long this takes or who does it. The proposal is small enough that a single developer could likely do it in 1-2 days, but this is not stated. |
| Dependency readiness | 25/30 | The forge distribution model already supports subdirectories (e.g., `templates/`, `experts/`). Adding a `rules/` subdirectory follows the same pattern. No external dependencies required. The LLM's ability to execute conditional file reads is the main dependency risk, acknowledged in the risk table. |

**Dimension Total: 82/100**

### 7. Scope Definition (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 28/30 | Four specific deliverables: restructure SKILL.md, merge conditional tags into matrix, extract 4 rule files, validate output. Each is a tangible deliverable. |
| Out-of-scope explicitly listed | 22/25 | Four items listed as out of scope: quick-tasks, CLI commands, task templates, task generation logic changes. Clear and specific. The only gap is that "quick-tasks" is mentioned without explaining its relevance (is there a risk of scope creep into it?). |
| Scope is bounded | 18/25 | The scope is reasonably bounded — it is a single skill refactoring. However, no explicit timeline or sprint boundary is stated. The proposal could be executed in a defined timeframe but does not commit to one. |

**Dimension Total: 68/80**

### 8. Risk Assessment (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 24/30 | Four risks are identified, all meaningful. The risks cover LLM behavioral failures (loading unconditionally, skipping when needed), infrastructure issues (distribution paths), and regression (wording changes). Missing: risk of the condition-rule matrix itself being too complex for the LLM to execute correctly, risk of rule files becoming stale if SKILL.md is updated without updating them, risk of the "8 variants" problem being replaced by a harder-to-test "NxM variants" problem (skeleton x rule file combinations). |
| Likelihood + impact rated | 22/30 | Likelihoods are rated (Medium, Low, Low, Medium). Impact levels are implied but not explicitly rated on a scale. The assessment is somewhat honest — Medium likelihood for LLM behavioral failures is reasonable. However, no impact severity rating (e.g., High/Medium/Low) is provided alongside likelihood. |
| Mitigations are actionable | 22/30 | Mitigations are somewhat actionable. Quote: "Condition-Rule Matrix is the FIRST section in the skeleton, before any step" — this is a specific design decision. Quote: "Run `forge task validate-index` on generated output; compare task structure, dependencies, types, and scopes against baseline" — this is a concrete validation step. However, no rollback plan is provided. If the refactored version produces worse results, what is the recovery path? |

**Dimension Total: 68/90**

### 9. Success Criteria (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 42/55 | Most criteria are testable: "SKILL.md reduced from ~23KB to <=8KB" (measurable), "All 6 conditional tags replaced" (verifiable), "No rule file is loaded when its condition is false" (testable). However: "Generated tasks are functionally identical to current output for a representative test case" — "representative" is vague. Which test case? How many? What does "functionally identical" mean exactly (same task count? same dependencies? same wording?). "Each rule file is independently understandable" — "understandable" is subjective and not objectively testable. |
| Coverage is complete | 15/25 | Success criteria cover the structural changes well but miss two of the four original pain points: there is no criterion measuring "execution stability" improvement and no criterion measuring "learning curve" reduction. The criteria test the mechanism (file structure) rather than the outcomes (better LLM behavior, easier maintenance). |

**Dimension Total: 57/80**

### 10. Logical Consistency (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses stated problem | 25/35 | The solution directly addresses token waste (pain point #3) and partially addresses maintenance cost (pain point #2) and learning curve (pain point #4). However, it does not convincingly address execution instability (pain point #1). The proposal claims LLMs "lose coherence navigating dense branching logic" but replaces inline branching with a load-time dispatch table — the LLM must still navigate conditional logic, just in a different form. If the LLM loads all rule files unconditionally (Risk #1, Medium likelihood), the execution instability problem is not solved. |
| Scope <-> Solution <-> Success Criteria aligned | 24/30 | Scope items map to the solution (4 rule files extracted, matrix created). Success criteria map to scope items (file sizes, tag replacement, functional equivalence). The main gap: scope says "Validate that generated task output remains identical" but success criteria says "functionally identical for a representative test case" — these should be the same rigor but "representative" weakens the criterion relative to the scope commitment. |
| Requirements <-> Solution coherent | 20/25 | The requirements (implied by the problem statement) map to the solution structure. The token savings table is coherent with the extraction plan. The main inconsistency: the proposal claims "2^3 = 8 possible effective prompt variants that are impossible to test exhaustively" but the new model has skeleton + 4 conditional files = 2^4 = 16 possible combinations (if files are independent), which is worse, not better. However, the proposal argues that most combinations don't occur in practice, which is reasonable but not demonstrated. |

**Dimension Total: 69/90**

---

## Cross-Dimension Coherence Check

1. **Problem claims vs Success criteria gap**: The problem lists 4 pain points. Success criteria measure 2 of them (token waste, maintenance structure). Execution stability and learning curve are claimed as problems but not measured as outcomes. This is a coherence gap between Dimension 1 (Problem Definition) and Dimension 9 (Success Criteria).

2. **Risk assessment vs Solution confidence**: The risk table rates LLM unconditional loading as "Medium likelihood" — this directly undermines the core premise that the solution improves execution stability. If the LLM might load all files anyway, the token savings are the only guaranteed benefit, making this a much narrower proposal than presented.

3. **Alternatives analysis vs Industry Benchmarking**: The alternatives table is the only place where different approaches are discussed. The complete absence of industry benchmarking (Dimension 3) means the "incremental trim" approach was selected without evidence that it is the industry-standard approach to prompt modularization.

---

## Phase 3: Blindspot Hunt

1. **[blindspot]** The proposal does not address who or what mechanism enforces the conditional loading. Quote: "every rule file has a single load condition checked against artifact presence." The proposal assumes the LLM will reliably check artifact presence and conditionally read files. But there is no enforcement mechanism — no CLI pre-processing, no wrapper script. The condition-rule matrix is instructions to the LLM, not a programmatic constraint. This is the single most critical implementation risk and it is treated as a design detail rather than an architectural decision. If the LLM fails to follow the matrix, the entire proposal's benefits vanish.

2. **[blindspot]** The proposal claims "2^3 = 8 possible effective prompt variants that are impossible to test exhaustively" as a problem, but the new model introduces 4 conditional files = 16 possible effective combinations. Quote: "No nested conditions, no overlapping tag scopes. Each condition maps to exactly one rule file." The conditions are independent, so 2^4 = 16 combinations exist. The proposal argues this is simpler because most combinations don't occur, but the combinatorial explosion claim made against the current system applies equally (or more) to the proposed system. This is a mathematical inconsistency.

3. **[blindspot]** No rollback plan is provided. If the refactored skill produces worse task generation, there is no stated recovery procedure. For a refactoring proposal that explicitly acknowledges "subtle behavioral regression" as a medium-likelihood risk, the absence of a rollback strategy is a significant omission.

---

## Summary

| Dimension | Score |
|-----------|-------|
| Problem Definition | 83/110 |
| Solution Clarity | 101/120 |
| Industry Benchmarking | 47/120 |
| Requirements Completeness | 72/110 |
| Solution Creativity | 38/100 |
| Feasibility | 82/100 |
| Scope Definition | 68/80 |
| Risk Assessment | 68/90 |
| Success Criteria | 57/80 |
| Logical Consistency | 69/90 |
| **Total** | **685/1000** |

---

## Attack Points (All Dimensions)

1. **[Industry Benchmarking]**: Zero industry references — the proposal does not cite any prompt engineering frameworks, LLM orchestration tools, or similar refactoring patterns from the industry. Must add at least 3 external references or patterns for how other projects handle conditional prompt assembly.

2. **[Requirements Completeness]**: Missing error/edge case requirements — what happens when a rule file is missing, empty, or the LLM reads the wrong file? No error handling requirements are specified.

3. **[Success Criteria]**: "functionally identical for a representative test case" is under-specified — which test case? How is "functionally identical" defined? Must specify exact validation protocol.

4. **[Success Criteria]**: Two of four pain points are not measured — execution stability and learning curve are claimed as problems but have no corresponding success criterion. Must add criteria for these or remove them from the problem statement.

5. **[Solution Clarity]**: User-facing behavior is not described — the proposal describes internal restructuring without explicitly stating that the end user should see no change. Must frame the solution from the user's perspective.

6. **[Risk Assessment]**: No rollback plan — for a refactoring with medium-likelihood regression risk, a rollback procedure must be specified.

7. **[Solution Creativity]**: No novel technique — the proposal applies basic file splitting to prompts. Must demonstrate why this approach is better than alternatives like CLI-driven prompt assembly or RAG-based rule retrieval.

8. **[Logical Consistency]**: Combinatorial explosion claim is inconsistent — the current 2^3=8 variants are called "impossible to test exhaustively" but the proposed 2^4=16 combinations are not acknowledged as a larger testing surface.

9. **[Problem Definition]**: "Execution instability" is claimed without evidence — no specific failure example, user report, or metric demonstrates this problem. Must provide concrete evidence or downgrade the claim.

10. **[Feasibility]**: No timeline or resource estimate — the proposal does not state how long implementation takes or who will do it.

11. **[blindspot]** No enforcement mechanism for conditional loading — the condition-rule matrix is instructions to the LLM, not a programmatic constraint. This is the critical path for the entire proposal's success.

12. **[blindspot]** Combinatorial count is 2^4=16, not smaller than the current 2^3=8. The proposal's own argument against the current system applies equally to the proposed system.
