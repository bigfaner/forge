---
id: "9"
title: "Clean up remaining tests, contracts, and verify final structure"
priority: "P0"
estimated_time: "1h"
dependencies: ["3", "4", "5", "6", "7", "8"]
scope: "backend"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 9: Clean up remaining tests, contracts, and verify final structure

## Description

Final cleanup pass: remove/update any remaining test files and contracts affected by the deletions and moves. Verify the final directory structure matches the proposal's target layout.

## Reference Files
- `docs/proposals/cli-command-restructure/proposal.md` — Source proposal

## Acceptance Criteria

- `internal/cmd/` contains only: root.go, errors.go, output.go, cleanup.go, quality_gate.go, verify_task_done.go, config.go, init.go, version.go, claude.go, proposal.go, lesson.go, and subdirectories (task/, test/, feature/, worktree/, forensic/, prompt/)
- No orphaned test files remain in `internal/cmd/` for moved commands
- `go build ./...` passes
- `go test ./...` passes
- `forge --help` shows expected commands (no e2e, no probe)
- `internal/cmd/` has no flat command-group files (only top-level commands and subdirectories)

## Hard Rules

- Do NOT change any command's runtime behavior
- Run full test suite to confirm zero regressions
- Verify no circular imports between sub-packages

## Implementation Notes

Check for and clean up:
- `integration_test.go` and `characterization_test.go` in internal/cmd — verify they don't reference deleted code
- Any testdata files related to deleted commands
- Run `grep -r "e2e\|probe" forge-cli/internal/cmd/` to verify no stale references
- Verify import graph: `go vet ./...` and architecture lint pass
- Confirm `internal/cmd/` flat file list matches target structure from proposal
