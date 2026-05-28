---
status: "completed"
started: "2026-05-28 15:22"
completed: "2026-05-28 15:45"
time_spent: "~23m"
---

# Task Record: 5 Migrate all template frontmatters + PhaseSummary

## Summary
Migrated all 41 template files from flat variables list to grouped frontmatter format (identity/context/conditional). Migrated PhaseSummary from frontmatter variable to independent body section in all 21 prompt templates. Updated prompt.go to remove PHASE_SUMMARY: prefix from phaseSummaryLine.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/prompt.go
- forge-cli/pkg/prompt/templates/code-quality-simplify.md
- forge-cli/pkg/prompt/templates/coding-cleanup.md
- forge-cli/pkg/prompt/templates/coding-enhancement.md
- forge-cli/pkg/prompt/templates/coding-feature.md
- forge-cli/pkg/prompt/templates/coding-fix.md
- forge-cli/pkg/prompt/templates/coding-refactor.md
- forge-cli/pkg/prompt/templates/doc-consolidate.md
- forge-cli/pkg/prompt/templates/doc-drift.md
- forge-cli/pkg/prompt/templates/doc-review.md
- forge-cli/pkg/prompt/templates/doc-summary.md
- forge-cli/pkg/prompt/templates/doc.md
- forge-cli/pkg/prompt/templates/eval-contract.md
- forge-cli/pkg/prompt/templates/eval-journey.md
- forge-cli/pkg/prompt/templates/fix-record-missed.md
- forge-cli/pkg/prompt/templates/gate.md
- forge-cli/pkg/prompt/templates/test-gen-contracts.md
- forge-cli/pkg/prompt/templates/test-gen-journeys.md
- forge-cli/pkg/prompt/templates/test-gen-scripts.md
- forge-cli/pkg/prompt/templates/test-run.md
- forge-cli/pkg/prompt/templates/validation-code.md
- forge-cli/pkg/prompt/templates/validation-ux.md
- forge-cli/pkg/task/templates/code-quality-simplify.md
- forge-cli/pkg/task/templates/coding.cleanup.md
- forge-cli/pkg/task/templates/coding.fix.md
- forge-cli/pkg/task/templates/doc-consolidate.md
- forge-cli/pkg/task/templates/doc-drift.md
- forge-cli/pkg/task/templates/doc-review.md
- forge-cli/pkg/task/templates/eval-contract.md
- forge-cli/pkg/task/templates/eval-journey.md
- forge-cli/pkg/task/templates/test-gen-contracts.md
- forge-cli/pkg/task/templates/test-gen-journeys.md
- forge-cli/pkg/task/templates/test-gen-scripts.md
- forge-cli/pkg/task/templates/test-run.md
- forge-cli/pkg/task/templates/validation-code.md
- forge-cli/pkg/task/templates/validation-ux.md
- forge-cli/pkg/task/records/coding.md
- forge-cli/pkg/task/records/doc.md
- forge-cli/pkg/task/records/eval.md
- forge-cli/pkg/task/records/gate.md
- forge-cli/pkg/task/records/test.md
- forge-cli/pkg/task/records/validation.md

### Key Decisions
- Complexity placed in context group for all 5 coding templates that use it
- coding.cleanup and coding.fix task templates use identity/context groups with their unique fields (ID, Title, SourceFiles, etc.)
- PhaseSummary body section wrapped in {{if .PhaseSummary}}...{{end}} conditional block

## Test Results
- **Tests Executed**: Yes
- **Passed**: 57
- **Failed**: 0
- **Coverage**: 75.0%

## Acceptance Criteria
- [x] SC-FM-1: All 41 templates migrated with correct groups
- [x] SC-FM-4: PhaseSummary removed from frontmatter, added as body section in 21 prompt templates
- [x] SC-FM-5: All grouped field names use PascalCase
- [x] SC-FM-6: Rendered frontmatter (second --- block) unchanged
- [x] prompt.go phaseSummaryLine updated: PHASE_SUMMARY prefix removed

## Notes
Batch migration via Python script. All existing tests pass including template validation, synthesis, and metadata parsing tests.
