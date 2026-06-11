---
id: "4"
title: "Split worktree.go into command files"
priority: "P1"
estimated_time: "2h"
dependencies: [2]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 4: Split worktree.go into command files

## Description
Split `worktree.go` (1068 lines, 6 commands + completion + TUI + file operations) into per-command files plus shared helpers. This is Phase 2 (oversized file splitting) of the cleanup proposal.

## Reference Files
- `docs/proposals/forge-cli-clean-code/proposal.md` — Source proposal
- `forge-cli/internal/cmd/worktree/worktree.go` — Source file to split
- `forge-cli/internal/cmd/worktree/register.go` — Existing registration file

## Acceptance Criteria
- [ ] `worktree.go` reduced to <300 lines
- [ ] Each command has its own file (or grouped logically)
- [ ] Shared helpers extracted to `helpers.go` or similar
- [ ] All symbols remain in the `worktree` package
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes
- [ ] No behavioral changes

## Hard Rules
- All symbols must remain in the `worktree` package
- Each commit must pass `go build ./...` and `go test ./...`
- Preserve all existing test files unchanged

## Implementation Notes
- 6 commands suggest 6 command files + 1 helpers file + register.go
- TUI-related code could be grouped together
- File operations (create, remove, list) could share a helpers file
