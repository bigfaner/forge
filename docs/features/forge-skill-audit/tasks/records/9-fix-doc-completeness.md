---
status: "completed"
started: "2026-06-10 19:34"
completed: "2026-06-10 19:36"
time_spent: "~2m"
---

# Task Record: 9 Fix small documentation completeness issues (M-3, M-5, M-6, M-7)

## Summary
Fixed 4 documentation completeness issues: (M-3) breakdown-tasks intent read path now specifies full path docs/proposals/<slug>/proposal.md; (M-5) test-isolation.md has OWNER/CONSUMERS header comment; (M-6) brainstorm SKILL.md Step 5 includes {{AUTHOR}} assignment guidance; (M-7) write-prd manifest template uses {{SLUG}} instead of {{FEATURE_SLUG}}.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/run-tests/rules/test-isolation.md
- plugins/forge/skills/brainstorm/SKILL.md
- plugins/forge/skills/write-prd/templates/manifest.md

### Key Decisions
无

## Document Metrics
4 files modified, 4/4 AC passed, 0 FEATURE_SLUG residuals

## Referenced Documents
- docs/proposals/forge-skill-audit/proposal.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] breakdown-tasks SKILL.md intent read specifies full path docs/proposals/<slug>/proposal.md
- [x] test-isolation.md has OWNER: run-tests | CONSUMERS: gen-test-scripts (INLINE) comment
- [x] brainstorm SKILL.md Step 5 contains {{AUTHOR}} assignment guidance
- [x] write-prd/templates/manifest.md uses {{SLUG}} not {{FEATURE_SLUG}}

## Notes
M-7 regression verified: grep -r FEATURE_SLUG plugins/forge/skills/write-prd/templates/ returns no results. All changes are markdown-only per Hard Rules.
