---
status: "completed"
started: "2026-05-15 01:27"
completed: "2026-05-15 01:43"
time_spent: "~16m"
---

# Task Record: 3 Wire ensureJust into forge init + add --skip-just flag

## Summary
Wire ensureJust into forge init command: added --skip-just flag, inserted ensureJust step between gitignore (step 3) and justfile (step 4), added INSTALLED status support to initAction, and exported function variables in pkg/just for testability.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/init.go
- forge-cli/internal/cmd/init_test.go
- forge-cli/pkg/just/ensure.go
- forge-cli/pkg/just/ensure_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- Inserted ensureJust step between step 3 (gitignore) and step 4 (justfile) as required by hard rule
- Created ensureJustFunc variable in init.go for testability rather than calling just.EnsureJust directly, allowing clean mock injection in integration tests
- Exported DetectJustFunc, InstallViaPackageManagerFunc, ExtractEmbeddedBinaryFunc, EmbeddedBinaryFunc from pkg/just for cross-package testability
- Installation failure is non-blocking: prints WARNING to stderr but init continues

## Test Results
- **Tests Executed**: Yes
- **Passed**: 26
- **Failed**: 0
- **Coverage**: 80.7%

## Acceptance Criteria
- [x] forge init runs ensureJust as a step before the justfile update step
- [x] forge init --skip-just skips the ensureJust step entirely, reporting SKIPPED
- [x] initAction.status supports INSTALLED in addition to existing values
- [x] Init summary output includes the just installation result (INSTALLED/SKIPPED/FAILED)
- [x] When just is already installed and >= 1.40.0, the step reports SKIPPED with version detail
- [x] When just is installed successfully, the step reports INSTALLED with method detail
- [x] When installation fails, the step reports FAILED but init continues (non-blocking)
- [x] forge init --help documents the --skip-just flag

## Notes
Pre-existing test failure TestSaveIndexAndSignalCompletion_SaveIndexError in integration_test.go is unrelated to this change (confirmed by running against pre-change HEAD). Coverage for pkg/just is 85.1%. Version bumped from 3.9.0 to 3.10.0 (minor: new feature).
