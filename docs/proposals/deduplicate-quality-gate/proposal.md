---
created: 2026-05-19
author: "fanhuifeng"
status: Draft
---

# Proposal: Deduplicate Quality Gate — Single-Authority Test Execution

## Problem

The project-wide unit test suite (`just test`) runs 2-3 times per task. Four independent layers each execute the quality gate without coordination, compounding token and wall-clock cost. E2E tests (`just e2e-test`) are a separate pipeline and unaffected.

### Evidence

| Task Type | `just test` Runs | Layers |
|-----------|-------------------|--------|
| Normal (coding) | 2x | Template verification + CLI submit |
| BREAKING task | 3x | Template verification + Dispatcher gate + CLI submit |

### Root Cause

The CLI quality gate was added on 2026-05-03 (`1ce5b74`) as the **unified authoritative checkpoint** — but the pre-existing layers were never removed:

| Layer | File | What it does | Origin |
|-------|------|-------------|--------|
| Task type template | `coding-*.md` | Agent runs `just test` before submitting | Pre-unification |
| submit-task skill | `SKILL.md` lines 88-103 | Agent runs `just test` for metrics collection | Pre-unification |
| CLI quality gate | `submit.go` line 139 | `validateQualityGate()` runs `compile→fmt→lint→test` | 2026-05-03 unification |
| Dispatcher breaking gate | `execute-task.md` Step 3a | Dispatcher runs `just test` for BREAKING tasks | Pre-unification |

The CLI gate was designed to be the sole authority. The other three layers are leftovers.

## Proposed Solution

Keep the CLI quality gate as the **sole authority**. Remove `just test` from the other three layers.

### What Changes

| Layer | Change | Rationale |
|-------|--------|-----------|
| Task type templates (5 files) | Remove `just test` step; keep `compile`+`lint` for fast feedback | CLI gate catches test failures; agent gets faster compile/lint feedback |
| submit-task skill | Remove all `just *` invocations (quality gate + metrics collection) | CLI gate handles quality; metrics parsing is a separate concern |
| Dispatcher breaking gate | Remove Step 3a entirely | CLI gate runs immediately after in submit |
| CLI quality gate (`submit.go`) | **No change** — remains sole authority | This is the unification that was always intended |
| `validateRecordData()` | Skip "no test evidence" check when CLI gate passed | CLI gate passing `just test` IS the test evidence |

### What Does NOT Change

- Quality gate sequence: `compile→fmt→lint→test` (unchanged)
- `validateQualityGate()` function (unchanged)
- E2E test pipeline (`just e2e-test`, separate system)
- Agent still writes record.json — but with zero metrics (testsPassed=0, etc.)
- `--force` bypass behavior (unchanged)

### Metrics

Records will have `testsPassed=0, testsFailed=0, coverage=0` for coding tasks. This is acceptable because:
- The CLI quality gate already verified tests pass (exit code 0)
- Per-framework CLI metrics parsing is a separate optimization, not part of this proposal
- The "no test evidence" validation in `validateRecordData()` must be adjusted to accept zero metrics when the CLI gate has passed

## Scope

### In Scope

- Remove `just test` from 5 coding type templates (`coding-feature.md`, `coding-enhancement.md`, `coding-cleanup.md`, `coding-fix.md`, `coding-refactor.md`)
- Remove `just test` from `gate.md` and `validation-code.md` templates
- Remove all `just *` invocations from submit-task `SKILL.md`
- Remove Step 3a breaking gate from `execute-task.md`
- Adjust `validateRecordData()` in `submit.go`: skip "no test evidence" check when CLI quality gate has passed

### Out of Scope

- CLI metrics parsing from test output (separate proposal)
- Adding/removing frameworks
- Changing quality gate sequence
- E2E test pipeline changes
- Profile system changes

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Agent submits broken code that passes compile+lint but fails tests | L | L | CLI quality gate catches it. Agent just needs to re-submit after fix. |
| Records show zero metrics, losing per-task test data | H | L | Metrics were always project-wide (not task-specific). CLI gate pass = test evidence. Metrics parsing can restore data later. |
| `validateRecordData()` change allows false completions | L | M | Only skip check when CLI gate explicitly passed (gate ran and returned true), not when skipped (--force or non-testable type). |

## Success Criteria

- [ ] `just test` runs exactly once per task submission (CLI gate only)
- [ ] BREAKING tasks run tests exactly once
- [ ] No regression in quality gate pass/fail accuracy
- [ ] `validateRecordData()` accepts zero metrics when CLI gate passed
- [ ] Agent no longer runs any `just *` commands in submit-task flow

## Next Steps

- Proceed to `/write-prd` to formalize requirements
