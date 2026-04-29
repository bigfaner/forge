---
status: "completed"
started: "2026-04-29 20:15"
completed: "2026-04-29 20:17"
time_spent: "~2m"
---

# Task Record: disc-1 Fix: add e2e-setup and e2e-verify recipes to Justfile

## Summary
Added e2e-setup and e2e-verify recipes to the project Justfile. e2e-setup installs npm deps and playwright chromium idempotently, exits 1 if tests/e2e/package.json is missing. e2e-verify --feature <slug> scans tests/e2e/<slug>/*.spec.ts for // VERIFY: markers, exits 1 with file:line output if found, exits 0 with OK message if clean.

## Changes

### Files Created
无

### Files Modified
- Justfile

### Key Decisions
- Copied recipe definitions verbatim from init-justfile.md to ensure consistency with documented contract
- Used [arg('feature', long)] syntax requiring just >= 1.50.0 (confirmed installed at 1.50.0)

## Test Results
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] e2e-setup recipe added to Justfile
- [x] e2e-setup exits 1 if tests/e2e/package.json missing
- [x] e2e-setup installs npm deps and playwright chromium idempotently
- [x] e2e-verify --feature <slug> recipe added to Justfile
- [x] e2e-verify exits 1 with file:line output if // VERIFY: markers found
- [x] e2e-verify exits 0 with OK message if no markers found

## Notes
无
