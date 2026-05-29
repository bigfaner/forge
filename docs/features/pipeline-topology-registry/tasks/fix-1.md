---
id: "fix-1"
title: "Fix: pipeline.go deviates from proposal on 9 points"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
type: "coding.fix"
---

# Fix: pipeline.go deviates from proposal on 9 points

## Root Cause

Task 1 implementation deviates from proposal.md in 9 areas. Proposal at docs/proposals/pipeline-topology-registry/proposal.md was inaccessible to executor (confirmed: record says 'No proposal.md found'). Fix pipeline.go to match proposal exactly.

CRITICAL BEHAVIORAL (4):
1. IntentGate: T-clean-code/T-validate-code/T-validate-ux/T-specs-consolidate/T-quick-doc-drift should NOT have GateAllowAll. Only T-review-doc gets GateAllowAll. All others default to GateBlockSkipTest (nil IntentGate).
2. PipelineNode missing Mode field (string: quick/breakdown/empty). T-eval-journey/T-test-gen-contracts/T-eval-contract/T-test-gen-scripts = breakdown only. T-quick-doc-drift = quick only.
3. PipelineNode missing Expansion field (string: per-surface-key/per-surface-type/empty). Replace PerSurfaceKey bool. T-test-gen-scripts needs per-surface-type. T-test-run needs per-surface-key.
4. CondHasTestableTasks should check IsTestableType(businessTasks), NOT GateTest+!isSkipTestIntent. Separation of concerns: ConfigGate checks config, IntentGate checks intent, GenerateCondition checks task composition.

TYPE SIGNATURE (5):
5. GenerateCondFunc: func(tasks []Task) bool (proposal) not func(ctx GenContext) bool
6. DepResolveFunc: func(ctx *GenContext) []string (proposal) not func(ref DepRef, ctx GenContext) string
7. ConfigGateFunc: func(mode string, auto forgeconfig.AutoConfig) bool (proposal param order)
8. GenContext: add BusinessTasks []Task, UpstreamIDs []string, RunTestChain []string, AllGenerated []string. Remove GeneratedTasks []AutoGenTaskDef and SurfaceTypes []string.
9. DepRef: use Ref+Resolve field names (proposal) not TaskTemplate+Resolver. Resolve is DepResolveFunc, Ref is static ID string.

Reference: docs/proposals/pipeline-topology-registry/proposal.md sections Core Data Structure, Predefined Gate/Condition/Resolver Functions, Pipeline Registry

## Reference Files

- Source: forge-cli/pkg/task/pipeline.go
- Test script: go test ./forge-cli/pkg/task/...
- Test results: 9 deviations from proposal: IntentGate distribution (6 nodes), missing Mode field, missing per-surface-type expansion, CondHasTestableTasks wrong logic, 5 type signature changes

## Surface Inference

This fix-task was created by the quality-gate hook. If `surface-key` and `surface-type` above are empty, infer them at execution time:

1. Parse `forge-cli/pkg/task/pipeline.go` to extract the first file path (comma-separated).
2. Run `forge surfaces --json <file-path>` to resolve surface-key/type.
3. Use the resolved surface-type to load the appropriate `rules/surfaces/<type>.md` for test orchestration guidance.

If `forge surfaces --json` fails (no surfaces configured, command not found), proceed without surface information — this does not block the fix.

## Fix Boundaries

When fixing test failures, observe these boundaries:

**Forbidden:**
- Starting dev server (`npx expo start`, `npm run dev`, etc.)
- Running `npm install` more than 3 times — mark task as blocked if dependency installation fails 3 times
- Running full test suite — regression is verified by the dispatcher after fix completes
- Manually opening browser to verify rendering

**Correct workflow:**
1. Read failing test + corresponding component source
2. Compare test's expected testID/selectors vs actual DOM structure
3. Modify component (add testID) or test (adjust selectors/assertions)
4. Run targeted tests on affected packages — unit tests must pass
5. Record completion

## Verification

After fixing, verify the fix works:
1. Run targeted tests on changed packages: `go test -race ./affected/package/...`
2. Replace the path with the actual packages you modified

> **Note:** Full project-wide tests run at CLI submit (`forge task submit`) — agent runs targeted tests only.

Full regression is verified by the dispatcher, not by this fix task.

When this task is recorded as completed via `task record`, the source task 1 is automatically restored to pending if all its dependencies are completed.
