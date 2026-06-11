---
status: "completed"
started: "2026-05-25 00:21"
completed: "2026-05-25 00:25"
time_spent: "~4m"
---

# Task Record: 7 Complete CLI flags documentation (8 missing)

## Summary
补全 README.md CLI flags 文档：新增 Flags 参考章节，覆盖 19 个子命令/子命令组的全部业务 flags（task add 12个、task list/query/submit/index/transition、feature/feature complete、worktree start/remove、init、quality-gate、prompt get-by-task-id、surfaces/surfaces detect、config、fact list、forensic search/extract），并补充 4 个遗漏的 task 子命令（check-deps、reopen、transition、validate-index）到常用 task 子命令表

## Changes

### Files Created
无

### Files Modified
- README.md

### Key Decisions
无

## Document Metrics
1 file modified, +150 lines flags reference covering 19 command groups, 4 subcommands added to task table

## Referenced Documents
- docs/proposals/v3-release-audit/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] 所有 forge <subcommand> --help 输出的 flags 在文档中有对应描述
- [x] 无文档中存在但 --help 不输出的幽灵 flags

## Notes
逐一运行 forge <subcommand> --help 收集实际 flags，与 README.md 原有内容比对。docs/ARCHITECTURE.md 无 flag 描述章节，无需修改。forge-cli/docs/OVERVIEW.md 存在幽灵命令和 flags（如 forge test detect/get/interfaces/framework、--force on task status），但不在本任务 Affected Files 范围内。
