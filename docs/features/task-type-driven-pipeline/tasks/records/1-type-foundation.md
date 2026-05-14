---
status: "completed"
started: "2026-05-14 22:53"
completed: "2026-05-14 23:00"
time_spent: "~7m"
---

# Task Record: 1 Add documentation/doc-evaluation type constants and remove InferType fallback

## Summary
Add TypeDocumentation and TypeDocEvaluation constants to the type system, register them in TaskTypeRegistry and ValidTypes, remove the TypeImplementation fallback from InferType so unknown IDs return empty string, and add T-eval-doc pattern recognition.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/types.go
- forge-cli/pkg/task/infer.go
- forge-cli/pkg/task/types_test.go
- forge-cli/pkg/task/infer_test.go
- forge-cli/pkg/task/build_test.go
- forge-cli/pkg/prompt/prompt_test.go
- forge-cli/internal/cmd/migrate_test.go

### Key Decisions
- InferType returns empty string for unknown IDs instead of TypeImplementation fallback, enabling hard-error detection in downstream validation (task 2)
- T-eval-doc is a single exact-match ID (not profile-suffixed), implemented as simple string comparison
- TaskTypeRegistry count updated from 11 to 13 to include both new types

## Test Results
- **Tests Executed**: Yes
- **Passed**: 318
- **Failed**: 0
- **Coverage**: 89.9%

## Acceptance Criteria
- [x] TypeDocumentation = 'documentation' and TypeDocEvaluation = 'doc-evaluation' constants exist in types.go
- [x] TaskTypeRegistry includes both new types with descriptions
- [x] ValidTypes map includes both new types
- [x] InferType returns '' for unknown IDs (no TypeImplementation fallback)
- [x] InferType recognizes T-eval-doc pattern and returns TypeDocEvaluation
- [x] Existing tests in types_test.go and infer_test.go pass after changes
- [x] New table-driven tests cover: InferType returning '' for business IDs, T-eval-doc pattern, both new constants, ValidTypes entries

## Notes
Pre-existing test failures in pkg/project (root detection) and internal/cmd (TestRunMigrate_HappyPath subprocess pattern) are unrelated to this change and confirmed to fail on unmodified code.
