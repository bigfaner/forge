---
id: "3"
title: "Update command/agent files with discovery instruction"
priority: "P0"
estimated_time: "30m"
dependencies: ["1"]
type: "documentation"
mainSession: false
---

# 3: Update command/agent files with discovery instruction

## Description

Remove the hardcoded keyword→filename mapping tables from `fix-bug.md`, replacing them with the same discovery instruction used in prompt templates.

> **Note**: This task originally also targeted `plugins/forge/agents/error-fixer.md`, but that agent was removed during the `typed-task-dispatch` feature. The discovery instruction was applied to `fix-bug.md` only.

## Reference Files
- `docs/proposals/knowledge-discovery/proposal.md` — Source proposal (discovery instruction text)

## Affected Files

### Create

| File | Description |
|------|-------------|
| (none) | |

### Modify

| File | Changes |
|------|---------|
| `plugins/forge/commands/fix-bug.md` | Remove mapping table, add discovery instruction |

### Delete

| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] `fix-bug.md` contains the discovery instruction (same as in prompt templates)
- [ ] No file contains the hardcoded mapping pattern `"auth"/"login"/"permission" → business-rules/auth.md`
- [ ] No file contains any keyword→filename mapping table
- [ ] The "Project Knowledge" section structure is preserved — only the mapping content changes

## Hard Rules

- The mapping table in `fix-bug.md` needs to be replaced with the discovery instruction
- Do not change the surrounding instructions (the "Infer relevant domains" bullet stays, only the "Example mappings" bullet is removed)

## Implementation Notes

- `fix-bug.md` has the mapping in "Project Knowledge" section
- Replace the "Example mappings: ..." line with the discovery instruction
