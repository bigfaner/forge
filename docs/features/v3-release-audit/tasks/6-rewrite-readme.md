---
id: "6"
title: "Full README.md rewrite for v3.0.0"
priority: "P0"
estimated_time: "2h"
dependencies: ["3", "4", "5"]
type: "doc"
mainSession: false
---

# 6: Full README.md rewrite for v3.0.0

## Description
README.md 版本号仍为 2.16.1（实际 3.0.0-rc.25），技能计数、任务类型表、Pipeline ID、命令名全部过时，含幽灵命令和 web 引用。需要全面重写以反映 v3.0.0 状态。此任务必须最后执行，依赖前面任务提供的准确计数。

## Reference Files
- `proposal.md#Problem` — Evidence table: README 7 Critical + 6 Major + 1 Minor errors
- `proposal.md#Proposed-Solution` — Target State: defines exact README structure requirements
- `proposal.md#Scope` — P0.1 defines full rewrite scope
- `proposal.md#Success-Criteria` — "README 事实性声明 100% 一致" P0 gate

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `README.md` | Full rewrite: version, skill count, task types, commands, paths, structure |

## Acceptance Criteria
- [ ] 版本号 = `cat plugins/forge/scripts/version.txt` 或当前 RC 版本
- [ ] 技能计数 = `ls plugins/forge/skills/ | wc -l`
- [ ] 任务类型表覆盖所有 dot-notation 类型
- [ ] 命令速查与 `forge --help` 一一对应
- [ ] 无幽灵命令（如已移除的 web/raycast 引用）
- [ ] 安装步骤指向正确 Go 版本要求
- [ ] 路径引用与实际目录匹配

## Hard Rules
- 逐条交叉验证每个事实性声明
- 不改未审计部分的结构
- 计数和命令名必须在任务 4、5 完成后从代码库实时获取

## Implementation Notes
README 是用户和 agent 的第一入口。重写后需与 Task 4（ARCHITECTURE.md）和 Task 5（CLI refs）的结果保持一致。建议先 `forge --help`、`ls plugins/forge/skills/`、`ls plugins/forge/agents/` 等获取实时计数再编写。
