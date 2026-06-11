---
id: "15"
title: "Add v3.0.0 new features section to README"
priority: "P2"
estimated_time: "30m"
dependencies: ["6"]
type: "doc"
mainSession: false
---

# 15: Add v3.0.0 new features section to README

## Description
README 重写后需补充 v3.0.0 新特性亮点：CLI 命令分组（forge-cli-v3 特性）、surface 自动检测、test profile 可插拔、worktree 支持等。帮助现有用户快速了解 v3 变化。

## Reference Files
- `proposal.md#Scope` — P2.21: README v3.0.0 new features
- `docs/proposals/forge-cli-v3/proposal.md` — Reference for CLI restructure changes

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `README.md` | Add "What's New in v3.0.0" section |

## Acceptance Criteria
- [ ] README 含 v3.0.0 新特性段落
- [ ] 列举的主要特性与 forge-cli-v3 proposal 和 v3 实际变更一致
- [ ] 段落长度适中（~20-30 行），不喧宾夺主

## Hard Rules
- 仅在 Task 6 完成后执行
- 特性列表需与实际实现匹配，不包含未发布功能

## Implementation Notes
可参考 forge-cli-v3 proposal 和 todo.txt 中标记 done 的 v3 相关条目确定特性列表。
