---
id: "5"
title: "Deduplicate breakdown-tasks and quick-tasks shared logic"
priority: "P2"
estimated_time: "3h"
dependencies: [3]
scope: "all"
breaking: false
type: "refactor"
mainSession: false
---

# 5: Deduplicate breakdown-tasks and quick-tasks shared logic

## Description

`breakdown-tasks` and `quick-tasks` share identical Type Assignment tables, Intent Propagation sections, Step 0 profile resolution, and 6+ near-duplicate template files. Extract shared logic into reference files that both skills include, eliminating manual sync burden.

## Reference Files
- `docs/proposals/skill-ecosystem-audit/proposal.md` — P1 finding #5

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/references/shared/type-assignment.md` | Shared Type Assignment table |
| `plugins/forge/references/shared/intent-propagation.md` | Shared Intent Propagation rules |
| `plugins/forge/references/shared/profile-resolution.md` | Shared Step 0 profile resolution instructions |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Replace inline Type Assignment (lines 310-319) and Intent Propagation (lines 322-330) with references to shared files |
| `plugins/forge/skills/quick-tasks/SKILL.md` | Replace inline Type Assignment (lines 96-103) and Intent Propagation (lines 109-116) with references to shared files |
| `plugins/forge/skills/breakdown-tasks/templates/task.md` | If content differs from quick-tasks version, keep separate |
| `plugins/forge/skills/quick-tasks/templates/task.md` | If content differs, keep separate |

## Acceptance Criteria
- `diff <(cat references/shared/type-assignment.md) <(grep -A20 'Type Assignment' breakdown-tasks/SKILL.md)` shows the table content is now sourced from the shared file
- Both breakdown-tasks and quick-tasks reference the same shared files
- Each shared reference file lists its consuming skills in a header comment
- Task generation output is identical for both `/breakdown-tasks` and `/quick-tasks` pipelines (verify with a test PRD)

## Hard Rules
- Shared reference files must be discoverable from the skill's installed location (relative path convention from Task 3)
- Do NOT merge templates that have different content (breakdown-tasks uses `<phase>.<sub>` IDs, quick-tasks uses simple integers) — only extract truly identical content
- Keep each shared file self-contained and atomic

## Implementation Notes
- The 6 near-duplicate template pairs (task.md, task-doc.md, validate-code-task.md, validate-ux-task.md, index.json, index.schema.json) — verify content is identical before merging. breakdown-tasks task.md uses phase prefixes; quick-tasks does not. These should stay separate.
- Focus on extracting the 3 shared SKILL.md sections (Type Assignment, Intent Propagation, Step 0) which are verbatim copies.
- The shared files should use the same relative-from-self resolution pattern established in Task 3.
