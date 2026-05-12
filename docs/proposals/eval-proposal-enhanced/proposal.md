# Design: Eval-Proposal Rubric Overhaul

**Date:** 2026-05-12
**Status:** Approved

---

## Summary

Overhaul the `eval-proposal` rubric from 6 dimensions / 100 points to 10 dimensions / 1000 points, addressing the root cause of "proposals score high but miss key points." Simultaneously update the proposal template to add new sections that support the expanded rubric.

---

## Problem

Proposals consistently pass eval-proposal (90+ / 100) but are later found to have missed key requirements, edge cases, and industry-standard solutions. Root cause analysis:

**Lack of divergent thinking → Insufficient research → Requirements gaps**

The current rubric checks "is it written clearly" but never checks "is it thought comprehensively." Specifically:

1. **Alternatives Analysis is toothless**: proposals list "do nothing" + one straw man, scorer gives high marks
2. **No completeness check**: no dimension verifies "did you miss any scenarios, constraints, or NFRs?"
3. **No feasibility gate**: unrealistic proposals pass through to PRD/tech-design, wasting downstream effort
4. **No creativity assessment**: proposals default to the first obvious solution without exploring innovative approaches
5. **Consistency is reactive only**: cross-section inconsistency exists as a -3 deduction, not a positive quality check

---

## Design Decisions

### Decision 1: Expand to 1000-point scale (absolute)

100-point scale doesn't provide enough granularity to distinguish between proposals of varying quality. 1000 points allow finer criterion differentiation and clearer gap analysis.

### Decision 2: Add 5 new/strengthened dimensions

- **Industry Benchmarking** (strengthened from Alternatives Analysis): force research of real-world solutions
- **Requirements Completeness** (new): explicit check for scenario coverage, NFRs, and constraints
- **Solution Creativity** (new): assess innovation beyond industry baseline
- **Feasibility** (new): gate unrealistic proposals before they waste downstream effort
- **Logical Consistency** (new, promoted from deduction rule): positive quality check for cross-section coherence

### Decision 3: Update proposal template

Add 3 new sections + 1 sub-section to provide structure for the new dimensions:

- `## Requirements Analysis` — supports Requirements Completeness
- `## Alternatives & Industry Benchmarking` — supports Industry Benchmarking (expands existing Alternatives Considered)
- `## Feasibility Assessment` — supports Feasibility
- `### Innovation Highlights` (under Proposed Solution) — supports Solution Creativity

### Decision 4: Migrate all eval skills to 1000-point scale (Phase 2)

Shared `doc-scorer` and `doc-reviser` agents require a uniform scoring model. Phase 1 implements eval-proposal; Phase 2 migrates eval-prd, eval-design, eval-ui, eval-test-cases, eval-consistency.

---

## Scoring Model

### eval-proposal (1000 pts)

| # | Dimension | Pts | Priority |
|---|-----------|-----|----------|
| 1 | Problem Definition | 110 | Backbone |
| 2 | Solution Clarity | 120 | Backbone |
| 3 | Industry Benchmarking | 160 | Root cause fix (input) |
| 4 | Requirements Completeness | 140 | Pain point fix |
| 5 | Solution Creativity | 130 | Root cause fix (output) |
| 6 | Feasibility | 100 | Waste prevention |
| 7 | Scope Definition | 80 | Necessary |
| 8 | Risk Assessment | 90 | Necessary |
| 9 | Success Criteria | 80 | Necessary |
| 10 | Logical Consistency | 90 | Safety net |

Priority rationale: Benchmarking > Completeness > Creativity because insufficient research is the root cause — if you research broadly, you discover missing requirements naturally.

### Default Parameters

| Parameter | Old | New |
|-----------|-----|-----|
| `--target` | 90 | 900 |
| `--iterations` | 3 | 3 (unchanged) |

---

## Complete Rubric

```markdown
# Proposal Evaluation Rubric

**Total: 1000 points**
**Report template:** `plugins/forge/skills/eval-proposal/templates/report.md`

## Dimensions

### 1. Problem Definition (110 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Problem stated clearly | 0-40 | Is the core problem unambiguous? Could two readers interpret it differently? |
| Evidence provided | 0-40 | Is there data, user feedback, or concrete examples backing the problem? Not just "we think..." |
| Urgency justified | 0-30 | Why solve this now? What happens if we don't? What's the cost of delay? |

### 2. Solution Clarity (120 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Approach is concrete | 0-40 | Can a reader explain back what will be built? Or is it vague hand-waving? |
| User-facing behavior described | 0-45 | What does the end user experience? Not internals — the observable behavior |
| Technical direction clear | 0-35 | Is there enough technical hint to know the general implementation approach? |

### 3. Industry Benchmarking (160 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Industry solutions referenced | 0-50 | Are real-world solutions/patterns for this type of problem cited? Product names, open-source projects, or published patterns? Not just self-invented options |
| At least 3 meaningful alternatives | 0-40 | Including "do nothing". Each alternative must be a genuinely different approach, not a straw man. At least one must be an industry-validated solution |
| Honest trade-off comparison | 0-35 | Are pros/cons based on actual project constraints? Not cherry-picked? |
| Chosen approach justified against benchmarks | 0-35 | Why this approach over industry standards? Or why adopt the standard? Clear rationale required |

### 4. Requirements Completeness (140 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Scenario coverage | 0-50 | Happy path + edge cases + error scenarios all identified? Or only the golden path? |
| Non-functional requirements | 0-45 | Performance, security, compatibility, accessibility — are relevant NFRs called out? |
| Constraints & dependencies | 0-45 | External systems, tech constraints, prerequisite conditions named? |

### 5. Solution Creativity (130 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Novelty over industry baseline | 0-50 | Is the proposal just copying an industry solution, or does it innovate beyond? Is the differentiation from the benchmark clearly articulated? |
| Cross-domain inspiration | 0-40 | Does it borrow ideas from other domains/products that face similar problems? |
| Simplicity of insight | 0-40 | Is the creative leap elegant ("why didn't I think of that") or forced/overengineered? |

### 6. Feasibility (100 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Technical feasibility | 0-40 | Can current tech stack support this? Any showstopper dependencies? |
| Resource & timeline feasibility | 0-30 | Does the team have skills/bandwidth? Is the scope realistic for the proposed timeline? |
| Dependency readiness | 0-30 | Are external APIs/services available? Are upstream prerequisites in place? |

### 7. Scope Definition (80 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| In-scope items are concrete | 0-30 | Each item is a deliverable, not a vague area |
| Out-of-scope explicitly listed | 0-25 | Are deferred items named, not just implied? |
| Scope is bounded | 0-25 | Can a team execute this in a defined timeframe? Or is it open-ended? |

### 8. Risk Assessment (90 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Risks identified | 0-30 | At least 3 meaningful risks, not trivial ones |
| Likelihood + impact rated | 0-30 | Is the assessment honest? Not all "low likelihood, high impact"? |
| Mitigations are actionable | 0-30 | Can someone act on the mitigation? Or is it "we'll handle it"? |

### 9. Success Criteria (80 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Criteria are measurable | 0-30 | Can you objectively verify each criterion? "Works well" is not measurable |
| Coverage is complete | 0-25 | Do criteria cover all in-scope items? Any gaps? |
| Criteria are testable | 0-25 | Could you write a test or checklist for each criterion? |

### 10. Logical Consistency (90 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Solution addresses the stated problem | 0-35 | Does the proposed solution actually solve the problem defined in section 1? Or is there a gap? |
| Scope ↔ Solution ↔ Success Criteria aligned | 0-30 | No contradictions between what we're building, what's in scope, and how we measure success |
| Requirements ↔ Solution coherent | 0-25 | Do the identified requirements map cleanly to the proposed solution? No orphan requirements or solution features with no requirement |

## Section → Dimension Mapping

| Section | Primary Dimension(s) |
|---------|---------------------|
| Problem (Evidence, Urgency) | 1. Problem Definition |
| Proposed Solution | 2. Solution Clarity |
| Proposed Solution > Innovation Highlights | 5. Solution Creativity |
| Requirements Analysis | 4. Requirements Completeness |
| Alternatives & Industry Benchmarking | 3. Industry Benchmarking |
| Feasibility Assessment | 6. Feasibility |
| Scope (In/Out) | 7. Scope Definition |
| Key Risks | 8. Risk Assessment |
| Success Criteria | 9. Success Criteria |
| *(cross-section)* | 10. Logical Consistency |

## Required Sections

| Section | Required |
|---------|----------|
| Problem (Evidence, Urgency) | ✓ |
| Proposed Solution | ✓ |
| Proposed Solution > Innovation Highlights | ✓ |
| Requirements Analysis | ✓ |
| Alternatives & Industry Benchmarking | ✓ |
| Feasibility Assessment | ✓ |
| Scope (In/Out) | ✓ |
| Key Risks | ✓ |
| Success Criteria | ✓ |

## Deduction Rules

- **Missing required section**: 0 pts for that dimension
- **Vague language without quantification**: -20 pts per instance ("better", "improved", "enhanced", "faster performance")
- **Cross-section inconsistency**: -30 pts per conflict (e.g., scope contradicts solution, success criteria don't cover scope)
- **Placeholder text ("TBD", "TODO")**: -20 pts per instance
- **Straw-man alternative**: -20 pts per instance (an alternative presented only to be rejected)
```

---

## Proposal Template (Updated)

```markdown
---
created: YYYY-MM-DD
author: "<!-- who proposed this -->"
status: Draft
---

# Proposal: {{PROPOSAL_TITLE}}

## Problem

<!-- State the core problem in one unambiguous sentence. Two readers should interpret it the same way. -->

### Evidence
<!-- Data, user feedback, or concrete examples proving this problem exists. Avoid "we think..." — cite specifics. -->

### Urgency
<!-- Why solve this now? What happens if we don't? What's the cost of delay? -->

## Proposed Solution

<!-- High-level approach. What will be built? What does the end user experience? -->

### Innovation Highlights

<!-- How does this solution differ from or improve upon industry-standard approaches?
What's the creative insight? What ideas are borrowed from other domains?
Keep it honest — if the approach is straightforward adoption of a standard, say so and explain why. -->

## Requirements Analysis

### Key Scenarios
<!-- List the main user scenarios this feature must handle:
- Happy path scenarios
- Edge cases and boundary conditions
- Error scenarios and failure modes
-->

### Non-Functional Requirements
<!-- Relevant NFRs: performance, security, compatibility, accessibility, scalability, etc.
Skip NFRs that are clearly not relevant — don't pad. -->

### Constraints & Dependencies
<!-- External systems this depends on
Technical constraints (platform, language, framework)
Prerequisite conditions that must be met
-->

## Alternatives & Industry Benchmarking

### Industry Solutions

<!-- How is this type of problem typically solved in the industry?
Reference specific products, open-source projects, or published patterns. -->

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | <!-- cost of inaction --> | <!-- --> | Rejected: <!-- why --> |
| <!-- industry standard --> | <!-- product/pattern name --> | <!-- --> | <!-- --> | <!-- --> |
| <!-- alternative approach --> | <!-- source --> | <!-- --> | <!-- --> | <!-- --> |
| **Chosen approach** | <!-- --> | <!-- --> | <!-- --> | **Selected: <!-- why -->** |

## Feasibility Assessment

### Technical Feasibility
<!-- Can current tech stack support this? Any showstopper technical dependencies? -->

### Resource & Timeline
<!-- Does the team have the skills? Is the scope realistic for the timeline? -->

### Dependency Readiness
<!-- Are external APIs/services available and stable? Are upstream prerequisites in place? -->

## Scope

### In Scope
- <!-- Each item is a specific deliverable -->

### Out of Scope
- <!-- Explicitly list deferred items -->

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| <!-- --> | <!-- H/M/L --> | <!-- H/M/L --> | <!-- actionable mitigation --> |

## Success Criteria

- [ ] <!-- Each criterion must include a number, percentage, or observable behavior that can be objectively verified -->

## Next Steps

- Proceed to `/write-prd` to formalize requirements
```

---

## Report Template (Updated)

The report template must accommodate 10 dimensions instead of 6, and use 1000-point scale.

Key changes from current template:
- Scorecard table expands from 6 to 10 dimensions
- All scores use `/1000` instead of `/100`
- Target default displays as `/1000`
- Dimension breakdown in final report lists 10 dimensions

---

## SKILL.md Changes

| Section | Change |
|---------|--------|
| `description` frontmatter | Update "100-point scoring" → "1000-point scoring" |
| Parameters | `--target` default: 90 → 900, description: "(0-100)" → "(0-1000)" |
| Step 2 (Scorer) | Parse `SCORE: X/1000` instead of `/100` |
| Step 3 (Gate) | Update all `/100` references to `/1000` |
| Step 4 (Reviser) | No structural change |
| Step 5 (Final Report) | Dimension table lists 10 dimensions with new max values |

---

## File Changeset

| Action | Path | Notes |
|--------|------|-------|
| Update | `plugins/forge/skills/eval-proposal/templates/rubric.md` | Complete rewrite: 10 dimensions, 1000 pts |
| Update | `plugins/forge/skills/eval-proposal/templates/report.md` | 10-dimension scorecard, 1000-pt scale |
| Update | `plugins/forge/skills/eval-proposal/SKILL.md` | Parameters, parsing, dimension table |
| Update | `plugins/forge/skills/brainstorm/templates/proposal.md` | New sections: Requirements Analysis, Feasibility Assessment, Innovation Highlights, expanded Alternatives |

---

## Phase 2: All Eval Skills Migration (1000 pts)

After eval-proposal is validated, migrate remaining eval skills:

| Skill | Current | Migration |
|-------|---------|-----------|
| eval-prd | 100 pts, 5 dimensions | Scale criterion points ×10, keep dimensions |
| eval-design | 100 pts, 6 dimensions | Scale criterion points ×10, keep dimensions |
| eval-ui | 100 pts, 4 perspectives | Scale criterion points ×10, keep perspectives |
| eval-test-cases | 100 pts | Scale criterion points ×10 |
| eval-consistency | 100 pts | Scale criterion points ×10 |

All skills update `--target` defaults from 90 → 900 (or their respective current default ×10).

---

## Constraints

- `doc-scorer` and `doc-reviser` agents require no changes — they read rubrics at runtime
- Rubric file is plain markdown, no frontmatter required
- The scorer must never be told what the reviser changed (existing isolation rule)
- Proposal template changes are backward-compatible: existing proposals missing new sections score 0 on the corresponding dimensions (existing "missing required section" rule)
- `brainstorm` SKILL.md may need updates to guide authors toward new template sections
