---
id: "2.gate"
title: "Phase 2 Gate: Command Reorganization Verification"
priority: "P0"
estimated_time: "1h"
dependencies: ["2.summary"]
breaking: true
type: "gate"
mainSession: false
---

# 2.gate: Phase 2 Gate: Command Reorganization Verification

## Description

Exit verification gate for Phase 2. Confirms that command group structure, renames, new commands, and behavioral additions are correctly implemented and match the design specification.

## Verification Checklist

1. [ ] `forge --help` shows exactly 10 visible entries (5 groups + 5 top-level, version hidden)
2. [ ] `forge task --help` shows exactly 10 subcommands (not 11 — verify-task-done is top-level)
3. [ ] `forge e2e --help` shows exactly 6 subcommands (including validate-specs)
4. [ ] `forge version` works but is hidden from `--help`
5. [ ] Unknown command suggestions work: `forge taks` → suggests "task"
6. [ ] All renamed commands work with new names (submit, check-deps, validate-index, verify-task-done, quality-gate)
7. [ ] `forge task list-types` outputs 11 types with descriptions
8. [ ] Quality-gate cap logic: 3 active fix-tasks → cap reached error
9. [ ] Concurrent write locking works for submit
10. [ ] `template` command removed — `forge template` returns "unknown command"
11. [ ] `go build ./...` compiles without errors
12. [ ] All existing tests pass (with updated command names)
13. [ ] No deviations from design spec (or deviations are documented as decisions)

## Reference Files

- `design/tech-design.md` — All interface sections, PRD Divergences, Resolved Design Decisions
- Phase 2 task records — `records/2.*.md`
- Phase 2 summary — `records/2-summary.md`

## Acceptance Criteria

- [ ] All applicable verification checklist items pass
- [ ] Any deviations from design are documented as decisions in the record
- [ ] Record created via `/record-task` with test evidence

## Implementation Notes

This is a verification-only task. No new feature code should be written.
If issues are found:
1. Fix inline if trivial (e.g., wrong Short description)
2. Document non-trivial issues as decisions in the record
3. Set status to `blocked` if a blocking issue cannot be resolved
