---
status: "completed"
started: "2026-05-27 01:11"
completed: "2026-05-27 01:13"
time_spent: "~2m"
---

# Task Record: 6 Clean up ghost fields and stale references in skill docs

## Summary
Cleaned up ghost fields and stale references across 12 skill doc files: deleted deprecated scope-assignment.md, fixed interfaces→surfaces, SCOPE→SURFACE_KEY/SURFACE_TYPE, decision-logging.md→decision-entry.md, scope→surface-type, removed undefined macros and test-template-dir reference

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/quick-tasks/SKILL.md
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/gen-test-scripts/SKILL.md
- plugins/forge/commands/run-tasks.md
- plugins/forge/commands/execute-task.md
- plugins/forge/commands/fix-bug.md
- plugins/forge/skills/write-prd/rules/knowledge-extraction.md
- plugins/forge/skills/tech-design/rules/knowledge-extraction.md
- plugins/forge/skills/breakdown-tasks/rules/db-schema.md
- plugins/forge/skills/breakdown-tasks/rules/existing-code-split.md
- plugins/forge/skills/breakdown-tasks/rules/ui-placement.md

### Key Decisions
无

## Document Metrics
1 file deleted, 11 files modified, 7 categories of stale references fixed

## Referenced Documents
- docs/proposals/pipeline-spec-code-alignment/proposal.md

## Review Status
completed

## Acceptance Criteria
- [x] scope-assignment.md deleted
- [x] No skill doc references interfaces config field
- [x] No skill doc extracts SCOPE from claim output
- [x] No skill doc references decision-logging.md
- [x] No skill doc references test-template-dir config field
- [x] db-schema.md uses surface-type fields instead of scope
- [x] existing-code-split.md references surface inference, not scope assignment

## Notes
scope-to-surface-key.md forge version reference verified correct (command exists and works). consolidate-specs/rules/overlap-detection.md references 'decision-logging protocol' generically, not a file reference — out of scope for this task.
