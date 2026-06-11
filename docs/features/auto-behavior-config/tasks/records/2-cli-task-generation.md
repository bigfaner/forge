---
status: "completed"
started: "2026-05-17 01:59"
completed: "2026-05-17 02:27"
time_spent: "~28m"
---

# Task Record: 2 Implement config-driven task generation with mode scoping

## Summary
Implemented config-driven task generation with mode scoping for the auto-behavior-config feature. Added AutoConfig/ModeToggle structs to profile package with YAML-aware default filling (e2eTest defaults true, consolidateSpecs defaults true, cleanCode defaults false). Modified GetBreakdownTestTasks and GetQuickTestTasks to accept AutoConfig and conditionally gate task generation per mode. Renamed T-test-5 to T-specs-1 and T-quick-5 to T-quick-specs-1. Added TypeCleanCode constant and T-clean-code-1 task type. Created code-quality-simplify.md prompt template. Updated BuildIndex to read AutoConfig from config and pass through to task generation. Broke import cycle by mirroring path constants in profile package. All existing tests updated and passing.

## Changes

### Files Created
- forge-cli/pkg/profile/autoconfig_test.go
- forge-cli/pkg/task/autoconfig_test.go
- forge-cli/pkg/prompt/data/code-quality-simplify.md

### Files Modified
- forge-cli/pkg/profile/config.go
- forge-cli/pkg/task/types.go
- forge-cli/pkg/task/infer.go
- forge-cli/pkg/task/testgen.go
- forge-cli/pkg/task/build.go
- forge-cli/internal/cmd/index.go
- forge-cli/pkg/task/testgen_test.go
- forge-cli/pkg/task/infer_test.go
- forge-cli/pkg/task/types_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- Used yaml.Node tree parsing (parseAutoRaw) to detect which sub-fields were explicitly present in YAML, solving the Go zero-value ambiguity for bool fields
- Mirrored .forge/config.yaml path constants in profile package instead of importing from feature package to break the import cycle (feature_test -> task -> profile -> feature)
- WithDefaults() on AutoConfig handles both fully-zero (no config) and partially-set cases, ensuring backward compatibility
- isTestTaskID extended to recognize T-clean-code- prefix for auto-generated task exclusion

## Test Results
- **Tests Executed**: Yes
- **Passed**: 671
- **Failed**: 0
- **Coverage**: 89.5%

## Acceptance Criteria
- [x] auto.e2eTest.quick=false -> quick mode generates zero T-quick test tasks
- [x] auto.e2eTest.full=false -> full mode generates zero T-test tasks
- [x] auto.consolidateSpecs.quick=false -> no T-quick-specs-1 generated
- [x] auto.consolidateSpecs.full=false -> no T-specs-1 generated
- [x] auto.cleanCode.quick=true -> T-clean-code-1 generated in quick mode
- [x] auto.cleanCode.full=true -> T-clean-code-1 generated in full mode
- [x] Projects without auto config behave identically to before (backward compat)
- [x] T-test-5 renamed to T-specs-1 in all Go code
- [x] T-quick-5 renamed to T-quick-specs-1 in all Go code
- [x] New type constant + prompt template for T-clean-code-1
- [x] T-clean-code-1 depends on last business task; first test task depends on T-clean-code-1 (when both exist)

## Notes
Breaking change: existing index.json files with T-test-5/T-quick-5 IDs are incompatible. Document in CHANGELOG. The gitPush config field is stored but not yet consumed by run-tasks (out of scope for this task).
