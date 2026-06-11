---
id: "6"
title: "Move forge worktree group to worktree/ subdirectory"
priority: "P2"
estimated_time: "1h"
dependencies: ["2"]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 6: Move forge worktree group to worktree/ subdirectory

## Description

Move all forge worktree subcommand files from `internal/cmd/` to `internal/cmd/worktree/` subdirectory. Create a `Register()` function for the worktree command group.

## Reference Files
- `docs/proposals/cli-command-restructure/proposal.md` — Source proposal

## Acceptance Criteria

- All worktree subcommand files are in `internal/cmd/worktree/`
- New package exports a `Register()` function
- root.go updated to use the new package
- `go build ./...` passes
- `go test ./...` passes
- `forge worktree` subcommands work identically

## Hard Rules

- The worktree sub-package must NOT import `internal/cmd` (no circular deps)

## Implementation Notes

Files to move to `internal/cmd/worktree/`:
- worktree.go, worktree_test.go
