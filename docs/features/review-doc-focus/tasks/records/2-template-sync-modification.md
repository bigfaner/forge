---
status: "completed"
started: "2026-05-24 19:18"
completed: "2026-05-24 19:22"
time_spent: "~4m"
---

# Task Record: 2 Template sync modification (autogen + prompt templates)

## Summary
Synced autogen template (task/data/doc-review.md) and agent prompt template (prompt/data/doc-review.md): autogen template now includes {{DOC_TASK_AC}} placeholder and AC Summary section with allowlist Discovery Strategy excluding tasks/records/manifest.md/index.json; prompt template refactored to 4-step workflow (Load Pre-extracted AC -> Discover Target Documents via allowlist -> Review & Fix with docs-only constraint -> Report Summary) with explicit prohibition on scanning tasks directory and modifying non-docs files.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/data/doc-review.md
- forge-cli/pkg/prompt/data/doc-review.md
- forge-cli/pkg/task/autogen_test.go
- forge-cli/pkg/prompt/prompt_test.go

### Key Decisions
- Prompt template reduced from 5 steps to 4 steps (merged AC reading and deliverable reading into combined discover step) while preserving workflow logic per Hard Rules
- Used IMPORTANT block in Step 3 for docs-only constraint to ensure visibility to agent
- Autogen template uses EXCLUDE block pattern for clarity instead of implicit allowlist-only exclusion

## Test Results
- **Tests Executed**: Yes
- **Passed**: 10
- **Failed**: 0
- **Coverage**: 89.3%

## Acceptance Criteria
- [x] autogen template contains {{DOC_TASK_AC}} placeholder
- [x] autogen template Discovery Strategy uses allowlist: only scan docs/ subtree .md files
- [x] autogen template excludes tasks/, records/, manifest.md, index.json
- [x] prompt template Step 1 is 'Load Pre-extracted AC' (no scan tasks directive)
- [x] prompt template Step 2 uses docs/ allowlist document discovery
- [x] prompt template Step 3 contains 'only modify docs/ files' explicit constraint
- [x] prompt template has no 'scan tasks directory' or 'Read the task acceptance criteria from its .md file' directives

## Notes
Both templates modified in same task per coupling constraint. Existing tests (TestGenerateTestTaskMD_EmbedTemplate_LoadsContent doc-review, TestSynthesize_AllTypes doc.review) continue to pass. Task pkg coverage: 87.7%, prompt pkg coverage: 90.9%.
