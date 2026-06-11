---
status: "completed"
started: "2026-05-27 23:22"
completed: "2026-05-27 23:26"
time_spent: "~4m"
---

# Task Record: 3 Update quick-tasks split rules, complexity判定, Reference Files inline, remove task cap

## Summary
Updated quick-tasks SKILL.md with split rules (independently verifiable standard), AC max 6 rule, multi-verb detection, complexity判定 logic, inline Reference Files format, and removed 15 coding task cap. Added complexity field to task.md template. Removed 15 task cap from quick.md.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/quick-tasks/SKILL.md
- plugins/forge/skills/quick-tasks/templates/task.md
- plugins/forge/commands/quick.md

### Key Decisions
无

## Document Metrics
3 files modified, 7 spec-code changes applied

## Referenced Documents
- docs/proposals/task-pipeline-precision/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] Merge rule changed from time estimation to independently verifiable standard
- [x] AC max 6 rule added
- [x] Multi-verb detection rule added
- [x] Complexity判定 logic with default heuristic + LLM override
- [x] Reference Files changed to inline precise info format
- [x] 15 coding task cap removed from SKILL.md HARD-GATE
- [x] templates/task.md frontmatter has complexity field with default medium
- [x] quick.md 15 task cap removed from HARD-GATE

## Notes
All 8 AC items pass. Split rules section uses 3-tier priority: independently verifiable standard, multi-verb detection, AC ceiling. Complexity heuristic: low (AC<=3 AND no Hard Rules AND RefFiles<=1), high (AC>6 OR has Hard Rules), medium (else). Reference Files now use inline file-path + change-description format with source traceability.
