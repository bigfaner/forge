---
id: "3"
title: "Remove EXTREMELY-IMPORTANT redundancy from skills"
priority: "P1"
estimated_time: "1.5h"
dependencies: [2]
type: "doc"
mainSession: false
---

# 3: Remove EXTREMELY-IMPORTANT redundancy from skills

## Description

Trim EXTREMELY-IMPORTANT blocks where items duplicate body constraints. Apply constraint-level audit: keep E-I items where body lacks the constraint, or body has it at lower enforcement level.

## Reference Files
- `docs/proposals/skill-instruction-audit/proposal.md#Success-Criteria`: SC-2 constraint-level audit
- `plugins/forge/skills/test-guide/SKILL.md`: 10-item E-I block (8 duplicates)
- `plugins/forge/skills/eval/SKILL.md`: E-I block with 3 overlapping rules
- `plugins/forge/skills/init-justfile/SKILL.md`: Notes overlapping E-I block
- `plugins/forge/skills/clean-code/SKILL.md`: 3 "preserve scope" principles + HARD-RULE
- `plugins/forge/skills/deep-research/SKILL.md`: Report Structure key points duplicating template

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/test-guide/SKILL.md` | Remove 8 of 10 E-I items; keep cross-step guardrails |
| `plugins/forge/skills/eval/SKILL.md` | Merge 3 overlapping rules into 1 |
| `plugins/forge/skills/init-justfile/SKILL.md` | Remove Notes items overlapping E-I block |
| `plugins/forge/skills/clean-code/SKILL.md` | Merge 3 "preserve scope" into 1; remove redundant HARD-RULE |
| `plugins/forge/skills/deep-research/SKILL.md` | Delete Report Structure key points duplicating template |

## Acceptance Criteria

- [ ] `test-guide/SKILL.md` E-I has ≤4 items; each passes constraint-level audit
- [ ] `eval/SKILL.md` has 1 concise E-I rule
- [ ] `init-justfile/SKILL.md` Notes has no E-I overlap
- [ ] `clean-code/SKILL.md` has 1 "preserve scope" statement
- [ ] `deep-research/SKILL.md` has no key points duplicating template

## Hard Rules

- 仅修改上述 5 个文件
- Retained E-I items must satisfy: body lacks constraint, or body has it at lower enforcement level

## Implementation Notes

Constraint-level audit for each E-I item: extract key verb → grep body → same level = delete, lower/absent = keep.

Preserve in test-guide: "preserve approved sections verbatim" (cross-step guardrail), "After 2 retry rejections, write as .draft.md" (MUST NOT force-apply).

Delete in test-guide: "MANUAL-ONLY" (line 11), "Convention file MUST include marker" (Step 5b), "Do NOT execute test commands" (Notes), "Drafts MUST be reviewed" (Step 4).
