---
id: "3"
title: "gen-journeys skill"
priority: "P0"
estimated_time: "1.5h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "feature"
mainSession: false
---

# 3: gen-journeys skill

## Description

实现 gen-journeys skill：从 PRD 用户故事提取 Journey（叙述性工作流文档），包含 happy path + 边缘场景 + Risk 分级。输出按 Journey 分文件，一个用户工作流一个文件。

来源：proposal Pipeline 第 1 步和 Scope "gen-journeys skill"。

## Reference Files
- `docs/proposals/contract-journey-test-model/proposal.md` — Source proposal
- `plugins/forge/skills/` — 现有 skill 结构参考
- `docs/conventions/forge-distribution.md` — Forge 分发模型

## Acceptance Criteria

- [ ] gen-journeys 从 PRD 用户故事提取 Journey，输出 Markdown 文档
- [ ] 每个 Journey 包含：名称、风险等级（High/Medium/Low）、至少 1 个 happy path 步骤、至少 1 个边缘场景步骤
- [ ] 输出按 Journey 分文件（一个用户工作流一个文件），而非按接口类型分文件
- [ ] 文档格式符合 gen-contracts 可直接解析的结构（Journey 名称 + Step 序列 + 每步的用户操作和期望结果）
- [ ] 高风险 Journey 的边缘场景数量 ≥ happy path 步骤数量

## Hard Rules

- gen-journeys 不需要代码侦察（纯叙述性提取）
- 单次生成一个 Journey；如果用例数超出上下文窗口则自动分批（同一 Journey 内）
- gen-journeys 输出必须是 gen-contracts 可直接消费的结构化格式

## Implementation Notes

- gen-journeys 是 4 步 pipeline 的第 1 步，认知任务最简单（纯叙述性提取）
- Risk 分级由 gen-journeys 根据 PRD 的严重性和失败影响推断
- 分批生成：Contract 数超过 15 或 token 数超过 50k 时自动拆分
