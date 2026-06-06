---
status: "completed"
started: "2026-06-06 13:21"
completed: "2026-06-06 13:22"
time_spent: "~1m"
---

# Task Record: 3 Update gen-test-scripts rule files

## Summary
Updated gen-test-scripts rule files to align with surface-key adaptive output directory structure. step-1-contract-loading.md Fact Table example paths updated to reflect multi-surface (tests/<surfaceKey>/<journey>/) and single-surface (tests/<journey>/) dual paths. run-to-learn.md skeleton test directory guidance updated to use adaptive wording. step-0.5-validation.md and quality-gates.md confirmed clean (no tests/ path references, no changes needed).

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-test-scripts/rules/step-1-contract-loading.md
- plugins/forge/skills/gen-test-scripts/rules/run-to-learn.md

### Key Decisions
无

## Document Metrics
2 of 4 rule files modified, 2 confirmed clean via grep scan, 0 new files

## Referenced Documents
- docs/proposals/surface-key-test-alignment/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] step-0.5-validation.md directory references consistent with SKILL.md output rules
- [x] step-1-contract-loading.md contract-to-directory mapping reflects surface-key directory structure
- [x] quality-gates.md test file path validation rules aligned
- [x] run-to-learn.md learning references use correct paths

## Notes
step-0.5-validation.md and quality-gates.md contained no tests/ path references (verified via grep), so no modifications were needed. Implementation Notes check also confirmed convention-guide.md has no tests/ references.
