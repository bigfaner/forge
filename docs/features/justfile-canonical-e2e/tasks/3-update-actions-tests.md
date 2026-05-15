---
id: "3"
title: "Update e2e/actions_test.go for just delegation"
priority: "P0"
estimated_time: "1.5h"
dependencies: ["2"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 3: Update e2e/actions_test.go for just delegation

## Description

Rewrite `actions_test.go` tests to assert the new `exec.Command("just", ...)` delegation behavior. Replace profile-specific command assertions (go test, npx playwright) with generic `just <recipe>` assertions. Add tests for `just` not-found error and exit code propagation.

## Reference Files
- `docs/proposals/justfile-canonical-e2e/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/e2e/actions_test.go` | Rewrite test assertions for `just` delegation; add exit code propagation tests; add `just` not-found test |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] TestRun asserts `exec.Command` called with `["just", "test-e2e"]` (or with feature arg `["just", "test-e2e", "feature=X"]`)
- [ ] TestSetup asserts `["just", "e2e-setup"]`
- [ ] TestCompile asserts `["just", "e2e-compile"]`
- [ ] TestDiscover asserts `["just", "e2e-discover"]`
- [ ] TestVerify unchanged (no just delegation)
- [ ] Non-zero `just` exit produces non-nil error; zero exit produces nil error
- [ ] `just` not on PATH produces clear error message
- [ ] `go test ./pkg/e2e/...` passes with 80%+ coverage on modified files
- [ ] No test references specific profile names for Run/Setup/Compile/Discover (they are now profile-agnostic)

## Hard Rules

- Keep the existing `stubExec` mock pattern. Do not introduce testify or gomock.
- The `setupProfile` and `setupGoTestProfile` helpers can be simplified — only `setupProfile` is needed since profile name no longer affects command dispatch.

## Implementation Notes

- The existing test structure has subtests per-profile (go-test dispatches `go test`, web-playwright dispatches `npx playwright test`). These collapse into a single assertion per function: "just test-e2e is called". The profile-specific subtests can be merged or simplified.
- The "unsupported profile" error subtests are no longer needed — delegation to `just` is profile-agnostic.
- The "external tool failure" subtests remain relevant — they test non-zero exit code propagation, which still applies.
- Coverage target is 80%+ on `actions.go`. The current tests likely exceed this; the rewrite should maintain or improve coverage.
