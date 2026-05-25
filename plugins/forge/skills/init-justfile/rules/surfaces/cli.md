# Surface: cli

## Orchestration Sequence

| Step    | Exit 0              | Exit 1                      | Exit 2 | Next  |
|---------|---------------------|-----------------------------|--------|-------|
| test    | All tests pass      | At least one test fails     | Env error (retry suggested) | teardown |
| teardown| Cleanup complete    | Cleanup failed              | --     | end   |

Notes:
- **No dev step**: CLI surfaces do not start a persistent service.
- **No probe step**: No HTTP health check needed for CLI tools.
- **No aggregate recipe**: CLI surface does not generate a `cli` aggregate recipe.
- test exit 2: environment error, skill should prompt "Test environment error, suggest retry".

## Recipe Contracts

| Recipe       | Signature             | Exit 0                  | Exit 1                      |
|--------------|-----------------------|-------------------------|-----------------------------|
| cli-test     | `just cli-test`       | All test cases pass     | At least one test fails     |
| cli-teardown | `just cli-teardown`   | Cleanup complete        | Cleanup failed              |

Implementation constraints:
- Each recipe MUST support `[linux]` and `[windows]` dual-platform variants.
- `cli-teardown` MUST pass `just --dry-run` syntax verification.
- **Do NOT generate** `cli-dev`, `cli-probe`, or `cli` aggregate recipes.

## Journey Filter Strategy

| Journey Tag | Match Rule  | Description                   |
|-------------|-------------|-------------------------------|
| `@cli`      | Exact match | CLI surface dedicated journey |
| Other       | Ignore      | Non-cli journeys not handled  |

## Recipe Template (dual-platform)

```just
# user-customized
cli-test:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement cli-test" >&2; exit 1

# user-customized
cli-test:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement cli-test" >&2; exit 1

# user-customized
cli-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement cli-teardown" >&2; exit 1

# user-customized
cli-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement cli-teardown" >&2; exit 1
```

**LLM instruction**: Replace the TODO stubs with actual commands derived from the language template and Convention knowledge. The stubs above show the required recipe structure and dual-platform attribute pattern. Do NOT generate `cli-dev`, `cli-probe`, or `cli` aggregate recipes.
