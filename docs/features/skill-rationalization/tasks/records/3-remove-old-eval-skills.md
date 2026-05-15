---
status: "blocked"
started: "2026-05-16 00:58"
completed: "N/A"
time_spent: ""
---

# Task Record: 3 Remove old eval skill directories, record-task skill, and simplify-skill command

## Summary
Deleted 7 old eval skill directories (eval-proposal, eval-prd, eval-design, eval-ui, eval-test-cases, eval-consistency, eval-harness), record-task skill directory, and simplify-skill command. Fixed all dangling references in remaining files: updated rubric report template paths, improve-harness rubric references, README.md, ARCHITECTURE.md, and Go test files.

## Changes

### Files Created
无

### Files Modified
- README.md
- docs/ARCHITECTURE.md
- plugins/forge/skills/eval/rubrics/proposal.md
- plugins/forge/skills/eval/rubrics/prd.md
- plugins/forge/skills/eval/rubrics/design.md
- plugins/forge/skills/eval/rubrics/ui-web.md
- plugins/forge/skills/eval/rubrics/ui-mobile.md
- plugins/forge/skills/eval/rubrics/ui-tui.md
- plugins/forge/skills/eval/rubrics/test-cases.md
- plugins/forge/skills/eval/rubrics/consistency.md
- plugins/forge/skills/eval/rubrics/harness.md
- plugins/forge/skills/improve-harness/SKILL.md
- plugins/forge/skills/improve-harness/templates/improvements.md
- tests/e2e/tui-ui-design/tui_ui_design_cli_test.go
- forge-cli/tests/e2e/justfile_mixed_cli_cli_test.go
- forge-cli/tests/e2e/plugin_content_cli_test.go

### Key Decisions
- Removed 'Report template' lines from rubric files since the generic eval skill defines report format inline in SKILL.md Step 5
- Updated improve-harness references from old eval-harness/templates/rubric.md to new eval/rubrics/harness.md
- Updated Go test files to point to new eval/rubrics/ paths and submit-task instead of record-task
- Left historical docs (self-evolution reports, proposals, lessons) untouched as they are archival records
- Left forensic test data (JSONL mock files) untouched as they are test fixtures, not live references

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 4
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] All 7 eval skill directories fully removed
- [x] plugins/forge/skills/record-task/ no longer exists
- [x] plugins/forge/commands/simplify-skill.md no longer exists
- [x] No dangling references to deleted skills/commands in remaining files

## Notes
4 pre-existing test failures in forge-cli/pkg/task/testgen_test.go (TestGetQuickTestTasks_PerType_*) are unrelated to this task -- verified by running tests before and after changes with identical failures. Task type is documentation with no dedicated test suite.
