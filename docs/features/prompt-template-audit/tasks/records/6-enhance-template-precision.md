---
status: "completed"
started: "2026-05-20 17:02"
completed: "2026-05-20 17:09"
time_spent: "~7m"
---

# Task Record: 6 P2: 增强模板指令精确度

## Summary
Enhanced template precision across 6 prompt templates: (1) Added AC extraction guidance to coding-feature.md and coding-enhancement.md Step 2, guiding agents to extract test requirements from Acceptance Criteria before TDD; (2) Replaced all Go-specific test commands (go test -race -cover ./changed/package/...) with language-agnostic descriptions in all 5 coding templates; (3) Added task-type-specific execution guidance (Create/Modify/Delete) to doc.md Step 2.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/coding-feature.md
- forge-cli/pkg/prompt/data/coding-enhancement.md
- forge-cli/pkg/prompt/data/coding-cleanup.md
- forge-cli/pkg/prompt/data/coding-fix.md
- forge-cli/pkg/prompt/data/coding-refactor.md
- forge-cli/pkg/prompt/data/doc.md

### Key Decisions
- Kept 'go test' as an example in the language-agnostic description (alongside pytest, jest) rather than removing all framework mentions entirely
- Replaced 4 additional Go-specific commands in coding-refactor.md (Pre-check, Phase A, Phase B, Behavioral sections) beyond the Targeted Tests section, since the task said 'all Go-specific examples'
- For doc.md, chose a simple 3-type taxonomy (Create/Modify/Delete) that covers all document task variations without over-engineering

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] coding-feature.md Step 2 contains AC extraction guidance
- [x] coding-enhancement.md Step 2 contains same AC extraction guidance
- [x] All 5 coding templates replace Go-specific test examples with language-agnostic descriptions
- [x] doc.md Step 2 contains task-type-specific execution guidance (identify type -> execute by type)

## Notes
无
