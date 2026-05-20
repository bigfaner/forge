---
id: "2"
title: "Add config check logic to quick.md"
priority: "P0"
estimated_time: "15m"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 2: Add config check logic to quick.md

## Description

修改 `quick.md` 命令文件，添加 `auto.runTasks` 配置检查逻辑。当 `quick: true`（默认）时跳过 Step 2 确认门，直接进入任务生成+执行；当 `quick: false` 时保留当前行为（暂停展示摘要等待用户确认）。

## Reference Files
- `docs/proposals/auto-execute-tasks/proposal.md` — Source proposal
- `plugins/forge/skills/quick/SKILL.md` — quick 命令 skill 文件

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/quick/SKILL.md` | Step 2 添加配置检查，条件性跳过确认门 |

## Acceptance Criteria

- [ ] `auto.runTasks.quick: true`（默认）时，`/quick` 在 brainstorm 确认后直接生成+执行任务，无 Step 2 暂停
- [ ] `auto.runTasks.quick: false` 时，保留 Step 2 确认门（当前行为）
- [ ] 使用 `forge config get auto.runTasks` 读取配置值
- [ ] 配置检查逻辑用 `EXTREMELY-IMPORTANT` 标注，确保 AI agent 正确执行

## Hard Rules

- 不能删除 Step 2 的确认门逻辑，只能条件性跳过
- 必须在 Step 2 开头进行配置检查，不是在 brainstorm 阶段

## Implementation Notes

- 配置读取失败时应 fallback 到默认行为（跳过确认），符合 quick 模式精简定位
- 参考 `quick.md` 中其他配置检查的模式（如有）
