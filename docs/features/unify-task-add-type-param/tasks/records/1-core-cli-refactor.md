---
status: "completed"
started: "2026-05-26 22:43"
completed: "2026-05-26 22:59"
time_spent: "~16m"
---

# Task Record: 1 Core CLI refactor: remove --template, add --type auto-discovery

## Summary
Removed --template flag from forge task add, renamed template files to match type values (fix-task.md -> coding.fix.md, cleanup-task.md -> coding.cleanup.md), added --type auto-discovery for template loading, updated quality_gate.go to use type values instead of template names

## Changes

### Files Created
- forge-cli/pkg/template/data/coding.fix.md
- forge-cli/pkg/template/data/coding.cleanup.md

### Files Modified
- forge-cli/pkg/template/template.go
- forge-cli/pkg/template/template_test.go
- forge-cli/internal/cmd/task/add.go
- forge-cli/internal/cmd/task/add_cmd_test.go
- forge-cli/internal/cmd/task/testbridge.go
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/internal/cmd/quality_gate_test.go
- forge-cli/internal/cmd/characterization_test.go
- forge-cli/pkg/task/add.go
- forge-cli/pkg/task/add_test.go

### Key Decisions
- Used nil-safe changed() closure for cmd.Flags().Changed() calls to support direct executeAdd() calls with nil cmd (test exports)
- Removed ExportAddTemplate from testbridge.go since --template flag no longer exists
- Provided template vars in characterization_test.go since --type coding.fix now auto-discovers the template which requires SOURCE_FILES/TEST_SCRIPT/TEST_RESULTS

## Test Results
- **Tests Executed**: Yes
- **Passed**: 31
- **Failed**: 0
- **Coverage**: 87.6%

## Acceptance Criteria
- [x] --type coding.fix loads coding.fix.md template and applies defaults (P0, breaking, 30min, fix prefix)
- [x] --type coding.feature works without template (no matching file, no error)
- [x] --template fix-task returns error (flag removed)
- [x] forge task add -h shows no --template flag
- [x] Quality gate auto-created fix tasks use type value coding.fix instead of template name fix-task
- [x] All existing tests pass after rename (go test ./...)

## Notes
All 31 packages pass with race detection. The handleGateFailure manual instruction string changed from '--template fix-task' to '--type coding.fix'. Old template files (fix-task.md, cleanup-task.md) deleted.
