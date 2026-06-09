# Surface: api

> **Test Type Reference**: The test type for API surface is **API Functional Test**, which verifies HTTP status codes + response body JSON + response headers via HTTP client.

## Orchestration Sequence

| Step | Exit Code 0 | Exit Code 1 | Exit Code 2 | Next Action |
|------|-------------|-------------|-------------|-------------|
| dev | API service started successfully, waiting for readiness | Startup failed (missing dependencies / port in use) | — | Proceed to probe |
| probe | Health check passed (GET /healthz returns 2xx) | Health check timed out (service not ready) | — | Proceed to test |
| test | Tests passed | Tests failed | Test environment error (retryable) | Proceed to teardown |
| teardown | Cleanup complete | Cleanup failed (residual processes) | — | End |

Notes:
- When dev fails, **do not continue** with subsequent steps; proceed directly to teardown (which is safe/idempotent — no process to clean if dev never started) and exit
- Probe retries up to 3 times with 5-second intervals; if all 3 attempts fail, treat as exit code 1
- Exit code 2 for test step allows re-running; the skill should prompt the user "Test environment error, consider retrying"

## Recipe Invocation Contract

> **Naming convention**: Recipe names below use the surface type (`api-`) as prefix for illustration. For **named surfaces**, replace the type prefix with the surface key (e.g., `api-dev` → `backend-dev` for `backend=api`). For **scalar surfaces**, the prefix is omitted (e.g., `api-dev` → `dev`). See SKILL.md Standard Target Contract for the `<prefix>` definition.

| Recipe Name | just Signature | Exit Code 0 Semantics | Exit Code 1 Semantics |
|-------------|---------------|----------------------|----------------------|
| api-dev | `just api-dev` | API server ready, listening on port | Startup failed, stderr contains error details |
| api-probe | `just api-probe` | HTTP GET /healthz returns 2xx | Connection refused or timed out |
| api-test | `just api-test [journey]` | All API functional tests passed | At least one test failed |
| api-teardown | `just api-teardown` | Processes terminated, port released | Residual processes or cleanup error |
| api | `just api` | Aggregate recipe: dev->probe->test->teardown complete flow | Any sub-step failed |
| api-compile | `just api-compile` | API surface code compiled successfully | Compilation failed, stderr contains error details |
| api-fmt | `just api-fmt` | API surface code formatted (no changes needed or changes applied) | Formatting failed or check-only mode found unformatted code |
| api-lint | `just api-lint` | API surface code passed all lint checks | Lint violations found, stderr contains rule violations |
| api-unit-test | `just api-unit-test` | All API surface unit tests passed | At least one unit test failed |

Implementation constraints:
- Each recipe must support both `[linux]` and `[windows]` platform variants
- The `api` aggregate recipe calls sub-recipes in orchestration sequence order, stopping immediately on a non-zero exit code
- `api-teardown` must be validated with `just --dry-run`
- Gate recipes (`api-compile`, `api-fmt`, `api-lint`, `api-unit-test`) are invoked by the quality gate per-task scoping mechanism; they operate ONLY on the api surface code, not other surfaces

## Journey Filter Strategy

| Journey Tag | Match Rule | Description |
|-------------|-----------|-------------|
| `@api` | Exact match | Journey dedicated to api surface |
| Other | Ignore | Non-api journeys are not handled by this rule |

## Recipe Generation Requirements

When generating recipes for the api surface, the agent must follow these structural constraints. Shared constraints (naming, dual platform, exit code semantics, test directory path, gate recipes) are defined in SKILL.md's **Standard Target Contract** section — follow those rules. Below are api-specific constraints.

### Form → Naming

- Named surface (key present, e.g., `backend=api`): `<key>-<verb>` — e.g., `backend-dev`, `backend-test`
- Scalar surface (no key, e.g., bare `api`): `<verb>` — e.g., `dev`, `test`

### Aggregate Recipe

The `<key>` aggregate recipe (e.g., `api` for scalar, `backend` for named) must follow the pattern:

```
just <key>-dev && just <key>-probe && just <key>-test; rc=$?; just <key>-teardown; exit $rc
```

### Server Lifecycle

Recipes for dev, probe, and teardown involve server process management (PID tracking, idempotent startup, health check polling). Follow the patterns defined in `rules/server-lifecycle.md` — do not inline server lifecycle bash code in the generated recipes.
