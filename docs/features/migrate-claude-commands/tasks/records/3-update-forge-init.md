---
status: "completed"
started: "2026-05-16 23:39"
completed: "2026-05-16 23:51"
time_spent: "~12m"
---

# Task Record: 3 Update forge init to stop generating claude justfile recipes

## Summary
Removed all claude-related justfile recipe generation from forge init. Removed justfileRecipes var, updateJustfile/buildJustfileAppend/recipeExists/collectAddedRecipeNames functions, and the Step 4 justfile update block from runInit. Updated command description and all tests.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/init.go
- forge-cli/internal/cmd/init_test.go

### Key Decisions
- Removed the entire justfile update step (Step 4) since justfileRecipes became empty after removing claude/claude-c recipes
- Kept ensureJustStep (Step 3.5, renumbered to Step 4) since it checks just installation independently of recipe generation
- Replaced two justfile assertion tests with a single 'does not create justfile' test

## Test Results
- **Tests Executed**: Yes
- **Passed**: 20
- **Failed**: 0
- **Coverage**: 81.1%

## Acceptance Criteria
- [x] justfileRecipes slice in init.go is empty or removed
- [x] forge init no longer appends claude-related recipes to justfile
- [x] Init command description updated to remove claude/claude-c recipe mention
- [x] Existing tests for forge init updated
- [x] buildJustfileAppend and related helper functions cleaned up

## Notes
Also removed TestBuildJustfileAppend test suite since the tested functions no longer exist.
