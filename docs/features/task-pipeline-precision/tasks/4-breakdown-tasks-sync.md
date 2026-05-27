---
id: "4"
title: "Sync breakdown-tasks with split rules, complexity判定, Reference Files inline"
priority: "P0"
estimated_time: "1h"
dependencies: [3]
type: "doc"
mainSession: false
---

# 4: Sync breakdown-tasks with split rules, complexity判定, Reference Files inline

## Description

Apply the same task generation improvements from Task 3 to breakdown-tasks SKILL.md. Ensure the two skill files use consistent rules despite different source documents (proposal vs tech-design).

## Reference Files
- `docs/proposals/task-pipeline-precision/proposal.md#Scope` — breakdown-tasks specific paragraphs to modify
- `docs/proposals/task-pipeline-precision/proposal.md#Key-Risks` — risk of breakdown-tasks sync inconsistency and mitigation

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Task Splitting Rules, complexity判定, Reference Files generation |
| `plugins/forge/skills/breakdown-tasks/templates/task.md` | Add complexity field to frontmatter |

## Acceptance Criteria

- [ ] Task Splitting Rules paragraph uses "independently verifiable" as the merge standard (same as quick-tasks)
- [ ] AC max 6 rule added (same wording as quick-tasks)
- [ ] Multi-verb detection rule added (same wording as quick-tasks)
- [ ] Complexity判定 logic with LLM override matches quick-tasks wording
- [ ] Reference Files generation changed to inline precise info format (same as quick-tasks)
- [ ] `templates/task.md` frontmatter has `complexity: "{{COMPLEXITY}}"` field with default "medium"
- [ ] breakdown-tasks specific features (Phase & Gate Detection, PRD Coverage Verification) are NOT affected by these changes

## Hard Rules
{{HARD_RULES}}

## Implementation Notes

- breakdown-tasks reads from tech-design (not proposal), so Reference Files inline content will differ in practice, but the generation rule format should be identical
- The complexity判定 threshold may need adjustment note: "breakdown-tasks tasks may have naturally finer AC granularity due to tech-design decomposition; the same thresholds apply but LLM override is expected to be more common"
