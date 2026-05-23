---
id: "7"
title: "Move forge forensic group to forensic/ subdirectory"
priority: "P2"
estimated_time: "1h"
dependencies: ["2"]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 7: Move forge forensic group to forensic/ subdirectory

## Description

Move all forge forensic subcommand files from `internal/cmd/` to `internal/cmd/forensic/` subdirectory. Create a `Register()` function for the forensic command group.

## Reference Files
- `docs/proposals/cli-command-restructure/proposal.md` — Source proposal

## Acceptance Criteria

- All forensic subcommand files are in `internal/cmd/forensic/`
- New package exports a `Register()` function
- root.go updated to use the new package
- `go build ./...` passes
- `go test ./...` passes
- `forge forensic` behavior unchanged

## Hard Rules

- The forensic sub-package must NOT import `internal/cmd` (no circular deps)

## Implementation Notes

Files to move to `internal/cmd/forensic/`:
- forensic.go, forensic_test.go
