---
status: "completed"
started: "2026-05-15 00:51"
completed: "2026-05-15 00:57"
time_spent: "~6m"
---

# Task Record: 3 Update e2e/actions_test.go for just delegation

## Summary
Simplified test helpers after just delegation refactor. Replaced setupGoTestProfile with profile-agnostic setupProfile for Run/Setup/Compile/Discover tests (profile name no longer affects command dispatch). Retained setupProfileWithE2E only for Verify tests which need the e2e directory structure for file scanning. All existing just-delegation assertions were already in place from task 2.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/e2e/actions_test.go

### Key Decisions
- Kept setupProfileWithE2E as a separate helper (not inlined) because three Verify subtests need the e2e directory structure
- Used generic profile name 'go-test' in setupProfile calls to match project conventions while being profile-agnostic for command dispatch

## Test Results
- **Tests Executed**: Yes
- **Passed**: 27
- **Failed**: 0
- **Coverage**: 88.9%

## Acceptance Criteria
- [x] TestRun asserts just test-e2e (with and without feature arg)
- [x] TestSetup asserts just e2e-setup
- [x] TestCompile asserts just e2e-compile
- [x] TestDiscover asserts just e2e-discover
- [x] TestVerify unchanged (no just delegation)
- [x] Non-zero just exit produces non-nil error; zero exit produces nil error
- [x] just not on PATH produces clear error message
- [x] go test ./pkg/e2e/... passes with 80%+ coverage
- [x] No test references specific profile names for command dispatch

## Notes
The tests were already updated for just delegation in task 2. This task simplified the test helper structure per Hard Rules: setupGoTestProfile was replaced with setupProfile for profile-agnostic tests and setupProfileWithE2E for Verify tests that need directory structure.
