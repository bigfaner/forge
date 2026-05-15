---
status: "completed"
started: "2026-05-15 22:04"
completed: "2026-05-15 22:21"
time_spent: "~17m"
---

# Task Record: 3 Update breakdown test task generation for per-type tasks

## Summary
Modified GetBreakdownTestTasks() to accept detectedTypes parameter and create per-type gen-scripts tasks (e.g., T-test-2-tui, T-test-2-api) instead of a single T-test-2 per profile. Updated resolveBreakdownDeps() so T-test-3 depends on ALL per-type T-test-2-* tasks. Added TestType field to TestTaskDef. Updated GenerateTestTaskMD() to include type info in .md output. Updated generateTestTasks() and BuildIndex call site. Legacy behavior preserved when detectedTypes is nil.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/testgen.go
- forge-cli/pkg/task/testgen_test.go
- forge-cli/pkg/task/build.go
- forge-cli/scripts/version.txt

### Key Decisions
- Added detectedTypes []string parameter to GetBreakdownTestTasks() - nil preserves legacy single-gen-scripts behavior
- Per-type tasks use key format gen-test-scripts-<profile>-<type> and ID format T-test-2<suffix>-<type>
- T-test-3 depends on ALL per-type T-test-2-* tasks for its profile (fan-in dependency)
- resolveBreakdownDeps uses blockSize = nTypes + 2 to handle variable per-type gen task counts
- BuildIndex passes nil for detectedTypes now; type detection wiring is a separate concern

## Test Results
- **Tests Executed**: Yes
- **Passed**: 10
- **Failed**: 0
- **Coverage**: 90.0%

## Acceptance Criteria
- [x] GetBreakdownTestTasks() creates separate tasks per type: e.g., T-test-2-tui, T-test-2-api, T-test-2-cli for go-test profile
- [x] Only types with test cases in test-cases.md get tasks (no empty tasks for types without cases)
- [x] T-test-3 depends on ALL per-type T-test-2-* tasks for its profile
- [x] Multi-profile works: T-test-2a-tui, T-test-2a-api, T-test-2b-tui, T-test-2b-api
- [x] Task keys and file names include type suffix: gen-test-scripts-go-test-tui, gen-test-scripts-go-test-api
- [x] GenerateTestTaskMD() produces correct .md content with profile + type info
- [x] Existing breakdown test tasks (T-test-1, T-test-1b, T-test-3, T-test-4, T-test-4.5, T-test-5) are unchanged

## Notes
Legacy behavior fully preserved when detectedTypes is nil/empty. Quick mode (GetQuickTestTasks) not modified - that is task 4.
