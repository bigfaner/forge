---
id: "5"
title: "Extend TC-004 contract to cover forge-cli/tests/"
priority: "P2"
estimated_time: "15m"
dependencies: [1, 2, 3, 4]
type: "doc"
mainSession: false
---

# 5: Extend TC-004 contract to cover forge-cli/tests/

## Description
Update the `test-suite-health` contract to specify that TC-004 ("Zero unconditional t.Skip") should also cover `forge-cli/tests/` integration tests, not just `tests/` e2e tests.

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `tests/test-suite-health/contracts/step-1-test-suite-health.md` | Add `forge-cli/tests/` to TC-004 scope |

### Delete
| File | Reason |
|------|--------|

### Create
| File | Description |
|------|-------------|

## Reference Files
- `docs/proposals/clean-invalid-tests/proposal.md#Scope` — In Scope item 14: extend TC-004 coverage
- `docs/proposals/clean-invalid-tests/proposal.md#Constraints-&-Dependencies` — "TC-004 规则目前只覆盖 tests/，清理后应扩展到 forge-cli/tests/"
- `tests/test-suite-health/contracts/step-1-test-suite-health.md` — current contract defining TC-004 scope

## Acceptance Criteria
- [ ] Contract file updated to include `forge-cli/tests/` in TC-004 scope
- [ ] "Zero unconditional t.Skip() calls" assertion explicitly lists both `tests/` and `forge-cli/tests/` as target directories

## Hard Rules
- Do NOT modify the Go meta-test file (`tests/test-suite-health/e2e_test_quality_cleanup_test.go`) — that is out of scope per proposal

## Implementation Notes
- The contract currently says "E2E test suite in tests/e2e/ and tests/*/ directories" — extend to also include `forge-cli/tests/`
- Keep the change minimal: just expand the scope description, don't rewrite the contract
