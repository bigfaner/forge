# Surface: api

## Orchestration Sequence

| Step    | Exit 0                        | Exit 1                                          | Exit 2 | Next      |
|---------|-------------------------------|-------------------------------------------------|--------|-----------|
| dev     | API server ready, listening   | Startup failed (deps missing / port in use)     | --     | probe     |
| probe   | HTTP GET /health returns 2xx  | Health check timeout (service not ready)        | --     | test      |
| test    | All tests pass                | At least one test fails                         | Env error (retry suggested) | teardown |
| teardown| Cleanup complete              | Cleanup failed (residual processes)             | --     | end       |

Notes:
- dev failure: do NOT continue to later steps; go directly to teardown and exit.
- probe: retry up to 3 times, 5-second interval; 3 consecutive failures = exit 1.
- test exit 2: environment error, skill should prompt "Test environment error, suggest retry".

## Recipe Contracts

| Recipe       | Signature             | Exit 0                          | Exit 1                          |
|--------------|-----------------------|---------------------------------|---------------------------------|
| api-dev      | `just api-dev`        | API server ready, port listening| Startup failed, stderr has detail |
| api-probe    | `just api-probe`      | HTTP GET /health returns 2xx    | Connection refused or timeout   |
| api-test     | `just api-test`       | All test cases pass             | At least one test fails         |
| api-teardown | `just api-teardown`   | Processes stopped, port freed   | Residual processes or cleanup error |
| api          | `just api`            | Aggregate: dev->probe->test->teardown complete | Any sub-step failed |

Implementation constraints:
- Each recipe MUST support `[linux]` and `[windows]` dual-platform variants.
- `api` aggregate recipe calls sub-recipes in orchestration order, stops on first non-zero exit.
- `api-teardown` MUST pass `just --dry-run` syntax verification.

## Journey Filter Strategy

| Journey Tag | Match Rule  | Description                   |
|-------------|-------------|-------------------------------|
| `@api`      | Exact match | API surface dedicated journey |
| `@smoke`    | Exact match | Smoke tests, API surface      |
| Other       | Ignore      | Non-api journeys not handled  |

## Recipe Template (dual-platform)

```just
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

# user-customized
api-probe:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement api-probe (HTTP GET /health)" >&2; exit 1

# user-customized
api-probe:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement api-probe (HTTP GET /health)" >&2; exit 1

# user-customized
api-test:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement api-test" >&2; exit 1

# user-customized
api-test:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement api-test" >&2; exit 1

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

**LLM instruction**: Replace the TODO stubs with actual commands derived from the language template and Convention knowledge. The stubs above show the required recipe structure and dual-platform attribute pattern.
