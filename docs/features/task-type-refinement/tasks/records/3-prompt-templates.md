---
status: "completed"
started: "2026-05-16 20:56"
completed: "2026-05-16 21:04"
time_spent: "~8m"
---

# Task Record: 3 Add prompt templates for feature, enhancement, cleanup, and refactor types

## Summary
Created four new prompt templates (feature.md, enhancement.md, cleanup.md, refactor.md) replacing implementation.md with type-specific execution strategies. feature/enhancement follow TDD, cleanup follows improve-then-verify, refactor follows behavior-preservation-verify. All templates support standard variables and are self-contained.

## Changes

### Files Created
- forge-cli/pkg/prompt/data/feature.md
- forge-cli/pkg/prompt/data/enhancement.md
- forge-cli/pkg/prompt/data/cleanup.md
- forge-cli/pkg/prompt/data/refactor.md

### Files Modified
- forge-cli/pkg/prompt/data/implementation.md
- forge-cli/pkg/prompt/prompt.go
- forge-cli/pkg/prompt/prompt_test.go

### Key Decisions
- feature.md and enhancement.md use TDD workflow (RED/GREEN/REFACTOR) matching their testable runtime behavior
- cleanup.md skips TDD - workflow is read task, make improvements, run quality gate
- refactor.md skips TDD - workflow is read task, restructure code preserving behavior, run quality gate
- implementation.md kept with Deprecated header comment for backward compatibility

## Test Results
- **Tests Executed**: Yes
- **Passed**: 21
- **Failed**: 0
- **Coverage**: 90.6%

## Acceptance Criteria
- [x] data/feature.md created: implement functionality to quality gate
- [x] data/enhancement.md created: enhance existing behavior to quality gate
- [x] data/cleanup.md created: improve technical debt to quality gate (no TDD requirement)
- [x] data/refactor.md created: restructure code to quality gate (behavior preservation check)
- [x] typeToTemplate map updated with 4 new entries mapping to new templates
- [x] implementation.md kept but marked deprecated (header comment)
- [x] All templates support standard variables: TASK_ID, TASK_FILE, SCOPE, FEATURE_SLUG

## Notes
无
