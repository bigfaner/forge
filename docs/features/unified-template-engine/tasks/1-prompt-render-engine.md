---
id: "1"
title: "定义 promptTemplateData 并迁移 renderTemplate 到 text/template"
priority: "P0"
estimated_time: "2h"
complexity: "medium"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: true
type: "coding.refactor"
mainSession: false
---

# 1: 定义 promptTemplateData 并迁移 renderTemplate 到 text/template

## Description
将 `pkg/prompt/prompt.go` 的 `renderTemplate()` 从 `strings.ReplaceAll` 迁移到 `text/template.Execute()`。定义 `promptTemplateData` struct 暴露模板所需的全部字段。将 `TASK_CATEGORY` 字符串拼接逻辑迁移到 struct 的 `TaskCategory` 字段，由模板条件块渲染。这是 prompt 渲染层的核心引擎变更，后续模板文件迁移和后处理简化均依赖此任务。

## Reference Files
- `forge-cli/pkg/prompt/prompt.go`: renderTemplate() 从 strings.ReplaceAll 迁移到 text/template.Execute()，新增 promptTemplateData struct (source: proposal.md#pkg/prompt-渲染层)
- `forge-cli/pkg/prompt/prompt.go`: TASK_CATEGORY 字符串拼接逻辑（第136-137行）迁移到 promptTemplateData.TaskCategory 字段 (source: proposal.md#renderTemplate-中-TASK_CATEGORY-注入迁移)

## Acceptance Criteria
- [ ] `promptTemplateData` struct 包含提案中定义的全部 11 个字段（TaskID, TaskFile, TaskCategory, FeatureSlug, PhaseSummary, CoverageStrategy, CoverageTarget, TestTypeArg, SurfaceKey, SurfaceType, Complexity）
- [ ] `renderTemplate()` 使用 `text/template.Parse()` + `Execute()` 渲染，不再调用 `strings.ReplaceAll`
- [ ] `TASK_CATEGORY` 字符串拼接逻辑已从 `renderTemplate()` 移入 `promptTemplateData.TaskCategory`，由调用端设置

## Implementation Notes
- 模板数据 struct 使用值类型（string），所有字段零值为空字符串，模板中 `{{if .Field}}` 对空字符串为 false，无需 nil 检查
- `renderTemplate()` 需使用 `template.Option("missingkey=error")` 配置，确保字段拼写错误在开发期暴露
- 此任务仅做引擎层迁移，模板文件的占位符语法和条件块在 Task 2 中处理
