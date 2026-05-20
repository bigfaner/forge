---
created: 2026-05-20
author: "fanhuifeng"
status: Approved
---

# Proposal: Brainstorm Adversarial Challenge Enhancement

## Problem

Brainstorm skill 对用户想法的挑战力度不足。当前 Challenge 集群是三个并列决策集群之一，AI 可以自由跳过或走过场，导致生成的提案缺乏深度审视，不够健壮的方案直接进入下游流水线。

### Evidence

- Decision Clusters 中 Challenge 是可选的"自由遍历"节点，无强制触发机制
- AI 天然倾向配合用户，软性建议无法克服顺从惯性
- 用户反馈：对话阶段太顺从、质疑不够深入、缺乏系统性挑战框架

### Urgency

v3.0.0 正在开发中，brainstorm 是整个 Forge 流水线的入口。入口质量直接决定下游 PRD、设计、任务的质量。现在不改，v3 发布后所有通过 quick mode 进来的功能都会缺乏充分审视。

## Proposed Solution

将挑战机制从"可选集群"升级为"融入每个决策点的强制行为"。具体：重写 SKILL.md 的 Decision Clusters，在每个集群内嵌入对应的挑战工具；同时更新 proposal.md 模板，新增 Assumptions Challenged 章节记录被验证/推翻的假设。

### Innovation Highlights

**框架先行、事实打击**：先用结构化挑战工具（5 Whys、假设翻转等）定位薄弱点，再用事实证据（代码库、逻辑一致性、领域常识）精确反驳。不是空泛质疑，是有方法的挑战。

## Requirements Analysis

### Key Scenarios

- 有代码库的项目：AI 查代码库事实反驳用户假设
- 冷启动/新项目：无代码可查，通过逻辑一致性检查和领域常识挑战
- 用户想法不清晰：通过 XY 问题检测引导用户找到真正需求
- 用户想法看似完美：通过极端场景推演暴露隐藏风险

### Non-Functional Requirements

- 挑战不能导致对话无限延长——每个集群有明确的挑战终止条件
- 挑战语气应是"理性审慎"而非"敌意攻击"
- 不增加用户的认知负担——AI 主动发起挑战，用户只需回应

### Constraints & Dependencies

- 不改变 7 步流程的整体结构
- 不影响 eval-proposal 的独立对抗机制
- 改动限于 SKILL.md 和 proposal.md 两个文件

## Alternatives & Industry Benchmarking

### Industry Solutions

- Six Thinking Hats（德博诺）：结构化多角度思考
- Red Team / Devil's Advocate：独立对抗角色
- 5 Whys（丰田）：根因分析

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 当前挑战形同虚设 | Rejected: 用户已确认问题存在 |
| 最小增强：只加强 Challenge 集群措辞 | — | 改动小 | 仍是软建议，AI 可跳过 | Rejected: 不解决根本问题 |
| 独立 Adversarial 步骤 | Red Team 模式 | 挑战彻底 | 增加流程步骤，可能拖延对话 | Rejected: 用户选择融入式 |
| **融入式挑战 + 事实驱动** | 5 Whys + 假设翻转 + 奥卡姆剃刀 | 每个决策点自带挑战，不可跳过 | SKILL.md 改动较大 | **Selected: 结构性解决顺从问题** |

## Feasibility Assessment

### Technical Feasibility

完全可行。改动是 prompt 文本（SKILL.md 和模板），不涉及代码逻辑。

### Resource & Timeline

一个任务即可完成：修改两个文件。

### Dependency Readiness

无外部依赖。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "加个 Challenge 集群就够了" | 假设翻转：如果 Challenge 是可选的，AI 还会认真执行吗？ | 推翻：可选 = 可跳过，必须融入每个集群 |
| "挑战需要代码库证据" | 极端推演：冷启动项目无代码时怎么办？ | 修正：泛化为"事实驱动"，证据来源包括代码、逻辑一致性、领域常识 |
| "独立步骤更能保证挑战质量" | 奥卡姆剃刀：融入式是否更简单且同样有效？ | 确认：融入式避免流程膨胀，且每个决策点都受挑战 |

## Scope

### In Scope

- 重写 SKILL.md 的 Decision Clusters，挑战融入每个集群（Problem → 5 Whys + XY 检测，Solution → 假设翻转 + 极端推演）
- 新增 Challenge Protocol 章节，定义 5 个挑战工具的使用时机、触发条件和终止条件
- 将"证据驱动"泛化为"事实驱动"（代码库事实 + 逻辑一致性 + 领域常识）
- proposal.md 模板新增 Assumptions Challenged 章节
- 奥卡姆剃刀作为贯穿全程的元原则

### Out of Scope

- 7 步流程结构不变（Step 1,3,4,5,6,7 不动）
- eval-proposal 不改（已有独立对抗机制）
- 其他 skill 不动
- 代码改动

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| AI 过度挑战导致用户反感 | M | M | 挑战语气定义为"理性审慎"，每点挑战附带事实或逻辑依据 |
| 挑战工具增加对话轮次 | M | L | 定义终止条件：假设被验证或用户给出充分理由即停 |
| 模板新增章节增加提案写作负担 | L | L | Assumptions Challenged 用表格格式，结构化降低写作成本 |
| 事实驱动在冷启动项目中无事实可用 | L | M | 已泛化定义：逻辑一致性 + 领域常识也是事实来源 |

## Success Criteria

- [ ] 每个 Decision Cluster 都有明确绑定的挑战工具，不是可选的
- [ ] 挑战必须附带事实依据（代码库/逻辑一致性/领域常识），不允许空泛质疑
- [ ] 冷启动项目（无代码库）的挑战同样有效
- [ ] proposal.md 模板的 Assumptions Challenged 章节能记录假设验证过程
- [ ] 改动不增加 7 步流程的步骤数量

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
