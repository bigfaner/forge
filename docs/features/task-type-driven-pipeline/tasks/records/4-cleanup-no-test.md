---
status: "completed"
started: "2026-05-14 23:26"
completed: "2026-05-14 23:39"
time_spent: "~13m"
---

# Task Record: 4 Remove --no-test flag and update skill documentation

## Summary
Removed --no-test CLI flag and BuildIndexOpts.NoTest field. Auto-detection via isDocsOnlyFeature() replaces manual control. Updated skill docs (quick-tasks/SKILL.md, quick.md), WORKFLOW.md, and e2e test TC-017. Bumped version to 3.9.0 (minor: auto-detection feature).

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/index.go
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/build_test.go
- forge-cli/internal/cmd/index_test.go
- forge-cli/tests/e2e/features/task-stage-gates/task_stage_gates_cli_test.go
- forge-cli/docs/WORKFLOW.md
- plugins/forge/skills/quick-tasks/SKILL.md
- plugins/forge/commands/quick.md
- forge-cli/scripts/version.txt

### Key Decisions
- Rewrote TestBuildIndex_NoTestSkipsTestGen as TestBuildIndex_NoProfilesSkipsTestGen -- without NoTest field, the 'skip test gen' behavior is now achieved by simply not providing TestProfiles
- Rewrote TestBuildIndex_StageGatesWithNoTestFlag as TestBuildIndex_StageGatesWithNoProfiles -- same logic but without NoTest field
- Rewrote TC-017 e2e test to verify --no-test returns unknown flag error instead of testing its behavior
- Preserved Task.NoTest (per-task noTest frontmatter field) -- this is a different concept from BuildIndexOpts.NoTest and controls quality gate skip for individual tasks

## Test Results
- **Tests Executed**: Yes
- **Passed**: 163
- **Failed**: 0
- **Coverage**: 89.8%

## Acceptance Criteria
- [x] forge task index --no-test --feature <slug> returns unknown flag: --no-test error
- [x] BuildIndexOpts struct has no NoTest field
- [x] index.go has no reference to --no-test or indexNoTest
- [x] quick-tasks/SKILL.md has no mention of --no-test
- [x] quick.md command has no --no-test flag section or passthrough
- [x] build.go has no reference to opts.NoTest
- [x] Version bumped in forge-cli/scripts/version.txt
- [x] Existing tests in index_test.go pass after removing the flag
- [x] Running forge task index --feature task-type-driven-pipeline (no --no-test) works correctly

## Notes
无
