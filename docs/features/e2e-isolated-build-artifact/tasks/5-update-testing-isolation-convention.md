---
id: "5"
title: "Update TEST-isolation-004 scope to cover all test locations"
priority: "P2"
estimated_time: "30min"
dependencies: ["1", "2", "3"]
type: "doc"
mainSession: false
---

# 5: Update TEST-isolation-004 scope to cover all test locations

## Description

`docs/conventions/testing-isolation.md` contains TEST-isolation-004 which defines the isolated binary convention, but its scope field does not cover all test locations. After Tasks 1-3 implement the convention everywhere, update the scope to reflect full coverage and verify the description matches the implemented TestMain pattern.

## Reference Files
- `docs/proposals/e2e-isolated-build-artifact/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `docs/conventions/testing-isolation.md` | Update TEST-isolation-004 scope to cover all E2E test locations |

## Acceptance Criteria
- TEST-isolation-004 scope field lists all E2E test modules: `tests/e2e/`, `tests/e2e/justfile-canonical-e2e/`, `forge-cli/tests/e2e/`
- Description matches the implemented TestMain auto-build pattern
- No reference to old PATH-based or shared-path strategies in the convention

## Implementation Notes
- Verify the convention accurately describes the TestMain pattern implemented in Tasks 1-3
