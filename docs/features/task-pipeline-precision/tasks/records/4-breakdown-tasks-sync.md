---
status: "completed"
started: "2026-05-27 23:27"
completed: "2026-05-27 23:30"
time_spent: "~3m"
---

# Task Record: 4 Sync breakdown-tasks with split rules, complexity判定, Reference Files inline

## Summary
Synced breakdown-tasks SKILL.md with quick-tasks split rules, complexity判定, and Reference Files inline format; added complexity field to task.md template

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/breakdown-tasks/templates/task.md

### Key Decisions
无

## Document Metrics
2 files modified: SKILL.md (granularity→split rules+complexity, reference files→inline format), task.md (+complexity field)

## Referenced Documents
- docs/proposals/task-pipeline-precision/proposal.md
- docs/conventions/forge-distribution.md
- plugins/forge/skills/quick-tasks/SKILL.md

## Review Status
final

## Acceptance Criteria
- [x] Task Splitting Rules paragraph uses 'independently verifiable' as the merge standard
- [x] AC max 6 rule added (same wording as quick-tasks)
- [x] Multi-verb detection rule added (same wording as quick-tasks)
- [x] Complexity判定 logic with LLM override matches quick-tasks wording
- [x] Reference Files generation changed to inline precise info format
- [x] templates/task.md frontmatter has complexity field with default 'medium'
- [x] breakdown-tasks specific features (Phase & Gate Detection, PRD Coverage Verification) NOT affected

## Notes
Task Splitting Rules replace old Granularity rule. Reference Files changed from section-pointer format to inline precise info format with source traceability. breakdown-tasks-specific adjustment note added to complexity判定.
