---
status: "completed"
started: "2026-06-09 22:28"
completed: "2026-06-09 22:36"
time_spent: "~8m"
---

# Task Record: 2 api/web/mobile service surface 模板

## Summary
Implemented api/web/mobile service surface templates with full lifecycle management (dev, probe, test, teardown, orchestration) and quality recipes. Migrated PID file management, idempotent start, health check retry, and dual-platform teardown logic from server-lifecycle.md bash templates into Go scaffold generation.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/scaffold/types.go
- forge-cli/internal/cmd/scaffold/generate.go
- forge-cli/internal/cmd/scaffold/generate_test.go
- forge-cli/internal/cmd/scaffold/register.go

### Key Decisions
- Added IsScript field to RecipeSpec to distinguish multi-line bash recipes from one-liners, keeping writeRecipe as single rendering path
- serviceRecipes(false) shared by api/web, serviceRecipes(true) for mobile with extra test-setup
- Orchestration recipe generated separately in writeOrchestrationRecipe using orchestrationSteps() per surface type
- Recipe body templates stored as raw string constants in types.go for readability

## Test Results
- **Tests Executed**: Yes
- **Passed**: 48
- **Failed**: 0
- **Coverage**: 87.3%

## Acceptance Criteria
- [x] AC-1: api surface outputs backend-dev + backend-probe + backend-test + backend-teardown + backend (orchestration) + quality recipes, all with backend- prefix
- [x] AC-2: web same as api; mobile extra test-setup recipe, orchestration is test-setup->dev->probe->test->teardown
- [x] AC-3: All lifecycle recipes contain PID file management (.forge/<<URL_KEY>>.pid) and health check retry logic with <<PLACEHOLDER>> syntax
- [x] AC-4: All lifecycle and quality recipes marked # user-customized with [unix] + [windows] dual-platform variants

## Notes
Removed 'not yet supported' check from ValidateArgs since api/web/mobile are now fully supported. Updated register.go --type help text to list all 5 surface types.
