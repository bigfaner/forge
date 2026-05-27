---
id: "3"
title: "简化 cleanTemplateOutput 并添加模板校验"
priority: "P0"
estimated_time: "1h"
complexity: "low"
dependencies: [2]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 3: 简化 cleanTemplateOutput 并添加模板校验

## Description
将 `cleanTemplateOutput()` 从 4 种条件删除逻辑简化为仅保留空白行塌陷。条件逻辑已由 Task 2 的模板级 `{{if}}` 块处理，后处理不再需要。新增 `ValidatePromptTemplates()` 使用 `missingkey=error` + 零值 struct Execute 到 io.Discard，确保无字段拼写错误。

## Reference Files
- `forge-cli/pkg/prompt/prompt.go`: cleanTemplateOutput() 移除空标签行、空 backtick 条件句、just 命令尾部空白、<!-- IF NOT_LOW --> 段落块四种条件删除逻辑 (source: proposal.md#pkg/prompt-渲染层)
- `forge-cli/pkg/prompt/prompt.go`: 新增 ValidatePromptTemplates() 启动时校验函数 (source: proposal.md#Key-Risks)

## Acceptance Criteria
- [ ] `cleanTemplateOutput()` 仅保留空白行塌陷逻辑，四种条件删除模式（空标签行、空 backtick 条件句、`just` 命令尾部空白、`<!-- IF NOT_LOW -->` 段落块）已全部移除
- [ ] `ValidatePromptTemplates()` 使用 `template.Option("missingkey=error")` 配置，对零值 `promptTemplateData` 执行 `Execute()` 到 `io.Discard`
- [ ] `go build ./...` 通过，`forge prompt get-by-task-id` 输出与迁移前功能等价（允许空白行差异）

## Implementation Notes
- 复杂度判定覆盖：AC=3 且无 Hard Rules，但涉及多文件架构验证。考虑到主要是删除代码 + 添加简单校验函数，认定为 low。
- `text/template` 的 `{{-`/`-}}` 可消除单个 `{{if}}` 块周围的空行，但无法处理连续空行（多个条件块在同一位置省略时产生的多行空白），保留 Go 级空白行塌陷作为最终后处理步骤
- 旧的 `isLabelWithEmptyValue()` 对标签名有空格限制，迁移后此限制自然消除
