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

> **Naming convention**: Recipe names below use the surface type (`cli-`) as prefix for illustration. For **named surfaces**, replace the type prefix with the surface key (e.g., `cli-test` → `myapp-test` for `myapp=cli`). For **scalar surfaces**, the prefix is omitted (e.g., `cli-test` → `test`). See SKILL.md Standard Target Contract for the `<prefix>` definition.

| Recipe Name | just Signature | Exit Code 0 Semantics | Exit Code 1 Semantics |
|-------------|---------------|----------------------|----------------------|
| cli-test | `just cli-test [journey]` | All CLI functional tests passed | At least one test failed |
| cli-teardown | `just cli-teardown` | Cleanup complete | Cleanup failed |
| cli-compile | `just cli-compile` | CLI surface code compiled successfully | Compilation failed, stderr contains error details |
| cli-fmt | `just cli-fmt` | CLI surface code formatted (no changes needed or changes applied) | Formatting failed or check-only mode found unformatted code |
| cli-lint | `just cli-lint` | CLI surface code passed all lint checks | Lint violations found, stderr contains rule violations |
| cli-unit-test | `just cli-unit-test` | All CLI surface unit tests passed | At least one unit test failed |

Implementation constraints:
- Each recipe must support both `[linux]` and `[windows]` platform variants
- `cli-teardown` must be validated with `just --dry-run`
- **Do not generate** `cli-dev`, `cli-probe`, or `cli` aggregate recipes
- Gate recipes (`cli-compile`, `cli-fmt`, `cli-lint`, `cli-unit-test`) are invoked by the quality gate per-task scoping mechanism; they operate ONLY on the cli surface code, not other surfaces

## Journey Filter Strategy

| Journey Tag | Match Rule | Description |
|-------------|-----------|-------------|
| `@cli` | Exact match | Journey dedicated to cli surface |
| Other | Ignore | Non-cli journeys are not handled by this rule |

## Recipe Generation Requirements

When generating recipes for the cli surface, the agent must follow these structural constraints. Shared constraints (naming, dual platform, exit code semantics, test directory path, gate recipes) are defined in SKILL.md's **Standard Target Contract** section — follow those rules. Below are cli-specific constraints.

### Surface-Specific Behavior

CLI surface does **not** generate `dev`, `probe`, or aggregate recipes. The orchestration sequence is `test -> teardown` only. No server lifecycle patterns apply.

### Form → Naming

- Named surface (key present, e.g., `myapp=cli`): `<key>-<verb>` — e.g., `myapp-test`, `myapp-teardown`
- Scalar surface (no key, e.g., bare `cli`): `<verb>` — e.g., `test`, `teardown`
