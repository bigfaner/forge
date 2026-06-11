---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["1"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the eval-diagnostic-mode feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-update-eval-templates
- [ ] `eval-journey.md` calls `/eval-journey` (not `/eval --type journey`) and contains zero references to hardcoded `850`
- [ ] `eval-contract.md` calls `/eval-contract` (not `/eval --type contract`) and contains zero references to hardcoded `850`
- [ ] Both templates' AC sections require only eval report generation, not score thresholds
- [ ] `grep -r ">= 850" forge-cli/pkg/task/templates/eval-*.md` returns zero results


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/eval-diagnostic-mode/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/eval-diagnostic-mode/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.

## Acceptance Criteria
- [ ] Documentation review completed for all deliverable documents
