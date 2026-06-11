---
status: "completed"
started: "2026-06-02 23:48"
completed: "2026-06-02 23:53"
time_spent: "~5m"
---

# Task Record: 17 Translate init-justfile + run-tests + gen-contracts surface rules to English

## Summary
Translated 11 Chinese surface rule files to English: init-justfile (5 files), run-tests (5 files), gen-contracts journey-contract-model.md. All logic, rules, and constraints preserved exactly.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/init-justfile/rules/surfaces/cli.md
- plugins/forge/skills/init-justfile/rules/surfaces/api.md
- plugins/forge/skills/init-justfile/rules/surfaces/web.md
- plugins/forge/skills/init-justfile/rules/surfaces/tui.md
- plugins/forge/skills/init-justfile/rules/surfaces/mobile.md
- plugins/forge/skills/run-tests/rules/surfaces/cli.md
- plugins/forge/skills/run-tests/rules/surfaces/api.md
- plugins/forge/skills/run-tests/rules/surfaces/web.md
- plugins/forge/skills/run-tests/rules/surfaces/tui.md
- plugins/forge/skills/run-tests/rules/surfaces/mobile.md
- plugins/forge/skills/gen-contracts/rules/journey-contract-model.md

### Key Decisions
无

## Document Metrics
11 files translated, 0 Chinese characters remaining (verified via grep), surface types all lowercase

## Referenced Documents
- docs/proposals/surface-first-testing/proposal.md
- docs/conventions/forge-distribution.md
- plugins/forge/skills/gen-journeys/rules/surface-cli.md

## Review Status
final

## Acceptance Criteria
- [x] init-justfile 5 surface rule files translated to English with logic preserved
- [x] run-tests 5 surface rule files translated to English with logic preserved
- [x] gen-contracts journey-contract-model.md translated to English with logic preserved
- [x] No Chinese characters in translated files (except technical paths like tests/<journey>/)
- [x] Surface type terms remain lowercase: web/api/cli/tui/mobile

## Notes
Used gen-journeys surface-cli.md as English template reference for consistent terminology and style.
