---
id: "3"
title: "Add Task Types reference section to README.md"
priority: "P1"
estimated_time: "1.5h"
dependencies: [1, 2]
type: "documentation"
mainSession: false
---

# 3: Add Task Types reference section to README.md

## Description
Add a "Task Types & Pipeline 参考" section to `README.md` that captures the CLI task details removed from quick-tasks and breakdown-tasks SKILL.md files. This serves as developer reference and can be loaded by task-executor agents when needed.

## Reference Files
- `docs/proposals/simplify-skill-task-docs/proposal.md` — Source proposal
- `README.md` — Target file to add section
- `plugins/forge/skills/quick-tasks/SKILL.md` — Verify what was removed (task 1 output)
- `plugins/forge/skills/breakdown-tasks/SKILL.md` — Verify what was removed (task 2 output)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `README.md` | Add "Task Types & Pipeline 参考" section before "文档索引" |

## Acceptance Criteria
- [ ] New section "Task Types & Pipeline 参考" added before "文档索引"
- [ ] Contains 13 task types table (type + who generates + purpose)
- [ ] Contains Quick pipeline responsibility chain (T-quick-1~5 with brief descriptions)
- [ ] Contains Full pipeline responsibility chain (T-test-1~5 with brief descriptions)
- [ ] Contains fix-task command template with `--block-source` explanation
- [ ] Contains profile-suffix convention (single vs multiple profiles)
- [ ] Contains gate/summary auto-generation rules (phases with >=2 business tasks)
- [ ] No content lost from the SKILL.md files — everything removed in tasks 1 and 2 is present here
- [ ] Section is compact (tables, not prose) and scannable

## Implementation Notes
- Read the current SKILL.md files (after tasks 1 and 2 complete) to verify all removed content is captured.
- Use tables for compactness. Avoid verbose prose.
- Add a note: "以下内容由 forge task index 自动生成，以 CLI 行为为准"
