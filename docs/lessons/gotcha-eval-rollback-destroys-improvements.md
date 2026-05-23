---
created: "2026-05-23"
tags: [architecture, testing]
---

# eval 回滚机制会销毁有价值的修订

## Problem

eval 流水线的 Step 5 协议规定：当最终评分未达目标且所有迭代用尽时，回滚到原始文档。但在实际执行中，经过 3 轮修订的 proposal 从 597 分提升到 795 分（+33%），回滚后丢失了所有改进。

## Root Cause

1. **协议假设**：回滚的设计意图是防止"部分修订"留下比原始更差的文档
2. **实际场景**：reviser 子代理的每一轮改进都是正向的（597→725→795），不存在"修订后退"的情况
3. **机制缺陷**：协议没有区分"修订方向正确但未达标"和"修订使文档变差"两种情况，一律回滚
4. **信息丢失**：proposal 未提前 commit，回滚时 backup 被删除，修订版无法恢复

## Solution

修改 eval 协议的回滚策略：**不自动回滚，改为由用户决定**。

具体做法：当所有迭代用尽且未达标时，不执行 `rm -rf + cp -r backup`，而是：
1. 保留修订版文档
2. 在最终报告中展示分数演变（原始 → 最终）
3. 由用户判断是否保留修订版或手动恢复原始版本

如果需要回滚能力，保留 backup 但不自动执行恢复。

## Reusable Pattern

当 eval 结果是"未达标但有显著改善"时，应展示改善轨迹让用户决定，而非机械回滚。协议应区分方向性失败（修订使文档变差）和距离性失败（修订方向正确但距离目标仍有差距）。

## Related Files

- `plugins/forge/skills/eval/SKILL.md` — Step 5 Rollback on Failure
- `docs/proposals/auto-gen-journeys-contracts/` — 本次受影响的 proposal
