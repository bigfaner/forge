---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["2", "3"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the task-sizing-gate feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 2-breakdown-tasks-audit-step
- [ ] `breakdown-tasks` 的 SKILL.md 在写 task 文件后、`forge task index` 前包含独立的 task sizing audit step
- [ ] audit step 指示 LLM 对每个 task 检查 multi-verb 和跨域 AC，发现问题自动拆分并输出报告
- [ ] 所有 `validate-index` 引用已更新为 `validate`
- [ ] 所有 step 编号连续无跳跃


### 3-quick-tasks-audit-step
- [ ] `quick-tasks` 的 SKILL.md 在写 task 文件后、`forge task index` 前包含独立的 task sizing audit step
- [ ] audit step 指示 LLM 对每个 task 检查 multi-verb 和跨域 AC，发现问题自动拆分并输出报告
- [ ] 所有 `validate-index` 引用已更新为 `validate`
- [ ] 所有 step 编号连续无跳跃


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/task-sizing-gate/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/task-sizing-gate/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
