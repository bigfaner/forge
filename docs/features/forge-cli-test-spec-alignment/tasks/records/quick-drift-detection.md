---
status: "completed"
started: "2026-05-21 00:54"
completed: "2026-05-21 01:00"
time_spent: "~6m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Drift detection on all project-level specs. Found drift in forge-cli-reference.md (missing worktree status/push subcommands, missing forge version command, incorrect start usage). Fixed all drift and regenerated vocabulary index.

## Changes

### Files Created
无

### Files Modified
- docs/conventions/forge-cli-reference.md
- docs/.vocabulary.md

### Key Decisions
- Added forge worktree status and forge worktree push to CLI reference (detected from root.go AddCommand chain)
- Added forge version as hidden command to CLI reference
- Fixed forge worktree start usage from <slug> (required) to [slug] (optional)

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] All project-level specs validated against current codebase
- [x] Drifted specs updated to match code
- [x] Vocabulary index regenerated

## Notes
Drift-only mode (no PRD/design files). All business-rules and other conventions verified current. Only forge-cli-reference.md had drift: 2 missing worktree subcommands (status, push), 1 missing top-level command (version), and 1 incorrect usage syntax (start <slug> -> start [slug]).
