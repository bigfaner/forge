---
status: "completed"
started: "2026-05-25 00:13"
completed: "2026-05-25 00:15"
time_spent: "~2m"
---

# Task Record: 5 Fix broken CLI cross-references (5 locations)

## Summary
Fixed 4 broken CLI cross-references: replaced all 'forge config get surface' with 'forge surfaces' across 4 files in run-tests and gen-test-scripts skills. Confirmed 'test.execution' references are all config-path descriptors (not CLI commands) and require no changes.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/run-tests/rules/env-check.md
- plugins/forge/skills/run-tests/SKILL.md
- plugins/forge/skills/gen-test-scripts/rules/run-to-learn.md
- plugins/forge/skills/gen-test-scripts/rules/step-0.5-validation.md

### Key Decisions
无

## Document Metrics
4 files modified, 4 string replacements, 0 broken references remaining

## Referenced Documents
- docs/features/v3-release-audit/tasks/5-fix-cli-references.md

## Review Status
completed

## Acceptance Criteria
- [x] grep 'forge config get surface' plugins/forge/ returns 0 results
- [x] grep 'test\.execution' plugins/forge/ returns only config-schema definitional references
- [x] All replacement commands consistent with forge --help output

## Notes
Task listed 5 files but gen-test-scripts/SKILL.md did not contain 'forge config get surface'. The 'test.execution' references in run-tests/SKILL.md and init-justfile/SKILL.md are all config field path descriptors (e.g. test.execution.run), not CLI commands -- no modification needed. Task file's Affected Files table listed gen-test-scripts/SKILL.md and init-justfile/SKILL.md but these did not contain the broken patterns. Additionally found step-0.5-validation.md which was not listed but did contain 'forge config get surface'.
