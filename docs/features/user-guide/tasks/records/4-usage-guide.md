---
status: "completed"
started: "2026-05-30 21:06"
completed: "2026-05-30 21:09"
time_spent: "~3m"
---

# Task Record: 4 编写使用指南 usage-guide.md

## Summary
创建 docs/user-guide/usage-guide.md 使用指南，包含 Full Mode 和 Quick Mode 端到端实战示例、5 个单命令场景、7 条常见问题排错指引

## Changes

### Files Created
- docs/user-guide/usage-guide.md

### Files Modified
无

### Key Decisions
无

## Document Metrics
~260 lines, 2 e2e examples, 5 single-command scenarios, 7 troubleshooting items

## Referenced Documents
- README.md
- docs/ARCHITECTURE.md
- docs/business-rules/task-lifecycle.md
- docs/business-rules/quality-gate.md

## Review Status
final

## Acceptance Criteria
- [x] 包含 Full Mode 至少一个端到端实战示例（从 brainstorm 到任务执行完成）
- [x] 包含 Quick Mode 至少一个端到端实战示例（从 /quick 到任务执行完成）
- [x] 包含至少 2 个单命令场景示例（如 /learn、/consolidate-specs）
- [x] 包含 5 条以上常见问题及排错指引（涵盖安装失败、配置错误、工作流异常、任务阻塞、测试失败）
- [x] 所有代码示例可直接复制执行，无需额外修改

## Notes
Full Mode 示例使用用户通知系统场景，Quick Mode 示例使用登录超时 bug 修复场景，排错指引覆盖安装、配置、工作流、任务状态、测试、执行、worktree 共 7 类问题
