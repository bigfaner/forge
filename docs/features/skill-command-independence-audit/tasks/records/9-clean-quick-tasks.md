---
status: "completed"
started: "2026-06-04 00:58"
completed: "2026-06-04 01:03"
time_spent: "~5m"
---

# Task Record: 9 Delete Integration + reduce shared content in quick-tasks

## Summary
Deleted ## Integration section and condensed ~150 lines of shared content with breakdown-tasks in quick-tasks SKILL.md (347→198 lines, 43% reduction). All 7 HARD-RULE/HARD-GATE blocks preserved. Reference Files Generation section retained as template placeholder replacement rules.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/quick-tasks/SKILL.md

### Key Decisions
无

## Document Metrics
347→198 lines (-43%); 7/7 hard-rule blocks preserved; 0 cross-skill references

## Referenced Documents
- docs/proposals/skill-command-independence-audit/proposal.md
- plugins/forge/skills/breakdown-tasks/SKILL.md

## Review Status
final

## Acceptance Criteria
- [x] ## Integration section deleted
- [x] ## Reference Files section preserved (template placeholder replacement rules)
- [x] ~150 lines of shared content with breakdown-tasks condensed, quick-tasks-specific logic retained

## Notes
Condensed shared sections inline: Docs-Only Fast Path (removed mermaid diagram), Step 0 Resolve Language (single-line summary), Split Rules/Complexity判定 (compact format), Surface-Key/Type Inference (single paragraph), Reference Files Generation (removed example block), Type Assignment (compressed table descriptions), Intent Propagation (inline table), File Scope Boundary (condensed), Breaking Task Test Impact Assessment (condensed), Step 4 Task Sizing Audit (condensed protocol), Steps 5-9 and Output Checklist (all trimmed). Commit HARD-RULE tag restored as wrapping tag.
