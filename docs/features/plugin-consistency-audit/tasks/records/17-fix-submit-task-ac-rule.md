---
status: "completed"
started: "2026-05-30 06:11"
completed: "2026-05-30 06:11"
time_spent: ""
---

# Task Record: 17 Fix: submit-task AC vs testsFailed rule conflict

## Summary
Resolved C-30 conflict in submit-task SKILL.md: added priority hierarchy rule stating category-specific rules override Common Rules when they conflict. Also resolved C-31 by clarifying testsPassed/testsFailed/coverage are conditional on status:completed.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/submit-task/SKILL.md

### Key Decisions
无

## Document Metrics
1 rule added (priority hierarchy), 1 rule refined (coding-only conditional clarification)

## Referenced Documents
- docs/features/plugin-consistency-audit/reports/04-skills-batch-c.md
- plugins/forge/skills/submit-task/data/record-format-coding.md

## Review Status
final

## Acceptance Criteria
- [x] Common Rules explicitly state category-specific rules override Common Rules on conflict
- [x] AC vs testsFailed conflict resolved via priority hierarchy (no merge needed)

## Notes
C-31 (coding-only vs conditional ambiguity) was fixed as an in-scope bonus since it was in the same Common Rules section being edited.
