---
status: "completed"
started: "2026-05-24 22:28"
completed: "2026-05-24 22:34"
time_spent: "~6m"
---

# Task Record: T-review-doc Review Documentation Quality

## Summary
Reviewed documentation quality for test-recipe-unification feature. Fixed 4 docs: ARCHITECTURE.md (updated test lifecycle to two-layer model, added FullGateSequence/UnitGateSequence/NonBreakingGateSequence, replaced e2e-test references), quality-gate.md (reflected new gate steps and two-layer model), profile-authoring.md (replaced e2e-test/e2e-setup/e2e-verify with unit-test/test/test-setup/probe), plan/task-index-command.md (replaced run-e2e-tests with run-tests). All 3 prompt templates (gate.md, fix-record-missed.md, validation-code.md) already reference just unit-test. All skill markdown already references new recipe names. Residual e2e references in forge-cli/pkg/prompt/data/test-verify-regression.md and test-run.md are out of docs/ scope.

## Changes

### Files Created
无

### Files Modified
- docs/ARCHITECTURE.md
- docs/business-rules/quality-gate.md
- docs/profile-authoring.md
- docs/plan/task-index-command.md

### Key Decisions
无

## Document Metrics
AC-10: 4/4 pass, AC-8: 2/2 pass (docs scope), AC-9: 4/4 pass. Fixed 4 files, 0 blocked items within docs/ scope.

## Referenced Documents
- docs/ARCHITECTURE.md
- docs/business-rules/quality-gate.md
- docs/profile-authoring.md
- docs/plan/task-index-command.md
- docs/proposals/test-recipe-unification/proposal.md
- forge-cli/pkg/prompt/data/gate.md
- forge-cli/pkg/prompt/data/coding-fix.md
- forge-cli/pkg/prompt/data/coding-feature.md
- forge-cli/pkg/prompt/data/fix-record-missed.md
- forge-cli/pkg/prompt/data/validation-code.md
- plugins/forge/skills/init-justfile/SKILL.md
- plugins/forge/skills/run-tests/SKILL.md
- plugins/forge/commands/fix-bug.md

## Review Status
fixes-applied

## Acceptance Criteria
- [x] CLI docs reference unit-test, test (not e2e-test)
- [x] ARCHITECTURE.md describes FullGateSequence, UnitGateSequence, NonBreakingGateSequence
- [x] quality-gate.md reflects new gate steps and two-layer model
- [x] No residual e2eTest/e2e-test references in documentation (historical excluded)
- [x] All 3 prompt templates reference just unit-test for per-task gate
- [x] No residual just test references in gate/fix/validation prompt contexts
- [x] All skill/command markdown references new recipe names
- [x] No residual e2e-test/e2e-setup/e2e-verify references in skill markdown
- [x] init-justfile/SKILL.md Standard Target Contract reflects new recipe model
- [x] run-tests config schema examples use test key (not e2eTest)

## Notes
Residual e2e references found in forge-cli/pkg/prompt/data/test-verify-regression.md (just test-e2e) and test-run.md (run-e2e-tests skill name). These are in code directory, outside docs/ scope constraint. Should be addressed in a follow-up coding task if needed.
