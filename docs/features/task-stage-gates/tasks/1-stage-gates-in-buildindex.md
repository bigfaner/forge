---
id: "1"
title: "Extend BuildIndex with phase detection and stage-gate generation"
priority: "P0"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: true
type: "implementation"
mainSession: false
---

# 1: Extend BuildIndex with phase detection and stage-gate generation

## Description
Extend `forge task index` (specifically `BuildIndex` in `pkg/task/build.go`) to automatically detect numbered phases from existing business task IDs and generate `.summary` + `.gate` task files for each qualifying phase. This provides deterministic phase-level quality checkpoints for both quick mode and full mode.

The implementation follows the same pattern as test-task auto-generation already in `build.go`: detect → check existence → generate from embedded template → merge with existing.

## Reference Files
- `docs/proposals/task-stage-gates/proposal.md` — Source proposal with full algorithm description
- `forge-cli/pkg/task/build.go` — Current BuildIndex implementation (generation pattern to follow)
- `forge-cli/pkg/task/testgen.go` — Existing test-task generation pattern (programmatic MD generation to follow)
- `forge-cli/pkg/task/build_test.go` — Existing tests (extend with new test cases)
- `plugins/forge/skills/breakdown-tasks/templates/gate-task.md` — Gate task content reference (what the generated files should look like)
- `plugins/forge/skills/breakdown-tasks/templates/phase-summary-task.md` — Summary task content reference

## Affected Files

### Create
| File | Description |
|------|-------------|
| `forge-cli/pkg/task/stage_gates.go` | Phase detection, template embedding, generation logic |
| `forge-cli/pkg/task/stage_gates_test.go` | Unit tests for phase detection, generation, idempotency, edge cases |

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/task/build.go` | Call stage-gate generation between existing step 5 (scan .md) and step 7 (generate test tasks) |
| `forge-cli/pkg/task/build_test.go` | Add integration test: BuildIndex end-to-end with stage-gate generation |

## Acceptance Criteria
- [ ] `forge task index --feature <slug>` generates `.summary` and `.gate` for every numbered phase with >=2 business tasks
- [ ] Single-task phases (<2 business tasks) are skipped — no files generated
- [ ] T-test/T-quick task IDs excluded from phase business task count
- [ ] Generated `.gate` has `depends_on: ["<N>.summary"]` and `breaking: true`
- [ ] Generated `.summary` has `depends_on` set to all business task IDs in the same phase
- [ ] Re-running `forge task index` does not overwrite existing files (idempotent)
- [ ] Partial state handled: if `.summary` exists but `.gate` missing, only `.gate` is generated
- [ ] Malformed task IDs (e.g., `intro`, `1.2a`) silently skipped — no crash, no gate generated
- [ ] Generated tasks appear in `index.json` with correct `type` (`gate` / `doc-generation.summary`)
- [ ] Unit tests cover: phase detection, test-task exclusion, template generation, idempotency, malformed IDs, partial state, threshold check

## Hard Rules
- Generate `.md` content programmatically (same pattern as `GenerateTestTaskMD` in `testgen.go`). Do NOT use `go:embed` for templates — `go:embed` can't reach `plugins/` from `forge-cli/pkg/task/`. Reference `plugins/forge/skills/breakdown-tasks/templates/gate-task.md` and `phase-summary-task.md` for content structure, but generate the markdown as Go string literals.
- Phase detection MUST use `strings.Split(id, ".")` yielding exactly 2 segments, both parseable via `strconv.Atoi`
- `--no-test` flag MUST NOT affect stage-gate generation (gates are structural, not test tasks)

## Implementation Notes
- Phase detection regex: task ID matches pattern `<digit>.<digit>`. Exclude IDs starting with `T-test` or `T-quick`.
- Gate template needs: phase number, feature slug, list of same-phase business task IDs (for reference section)
- Summary template needs: phase number, feature slug, list of same-phase business task IDs (for dependencies)
- Add generation output to BuildIndexResult: `StageGatesGenerated int` field
- Follow TDD: write phase detection tests first, then template generation tests, then integration
