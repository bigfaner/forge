---
id: "5"
title: "Extract shared logic from breakdown-tasks ↔ quick-tasks"
priority: "P2"
estimated_time: "2h"
dependencies: []
type: "refactor"
scope: "all"
breaking: false
mainSession: false
---

# 5: Extract shared logic from breakdown-tasks ↔ quick-tasks

## Description

`breakdown-tasks` and `quick-tasks` share identical Type Assignment tables, Intent Propagation logic, and Step 0 profile resolution procedures. These are duplicated inline in both SKILL.md files. Extract them to shared reference files so changes only need to be made once.

## Reference Files
- `docs/proposals/skill-ecosystem-audit/proposal.md` — Source proposal (W4, item 10)
- `plugins/forge/skills/breakdown-tasks/SKILL.md` — Source of shared content
- `plugins/forge/skills/quick-tasks/SKILL.md` — Source of shared content

## Affected Files

### Create

| File | Description |
|------|-------------|
| `plugins/forge/references/shared/type-assignment.md` | Type Assignment table extracted from both skills |
| `plugins/forge/references/shared/intent-propagation.md` | Intent Propagation logic extracted from both skills |
| `plugins/forge/references/shared/step0-profile-resolution.md` | Step 0 profile resolution extracted from both skills |

### Modify

| File | Changes |
|------|---------|
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Replace inline Type Assignment table, Intent Propagation, and Step 0 sections with references to shared files |
| `plugins/forge/skills/quick-tasks/SKILL.md` | Replace inline duplicated sections with references to shared files |

## Acceptance Criteria

- [ ] Three shared reference files exist in `plugins/forge/references/shared/`
  `ls plugins/forge/references/shared/type-assignment.md plugins/forge/references/shared/intent-propagation.md plugins/forge/references/shared/step0-profile-resolution.md` succeeds
- [ ] Both breakdown-tasks and quick-tasks reference the shared files
  `grep -rn 'type-assignment.md\|intent-propagation.md\|step0-profile-resolution.md' plugins/forge/skills/breakdown-tasks/SKILL.md plugins/forge/skills/quick-tasks/SKILL.md` returns hits for all 3 files
- [ ] No content is lost — the extracted content matches the original inline versions

## Hard Rules

- Use `${CLAUDE_SKILL_DIR}/../../references/shared/<file>.md` for references from SKILL.md files.
- Each shared file must be self-contained — include any context needed to understand the table/logic.
- Do not modify the logic or content during extraction — verbatim extraction only.

## Implementation Notes

- The Type Assignment table maps task output types to `type` field values.
- Intent Propagation describes how proposal `intent` frontmatter maps to default task types.
- Step 0 profile resolution describes how to detect the project's test language(s).
- Both skills should load these files when they need them, not inline the content.
