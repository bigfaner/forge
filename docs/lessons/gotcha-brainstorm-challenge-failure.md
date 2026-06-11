---
created: "2026-05-20"
tags: [architecture, error-handling]
---

# Brainstorm 中未挑战伪需求：模糊搜索案例

## Problem

在 brainstorm 会话中，用户提出"支持模糊搜索"需求。Agent 花了 4 轮问答讨论搜索范围（slug/title/tags）、匹配算法（子串/fzf）、适用命令——但从未质疑模糊搜索本身是否必要。最终用户自己意识到 `forge proposal | grep <keyword>` 就够了，主动砍掉了这个需求。

## Root Cause

1. **表面原因**: 把用户需求当圣旨，执行"需求澄清"而非"需求质疑"。问的是"怎么做"（how），不是"为什么做"（why）。
2. **深层原因**: Brainstorm SKILL.md 明确要求 "Challenge is mandatory at every decision point"，并给出了工具箱（Occam's Razor、XY Detection、5 Whys）。Agent 读了这些规则但没有在决策点激活它们——知识和行为脱节。
3. **根本原因**: Agent 陷入"gather requirements → implement"的思维惯性，把 brainstorm 当成需求访谈而非协作设计。缺少一个强制检查点：在进入任何 how 讨论前，先用最简方案（Occam's Razor）检验需求是否成立。

## Solution

在 brainstorm 的每个决策点，强制执行"最简方案优先"检查：

1. 用户提出功能需求 → 先问"管道/现有工具能不能解决？"
2. 只有当现有工具明确不够时，才进入 how 讨论
3. 把 Occam's Razor 从"贯穿始终的哲学"升级为每个决策点的硬性检查门

具体到模糊搜索案例，正确的第一反应应该是：
> "`forge proposal | grep timeout` 不就能过滤吗？为什么需要在 CLI 内部实现搜索？"

## Reusable Pattern

**Brainstorm 决策点的强制检查顺序**:
1. Occam's Razor → 现有工具/管道能否满足？
2. XY Detection → 用户要的 X 是否真的是解决 Y 的最优路径？
3. 只有前两关通过后，才进入功能细节讨论（how）

**通用反模式**: 花时间讨论"如何实现"一个根本不需要实现的功能。

## Related Files

- plugins/forge/skills/brainstorm/SKILL.md (challenge protocol 定义)
- plugins/forge/skills/brainstorm/rules/challenge-protocol.md
