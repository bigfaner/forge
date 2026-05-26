---
id: "7"
title: "Update gen-journeys SKILL.md for multi-surface support"
priority: "P1"
estimated_time: "1h"
dependencies: ["2"]
type: "doc"
mainSession: false
---

# 7: Update gen-journeys SKILL.md for multi-surface support

## Description
更新 `plugins/forge/skills/gen-journeys/SKILL.md`，适配 gen-journeys 合并为单任务后的多 surface 遍历。新增多 surface 规则加载指导，确保输出 Journey 文件覆盖所有配置 surface type，每个 Journey 标注覆盖的 surface type 集合。

## Reference Files
- `proposal.md#Key-Risks` — gen-journeys SKILL.md 需适配多 surface 内部遍历的风险
- `proposal.md#Success-Criteria` — SC12: gen-journeys 输出 Journey 文件中每个 Journey 标注 surface type 集合

## Acceptance Criteria
- [ ] SKILL.md 包含多 surface 规则加载指导（按 surface type 分节组织）
- [ ] 输出格式要求：每个 Journey 标注覆盖的 surface type 集合（如 `[web, api]`）
- [ ] 所有配置的 surface type 至少被一个 Journey 覆盖

## Hard Rules
- 修改 SKILL.md 前必须加载 `docs/conventions/forge-distribution.md` 了解分发模型约束

## Implementation Notes
- SKILL.md 已有 surface 检测逻辑，需扩展为多 surface 遍历
- gen-journeys 以 PRD 为主要输入，surface 规则作为参考指导
