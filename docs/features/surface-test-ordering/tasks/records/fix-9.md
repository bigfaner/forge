---
status: "completed"
started: "2026-05-26 17:43"
completed: "2026-05-26 17:55"
time_spent: "~12m"
---

# Task Record: fix-9 Fix: task-lifecycle config fixtures

## Summary
Fixed config fixtures in task-lifecycle test suite: added config.yaml to stageGateTestDir helper (3 locations) so FindProjectRoot recognizes temp dirs as forge project roots; corrected TC_003 fixture type from 'fix' to 'coding.fix'; corrected TSG_011 expected type from 'doc-generation.summary' to 'doc.summary'

## Changes

### Files Created
无

### Files Modified
- tests/task-lifecycle/task_stage_gates_test.go
- tests/task-lifecycle/fix_task_claim_priority_test.go

### Key Decisions
- Added minimal 'surfaces: cli' config.yaml in stageGateTestDir and two inline test setups (TSG_019, TSG_020) to satisfy FindProjectRoot's isForgeProjectRoot check
- Changed TC_003 fix-1 type from 'fix' to 'coding.fix' to match TypeCodingFix constant used in claim dependency blocking logic
- Changed TSG_011 expected type from 'doc-generation.summary' to 'doc.summary' to match actual TypeDocSummary constant

## Test Results
- **Tests Executed**: Yes
- **Passed**: 24
- **Failed**: 0
- **Coverage**: 87.7%

## Acceptance Criteria
- [x] All task-lifecycle e2e tests pass
- [x] stageGateTestDir creates valid config.yaml for FindProjectRoot
- [x] TC_003 fix task claim priority test passes
- [x] TSG_011 index.json type assertion matches actual code

## Notes
Root cause: stageGateTestDir created .forge dir without config.yaml, so FindProjectRoot skipped temp dir. Two additional fixture issues: TC_003 used wrong type 'fix' instead of 'coding.fix', and TSG_011 expected wrong summary type name.
