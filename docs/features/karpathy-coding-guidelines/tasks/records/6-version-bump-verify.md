---
status: "completed"
started: "2026-05-20 00:55"
completed: "2026-05-20 01:01"
time_spent: "~6m"
---

# Task Record: 6 Bump CLI version and verify prompt generation includes principles

## Summary
Bumped CLI patch version from 4.4.1 to 4.4.2 and added test to verify all 5 coding templates contain CODING_PRINCIPLES text block in synthesized prompts

## Changes

### Files Created
无

### Files Modified
- forge-cli/scripts/version.txt
- forge-cli/pkg/prompt/prompt_test.go

### Key Decisions
- Patch bump per task AC (4.4.1 -> 4.4.2) rather than minor bump, as task explicitly requests patch
- Added table-driven subtest for all 5 coding types (feature/enhancement/fix/refactor/cleanup) to verify CODING_PRINCIPLES presence in synthesized output

## Test Results
- **Tests Executed**: Yes
- **Passed**: 24
- **Failed**: 0
- **Coverage**: 90.6%

## Acceptance Criteria
- [x] CLI version patch-bumped
- [x] Forge CLI compiles: just compile backend
- [x] For each coding type, forge prompt output contains CODING_PRINCIPLES text block
- [x] Existing prompt-related tests pass: go test -race -cover ./forge-cli/pkg/prompt/...
- [x] No template content modified (only version bump + test addition)

## Notes
Pre-existing failures in internal/docsync (extract_design_md tests) confirmed unrelated to this task via git stash verification. Version mechanism uses forge-cli/scripts/version.txt read by install scripts via ldflags injection.
