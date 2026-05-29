---
status: "completed"
started: "2026-05-29 17:06"
completed: "2026-05-29 17:11"
time_spent: "~5m"
---

# Task Record: 2 Update write-prd SKILL.md for refactor intent branch

## Summary
Updated write-prd SKILL.md with intent detection logic and spec-only PRD branch for refactor/cleanup intents

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/write-prd/SKILL.md

### Key Decisions
无

## Document Metrics
~110 lines added: Intent Detection section, dual Process Flow, dual Checklist, dual Output Documents, Step 7A with 3 mandatory fields, intent gates on Steps 7/8/9

## Referenced Documents
- docs/proposals/intent-driven-pipeline-branching/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] write-prd SKILL.md contains intent detection logic: when proposal.md frontmatter intent is refactor, execute spec-only PRD branch
- [x] spec-only PRD format contains three mandatory fields: change scope (affected modules/files), constraints (behavioral invariants), verification criteria (regression acceptance criteria)
- [x] refactor branch does not generate prd-user-stories.md file

## Notes
Added Intent Detection subsection in Prerequisites with detection bash command. Added EXTREMELY-IMPORTANT intent gates on Steps 7, 7A, and 8. Step 7A defines the three mandatory fields with table and markdown template. Process Flow, Checklist, and Output Documents sections all split into new-feature vs refactor/cleanup variants. Step 9 manifest updated to conditionally include User Stories row.
