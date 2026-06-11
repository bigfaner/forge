---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Pre-Revision Report — Iteration 0

## ATTACK_POINTS

### Factual Corrections

- **[high]** 占位符清单不完整，遗漏多服务编排专用占位符 | quote: "提案声明了 11 个占位符，但当前 server-lifecycle.md 的 Slot Placeholder Reference 表中实际有更多占位符未被纳入提案清单" | improvement: 补充 `<URL_KEY>`, `<SERVICE_LIST>`, 多服务编排的 6 个专用占位符，或在提案中明确声明多服务编排模式的处理策略

- **[high]** 删除 Phase 1 一致性检查存在 agent 填值引入错误的风险 | quote: "即使 CLI 生成的骨架是正确的，agent 在填占位符时可能引入错误：填错占位符值、遗漏某个占位符、或者破坏了 just 语法结构" | improvement: 保留轻量级 recipe 完整性验证（just --list 检查）替代完全删除 Phase 1

- **[high]** 未覆盖多服务编排模式的 recipes 生成 | quote: "server-lifecycle.md 第 309-425 行包含了完整的多服务编排模式（Port-Aware Startup Order、Multi-Service Teardown、test-setup with embedded multi-service lifecycle），但提案的每个 surface type 生成的 recipes 表格中没有任何关于多服务编排的提及" | improvement: 在 --aggregate 模式中增加跨 surface test 编排 recipe，或明确声明不支持

- **[high]** Convention cold start fallback 策略被删除但未说明是否仍需维护 | quote: "占位符解析表中部分占位符（如 {{START_CMD}}、{{PORT}}）需要 Convention 知识才能正确填写，而 Convention 文件可能不存在（cold start）" | improvement: 在 SKILL.md 中保留 5-10 行 cold start fallback 策略摘要

### Structural/Architectural Suggestions

- **[medium]** 取消 gate recipe 概念未说明对 consumer 端的影响 | quote: "提案取消 gate recipe 概念，统一为 surface recipe，但没有说明这对 consumer 端（run-tests skill、quality gate 机制）的影响" | improvement: 增加 Consumer Impact 小节，说明 run-tests/quality gate 的 recipe 名解析变更

- **[medium]** ci 聚合 recipe 不包含 surface-level test | quote: "聚合 recipe 的 ci 定义不包含 surface-level test" | improvement: 在提案中明确声明 ci 是否包含 surface test 及设计理由

- **[medium]** user-customized 标记策略不明确 | quote: "提案将所有 lifecycle recipes 标记 # user-customized，但当前 SKILL.md 中 # user-customized 有明确的作用域定义和分层语义" | improvement: 统一声明所有 lifecycle + quality recipes 标记 # user-customized，aggregate 不标记

## BORDERLINE_FINDINGS

(none)

## SKIPPED_FINDINGS

(none — all suggestions were classified as structural or factual)

## Rubric Data

All dimensions: N/A (freeform pre-revision, no rubric scoring)
