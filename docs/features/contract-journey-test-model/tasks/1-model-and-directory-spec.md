---
id: "1"
title: "Journey-Driven 测试模型定义与目录规范"
priority: "P0"
estimated_time: "3h"
dependencies: []
scope: "all"
breaking: false
type: "documentation"
mainSession: false
---

# 1: Journey-Driven 测试模型定义与目录规范

## Description

定义 Journey-Driven 测试模型的核心概念（Journey、Contract 六维度、Outcome、Risk 分级、语义描述符）和目录结构规范（`tests/<journey>/`、`_contracts/` 子目录）。这是后续所有实现任务的基础定义文档。

来源：proposal 的 Proposed Solution 和 Scope 中的"Journey-Driven 测试模型定义"与"目录规范"两项。

## Reference Files
- `docs/proposals/contract-journey-test-model/proposal.md` — Source proposal
- `plugins/forge/skills/` — 现有 skill 结构参考

## Acceptance Criteria

- [ ] 模型定义文档包含：Journey（用户真实工作流 + Risk 分级 High/Medium/Low）、Contract（六维度：Preconditions/Input/Output/State/Side-effect/Invariants）、Outcome（多 Outcome 按 Preconditions 互斥）、语义描述符（gen-contracts 阶段使用，gen-test-scripts 转换为 regex）
- [ ] 目录规范定义：`tests/<journey>/` 为测试目录（测试直接生成到最终目录），`tests/<journey>/_contracts/` 为 Contract 规范目录
- [ ] 模型定义文档符合 gen-contracts 可直接解析的结构（Journey 名称 + Step 序列 + 每步 Outcome）
- [ ] 定义配置驱动框架的 config.yaml schema（languages、test-framework、test-command、capabilities 字段）

## Hard Rules

- 不硬编码语言或框架名称到模型定义中
- Contract 六维度中 Preconditions/Input/Output/State 为必选，Side-effect 和步骤级 Invariants 为可选

## Implementation Notes

- 参考 proposal 的 Contract 规范文件示例（`tests/task-lifecycle/_contracts/step-2-task-claim.md`）作为目录规范的具体范例
- Risk 分级标准：High = 涉及状态突变或数据丢失风险；Medium = 涉及多步交互；Low = 只读操作
- Journey 级别 Invariants 为必选（跨步骤不变量），步骤级 Invariants 为可选
