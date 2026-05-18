---
created: 2026-05-18
author: fanhuifeng
status: Draft
---

# Proposal: Adversarial Scorer with Domain Expert Personas

## Problem

The doc-scorer agent evaluates documents by mechanically applying rubric dimensions — checking each box, deducting points for gaps — but never independently challenging whether the document's core reasoning holds up. This produces high scores with low insight.

### Evidence

1. **Disguised patches** (lesson: `gotcha-eval-rubric-misses-disguised-patches.md`): A refactoring design scored 958/1000 but reintroduced the pattern it aimed to eliminate (`ViaSubmit bool` — a disguised 5th path). The rubric measured document quality, not design quality.
2. **Page gap problem** (`arch-forge-skill-gap-analysis.md`): UI prototypes had systemic cross-page issues (inconsistent navigation, broken links, code layer violations). Each page looked fine in isolation — the scorer never checked the spaces between pages.
3. **Generic evaluator across domains**: The scorer uses the same "harsh document evaluator" persona whether it's reviewing a CTO-level proposal or a QA engineer's test cases. Domain expertise — knowing where real-world proposals actually fail, where production architectures actually break — is absent.

### Urgency

Every eval pass that scores high but misses fundamental flaws erodes trust in the eval pipeline. If eval can't catch the issues identified in the gap analysis, downstream skills (write-prd, ui-design, tech-design) will continue producing documents that look good on paper but fail in practice.

## Proposed Solution

Rewrite `doc-scorer.md` with a **layered adversarial protocol** — three phases that wrap around the existing rubric scoring:

1. **Pre-scoring: Reasoning Audit** — Before touching the rubric, form an independent judgment of the document's fundamental soundness. Does the argument hold up? Does the solution reintroduce what it claims to eliminate?
2. **Rubric Scoring** — Existing dimension-by-dimension evaluation, but infused with a "verify, don't trust" stance. Every assertion is treated as unverified until proven.
3. **Post-scoring: Blindspot Hunt** — After all dimensions are scored, ask "what did the rubric miss?" Produce rubric-independent attack points.

Additionally, the scorer adopts a **domain expert persona** auto-selected from the rubric's `type` frontmatter field. Each persona brings real-world failure intuition that a generic evaluator lacks.

### Innovation Highlights

The key insight is **adversarial reasoning as the scorer's primary job, not an add-on**. Current scorers in document evaluation tools (including AI code review tools) score first and explain second. This design flips the order: reason first, score second, then hunt for what you missed.

The domain expert persona mechanism is borrowed from red-team security practices — where the attacker's expertise determines which vulnerabilities they find, not just their adversarial intent.

## Requirements Analysis

### Key Scenarios

- Evaluating a refactoring proposal that claims to simplify but subtly reintroduces the old pattern → scorer catches the disguised patch
- Evaluating a UI design where each page is well-specified but navigation between pages is inconsistent → scorer catches cross-page coherence gaps
- Evaluating a tech design where error handling "looks complete" but doesn't cover the actual failure modes → scorer independently traces error paths
- Evaluating test cases where steps are detailed but not actually executable by a downstream agent → scorer catches actionability gaps

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

Document evaluation in AI tools typically uses either: (1) rubric-based scoring with fixed criteria (Anthropic's evals, OpenAI's eval framework), or (2) free-form adversarial red-teaming. Neither combines the two — rubric scoring provides structure but misses reasoning flaws; red-teaming catches reasoning flaws but is unstructured and hard to iterate.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | No implementation cost | Misses disguised patches, page gaps, coherence issues | Rejected: evidence shows current approach misses real flaws |
| Separate devil's advocate agent | Red-team practice | Full independence from rubric | Adds loop complexity, requires orchestrator changes, slower | Rejected: over-engineered for the problem |
| Full scorer rewrite | — | Clean slate | High regression risk, all 17 rubric types must be re-validated | Rejected: too risky |
| **Layered protocol in scorer** | Red-team + rubric hybrid | Structured adversarial reasoning, backward-compatible, single file change | Scorer prompt gets longer | **Selected: highest leverage, lowest risk** |

## Feasibility Assessment

### Technical Feasibility

Single file change (`doc-scorer.md`). The scorer already reads the rubric's `type` field. The output format is already extensible (ATTACKS lines are free-form). No new dependencies.

### Resource & Timeline

1 task: rewrite doc-scorer.md. Can be validated by re-running any existing eval (e.g., `/eval-proposal` on an existing proposal) and comparing output quality.

### Dependency Readiness

No external dependencies. The rubric `type` field exists in all 17 rubrics.

## Scope

### In Scope

- Rewrite `plugins/forge/agents/doc-scorer.md` with layered adversarial protocol
- Add pre-scoring reasoning audit phase
- Add post-scoring blindspot hunt phase
- Add cross-dimension coherence check
- Add "verify, don't trust" verification stance
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

## Success Criteria

- [ ] Re-evaluating the design doc that had `ViaSubmit bool` produces an attack point flagging the disguised patch
- [ ] Re-evaluating the UI design with page gap issues produces cross-page coherence attack points
- [ ] All existing eval types (proposal, prd, design, ui-*, test-cases-*, consistency, harness) work without rubric changes
- [ ] Output format remains `SCORE:`, `DIMENSIONS:`, `ATTACKS:` — parseable by existing orchestrator
- [ ] Scorer adopts different expert personas for different rubric types (visible in attack point language/perspective)

## Next Steps

- Proceed directly to task generation (feature is small enough for `/quick-tasks`)
