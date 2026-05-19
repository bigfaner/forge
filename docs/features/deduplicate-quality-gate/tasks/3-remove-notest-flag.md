---
id: "3"
title: "Remove noTest flag from all structs and logic"
priority: "P1"
estimated_time: "1-2h"
dependencies: ["1"]
scope: "backend"
breaking: true
type: "coding.cleanup"
mainSession: false
---

# 3: Remove noTest Flag (Full Removal)

## Description

The `noTest` flag is 100% redundant with `IsTestableType()`. All auto-generated tasks using `noTest: true` have non-`coding.*` types (e.g., `doc*`, `test.*`, `validation.*`), which are already handled by `IsTestableType()` returning false. Remove `NoTest` from all Go structs, logic, and generated output.

## Reference Files
- `docs/proposals/deduplicate-quality-gate/proposal.md` — Source proposal (item 8)

## Acceptance Criteria

- [ ] `NoTest` field removed from `Task` struct (`types.go`)
- [ ] `NoTest` field removed from `TaskState` struct (`types.go`)
- [ ] `NoTest` field removed from `FrontmatterData` struct (`frontmatter.go`)
- [ ] `NoTest` field removed from `TestTaskDef` struct (`testgen.go`)
- [ ] All references to `NoTest` / `noTest` removed from `submit.go` (auto-coverage logic, gate skip, `formatTestsExecuted()`)
- [ ] All references to `NoTest` removed from `claim.go` (state saving)
- [ ] All references to `NoTest` removed from `testgen.go` (task generation, `TaskFromFile()`, frontmatter output)
- [ ] `IsTestableType()` is sole authority for test requirement determination
- [ ] All existing tests pass; tests referencing NoTest updated
- [ ] `go build ./...` compiles without error

## Implementation Notes

- **`forge-cli/pkg/task/types.go`**: Remove `NoTest` from `Task` (line 118), `TaskState` (line 216).
- **`forge-cli/pkg/task/frontmatter.go`**: Remove `NoTest` from `FrontmatterData` (line 20).
- **`forge-cli/pkg/task/testgen.go`**: Remove `NoTest` from `TestTaskDef` (line 21). Update all task definitions that set `NoTest: true` — remove the field. Update `TaskFromFile()` (line 595) to not set NoTest. Update `GenerateTestTaskMD()` (lines 228-229) to not output `noTest: true`.
- **`forge-cli/internal/cmd/submit.go`**: Remove `|| t.NoTest` from lines 129 and 139. Update `formatTestsExecuted()` (lines 441-449) to remove `noTest` parameter and the noTest branch.
- **`forge-cli/internal/cmd/claim.go`**: Remove `NoTest: t.NoTest` from state creation (line 128).
- **Backward compatibility**: Existing `index.json` files with `"noTest": true` will simply have the field ignored on deserialization (Go's JSON unmarshalling ignores unknown fields). No migration needed.
- TDD: remove NoTest from test fixtures and assertions first, then remove from structs.
