---
status: "completed"
started: "2026-05-18 01:54"
completed: "2026-05-18 01:54"
time_spent: ""
---

# Task Record: 6 Extract eval validate-ux sub-pipeline to rubric file

## Summary
Extracted the validate-ux sub-pipeline (project-type detection, PRD-to-operation translation, ux-snapshot.md format, 7 impact types) from eval/SKILL.md into rubrics/validate-ux-pipeline.md. SKILL.md reduced from 368 to 277 lines. Also fixed pre-existing argument-hint to argument-hints in extract-design-md.md.

## Changes

### Files Created
- plugins/forge/skills/eval/rubrics/validate-ux-pipeline.md

### Files Modified
- plugins/forge/skills/eval/SKILL.md
- plugins/forge/commands/extract-design-md.md

### Key Decisions
- Placed rubric file in eval/rubrics/ per Hard Rules, not as a new skill
- Replaced the 90-line inline validate-ux block with a concise table row referencing ${CLAUDE_SKILL_DIR}/rubrics/validate-ux-pipeline.md
- Rubric file is self-contained with 4 sections matching the original inline content exactly
- Fixed argument-hint to argument-hints in extract-design-md.md to clear pre-existing test failure from task 1 scope

## Test Results
- **Tests Executed**: Yes
- **Passed**: 1061
- **Failed**: 0
- **Coverage**: 88.5%

## Acceptance Criteria
- [x] eval/SKILL.md under 280 lines (was 368)
- [x] validate-ux-pipeline.md exists in eval/rubrics/
- [x] No content is lost -- extracted content matches the original inline version
- [x] eval still correctly dispatches to validate-ux scoring

## Notes
Also fixed pre-existing argument-hint frontmatter in extract-design-md.md (task 1 scope) to clear the test suite blockage.
