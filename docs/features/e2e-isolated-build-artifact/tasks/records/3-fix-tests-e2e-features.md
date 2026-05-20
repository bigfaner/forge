---
status: "blocked"
started: "2026-05-20 17:00"
completed: "N/A"
time_spent: ""
---

# Task Record: 3 Fix tests/e2e/ feature tests to use TestMain-built binary

## Summary
Unified all e2e tests under tests/e2e/ to use TestMain-built isolated binary instead of system-installed forge via PATH. Created forge_binary.go with exported ForgeBinary var and init() builder. Updated 7 files in top-level package and 4 sub-packages. Removed duplicate TestMain from test-knowledge-convention-driven.

## Changes

### Files Created
- tests/e2e/forge_binary.go

### Files Modified
- tests/e2e/main_test.go
- tests/e2e/task_record_immutability_cli_test.go
- tests/e2e/quality_gate_fix_task_loop_breaker_cli_test.go
- tests/e2e/feature_set_command_cli_test.go
- tests/e2e/features/cli-list-reverse-chronological/cli_list_reverse_chronological_cli_test.go
- tests/e2e/features/proposal-status-lifecycle/proposal_status_lifecycle_cli_test.go
- tests/e2e/features/fix-task-claim-priority/fix_task_claim_priority_cli_test.go
- tests/e2e/features/task-type-refinement/task_type_refinement_cli_test.go
- tests/e2e/features/test-knowledge-convention-driven/test_helpers_test.go
- tests/e2e/features/test-knowledge-convention-driven/forge_commands_cli_test.go
- tests/e2e/features/test-knowledge-convention-driven/gen_test_scripts_cli_test.go
- tests/e2e/features/test-knowledge-convention-driven/integration_cli_test.go
- tests/e2e/features/test-knowledge-convention-driven/removed_commands_cli_test.go
- tests/e2e/features/test-knowledge-convention-driven/test_guide_cli_test.go

### Key Decisions
- Created forge_binary.go (non-_test.go) with exported ForgeBinary and ForgeCmd so sub-packages can import via e2etests alias
- Used sync.Once in init() to build binary exactly once, shared across parent and child packages
- Sub-packages renamed from package e2e to unique names (e2eclilistr, e2epsl, e2efixclaim, e2etasktype, e2etestconv) to allow importing parent package
- Removed duplicate TestMain and forgeCLIDir from test-knowledge-convention-driven/test_helpers_test.go, replaced with simple forgeCmd helper using e2etests.ForgeBinary

## Test Results
- **Tests Executed**: Yes
- **Passed**: 7
- **Failed**: 1
- **Coverage**: 66.7%

## Acceptance Criteria
- [x] All test files under tests/e2e/ use forgeBinary variable from main_test.go instead of bare 'forge'
- [x] grep -r 'exec.Command("forge"' returns zero results (excluding justfile-canonical-e2e)
- [x] Duplicate TestMain in test-knowledge-convention-driven/ removed; package uses parent's forgeBinary
- [x] Walk-up path resolution eliminated from all test files

## Notes
无
