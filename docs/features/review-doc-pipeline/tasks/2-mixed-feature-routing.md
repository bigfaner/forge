---
id: "2"
title: "Modify build.go routing for mixed features"
priority: "P0"
estimated_time: "1.5h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 2: Modify build.go routing for mixed features

## Description

Change `build.go` routing logic so mixed features (containing both `doc` and `coding.*` tasks) generate both review-doc AND test pipeline tasks. Currently the routing is mutually exclusive: `needsTestPipeline()` and `needsDocEval()` never both return true.

The new behavior:
- **Pure doc features**: generate T-review-doc only (no test pipeline)
- **Pure code features**: generate test pipeline only (no review-doc)
- **Mixed features**: generate T-review-doc + test pipeline tasks, with T-review-doc executing before gen-journeys

## Reference Files
- `docs/proposals/review-doc-pipeline/proposal.md` — Source proposal

## Acceptance Criteria

- [ ] `needsReviewDoc()` (renamed in Task 1) returns true when ANY non-auto-gen task has type `doc` (not only when ALL are doc)
- [ ] `needsTestPipeline()` unchanged — returns true when any task has testable type
- [ ] Both `needsReviewDoc` and `needsTestPipeline` can return true simultaneously for mixed features
- [ ] For pure doc features: `needsReviewDoc=true`, `needsTest=false` → only T-review-doc generated
- [ ] For pure code features: `needsReviewDoc=false`, `needsTest=true` → only test pipeline generated
- [ ] For mixed features: `needsReviewDoc=true`, `needsTest=true` → both T-review-doc and test pipeline tasks generated
- [ ] When both are generated for mixed features, T-review-doc is a dependency of T-gen-journeys (review-doc executes before test generation)
- [ ] `forge task index` produces correct `index.json` for all three scenarios
- [ ] Existing stage-gate logic (when `needsTest=true`) remains unchanged
- [ ] `go build ./...` and existing tests pass

## Hard Rules

- Do NOT modify `needsTestPipeline()` logic — only change `needsReviewDoc()`
- The dependency chain for mixed features: T-review-doc → T-gen-journeys → T-eval-journey → T-gen-contracts → T-eval-contract → T-gen-test-scripts
- Stage gates (.summary.md, .gate.md) are tied to the test pipeline only — review-doc does not create stage gates

## Implementation Notes

Key changes in `build.go` `BuildIndex` function:

1. `needsReviewDoc()` logic change:
   - Old: `return true` only when ALL business tasks are `TypeDoc`
   - New: `return true` when ANY business task is `TypeDoc`
   - Use `CategoryOf(t.Type) == CategoryDoc` to check (this covers `doc`, `doc.review`, etc.)

2. Remove mutual exclusivity in routing:
   - Old: `if needsEval { ... doc-eval path ... } else if needsTest { ... test pipeline ... }`
   - New: `if needsReview { ... generate T-review-doc ... }` then `if needsTest { ... generate test pipeline ... }` (both can run)

3. Dependency injection for mixed features:
   - When both `needsReview` and `needsTest` are true, the T-gen-journeys task should depend on T-review-doc
   - Update the dependency resolution in the test pipeline section to include T-review-doc as a prerequisite
