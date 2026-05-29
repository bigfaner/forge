---
id: "5"
title: "Add two-phase validation for pipeline registry"
priority: "P1"
estimated_time: "1.5h"
complexity: "high"
dependencies: [1]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 5: Add Two-Phase Validation for Pipeline Registry

## Description
Implement the two-phase validation system that ensures the PipelineRegistry is structurally correct before any task generation runs. Phase 1 is a static init-time validation via `init()` that panics on failure (replaces `ValidateAutogenTemplates`). Phase 2 is a dynamic runtime validation in `GenerateTestTasks` that returns errors.

## Reference Files
- `forge-cli/pkg/task/pipeline.go`: Add `ValidatePipelineRegistry` (Phase 1) and runtime checks in `GenerateTestTasks` (Phase 2) (source: docs/proposals/pipeline-topology-registry/proposal.md § Scope > In Scope, item 7)
- `forge-cli/pkg/task/autogen.go:120-175`: `ValidateAutogenTemplates` — replaced by Phase 1 validation (source: docs/proposals/pipeline-topology-registry/proposal.md § Scope > In Scope, item 7)

## Acceptance Criteria
- [ ] Phase 1 `ValidatePipelineRegistry()` runs in `init()`, validates: all `DependsOn.Ref` strings reference existing node IDs; `ResolveIfGenerated` references point to nodes declared before the caller; all expanded IDs are unique; `GenerateCondition` is non-nil; `Key`/`ID` template placeholders match `Expansion` setting; escape hatch count <= 5
- [ ] Phase 1 panics on failure with actionable error messages
- [ ] Phase 2 runs at start of `GenerateTestTasks`, validates: all resolver-returned IDs exist in generated task set; no circular dependencies; returns errors (does not panic)
- [ ] `ValidateAutogenTemplates` deleted after Phase 1 replaces its coverage
- [ ] `go build ./...` passes

## Hard Rules
- 仅修改以下文件：`forge-cli/pkg/task/pipeline.go`, `forge-cli/pkg/task/autogen.go`

## Implementation Notes

### Test Impact
- Affected test suite(s): `forge-cli/pkg/task/autogen_test.go`
- Expected fixture changes: test expectations for ValidateAutogenTemplates will need updating (deferred to task 6)
- Risk level: medium

- Phase 1 ordering invariants: `ResolveUpstream` users must appear after their expected upstream nodes
- Phase 2 circular dependency check can use standard topological sort on generated task graph
- `--no-validate` emergency bypass is out of scope for this task (documented as risk mitigation, not implemented)
