---
id: "2"
title: "迁移 21 个 prompt 模板占位符语法并添加条件块"
priority: "P0"
estimated_time: "2h"
complexity: "medium"
dependencies: [1]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 2: 迁移 21 个 prompt 模板占位符语法并添加条件块

## Description
将 `pkg/prompt/data/` 下全部 21 个模板文件的占位符语法从 `{{X}}` 迁移到 `{{.X}}`（dot-notation）。为每个模板添加 `{{if}}` 条件块替换当前靠后处理删除的段落。具体条件块包括：PhaseSummary（18 个模板，两个独立位置）、CoverageStrategy、SurfaceKey、Complexity（4 个 coding 模板）。4 个 doc 模板补齐 SURFACE_KEY header 声明。移除 `<!-- IF NOT_LOW -->` 标记，替换为 `{{if ne .Complexity "low"}}` 条件块。

## Reference Files
- `forge-cli/pkg/prompt/data/*.md`: 21 个模板文件占位符语法 `{{X}}` → `{{.X}}` (source: proposal.md#pkg/prompt-渲染层)
- `forge-cli/pkg/prompt/data/coding-*.md`: 4 个 coding 模板添加 `{{if ne .Complexity "low"}}...{{end}}` 替换 `<!-- IF NOT_LOW -->` 标记 (source: proposal.md#In-Scope)
- `forge-cli/pkg/prompt/data/doc-*.md`: 4 个 doc 模板补齐 SURFACE_KEY header 声明 (source: proposal.md#In-Scope)

## Acceptance Criteria
- [ ] 21 个模板文件中所有占位符为 `{{.X}}` 格式，无 `{{PLACEHOLDER}}` 格式残留
- [ ] `{{if .PhaseSummary}}...{{end}}` 条件块在 18 个使用 PHASE_SUMMARY 的模板中正确包裹两处独立位置（标签行 + 条件指令行），无值时两处均消失
- [ ] `{{if .CoverageStrategy}}...{{end}}` 条件块正确处理三态：空（省略段落）/ 策略文本 / 特殊 "No coverage..." 指令
- [ ] 4 个 doc 模板包含 SURFACE_KEY header 声明；4 个 coding 模板中 `<!-- IF NOT_LOW -->...<!-- END_IF -->` 已替换为 `{{if ne .Complexity "low"}}...{{end}}`

## Implementation Notes
- PhaseSummary 在模板中出现在两个独立位置（标签行 + 条件指令行，相隔 10-30 行），需分别用两个独立的 `{{if .PhaseSummary}}` 块包裹，不可合并为单个块
- `just compile` 命令尾部空格处理：`just compile{{if .SurfaceKey}} {{.SurfaceKey}}{{end}}`，当 SurfaceKey 为空时无尾随空格
- CoverageStrategy 三态由 Go 代码解析，模板侧仅判断 `{{if .CoverageStrategy}}`
