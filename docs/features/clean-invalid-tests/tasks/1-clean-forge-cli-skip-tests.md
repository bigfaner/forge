---
id: "1"
title: "Clean forge-cli/tests/ skip and dead tests"
priority: "P0"
estimated_time: "1h"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 1: Clean forge-cli/tests/ skip and dead tests

## Description
Delete all invalid integration tests in `forge-cli/tests/`: 44 unconditional `t.Skip` tests across 10 files, tests referencing deleted features/retired skills (`test-generation/gen_test_scripts_test.go`), and tests requiring interactive stdin that cannot be automated (`forge_info_commands_test.go` 4 skip tests).

These tests provide zero quality assurance value. Skip tests never execute; dead-references tests fail against non-existent code; interactive tests cannot run in CI.

## Reference Files
- `docs/proposals/clean-invalid-tests/proposal.md#Evidence` — table showing 44 skip tests across 10 files, 3 dead-reference tests, 4 interactive tests
- `docs/proposals/clean-invalid-tests/proposal.md#Scope` — In Scope items 1-3 defining exact deletion targets
- `docs/proposals/clean-invalid-tests/proposal.md#Key-Risks` — deletion risk mitigation (git history preserves everything)

## Acceptance Criteria
- [ ] Zero unconditional `t.Skip` remain in `forge-cli/tests/` (verify: `grep -r 't.Skip(' forge-cli/tests/ --include='*_test.go' | grep -v '_test.go:\s*//'`)
- [ ] `forge-cli/tests/test-generation/gen_test_scripts_test.go` deleted (or cleaned of dead references)
- [ ] Interactive-input skip tests removed from `forge-cli/tests/forge-commands/forge_info_commands_test.go`
- [ ] `go build ./forge-cli/tests/...` compiles successfully after deletions
- [ ] Remaining tests in modified files still compile and structurally valid

## Hard Rules
- Do NOT delete non-skip tests in affected files — only remove the skip functions themselves
- After deleting test functions, verify the file still compiles (import cleanup may be needed)

## Implementation Notes
- Start by identifying all `t.Skip(` calls in `forge-cli/tests/` — some may be conditional (inside `if` blocks or `Skipf`), only delete unconditional ones
- For files that become empty after deletions, leave them for Task 2 to clean up
- Run `go build ./forge-cli/tests/...` after each batch of deletions to catch compilation errors early
