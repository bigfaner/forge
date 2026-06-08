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
| tui-compile | `just tui-compile` | TUI surface code compiled successfully | Compilation failed, stderr contains error details |
| tui-fmt | `just tui-fmt` | TUI surface code formatted (no changes needed or changes applied) | Formatting failed or check-only mode found unformatted code |
| tui-lint | `just tui-lint` | TUI surface code passed all lint checks | Lint violations found, stderr contains rule violations |
| tui-unit-test | `just tui-unit-test` | All TUI surface unit tests passed | At least one unit test failed |

Implementation constraints:
- Each recipe must support both `[linux]` and `[windows]` platform variants
- `tui-teardown` must be validated with `just --dry-run`
- **Do not generate** `tui-dev`, `tui-probe`, or `tui` aggregate recipes
- Gate recipes (`tui-compile`, `tui-fmt`, `tui-lint`, `tui-unit-test`) are invoked by the quality gate per-task scoping mechanism; they operate ONLY on the tui surface code, not other surfaces

## Journey Filter Strategy

| Journey Tag | Match Rule | Description |
|-------------|-----------|-------------|
| `@tui` | Exact match | Journey dedicated to tui surface |
| Other | Ignore | Non-tui journeys are not handled by this rule |

## Recipe Generation Requirements

When generating recipes for the tui surface, the agent must follow these structural constraints.

### Naming

- Named surface (multi-surface project): `<key>-<verb>` — e.g. `terminal-test`, `terminal-teardown`
- Scalar surface (single-surface project): `<verb>` — e.g. `test`, `teardown`

### Dual Platform

Every recipe must have both `[linux]` and `[windows]` attribute variants. The `[linux]` variant must be preceded by a `# user-customized` comment on the line above its definition.

### Exit Code Semantics

Each recipe's exit codes must match the semantics defined in the **Recipe Invocation Contract** table above (exit code 0 = success, exit code 1 = failure).

### Test Directory Path

The `<surfaceKey>-test` recipe must resolve test scripts from:
- **Single surface** (project has 1 surface): `tests/<journey>/`
- **Multi surface** (project has 2+ surfaces): `tests/<surfaceKey>/<journey>/`

Use the surface's **key** (not type) for the `<surfaceKey>` segment. Example: for `terminal=tui`, the path is `tests/terminal/<journey>/`.

### No Server Lifecycle Recipes

TUI surface does **not** generate `dev`, `probe`, or aggregate recipes. The orchestration sequence is `test -> teardown` only.

### Gate Recipes

`compile`, `fmt`, `lint`, `unit-test` recipes are generated only in **multi-surface** projects. Each gate recipe must scope its operation to the tui surface code only.
