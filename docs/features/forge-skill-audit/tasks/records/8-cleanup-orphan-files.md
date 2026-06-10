---
status: "completed"
started: "2026-06-10 19:33"
completed: "2026-06-10 19:34"
time_spent: "~1m"
---

# Task Record: 8 Clean up test-guide orphan files (M-4)

## Summary
Moved 2 orphan rules files (draft-generation.md and pattern-extraction.md) from test-guide/rules/ to test-guide/rules/_deprecated/, matching the eval skill _deprecated/ convention

## Changes

### Files Created
- plugins/forge/skills/test-guide/rules/_deprecated/draft-generation.md
- plugins/forge/skills/test-guide/rules/_deprecated/pattern-extraction.md

### Files Modified
无

### Key Decisions
无

## Document Metrics
2 files moved to _deprecated/, 0 reference breakages

## Referenced Documents
- docs/proposals/forge-skill-audit/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] draft-generation.md and pattern-extraction.md moved to skills/test-guide/rules/_deprecated/
- [x] No broken references in SKILL.md or other rules files

## Notes
Confirmed via grep that SKILL.md and other rules files never referenced these files, so no reference updates needed. _deprecated/ directory follows eval skill convention.
