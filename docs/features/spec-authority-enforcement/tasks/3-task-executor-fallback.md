---
id: "3"
title: "Add Reference Files fallback rule to task-executor.md Hard Constraints"
priority: "P0"
estimated_time: "20m"
dependencies: []
type: "doc"
mainSession: false
---

# 3: Add Reference Files fallback rule to task-executor.md Hard Constraints

## Description

Add a new Hard Constraint rule to `plugins/forge/agents/task-executor.md` that serves as a fallback: even if the synthesized strategy template (coding.* template) doesn't include a Reference Files declaration, the agent must still proactively read the task file's `## Reference Files` and treat them as authoritative.

## Reference Files
- `docs/proposals/spec-authority-enforcement/proposal.md#Proposed-Solution` — Agent layer description and fallback rule purpose
- `docs/proposals/spec-authority-enforcement/proposal.md#Priority-Rules` — Priority: Hard Rules > Reference Files > existing code

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/agents/task-executor.md` | Add new rule 8 in `<EXTREMELY-IMPORTANT>` Hard Constraints block |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] A new rule is added as item 8 in the `<EXTREMELY-IMPORTANT>` Hard Constraints block (after existing rule 7)
- [ ] The rule states: agent must proactively read `## Reference Files` from the task file and treat them as authoritative, even if the synthesized strategy doesn't include a Reference Files declaration
- [ ] The rule references the priority order: task `## Hard Rules` > `## Reference Files` > existing code structure
- [ ] The rule does not duplicate the template-layer `<IMPORTANT>` block — it's a fallback, not a repeat
- [ ] The modification takes effect immediately (no compilation needed — task-executor.md is distributed as-is)

## Hard Rules
- MUST load `docs/conventions/forge-distribution.md` before modifying task-executor.md
- MUST add as a NEW numbered rule in the existing `<EXTREMELY-IMPORTANT>` block — do NOT create a second block
- MUST NOT modify existing rules 1-7

## Implementation Notes

The new rule should be concise — it's a fallback, not a full declaration. Suggested wording:

```
8. SPEC AUTHORITY FALLBACK — if the synthesized strategy does not include a Reference Files declaration, you MUST still:
   - Read the task file's `## Reference Files` section
   - Treat listed documents as authoritative sources (priority: `## Hard Rules` > `## Reference Files` > existing code)
   - Output a confirmation: "Fallback: Loaded Reference Files from task file: [list]"
```

This is a safety net for tasks using templates that haven't been updated with the `<IMPORTANT>` Reference Files block.
