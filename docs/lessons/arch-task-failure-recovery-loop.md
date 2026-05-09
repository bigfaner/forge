# Task Execution Must Block on Failures and Spawn Fix Tasks

## Problem

When `run-e2e-tests` (T-test-3/T-test-4) encounters failing test cases, the task executor:

1. Records the result with `testsPassed=N, testsFailed=M` but `status=completed`
2. Proceeds to the next task (e.g., `graduate-tests`) as if nothing is wrong
3. Never creates `fix-*` tasks to address the failures

The consequence: downstream tasks operate on a broken foundation (e.g., graduating tests that don't pass), and failures silently compound through the pipeline.

**Smoking gun**: `docs/features/e2e-test-scripts-rebuild/tasks/process/record.json` contains `testsPassed=3, testsFailed=3, status=completed`. Agent rationalized: "3 failures are environment-only constraints." No fix tasks created. Pipeline continued.

## Root Cause

**6 gaps across the forge pipeline (causal chain — failure flows through all 6):**

### Gap 1: task-executor agent — linear workflow, no post-record check

**File**: `~/.claude/plugins/cache/forge/forge/2.14.0/agents/task-executor.md` (lines 38-139)

The task-executor runs 5 steps linearly: read task → TDD implement → verify → record → commit. Step 3 has a failure loop ("fix tests then retry"), but once it passes verification, Steps 4-5 execute unconditionally. After `task record` writes a record with `testsFailed > 0`, the agent still commits and outputs `DONE`.

**Gap**: The executor treats recording as fire-and-forget. It trusts `task record` to reject bad data but never verifies the outcome.

### Gap 2: run-tasks dispatcher — existence-only check

**File**: `~/.claude/plugins/cache/forge/forge/2.14.0/commands/run-tasks.md` (lines 108-113)

Step 3 "Verify Record" checks only whether `records/*.md` **file exists**. It does NOT parse the record to check `testsFailed`, `testsPassed`, or acceptance criteria. After this existence check passes, the dispatcher proceeds to the breaking gate (Step 5) — which runs `just test` only for tasks marked `breaking: true`, and only AFTER the task is already marked "completed" in `index.json`.

**Gap**: "Verify record" is a filesystem check, not a quality gate.

### Gap 3: task record CLI — allows completed + testsFailed > 0 (ROOT CAUSE)

**File**: `~/.zcode-task-cli/task` (Go binary)

The CLI validates:
- `completed + testsPassed=0 + testsFailed=0 + coverage>=0` → rejected (no test evidence)
- `completed + any acceptanceCriteria.met=false` → rejected (unmet AC)

But there is **NO validation** for `completed + testsFailed > 0`. Partial failure (e.g., 3 passed, 3 failed) passes cleanly. The binary contains no string matching "completed with test failures".

**Gap**: This is the ROOT CAUSE. The CLI's validation gap allows recording a task as "completed" when tests are failing. Everything downstream trusts this status.

### Gap 4: run-e2e-tests skill — reports but doesn't block

**File**: `~/.claude/plugins/cache/forge/forge/2.14.0/skills/run-e2e-tests/SKILL.md` (lines 235-254)

The skill faithfully reports failures (`E2E Test Results: X/Y passed (Z failed)`), but does NOT write `record.json` or call `task record` — that is left to the invoking agent. The skill is prompt-based, so it relies on the agent reading `latest.md` and making the right decision.

**Gap**: "Report faithfully" is not the same as "block the pipeline." The skill has no mechanical enforcement.

### Gap 5: record-task skill — documents the gap

**File**: `~/.claude/plugins/cache/forge/forge/2.14.0/skills/record-task/SKILL.md` (lines 122-133)

The skill documents 3 validation rules enforced by the CLI. There is NO row for `status=completed + testsFailed > 0`. Agents reading this skill see the 3 rules and conclude anything not listed is acceptable.

**Gap**: The skill documentation reinforces the CLI gap by omission.

### Gap 6: T-test templates — correct words, no enforcement

**File**: `~/.claude/plugins/cache/forge/forge/2.14.0/skills/breakdown-tasks/templates/run-e2e-tests.md` (lines 43-65)

T-test-3 correctly describes failure handling: "Mark this task blocked, create fix task, re-claim after fix." The acceptance criteria say "All tests pass (status = PASS in latest.md)."

But enforcement relies entirely on **agent compliance**. The agent self-reports `met: true/false`. If it rationalizes failures as "environment-only constraints", the CLI only checks the `met` boolean, not the actual test results. There is no cross-check between `testsFailed > 0` and AC claiming "all tests pass" with `met: true`.

**Gap**: Correct policy, no mechanical guardrail. Agent discretion overrides system safety.

### Complete Failure Path

```
Agent runs tests → some tests fail
  → Agent rationalizes ("environment-only", "out of scope")
    → Agent writes record.json: testsFailed>0, status=completed
      → task record CLI accepts it (only rejects testsPassed=0 AND testsFailed=0)
        → index.json updated to status=completed
          → run-tasks dispatcher checks file existence only, continues
            → Next task claimed and dispatched
              → Failure silently compounds downstream
```

## Solution

### Fix #1 (CLI — highest priority): Add validation rule to `task record`

```
status=completed + testsFailed > 0 → reject with:
  "Cannot mark task completed with %d test failures. Fix failures or set status to 'blocked'."
  Override with --force.
```

This single change breaks the failure path at Gap #3 and forces all upstream components to handle failures correctly.

### Fix #2 (task-executor agent): Post-record verification

After `task record` succeeds, re-read the record. If `testsFailed > 0`, do NOT commit. Escalate to user.

### Fix #3 (run-tasks dispatcher): Content quality gate

Step 3 should parse `process/record.json` after the agent completes:
- If `testsFailed > 0` or `status != "completed"` → mark task blocked, spawn fix task
- Only proceed if `status == "completed" AND testsFailed == 0`

### Fix #4 (T-test templates): Mechanical AC cross-check

Add a rule: if `testsFailed > 0`, the CLI must reject any AC containing "all tests pass" with `met: true`. Cross-check numeric test results against AC claims.

## Key Takeaway

**The pipeline has 6 independent gaps. Any single fix helps; all 4 fixes are needed for a robust safety net.**

Priority order:
1. **Fix #1 (CLI validation)** — root cause, blocks the failure path mechanically
2. **Fix #3 (dispatcher quality gate)** — catches failures the CLI misses
3. **Fix #2 (agent post-record check)** — defense in depth
4. **Fix #4 (AC cross-check)** — prevents agent self-deception

**How to apply**: When tests fail, the agent must set `status: "blocked"` in record.json. The CLI must reject `completed + testsFailed > 0`. The dispatcher must parse record content, not just check file existence.

**Why**: `pattern-phase-gate-tasks.md` covers *between-phase* gates. This lesson covers the *within-task* failure path. `gotcha-ac-self-report-without-verification.md` covers the AC trust issue. Together, these three lessons form the complete safety net picture.
