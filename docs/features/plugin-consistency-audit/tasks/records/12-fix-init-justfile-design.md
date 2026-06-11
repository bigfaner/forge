---
status: "completed"
started: "2026-05-30 06:03"
completed: "2026-05-30 06:05"
time_spent: "~2m"
---

# Task Record: 12 Fix: resolve init-justfile template vs SKILL.md design contradiction

## Summary
Resolved init-justfile SKILL.md design contradiction: updated Step 0 HARD-RULE and Step 3a to define templates as active starting-point components with a three-layer generation process (template load -> Convention override -> LLM customization)

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/init-justfile/SKILL.md

### Key Decisions
无

## Document Metrics
4 sections updated: Step 0 HARD-RULE, Step 3a, Step 3b recipe content generation, EXTREMELY-IMPORTANT

## Referenced Documents
- docs/features/plugin-consistency-audit/reports/04-skills-batch-c.md
- plugins/forge/skills/init-justfile/templates/generic.just
- plugins/forge/skills/init-justfile/templates/go.just
- plugins/forge/skills/init-justfile/templates/node.just
- plugins/forge/skills/init-justfile/templates/python.just
- plugins/forge/skills/init-justfile/templates/rust.just
- plugins/forge/skills/init-justfile/templates/mixed.just
- plugins/forge/skills/init-justfile/rules/surfaces/web.md

## Review Status
final

## Acceptance Criteria
- [x] SKILL.md Step 0 HARD-RULE consistent with templates/ directory content
- [x] SKILL.md explicitly states template purpose (starting point / reference / not used)
- [x] Step 3a references templates with specific file-to-language mapping
- [x] Template extras (run, dev, test, test-setup, probe) documented as convenience targets

## Notes
Templates confirmed as active components (modern boundary markers, surface-aware recipes, surface rules explicitly reference them). Decision: templates serve as structural starting points; LLM customizes recipe bodies. No template files were deleted or moved. This decision directs Task 7 (go.just fix) to fix the template content rather than remove it.
