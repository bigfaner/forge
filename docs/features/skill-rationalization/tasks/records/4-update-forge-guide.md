---
status: "completed"
started: "2026-05-16 01:12"
completed: "2026-05-16 01:18"
time_spent: "~6m"
---

# Task Record: 4 Update forge guide and references for new eval structure

## Summary
Updated forge guide, eval-forge audit skill, scorer prompt, rubric, CLI prompt template, and README to reference the new consolidated eval skill structure (skills/eval/ + rubrics/ + command wrappers) instead of 7 separate eval skills. Skill count updated from 23 to 17 in README.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/hooks/guide.md
- .claude/skills/eval-forge/SKILL.md
- .claude/skills/eval-forge/templates/scorer-prompt.md
- .claude/skills/eval-forge/templates/rubric.md
- forge-cli/pkg/prompt/data/test-pipeline-eval-cases.md
- README.md

### Key Decisions
- Slash commands (/eval-prd, /eval-ui, etc.) remain unchanged in mermaid diagrams and workflow descriptions since command wrappers preserve the same UX
- Eval-forge audit now validates skills/eval/ + rubrics/ structure and command wrapper delegation instead of checking individual eval-* directories
- Updated rubric D6b criterion from checking eval-*/templates/rubric.md to checking skills/eval/rubrics/ and command wrapper delegation

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] Forge guide mentions generic eval skill instead of 7 separate eval skills
- [x] eval-forge audit skill correctly validates the new eval structure (skills/eval/ + rubrics/ + command wrappers)
- [x] CLI prompt templates reference correct paths for the new eval skill
- [x] Skill count in documentation reflects 17 (down from 24)

## Notes
Documentation-only task. No code changes. hooks.json had no eval references to update.
