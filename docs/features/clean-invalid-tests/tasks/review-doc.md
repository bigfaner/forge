---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["5"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the clean-invalid-tests feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 5-extend-tc004-scope
- [ ] Contract file updated to include `forge-cli/tests/` in TC-004 scope
- [ ] "Zero unconditional t.Skip() calls" assertion explicitly lists both `tests/` and `forge-cli/tests/` as target directories


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/clean-invalid-tests/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/clean-invalid-tests/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
