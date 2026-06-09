---
status: "blocked"
started: "2026-06-09 16:04"
completed: "N/A"
time_spent: ""
---

# Task Record: T-test-run Run CLI Functional Test

## Summary
Ran 37 CLI functional tests across 3 journeys. Fixed 5 test script bugs (worktree line counting, symlink error matching, non-TTY stdin detection). 35/37 tests pass. 2 tests blocked by production code defect: forge worktree remove cannot handle corrupted worktrees (git worktree remove --force fails validation when .git file is missing). Created fix task fix-1.

## Changes

### Files Created
- tests/results/latest.md

### Files Modified
- tests/idempotent-start/step3_verify_state_test.go
- tests/idempotent-start/idempotent_start_smoke_test.go
- tests/idempotent-start/step5_corrupted_directory_test.go
- tests/start-existing-flags/main_test.go
- tests/start-existing-flags/interactive_test.go

### Key Decisions
无

## Cases Generated
37

## Cases Evaluated
N/A

## Scripts Created
无

## Test Results
37 cases: 35 passed, 2 failed (production code defect in forge worktree remove). Test script bugs fixed: (1) strings.Count overcounted slug occurrences in git worktree list output - fixed to count matching lines instead; (2) symlink error assertion too narrow - broadened to match actual production error message; (3) non-TTY tests used inherited TTY stdin - added runForgeStartInteractiveNonTTY with os.Pipe to force non-char-device stdin on macOS.

## Acceptance Criteria
- [ ] All test cases MUST pass
- [x] Tests MUST verify actual functional behavior

## Notes
Blocked by fix-1: forge worktree remove needs fallback to os.RemoveAll + git worktree prune when git worktree remove --force fails due to missing .git file in corrupted worktrees. Journeys: idempotent-start (19/19 pass), start-existing-flags (9/9 pass), corrupted-worktree-recovery (7/9 pass, 2 blocked by production defect).
