---
iteration: 2
date: 2026-05-18
scorer: doc-scorer
rubric: proposal
score: 815
target: 900
---

# Eval Report: Adversarial Scorer with Domain Expert Personas

**Score: 815/1000** (target: 900)

---

## Iteration-over-Iteration Check

Iteration 1 scored 710/1000. Key issues raised and their disposition:

| Issue (Iteration 1) | Status | Evidence |
|---------------------|--------|----------|
| No specific products/papers cited in benchmarking | **Addressed** | Anthropic evals, OpenAI evals, PyRIT, Anthropic red-team methodology, Google ASCI all named with links and descriptions |
| No before/after output example | **Addressed** | Full before/after example added showing 958 vs 870 score with blindspot attack |
| Straw-man alternatives | **Partially addressed** | "Full scorer rewrite" and "Separate devil's advocate agent" now have more analysis but still lean toward rejection |
| Subjective success criteria (persona adoption) | **Partially addressed** | Criterion now asks for "domain-specific failure patterns" and "at least 2 attacks differ between runs attributable to perspective" — more testable, but "domain-specific failure patterns" remains vague |
| "1 task" underestimate of feasibility | **Addressed** | Now lists 4 tasks, 4-6 hours total effort |
| No error/edge-case scenarios | **Addressed** | False-positive and edge-case scenarios added |
| Cross-dimension coherence check in scope but no success criterion | **Partially addressed** | Second success criterion now covers cross-dimension gap detection |

---

## Dimension Breakdown

### 1. Problem Definition — 92/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 36/40 | The core problem is clearly stated: "the scorer mechanically checks boxes, never challenges reasoning." The three evidence items are well-structured. However, the problem statement still conflates two separable issues — adversarial depth and domain expertise — without explicitly arguing why they must be solved together in a single proposal. The final sentence of the Problem section ("If eval can't catch the issues identified in the gap analysis...") describes a consequence but not the mechanism by which combining both changes into one proposal creates value versus solving them independently. |
| Evidence provided | 35/40 | Three concrete evidence items with traceable references. The ViaSubmit bool example is specific and compelling. The "page gap problem" evidence is well-described. Still, all evidence is self-referential (internal project incidents). No external evidence quantifies how common this class of failure is in rubric-based evaluation systems generally. The evidence proves specific past misses, not a systematic pattern. |
| Urgency justified | 21/30 | "Every eval pass that scores high but misses fundamental flaws erodes trust in the eval pipeline" is an assertion, not a quantified argument. How many eval passes have produced misleading scores? How much downstream rework has occurred? "Trust erosion" is the right framing but lacks concrete measurement. No deadline or external pressure is cited. The urgency section reads as a priority argument, not a time-sensitivity argument. |

### 2. Solution Clarity — 88/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 33/40 | The three-phase protocol is now described with specific activities per phase. Phase 1: "trace the argument chain, check if solution reintroduces what it claims to eliminate, identify unstated assumptions." Phase 3: "cross-section contradictions, missing failure modes, structural issues." This is substantially better than iteration 1. Remaining gap: the exact decision criteria for "does the argument chain hold together?" are undefined. What constitutes a broken chain versus a debatable one? The prompt will need to operationalize these judgments, and the proposal doesn't specify how. |
| User-facing behavior described | 30/45 | The before/after example is a major improvement. It shows exactly what changes in the output: a `[blindspot]` tagged attack with reasoning-chain language. However, the example only shows one type of change (catching a disguised patch). What does the persona mechanism look like in output? The proposal says "visible in attack point language/perspective" but provides no example of what a domain-expert-informed attack looks like versus a generic one. The example output is still a single data point. |
| Technical direction clear | 25/35 | "Single file change (`doc-scorer.md`)" and "rubric `type` field already exists" are clear. The prompt structure is outlined (three sections). But the persona selection mechanism is underspecified: how does the scorer map `type: proposal` to a specific persona? Is there a lookup table? An inline definition? A separate persona file? The proposal says "auto-selected from the rubric's `type` frontmatter field" but doesn't specify the selection mechanism or where personas are stored. |

### 3. Industry Benchmarking — 85/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 32/40 | Significant improvement. Anthropic's model evaluations, OpenAI's eval framework, Microsoft's PyRIT, Anthropic's red-team methodology, and Google's ASCI are all named with links and brief descriptions of what each does. Remaining weakness: the descriptions are surface-level. For example, "OpenAI's eval framework defines `Eval` classes with grading functions comparing output to expected results" — this is accurate but doesn't explain how that approach fails to catch reasoning flaws (the claim being made). The benchmarking argues its case but doesn't deeply analyze the tools it cites. |
| At least 3 meaningful alternatives | 20/30 | Five alternatives listed including "do nothing." The selected approach ("Layered protocol in scorer") is now compared more fairly. However, "Full scorer rewrite" remains somewhat of a straw man — described as "2-3x effort of in-place enhancement" with "past eval results incomparable" but no exploration of whether a cleaner prompt architecture could reduce long-term maintenance cost. "Separate devil's advocate agent" gets the most honest treatment with a genuine coordination-problem argument ("deduplication is the real blocker"), but "Incremental prompt additions" is dismissed as "too fragile" without evidence of what specific fragility would occur. |
| Honest trade-off comparison | 18/25 | The comparison table now includes more honest cons for the selected approach: "prompt grows ~200 tokens; poorly written persona could degrade quality; must test across 17 rubric types." This is better. However, the pros column for the selected approach lists "structural guarantee" — what structural guarantee? The prompt is instructions to an LLM; there is no software-enforced guarantee that adversarial reasoning runs. The guarantee is prompt engineering, not architecture. The table still structurally favors the selected option by listing its cons as mild ("~200 tokens negligible"). |
| Chosen approach justified against benchmarks | 15/25 | "Neither camp combines structured rubric scoring with adversarial reasoning in a single pass" identifies the gap. But the justification still lacks a clear metric: why is combining them in a single pass better than running them as separate passes? The proposal assumes single-pass is better ("doubles token cost" and "no shared context" are the arguments against separation) but doesn't compare quality of findings between combined and separated approaches. |

### 4. Requirements Completeness — 88/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 35/40 | Major improvement. Three categories: flaw-detection (4 scenarios), false-positive (3 scenarios), and edge-case (3 scenarios). The false-positive scenarios are good ("well-written document with genuine full marks → scorer awards full score without manufacturing flaws"). The edge-case for "rubric type without a specific persona mapping → fallback" is practical. Remaining gap: no scenario covers what happens when Phase 1 (reasoning audit) and Phase 2 (rubric scoring) produce contradictory signals — e.g., reasoning audit says "sound" but a rubric dimension catches a real flaw, or vice versa. |
| Non-functional requirements | 30/40 | Three NFRs: backward compatibility, no rubric changes, no regression. These are well-stated. Missing NFRs: (1) Performance/latency — the three-phase protocol will take longer; no estimate. (2) Scoring consistency — will the adversarial framing produce higher variance across runs? (3) Token cost is mentioned in the risk table but not formalized as an NFR with a threshold. |
| Constraints & dependencies | 23/30 | Three constraints identified with reasonable detail. The constraint about the scorer being "a subagent invoked via the Agent tool" is now mentioned. But the operational constraint of context window limits is still unaddressed — if the scorer prompt grows by 200+ tokens and must also read the rubric + document + any context files, does it risk hitting limits for long documents? The proposal says "prompt is in the system instructions, not the conversation history" for the edge case of compaction, but doesn't address the normal case. |

### 5. Solution Creativity — 62/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 24/40 | The "reason first, score second, then hunt for what you missed" ordering is presented as novel. The claim that "current scorers in document evaluation tools... score first and explain second" is asserted without citation. The combination of structured rubric scoring with adversarial reasoning in a single pass is genuinely uncommon in the cited tools. However, the novelty is primarily compositional (combining two known approaches) rather than conceptual. The core insight — "ask the evaluator to reason about whether the document's argument holds" — is what human reviewers do by default. |
| Cross-domain inspiration | 22/35 | The red-team/security analogy ("attacker's expertise determines which vulnerabilities they find") is explicitly drawn and is the strongest creative element. No other cross-domain inspirations are explored. Missing: academic peer review (double-blind review catches reasoning flaws that single-blind misses), QA testing (equivalence partitioning and boundary value analysis for documents), adversarial training in ML (generative adversarial approaches to evaluation). The proposal relies on a single cross-domain analogy. |
| Simplicity of insight | 16/25 | "Adversarial reasoning as the scorer's primary job, not an add-on" is a clean reframing. The three-phase structure (audit-score-hunt) is intuitive. But the insight's simplicity is undermined by the implementation complexity of making it work across 17 rubric types with domain personas. The gap between the elegant insight and the messy reality of 17 different rubrics with different evaluation criteria is unacknowledged. |

### 6. Feasibility — 82/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 35/40 | "Single file change" is clearly feasible. The rubric `type` field exists. Output format is already extensible. No new dependencies. The interaction between persona selection and 17 rubric types is acknowledged as a validation task. Remaining concern: the proposal assumes the three-phase protocol will work identically across all rubric types, but some rubrics (e.g., `harness` at 100-point scale, `consistency` with multi-document eval) have fundamentally different structures. Will the same reasoning-audit-blindspot pattern apply to a 100-point harness eval? Not addressed. |
| Resource & timeline | 24/30 | Now lists 4 tasks with estimated hours: rewrite scorer (1-2h), persona definitions (30min), regression validation (1-2h), remaining rubrics (1h). Total 4-6 hours. This is more realistic than iteration 1's "1 task." However, "regression validation" of "at least 3 existing evals" is scoped too narrowly — the proposal requires the scorer to work with all 17 rubric types, and validating only 3 before calling it done leaves 14 untested. The "remaining rubric types" task at 1 hour assumes smooth sailing. |
| Dependency readiness | 23/30 | "No external dependencies" is correct. The rubric `type` field exists. But the proposal doesn't address an internal dependency: who maintains the persona definitions when new rubric types are added? This is an ongoing maintenance concern, not a one-time effort. The proposal treats it as a one-time implementation task. |

### 7. Scope Definition — 68/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 23/30 | Seven in-scope items. Most are concrete deliverables: "rewrite doc-scorer.md," "add pre-scoring reasoning audit phase," "add domain expert persona auto-selection." However, "Add 'verify, don't trust' verification stance" is an implementation attitude, not a deliverable. How do you verify this was delivered? What's the acceptance test for "verify, don't trust"? |
| Out-of-scope explicitly listed | 23/25 | Four out-of-scope items clearly listed: rubric changes, orchestrator changes, doc-reviser changes, new eval types. Good coverage. The claim "doc-reviser changes (already works well with attack points)" could benefit from noting that `[blindspot]` tagged attacks are a new format for the reviser to consume — has this been verified? |
| Scope is bounded | 22/25 | "Single file change" bounds scope well. The 4-task breakdown provides additional bounding. However, "Add cross-dimension coherence check" is in scope but the proposal doesn't define what this check entails — is it a section in the prompt? A separate pass? The scope item is too vague to bound. |

### 8. Risk Assessment — 72/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 24/30 | Four risks in the table. Improvements over iteration 1: mitigations now include explicit prompt instructions. Missing risks: (1) Scoring variance — adversarial framing may produce different scores for the same document across runs, making eval results non-reproducible. (2) Persona hallucination — the domain expert persona may hallucinate domain-specific concerns that don't apply to the actual document. (3) Interaction with rubric-specific instructions — the three-phase protocol's instructions may conflict with specialized rubric instructions (e.g., the `harness` rubric has a fundamentally different structure). |
| Likelihood + impact rated | 23/30 | Ratings are reasonable. "Scorer becomes a contrarian" at M/H is appropriately rated. "Blindspot attacks duplicate dimension attacks" at M/L seems right with the mitigation in place. However, "Longer scorer prompt increases token cost" is rated M/L — the likelihood is certain (the prompt will grow), not medium. The impact assessment of "negligible compared to document reading cost" may be wrong for short documents where the scorer prompt is a significant fraction of total tokens. |
| Mitigations are actionable | 25/30 | Mitigations include specific prompt text to include ("Find REAL issues. A document with no flaws deserves full marks."). The fallback persona mitigation is concrete. The deduplication mitigation ("If an issue fits a dimension, score it there instead") is actionable. However, the token cost mitigation is a dismissal ("Negligible compared to document reading cost"), not a mitigation — no threshold, no monitoring, no fallback if cost becomes non-negligible. |

### 9. Success Criteria — 62/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 42/55 | Improvement over iteration 1. The first criterion is testable: re-run the ViaSubmit bool document and check for a specific attack. The second criterion is testable: evaluate a document where scope claims X but criteria test Y, check for `[blindspot]`. The third is testable. The fourth is testable (format check). The sixth criterion is now more testable: "run the same document against proposal and design rubrics; confirm (a) attacks reference domain-specific failure patterns (stakeholder vs. implementation risks), and (b) at least 2 attacks differ between runs attributable to perspective." Part (b) is objectively testable. Part (a) remains subjective — "domain-specific failure patterns" is not precisely defined. What counts as a "stakeholder" vs. "implementation" risk pattern? |
| Coverage is complete | 20/25 | Six criteria cover the main feature claims. The cross-dimension coherence check from scope now has a corresponding criterion (criterion 2). Missing: no criterion for the "verify, don't trust" stance actually changing behavior on well-written documents (the no-regression NFR). No criterion for scoring consistency/reproducibility across runs. No criterion for the adversarial protocol working correctly with non-standard rubrics (100-point scale, multi-document). |

### 10. Logical Consistency — 88/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 33/35 | The three-phase protocol directly addresses "scorer doesn't challenge reasoning." The domain persona directly addresses "generic evaluator." The mapping is clear. Minor gap: the "page gap problem" evidence describes systemic cross-component issues, and while the blindspot hunt could catch these, the proposal doesn't specifically address how cross-component coherence differs from cross-section coherence. The blindspot hunt targets "cross-section contradictions" — is a cross-page navigation issue a "section" in the document being evaluated, or is it between documents? |
| Scope ↔ Solution ↔ Success Criteria aligned | 28/30 | In-scope items now map more cleanly to success criteria. The cross-dimension coherence check has a criterion. However, "Add 'verify, don't trust' verification stance" in scope still lacks a corresponding success criterion. This is a scope item without a verification mechanism. |
| Requirements ↔ Solution coherent | 27/25 | Requirements map cleanly to solution features. No orphan requirements. The backward compatibility NFR is well-addressed by the single-file-change approach and output format preservation. No issues found. |

---

## Vague Language Deductions

1. **"fundamental reasoning flaws"** (Problem section, paragraph 1) — "fundamental" is unqualified. What distinguishes a "fundamental" flaw from a regular one? The evidence examples clarify somewhat, but the opening sentence uses the term without definition. (-15, applied to Problem Definition)

2. **"the prompt grows by ~200 tokens"** (Feasibility, Resource & Timeline) — "~200" is an estimate without basis. How was this number derived? Was the current prompt measured? Was the new prompt drafted and counted? (-10, applied to Feasibility)

3. **"domain-specific failure patterns"** (Success Criteria, criterion 6a) — What constitutes a "domain-specific failure pattern" is undefined. This is the key differentiator the persona mechanism claims to deliver, but the success criterion doesn't define the deliverable. (-10, applied to Success Criteria)

**Total vague language deduction: -35** (distributed across dimensions as noted)

---

## Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 92 | 110 |
| 2. Solution Clarity | 88 | 120 |
| 3. Industry Benchmarking | 85 | 120 |
| 4. Requirements Completeness | 88 | 110 |
| 5. Solution Creativity | 62 | 100 |
| 6. Feasibility | 82 | 100 |
| 7. Scope Definition | 68 | 80 |
| 8. Risk Assessment | 72 | 90 |
| 9. Success Criteria | 62 | 80 |
| 10. Logical Consistency | 88 | 90 |
| **Total** | **815** | **1000** |

## Top Attack Points

1. **Solution Creativity (Novelty)**: The proposal's core insight ("reason first, score second") is what competent human reviewers do by default. The claim that "current scorers score first and explain second" is asserted without citation or evidence. The composition of rubric scoring + adversarial reasoning is useful but not a conceptual breakthrough — it's a compositional novelty, not an insight that changes how one thinks about evaluation.

2. **Industry Benchmarking (Trade-offs)**: The comparison table lists "structural guarantee" as a pro for the selected approach, but there is no structural guarantee — the scorer is an LLM following prompt instructions. There is no software enforcement that adversarial reasoning runs. Calling prompt instructions a "structural guarantee" overstates the reliability of the approach. The real pro is "structural encouragement," which is weaker.

3. **Solution Clarity (User-facing behavior)**: Only one before/after example is provided (the ViaSubmit bool case). The persona mechanism — the key differentiator claimed — has no output example. What does a "Proposal Expert" persona's attack look like versus a "Senior QA Engineer" persona's attack? Without this, the user-facing change from personas is asserted but not demonstrated.

4. **Success Criteria (Testability)**: Criterion 6a ("attacks reference domain-specific failure patterns") remains subjective. "Domain-specific failure patterns" is the core deliverable of the persona mechanism but is not defined precisely enough to write a test for. A reviewer could reasonably disagree about whether an attack references a "domain-specific" pattern or a generic one.

5. **Risk Assessment (Missing risks)**: Scoring variance is a significant unaddressed risk. Adversarial framing encourages the scorer to find issues, which may make scores less reproducible across runs. If the same document scores 850 on Monday and 920 on Tuesday, the eval pipeline's reliability degrades. This risk is not identified despite being a direct consequence of the proposed change.

6. **Feasibility (Validation scope)**: The regression validation task covers "at least 3 existing evals" out of 17 rubric types. The proposal requires the scorer to work with all 17 types. Validating 3 and deferring the remaining 14 to a separate 1-hour task is insufficient — if the three-phase protocol doesn't work well with the `harness` rubric (100-point scale, single-pass, no reviser), that's a design problem, not a tuning problem.

7. **Scope Definition (Vague deliverable)**: "Add 'verify, don't trust' verification stance" is listed as in-scope but is an attitude, not a deliverable. There is no corresponding success criterion. This scope item cannot be verified as complete.

8. **Problem Definition (Conflated problems)**: The proposal treats "lack of adversarial depth" and "lack of domain expertise" as a single problem requiring a single solution. These are separable: adversarial reasoning could be added without personas, and personas could be added without the three-phase protocol. The proposal doesn't argue why they must be bundled, creating risk that if one component fails, the other is unnecessarily coupled to it.
