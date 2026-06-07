---
status: "completed"
started: "2026-06-07 23:38"
completed: "2026-06-07 23:48"
time_spent: "~10m"
---

# Task Record: 1 RunGate() 增加 prefixed recipe 解析

## Summary
Added ResolvePrefixedRecipe() to just.go and refactored RunGate() to use prefixed recipe resolution instead of ResolveScope(). When scope (surface-key) is non-empty, RunGate tries <scope>-<recipe> first (e.g., backend-compile), falling back to generic recipe. Empty scope skips prefixed branch entirely, preserving feature-level gate behavior.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/just/just.go
- forge-cli/pkg/just/just_test.go

### Key Decisions
- ResolvePrefixedRecipe replaces ResolveScope in RunGate — prefixed recipe pattern (<key>-<recipe>) supersedes argument-passing pattern (recipe <scope>)
- onFail callback receives resolved recipe name (prefixed or generic) instead of step.Name, enabling callers to distinguish surface context
- Removed unused hasRecipeWithArg after refactoring — ResolveScope kept for backward compat but no longer called from critical path
- Updated legacy test from argument-mode (compile frontend:) to prefix-mode (frontend-compile:) to match new resolution strategy

## Test Results
- **Tests Executed**: Yes
- **Passed**: 21
- **Failed**: 0
- **Coverage**: 84.9%

## Acceptance Criteria
- [x] mixed.just project with surface-key:backend executes backend-compile/backend-lint, falls back to generic when prefixed absent
- [x] single surface project with empty scope skips prefixed branch, recipe names remain compile/lint
- [x] RunGate scope="" call path (feature-level gate) behavior unchanged
- [x] Prefixed failure onFail step name is <key>-<recipe>; generic fallback failure step name is original recipe name
- [x] All tests pass in ./pkg/just/... and ./internal/cmd/task/...

## Notes
Coverage measured at 84.9% for pkg/just package. ResolveScope and TestResolveScope retained in codebase for backward compatibility but no longer used by RunGate.
