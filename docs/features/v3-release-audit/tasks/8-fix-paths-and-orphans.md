---
id: "8"
title: "Fix cross-skill path violations and orphan rules"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["6"]
type: "doc"
mainSession: false
---

# 8: Fix cross-skill path violations and orphan rules

## Description
两个相关问题：(1) run-tests 中硬编码路径违规，需改为描述性引用；(2) 15 个 rules 文件未被父 SKILL.md 引用（6 个真孤儿需添加 Load，5 个参数化 surface rules 标注引用，4 处 `forge test run --tags` 降级归入此项）。

## Reference Files
- `proposal.md#Scope` — P1.7: cross-skill path violations; P1.8: orphan rules fix
- `proposal.md#Problem` — Evidence table: Skill-CLI cross-ref 2 Major, architecture health findings

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/run-tests/SKILL.md` | Replace hardcoded paths with descriptive references |
| Orphan rules' parent SKILL.md files | Add Load directives for 6 truly orphaned rules |

## Acceptance Criteria
- [ ] `grep -r "hardcoded/path/pattern" plugins/forge/skills/run-tests/` 返回 0（具体模式待确认）
- [ ] 所有 rules/ 文件被至少一个 SKILL.md 引用（入度 ≥ 1）
- [ ] 6 个真孤儿 rules 已添加 Load 指令
- [ ] 5 个参数化 surface rules 已标注引用关系

## Hard Rules
- 修改 run-tests 路径时需遵守 forge-distribution.md
- 添加 Load 前确认 rules 文件确实被需要（不是死规则）

## Implementation Notes
需先列出所有 rules/ 文件，逐一 grep 确认入度（被引用次数），区分真孤儿和参数化引用。
