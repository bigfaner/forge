---
id: "1"
title: "Remove dead code from run.go"
priority: "P0"
estimated_time: "30m"
dependencies: []
scope: "backend"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 1: Remove dead code from run.go

## Description
Delete `GetVersion()`, `GetName()`, `IsTestMode()` from `cmd/forge/run.go` — these functions have zero callers. This is Phase 1 (dead code elimination) of the cleanup proposal.

## Reference Files
- `docs/proposals/forge-cli-clean-code/proposal.md` — Source proposal
- `forge-cli/cmd/forge/run.go` — Contains dead functions

## Acceptance Criteria
- [ ] `GetVersion()`, `GetName()`, `IsTestMode()` removed from `run.go`
- [ ] `go build ./...` passes with zero errors
- [ ] `go test ./...` passes (all test files)
- [ ] `grep -r` confirms zero callers for each removed function

## Hard Rules
- Verify each function has zero callers via LSP `findReferences` or `grep -r` before deletion
- Run `go build ./...` and `go test ./...` after deletion to confirm no breakage

## Implementation Notes
- These are simple utility functions with no side effects — safe to delete
- Risk: audit data could be stale. Re-verify call count before deleting
