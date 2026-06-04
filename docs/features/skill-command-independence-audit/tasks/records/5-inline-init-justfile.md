---
status: "completed"
started: "2026-06-04 00:44"
completed: "2026-06-04 00:46"
time_spent: "~2m"
---

# Task Record: 5 Inline test-type model into init-justfile + reduce examples

## Summary
Inlined test-type model mapping table into init-justfile/SKILL.md with origin traceability marker, reduced justfile examples (aggregate, surface, verification output), eliminated cross-skill reference to test-guide/references/test-type-model.md

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/init-justfile/SKILL.md

### Key Decisions
无

## Document Metrics
5-row mapping table inlined, 4 example blocks reduced (aggregate: 12->1 line, surface: 18->2 lines, verification: 12->2 lines, output: 55->10 lines), 0 cross-skill internal references remaining

## Referenced Documents
- docs/proposals/skill-command-independence-audit/proposal.md
- plugins/forge/skills/test-guide/references/test-type-model.md

## Review Status
final

## Acceptance Criteria
- [x] init-justfile/SKILL.md contains inlined test-type mapping table with INLINE:origin marker
- [x] justfile examples reduced, key patterns retained
- [x] init-justfile no longer contains cross-skill reference to test-guide internal file

## Notes
Inline marker uses <!-- INLINE:origin=test-guide/references/test-type-model.md --> with matching END marker. Remaining /forge:test-guide references are skill command invocations (user hints), not internal file references.
