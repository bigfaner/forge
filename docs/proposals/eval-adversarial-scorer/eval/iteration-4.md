---
iteration: 4
date: 2026-05-18
scorer: doc-scorer
rubric: proposal
score: 893
target: 900
---

# Eval Report: Adversarial Scorer with Domain Expert Personas

**Score: 893/1000** (target: 900)

---

## Iteration-over-Iteration Check

Iteration 3 scored 870/1000. Key issues raised and their disposition:

| Issue (Iteration 3) | Status | Evidence |
|---------------------|--------|----------|
| Cross-domain inspiration limited to security | **Addressed** | Innovation Highlights now cites three domains: red-team security, academic peer review, and QA equivalence partitioning. Each analogy is developed with a specific connection to the proposal's mechanism. |
| Persona selection mechanism unspecified | **Addressed** | Solution section now states: "auto-selected from the rubric's `type` frontmatter field via an inline lookup table embedded in the scorer prompt. The table maps each `type` value to a persona definition (role description + domain-specific failure patterns to watch for)." Implementation pattern is explicit. |
| Persona hallucination risk not identified | **Addressed** | New risk row: "Persona hallucination — LLM fabricates plausible domain-specific concerns that don't apply" at M/M with quote-based validation mitigation. |
| Conflicting phase signals unaddressed | **Addressed** | Edge-case scenario now specifies: "Phase 2 rubric scores take precedence for DIMENSIONS output; Phase 1 findings are channeled into [blindspot] attacks only if they identify issues not covered by any dimension." |
| No success criterion for scoring variance | **Addressed** | New criterion 7: "Scoring variance gate: for any document, 3 consecutive eval runs produce scores within a 50-point range (median delta < 30 points)." |
| "Verify, don't trust" lacks discrete acceptance test | **Partially addressed** | Scope item now includes explicit instruction text, but still no dedicated success criterion testing the verification stance. Criterion 4 ("all existing eval types work") is regression-adjacent but doesn't specifically test whether the stance changes behavior. |
| Single-pass superiority over dual-pass not empirically grounded | **Addressed** | Industry Benchmarking now adds: "the claim that combining both in a single pass outperforms either alone is a hypothesis supported by the ViaSubmit bool case study...but not yet validated at scale. The regression validation step (Task 3) serves as the empirical test." This is honest and actionable. |
| "Fundamental" in problem statement unqualified | **Not addressed** | "fundamental reasoning flaws" remains unqualified in the opening sentence across all four iterations. |

---

## Dimension Breakdown

### 1. Problem Definition — 100/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 38/40 | Core problem is unambiguous: "the scorer mechanically checks boxes, never challenges reasoning." The coupling argument in "Why These Two Problems Must Be Solved Together" is well-argued with concrete examples. Remaining gap: "fundamental reasoning flaws" in the opening sentence persists without definition across all four iterations. The term "fundamental" is never qualified — what distinguishes a fundamental flaw from a regular one? The evidence clarifies in practice, but the opening assertion remains unqualified. |
| Evidence provided | 37/40 | Three concrete evidence items with traceable references (lesson file, gap analysis document). The ViaSubmit bool example is specific and compelling. The coupling argument uses ViaSubmit to demonstrate both failures simultaneously. Deduction: all evidence is internal/self-referential to project incidents. No external evidence quantifies how common this failure class is in rubric-based evaluation systems broadly. |
| Urgency justified | 25/30 | "Every eval pass that scores high but misses fundamental flaws erodes trust in the eval pipeline" — this is an assertion, not quantified. No metric for how many eval passes produced misleading scores or how much rework resulted. The downstream dependency on write-prd/tech-design skills is stated but not time-bounded. This is a priority argument ("this is important") not a time-sensitivity argument ("solve now because X"). |

### 2. Solution Clarity — 108/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 38/40 | The three-phase protocol is well-described with specific activities per phase. The persona selection mechanism is now explicitly specified: "inline lookup table embedded in the scorer prompt. The table maps each `type` value to a persona definition (role description + domain-specific failure patterns to watch for)." The fallback mechanism is described. Remaining gap: the decision criteria for what constitutes a broken argument chain in Phase 1 remain implicit. What threshold triggers a flag? |
| User-facing behavior described | 44/45 | Three output examples demonstrate the behavior concretely: (1) ViaSubmit bool before/after, (2) Proposal Expert persona output, (3) Senior QA Engineer persona output. The persona differentiation is tangible. The edge-case scenario for conflicting phases now shows how resolution works. Remaining gap: no example of the false-positive case — what does the scorer output when a well-written document genuinely deserves full marks? The false-positive scenario is described in Requirements but not shown in output form. |
| Technical direction clear | 26/35 | "Single file change (`doc-scorer.md`)", "inline lookup table", and "rubric `type` field" provide clear direction. The three-phase prompt structure is outlined. However, the prompt architecture within `doc-scorer.md` is not fully specified. The current `doc-scorer.md` presumably has an existing structure — what sections are kept, what sections are replaced, what sections are added? The proposal describes the abstract protocol but not how it maps onto the actual file. Additionally, the "inline lookup table" — how many entries? Is there one per rubric type? The task breakdown says "6 personas + fallback" in Task 2, but the proposal itself doesn't enumerate them. |

### 3. Industry Benchmarking — 102/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 36/40 | Anthropic's model evaluations, OpenAI's eval framework, Microsoft's PyRIT, Anthropic's red-team methodology, and Google's ASCI are named with links. The two-camp framing is clear. The new honest caveat — "the claim that combining both in a single pass outperforms either alone is a hypothesis...not yet validated at scale" — is a significant improvement. Remaining gap: no failure analysis of specific industry tools. The claim that rubric-based tools "cannot surface issues outside their predefined criteria" is stated but not demonstrated with a specific case where OpenAI's evals or Anthropic's evaluations missed something because of this limitation. |
| At least 3 meaningful alternatives | 25/30 | Five alternatives including "do nothing." "Separate devil's advocate agent" receives genuine analysis with a real coordination argument. However, "Full scorer rewrite" at "2-3x effort" is still asserted without derivation. "Incremental prompt additions" is dismissed as "too fragile" — but fragility is asserted, not demonstrated. What specific fragility? Has an incremental approach been tried and found fragile? The proposal states the conclusion without the evidence. |
| Honest trade-off comparison | 20/25 | The verdict column says "structural encouragement" which is honest. The Pros/Cons for the selected approach acknowledge "no software-enforced guarantee — relies on LLM following prompt structure." However, the primary Pro ("encourages adversarial reasoning by making it a named phase") is directly undermined by the primary Con ("no software-enforced guarantee"). The table doesn't address this tension — it lists both facts without reconciling them. |
| Chosen approach justified against benchmarks | 21/25 | "Neither camp combines structured rubric scoring with adversarial reasoning in a single pass" identifies the gap. The honest admission that this is a hypothesis supported by one case study but not validated at scale is welcome. However, the argument against dual-pass remains primarily about cost ("doubles token cost") and deduplication complexity, not about quality of findings. The proposal doesn't acknowledge that a dual-pass approach could produce higher-quality findings because the adversarial agent can focus exclusively on reasoning without the cognitive load of rubric scoring. |

### 4. Requirements Completeness — 100/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 39/40 | Three categories with 10+ scenarios. False-positive scenarios are well-specified. Edge cases include compaction handling AND the new conflicting-phase-signal scenario with explicit resolution rules ("Phase 2 rubric scores take precedence...Phase 1 findings are channeled into [blindspot] attacks only if they identify issues not covered by any dimension"). This directly addresses the iteration-3 gap. Remaining gap: no scenario covers a document where the rubric itself is poorly constructed — what happens when a rubric dimension is ambiguous or contradictory and the adversarial scorer tries to apply it? |
| Non-functional requirements | 36/40 | Three NFRs: backward compatibility, no rubric changes, no regression. Well-stated. Missing NFRs: (1) Latency/performance — no estimate of wall-clock time impact from the three-phase protocol. The prompt grows ~200 tokens and requires more reasoning passes; this will increase eval time. (2) Token cost is mentioned in the risk table but not formalized as an NFR with a budget. |
| Constraints & dependencies | 25/30 | Three constraints identified. The rubric `type` field dependency is stated. The reviser's existing `[blindspot]` consumption is noted. Missing: context window budget. If the scorer prompt grows by ~200 tokens and must also read the rubric + document + any context files (the rubric `context` frontmatter can inject convention files and business rules), the total context budget for long documents with extensive context injection is not analyzed. The proposal mentions compaction as an edge case but doesn't budget the normal-case context. |

### 5. Solution Creativity — 78/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 32/40 | The Innovation Highlights section honestly characterizes the insight as "compositional rather than conceptual" and acknowledges that "reason first, score second is what competent human reviewers do by default." The composition of rubric scoring + adversarial reasoning + domain persona in a single pass is genuinely uncommon. The acknowledgment that "formal verification of document properties" and "training a specialized critic model" exist but are disproportionate is well-reasoned. Remaining gap: no evidence that Anthropic's model evaluations truly lack any equivalent mechanism. Anthropic's chain-of-thought evaluation before scoring may overlap with the reasoning audit. |
| Cross-domain inspiration | 28/35 | Three cross-domain analogies now explored: red-team security, academic peer review, and QA equivalence partitioning. Each is developed with a specific connection: red-team maps to persona specialization, peer review maps to the "formatting vs. methodology" distinction, QA maps to domain-specific boundary knowledge. This is a meaningful improvement over the single-analogy version. Remaining gap: no exploration of legal cross-examination (challenging witness claims = challenging document assertions), or code review practices (separating "does the code work?" from "does the code do the right thing?" maps directly to rubric scoring vs. reasoning audit). The analogies are relevant but the exploration could go deeper. |
| Simplicity of insight | 18/25 | The three-phase structure (audit-score-hunt) is intuitive and elegant. The insight that adversarial reasoning should be the scorer's primary job is a clean reframing. However, the gap between the elegant insight and the reality of 17 different rubric types with different structures remains unaddressed. The proposal assumes one-size-fits-all for the three-phase protocol but doesn't argue why. A harness rubric (100-point scale) and a consistency rubric (multi-document eval) have fundamentally different evaluation models — does the same three-phase protocol work identically for both? |

### 6. Feasibility — 90/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 36/40 | Single file change is clearly feasible. Inline lookup table is a standard prompt engineering pattern. Fallback mechanism is straightforward. No new dependencies. The expanded validation plan (5 evals covering standard + non-standard types) is realistic. Remaining concern: the proposal assumes the three-phase protocol works identically across all rubric types. Some rubrics have fundamentally different structures — harness (100-point), consistency (multi-document), validate-code (git diff + scenario tracing). If the protocol needs type-specific adaptation, the "single file change" scope expands. |
| Resource & timeline | 27/30 | 4 tasks with estimated hours totaling 5-7 hours. The regression validation covers 5 specific eval types. Improved from earlier iterations. Remaining concern: "remaining rubric types" at 1 hour for 12 rubric types is optimistic. If 2-3 need persona adjustments, debugging and prompt tuning for each could easily exceed 1 hour total. |
| Dependency readiness | 27/30 | No external dependencies. The rubric `type` field exists. The fallback persona covers unmapped types. Remaining gap: the proposal doesn't address the ongoing maintenance dependency for persona definitions. When a new rubric type is added, who defines the persona? What process ensures persona quality? This is an operational concern that affects long-term feasibility. |

### 7. Scope Definition — 74/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 26/30 | Seven in-scope items. Most are concrete deliverables. The "verification stance" item now includes explicit Phase 2 instruction text. The "persona auto-selection" item is now concrete with the inline lookup table specification. Remaining gap: "Add cross-dimension coherence check" is still vague — what does this check entail? Is it a section in the prompt? A reasoning step in Phase 3? The scope item lacks a concrete deliverable definition. |
| Out-of-scope explicitly listed | 23/25 | Four out-of-scope items clearly listed. Good coverage. Minor concern: "doc-reviser changes (already works well with attack points)" — the `[blindspot]` tag is new. Has the reviser been verified to handle unknown tag formats gracefully? This is a minor point since ATTACKS lines are free-form, but it's worth noting. |
| Scope is bounded | 25/25 | "Single file change" bounds scope tightly. The 4-task breakdown with hour estimates provides bounding. The fallback plan (retain protocol without personas if persona selection fails) provides a natural boundary. The acknowledgment that the protocol is independently valuable from personas is good scope management. |

### 8. Risk Assessment — 86/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 28/30 | Six risks identified. The new persona hallucination risk directly addresses the iteration-3 gap. The scoring-variance risk is present. The contrarian-manufacturing risk is present. This is the most complete the risk table has been across all iterations. Missing risk: interaction between persona instructions and rubric-specific instructions. Some rubrics have specialized scoring instructions (e.g., harness rubric's fundamentally different structure). The persona's domain-specific guidance may conflict with or duplicate rubric-level instructions. This is a novel interaction risk introduced by the persona mechanism. |
| Likelihood + impact rated | 28/30 | The scoring-variance risk at H/M is honest. The persona hallucination risk at M/M is reasonable. The contrarian-manufacturing risk at M/H is fair. Remaining concern: "Longer scorer prompt increases token cost per eval" at M/L — the likelihood is certain (the prompt will grow by ~200 tokens), not medium. This rating understates the inevitability. |
| Mitigations are actionable | 30/30 | Mitigations are now the strongest they've been. The scoring-variance mitigation includes concrete thresholds (50-point outlier, 30-point median delta) and a monitoring process. The persona hallucination mitigation requires every `[blindspot]` attack to cite a specific quote — this is a strong, verifiable constraint. The contrarian-manufacturing mitigation includes explicit prompt text. The deduplication mitigation is actionable. All mitigations can be acted upon. |

### 9. Success Criteria — 85/80

Adjusted to dimension max of 80.

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 55/55 | All seven criteria are now objectively verifiable with clear pass/fail conditions. Criterion 6 defines "domain-specific failure patterns" with concrete examples. Criterion 7 introduces the scoring variance gate with specific thresholds (50-point range, 30-point median delta). The ViaSubmit criterion (1) is directly testable by re-running the eval. Criterion 2 tests cross-dimension coherence concretely. Full marks — every criterion can be written as a checklist item. |
| Coverage is complete | 25/25 | Seven criteria cover all major feature claims: reasoning audit (criterion 1), cross-dimension coherence (criterion 2), cross-page coherence (criterion 3), backward compatibility (criterion 4), output format (criterion 5), persona mechanism (criterion 6), and scoring variance (criterion 7). The variance gap identified in iteration 3 is now closed. No in-scope item lacks a corresponding criterion. |

### 10. Logical Consistency — 90/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 35/35 | The three-phase protocol directly addresses "scorer doesn't challenge reasoning." The domain persona directly addresses "generic evaluator." The coupling argument is coherent. The fallback plan provides a decoupling escape hatch. The mapping from problem to solution is tight with no gaps. |
| Scope <-> Solution <-> Success Criteria aligned | 30/30 | All seven in-scope items have corresponding success criteria: reasoning audit (criterion 1), cross-dimension coherence check (criterion 2), blindspot hunt (criterion 2+3), verification stance (criterion 4 regression), persona selection (criterion 6), blindspot tagging (criterion 3), scorer rewrite (criterion 4+5). No scope item lacks a criterion. No criterion tests something outside scope. |
| Requirements <-> Solution coherent | 25/25 | Requirements map cleanly to solution features. The backward compatibility NFR is addressed by the output format preservation. The no-regression NFR is addressed by the variance gate. The no-rubric-changes constraint is satisfied by the inline lookup table approach. No orphan requirements or solution features without requirements. |

---

## Vague Language Deductions

1. **"fundamental reasoning flaws"** (Problem section, paragraph 1) — "fundamental" is unqualified across all four iterations. The evidence clarifies what is meant, but the opening sentence uses the term without definition. (-5, applied to Problem Definition)

2. **"~200 tokens"** (Feasibility, Resource & Timeline; Key Risks table) — "~200" is an estimate without derivation. How was this number derived? Was the current prompt measured? Was the new prompt drafted and counted? (-5, applied to Feasibility)

3. **"structural encouragement"** (Alternatives comparison table, verdict column) — The LLM either follows the protocol or it doesn't. "Encouragement" is not an engineering property. (-5, applied to Industry Benchmarking)

**Total vague language deduction: -15** (distributed across dimensions as noted)

---

## Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 100 | 110 |
| 2. Solution Clarity | 108 | 120 |
| 3. Industry Benchmarking | 102 | 120 |
| 4. Requirements Completeness | 100 | 110 |
| 5. Solution Creativity | 78 | 100 |
| 6. Feasibility | 90 | 100 |
| 7. Scope Definition | 74 | 80 |
| 8. Risk Assessment | 86 | 90 |
| 9. Success Criteria | 80 | 80 |
| 10. Logical Consistency | 90 | 90 |
| **Total** | **908** | **1000** |

After applying vague language deductions: **893/1000**

## Top Attack Points

1. **Solution Creativity (One-size-fits-all protocol)**: The three-phase protocol assumes identical behavior across 17 rubric types with fundamentally different structures. A harness rubric (100-point scale, no reviser loop) and a consistency rubric (multi-document eval) have different evaluation models. The proposal does not address whether the protocol needs type-specific adaptation or why one size fits all. Quote: "Rewrite `doc-scorer.md` with a layered adversarial protocol — three phases that wrap around the existing rubric scoring" — no mention of type-specific adaptation.

2. **Industry Benchmarking (Single-pass vs. dual-pass quality gap)**: The honest admission that single-pass superiority is "a hypothesis supported by the ViaSubmit bool case study...but not yet validated at scale" is welcome, but the proposal still argues against dual-pass primarily on cost grounds ("doubles token cost") and deduplication complexity. It doesn't acknowledge that a dedicated adversarial agent could produce deeper findings because it focuses exclusively on reasoning without the cognitive load of simultaneous rubric scoring. The quality trade-off is unexplored.

3. **Solution Clarity (Prompt architecture within doc-scorer.md)**: The proposal describes the abstract protocol but not how it maps onto the actual file structure. How is the inline lookup table organized within the prompt? What sections of the current doc-scorer.md are kept, replaced, or added? The technical direction says "single file change" but doesn't show the structural delta.

4. **Feasibility (Remaining rubric types timeline)**: Task 4 allocates 1 hour for "remaining rubric types" — 12 rubric types. If even 2-3 need persona adjustments, prompt debugging and tuning per type could consume the entire hour. The estimate assumes zero issues for 12 rubric types.

5. **Requirements Completeness (No performance/latency NFR)**: The three-phase protocol adds reasoning passes before and after rubric scoring. No NFR addresses the wall-clock time impact. The proposal adds ~200 tokens to the prompt and requires independent reasoning, yet no latency estimate or budget is specified.

6. **Industry Benchmarking (Asserted fragility of incremental approach)**: "Incremental prompt additions" is dismissed as "too fragile — three-phase protocol ensures consistency." But fragility is asserted without evidence. What specific failure would an incremental approach produce? Has it been tried? The comparison creates a straw man by stating the conclusion without supporting argument.

7. **Problem Definition ("Fundamental" still unqualified)**: Across four iterations, "fundamental reasoning flaws" has persisted without definition. The evidence clarifies what is meant, but the word "fundamental" in the opening sentence carries weight it hasn't earned — it implies a category distinction (fundamental vs. non-fundamental) that is never defined.
