---
id: "4"
title: "gen-contracts skill + TUI await 形式化"
priority: "P0"
estimated_time: "4h"
dependencies: ["3"]
scope: "backend"
breaking: false
type: "feature"
mainSession: false
---

# 4: gen-contracts skill + TUI await 形式化

## Description

实现 gen-contracts skill：从 Journey 文档 + 代码侦察（Fact Table）推导 Contract 规范（六维度 + 语义描述符 + 多 Outcome + Invariants）。同时形式化 TUI 异步 Cmd 的 await 语义，定义为 Contract 维度的 await 规范。

来源：proposal Pipeline 第 2 步、Scope "gen-contracts skill"和"TUI await 语义形式化"。

## Reference Files
- `docs/proposals/contract-journey-test-model/proposal.md` — Source proposal（含 Contract 规范文件示例）
- `plugins/forge/skills/gen-test-scripts/references/` — 现有 Fact Table 使用方式
- `forge-cli/pkg/e2e/` — Fact Table 机制

## Acceptance Criteria

- [ ] gen-contracts 从 Journey + Fact Table 推导 Contract，每个 Outcome 包含 4 个必选维度（Preconditions/Input/Output/State），可省略 Side-effect 和步骤级 Invariants
- [ ] 所有生成的 Contract 通过六维度完整性校验（必选维度非空，语义描述符不含 regex 语法）
- [ ] 每个 Journey 生成 Journey 级别 Invariants（跨步骤不变量声明，至少 1 条）
- [ ] 多 Outcome Contract 正确性：每个 Outcome 的 Preconditions 互斥（同一输入不会同时满足两个 Outcome 的 Preconditions）
- [ ] TUI await 语义：异步 Cmd 等待超时时 fail-fast 并报告超时 Cmd 名称；`tea.Batch(cmd1, cmd2)` 等待全部完成后再进入下一步
- [ ] Contract 规范存储为结构化文档（markdown with schema），存放在 `tests/<journey>/_contracts/` 目录

## Hard Rules

- 语义描述符不含 regex 语法（gen-contracts 阶段不生成正则）
- Outcome 按 Preconditions 互斥，消除组合爆炸
- 超过 5 个 Outcome 的步骤触发 LLM 检查点
- 无法归入已有维度的验证点自动归入 Invariants 并标注 `dimension: unclassified`

## Implementation Notes

- gen-contracts 是最复杂的 skill，需要代码侦察（Fact Table）推导六维度
- State 维度降级路径：项目未暴露状态查询接口时 → `state-verification: partial`（仅从 Output 推断）或 `state-verification: deferred`
- TUI await 形式化：`await` 表示等待所有 pending Cmd 完成（最长 N 毫秒），超时 fail-fast
- 分批生成：Contract 数超过 15 或 token 数超过 50k 时自动拆分，合并后内容一致
