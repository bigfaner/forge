---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["2"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the per-task-surface-scoped-gate feature (quick mode).

## Acceptance Criteria
- [ ] 所有 doc task 产出的文档符合 AC baseline 中的要求
- [ ] 文档无错别字、格式错误或断链

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 2-surface-rules-gate-recipes
- [ ] 5 个 surface rule 文件（api.md、web.md、cli.md、tui.md、mobile.md）均包含 `<key>-compile`/`<key>-fmt`/`<key>-lint`/`<key>-unit-test` 的 stub recipe 定义
- [ ] 5 个 surface rule 文件均包含 compile/fmt/lint/unit-test 的 Recipe Invocation Contract 条目


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/per-task-surface-scoped-gate/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/per-task-surface-scoped-gate/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
