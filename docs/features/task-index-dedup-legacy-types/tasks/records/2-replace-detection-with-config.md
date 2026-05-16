---
status: "completed"
started: "2026-05-16 23:38"
completed: "2026-05-17 00:03"
time_spent: "~25m"
---

# Task Record: 2 Replace test-type detection with config-driven capabilities and remove legacy code

## Summary
Replace DetectTypesFromTestCases() runtime detection with config-driven capabilities. Add TestCapabilities field to BuildIndexOpts, resolve capabilities from config.yaml with UnionCapabilities fallback in callers (index.go, add.go). Remove DetectTypesFromTestCases, summaryTableRow regex, and all legacy generic else-branches from GetBreakdownTestTasks, GetQuickTestTasks, resolveBreakdownDeps, resolveQuickDeps. All test generation is now deterministic and config-driven. Bump version to 3.17.0.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/testgen.go
- forge-cli/pkg/task/build_test.go
- forge-cli/pkg/task/testgen_test.go
- forge-cli/internal/cmd/index.go
- forge-cli/internal/cmd/add.go
- forge-cli/scripts/version.txt

### Key Decisions
- Empty capabilities returns nil (no test tasks) instead of falling back to legacy generic path
- BuildIndex condition changed to require len(opts.TestCapabilities) > 0 alongside profiles and mode
- Capabilities resolution: config.yaml Capabilities > UnionCapabilities(profiles) > empty
- Removed regexp and strconv imports from testgen.go after deleting DetectTypesFromTestCases

## Test Results
- **Tests Executed**: Yes
- **Passed**: 398
- **Failed**: 0
- **Coverage**: 90.0%

## Acceptance Criteria
- [x] BuildIndexOpts has TestCapabilities []string field
- [x] BuildIndex reads opts.TestCapabilities instead of calling DetectTypesFromTestCases()
- [x] Callers (index.go, add.go) resolve capabilities from config.yaml with UnionCapabilities fallback
- [x] DetectTypesFromTestCases function and summaryTableRow regex deleted entirely
- [x] detectedTypes parameter removed from all 4 functions (GetQuickTestTasks, GetBreakdownTestTasks, resolveQuickDeps, resolveBreakdownDeps)
- [x] All legacy else branches (generic task creation) removed from all 4 functions
- [x] forge task index produces identical output regardless of whether test-cases.md exists
- [x] All unit tests updated and passing
- [x] No orphaned generic tasks in any scenario

## Notes
Version bumped from 3.16.0 to 3.17.0 (minor: new config-driven behavior). Breaking change to BuildIndexOpts struct and public function signatures as indicated by task metadata.
