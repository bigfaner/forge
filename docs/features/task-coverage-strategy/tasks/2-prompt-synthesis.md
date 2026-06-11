---
id: "2"
title: "Inject coverage target into prompt synthesis"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 2: Inject coverage target into prompt synthesis

## Description

修改 prompt 合成流程，将覆盖率目标从配置解析结果注入到 prompt 模板中。合成时按优先级解析有效覆盖率：任务 frontmatter `coverage` > 全局 config `coverage` > 内置默认值。

## Reference Files
- `docs/proposals/task-coverage-strategy/proposal.md` — Source proposal
- `forge-cli/pkg/prompt/prompt.go` — Prompt 合成核心逻辑
- `forge-cli/pkg/prompt/data/*.md` — 嵌入的 prompt 模板文件
- `forge-cli/pkg/task/types.go` — Task struct 定义

## Acceptance Criteria

- `renderTemplate()` 新增 `{{COVERAGE_TARGET}}` 占位符替换
- 新增 `{{COVERAGE_STRATEGY}}` 占位符，值为 `percentage` 或 `maintain`
- `{{COVERAGE_TARGET}}` 对于 percentage 策略渲染为 `"达到 N% 测试覆盖率"`，对于 maintain 策略渲染为 `"保持现有覆盖率，下降不超过 2%"`
- 优先级正确：frontmatter `coverage` > config per-type > built-in default
- 对于非 testable 类型（`doc*`、`gate` 等），不注入覆盖率相关占位符
- 现有测试通过

## Hard Rules

- 占位符替换使用现有的 `strings.ReplaceAll` 方式，不引入模板引擎
- `SynthesizeOpts` 可扩展但不要添加过多字段，优先从 config 和 task 数据中实时读取

## Implementation Notes

- `SynthesizeOpts` 可能需要新增 `Config` 或在合成函数内部读取 config
- 对于 maintain 策略，`{{COVERAGE_TARGET}}` 和 `{{COVERAGE_STRATEGY}}` 的组合应产生明确的"不新增测试，保持覆盖率"指令
- 需要修改 `Task` struct（`types.go`）来携带 coverage 信息（从 frontmatter 传入），或者在合成时实时读取 frontmatter
