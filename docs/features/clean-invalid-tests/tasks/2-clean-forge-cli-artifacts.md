---
id: "2"
title: "Clean forge-cli/tests/ empty artifacts and dead helpers"
priority: "P1"
estimated_time: "45m"
dependencies: [1]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 2: Clean forge-cli/tests/ empty artifacts and dead helpers

## Description
After Task 1 removes invalid tests, clean up the remaining debris: delete test files that became empty, remove empty suite directories, and remove testkit helper functions that are no longer referenced by any surviving test.

## Reference Files
- `docs/proposals/clean-invalid-tests/proposal.md#Scope` — In Scope items 4-6 defining empty file, empty dir, and dead helper cleanup
- `docs/proposals/clean-invalid-tests/proposal.md#Key-Risks` — risk of deleting helpers still indirectly referenced; mitigation is grep verification
- `docs/proposals/clean-invalid-tests/proposal.md#Success-Criteria` — "清理后无空测试文件" and "无空 suite 目录" and "testkit 辅助函数无死引用"

## Acceptance Criteria
- [ ] No empty `*_test.go` files in `forge-cli/tests/` (verify: `find forge-cli/tests/ -name '*_test.go' -empty`)
- [ ] No empty suite directories (verify: `find forge-cli/tests/ -type d -empty`)
- [ ] `forge-cli/tests/testkit/` helper functions all referenced by surviving tests (verify: grep each exported func)
- [ ] `go build ./forge-cli/tests/...` compiles successfully
- [ ] If testkit itself becomes empty, delete the entire `testkit/` directory

## Hard Rules
- Before deleting any testkit helper function, `grep` across ALL `forge-cli/tests/` to confirm no remaining references
- Only delete empty directories — do not remove directories containing non-test files (e.g., testdata)

## Implementation Notes
- Check if `main_test.go` files in emptied suites are still needed — if the suite has no test functions, `main_test.go` (typically containing `TestMain`) is dead code too
- A suite directory is "empty" if it contains only `main_test.go` with a `TestMain` and no other test functions
- For testkit helpers, check both direct function calls and any references in struct literals or interface implementations
