---
status: "completed"
started: "2026-05-23 11:17"
completed: "2026-05-23 11:20"
time_spent: "~3m"
---

# Task Record: 2 Prompt 模板错误处理与行为对齐

## Summary
Align 8 prompt templates' error handling semantics with task-executor agent's Execution Protocol: fmt failures changed from Stop to WARNING (non-blocking), lint/stop failures now reference Complex Error Pause Flow, coding-refactor Pre-check failure sets blocked status, and coding-cleanup/coding-refactor coverage directives aligned with 'no new tests' policy.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/coding-cleanup.md
- forge-cli/pkg/prompt/data/coding-refactor.md
- forge-cli/pkg/prompt/data/coding-enhancement.md
- forge-cli/pkg/prompt/data/coding-feature.md
- forge-cli/pkg/prompt/data/coding-fix.md
- forge-cli/pkg/prompt/data/validation-code.md
- forge-cli/pkg/prompt/data/gate.md

### Key Decisions
- fmt WARNING pattern: check if affected files are ones you modified; if yes fix, if pre-existing then continue and log warning
- Complex Error Pause Flow reference added to lint failure actions: ~3 attempts threshold before creating fix task
- coding-refactor Pre-check now explicitly sets blocked status via 'forge task status' instead of vague 'stop and report'
- Coverage strategy text for cleanup/refactor: 'maintain existing coverage, no new tests required' as lead sentence, with injected placeholders as secondary context
- validation-ux.md unchanged: already had correct blocked behavior and no fmt step
- gate.md mermaid diagram updated to reflect new fmt WARNING branching and lint Complex Error Pause Flow

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] 8 templates: fmt failure changed from Stop to WARNING (non-blocking)
- [x] All 'stop' directives include 'evaluate Complex Error Pause Flow' semantics
- [x] coding-refactor Pre-check failure sets blocked status
- [x] coding-cleanup and coding-refactor coverage aligned with 'no new tests'
- [x] All-Completed Hook quality gate still runs fmt as safety net

## Notes
validation-ux.md was listed in Affected Files but requires no changes: it has no fmt/lint step and already uses blocked status correctly. Only 7 of 8 files were modified.
