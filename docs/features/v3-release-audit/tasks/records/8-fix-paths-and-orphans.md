---
status: "completed"
started: "2026-05-25 00:25"
completed: "2026-05-25 00:35"
time_spent: "~10m"
---

# Task Record: 8 Fix cross-skill path violations and orphan rules

## Summary
Fixed cross-skill path violations (8 hardcoded paths replaced with descriptive references) and resolved 12 orphan rules files (6 true orphans added Load/reference directives, 5 parameterized surface rules annotated, 4 forge test run --tags corrected to framework-native tag filters, 1 deprecated rule annotated).

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/run-tests/SKILL.md
- plugins/forge/skills/run-tests/rules/env-check.md
- plugins/forge/skills/eval/SKILL.md
- plugins/forge/skills/eval/rules/pre-processing.md
- plugins/forge/skills/eval/rules/scorer-composition.md
- plugins/forge/skills/gen-contracts/SKILL.md
- plugins/forge/skills/gen-contracts/rules/risk-density.md
- plugins/forge/skills/gen-contracts/rules/journey-contract-model.md
- plugins/forge/skills/gen-journeys/SKILL.md
- plugins/forge/skills/gen-journeys/rules/journey-contract-model.md
- plugins/forge/skills/gen-test-scripts/SKILL.md

### Key Decisions
无

## Document Metrics
8 path violations fixed, 12 orphan rules resolved (6 Load added, 5 surface annotated, 4 incorrect commands corrected, 1 deprecated annotated)

## Referenced Documents
- docs/proposals/v3-release-audit/proposal.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] grep -r hardcoded/path/pattern plugins/forge/skills/run-tests/ returns 0
- [x] All rules/ files referenced by at least one SKILL.md (in-degree >= 1)
- [x] 6 true orphan rules have Load directives
- [x] 5 parameterized surface rules annotated with reference relationship

## Notes
Cross-skill paths (skills/gen-journeys/rules/surface-<type>.md) replaced with descriptive references per forge-distribution.md Section 5. freeform-injection.md kept as deprecated reference with annotation. test-isolation.md referenced from gen-test-scripts SKILL.md (not parent run-tests) as it is a code generation concern. 4 occurrences of nonexistent 'forge test run --tags' corrected to framework-native tag filters in journey-contract-model.md (both gen-contracts and gen-journeys copies).
