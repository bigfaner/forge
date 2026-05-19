---
id: "4"
title: "Expand testableTypes and add type-based quality-gate skip in Go"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: true
type: "feature"
mainSession: false
---

# 4: Expand testableTypes and add type-based quality-gate skip in Go

## Description
Two Go code changes to make the type-based quality-gate behavior work correctly:

**1. Expand `testableTypes`** (`forge-cli/pkg/task/build.go`): The `testableTypes` map currently only includes `feature`, `enhancement`, and `fix`. Missing `cleanup` and `refactor` causes features with only cleanup/refactor tasks to be treated as docs-only — no quality-gate, no test pipeline generated.

**2. Type-based quality-gate skip** (`forge-cli/internal/cmd/submit.go`): The quality-gate pre-check at submit time uses only `t.NoTest` to decide whether to skip. Add `IsTestableType(t.Type)` as an additional condition: only run quality-gate when `!t.NoTest && IsTestableType(t.Type)`. Also update the coverage auto-set logic to trigger for non-testable types.

## Reference Files
- `docs/proposals/task-type-code-docs-boundary/proposal.md` — Source proposal
- `plugins/forge/references/shared/type-assignment.md` — Classification rule (task 1 output)
- `forge-cli/pkg/task/build.go` — `testableTypes` map, `IsTestableType`, `needsTestPipeline`
- `forge-cli/internal/cmd/submit.go` — quality-gate skip logic
- `forge-cli/internal/cmd/quality_gate.go` — `isDocsOnly` (already uses `IsTestableType`, no change needed)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/task/build.go` | Add `TypeCleanup` and `TypeRefactor` to `testableTypes` map |
| `forge-cli/internal/cmd/submit.go` | Quality-gate skip: `!t.NoTest` → `!t.NoTest && task.IsTestableType(t.Type)`; coverage auto-set: also trigger for non-testable types |

## Acceptance Criteria
- [ ] `testableTypes` map includes `TypeCleanup` and `TypeRefactor`
- [ ] `IsTestableType("cleanup")` returns true
- [ ] `IsTestableType("refactor")` returns true
- [ ] `needsTestPipeline()` returns true when any task has type cleanup or refactor
- [ ] `submit.go` skips quality-gate for tasks where `!IsTestableType(t.Type)` (e.g., `type: "documentation"`)
- [ ] Coverage auto-set (`coverage = -1.0`) triggers for non-testable types in addition to `noTest`
- [ ] All existing tests pass
- [ ] New unit tests cover: expanded testableTypes, type-based quality-gate skip condition

## Hard Rules
- Follow TDD: write failing tests first, then implement
- Run `go build ./...` and `go test -race -cover ./...` after changes
- Bump version in `scripts/version.txt` (minor: new behavior for existing types)

## Implementation Notes
- `isDocsOnly()` in quality_gate.go already calls `task.IsTestableType()` — will automatically benefit from the expanded map
- `needsDocEval()` in build.go checks for ALL tasks being documentation type — unaffected by this change
- Backward compatible: adding types to testableTypes only ADDS quality-gate enforcement, never removes it
- Risk: cleanup/refactor tasks that previously skipped quality-gate will now run it. This is the intended fix.
