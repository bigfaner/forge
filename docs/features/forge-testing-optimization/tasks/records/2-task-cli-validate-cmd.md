---
status: "completed"
started: "2026-05-10 13:19"
completed: "2026-05-10 13:28"
time_spent: "~9m"
---

# Task Record: 2 Add task-cli validate-specs command

## Summary
Add validate-specs command to task-cli that spawns the Node.js validation script against spec files in tests/e2e/features/<slug>/, with structured pass/fail output, graceful degradation when Node/ts-morph unavailable, and Windows-compatible path handling

## Changes

### Files Created
- task-cli/internal/cmd/validate_specs.go
- task-cli/internal/cmd/validate_specs_test.go

### Files Modified
- task-cli/internal/cmd/root.go

### Key Decisions
- Used exitFunc variable for os.Exit override in tests to avoid killing test process
- Graceful degradation: exit 0 with WARNING when Node.js or ts-morph not available (per proposal risk mitigation)
- Spec discovery uses filepath.Glob with *.spec.ts pattern for platform-agnostic path handling
- Auto-detect test-cases.md path from feature context, skip E2 check if not found

## Test Results
- **Tests Executed**: Yes
- **Passed**: 16
- **Failed**: 0
- **Coverage**: 85.4%

## Acceptance Criteria
- [x] task validate-specs command executes the validation script against spec files
- [x] Command discovers spec files from tests/e2e/features/<slug>/ based on current feature context
- [x] Structured output: prints validation results (errors/warnings) to stdout
- [x] Exit code: 0 if no errors, 1 if errors found, 2 if script fails to run
- [x] Unit tests cover: spec discovery, output parsing, error handling
- [x] Works on Windows (path separator handling in exec.Command)

## Notes
Integration tests use mock .mjs scripts executed via node (no shell dependency). All tests pass on Windows.
