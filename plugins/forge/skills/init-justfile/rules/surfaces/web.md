# Surface: web

## Orchestration Sequence

| Step    | Exit 0                        | Exit 1                                          | Exit 2 | Next      |
|---------|-------------------------------|-------------------------------------------------|--------|-----------|
| dev     | Server ready, listening       | Startup failed (deps missing / port in use)     | --     | probe     |
| probe   | HTTP health check 2xx         | Health check timeout (service not ready)        | --     | test      |
| test    | All tests pass                | At least one test fails                         | Env error (retry suggested) | teardown |
| teardown| Cleanup complete              | Cleanup failed (residual processes)             | --     | end       |

Notes:
- dev failure: do NOT continue to later steps; go directly to teardown and exit.
- probe: retry up to 3 times, 5-second interval; 3 consecutive failures = exit 1.
- test exit 2: environment error, skill should prompt "Test environment error, suggest retry".

## Recipe Contracts

| Recipe    | Signature          | Exit 0                          | Exit 1                          |
|-----------|--------------------|---------------------------------|---------------------------------|
| web-dev   | `just web-dev`     | Dev server ready, port listening| Startup failed, stderr has detail |
| web-probe | `just web-probe`   | HTTP health check returns 2xx   | Connection refused or timeout   |
| web-test  | `just web-test`    | All test cases pass             | At least one test fails         |
| web-teardown | `just web-teardown` | Processes stopped, port freed | Residual processes or cleanup error |
| web       | `just web`         | Aggregate: dev->probe->test->teardown complete | Any sub-step failed |

Implementation constraints:
- Each recipe MUST support `[linux]` and `[windows]` dual-platform variants.
- `web` aggregate recipe calls sub-recipes in orchestration order, stops on first non-zero exit.
- `web-teardown` MUST pass `just --dry-run` syntax verification.

## Journey Filter Strategy

| Journey Tag | Match Rule  | Description                   |
|-------------|-------------|-------------------------------|
| `@web`      | Exact match | Web surface dedicated journey |
| `@e2e`      | Exact match | End-to-end tests, web surface |
| `@smoke`    | Exact match | Smoke tests, web surface      |
| Other       | Ignore      | Non-web journeys not handled  |

## Recipe Template (dual-platform)

```just
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

# user-customized
web-test:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement web-test" >&2; exit 1

# user-customized
web-test:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement web-test" >&2; exit 1

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

**LLM instruction**: Replace the TODO stubs with actual commands derived from the language template and Convention knowledge. The stubs above show the required recipe structure and dual-platform attribute pattern.
