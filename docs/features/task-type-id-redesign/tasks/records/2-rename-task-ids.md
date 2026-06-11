---
status: "completed"
started: "2026-05-19 14:17"
completed: "2026-05-19 14:42"
time_spent: "~25m"
---

# Task Record: 2 Rename task IDs and update generation/inference

## Summary
Renamed all auto-generated task IDs from numeric format (T-test-1, T-quick-2, etc.) to readable format (T-test-gen-cases, T-quick-gen-and-run, etc.) per the ID mapping table. Updated InferType() pattern matching in infer.go with new base strings. Added Validation ModeToggle to AutoConfig with default false. Added T-validate-code and T-validate-ux task generation gated behind auto.Validation config. Updated dependency chains, genScriptBases, validate_index.go, quality_gate.go, and all test files to use new IDs. Full-text searched all Go files to ensure no orphaned references.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/testgen.go
- forge-cli/pkg/task/infer.go
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/types.go
- forge-cli/pkg/profile/config.go
- forge-cli/pkg/prompt/prompt.go
- forge-cli/internal/cmd/validate_index.go
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/pkg/task/infer_test.go
- forge-cli/pkg/task/testgen_test.go
- forge-cli/pkg/task/build_test.go
- forge-cli/pkg/task/frontmatter_test.go
- forge-cli/pkg/task/autoconfig_test.go
- forge-cli/pkg/task/stage_gates_test.go
- forge-cli/pkg/task/add_test.go
- forge-cli/pkg/task/index_test.go
- forge-cli/pkg/profile/autoconfig_test.go
- forge-cli/pkg/prompt/prompt_test.go
- forge-cli/internal/cmd/migrate_test.go
- forge-cli/internal/cmd/validate_index_test.go
- forge-cli/internal/cmd/claim_test.go
- forge-cli/internal/cmd/submit_test.go
- forge-cli/internal/cmd/status_test.go
- forge-cli/internal/cmd/quality_gate_test.go
- forge-cli/internal/cmd/integration_test.go
- forge-cli/tests/e2e/spec_drift_detection_cli_test.go
- forge-cli/tests/e2e/features/task-stage-gates/task_stage_gates_cli_test.go

### Key Decisions
- Profile suffix pattern preserved: single lowercase letter appended directly to base ID (e.g., T-test-gen-scriptsa for profile 'a')
- Type suffix uses hyphen separator after optional profile letter (e.g., T-test-gen-scriptsa-api)
- Validation tasks gated behind auto.Validation config (defaults to false for both Quick and Full)
- T-validate-code placed after test pipeline verify-regression, before specs-consolidate
- T-validate-ux placed after T-validate-code, before specs-consolidate; MainSession: true

## Test Results
- **Tests Executed**: Yes
- **Passed**: 1202
- **Failed**: 0
- **Coverage**: 89.2%

## Acceptance Criteria
- [x] All task IDs renamed per mapping table
- [x] InferType("T-test-gen-cases") returns "test.gen-cases"
- [x] InferType("T-validate-code") returns "validation.code"
- [x] profileSuffixedID/typeSuffixedID base strings updated to new ID prefixes
- [x] Dependency chain references use new IDs
- [x] T-validate-code generated when auto.Validation enabled
- [x] T-validate-ux generated when auto.Validation enabled

## Notes
All acceptance criteria met. T-validate-code and T-validate-ux are both generated when auto.Validation is enabled. Dependency wiring: T-validate-code depends on T-test-verify-regression (breakdown) or T-quick-verify-regression (quick); T-validate-ux depends on T-validate-code.
