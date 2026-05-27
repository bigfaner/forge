---
status: "completed"
started: "2026-05-27 00:43"
completed: "2026-05-27 00:45"
time_spent: "~2m"
---

# Task Record: T-review-doc Review Documentation Quality

## Summary
Reviewed documentation quality for run-tests-journey-isolation feature. AC-1 through AC-5 passed. AC-6 required fixes: added missing Per-Journey execution sections to mobile.md and tui.md surface rule files.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/run-tests/rules/surfaces/mobile.md
- plugins/forge/skills/run-tests/rules/surfaces/tui.md

### Key Decisions
无

## Document Metrics
AC pass: 6/6 (5 pre-existing pass, 2 fixes applied)

## Referenced Documents
- docs/proposals/run-tests-journey-isolation/proposal.md

## Review Status
fixes-applied

## Acceptance Criteria
- [x] test-run.md references forge:run-tests (not forge:run-e2e-tests)
- [x] SKILL.md includes journey discovery step (ls docs/features/<slug>/testing/)
- [x] SKILL.md specifies per-journey execution: just test <journey>
- [x] SKILL.md specifies dev/probe once, per-journey loop, teardown once
- [x] SKILL.md handles no-journey edge case with gen-journeys suggestion
- [x] Surface rule files updated for per-journey test execution pattern

## Notes
mobile.md and tui.md were missing Per-Journey execution sections. Added sections consistent with web.md and api.md patterns.
