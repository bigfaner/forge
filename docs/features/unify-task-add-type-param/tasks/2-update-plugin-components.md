---
id: "2"
title: "Update plugin components to use --type coding.fix"
priority: "P1"
estimated_time: "30min"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 2: Update plugin components to use --type coding.fix

## Description
Update all 6 Forge plugin components (1 agent, 2 commands, 3 skills) that reference `--template fix-task` to use the new `--type coding.fix` syntax. This is a mechanical search-and-replace across markdown files.

## Reference Files
- `proposal.md#Impact-Analysis` — lists all 6 affected plugin components with file paths and line numbers
- `proposal.md#Proposed-Solution` — defines the new unified --type syntax

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/agents/task-executor.md` | Line 106: `--template fix-task` → `--type coding.fix` |
| `plugins/forge/commands/execute-task.md` | Lines 70, 115: `--template fix-task` → `--type coding.fix` |
| `plugins/forge/commands/run-tasks.md` | Lines 64, 68, 88, 113: `--template fix-task` → `--type coding.fix` |
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Line 166: `--template fix-task` → `--type coding.fix` |
| `plugins/forge/skills/quick-tasks/SKILL.md` | Line 171: `--template fix-task` → `--type coding.fix` |
| `plugins/forge/skills/submit-task/SKILL.md` | Line 102: `--template fix-task` → `--type coding.fix` |

## Acceptance Criteria
- [ ] No plugin markdown file contains `--template fix-task` or `--template cleanup-task`
- [ ] All 6 files use `--type coding.fix` (or `--type coding.cleanup` where applicable)
- [ ] No changes to template variable flags (`--var`, `SOURCE_FILES`, `TEST_SCRIPT`, `TEST_RESULTS`)

## Hard Rules
- Only replace `--template fix-task` → `--type coding.fix` and `--template cleanup-task` → `--type coding.cleanup`
- Do NOT modify `--source-task-id`, `--block-source`, `--description`, `--var` flags — these remain unchanged

## Implementation Notes
- Use project-wide search for `--template fix-task` across `plugins/forge/` to catch any additional occurrences not listed in the proposal
- The `quick-tasks` skill template (Step 4) also references `--template fix-task` in its instruction text — update this too
