---
status: "completed"
started: "2026-06-09 22:51"
completed: "2026-06-09 23:09"
time_spent: "~18m"
---

# Task Record: 6 更新 quality gate 下游 consumer

## Summary
Removed fallback logic from ResolvePrefixedRecipe and resolveRecipe. When scope/surfaceType is set, functions now only return prefixed recipe names (<scope>-<recipe>) or empty string — no fallback to generic recipe. Scalar surfaces (empty scope) still use unprefixed recipes unchanged. Updated 4 tests in just_test.go and 5 tests in quality_gate_test.go to match new behavior.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/just/just.go
- forge-cli/pkg/just/just_test.go
- forge-cli/internal/cmd/qualitygate/quality_gate_lifecycle.go
- forge-cli/internal/cmd/qualitygate/quality_gate_test.go

### Key Decisions
- Both ResolvePrefixedRecipe and resolveRecipe share identical no-fallback semantics: prefixed only or empty
- No SKILL.md changes needed — recipe fallback was entirely in Go code, not in skill definition text
- Updated test justfiles to use prefixed recipe names (web-dev, web-probe, cli-test, etc.) to match the new resolution behavior

## Test Results
- **Tests Executed**: Yes
- **Passed**: 134
- **Failed**: 0
- **Coverage**: 84.7%

## Acceptance Criteria
- [x] ResolvePrefixedRecipe removes generic fallback: scoped calls only return <scope>-<recipe>
- [x] resolveRecipe removes generic fallback, consistent with ResolvePrefixedRecipe
- [x] run-tests SKILL.md has no fallback chain descriptions
- [x] Single surface scalar projects (no scope) still use unprefixed recipes

## Notes
Tests passed: just package 84.7% coverage, qualitygate 73.9% coverage. No SKILL.md changes were needed because fallback descriptions were not present in the skill definition — the fallback logic was purely in Go code.
