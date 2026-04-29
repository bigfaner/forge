---
status: "completed"
started: "2026-04-29 19:53"
completed: "2026-04-29 20:08"
time_spent: "~15m"
---

# Task Record: T-test-2 Generate e2e Test Scripts

## Summary
Generated executable TypeScript e2e test scripts from 20 CLI test cases. Created cli.spec.ts with full traceability comments, helpers.ts with runCli/readProjectFile utilities, package.json and tsconfig.json. Scripts placed in both testing/scripts/ (development) and tests/e2e/justfile-e2e-integration/ (canonical location). 13 tests pass immediately; 7 are expected RED (testing justfile recipes e2e-setup/e2e-verify not yet added, and run-e2e-tests skill update not yet done).

## Changes

### Files Created
- docs/features/justfile-e2e-integration/testing/scripts/cli.spec.ts
- docs/features/justfile-e2e-integration/testing/scripts/helpers.ts
- docs/features/justfile-e2e-integration/testing/scripts/package.json
- docs/features/justfile-e2e-integration/testing/scripts/tsconfig.json
- tests/e2e/justfile-e2e-integration/cli.spec.ts
- tests/e2e/helpers.ts
- tests/e2e/package.json
- tests/e2e/tsconfig.json

### Files Modified
无

### Key Decisions
- All 20 test cases are CLI type — no UI or API spec files needed
- File-content tests (TC-001 to TC-019 except TC-003/004/009-013) read plugin markdown files directly via readProjectFile helper
- just command tests (TC-003/004/009-012/020) are expected RED until justfile e2e-setup/e2e-verify recipes are added by task 1.1
- TC-013 is expected RED until run-e2e-tests SKILL.md is updated to reference justfile check
- PROJECT_ROOT computed as 2 levels up from tests/e2e/ to reach forge project root

## Test Results
- **Passed**: 13
- **Failed**: 7
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] tests/e2e/<feature>/ contains at least one spec file (cli.spec.ts)
- [x] tests/e2e/helpers.ts exists (shared infrastructure)
- [x] Each test() includes traceability comment // Traceability: TC-NNN → {PRD Source}

## Notes
无
