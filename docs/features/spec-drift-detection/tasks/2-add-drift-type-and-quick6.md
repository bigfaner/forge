---
id: "2"
title: "Add doc-generation.drift type and T-quick-6 to Go test pipeline"
priority: "P0"
estimated_time: "1-2h"
dependencies: ["1"]
scope: "backend"
breaking: true
type: "implementation"
mainSession: false
---

# 2: Add doc-generation.drift type and T-quick-6 to Go test pipeline

## Description

Add a new task type `doc-generation.drift` to the Go CLI's type system, create its strategy template, and register T-quick-6 in the quick test pipeline so that quick-mode features automatically get a drift detection test task after verify-regression.

## Reference Files
- `docs/proposals/spec-drift-detection/proposal.md` — Source proposal
- `forge-cli/pkg/task/types.go` — Type constants, registry, valid types map
- `forge-cli/pkg/task/testgen.go` — Test task generation (T-quick-1..5, T-test-1..5)
- `forge-cli/pkg/prompt/prompt.go` — Type-to-strategy-template mapping
- `forge-cli/pkg/prompt/data/doc-generation-consolidate.md` — Reference template

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/task/types.go` | Add `TypeDocGenerationDrift` constant, registry entry, valid type |
| `forge-cli/pkg/task/testgen.go` | Add T-quick-6 definition in `generateQuickTestTasks`, update `resolveQuickDeps` |
| `forge-cli/pkg/prompt/prompt.go` | Add mapping: `doc-generation.drift` → `data/doc-generation-drift.md` |

### Create
| File | Description |
|------|-------------|
| `forge-cli/pkg/prompt/data/doc-generation-drift.md` | Strategy template for drift-only mode |

## Acceptance Criteria

- [ ] `TypeDocGenerationDrift = "doc-generation.drift"` added to types.go constants
- [ ] Registry entry added with description "detect and fix spec drift against codebase"
- [ ] Valid type map includes `doc-generation.drift`
- [ ] `doc-generation-drift.md` strategy template created — invokes `consolidate-specs` skill in drift-only mode (skip extraction, run Steps 9-11 only)
- [ ] prompt.go maps new type to strategy template path
- [ ] T-quick-6 added after T-quick-5 in `generateQuickTestTasks` with `TypeDocGenerationDrift`
- [ ] T-quick-6 depends on T-quick-5 in `resolveQuickDeps`
- [ ] All existing tests pass (`go test ./...`)
- [ ] Version bumped in `scripts/version.txt` (minor: new feature)

## Hard Rules

- Follow dependency direction: `cmd → internal → pkg` (no reverse)
- TDD: write tests first for new type constant, registry entry, and T-quick-6 generation
- Table-driven tests for new test cases

## Implementation Notes

- The strategy template for drift should be minimal — it just tells the executor to invoke `consolidate-specs` skill, which will detect drift-only mode automatically when no PRD/design files exist
- T-quick-6 has `NoTest: true` like T-test-5 (spec consolidation tasks don't run unit tests)
- T-quick-6 has `Scope: "all"` — drift detection touches all spec files
