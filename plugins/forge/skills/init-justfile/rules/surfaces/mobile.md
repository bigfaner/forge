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
| mobile-compile | `just mobile-compile` | Mobile surface code compiled successfully | Compilation failed, stderr contains error details |
| mobile-fmt | `just mobile-fmt` | Mobile surface code formatted (no changes needed or changes applied) | Formatting failed or check-only mode found unformatted code |
| mobile-lint | `just mobile-lint` | Mobile surface code passed all lint checks | Lint violations found, stderr contains rule violations |
| mobile-unit-test | `just mobile-unit-test` | All mobile surface unit tests passed | At least one unit test failed |

Implementation constraints:
- Each recipe must support both `[linux]` and `[windows]` platform variants
- The `mobile` aggregate recipe calls sub-recipes in orchestration sequence order, stopping immediately on a non-zero exit code
- `mobile-teardown` must be validated with `just --dry-run`
- Gate recipes (`mobile-compile`, `mobile-fmt`, `mobile-lint`, `mobile-unit-test`) are invoked by the quality gate per-task scoping mechanism; they operate ONLY on the mobile surface code, not other surfaces

## Journey Filter Strategy

| Journey Tag | Match Rule | Description |
|-------------|-----------|-------------|
| `@mobile` | Exact match | Journey dedicated to mobile surface |
| Other | Ignore | Non-mobile journeys are not handled by this rule |

## Recipe Generation Requirements

When generating recipes for the mobile surface, the agent must follow these structural constraints. Shared constraints (naming, dual platform, exit code semantics, test directory path, gate recipes) are defined in SKILL.md's **Standard Target Contract** section — follow those rules. Below are mobile-specific constraints.

### Form → Naming

- Named surface (key present, e.g., `app=mobile`): `<key>-<verb>` — e.g., `app-dev`, `app-test`
- Scalar surface (no key, e.g., bare `mobile`): `<verb>` — e.g., `dev`, `test`

### Aggregate Recipe

The `<key>` aggregate recipe (e.g., `mobile` for scalar, `app` for named) must follow the pattern:

```
just <key>-test-setup && just <key>-dev && just <key>-probe && just <key>-test; rc=$?; just <key>-teardown; exit $rc
```

Note: mobile's aggregate includes `test-setup` as the first step, preceding `dev`.

### Server Lifecycle

Recipes for dev, probe, and teardown involve server process management (PID tracking, idempotent startup, health check polling). Follow the patterns defined in `rules/server-lifecycle.md` — do not inline server lifecycle bash code in the generated recipes.
