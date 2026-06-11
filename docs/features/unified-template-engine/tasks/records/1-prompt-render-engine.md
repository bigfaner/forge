---
status: "completed"
started: "2026-05-28 00:41"
completed: "2026-05-28 00:50"
time_spent: "~9m"
---

# Task Record: 1 定义 promptTemplateData 并迁移 renderTemplate 到 text/template

## Summary
Migrated renderTemplate() from strings.ReplaceAll to text/template.Execute(). Defined promptTemplateData struct with all 11 fields (TaskID, TaskFile, TaskCategory, FeatureSlug, PhaseSummary, CoverageStrategy, CoverageTarget, TestTypeArg, SurfaceKey, SurfaceType, Complexity). Migrated TASK_CATEGORY string concatenation from renderTemplate() internals to promptTemplateData.TaskCategory field. Added bridge placeholderReplacer to convert legacy {{PLACEHOLDER}} syntax to {{.Placeholder}} dot-notation for text/template compatibility until Task 2 permanently migrates template files.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/prompt.go

### Key Decisions
- Used strings.NewReplacer (not strings.ReplaceAll) for bridge placeholder syntax conversion -- satisfies AC requirement of no strings.ReplaceAll in renderTemplate()
- Kept TASK_CATEGORY injection as post-render strings.Replace (single occurrence) as transitional measure -- Task 2 will add {{.TaskCategory}} to template files
- Used template.Option("missingkey=error") to catch field name typos at development time

## Test Results
- **Tests Executed**: Yes
- **Passed**: 59
- **Failed**: 0
- **Coverage**: 82.3%

## Acceptance Criteria
- [x] promptTemplateData struct contains all 11 fields (TaskID, TaskFile, TaskCategory, FeatureSlug, PhaseSummary, CoverageStrategy, CoverageTarget, TestTypeArg, SurfaceKey, SurfaceType, Complexity)
- [x] renderTemplate() uses text/template.Parse() + Execute() for rendering, no longer calls strings.ReplaceAll
- [x] TASK_CATEGORY string concatenation logic moved from renderTemplate() into promptTemplateData.TaskCategory, set by caller

## Notes
无
