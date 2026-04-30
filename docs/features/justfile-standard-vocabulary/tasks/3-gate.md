---
id: "3.gate"
title: "Phase 3 Exit Gate"
priority: "P0"
estimated_time: "1h"
dependencies: ["3.summary"]
status: pending
breaking: true
---

# 3.gate: Phase 3 Exit Gate

## Description

Exit verification gate for Phase 3 (Reference Implementation). Confirms the forge project justfile works correctly as the mixed project reference implementation.

## Verification Checklist

1. [ ] Forge project justfile contains all 15 standard commands
2. [ ] `just project-type` returns `mixed` with exit code 0
3. [ ] `just compile frontend` executes without error
4. [ ] `just compile backend` executes without error
5. [ ] `just compile` (no scope) executes both frontend and backend
6. [ ] Invalid scope (`just build foo`) exits with code 1 and stderr message
7. [ ] Boundary markers present in justfile
8. [ ] Existing e2e tests pass: `just test-e2e --feature justfile-e2e-integration`
9. [ ] No deviations from design spec (or deviations are documented as decisions)

## Reference Files

- `docs/features/justfile-standard-vocabulary/design/tech-design.md` — Interface 1-2, Model 5
- Phase 3 task records: `records/3.*.md`
- Phase 3 summary: `records/3-summary.md`
- `justfile` — Reference implementation under test

## Acceptance Criteria

- [ ] All applicable verification checklist items pass
- [ ] Any deviations from design are documented as decisions in the record
- [ ] Record created via `/record-task` with test evidence

## Implementation Notes

This is a verification-only task. No new feature code should be written.
If issues are found:
1. Fix inline if trivial
2. Document non-trivial issues as decisions in the record
3. Set status to `blocked` if a blocking issue cannot be resolved
