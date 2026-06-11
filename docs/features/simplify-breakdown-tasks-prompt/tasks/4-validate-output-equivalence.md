---
id: "4"
title: "Validate output equivalence"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["3"]
type: "doc"
mainSession: false
---

# 4: Validate output equivalence

## Description

Validate that the refactored skill produces functionally identical task output compared to the original. This covers the success criteria from the proposal: same task count, same dependency graph, same type/scope assignments, same PRD coverage.

## Reference Files
- `docs/proposals/simplify-breakdown-tasks-prompt/proposal.md` — Source proposal (Success Criteria, Validation Protocol)
- `plugins/forge/skills/breakdown-tasks/SKILL.md` — Refactored skeleton (from task 3)
- `plugins/forge/skills/breakdown-tasks/rules/` — Rule files (from tasks 1, 2)

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | Validation-only task — no files created |

### Modify
| File | Changes |
|------|---------|
| (none) | |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] **Test case 1 (backend+phases+DB)**: Run the refactored skill against a backend feature with phases and DB schema. Output passes `forge task validate-index`. Task count, dependency graph, types, and scopes match baseline (original SKILL.md output) with +/-1 task count tolerance.
- [ ] **Test case 2 (full-stack, if available)**: Run the refactored skill against a feature loading all 4 rule files (UI + phases + DB + existing code). Output passes `forge task validate-index`. If no full-stack baseline exists, create one with the original SKILL.md first, then compare.
- [ ] **Structural equivalence verified**: diff the two `index.json` files — task IDs, dependency edges, type assignments, and scope assignments must match. Wording differences are acceptable.
- [ ] **PRD coverage maintained**: every user story from the test case PRD is mapped to at least one task.
- [ ] **Error handling**: temporarily rename (remove) one rule file (e.g., `rules/db-schema.md` → `rules/db-schema.md.bak`), run the skill against a feature that would normally load that file. Output still passes `forge task validate-index` (simpler but structurally valid). Restore the file after test.
- [ ] **Execution stability**: run the same backend+phases+DB input 3 times. All 3 runs produce structurally equivalent output (same task count, same types per task).
- [ ] **Token savings confirmed**: skeleton file size is ≤8KB. Document the actual size vs the original 23KB.

## Hard Rules

- Do NOT modify any skill or rule files during validation — only rename temporarily for error handling test
- If any test case fails, document the specific failure (which task differs, which field is wrong) as implementation notes for fixing in a follow-up

## Implementation Notes

"Functionally identical" means: same task count (+/-1 tolerance), same dependency graph structure, same type/scope assignments for each task, same PRD coverage. Wording may differ but structural elements must match.

The proposal's success criteria (lines 178-189) specify:
- 3 consecutive structurally equivalent runs for execution stability
- 5-minute rule location test for learning curve (optional for this task — record time as informational)
- Error handling validation via temporary rule file removal

If baseline outputs from the original SKILL.md don't exist for the test cases, generate them first by running the original SKILL.md against the same inputs, saving the index.json as baseline.
