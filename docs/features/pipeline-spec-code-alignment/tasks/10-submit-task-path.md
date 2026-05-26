---
id: "10"
title: "Fix submit-task record format path resolution"
priority: "P1"
estimated_time: "30min"
dependencies: [9]
type: "doc"
mainSession: false
---

# 10: Fix submit-task record format path resolution

## Description

The submit-task SKILL.md references record format files using a path that cannot be resolved from the subagent's working directory. The subagent runs in the project root, but the path in the skill doc assumes a different working directory context.

Fix: Change the record format file path to an absolute path within the plugin distribution, or use a path that resolves correctly from any working directory. For example: `plugins/forge/skills/submit-task/data/record-format-{TASK_CATEGORY}.md` or the equivalent resolved path based on the forge distribution model.

Also check the record-format-doc.md content for the `doc.eval` ghost type fix (done in Task 7) and ensure path references are consistent.

## Reference Files
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Problem` — Evidence D7 (submit-task path resolution failure — Critical)
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Proposed-Solution` — Cluster 6 description
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Scope` — Cluster 6 In Scope bullet

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/submit-task/SKILL.md` | Record format path → resolvable path from any working directory |

## Acceptance Criteria
- [ ] Record format file path can be resolved from the subagent's working directory
- [ ] Path follows the forge distribution model conventions
- [ ] Both `record-format-coding.md` and `record-format-doc.md` paths resolve correctly

## Hard Rules
- Path must work from the project root (where the subagent runs)
- Follow the path resolution rules in `docs/conventions/forge-distribution.md`

## Implementation Notes
- Check `docs/conventions/forge-distribution.md` for the correct path pattern for distributed skill data files
- The subagent runs in the project root, so relative paths must be from that context
