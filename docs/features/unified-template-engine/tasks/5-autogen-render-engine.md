---
id: "5"
title: "迁移 pkg/task/autogen 渲染引擎并清理 scope"
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

# 5: 迁移 pkg/task/autogen 渲染引擎并清理 scope

## Description
将 `pkg/task/autogen.go` 的 `renderBody()` 从 `strings.ReplaceAll` 迁移到 `text/template.Execute()`，定义 `autogenTemplateData` struct。迁移 12 个非 record 模板的占位符语法。`{{SCOPE}}` 两种使用模式统一替换：段落级用 `{{if .SurfaceKey}}` 包裹，行内值替换为 `{{.SurfaceKey}}`。直接删除 `BodyContext.Scope` 字段。移除 `removeLineContaining()` 和 `removeSection()` 后处理函数。新增 `ValidateAutogenTemplates()` 校验。

## Reference Files
- `forge-cli/pkg/task/autogen.go`: renderBody() 迁移到 text/template，新增 autogenTemplateData struct，移除 removeLineContaining()/removeSection() (source: proposal.md#pkg/task/autogen-自动生成任务层)
- `forge-cli/pkg/task/data/test-gen-contracts.md`: {{SCOPE}} 段落级迁移 (source: proposal.md#pkg/task/autogen-自动生成任务层)
- `forge-cli/pkg/task/data/test-gen-journeys.md`: {{SCOPE}} 段落级迁移 (source: proposal.md#pkg/task/autogen-自动生成任务层)
- `forge-cli/pkg/task/data/doc-consolidate.md`: {{SCOPE}} 行内值迁移 (source: proposal.md#pkg/task/autogen-自动生成任务层)
- `forge-cli/pkg/task/autogen.go`: BodyContext.Scope 字段删除 (source: proposal.md#In-Scope)

## Acceptance Criteria
- [ ] `autogenTemplateData` struct 包含提案定义的全部字段，`renderBody()` 使用 `text/template.Execute()` 渲染
- [ ] 12 个非 record 模板文件占位符已迁移到 `{{.X}}` 格式，`pkg/task/data/` 中无 `{{SCOPE}}` 残留
- [ ] `BodyContext` struct 中 `Scope` 字段已删除，语义由 `SurfaceKey`/`SurfaceType` 覆盖
- [ ] `removeLineContaining()` 和 `removeSection()` 函数已从 `autogen.go` 中移除
- [ ] `ValidateAutogenTemplates()` 使用 `missingkey=error` + 零值 struct Execute 校验通过

## Implementation Notes
- `{{SCOPE}}` 两种使用模式的迁移需注意区分：段落级（test-gen-contracts, test-gen-journeys）用 `{{if .SurfaceKey}}` 整段包裹；行内值（doc-consolidate, test-run 等）直接替换为 `{{.SurfaceKey}}`
- `BodyContext.Scope` 是 `[]string` 类型，迁移后统一为 `SurfaceKey`（string），无 range 循环
- `Scope` 仅在 `autogen.go` 内部消费，调用链为 `BuildIndex()` → `renderBody()`，影响面可控

### Test Impact
- Affected test suite(s): `forge-cli/pkg/task/`, `forge-cli/internal/`
- Expected fixture changes: autogen golden files, BuildIndex test fixtures
- Risk level: medium
