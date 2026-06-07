---
status: "completed"
started: "2026-06-07 23:52"
completed: "2026-06-07 23:54"
time_spent: "~2m"
---

# Task Record: 3 Prompt templates gate recipe 从参数模式切换为前缀模式

## Summary
8 个 prompt template 文件中共 33 处 gate recipe 调用从参数模式 `just <recipe>{{if .SurfaceKey}} {{.SurfaceKey}}{{end}}` 切换为前缀模式 `just {{if .SurfaceKey}}{{.SurfaceKey}}-{{end}}<recipe>`，与 Task 1 的 RunGate() prefixed resolution 形成互补

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/templates/coding-feature.md
- forge-cli/pkg/prompt/templates/coding-enhancement.md
- forge-cli/pkg/prompt/templates/coding-fix.md
- forge-cli/pkg/prompt/templates/coding-cleanup.md
- forge-cli/pkg/prompt/templates/coding-refactor.md
- forge-cli/pkg/prompt/templates/gate.md
- forge-cli/pkg/prompt/templates/validation-code.md
- forge-cli/pkg/prompt/templates/fix-record-missed.md

### Key Decisions
无

## Document Metrics
8 files, 33 replacements, 0 regressions

## Referenced Documents
- docs/proposals/per-task-surface-scoped-gate/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] 8 个 prompt template 文件中所有 33 处 gate recipe 调用从参数模式切换为前缀模式
- [x] SurfaceKey 为空时渲染结果仍为 `just compile`（无前缀），与改动前一致
- [x] grep 验证所有模板中不再包含旧参数模式出现在 recipe 调用上下文中

## Notes
coding-refactor.md 含 10 处（6 增量 compile + 3 final gate + 1 lint 表格行），其余文件各 3-4 处。identity 段落的 SURFACE_KEY 声明未被误改。
