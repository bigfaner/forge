---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["3"]
type: "doc.review"
scope: "all"
---

Review documentation quality for the run-tests-journey-isolation feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-fix-prompt-and-journey-isolation

- [ ] `forge-cli/pkg/prompt/data/test-run.md` references `forge:run-tests` (not `forge:run-e2e-tests`)
- [ ] SKILL.md includes a journey discovery step that runs `ls docs/features/<slug>/testing/` before test execution
- [ ] SKILL.md specifies per-journey execution: `just test <journey>` for each discovered journey directory
- [ ] SKILL.md specifies dev/probe execute once, per-journey loop for test, teardown once (for web/api/mobile)
- [ ] SKILL.md handles the "no journey" edge case: `docs/features/<slug>/testing/` missing or empty → error message suggesting run gen-journeys first
- [ ] Surface rule files updated to reflect per-journey test execution pattern


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/run-tests-journey-isolation/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/run-tests-journey-isolation/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
