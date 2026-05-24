---
status: "completed"
started: "2026-05-24 22:21"
completed: "2026-05-24 22:27"
time_spent: "~6m"
---

# Task Record: 9 Update skill and command markdown for new recipe model

## Summary
Updated 6 skill/command markdown files to reflect new two-layer test recipe model (unit-test/test/test-setup/probe), removed all e2e-test/e2e-setup/e2e-verify references

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/fix-bug.md
- plugins/forge/skills/clean-code/SKILL.md
- plugins/forge/skills/gen-test-scripts/rules/run-to-learn.md
- plugins/forge/skills/init-justfile/SKILL.md
- plugins/forge/skills/run-tests/SKILL.md
- plugins/forge/skills/run-tests/references/config-schema.md

### Key Decisions
无

## Document Metrics
6 files modified, 0 e2e-test/e2e-setup/e2e-verify residuals

## Referenced Documents
- docs/proposals/test-recipe-unification/proposal.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] All skill/command markdown references unit-test, test, test-setup, probe recipe names
- [x] No residual e2e-test, e2e-setup, e2e-verify references
- [x] init-justfile/SKILL.md Standard Target Contract reflects new recipe model with per-language/per-surface generation
- [x] run-tests config schema examples use test key (not e2eTest)

## Notes
Major rewrite of init-justfile/SKILL.md Standard Target Contract (Step 3 and output examples). Minor text replacements in fix-bug.md, clean-code/SKILL.md, run-to-learn.md. Config schema examples in run-tests fully updated.
