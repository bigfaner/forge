---
id: "6"
title: "Update existing tests and add registry-specific tests"
priority: "P1"
estimated_time: "2.5h"
complexity: "high"
dependencies: [2, 3, 4, 5]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 6: Update Existing Tests and Add Registry-Specific Tests

## Description
Update all existing test suites affected by the refactoring, and add new tests for registry validation and derived function correctness. Ensure all tests pass with zero regressions.

## Reference Files
- `forge-cli/pkg/task/autogen_test.go`: Update tests for deleted functions (GetBreakdownTestTasks, GetQuickTestTasks, resolveBreakdownDeps, resolveQuickDeps) (source: proposal.md#Scope item 8)
- `forge-cli/pkg/task/build_test.go`: Update tests for step 7/7.5/7.6 rewrite (source: proposal.md#Scope item 8)
- `forge-cli/pkg/task/claim_test.go`: Update tests for claim logic changes (source: proposal.md#Scope item 8)
- `forge-cli/pkg/task/pipeline.go`: Add tests for registry validation, GenerateTestTasks, resolver functions (source: proposal.md#Scope item 9)
- `forge-cli/pkg/task/infer_test.go`: Update InferType tests for registry-driven matching (source: proposal.md#Scope item 8)

## Acceptance Criteria
- [ ] All existing tests pass: `autogen_test.go`, `infer_test.go`, `build_test.go`, `claim_test.go`
- [ ] Snapshot tests capture generated task lists for all mode/config/intent combinations before and after refactor
- [ ] New tests for `GenerateTestTasks` cover: empty business tasks, mixed coding+doc tasks, quick vs breakdown mode, skip-test intent, all config gates off, single surface degenerate case
- [ ] New tests for resolver functions: `ResolveHighestGateOrLastBiz` (two-pass gate priority), `ResolveIfGenerated` (nil when not generated), `ResolveLastRunTestOrBusiness` (fallback to business), `ResolveUpstream` (serial chain)
- [ ] New tests for `InferType` registry-driven matching: surface-expanded variants, T-review-doc, stage-gate suffixes, runtime task prefixes
- [ ] New tests for `ValidatePipelineRegistry` Phase 1: missing ref, forward ResolveIfGenerated reference, duplicate expanded IDs, nil GenerateCondition
- [ ] Test coverage >= 80% for `pipeline.go` and refactored functions

## Implementation Notes

### Test Impact
- This IS the test task — no deferred fixture changes
- Risk level: high

- Pre-refactoring snapshot strategy: run existing tests to capture current task generation output, then compare with post-refactoring output
- Mixed coding+doc test case verifies T-review-doc and T-clean-code both appear as T-test-gen-journeys dependencies
- Refactor/cleanup intent test verifies only T-review-doc (when doc tasks exist) generates — no test/validation/consolidate/clean-code tasks
