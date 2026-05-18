---
id: "3"
title: "Inline step0-profile-resolution, type-assignment, and intent-propagation into consuming skills"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "documentation"
mainSession: false
---

# 3: Inline step0-profile-resolution, type-assignment, and intent-propagation into consuming skills

## Description
Replace 3 reference file read instructions in breakdown-tasks and quick-tasks skills with inlined content. These references define language resolution (step0-profile-resolution), type-to-assignment mapping (type-assignment), and proposal intent propagation (intent-propagation).

## Reference Files
- `docs/proposals/remove-references-dir/proposal.md` — Source proposal
- `plugins/forge/references/shared/step0-profile-resolution.md` — Language detection protocol
- `plugins/forge/references/shared/type-assignment.md` — Type mapping table
- `plugins/forge/references/shared/intent-propagation.md` — Intent propagation logic

> **Note:** Line numbers are approximate and may drift. Search for `references/shared/step0-profile-resolution`, `references/shared/type-assignment`, or `references/shared/intent-propagation` to locate exact reference sites.

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Replace lines 44, 302, 306 — inline all 3 reference contents |
| `plugins/forge/skills/quick-tasks/SKILL.md` | Replace lines 39, 85, 89 — inline all 3 reference contents |

## Acceptance Criteria
- [ ] No occurrence of `references/shared/step0-profile-resolution` in either file
- [ ] No occurrence of `references/shared/type-assignment` in either file
- [ ] No occurrence of `references/shared/intent-propagation` in either file
- [ ] Each file contains the full language detection protocol, type mapping table, and intent propagation logic

## Hard Rules
- Inline the full content of each reference file at the point where the Read instruction appears
- The step0-profile-resolution content includes a HARD-RULE about not silently defaulting language — preserve this constraint

## Implementation Notes
- breakdown-tasks and quick-tasks are structurally similar (quick-tasks derives from breakdown-tasks)
- All 3 references appear at the same locations in both files (Step 0, Step 3 "Type Assignment", Step 3 "Intent Propagation")
