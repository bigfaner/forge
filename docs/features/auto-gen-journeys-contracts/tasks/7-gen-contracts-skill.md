---
id: "7"
title: "gen-contracts SKILL.md 适配：SKIP_EVAL_GATE 路径"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 7: gen-contracts SKILL.md 适配：SKIP_EVAL_GATE 路径

## Description

修改 gen-contracts SKILL.md 的 Prerequisites 部分，添加 SKIP_EVAL_GATE 条件路径：当任务上下文包含 SKIP_EVAL_GATE=true 时，跳过 eval-journey 报告的前置检查。

## Reference Files
- `docs/proposals/auto-gen-journeys-contracts/proposal.md` — Source proposal
- `plugins/forge/skills/gen-contracts/SKILL.md` — 当前 SKILL.md

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-contracts/SKILL.md` | Prerequisites 新增 SKIP_EVAL_GATE 条件路径 |

## Acceptance Criteria

- [ ] Prerequisites 中 eval-journey 报告（`.eval-report.md`）前置条件新增条件豁免：若任务上下文包含 SKIP_EVAL_GATE=true 指令，跳过此 Blocker 检查
- [ ] SKIP_EVAL_GATE 路径下，gen-contracts 直接进入 Step 1（Read Journeys）和 Step 2（Code Reconnaissance）
- [ ] 非 SKIP_EVAL_GATE 模式（Breakdown 模式 / 手动调用 `/gen-contracts`）时行为不变：仍要求 eval 报告
- [ ] SKIP_EVAL_GATE 路径下生成的 Contract 文件应有标记表明未经 eval 验证（如注释或 metadata）

## Hard Rules

- 不修改核心 Contract 生成逻辑（Step 3 的六维度生成、Step 4 的验证）
- 仅修改 Prerequisites 的前置检查逻辑

## Implementation Notes

- 当前 gen-contracts Prerequisites (L38-41) 要求 eval report for all Journeys，标记为 "Blocker: do not proceed if any Journey scored below target"
- SKIP_EVAL_GATE 路径是 Quick 模式的必要妥协：Quick 模式跳过 eval-journey 阶段，因此不可能有 eval 报告
