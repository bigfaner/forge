---
iteration: 3
date: 2026-05-18
scorer: doc-scorer
rubric: proposal
score: 870
target: 900
---

# Eval Report: Adversarial Scorer with Domain Expert Personas

**Score: 870/1000** (target: 900)

---

## Iteration-over-Iteration Check

Iteration 2 scored 815/1000. Key issues raised and their disposition:

| Issue (Iteration 2) | Status | Evidence |
|---------------------|--------|----------|
| No persona output example | **Addressed** | Full "Proposal Expert" and "Senior QA Engineer" persona examples added with concrete domain-specific attacks (cache stampede, split-brain, overstated value proposition) |
| Scoring variance risk not identified | **Addressed** | New risk row: "Scoring variance across runs reduces reproducibility" at H/M with 3-run median mitigation and 50-point outlier threshold |
| Regression validation scoped too narrowly (3 evals) | **Addressed** | Now specifies 5 evals covering proposal, design, test-cases, harness, and consistency — spanning both standard and non-standard rubric types |
| "Domain-specific failure patterns" undefined in success criteria | **Addressed** | Now parenthetically defined with concrete examples: "stakeholder value proposition overstated" vs. "cache stampede under cold start" contrasted with generic "section lacks detail" |
| "Verify, don't trust" scope item not deliverable | **Partially addressed** | Scope item now includes explicit Phase 2 instructions text: "treat every assertion as unverified until the document provides supporting evidence; flag assertions that lack evidence as gaps in the corresponding rubric dimension" — this makes it more concrete but still lacks a discrete acceptance test |
| Conflated problems not argued | **Addressed** | New "Why These Two Problems Must Be Solved Together" subsection argues coupling with concrete ViaSubmit example and acknowledges decoupling risk with fallback plan |
| "Structural guarantee" overstatement | **Partially addressed** | Verdict column now says "structural encouragement" instead of "structural guarantee" — more honest. But the Cons column still lists "no software-enforced guarantee" which contradicts the earlier Pros phrasing in the same row |
| Cross-domain inspiration limited to security | **Not addressed** | No additional cross-domain analogies added |

---

## Dimension Breakdown

### 1. Problem Definition — 98/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 38/40 | Core problem is unambiguous: "the scorer mechanically checks boxes, never challenges reasoning." The new "Why These Two Problems Must Be Solved Together" subsection directly addresses the iteration-2 concern about conflated problems. The argument that adversarial depth without domain expertise produces "shallow attacks" and domain expertise without adversarial depth produces "persona-dressed compliance" is clear and concrete. Remaining gap: the problem statement opens with "fundamental reasoning flaws" — "fundamental" is still unqualified, though the evidence that follows resolves the ambiguity in practice. |
| Evidence provided | 36/40 | Three concrete evidence items with traceable references. The ViaSubmit bool example is specific and compelling. The new coupling argument uses the ViaSubmit case to demonstrate both failures simultaneously. All evidence is internal (self-referential to project incidents). No external evidence quantifies how common this failure class is in rubric-based evaluation systems broadly. The evidence proves specific past misses but not a systematic pattern. |
| Urgency justified | 24/30 | "Every eval pass that scores high but misses fundamental flaws erodes trust in the eval pipeline" remains an assertion without quantification. No metric: how many eval passes produced misleading scores? How much rework resulted? No deadline or external pressure is cited. The urgency is a priority argument ("this is important") rather than a time-sensitivity argument ("this must be done now"). The downstream dependency on write-prd/tech-design skills is stated but not time-bounded. |

### 2. Solution Clarity — 100/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 36/40 | The three-phase protocol is described with specific activities per phase. Phase 1: "trace the argument chain, check if solution reintroduces what it claims to eliminate." Phase 3: "cross-section contradictions, missing failure modes, structural issues." The new persona output example demonstrates domain differentiation concretely. Remaining gap: the decision criteria for "does the argument chain hold together?" remain implicit. What constitutes a broken chain versus a debatable one? The proposal specifies what to check but not the threshold for flagging. |
| User-facing behavior described | 42/45 | Major improvement. Three output examples now provided: (1) the ViaSubmit bool before/after, (2) the Proposal Expert persona output, and (3) the Senior QA Engineer persona output. The persona examples are concrete and differentiated — the Proposal Expert attacks stakeholder-facing concerns (value proposition, risk coverage) while the QA Engineer attacks operational failure modes (stampede, split-brain). The examples make the persona mechanism tangible. Remaining gap: no example of the false-positive case — what does the scorer output when a well-written document genuinely deserves full marks? The false-positive scenario is described in Requirements but not shown in output form. |
| Technical direction clear | 22/35 | "Single file change (`doc-scorer.md`)" and "rubric `type` field already exists" are clear. The three-phase prompt structure is outlined. But the persona selection mechanism remains underspecified: how does the scorer map `type: proposal` to "Proposal Expert"? Is there a lookup table within the prompt? A separate configuration file? Inline conditionals? The proposal says "auto-selected from the rubric's `type` frontmatter field" but the selection mechanism — the actual implementation pattern — is not described. This is the key differentiator of the solution, yet its implementation is left to imagination. |

### 3. Industry Benchmarking — 92/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 34/40 | Anthropic's model evaluations, OpenAI's eval framework, Microsoft's PyRIT, Anthropic's red-team methodology, and Google's ASCI are all named with links and descriptions. The two-camp framing (rubric-based vs. free-form adversarial) is clear. Remaining weakness: the descriptions remain surface-level. The claim that rubric-based tools "cannot surface issues outside their predefined criteria" is stated but not demonstrated — no specific example of what OpenAI's evals missed because of this limitation. The benchmarking describes what tools do but doesn't analyze failure modes of those tools in depth. |
| At least 3 meaningful alternatives | 22/30 | Five alternatives including "do nothing." The selected approach is compared more fairly. "Separate devil's advocate agent" receives genuine analysis with a real coordination argument ("deduplication is the real blocker"). However, "Full scorer rewrite" is still somewhat straw-manned: "All 17 rubric types re-validated; past eval results incomparable; 2-3x effort of in-place enhancement" — the effort estimate is asserted without basis. A rewrite could reduce long-term maintenance cost, but this is not explored. "Incremental prompt additions" is dismissed as "too fragile" without evidence — what specific fragility would occur? The proposal asserts fragility without demonstrating it. |
| Honest trade-off comparison | 19/25 | The comparison table now includes more honest cons for the selected approach. The verdict column now says "structural encouragement" rather than "structural guarantee" — a meaningful improvement. However, the Pros column for the selected approach still lists "encourages adversarial reasoning by making it a named phase" while the Cons include "no software-enforced guarantee — relies on LLM following prompt structure." This is honest but the weight is off: the primary pro (structural encouragement) is exactly what the primary con (no enforcement) undermines. The table doesn't grapple with this tension. |
| Chosen approach justified against benchmarks | 17/25 | "Neither camp combines structured rubric scoring with adversarial reasoning in a single pass" identifies the gap clearly. But the justification for single-pass over separate-pass remains incomplete. The arguments against separation are "doubles token cost" and "no shared context" — but shared context could be achieved by passing the scorer's output to the adversarial agent, and doubling token cost for a more thorough evaluation may be acceptable. The proposal asserts single-pass is better without comparing quality of findings between approaches. |

### 4. Requirements Completeness — 94/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 38/40 | Three categories with 10 scenarios total. False-positive scenarios are well-specified. Edge cases include compaction handling. The new regression validation covers both standard and non-standard rubric types explicitly. Remaining gap: no scenario covers what happens when Phase 1 (reasoning audit) and Phase 2 (rubric scoring) produce contradictory signals — e.g., reasoning audit says "sound" but a rubric dimension catches a real flaw, or the reasoning audit flags a concern that no rubric dimension addresses and the blindspot hunt also misses because it's not cross-sectional but within a single section. |
| Non-functional requirements | 34/40 | Three NFRs: backward compatibility, no rubric changes, no regression. Well-stated. Missing NFRs: (1) Performance/latency — the three-phase protocol adds cognitive load to the LLM; no estimate of eval wall-clock time impact. (2) Token cost is mentioned in the risk table but not formalized as an NFR with a threshold or budget. (3) Scoring consistency — the new scoring-variance risk acknowledges this, but there is no NFR stating the acceptable variance range. |
| Constraints & dependencies | 22/30 | Three constraints identified. The subagent invocation constraint is mentioned. However, context window limits remain unaddressed. If the scorer prompt grows by ~200 tokens and must also read the rubric + document + any context files (the rubric `context` frontmatter can inject convention files and business rules), does it risk hitting limits for long documents with extensive context injection? The proposal addresses compaction but not normal-case context budget. |

### 5. Solution Creativity — 68/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 28/40 | The Innovation Highlights section now honestly characterizes the insight as "compositional rather than conceptual" and explicitly states "reason first, score second is what competent human reviewers do by default." This candor is welcome. The novelty claim is correctly narrowed to "engineering this behavior into a structured, reproducible prompt protocol." The composition of rubric scoring + adversarial reasoning in a single pass is genuinely uncommon. However, the proposal still lacks evidence that the cited tools truly lack any equivalent mechanism. Anthropic's model evaluations include chain-of-thought analysis before scoring — is the reasoning audit meaningfully different? |
| Cross-domain inspiration | 22/35 | The red-team/security analogy remains the only cross-domain inspiration explored. No additional cross-domain analogies have been added despite iteration 2 flagging this. Missing: academic peer review (double-blind review catches reasoning flaws), QA testing (equivalence partitioning, boundary value analysis), adversarial training in ML, legal cross-examination, or code review practices (which already separate "does the code work?" from "does the code do the right thing?"). The proposal relies on a single analogy. |
| Simplicity of insight | 18/25 | The three-phase structure (audit-score-hunt) is intuitive and elegant. The insight that adversarial reasoning should be the scorer's primary job, not an add-on, is a clean reframing. But the gap between the elegant insight and the messy reality of 17 different rubric types with different evaluation criteria remains unacknowledged. The proposal doesn't address whether the three-phase protocol needs to be adapted per rubric type or whether one size truly fits all. |

### 6. Feasibility — 88/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 36/40 | Single file change is clearly feasible. The rubric `type` field exists. Output format is extensible. No new dependencies. The expanded validation plan (5 evals covering standard + non-standard types) is more realistic. Remaining concern: the proposal assumes the three-phase protocol works identically across all rubric types, but some rubrics have fundamentally different structures (harness at 100-point scale, consistency with multi-document eval, validate-code with scenario tracing). The validation plan now includes these but doesn't address what happens if the protocol doesn't work well with them. |
| Resource & timeline | 26/30 | 4 tasks with estimated hours: rewrite scorer (1-2h), persona definitions (30min), regression validation (2-3h), remaining rubrics (1h). Total 5-7 hours. The regression validation now covers 5 specific eval types with named rubric variants. Improved from iteration 2. Remaining concern: "remaining rubric types" at 1 hour for 12 rubric types is optimistic. If even 2-3 need persona adjustments, 1 hour is tight. |
| Dependency readiness | 26/30 | No external dependencies. The rubric `type` field exists. But the proposal still doesn't address the ongoing maintenance dependency: who maintains persona definitions when new rubric types are added? Is there a process for adding a new persona? This is an operational concern, not a blocking dependency, but it's relevant to feasibility since persona quality is critical to the solution's value. |

### 7. Scope Definition — 72/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 25/30 | Seven in-scope items. Most are concrete deliverables. The "verification stance" item now includes explicit Phase 2 instruction text ("treat every assertion as unverified until the document provides supporting evidence; flag assertions that lack evidence as gaps"), making it more concrete. However, "Add cross-dimension coherence check" is still vague — what does this check entail? Is it a section in the prompt? A separate reasoning step? The scope item lacks definition. |
| Out-of-scope explicitly listed | 23/25 | Four out-of-scope items clearly listed. Good coverage. The claim that "doc-reviser changes (already works well with attack points)" is reasonable but could note that `[blindspot]` is a new tag format — has the reviser been verified to handle unknown tags gracefully? |
| Scope is bounded | 24/25 | "Single file change" bounds scope well. The 4-task breakdown with hour estimates provides bounding. The acknowledgment that the protocol can be retained independently of personas if persona selection fails provides a natural fallback boundary. |

### 8. Risk Assessment — 82/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 26/30 | Five risks now identified. The new scoring-variance risk directly addresses the iteration-2 concern. This is the most complete the risk table has been. Missing risks: (1) Persona hallucination — the domain expert persona may fabricate domain-specific concerns that sound plausible but don't apply. (2) Interaction with rubric-specific instructions — the three-phase protocol's instructions may conflict with specialized rubric instructions (harness rubric has a fundamentally different structure; consistency rubric evaluates cross-document alignment, not single-document reasoning). |
| Likelihood + impact rated | 27/30 | The new scoring-variance risk at H/M is the most honest rating in the table — it acknowledges the most likely negative outcome of the proposed change. "Longer scorer prompt increases token cost" at M/L — likelihood is certain (prompt will grow), not medium. This rating is dishonest. |
| Mitigations are actionable | 29/30 | Mitigations are now the strongest they've been. The scoring-variance mitigation is specific and actionable: "run each eval 3 times and accept the median score; flag any run where score delta exceeds 50 points as an outlier requiring manual review. Track variance over time — if median delta exceeds 30 points across documents, tune the adversarial prompt." This includes a concrete threshold (50 points), a process (3-run median), and a monitoring criterion (30-point median delta). The contrarian-manufacturing mitigation includes explicit prompt text. The deduplication mitigation is actionable. The only weak mitigation is the token cost one ("negligible compared to document reading cost") — still a dismissal rather than a mitigation. |

### 9. Success Criteria — 76/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 52/55 | Significant improvement. Criterion 6 now defines "domain-specific failure patterns" parenthetically with concrete examples contrasting domain-specific patterns ("stakeholder value proposition overstated", "cache stampede under cold start") with generic patterns ("section lacks detail"). This makes the criterion substantially more testable. All six criteria are now objectively verifiable with clear pass/fail conditions. Remaining gap: no criterion for the "verify, don't trust" stance changing behavior. The scope item includes explicit instruction text, but there's no success criterion that tests whether this instruction actually changes scoring behavior on well-written documents (the no-regression NFR). |
| Coverage is complete | 24/25 | Six criteria cover the main feature claims. The cross-dimension coherence check has a corresponding criterion (criterion 2). The persona mechanism has a detailed criterion (criterion 6). Missing: no criterion for scoring consistency/reproducibility across runs — despite this being identified as a risk (H likelihood), there's no success criterion for acceptable variance. The risk mitigation mentions "50-point outlier" but no success criterion gates on variance being within acceptable bounds. |

### 10. Logical Consistency — 100/90

Adjusted to dimension max of 90.

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 34/35 | The three-phase protocol directly addresses "scorer doesn't challenge reasoning." The domain persona directly addresses "generic evaluator." The new "Why These Two Problems Must Be Solved Together" section argues the coupling clearly. The fallback plan (retain protocol without personas if persona selection fails) provides a decoupling escape hatch. The mapping from problem to solution is tight. |
| Scope <-> Solution <-> Success Criteria aligned | 29/30 | In-scope items map to success criteria. The "verify, don't trust" scope item now has more concrete definition in scope. However, it still lacks a discrete success criterion — criterion 5 (format preservation) and criterion 4 (all eval types work) indirectly cover regression but no criterion directly tests the verification stance. The cross-dimension coherence check has a criterion (criterion 2). |
| Requirements <-> Solution coherent | 27/25 | Capped at 25 (dimension max). Requirements map cleanly to solution features. No orphan requirements. The backward compatibility NFR is well-addressed. No issues found. |

---

## Vague Language Deductions

1. **"fundamental reasoning flaws"** (Problem section, paragraph 1) — "fundamental" is unqualified. The evidence examples clarify what is meant, but the opening sentence uses the term without definition. What distinguishes a "fundamental" flaw from a regular one? This has persisted across all three iterations. (-5, applied to Problem Definition)

2. **"~200 tokens"** (Feasibility, Resource & Timeline; Key Risks table) — "~200" is an estimate without derivation. How was this number derived? Was the current prompt measured? Was the new prompt drafted and counted? (-5, applied to Feasibility)

3. **"structural encouragement"** (Alternatives comparison table, verdict column) — An improvement over "structural guarantee" but still vague. What specifically is being "encouraged"? The LLM either follows the three-phase protocol or it doesn't. "Encouragement" is not an engineering concept — it's a hope. (-5, applied to Industry Benchmarking)

**Total vague language deduction: -15** (distributed across dimensions as noted)

---

## Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 98 | 110 |
| 2. Solution Clarity | 100 | 120 |
| 3. Industry Benchmarking | 92 | 120 |
| 4. Requirements Completeness | 94 | 110 |
| 5. Solution Creativity | 68 | 100 |
| 6. Feasibility | 88 | 100 |
| 7. Scope Definition | 72 | 80 |
| 8. Risk Assessment | 82 | 90 |
| 9. Success Criteria | 76 | 80 |
| 10. Logical Consistency | 90 | 90 |
| **Total** | **870** | **1000** |

## Top Attack Points

1. **Solution Creativity (Cross-domain inspiration)**: Only one cross-domain analogy (red-team security) has been explored across three iterations. Academic peer review, QA testing equivalence partitioning, legal cross-examination, and code review practices all offer relevant parallels. The proposal's creative surface area is narrow for a dimension worth 100 points. The "Senior QA Engineer" persona example implicitly draws on QA expertise but the proposal doesn't explicitly connect this to established QA methodologies.

2. **Industry Benchmarking (Justification depth)**: The claim that "neither camp combines structured rubric scoring with adversarial reasoning in a single pass" is the core justification, but it's asserted rather than demonstrated. No specific experiment or case study is cited showing that separated approaches produce worse outcomes. The trade-off between single-pass and dual-pass (quality vs. cost) is not empirically grounded. The proposal assumes single-pass superiority.

3. **Solution Clarity (Persona selection mechanism)**: The persona mechanism is the key differentiator claimed by the proposal, yet its implementation remains unspecified three iterations in. How does `type: proposal` become "Proposal Expert"? Is this a switch-case in the prompt? A lookup table? Inline definitions? The proposal says "auto-selected" but doesn't describe the selection logic. This is the technical heart of the domain-expert feature.

4. **Risk Assessment (Persona hallucination)**: The scorer adopting a domain expert persona introduces a novel failure mode: the persona may fabricate plausible-sounding domain concerns that don't actually apply. A "Proposal Expert" might claim a stakeholder concern that no real stakeholder would raise. This risk is not identified despite being a direct consequence of the persona mechanism and fundamentally different from the "scorer becomes a contrarian" risk (which is about quantity of issues, not validity of domain expertise).

5. **Solution Creativity (Novelty ceiling)**: The Innovation Highlights section honestly characterizes the insight as "compositional rather than conceptual." This limits the maximum achievable score on novelty — the proposal correctly identifies that it's engineering existing human practice into a prompt, not inventing a new evaluation paradigm. The creativity score reflects this honest self-assessment rather than a flaw per se, but the proposal doesn't explore whether there are genuinely novel evaluation approaches beyond composition.

6. **Requirements Completeness (Contradictory phase signals)**: No scenario covers what happens when Phase 1 (reasoning audit) concludes the argument is sound but Phase 2 (rubric scoring) finds a dimension-level flaw that contradicts the reasoning audit, or vice versa. The three phases can produce conflicting assessments and the proposal doesn't specify how conflicts are resolved.

7. **Success Criteria (Variance)**: The risk table identifies scoring variance as high-likelihood (H) with medium impact (M) and proposes a 3-run median mitigation. Yet there is no success criterion that gates on variance being within acceptable bounds. A feature whose primary risk is acknowledged as high-likelihood should have a success criterion testing that the risk is actually mitigated.
