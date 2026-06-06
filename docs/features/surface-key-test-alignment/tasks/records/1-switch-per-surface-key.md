---
status: "completed"
started: "2026-06-06 12:54"
completed: "2026-06-06 13:17"
time_spent: "~23m"
---

# Task Record: 1 Switch gen-test-scripts to per-surface-key expansion in pipeline.go

## Summary
Switched gen-test-scripts pipeline registry entry from per-surface-type to per-surface-key expansion, aligning it with run-tests. Updated Key/ID templates from {surface-type} to {surface-key}. Fixed isTestTaskID to handle per-surface-key prefix matching. Updated migrate command to load surfaces config for type inference. Updated all affected tests.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/pipeline.go
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/autogen_test.go
- forge-cli/pkg/task/pipeline_test.go
- forge-cli/pkg/task/infer_test.go
- forge-cli/pkg/task/autoconfig_test.go
- forge-cli/internal/cmd/task/migrate.go
- forge-cli/internal/cmd/task/migrate_test.go

### Key Decisions
- Added isAutoGenRegistryPrefix fallback to isTestTaskID so per-surface-key IDs like T-test-gen-scripts-backend are recognized without a surfaces map
- Updated migrate.go to load surfaces config before type inference, fixing InferType for per-surface-key IDs
- Gen-scripts serial chain: first task depends on eval-contract, subsequent tasks chain to predecessor (matching run-tests pattern)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 189
- **Failed**: 0
- **Coverage**: 86.4%

## Acceptance Criteria
- [x] gen-test-scripts registry entry Expansion field changed from per-surface-type to per-surface-key
- [x] Key and ID templates changed from {surface-type} to {surface-key}: gen-test-scripts-{surface-key}, T-test-gen-scripts-{surface-key}
- [x] Single surface project produces suffix-less gen-test-scripts.md via isSingleSurface logic (regression verified)
- [x] go test ./pkg/task/... all pass

## Notes
Lint failure (upgrade_test.go unused parameter) is pre-existing and unrelated to this change.
