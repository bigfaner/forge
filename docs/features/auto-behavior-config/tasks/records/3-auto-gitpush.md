---
status: "completed"
started: "2026-05-17 02:28"
completed: "2026-05-17 02:39"
time_spent: "~11m"
---

# Task Record: 3 Add auto.gitPush step to run-tasks skill

## Summary
Add auto.gitPush post-completion step to run-tasks skill. Added `getAutoKeyValue` accessor in profile package to support `forge config get auto.gitPush` dot-notation key. Updated run-tasks.md with Step 5 that conditionally runs git push after all-completed hook passes, with graceful error handling for missing remote, auth failures, and diverged branches.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/profile/config.go
- forge-cli/pkg/profile/config_test.go
- forge-cli/internal/cmd/config_test.go
- plugins/forge/commands/run-tasks.md

### Key Decisions
- Used separate getAutoKeyValue() function for dot-notation auto keys instead of adding to flat configKeyAccessors map, keeping concerns separated
- auto.gitPush returns 'false' as default (not ErrKeyNotFound) so run-tasks can cleanly check the value without needing to handle missing-key errors
- Push uses `git push $REMOTE HEAD` (not force push) to push only current branch, matching Hard Rule of never force pushing
- Error messages include actionable hints (auth error, no remote write access, diverged branches) per error-reporting business rule

## Test Results
- **Tests Executed**: Yes
- **Passed**: 6
- **Failed**: 0
- **Coverage**: 85.7%

## Acceptance Criteria
- [x] auto.gitPush=true triggers git push after all-completed hook succeeds
- [x] auto.gitPush=false or absent causes no push (backward compatible)
- [x] Push failures produce clear error message (auth, no remote) with no crash
- [x] Pushes to configured remote (default: origin)

## Notes
The Go accessor (getAutoKeyValue) currently handles only 'auto.gitPush' as that is the only auto dot-notation key needed by the CLI. The run-tasks.md Step 5 is purely instructional (agent reads config via forge CLI and conditionally pushes), so no Go integration tests for the push logic itself.
