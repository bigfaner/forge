---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["4"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the init-justfile-slim feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 4-simplify-skill-and-delete-rules
- [ ] SKILL.md + 保留的 rules（self-correction.md）总行数 ≤ 280 行
- [ ] SKILL.md 引用 `forge justfile scaffold` CLI 命令替代手动模板生成，Agent 工作流更新为提案的 Step 0-5 新流程
- [ ] 删除 6 个文件：`rules/server-lifecycle.md`、`rules/surfaces/api.md`、`rules/surfaces/cli.md`、`rules/surfaces/mobile.md`、`rules/surfaces/tui.md`、`rules/surfaces/web.md`
- [ ] `rules/self-correction.md`（34 行）保留不动
- [ ] Convention Cold Start Fallback 策略以 5-10 行摘要保留在 SKILL.md 中


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/init-justfile-slim/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/init-justfile-slim/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.

## Acceptance Criteria

- [ ] All acceptance criteria met
