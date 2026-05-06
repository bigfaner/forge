---
created: "2026-05-05"
source-lesson: pm-work-tracker/docs/lessons/arch-task-executor-missing-e2e-step.md
status: implemented
---

# Harden E2E Testing Lifecycle

## Problem

After implementation tasks modify UI/API components, e2e test breakage goes undetected until the Stop hook's `task all-completed` fires — by then the agent's context is lost and fixing requires re-reading all changes.

### Current architecture (3-layer, with gap)

```
Task-level:      just test [scope]              ← agent has context, but no e2e
     ↓
  [GAP]          feature e2e not verified per task
     ↓
Project-level:   just test-e2e                  ← Stop hook, agent context lost
```

### Causal chain

1. **Symptom**: E2e breakage from UI changes goes undetected until user manually runs `just test-e2e`
2. **Direct cause**: Task quality gate only runs `just test` (unit/integration), not `just test-e2e`
3. **Code cause**: `run-tasks.md` Step 5 defines Breaking Task Gate as `just test` with no e2e step
4. **Process cause**: Testing lifecycle has a gap between task-level and project-level verification
5. **Root cause**: The testing lifecycle was designed assuming e2e changes only happen during T-test-3. But implementation tasks (Phase 2-4) also modify UI/API, breaking e2e specs generated in T-test-1/2.

## Implementation

### 1. Extend Breaking Task Gate (run-tasks.md Step 5)

**File**: `plugins/forge/commands/run-tasks.md`

Split Step 5 into two independent sub-gates with a boolean routing table:

| BREAKING=true? | SCOPE frontend\|all + specs exist? | Run 5a? | Run 5b? |
|----------------|-------------------------------------|---------|---------|
| Yes | No | Yes | No |
| No | Yes | No | Yes |
| Yes | Yes | Yes | Yes |
| No | No | Skip Step 5 entirely | Skip Step 5 entirely |

- **5a** (Unit/Integration Gate): runs `just test [scope]` when `BREAKING: true`
- **5b** (Feature E2E Gate): runs `just test-e2e --feature "$FEATURE"` when SCOPE is frontend|all and specs exist

The dispatcher evaluates SCOPE and FEATURE from claim output (Step 1) before executing any bash commands in 5b. FEATURE is extracted directly from claim output — no separate `task feature` call needed (added to `task claim` output via `PrintFieldIfNotEmpty("FEATURE", featureSlug)` in claim.go).

5b bash block uses a `SKIP` variable pattern for pre-flight checks (recipe existence, spec file existence), avoiding non-portable `continue` in pseudocode context.

**On failure**: `task add` a fix task with unblock instruction, then continue loop.

**Iron Law #3** updated with explicit exception scope: "NO running tests directly — EXCEPT in Step 5 (Breaking Task Gate)".

### 2. T-test-3: Blocked+Fix+Retry Pattern

**File**: `plugins/forge/skills/breakdown-tasks/templates/run-e2e-tests.md`

Change T-test-3 failure handling from "mark completed + add fix tasks" to "mark blocked + add fix tasks + auto-unblock after fix":

**Before** (current):
```
Fail → add fix tasks (P0) → mark completed → fix tasks run → continue to T-test-4
```

**After** (implemented):
```
Fail → add fix tasks (P0, with unblock instruction) → mark blocked → fix tasks run →
fix task unblocks T-test-3 → T-test-3 re-claimed → re-run e2e
```

### 3. New T-test-4.5: Post-Graduation Regression

**File**: New `plugins/forge/skills/breakdown-tasks/templates/verify-regression.md`

Insert between T-test-4 (graduate) and T-test-5 (consolidate-specs). Uses the same blocked+fix+retry pattern as T-test-3. Failure handling format is aligned: multi-line `task add` with `--title "Fix: <concise description>"`, explicit "Do NOT fix inline" warning.

### 4. Simplify T-test-4 (Graduate Only)

**File**: `plugins/forge/skills/breakdown-tasks/templates/graduate-tests.md`

Post-graduation verification removed (moved to T-test-4.5). T-test-4 is now pure graduation:
1. Verify e2e passed (check `latest.md`)
2. Graduate scripts to `tests/e2e/<module>/`
3. Record completion — T-test-4.5 handles regression

### 5. Execute-Task E2E Reminder

**File**: `plugins/forge/commands/execute-task.md`

Added `### E2E Reminder (After Step 3)` sub-section with explicit steps:
1. Check SCOPE is `frontend` or `all`
2. Glob for `tests/e2e/features/<FEATURE>/*.spec.ts`
3. If both true, print reminder: `just test-e2e --feature <FEATURE>`

Informational only for manual `/execute-task` invocation. When dispatched by `/run-tasks`, e2e gates are handled by the dispatcher's Step 5b automatically.

### 6. Updated Task Chain

**File**: `plugins/forge/skills/breakdown-tasks/SKILL.md` (Step 4d)

```
T-test-1 (gen-test-cases)
  → T-test-2 (gen-test-scripts)
    → T-test-3 (run e2e + blocked+fix+retry)   ← modified pattern
      → T-test-4 (graduate)                     ← simplified
        → T-test-4.5 (full regression + blocked+fix+retry)  ← new
          → T-test-5 (consolidate-specs)         ← dependency updated
```

**File**: `plugins/forge/skills/breakdown-tasks/templates/index.json`

T-test-4.5 entry added with dependency on T-test-4. T-test-5 dependency updated from T-test-4 to T-test-4.5.

### 7. Idempotent e2e-setup

**File**: All 6 justfile templates in `plugins/forge/references/justfile-templates/`

No code change needed — `e2e-setup` is already idempotent: `npm install` skips if `node_modules` exists, and `playwright install chromium` skips if the browser is already in the global cache. The recipe can be called unconditionally.

### 8. FEATURE Field in Claim Output

**File**: `task-cli/internal/cmd/claim.go`

Added `PrintFieldIfNotEmpty("FEATURE", featureSlug)` to `printTaskDetails()`. The `featureSlug` parameter was already available — it was only used for building FILE and RECORD paths. Now also emitted as a standalone field so the dispatcher can use it without a separate `task feature` CLI call.

Version bumped 1.3.5 → 1.4.0 (minor: new output field).

## After Architecture

```
Implementation task (SCOPE=frontend, specs exist)
  → just test [scope]                ← Step 5a quality gate
  → just test-e2e --feature $FEATURE ← Step 5b feature e2e gate
  → on failure: task add fix          ← error recovery while context is fresh

T-test-3 (scoped e2e)
  → just test-e2e --feature <slug>
  → fail → blocked → fix → unblock → re-run

T-test-4 (graduate)
  → migrate specs to tests/e2e/<module>/

T-test-4.5 (full regression)
  → just test-e2e
  → fail → blocked → fix → unblock → re-run

T-test-5 (consolidate-specs)
  → extract business rules and tech specs

Stop hook (safety net, unchanged)
```

## Scope Resolution Table

| Condition | Run e2e? | Reason |
|-----------|----------|--------|
| scope=frontend, specs exist | Yes | UI changes may break selectors |
| scope=all, specs exist | Yes | Full-stack changes may break API/UI |
| scope=backend, specs exist | No | Backend changes caught by unit tests |
| scope=frontend, no specs yet | No | T-test-2 hasn't generated specs yet |
| `just test-e2e` recipe missing | No | Project doesn't support e2e |
| FEATURE empty | No | Claim always provides FEATURE (RequireFeature guards) |

## Files Modified

| File | Change |
|------|--------|
| `task-cli/internal/cmd/claim.go` | Add FEATURE field to claim output |
| `task-cli/internal/cmd/claim_test.go` | FEATURE assertion in 3 test functions |
| `task-cli/scripts/version.txt` | Bump 1.3.5 → 1.4.0 |
| `plugins/forge/commands/run-tasks.md` | Step 5: split into 5a/5b with routing table, Iron Law #3 exception |
| `plugins/forge/commands/execute-task.md` | Add e2e reminder after Step 3 |
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Step 4d: add T-test-4.5, update chain (6 tasks) |
| `plugins/forge/skills/breakdown-tasks/templates/run-e2e-tests.md` | T-test-3: blocked+fix+retry |
| `plugins/forge/skills/breakdown-tasks/templates/graduate-tests.md` | T-test-4: simplify (graduate only) |
| `plugins/forge/skills/breakdown-tasks/templates/verify-regression.md` | New: T-test-4.5 template |
| `plugins/forge/skills/breakdown-tasks/templates/index.json` | Add T-test-4.5 entry, update T-test-5 dependency |

## Out of Scope

- **task-executor e2e awareness**: Lesson doc suggested making the executor proactively update e2e specs. Deferred — dispatcher gate catches breakage, error-fixer fixes it. Revisit if fix rate is too high.
