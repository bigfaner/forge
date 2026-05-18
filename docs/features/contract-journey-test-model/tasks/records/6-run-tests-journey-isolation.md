---
status: "completed"
started: "2026-05-18 01:59"
completed: "2026-05-18 02:29"
time_spent: "~30m"
---

# Task Record: 6 run-tests 重写 + Journey 隔离

## Summary
Implemented run-tests Journey isolation: added forge testing run-journey CLI command that reads test-command from .forge/config.yaml, creates isolated temp directories per journey (with journey name + random suffix), executes tests in isolation, collects results with structured reporting, and cleans up temp dirs regardless of success or failure. Contract validation failure output includes dimension name, contract file path, expected value, and actual value.

## Changes

### Files Created
- forge-cli/internal/cmd/journey_isolation.go
- forge-cli/internal/cmd/journey_isolation_test.go

### Files Modified
- forge-cli/internal/cmd/testing.go
- forge-cli/internal/cmd/testing_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- Journey isolation uses system temp dir (not project dir) to guarantee no filesystem interference between parallel journeys
- Cleanup uses defer pattern to guarantee temp dir removal regardless of execution success or failure
- ContractFailure includes dimension name, contract path, expected and actual values for precise failure diagnosis
- readTestCommand reads from ForgeConfig.TestCommand which projects declare in .forge/config.yaml
- resolveJourneyExecutionConfig falls back to profile.ReadLanguages when resolveLanguageFromFlags fails (no --language flag)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 26
- **Failed**: 0
- **Coverage**: 80.5%

## Acceptance Criteria
- [x] run-tests calls project-declared test command from .forge/config.yaml test-command
- [x] Each Journey runs in isolated temp dir (path includes journey name + random suffix), cleaned up after execution
- [x] 3 Journeys running in parallel produce identical results to sequential (filesystem state, exit codes, output diff empty)
- [x] Contract validation failure output includes dimension name, contract file path, expected value, actual value
- [x] Execution failure does not block other Journeys

## Notes
Hard rule compliance: (1) Journey isolation via system temp dirs, not shared state; (2) Failure in one goroutine does not block others; (3) Cleanup runs via defer after execution regardless of outcome. executeJourneyInIsolation tested with real shell scripts for success, failure, and cwd verification.
