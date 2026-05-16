---
id: "1"
title: "Delete 4 entire test files with no meaningful tests"
priority: "P1"
estimated_time: "15m"
dependencies: []
scope: "all"
breaking: false
type: "implementation"
mainSession: false
---

# 1: Delete 4 entire test files with no meaningful tests

## Description

Delete 4 test files from `tests/e2e/` that contain zero meaningful tests:
- `extract_design_md_platform_adapters_cli_test.go` — 18 tests that read a static `.md` command file and grep for text strings. No CLI command is ever executed.
- `cli_list_reverse_chronological_cli_test.go` — byte-for-byte duplicate of `features/cli-list-reverse-chronological/cli_list_reverse_chronological_cli_test.go`.
- `fix_task_claim_priority_cli_test.go` — byte-for-byte duplicate of `features/fix-task-claim-priority/fix_task_claim_priority_cli_test.go`.
- `cli_lean_output_cli_test.go` — 19 tests: 8 vacuous assertions (`if cond { assert }` never fires), 11 conditional-skip without fixture (skips when no tasks exist).

## Reference Files
- `docs/proposals/e2e-test-quality-cleanup/proposal.md` — Source proposal

## Acceptance Criteria
- [ ] The 4 files do not exist in `tests/e2e/`
- [ ] `tests/e2e/features/cli-list-reverse-chronological/cli_list_reverse_chronological_cli_test.go` still exists (features/ copy preserved)
- [ ] `tests/e2e/features/fix-task-claim-priority/fix_task_claim_priority_cli_test.go` still exists (features/ copy preserved)
- [ ] `just test-e2e` compiles and passes (remaining tests unaffected)

## Hard Rules
- Do NOT delete files in `tests/e2e/features/` — only root-level duplicates
- Do NOT delete `tests/e2e/main_test.go` or `tests/e2e/helpers.go`

## Implementation Notes
- The duplicates exist because `graduate-tests` copied to root without removing features/ source. This task removes the root copies.
- After deletion, verify `go test -tags=e2e ./...` from `tests/e2e/` still compiles.
