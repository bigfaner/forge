# Surface: cli

> **Test Type Reference**: The test type for CLI surface is **CLI Functional Test**, which verifies process exit codes + stdout/stderr output via subprocess execution.

## Orchestration Sequence

| Step | Exit Code 0 | Exit Code 1 | Exit Code 2 | Next Action |
|------|-------------|-------------|-------------|-------------|
| test | Tests passed | Tests failed | Test environment error (retryable) | Proceed to teardown |
| teardown | Cleanup complete | Cleanup failed (residual processes) | — | End |

Notes:
- **No dev step**: CLI surface does not start a persistent service
- **No probe step**: CLI tools do not require HTTP health checks
- **No aggregate recipe**: CLI surface does not generate a `cli` aggregate recipe
- Exit code 2 for test step allows re-running; the skill should prompt the user "Test environment error, consider retrying"

## Recipe Invocation Contract

| Recipe Name | just Signature | Exit Code 0 Semantics | Exit Code 1 Semantics |
|-------------|---------------|----------------------|----------------------|
| cli-test | `just cli-test [journey]` | All CLI functional tests passed | At least one test failed |
| cli-teardown | `just cli-teardown` | Cleanup complete | Cleanup failed |

Implementation constraints:
- Each recipe must support both `[linux]` and `[windows]` platform variants
- `cli-teardown` must be validated with `just --dry-run`
- **Do not generate** `cli-dev`, `cli-probe`, or `cli` aggregate recipes

## Journey Filter Strategy

| Journey Tag | Match Rule | Description |
|-------------|-----------|-------------|
| `@cli` | Exact match | Journey dedicated to cli surface |
| Other | Ignore | Non-cli journeys are not handled by this rule |

## Recipe Template (Dual Platform)

<Test-Dir-Path>
The `cli-test` recipe must resolve test scripts from the correct directory:
- **Single surface** (project has 1 surface): `tests/<journey>/`
- **Multi surface** (project has 2+ surfaces): `tests/<surfaceKey>/<journey>/`

When filling the recipe body, use the surface's **key** (not type) for the `<surfaceKey>` segment. Example: for `myapp=cli`, the path is `tests/myapp/<journey>/`.
</Test-Dir-Path>

```just
# Run CLI functional tests (optionally filter by journey)
# user-customized
cli-test journey='':
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement cli-test" >&2; exit 1

# user-customized
cli-test journey='':
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement cli-test" >&2; exit 1


# Clean up CLI test artifacts
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

**LLM Instruction**: Replace the TODO stubs with actual commands derived from language templates and Convention knowledge. The stubs above demonstrate the required recipe structure and dual-platform attribute pattern. **Do not generate** `cli-dev`, `cli-probe`, or `cli` aggregate recipes.
