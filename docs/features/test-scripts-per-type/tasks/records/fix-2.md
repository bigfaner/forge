---
status: "completed"
started: "2026-05-15 23:09"
completed: "2026-05-15 23:21"
time_spent: "~12m"
---

# Task Record: fix-2 Fix: rewrite e2e tests to use actual CLI interfaces

## Summary
Rewrote e2e tests to use actual CLI interfaces (forge task index) instead of non-existent commands (forge gen-test-scripts, forge breakdown-tasks, forge quick-tasks). Created 12 new tests that verify per-type task generation via forge task index: multi-type creates per-type gen-scripts tasks, single-type creates one task, no test-cases falls back to legacy, zero-type falls back to legacy, run task depends on all per-type gen tasks, multi-profile per-type tasks, quick mode per-type, per-type MD mentions test type, idempotent re-run, correct task IDs, and shared infrastructure not duplicated.

## Changes

### Files Created
- tests/e2e/features/test-scripts-per-type/test_scripts_per_type_cli_test.go
- tests/e2e/features/test-scripts-per-type/go.mod
- tests/e2e/features/test-scripts-per-type/go.sum

### Files Modified
无

### Key Decisions
- Tests use forge task index CLI command which is the actual interface for per-type task generation
- Tests create temp project fixtures with test-cases.md containing different type distributions
- Each test verifies index.json output and generated .md file content rather than calling non-existent CLI commands
- Test module is standalone (own go.mod) since it lives in a subdirectory like justfile-canonical-e2e

## Test Results
- **Tests Executed**: Yes
- **Passed**: 818
- **Failed**: 0
- **Coverage**: 90.3%

## Acceptance Criteria
- [x] E2E tests call forge task index instead of non-existent commands
- [x] Tests verify per-type gen-scripts task generation from forge task index
- [x] Tests verify type detection from test-cases.md
- [x] Tests verify fallback to legacy when no types detected
- [x] just test passes (unit tests)

## Notes
Original 12 tests (TC-001 to TC-012) called forge gen-test-scripts, forge breakdown-tasks, and forge quick-tasks which don't exist. Rewrote all 12 to test the actual CLI interface: forge task index. Go unit tests for DetectTypesFromTestCases already existed in forge-cli/pkg/task/testgen_test.go and remain passing.
