---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["3.6"]
type: "doc.review"
scope: "all"
---

Review documentation quality for the surface-aware-justfile feature (breakdown mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 3.2-init-justfile-rules

- web.md 包含完整的编排序列表（4 步）、配方调用契约表（5 个配方含聚合配方）、journey 过滤策略表
- api.md 与 web.md 结构相同，probe 目标为 /healthz，支持 @api journey
- cli.md 和 tui.md 无 dev/probe 步骤，无聚合配方
- mobile.md 包含 test-setup 步骤（模拟器准备），支持 @mobile journey
- 每个文件的格式严格遵循 Interface 2 定义的 markdown 结构


### 3.4-run-tests-rules

- web/api 规则文件包含 4 步编排序列及每步的 exit code 语义（0/1/2）
- cli/tui 规则文件包含 2 步编排（无 dev/probe）
- mobile 规则文件包含 test-setup 前置步骤
- 每个文件定义 probe 失败 HARD-GATE 约束
- 规则文件可被 run-tests SKILL.md 直接加载消费


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/surface-aware-justfile/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/surface-aware-justfile/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
