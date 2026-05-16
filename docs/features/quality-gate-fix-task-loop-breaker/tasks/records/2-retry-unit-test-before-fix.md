---
status: "completed"
started: "2026-05-16 20:50"
completed: "2026-05-16 21:04"
time_spent: "~14m"
---

# Task Record: 2 P1: Retry unit-test once before creating fix task

## Summary
Added retry-once policy for the unit-test step in quality_gate.go. When unit tests fail, they are re-run once. If retry passes, a warning is logged and no fix task is created. If retry also fails, a fix task is created with 'retried once, both attempts failed' in the description plus the retry-run output. Retry logic only applies to unit-test, not compile/fmt/lint.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/internal/cmd/quality_gate_test.go

### Key Decisions
- Extracted runUnitTestStep function with testRunFunc parameter for testability, allowing mock injection without modifying testrunner package
- Retry logic encapsulated in runUnitTestStep as a gate policy (not in testrunner), matching task requirement
- On double-failure, combined output (first attempt + retry) is written to raw output file and passed to addFixTask, with 'retried once, both attempts failed' prepended

## Test Results
- **Tests Executed**: Yes
- **Passed**: 4
- **Failed**: 0
- **Coverage**: 81.0%

## Acceptance Criteria
- [x] When unit tests fail, they are retried once before creating a fix task
- [x] If retry passes, a warning is logged to stderr, no fix task is created, gate continues to e2e step
- [x] If retry also fails, a fix task is created with description including 'retried once, both attempts failed' plus the retry-run output
- [x] Retry logic only applies to the unit-test step, not compile/fmt/lint
- [x] New tests added for: retry pass (no fix task), retry fail (fix task with retry mention)

## Notes
runUnitTestStep accepts a testRunFunc parameter for dependency injection. In production, testrunner.RunProjectTests is passed. In tests, a mock function is injected to control pass/fail behavior without subprocess execution.
