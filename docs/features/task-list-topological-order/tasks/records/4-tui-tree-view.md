---
status: "completed"
started: "2026-05-28 00:45"
completed: "2026-05-28 01:02"
time_spent: "~17m"
---

# Task Record: 4 Add TUI tree view with --tree flag

## Summary
Implemented interactive TUI tree view for `forge task list --tree` using bubbletea. Added dependency tree construction from TaskIndex, keyboard navigation (up/down/j/k), expand/collapse (left/right/h/l/enter), dual-encoded status indicators (color + symbol), terminal capability detection with graceful fallback to table mode, and `--tree --sort id` sibling ordering.

## Changes

### Files Created
- forge-cli/internal/cmd/task/tree.go
- forge-cli/internal/cmd/task/tree_test.go

### Files Modified
- forge-cli/internal/cmd/task/list.go
- forge-cli/scripts/version.txt

### Key Decisions
- Used bubbletea as TUI library per Hard Rules, with lipgloss for color styling
- Cycle nodes detected via TopologicalSort and rendered as roots without children to avoid infinite recursion
- Terminal capability detection (TTY + TERM env) happens before entering TUI mode; non-TTY silently falls back to table output
- Tree model (treeNode/flatItem) separate from TUI model (treeModel) for testability
- renderTreeFallback provides plain-text tree output as intermediate between full TUI and table

## Test Results
- **Tests Executed**: Yes
- **Passed**: 30
- **Failed**: 0
- **Coverage**: 74.1%

## Acceptance Criteria
- [x] forge task list --tree launches TUI tree view showing task dependency hierarchy
- [x] Keyboard navigation: up/down move between nodes, right/left or enter expand/collapse children
- [x] Status indicators use dual encoding: color + symbol
- [x] --tree --sort id preserves tree structure but sorts sibling nodes by natural ID
- [x] Non-TTY or unsupported terminal: gracefully falls back to table mode
- [x] TUI exit: q or Ctrl+C returns to shell
- [x] Test coverage for tree model construction, status encoding, sibling sorting

## Notes
30 new tests covering buildForest (6), statusSymbol (9), statusColor (9), canUseTUI (4), renderTreePlain (3), flattenTree (2), treeModel navigation/view (10), integration (2), renderTreeFallback (1). Functions not testable via unit tests (runTreeTUI, getTerminalInfo, Init) are thin wrappers around terminal I/O.
