---
status: "completed"
started: "2026-05-23 03:26"
completed: "2026-05-23 03:38"
time_spent: "~12m"
---

# Task Record: 7 Move forge forensic group to forensic/ subdirectory

## Summary
Moved all forge forensic subcommand files from internal/cmd/ to internal/cmd/forensic/ subdirectory. Created forensic/register.go with a Register() function that adds all subcommands (search, extract, subagents) to the parent Cmd. Updated root.go to import and use the new forensic package. Updated root_test.go to reference forensicpkg.Cmd. Migrated testdata to the new package. Updated forge-cli-reference.md convention doc to reflect new source file locations.

## Changes

### Files Created
- forge-cli/internal/cmd/forensic/register.go
- forge-cli/internal/cmd/forensic/forensic.go
- forge-cli/internal/cmd/forensic/forensic_test.go
- forge-cli/internal/cmd/forensic/testdata/fix-bug.jsonl
- forge-cli/internal/cmd/forensic/testdata/history.jsonl
- forge-cli/internal/cmd/forensic/testdata/lesson-fix.jsonl
- forge-cli/internal/cmd/forensic/testdata/quick-mode.jsonl
- forge-cli/internal/cmd/forensic/testdata/subagent-eval.jsonl
- forge-cli/internal/cmd/forensic/testdata/subagents-session/subagents/agent-a1b2c3d4.jsonl
- forge-cli/internal/cmd/forensic/testdata/subagents-session/subagents/agent-a1b2c3d4.meta.json
- forge-cli/internal/cmd/forensic/testdata/subagents-session/subagents/agent-e5f6a7b8.jsonl
- forge-cli/internal/cmd/forensic/testdata/subagents-session/subagents/agent-e5f6a7b8.meta.json

### Files Modified
- forge-cli/internal/cmd/root.go
- forge-cli/internal/cmd/root_test.go
- docs/conventions/forge-cli-reference.md

### Key Decisions
- New package imports base directly (not cmd) to satisfy hard rule: no circular deps
- Variables renamed from forensicKeyword/session/skill/etc to keyword/session/skill since they're now in the forensic package namespace
- searchWithProjectPath signature simplified: removed unused histPath parameter, takes only projectPath and file handle
- captureStdout duplicated in forensic_test.go since it's a test helper in the cmd package that can't be imported across packages

## Test Results
- **Tests Executed**: Yes
- **Passed**: 60
- **Failed**: 0
- **Coverage**: 90.2%

## Acceptance Criteria
- [x] All forensic subcommand files are in internal/cmd/forensic/
- [x] New package exports a Register() function
- [x] root.go updated to use the new package
- [x] go build ./... passes
- [x] go test ./... passes
- [x] forge forensic behavior unchanged

## Notes
Hard rule verified: forensic package does NOT import internal/cmd (only imports internal/cmd/base).
