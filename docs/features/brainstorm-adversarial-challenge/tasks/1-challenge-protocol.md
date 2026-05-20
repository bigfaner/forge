---
id: "1"
title: "Embed adversarial challenge tools into brainstorm skill"
priority: "P1"
estimated_time: "2h"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Embed adversarial challenge tools into brainstorm skill

## Description

Brainstorm skill 的 Challenge 集群是三个并列决策集群之一，AI 可自由跳过。将挑战机制从"可选集群"升级为"融入每个决策点的强制行为"：重写 SKILL.md 的 Decision Clusters，新增 Challenge Protocol 章节，更新 proposal 模板新增 Assumptions Challenged 章节。

## Reference Files
- `docs/proposals/brainstorm-adversarial-challenge/proposal.md` — Source proposal
- `docs/conventions/forge-distribution.md` — Forge plugin distribution model (MUST read before modifying plugin files)

## Affected Files

### Create
(none)

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/brainstorm/SKILL.md` | Rewrite Decision Clusters to embed challenge tools per cluster; add Challenge Protocol section; add fact-driven principle; add Occam's razor meta-principle |
| `plugins/forge/skills/brainstorm/templates/proposal.md` | Add Assumptions Challenged section between Feasibility Assessment and Scope |

### Delete
(none)

## Acceptance Criteria

- [ ] 每个 Decision Cluster（Problem、Solution）都有明确绑定的挑战工具，不是可选的
- [ ] Challenge Protocol 章节定义 5 个挑战工具（5 Whys、XY 检测、假设翻转、极端推演、奥卡姆剃刀）的使用时机、触发条件和终止条件
- [ ] 挑战工具要求附带事实依据（代码库事实/逻辑一致性/领域常识），不允许空泛质疑
- [ ] 冷启动项目（无代码库）的挑战同样有效——证据来源泛化为三类
- [ ] proposal.md 模板包含 Assumptions Challenged 章节（表格格式：Assumption / Challenge Tool / Finding）
- [ ] 奥卡姆剃刀作为贯穿全程的元原则明确写入
- [ ] 7 步流程结构不变（Step 1-7 数量和名称不变），只改动 Step 2 内部内容
- [ ] 不引入新的流程步骤

## Hard Rules

- **MUST read `docs/conventions/forge-distribution.md` before modifying any plugin file** — understand path resolution constraints
- Skill internal file references must follow forge-distribution path rules
- `${CLAUDE_SKILL_DIR}` does NOT apply in template files — only in SKILL.md

## Implementation Notes

- 当前 SKILL.md 的 Decision Clusters 是一个简单的三行表格，Challenge 作为独立集群。新设计删除 Challenge 集群，将其工具分配到 Problem 和 Solution 集群
- Problem 集群绑定：5 Whys（根因挖掘）+ XY 问题检测（识别表面需求）
- Solution 集群绑定：假设翻转（验证关键假设）+ 极端推演（暴露隐藏风险）
- 奥卡姆剃刀作为元原则贯穿所有集群——当多个解释/方案并存时，优先选择最简方案
- 事实驱动泛化：证据来源 = 代码库事实 + 逻辑一致性 + 领域常识。在冷启动项目中，后两者是主要证据来源
- Challenge Protocol 需定义每个工具的终止条件，防止对话无限延长
- 挑战语气定义为"理性审慎"——每点挑战必须附带事实或逻辑依据
- proposal.md 模板新增 Assumptions Challenged 章节位于 Feasibility Assessment 和 Scope 之间
- Key Risk: AI 过度挑战导致用户反感 → 每点挑战附带事实或逻辑依据，不空泛质疑
