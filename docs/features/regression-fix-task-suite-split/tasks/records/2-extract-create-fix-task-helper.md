---
status: "completed"
started: "2026-05-29 14:32"
completed: "2026-05-29 14:42"
time_spent: "~10m"
---

# Task Record: 2 提取 createFixTask 共享 helper 函数

## Summary
Extracted createFixTask shared helper function from addSingleFixTask, encapsulating surface inference, opts construction, AddTask, CreateTaskMarkdown, and EnsureForgeState. addSingleFixTask refactored to cap check + delegation. Added fixTaskOverride struct for title/description/extraVars customization needed by Task 4.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/internal/cmd/quality_gate_test.go

### Key Decisions
- Used variadic fixTaskOverride parameter for optional overrides (Title, Description, ExtraVars) to support both addSingleFixTask defaults and Task 4's custom per-file tasks
- Cap check stays in addSingleFixTask, not in createFixTask, so Task 4's addRegressionFixTasks can bypass the per-step cap with its own soft cap

## Test Results
- **Tests Executed**: Yes
- **Passed**: 325
- **Failed**: 0
- **Coverage**: 68.3%

## Acceptance Criteria
- [x] New createFixTask function encapsulates task creation core logic (surface inference, opts construction, AddTask, CreateTaskMarkdown, EnsureForgeState) with parameters title, sourceFiles, output, errorDocPath, step
- [x] addSingleFixTask refactored to: cap check (countFixTasks + maxFixTasksPerStep) + calling createFixTask
- [x] addSingleFixTask original behavior unchanged (compile/fmt/lint/unit-test path fix task creation results consistent with pre-refactor)
- [x] createFixTask has independent unit test coverage for task field population, markdown generation, state update

## Notes
8 new tests added for createFixTask: DefaultFields, MarkdownGenerated, ForgeStateUpdated, TitleOverride, DescriptionOverride, ExtraVarsOverride, SurfaceInferenceInHelper, plus AddSingleFixTask_DelegatesToCreateFixTask. All 325 tests pass with 0 failures.
