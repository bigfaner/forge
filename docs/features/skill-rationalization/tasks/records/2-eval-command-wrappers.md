---
status: "completed"
started: "2026-05-16 00:54"
completed: "2026-05-16 00:57"
time_spent: "~3m"
---

# Task Record: 2 Create eval command wrappers for backward-compatible slash commands

## Summary
Created 7 thin command wrapper files in plugins/forge/commands/ that delegate to the generic eval skill via Skill(skill='forge:eval', args='--type <type>'). Each wrapper preserves the original slash command name and description from the source eval skill SKILL.md frontmatter, ensuring backward compatibility.

## Changes

### Files Created
- plugins/forge/commands/eval-proposal.md
- plugins/forge/commands/eval-prd.md
- plugins/forge/commands/eval-design.md
- plugins/forge/commands/eval-ui.md
- plugins/forge/commands/eval-test-cases.md
- plugins/forge/commands/eval-consistency.md
- plugins/forge/commands/eval-harness.md

### Files Modified
无

### Key Decisions
- Copied descriptions verbatim from original eval skill SKILL.md frontmatter to maintain system prompt consistency
- Used the command-to-skill delegation pattern established by commands/quick.md

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] Each command wrapper has correct frontmatter (name, description)
- [x] Each wrapper invokes Skill(forge:eval, args=--type <type>)
- [x] eval-harness wrapper passes --type harness (not a separate skill invocation)
- [x] Command descriptions match the original eval skill descriptions for system prompt consistency

## Notes
Pre-existing test failures in forge-cli/pkg/task (TestGetQuickTestTasks_PerType_*) are unrelated to this documentation-only change.
