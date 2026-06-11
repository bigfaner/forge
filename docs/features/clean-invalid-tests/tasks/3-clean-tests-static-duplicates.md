---
id: "3"
title: "Clean tests/ static-file grep and duplicate tests"
priority: "P1"
estimated_time: "1h"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 3: Clean tests/ static-file grep and duplicate tests

## Description
Delete two categories of invalid e2e tests in `tests/`:

1. **Static-file grep tests** (23 total): Tests that read `.md`/`types.go` files and check for string content via regex/grep. These never invoke any command and provide zero runtime verification. Includes all 18 tests in `extract_design_md` suite and 5 tests in `quick_test_slim`.

2. **Duplicate root copies** (12 total): Tests that were graduated (moved to `.graduated/`) but the source files were never deleted. Includes 6 tests in `cli_list_reverse_chronological` and 6 in `fix_task_claim_priority`.

## Reference Files
- `docs/proposals/clean-invalid-tests/proposal.md#Evidence` — table row "读静态文件 grep 文本" (23 tests) and "重复测试（root 副本）" (12 tests)
- `docs/proposals/clean-invalid-tests/proposal.md#Scope` — In Scope items 7-8
- `docs/proposals/clean-invalid-tests/proposal.md#Assumptions-Challenged` — "文本验证测试提供了某种覆盖" challenged: removing them won't miss any bugs

## Acceptance Criteria
- [ ] Zero tests in `tests/` that read static source files for text matching (verify: `grep -rn 'os.ReadFile\|ioutil.ReadFile' tests/ --include='*_test.go' | grep -v test-suite-health`)
- [ ] `tests/feature-management/cli_list_reverse_chronological_test.go` deleted (duplicate of graduated version)
- [ ] Duplicate `fix_task_claim_priority` tests removed from root `tests/` directory
- [ ] Empty files and directories cleaned up after deletions
- [ ] `go build -tags=e2e ./tests/...` compiles successfully

## Hard Rules
- Do NOT delete tests in `tests/.graduated/` — only delete root-level duplicates
- Do NOT delete `tests/test-suite-health/` meta-tests — those enforce quality invariants
- Verify each deletion target is actually a duplicate by checking the corresponding `.graduated/` entry exists

## Implementation Notes
- For `extract_design_md` tests: likely in `tests/.graduated/extract-design-md-platform-adapters/` or referenced from there. The proposal says 18 tests read static files — find and delete the test functions or entire files
- For `quick_test_slim`: `tests/test-generation/quick_test_slim_test.go` — delete the 5 grep-text test functions; if all functions in file are grep-text, delete entire file
- For duplicates: `tests/feature-management/cli_list_reverse_chronological_test.go` is the root copy; check `.graduated/` for the graduated version before deleting
- After deletions, check if parent directories became empty and clean up
