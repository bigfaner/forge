---
status: "completed"
started: "2026-05-15 22:22"
completed: "2026-05-15 22:35"
time_spent: "~13m"
---

# Task Record: 4 Update quick test task generation for per-type tasks

## Summary
Updated GetQuickTestTasks() to accept detectedTypes parameter and create per-type gen-scripts tasks (e.g., T-quick-2-tui, T-quick-2-api) instead of a single T-quick-2 per profile. Updated resolveQuickDeps() to use dynamic block-size arithmetic so T-quick-3 depends on ALL per-type T-quick-2-* tasks for its profile. Updated generateTestTasks() in build.go to pass detectedTypes to quick mode. All changes follow the same pattern as Task 3 (breakdown mode).

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/testgen.go
- forge-cli/pkg/task/testgen_test.go
- forge-cli/pkg/task/build.go

### Key Decisions
- Used same per-type pattern as breakdown mode: blockSize = nTypes + 3 (gen-cases + per-type-gen + run + graduate) for dynamic dependency resolution
- Preserved legacy code path when detectedTypes is nil/empty for backward compatibility
- resolveQuickDeps now accepts detectedTypes parameter mirroring resolveBreakdownDeps signature

## Test Results
- **Tests Executed**: Yes
- **Passed**: 173
- **Failed**: 0
- **Coverage**: 90.3%

## Acceptance Criteria
- [x] GetQuickTestTasks() creates separate tasks per type: e.g. T-quick-2-tui, T-quick-2-api, T-quick-2-cli for go-test profile
- [x] Only types with test cases in test-cases.md get tasks (no empty tasks)
- [x] T-quick-3 depends on ALL per-type T-quick-2-* tasks for its profile
- [x] Multi-profile works: T-quick-2a-tui, T-quick-2a-api, T-quick-2b-tui, T-quick-2b-api
- [x] Task keys and file names include type suffix: quick-gen-scripts-go-test-tui, quick-gen-scripts-go-test-api
- [x] GenerateTestTaskMD() produces correct .md content with profile + type info
- [x] Existing quick test tasks (T-quick-1, T-quick-3, T-quick-4, T-quick-5) are unchanged

## Notes
4 new test functions added covering single-profile/per-type, multi-profile/per-type, single-type, and three-types scenarios. All existing tests updated for new signature. Coverage at 90.3% for pkg/task.
