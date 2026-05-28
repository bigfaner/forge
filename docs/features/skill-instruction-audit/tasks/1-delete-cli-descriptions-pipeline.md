---
id: "1"
title: "Delete CLI behavior descriptions from pipeline skills/commands"
priority: "P0"
estimated_time: "2h"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Delete CLI behavior descriptions from pipeline skills/commands

## Description

Remove all CLI behavior descriptions from the core pipeline skill/command files. These files explain what forge CLI commands do internally instead of just instructing the agent to run the command.

Apply the three-category boundary rule from proposal:
- **Keep**: imperative instructions, output contracts (field names + missing/anomaly meanings), exit code contracts
- **Delete**: behavioral explanations (internal branch logic, default value derivation, implementation details)

## Reference Files
- `docs/proposals/skill-instruction-audit/proposal.md#CLI-描述删除边界规则`: Three-category boundary rule table with examples
- `plugins/forge/skills/submit-task/SKILL.md`: "## What forge task submit Does" section (source: proposal.md#Evidence)
- `plugins/forge/skills/breakdown-tasks/SKILL.md`: "Auto-generated tasks by forge task index" block
- `plugins/forge/skills/quick-tasks/SKILL.md`: "Auto-generated tasks by forge task index" block
- `plugins/forge/commands/execute-task.md`: Output parsing field semantic explanations, "forge task add deduplicates" behavioral explanation

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/submit-task/SKILL.md` | Delete "## What forge task submit Does" section; keep command + exit code contract |
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Delete "Auto-generated tasks by forge task index" explanatory block; keep Step 5 command |
| `plugins/forge/skills/quick-tasks/SKILL.md` | Delete "Auto-generated tasks by forge task index" explanatory block; keep Step 5 command |
| `plugins/forge/commands/execute-task.md` | Delete field semantic explanations in output parsing; keep field name list; delete "forge task add deduplicates" explanations |
| `plugins/forge/commands/run-tasks.md` | Same pattern as execute-task; delete "Subagent calls forge prompt get-by-task-id internally" |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] `submit-task/SKILL.md` has no "What .* Does" section; `forge task submit` command and exit code contract remain
- [ ] `breakdown-tasks/SKILL.md` and `quick-tasks/SKILL.md` have no "Auto-generated tasks by forge task index" block; Step 5 command remains
- [ ] `execute-task.md` retains field name list (TASK_ID, FILE, etc.) but removes per-field semantic explanations and example values
- [ ] `run-tasks.md` retains field name list; no "Subagent calls forge prompt get-by-task-id internally"
- [ ] No output contract field names lost (SURFACE_KEY, SURFACE_TYPE, TASK_ID, FILE, MAIN_SESSION remain)

## Hard Rules

- 仅修改上述 5 个文件
- 保留 exit code 契约和输出字段名列表，只删除行为解释

## Implementation Notes

- **submit-task**: Remove entire "## What forge task submit Does" section. Keep command invocation and STATUS check.
- **breakdown-tasks/quick-tasks**: Remove "Auto-generated tasks by forge task index" block (~4 lines each). Keep Step 5 bash command.
- **execute-task**: Remove "(e.g., ...)" example annotations and "defaults to ..." from output parsing. Remove "forge task add automatically deduplicates — check output" explanations (2 occurrences).
- **run-tasks**: Same as execute-task. Remove "Subagent calls forge prompt get-by-task-id internally" (2 occurrences). Remove "Subagent detects 'Fix record for' prefix and calls forge prompt get-by-task-id --fix-record-missed internally".
