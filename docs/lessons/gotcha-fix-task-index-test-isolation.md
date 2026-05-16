---
created: "2026-05-16"
tags: [testing, architecture]
---

# Fix Tasks in Index Break Integration Test Isolation

## Problem

After running `forge quality-gate` (stop hook), the test suite in `forge-cli/internal/cmd` starts failing with errors like:

- `TestAddCmd_DedupSkipsActiveFix`: "active fix tasks already exist for source 1.1: fix-4"
- `TestExecuteClaim`: claims `fix-2` instead of expected `task1`
- `TestRunFeature_Display`: shows `simplify-e2e-tests` instead of `test-feature`
- `TestFeatureList_WithFeatures`: lists 33 real features instead of test fixtures

Each hook cycle adds another fix task to the index, making more tests fail, triggering more fix tasks.

## Root Cause

Causal chain (4 levels):

1. **Symptom**: 14 tests fail in `forge-cli/internal/cmd` after quality-gate hook runs
2. **Direct cause**: Fix tasks (fix-2, fix-3, fix-4, ...) exist in `docs/features/<slug>/tasks/index.json`, and `.forge/state.json` points to an active feature with `allCompleted: false`
3. **Root cause**: Integration tests in `forge-cli/internal/cmd` use `exec.Command(os.Args[0], ...)` to spawn subprocesses or call `rootCmd.Execute()` directly. Some code paths (dedup check, claim logic, feature resolution) resolve project state from the real filesystem instead of the test's temp dir — they see the fix tasks and active feature from the real project
4. **Trigger condition**: Any quality-gate failure that creates fix tasks in the active feature's index, while the active feature's `.forge/state.json` has `allCompleted: false`

## Solution

Clean up stale fix tasks and project state:

1. Remove orphan fix task `.md` files: `rm docs/features/<slug>/tasks/fix-N.md`
2. Remove orphan entries from `index.json` (they survive `forge task index` re-runs)
3. Remove `.forge/state.json` (stale `allCompleted: false` triggers the hook)
4. Run `forge task index --feature <slug>` to rebuild clean index
5. Verify with `just test`

## Reusable Pattern

**When quality-gate creates fix tasks and tests start failing in `forge-cli/internal/cmd`, first check if the failures are caused by index pollution before investigating code changes.**

Diagnostic pattern:
1. Check if fix tasks exist: `ls docs/features/*/tasks/fix-*.md`
2. Check `.forge/state.json`: `cat .forge/state.json`
3. If fix tasks exist and state shows `allCompleted: false`, clean them up first
4. Re-run tests — if they pass, the fix tasks were the cause, not code changes

Prevention: after completing a feature's tasks, always ensure `.forge/state.json` is cleared (`allCompleted: true` + consumed by `forge all-completed`) before ending the session.

## Example

```bash
# Diagnostic: check for fix task accumulation
ls docs/features/simplify-e2e-tests/tasks/fix-*.md
# fix-1.md  fix-2.md  fix-3.md  fix-4.md  fix-5.md

# Clean up stale fix tasks (keep fix-1 if it was legitimate)
rm docs/features/simplify-e2e-tests/tasks/fix-{2,3,4,5}.md

# Remove orphan entries from index
node -e "
const d = JSON.parse(require('fs').readFileSync('.../index.json','utf8'));
Object.keys(d.tasks).filter(k => k.startsWith('fix-') && k !== 'fix-1')
  .forEach(k => delete d.tasks[k]);
require('fs').writeFileSync('.../index.json', JSON.stringify(d,null,2));
"

# Clear stale state
rm .forge/state.json

# Verify
forge task index --feature simplify-e2e-tests
just test
```

## Related Files

- `forge-cli/internal/cmd/integration_test.go` — `setupFullProject` test isolation helper
- `forge-cli/internal/cmd/add_cmd_test.go` — dedup tests that pick up real fix tasks
- `forge-cli/internal/cmd/claim_test.go` — claim tests that see real feature state
- `forge-cli/internal/cmd/feature_test.go` — feature list tests that see real features
- `docs/lessons/gotcha-quality-gate-fix-task-loop.md` — related: stale output file loop
- `docs/lessons/gotcha-quality-gate-cross-feature-pollution.md` — related: fix tasks polluting unrelated features
