---
status: "completed"
started: "2026-05-27 00:37"
completed: "2026-05-27 00:59"
time_spent: "~22m"
---

# Task Record: 2 Fix quality_gate.go and template.go for cleanup-task and fix-task

## Summary
Fixed quality_gate.go and template.go: inferSurface uses all source files, fix-tasks grouped by test suite (directory), cleanup-tasks use Breaking=false with EstimatedTime=15min from template defaults (dual-source truth), handleGateFailure reports actual breaking status

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/pkg/template/template.go
- forge-cli/pkg/template/data/coding.cleanup.md
- forge-cli/pkg/template/data/coding.fix.md
- forge-cli/pkg/task/add.go
- forge-cli/internal/cmd/quality_gate_test.go
- forge-cli/pkg/template/template_test.go
- forge-cli/docs/WORKFLOW.md

### Key Decisions
- Dual-source truth: template.go Defaults is the authoritative source for Breaking/EstimatedTime; addFixTask reads from tmpl.GetDefaults() instead of hardcoding
- inferSurface iterates all comma-separated source files instead of just the first, falling back through each until a surface match is found
- Fix-tasks grouped by directory (test suite): groupFilesByDir splits source files by directory, creating one fix-task per directory group for parallel execution
- handleGateFailure now accepts a breaking bool parameter so the hook JSON message accurately reflects whether the fix-task blocks downstream
- Both coding.fix and coding.cleanup templates now use {{ESTIMATED_TIME}} placeholder, resolved via ApplyVars builtin

## Test Results
- **Tests Executed**: Yes
- **Passed**: 29
- **Failed**: 0
- **Coverage**: 87.6%

## Acceptance Criteria
- [x] inferSurface uses all source files, not just the first
- [x] Fix-tasks are grouped by test suite (same directory), not problem type
- [x] EstimatedTime comes from Go opts as authoritative source
- [x] Cleanup-task uses Breaking: false and EstimatedTime: 15min
- [x] coding.cleanup.md template frontmatter has breaking: false
- [x] Existing tests pass (go test ./...)

## Notes
All 6 acceptance criteria met. Fix-task by-suite grouping implemented via groupFilesByDir: when source files span multiple directories, one fix-task is created per directory. Same-directory files stay in one task (bottom-line rule from spec).
