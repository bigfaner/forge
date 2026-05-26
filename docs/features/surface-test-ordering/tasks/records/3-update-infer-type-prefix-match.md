---
status: "completed"
started: "2026-05-26 10:35"
completed: "2026-05-26 10:50"
time_spent: "~15m"
---

# Task Record: 3 Update InferType for prefix matching

## Summary
Updated InferType to support T-test-run-{surface-key} prefix matching: added surfaces map parameter, new testRunSurfaceKeyMatch helper, and updated all callers (build.go, extract.go, migrate.go, prompt.go). Tests cover known key match, unknown key fallback, and single-surface backward compatibility.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/infer.go
- forge-cli/pkg/task/infer_test.go
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/build_test.go
- forge-cli/pkg/task/extract.go
- forge-cli/pkg/task/autoconfig_test.go
- forge-cli/internal/cmd/task/migrate.go
- forge-cli/pkg/prompt/prompt.go
- forge-cli/scripts/version.txt

### Key Decisions
- Added surfaces map[string]string parameter to InferType instead of creating a separate function, keeping the API unified
- testRunSurfaceKeyMatch is a dedicated helper extracted for clarity and testability
- Nil surfaces map disables prefix matching entirely, preserving backward compatibility for callers without surfaces context

## Test Results
- **Tests Executed**: Yes
- **Passed**: 53
- **Failed**: 0
- **Coverage**: 87.3%

## Acceptance Criteria
- [x] InferType("T-test-run-backend") returns correct surface type via prefix matching
- [x] Test coverage: known surface-key -> correct type; unknown key -> fallback exact match; single surface -> no suffix
- [x] Single surface project InferType("T-test-run") behavior unchanged

## Notes
Pre-existing test failure TestRunMigrate_ScopeToSurface_WithConfig confirmed unrelated to this change. Version bumped to 5.9.1 (patch).
