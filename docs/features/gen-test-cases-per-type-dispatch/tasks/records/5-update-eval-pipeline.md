---
status: "completed"
started: "2026-05-17 00:38"
completed: "2026-05-17 00:40"
time_spent: "~2m"
---

# Task Record: 5 Update eval skill table + eval-test-cases command

## Summary
Added 5 new test-cases-* entries to eval SKILL.md (prerequisites, default doc dir, parameters enum, rubric reference, pre-processing, scorer inputs, reviser constraints, final report, next step tables) and refactored eval-test-cases command into a per-type dispatcher with legacy monolithic fallback.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/eval/SKILL.md
- plugins/forge/commands/eval-test-cases.md

### Key Decisions
- Grouped test-cases-* types on same lines where behavior is identical (e.g. pre-processing, scorer inputs) to avoid table bloat while maintaining discoverability
- eval-test-cases dispatcher aggregates per-type scores into a combined pass/fail summary table

## Test Results
- **Tests Executed**: No (noTest task)
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Eval skill Prerequisites table has entries for test-cases-ui/tui/mobile/api/cli
- [x] Eval skill Default Doc Dir table maps each test-cases-* type to testing/ directory
- [x] Eval skill Parameters --type enum includes all 5 new test-cases-* values
- [x] Eval skill Rubric Reference table has entries for all 5 new rubrics (1000 scale, 900 target, 6 iterations)
- [x] eval-test-cases command loops per-type when per-type files exist
- [x] eval-test-cases command falls back to --type test-cases for legacy mode
- [x] Pre-processing for test-cases-* types resolves test profile and passes capabilities to scorer

## Notes
Hard rules verified: core scorer-gate-revise loop untouched; dispatcher passes single {type}-test-cases.md file per invocation; legacy fallback works when only testing/test-cases.md exists.
