---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Pre-Revision Report (Iteration 0)

## ATTACK_POINTS

Triage: factual correction / structural suggestion

### Factual Corrections

- **[high]** 领域级专家评审深度退化的缓解策略是循环论证 | quote: "LLM 在领域内细化专业方向时参考 proposal 内容，保持针对性。但这恰好是当前系统已经在做的事情" | improvement: 承认缓解策略无效，给出同时实现覆盖面和深度的具体机制
- **[high]** 分类表"8-12个大领域"假设缺乏实证支撑 | quote: "未提供任何关于这个数字如何得出的依据。当前 11 个 proposal 横跨了多少个自然领域？提案没有做这个基础分析" | improvement: 增加现有 proposal 领域分布分析章节，用数据验证范围
- **[medium]** "改动仅涉及 prompt 文件（Markdown），不涉及代码变更"过于简化 | quote: "三个 prompt 文件的联动关系并非简单：expert-inference.md 的新两步流程需要与 freeform-expert-persistence.md 的 Jaccard 匹配逻辑协调" | improvement: 修改 Feasibility Assessment，承认 prompt 链联动复杂度
- **[medium]** extraction-prompt.md 未被纳入影响范围 | quote: "extraction-prompt.md（提案未提及但存在于 freeform 目录中）也可能需要调整，因为它从 freeform review 中提取 findings 时需要知道专家的 scope 级别" | improvement: 补充 extraction-prompt.md 到 In Scope

### Structural Suggestions

- **[high]** 跨领域 proposal 被强制归入单一领域导致评审盲区 | quote: "越是跨领域的复杂 proposal（往往是最需要深度评审的），评审质量越可能不足" | improvement: 为跨领域 proposal 设计显式的多专家策略或至少承认此限制
- **[medium]** 分类表对"一致性"的承诺过度 | quote: "如果 LLM 对同一领域的不同 proposal 能可靠地映射到同一个分类表条目，那么理论上它也应该能可靠地生成一致的 domain 关键词——提案没有解释为什么前者比后者更可靠" | improvement: 将 Innovation Highlights 中的一致性承诺降级为"缩小搜索空间"
- **[medium]** 覆盖范围与评审深度零和关系未给出机制 | quote: "覆盖范围扩大和评审深度之间是零和关系，提案没有给出如何同时实现两者的机制" | improvement: 重新定义 SC-1，或承认 trade-off 并在 Key Risks 中具体化
- **[medium]** 旧专家文件持续参与匹配成为噪声数据 | quote: "旧的 proposal-specific 专家本质上成为了新系统中的噪声数据，而提案没有计划处理这个噪声" | improvement: 增加 In Scope 或将旧专家处理纳入设计

## SKIPPED_FINDINGS (Subjective Preference)

- 建议将分类表从 prompt 内嵌改为外部数据文件（架构偏好，非方案缺陷）
- 建议旧专家匹配时降权（实现策略，非 proposal 层面问题）
- 建议多专家评审策略（扩展方案，超出当前 scope）

## BORDERLINE_FINDINGS

- "分类表嵌入 prompt 形成单点维护负担"：介于 structural（设计选择有隐含风险）和 subjective（外部文件是替代架构偏好）之间。按 structural 处理。

## Classification Audit

- Factual corrections: 4 (2 high, 2 medium)
- Structural suggestions: 4 (1 high, 3 medium)
- Subjective preferences: 3 (all low severity suggestions)
- Total findings triaged: 11

## Rubric

(All dimensions): N/A
