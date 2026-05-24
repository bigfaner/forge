---
status: "completed"
started: "2026-05-24 09:51"
completed: "2026-05-24 09:52"
time_spent: "~1m"
---

# Task Record: 4 Improve quick-tasks SKILL.md Reference Files generation

## Summary
Added Reference Files generation guidance to quick-tasks SKILL.md Step 2 with section-level precision, extraction logic, and format requirements

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/quick-tasks/SKILL.md

### Key Decisions
无

## Document Metrics
1 file modified, ~25 lines added in Step 2 subsection

## Referenced Documents
- docs/proposals/spec-authority-enforcement/proposal.md
- docs/conventions/forge-distribution.md
- plugins/forge/skills/quick-tasks/templates/task.md
- plugins/forge/skills/quick-tasks/templates/task-doc.md

## Review Status
completed

## Acceptance Criteria
- [x] SKILL.md Step 2 includes explicit instructions for generating Reference Files with section-level precision
- [x] Instructions require format: proposal.md#Section-Title — brief description
- [x] Instructions specify 2-5 specific sections from proposal.md relevant to each task
- [x] Instructions specify extraction logic: identify relevant proposal sections based on task description and affected files
- [x] If proposal.md references external design documents that exist on disk, include those sections as additional Reference Files
- [x] Each generated coding task must have >=1 precise section reference

## Notes
Added as a subsection within existing Step 2, between Scope Inference and Priority. Includes concrete example. Did not change overall SKILL.md structure per Hard Rules.
