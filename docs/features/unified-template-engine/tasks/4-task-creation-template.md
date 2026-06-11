---
id: "4"
title: "迁移 pkg/template 任务创建层到 text/template"
priority: "P0"
estimated_time: "1h"
complexity: "medium"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: true
type: "coding.refactor"
mainSession: false
---

# 4: 迁移 pkg/template 任务创建层到 text/template

## Description
将 `pkg/template/template.go` 的 `ApplyVars()` 替换为 `text/template.Execute()`，定义 `taskTemplateData` struct。迁移 2 个模板文件（coding-fix.md, coding-cleanup.md）的占位符语法。移除 `pkg/task/add.go` 中的 `injectSurfaceFrontmatter()`——该函数用 `strings.Replace` 替换模板中硬编码的 `surface-key: ""` 和 `surface-type: ""` 字面值，统一后由模板渲染直接填充。任务创建模板条件化：surface 有值时渲染字段 + 省略 Surface Inference 段落。

## Reference Files
- `forge-cli/pkg/template/template.go`: ApplyVars() 替换为 text/template.Execute()，新增 taskTemplateData struct (source: proposal.md#pkg/template-任务创建层)
- `forge-cli/pkg/template/data/coding-fix.md`: 占位符 `{{X}}` → `{{.X}}`，添加 surface 条件块 (source: proposal.md#pkg/template-任务创建层)
- `forge-cli/pkg/template/data/coding-cleanup.md`: 同 coding-fix.md 迁移 (source: proposal.md#pkg/template-任务创建层)
- `forge-cli/pkg/task/add.go`: 移除 injectSurfaceFrontmatter() 函数 (source: proposal.md#pkg/template-任务创建层)

## Acceptance Criteria
- [ ] `taskTemplateData` struct 包含提案定义的 5 个字段（TaskName, SurfaceKey, SurfaceType, TaskGoal, ScopeDescription）
- [ ] `CreateTaskMarkdown()` 使用 `text/template.Execute()` 渲染，`ApplyVars()` 和 `injectSurfaceFrontmatter()` 已移除
- [ ] 2 个任务创建模板占位符已迁移到 `{{.X}}` 格式，surface 有值时渲染字段并省略 Surface Inference 段落

## Implementation Notes
- `pkg/template` 的 `ApplyVars()` 同时被 `forge task add` 命令和 quality gate hook 调用，迁移需确保两个调用路径均正确
- `forge task add` 命令路径保持软性行为——推断失败时使用空字符串而非报错
- ScopeDescription 是 task-level 语境描述，非 deprecated 的 surface-level SurfaceKey 概念
