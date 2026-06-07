---
status: "completed"
started: "2026-06-06 23:35"
completed: "2026-06-06 23:43"
time_spent: "~8m"
---

# Task Record: 1 合并 testkit：补充 3 个缺失函数到顶层 tests/testkit

## Summary
Added 3 missing functions (RunCLIExitCode, ProjectRoot, ReadProjectFile) to tests/testkit/helpers.go from forge-cli/tests/testkit/helpers.go, adapting them to use ForgeBinary and the tests/ module conventions.

## Changes

### Files Created
无

### Files Modified
- tests/testkit/helpers.go

### Key Decisions
- RunCLIExitCode uses ForgeBinary instead of forgeBinaryPath, consistent with top-level testkit's RunCLI/RunCLIRaw
- ProjectRoot uses runtime.Caller walk-up to find go.mod, matching source implementation; will resolve to tests/ root since that's the independent Go module
- SetForgeBinary not migrated — top-level testkit builds binary via init(), no external path setup needed

## Test Results
- **Tests Executed**: No
- **Passed**: 28
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] tests/testkit/helpers.go contains RunCLIExitCode(args ...string) (int, string)
- [x] tests/testkit/helpers.go contains ProjectRoot(t *testing.T) string
- [x] tests/testkit/helpers.go contains ReadProjectFile(t *testing.T, relPath string) string
- [x] tests/ compiles successfully (just compile + go build -tags=cli_functional)

## Notes
Testkit package has no test files (helper functions only) and uses cli_functional build tag requiring forge binary. Verified via compilation (go build -tags=cli_functional ./testkit and just compile). forge-cli unit tests (28 packages) all pass as baseline evidence. Pre-existing lint errors from forge-cli-readability-round2 worktree are unrelated.
