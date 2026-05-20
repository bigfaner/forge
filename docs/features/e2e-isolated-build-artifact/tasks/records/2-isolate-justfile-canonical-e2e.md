---
status: "blocked"
started: "2026-05-20 16:57"
completed: "N/A"
time_spent: ""
---

# Task Record: 2 Isolate justfile-canonical-e2e/ to TestMain auto-build

## Summary
Isolated justfile-canonical-e2e module to TestMain auto-build pattern: added main_test.go with TestMain that compiles forge into a temp directory, simplified helpers_test.go to use forgeBinaryPath package variable instead of hardcoded relative path with lazy build-if-missing fallback. Temp directory is cleaned up on exit via defer os.RemoveAll.

## Changes

### Files Created
- tests/e2e/justfile-canonical-e2e/main_test.go

### Files Modified
- tests/e2e/justfile-canonical-e2e/helpers_test.go

### Key Decisions
- Followed the same TestMain pattern from tests/e2e/main_test.go for consistency
- forgeCLIDir walks up from test module to find forge-cli/cmd/forge, respecting the separate Go module constraint
- Kept forgeBinary() as a thin wrapper for consistency with existing test callers

## Test Results
- **Tests Executed**: No
- **Passed**: 16
- **Failed**: 4
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] main_test.go added with TestMain that builds forge into temp dir
- [x] forgeBinaryPath package variable set to temp binary path
- [x] helpers_test.go uses forgeBinaryPath instead of ../../../forge-cli/bin/forge
- [x] No reference to forge-cli/bin/forge in this module
- [x] Temp directory cleaned up on exit

## Notes
4 pre-existing test failures (TC-016 through TC-019) are unrelated to this refactor — they test 'no profile configured' error but get 'NO_PROJECT' error instead. These failures existed before this change.
