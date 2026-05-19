# Evaluation Report: Simplify breakdown-tasks Prompt

**Proposal**: `docs/proposals/simplify-breakdown-tasks-prompt/proposal.md`
**Iteration**: 2
**Date**: 2026-05-19
**Evaluator**: CTO Expert Persona

---

## Previous Iteration Issues Addressed

| Issue from Iteration 1 | Status | How Addressed |
|------------------------|--------|---------------|
| Zero industry references | Resolved | Added 3 industry patterns: LangChain/DSPy conditional prompt assembly, RAG-based rule retrieval, feature-flag/conditional compilation (lines 121-131) |
| Missing error/edge case requirements | Resolved | Added Error Handling & Edge Cases section (lines 159-164): rule file missing, empty, malformed artifact, wrong file loaded |
| "Functionally identical" under-specified | Resolved | Defined validation protocol with specific test case, tolerance (+/-1 task count), structural equality via JSON diff (lines 179-181) |
| Two pain points not measured (execution stability, learning curve) | Resolved | Added measurable criteria: 3 consecutive structurally equivalent runs (line 184), 5-minute rule location (line 185) |
| User-facing behavior not described | Resolved | Added explicit User-Facing Behavior section (lines 26-28) |
| No rollback plan | Resolved | Added Rollback Plan section (lines 199-203): git revert, re-validate, post-mortem |
| Combinatorial explosion inconsistency (2^4 vs 2^3) | Resolved | Added explicit Combinatorial Note (line 87) with 3 arguments for why 16 combinations are better than 8 |
| "Execution instability" claimed without evidence | Resolved | Added specific PR references (#117, #119, #121) and failure mode descriptions (line 13) |
| No timeline or resource estimate | Resolved | Added Timeline & Resources section (lines 167-171): 1-2 days, breakdown, delivery strategy |
| No enforcement mechanism discussion | Resolved | Added Enforcement Mechanism section (lines 207-213): explicit trade-off analysis |
| Only 3 shallow alternatives | Resolved | Added CLI-driven assembly and RAG-based retrieval as alternatives (lines 139-140) |

---

## Phase 1: Reasoning Audit (Pre-Score Anchors)

### Argument Chain Trace

1. **Problem -> Solution**: The problem lists 4 pain points with specific PR evidence. The solution (on-demand decomposition) directly addresses token waste through size reduction, execution stability through reduced inline complexity, maintenance through independent rule files, and learning curve through smaller files. The mapping is now complete across all 4 pain points.

2. **Solution -> Evidence**: Token savings are estimated (~8KB to ~15.5KB vs 23KB) but remain projections, not measurements from a prototype. The combinatorial argument (line 87) is well-reasoned: deterministic artifact-based loading is a stronger guarantee than tag-based conditional evaluation. The enforcement mechanism trade-off (line 213) is honest about limitations.

3. **Evidence -> Success Criteria**: Success criteria now cover all 4 pain points with measurable tests. The validation protocol (line 179-181) is specific enough to execute. The execution stability criterion (3 consecutive equivalent runs) is a reasonable pragmatic test.

4. **Self-contradiction check**: The proposal explicitly addresses the 2^4=16 vs 2^3=8 objection (line 87) with three coherent arguments. The rollback plan is consistent with the single-commit delivery strategy. One remaining tension: line 161 accepts "simpler output" when a rule file is missing, while line 178 promises "functionally identical" output — these are slightly contradictory in edge cases.

### Pre-Score Anchors

- This is a significantly stronger proposal than iteration 1. Nearly all identified gaps have been addressed substantively.
- The remaining weaknesses are in areas where the proposal is good but not perfect: industry references are present but shallow, creativity is still limited, and some minor inconsistencies remain.
- The core argument chain is now coherent and honest about trade-offs.

---

## Phase 2: Rubric Scoring

### 1. Problem Definition (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 38/40 | Four pain points with specific failure descriptions. Line 13: "LLMs applied UI-placement rules to backend-only features, omitted phase gates when PRD contained phase sections, and produced inconsistent scope assignments across runs with identical inputs." Concrete and verifiable via the cited PRs. Minor gap: "inconsistent scope assignments" is slightly vague — what kind of inconsistency? |
| Evidence provided | 35/40 | Line 13 cites 3 of 8 recent runs (PR #117, #119, #121) with specific failure modes. Line 18 explains the conditional tag structure (2^3=8 variants). The PR references make the claims verifiable. Gap: no attached failure output or logs, but the specificity is sufficient. |
| Urgency justified | 25/30 | Line 137: "growing worse with each addition" and line 167-171 shows the fix is cheap (1-2 days). The cost-benefit ratio is implicitly strong. Gap: no quantified cost of current 23KB per execution (dollar amount or daily volume). |

**Dimension Total: 98/110**

### 2. Solution Clarity (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 40/40 | The on-demand decomposition model is described with diagrams, file sizes, load conditions, and extraction plan. A reader could implement this without ambiguity. |
| User-facing behavior described | 42/45 | Line 26-28 explicitly states: "nothing changes" for output, "faster execution on simple features" and "reduced token consumption." The key user-facing impact is clear. Minor gap: no mention of whether error messages or progress indicators change. |
| Technical direction clear | 35/35 | Condition-rule matrix, file extraction plan, enforcement mechanism, and error handling provide complete technical direction. The trade-off between LLM-instructed and programmatic enforcement (line 213) is explicitly stated. |

**Dimension Total: 117/120**

### 3. Industry Benchmarking (120 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 30/40 | Lines 121-131 cite 3 patterns: LangChain PipelinePromptTemplate / DSPy module composition, RAG-based retrieval (LlamaIndex, Semantic Kernel), and feature-flag/conditional compilation (#ifdef). These are real, relevant patterns. Gap: references are shallow — 1-2 sentences each, no links to documentation, no version-specific behavior cited. |
| At least 3 meaningful alternatives | 28/30 | Five alternatives: do nothing, rewrite, CLI-driven assembly, RAG-based retrieval, incremental trim. The CLI-driven and RAG alternatives are genuinely different technical approaches that were missing in iteration 1. |
| Honest trade-off comparison | 22/25 | Line 139: CLI-driven assembly "breaks the 'skill is self-contained' model" — an honest, specific con. Line 140: RAG "over-engineered for 4 rules" — fair assessment. Gap: no quantitative trade-off (e.g., "CLI approach would eliminate the Medium-likelihood LLM non-compliance risk at the cost of X days of CLI modification"). |
| Chosen approach justified against benchmarks | 20/25 | Line 131 gives 3 specific reasons for choosing file-existence over RAG (small rule set, boolean check, consistent with forge patterns). Line 213 explicitly acknowledges the trade-off. Gap: "best risk/reward ratio" (line 141) remains a qualitative assertion without evidence that alternatives were prototyped or compared empirically. |

**Dimension Total: 100/120**

### 4. Requirements Completeness (110 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 35/40 | Token savings table covers 5 scenarios (lines 76-84). Error handling section (lines 159-164) covers 4 edge cases: missing file, empty file, malformed artifact, wrong file. Gap: no scenario for conflicting rules between two loaded files (e.g., ui-placement says task A before B, db-schema says B before A). |
| Non-functional requirements | 35/40 | Token cost reduction is well-analyzed. Execution speed improvement is mentioned (line 28). Timeline constraint is stated (1-2 days). Error handling NFRs are specified. Gap: no explicit latency requirement ("loading 4 files should not add more than X seconds"), no concurrency requirement, no security requirement. |
| Constraints & dependencies | 27/30 | Forge distribution model referenced (line 131). Skill-relative paths specified. Templates unchanged. Timeline bounded. LLM reliability dependency acknowledged (line 213). Gap: no specific CLI version dependency or skill version compatibility constraint stated. |

**Dimension Total: 97/110**

### 5. Solution Creativity (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 18/40 | The proposal applies standard on-demand loading — a well-established pattern in software engineering and prompt engineering (LangChain already does this). The deterministic file-existence check is simpler than RAG, not more innovative. No novel technique is introduced beyond standard conditional module loading. |
| Cross-domain inspiration | 10/35 | Line 129 references compiler conditional compilation (#ifdef) and feature-flag configuration as inspiration. This is cross-domain (compilers -> prompts) but remains at the naming level. No evidence that specific techniques from these domains (e.g., dead code elimination, dependency graph analysis, toggle orchestration) were adapted. |
| Simplicity of insight | 22/25 | "Only load rules you need" is elegant. The proposal resists over-engineering (explicitly rejects RAG as "unjustified for current scale"). The skeleton + overlay design is minimal. Minor deduction: the enforcement mechanism adds complexity (redundant inline checks, matrix-first positioning) that partially offsets the simplicity of the core insight. |

**Dimension Total: 50/100**

### 6. Feasibility (100 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 38/40 | Fully feasible with current tech stack. Forge supports subdirectories. Condition-rule matrix uses file-existence checks that LLMs can execute. Risk of LLM non-compliance is acknowledged as Medium-likelihood with specific mitigations (matrix-first, redundant inline checks, fail-safe design). |
| Resource & timeline feasibility | 28/30 | Lines 167-171: "1-2 days (single developer)" with hour-level breakdown and delivery strategy. Specific and reasonable for the scope. Minor gap: "project maintainer" is a role, not a named individual — if the maintainer is unavailable, timeline slips. |
| Dependency readiness | 27/30 | Forge distribution model supports subdirectories. Skill-relative paths are existing patterns. No external dependencies. LLM conditional read capability is the main dependency risk, acknowledged. No upstream changes needed. |

**Dimension Total: 93/100**

### 7. Scope Definition (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 30/30 | Four deliverables: restructure SKILL.md, merge 6 tags into matrix, extract 4 rule files, validate output (lines 147-150). Each is tangible and specific. |
| Out-of-scope explicitly listed | 23/25 | Four items: quick-tasks, CLI commands, task templates, task generation logic (lines 154-157). Clear. Minor gap: "quick-tasks" is listed without explaining its relevance (why mention it at all if there is no risk of scope creep into it?). |
| Scope is bounded | 23/25 | Timeline (1-2 days), single atomic commit, explicit out-of-scope list. Bounded enough to execute. Minor gap: no sprint/iteration assignment (e.g., "complete by end of week X"). |

**Dimension Total: 76/80**

### 8. Risk Assessment (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 28/30 | Five risks identified (lines 189-196), up from 4 in iteration 1. New risk: "Condition-rule matrix too complex for LLM to execute reliably" (line 195). Covers LLM behavioral failures, infrastructure, regression, and design complexity. Gap: no risk of rule files drifting out of sync with skeleton as SKILL.md evolves over time. |
| Likelihood + impact rated | 26/30 | Each risk now has explicit Likelihood (Medium/Low) and Impact (High/Medium) ratings. Improvement over iteration 1 where impact was implied but not explicit. Gap: no numeric scale (e.g., 1-5) for precision — "Medium" is interpretable but not precise. |
| Mitigations are actionable | 26/30 | Mitigations are concrete: "Matrix is FIRST section" (design decision), "Each step re-prints load instruction inline" (redundancy mechanism), "Run forge task validate-index" (validation step). Rollback plan (lines 199-203) provides recovery path: git revert, re-validate, post-mortem. Gap: rollback assumes single commit — if implementation requires multiple commits during development, revert is more complex. |

**Dimension Total: 80/90**

### 9. Success Criteria (80 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 48/55 | Most criteria are objectively testable: file size <=8KB (measurable), all 6 tags replaced (verifiable), zero conditional load when false (testable). Validation protocol (lines 179-181) specifies test case type, "functionally identical" definition (task count +/-1, dependency graph, types/scopes, PRD coverage), and validation method (validate-index + JSON diff + manual review). Execution stability (3 consecutive runs) is testable. Learning curve (5 minutes) is testable with a timer. Gap: "Each rule file is independently understandable" (line 182) — "validated by having an unfamiliar team member explain the rule file's logic" is a soft, human-dependent test. What if no unfamiliar team member is available? What constitutes a passing explanation? |
| Coverage is complete | 22/25 | Success criteria now cover all 4 pain points: token waste (file size), maintenance (independent understandability), execution stability (3 runs), learning curve (5 minutes). Major improvement over iteration 1. Minor gap: no success criterion specifically testing the error handling edge cases defined in lines 159-164 (e.g., "skill still produces valid output when a rule file is temporarily removed"). |

**Dimension Total: 70/80**

### 10. Logical Consistency (90 pts)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses stated problem | 30/35 | Token waste: directly addressed (file size reduction). Maintenance: addressed (independent rule files). Learning curve: addressed (smaller files, 5-min criterion). Execution stability: addressed through reduced inline complexity, but the causal chain is indirect — the proposal argues that loading fewer rules reduces confusion, but doesn't empirically demonstrate this. The 3-run stability criterion will validate it post-implementation, but the proposal's claim that the solution *will* improve stability is an assumption, not a guarantee. |
| Scope <-> Solution <-> Success Criteria aligned | 27/30 | Scope items map to solution (4 rule files, matrix, validation). Success criteria map to scope items. Gap: scope line 150 says "Validate that generated task output remains identical for existing test cases" (plural), but success criteria line 179 specifies a single test case ("a backend+phases+DB feature"). The scope commitment implies broader validation than the criterion measures. |
| Requirements <-> Solution coherent | 23/25 | Token savings table aligns with extraction plan. Error handling aligns with design. Condition-rule matrix aligns with load model. Gap: line 161 says "the LLM should proceed with skeleton-only rules" when a file is missing, producing "valid (simpler) output." But line 178 promises "functionally identical to current output." If ui-placement.md is missing, UI tasks will lack placement constraints — this is a quality degradation, not functional equivalence. The proposal accepts this edge case as tolerable but does not reconcile it with the functional equivalence promise. |

**Dimension Total: 80/90**

---

## Cross-Dimension Coherence Check

1. **Scope vs Success Criteria precision gap**: Scope (line 150) says "existing test cases" (plural); success criteria (line 179) tests one case. If this is a single-commit refactor validated against one scenario, the scope overpromises relative to validation breadth.

2. **Error handling vs functional equivalence tension**: The error handling section (line 161) accepts degraded output when files are missing, but success criteria (line 178) promises functional equivalence. These two positions are slightly contradictory — the proposal should either state that functional equivalence only holds when all rule files are present, or adjust the error handling stance.

3. **Execution stability claim vs evidence gap**: The problem claims execution instability (3/8 runs) and the solution claims to fix it, but the causal mechanism (smaller prompts = more consistent LLM behavior) is asserted rather than evidenced. The 3-run success criterion will test this, but the proposal's reasoning treats the improvement as given.

---

## Phase 3: Blindspot Hunt

1. **[blindspot]** No version coupling or drift detection between skeleton and rule files. Quote: "Stays in skeleton (always loaded): element mapping base rows, scope algorithm, type assignment, intent propagation, template selection, PRD coverage check, granularity basics." If a contributor adds a Step 4c to the skeleton that affects UI placement, the `rules/ui-placement.md` file may become stale because it was written for Steps 0-4b. There is no mechanism (version number, checksum, or reference key) to detect when a rule file is out of sync with the skeleton version. The proposal addresses the initial extraction but not ongoing synchronization.

2. **[blindspot]** The test case selection may miss the most failure-prone scenario. Quote: "a backend+phases+DB feature (the most common scenario per usage logs)." The most common scenario is not necessarily the most failure-prone. The full-stack scenario (all 4 rule files loaded) is the most complex combination and the most likely to expose LLM non-compliance with the condition-rule matrix. Testing only the most common scenario provides false confidence.

3. **[blindspot]** The learning curve success criterion depends on team composition. Quote: "validated by having an unfamiliar team member explain the rule file's logic." If no unfamiliar team member exists (e.g., solo developer or small team where everyone has seen the files), this criterion is unverifiable. The proposal should specify an alternative validation method (e.g., time-boxed documentation-aided comprehension test, or comparison against a baseline time measurement).

---

## Summary

| Dimension | Score |
|-----------|-------|
| Problem Definition | 98/110 |
| Solution Clarity | 117/120 |
| Industry Benchmarking | 100/120 |
| Requirements Completeness | 97/110 |
| Solution Creativity | 50/100 |
| Feasibility | 93/100 |
| Scope Definition | 76/80 |
| Risk Assessment | 80/90 |
| Success Criteria | 70/80 |
| Logical Consistency | 80/90 |
| **Total** | **861/1000** |

---

## Attack Points (All Dimensions)

1. **[Solution Creativity]**: No novel technique beyond industry-standard conditional module loading — the proposal applies LangChain-style prompt decomposition without innovation. Quote: "Our condition-rule matrix serves the same role as LangChain's pipeline stages." Must demonstrate what this approach does that existing frameworks do not, or accept that creativity is not this proposal's strength.

2. **[Success Criteria]**: "Validated by having an unfamiliar team member explain the rule file's logic" is a soft, human-dependent criterion. Quote from line 182. Must define objective passing criteria (e.g., "correctly identifies the load condition, lists the rules enforced, and describes the expected output behavior") or provide an alternative validation method for solo/small teams.

3. **[Success Criteria]**: Test case selection covers most common scenario, not most complex. Quote from line 179: "a backend+phases+DB feature (the most common scenario per usage logs)." The full-stack scenario (all 4 files loaded) is the most failure-prone combination and should be tested explicitly.

4. **[Success Criteria]**: No success criterion tests error handling edge cases. Lines 159-164 define 4 error scenarios but no success criterion validates them. Should add: "skill produces valid output when a rule file is temporarily removed."

5. **[Logical Consistency]**: Error handling accepts degraded output but success criteria promise equivalence. Quote line 161: "the LLM should proceed with skeleton-only rules" vs quote line 178: "Generated tasks are functionally identical to current output." Must clarify that functional equivalence only holds when all applicable rule files are present.

6. **[Logical Consistency]**: Scope says "test cases" (plural) but success criteria specifies one test case. Quote line 150: "existing test cases" vs quote line 179: "a backend+phases+DB feature." Must align scope and criteria — either test multiple scenarios or adjust scope language to "a representative test case."

7. **[Industry Benchmarking]**: Industry references are shallow — 1-2 sentences each, no links, no version-specific behavior. Quote from line 125: "LangChain's PipelinePromptTemplate, DSPy's module composition." Must provide specific references (URLs, version numbers, or paper citations) to demonstrate genuine engagement with prior art rather than name-dropping.

8. **[Risk Assessment]**: Missing risk of rule file drift over time. Quote line 220: extraction plan covers initial creation only. No mechanism to detect when rule files become stale as the skeleton evolves. Must add a drift risk or a maintenance convention.

9. **[Problem Definition]**: Execution instability claim is stronger but still lacks attached evidence. Quote line 13: "observed in 3 of 8 recent task-generation runs (PR #117, #119, #121)." The PR references are verifiable but no actual failure output or comparison is included in the proposal. Must either attach a concrete failure example or acknowledge this as a limitation.

10. **[Solution Creativity]**: Cross-domain inspiration is limited to naming. Quote line 129: "structurally equivalent to a compiler's conditional compilation." No specific technique from compiler design (dead code elimination, dependency analysis) is adapted. Must either deepen the cross-domain application or remove the claim.
