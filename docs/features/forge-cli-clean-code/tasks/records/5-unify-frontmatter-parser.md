---
status: "completed"
started: "2026-05-24 01:58"
completed: "2026-05-24 02:04"
time_spent: "~6m"
---

# Task Record: 5 Extend ParseFrontmatter as shared YAML parser

## Summary
Unified frontmatter parsing by extracting shared core logic into infocmd.ExtractFrontmatter, making both task.ParseFrontmatter and infocmd.ParseFrontmatter delegates of the single implementation. Removed dead code (cutLine helper, unused trParseFrontmatter in e2e tests).

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/infocmd/infocmd.go
- forge-cli/pkg/task/frontmatter.go
- forge-cli/tests/task-type-system/task_type_refinement_test.go

### Key Decisions
- Placed ExtractFrontmatter in infocmd package (already the shared utility package used by lesson/research/proposal) rather than task package, avoiding a new reverse dependency
- Both ParseFrontmatter signatures preserved unchanged for backward compatibility (Hard Rule 1)
- ExtractFrontmatter returns (rawYAML, body, error) enabling callers to unmarshal into their own structs (Hard Rule 2)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 835
- **Failed**: 0
- **Coverage**: 92.3%

## Acceptance Criteria
- [x] ParseFrontmatter() returns two-layer API via ExtractFrontmatter
- [x] All duplicate frontmatter parsing sites replaced with shared parser
- [x] Zero duplicate frontmatter parsing implementations remain
- [x] go build ./... passes
- [x] go test ./... passes

## Notes
Coverage per package: infocmd 81.0%, task 92.3%, lesson 94.4%, research 100.0%, proposal 96.0%. The trParseFrontmatter in task_type_refinement_test.go was dead code (defined but never called) and removed.
