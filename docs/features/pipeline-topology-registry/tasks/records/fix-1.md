---
status: "completed"
started: "2026-05-29 21:45"
completed: "2026-05-29 21:56"
time_spent: "~11m"
---

# Task Record: fix-1 Fix: pipeline.go deviates from proposal on 9 points

## Summary
Fixed all 9 deviations in pipeline.go from the authoritative proposal: (1) IntentGate distribution corrected so only T-review-doc uses GateAllowAll, all others default to GateBlockSkipTest; (2) added Mode field to PipelineNode with correct breakdown/quick assignments; (3) replaced PerSurfaceKey bool with Expansion string field (per-surface-type/per-surface-key); (4) CondHasTestableTasks now checks IsTestableType(businessTasks) instead of GateTest+!isSkipTestIntent; (5) GenerateCondFunc signature changed to func(tasks []Task) bool; (6) DepResolveFunc signature changed to func(ctx *GenContext) []string; (7) ConfigGateFunc parameter order corrected to func(mode string, auto ...); (8) GenContext fields updated to BusinessTasks/UpstreamIDs/RunTestChain/AllGenerated, removed GeneratedTasks/SurfaceTypes; (9) DepRef fields renamed to Ref+Resolve.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/pipeline.go

### Key Decisions
- Used proposal as authoritative spec over existing code for all 9 deviations
- Resolver functions changed from method-style to var-style (DepResolveFunc assigned to package-level vars) matching proposal pattern
- ResolveIfGenerated changed from direct function to factory function returning DepResolveFunc, matching proposal signature

## Test Results
- **Tests Executed**: Yes
- **Passed**: 470
- **Failed**: 0
- **Coverage**: 83.9%

## Acceptance Criteria
- [x] AC-1: IntentGate distribution - only T-review-doc has GateAllowAll
- [x] AC-2: PipelineNode has Mode field with correct breakdown/quick assignments
- [x] AC-3: Expansion string field replaces PerSurfaceKey bool
- [x] AC-4: CondHasTestableTasks checks IsTestableType(businessTasks)
- [x] AC-5: GenerateCondFunc signature is func(tasks []Task) bool
- [x] AC-6: DepResolveFunc signature is func(ctx *GenContext) []string
- [x] AC-7: ConfigGateFunc parameter order is (mode, auto)
- [x] AC-8: GenContext has BusinessTasks/UpstreamIDs/RunTestChain/AllGenerated
- [x] AC-9: DepRef uses Ref+Resolve field names

## Notes
All 9 deviations fixed. Compilation, go vet, and existing test suite all pass. The pipeline.go types are currently only consumed within the file itself (no external consumers of GenContext/DepRef/PipelineRegistry yet), so the type signature changes are safe.
