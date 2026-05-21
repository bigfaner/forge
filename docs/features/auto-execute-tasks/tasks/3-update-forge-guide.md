---
id: "3"
title: "Update forge guide for runTasks"
priority: "P1"
estimated_time: "10m"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 3: Update forge guide for runTasks

## Description

更新 forge guide（`plugins/forge/hooks/guide.md`）的 Automation Config 部分，添加 `runTasks` 配置项说明。

## Reference Files
- `docs/proposals/auto-execute-tasks/proposal.md` — Source proposal
- `plugins/forge/hooks/guide.md` — Forge guide 文档

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/hooks/guide.md` | Automation Config 部分添加 runTasks 说明 |

## Acceptance Criteria

- [ ] Automation Config 部分列出 `runTasks` 配置项及其 `quick`/`full` 子键
- [ ] 说明默认值：`quick: true, full: false`
- [ ] 简述作用：控制 `/quick` 流水线是否跳过确认门自动执行

## Hard Rules

- 保持现有 Automation Config 段落的简洁风格

## Implementation Notes

- 当前 Automation Config 段落位于 `guide.md` 第 80 行附近
- 现有格式：列出配置项名，说明 quick/full 子键和默认值
