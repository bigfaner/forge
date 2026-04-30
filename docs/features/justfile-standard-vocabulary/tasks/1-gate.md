---
id: "1.gate"
title: "Phase 1 Exit Gate"
priority: "P0"
estimated_time: "1h"
dependencies: ["1.summary"]
status: pending
breaking: true
---

# 1.gate: Phase 1 Exit Gate

## Description

Exit verification gate for Phase 1 (Foundation). Confirms that all foundation outputs are complete, internally consistent, and match the design specification before integration work begins.

## Verification Checklist

1. [ ] `task-cli/pkg/task/types.go`: `Task` and `TaskState` structs include `Scope` field with correct json tags
2. [ ] `task-cli/pkg/task/types_test.go`: Unit tests for scope serialization pass (`go test ./...`)
3. [ ] `index.schema.json`: `scope` property with enum `["frontend", "backend", "all"]` present in tasks schema
4. [ ] `init-justfile.md`: Backend template contains all 15 recipes with correct Go commands
5. [ ] `init-justfile.md`: Frontend template contains all 15 recipes with correct npm commands
6. [ ] `init-justfile.md`: Mixed template contains 10 scoped recipes with bash case dispatch + 5 unscoped
7. [ ] `run-e2e-tests/SKILL.md`: Contains `just run`, no `npx serve`
8. [ ] `execute-task.md`, `task-executor.md`, `error-fixer.md`: Contain `just compile && just test`, no `just build && just test`
9. [ ] 10 validation files contain expected standard commands (grep checks pass)
10. [ ] No deviations from design spec (or deviations are documented as decisions)

## Reference Files

- `docs/features/justfile-standard-vocabulary/design/tech-design.md` — Cross-Layer Data Map, Model 1-5
- Phase 1 task records: `records/1.*.md`
- Phase 1 summary: `records/1-summary.md`

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
