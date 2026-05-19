---
status: "completed"
started: "2026-05-19 01:00"
completed: "2026-05-19 01:02"
time_spent: "~2m"
---

# Task Record: 1 Extract scorer and reviser protocols

## Summary
Extracted generic scoring and revision protocols from doc-scorer.md and doc-reviser.md into standalone protocol files under agents/experts/protocol/. Scorer protocol contains the three-phase adversarial workflow (Reasoning Audit, Rubric Scoring, Blindspot Hunt) without persona/domain-specific content. Reviser protocol contains attack-point-driven revision workflow without rubric path input.

## Changes

### Files Created
- plugins/forge/agents/experts/protocol/scorer-protocol.md
- plugins/forge/agents/experts/protocol/reviser-protocol.md

### Files Modified
无

### Key Decisions
- Scorer protocol excludes the entire Persona Selection section and domain-specific failure patterns — those belong in expert files (Task 2)
- Reviser protocol omits RUBRIC_PATH entirely — attack points already contain domain-informed prescriptions; structural issues caught by scorer
- Both protocol files use template variables ({{DOC_DIR}}, {{RUBRIC_PATH}}, {{REPORT_PATH}}, etc.) compatible with eval SKILL.md prompt composition

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] scorer-protocol.md contains the three-phase workflow (Reasoning Audit, Rubric Scoring, Blindspot Hunt) extracted from doc-scorer.md
- [x] scorer-protocol.md does NOT contain persona selection logic or domain-specific failure patterns
- [x] scorer-protocol.md ends with the HARD-RULE output format specification (SCORE/DIMENSIONS/ATTACKS)
- [x] reviser-protocol.md contains the attack-point-driven revision workflow extracted from doc-reviser.md
- [x] reviser-protocol.md does NOT reference rubric path as an input
- [x] Both protocol files are self-contained — no references to doc-scorer.md or doc-reviser.md

## Notes
Documentation-only task, no test metrics applicable. Protocol files are stable and generic — changes propagate to all experts via single file.
