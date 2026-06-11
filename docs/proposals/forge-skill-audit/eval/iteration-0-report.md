---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Iteration 0: Pre-Revision Report

## ATTACK_POINTS

- **[high]** eval-journey 描述遗漏第 7 维度 "Workflow Coverage"，LLM scorer 可能跳过该维度 | quote: "eval-journey.md description reads 'Scores completeness, semantic purity, precondition exclusivity, fact alignment, surface fitness, and internal consistency' — that is 6 dimensions. The actual journey rubric has 7 dimensions, with 'Workflow Coverage' at 150 points being the omitted dimension" | improvement: H-1 修复范围扩展，更新 eval-journey.md 描述字段补充第 7 维度

- **[high]** eval-contract 描述声称 "six-dimension" 但实际有 8 维度，遗漏 Anchor Integrity 和 Fixture Specification | quote: "eval-contract.md description reads 'Scores six-dimension structural integrity' and lists 6 evaluation areas, but the contract rubric has 8 dimensions" | improvement: H-1 修复范围扩展，更新 eval-contract.md 描述字段补充缺失维度

- **[high]** eval/SKILL.md 声称 "Supports 100-point and 1000-point scales" 但无 100 分制 rubric，且 1100/1150 超出 1000 | quote: "If the scorer believes the scale ceiling is 1000, it may cap scores for journey/contract evaluations" | improvement: H-1 修复范围扩展，更新 eval/SKILL.md 描述字段

- **[medium]** M-9 仅列出 3 处 INLINE 引用但实际有 4 处，遗漏 gen-contracts 对 gen-journeys 的双向依赖 | quote: "Actual grep of <!-- INLINE:origin= across all skill SKILL.md files reveals 4 INLINE references. The missing one is gen-contracts/SKILL.md line 58" | improvement: 更新 M-9 描述为 4 处，标注 gen-contracts↔gen-journeys 双向依赖

- **[medium]** Problem 段声称 "22 个 skill" 但实际目录只有 21 个 | quote: "The proposal's Problem section states '22 个 skill', but the actual plugins/forge/skills/ directory contains 21 subdirectories" | improvement: 修正为 "21 个 skill"

- **[medium]** H-1 修复仅覆盖 5 处同步点中的 4 处，config 键默认值未提及 | quote: "The H-1 analysis itself identifies '5 处' that need synchronization. The proposed fix addresses 4 of the 5 but does not mention updating config key defaults" | improvement: H-1 修复描述中显式标注第 5 处为已知缺口

- **[medium]** M-1 成功标准含条件前置，可能导致 M-1 被完全跳过 | quote: "This creates a conditional success criterion that may cause M-1 to be skipped entirely if the Go config reader check fails" | improvement: 澄清 M-1 是承诺执行还是有条件推迟

- **[medium]** C 维度审计范围声称 "约40 模板文件" 但实际有 56 个 | quote: "The Summary Statistics table for dimension C states '约40 模板文件' as the audit scope, but the actual path contains 56 files" | improvement: 更新为准确计数

## BORDERLINE_FINDINGS

- H-2 仅移除一条死路径但未分析 docs/features/ 与 docs/proposals/ 的路径拓扑关系 — 需要更广泛的路径分析，超出当前 fix scope
- 提案未枚举哪些下游消费者直接读取 rubric-reference.md — 需要全量 grep 分析，标记为后续验证步骤

## SKIPPED_FINDINGS

- M-9 版本戳改用含行数校验的结构化格式 — 主观偏好改进方案
- 添加 eval/SKILL.md 和 eval-contract.md 的回归验证 grep 步骤 — 已在 Regression Verification 中有类似覆盖
- H-1 修复描述显式枚举全部 5 处真实来源 — 已在 ATTACK_POINTS 中要求标注缺口

## Rubric

All dimensions: N/A (pre-revision, no rubric scoring)
