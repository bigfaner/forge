---
status: "completed"
started: "2026-06-09 14:16"
completed: "2026-06-09 14:28"
time_spent: "~12m"
---

# Task Record: 1 Rename `copy-files` → `includes` across codebase

## Summary
Renamed copy-files to includes across the entire codebase: struct field CopyFiles->Includes, YAML tag copy-files->includes, functions validateCopyFiles->validateIncludes, copyFilesToWorktree->copyIncludesToWorktree, validateCopyFilePath->validateIncludePath, and all test references updated

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/forgeconfig/config.go
- forge-cli/internal/cmd/worktree/helpers.go
- forge-cli/internal/cmd/worktree/cmd_start.go
- forge-cli/internal/cmd/init_config.go
- forge-cli/internal/cmd/testdata/forge-config.schema.json
- forge-cli/internal/cmd/testdata/forge-config.example.yaml
- forge-cli/pkg/forgeconfig/config_test.go
- forge-cli/internal/cmd/init_test.go
- forge-cli/internal/cmd/worktree/worktree_test.go

### Key Decisions
- Direct replacement without backward compatibility for old copy-files field name
- Renamed all function names including validateCopyFilePath->validateIncludePath for consistency
- Updated all user-facing error messages and hints to reference includes

## Test Results
- **Tests Executed**: Yes
- **Passed**: 62
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] WorktreeConfig.CopyFiles renamed to WorktreeConfig.Includes, YAML tag changed to includes
- [x] grep -r 'copy-files|CopyFiles|copy_files' forge-cli/ --include='*.go' returns zero results
- [x] JSON schema and example YAML updated to use includes
- [x] All existing tests pass after rename

## Notes
Pure mechanical rename across 9 files. go build, go vet, and all targeted tests pass. AC-2 grep verification confirmed zero results.
