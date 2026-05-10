---
id: "{{ID}}"
title: "{{TITLE}}"
priority: "P0"
estimated_time: "1h"
dependencies: [{{DEPENDENCIES}}]
status: pending
breaking: true
mainSession: false
---

# {{ID}}: {{TITLE}}

## Description

Exit verification gate for this phase. Confirms that all outputs are complete, internally consistent, and match the design specification before the next phase begins.

## Verification Checklist

1. [ ] All interfaces from this phase compile without errors
2. [ ] Data models match `design/tech-design.md` (skip if single-layer feature — mark N/A)
3. [ ] No type mismatches between adjacent layers (skip if single-layer feature — mark N/A)
4. [ ] Project builds successfully
5. [ ] All existing tests pass
6. [ ] No deviations from design spec (or deviations are documented as decisions)
7. [ ] All Integration Specs from `tech-design.md` have corresponding code changes (for each Integration Spec: verify target file was modified since feature branch started; if branch point cannot be determined, verify target file was modified per the task record)
8. [ ] All integration test cases pass (if gen-test-cases already ran)

## Reference Files

- `design/tech-design.md` — Cross-Layer Data Map section (if exists)
- This phase's task records — `records/{{PHASE}}.*.md`
- This phase's summary — `records/{{PHASE}}-summary.md`

## Acceptance Criteria

- [ ] All applicable verification checklist items pass
- [ ] Any deviations from design are documented as decisions in the record
- [ ] Record created via `/record-task` with test evidence

## Implementation Notes

This is a verification-only task. No new feature code should be written.
If issues are found:
1. Fix inline if trivial (e.g., type mismatch in a single file)
2. Document non-trivial issues as decisions in the record
3. Set status to `blocked` if a blocking issue cannot be resolved

## Execution Workflow

1. Verify all interfaces compile without errors.
   - Command: `just compile [scope]`
   - Success: exit 0.
   - Failure: fix compilation errors inline if trivial; if not trivial, set status to `blocked`.
2. Run the project build.
   - Command: `just build [scope]` (or `just compile [scope]` if no separate build command)
   - Success: exit 0.
   - Failure: fix build errors inline if trivial; if not trivial, set status to `blocked`.
3. Run existing tests.
   - Command: `just test [scope]`
   - Success: all tests pass (exit 0).
   - Failure: fix test failures inline if trivial; if not trivial, set status to `blocked`.
4. Run quality gate in strict sequence:
   - `just fmt [scope]` — Success: exit 0. Failure: task is blocked.
   - `just lint [scope]` — Success: exit 0. Failure: self-fix once; if still failing, task is blocked.
5. Verify each checklist item from the Verification Checklist section above.
   - Success: all applicable items pass.
   - Failure: document deviation as a decision in the record; set status to `blocked` for blocking issues.
6. Stop. Proceed to Step 3 (Record).
