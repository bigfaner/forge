---
created: "2026-05-19"
tags: [testing, architecture]
---

# Gotcha: "docs-only" proposals need code path audit

## Problem

A proposal classified a feature as "docs-only — no Go code changes needed" based on surface-level verification. During task generation, the reviewer noticed Go functions like `IsTestableType` that should be affected. Two actual Go code gaps were found:

1. `testableTypes` map only had `{feature, enhancement, fix}` — missing `cleanup` and `refactor`
2. `submit.go` quality-gate skip logic was inconsistent with the new type-based approach — it relied on the old field instead of using `IsTestableType`

## Root Cause

The proposal author verified that `isDocsOnlyFeature()` already used `IsTestableType` (correct), but stopped there. They did not audit:

1. Whether `testableTypes` included all types listed in the documentation's classification table
2. Whether all code paths using type for decision-making were consistent with the documented behavior (some code paths still used the old mechanism exclusively)

The assumption was "Go code already implements what the docs describe" without verifying each code path.

## Solution

1. Added `TypeCleanup` and `TypeRefactor` to `testableTypes` map
2. Changed submit.go quality-gate skip to use `task.IsTestableType(t.Type)` consistently
3. Updated proposal to include Go code changes in scope

## Reusable Pattern

When proposing a "docs-only" feature that references existing code behavior (e.g., "Go code already does X"), audit every code path that implements the referenced behavior — not just the most obvious one. Specifically:

- Grep for all usages of the referenced function/constant/map
- Verify each usage site matches the proposed documentation
- Check for gaps between documented type classification and actual code (maps, switches, conditionals)

## Example

```
# Proposal claimed: "isDocsOnlyFeature() 已正确检查 type"
# But grep revealed:
#   - build.go: testableTypes only had 3 of 5 code types
#   - submit.go: quality-gate skip used old mechanism, not IsTestableType
#   - quality_gate.go: isDocsOnly() correctly used IsTestableType ✓
```

## Related Files

- `forge-cli/pkg/task/build.go` — `testableTypes` map
- `forge-cli/internal/cmd/submit.go` — quality-gate skip logic
- `plugins/forge/references/shared/type-assignment.md` — type classification rule

## References

- Proposal: `docs/proposals/task-type-code-docs-boundary/proposal.md`
