---
status: "completed"
started: "2026-05-20 13:07"
completed: "2026-05-20 13:17"
time_spent: "~10m"
---

# Task Record: 2 Add system type interception in BuildIndex and validate-index

## Summary
Add system type interception in BuildIndex and validate-index: non-auto-gen tasks using system types are now rejected with descriptive error messages. Exported IsAutoGenTaskID and added FormatSystemTypes helper. Both BuildIndex and validateTasks enforce consistent rules.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/types.go
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/build_test.go
- forge-cli/internal/cmd/validate_index.go
- forge-cli/internal/cmd/validate_index_test.go

### Key Decisions
- Exported IsAutoGenTaskID (was isAutoGenTaskID) to allow validate_index.go to reuse the same auto-gen detection logic without duplication
- Added FormatSystemTypes() helper for deterministic, sorted error messages listing all reserved types
- System type check placed after existing type validation in both BuildIndex (section 5.5.0) and validateTasks (after type switch)
- Error format: task '%s': type '%s' is a system-reserved type (reserved: %s) - consistent between both enforcement points

## Test Results
- **Tests Executed**: Yes
- **Passed**: 18
- **Failed**: 0
- **Coverage**: 89.9%

## Acceptance Criteria
- [x] BuildIndex rejects non-auto-gen tasks using system types
- [x] BuildIndex allows auto-gen tasks using system types
- [x] validate-index rejects non-auto-gen tasks using system types
- [x] Error messages include specific type and full system type list
- [x] Quality gate fix tasks (coding.fix, coding.cleanup) pass validation
- [x] go test ./forge-cli/... passes

## Notes
Also fixed TestBuildIndex_TypeInference test which used a business task ID with gate type - now uses coding.feature which is a valid business type.
