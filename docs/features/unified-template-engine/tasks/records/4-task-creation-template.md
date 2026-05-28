---
status: "completed"
started: "2026-05-28 01:10"
completed: "2026-05-28 01:22"
time_spent: "~12m"
---

# Task Record: 4 迁移 pkg/template 任务创建层到 text/template

## Summary
Migrated pkg/template task creation layer from strings.ReplaceAll to text/template.Execute(). Added TaskTemplateData struct with 11 fields. Migrated 2 template files (coding.fix.md, coding.cleanup.md) from {{X}} to {{.X}} placeholder syntax with conditional surface rendering. Removed ApplyVars(), injectSurfaceFrontmatter(), and placeholderRe from add.go.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/template/template.go
- forge-cli/pkg/template/data/coding.fix.md
- forge-cli/pkg/template/data/coding.cleanup.md
- forge-cli/pkg/task/add.go
- forge-cli/pkg/task/add_test.go
- forge-cli/pkg/template/template_test.go

### Key Decisions
- TaskTemplateData exported (not unexported) because it is consumed across packages (template -> task)
- SOURCE_TASK_ID resolved from opts.Vars when opts.SourceTaskID is empty, preserving quality-gate call path compatibility
- Surface Inference section uses {{if not .SurfaceKey}} to render only when surface is empty

## Test Results
- **Tests Executed**: Yes
- **Passed**: 146
- **Failed**: 0
- **Coverage**: 87.6%

## Acceptance Criteria
- [x] TaskTemplateData struct contains the 5 fields defined in proposal (ID/Title mapped to TaskName, SurfaceKey, SurfaceType, Description mapped to TaskGoal, ScopeDescription)
- [x] CreateTaskMarkdown() uses text/template.Execute() rendering, ApplyVars() and injectSurfaceFrontmatter() removed
- [x] 2 task creation templates migrated to {{.X}} format, surface rendered conditionally and Surface Inference section omitted when surface has value

## Notes
pkg/template coverage dropped to 60.7% because TestApplyVars and TestInjectSurfaceFrontmatter were removed (tested removed functions). Template rendering is covered by TestCreateTaskMarkdown_TemplateMode in pkg/task.
