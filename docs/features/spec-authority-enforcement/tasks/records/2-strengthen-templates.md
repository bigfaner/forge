---
status: "completed"
started: "2026-05-24 09:44"
completed: "2026-05-24 09:48"
time_spent: "~4m"
---

# Task Record: 2 Strengthen templates with Reference Files declaration + AC validation

## Summary
Inserted <IMPORTANT> Reference Files authority declaration and AC per-item validation into all 9 templates identified by Task 1 audit: coding-feature, coding-enhancement, coding-refactor, coding-fix, coding-cleanup, doc, doc-review, gate, validation-code, validation-ux. Build and all existing tests verified successfully.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/data/coding-feature.md
- forge-cli/pkg/prompt/data/coding-enhancement.md
- forge-cli/pkg/prompt/data/coding-refactor.md
- forge-cli/pkg/prompt/data/coding-fix.md
- forge-cli/pkg/prompt/data/coding-cleanup.md
- forge-cli/pkg/prompt/data/doc.md
- forge-cli/pkg/prompt/data/doc-review.md
- forge-cli/pkg/prompt/data/gate.md
- forge-cli/pkg/prompt/data/validation-code.md
- forge-cli/pkg/prompt/data/validation-ux.md

### Key Decisions
- Used <IMPORTANT> tag (not <EXTREMELY-IMPORTANT>) per Hard Rules to avoid marker dilution
- Inserted Reference Files declaration as a separate <IMPORTANT> block before existing Hard Rules <IMPORTANT> blocks
- Inserted AC validation as first sub-step in Verify/Self-Check steps before static checks
- coding-feature.md and coding-enhancement.md received identical modifications per SYNC NOTICE

## Test Results
- **Tests Executed**: Yes
- **Passed**: 136
- **Failed**: 0
- **Coverage**: 90.9%

## Acceptance Criteria
- [x] Each template identified by Task 1 has a <IMPORTANT> Reference Files authority declaration inserted in Step 1
- [x] Declaration text matches the proposal's exact template (4 MUST items + 2 fallback outputs)
- [x] Each template has AC per-item validation inserted in its Verify/Self-Check step as the first sub-step
- [x] AC validation text matches the proposal's exact template (per-item PASS/FAIL + skip condition)
- [x] Use <IMPORTANT> tag (not <EXTREMELY-IMPORTANT>) to avoid marker dilution
- [x] go build ./... succeeds after all modifications
- [x] Existing <IMPORTANT> Hard Rules blocks in templates are NOT modified or removed

## Notes
All 136 existing tests in pkg/prompt passed (90.9% coverage). Template-only changes verified via go build and existing test suite. Existing <IMPORTANT> Hard Rules blocks preserved unchanged in all 9 templates.
