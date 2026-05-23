---
status: "completed"
started: "2026-05-23 13:38"
completed: "2026-05-23 13:43"
time_spent: "~5m"
---

# Task Record: 8 Update prompt templates for type-specific record field awareness

## Summary
Add category-specific record field hints to 19 prompt templates in forge-cli/pkg/prompt/data/. Coding templates now mention testsPassed/testsFailed/coverage, doc templates mention referencedDocs/reviewStatus/docMetrics, test templates mention casesGenerated/casesEvaluated/scriptsCreated, validation templates mention validationPassed/issuesFound, and gate template mentions gatePassed/gateChecks. Changes are purely additive (appended before the final Output line) with no restructuring of existing content.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/coding-feature.md
- forge-cli/pkg/prompt/data/coding-enhancement.md
- forge-cli/pkg/prompt/data/coding-cleanup.md
- forge-cli/pkg/prompt/data/coding-refactor.md
- forge-cli/pkg/prompt/data/coding-fix.md
- forge-cli/pkg/prompt/data/doc.md
- forge-cli/pkg/prompt/data/doc-eval.md
- forge-cli/pkg/prompt/data/doc-summary.md
- forge-cli/pkg/prompt/data/doc-consolidate.md
- forge-cli/pkg/prompt/data/doc-drift.md
- forge-cli/pkg/prompt/data/test-gen-cases.md
- forge-cli/pkg/prompt/data/test-eval-cases.md
- forge-cli/pkg/prompt/data/test-gen-scripts.md
- forge-cli/pkg/prompt/data/test-run.md
- forge-cli/pkg/prompt/data/test-gen-and-run.md
- forge-cli/pkg/prompt/data/test-graduate.md
- forge-cli/pkg/prompt/data/test-verify-regression.md
- forge-cli/pkg/prompt/data/validation-code.md
- forge-cli/pkg/prompt/data/validation-ux.md
- forge-cli/pkg/prompt/data/gate.md
- forge-cli/pkg/prompt/prompt_test.go

### Key Decisions
- Appended Record Fields section before final Output line in each template (additive only)
- fix-record-missed and clean-code templates excluded as recovery/delegate tasks
- Test templates use flexible assertion (any of casesGenerated/casesEvaluated/scriptsCreated) since not all test types use all fields

## Test Results
- **Tests Executed**: Yes
- **Passed**: 19
- **Failed**: 0
- **Coverage**: 90.0%

## Acceptance Criteria
- [x] Coding prompt templates mention: testsPassed, testsFailed, coverage are required for completed tasks
- [x] Doc prompt templates mention: referencedDocs, reviewStatus, docMetrics are recommended fields
- [x] Test prompt templates mention: casesGenerated, casesEvaluated, scriptsCreated are relevant fields
- [x] Validation prompt templates mention: validationPassed, issuesFound are relevant fields
- [x] Gate prompt template mentions: gatePassed, gateChecks are relevant fields
- [x] Changes are additive (append record-field hints) — no restructuring of existing prompt content

## Notes
19 templates updated (5 coding, 5 doc, 7 test, 2 validation, 1 gate). 2 templates excluded (fix-record-missed, clean-code). 5 new test functions added covering all categories. Package coverage: 90.0%.
