---
status: "completed"
started: "2026-05-20 16:49"
completed: "2026-05-20 16:56"
time_spent: "~7m"
---

# Task Record: 1 Isolate forge-cli/tests/e2e/ to TestMain auto-build

## Summary
Isolated forge-cli/tests/e2e/ to TestMain auto-build pattern: compile forge binary from source into temp directory in TestMain, propagate path to tests via forgeBinary package variable and testkit.SetForgeBinary(). All bare exec.Command("forge",...) calls replaced with exec.Command(forgeBinary,...) or testkit's forgeBinaryPath.

## Changes

### Files Created
- forge-cli/tests/e2e/features/fix-task-claim-priority/main_test.go

### Files Modified
- forge-cli/tests/e2e/main_test.go
- forge-cli/tests/e2e/helpers_test.go
- forge-cli/tests/e2e/submit_cli_test.go
- forge-cli/tests/e2e/features/fix-task-claim-priority/fix_task_claim_priority_cli_test.go
- forge-cli/tests/e2e/testkit/helpers.go

### Key Decisions
- Used testkit.SetForgeBinary() to propagate binary path from TestMain to testkit package, keeping testkit's public API unchanged (callers still call RunCLI/RunCLIExitCode/RunCLIWithResult with same signatures)
- Added separate main_test.go in features/fix-task-claim-priority/ because Go treats subdirectories as separate packages -- the parent's forgeBinary variable is not accessible
- forgeBinaryPath defaults to "forge" for backward compatibility if SetForgeBinary is never called

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] forge-cli/tests/e2e/ has a main_test.go with TestMain that builds forge binary from ../../cmd/forge (relative to forge-cli module) into a temp directory
- [x] forgeBinary package variable holds the temp binary path
- [x] helpers_test.go and testkit/helpers.go use the forgeBinary variable instead of bare "forge"
- [x] Temp directory cleaned up in TestMain cleanup
- [x] grep -r 'exec.Command("forge"' forge-cli/tests/e2e/ returns zero results

## Notes
This is test infrastructure only -- no source code changes, so coverage is not applicable. All e2e test packages compile successfully with the new pattern.
