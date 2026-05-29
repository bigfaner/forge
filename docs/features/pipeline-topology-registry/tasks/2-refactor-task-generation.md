---
id: "2"
title: "Refactor task generation to use registry-driven GenerateTestTasks"
priority: "P1"
estimated_time: "2h"
complexity: "high"
dependencies: [1]
surface-key: "."
surface-type: "cli"
breaking: true
type: "coding.refactor"
mainSession: false
---

# 2: Refactor Task Generation to Use Registry-Driven GenerateTestTasks

## Description
Replace the procedural `GetBreakdownTestTasks` and `GetQuickTestTasks` functions with a single registry-driven `GenerateTestTasks` that filters PipelineRegistry by Mode/ConfigGate/IntentGate/GenerateCondition/UISurfaceOnly, expands per-surface nodes with serial chain support, and resolves dependencies via GenContext progressive population.

Delete `GetBreakdownTestTasks`, `GetQuickTestTasks`, `resolveBreakdownDeps`, `resolveQuickDeps`, `wireRunTestChain`, `wireQuickRunTestChain`, and related helper functions.

## Reference Files
- `forge-cli/pkg/task/autogen.go:188-310`: `GetBreakdownTestTasks` — replaced by registry-driven generation (source: proposal.md#Derived-Functions)
- `forge-cli/pkg/task/autogen.go:324-414`: `GetQuickTestTasks` — replaced by registry-driven generation
- `forge-cli/pkg/task/autogen.go:603-702`: `resolveBreakdownDeps`/`resolveQuickDeps` — replaced by DepRef resolution
- `forge-cli/pkg/task/autogen.go:709-784`: `wireRunTestChain`/`wireQuickRunTestChain` — replaced by per-surface-key serial chain expansion
- `forge-cli/pkg/task/pipeline.go`: Registry-driven GenerateTestTasks implementation (from task 1)

## Acceptance Criteria
- [ ] `GenerateTestTasks(mode, surfaces, executionOrder, auto, intent, businessTasks, existingTasks)` implements the 5-step filter/expand/resolve/update/return algorithm
- [ ] Filtering correctly applies Mode, ConfigGate, IntentGate (nil defaults to GateBlockSkipTest), GenerateCondition, UISurfaceOnly
- [ ] per-surface-key expansion creates serial chains; per-surface-type expansion creates parallel tasks
- [ ] GenContext progressively populated: AllGenerated, UpstreamIDs, RunTestChain
- [ ] ResolveIfGenerated correctly returns nil when referenced node was not generated
- [ ] `GetBreakdownTestTasks`, `GetQuickTestTasks`, `resolveBreakdownDeps`, `resolveQuickDeps`, `wireRunTestChain`, `wireQuickRunTestChain` deleted
- [ ] `go build ./...` passes

## Implementation Notes

### Test Impact
- Affected test suite(s): `forge-cli/pkg/task/autogen_test.go`
- Expected fixture changes: test expectations for GetBreakdownTestTasks/GetQuickTestTasks will need updating (deferred to task 6)
- Risk level: high

- When isSingleSurface is true, per-surface-key expansion strips surface suffix from ID
- ResolveHighestGateOrLastBiz requires ExistingTasks to contain gates/summaries from index
- No post-generation injection needed (escape hatch count: 0) — all deps resolved via registry
