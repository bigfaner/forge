---
status: "completed"
started: "2026-05-20 16:50"
completed: "2026-05-20 17:01"
time_spent: "~11m"
---

# Task Record: 1 Add auto.validation and worktree prompts to forge init (huh TUI)

## Summary
Added auto.validation (ModeToggle quick/full) prompts to askAutoBehavior() between cleanCode and gitPush, and added askWorktreeConfig() function with optional source-branch (huh.NewInput) and copy-files (huh.NewMultiSelect) prompts integrated into runConfigInitIfNeeded().

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/init.go
- forge-cli/internal/cmd/init_test.go

### Key Decisions
- Validation prompts follow the same askConfirm pattern as other auto fields (cleanCode, consolidateSpecs) for consistency
- Worktree copy-files multi-select only appears when source-branch is provided, keeping the flow minimal
- Both worktree fields are skippable: empty source-branch + no copy-files = nil WorktreeConfig (no block in config YAML)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 601
- **Failed**: 0
- **Coverage**: 81.0%

## Acceptance Criteria
- [x] askAutoBehavior() prompts for auto.validation quick/full using huh.Select pattern
- [x] Validation prompts appear between cleanCode and gitPush prompts
- [x] After auto behavior, a worktree section prompts for optional source-branch and copy-files
- [x] Both worktree fields are skippable (empty = no worktree block)
- [x] Existing tests in init_test.go still pass

## Notes
无
