# Surface: web

> **Test Type Reference**: The test type for web surface is **Web E2E Test**, which verifies DOM element visibility + user interaction response + page URL changes + element attribute values via browser automation.

## Orchestration Sequence

| Step | Exit Code 0 | Exit Code 1 | Exit Code 2 | Next Action |
|------|-------------|-------------|-------------|-------------|
| dev | Service started successfully, waiting for readiness | Startup failed (missing dependencies / port in use) | — | Proceed to probe |
| probe | Health check passed | Health check timed out (service not ready) | — | Proceed to test |
| test | Tests passed | Tests failed | Test environment error (retryable) | Proceed to teardown |
| teardown | Cleanup complete | Cleanup failed (residual processes) | — | End |

Notes:
- When dev fails, **do not continue** with subsequent steps; proceed directly to teardown (which is safe/idempotent — no process to clean if dev never started) and exit
- Probe retries up to 3 times with 5-second intervals; if all 3 attempts fail, treat as exit code 1
- Exit code 2 for test step allows re-running; the skill should prompt the user "Test environment error, consider retrying"

## Recipe Invocation Contract

> **Naming convention**: Recipe names below use the surface type (`web-`) as prefix for illustration. For **named surfaces**, replace the type prefix with the surface key (e.g., `web-dev` → `frontend-dev` for `frontend=web`). For **scalar surfaces**, the prefix is omitted (e.g., `web-dev` → `dev`). See SKILL.md Standard Target Contract for the `<prefix>` definition.

| Recipe Name | just Signature | Exit Code 0 Semantics | Exit Code 1 Semantics |
|-------------|---------------|----------------------|----------------------|
| web-dev | `just web-dev` | Development server ready, listening on port | Startup failed, stderr contains error details |
| web-probe | `just web-probe` | HTTP health check returns 2xx | Connection refused or timed out |
| web-test | `just web-test [journey]` | All Web E2E tests passed | At least one test failed |
| web-teardown | `just web-teardown` | Processes terminated, port released | Residual processes or cleanup error |
| web | `just web` | Aggregate recipe: dev->probe->test->teardown complete flow | Any sub-step failed |
| web-compile | `just web-compile` | Web surface code compiled successfully | Compilation failed, stderr contains error details |
| web-fmt | `just web-fmt` | Web surface code formatted (no changes needed or changes applied) | Formatting failed or check-only mode found unformatted code |
| web-lint | `just web-lint` | Web surface code passed all lint checks | Lint violations found, stderr contains rule violations |
| web-unit-test | `just web-unit-test` | All Web surface unit tests passed | At least one unit test failed |

Implementation constraints:
- Each recipe must support both `[linux]` and `[windows]` platform variants
- The `web` aggregate recipe calls sub-recipes in orchestration sequence order, stopping immediately on a non-zero exit code
- `web-teardown` must be validated with `just --dry-run`
- Gate recipes (`web-compile`, `web-fmt`, `web-lint`, `web-unit-test`) are invoked by the quality gate per-task scoping mechanism; they operate ONLY on the web surface code, not other surfaces

## Journey Filter Strategy

| Journey Tag | Match Rule | Description |
|-------------|-----------|-------------|
| `@web` | Exact match | Journey dedicated to web surface |
| `@web-e2e` | Exact match | Web E2E test, assigned to web surface |
| `@smoke` | Exact match | Smoke test, assigned to web surface |
| Other | Ignore | Non-web journeys are not handled by this rule |

## Recipe Generation Requirements

When generating recipes for the web surface, the agent must follow these structural constraints. Shared constraints (naming, dual platform, exit code semantics, test directory path, gate recipes) are defined in SKILL.md's **Standard Target Contract** section — follow those rules. Below are web-specific constraints.

### Form → Naming

- Named surface (key present, e.g., `frontend=web`): `<key>-<verb>` — e.g., `frontend-dev`, `frontend-test`
- Scalar surface (no key, e.g., bare `web`): `<verb>` — e.g., `dev`, `test`

### Aggregate Recipe

The `<key>` aggregate recipe (e.g., `web` for scalar, `frontend` for named) must follow the pattern:

```
just <key>-dev && just <key>-probe && just <key>-test; rc=$?; just <key>-teardown; exit $rc
```

### Server Lifecycle

Recipes for dev, probe, and teardown involve server process management (PID tracking, idempotent startup, health check polling). Follow the patterns defined in `rules/server-lifecycle.md` — do not inline server lifecycle bash code in the generated recipes.
