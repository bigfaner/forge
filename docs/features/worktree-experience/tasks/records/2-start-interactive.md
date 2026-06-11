---
status: "completed"
started: "2026-05-20 17:17"
completed: "2026-05-20 17:37"
time_spent: "~20m"
---

# Task Record: 2 Start interactive mode: -i flag for proposal/feature selection

## Summary
Added -i/--interactive flag to `forge worktree start` that presents a selectable list of unfinished proposals and features. Users pick one via numbered list + stdin, and the slug is auto-filled. No external TUI dependencies used - only fmt.Scanln-style approach via bufio.Reader.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/worktree.go
- forge-cli/internal/cmd/worktree_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- Used bufio.NewReader + stdinFunc variable for testability instead of external TUI library
- Proposals take priority over features when same slug exists in both directories
- Changed worktreeStartCmd.Args from ExactArgs(1) to MaximumNArgs(1) to allow -i without slug
- isTerminalFunc as overridable var for TTY detection testing
- Empty list prints helpful message and exits cleanly (no error)
- When both -i and slug provided, slug takes precedence (ignores interactive mode)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 22
- **Failed**: 0
- **Coverage**: 81.8%

## Acceptance Criteria
- [x] forge worktree start -i lists all unfinished proposals and features with their status
- [x] User can select one item from the list; the slug is extracted and used as the argument
- [x] When -i is used, the <slug> positional arg becomes optional (one or the other)
- [x] When both -i and <slug> are provided, <slug> takes precedence (ignore -i)
- [x] Empty list (no proposals or features) prints a helpful message and exits
- [x] Selection prompt works in a terminal context (TTY detection for non-interactive environments)

## Notes
No external TUI dependencies added per hard rules. Used fmt.Fprintf for numbered list and bufio.Reader for stdin. Version bumped from 4.4.3 to 4.5.0 (minor: new feature).
