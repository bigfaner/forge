---
id: "5"
title: "eval Skill Phase 0 Integration"
priority: "P0"
estimated_time: "2h"
dependencies: ["1", "2", "3", "4"]
type: "doc"
mainSession: false
---

# 5: eval Skill Phase 0 Integration

## Description

将 Phase 0 自由专家评审集成到 eval skill 的 proposal 类型处理流程中。这是最终集成任务，将 Task 1-4 创建的所有组件串联为完整的 Phase 0 流程。

需要修改 eval SKILL.md 和 eval-proposal command，添加 `--freeform-expert` 参数解析和 Phase 0 编排逻辑。不传参数时行为与现有 eval-proposal 完全一致。

## Reference Files
- `docs/proposals/eval-freeform-expert-review/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/SKILL.md` | 添加 `--freeform-expert` 参数到 Parameters 表；添加 Phase 0 架构到流程图；添加 Phase 0 编排步骤（专家推断 → 用户确认 → 自由评审 → 提取 → 注入 → 条件分支到 Step 2 或降级） |
| `plugins/forge/commands/eval-proposal.md` | 添加 `--freeform-expert` 到 argument-hint 和 Skill 调用参数 |
| `plugins/forge/skills/eval/rules/scorer-composition.md` | 更新 scorer prompt 组合逻辑，当 Phase 0 产出存在时追加注入的 findings |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria

- [ ] eval SKILL.md Parameters 表包含 `--freeform-expert` 参数（默认不启用）
- [ ] eval-proposal command 的 argument-hint 包含 `[--freeform-expert]`
- [ ] eval SKILL.md 架构流程图包含 Phase 0（在 rubric 循环之前）
- [ ] Phase 0 编排步骤完整：参数检测 → 专家复用检查 → 专家推断/确认 → 自由评审 → 提取 → JSON 校验 → 注入/降级
- [ ] 不传 `--freeform-expert` 时，eval 流程跳过 Phase 0 直接进入 Step 1（零回归）
- [ ] scorer-composition.md 更新：当 Phase 0 产出的 findings 文件存在时，将 findings 追加到 scorer prompt 末尾
- [ ] 错误场景覆盖：专家生成失败（降级）、自由评审产出为空（降级）、提取失败（降级）、部分提取失败（命中率告警）、注入无效果（报告标注）

## Hard Rules

- 不传 `--freeform-expert` 时行为必须与当前完全一致——任何行为差异视为回归
- Phase 0 的所有降级路径最终都回到标准 rubric 流程，不会中断 eval pipeline
- 遵守 `docs/conventions/forge-distribution.md` 的分发规范

## Implementation Notes

- 这是最终集成任务，需要理解 Task 1-4 创建的所有文件及其职责
- Phase 0 编排逻辑在 SKILL.md 中以结构化步骤描述（Claude 主 session 执行）
- 自由评审子 agent 通过 Agent tool 调用（model: sonnet）
- 提取子 agent 也通过 Agent tool 调用（model: sonnet）
