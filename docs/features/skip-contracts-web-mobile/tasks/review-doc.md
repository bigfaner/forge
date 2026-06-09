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

Review documentation quality for the skip-contracts-web-mobile feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 2-skill-direct-path-web-mobile
- [ ] SC-5: 直达路径生成的脚本包含与 journey 步骤对应的用户动作调用（click/type/navigate）和至少一个可视化断言（非空、非骨架）
- [ ] SC-6: 按 surface type 自检覆盖率（count journeys_of_type == count test-scripts_of_type），缺口或类型不匹配时 FAIL 并输出缺口列表
- [ ] SC-8: types/web.md 和 types/mobile.md 直达规则能产出包含有意义断言的测试脚本


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/skip-contracts-web-mobile/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/skip-contracts-web-mobile/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.

## Acceptance Criteria

- [ ] All acceptance criteria met
