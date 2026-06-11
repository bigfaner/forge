---
status: "completed"
started: "2026-05-18 23:07"
completed: "2026-05-18 23:18"
time_spent: "~11m"
---

# Task Record: 2 Task templates — remove user confirmation, adapt for auto mode

## Summary
Updated consolidate-specs task template and both prompt templates (consolidate + drift) for fully automated pipeline execution. Removed the 'blocked on CROSS items' skip condition from the task template, replaced with auto-integration instruction referencing SKILL.md Step 6 non-interactive mode. Updated doc-generation-consolidate.md to instruct non-interactive mode with auto-integrate and [auto-specs] commit. Updated doc-generation-drift.md similarly for drift-only non-interactive mode. All-LOCAL items still auto-proceed without blocking. Interactive/manual invocation behavior preserved unchanged.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/doc-generation-consolidate.md
- forge-cli/pkg/prompt/data/doc-generation-drift.md
- plugins/forge/skills/breakdown-tasks/templates/consolidate-specs.md
- forge-cli/pkg/prompt/prompt_test.go

### Key Decisions
- Used 'non-interactive' keyword in prompt templates so the consolidate-specs SKILL.md Step 6 non-interactive path is triggered automatically
- Avoided the literal word 'blocked' in prompt templates to prevent confusion with task status blocked
- Task template now references SKILL.md Step 6 non-interactive mode instead of duplicating the auto-integration logic

## Test Results
- **Tests Executed**: Yes
- **Passed**: 79
- **Failed**: 0
- **Coverage**: 90.6%

## Acceptance Criteria
- [x] consolidate-specs.md task template: 'early exit or user review' step skips confirmation in pipeline mode
- [x] consolidate-specs.md task template: all-LOCAL items still auto-proceed without blocking
- [x] doc-generation-consolidate.md: prompt template instructs agent to run in non-interactive mode
- [x] doc-generation-drift.md: prompt template instructs agent to run in non-interactive mode
- [x] Full pipeline /run-tasks no longer blocks on consolidate-specs waiting for user input
- [x] Quick pipeline /run-tasks executes drift-only consolidate-specs without blocking

## Notes
无
