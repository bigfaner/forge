---
id: "13"
title: "Unify error handling pattern across commands"
priority: "P2"
estimated_time: "1h"
dependencies: [11, 12]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 13: Unify error handling pattern across commands

## Description
After tasks 11 and 12 fix specific anti-patterns, scan all commands for remaining `os.Exit` calls in non-top-level functions and unify the pattern: internal functions return `error`, only top-level RunE handlers may call `os.Exit`. Phase 4 (anti-pattern fix).

## Reference Files
- `docs/proposals/forge-cli-clean-code/proposal.md` — Source proposal
- `forge-cli/internal/cmd/` — All command files

## Acceptance Criteria
- [ ] 0 `os.Exit` calls in non-top-level functions (top-level = direct RunE/Error handlers for cobra.Command)
- [ ] All internal functions return `error` and let callers decide on exit behavior
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes

## Hard Rules
- The 2 `os.Exit(0)` in `quality_gate.go` RunE handlers are top-level and MUST be preserved
- Do not introduce new error types or wrappers — just change control flow

## Implementation Notes
- Run `grep -r "os.Exit" forge-cli/internal/cmd/` to find remaining instances
- For each hit, determine if it's in a top-level RunE handler or an internal function
- Only fix internal functions — top-level handlers keep their exit behavior
