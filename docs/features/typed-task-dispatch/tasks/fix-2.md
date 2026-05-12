---
id: "fix-2"
title: "Fix: e2e test failures in typed-task-dispatch cli.spec.ts (16/20 failed)"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
---

# Fix: e2e test failures in typed-task-dispatch cli.spec.ts (16/20 failed)

## Root Cause

16/20 e2e tests failed. Root causes: (1) Wrong skill file paths — tests reference plugins/forge/skills/*.md but actual paths differ; (2) task prompt exits 0 for missing/invalid type instead of non-zero; (3) task prompt --fix-record-missed leaks TDD steps; (4) phase boundary detection not injecting summary path; (5) task prompt exits 0 when .forge/state.json missing; (6) task migrate blocked by active feature state; (7) task validate schema mismatch. Fix the spec file paths and verify CLI behavior matches test expectations.

## Reference Files

- Source: tests/e2e/features/typed-task-dispatch/cli.spec.ts
- Test script: tests/e2e/features/typed-task-dispatch/cli.spec.ts
- Test results: tests/e2e/features/typed-task-dispatch/results/latest.md

## E2E Fix Boundaries

When fixing E2E test failures, observe these boundaries:

**Forbidden:**
- Starting dev server (`npx expo start`, `npm run dev`, etc.)
- Running `npm install` more than 3 times — mark task as blocked if dependency installation fails 3 times
- Running e2e tests (`just test-e2e`) — regression is verified by the dispatcher after fix completes
- Manually opening browser to verify rendering

**Correct workflow:**
1. Read failing test + corresponding component source
2. Compare test's expected testID/selectors vs actual DOM structure
3. Modify component (add testID) or test (adjust selectors/assertions)
4. `just test` — unit tests must pass
5. Record completion

## Verification

After fixing, verify the fix works:
1. `just test [scope]` — must pass

E2e regression is verified by the dispatcher, not by this fix task.

When this task is recorded as completed via `task record`, the source task T-test-3 is automatically restored to pending if all its dependencies are completed.
