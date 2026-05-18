---
created: 2026-05-18
author: fanhuifeng
status: Draft
---

# Proposal: Adversarial Scorer with Domain Expert Personas

## Problem

The doc-scorer agent evaluates documents by mechanically applying rubric dimensions — checking each box, deducting points for gaps — but never independently challenging whether the document's core reasoning holds up. This produces scores of 900+ on documents that contain fundamental reasoning flaws (see the 958/1000 ViaSubmit bool example below).

### Evidence

1. **Disguised patches** (lesson: `gotcha-eval-rubric-misses-disguised-patches.md`): A refactoring design scored 958/1000 but reintroduced the pattern it aimed to eliminate (`ViaSubmit bool` — a disguised 5th path). The rubric measured document quality, not design quality.
2. **Page gap problem** (`arch-forge-skill-gap-analysis.md`): UI prototypes had systemic cross-page issues (inconsistent navigation, broken links, code layer violations). Each page looked fine in isolation — the scorer never checked the spaces between pages.
3. **Generic evaluator across domains**: The scorer uses the same "harsh document evaluator" persona whether it's reviewing a CTO-level proposal or a QA engineer's test cases. Domain expertise — knowing where real-world proposals actually fail, where production architectures actually break — is absent.

### Urgency

Every eval pass that scores high but misses fundamental flaws erodes trust in the eval pipeline. If eval can't catch the issues identified in the gap analysis, downstream skills (write-prd, ui-design, tech-design) will continue producing documents that look good on paper but fail in practice.

### Why These Two Problems Must Be Solved Together

"Lack of adversarial depth" (the scorer never challenges reasoning) and "lack of domain expertise" (the scorer lacks real-world failure intuition) are separable in theory but coupled in practice. Adversarial depth without domain expertise produces attacks that are adversarial but shallow — the scorer challenges the document aggressively but doesn't know where real-world failures cluster. Domain expertise without adversarial depth produces persona-dressed compliance — the scorer speaks like an expert but still mechanically checks boxes. The ViaSubmit bool case demonstrates both failures simultaneously: no adversarial challenge to the reasoning chain (depth), and no architectural expertise to recognize a disguised code path (domain). Solving one without the other would leave the core gap partially unaddressed. The coupling risk is acknowledged: if the three-phase protocol proves effective but persona selection doesn't work across all rubric types, the protocol can be retained while falling back to a generic expert persona.

## Proposed Solution

Rewrite `doc-scorer.md` with a **layered adversarial protocol** — three phases that wrap around the existing rubric scoring:

1. **Pre-scoring: Reasoning Audit** — Before touching the rubric, form an independent judgment: does the document's argument chain (problem → solution → evidence → success criteria) hold together? Does the solution reintroduce what it claims to eliminate?
2. **Rubric Scoring** — Existing dimension-by-dimension evaluation, but infused with a "verify, don't trust" stance. Every assertion is treated as unverified until proven.
3. **Post-scoring: Blindspot Hunt** — After all dimensions are scored, ask "what did the rubric miss?" Produce rubric-independent attack points.

Additionally, the scorer adopts a **domain expert persona** auto-selected from the rubric's `type` frontmatter field via an inline lookup table embedded in the scorer prompt. The table maps each `type` value to a persona definition (role description + domain-specific failure patterns to watch for). When a type is not in the table, the scorer falls back to a generic "Senior Technical Reviewer" persona. This lookup happens at the top of the scorer's execution, before any reading or scoring — the persona shapes how the scorer reads the document from the very beginning.

### Prompt Structure

The revised `doc-scorer.md` prompt will have three sections:

**Phase 1: Reasoning Audit** — Read end-to-end before rubric. Trace the argument chain (problem → solution → evidence). Check if the solution reintroduces what it claims to eliminate. Identify unstated assumptions. Record findings as pre-score anchors independent of any dimension.

**Phase 2: Rubric Scoring** — Existing dimension-by-dimension evaluation with added stance: "verify every claim against evidence; treat unverified assertions as gaps."

**Phase 3: Blindspot Hunt** — After scoring, look for: cross-section contradictions no dimension covers, missing failure modes the domain expert would expect, structural issues (scope/solution/criteria alignment). Tag each as `[blindspot]` in ATTACKS — strictly for issues outside all rubric dimensions.

### Before/After Example

**Before** (current scorer evaluating a refactoring proposal that re-introduces the pattern it eliminates):

```
ATTACKS:
- [clarity] Section 3.2 could be more specific about the migration strategy (-15)
- [feasibility] Resource estimate lacks breakdown by component (-10)
Total: 958/1000
```

The scorer scores each rubric dimension well, producing a high score. It never notices that the design re-introduces `ViaSubmit bool` — a disguised 5th path — because no rubric dimension asks "does this design achieve its stated goal?"

**After** (adversarial scorer with domain expert persona):

```
ATTACKS:
- [clarity] Section 3.2 could be more specific about the migration strategy (-15)
- [feasibility] Resource estimate lacks breakdown by component (-10)
- [blindspot] The proposal claims to eliminate all submit paths except the 4 canonical ones,
  but section 4.3 introduces a ViaSubmit bool flag that acts as a disguised 5th path.
  This directly contradicts the proposal's stated goal. The rubric has no dimension for
  "solution achieves its stated objective" — this is a reasoning failure, not a
  documentation failure. (-40)
Total: 870/1000
```

The same rubric dimensions score similarly, but the blindspot hunt catches the core reasoning flaw that no dimension covers.

### Persona Output Example

The same document (a proposal for a caching layer) evaluated against different rubric types produces different attack profiles:

**Generic scorer (no persona):**
```
ATTACKS:
- [clarity] Cache invalidation strategy is underspecified (-15)
- [feasibility] No latency benchmarks provided (-10)
Total: 940/1000
```

**"Proposal Expert" persona (evaluating as `type: proposal`):**
```
ATTACKS:
- [clarity] Cache invalidation strategy is underspecified (-15)
- [feasibility] No latency benchmarks provided (-10)
- [blindspot] The proposal claims 99.9% cache hit rate will "eliminate database load," but
  the invalidation policy (write-through on every mutation) means every write still hits the
  database synchronously. The value proposition is overstated — the cache reduces read load,
  not write load. A stakeholder deciding on funding based on "eliminates DB load" would be
  misled. (-35)
- [blindspot] No rollback plan if the cache layer introduces consistency issues in production.
  For an infrastructure change touching all read paths, this is a critical omission in the
  risk section. (-20)
Total: 870/1000
```

**"Senior QA Engineer" persona (same document evaluated as `type: test-cases`):**
```
ATTACKS:
- [blindspot] Cache stampede scenario not covered — if the cache cold-starts under load,
  every request bypasses to the database simultaneously. No circuit breaker or
  request-coalescing mentioned. (-25)
- [blindspot] The proposal assumes a single cache cluster. What happens during a cache
  partition (split-brain)? Two application nodes could see different data with no
  reconciliation strategy. (-30)
Total: 885/1000
```

The Proposal Expert attacks stakeholder-facing concerns (value proposition accuracy, risk coverage). The QA Engineer attacks operational failure modes (stampede, split-brain). Both catch issues the generic scorer misses, but from different angles.

### Innovation Highlights

The core insight is **compositional rather than conceptual**: "reason first, score second" is what competent human reviewers do by default. What is new is engineering this behavior into a structured, reproducible prompt protocol that runs every time, not only when a human reviewer happens to think of it. The cited industry tools (Anthropic model evaluations, OpenAI evals) are architected as dimension-by-dimension scorers — each criterion is graded independently, and the final score is an aggregation. Neither framework includes an independent reasoning-pass before or after the rubric. PyRIT and red-teaming frameworks are purely adversarial with no scoring structure. The novelty here is the composition: a three-phase protocol that enforces the adversarial reasoning human reviewers naturally apply but automated evaluators currently skip.

The domain expert persona mechanism borrows from multiple domains where expertise determines detection quality:

- **Red-team security**: A crypto specialist and a web-app specialist running the same tool find different bugs. Similarly, a "Proposal Expert" and a "QA Engineer" running the same rubric surface different failure modes.
- **Academic peer review**: Reviewers don't just check whether a paper follows journal format — they evaluate whether the methodology supports the conclusions. A reviewer who only checks formatting would give high marks to a well-formatted paper with flawed methodology. This is precisely the current scorer's failure mode.
- **QA equivalence partitioning**: Expert QA engineers don't follow test plans mechanically — they intuitively know which boundary conditions matter for their specific domain (cache invalidation for distributed systems, race conditions for concurrent code, input sanitization for web APIs). The persona mechanism encodes this domain-specific boundary knowledge.

The composition ceiling — combining known techniques rather than inventing new ones — is acknowledged. Genuinely novel evaluation approaches do exist (e.g., formal verification of document properties, or training a specialized critic model), but these require infrastructure investment (model training, formal language design) that is disproportionate to the problem. The three-phase protocol achieves 80% of the benefit at 5% of the cost.

## Requirements Analysis

### Key Scenarios

**Flaw-detection scenarios (the scorer should catch these):**

- Evaluating a refactoring proposal that claims to simplify but subtly reintroduces the old pattern → scorer catches the disguised patch
- Evaluating a UI design where each page is well-specified but navigation between pages is inconsistent → scorer catches cross-page coherence gaps
- Evaluating a tech design where error handling "looks complete" but doesn't cover the actual failure modes → scorer independently traces error paths
- Evaluating test cases where steps are detailed but not actually executable by a downstream agent → scorer catches actionability gaps

**False-positive scenarios (the scorer should NOT flag these):**

- Well-written document with genuine full marks → scorer awards full score without manufacturing flaws. Reasoning audit concludes "sound", blindspot hunt finds nothing.
- Issue already captured by a rubric dimension → scored in that dimension, NOT duplicated as a `[blindspot]` attack.
- Rubric type without a specific persona mapping → fallback to "Senior Technical Reviewer", produces valid output.

**Edge-case scenarios:**

- Rubric type that doesn't exist in the persona mapping → fallback persona activates, eval completes normally
- Document that is mostly sound but has one subtle cross-section contradiction → blindspot hunt catches the contradiction even though each section individually scores well
- Very long document (multiple eval iterations with compaction) → adversarial protocol still activates even when context is compacted, because the prompt is in the system instructions, not the conversation history
- Phase 1 (reasoning audit) and Phase 2 (rubric scoring) produce conflicting assessments → Phase 2 rubric scores take precedence for `DIMENSIONS:` output; Phase 1 findings are channeled into `[blindspot]` attacks only if they identify issues not covered by any dimension. If Phase 1 flags a fundamental flaw but Phase 2 scores that dimension well, the Phase 1 finding appears as a `[blindspot]` attack with explicit notation: "Reasoning audit flagged this independently of dimension scoring."

### Non-Functional Requirements

- **Backward compatibility**: Output format (`SCORE:`, `DIMENSIONS:`, `ATTACKS:`) must remain parseable by the existing orchestrator
- **No rubric changes**: The scorer enhancement must work with all 17 existing rubric files without modification
- **No regression**: Existing high-quality documents should still score high; the adversarial stance should catch real flaws, not manufacture them

### Constraints & Dependencies

- The scorer reads the rubric's `type` frontmatter field to select persona — this field already exists in all rubrics
- The scorer is a subagent invoked via the Agent tool — no orchestrator changes needed
- The reviser already consumes `ATTACKS:` lines — `[blindspot]` tagged attacks integrate naturally

## Alternatives & Industry Benchmarking

### Industry Solutions

Document evaluation in AI tools falls into two camps:

1. **Rubric-based scoring with fixed criteria.** Anthropic's [model evaluations](https://docs.anthropic.com/en/docs/build-with-claude/evaluations) use graded rubrics mapping criteria to numeric scores. OpenAI's [eval framework](https://github.com/openai/evals) defines `Eval` classes with grading functions comparing output to expected results. Both produce reproducible scores but cannot surface issues outside their predefined criteria.

2. **Free-form adversarial red-teaming.** Microsoft's [PyRIT](https://github.com/Azure/PyRIT) automates red-team attacks via composable strategies. Anthropic's [red-team methodology](https://www.anthropic.com/news/announcing-our-responsible-scaling-policy) uses human red-teamers with domain expertise — the expertise determines what they find. Google's [ASCI](https://github.com/google/ASCI) targets adversarial testing of AI agents.

Neither camp combines structured rubric scoring with adversarial reasoning in a single pass. PyRIT and red-teaming produce unstructured findings that are hard to iterate against; rubric scoring produces structured scores that miss reasoning flaws. This proposal fills that gap. Note: the claim that combining both in a single pass outperforms either alone is a hypothesis supported by the ViaSubmit bool case study (rubric-only scored 958/1000 and missed the disguised patch) but not yet validated at scale. The regression validation step (Task 3) serves as the empirical test.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | No cost, zero regression risk | Misses disguised patches, page gaps; 3+ real flaws passed undetected | Rejected: evidence proves misses |
| Separate devil's advocate agent | PyRIT composable strategies; Anthropic red-team | Full independence from rubric scoring; can run in parallel | Orchestrator changes (new agent, invocation step, output merge). ~3-5s slower per eval. Doubles token cost. No shared context between agents — adversarial agent can't see what scorer already caught, so findings may duplicate. | Viable but coordination problem (two agents, deduplication) is the real blocker |
| Full scorer rewrite | Clean-slate design | Could restructure to claim-evidence-verification model; eliminates prompt cruft | All 17 rubric types re-validated; past eval results incomparable; 2-3x effort of in-place enhancement | Rejected: prompt structure isn't the problem — missing adversarial reasoning is |
| Incremental prompt additions | Common prompt engineering approach | Simplest change — add a few "also check reasoning" lines | No structural encouragement that adversarial reasoning runs; LLM may treat it as optional; scattered instructions harder to maintain | Too fragile — three-phase protocol ensures consistency |
| **Layered protocol in scorer** | PyRIT composable attacks + rubric hybrid | Structured phases; backward-compatible; single file; encourages adversarial reasoning by making it a named phase | Prompt grows ~200 tokens; poorly written persona could degrade quality; must test across 17 rubric types; no software-enforced guarantee — relies on LLM following prompt structure | **Selected: structural encouragement with minimal integration cost** |

## Feasibility Assessment

### Technical Feasibility

Single file change (`doc-scorer.md`). The scorer already reads the rubric's `type` field. The output format is already extensible (ATTACKS lines are free-form). No new dependencies.

### Resource & Timeline

Estimated **3-4 tasks**, not 1:

1. **Rewrite `doc-scorer.md`** — draft the three-phase prompt structure and persona selection logic (~1-2 hours of prompt engineering)
2. **Write persona definitions** — define domain expert personas for the most common rubric types (proposal, prd, design, ui, test-cases, harness) plus the fallback "Senior Technical Reviewer" (~30 min)
3. **Regression validation** — re-run at least 5 existing evals covering both standard and non-standard rubric types: one proposal (1000-point), one design (1000-point), one test-cases (1000-point), one harness (100-point scale, no reviser loop), and one consistency (multi-document eval). Compare output format and scores, fix any breakage (~2-3 hours)
4. **Remaining rubric types** — verify the scorer works with the other 12 rubric types, adjust persona mapping as needed (~1 hour)

Total estimated effort: 5-7 hours. No external dependencies block start.

### Dependency Readiness

No external dependencies. The rubric `type` field exists in all 17 rubrics.

## Scope

### In Scope

- Rewrite `plugins/forge/agents/doc-scorer.md` with layered adversarial protocol
- Add pre-scoring reasoning audit phase
- Add post-scoring blindspot hunt phase
- Add cross-dimension coherence check
- Add verification stance as explicit Phase 2 instructions: "treat every assertion as unverified until the document provides supporting evidence; flag assertions that lack evidence as gaps in the corresponding rubric dimension"
- Add domain expert persona auto-selection from rubric `type`
- Tag blindspot attacks with `[blindspot]` for reviser consumption

### Out of Scope

- Rubric file changes (all 17 rubrics stay as-is)
- eval SKILL.md orchestrator changes
- doc-reviser changes (already works well with attack points)
- New eval types or rubric additions

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Scorer becomes a contrarian, manufacturing flaws | M | H | Explicit instruction: "Find REAL issues. A document with no flaws deserves full marks. Manufacturing issues wastes the reviser's time and yours." |
| Domain personas don't match all 17 rubric types | L | M | Provide a fallback "Senior Technical Reviewer" persona for types without a specific expert mapping |
| Blindspot attacks duplicate dimension attacks | M | L | Explicit instruction: "Blindspot attacks must identify issues OUTSIDE any rubric dimension. If an issue fits a dimension, score it there instead." |
| Longer scorer prompt increases token cost per eval | M | L | The prompt grows by ~200 tokens. Negligible compared to document reading cost. |
| Scoring variance across runs reduces reproducibility | H | M | Adversarial framing encourages finding issues, which may produce different scores for the same document across runs. Mitigation: run each eval 3 times and accept the median score; flag any run where score delta exceeds 50 points as an outlier requiring manual review. Track variance over time — if median delta exceeds 30 points across documents, tune the adversarial prompt to reduce randomness. |
| Persona hallucination — LLM fabricates plausible domain-specific concerns that don't apply | M | M | The persona instructs the scorer to look for specific failure patterns, but the LLM may fabricate plausible-sounding domain expertise. Mitigation: every `[blindspot]` attack must cite a specific quote from the document; blindspot attacks without quotes are discarded. Cross-reference persona-specific claims against the actual document content. |

## Success Criteria

- [ ] Re-evaluating the design doc that had `ViaSubmit bool` produces an attack point flagging the disguised patch
- [ ] Evaluating a document where scope claims X but success criteria only test Y produces a `[blindspot]` attack flagging the cross-dimension gap (covers the cross-dimension coherence check)
- [ ] Re-evaluating the UI design with page gap issues produces cross-page coherence attack points
- [ ] All existing eval types (proposal, prd, design, ui-*, test-cases-*, consistency, harness) work without rubric changes
- [ ] Output format remains `SCORE:`, `DIMENSIONS:`, `ATTACKS:` — parseable by existing orchestrator
- [ ] Persona adoption verified: run the same document against `proposal` and `design` rubrics; confirm (a) attacks reference domain-specific failure patterns (defined as: attacks that name a failure mode specific to the document's domain, such as "stakeholder value proposition overstated" for proposals or "cache stampede under cold start" for designs — contrasted with generic patterns like "section lacks detail" that apply to any document type), and (b) at least 2 attacks differ between runs attributable to perspective
- [ ] Scoring variance gate: for any document, 3 consecutive eval runs produce scores within a 50-point range (median delta < 30 points). If variance exceeds this threshold, the adversarial prompt must be tuned before shipping.

## Next Steps

- Proceed to task generation via `/breakdown-tasks` (estimated 3-4 tasks, ~5-7 hours total effort)
