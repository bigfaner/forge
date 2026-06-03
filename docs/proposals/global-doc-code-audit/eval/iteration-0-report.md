---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Pre-Revision (Freeform Findings)

## ATTACK_POINTS

- **[high]** 提案声称 conventions 下有 22 份规范文档，实际仅 18 份（差额 4 份） | quote: "docs/conventions/ 下 22 份规范文档（顶层15份 + testing/7份）" | improvement: 验证实际文件数量并修正为 "18 份规范文件（顶层15份 + testing/ 子目录3份）"

- **[high]** docs/reference/ 目录不存在，提案引用的 test-type-model.md 实际位于 plugins/ 内部 | quote: "docs/reference/ 仅含 1 个文件 test-type-model.md，不单独建 Task" | improvement: 删除对 docs/reference/ 的引用，修正 L2 范围描述

- **[high]** 146 个 task 的数字无法验证，5 个审计提案目录中均无 tasks/ 子目录 | quote: "执行现有5个提案的146个task" | improvement: 核实数字来源或删除此具体数字，改为概数

- **[high]** docs/features/ 被排除但 conventions 是从 feature 提取的派生产物，排除导致无法追溯根因 | quote: "提案将 docs/features/（182 个子目录）和 docs/proposals/（149 个子目录）完全排除在审计范围之外" | improvement: 在 Out of Scope 部分补充排除理由，说明 L2 发现不一致时的根因追溯路径

- **[high]** 16 小时时间预算对语义比对极其紧张（每文件不到 30 分钟，每条目不到 4 分钟） | quote: "不超过 2 个工作日（约 16 小时有效工作时间）" | improvement: 调整时间估算或减少单次审计范围，提供更合理的时间预算

- **[high]** 层级间反馈机制与紧凑时间线矛盾，跨层协调可操作性存疑 | quote: "三层审计之间存在交叉依赖。但提案同时要求 2 个工作日完成全部三层" | improvement: 量化反馈机制的时间开销，纳入总时间预算，或放松时间约束

- **[high]** docs/ 目录不分发到用户环境，L2 审计忽略了分发模型约束 | quote: "docs/ 目录的内容不分发到用户环境——分发的仅是 plugins/forge/ 下的内容" | improvement: 在 Problem 或 Constraints 部分区分源码仓库维护者和终端用户的文档消费路径

- **[medium]** 提案声称根目录除 README.md 和 DESIGN.md 外无其他文档，实际存在 CLAUDE.md | quote: "根目录下除 README.md 和 DESIGN.md 外无其他面向用户的文档文件" | improvement: 在范围完整性说明中提及 CLAUDE.md 并说明排除/纳入理由

- **[medium]** 遗漏率指标需要已知全集才能测量，但审计目的是发现全集 | quote: "遗漏率不超过 20%" | improvement: 改为可验证的过程性标准

- **[medium]** 20% 遗漏率阈值允许每 5 个问题遗漏 1 个 | quote: "20% 的阈值意味着允许每 5 个问题中遗漏 1 个" | improvement: 降低阈值或改用交叉验证重合度作为替代指标

- **[medium]** "可由 task-executor 独立执行" 与 "知识库清理需人工确认" 存在措辞张力 | quote: "修复类 Task 可由 task-executor 独立执行" vs "知识库清理需人工确认，不可自动删除" | improvement: 区分两类 Task 模板——修复类用标准模板，审查类用带人工确认节点的模板

- **[medium]** 审计产出英文但源文档中文，语言切换增加误读风险，动机未说明 | quote: "审计产出的所有 Skill、Command、任务模板、提示词模板等统一采用英文撰写" | improvement: 补充英文约束的动机说明

- **[low]** 建议：补充三类比对的具体操作协议 | quote: "声明提取 → 代码定位 → 逐条比对" | improvement: 为路径引用验证、行为描述验证、配置声明验证分别提供操作步骤

- **[low]** 建议：明确与 consolidate-specs skill 的关系 | quote: 未提及 consolidate-specs | improvement: 说明审计与 drift detection 的差异

- **[low]** 建议：对 L2 范围做完整文件系统扫描验证，明确排除理由 | quote: "business-rules/4 + conventions/22 + reference/1" | improvement: 列出 docs/ 下所有子目录并逐一说明排除理由

## BORDERLINE_FINDINGS

- Finding: "建议重新审视 docs/features/ 的排除决定" — 已在 ATTACK_POINTS 中处理（high severity finding 4）。保留排除但补充排除理由。

## SKIPPED_FINDINGS

- Finding: "建议将 L3 知识库审查与 L1/L2 在方法论上明确区分" — Subjective preference: 方法论区分属于执行细节。已有 L3 单独流程描述。

## Classification Audit

- Factual correction: 5 findings (conventions 数量、reference 目录、146 task 数字、CLAUDE.md 遗漏、英文约束动机)
- Structural suggestion: 8 findings (features 排除理由、遗漏率指标、20% 阈值、措辞张力、时间预算、反馈机制、分发模型、操作协议)
- Subjective preference: 1 finding (L3 方法论区分), skipped

## Rubric

All dimensions: N/A (freeform findings, not rubric-scored)
