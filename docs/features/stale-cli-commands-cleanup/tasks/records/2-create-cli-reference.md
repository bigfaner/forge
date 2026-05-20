---
status: "completed"
started: "2026-05-20 17:11"
completed: "2026-05-20 17:14"
time_spent: "~3m"
---

# Task Record: 2 创建 forge CLI 命令参考文档

## Summary
创建 docs/conventions/forge-cli-reference.md，记录所有 43 个有效 forge CLI 命令（10 个顶层命令 + 33 个子命令），覆盖 8 个命令组（task、e2e、test、prompt、worktree、config、feature、forensic），每个命令标注用途和源文件路径。明确列出 4 个已移除的命令（detect、interfaces、framework、get）以防误用。

## Changes

### Files Created
- docs/conventions/forge-cli-reference.md

### Files Modified
无

### Key Decisions
- 不包含 cobra 框架自动生成的 completion 和 help 命令，严格按 root.go 注册清单记录
- 将 config 和 feature 的子命令单独列出（而非与顶层命令合并），保持与 root.go 中 AddCommand 注册结构一致

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] 文档覆盖所有有效 CLI 命令（约 46 个），包括命令组、子命令和简要说明
- [x] 每个命令标注其所在的源文件路径（便于维护）
- [x] 文档包含 frontmatter domains 字段用于 consolidate-specs 自动加载
- [x] 明确标注已移除的命令（detect、interfaces、framework、get），防止误用

## Notes
无
