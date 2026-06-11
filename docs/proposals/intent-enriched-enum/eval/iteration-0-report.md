---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Pre-Revision Report (Iteration 0)

## ATTACK_POINTS

1. **[high]** 6值intent无法覆盖8值task type，"干净的1:1映射"表述与实际不符 | quote: "与 task type 形成干净的 1:1 映射" | improvement: 将"干净的1:1映射"改为明确说明覆盖范围，承认doc.consolidate/doc.drift无对应intent并说明理由

2. **[high]** intent:fix映射到coding.fix会打破"do not assign manually"规则 | quote: "coding.fix 在 Type Assignment 中被标注为 'do not assign manually'" | improvement: 决定fix的映射策略——要么允许手动分配coding.fix（更新Type Assignment表），要么fix映射到其他类型并在提案中解释偏离1:1的理由

3. **[high]** coding.fix的"do not assign manually"约束与1:1映射矛盾——这是提案必须解决的核心架构决策 | quote: "coding.fix 在 breakdown-tasks 和 quick-tasks 中被标注为 'Auto-generated for test failures via forge task add; do not assign manually'" | improvement: 在提案的Requirements Analysis中显式处理此矛盾，给出设计决策和理由

4. **[medium]** Override Signals仅停留在概念层面，缺乏结构化检测条件 | quote: "Override Signals 的具体规则在提案中只停留在概念层面，没有给出信号检测的具体条件" | improvement: 定义结构化检测条件表（信号类型 → 关键词/模式 → 覆盖动作），使LLM可可靠执行

5. **[medium]** enhancement的完整pipeline行为未定义——PRD格式、tech-design决策类型、self-check规则均缺失 | quote: "enhancement 的 pipeline 行为没有明确定义...提案没有讨论 enhancement 是否需要自己的 PRD 格式变体" | improvement: 在提案中定义enhancement的完整pipeline行为（建议：简化版full PRD——保留Background/Goals，跳过User Stories，保留Test Pipeline）

6. **[medium]** Pipeline Configuration表缺少enhancement和fix行的行为定义 | quote: "enhancement 和 fix 行的行为在提案中没有定义" | improvement: 在提案中给出完整的6行Pipeline Configuration表

7. **[medium]** scope遗漏rules/子目录中可能引用intent值的文件 | quote: "write-prd 和 tech-design 的 rules/ 目录中可能还有其他文件引用了 intent 值" | improvement: 执行grep扫描所有rules/文件中的intent引用，将遗漏文件加入In Scope

8. **[medium]** Industry Benchmarking缺少具体引用——"分类系统的精确度随场景增长自然演进"是常识陈述而非benchmarking | quote: "分类系统的精确度随场景增长自然演进是常见模式" | improvement: 引用至少一个真实系统的分类演进案例

9. **[medium]** Override Signals否决了"完全内容驱动"但本质上是缩小版的内容驱动，未论证范围缩小是否足以补偿稳定性风险 | quote: "完全内容驱动 pipeline...依赖 LLM 判断力，不稳定" vs "PRD 内容中的明确信号可以覆盖默认值" | improvement: 在提案中明确论证混合模式的稳定性边界——结构化条件表的LLM遵守度显著高于prose判断

## BORDERLINE_FINDINGS

- enhancement和fix到task type的映射存在未解决的矛盾（与finding #2/#3相关但不完全相同——此条关注映射表本身的定义缺失而非coding.fix的手动分配限制）

## SKIPPED_FINDINGS

- proposal template的intent默认值和有效值注释格式（low severity，属于实现细节，可在tech-design阶段处理）
- Next Steps应走write-prd而非quick-tasks（low severity，属于流程建议，不影响提案质量）
- rules/子目录完整扫描（已作为finding #7的改进方向包含）

## Rubric Data

All dimensions: N/A (pre-revision does not use rubric scoring)
