---
id: "3"
title: "Move forge task group to task/ subdirectory"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["2"]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 3: Move forge task group to task/ subdirectory

## Description

Move all forge task subcommand files from `internal/cmd/` to `internal/cmd/task/` subdirectory. This is the largest command group (~20 files). Create a `Register()` function to inject task commands into the parent command from the new package.

## Reference Files
- `docs/proposals/cli-command-restructure/proposal.md` — Source proposal

## Acceptance Criteria

- All task subcommand files are in `internal/cmd/task/`
- New package `task` has a `Register(parent *cobra.Command)` function
- root.go imports the task package and calls `task.Register(taskCmd)` (or similar pattern)
- `go build ./...` passes
- `go test ./...` passes
- `forge task` subcommands all work identically

## Hard Rules

- The task sub-package must NOT import `internal/cmd` (no circular deps)
- Shared utilities (errors.go, output.go) stay in cmd — task package receives them as parameters if needed, or the Register pattern passes them in
- All task-related tests move to the subdirectory

## Implementation Notes

Files to move to `internal/cmd/task/`:
- add.go, add_cmd_test.go
- claim.go, claim_integration_test.go, claim_test.go
- check_deps.go, check_deps_test.go
- index.go, index_test.go
- list_types.go, list_types_test.go
- migrate.go, migrate_test.go
- query.go, query_test.go
- reopen.go
- runners_test.go (if task-specific)
- slug_width_test.go (if task-specific)
- status.go, status_test.go
- submit.go
- task_parent.go
- transition.go
- validate_index.go, validate_index_test.go

Pattern: Create `Register(rootCmd *cobra.Command)` in the task package. In root.go, replace direct variable references with `task.Register(rootCmd)`. Remove the "Task group subcommands" block from root.go init().
