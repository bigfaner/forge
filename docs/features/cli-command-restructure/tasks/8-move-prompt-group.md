---
id: "8"
title: "Move forge prompt group to prompt/ subdirectory"
priority: "P2"
estimated_time: "1h"
dependencies: ["2"]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 8: Move forge prompt group to prompt/ subdirectory

## Description

Move all forge prompt subcommand files from `internal/cmd/` to `internal/cmd/prompt/` subdirectory. Create a `Register()` function for the prompt command group.

## Reference Files
- `docs/proposals/cli-command-restructure/proposal.md` — Source proposal

## Acceptance Criteria

- All prompt subcommand files are in `internal/cmd/prompt/`
- New package exports a `Register()` function
- root.go updated to use the new package
- `go build ./...` passes
- `go test ./...` passes
- `forge prompt` behavior unchanged

## Hard Rules

- The prompt sub-package must NOT import `internal/cmd` (no circular deps)

## Implementation Notes

Files to move to `internal/cmd/prompt/`:
- prompt_get.go, prompt_parent.go, prompt_test.go
