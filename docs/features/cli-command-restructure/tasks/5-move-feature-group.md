---
id: "5"
title: "Move forge feature group to feature/ subdirectory"
priority: "P2"
estimated_time: "1h"
dependencies: ["2"]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 5: Move forge feature group to feature/ subdirectory

## Description

Move all forge feature subcommand files from `internal/cmd/` to `internal/cmd/feature/` subdirectory. Create a `Register()` function for the feature command group.

## Reference Files
- `docs/proposals/cli-command-restructure/proposal.md` — Source proposal

## Acceptance Criteria

- All feature subcommand files are in `internal/cmd/feature/`
- New package exports a `Register()` function
- root.go updated to use the new package
- `go build ./...` passes
- `go test ./...` passes
- `forge feature` subcommands work identically

## Hard Rules

- The feature sub-package must NOT import `internal/cmd` (no circular deps)

## Implementation Notes

Files to move to `internal/cmd/feature/`:
- feature.go, feature_test.go
- feature_complete.go, feature_complete_test.go
