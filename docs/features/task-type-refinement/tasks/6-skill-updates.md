---
id: "6"
title: "Update skill type assignment tables for new task types"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "documentation"
mainSession: false
---

# 6: Update skill type assignment tables for new task types

## Description

Update the type assignment documentation in `/quick-tasks` and `/breakdown-tasks` skills to reflect the new type taxonomy. The skills need to know that `implementation` is deprecated and tasks should be classified as `feature`, `enhancement`, `cleanup`, or `refactor` based on the task's primary output.

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/quick-tasks/SKILL.md` | Update Type Assignment table and template selection logic |
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Update Type Assignment table and template selection logic |
| `plugins/forge/skills/quick-tasks/templates/task.md` | Change default type from `implementation` to `feature` |
| `plugins/forge/skills/breakdown-tasks/templates/task.md` | Change default type from `implementation` to `feature` |

## Acceptance Criteria
- [ ] `quick-tasks/SKILL.md` Type Assignment table lists: `feature`, `enhancement`, `cleanup`, `refactor`, `documentation`, `gate`
- [ ] `breakdown-tasks/SKILL.md` Type Assignment table lists: `feature`, `enhancement`, `cleanup`, `refactor`, `doc-generation`, `gate`, `test-pipeline`
- [ ] Type Assignment tables include clear "When to assign" guidance for each type
- [ ] Template selection logic updated: tasks with testable runtime behavior → `task.md`, doc-only → `task-doc.md` (unchanged), but default type in `task.md` frontmatter is `feature`
- [ ] Skills describe how to read proposal `intent` field and use it as default type, with per-task override

## Hard Rules
- `implementation` should NOT appear in the new type tables. Mark it as deprecated with a note that existing index.json files should be migrated.

## Implementation Notes
- The proposal D1 describes the intent propagation: proposal `intent` → default task type, individual task frontmatter `type` → override.
- Both skills need to know about the mapping: proposal intent value → task type constant (1:1 mapping since they use the same names).
