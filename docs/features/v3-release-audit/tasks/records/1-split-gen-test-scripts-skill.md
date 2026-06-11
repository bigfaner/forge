---
status: "completed"
started: "2026-05-24 23:56"
completed: "2026-05-25 00:02"
time_spent: "~6m"
---

# Task Record: 1 Split gen-test-scripts SKILL.md (527→≤350 lines)

## Summary
Split gen-test-scripts SKILL.md from 527 to 331 lines by extracting Step 0.5 (Surface Detection) and Step 1 (Code Reconnaissance) into rules/ files, plus moving Step 3.0 Surface-Driven Generation Strategy into the surface rules file to meet the 350-line constraint

## Changes

### Files Created
- plugins/forge/skills/gen-test-scripts/rules/step-0.5-validation.md
- plugins/forge/skills/gen-test-scripts/rules/step-1-contract-loading.md

### Files Modified
- plugins/forge/skills/gen-test-scripts/SKILL.md

### Key Decisions
无

## Document Metrics
SKILL.md: 527->331 lines (-37%); 2 new rules files (95+72=167 lines); all AC met

## Referenced Documents
- docs/proposals/v3-release-audit/proposal.md
- docs/conventions/skill-self-containment.md
- docs/conventions/skill-structure.md

## Review Status
final

## Acceptance Criteria
- [x] SKILL.md <= 350 lines
- [x] New rules files referenced by SKILL.md via Load directives (in-degree >= 1)
- [x] Split SKILL.md flow is complete with no broken references
- [x] wc -l SKILL.md <= 350

## Notes
Step 3.0 Surface-Driven Generation Strategy was also moved into rules/step-0.5-validation.md because extracting only Step 0.5+Step 1 (108 lines) would leave SKILL.md at 419 lines, still exceeding 350. This additional extraction is consistent with skill-structure.md guidance (rules definitions >5 lines belong in rules/).
