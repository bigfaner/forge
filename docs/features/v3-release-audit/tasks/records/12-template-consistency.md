---
status: "completed"
started: "2026-05-25 00:52"
completed: "2026-05-25 00:59"
time_spent: "~7m"
---

# Task Record: 12 Fix template variable naming and frontmatter consistency

## Summary
Unified all template variable naming to {{VAR}} format and standardized frontmatter across 20+ template files

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/consolidate-specs/templates/biz-specs.md
- plugins/forge/skills/consolidate-specs/templates/tech-specs.md
- plugins/forge/skills/consolidate-specs/templates/markers.md
- plugins/forge/skills/consolidate-specs/templates/review-choices.md
- plugins/forge/skills/consolidate-specs/templates/vocabulary-index.md
- plugins/forge/skills/deep-research/templates/research-report.md
- plugins/forge/skills/forensic/templates/report.md
- plugins/forge/skills/learn/templates/convention-entry.md
- plugins/forge/skills/learn/templates/decision-entry.md
- plugins/forge/skills/learn/templates/lesson-entry.md
- plugins/forge/skills/tech-design/templates/api-handbook.md
- plugins/forge/skills/tech-design/templates/tech-design.md
- plugins/forge/skills/tech-design/templates/er-diagram.md
- plugins/forge/skills/tech-design/templates/decision-entry.md
- plugins/forge/skills/brainstorm/templates/proposal.md
- plugins/forge/skills/ui-design/templates/ui-design.md
- plugins/forge/skills/ui-design/templates/prototype.md
- plugins/forge/skills/quick-tasks/templates/task-doc.md
- plugins/forge/skills/quick-tasks/templates/task.md
- plugins/forge/skills/breakdown-tasks/templates/manifest-update-tasks.md
- plugins/forge/skills/gen-contracts/templates/contract.md
- plugins/forge/skills/write-prd/templates/prd-spec.md

### Key Decisions
无

## Document Metrics
22 files modified, ~120 variable substitutions, 3 format types unified to 1

## Referenced Documents
- docs/proposals/v3-release-audit/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] All template variables use unified format
- [x] Template frontmatter fields complete and format consistent

## Notes
Unified 3 variable formats (<var>, ${var}, YYYY-MM-DD) to {{VAR}}. Frontmatter values now consistently quote variables. ${panel} in prototype.md left unchanged as it is JavaScript template literal code, not a Forge template variable.
