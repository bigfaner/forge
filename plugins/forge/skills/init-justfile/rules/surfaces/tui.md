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

> **Naming convention**: Recipe names below use the surface type (`tui-`) as prefix for illustration. For **named surfaces**, replace the type prefix with the surface key (e.g., `tui-test` → `terminal-test` for `terminal=tui`). For **scalar surfaces**, the prefix is omitted (e.g., `tui-test` → `test`). See SKILL.md Standard Target Contract for the `<prefix>` definition.

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

When generating recipes for the tui surface, the agent must follow these structural constraints. Shared constraints (naming, dual platform, exit code semantics, test directory path, gate recipes) are defined in SKILL.md's **Standard Target Contract** section — follow those rules. Below are tui-specific constraints.

### Surface-Specific Behavior

TUI surface does **not** generate `dev`, `probe`, or aggregate recipes. The orchestration sequence is `test -> teardown` only. No server lifecycle patterns apply.

### Form → Naming

- Named surface (key present, e.g., `terminal=tui`): `<key>-<verb>` — e.g., `terminal-test`, `terminal-teardown`
- Scalar surface (no key, e.g., bare `tui`): `<verb>` — e.g., `test`, `teardown`
