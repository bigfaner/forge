---
created: "2026-05-09"
tags: [testing, local-dev-deployment]
---

# Don't Repeat Verification When Unit Tests Already Prove It

## Problem

After implementing task-cli changes (adding `Proposal` field + quick-mode validation), all verification was complete:

1. `go test -race -cover` — all tests pass (including 4 new quick-mode tests)
2. `task validate` on a quick-mode index.json — PASS, warnings eliminated
3. `make check-docs` — PASS

Then proceeded to create temp directories and write manual e2e test files to "verify end-to-end". User interrupted — the process appeared stuck doing unnecessary work.

## Root Cause

**Symptom**: User interrupted two Write calls creating temp test files for manual e2e verification.

**Direct cause**: Proceeded to manual e2e verification without checking whether existing tests already covered the scenario.

**Root cause**: Habit of "prove it works end-to-end" without evaluating whether unit tests + a single smoke test already provide sufficient confidence. The new `TestValidator_QuickMode` tests cover:
- Quick-mode suppression of prd/design/summary warnings
- Proposal field deserialization
- Marshal/unmarshal round-trip
- Full-mode with all fields present

**Trigger**: When tests pass AND a representative manual smoke test passes, stop. Don't escalate to full manual e2e setup.

## Solution

Stopped the redundant manual verification. The changes were already verified by:
- 4 new unit tests in `validate_test.go`
- `task validate` smoke test on quick-mode index.json
- `make check-docs` for doc freshness

## Reusable Pattern

**Trust the test pyramid.** When:

1. Unit tests cover the new behavior (including edge cases)
2. One representative smoke test confirms real-world behavior (`task validate` on actual file)
3. Doc freshness checks pass

→ **Stop.** Don't create temp directories and write manual e2e test files. The additional manual verification provides zero additional confidence while costing user patience.

**Exception**: If the change involves file I/O paths, CLI argument parsing, or cross-process communication that unit tests can't easily cover, one smoke test is warranted — but keep it to a single `task validate` call, not a full temp directory setup.

## Example

```bash
# Sufficient verification chain:
go test -race -cover ./pkg/task/... ./internal/cmd/...   # Unit tests
task validate /path/to/quick-mode-index.json              # One smoke test
make check-docs                                           # Doc freshness

# STOP HERE. Don't proceed to:
mkdir -p /tmp/e2e-test/...
echo "..." > /tmp/e2e-test/proposal.md
# ... unnecessary manual setup
```

## Related Files

- `task-cli/internal/cmd/validate.go` — Validation logic with quick-mode support
- `task-cli/internal/cmd/validate_test.go` — 4 quick-mode test cases
- `task-cli/pkg/task/types.go` — `Proposal` field addition
