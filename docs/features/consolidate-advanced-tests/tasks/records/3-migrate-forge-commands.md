---
status: "completed"
started: "2026-06-06 23:59"
completed: "2026-06-07 00:05"
time_spent: "~6m"
---

# Task Record: 3 迁移 forge-commands journey

## Summary
Migrated forge-commands journey: 3 test files (discovery, e2e_commands, forge_info_commands) + main_test.go + 3 contracts moved to tests/forge-commands/. Updated imports to forge-tests/testkit, rewrote main_test.go with ForgeBinary init pattern. forge_init_install_just_test.go excluded — imports forge-cli/pkg/just which is inaccessible from independent tests/ module.

## Changes

### Files Created
- tests/forge-commands/main_test.go
- tests/forge-commands/discovery_test.go
- tests/forge-commands/e2e_commands_test.go
- tests/forge-commands/forge_info_commands_test.go
- tests/forge-commands/contracts/step-1-discovery.md
- tests/forge-commands/contracts/step-2-info-commands.md
- tests/forge-commands/contracts/step-3-e2e-runner.md

### Files Modified
无

### Key Decisions
- Excluded forge_init_install_just_test.go from migration because it imports forge-cli/pkg/just which cannot be resolved from the independent tests/ Go module (module forge-tests has no replace directive for forge-cli)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 16
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] tests/forge-commands/ contains all migrated test files with import testkit forge-tests/testkit
- [x] main_test.go uses ForgeBinary init pattern
- [x] contracts/ directory has 3 contract files correctly migrated
- [x] just test includes this journey and passes

## Notes
3 of 4 source test files migrated. forge_init_install_just_test.go remains in forge-cli/tests/forge-commands/ because it directly calls forge-cli/pkg/just internal APIs (just.EnsureJust, just.DetectJust, etc.) which cannot be imported from the independent forge-tests module. 16 tests pass, 1 skip (TestTC_032 cannot locate justfile from test working dir — pre-existing behavior).
