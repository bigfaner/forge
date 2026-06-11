---
status: "completed"
started: "2026-05-17 17:34"
completed: "2026-05-17 17:45"
time_spent: "~11m"
---

# Task Record: 2 Migrate CLI template and fix wiring

## Summary
Migrated CLI prompt template from code-quality-simplify.md to code-quality-clean-code.md calling forge:clean-code skill. Deleted old template, added TypeCleanCode entry to typeToTemplate mapping using task.TypeCleanCode constant, and bumped version to 3.18.3.

## Changes

### Files Created
- forge-cli/pkg/prompt/data/code-quality-clean-code.md

### Files Modified
- forge-cli/pkg/prompt/prompt.go
- forge-cli/pkg/prompt/prompt_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- New template follows existing pattern (test-pipeline-gen-cases.md) with Hard Rules section mandating forge:clean-code skill invocation
- typeToTemplate entry uses task.TypeCleanCode constant as required by Hard Rules, not string literal

## Test Results
- **Tests Executed**: Yes
- **Passed**: 974
- **Failed**: 0
- **Coverage**: 83.4%

## Acceptance Criteria
- [x] forge-cli/pkg/prompt/data/code-quality-simplify.md deleted
- [x] forge-cli/pkg/prompt/data/code-quality-clean-code.md created, calls Skill(skill="forge:clean-code")
- [x] typeToTemplate has entry: task.TypeCleanCode: "data/code-quality-clean-code.md"
- [x] forge prompt get-by-task-id T-clean-code-1 returns valid prompt (no unknown type error)
- [x] scripts/version.txt bumped (patch)
- [x] go test -race -cover ./... passes from forge-cli/ directory

## Notes
TestSynthesize_CleanCodeTemplate_InvokesSkill added as dedicated test. TypeCleanCode also added to TestSynthesize_AllTypes table. All 974 tests pass with 83.4% total coverage.
