# Surface: tui

> **Test Type Reference**: The test type for TUI surface is **Terminal Functional Test**, which verifies terminal rendered output + interactive response sequences via subprocess + stdin pipe.

## Orchestration Sequence

| Step | Exit Code 0 | Exit Code 1 | Exit Code 2 | Next Action |
|------|-------------|-------------|-------------|-------------|
| test | Tests passed | Tests failed | Test environment error (retryable) | Proceed to teardown |
| teardown | Cleanup complete | Cleanup failed (residual processes) | — | End |

Notes:
- **No dev step**: TUI surface does not start a persistent service
- **No probe step**: TUI applications do not require HTTP health checks
- **No aggregate recipe**: TUI surface does not generate a `tui` aggregate recipe
- Exit code 2 for test step allows re-running; the skill should prompt the user "Test environment error, consider retrying"

## Recipe Invocation Contract

| Recipe Name | just Signature | Exit Code 0 Semantics | Exit Code 1 Semantics |
|-------------|---------------|----------------------|----------------------|
| tui-test | `just tui-test [journey]` | All terminal functional tests passed | At least one test failed |
| tui-teardown | `just tui-teardown` | Cleanup complete | Cleanup failed |

Implementation constraints:
- Each recipe must support both `[linux]` and `[windows]` platform variants
- `tui-teardown` must be validated with `just --dry-run`
- **Do not generate** `tui-dev`, `tui-probe`, or `tui` aggregate recipes

## Journey Filter Strategy

| Journey Tag | Match Rule | Description |
|-------------|-----------|-------------|
| `@tui` | Exact match | Journey dedicated to tui surface |
| Other | Ignore | Non-tui journeys are not handled by this rule |

## Recipe Template (Dual Platform)

```just
# Run terminal functional tests (optionally filter by journey)
# user-customized
tui-test journey='':
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement tui-test" >&2; exit 1

# user-customized
tui-test journey='':
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement tui-test" >&2; exit 1


# Clean up TUI test artifacts
# user-customized
tui-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement tui-teardown" >&2; exit 1

# user-customized
tui-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement tui-teardown" >&2; exit 1
```

**LLM Instruction**: Replace the TODO stubs with actual commands derived from language templates and Convention knowledge. The stubs above demonstrate the required recipe structure and dual-platform attribute pattern. **Do not generate** `tui-dev`, `tui-probe`, or `tui` aggregate recipes.
