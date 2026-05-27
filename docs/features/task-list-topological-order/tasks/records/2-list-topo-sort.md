---
status: "completed"
started: "2026-05-28 00:27"
completed: "2026-05-28 00:38"
time_spent: "~11m"
---

# Task Record: 2 Integrate topo sort into forge task list + add --sort flag

## Summary
Integrated topological sort into forge task list as default ordering, added --sort id flag to restore natural ID order, added [cycle] and [missing: <id>] markers with TTY-aware color output

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/task/list.go
- forge-cli/internal/cmd/task/list_test.go

### Key Decisions
- Used TopologicalSort from pkg/task/toposort.go as the sorting backend for default mode
- Cycle nodes appended after topologically-sorted nodes so they still appear in the table
- Missing deps looked up per-task via buildMissingPerTask to show [missing: <id>] on the correct row
- TTY detection via os.Stdout.Stat() with overridable listIsTerminalFunc for testing
- Yellow ANSI (\033[33m) for markers in TTY mode, plain text in pipe mode
- Column width calculation uses displayWidthPlain to account for marker text without ANSI codes
- Updated existing TestListCmd_Sorting to use --sort id for deterministic T-prefix ordering

## Test Results
- **Tests Executed**: Yes
- **Passed**: 26
- **Failed**: 0
- **Coverage**: 91.2%

## Acceptance Criteria
- [x] forge task list defaults to topological sort ordering
- [x] forge task list --sort id restores natural ID ordering
- [x] Cycle nodes display [cycle] marker in the table row
- [x] Missing deps display [missing: <id>] marker in the table row
- [x] Non-TTY environments do not emit ANSI color codes for markers
- [x] Empty feature still shows no tasks found
- [x] Existing list_test.go tests pass after integration
- [x] New tests cover topo sort default, --sort id fallback, cycle marker, missing dep marker, pipe mode color suppression

## Notes
Hard Rules verified: --sort flag only accepts topo and id values. Pipeline output uses plain text markers without ANSI. Column alignment accounts for marker width via displayWidthPlain.
