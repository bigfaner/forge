---
status: "blocked"
started: "2026-06-07 00:36"
completed: "N/A"
time_spent: ""
---

# Task Record: 5 迁移 skill-ops journey

## Summary
Migrated forge-cli/tests/skill-ops/ to tests/skill-ops/: 4 test files with import path replacement (forge-cli/tests/testkit -> forge-tests/testkit), main_test.go rewritten to ForgeBinary init pattern, 3 contract files copied verbatim. All tests produce identical results to source location.

## Changes

### Files Created
- tests/skill-ops/main_test.go
- tests/skill-ops/clean_code_skill_test.go
- tests/skill-ops/forensic_test.go
- tests/skill-ops/plugin_content_test.go
- tests/skill-ops/prompt_test.go
- tests/skill-ops/contracts/step-1-plugin-validation.md
- tests/skill-ops/contracts/step-2-forensic.md
- tests/skill-ops/contracts/step-3-prompt.md

### Files Modified
无

### Key Decisions
- Reused ForgeBinary init pattern from tests/task-lifecycle/main_test.go instead of custom build logic in source main_test.go
- Kept package name as 'skillops' for consistency with source
- plugin_content_test.go ReadProjectFile path resolution ('../'+relPath) verified to work correctly from new location -- ProjectRoot returns tests/ via runtime.Caller, filepath.Join resolves ../ to project root

## Test Results
- **Tests Executed**: Yes
- **Passed**: 8
- **Failed**: 1
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] tests/skill-ops/ contains all 4 migrated test files with testkit import from forge-tests/testkit
- [x] main_test.go uses ForgeBinary init pattern
- [x] contracts/ directory with 3 contract files correctly migrated
- [ ] just test includes this journey and passes

## Notes
Pre-existing FAIL: TestCleanCode_SkillFile_ContainsFivePrinciples fails because SKILL.md does not contain 'Focus Scope' principle. This failure existed at source location (forge-cli/tests/skill-ops/) before migration -- identical error. Not caused by migration. 1 SKIP (TestTC_005) also matches source behavior. Coverage not measured (0.0) as cli_functional tests do not collect coverage.
