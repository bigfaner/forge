---
created: "2026-05-18"
tags: [architecture, testing]
---

# Duplicate Test Runs in Task Execution Pipeline

## Problem

Tests run 2-3 times per task, wasting significant token and wall-clock time. For a typical feature task, `just test` executes at least twice; for BREAKING tasks, at least three times.

## Root Cause

Three independent layers each run the full quality gate (`compile → fmt → lint → test`) without coordination:

1. **Task type template** (e.g., `forge-cli/pkg/prompt/data/feature.md` line 46-50): Every synthesized strategy includes a final "Full Verification" step that runs `just test`. This is baked into all type templates (feature, enhancement, cleanup, refactor, fix).

2. **submit-task skill metrics collection** (`plugins/forge/skills/submit-task/SKILL.md` lines 88-103): The HARD-RULE instructs the agent to run `just test` and manually parse stdout text to extract `testsPassed`, `testsFailed`, `coverage` — no structured output, just the agent reading raw console text and guessing numbers.

3. **`forge task submit` CLI** (`forge-cli/internal/cmd/submit.go` line 137): `validateQualityGate()` unconditionally runs `DefaultGateSequence()` (compile → fmt → lint → test). It only checks pass/fail via exit code — does NOT parse or extract any metrics from the output.

4. **Dispatcher breaking gate** (`plugins/forge/commands/execute-task.md` Step 3a): For BREAKING tasks, the dispatcher runs `just test` again after the subagent returns.

None of these layers are aware of what the others have already executed. Each layer independently decides "I must verify before proceeding."

## Solution

**Short-term**: Remove the test step from task type templates. The `forge task submit` CLI is the authoritative quality gate — it always runs and cannot be skipped by the agent. The template's verification step is redundant because:
- If the agent's code is bad, `forge task submit` will catch it
- If the agent tries to skip tests, the CLI blocks the submission

**Medium-term**: Let `validateQualityGate()` parse metrics from the test output it already captures. The CLI already uses `RunCapture` which returns `(output string, success bool)` — currently it only checks `success` (exit code) and discards `output`. Instead, parse `output` with per-framework regex to extract `testsPassed`, `testsFailed`, `coverage`, and return them to the caller. This eliminates the need for the agent to run `just test` separately for metrics collection.

No special flags or file outputs needed — parse the **default console output** that each framework already produces:

| Framework | Default output pattern | Extractable metrics |
|-----------|----------------------|-------------------|
| go-testing | `ok  pkg  1.2s  coverage: 85.6%` | pass/fail counts, coverage % |
| pytest | `5 passed, 1 failed in 2.3s` | pass/fail/skip counts |
| jest | `Tests: 12 passed, 1 failed` | pass/fail counts, coverage % |
| junit5 | `Tests run: 12, Failures: 1` | pass/fail/error counts |

The `profile.FrameworkInfo` registry (`forge-cli/pkg/profile/framework.go`) already knows 6 frameworks — add an output parser per framework there, no new infrastructure needed.

## Reusable Pattern

When multiple layers of a pipeline each independently verify the same property, they create compounding overhead. Designate **one authoritative checkpoint** and make all other layers trust it. In forge's case: `validateQualityGate()` should be the single point that both gates quality AND collects metrics — one test run, two outcomes.

## Related Files

- `forge-cli/pkg/prompt/data/feature.md` — template with redundant test step
- `forge-cli/pkg/prompt/data/enhancement.md` — same pattern
- `forge-cli/pkg/prompt/data/cleanup.md` — same pattern
- `forge-cli/pkg/prompt/data/fix.md` — same pattern
- `forge-cli/internal/cmd/submit.go` — authoritative quality gate (line 137)
- `forge-cli/pkg/just/just.go` — DefaultGateSequence (line 23), RunCapture (line 57) already returns output
- `forge-cli/pkg/profile/framework.go` — framework registry (6 frameworks), needs output parsers
- `plugins/forge/skills/submit-task/SKILL.md` — metrics collection with `just test`
- `plugins/forge/commands/execute-task.md` — dispatcher breaking gate
