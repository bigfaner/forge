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
- When dev fails, **do not continue** with subsequent steps; proceed directly to teardown and exit
- Probe retries up to 3 times with 5-second intervals; if all 3 attempts fail, treat as exit code 1
- Exit code 2 for test step allows re-running; the skill should prompt the user "Test environment error, consider retrying"

## Recipe Invocation Contract

| Recipe Name | just Signature | Exit Code 0 Semantics | Exit Code 1 Semantics |
|-------------|---------------|----------------------|----------------------|
| api-dev | `just api-dev` | API server ready, listening on port | Startup failed, stderr contains error details |
| api-probe | `just api-probe` | HTTP GET /healthz returns 2xx | Connection refused or timed out |
| api-test | `just api-test [journey]` | All API functional tests passed | At least one test failed |
| api-teardown | `just api-teardown` | Processes terminated, port released | Residual processes or cleanup error |
| api | `just api` | Aggregate recipe: dev->probe->test->teardown complete flow | Any sub-step failed |

Implementation constraints:
- Each recipe must support both `[linux]` and `[windows]` platform variants
- The `api` aggregate recipe calls sub-recipes in orchestration sequence order, stopping immediately on a non-zero exit code
- `api-teardown` must be validated with `just --dry-run`

## Journey Filter Strategy

| Journey Tag | Match Rule | Description |
|-------------|-----------|-------------|
| `@api` | Exact match | Journey dedicated to api surface |
| Other | Ignore | Non-api journeys are not handled by this rule |

## Recipe Template (Dual Platform)

<Test-Dir-Path>
The `api-test` recipe must resolve test scripts from the correct directory:
- **Single surface** (project has 1 surface): `tests/<journey>/`
- **Multi surface** (project has 2+ surfaces): `tests/<surfaceKey>/<journey>/`

When filling the recipe body, use the surface's **key** (not type) for the `<surfaceKey>` segment. Example: for `backend=api`, the path is `tests/backend/<journey>/`.
</Test-Dir-Path>

```just
# Start API development server
# user-customized
api-dev:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement api-dev (start API server)" >&2; exit 1

# user-customized
api-dev:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement api-dev (start API server)" >&2; exit 1

# Health check for API server
# user-customized
api-probe:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement api-probe (HTTP GET /healthz)" >&2; exit 1

# user-customized
api-probe:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement api-probe (HTTP GET /healthz)" >&2; exit 1

# Run API functional tests (optionally filter by journey)
# user-customized
api-test journey='':
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement api-test" >&2; exit 1

# user-customized
api-test journey='':
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement api-test" >&2; exit 1


# Clean up API test artifacts
# user-customized
api-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement api-teardown" >&2; exit 1

# user-customized
api-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement api-teardown" >&2; exit 1

# api aggregate: dev -> probe -> test -> teardown
api:
    #!/usr/bin/env bash
    set -euo pipefail
    just api-dev && just api-probe && just api-test; rc=$?; just api-teardown; exit $rc
```

**LLM Instruction**: Replace the TODO stubs with actual commands derived from language templates and Convention knowledge. The stubs above demonstrate the required recipe structure and dual-platform attribute pattern.
