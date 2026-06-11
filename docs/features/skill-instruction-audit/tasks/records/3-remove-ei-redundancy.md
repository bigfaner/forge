---
status: "completed"
started: "2026-05-28 22:55"
completed: "2026-05-28 22:58"
time_spent: "~3m"
---

# Task Record: 3 Remove EXTREMELY-IMPORTANT redundancy from skills

## Summary
Removed EXTREMELY-IMPORTANT redundancy from 5 skill files via constraint-level audit: test-guide (10->2 E-I items), eval (3->1), init-justfile (removed 3 Notes overlapping E-I), clean-code (merged 3 preserve-scope statements into 1 core principle, removed 2 HARD-RULEs), deep-research (deleted Report Structure key points duplicating template).

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/test-guide/SKILL.md
- plugins/forge/skills/eval/SKILL.md
- plugins/forge/skills/init-justfile/SKILL.md
- plugins/forge/skills/clean-code/SKILL.md
- plugins/forge/skills/deep-research/SKILL.md

### Key Decisions
无

## Document Metrics
5 files modified, 8 E-I items removed, 3 Notes items removed, 2 HARD-RULEs removed, 1 principle consolidated

## Referenced Documents
- docs/proposals/skill-instruction-audit/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] test-guide/SKILL.md E-I has <=4 items; each passes constraint-level audit
- [x] eval/SKILL.md has 1 concise E-I rule
- [x] init-justfile/SKILL.md Notes has no E-I overlap
- [x] clean-code/SKILL.md has 1 preserve scope statement
- [x] deep-research/SKILL.md has no key points duplicating template

## Notes
Constraint-level audit applied: each retained E-I item verified against body for enforcement level differential. test-guide retained 2 cross-step guardrails (verbatim preservation + draft fallback). eval merged 3 rules into 1. clean-code merged scope constraints into core principle, eliminated Principle 5 and 2 HARD-RULEs as body-level duplicates.
