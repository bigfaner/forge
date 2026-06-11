---
status: "completed"
started: "2026-05-15 22:02"
completed: "2026-05-15 22:04"
time_spent: "~2m"
---

# Task Record: 4 Update guide.md to reflect spec drift detection flow

## Summary
Updated guide.md to document spec drift detection: updated T-test-5 mermaid node to include drift audit, expanded Quick Mode from T-quick-1~5 to T-quick-1~6 with drift detection as final test step, updated specs/ directory rule to mention drift audit, and updated agent note about docs/business-rules/ and docs/conventions/ to mention drift verification.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/hooks/guide.md

### Key Decisions
- Kept changes minimal and additive per hard rules - only extended existing labels/text rather than restructuring sections

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] Skill Workflow mermaid diagram updated: T-test-5 node shows consolidate-specs + drift audit
- [x] Quick Mode section updated: T-quick-1~6 (was T-quick-1~5), mentions drift detection as the final test step
- [x] specs/ rule in Directory Conventions updated to mention drift detection
- [x] Agent note about docs/business-rules/ and docs/conventions/ updated to mention drift verification

## Notes
Documentation-only task. All four acceptance criteria met with minimal additive changes.
