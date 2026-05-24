---
id: "7"
title: "Extract defaultRunClaude to shared location"
priority: "P1"
estimated_time: "30m"
dependencies: [4]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 7: Extract defaultRunClaude to shared location

## Description
`defaultRunClaude()` is identically defined in both `claude.go` and `worktree.go`. Extract to a single shared location and have both call sites reference it. Phase 3 (duplicate logic consolidation).

## Reference Files
- `docs/proposals/forge-cli-clean-code/proposal.md` — Source proposal
- `forge-cli/internal/cmd/claude.go` — First copy
- `forge-cli/internal/cmd/worktree/worktree.go` — Second copy (after split in task 4)

## Acceptance Criteria
- [ ] `defaultRunClaude()` exists in exactly 1 location
- [ ] Both `claude.go` and `worktree` package reference the shared definition
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes

## Hard Rules
- The shared location must be importable by both the `cmd` and `worktree` packages without circular imports

## Implementation Notes
- If the function is only used by `cmd` and `worktree` (a sub-package of `cmd`), placing it in `cmd/` or a shared `internal/` package works
- Verify the two copies are truly identical before merging
