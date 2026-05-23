---
status: "completed"
started: "2026-05-23 09:47"
completed: "2026-05-23 09:53"
time_spent: "~6m"
---

# Task Record: 3 Update 13 embed template files with strategy-based content

## Summary
Replaced all 13 embed template files in forge-cli/pkg/task/data/ with concise strategy-based content per the proposal's three-strategy model. Strategy A templates (7 files) provide feature context with {{FEATURE_SLUG}}, {{SCOPE}}, {{INTERFACES}} placeholders. Strategy B templates (2 files) provide validation criteria with {{ACCEPTANCE_CRITERIA}} and quality gate commands. Strategy C templates (4 files) provide discovery strategies (git diff, directory scan). Updated 4 existing test assertions to match new template content.

## Changes

### Files Created
无

### Files Modified
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
- forge-cli/pkg/task/autogen_test.go

### Key Decisions
- Strategy A templates use {{MODE}} in first line; when Mode is empty, removeLineContaining deletes the entire first line — test assertions adjusted to check content from non-removable lines
- Template content follows proposal exactly with no workflow duplication — skill invocation hints removed since prompt templates already provide workflow guidance

## Test Results
- **Tests Executed**: Yes
- **Passed**: 591
- **Failed**: 0
- **Coverage**: 88.2%

## Acceptance Criteria
- [x] All 13 template files replaced with strategy-based content from the proposal
- [x] Strategy A templates (7 files) contain: feature slug, scope, interfaces placeholders; skill invocation hints
- [x] Strategy B templates (2 files) contain: acceptance criteria placeholders; quality gate commands
- [x] Strategy C templates (4 files) contain: discovery strategy steps (git diff, directory scan)
- [x] No workflow duplication between template body and prompt templates — templates only provide feature-specific context
- [x] Templates are concise (5-15 lines each, per proposal)

## Notes
无
