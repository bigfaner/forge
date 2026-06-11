---
status: "completed"
started: "2026-05-20 17:14"
completed: "2026-05-20 17:21"
time_spent: "~7m"
---

# Task Record: 3 Fix tests/e2e/ feature tests to use TestMain-built binary

## Summary
Unified all e2e test files under tests/e2e/ to use the TestMain-built forge binary via the forgeBinary variable, eliminating duplicate binary build functions and walk-up path resolution.

## Changes

### Files Created
无

### Files Modified
- tests/e2e/main_test.go
- tests/e2e/test_scripts_per_type_cli_test.go
- tests/e2e/quick_test_slim_cli_test.go
- tests/e2e/task_lifecycle_hardening_cli_test.go

### Key Decisions
- Removed forgeBin() and forgeBinPath from test_scripts_per_type_cli_test.go, replaced with package-level forgeBinary
- Removed quickSlimBin() and quickSlimBinPath from quick_test_slim_cli_test.go, replaced with package-level forgeBinary
- Fixed main_test.go initialization order: moved forgeBinary assignment into TestMain (after init()) to avoid capturing empty ForgeBinary before init() runs

## Test Results
- **Tests Executed**: Yes
- **Passed**: 10
- **Failed**: 0
- **Coverage**: 66.7%

## Acceptance Criteria
- [x] All test files under tests/e2e/ use the forgeBinary variable from main_test.go instead of bare 'forge'
- [x] grep -r 'exec.Command("forge"' tests/e2e/ returns zero results (excluding justfile-canonical-e2e)
- [x] No duplicate TestMain in test-knowledge-convention-driven/; it uses parent's ForgeBinary
- [x] Walk-up path resolution (forge-cli/bin/forge) eliminated from all test files

## Notes
无
