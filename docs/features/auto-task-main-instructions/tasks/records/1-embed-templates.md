---
status: "completed"
started: "2026-05-22 22:29"
completed: "2026-05-22 22:41"
time_spent: "~12m"
---

# Task Record: 1 Create embed template files for auto-gen task types and update GenerateTestTaskMD()

## Summary
Created 13 embed template files in forge-cli/pkg/task/data/ for each auto-gen task type. Updated GenerateTestTaskMD() to use embed.FS to load template files as task body content, replacing the previous sparse body generation. StrategyContent is appended after template content, TestType is noted in body, and missing template files fall back to legacy behavior. Frontmatter generation logic is unchanged.

## Changes

### Files Created
- forge-cli/pkg/task/data/test-gen-cases.md
- forge-cli/pkg/task/data/test-eval-cases.md
- forge-cli/pkg/task/data/test-gen-scripts.md
- forge-cli/pkg/task/data/test-gen-and-run.md
- forge-cli/pkg/task/data/test-run.md
- forge-cli/pkg/task/data/test-graduate.md
- forge-cli/pkg/task/data/test-verify-regression.md
- forge-cli/pkg/task/data/validation-code.md
- forge-cli/pkg/task/data/validation-ux.md
- forge-cli/pkg/task/data/doc-eval.md
- forge-cli/pkg/task/data/doc-consolidate.md
- forge-cli/pkg/task/data/doc-drift.md
- forge-cli/pkg/task/data/code-quality-simplify.md

### Files Modified
- forge-cli/pkg/task/autogen.go
- forge-cli/pkg/task/autogen_test.go

### Key Decisions
- Template content mirrors the prompt/data/*.md templates but adapted for task file body (no {{PLACEHOLDER}} substitution needed)
- Type-to-filename mapping follows same convention as prompt.go: type '.' replaced by '-'
- Fallback to legacy behavior when template file is missing or unreadable

## Test Results
- **Tests Executed**: Yes
- **Passed**: 8
- **Failed**: 0
- **Coverage**: 87.6%

## Acceptance Criteria
- [x] forge-cli/pkg/task/data/ directory created with 13 .md template files
- [x] Each template file contains meaningful Main Instructions for its task type
- [x] Type-to-filename mapping added (convention: type with '.' replaced by '-')
- [x] GenerateTestTaskMD() uses embed.FS to load template file as body
- [x] When StrategyContent is non-empty, it's appended AFTER the template content
- [x] When TestType is non-empty, it's noted in the body (as current behavior)
- [x] forge task index --feature <slug> generates task files with the new body content
- [x] Existing tests pass (backward compatible frontmatter)

## Notes
Hard Rules verified: frontmatter generation logic unchanged, AutoGenTaskDef struct unchanged, embed.FS used (same pattern as prompt.go), missing template file falls back gracefully to current behavior.
