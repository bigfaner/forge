---
id: "6"
title: "Extract eval validate-ux sub-pipeline to rubric file"
priority: "P2"
estimated_time: "1h"
dependencies: []
type: "refactor"
scope: "all"
breaking: false
mainSession: false
---

# 6: Extract eval validate-ux sub-pipeline to rubric file

## Description

`eval/SKILL.md` has 90+ lines of validate-ux sub-pipeline inline (lines 134-225): project-type detection, PRD-to-operation translation, snapshot format, 7 impact types. This is cognitive overload and inconsistent with the gen-test-scripts extraction pattern (which uses `types/` dispatch files). Extract to a rubric file.

## Reference Files
- `docs/proposals/skill-ecosystem-audit/proposal.md` — Source proposal (W4, item 11)
- `plugins/forge/skills/eval/SKILL.md` — Current 366-line file

## Affected Files

### Create

| File | Description |
|------|-------------|
| `plugins/forge/skills/eval/rubrics/validate-ux-pipeline.md` | Extracted validate-ux sub-pipeline (project-type detection, operation translation, snapshot format, impact types) |

### Modify

| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/SKILL.md` | Replace inline validate-ux sub-pipeline with reference to rubric file. Target: under 280 lines. |

## Acceptance Criteria

- [ ] `eval/SKILL.md` under 280 lines (currently 366)
  `wc -l plugins/forge/skills/eval/SKILL.md` returns ≤ 280
- [ ] `validate-ux-pipeline.md` exists in `eval/rubrics/`
  `ls plugins/forge/skills/eval/rubrics/validate-ux-pipeline.md` succeeds
- [ ] No content is lost — extracted content matches the original inline version
- [ ] eval still correctly dispatches to validate-ux scoring

## Hard Rules

- This is a rubric file consumed by eval, NOT a new skill. Place in `eval/rubrics/`.
- Use `${CLAUDE_SKILL_DIR}/rubrics/validate-ux-pipeline.md` for reference from SKILL.md.
- Do not modify the scoring logic or impact types during extraction.

## Implementation Notes

- The inline content includes: project-type detection (web/tui/mobile), PRD-to-operation translation, snapshot format spec, and 7 impact type definitions.
- The rubric file should be self-contained so eval can load it independently.
- The SKILL.md should reference it with a one-line pointer and a brief summary of what it contains.
