# Test Script Staging vs Graduation Location Confusion

## Problem

During T-test-2 (gen-test-scripts), the subagent wrote spec files to **both**:
- `docs/features/<slug>/testing/scripts/` (staging)
- `tests/e2e/<slug>/` (regression suite)

This caused T-test-3 (run-e2e-tests) to run tests from `tests/e2e/` directly, and T-test-4 (graduate-tests) to find files already there — making graduation a no-op or a confusing gate check.

## Root Cause

Three-level contradiction in the workflow:

1. **SKILL.md** (`gen-test-scripts`) explicitly outputs to `docs/features/<slug>/testing/scripts/` — this is the feature-scoped staging area
2. **Task template T-test-2 AC** checks `tests/e2e/<feature>/` — the regression suite location (T-test-4's job)
3. **Task template T-test-3** references `tests/e2e/<feature>/` as the source for running tests

The subagent saw the AC requiring `tests/e2e/<feature>/` and resolved the contradiction by writing to both locations. This bypassed the graduation gate entirely.

## Solution

Fix the task templates to respect the two-stage design:

- **T-test-2 AC**: check `docs/features/<slug>/testing/scripts/` (staging), not `tests/e2e/`
- **T-test-3**: run tests from `docs/features/<slug>/testing/scripts/`, write results to `docs/features/<slug>/testing/results/`
- **T-test-4**: graduate passing scripts from staging → `tests/e2e/<slug>/`

## Key Takeaway

The `gen-test-scripts` → `graduate-tests` pipeline has two intentionally separate locations:

| Stage | Location | Owner |
|-------|----------|-------|
| Staging (development) | `docs/features/<slug>/testing/scripts/` | T-test-2, T-test-3 |
| Regression suite | `tests/e2e/<slug>/` | T-test-4 only |

When T-test-2 AC or T-test-3 references `tests/e2e/`, subagents will short-circuit graduation by writing there early. Always keep AC and implementation notes aligned with the staging location until T-test-4.
