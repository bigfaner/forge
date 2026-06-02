# Surface: mobile

> **Test Type Reference**: The test type for mobile surface is **Mobile E2E Test**, which verifies UI element visibility + user interaction response + screen ID changes via Maestro YAML / device automation.

## Orchestration Sequence

| Step | Exit Code 0 | Exit Code 1 | Exit Code 2 | Next Action |
|------|-------------|-------------|-------------|-------------|
| test-setup | Emulator ready, test environment prepared | Emulator startup failed or environment unavailable | — | Proceed to dev |
| dev | Emulator running, app deployed and ready | Startup failed (emulator unavailable) | — | Proceed to probe |
| probe | Appium health check passed | Appium unresponsive | — | Proceed to test |
| test | Tests passed | Tests failed | Test environment error (retryable) | Proceed to teardown |
| teardown | Cleanup complete (emulator stopped, processes cleaned) | Cleanup failed (residual emulators / processes) | — | End |

Notes:
- test-setup is responsible for emulator preparation and is a prerequisite step for mobile surface; if test-setup fails, exit immediately without continuing to subsequent steps
- When dev fails, **do not continue** with subsequent steps; proceed directly to teardown and exit
- Probe retries up to 3 times with 5-second intervals; if all 3 attempts fail, treat as exit code 1
- Exit code 2 for test step allows re-running; the skill should prompt the user "Test environment error, consider retrying"

## Recipe Invocation Contract

| Recipe Name | just Signature | Exit Code 0 Semantics | Exit Code 1 Semantics |
|-------------|---------------|----------------------|----------------------|
| mobile-test-setup | `just mobile-test-setup` | Emulator ready, test environment prepared | Emulator startup failed, stderr contains error details |
| mobile-dev | `just mobile-dev` | Emulator running, app deployed and ready | Startup failed, stderr contains error details |
| mobile-probe | `just mobile-probe` | Appium health check passed | Appium unresponsive |
| mobile-test | `just mobile-test [journey]` | All mobile E2E tests passed | At least one test failed |
| mobile-teardown | `just mobile-teardown` | Emulator stopped, process cleanup complete | Residual emulators or cleanup error |
| mobile | `just mobile` | Aggregate recipe: test-setup->dev->probe->test->teardown complete flow | Any sub-step failed |

Implementation constraints:
- Each recipe must support both `[linux]` and `[windows]` platform variants
- The `mobile` aggregate recipe calls sub-recipes in orchestration sequence order, stopping immediately on a non-zero exit code
- `mobile-teardown` must be validated with `just --dry-run`

## Journey Filter Strategy

| Journey Tag | Match Rule | Description |
|-------------|-----------|-------------|
| `@mobile` | Exact match | Journey dedicated to mobile surface |
| Other | Ignore | Non-mobile journeys are not handled by this rule |

## Recipe Template (Dual Platform)

```just
# Prepare emulator and test environment for mobile tests
# user-customized
mobile-test-setup:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-test-setup (prepare emulator)" >&2; exit 1

# user-customized
mobile-test-setup:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-test-setup (prepare emulator)" >&2; exit 1

# Start emulator and deploy app
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

# Health check for Appium
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

# Run Mobile E2E tests (optionally filter by journey)
# user-customized
mobile-test journey='':
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-test" >&2; exit 1

# user-customized
mobile-test journey='':
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-test" >&2; exit 1


# Clean up mobile test artifacts
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

# mobile aggregate: test-setup -> dev -> probe -> test -> teardown
mobile:
    #!/usr/bin/env bash
    set -euo pipefail
    just mobile-test-setup && just mobile-dev && just mobile-probe && just mobile-test; rc=$?; just mobile-teardown; exit $rc
```

**LLM Instruction**: Replace the TODO stubs with actual commands derived from language templates and Convention knowledge. The stubs above demonstrate the required recipe structure and dual-platform attribute pattern.
