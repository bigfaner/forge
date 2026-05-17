---
status: "completed"
started: "2026-05-17 16:52"
completed: "2026-05-17 17:09"
time_spent: "~17m"
---

# Task Record: 1 Implement forge feature complete --if-done CLI command

## Summary
Implement forge feature complete --if-done CLI command: a Stop hook that checks index.json for all completed/skipped tasks, updates manifest.md (and proposal.md in quick mode), commits via explicit git add paths, and optionally pushes to remote.

## Changes

### Files Created
- forge-cli/internal/cmd/feature_complete.go
- forge-cli/internal/cmd/feature_complete_test.go

### Files Modified
- plugins/forge/hooks/hooks.json
- forge-cli/scripts/version.txt

### Key Decisions
- Direct function testing (checkFeatureCompletion + completeFeature) instead of subprocess-based integration tests — more reliable, same coverage
- updateFileStatus takes only filePath and value (not key) since it always updates 'status' field — avoids unparam lint warning
- gitPush returns error for caller to decide blocking behavior — completeFeature treats push failure as non-blocking per hook protocol

## Test Results
- **Tests Executed**: Yes
- **Passed**: 19
- **Failed**: 0
- **Coverage**: 83.2%

## Acceptance Criteria
- [x] forge feature complete --if-done exits 0 silently when no feature is active
- [x] forge feature complete --if-done exits 0 silently when index.json has pending/in_progress tasks
- [x] When all tasks are completed/skipped, command updates manifest.md status to completed
- [x] When proposal.md exists (quick mode), command also updates its status to Completed
- [x] When proposal.md is absent (full pipeline), only manifest.md is updated
- [x] Git commit targets only manifest.md and proposal.md by explicit path
- [x] When auto.gitPush is true in config, command pushes to remote after commit
- [x] When auto.gitPush is false or absent, no push occurs
- [x] Push failure is logged to stderr but does not block (exit 0)
- [x] hooks.json Stop array contains forge quality-gate first, then forge feature complete --if-done second
- [x] Go tests pass with >= 80% coverage on new code

## Notes
runFeatureCompleteCmd (cobra entry point) has 0% coverage because it calls os.Exit — testing requires fragile subprocess tests. All actual logic functions have 75-91% coverage.
