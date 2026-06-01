---
status: "completed"
started: "2026-05-31 15:09"
completed: "2026-05-31 15:09"
time_spent: ""
---

# Task Record: 1 Rename validate-index to validate and add AC count validation

## Summary
Renamed forge task validate-index to forge task validate and added AC count validation (1-6 range). File renamed validate_index.go → validate.go, internal symbols updated (validateCmd, runValidate, ExportRunValidate). Added validateACCount method and countACItems parser. 20 new tests covering AC count validation. All call sites in test files updated.

## Changes

### Files Created
- forge-cli/internal/cmd/task/validate.go
- forge-cli/internal/cmd/task/validate_test.go

### Files Modified
- forge-cli/internal/cmd/task/register.go
- forge-cli/internal/cmd/task/testbridge.go
- forge-cli/internal/cmd/root_test.go
- forge-cli/internal/cmd/integration_test.go
- forge-cli/internal/cmd/runners_test.go
- forge-cli/tests/task-type-system/task_types_dispatch_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- Direct rename (no alias) as breaking change is acceptable for internal CLI
- AC parsing uses ## Acceptance Criteria section + - [ ] prefix counting, matching template format
- File-level rename (validate_index.go → validate.go) for consistency with command name

## Test Results
- **Tests Executed**: Yes
- **Passed**: 173
- **Failed**: 0
- **Coverage**: 75.1%

## Acceptance Criteria
- [x] forge task validate 校验 index.json 结构 + 所有 task 文件的 AC 数量（向后兼容原有校验）
- [x] AC > 6 时返回 exit 1 + 错误信息（包含 task 文件名和 AC 数量）
- [x] AC = 0 时返回 exit 1 + 错误信息
- [x] forge task validate-index 不再存在（直接替换，breaking change）

## Notes
Quality gate initially blocked by stale golangci-lint cache referencing deleted worktree forge-cli-codebase-standards. Cache cleared, 0 lint issues now.
