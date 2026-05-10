---
status: "completed"
started: "2026-05-10 13:11"
completed: "2026-05-10 13:18"
time_spent: "~7m"
---

# Task Record: 1 Create validate-specs.mjs validation script + update package.json

## Summary
Created validate-specs.mjs ts-morph-based AST validation script with 8 rules (4 ERROR + 4 WARNING) for generated Playwright spec files. Updated package.json template to include ts-morph as devDependency. Created comprehensive test suite with 15 tests covering all rules.

## Changes

### Files Created
- plugins/forge/skills/gen-test-scripts/templates/validate-specs.mjs
- plugins/forge/skills/gen-test-scripts/templates/validate-specs.test.mjs
- plugins/forge/skills/gen-test-scripts/templates/__test_fixtures__/clean-spec.ts
- plugins/forge/skills/gen-test-scripts/templates/__test_fixtures__/e1-waitfor-timeout.ts
- plugins/forge/skills/gen-test-scripts/templates/__test_fixtures__/e3-no-traceability.ts
- plugins/forge/skills/gen-test-scripts/templates/__test_fixtures__/e4-dom-traversal.ts
- plugins/forge/skills/gen-test-scripts/templates/__test_fixtures__/w1-large-serial.ts
- plugins/forge/skills/gen-test-scripts/templates/__test_fixtures__/w3-before-each-login.ts
- plugins/forge/skills/gen-test-scripts/templates/__test_fixtures__/w4-css-class-selector.ts
- plugins/forge/skills/gen-test-scripts/templates/__test_fixtures__/test-cases.md
- plugins/forge/skills/gen-test-scripts/templates/__test_fixtures__/empty-dir

### Files Modified
- plugins/forge/skills/gen-test-scripts/templates/package.json

### Key Decisions
- Used ts-morph ^21.0.0 for AST analysis — compatible with TS 5.x/6.x per proposal
- Fallback on ts-morph parse failure: reports WARNING not ERROR, per proposal risk mitigation
- E3 Traceability check looks both for leading comments and inline string literals containing 'Traceability:'
- E2 TC ID coverage uses regex matching across all spec files, optional via --test-cases flag
- W1 counts all test() calls (not test.describe/test.skip) inside serial describe blocks
- Script accepts spec directory as first arg, optional --test-cases path for E2 coverage check

## Test Results
- **Tests Executed**: Yes
- **Passed**: 15
- **Failed**: 0
- **Coverage**: 100.0%

## Acceptance Criteria
- [x] validate-specs.mjs detects E1: waitForTimeout / setTimeout usage via AST CallExpression analysis
- [x] validate-specs.mjs detects E2: TC ID full coverage — all TC-\d+ from test-cases.md must appear in spec files
- [x] validate-specs.mjs detects E3: every test() has a Traceability comment (// Traceability:)
- [x] validate-specs.mjs detects E4: no DOM parent traversal locator('..')
- [x] validate-specs.mjs detects W1: serial suite with >15 test() calls
- [x] validate-specs.mjs detects W2: serial suite without afterAll cleanup
- [x] validate-specs.mjs detects W3: beforeEach containing login/loginViaUI calls
- [x] validate-specs.mjs detects W4: CSS class selectors (.xxx pattern in locator strings)
- [x] Script outputs structured JSON with { errors: [{id, rule, file, line, message}], warnings: [...] }
- [x] Exit code: 0 if no errors (warnings OK), 1 if any errors present
- [x] ts-morph listed as devDependency in package.json template
- [x] Script handles missing files gracefully (reports ERROR, doesn't crash)

## Notes
Test fixtures are included alongside the script for future regression testing. package-lock.json and node_modules were created by npm install during testing — only package.json change (ts-morph devDependency) should be committed.
