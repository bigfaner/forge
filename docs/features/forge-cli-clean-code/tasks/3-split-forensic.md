---
id: "3"
title: "Split forensic.go into functional files"
priority: "P1"
estimated_time: "2h"
dependencies: [2]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 3: Split forensic.go into functional files

## Description
Split `forensic.go` (993 lines, 20+ structs, 3 commands, 300-line nested functions) into focused files by responsibility. This is Phase 2 (oversized file splitting) of the cleanup proposal.

## Reference Files
- `docs/proposals/forge-cli-clean-code/proposal.md` — Source proposal
- `forge-cli/internal/cmd/forensic/forensic.go` — Source file to split
- `forge-cli/internal/cmd/forensic/register.go` — Existing registration file

## Acceptance Criteria
- [ ] `forensic.go` reduced to <300 lines
- [ ] New files created with single-responsibility groupings (e.g. types.go, search.go, extract.go, subagents.go)
- [ ] All structs, functions, and types remain in the same package (`forensic`)
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes
- [ ] No behavioral changes — pure file reorganization

## Hard Rules
- All symbols must remain in the `forensic` package — no package renames or moves
- Each commit must pass `go build ./...` and `go test ./...`
- Preserve all existing test files unchanged

## Implementation Notes
- Start by listing all type definitions and function groups to plan the split
- `register.go` already exists for command registration — don't duplicate
- 20+ struct definitions suggest a `types.go` file for data types
- 300-line nested functions should be extracted as standalone functions in their own file
