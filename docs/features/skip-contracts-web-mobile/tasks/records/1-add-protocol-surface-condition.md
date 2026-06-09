---
status: "completed"
started: "2026-06-09 18:13"
completed: "2026-06-09 18:22"
time_spent: "~9m"
---

# Task Record: 1 Add protocol surface condition to skip gen-contracts pipeline

## Summary
Added CondHasProtocolSurfaceTask condition to PipelineRegistry. Pure Web/Mobile features now skip T-test-gen-contracts and T-eval-contract nodes; dependency chain for gen-scripts naturally degrades to gen-journeys/eval-journey via ResolveUpstream. Mixed surface and existing API/CLI/TUI pipelines unchanged.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/pipeline.go
- forge-cli/pkg/task/pipeline_validate.go
- forge-cli/pkg/task/pipeline_test.go
- forge-cli/pkg/task/autogen_test.go

### Key Decisions
- Used existing ResolveUpstream mechanism for dependency degradation — no new resolver needed because when gen-contracts/eval-contract are skipped, ResolveUpstream naturally resolves to the last generated node (eval-journey or gen-journeys)
- Defined protocolSurfaceTypes as a package-level map constant following the uiSurfaceTypes pattern
- Conservative approach: nil/empty/unknown surface-type returns true (don't skip) to avoid breaking existing features

## Test Results
- **Tests Executed**: Yes
- **Passed**: 8
- **Failed**: 0
- **Coverage**: 88.5%

## Acceptance Criteria
- [x] SC-1: Pure Web feature skips T-test-gen-contracts and T-eval-contract
- [x] SC-2: Pure Mobile feature skips gen-contracts/eval-contract
- [x] SC-3: Mixed surface feature (api+web) generates gen-contracts/eval-contract normally
- [x] SC-4: Frontend-only feature in multi-surface project skips gen-contracts
- [x] SC-7: Existing API/CLI/TUI pipeline behavior unchanged
- [x] surface-type missing/empty/unknown: conservative no-skip with WARN log

## Notes
Dependency chain degradation works naturally through ResolveUpstream — no custom resolver needed. All 8 new tests pass plus all existing tests (no regressions).
