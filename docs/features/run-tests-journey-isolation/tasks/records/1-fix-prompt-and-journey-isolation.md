---
status: "completed"
started: "2026-05-27 00:07"
completed: "2026-05-27 00:10"
time_spent: "~3m"
---

# Task Record: 1 Fix prompt template + add journey isolation to run-tests skill

## Summary
Fixed P0 bug (forge:run-e2e-tests -> forge:run-tests) in prompt template and added journey isolation (Step 1.5 discovery + per-journey test loop) to run-tests SKILL.md and all surface rule files

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/test-run.md
- plugins/forge/skills/run-tests/SKILL.md
- plugins/forge/skills/run-tests/rules/surfaces/cli.md
- plugins/forge/skills/run-tests/rules/surfaces/web.md
- plugins/forge/skills/run-tests/rules/surfaces/api.md

### Key Decisions
无

## Document Metrics
1 prompt template fixed, 1 SKILL.md updated with journey discovery + per-journey execution, 3 surface rule files updated with Per-Journey execution section

## Referenced Documents
- docs/proposals/run-tests-journey-isolation/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] test-run.md references forge:run-tests (not forge:run-e2e-tests)
- [x] SKILL.md includes journey discovery step via ls docs/features/<slug>/testing/
- [x] SKILL.md specifies per-journey execution: just test <journey> for each discovered journey
- [x] SKILL.md specifies dev/probe once, per-journey loop for test, teardown once (web/api/mobile)
- [x] SKILL.md handles no-journey edge case with error suggesting gen-journeys
- [x] Surface rule files updated to reflect per-journey test execution pattern

## Notes
All 6 AC items pass. P0 fix confirmed on lines 11 and 31 of test-run.md. Step 1.5 added between Step 1 and Step 2 in SKILL.md. Step 4 rewritten with per-journey loop for both 4a (web/api/mobile) and 4b (cli/tui) sequences. Error handling table updated with no-journeys entry.
