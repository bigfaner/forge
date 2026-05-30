---
status: "completed"
started: "2026-05-30 06:18"
completed: "2026-05-30 06:20"
time_spent: "~2m"
---

# Task Record: 22 Fix: guide.md surface coverage + CLI reference gaps

## Summary
Fixed guide.md gaps: added mobile surface orchestration, tui/mobile test types, expanded CLI commands (task/feature/config/pipeline), and added Configuration subsection documenting auto.* config keys

## Changes

### Files Created
无

### Files Modified
- plugins/forge/hooks/guide.md

### Key Decisions
无

## Document Metrics
84 lines (was 48), 4 AC items all met, 5 subsections in Forge CLI section

## Referenced Documents
- docs/features/plugin-consistency-audit/reports/05-commands-agent-hooks.md
- docs/reference/test-type-model.md
- plugins/forge/skills/run-tests/SKILL.md
- plugins/forge/skills/run-tests/rules/surfaces/mobile.md

## Review Status
final

## Acceptance Criteria
- [x] Surface Type orchestration includes mobile (probe + teardown + test-setup)
- [x] Test Type examples include tui -> Terminal Functional Test and mobile -> Mobile E2E Test
- [x] Forge CLI section adds task claim/status/add/submit, feature set/complete, config get, quality-gate, surfaces detect
- [x] Configuration subsection documents forge config get auto.* keys

## Notes
All 4 HOOK findings (HOOK-01 through HOOK-04) addressed in a single file edit. Mobile orchestration pattern sourced from run-tests SKILL.md line 163 and rules/surfaces/mobile.md. CLI command list sourced from usage across commands/ and skills/ directories.
