---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["9"]
type: "doc.review"
scope: "all"
---

Review documentation quality for the test-recipe-unification feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 10-documentation
- CLI docs reference `unit-test`, `test` (not `e2e-test`)
- ARCHITECTURE.md describes FullGateSequence, UnitGateSequence, NonBreakingGateSequence
- quality-gate.md reflects new gate steps and two-layer model
- No residual `e2eTest` or `e2e-test` references in documentation (historical lessons/proposals excluded)


### 8-prompt-templates
- All 3 prompt templates reference `just unit-test` instead of `just test` for per-task gate scenarios
- No residual `just test` references in gate/fix/validation prompt contexts


### 9-skill-markdown
- All skill/command markdown references `unit-test`, `test`, `test-setup`, `probe` recipe names
- No residual `e2e-test`, `e2e-setup`, `e2e-verify` references
- `init-justfile/SKILL.md` Standard Target Contract reflects new recipe model with per-language/per-surface generation
- `run-tests` config schema examples use `test` key (not `e2eTest`)


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/test-recipe-unification/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/test-recipe-unification/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
