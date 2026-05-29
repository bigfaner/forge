---
status: "completed"
started: "2026-05-29 21:20"
completed: "2026-05-29 21:29"
time_spent: "~9m"
---

# Task Record: 1 Define pipeline registry core types and declarations

## Summary
Created pkg/task/pipeline.go with full Pipeline Topology Registry — the single source of truth for all auto-generated task definitions. Includes core types (PipelineNode, DepRef, GenContext, ConfigGateFunc, IntentGateFunc, GenerateCondFunc, DepResolveFunc), all gate/condition/resolver functions, and the PipelineRegistry slice with 12 nodes in declaration order.

## Changes

### Files Created
- forge-cli/pkg/task/pipeline.go

### Files Modified
无

### Key Decisions
- ResolveHighestGateOrLastBiz uses two-pass logic (gate→summary) then compares with lastBiz by phase, matching existing findHighestGateOrSummary behavior
- ResolveDocTasks returns comma-joined sorted doc-task IDs for multi-dependency resolution
- PipelineRegistry includes both T-specs-consolidate (breakdown) and T-quick-doc-drift (quick) as separate nodes, allowing mode-based filtering by callers

## Test Results
- **Tests Executed**: Yes
- **Passed**: 12
- **Failed**: 0
- **Coverage**: 84.1%

## Acceptance Criteria
- [x] pkg/task/pipeline.go exists with PipelineNode, DepRef, GenContext, ConfigGateFunc, IntentGateFunc, GenerateCondFunc, DepResolveFunc types
- [x] All 4 ConfigGate functions defined (GateTest, GateValidation, GateConsolidateSpecs, GateCleanCode) with correct mode/config mapping
- [x] Both IntentGate functions defined (GateAllowAll, GateBlockSkipTest)
- [x] All 3 GenerateCondition functions defined (CondHasTestableTasks, CondHasDocTasks, CondAlways)
- [x] All resolver functions defined including ResolveIfGenerated, ResolveHighestGateOrLastBiz, ResolveLastBusinessTask, ResolveLastRunTestOrBusiness, ResolveUpstream, ResolveDocTasks
- [x] PipelineRegistry slice contains all 12 nodes in correct declaration order
- [x] T-review-doc has IntentGate: GateAllowAll; T-clean-code has DependsOn: ResolveHighestGateOrLastBiz; T-test-gen-journeys has DependsOn: ResolveIfGenerated for both T-review-doc and T-clean-code
- [x] go build ./... passes with the new file

## Notes
No proposal.md found — task AC and implementation notes used as authoritative spec. New file only, no existing files modified.
