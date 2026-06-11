---
status: "completed"
started: "2026-05-15 22:46"
completed: "2026-05-15 23:06"
time_spent: "~20m"
---

# Task Record: T-quick-5 Verify Quick E2E Regression

## Summary
Fix e2e regression failures: add feature-directory validation in e2e.Run(), create record-task skill alias, add type-assignment sections to breakdown-tasks and quick-tasks SKILL.md, add Element field and sitemap-missing fallback to gen-test-cases SKILL.md, update task-executor.md to reference record-task, fix unit test fixture for feature existence check.

## Changes

### Files Created
- plugins/forge/skills/record-task/SKILL.md

### Files Modified
- forge-cli/pkg/e2e/actions.go
- forge-cli/pkg/e2e/actions_test.go
- plugins/forge/agents/task-executor.md
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/quick-tasks/SKILL.md
- plugins/forge/skills/gen-test-cases/SKILL.md

### Key Decisions
- Added feature directory existence check in e2e.Run() before calling just test-e2e, matching the pattern already used in Verify()
- Created record-task as a thin alias skill that delegates to submit-task, since record-task was renamed to submit-task but tests still reference the old name
- Added type-assignment sections to SKILL.md files to document the type taxonomy (implementation, doc-generation, gate, test-pipeline) expected by e2e tests

## Test Results
- **Tests Executed**: Yes
- **Passed**: 804
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] All e2e regression tests pass
- [x] All unit tests still pass after changes

## Notes
Fixed 3 categories of regression failures: (1) e2e.Run() missing feature-not-found validation, (2) missing record-task skill (renamed to submit-task), (3) missing type-assignment/Element sections in SKILL.md files.
