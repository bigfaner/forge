# Surface: mobile

## Orchestration Sequence

| Step    | Exit 0                        | Exit 1                                          | Exit 2 | Next      |
|---------|-------------------------------|-------------------------------------------------|--------|-----------|
| dev     | Emulator running, app ready   | Startup failed (emulator unavailable)           | --     | probe     |
| probe   | Appium health check pass      | Appium not responding                           | --     | test      |
| test    | All tests pass                | At least one test fails                         | Env error (retry suggested) | teardown |
| teardown| Cleanup complete              | Cleanup failed (residual emulator/processes)    | --     | end       |

Notes:
- dev failure: do NOT continue to later steps; go directly to teardown and exit.
- probe: retry up to 3 times, 5-second interval; 3 consecutive failures = exit 1.
- test exit 2: environment error, skill should prompt "Test environment error, suggest retry".

## Recipe Contracts

| Recipe          | Signature                | Exit 0                              | Exit 1                              |
|-----------------|--------------------------|-------------------------------------|-------------------------------------|
| mobile-dev      | `just mobile-dev`        | Emulator running, app deployed      | Startup failed, stderr has detail   |
| mobile-probe    | `just mobile-probe`      | Appium health check returns ok      | Appium not responding               |
| mobile-test     | `just mobile-test`       | All test cases pass                 | At least one test fails             |
| mobile-teardown | `just mobile-teardown`   | Emulator stopped, processes cleaned | Residual processes or cleanup error |
| mobile          | `just mobile`            | Aggregate: dev->probe->test->teardown complete | Any sub-step failed |

Implementation constraints:
- Each recipe MUST support `[linux]` and `[windows]` dual-platform variants.
- `mobile` aggregate recipe calls sub-recipes in orchestration order, stops on first non-zero exit.
- `mobile-teardown` MUST pass `just --dry-run` syntax verification.

## Journey Filter Strategy

| Journey Tag | Match Rule  | Description                      |
|-------------|-------------|----------------------------------|
| `@mobile`   | Exact match | Mobile surface dedicated journey |
| Other       | Ignore      | Non-mobile journeys not handled  |

## Recipe Template (dual-platform)

```just
# user-customized
mobile-dev:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-dev (start emulator + deploy app)" >&2; exit 1

# user-customized
mobile-dev:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-dev (start emulator + deploy app)" >&2; exit 1

# user-customized
mobile-probe:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-probe (Appium health check)" >&2; exit 1

# user-customized
mobile-probe:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-probe (Appium health check)" >&2; exit 1

# user-customized
mobile-test:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-test" >&2; exit 1

# user-customized
mobile-test:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-test" >&2; exit 1

# user-customized
mobile-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-teardown" >&2; exit 1

# user-customized
mobile-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-teardown" >&2; exit 1

# mobile aggregate: dev -> probe -> test -> teardown
mobile:
    #!/usr/bin/env bash
    set -euo pipefail
    just mobile-dev && just mobile-probe && just mobile-test; rc=$?; just mobile-teardown; exit $rc
```

**LLM instruction**: Replace the TODO stubs with actual commands derived from the language template and Convention knowledge. The stubs above show the required recipe structure and dual-platform attribute pattern.
