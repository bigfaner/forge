---
id: "4"
title: "Add TUI tree view with --tree flag"
priority: "P1"
estimated_time: "1d"
dependencies: ["2"]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.feature"
mainSession: false
---

# 4: Add TUI tree view with --tree flag

## Description

Implement an interactive TUI tree view for `forge task list --tree` using `bubbletea`. The tree displays tasks in a dependency hierarchy with expand/collapse nodes, keyboard navigation, and dual-encoded status indicators (color + symbol). Falls back to table mode when terminal capabilities are insufficient (SSH, dumb terminals).

## Reference Files
- `proposal.md#TUI-树视图（--tree）` — defines interaction model: keyboard navigation, expand/collapse, status indicators, --tree --sort id behavior
- `proposal.md#Non-Functional-Requirements` — pure Go TUI library (bubbletea), no external tools; color auto-disable in non-TTY
- `proposal.md#Key-Risks` — TUI rendering risk in SSH/remote terminals; mitigation: fallback to table mode with terminal capability detection
- `proposal.md#Scope` — In Scope: TUI tree (Phase 2), status indicators (color+symbol), non-TTY auto-disable
- `proposal.md#Success-Criteria` — Phase 2 acceptance: --tree enters TUI, navigation, status encoding, --tree --sort id interaction, test coverage

## Acceptance Criteria

- [ ] `forge task list --tree` launches TUI tree view showing task dependency hierarchy
- [ ] Keyboard navigation: up/down move between nodes, right/left or enter expand/collapse children
- [ ] Status indicators use dual encoding: color (green=completed, yellow=in_progress, red=blocked/failed, gray=pending) + symbol (✓, ~, ✗, ○)
- [ ] `--tree --sort id` preserves tree structure (parent-child by dependency) but sorts sibling nodes by natural ID within each level
- [ ] Non-TTY or unsupported terminal: gracefully falls back to table mode (same as `forge task list` without `--tree`)
- [ ] TUI exit: `q` or `Ctrl+C` returns to shell
- [ ] Test coverage for tree model construction, status encoding, sibling sorting

## Hard Rules

- Use `bubbletea` as the TUI library — no other TUI frameworks
- Terminal capability detection must happen before entering TUI mode; fall back silently (no error) to table output
- No TUI editing capabilities — view-only (no status modification from tree)

## Implementation Notes

- New file: `internal/cmd/task/tree.go` for the TUI model (bubbletea Model/View/Update)
- Tree data model: build from TaskIndex using dependency edges; root nodes are tasks with no dependencies
- For terminal detection: check `os.Stdout` is a terminal via `golang.org/x/term.IsTerminal()` and `$TERM` env var
- The `--tree` flag conflicts with `--sort` only in sibling ordering — both flags can coexist
- Estimate 5-7 days including cross-platform testing (macOS/iTerm2, Linux/gnome-terminal, Windows Terminal)
