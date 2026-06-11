---
status: "completed"
started: "2026-05-15 22:39"
completed: "2026-05-15 22:49"
time_spent: "~10m"
---

# Task Record: 3 Update breakdown test task generation for per-type tasks

## Summary
Wire DetectTypesFromTestCases into BuildIndex so per-type gen-scripts tasks are created when test-cases.md has non-zero counts for specific types. Added DetectTypesFromTestCases function that parses the Summary table in test-cases.md to extract type names with non-zero counts, and updated BuildIndex to call it and pass detected types to generateTestTasks instead of nil.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/testgen.go
- forge-cli/pkg/task/testgen_test.go
- forge-cli/pkg/task/build.go

### Key Decisions
- DetectTypesFromTestCases parses the Summary table (not individual test cases) because the table is canonical and already aggregates per-type counts
- Types with count > 0 are included; types with count 0 or absent are excluded, satisfying the hard rule of no empty tasks
- Reading test-cases.md is non-fatal: if file missing, detectedTypes stays nil and legacy single-gen-scripts behavior is preserved

## Test Results
- **Tests Executed**: Yes
- **Passed**: 174
- **Failed**: 0
- **Coverage**: 90.3%

## Acceptance Criteria
- [x] GetBreakdownTestTasks() creates separate tasks per type
- [x] Only types with test cases in test-cases.md get tasks
- [x] T-test-3 depends on ALL per-type T-test-2-* tasks for its profile
- [x] Multi-profile works: T-test-2a-tui, T-test-2a-api, T-test-2b-tui, T-test-2b-api
- [x] Task keys and file names include type suffix
- [x] GenerateTestTaskMD() produces correct .md content with profile + type info
- [x] Existing breakdown test tasks unchanged

## Notes
The per-type task generation logic (GetBreakdownTestTasks, resolveBreakdownDeps, GetQuickTestTasks, resolveQuickDeps) was already implemented in a prior task. This task added the missing piece: reading test-cases.md to detect present types and passing them through BuildIndex.
