---
id: "1.gate"
title: "Phase 1 Gate: Base Rename Verification"
priority: "P0"
estimated_time: "1h"
dependencies: ["1.summary"]
breaking: true
type: "gate"
mainSession: false
---

# 1.gate: Phase 1 Gate: Base Rename Verification

## Description

Exit verification gate for Phase 1. Confirms that the module/directory/binary rename is complete and all code compiles and tests pass with new import paths.

## Verification Checklist

1. [ ] `forge-cli/` directory exists (no `task-cli/` remains)
2. [ ] `go.mod` module path is `forge-cli`
3. [ ] `go build ./...` compiles without errors
4. [ ] All existing tests pass with new import paths
5. [ ] `grep -r "task-cli" forge-cli/` returns zero matches (no stale import paths)
6. [ ] `forge-cli/pkg/version/version.go` has `Name = "forge"`
7. [ ] No deviations from design spec (or deviations are documented as decisions)

## Reference Files

- `design/tech-design.md` — File Structure Changes, Version Name Change, Module Path Change
- Phase 1 task records — `records/1.*.md`
- Phase 1 summary — `records/1-summary.md`

## Acceptance Criteria

- [ ] All applicable verification checklist items pass
- [ ] Any deviations from design are documented as decisions in the record
- [ ] Record created via `/record-task` with test evidence

## Implementation Notes

This is a verification-only task. No new feature code should be written.
If issues are found:
1. Fix inline if trivial (e.g., missed import path)
2. Document non-trivial issues as decisions in the record
3. Set status to `blocked` if a blocking issue cannot be resolved
