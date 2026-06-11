---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["2"]
type: "doc.review"
scope: "all"
---

Review documentation quality for the remove-forge-test-command feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 2-clean-doc-references
- [ ] Full-text search for `forge test promote`, `forge test run-journey`, `forge test verify` returns zero results (excluding `docs/features/` history docs)
- [ ] No documentation file instructs users or agents to run `forge test` subcommands


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/remove-forge-test-command/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/remove-forge-test-command/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
