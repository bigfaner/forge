# Surface: web — Web E2E Test Orchestration

This rule file defines the Web E2E test orchestration sequence for the web surface in the run-tests skill. The consumer is the SKILL.md dispatcher.

## Orchestration Sequence

| Step | just Recipe | Exit Code 0 | Exit Code 1 | Exit Code 2 | Next Action |
|------|------------|-------------|-------------|-------------|-------------|
| dev | `just <recipe-prefix>-dev` | Service started successfully, waiting for readiness | Startup failed (missing dependencies / port in use) | — | Proceed to probe |
| probe | `just <recipe-prefix>-probe` | Health check passed | Health check timed out (service not ready) | — | Proceed to test |
| test | `just <recipe-prefix>-test <journey>` | Web E2E tests passed | Web E2E tests failed | Test environment error (retryable) | Proceed to teardown |
| teardown | `just <recipe-prefix>-teardown` | Cleanup complete | Cleanup failed (residual processes) | — | End |

## Probe Retry Strategy

- Maximum 3 retries with 5-second intervals
- If all 3 attempts fail, treat as exit code 1 (retryable)

## Failure Handling

### dev failure

When dev exits non-zero, **do not continue** with subsequent steps; proceed directly to teardown and exit with dev's exit code.

### probe failure (HARD-GATE)

<HARD-GATE>
After probe fails, within the same orchestration cycle:
- **MUST NOT** retry probe (retries are handled by the probe retry strategy within limits, not cycle-level retries)
- **MUST NOT** restart dev
- MUST execute teardown before exiting
</HARD-GATE>

After probe ultimately fails:
- Exit code 1 (retryable): Execute teardown, exit with exit code 1
- Exit code 2 (blocking): Execute teardown, exit with exit code 2

### test failure

- Exit code 1: Execute teardown, exit with exit code 1
- Exit code 2 (retryable): Execute teardown, prompt the user "Test environment error, consider retrying", exit with exit code 2

### teardown failure

When teardown fails, log the error and preserve `.forge/test-state.json` for recovery. Exit with the current step's exit code.

## Journey Filter

| Tag | Match Rule |
|-----|-----------|
| `@web` | Exact match |
| `@web-e2e` | Exact match |
| `@smoke` | Exact match |

## Per-Journey Execution

The dev/probe lifecycle for web surface wraps all journey tests. Use the `recipe-prefix` determined in SKILL.md Step 1 (for single-surface projects, the surface-type "web"; for multi-surface projects, the surface-key) to construct recipe names:

```
just <recipe-prefix>-dev
just <recipe-prefix>-probe (with retry)
for each journey in JOURNEYS:
    just <recipe-prefix>-test <journey>
    record results
    on failure: just <recipe-prefix>-teardown, exit
just <recipe-prefix>-teardown
```

dev and probe execute once, test runs in a per-journey loop, teardown executes once. The test recipe invocation format is `just <recipe-prefix>-test <journey>`, where `<journey>` is a directory name discovered from `docs/features/<slug>/testing/`. `<recipe-prefix>` is "web" for single-surface projects, or the corresponding surface-key for multi-surface projects.

## Suite Name

Test report suite names use the `web-e2e/<journey-name>` format.
