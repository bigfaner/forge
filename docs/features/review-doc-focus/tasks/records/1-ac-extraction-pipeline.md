---
status: "completed"
started: "2026-05-24 19:11"
completed: "2026-05-24 19:17"
time_spent: "~6m"
---

# Task Record: 1 AC extraction pipeline (BodyContext + extract.go + renderBody)

## Summary
Built AC extraction pipeline: added DocTaskCriteria map[string]string field to BodyContext, implemented extractDocTaskCriteria() in extract.go to scan doc task files and extract ## Acceptance Criteria sections (respecting fenced code blocks), added {{DOC_TASK_AC}} placeholder support in renderBody() with sorted markdown sub-section serialization, and wired BuildIndex() to populate BodyContext.DocTaskCriteria when review-doc is needed.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/autogen.go
- forge-cli/pkg/task/extract.go
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/extract_test.go
- forge-cli/pkg/task/autogen_test.go

### Key Decisions
- Used strings.ReplaceAll for DOC_TASK_AC placeholder (consistent with existing 12+ templates, per Hard Rules)
- Named field DocTaskCriteria to avoid collision with existing AcceptanceCriteria []string (per Hard Rules)
- extractACSection tracks fenced code block state to avoid treating ## inside code blocks as section boundaries
- extractDocTaskCriteria filters by CategoryForType == CategoryDoc and excludes system types
- serializeDocTaskAC sorts keys alphabetically for deterministic output

## Test Results
- **Tests Executed**: Yes
- **Passed**: 14
- **Failed**: 0
- **Coverage**: 87.7%

## Acceptance Criteria
- [x] BodyContext struct 新增 DocTaskCriteria map[string]string 字段
- [x] extract.go 新增 extractDocTaskCriteria(taskDir string) map[string]string 函数
- [x] 提取算法：逐行扫描 .md 文件，找到 ## Acceptance Criteria 后收集至下一个 ## 行之前的所有内容
- [x] renderBody() 支持 {{DOC_TASK_AC}} 占位符，将 map 序列化为 markdown 后通过 strings.ReplaceAll 注入
- [x] BuildIndex() 在生成 review-doc 任务时调用提取函数并填充 BodyContext
- [x] 渲染策略使用 strings.ReplaceAll（非 Go text/template），保持与现有模板一致

## Notes
无
