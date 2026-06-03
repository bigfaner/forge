---
status: "completed"
started: "2026-06-04 00:38"
completed: "2026-06-04 00:40"
time_spent: "~2m"
---

# Task Record: 3 Inline test isolation into gen-test-scripts + delete Related sections

## Summary
Inline test-isolation.md decision table into gen-test-scripts SKILL.md and delete Related Skills section

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-test-scripts/SKILL.md

### Key Decisions
无

## Document Metrics
5-row isolation decision table inlined (TEST-isolation-000 to 004), Related Skills section (8 lines) removed, net ~+4 lines, 0 cross-skill references remaining

## Referenced Documents
- docs/proposals/skill-command-independence-audit/proposal.md
- plugins/forge/skills/run-tests/rules/test-isolation.md

## Review Status
final

## Acceptance Criteria
- [x] gen-test-scripts/SKILL.md contains inlined isolation decision table with INLINE origin marker
- [x] Related Skills, Integration, References sections deleted
- [x] gen-test-scripts no longer contains cross-skill references to run-tests internal files

## Notes
HARD-RULE/HARD-GATE/EXTREMELY-IMPORTANT/PROHIBITIONS block count preserved at 13 (12 HARD-RULE + 1 EXTREMELY-IMPORTANT). Pipeline Position diagram reference to run-tests retained as flow description, not knowledge dependency.
