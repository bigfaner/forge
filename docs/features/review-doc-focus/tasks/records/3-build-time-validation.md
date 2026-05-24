---
status: "completed"
started: "2026-05-24 19:22"
completed: "2026-05-24 19:31"
time_spent: "~9m"
---

# Task Record: 3 Build-time AC validation (warnings + empty AC + title tolerance)

## Summary
Added build-time AC validation: title matching tolerance (case-insensitive + Chinese alias), warning logs for missing AC, placeholder text for empty AC in review-doc, zero-AC feature warning, and DocTaskCriteria key-set verification in BuildIndex()

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/extract.go
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/autogen.go
- forge-cli/pkg/task/extract_test.go
- forge-cli/pkg/task/build_test.go

### Key Decisions
- Used strings.EqualFold for case-insensitive AC heading match plus hardcoded Chinese alias check in isACHeading()
- extractDocTaskCriteria now includes doc tasks with empty AC content (instead of skipping) so callers can emit per-task warnings
- serializeDocTaskAC shows '> No acceptance criteria defined.' placeholder for empty entries
- BuildIndex step 5.5.2 validates DocTaskCriteria coverage and emits both per-task and feature-level warnings

## Test Results
- **Tests Executed**: Yes
- **Passed**: 10
- **Failed**: 0
- **Coverage**: 87.8%

## Acceptance Criteria
- [x] Title matching supports ## Acceptance Criteria, ## Acceptance criteria, ## 验收标准
- [x] Section missing outputs warning log [WARN] task <name> has no Acceptance Criteria section
- [x] Empty AC shows > No acceptance criteria defined. in summary area
- [x] All doc tasks missing AC outputs [WARN] feature has no AC for any doc task
- [x] BuildIndex() validates DocTaskCriteria key set matches doc task list

## Notes
All existing tests pass with no regressions. Coverage 87.8% exceeds 60% target.
