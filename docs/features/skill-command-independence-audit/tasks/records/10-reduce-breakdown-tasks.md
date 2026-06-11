---
status: "completed"
started: "2026-06-04 01:03"
completed: "2026-06-04 01:09"
time_spent: "~6m"
---

# Task Record: 10 Reduce shared content in breakdown-tasks

## Summary
Condensed breakdown-tasks SKILL.md shared content with quick-tasks (324→226 lines, 30% reduction). All breakdown-tasks-specific logic preserved.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/breakdown-tasks/SKILL.md

### Key Decisions
无

## Document Metrics
324→226 lines (30% reduction), HARD-RULE: 5 (unchanged), HARD-GATE: 2 (unchanged), Steps 0-9 intact, decision tables intact

## Referenced Documents
- docs/proposals/skill-command-independence-audit/proposal.md
- plugins/forge/skills/quick-tasks/SKILL.md

## Review Status
final

## Acceptance Criteria
- [x] Shared ~150 lines with quick-tasks condensed, breakdown-tasks-specific logic preserved
- [x] breakdown-tasks still fully guides AI agent task breakdown execution

## Notes
Condensed: Docs-Only Fast Path, Step 0, Split Rules, Complexity, File Scope Boundary, Test Impact Assessment, Surface-Key/Type Inference, Priority Assignment, Type Assignment table, Intent Propagation, Reference Files Generation, Test Tasks, Task Sizing Audit, Steps 6-9, Output Checklist. All HARD-RULE/HARD-GATE blocks preserved. All breakdown-tasks-unique content (Condition-Rule Matrix, Step 1 multi-artifact reading, Step 2 Element Mapping/PRD Coverage/Phase Detection, Step 3 Phases & Dependencies, Step 4 phase.sub naming, 4a Business Tasks, User Stories, Step 8 Manifest traceability) unchanged.
