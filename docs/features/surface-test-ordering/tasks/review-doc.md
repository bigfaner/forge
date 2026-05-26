---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["7"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the surface-test-ordering feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 7-update-gen-journeys-skill
- [ ] SKILL.md 包含多 surface 规则加载指导（按 surface type 分节组织）
- [ ] 输出格式要求：每个 Journey 标注覆盖的 surface type 集合（如 `[web, api]`）
- [ ] 所有配置的 surface type 至少被一个 Journey 覆盖


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/surface-test-ordering/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/surface-test-ordering/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
