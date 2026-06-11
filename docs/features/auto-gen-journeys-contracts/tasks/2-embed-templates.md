---
id: "2"
title: "创建 embed 模板 test-gen-journeys.md 和 test-gen-contracts.md"
priority: "P0"
estimated_time: "1.5h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 2: 创建 embed 模板 test-gen-journeys.md 和 test-gen-contracts.md

## Description

在 `forge-cli/pkg/task/data/` 目录下创建两个新的 embed 模板文件，定义 gen-journeys 和 gen-contracts 自动任务的执行指令。模板需包含 AUTO_COMMIT 和 SKIP_EVAL_GATE 条件指令。

## Reference Files
- `docs/proposals/auto-gen-journeys-contracts/proposal.md` — Source proposal
- `forge-cli/pkg/task/data/test-gen-and-run.md` — 现有模板参考（简洁风格）
- `forge-cli/pkg/task/data/eval-journey.md` — 现有模板参考（复杂模板风格）
- `forge-cli/pkg/task/autogen.go` — renderBody() 占位符逻辑 (L246-304)

## Acceptance Criteria

- [ ] `data/test-gen-journeys.md` 模板存在，使用 `{{FEATURE_SLUG}}`、`{{MODE}}` 占位符
- [ ] `data/test-gen-contracts.md` 模板存在，使用 `{{FEATURE_SLUG}}`、`{{MODE}}` 占位符
- [ ] test-gen-journeys.md 包含条件指令："若 AUTO_COMMIT=true，跳过用户审批步骤，直接 git add + commit 生成的 Journey 文件"
- [ ] test-gen-contracts.md 包含条件指令："若 SKIP_EVAL_GATE=true，跳过 eval-journey 前置检查，直接进行代码侦察和 Contract 生成"
- [ ] 模板通过 `go:embed data/*.md` 自动包含（无需修改 embed 声明）
- [ ] 模板内容清晰指导 AI 执行器完成 gen-journeys/gen-contracts skill 的流程

## Hard Rules

- 模板使用现有 renderBody() 支持的占位符（{{FEATURE_SLUG}}, {{MODE}}, {{TEST_TYPE}} 等）
- AUTO_COMMIT 和 SKIP_EVAL_GATE 不作为代码级变量——它们通过模板文本以自然语言条件指令的形式嵌入任务 .md
- 不修改 renderBody() 函数

## Implementation Notes

- test-gen-journeys.md 需覆盖 proposal.md 输入场景（Quick 模式）和 PRD 输入场景（Breakdown 模式）
- test-gen-contracts.md 需覆盖有 eval 报告场景（Breakdown）和无 eval 报告场景（Quick）
- 参考 eval-journey.md 的模板结构和详细程度
