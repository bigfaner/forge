---
id: "2"
title: "Rewrite SKILL.md: Scope Inference, split rules, cleanup"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
type: "documentation"
mainSession: false
---

# 2: Rewrite SKILL.md: Scope Inference, split rules, cleanup

## Description
Rewrite quick-tasks/SKILL.md to match the proposal's information density. Three changes:

1. **Step 2 rewrite**: Replace Scope Assignment algorithm (file-path-based) with Scope Inference (semantic). Add "split by functional steps" rule. Add Reference Files directory-level fill rule. Remove "Determine affected file paths" instruction.

2. **Step 3 cleanup**: Remove "Fill Affected Files from the solution description" instruction (L84), since implementation tasks no longer have an Affected Files section.

3. **Output Checklist cleanup**: Remove the "Each task file includes `## Affected Files` section" check item (L179).

## Reference Files
- `docs/proposals/optimize-quick-tasks-split/proposal.md` — Source proposal
- `plugins/forge/skills/quick-tasks/SKILL.md` — Skill file to modify
- `plugins/forge/skills/quick-tasks/templates/task.md` — Updated template (Task 1 output)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/quick-tasks/SKILL.md` | Rewrite Step 2 (L48-69), clean Step 3 (L84), clean Output Checklist (L179) |

## Acceptance Criteria
- [ ] Step 2 no longer contains "Determine affected file paths" instruction or Scope Assignment algorithm
- [ ] Step 2 contains Scope Inference rules (frontend/backend/all from task description semantics)
- [ ] Step 2 contains split-by-functional-steps rule (multi-step In Scope items → multiple tasks, total ≤ 10)
- [ ] Step 2 contains Reference Files fill rule (directory-level paths from proposal In Scope)
- [ ] Step 3 no longer references "Fill Affected Files"
- [ ] Output Checklist no longer has Affected Files check item
- [ ] `forge task index` still generates valid index.json after changes
- [ ] executor agents can work from Description + Reference Files without needing file-level paths

## Hard Rules
- Do NOT change Step 0, Step 1, Step 4, Step 5, Step 6, Step 7, or the Integration section
- 10 task hard limit must remain enforced

## Implementation Notes
- Scope Inference rules from proposal:
  - Description mentions UI/pages/components/styles → scope: "frontend"
  - Description mentions API/server/database/CLI → scope: "backend"
  - Mixed or unclear → scope: "all"
- Reference Files example: proposal says "Add --type argument to gen-test-scripts" → fill `plugins/forge/skills/gen-test-scripts/`
- The split rule: if one In Scope bullet contains multiple independently verifiable functional steps, split into separate tasks. Each task is an independently verifiable functional unit.
- Scope field still needed in task frontmatter (from task.md template) but now inferred semantically
