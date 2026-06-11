---
id: "1"
title: "Define pipeline registry core types and declarations"
priority: "P0"
estimated_time: "2h"
complexity: "high"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 1: Define Pipeline Registry Core Types and Declarations

## Description
Create `pkg/task/pipeline.go` with the full Pipeline Topology Registry — the single source of truth for all auto-generated task definitions. This includes core data structures (PipelineNode, DepRef, GenContext), all gate/condition/resolver functions, and the PipelineRegistry slice in declaration order.

This is the foundation task — all subsequent refactoring tasks depend on the registry being defined.

## Reference Files
- `forge-cli/pkg/task/pipeline.go`: New file — all registry types and declarations (source: docs/proposals/pipeline-topology-registry/proposal.md § Core Data Structure)
- `forge-cli/pkg/task/autogen.go:178-181`: `isSkipTestIntent` — referenced by `GateBlockSkipTest` (source: docs/proposals/pipeline-topology-registry/proposal.md § Predefined Gate, Condition & Resolver Functions)
- `forge-cli/pkg/task/types.go:42-64`: Type constants referenced by registry entries (source: docs/proposals/pipeline-topology-registry/proposal.md § Pipeline Registry)
- `forge-cli/pkg/task/autogen.go:922-994`: `findHighestGateOrSummary`/`findMaxBusinessTaskID`/`phaseFromID`/`numericID` — logic migrates into resolvers (source: docs/proposals/pipeline-topology-registry/proposal.md § Predefined Gate, Condition & Resolver Functions)

## Acceptance Criteria
- [ ] `pkg/task/pipeline.go` exists with `PipelineNode`, `DepRef`, `GenContext`, `ConfigGateFunc`, `IntentGateFunc`, `GenerateCondFunc`, `DepResolveFunc` types
- [ ] All 4 ConfigGate functions defined (GateTest, GateValidation, GateConsolidateSpecs, GateCleanCode) with correct mode/config mapping
- [ ] Both IntentGate functions defined (GateAllowAll, GateBlockSkipTest)
- [ ] All 3 GenerateCondition functions defined (CondHasTestableTasks, CondHasDocTasks, CondAlways)
- [ ] All resolver functions defined including ResolveIfGenerated, ResolveHighestGateOrLastBiz (two-pass gate priority), ResolveLastBusinessTask (using numericID), ResolveLastRunTestOrBusiness, ResolveUpstream, ResolveDocTasks
- [ ] PipelineRegistry slice contains all 12 nodes in correct declaration order: T-review-doc, T-clean-code, T-test-gen-journeys, T-eval-journey, T-test-gen-contracts, T-eval-contract, T-test-gen-scripts, T-test-run, T-validate-code, T-validate-ux, T-specs-consolidate/T-quick-doc-drift
- [ ] T-review-doc has IntentGate: GateAllowAll; T-clean-code has DependsOn: ResolveHighestGateOrLastBiz; T-test-gen-journeys has DependsOn: ResolveIfGenerated for both T-review-doc and T-clean-code
- [ ] `go build ./...` passes with the new file

## Implementation Notes
- GenContext.ExistingTasks must be populated by caller with full index (including gates/summaries)
- ResolveHighestGateOrLastBiz must use two-pass logic: first find highest-phase gate, then fallback to summary only if no gate found (matches current findHighestGateOrSummary)
- ResolveLastBusinessTask must use numericID for sorting (not slice order)
- per-surface-key expansion creates serial chains (first task uses DependsOn resolvers, subsequent tasks depend on previous expanded task)
