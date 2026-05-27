---
status: "completed"
started: "2026-05-27 01:14"
completed: "2026-05-27 01:16"
time_spent: "~2m"
---

# Task Record: 7 Fix surface and type consistency in docs

## Summary
Fix surface and type consistency across 6 documentation files: two-layer surface resolution strategy, task-doc.md exemption note, webui->web, doc.eval->doc.review

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/quick-tasks/SKILL.md
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/breakdown-tasks/rules/scope-to-surface-key.md
- plugins/forge/skills/quick-tasks/templates/task-doc.md
- plugins/forge/skills/gen-test-scripts/rules/step-0.5-validation.md
- plugins/forge/skills/submit-task/data/record-format-doc.md

### Key Decisions
无

## Document Metrics
6 files modified, 5 AC items passed

## Referenced Documents
- docs/proposals/pipeline-spec-code-alignment/proposal.md

## Review Status
completed

## Acceptance Criteria
- [x] Surface resolution docs describe two-layer strategy (project-level shortcut + file-level query)
- [x] task-doc.md either has surface fields or documents why they're absent
- [x] No reference to webui surface type — all use canonical web
- [x] record-format-doc.md lists doc.review and does not list doc.eval
- [x] Single-surface project surface-type is non-empty (not left blank as placeholder)

## Notes
All changes are doc-only per Hard Rules. No Go code behavior changed. Two-layer strategy eliminates N*M forge surfaces calls in single-surface projects.
