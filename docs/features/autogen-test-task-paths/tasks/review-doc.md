---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["1", "2"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the autogen-test-task-paths feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-add-feature-paths-embed
- [ ] 6 个 embed 模板均包含 `## Feature Paths` 区域，含 journeys (`ls docs/features/{{.FeatureSlug}}/testing/`) 和 contracts (`ls docs/features/{{.FeatureSlug}}/testing/<journey>/contracts/`) 两个 discovery 命令
- [ ] 富模板（test-gen-journeys、test-gen-contracts）若已有等价路径引用则不重复添加
- [ ] `go build ./...` 和 `go test ./...` 通过


### 2-add-feature-slug-prompt
- [ ] 6 个 prompt 模板均输出 `FEATURE_SLUG: {{.FeatureSlug}}` 行，位于 `TASK_FILE` 行之后
- [ ] `go build ./...` 和 `go test ./...` 通过


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/autogen-test-task-paths/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/autogen-test-task-paths/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
