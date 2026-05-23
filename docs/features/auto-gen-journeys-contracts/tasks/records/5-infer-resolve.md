---
status: "completed"
started: "2026-05-24 00:33"
completed: "2026-05-24 00:42"
time_spent: "~9m"
---

# Task Record: 5 更新 infer.go 和 ResolveFirstTestDep 适配新任务类型

## Summary
Updated InferType pattern ordering (gen-journeys/gen-contracts before gen-scripts per Hard Rule), added panic on missing gen-journeys in ResolveFirstTestDep via findTaskIndexByPrefixOrPanic helper, and added 6 new tests covering panic behavior, clean-code integration, and edge cases

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/infer.go
- forge-cli/pkg/task/autogen.go
- forge-cli/pkg/task/infer_test.go
- forge-cli/pkg/task/autogen_test.go

### Key Decisions
- Reordered InferType switch cases: gen-contracts and gen-journeys now precede gen-scripts (Hard Rule compliance, functionally equivalent since prefixes don't overlap)
- Extracted findTaskIndexByPrefixOrPanic helper to DRY the panic-on-missing behavior in ResolveFirstTestDep
- ResolveFirstTestDep breakdown branch: replaced 3-level fallback chain (gen-journeys -> eval-journey -> gen-scripts) with direct findTaskIndexByPrefixOrPanic for gen-journeys
- ResolveFirstTestDep quick branch: replaced conditional check with findTaskIndexByPrefixOrPanic for gen-journeys

## Test Results
- **Tests Executed**: Yes
- **Passed**: 6
- **Failed**: 0
- **Coverage**: 91.8%

## Acceptance Criteria
- [x] InferType() identifies T-test-gen-journeys (with type suffix variants like T-test-gen-journeys-tui) returning correct type
- [x] InferType() identifies T-test-gen-contracts returning correct type
- [x] ResolveFirstTestDep() Breakdown branch: first task updated to T-test-gen-journeys via findTaskIndexByPrefix
- [x] ResolveFirstTestDep() Quick branch: first task updated to T-test-gen-journeys via findTaskIndexByPrefix
- [x] ExtractTypeSuffix() correctly handles T-test-gen-journeys-{type} type suffix extraction
- [x] All existing tests pass

## Notes
Hard Rule for pattern ordering enforced: gen-journeys/gen-contracts cases now precede gen-scripts in InferType. Hard Rule for panic on -1 enforced: findTaskIndexByPrefixOrPanic panics with descriptive message including all task IDs. Coverage increased from 90.8% to 91.8%.
