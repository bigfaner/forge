---
status: "completed"
started: "2026-05-20 10:12"
completed: "2026-05-20 10:16"
time_spent: "~4m"
---

# Task Record: T-test-gen-cases Generate e2e Test Cases

## Summary
Generated 20 CLI test cases from PRD acceptance criteria for test-knowledge-convention-driven feature. Covered all 6 user stories and key spec functional specs (FS-1 through FS-9). Created per-type test case file and manifest with full traceability.

## Changes

### Files Created
- docs/features/test-knowledge-convention-driven/testing/cli-test-cases.md
- docs/features/test-knowledge-convention-driven/testing/manifest.md

### Files Modified
无

### Key Decisions
- Classified all test cases as CLI type since the project interface is cli-only (no web-ui, api, tui, or mobile interfaces detected)
- Mapped PRD stories to concrete CLI commands: gen-test-scripts, test-guide, forge-commands, forge-build, config-init, task-index, init-justfile, consolidate-specs
- Included spec-level acceptance criteria (FS-1 validation, FS-2 error handling, FS-3 reliability, FS-7 import audit, FS-6/FS-8 config cleanup, FS-9 drift detection) as additional test cases beyond story-level ACs
- Added performance test cases (TC-018 generation time budget, TC-019 first-pass compile rate) from spec goals section

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] All 6 PRD user stories covered by at least one test case each
- [x] Every test case traces to a specific PRD source (Story/AC or Spec section)
- [x] Test cases follow CLI format with Target, Test ID, Pre-conditions, Steps, Expected, Priority
- [x] Manifest generated with Summary and Cross-Type Traceability tables
- [x] Priority assignment follows P0 for core stories, P1 for error/edge cases, P2 for nice-to-have

## Notes
Task type is test.gen-cases (noTest: true). No test execution required. CLI-only interface detected via forge test interfaces.
