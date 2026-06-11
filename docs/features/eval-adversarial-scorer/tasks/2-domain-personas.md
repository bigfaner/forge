---
id: "2"
title: "Add domain expert persona auto-selection"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "enhancement"
mainSession: false
---

# 2: Add domain expert persona auto-selection

## Description

Add domain expert persona auto-selection to the rewritten `doc-scorer.md`. The scorer reads the rubric's `type` frontmatter field and selects a persona from an inline lookup table embedded in the scorer prompt. Each persona definition includes a role description and domain-specific failure patterns to watch for. The persona shapes how the scorer reads the document from the very beginning — it is selected at the top of execution, before any reading or scoring.

When a rubric type is not in the lookup table, the scorer falls back to a generic "Senior Technical Reviewer" persona.

## Reference Files
- `docs/proposals/eval-adversarial-scorer/proposal.md` — Source proposal
- `plugins/forge/agents/doc-scorer.md` — Scorer being modified (rewritten in task 1)
- `docs/conventions/forge-distribution.md` — MUST read before modifying plugins/forge/ files

## Acceptance Criteria

- [ ] Persona auto-selected from rubric `type` frontmatter field — no manual configuration needed
- [ ] At minimum, personas defined for: `proposal`, `prd`, `design`, `ui`, `test-cases`, `harness` rubric types
- [ ] Fallback "Senior Technical Reviewer" persona activates when rubric type is not in the lookup table
- [ ] Persona adoption verified: run the same document against `proposal` and `design` rubrics; confirm (a) attacks reference domain-specific failure patterns, and (b) at least 2 attacks differ between runs attributable to perspective
- [ ] All existing eval types (proposal, prd, design, ui-*, test-cases-*, consistency, harness) work without rubric changes
- [ ] Scoring variance gate: for any document, 3 consecutive eval runs produce scores within a 50-point range (median delta < 30 points)

## Hard Rules

- MUST read `docs/conventions/forge-distribution.md` before modifying `plugins/forge/agents/doc-scorer.md`
- No rubric file changes — the scorer reads the existing `type` field that all 17 rubrics already have
- No orchestrator changes — persona selection is internal to the scorer prompt

## Implementation Notes

- Persona definitions are inline in the scorer prompt (lookup table), not external files — keeps the change self-contained
- Each persona should specify: role description, domain-specific failure patterns, what the expert would intuitively check that a generic reviewer would miss
- Persona hallucination risk: the LLM may fabricate plausible-sounding domain concerns. Mitigation: every `[blindspot]` attack must cite a specific quote from the document
- Key risk from proposal: "Domain personas don't match all 17 rubric types" — the fallback persona handles this
