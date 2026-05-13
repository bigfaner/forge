---
id: "4.gate"
title: "Phase 4 Gate: Reference Update Verification"
priority: "P0"
estimated_time: "1h"
dependencies: ["4.summary"]
breaking: true
type: "gate"
mainSession: false
---

# 4.gate: Phase 4 Gate: Reference Update Verification

## Description

Final gate. Verifies that all old `task <command>` references have been replaced with `forge` equivalents across hooks, skills, agents, commands, docs, tests, and scripts. Runs the full quality gate to confirm project health.

## Verification Checklist

1. [ ] `grep -rE '\btask (claim|submit|status|query|check-deps|validate-index|verify-task-done|quality-gate|cleanup|feature|prompt|add|index|migrate|validate-specs|record|all-completed|verify-completion|check|validate)\b' plugins/ forge-cli/docs/ .claude/` returns zero matches
2. [ ] `just check-stale-refs` passes (or equivalent grep-based check)
3. [ ] hooks.json parses as valid JSON with `forge` command references
4. [ ] All 12 modified skill files are valid markdown
5. [ ] All 4 doc files contain only `forge` command references
6. [ ] Go tests pass with new binary name and command paths
7. [ ] `go build ./...` compiles without errors
8. [ ] `go test ./...` passes
9. [ ] No deviations from design spec (or deviations are documented as decisions)

## Reference Files

- `design/tech-design.md` — §Appendix — Phase 4 Reference Update Map (complete enumeration)
- Phase 4 task records — `records/4.*.md`
- Phase 4 summary — `records/4-summary.md`

## Acceptance Criteria

- [ ] All applicable verification checklist items pass
- [ ] Any deviations from design are documented as decisions in the record
- [ ] Record created via `/record-task` with test evidence

## Implementation Notes

This is a verification-only task. No new feature code should be written.
If issues are found:
1. Fix inline if trivial (e.g., missed reference in a single file)
2. Document non-trivial issues as decisions in the record
3. Set status to `blocked` if a blocking issue cannot be resolved
