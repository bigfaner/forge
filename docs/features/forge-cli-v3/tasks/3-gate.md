---
id: "3.gate"
title: "Phase 3 Gate: Feature Migration Verification"
priority: "P0"
estimated_time: "1h"
dependencies: ["3.summary"]
breaking: true
type: "gate"
mainSession: false
---

# 3.gate: Phase 3 Gate: Feature Migration Verification

## Description

Exit verification gate for Phase 3. Confirms that e2e subcommands and probe command work correctly, and that behavior is equivalent to the justfile bash recipes they replace.

## Verification Checklist

1. [ ] `forge e2e run` reads profile from config.yaml and dispatches correctly
2. [ ] `forge e2e setup` installs dependencies idempotently
3. [ ] `forge e2e verify --feature <slug>` checks for VERIFY markers
4. [ ] `forge e2e compile` performs compile-check for active profile
5. [ ] `forge e2e discover` lists test cases without running
6. [ ] `forge probe` performs HTTP health checks
7. [ ] Profile error cases: no profile → "no e2e profile configured", unknown → "unknown profile: <value>"
8. [ ] External tool failures normalized to exit code 1 with descriptive stderr
9. [ ] Justfile: migrated recipes removed, no broken references
10. [ ] `go build ./...` compiles without errors
11. [ ] All existing tests pass
12. [ ] No deviations from design spec (or deviations are documented as decisions)

## Reference Files

- `design/tech-design.md` — §Interfaces 4-6 (probe, e2e subcommands, pkg/e2e)
- `design/tech-design.md` — §Error Handling — E2E External Tool Failures
- Phase 3 task records — `records/3.*.md`
- Phase 3 summary — `records/3-summary.md`

## Acceptance Criteria

- [ ] All applicable verification checklist items pass
- [ ] Any deviations from design are documented as decisions in the record
- [ ] Record created via `/record-task` with test evidence

## Implementation Notes

This is a verification-only task. No new feature code should be written.
If issues are found:
1. Fix inline if trivial (e.g., wrong error message format)
2. Document non-trivial issues as decisions in the record
3. Set status to `blocked` if a blocking issue cannot be resolved
