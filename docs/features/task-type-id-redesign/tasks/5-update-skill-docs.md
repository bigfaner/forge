---
id: "5"
title: "Update skill documentation"
priority: "P2"
estimated_time: "30min"
dependencies: ["1"]
type: "documentation"
mainSession: false
---

# 5: Update skill documentation

## Description
Update three skill documentation files to reflect the new type system with prefix-based naming and updated quality-gate protocol.

### Files to update
1. **`references/shared/type-assignment.md`**: Update type table with new prefix-based types (`coding.*`, `doc*`, `test.*`, `validation.*`), add prefix rules for quality-gate routing
2. **`hooks/guide.md`**: Update quality-gate protocol to use `IsTestableType` prefix check instead of hardcoded type list
3. **`skills/submit-task/SKILL.md`**: Update quality-gate check to use prefix-based type determination

## Reference Files
- `docs/proposals/task-type-id-redesign/proposal.md` — Source proposal
- `plugins/forge/references/shared/type-assignment.md` — Type assignment reference
- `plugins/forge/hooks/guide.md` — Quality-gate protocol
- `plugins/forge/skills/submit-task/SKILL.md` — Submit-task skill

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/references/shared/type-assignment.md` | Update type table, add prefix rules |
| `plugins/forge/hooks/guide.md` | Update quality-gate protocol |
| `plugins/forge/skills/submit-task/SKILL.md` | Update quality-gate check |

## Acceptance Criteria
- [ ] Type table shows all new prefix-based types with correct categories
- [ ] Quality-gate routing rule: `coding.*` → run gate, `doc*` → skip gate, `test.*`/`validation.*`/`gate` → special handling
- [ ] submit-task SKILL.md references new type names for quality-gate decision

## Hard Rules
- Must load `docs/conventions/forge-distribution.md` before modifying any plugin files
- Documentation must be consistent with the actual Go implementation from task 1

## Implementation Notes
- The key change in all three files is replacing the hardcoded type list with prefix-based rules
- `coding.*` = quality-gate runs; everything else = skip or special handling
