---
id: "2"
title: "Remove vacuous, recursive, and dead-skip test functions from 3 files"
priority: "P1"
estimated_time: "30m"
dependencies: ["1"]
scope: "all"
breaking: false
type: "implementation"
mainSession: false
---

# 2: Remove vacuous, recursive, and dead-skip test functions from 3 files

## Description

Remove specific test functions from 3 files that remain after Task 1:

**`tests/e2e/simplify_e2e_tests_cli_test.go`** — remove TC-003 and TC-004:
- `TestTC_003_VerifyE2eTestSuiteCompiles` — spawns `go test -tags=e2e ./... -run=^$` from within a test (recursive)
- `TestTC_004_VerifyRemainingCliBehaviorTestsPass` — spawns `go test -tags=e2e ./features/... -timeout 300s` from within a test (recursive, causes process explosion on Windows)

**`tests/e2e/feature_set_command_cli_test.go`** — remove TC-016 and TC-017:
- `TestTC_016_VerboseShowsWorktreeSource` — unconditional `t.Skip("requires real git worktree environment setup")`
- `TestTC_017_VerboseShowsBranchSource` — unconditional `t.Skip("requires real git feature branch environment setup")`

**`tests/e2e/quick_test_slim_cli_test.go`** — remove TC-003, TC-009, TC-010, TC-013, TC-016:
- `TestTC_003_QuickModeMergedTaskPromptTemplate` — reads static template file, greps for "gen-test-scripts"/"run-e2e-tests"
- `TestTC_009_MergedTemplateCallsBothSkills` — reads same template file, greps same strings (duplicates TC-003)
- `TestTC_010_TypeConstantRegisteredInTypesGo` — reads `types.go` source, greps for constant name
- `TestTC_013_ExistingGenScriptsAndRunTemplatesRemainIntact` — reads template files, checks they exist and are non-empty
- `TestTC_016_SubagentCompletesGenAndRunInSingleSession` — reads template file, greps for "Step 1/2/3", "combined" (superset of TC-003/009)

## Reference Files
- `docs/proposals/e2e-test-quality-cleanup/proposal.md` — Source proposal

## Acceptance Criteria
- [ ] `simplify_e2e_tests_cli_test.go` contains only TC-001 and TC-002 (plus helper functions)
- [ ] `feature_set_command_cli_test.go` has no unconditional `t.Skip` calls
- [ ] `quick_test_slim_cli_test.go` has no tests that read static source files (no `os.ReadFile` on `.go` or template files outside the temp test project)
- [ ] All remaining tests in the 3 files compile and pass via `just test-e2e`
- [ ] No orphaned imports after removal (e.g., `runtime` in simplify_e2e_tests if only used by deleted functions — but it's also used by `projectRoot`, so keep it)

## Hard Rules
- Keep all helper functions (`projectRoot`, `e2eRoot`, `quickSlim*` helpers) — they may be used by remaining tests
- Keep TC-001 and TC-002 in simplify_e2e_tests — they verify real behavioral outcomes (directory deleted, function removed)
- Keep TC-018/TC-019/TC-020 in feature_set_command — they test real CLI behavior with self-contained fixtures

## Implementation Notes
- When removing a test function, also remove the comment block above it (the `// ===` separator and traceability comment)
- For quick_test_slim, the removed tests all read files from the project source tree. The remaining tests (TC-001, TC-002, TC-004-008, TC-011, TC-012, TC-014, TC-015) all run `forge task index` and verify the generated output — these are real integration tests.
