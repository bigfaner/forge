---
status: "completed"
started: "2026-05-24 21:23"
completed: "2026-05-24 21:48"
time_spent: "~25m"
---

# Task Record: 5 Add forge surfaces detect subcommand

## Summary
Implemented `forge surfaces detect` subcommand with read-only default, --apply flag for TUI confirmation and config writing, and non-interactive mode that prints to stdout. Registered as subcommand of existing surfacesCmd. Reuses DetectSurfacesWithConflicts for detection+inference pipeline and askSurfaceConfirmation for TUI flow.

## Changes

### Files Created
- forge-cli/internal/cmd/surfaces_detect.go
- forge-cli/internal/cmd/surfaces_detect_test.go

### Files Modified
无

### Key Decisions
- Added isInteractiveTerminalFunc variable for testability instead of direct os.Stdin.Stat check
- Used detectEmptyError sentinel type for exit code 1 on empty detection (SilenceErrors=true suppresses output)
- Reused askSurfaceConfirmation in --apply mode (re-runs detection internally, consistent with init flow)
- stdout format uses detected:/inferred: prefix (mapping from internal dependency:/inference: format)
- Help text documents unstable stdout format per Hard Rule

## Test Results
- **Tests Executed**: Yes
- **Passed**: 15
- **Failed**: 0
- **Coverage**: 85.0%

## Acceptance Criteria
- [x] forge surfaces detect runs read-only, shows results with source annotations, exits without writing config
- [x] forge surfaces detect --apply shows TUI confirmation, writes to config on confirm, exit code 0
- [x] Non-interactive terminal: prints to stdout, no TUI, no config write; exit 0 on success, 1 on empty
- [x] Stdout format: one line per surface path=type (source) with detected:/inferred: prefix
- [x] Empty detection: prints nothing, exits with code 1
- [x] After --apply confirm, config file contains detected surfaces
- [x] --project-root flag supported

## Notes
Non-interactive stdout format documented as UNSTABLE in command help text per Hard Rule. No config write without explicit --apply flag enforced in all code paths.
