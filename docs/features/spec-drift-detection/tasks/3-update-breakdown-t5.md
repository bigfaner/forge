---
id: "3"
title: "Update breakdown-tasks SKILL.md T-test-5 description for drift verification"
priority: "P2"
estimated_time: "15-30min"
dependencies: ["1"]
type: "documentation"
mainSession: false
---

# 3: Update breakdown-tasks SKILL.md T-test-5 description for drift verification

## Description

Update the breakdown-tasks SKILL.md to reflect that T-test-5 (Consolidate Specs) now includes drift verification as part of its workflow, not just spec extraction and integration.

## Reference Files
- `docs/proposals/spec-drift-detection/proposal.md` — Source proposal
- `plugins/forge/skills/breakdown-tasks/SKILL.md` — File to modify

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Update T-test-5 description to include drift detection |

## Acceptance Criteria

- [ ] breakdown-tasks SKILL.md references that T-test-5 now includes drift detection after spec consolidation
- [ ] Description clarifies the full pipeline: extract → integrate → detect drift → auto-fix

## Hard Rules

- Only modify `plugins/forge/skills/breakdown-tasks/SKILL.md`
- Do not change task generation logic (that's in Go code, handled by Task 2)

## Implementation Notes

- The current text at line 298 mentions "test pipeline (T-test-1 through T-test-5)" — ensure T-test-5's expanded scope is documented
- Keep the change minimal — this is descriptive documentation, not prescriptive workflow
