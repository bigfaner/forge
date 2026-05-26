---
status: "completed"
started: "2026-05-27 01:08"
completed: "2026-05-27 01:10"
time_spent: "~2m"
---

# Task Record: 5 Fix prompt template references (test-run.md and SCOPE labels)

## Summary
Fixed two critical prompt template issues: (1) test-run.md now references forge:run-tests instead of non-existent forge:run-e2e-tests, (2) renamed SCOPE: label to SURFACE_KEY: across all 18 prompt templates and updated Go code comments in prompt.go and prompt_test.go

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/test-run.md
- forge-cli/pkg/prompt/data/code-quality-simplify.md
- forge-cli/pkg/prompt/data/coding-cleanup.md
- forge-cli/pkg/prompt/data/coding-enhancement.md
- forge-cli/pkg/prompt/data/coding-feature.md
- forge-cli/pkg/prompt/data/coding-fix.md
- forge-cli/pkg/prompt/data/coding-refactor.md
- forge-cli/pkg/prompt/data/doc-summary.md
- forge-cli/pkg/prompt/data/eval-contract.md
- forge-cli/pkg/prompt/data/eval-journey.md
- forge-cli/pkg/prompt/data/fix-record-missed.md
- forge-cli/pkg/prompt/data/gate.md
- forge-cli/pkg/prompt/data/test-gen-contracts.md
- forge-cli/pkg/prompt/data/test-gen-journeys.md
- forge-cli/pkg/prompt/data/test-gen-scripts.md
- forge-cli/pkg/prompt/data/test-verify-regression.md
- forge-cli/pkg/prompt/data/validation-code.md
- forge-cli/pkg/prompt/data/validation-ux.md
- forge-cli/pkg/prompt/prompt.go
- forge-cli/pkg/prompt/prompt_test.go

### Key Decisions
- Used SURFACE_KEY: as label prefix (not SURFACE:) to match the variable name for consistency
- Updated isLabelWithEmptyValue comments but not its logic since it already handles any label prefix generically

## Test Results
- **Tests Executed**: Yes
- **Passed**: 62
- **Failed**: 0
- **Coverage**: 82.7%

## Acceptance Criteria
- [x] test-run.md references forge:run-tests (not forge:run-e2e-tests)
- [x] All 18 prompt templates use SURFACE_KEY: label instead of SCOPE:
- [x] Go code that injects/replaces these labels is updated to match
- [x] Existing tests pass (go test ./...)

## Notes
18 templates had SCOPE: label (proposal said 17 but actual count was 18). doc-review.md has a different 'SCOPE CONSTRAINT:' directive which is semantically unrelated and was left unchanged.
