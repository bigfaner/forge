---
id: "4"
title: "Refactor build.go callers and delete legacy functions"
priority: "P1"
estimated_time: "2h"
complexity: "high"
dependencies: [2, 3]
surface-key: "."
surface-type: "cli"
breaking: true
type: "coding.refactor"
mainSession: false
---

# 4: Refactor build.go Callers and Delete Legacy Functions

## Description
Update build.go steps 7/7.5/7.6 to call the new registry-driven `GenerateTestTasks`. Delete all legacy functions replaced by the registry: `ResolveFirstTestDep`, `resolveTestDepsAndInjectReviewDoc`, `findHighestGateOrSummary`, `findMaxBusinessTaskID`, `ResolveDriftFallbackDep`, `GetReviewDocTask`, `ResolveReviewDocDep`, `findTaskIndexOrPanic`, `findTaskIndexByPrefixOrPanic`. Preserve `needsTestPipeline` and `needsReviewDoc` in caller for stage-gate/doc-criteria control.

## Reference Files
- `forge-cli/pkg/task/build.go:331-416`: Steps 7/7.5/7.6 — rewrite to call GenerateTestTasks with businessTasks + existingTasks (source: docs/proposals/pipeline-topology-registry/proposal.md § Functions Relationship to Registry)
- `forge-cli/pkg/task/autogen.go:843-919`: `ResolveFirstTestDep` — delete, logic covered by ResolveHighestGateOrLastBiz + ResolveIfGenerated (source: docs/proposals/pipeline-topology-registry/proposal.md § Functions Relationship to Registry)
- `forge-cli/pkg/task/build.go:588-610`: `resolveTestDepsAndInjectReviewDoc` — delete, covered by ResolveIfGenerated (source: docs/proposals/pipeline-topology-registry/proposal.md § Functions Relationship to Registry)
- `forge-cli/pkg/task/build.go:502`: `GenerateTestTasks` dispatcher — update to call registry-driven version (source: docs/proposals/pipeline-topology-registry/proposal.md § Derived Functions)
- `forge-cli/pkg/task/build.go:419`: `ResolveDriftFallbackDep` call at step 7.6 — delete, covered by ResolveLastRunTestOrBusiness (source: docs/proposals/pipeline-topology-registry/proposal.md § Functions Relationship to Registry)

## Acceptance Criteria
- [ ] Step 7/7.5 unified: single call to registry-driven `GenerateTestTasks` with businessTasks + existingTasks (index.TasksMap())
- [ ] T-review-doc generated as part of registry pipeline (not separate step), gated by CondHasDocTasks
- [ ] Step 7.6 (`ResolveDriftFallbackDep`) removed — covered by registry resolvers
- [ ] All functions in deletion table deleted: ResolveFirstTestDep, resolveTestDepsAndInjectReviewDoc, findHighestGateOrSummary, findMaxBusinessTaskID, ResolveDriftFallbackDep, GetReviewDocTask, ResolveReviewDocDep, findTaskIndexOrPanic, findTaskIndexByPrefixOrPanic
- [ ] `needsTestPipeline` preserved for stage-gate generation control (build.go step 6.5)
- [ ] `needsReviewDoc` preserved for doc task criteria extraction (build.go step 5.5.2)
- [ ] `go build ./...` passes

## Hard Rules
- 仅修改以下文件：`forge-cli/pkg/task/build.go`, `forge-cli/pkg/task/autogen.go`

## Implementation Notes

### Test Impact
- Affected test suite(s): `forge-cli/pkg/task/build_test.go`, `forge-cli/pkg/task/claim_test.go`
- Expected fixture changes: test expectations for step 7/7.5/7.6 will need updating (deferred to task 6)
- Risk level: high

- The GenerateTestTasks dispatcher (build.go:502) becomes a simple pass-through to the registry version
- needsTestPipeline remains as caller-side gate for stage-gate generation only — its test-pipeline control is replaced by registry IntentGate
- When calling GenerateTestTasks, pass `index.TasksMap()` as existingTasks for ResolveHighestGateOrLastBiz
