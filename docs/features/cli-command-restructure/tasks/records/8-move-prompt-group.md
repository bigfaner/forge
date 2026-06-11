---
status: "completed"
started: "2026-05-23 03:39"
completed: "2026-05-23 03:47"
time_spent: "~8m"
---

# Task Record: 8 Move forge prompt group to prompt/ subdirectory

## Summary
Moved all prompt subcommand files from internal/cmd/ to internal/cmd/prompt/ subdirectory. Created Register() function following the existing sub-package pattern (forensic, task, test). Updated root.go and root_test.go to import promptpkg. Tests self-contained in new package with own helpers (no import of parent cmd package). Updated forge-cli-reference.md convention doc.

## Changes

### Files Created
- forge-cli/internal/cmd/prompt/register.go
- forge-cli/internal/cmd/prompt/prompt_get.go
- forge-cli/internal/cmd/prompt/prompt_test.go

### Files Modified
- forge-cli/internal/cmd/root.go
- forge-cli/internal/cmd/root_test.go
- docs/conventions/forge-cli-reference.md

### Key Decisions
- Used promptpkg alias for import to avoid name collision with pkg/prompt package
- Test helpers duplicated in prompt package to satisfy hard rule: no import of internal/cmd

## Test Results
- **Tests Executed**: Yes
- **Passed**: 7
- **Failed**: 0
- **Coverage**: 86.7%

## Acceptance Criteria
- [x] All prompt subcommand files are in internal/cmd/prompt/
- [x] New package exports a Register() function
- [x] root.go updated to use the new package
- [x] go build ./... passes
- [x] go test ./... passes
- [x] forge prompt behavior unchanged

## Notes
Hard rule verified: prompt sub-package does NOT import internal/cmd (no circular deps). All 7 prompt tests pass, all cmd package tests pass.
