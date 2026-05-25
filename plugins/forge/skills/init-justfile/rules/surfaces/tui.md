# Surface: tui

## Orchestration Sequence

| Step    | Exit 0              | Exit 1                      | Exit 2 | Next  |
|---------|---------------------|-----------------------------|--------|-------|
| test    | All tests pass      | At least one test fails     | Env error (retry suggested) | teardown |
| teardown| Cleanup complete    | Cleanup failed              | --     | end   |

Notes:
- **No dev step**: TUI surfaces do not start a persistent service.
- **No probe step**: No HTTP health check needed for TUI applications.
- **No aggregate recipe**: TUI surface does not generate a `tui` aggregate recipe.
- test exit 2: environment error, skill should prompt "Test environment error, suggest retry".

## Recipe Contracts

| Recipe       | Signature             | Exit 0                  | Exit 1                      |
|--------------|-----------------------|-------------------------|-----------------------------|
| tui-test     | `just tui-test`       | All test cases pass     | At least one test fails     |
| tui-teardown | `just tui-teardown`   | Cleanup complete        | Cleanup failed              |

Implementation constraints:
- Each recipe MUST support `[linux]` and `[windows]` dual-platform variants.
- `tui-teardown` MUST pass `just --dry-run` syntax verification.
- **Do NOT generate** `tui-dev`, `tui-probe`, or `tui` aggregate recipes.

## Journey Filter Strategy

| Journey Tag | Match Rule  | Description                   |
|-------------|-------------|-------------------------------|
| `@tui`      | Exact match | TUI surface dedicated journey |
| Other       | Ignore      | Non-tui journeys not handled  |

## Recipe Template (dual-platform)

```just
# user-customized
tui-test:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement tui-test" >&2; exit 1

# user-customized
tui-test:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement tui-test" >&2; exit 1

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

**LLM instruction**: Replace the TODO stubs with actual commands derived from the language template and Convention knowledge. The stubs above show the required recipe structure and dual-platform attribute pattern. Do NOT generate `tui-dev`, `tui-probe`, or `tui` aggregate recipes.
