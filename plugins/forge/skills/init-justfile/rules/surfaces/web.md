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
- When dev fails, **do not continue** with subsequent steps; proceed directly to teardown and exit
- Probe retries up to 3 times with 5-second intervals; if all 3 attempts fail, treat as exit code 1
- Exit code 2 for test step allows re-running; the skill should prompt the user "Test environment error, consider retrying"

## Recipe Invocation Contract

| Recipe Name | just Signature | Exit Code 0 Semantics | Exit Code 1 Semantics |
|-------------|---------------|----------------------|----------------------|
| web-dev | `just web-dev` | Development server ready, listening on port | Startup failed, stderr contains error details |
| web-probe | `just web-probe` | HTTP health check returns 2xx | Connection refused or timed out |
| web-test | `just web-test [journey]` | All Web E2E tests passed | At least one test failed |
| web-teardown | `just web-teardown` | Processes terminated, port released | Residual processes or cleanup error |
| web | `just web` | Aggregate recipe: dev->probe->test->teardown complete flow | Any sub-step failed |

Implementation constraints:
- Each recipe must support both `[linux]` and `[windows]` platform variants
- The `web` aggregate recipe calls sub-recipes in orchestration sequence order, stopping immediately on a non-zero exit code
- `web-teardown` must be validated with `just --dry-run`

## Journey Filter Strategy

| Journey Tag | Match Rule | Description |
|-------------|-----------|-------------|
| `@web` | Exact match | Journey dedicated to web surface |
| `@web-e2e` | Exact match | Web E2E test, assigned to web surface |
| `@smoke` | Exact match | Smoke test, assigned to web surface |
| Other | Ignore | Non-web journeys are not handled by this rule |

## Recipe Template (Dual Platform)

```just
# Start web development server
# user-customized
web-dev:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement web-dev (start web dev server)" >&2; exit 1

# user-customized
web-dev:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement web-dev (start web dev server)" >&2; exit 1

# Health check for web server
# user-customized
web-probe:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement web-probe (HTTP health check)" >&2; exit 1

# user-customized
web-probe:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement web-probe (HTTP health check)" >&2; exit 1

# Run Web E2E tests (optionally filter by journey)
# user-customized
web-test journey='':
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement web-test" >&2; exit 1

# user-customized
web-test journey='':
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement web-test" >&2; exit 1


# Clean up web test artifacts
# user-customized
web-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement web-teardown" >&2; exit 1

# user-customized
web-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement web-teardown" >&2; exit 1

# web aggregate: dev -> probe -> test -> teardown
web:
    #!/usr/bin/env bash
    set -euo pipefail
    just web-dev && just web-probe && just web-test; rc=$?; just web-teardown; exit $rc
```

**LLM Instruction**: Replace the TODO stubs with actual commands derived from language templates and Convention knowledge. The stubs above demonstrate the required recipe structure and dual-platform attribute pattern.
