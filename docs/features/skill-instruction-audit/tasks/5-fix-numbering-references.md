---
id: "5"
title: "Fix numbering/reference errors in skills"
priority: "P1"
estimated_time: "1h"
dependencies: [4]
type: "doc"
mainSession: false
---

# 5: Fix numbering/reference errors in skills

## Description

Fix step numbering gaps, misleading cross-references, invalid section references. "编号/引用修复" subcategory (~12 instances).

## Reference Files
- `plugins/forge/skills/tech-design/SKILL.md`: Process Flow skips Step 9-10
- `plugins/forge/skills/run-tests/SKILL.md`: Step 5 references "Convention loaded in Step 0" (wrong)
- `plugins/forge/skills/write-prd/SKILL.md`: Step 9.5 decimal numbering
- `plugins/forge/skills/quick-tasks/SKILL.md`: Step 1.5 decimal numbering
- `plugins/forge/skills/breakdown-tasks/SKILL.md`: Step 4b null-operation step
- `plugins/forge/skills/gen-contracts/SKILL.md`: Invalid Section number references

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/tech-design/SKILL.md` | Fix Process Flow: insert Step 9→10 between 8→11 |
| `plugins/forge/skills/run-tests/SKILL.md` | Fix Step 5: replace "Convention loaded in Step 0" with correct reference |
| `plugins/forge/skills/write-prd/SKILL.md` | Renumber 9.5→10, push current 10→11 |
| `plugins/forge/skills/quick-tasks/SKILL.md` | Renumber 1.5→2, push subsequent steps |
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Convert Step 4b to informational note |
| `plugins/forge/skills/gen-contracts/SKILL.md` | Remove Section X.Y references; keep concept definitions |

## Acceptance Criteria

- [ ] `tech-design` Process Flow: 0→1→...→8→9→10→11 no gaps
- [ ] `run-tests` Step 5 does not reference "Convention loaded in Step 0"
- [ ] `write-prd` has no decimal step numbers
- [ ] `quick-tasks` has no decimal step numbers
- [ ] `breakdown-tasks` Step 4b is informational note, not numbered step
- [ ] `gen-contracts` has no "Section X.Y" references

## Hard Rules

- 仅修改上述 6 个文件
- After renumbering, verify all cross-references point to correct steps
