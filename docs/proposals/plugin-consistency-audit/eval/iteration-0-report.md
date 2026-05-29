---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Pre-Revision (Freeform Findings)

## Triage Summary

9 findings triaged: 3 factual correction, 5 structural suggestion, 1 borderline (deferred), 0 subjective (skipped)

## ATTACK_POINTS

### Factual Corrections

- **[high]** 提案声称有22个skill但实际只有21个，导致成功标准不可验证 | quote: "22 个 skill 100% 覆盖审计，每个 skill 的 SKILL.md 与其 templates/rules/data 逐一对比" | improvement: 核实实际skill数量并更新提案中所有出现位置

- **[high]** 文件计数"172+"与实际不符（实际208或182），削弱紧迫性论证 | quote: "172+ 个 .md 文件通过手动维护交叉引用，重构过程中依赖局部修改而非全局验证" | improvement: 核实实际文件计数并更新，或明确指代范围

- **[low]** 提案自身 frontmatter 的 status 和 intent 元数据与实际状态不匹配 | quote: "status: Draft / intent: refactor" | improvement: 修正 metadata 使其与提案实际状态一致

### Structural Suggestions

- **[high]** 子目录分类遗漏 examples/、types/、experts/，审计可能漏检关键文件 | quote: "22 个 skill 的 SKILL.md 与其各自的 templates/rules/data 之间的逻辑自洽性" | improvement: 扩展子目录分类覆盖所有实际存在的目录类型

- **[medium]** eval skill 的 rubrics/experts 被排除在外，约19%的skill文件未被审计 | quote: "rules/rubrics/experts 的功能性质量审查" | improvement: 细化排除边界，区分"功能性质量审查"与"交叉引用校验"

- **[medium]** hooks/guide.md 作为单文件的"内部一致性"含义未定义 | quote: "hooks/guide.md 的内部一致性" | improvement: 明确定义 guide.md 审计的操作含义

- **[medium]** 四类分类方案未覆盖 INCOMPLETE、DEPRECATED 等失败模式 | quote: "问题分类: 矛盾(CONFLICT)、冗余(REDUNDANT)、时序(TIMING)、引用(REFERENCE)" | improvement: 增加 INCOMPLETE 分类覆盖结构性缺口

- **[medium]** P0-P3 严重等级缺乏定义，导致报告不可复现 | quote: "每个问题包含: 文件路径、问题描述、严重等级(P0-P3)、修复建议" | improvement: 为每个严重等级定义明确的判定标准

## BORDERLINE_FINDINGS

- **[medium]** 仅审计单组件自洽性可能遗漏跨组件接口不一致 | Defer rationale: 该发现质疑提案的核心设计决策（单组件自洽），属于战略层面取舍而非内部一致性问题。用户已明确选择此范围。交由 Scorer 评分流程评估。

## SKIPPED_FINDINGS

(None — all suggestions correspond to risk/problem findings above)

## Classification Audit

- Factual correction: 3 items (wrong skill count, wrong file count, metadata mismatch)
- Structural suggestion: 5 items (taxonomy gap, scope boundary, undefined operation, classification gap, undefined severity)
- Subjective preference: 0 items
