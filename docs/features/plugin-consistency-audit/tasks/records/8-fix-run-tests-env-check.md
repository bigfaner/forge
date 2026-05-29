---
status: "completed"
started: "2026-05-30 05:58"
completed: "2026-05-30 05:59"
time_spent: "~1m"
---

# Task Record: 8 Fix: run-tests env-check.md Playwright hardcodes

## Summary
Replaced Playwright hardcodes in env-check.md Web surface check #3 with Convention-derived generic framework commands

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/run-tests/rules/env-check.md

### Key Decisions
无

## Document Metrics
1 table row updated (2 cells), 0 playwright references remaining

## Referenced Documents
- docs/features/plugin-consistency-audit/reports/02-skills-batch-a.md
- plugins/forge/skills/run-tests/SKILL.md

## Review Status
final

## Acceptance Criteria
- [x] Web surface check #3 no longer hardcodes npx playwright install
- [x] Replaced with Convention-derived generic description (per Convention file framework section)
- [x] Repair Suggestion column also replaced with generic description
- [x] No other Playwright hardcode remnants in file (global search verified)

## Notes
RT-01 fix only. RT-02 (gen-journeys reference) is P2 and out of scope per Implementation Notes.
