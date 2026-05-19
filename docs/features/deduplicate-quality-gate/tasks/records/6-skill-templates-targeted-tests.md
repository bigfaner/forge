---
status: "completed"
started: "2026-05-20 00:16"
completed: "2026-05-20 00:18"
time_spent: "~2m"
---

# Task Record: 6 Update skill templates with targeted test instructions

## Summary
Replaced `just test` instructions with targeted test instructions (framework-native `go test -race -cover ./changed/package/...`) across all 7 template files. Added CLI submit gate notes. Preserved static checks (compile/fmt/lint).

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/coding-feature.md
- forge-cli/pkg/prompt/data/coding-enhancement.md
- forge-cli/pkg/prompt/data/coding-cleanup.md
- forge-cli/pkg/prompt/data/coding-refactor.md
- forge-cli/pkg/prompt/data/coding-fix.md
- forge-cli/pkg/template/data/fix-task.md
- forge-cli/pkg/template/data/cleanup-task.md

### Key Decisions
- Used `go test -race -cover ./changed/package/...` as the targeted test command pattern since forge-cli is a Go project
- Kept `just compile/fmt/lint` for static verification -- only `just test` was replaced per proposal Tier 1
- Added explicit Note blocks in all 7 files clarifying that full project-wide tests run at CLI submit

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] No template references `just test` for task-level verification (grep confirms zero matches in all 7 files)
- [x] All 5 coding-*.md templates instruct agents to run targeted tests
- [x] fix-task.md and cleanup-task.md templates updated to targeted tests
- [x] Templates still instruct agents to run `just compile`, `just fmt`, `just lint` for static checks
- [x] Templates note that the CLI submit gate handles full verification

## Notes
无
