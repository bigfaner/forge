---
status: "completed"
started: "2026-06-06 13:39"
completed: "2026-06-06 13:40"
time_spent: "~1m"
---

# Task Record: 10 Update hooks/guide.md and submit-task record format

## Summary
Updated hooks/guide.md test directory convention to reflect surface-key adaptive rules (multi-surface: tests/<surfaceKey>/<journey>/, single: tests/<journey>/); updated submit-task record-format-test.md example paths to show multi-surface directory structure

## Changes

### Files Created
无

### Files Modified
- plugins/forge/hooks/guide.md
- plugins/forge/skills/submit-task/data/record-format-test.md

### Key Decisions
无

## Document Metrics
2 files updated, 2 ACs met

## Referenced Documents
- docs/proposals/surface-key-test-alignment/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] hooks/guide.md test directory convention reflects surface-key adaptive rules
- [x] submit-task/data/record-format-test.md example paths match new directory structure

## Notes
guide.md line 77 updated from fixed tests/<journey>/ to adaptive rule matching test-guide SKILL.md. record-format-test.md example paths updated from tests/login/ to tests/api/login/ to illustrate multi-surface directory structure.
